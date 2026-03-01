package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupAlertChannelReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.AlertChannel]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.AlertChannel { return &v1alpha1.AlertChannel{} },
		Adapter: ResourceAdapter[*v1alpha1.AlertChannel]{
			APIEndpoint: "/alert-channels/",
			BuildPayload: func(obj *v1alpha1.AlertChannel) map[string]interface{} {
				s := obj.Spec
				return map[string]interface{}{
					"name":         s.Name,
					"channel_type": s.ChannelType,
					"config":       s.Config,
					"is_active":    s.IsActive,
				}
			},
			GetResourceID: func(obj *v1alpha1.AlertChannel) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.AlertChannel, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.AlertChannel, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.AlertChannel) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.AlertChannel, c []metav1.Condition) { obj.Status.Conditions = c },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AlertChannel{}).
		Named("alertchannel").
		Complete(r)
}
