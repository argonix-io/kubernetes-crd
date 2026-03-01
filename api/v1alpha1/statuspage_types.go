package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatusPageSpec defines the desired state of a StatusPage.
type StatusPageSpec struct {
	// Name is the display name of the status page.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Slug is the URL slug (must be unique).
	// +kubebuilder:validation:Required
	Slug string `json:"slug"`

	// CustomDomain is an optional custom domain for the status page.
	// +optional
	CustomDomain string `json:"customDomain,omitempty"`

	// Visibility controls who can see the page: public or private.
	// +kubebuilder:default="public"
	// +kubebuilder:validation:Enum=public;private
	// +optional
	Visibility string `json:"visibility,omitempty"`

	// LogoURL is the URL to the logo image.
	// +optional
	LogoURL string `json:"logoUrl,omitempty"`

	// AccentColor is the hex color code for the accent color.
	// +kubebuilder:default="#3B82F6"
	// +optional
	AccentColor string `json:"accentColor,omitempty"`

	// CustomCSS is custom CSS to apply to the status page.
	// +optional
	CustomCSS string `json:"customCss,omitempty"`

	// ShowHealthGraph indicates whether to display the health graph.
	// +kubebuilder:default=true
	// +optional
	ShowHealthGraph bool `json:"showHealthGraph,omitempty"`

	// IsActive indicates whether the status page is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// StatusPageStatus defines the observed state of a StatusPage.
type StatusPageStatus struct {
	// ID is the UUID of the status page in the Argonix API.
	ID string `json:"id,omitempty"`

	// PublicURL is the full URL of the status page.
	PublicURL string `json:"publicUrl,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Slug",type=string,JSONPath=`.spec.slug`
// +kubebuilder:printcolumn:name="Visibility",type=string,JSONPath=`.spec.visibility`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// StatusPage is the Schema for the statuspages API.
type StatusPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StatusPageSpec   `json:"spec,omitempty"`
	Status            StatusPageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StatusPageList contains a list of StatusPage.
type StatusPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatusPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatusPage{}, &StatusPageList{})
}
