package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DetectionRuleSpec defines the desired state of a DetectionRule.
type DetectionRuleSpec struct {
	// Name is the display name of the detection rule.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is an optional description of the rule purpose.
	// +optional
	Description string `json:"description,omitempty"`

	// RuleType is the type of detection rule.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=threshold;pattern;sequence
	RuleType string `json:"ruleType"`

	// Severity is the severity level assigned to detections from this rule.
	// +kubebuilder:validation:Enum=critical;high;medium;low
	// +kubebuilder:default=medium
	// +optional
	Severity string `json:"severity,omitempty"`

	// IsActive indicates whether the detection rule is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// Config is a JSON-encoded rule configuration specific to the rule type.
	// +kubebuilder:validation:Required
	Config string `json:"config"`

	// MitreTactic is the MITRE ATT&CK tactic ID (e.g. initial-access, lateral-movement).
	// +optional
	MitreTactic string `json:"mitreTactic,omitempty"`

	// MitreTechnique is the MITRE ATT&CK technique ID (e.g. T1078, T1059).
	// +optional
	MitreTechnique string `json:"mitreTechnique,omitempty"`
}

// DetectionRuleStatus defines the observed state of a DetectionRule.
type DetectionRuleStatus struct {
	// ID is the UUID of the detection rule in the Argonix API.
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
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.ruleType`
// +kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.spec.severity`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// DetectionRule is the Schema for the detectionrules API.
type DetectionRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DetectionRuleSpec   `json:"spec,omitempty"`
	Status            DetectionRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DetectionRuleList contains a list of DetectionRule.
type DetectionRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DetectionRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DetectionRule{}, &DetectionRuleList{})
}
