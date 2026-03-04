package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupConnectorReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Connector]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Connector { return &v1alpha1.Connector{} },
		Adapter: ResourceAdapter[*v1alpha1.Connector]{
			APIEndpoint: "/argos/connectors/",
			BuildPayload: func(obj *v1alpha1.Connector) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":           s.Name,
					"connector_type": s.ConnectorType,
					"config":         s.Config,
					"is_active":      s.IsActive,
				}
				if s.Capabilities != "" {
					var caps interface{}
					if err := json.Unmarshal([]byte(s.Capabilities), &caps); err == nil {
						payload["capabilities"] = caps
					}
				}
				if s.Tags != "" {
					var tags interface{}
					if err := json.Unmarshal([]byte(s.Tags), &tags); err == nil {
						payload["tags"] = tags
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Connector) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Connector, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Connector, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Connector) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.Connector, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Connector{}).
		Named("connector").
		Complete(r)
}
