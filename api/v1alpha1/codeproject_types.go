package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CodeProjectSpec defines the desired state of a CodeProject.
type CodeProjectSpec struct {
	// RepositoryURL is the Git repository URL for defensive code analysis.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000
	// +kubebuilder:validation:Required
	RepositoryURL string `json:"repositoryUrl"`

	// Name is the display name of the project. The API can derive it from the repository URL.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +optional
	Name *string `json:"name,omitempty"`

	// Key is the project key. The API can derive it from the project name.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +optional
	Key *string `json:"key,omitempty"`

	// Connector selects the Git connector UUID. Use an empty string to clear it.
	// +optional
	Connector *string `json:"connector,omitempty"`

	// ProviderProjectID is the source control provider project identifier.
	// +kubebuilder:validation:MaxLength=255
	// +optional
	ProviderProjectID *string `json:"providerProjectId,omitempty"`

	// DefaultBranch is the project default branch.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +optional
	DefaultBranch *string `json:"defaultBranch,omitempty"`

	// Visibility sets the visibility of the project.
	// +kubebuilder:validation:Enum=private;internal;public
	// +optional
	Visibility *string `json:"visibility,omitempty"`

	// QualityGate selects a quality gate UUID. Use an empty string to clear it.
	// +optional
	QualityGate *string `json:"qualityGate,omitempty"`

	// LongLivedQualityGate selects a quality gate UUID for long-lived branches. Use an empty string to clear it.
	// +optional
	LongLivedQualityGate *string `json:"longLivedQualityGate,omitempty"`

	// EnabledEngines lists the defensive code analysis engines to run.
	// +kubebuilder:validation:MinItems=1
	// +optional
	EnabledEngines *[]string `json:"enabledEngines,omitempty"`

	// PublishQualityGate enables Git quality gate status publication.
	// +optional
	PublishQualityGate *bool `json:"publishQualityGate,omitempty"`

	// PublishReviewFeedback enables Git review feedback publication.
	// +optional
	PublishReviewFeedback *bool `json:"publishReviewFeedback,omitempty"`

	// ExcludedPaths lists repository paths or glob patterns to exclude.
	// +kubebuilder:validation:MaxItems=200
	// +optional
	ExcludedPaths *[]string `json:"excludedPaths,omitempty"`

	// WebhookAnalysisMode selects which code analysis to run on repository events.
	// +kubebuilder:validation:Enum=none;standard;ai;both
	// +optional
	WebhookAnalysisMode *string `json:"webhookAnalysisMode,omitempty"`

	// WebhookStandardScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookStandardScope *string `json:"webhookStandardScope,omitempty"`

	// WebhookStandardNonDefaultScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookStandardNonDefaultScope *string `json:"webhookStandardNonDefaultScope,omitempty"`

	// WebhookStandardLongLivedScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookStandardLongLivedScope *string `json:"webhookStandardLongLivedScope,omitempty"`

	// WebhookAIScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookAIScope *string `json:"webhookAIScope,omitempty"`

	// WebhookAINonDefaultScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookAINonDefaultScope *string `json:"webhookAINonDefaultScope,omitempty"`

	// WebhookAILongLivedScope controls whether repository webhook analysis uses full or changed code.
	// +kubebuilder:validation:Enum=full;changed
	// +optional
	WebhookAILongLivedScope *string `json:"webhookAILongLivedScope,omitempty"`

	// WebhookLongLivedPatterns lists patterns matching long-lived branches.
	// +optional
	WebhookLongLivedPatterns *[]string `json:"webhookLongLivedPatterns,omitempty"`

	// WebhookDefaultGateScope sets which gate conditions repository webhooks evaluate.
	// +kubebuilder:validation:Enum=all;new_code
	// +optional
	WebhookDefaultGateScope *string `json:"webhookDefaultGateScope,omitempty"`

	// WebhookNonDefaultGateScope sets which gate conditions repository webhooks evaluate.
	// +kubebuilder:validation:Enum=all;new_code
	// +optional
	WebhookNonDefaultGateScope *string `json:"webhookNonDefaultGateScope,omitempty"`

	// WebhookLongLivedGateScope sets which gate conditions repository webhooks evaluate.
	// +kubebuilder:validation:Enum=all;new_code
	// +optional
	WebhookLongLivedGateScope *string `json:"webhookLongLivedGateScope,omitempty"`

	// SQLDialect sets the explicit SQLFluff dialect when enabled.
	// +kubebuilder:validation:MaxLength=20
	// +optional
	SQLDialect *string `json:"sqlDialect,omitempty"`

	// AutoAssign enables automatic issue assignment.
	// +optional
	AutoAssign *bool `json:"autoAssign,omitempty"`

	// IsActive enables the project.
	// +optional
	IsActive *bool `json:"isActive,omitempty"`
}

// CodeProjectStatus contains the observed Argonix state.
type CodeProjectStatus struct {
	ID string `json:"id,omitempty"`
	// LatestAnalysis is the UUID of the most recent analysis.
	LatestAnalysis string `json:"latestAnalysis,omitempty"`
	// Key is the effective project key.
	Key string `json:"key,omitempty"`
	// Name is the effective project name.
	Name         string             `json:"name,omitempty"`
	DateCreated  string             `json:"dateCreated,omitempty"`
	DateModified string             `json:"dateModified,omitempty"`
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Project",type=string,JSONPath=`.status.key`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// CodeProject is the schema for the codeprojects API.
type CodeProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CodeProjectSpec   `json:"spec,omitempty"`
	Status            CodeProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CodeProjectList contains a list of CodeProject objects.
type CodeProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CodeProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CodeProject{}, &CodeProjectList{})
}
