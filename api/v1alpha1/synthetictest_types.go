package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyntheticTestSpec defines the desired state of a SyntheticTest.
type SyntheticTestSpec struct {
	// Name is the display name of the synthetic test.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is a description of the synthetic test.
	// +optional
	Description string `json:"description,omitempty"`

	// IsActive indicates whether the synthetic test is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// TestType is the type of synthetic test (api or browser).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=api;browser
	TestType string `json:"testType"`

	// Steps is a JSON-encoded array of step objects.
	// +kubebuilder:validation:Required
	Steps string `json:"steps"`

	// CheckInterval is the number of seconds between checks.
	// +kubebuilder:default=300
	// +optional
	CheckInterval int64 `json:"checkInterval,omitempty"`

	// Timeout is the request timeout in seconds.
	// +kubebuilder:default=30
	// +optional
	Timeout int64 `json:"timeout,omitempty"`

	// Tags is a list of tags.
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Locations is a list of region codes.
	// +optional
	Locations []string `json:"locations,omitempty"`
}

// SyntheticTestStatus defines the observed state of a SyntheticTest.
type SyntheticTestStatus struct {
	// ID is the UUID of the synthetic test in the Argonix API.
	ID string `json:"id,omitempty"`

	// CurrentStatus is the current status.
	CurrentStatus string `json:"currentStatus,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.testType`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.currentStatus`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SyntheticTest is the Schema for the synthetictests API.
type SyntheticTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyntheticTestSpec   `json:"spec,omitempty"`
	Status SyntheticTestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SyntheticTestList contains a list of SyntheticTest.
type SyntheticTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyntheticTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyntheticTest{}, &SyntheticTestList{})
}
