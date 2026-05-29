package controller

import (
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupCostBudgetReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.CostBudget]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.CostBudget { return &v1alpha1.CostBudget{} },
		Adapter: ResourceAdapter[*v1alpha1.CostBudget]{
			APIEndpoint: "/finops/budgets/",
			BuildPayload: func(obj *v1alpha1.CostBudget) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":      s.Name,
					"period":    s.Period,
					"currency":  s.Currency,
					"is_active": s.IsActive,
				}
				if s.Amount != "" {
					if v, err := strconv.ParseFloat(s.Amount, 64); err == nil {
						payload["amount"] = v
					}
				}
				parseJSONStringField(s.Filters, "filters", payload)
				parseJSONStringField(s.AlertThresholdsPct, "alert_thresholds_pct", payload)
				parseJSONStringField(s.AlertChannelIDs, "alert_channels", payload)
				return payload
			},
			GetResourceID: func(obj *v1alpha1.CostBudget) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.CostBudget, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.CostBudget, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateUpdated = getString(data, "date_updated")
			},
			GetConditions: func(obj *v1alpha1.CostBudget) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.CostBudget, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.CostBudget{}).
		Named("costbudget").
		Complete(r)
}
