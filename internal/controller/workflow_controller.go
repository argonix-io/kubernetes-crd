package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupWorkflowReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Workflow]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Workflow { return &v1alpha1.Workflow{} },
		Adapter: ResourceAdapter[*v1alpha1.Workflow]{
			APIEndpoint: "/argos/workflows/",
			BuildPayload: func(obj *v1alpha1.Workflow) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":                  s.Name,
					"slug":                  s.Slug,
					"description":           s.Description,
					"category":              s.Category,
					"requires_confirmation": s.RequiresConfirmation,
					"schedule":              s.Schedule,
					"is_active":             s.IsActive,
				}
				if s.Steps != "" {
					var steps interface{}
					if err := json.Unmarshal([]byte(s.Steps), &steps); err == nil {
						payload["steps"] = steps
					}
				}
				if s.RequiredConnectorTypes != "" {
					var rct interface{}
					if err := json.Unmarshal([]byte(s.RequiredConnectorTypes), &rct); err == nil {
						payload["required_connector_types"] = rct
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Workflow) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Workflow, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Workflow, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Workflow) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.Workflow, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Workflow{}).
		Named("workflow").
		Complete(r)
}
