package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KnowledgeBaseSpec defines the desired state of a KnowledgeBase.
type KnowledgeBaseSpec struct {
	// Name is the display name of the knowledge base.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// SourceType is the source type for the knowledge base.
	// +kubebuilder:validation:Enum=manual;confluence;notion;github;gitlab;web
	// +kubebuilder:default=manual
	// +optional
	SourceType string `json:"sourceType,omitempty"`

	// ConnectorID is the UUID of the connector used for syncing.
	// +optional
	ConnectorID string `json:"connectorId,omitempty"`

	// IsActive indicates whether the knowledge base is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// SyncConfig is a JSON-encoded sync configuration.
	// +optional
	SyncConfig string `json:"syncConfig,omitempty"`
}

// KnowledgeBaseStatus defines the observed state of a KnowledgeBase.
type KnowledgeBaseStatus struct {
	// ID is the UUID of the knowledge base in the Argonix API.
	ID string `json:"id,omitempty"`

	// LastSyncedAt is the timestamp of the last sync.
	LastSyncedAt string `json:"lastSyncedAt,omitempty"`

	// DocumentCount is the number of documents in the knowledge base.
	DocumentCount int64 `json:"documentCount,omitempty"`

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
// +kubebuilder:printcolumn:name="Docs",type=integer,JSONPath=`.status.documentCount`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// KnowledgeBase is the Schema for the knowledgebases API.
type KnowledgeBase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KnowledgeBaseSpec   `json:"spec,omitempty"`
	Status            KnowledgeBaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KnowledgeBaseList contains a list of KnowledgeBase.
type KnowledgeBaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KnowledgeBase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KnowledgeBase{}, &KnowledgeBaseList{})
}
