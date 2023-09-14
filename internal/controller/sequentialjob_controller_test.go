/*

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
// +kubebuilder:docs-gen:collapse=Apache License

/*
Ideally, we should have one `<kind>_controller_test.go` for each controller scaffolded and called in the `suite_test.go`.
So, let's write our example test for the CronJob controller (`cronjob_controller_test.go.`)
*/

/*
As usual, we start with the necessary imports. We also define some utility variables.
*/
package controller

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	batchv1 "k8s.io/api/batch/v1"
	zhewtangbatchv1 "zhewtang.github.io/ksov2/api/v1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("SequentialJob controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		SequentialJobNameGood  = "test-sequentialjob-good"
		SequentialJobNameBad   = "test-sequentialjob-bad"
		SequentialJobNamespace = "default"
		JobName                = "test-job"

		timeout  = time.Second * 20
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Simulate happy path", func() {
		It("Should eventually lead to overallstate complete", func() {
			By("By creating a new SequentialJob resource")
			ctx := context.TODO()
			sequentialJob := &zhewtangbatchv1.SequentialJob{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "zhewtang.github.io/ksov2/api/v1",
					Kind:       "SequentialJob",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      SequentialJobNameGood,
					Namespace: SequentialJobNamespace,
				},
				Spec: setupSequentialJobSpec(),
			}

			Expect(k8sClient.Create(ctx, sequentialJob)).Should(Succeed())

			sequentialJobLookupKey := types.NamespacedName{Name: SequentialJobNameGood, Namespace: SequentialJobNamespace}
			createdSequentialJob := &zhewtangbatchv1.SequentialJob{}

			// We'll need to retry getting this newly created SequentialJob, given that creation may not immediately happen.
			Eventually(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())

			// Let's make sure our config string value was properly converted/handled.
			commonExpect(createdSequentialJob, true)

			// reconcile should create the job1
			By("By creating one child job")
			createdChildJob1 := &batchv1.Job{}
			createdChildJob2 := &batchv1.Job{}

			childJob1LookupKey := types.NamespacedName{Name: getJobName(sequentialJob, 0), Namespace: SequentialJobNamespace}
			childJob2LookupKey := types.NamespacedName{Name: getJobName(sequentialJob, 1), Namespace: SequentialJobNamespace}
			Eventually(jobExists(childJob1LookupKey, createdChildJob1, ctx), timeout, interval).Should(BeTrue())
			Consistently(jobExists(childJob2LookupKey, createdChildJob2, ctx)).Should(BeFalse())

			By("Overall and job1 status should be unknown")
			Expect(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())
			Expect(createdSequentialJob.Status.OverallState).Should(Equal(zhewtangbatchv1.Unknown))
			Expect(createdSequentialJob.Status.ChildJobStates[0].JobState).Should(Equal(zhewtangbatchv1.Unknown))

			By("Simulating job1 completed successfully")
			createdChildJob1.Status.Conditions = append(createdChildJob1.Status.Conditions, batchv1.JobCondition{
				Type:   batchv1.JobComplete,
				Status: corev1.ConditionTrue,
			})
			Expect(k8sClient.Status().Update(ctx, createdChildJob1)).Should(Succeed())

			By("Job1 state should be complted and job 2 should be created")
			Eventually(func() zhewtangbatchv1.JobState {
				Expect(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())
				return createdSequentialJob.Status.ChildJobStates[0].JobState
			}, timeout, interval).Should(Equal(zhewtangbatchv1.Comleted))
			Eventually(jobExists(childJob2LookupKey, createdChildJob2, ctx), timeout, interval).Should(BeTrue())

			By("Simulating job2 completed successfully")
			createdChildJob2.Status.Conditions = append(createdChildJob2.Status.Conditions, batchv1.JobCondition{
				Type:   batchv1.JobComplete,
				Status: corev1.ConditionTrue,
			})
			Expect(k8sClient.Status().Update(ctx, createdChildJob2)).Should(Succeed())

			By("Job2 state should be complted and overall state should be completed")
			Eventually(func() zhewtangbatchv1.JobState {
				Expect(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())
				return createdSequentialJob.Status.ChildJobStates[1].JobState
			}, timeout, interval).Should(Equal(zhewtangbatchv1.Comleted))
			Expect(createdSequentialJob.Status.OverallState).Should(Equal(zhewtangbatchv1.Comleted))
		})
	})

	Context("Simulate unhappy path", func() {
		It("Should eventually lead overallstate to failure state", func() {
			By("By creating a new SequentialJob resource")
			ctx := context.TODO()
			sequentialJob := &zhewtangbatchv1.SequentialJob{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "zhewtang.github.io/ksov2/api/v1",
					Kind:       "SequentialJob",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      SequentialJobNameBad,
					Namespace: SequentialJobNamespace,
				},
				Spec: setupSequentialJobSpec(),
			}

			Expect(k8sClient.Create(ctx, sequentialJob)).Should(Succeed())

			sequentialJobLookupKey := types.NamespacedName{Name: SequentialJobNameBad, Namespace: SequentialJobNamespace}
			createdSequentialJob := &zhewtangbatchv1.SequentialJob{}

			// We'll need to retry getting this newly created SequentialJob, given that creation may not immediately happen.
			Eventually(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())

			// Let's make sure our config string value was properly converted/handled.
			commonExpect(createdSequentialJob, false)

			// reconcile should create the job1
			By("By creating one child job")
			createdChildJob1 := &batchv1.Job{}
			createdChildJob2 := &batchv1.Job{}

			childJob1LookupKey := types.NamespacedName{Name: getJobName(sequentialJob, 0), Namespace: SequentialJobNamespace}
			childJob2LookupKey := types.NamespacedName{Name: getJobName(sequentialJob, 1), Namespace: SequentialJobNamespace}
			Eventually(jobExists(childJob1LookupKey, createdChildJob1, ctx), timeout, interval).Should(BeTrue())
			Consistently(jobExists(childJob2LookupKey, createdChildJob2, ctx)).Should(BeFalse())

			By("Overall and job1 status should be unknown")
			Expect(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())
			Expect(createdSequentialJob.Status.OverallState).Should(Equal(zhewtangbatchv1.Unknown))
			Expect(createdSequentialJob.Status.ChildJobStates[0].JobState).Should(Equal(zhewtangbatchv1.Unknown))

			By("Simulating job1 to failure")
			createdChildJob1.Status.Conditions = append(createdChildJob1.Status.Conditions, batchv1.JobCondition{
				Type:   batchv1.JobFailed,
				Status: corev1.ConditionTrue,
			})
			Expect(k8sClient.Status().Update(ctx, createdChildJob1)).Should(Succeed())

			By("Job1 state should be Failure, job 2 should not be created and overallstate should be Failure")
			Eventually(func() zhewtangbatchv1.JobState {
				Expect(k8sClient.Get(ctx, sequentialJobLookupKey, createdSequentialJob)).Should(Succeed())
				return createdSequentialJob.Status.ChildJobStates[0].JobState
			}, timeout, interval).Should(Equal(zhewtangbatchv1.Failure))
			Eventually(jobExists(childJob2LookupKey, createdChildJob2, ctx), timeout, interval).Should(BeFalse())
			Expect(createdSequentialJob.Status.OverallState).Should(Equal(zhewtangbatchv1.Failure))
		})
	})
})

