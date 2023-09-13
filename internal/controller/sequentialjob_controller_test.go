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
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	zhewtangbatchv1 "zhewtang.github.io/ksov2/api/v1"
)

func TestSequentialJobReconciler_reconcileJobState(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx context.Context
		sj  *zhewtangbatchv1.SequentialJob
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SequentialJobReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.reconcileJobState(tt.args.ctx, tt.args.sj); (err != nil) != tt.wantErr {
				t.Errorf("SequentialJobReconciler.reconcileJobState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
