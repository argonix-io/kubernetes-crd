package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorGroupSpec defines the desired state of a MonitorGroup.
type MonitorGroupSpec struct {
	// Name is the display name of the monitoring folder.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description describes the monitoring folder.
	// +optional
	Description *string `json:"description,omitempty"`
}

// MonitorGroupStatus contains the observed Argonix state.
type MonitorGroupStatus struct {
	ID string `json:"id,omitempty"`

	DateCreated  string             `json:"dateCreated,omitempty"`
	DateModified string             `json:"dateModified,omitempty"`
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// MonitorGroup is the schema for the monitorgroups API.
type MonitorGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MonitorGroupSpec   `json:"spec,omitempty"`
	Status            MonitorGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MonitorGroupList contains a list of MonitorGroup objects.
type MonitorGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MonitorGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MonitorGroup{}, &MonitorGroupList{})
}
