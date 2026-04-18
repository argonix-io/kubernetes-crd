package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecurityScanSpec defines the desired state of a SecurityScan.
type SecurityScanSpec struct {
	// ScanType is the type of scan to run: cspm, vulnerability, or secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=cspm;vulnerability;secret
	ScanType string `json:"scanType"`

	// Connector is the UUID of the connector to scan. Empty scans all connectors.
	// +optional
	Connector string `json:"connector,omitempty"`

	// TriggeredBy indicates how the scan was triggered: manual or scheduled.
	// +kubebuilder:default=manual
	// +kubebuilder:validation:Enum=manual;scheduled
	// +optional
	TriggeredBy string `json:"triggeredBy,omitempty"`
}

// SecurityScanStatus defines the observed state of a SecurityScan.
type SecurityScanStatus struct {
	// ID is the UUID of the security scan in the Argonix API.
	ID string `json:"id,omitempty"`

	// Status is the current scan status: pending, running, completed, or failed.
	Status string `json:"status,omitempty"`

	// FindingsCount is the number of findings from the scan.
	FindingsCount int64 `json:"findingsCount,omitempty"`

	// StartedAt is the timestamp when the scan started.
	StartedAt string `json:"startedAt,omitempty"`

	// CompletedAt is the timestamp when the scan completed.
	CompletedAt string `json:"completedAt,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.scanType`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Findings",type=integer,JSONPath=`.status.findingsCount`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SecurityScan is the Schema for the securityscans API.
type SecurityScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecurityScanSpec   `json:"spec,omitempty"`
	Status            SecurityScanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityScanList contains a list of SecurityScan.
type SecurityScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityScan{}, &SecurityScanList{})
}
