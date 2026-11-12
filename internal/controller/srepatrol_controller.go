package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func buildSREPatrolPayload(obj *v1alpha1.SREPatrol) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{
		"name": s.Name,
	}
	setOptional(payload, "is_active", s.IsActive)
	setOptional(payload, "cadence_minutes", s.CadenceMinutes)
	setOptional(payload, "autonomy", s.Autonomy)
	setOptional(payload, "min_confidence", s.MinConfidence)
	setOptional(payload, "perspective", s.Perspective)
	setOptional(payload, "use_applications", s.UseApplications)
	setOptional(payload, "include_unattributed", s.IncludeUnattributed)
	setOptional(payload, "minimum_attribution_confidence", s.MinimumAttributionConfidence)
	setOptional(payload, "connector_ids", s.ConnectorIDs)
	setOptional(payload, "application_root_ids", s.ApplicationRootIDs)
	setOptional(payload, "notification_channel_ids", s.NotificationChannelIDs)
	setOptional(payload, "allowed_workflow_ids", s.AllowedWorkflowIDs)
	setOptional(payload, "workflow_bindings", s.WorkflowBindings)
	setOptional(payload, "delivery_channel_ids", s.DeliveryChannelIDs)
	setOptional(payload, "alert_channel_ids", s.AlertChannelIDs)
	setOptional(payload, "ticketing_min_severity", s.TicketingMinSeverity)
	setOptional(payload, "allow_mr_authoring", s.AllowMRAuthoring)
	setOptional(payload, "custom_prompt", s.CustomPrompt)

	return payload
}

func SetupSREPatrolReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.SREPatrol]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.SREPatrol { return &v1alpha1.SREPatrol{} },
		Adapter: ResourceAdapter[*v1alpha1.SREPatrol]{
			APIEndpoint:     "/argos/patrols/",
			PatchDriftOnly:  true,
			NormalizeRemote: normalizePatrolRelations,

			BuildPayload:  buildSREPatrolPayload,
			GetResourceID: func(obj *v1alpha1.SREPatrol) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.SREPatrol, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.SREPatrol, data map[string]interface{}) {
				obj.Status.LastRunStatus = getString(data, "last_run_status")
				obj.Status.LastRunDate = getString(data, "last_run_date")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.SREPatrol) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.SREPatrol, conditions []metav1.Condition) { obj.Status.Conditions = conditions },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SREPatrol{}).
		Named("srepatrol").
		Complete(r)
}

func normalizePatrolRelations(desired, remote map[string]interface{}) {
	normalizeArgosRelations(desired, remote)
	normalizeIDList(desired, remote, "allowed_workflow_ids", "allowed_workflows")
	normalizeIDList(desired, remote, "delivery_channel_ids", "delivery_channels")
}
