package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupEnvironmentReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Environment]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Environment { return &v1alpha1.Environment{} },
		Adapter: ResourceAdapter[*v1alpha1.Environment]{
			APIEndpoint: "/environments/",
			BuildPayload: func(obj *v1alpha1.Environment) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":       s.Name,
					"is_default": s.IsDefault,
				}
				if len(s.Variables) > 0 {
					varsJSON, err := json.Marshal(s.Variables)
					if err == nil {
						var vars interface{}
						json.Unmarshal(varsJSON, &vars)
						payload["variables"] = vars
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Environment) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Environment, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Environment, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Environment) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.Environment, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Environment{}).
		Named("environment").
		Complete(r)
}
