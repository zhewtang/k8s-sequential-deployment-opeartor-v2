/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ptr "k8s.io/utils/pointer"
	zhewtangbatchv1 "zhewtang.github.io/ksov2/api/v1"
)

// SequentialJobReconciler reconciles a SequentialJob object
type SequentialJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.zhewtang.github.io,resources=sequentialjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.zhewtang.github.io,resources=sequentialjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.zhewtang.github.io,resources=sequentialjobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SequentialJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SequentialJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling SequentialJob")

	var sj zhewtangbatchv1.SequentialJob
	if err := r.Get(ctx, req.NamespacedName, &sj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !sj.DeletionTimestamp.IsZero() {
		logger.Info("SequentialJob deleted")
		return ctrl.Result{}, nil
	}

	if err := r.reconcileChildJobs(ctx, &sj); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile child jobs: %w", err)
	}

	if err := r.reconcileJobState(ctx, &sj); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile status: %w", err)
	}

	logger.Info("Reconciled SequentialJob")
	return ctrl.Result{}, nil
}

func (r *SequentialJobReconciler) reconcileChildJobs(ctx context.Context, sj *zhewtangbatchv1.SequentialJob) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling child jobs")

	for i := range sj.Spec.Jobs {
		var j batchv1.Job
		jobName := getJobName(sj, i)
		err := r.Get(ctx, client.ObjectKey{
			Namespace: sj.Namespace,
			Name:      jobName,
		}, &j)
		if apierrors.IsNotFound(err) {
			logger.Info("Child job does not exist. Creating child job", "job", jobName, "index", i)
			return r.reconcileChildJob(ctx, sj, i)
		} else if err != nil {
			return fmt.Errorf("failed to get child job[%d] %q: %w", i, jobName, err)
		}

		switch jobStateFor(j) {
		case zhewtangbatchv1.Comleted:
			logger.Info("Child job completed", "job", jobName, "index", i)
		case zhewtangbatchv1.Failure:
			logger.Info("Child job failed", "job", jobName, "index", i)
			return nil
		case zhewtangbatchv1.Suspended:
			logger.Info("Child job suspended", "job", jobName, "index", i)
			return nil
		case zhewtangbatchv1.Unknown:
			logger.Info("Child job unknown", "job", jobName, "index", i)
			return nil
		}
	}
	logger.Info("All child jobs completed")
	return nil
}

func getJobName(sj *zhewtangbatchv1.SequentialJob, i int) string {
	return fmt.Sprintf("%s-%d", sj.Name, i)
}

func (r *SequentialJobReconciler) reconcileJobState(ctx context.Context, sj *zhewtangbatchv1.SequentialJob) error {
	logger := log.FromContext(ctx)
	logger.Info("Updating state")

	var childStates []zhewtangbatchv1.ChildJobState
	for i := range sj.Spec.Jobs {
		jobName := getJobName(sj, i)
		var j batchv1.Job
		err := r.Get(ctx, client.ObjectKey{
			Namespace: sj.Namespace,
			Name:      jobName,
		}, &j)
		err = client.IgnoreNotFound(err)
		if err != nil {
			return fmt.Errorf("failed to get child job[%d] %q: %w", i, jobName, err)
		}

		childStates = append(childStates, zhewtangbatchv1.ChildJobState{
			JobName:  jobName,
			JobState: jobStateFor(j),
		})
	}

	newStatus := zhewtangbatchv1.SequentialJobStatus{
		ChildJobStates: childStates,
		OverallState:   overallJobState(childStates),
	}

	if equality.Semantic.DeepEqual(sj.Status, newStatus) {
		return nil
	}
	sj.Status = newStatus
	logger.Info("Updating status", "status", sj.Status)

	if err := r.Status().Update(ctx, sj); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// reconcileChildJob creates a child job if it does not exist
func (r *SequentialJobReconciler) reconcileChildJob(ctx context.Context, sj *zhewtangbatchv1.SequentialJob, i int) error {
	logger := log.FromContext(ctx)
	logger.Info("Create a child job", "index", i)

	jobName := getJobName(sj, i)
	j := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: sj.GetNamespace(),
			Name:      jobName,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: ptr.Int32(0),
			Template:     sj.Spec.Jobs[i],
		},
	}

	if err := ctrl.SetControllerReference(sj, &j, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	if err := r.Create(ctx, &j); err != nil {
		return fmt.Errorf("failed to create child job: %w", err)
	}

	logger.Info("Created child job", "job", jobName, "index", i)

	return nil
}

func jobStateFor(j batchv1.Job) zhewtangbatchv1.JobState {
	for _, c := range j.Status.Conditions {
		if c.Type == "Complete" && c.Status == "True" {
			return zhewtangbatchv1.Comleted
		} else if c.Type == "Failed" && c.Status == "True" {
			return zhewtangbatchv1.Failure
		} else if c.Type == "Suspended" && c.Status == "True" {
			return zhewtangbatchv1.Suspended
		}
	}
	return zhewtangbatchv1.Unknown
}

func overallJobState(cjs []zhewtangbatchv1.ChildJobState) zhewtangbatchv1.JobState {
	for _, js := range cjs {
		if js.JobState == zhewtangbatchv1.Suspended {
			return zhewtangbatchv1.Suspended
		} else if js.JobState == zhewtangbatchv1.Failure {
			return zhewtangbatchv1.Failure
		} else if js.JobState == zhewtangbatchv1.Unknown {
			return zhewtangbatchv1.Unknown
		}
	}
	return zhewtangbatchv1.Comleted
}

// SetupWithManager sets up the controller with the Manager.
func (r *SequentialJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zhewtangbatchv1.SequentialJob{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