func commonExpect(createdSequentialJob *zhewtangbatchv1.SequentialJob, goodJob bool) {
	// Let's make sure our config string value was properly converted/handled.
	if goodJob {
		Expect(createdSequentialJob.Name).Should(Equal("test-sequentialjob-good"))
	} else {
		Expect(createdSequentialJob.Name).Should(Equal("test-sequentialjob-bad"))
	}
	Expect(createdSequentialJob.Spec.Jobs[0].Spec.Containers[0].Image).Should(Equal("test-image--0"))
	Expect(createdSequentialJob.Spec.Jobs[0].Spec.Containers[0].Name).Should(Equal("test-container--0"))
	Expect(createdSequentialJob.Spec.Jobs[0].Spec.RestartPolicy).Should(Equal(corev1.RestartPolicyOnFailure))

	Expect(createdSequentialJob.Spec.Jobs[1].Spec.Containers[0].Image).Should(Equal("test-image--1"))
	Expect(createdSequentialJob.Spec.Jobs[1].Spec.Containers[0].Name).Should(Equal("test-container--1"))
	Expect(createdSequentialJob.Spec.Jobs[1].Spec.RestartPolicy).Should(Equal(corev1.RestartPolicyOnFailure))
}

func jobExists(namespacedName types.NamespacedName, job *batchv1.Job, ctx context.Context) bool {
	return k8sClient.Get(ctx, namespacedName, job) == nil
}

func setupSequentialJobSpec() zhewtangbatchv1.SequentialJobSpec {
	return zhewtangbatchv1.SequentialJobSpec{
		Jobs: setupJobs(2),
	}
}

func setupJobs(count int) []corev1.PodTemplateSpec {
	jobs := []corev1.PodTemplateSpec{}
	for i := 0; i < count; i++ {
		jobs = append(jobs, setupJob(i))
	}
	return jobs
}

func setupJob(index int) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  fmt.Sprintf("%s--%d", "test-container", index),
					Image: fmt.Sprintf("%s--%d", "test-image", index),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
