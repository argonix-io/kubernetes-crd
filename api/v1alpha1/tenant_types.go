package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantSpec defines the desired state of a FinOps Tenant — a logical
// cost owner (business unit, team, product, environment, client).
type TenantSpec struct {
	// Name is the display name of the tenant.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Kind is the category of the tenant.
	// +kubebuilder:validation:Enum=business_unit;client;team;product;environment;system
	// +kubebuilder:default=team
	// +optional
	Kind string `json:"kind,omitempty"`

	// Color is a hex color used to render the tenant in dashboards.
	// +optional
	Color string `json:"color,omitempty"`

	// MonthlyBudget is the monthly budget in the org currency (e.g. 1500.00).
	// +optional
	MonthlyBudget string `json:"monthlyBudget,omitempty"`

	// ContactEmail is the owner contact for budget alerts.
	// +optional
	ContactEmail string `json:"contactEmail,omitempty"`

	// Description is a free-form description of the tenant.
	// +optional
	Description string `json:"description,omitempty"`
}

// TenantStatus defines the observed state of a Tenant.
type TenantStatus struct {
	// ID is the UUID of the tenant in the Argonix API.
	ID string `json:"id,omitempty"`

	// Slug is the API-generated slug for the tenant.
	Slug string `json:"slug,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Kind",type=string,JSONPath=`.spec.kind`
// +kubebuilder:printcolumn:name="Slug",type=string,JSONPath=`.status.slug`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Tenant is the Schema for the FinOps tenants API.
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TenantSpec   `json:"spec,omitempty"`
	Status            TenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant.
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
