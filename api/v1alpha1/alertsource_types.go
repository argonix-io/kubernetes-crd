package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AlertSourceSpec defines the desired state of an AlertSource.
type AlertSourceSpec struct {
	// Name is the display name of the alert source.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// SourceType is the type of external alerting system.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=alertmanager;datadog;grafana;pagerduty;opsgenie;generic
	SourceType string `json:"sourceType"`

	// IsActive indicates whether the alert source is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// Connector is the UUID of the connector to use for investigation/remediation.
	// +optional
	Connector string `json:"connector,omitempty"`

	// Filters is a JSON-encoded filter configuration for incoming alerts.
	// +optional
	Filters string `json:"filters,omitempty"`

	// AutoInvestigate enables Argos AI auto-investigation for ingested alerts.
	// +kubebuilder:default=false
	// +optional
	AutoInvestigate bool `json:"autoInvestigate,omitempty"`

	// AutoRemediate enables Argos AI auto-remediation for ingested alerts.
	// +kubebuilder:default=false
	// +optional
	AutoRemediate bool `json:"autoRemediate,omitempty"`

	// Channels is a list of alert channel UUIDs to notify when alerts are ingested.
	// +optional
	Channels []string `json:"channels,omitempty"`
}

// AlertSourceStatus defines the observed state of an AlertSource.
type AlertSourceStatus struct {
	// ID is the UUID of the alert source in the Argonix API.
	ID string `json:"id,omitempty"`

	// WebhookSecret is the secret token for the webhook endpoint.
	WebhookSecret string `json:"webhookSecret,omitempty"`

	// WebhookURL is the full webhook URL for sending alerts to this source.
	WebhookURL string `json:"webhookUrl,omitempty"`

	// LastReceivedAt is the timestamp of the last received alert.
	LastReceivedAt string `json:"lastReceivedAt,omitempty"`

	// TotalReceived is the total number of alerts received.
	TotalReceived int64 `json:"totalReceived,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Source",type=string,JSONPath=`.spec.sourceType`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Received",type=integer,JSONPath=`.status.totalReceived`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// AlertSource is the Schema for the alertsources API.
type AlertSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AlertSourceSpec   `json:"spec,omitempty"`
	Status            AlertSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AlertSourceList contains a list of AlertSource.
type AlertSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertSource{}, &AlertSourceList{})
}
