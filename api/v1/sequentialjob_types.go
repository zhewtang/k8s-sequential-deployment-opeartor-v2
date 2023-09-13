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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// define the job state
type JobState string

const (
	Unknown   JobState = "Unknown"
	Suspended JobState = "Suspended"
	Comleted  JobState = "Completed"
	Failure   JobState = "Failure"
)

// ChildJobState defines the observed state of ChildJob, which is an element of SequentialJob
type ChildJobState struct {
	JobName  string   `json:"jobName"`
	JobState JobState `json:"jobState"`
}

// SequentialJobSpec defines the desired state of SequentialJob
type SequentialJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Jobs []corev1.PodTemplateSpec `json:"jobs"`
}

// SequentialJobStatus defines the observed state of SequentialJob
type SequentialJobStatus struct {
	OverallState   JobState        `json:"overallState,omitempty"`
	ChildJobStates []ChildJobState `json:"childJobStates,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sj

// SequentialJob is the Schema for the sequentialjobs API
type SequentialJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SequentialJobSpec   `json:"spec,omitempty"`
	Status SequentialJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SequentialJobList contains a list of SequentialJob
type SequentialJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SequentialJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SequentialJob{}, &SequentialJobList{})
}
