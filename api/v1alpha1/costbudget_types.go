package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CostBudgetSpec defines the desired state of a FinOps CostBudget — a
// spend threshold over a cost slice that fires alerts at given percentages.
type CostBudgetSpec struct {
	// Name is the display name of the budget.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Period is the budget window.
	// +kubebuilder:validation:Enum=monthly;quarterly
	// +kubebuilder:default=monthly
	// +optional
	Period string `json:"period,omitempty"`

	// Amount is the budget amount in currency units (e.g. 5000.00).
	// +kubebuilder:validation:Required
	Amount string `json:"amount"`

	// Currency is the ISO currency code of the budget amount.
	// +kubebuilder:default=EUR
	// +optional
	Currency string `json:"currency,omitempty"`

	// Filters is a JSON-encoded object scoping the budget to a cost slice.
	// Allowed keys: connector, tenant, provider, service, k8s_cluster,
	// k8s_namespace, tag_key, tag_value. Empty means the whole org.
	// +optional
	Filters string `json:"filters,omitempty"`

	// AlertThresholdsPct is a JSON-encoded list of % thresholds at which to
	// fire alerts (e.g. [80, 100, 120]).
	// +optional
	AlertThresholdsPct string `json:"alertThresholdsPct,omitempty"`

	// AlertChannelIDs is a JSON-encoded list of alert channel UUIDs.
	// +optional
	AlertChannelIDs string `json:"alertChannelIds,omitempty"`

	// IsActive indicates whether the budget is enforced.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`
}

// CostBudgetStatus defines the observed state of a CostBudget.
type CostBudgetStatus struct {
	// ID is the UUID of the budget in the Argonix API.
	ID string `json:"id,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateUpdated is the last modification timestamp.
	DateUpdated string `json:"dateUpdated,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Period",type=string,JSONPath=`.spec.period`
// +kubebuilder:printcolumn:name="Amount",type=string,JSONPath=`.spec.amount`
// +kubebuilder:printcolumn:name="Currency",type=string,JSONPath=`.spec.currency`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// CostBudget is the Schema for the FinOps budgets API.
type CostBudget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CostBudgetSpec   `json:"spec,omitempty"`
	Status            CostBudgetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CostBudgetList contains a list of CostBudget.
type CostBudgetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CostBudget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CostBudget{}, &CostBudgetList{})
}
