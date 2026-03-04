package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupPersonaReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Persona]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Persona { return &v1alpha1.Persona{} },
		Adapter: ResourceAdapter[*v1alpha1.Persona]{
			APIEndpoint: "/argos/personas/",
			BuildPayload: func(obj *v1alpha1.Persona) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":          s.Name,
					"description":   s.Description,
					"template":      s.Template,
					"system_prompt": s.SystemPrompt,
					"is_active":     s.IsActive,
				}
				if s.ConnectorIDs != "" {
					var ids interface{}
					if err := json.Unmarshal([]byte(s.ConnectorIDs), &ids); err == nil {
						payload["connector_ids"] = ids
					}
				}
				if s.AllowedTools != "" {
					var tools interface{}
					if err := json.Unmarshal([]byte(s.AllowedTools), &tools); err == nil {
						payload["allowed_tools"] = tools
					}
				}
				if s.ApprovalRules != "" {
					var rules interface{}
					if err := json.Unmarshal([]byte(s.ApprovalRules), &rules); err == nil {
						payload["approval_rules"] = rules
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Persona) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Persona, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Persona, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Persona) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.Persona, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Persona{}).
		Named("persona").
		Complete(r)
}
