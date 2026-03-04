package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChatChannelSpec defines the desired state of a ChatChannel.
type ChatChannelSpec struct {
	// ChannelType is the type of chat channel.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=slack;teams;jira
	ChannelType string `json:"channelType"`

	// ChannelID is the external channel identifier.
	// +kubebuilder:validation:Required
	ChannelID string `json:"channelId"`

	// ConnectorID is the UUID of the connector for this channel.
	// +kubebuilder:validation:Required
	ConnectorID string `json:"connectorId"`

	// ChannelName is the display name for the channel.
	// +optional
	ChannelName string `json:"channelName,omitempty"`

	// PersonaID is the UUID of the persona assigned to this channel.
	// +optional
	PersonaID string `json:"personaId,omitempty"`

	// Config is a JSON-encoded additional configuration.
	// +optional
	Config string `json:"config,omitempty"`

	// IsActive indicates whether the channel is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// ChatChannelStatus defines the observed state of a ChatChannel.
type ChatChannelStatus struct {
	// ID is the UUID of the chat channel in the Argonix API.
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
// +kubebuilder:printcolumn:name="Channel",type=string,JSONPath=`.spec.channelName`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ChatChannel is the Schema for the chatchannels API.
type ChatChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ChatChannelSpec   `json:"spec,omitempty"`
	Status            ChatChannelStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChatChannelList contains a list of ChatChannel.
type ChatChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChatChannel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChatChannel{}, &ChatChannelList{})
}
