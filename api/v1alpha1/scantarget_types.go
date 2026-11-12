package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScanTargetSpec defines the desired state of a ScanTarget.
type ScanTargetSpec struct {
	// Kind selects a Git repository or OCI container image.
	// +kubebuilder:validation:Enum=repository;image
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// Identifier is the HTTPS repository URL or OCI image reference.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000
	// +kubebuilder:validation:Required
	Identifier string `json:"identifier"`

	// EnabledScanners lists scanners supported for this target kind.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	EnabledScanners []string `json:"enabledScanners"`

	// Name is the display name of the scan target.
	// +kubebuilder:validation:MaxLength=255
	// +optional
	Name *string `json:"name,omitempty"`

	// Connector selects a Git or registry connector UUID. Use an empty string to clear it.
	// +optional
	Connector *string `json:"connector,omitempty"`

	// Project associates the target with a CodeProject UUID. Use an empty string to clear it.
	// +optional
	Project *string `json:"project,omitempty"`

	// Branch is the Git branch to scan.
	// +kubebuilder:validation:MaxLength=255
	// +optional
	Branch *string `json:"branch,omitempty"`

	// DefaultRef is the default tag or digest for an unpinned image.
	// +kubebuilder:validation:MaxLength=255
	// +optional
	DefaultRef *string `json:"defaultRef,omitempty"`

	// Platform is the OCI platform to scan.
	// +kubebuilder:validation:MaxLength=100
	// +optional
	Platform *string `json:"platform,omitempty"`

	// Frequency selects the scan schedule.
	// +kubebuilder:validation:Enum=manual;daily;weekly
	// +optional
	Frequency *string `json:"frequency,omitempty"`

	// ScheduleHour is the scan hour in UTC.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	// +optional
	ScheduleHour *int `json:"scheduleHour,omitempty"`

	// ScheduleMinute is the scan minute in UTC.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=59
	// +optional
	ScheduleMinute *int `json:"scheduleMinute,omitempty"`

	// ScheduleDays lists weekdays with Monday at 0. A weekly scan defaults to Monday.
	// +kubebuilder:validation:MaxItems=7
	// +kubebuilder:validation:items:Minimum=0
	// +kubebuilder:validation:items:Maximum=6
	// +optional
	ScheduleDays *[]int `json:"scheduleDays,omitempty"`

	// IsActive enables the target.
	// +optional
	IsActive *bool `json:"isActive,omitempty"`

	// WebhookEnabled enables scanning in response to signed Git webhooks.
	// +optional
	WebhookEnabled *bool `json:"webhookEnabled,omitempty"`

	// EnforceMergeRequests publishes scan outcomes to the merge request gate.
	// +optional
	EnforceMergeRequests *bool `json:"enforceMergeRequests,omitempty"`

	// GateEnvironment selects the environment used to evaluate the merge request gate.
	// +kubebuilder:validation:MaxLength=100
	// +optional
	GateEnvironment *string `json:"gateEnvironment,omitempty"`
}

// ScanTargetStatus contains the observed Argonix state.
type ScanTargetStatus struct {
	ID string `json:"id,omitempty"`
	// LastDispatchedAt is when the most recent scan was dispatched.
	LastDispatchedAt string             `json:"lastDispatchedAt,omitempty"`
	DateCreated      string             `json:"dateCreated,omitempty"`
	DateModified     string             `json:"dateModified,omitempty"`
	Conditions       []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.kind`
// +kubebuilder:printcolumn:name="Identifier",type=string,JSONPath=`.spec.identifier`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ScanTarget is the schema for the scantargets API.
type ScanTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ScanTargetSpec   `json:"spec,omitempty"`
	Status            ScanTargetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScanTargetList contains a list of ScanTarget objects.
type ScanTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScanTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScanTarget{}, &ScanTargetList{})
}
