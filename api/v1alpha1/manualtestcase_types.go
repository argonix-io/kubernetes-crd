package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManualTestStep represents a single step in a manual test case.
type ManualTestStep struct {
	// Description is the step description.
	Description string `json:"description"`

	// Expected is the expected result.
	Expected string `json:"expected"`
}

// ManualTestCaseSpec defines the desired state of a ManualTestCase.
type ManualTestCaseSpec struct {
	// Title is the title of the manual test case.
	// +kubebuilder:validation:Required
	Title string `json:"title"`

	// Description is a description of the test case.
	// +optional
	Description string `json:"description,omitempty"`

	// Preconditions describes the preconditions for the test.
	// +optional
	Preconditions string `json:"preconditions,omitempty"`

	// Steps is the list of test steps.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Steps []ManualTestStep `json:"steps"`

	// Priority is the test priority.
	// +kubebuilder:default="medium"
	// +kubebuilder:validation:Enum=critical;high;medium;low
	// +optional
	Priority string `json:"priority,omitempty"`

	// Tags is a list of tags.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// ManualTestCaseStatus defines the observed state of a ManualTestCase.
type ManualTestCaseStatus struct {
	// ID is the UUID of the manual test case in the Argonix API.
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
// +kubebuilder:printcolumn:name="Priority",type=string,JSONPath=`.spec.priority`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ManualTestCase is the Schema for the manualtestcases API.
type ManualTestCase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManualTestCaseSpec   `json:"spec,omitempty"`
	Status ManualTestCaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ManualTestCaseList contains a list of ManualTestCase.
type ManualTestCaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManualTestCase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManualTestCase{}, &ManualTestCaseList{})
}
