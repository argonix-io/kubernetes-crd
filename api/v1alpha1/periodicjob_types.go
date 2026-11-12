package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PeriodicJobSpec defines the desired state of a PeriodicJob.
// +kubebuilder:validation:XValidation:rule="self.kind == oldSelf.kind",message="kind is immutable"
type PeriodicJobSpec struct {
	// Name is the display name of the periodic job.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Kind selects the job type.
	// +kubebuilder:validation:Enum=agent;application_reconcile
	// +kubebuilder:default=agent
	// +optional
	Kind *string `json:"kind,omitempty"`

	// Prompt is the instruction for an agent job.
	// +optional
	Prompt *string `json:"prompt,omitempty"`

	// ScheduleHour is the hour of the day in UTC.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	// +optional
	ScheduleHour *int `json:"scheduleHour,omitempty"`

	// ScheduleMinute is the minute of the hour.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=59
	// +optional
	ScheduleMinute *int `json:"scheduleMinute,omitempty"`

	// ScheduleDays selects weekdays, with Monday at 0. Empty runs every day.
	// +kubebuilder:validation:MaxItems=7
	// +optional
	ScheduleDays *[]int `json:"scheduleDays,omitempty"`

	// ConnectorIDs lists the scoped Argonix connector UUIDs.
	// +optional
	ConnectorIDs *[]string `json:"connectorIds,omitempty"`

	// ApplicationRootIDs lists the scoped inventory application UUIDs.
	// +kubebuilder:validation:MaxItems=3
	// +optional
	ApplicationRootIDs *[]string `json:"applicationRootIds,omitempty"`

	// NotificationChannelIDs lists the Argos chat channel UUIDs.
	// +optional
	NotificationChannelIDs *[]string `json:"notificationChannelIds,omitempty"`

	// AlertChannelIDs lists the Argonix alert channel UUIDs.
	// +optional
	AlertChannelIDs *[]string `json:"alertChannelIds,omitempty"`

	// AutoCreateTicket enables a validation ticket after a successful run.
	// +optional
	AutoCreateTicket *bool `json:"autoCreateTicket,omitempty"`

	// DeliverableRetentionDays sets retention in days. Set zero to clear retention.
	// +kubebuilder:validation:Minimum=0
	// +optional
	DeliverableRetentionDays *int `json:"deliverableRetentionDays,omitempty"`

	// DeliverableRetentionAction selects what happens to expired deliverables.
	// +kubebuilder:validation:Enum=delete;move
	// +optional
	DeliverableRetentionAction *string `json:"deliverableRetentionAction,omitempty"`

	// DeliverableRetentionDestination is the destination for moved deliverables.
	// +kubebuilder:validation:MaxLength=500
	// +optional
	DeliverableRetentionDestination *string `json:"deliverableRetentionDestination,omitempty"`

	// IsActive enables the job schedule.
	// +optional
	IsActive *bool `json:"isActive,omitempty"`

	// Mode sets the permissions of the job. The API defaults to discussion.
	// +kubebuilder:validation:Enum=discussion;safe_write;remediation
	// +optional
	Mode *string `json:"mode,omitempty"`

	// TimeoutSeconds sets the maximum run time in seconds.
	// +kubebuilder:validation:Minimum=1
	// +optional
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// PeriodicJobStatus contains the observed Argonix state.
type PeriodicJobStatus struct {
	ID string `json:"id,omitempty"`
	// LatestRunStatus is the most recent run status.
	LatestRunStatus string `json:"latestRunStatus,omitempty"`
	// LatestRunDate is when the most recent run was created.
	LatestRunDate string             `json:"latestRunDate,omitempty"`
	DateCreated   string             `json:"dateCreated,omitempty"`
	DateModified  string             `json:"dateModified,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Latest run",type=string,JSONPath=`.status.latestRunStatus`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// PeriodicJob is the schema for the periodicjobs API.
type PeriodicJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PeriodicJobSpec   `json:"spec,omitempty"`
	Status            PeriodicJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PeriodicJobList contains a list of PeriodicJob objects.
type PeriodicJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PeriodicJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PeriodicJob{}, &PeriodicJobList{})
}
