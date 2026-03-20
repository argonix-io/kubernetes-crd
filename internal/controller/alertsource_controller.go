package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupAlertSourceReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.AlertSource]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.AlertSource { return &v1alpha1.AlertSource{} },
		Adapter: ResourceAdapter[*v1alpha1.AlertSource]{
			APIEndpoint: "/alert-sources/",
			BuildPayload: func(obj *v1alpha1.AlertSource) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":                 s.Name,
					"source_type":          s.SourceType,
					"is_active":            s.IsActive,
					"auto_investigate":     s.AutoInvestigate,
					"auto_remediate":       s.AutoRemediate,
				}
				if s.Connector != "" {
					payload["connector"] = s.Connector
				}
				if s.Filters != "" {
					payload["filters"] = s.Filters
				}
				if len(s.Channels) > 0 {
					payload["channels"] = s.Channels
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.AlertSource) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.AlertSource, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.AlertSource, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
				obj.Status.WebhookSecret = getString(data, "webhook_secret")
				obj.Status.WebhookURL = getString(data, "webhook_url")
				obj.Status.LastReceivedAt = getString(data, "last_received_at")
				if v, ok := data["total_received"].(float64); ok {
					obj.Status.TotalReceived = int64(v)
				}
			},
			GetConditions: func(obj *v1alpha1.AlertSource) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.AlertSource, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AlertSource{}).
		Named("alertsource").
		Complete(r)
}
