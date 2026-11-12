package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func buildMonitorGroupPayload(obj *v1alpha1.MonitorGroup) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{
		"name": s.Name,
	}
	setOptional(payload, "description", s.Description)

	return payload
}

func SetupMonitorGroupReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.MonitorGroup]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.MonitorGroup { return &v1alpha1.MonitorGroup{} },
		Adapter: ResourceAdapter[*v1alpha1.MonitorGroup]{
			APIEndpoint:    "/monitor-groups/",
			PatchDriftOnly: true,

			BuildPayload:  buildMonitorGroupPayload,
			GetResourceID: func(obj *v1alpha1.MonitorGroup) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.MonitorGroup, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.MonitorGroup, data map[string]interface{}) {

				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.MonitorGroup) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.MonitorGroup, conditions []metav1.Condition) { obj.Status.Conditions = conditions },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MonitorGroup{}).
		Named("monitorgroup").
		Complete(r)
}
