package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HealthNotebookSpec defines the desired state of a HealthNotebook —
// a scheduled AI-generated infrastructure health report ("carnet de santé").
type HealthNotebookSpec struct {
	// Name is the display name of the health notebook.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// IsActive indicates whether the notebook runs on its schedule.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// Frequency is how often the notebook runs.
	// +kubebuilder:validation:Enum=daily;weekly;monthly
	// +kubebuilder:default=daily
	// +optional
	Frequency string `json:"frequency,omitempty"`

	// ScheduleHour is the hour of day (0-23) the notebook runs.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	// +optional
	ScheduleHour int `json:"scheduleHour,omitempty"`

	// ScheduleMinute is the minute of the hour (0-59) the notebook runs.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=59
	// +optional
	ScheduleMinute int `json:"scheduleMinute,omitempty"`

	// ScheduleDays is a JSON-encoded list of weekday indexes (weekly frequency).
	// +optional
	ScheduleDays string `json:"scheduleDays,omitempty"`

	// ScheduleDayOfMonth is the day of month (1-31) for the monthly frequency.
	// +optional
	ScheduleDayOfMonth int `json:"scheduleDayOfMonth,omitempty"`

	// CustomPrompt overrides the default health-check prompt.
	// +optional
	CustomPrompt string `json:"customPrompt,omitempty"`

	// IncludeThreats includes threat detections in the report.
	// +kubebuilder:default=true
	// +optional
	IncludeThreats bool `json:"includeThreats,omitempty"`

	// IncludeFindings includes security findings in the report.
	// +kubebuilder:default=true
	// +optional
	IncludeFindings bool `json:"includeFindings,omitempty"`

	// IncludeMonitors includes monitor/uptime status in the report.
	// +kubebuilder:default=true
	// +optional
	IncludeMonitors bool `json:"includeMonitors,omitempty"`

	// IncludeIncidents includes incidents in the report.
	// +kubebuilder:default=true
	// +optional
	IncludeIncidents bool `json:"includeIncidents,omitempty"`

	// IncludeCost includes FinOps cost anomalies in the report.
	// +kubebuilder:default=true
	// +optional
	IncludeCost bool `json:"includeCost,omitempty"`

	// ConnectorIDs is a JSON-encoded list of connector UUIDs to scope the report.
	// Empty means all organization connectors.
	// +optional
	ConnectorIDs string `json:"connectorIds,omitempty"`

	// NotificationChannelIDs is a JSON-encoded list of chat channel UUIDs.
	// +optional
	NotificationChannelIDs string `json:"notificationChannelIds,omitempty"`

	// AlertChannelIDs is a JSON-encoded list of alert channel UUIDs.
	// +optional
	AlertChannelIDs string `json:"alertChannelIds,omitempty"`

	// TicketingMinSeverity is the minimum severity that opens a ticket
	// (info, warning, error, critical). Empty disables ticketing.
	// +kubebuilder:validation:Enum="";info;warning;error;critical
	// +optional
	TicketingMinSeverity string `json:"ticketingMinSeverity,omitempty"`

	// TicketingAction is what to do when a ticket-worthy finding recurs.
	// +kubebuilder:validation:Enum=none;investigate;remediate
	// +kubebuilder:default=none
	// +optional
	TicketingAction string `json:"ticketingAction,omitempty"`
}

// HealthNotebookStatus defines the observed state of a HealthNotebook.
type HealthNotebookStatus struct {
	// ID is the UUID of the health notebook in the Argonix API.
	ID string `json:"id,omitempty"`

	// LatestScore is the most recent completed health score (0-100).
	LatestScore int `json:"latestScore,omitempty"`

	// LatestSeverity is the highest severity of the most recent run.
	LatestSeverity string `json:"latestSeverity,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Frequency",type=string,JSONPath=`.spec.frequency`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Score",type=integer,JSONPath=`.status.latestScore`
// +kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.status.latestSeverity`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// HealthNotebook is the Schema for the healthnotebooks API.
type HealthNotebook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HealthNotebookSpec   `json:"spec,omitempty"`
	Status            HealthNotebookStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HealthNotebookList contains a list of HealthNotebook.
type HealthNotebookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HealthNotebook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HealthNotebook{}, &HealthNotebookList{})
}
