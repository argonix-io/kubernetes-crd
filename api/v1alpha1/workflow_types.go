package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

	// InputHints is a JSON-encoded map annotating each {{input.x}} referenced
	// in Steps with an optional description and example. It does not define the
	// input set (that is always derived from the {{input.x}} references in
	// Steps); orphan hints are ignored.
	// +optional
	InputHints string `json:"inputHints,omitempty"`

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

	// ExtractInputsWithAI allows input extraction from free text.
	// +optional
	ExtractInputsWithAI *bool `json:"extractInputsWithAI,omitempty"`

	// ScheduleTimezone is the IANA time zone for cron execution.
	// +optional
	ScheduleTimezone *string `json:"scheduleTimezone,omitempty"`

	// ScheduleInputs provides values for scheduled workflow inputs.
	// +optional
	ScheduleInputs *map[string]apiextensionsv1.JSON `json:"scheduleInputs,omitempty"`

	// SecurityTrigger selects when a security event should run the workflow.
	// +optional
	SecurityTrigger *string `json:"securityTrigger,omitempty"`

	// SecurityMinSeverity is the minimum severity for a security trigger.
	// +optional
	SecurityMinSeverity *string `json:"securityMinSeverity,omitempty"`

	// SecurityDetectionRule scopes the trigger to a detection rule UUID. Use an empty string to clear it.
	// +optional
	SecurityDetectionRule *string `json:"securityDetectionRule,omitempty"`

	// IncidentTrigger selects the incident lifecycle trigger.
	// +optional
	IncidentTrigger *string `json:"incidentTrigger,omitempty"`

	// IncidentAutonomy sets the approved incident response mode.
	// +optional
	IncidentAutonomy *string `json:"incidentAutonomy,omitempty"`

	// IncidentMinSeverity is the minimum incident severity.
	// +optional
	IncidentMinSeverity *string `json:"incidentMinSeverity,omitempty"`

	// IncidentMinConfidence is the minimum incident confidence percentage.
	// +optional
	IncidentMinConfidence *int `json:"incidentMinConfidence,omitempty"`

	// IncidentMaxAffectedResources limits the incident blast radius.
	// +optional
	IncidentMaxAffectedResources *int `json:"incidentMaxAffectedResources,omitempty"`

	// IncidentAllowedEnvironments lists permitted incident environments.
	// +optional
	IncidentAllowedEnvironments *[]string `json:"incidentAllowedEnvironments,omitempty"`

	// IncidentAllowedTools lists explicitly approved incident tools.
	// +optional
	IncidentAllowedTools *[]string `json:"incidentAllowedTools,omitempty"`

	// IncidentAllowedConnectorTypes lists explicitly approved incident connector types.
	// +optional
	IncidentAllowedConnectorTypes *[]string `json:"incidentAllowedConnectorTypes,omitempty"`

	// IncidentRequireCompensation requires a compensation step for mutations.
	// +optional
	IncidentRequireCompensation *bool `json:"incidentRequireCompensation,omitempty"`

	// IncidentRequireVerification requires verification after incident actions.
	// +optional
	IncidentRequireVerification *bool `json:"incidentRequireVerification,omitempty"`

	// IncidentStopOnAmbiguousInput stops the workflow on ambiguous incident input.
	// +optional
	IncidentStopOnAmbiguousInput *bool `json:"incidentStopOnAmbiguousInput,omitempty"`

	// IncidentTimeoutMinutes limits the duration of an incident run.
	// +optional
	IncidentTimeoutMinutes *int `json:"incidentTimeoutMinutes,omitempty"`
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

	// ScheduleLastRunAt is the last scheduled execution time.
	// +optional
	ScheduleLastRunAt string `json:"scheduleLastRunAt,omitempty"`

	// ScheduleNextRunAt is the next scheduled execution time.
	// +optional
	ScheduleNextRunAt string `json:"scheduleNextRunAt,omitempty"`
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
