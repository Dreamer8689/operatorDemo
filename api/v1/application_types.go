/*
Copyright 2025 dreamer.

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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationSpec defines the desired state of Application.
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Application. Edit application_types.go to remove/update
	Replicas   int32              `json:"replicas,omitempty"`
	Deployment DeploymentTemplate `json:"deployment,omitempty"`
	Service    ServiceTemplate    `json:"service,omitempty"`
}

type DeploymentTemplate struct {
	appsv1.DeploymentSpec `json:",inline"`
}

type ServiceTemplate struct {
	v1.ServiceSpec `json:",inline"`
}

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	Workflow appsv1.DeploymentStatus `json:"workflow,omitempty"`
	Network  v1.ServiceStatus        `json:"network,omitempty"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

//+kubebuilder:resource:path=applications,singular=application,scope=Namespaced,shortName=app

// Application is the Schema for the applications API.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

func (r *Application) DeepCopyObject() runtime.Object {
	if c := r.DeepCopy(); c != nil {
		return c
	}
	return nil
}

var applicationLog = logf.Log.WithName("application-resource")

func (r *Application) Default() {
	applicationLog.Info("Defaulting for Application", "name", r.GetName())
	if r.Spec.Deployment.Replicas == nil {
		r.Spec.Deployment.Replicas = new(int32)
		*r.Spec.Deployment.Replicas = 3
	}
}

func (r *Application) validateApplication() error {

	if r.Spec.Deployment.Replicas == nil {
		return nil
	}

	if *r.Spec.Deployment.Replicas < 1 {
		return nil
	}

	if *r.Spec.Deployment.Replicas > 10 {
		return fmt.Errorf("replicas must be less than 10")
	}

	return nil

}

func (r *Application) ValidateCreate() error {
	applicationLog.Info("Validating for Application upon creation", "name", r.GetName())
	return r.validateApplication()
}

func (r *Application) ValidateUpdate(old runtime.Object) error {
	applicationLog.Info("Validating for Application upon update", "name", r.GetName())
	return r.validateApplication()
}

func (r *Application) ValidateDelete() error {
	applicationLog.Info("Validating for Application upon deletion", "name", r.GetName())
	return nil
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
