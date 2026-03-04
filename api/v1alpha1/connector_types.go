package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConnectorSpec defines the desired state of a Connector.
type ConnectorSpec struct {
	// Name is the display name of the connector.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// ConnectorType is the type of external service.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=slack;teams;pagerduty;opsgenie;jira;servicenow;github;gitlab;datadog;grafana;prometheus;cloudwatch;elastic;splunk;sentry;new_relic;aws;gcp;azure;kubernetes;terraform;ansible;jenkins;confluence;notion;linear;zendesk;okta;custom_webhook
	ConnectorType string `json:"connectorType"`

	// Config is a JSON-encoded configuration for the connector (API keys, tokens, etc.).
	// +kubebuilder:validation:Required
	Config string `json:"config"`

	// IsActive indicates whether the connector is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// Capabilities is a JSON-encoded list of connector capabilities.
	// +optional
	Capabilities string `json:"capabilities,omitempty"`

	// Tags is a JSON-encoded list of tags.
	// +optional
	Tags string `json:"tags,omitempty"`
}

// ConnectorStatus defines the observed state of a Connector.
type ConnectorStatus struct {
	// ID is the UUID of the connector in the Argonix API.
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
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.connectorType`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Connector is the Schema for the connectors API.
type Connector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ConnectorSpec   `json:"spec,omitempty"`
	Status            ConnectorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConnectorList contains a list of Connector.
type ConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Connector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Connector{}, &ConnectorList{})
}
