package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Alchemist struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlchemistSpec   `json:"spec"`
	Status AlchemistStatus `json:"status,omitempty"`
}

// spec defines the desired state
type AlchemistSpec struct {
	Name          string `json:"name"`
	Replicas      *int32 `json:"replicas"`
	ContainerPort int32  `json:"containerPort"`
}

// observed state
type AlchemistStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlchemistList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Alchemist `json:"items"`
}
