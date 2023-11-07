/*
Copyright 2023 Kenji Tanaka.

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

package controllers

import (
	"context"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/engrkenjitanaka/time-scaler-operator/api/v1alpha1"
)

// TimeScalerReconciler reconciles a TimeScaler object
type TimeScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.engineerkenji.com,resources=timescalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.engineerkenji.com,resources=timescalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.engineerkenji.com,resources=timescalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TimeScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *TimeScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	timeScaler := &apiv1alpha1.TimeScaler{}
	err := r.Get(ctx, req.NamespacedName, timeScaler)
	if err != nil {
		return ctrl.Result{}, err
	}

	startTime := timeScaler.Spec.StartTime
	endTime := timeScaler.Spec.EndTime
	replicas := timeScaler.Spec.Replicas

	currentHour := time.Now().UTC().Hour()

	if currentHour >= startTime && currentHour <= endTime {
		for _, deploy := range timeScaler.Spec.Deployments {
			dep := &v1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{
				Name:      deploy.Name,
				Namespace: deploy.Namespace,
			}, dep)
			if err != nil {
				return ctrl.Result{}, err
			}

			if dep.Spec.Replicas != &replicas {
				dep.Spec.Replicas = &replicas
				err := r.Update(ctx, dep)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TimeScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.TimeScaler{}).
		Complete(r)
}
