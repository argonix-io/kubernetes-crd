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
			APIEndpoint:     "/argos/health-notebooks/",
			PatchDriftOnly:  true,
			NormalizeRemote: normalizeArgosRelations,
			ValidateSpec: func(obj *v1alpha1.HealthNotebook) error {
				s := obj.Spec
				for _, field := range []struct{ name, value string }{
					{"scheduleDays", s.ScheduleDays},
					{"connectorIds", s.ConnectorIDs},
					{"notificationChannelIds", s.NotificationChannelIDs},
					{"alertChannelIds", s.AlertChannelIDs},
				} {
					if err := validateJSONField(field.name, field.value, '['); err != nil {
						return err
					}
				}
				return nil
			},
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
				payload["ticketing_min_severity"] = s.TicketingMinSeverity
				parseJSONStringField(s.ScheduleDays, "schedule_days", payload)
				parseJSONStringField(s.ConnectorIDs, "connector_ids", payload)
				parseJSONStringField(s.NotificationChannelIDs, "notification_channel_ids", payload)
				parseJSONStringField(s.AlertChannelIDs, "alert_channel_ids", payload)
				setOptional(payload, "report_kind", s.ReportKind)
				setOptional(payload, "use_applications", s.UseApplications)
				setOptional(payload, "perspective", s.Perspective)
				setOptional(payload, "application_root_ids", s.ApplicationRootIDs)
				setOptional(payload, "application_label_selector", s.ApplicationLabelSelector)
				setOptional(payload, "include_unattributed", s.IncludeUnattributed)
				setOptional(payload, "minimum_attribution_confidence", s.MinimumAttributionConfidence)
				setOptional(payload, "lookback_hours", s.LookbackHours)
				setOptional(payload, "include_log_anomalies", s.IncludeLogAnomalies)
				setOptional(payload, "include_vulnerabilities", s.IncludeVulnerabilities)
				setOptional(payload, "include_runtime", s.IncludeRuntime)
				setOptional(payload, "include_dns", s.IncludeDNS)
				setOptional(payload, "include_network", s.IncludeNetwork)
				setOptional(payload, "include_operational_signals", s.IncludeOperationalSignals)
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
