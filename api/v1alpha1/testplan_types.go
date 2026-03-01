package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestPlanSpec defines the desired state of a TestPlan.
type TestPlanSpec struct {
	// Name is the display name of the test plan.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is a description of the test plan.
	// +optional
	Description string `json:"description,omitempty"`

	// Suites is a list of test suite UUIDs included in the plan.
	// +optional
	Suites []string `json:"suites,omitempty"`

	// Tags is a list of tags.
	// +optional
	Tags []string `json:"tags,omitempty"`

	// EndDate is the target completion date (YYYY-MM-DD).
	// +optional
	EndDate string `json:"endDate,omitempty"`
}

// TestPlanStatus defines the observed state of a TestPlan.
type TestPlanStatus struct {
	// ID is the UUID of the test plan in the Argonix API.
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
// +kubebuilder:printcolumn:name="End Date",type=string,JSONPath=`.spec.endDate`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// TestPlan is the Schema for the testplans API.
type TestPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestPlanSpec   `json:"spec,omitempty"`
	Status TestPlanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestPlanList contains a list of TestPlan.
type TestPlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestPlan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestPlan{}, &TestPlanList{})
}
