package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func buildPeriodicJobPayload(obj *v1alpha1.PeriodicJob) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{
		"name": s.Name,
	}
	setOptional(payload, "kind", s.Kind)
	setOptional(payload, "prompt", s.Prompt)
	setOptional(payload, "schedule_hour", s.ScheduleHour)
	setOptional(payload, "schedule_minute", s.ScheduleMinute)
	setOptional(payload, "schedule_days", s.ScheduleDays)
	setOptional(payload, "connector_ids", s.ConnectorIDs)
	setOptional(payload, "application_root_ids", s.ApplicationRootIDs)
	setOptional(payload, "notification_channel_ids", s.NotificationChannelIDs)
	setOptional(payload, "alert_channel_ids", s.AlertChannelIDs)
	setOptional(payload, "auto_create_ticket", s.AutoCreateTicket)
	setOptional(payload, "deliverable_retention_action", s.DeliverableRetentionAction)
	setOptional(payload, "deliverable_retention_destination", s.DeliverableRetentionDestination)
	setOptional(payload, "is_active", s.IsActive)
	setOptional(payload, "mode", s.Mode)
	setOptional(payload, "timeout_seconds", s.TimeoutSeconds)
	if s.DeliverableRetentionDays != nil {
		payload["deliverable_retention_days"] = nil
		if *s.DeliverableRetentionDays != 0 {
			payload["deliverable_retention_days"] = *s.DeliverableRetentionDays
		}
	}
	return payload
}

func SetupPeriodicJobReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.PeriodicJob]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.PeriodicJob { return &v1alpha1.PeriodicJob{} },
		Adapter:       periodicJobAdapter(),
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PeriodicJob{}).
		Named("periodicjob").
		Complete(r)
}

func periodicJobAdapter() ResourceAdapter[*v1alpha1.PeriodicJob] {
	return ResourceAdapter[*v1alpha1.PeriodicJob]{
		APIEndpoint:     "/argos/jobs/",
		PatchDriftOnly:  true,
		NormalizeRemote: normalizeArgosRelations,

		BuildPayload:  buildPeriodicJobPayload,
		GetResourceID: func(obj *v1alpha1.PeriodicJob) string { return obj.Status.ID },
		SetResourceID: func(obj *v1alpha1.PeriodicJob, id string) { obj.Status.ID = id },
		SetStatusFromResponse: func(obj *v1alpha1.PeriodicJob, data map[string]interface{}) {
			obj.Status.LatestRunStatus = getString(data, "latest_run_status")
			obj.Status.LatestRunDate = getString(data, "latest_run_date")
			obj.Status.DateCreated = getString(data, "date_created")
			obj.Status.DateModified = getString(data, "date_modified")
		},
		GetConditions: func(obj *v1alpha1.PeriodicJob) []metav1.Condition { return obj.Status.Conditions },
		SetConditions: func(obj *v1alpha1.PeriodicJob, conditions []metav1.Condition) { obj.Status.Conditions = conditions },
	}
}
