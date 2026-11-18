package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupHealthNotebookReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.HealthNotebook]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.HealthNotebook { return &v1alpha1.HealthNotebook{} },
		Adapter: ResourceAdapter[*v1alpha1.HealthNotebook]{
			APIEndpoint: "/argos/health-notebooks/",
			BuildPayload: func(obj *v1alpha1.HealthNotebook) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":              s.Name,
					"is_active":         s.IsActive,
					"frequency":         s.Frequency,
					"schedule_hour":     s.ScheduleHour,
					"schedule_minute":   s.ScheduleMinute,
					"custom_prompt":     s.CustomPrompt,
					"include_threats":   s.IncludeThreats,
					"include_findings":  s.IncludeFindings,
					"include_monitors":  s.IncludeMonitors,
					"include_incidents": s.IncludeIncidents,
					"include_cost":      s.IncludeCost,
					"ticketing_action":  s.TicketingAction,
				}
				if s.ScheduleDayOfMonth != 0 {
					payload["schedule_day_of_month"] = s.ScheduleDayOfMonth
				}
				if s.TicketingMinSeverity != "" {
					payload["ticketing_min_severity"] = s.TicketingMinSeverity
				}
				parseJSONStringField(s.ScheduleDays, "schedule_days", payload)
				parseJSONStringField(s.ConnectorIDs, "connector_ids", payload)
				parseJSONStringField(s.NotificationChannelIDs, "notification_channel_ids", payload)
				parseJSONStringField(s.AlertChannelIDs, "alert_channel_ids", payload)
				return payload
			},
			GetResourceID: func(obj *v1alpha1.HealthNotebook) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.HealthNotebook, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.HealthNotebook, data map[string]interface{}) {
				obj.Status.LatestScore = getInt(data, "latest_score")
				obj.Status.LatestSeverity = getString(data, "latest_severity")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.HealthNotebook) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.HealthNotebook, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.HealthNotebook{}).
		Named("healthnotebook").
		Complete(r)
}
