package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupKnowledgeBaseReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.KnowledgeBase]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.KnowledgeBase { return &v1alpha1.KnowledgeBase{} },
		Adapter: ResourceAdapter[*v1alpha1.KnowledgeBase]{
			APIEndpoint: "/argos/knowledge-bases/",
			BuildPayload: func(obj *v1alpha1.KnowledgeBase) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":        s.Name,
					"source_type": s.SourceType,
					"is_active":   s.IsActive,
				}
				if s.ConnectorID != "" {
					payload["connector_id"] = s.ConnectorID
				}
				if s.SyncConfig != "" {
					var sc interface{}
					if err := json.Unmarshal([]byte(s.SyncConfig), &sc); err == nil {
						payload["sync_config"] = sc
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.KnowledgeBase) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.KnowledgeBase, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.KnowledgeBase, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
				obj.Status.LastSyncedAt = getString(data, "last_synced_at")
				if dc, ok := data["document_count"].(float64); ok {
					obj.Status.DocumentCount = int64(dc)
				}
			},
			GetConditions: func(obj *v1alpha1.KnowledgeBase) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.KnowledgeBase, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.KnowledgeBase{}).
		Named("knowledgebase").
		Complete(r)
}
