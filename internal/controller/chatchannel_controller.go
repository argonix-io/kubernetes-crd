package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupChatChannelReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.ChatChannel]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.ChatChannel { return &v1alpha1.ChatChannel{} },
		Adapter: ResourceAdapter[*v1alpha1.ChatChannel]{
			APIEndpoint: "/argos/chat-channels/",
			BuildPayload: func(obj *v1alpha1.ChatChannel) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"channel_type": s.ChannelType,
					"channel_id":   s.ChannelID,
					"channel_name": s.ChannelName,
					"connector_id": s.ConnectorID,
					"is_active":    s.IsActive,
				}
				if s.PersonaID != "" {
					payload["persona_id"] = s.PersonaID
				}
				if s.Config != "" {
					var cfg interface{}
					if err := json.Unmarshal([]byte(s.Config), &cfg); err == nil {
						payload["config"] = cfg
					}
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.ChatChannel) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.ChatChannel, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.ChatChannel, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.ChatChannel) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.ChatChannel, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ChatChannel{}).
		Named("chatchannel").
		Complete(r)
}
