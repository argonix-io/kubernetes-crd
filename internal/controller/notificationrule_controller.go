package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupNotificationRuleReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.NotificationRule]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.NotificationRule { return &v1alpha1.NotificationRule{} },
		Adapter: ResourceAdapter[*v1alpha1.NotificationRule]{
			APIEndpoint: "/notification-rules/",
			BuildPayload: func(obj *v1alpha1.NotificationRule) map[string]interface{} {
				s := obj.Spec
				return map[string]interface{}{
					"name":                 s.Name,
					"trigger_condition":    s.TriggerCondition,
					"consecutive_failures": s.ConsecutiveFailures,
					"cooldown_minutes":     s.CooldownMinutes,
					"all_monitors":         s.AllMonitors,
					"monitor_tags":         s.MonitorTags,
					"monitors":             s.Monitors,
					"synthetic_tests":      s.SyntheticTests,
					"channels":             s.Channels,
					"auto_investigate":     s.AutoInvestigate,
				}
			},
			GetResourceID: func(obj *v1alpha1.NotificationRule) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.NotificationRule, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.NotificationRule, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.NotificationRule) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.NotificationRule, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NotificationRule{}).
		Named("notificationrule").
		Complete(r)
}
