package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupSecurityPolicyReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.SecurityPolicy]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.SecurityPolicy { return &v1alpha1.SecurityPolicy{} },
		Adapter: ResourceAdapter[*v1alpha1.SecurityPolicy]{
			APIEndpoint: "/security/policies/",
			BuildPayload: func(obj *v1alpha1.SecurityPolicy) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":        s.Name,
					"description": s.Description,
					"environment": s.Environment,
					"is_active":   s.IsActive,
				}
				// Parse rules JSON string into actual JSON object.
				parseJSONStringField(s.Rules, "rules", payload)
				return payload
			},
			GetResourceID: func(obj *v1alpha1.SecurityPolicy) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.SecurityPolicy, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.SecurityPolicy, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.SecurityPolicy) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.SecurityPolicy, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SecurityPolicy{}).
		Named("securitypolicy").
		Complete(r)
}
