package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WorkflowSpec defines the desired state of a Workflow.
type WorkflowSpec struct {
	// Name is the display name of the workflow.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Slug is a URL-friendly identifier.
	// +optional
	Slug string `json:"slug,omitempty"`

	// Description is a description of the workflow.
	// +optional
	Description string `json:"description,omitempty"`

	// Category is the workflow category.
	// +kubebuilder:validation:Enum=identity;incident;onboarding;devops;security;general
	// +kubebuilder:default=general
	// +optional
	Category string `json:"category,omitempty"`

	// Steps is a JSON-encoded list of workflow steps.
	// +optional
	Steps string `json:"steps,omitempty"`

	// RequiredConnectorTypes is a JSON-encoded list of required connector types.
	// +optional
	RequiredConnectorTypes string `json:"requiredConnectorTypes,omitempty"`

	// RequiresConfirmation indicates whether human approval is needed.
	// +kubebuilder:default=true
	// +optional
	RequiresConfirmation bool `json:"requiresConfirmation,omitempty"`

	// Schedule is a cron expression for scheduled execution.
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// IsActive indicates whether the workflow is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// WorkflowStatus defines the observed state of a Workflow.
type WorkflowStatus struct {
	// ID is the UUID of the workflow in the Argonix API.
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
// +kubebuilder:printcolumn:name="Category",type=string,JSONPath=`.spec.category`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Workflow is the Schema for the workflows API.
type Workflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkflowSpec   `json:"spec,omitempty"`
	Status            WorkflowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkflowList contains a list of Workflow.
type WorkflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workflow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workflow{}, &WorkflowList{})
}
