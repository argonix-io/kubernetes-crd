package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecurityPolicySpec defines the desired state of a SecurityPolicy.
type SecurityPolicySpec struct {
	// Name is the display name of the security policy.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is an optional description of the policy purpose.
	// +optional
	Description string `json:"description,omitempty"`

	// Rules is a JSON-encoded policy rules object.
	// Example: {"max_critical": 0, "max_high": 5, "scan_max_age_hours": 24, "required_scan_types": ["trivy"]}
	// +kubebuilder:validation:Required
	Rules string `json:"rules"`

	// Environment is the target environment (e.g. production, staging). Empty applies to all.
	// +optional
	Environment string `json:"environment,omitempty"`

	// IsActive indicates whether the policy is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// SecurityPolicyStatus defines the observed state of a SecurityPolicy.
type SecurityPolicyStatus struct {
	// ID is the UUID of the security policy in the Argonix API.
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
// +kubebuilder:printcolumn:name="Environment",type=string,JSONPath=`.spec.environment`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SecurityPolicy is the Schema for the securitypolicies API.
type SecurityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecurityPolicySpec   `json:"spec,omitempty"`
	Status            SecurityPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityPolicyList contains a list of SecurityPolicy.
type SecurityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityPolicy{}, &SecurityPolicyList{})
}
