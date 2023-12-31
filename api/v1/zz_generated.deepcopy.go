//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChildJobState) DeepCopyInto(out *ChildJobState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChildJobState.
func (in *ChildJobState) DeepCopy() *ChildJobState {
	if in == nil {
		return nil
	}
	out := new(ChildJobState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SequentialJob) DeepCopyInto(out *SequentialJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SequentialJob.
func (in *SequentialJob) DeepCopy() *SequentialJob {
	if in == nil {
		return nil
	}
	out := new(SequentialJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SequentialJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SequentialJobList) DeepCopyInto(out *SequentialJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SequentialJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SequentialJobList.
func (in *SequentialJobList) DeepCopy() *SequentialJobList {
	if in == nil {
		return nil
	}
	out := new(SequentialJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SequentialJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SequentialJobSpec) DeepCopyInto(out *SequentialJobSpec) {
	*out = *in
	if in.Jobs != nil {
		in, out := &in.Jobs, &out.Jobs
		*out = make([]corev1.PodTemplateSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SequentialJobSpec.
func (in *SequentialJobSpec) DeepCopy() *SequentialJobSpec {
	if in == nil {
		return nil
	}
	out := new(SequentialJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SequentialJobStatus) DeepCopyInto(out *SequentialJobStatus) {
	*out = *in
	if in.ChildJobStates != nil {
		in, out := &in.ChildJobStates, &out.ChildJobStates
		*out = make([]ChildJobState, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SequentialJobStatus.
func (in *SequentialJobStatus) DeepCopy() *SequentialJobStatus {
	if in == nil {
		return nil
	}
	out := new(SequentialJobStatus)
	in.DeepCopyInto(out)
	return out
}
