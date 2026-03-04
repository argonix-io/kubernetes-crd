package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MaintenanceWindowSpec defines the desired state of a MaintenanceWindow.
type MaintenanceWindowSpec struct {
	// Name is the display name of the maintenance window.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// GroupID is the UUID of the group to apply the window to.
	// +kubebuilder:validation:Required
	GroupID string `json:"groupId"`

	// StartsAt is the start time (ISO 8601). Used with "once" repeat.
	// +optional
	StartsAt string `json:"startsAt,omitempty"`

	// EndsAt is the end time (ISO 8601). Used with "once" repeat.
	// +optional
	EndsAt string `json:"endsAt,omitempty"`

	// Repeat is the recurrence type.
	// +kubebuilder:validation:Enum=once;daily;weekly;monthly;cron
	// +kubebuilder:default=once
	// +optional
	Repeat string `json:"repeat,omitempty"`

	// TimeFrom is the daily start time (HH:MM). Used with recurring schedules.
	// +optional
	TimeFrom string `json:"timeFrom,omitempty"`

	// TimeTo is the daily end time (HH:MM). Used with recurring schedules.
	// +optional
	TimeTo string `json:"timeTo,omitempty"`

	// Weekdays is a comma-separated weekday numbers (1=Mon, 7=Sun). Used with "weekly".
	// +optional
	Weekdays string `json:"weekdays,omitempty"`

	// DayOfMonth is the day of month (1-31). Used with "monthly".
	// +optional
	DayOfMonth int64 `json:"dayOfMonth,omitempty"`

	// CronExpression is a cron expression. Used with "cron" repeat.
	// +optional
	CronExpression string `json:"cronExpression,omitempty"`

	// IsActive indicates whether the maintenance window is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// MaintenanceWindowStatus defines the observed state of a MaintenanceWindow.
type MaintenanceWindowStatus struct {
	// ID is the UUID of the maintenance window in the Argonix API.
	ID string `json:"id,omitempty"`

	// ScheduleSummary is a human-readable schedule summary.
	ScheduleSummary string `json:"scheduleSummary,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Repeat",type=string,JSONPath=`.spec.repeat`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// MaintenanceWindow is the Schema for the maintenancewindows API.
type MaintenanceWindow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MaintenanceWindowSpec   `json:"spec,omitempty"`
	Status            MaintenanceWindowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MaintenanceWindowList contains a list of MaintenanceWindow.
type MaintenanceWindowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaintenanceWindow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaintenanceWindow{}, &MaintenanceWindowList{})
}
