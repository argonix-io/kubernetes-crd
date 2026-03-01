package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GroupSpec defines the desired state of a Group.
type GroupSpec struct {
	// Name is the display name of the group.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is a description of the group.
	// +optional
	Description string `json:"description,omitempty"`

	// Tags is a map of key-value tags.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// GroupStatus defines the observed state of a Group.
type GroupStatus struct {
	// ID is the UUID of the group in the Argonix API.
	ID string `json:"id,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Group is the Schema for the groups API.
type Group struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GroupSpec   `json:"spec,omitempty"`
	Status GroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GroupList contains a list of Group.
type GroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Group `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Group{}, &GroupList{})
}
