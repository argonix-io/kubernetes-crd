package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSuiteSpec defines the desired state of a TestSuite.
type TestSuiteSpec struct {
	// Name is the display name of the test suite.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is a description of the test suite.
	// +optional
	Description string `json:"description,omitempty"`

	// Tags is a list of tags.
	// +optional
	Tags []string `json:"tags,omitempty"`

	// SyntheticTests is a list of synthetic test UUIDs included in the suite.
	// +optional
	SyntheticTests []string `json:"syntheticTests,omitempty"`

	// ManualTestCases is a list of manual test case UUIDs included in the suite.
	// +optional
	ManualTestCases []string `json:"manualTestCases,omitempty"`
}

// TestSuiteStatus defines the observed state of a TestSuite.
type TestSuiteStatus struct {
	// ID is the UUID of the test suite in the Argonix API.
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

// TestSuite is the Schema for the testsuites API.
type TestSuite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestSuiteSpec   `json:"spec,omitempty"`
	Status TestSuiteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestSuiteList contains a list of TestSuite.
type TestSuiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestSuite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestSuite{}, &TestSuiteList{})
}
