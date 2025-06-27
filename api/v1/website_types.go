//go:generate controller-gen object paths="."

// +groupName=apps.dreamer123.com
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebsiteSpec defines the desired state of Website
type WebsiteSpec struct {
	// Host 是网站的域名
	Host string `json:"host"`

	// Image 是要部署的网站容器镜像
	Image string `json:"image"`

	// Replicas 是要运行的副本数
	Replicas int32 `json:"replicas"`
}

// WebsiteStatus defines the observed state of Website
type WebsiteStatus struct {
	// AvailableReplicas 表示当前可用的副本数
	AvailableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Website is the Schema for the websites API
type Website struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebsiteSpec   `json:"spec,omitempty"`
	Status WebsiteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WebsiteList contains a list of Website
type WebsiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Website `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Website{}, &WebsiteList{})
}
