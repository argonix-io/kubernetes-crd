package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NotificationRuleSpec defines the desired state of a NotificationRule.
type NotificationRuleSpec struct {
	// Name is the display name of the notification rule.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// IsActive indicates whether the notification rule is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// TriggerCondition is the condition that triggers the notification.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=status_change;goes_down;goes_up;degraded;ssl_expiry;test_failing;test_passing;test_run_complete
	TriggerCondition string `json:"triggerCondition"`

	// ConsecutiveFailures is the number of consecutive failures before triggering.
	// +kubebuilder:default=1
	// +optional
	ConsecutiveFailures int64 `json:"consecutiveFailures,omitempty"`

	// CooldownMinutes is the cooldown period in minutes between notifications.
	// +kubebuilder:default=5
	// +optional
	CooldownMinutes int64 `json:"cooldownMinutes,omitempty"`

	// AllMonitors indicates whether the rule applies to all monitors.
	// +kubebuilder:default=false
	// +optional
	AllMonitors bool `json:"allMonitors,omitempty"`

	// AllSyntheticTests indicates whether the rule applies to all synthetic tests.
	// +kubebuilder:default=false
	// +optional
	AllSyntheticTests bool `json:"allSyntheticTests,omitempty"`

	// MonitorTags is a list of tags to match monitors.
	// +optional
	MonitorTags []string `json:"monitorTags,omitempty"`

	// Monitors is a list of monitor UUIDs this rule applies to.
	// +optional
	Monitors []string `json:"monitors,omitempty"`

	// SyntheticTests is a list of synthetic test UUIDs this rule applies to.
	// +optional
	SyntheticTests []string `json:"syntheticTests,omitempty"`

	// Channels is a list of alert channel UUIDs to notify.
	// +kubebuilder:validation:Required
	Channels []string `json:"channels"`

	// AutoInvestigate enables Argos AI auto-investigation when the rule triggers.
	// When enabled, Argos will automatically investigate the root cause and post analysis to channels.
	// +kubebuilder:default=false
	// +optional
	AutoInvestigate bool `json:"autoInvestigate,omitempty"`
}

// NotificationRuleStatus defines the observed state of a NotificationRule.
type NotificationRuleStatus struct {
	// ID is the UUID of the notification rule in the Argonix API.
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
// +kubebuilder:printcolumn:name="Trigger",type=string,JSONPath=`.spec.triggerCondition`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// NotificationRule is the Schema for the notificationrules API.
type NotificationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NotificationRuleSpec   `json:"spec,omitempty"`
	Status NotificationRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NotificationRuleList contains a list of NotificationRule.
type NotificationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NotificationRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NotificationRule{}, &NotificationRuleList{})
}
