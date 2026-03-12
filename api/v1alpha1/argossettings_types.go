package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ArgosSettingsSpec defines the desired state of ArgosSettings.
type ArgosSettingsSpec struct {
	// LLMProvider is the AI model provider: local, google, anthropic, or openai.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=local;google;anthropic;openai
	// +kubebuilder:default=google
	LLMProvider string `json:"llmProvider"`

	// LLMModel is the model override (blank = default for the provider).
	// Any model name is accepted when llmBaseURL is set.
	// +optional
	LLMModel string `json:"llmModel,omitempty"`

	// LLMApiKeySecretRef references a Kubernetes Secret containing the custom API key.
	// The key in the Secret must be "api-key".
	// +optional
	LLMApiKeySecretRef *SecretKeyRef `json:"llmApiKeySecretRef,omitempty"`

	// LLMBaseURL is a custom base URL for the LLM endpoint (e.g. vLLM, Azure OpenAI, self-hosted Ollama).
	// +optional
	LLMBaseURL string `json:"llmBaseURL,omitempty"`

	// CustomInstructions are prepended to every Argos conversation system prompt.
	// +optional
	CustomInstructions string `json:"customInstructions,omitempty"`

	// DemoMode enables scripted demo responses instead of calling the LLM.
	// +kubebuilder:default=false
	// +optional
	DemoMode bool `json:"demoMode,omitempty"`
}

// SecretKeyRef references a key within a Kubernetes Secret.
type SecretKeyRef struct {
	// Name of the Secret.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Key within the Secret data. Defaults to "api-key".
	// +kubebuilder:default="api-key"
	// +optional
	Key string `json:"key,omitempty"`
}

// ArgosSettingsStatus defines the observed state of ArgosSettings.
type ArgosSettingsStatus struct {
	// ID is the UUID of the settings in the Argonix API.
	ID string `json:"id,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.llmProvider`
// +kubebuilder:printcolumn:name="Model",type=string,JSONPath=`.spec.llmModel`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ArgosSettings is the Schema for the argossettings API.
type ArgosSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ArgosSettingsSpec   `json:"spec,omitempty"`
	Status            ArgosSettingsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArgosSettingsList contains a list of ArgosSettings.
type ArgosSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArgosSettings `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArgosSettings{}, &ArgosSettingsList{})
}
