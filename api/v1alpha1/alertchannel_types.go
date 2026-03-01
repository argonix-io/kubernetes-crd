package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AlertChannelSpec defines the desired state of an AlertChannel.
type AlertChannelSpec struct {
	// Name is the display name of the alert channel.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// ChannelType is the type of alert channel.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=email;slack;webhook;pagerduty;opsgenie;telegram;discord;microsoft_teams
	ChannelType string `json:"channelType"`

	// Config is a JSON-encoded configuration object specific to the channel type.
	// +kubebuilder:validation:Required
	Config string `json:"config"`

	// IsActive indicates whether the alert channel is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// AlertChannelStatus defines the observed state of an AlertChannel.
type AlertChannelStatus struct {
	// ID is the UUID of the alert channel in the Argonix API.
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
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.channelType`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// AlertChannel is the Schema for the alertchannels API.
type AlertChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AlertChannelSpec   `json:"spec,omitempty"`
	Status            AlertChannelStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AlertChannelList contains a list of AlertChannel.
type AlertChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertChannel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertChannel{}, &AlertChannelList{})
}
