/*
Copyright 2025 Iaroslav Vorozhko.

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	toolsv1beta1 "github.com/vorozhko/app-operator/api/v1beta1"
)

// AppoperatorReconciler reconciles a Appoperator object
type AppoperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tools.vorozhko.net,resources=appoperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tools.vorozhko.net,resources=appoperators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tools.vorozhko.net,resources=appoperators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Appoperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *AppoperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the Appoperator instance
	var appoperator toolsv1beta1.Appoperator
	if err := r.Get(ctx, req.NamespacedName, &appoperator); err != nil {
		if apierrors.IsNotFound(err) {
			// The resource was deleted, nothing to do
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to fetch Appoperator")
		return ctrl.Result{}, err
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appoperator.Name,
			Namespace: appoperator.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: appoperator.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": appoperator.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": appoperator.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app-container",
							Image: appoperator.Spec.Image,
						},
					},
				},
			},
		},
	}

	// Check if the Deployment already exists
	var existingDeployment appsv1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: appoperator.Name, Namespace: appoperator.Namespace}, &existingDeployment)
	if err != nil && apierrors.IsNotFound(err) {
		// Deployment does not exist, create it
		if err := r.Create(ctx, deployment); err != nil {
			log.Error(err, "Failed to create Deployment")
			return ctrl.Result{}, err
		}
		log.Info("Created Deployment", "deployment", deployment.Name)
	} else if err != nil {
		// Error fetching the Deployment
		log.Error(err, "Failed to fetch Deployment")
		return ctrl.Result{}, err
	} else {
		// Deployment exists, update it if necessary
		// todo: improve Image check for Deployments with multiple images
		if *existingDeployment.Spec.Replicas != *deployment.Spec.Replicas ||
			existingDeployment.Spec.Template.Spec.Containers[0].Image != deployment.Spec.Template.Spec.Containers[0].Image {
			existingDeployment.Spec = deployment.Spec
			if err := r.Update(ctx, &existingDeployment); err != nil {
				log.Error(err, "Failed to update Deployment")
				return ctrl.Result{}, err
			}
			log.Info("Updated Deployment", "deployment", deployment.Name)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppoperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1beta1.Appoperator{}).
		Named("appoperator").
		Complete(r)
}
