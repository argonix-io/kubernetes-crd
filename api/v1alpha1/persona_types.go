package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersonaSpec defines the desired state of a Persona.
type PersonaSpec struct {
	// Name is the display name of the persona.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is a description of the persona.
	// +optional
	Description string `json:"description,omitempty"`

	// Template is the persona template preset.
	// +kubebuilder:validation:Enum=devops;it_support;hr;security;custom
	// +kubebuilder:default=custom
	// +optional
	Template string `json:"template,omitempty"`

	// SystemPrompt is the custom system prompt for the AI agent.
	// +optional
	SystemPrompt string `json:"systemPrompt,omitempty"`

	// IsActive indicates whether the persona is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// ConnectorIDs is a JSON-encoded list of connector UUIDs.
	// +optional
	ConnectorIDs string `json:"connectorIds,omitempty"`

	// AllowedTools is a JSON-encoded list of allowed tool names.
	// +optional
	AllowedTools string `json:"allowedTools,omitempty"`

	// ApprovalRules is a JSON-encoded approval rules configuration.
	// +optional
	ApprovalRules string `json:"approvalRules,omitempty"`
}

// PersonaStatus defines the observed state of a Persona.
type PersonaStatus struct {
	// ID is the UUID of the persona in the Argonix API.
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
// +kubebuilder:printcolumn:name="Template",type=string,JSONPath=`.spec.template`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Persona is the Schema for the personas API.
type Persona struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PersonaSpec   `json:"spec,omitempty"`
	Status            PersonaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PersonaList contains a list of Persona.
type PersonaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Persona `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Persona{}, &PersonaList{})
}
