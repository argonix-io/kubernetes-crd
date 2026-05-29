package controller

import (
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupTenantReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Tenant]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Tenant { return &v1alpha1.Tenant{} },
		Adapter: ResourceAdapter[*v1alpha1.Tenant]{
			APIEndpoint: "/finops/tenants/",
			BuildPayload: func(obj *v1alpha1.Tenant) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":          s.Name,
					"kind":          s.Kind,
					"color":         s.Color,
					"contact_email": s.ContactEmail,
					"description":   s.Description,
				}
				if s.MonthlyBudget != "" {
					if v, err := strconv.ParseFloat(s.MonthlyBudget, 64); err == nil {
						payload["monthly_budget"] = v
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Tenant) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Tenant, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Tenant, data map[string]interface{}) {
				obj.Status.Slug = getString(data, "slug")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Tenant) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.Tenant, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Tenant{}).
		Named("tenant").
		Complete(r)
}
