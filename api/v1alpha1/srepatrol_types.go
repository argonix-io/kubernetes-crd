package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SREPatrolSpec defines the desired state of a SREPatrol.
type SREPatrolSpec struct {
	// Name is the display name of the patrol.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// IsActive enables scheduled sweeps.
	// +optional
	IsActive *bool `json:"isActive,omitempty"`

	// CadenceMinutes is the number of minutes between sweeps.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=60
	// +optional
	CadenceMinutes *int `json:"cadenceMinutes,omitempty"`

	// Autonomy sets the approved response mode for operational signals.
	// +kubebuilder:validation:Enum=observe;propose;act
	// +optional
	Autonomy *string `json:"autonomy,omitempty"`

	// MinConfidence sets the reporting confidence threshold.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	MinConfidence *int `json:"minConfidence,omitempty"`

	// Perspective selects the application perspective.
	// +kubebuilder:validation:Enum=all;business;platform
	// +optional
	Perspective *string `json:"perspective,omitempty"`

	// UseApplications enables application scoping.
	// +optional
	UseApplications *bool `json:"useApplications,omitempty"`

	// IncludeUnattributed includes signals without attribution.
	// +optional
	IncludeUnattributed *bool `json:"includeUnattributed,omitempty"`

	// MinimumAttributionConfidence sets the minimum attribution confidence.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	MinimumAttributionConfidence *int `json:"minimumAttributionConfidence,omitempty"`

	// ConnectorIDs lists scoped connector UUIDs.
	// +optional
	ConnectorIDs *[]string `json:"connectorIds,omitempty"`

	// ApplicationRootIDs lists scoped application UUIDs.
	// +optional
	ApplicationRootIDs *[]string `json:"applicationRootIds,omitempty"`

	// NotificationChannelIDs lists Argos chat channel UUIDs.
	// +optional
	NotificationChannelIDs *[]string `json:"notificationChannelIds,omitempty"`

	// AllowedWorkflowIDs lists explicitly approved workflows.
	// +optional
	AllowedWorkflowIDs *[]string `json:"allowedWorkflowIds,omitempty"`

	// WorkflowBindings lists the trigger rules and input mappings for approved workflows.
	// +optional
	WorkflowBindings *[]map[string]apiextensionsv1.JSON `json:"workflowBindings,omitempty"`

	// DeliveryChannelIDs lists result delivery channel UUIDs.
	// +optional
	DeliveryChannelIDs *[]string `json:"deliveryChannelIds,omitempty"`

	// AlertChannelIDs lists ticketing channel UUIDs.
	// +optional
	AlertChannelIDs *[]string `json:"alertChannelIds,omitempty"`

	// TicketingMinSeverity is the minimum severity for creating a ticket.
	// +optional
	TicketingMinSeverity *string `json:"ticketingMinSeverity,omitempty"`

	// AllowMRAuthoring allows drafting merge requests for review.
	// +optional
	AllowMRAuthoring *bool `json:"allowMrAuthoring,omitempty"`

	// CustomPrompt provides standing instructions for the patrol report.
	// +optional
	CustomPrompt *string `json:"customPrompt,omitempty"`
}

// SREPatrolStatus contains the observed Argonix state.
type SREPatrolStatus struct {
	ID string `json:"id,omitempty"`
	// LastRunStatus is the most recent run status.
	LastRunStatus string `json:"lastRunStatus,omitempty"`
	// LastRunDate is when the most recent run was created.
	LastRunDate  string             `json:"lastRunDate,omitempty"`
	DateCreated  string             `json:"dateCreated,omitempty"`
	DateModified string             `json:"dateModified,omitempty"`
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Last run",type=string,JSONPath=`.status.lastRunStatus`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SREPatrol is the schema for the srepatrols API.
type SREPatrol struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SREPatrolSpec   `json:"spec,omitempty"`
	Status            SREPatrolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SREPatrolList contains a list of SREPatrol objects.
type SREPatrolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SREPatrol `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SREPatrol{}, &SREPatrolList{})
}
