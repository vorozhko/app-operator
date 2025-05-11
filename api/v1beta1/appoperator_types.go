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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppoperatorSpec defines the desired state of Appoperator.
type AppoperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image specify the container image for the deployment
	Image string `json:"image,omitempty"`

	// Replicas specify the container replica count for the deployment
	Replicas *int32 `json:"replicas,omitempty"`
}

// AppoperatorStatus defines the observed state of Appoperator.
type AppoperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions represent the latest available observations of the Appoperator's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Appoperator is the Schema for the appoperators API.
type Appoperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppoperatorSpec   `json:"spec,omitempty"`
	Status AppoperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppoperatorList contains a list of Appoperator.
type AppoperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Appoperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Appoperator{}, &AppoperatorList{})
}
