package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupWorkflowReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Workflow]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Workflow { return &v1alpha1.Workflow{} },
		Adapter: ResourceAdapter[*v1alpha1.Workflow]{
			APIEndpoint:    "/argos/workflows/",
			PatchDriftOnly: true,
			ValidateSpec: func(obj *v1alpha1.Workflow) error {
				s := obj.Spec
				if err := validateJSONField("steps", s.Steps, '['); err != nil {
					return err
				}
				if err := validateJSONField("requiredConnectorTypes", s.RequiredConnectorTypes, '['); err != nil {
					return err
				}
				return validateJSONField("inputHints", s.InputHints, '{')
			},
			BuildPayload: func(obj *v1alpha1.Workflow) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":                  s.Name,
					"slug":                  s.Slug,
					"description":           s.Description,
					"category":              s.Category,
					"requires_confirmation": s.RequiresConfirmation,
					"schedule":              s.Schedule,
					"is_active":             s.IsActive,
				}
				parseJSONStringField(s.Steps, "steps", payload)
				parseJSONStringField(s.InputHints, "input_hints", payload)
				parseJSONStringField(s.RequiredConnectorTypes, "required_connector_types", payload)
				setOptional(payload, "extract_inputs_with_ai", s.ExtractInputsWithAI)
				setOptional(payload, "schedule_timezone", s.ScheduleTimezone)
				setOptional(payload, "schedule_inputs", s.ScheduleInputs)
				setOptional(payload, "security_trigger", s.SecurityTrigger)
				setOptional(payload, "security_min_severity", s.SecurityMinSeverity)
				setNullableID(payload, "security_detection_rule", s.SecurityDetectionRule)
				setOptional(payload, "incident_trigger", s.IncidentTrigger)
				setOptional(payload, "incident_autonomy", s.IncidentAutonomy)
				setOptional(payload, "incident_min_severity", s.IncidentMinSeverity)
				setOptional(payload, "incident_min_confidence", s.IncidentMinConfidence)
				setOptional(payload, "incident_max_affected_resources", s.IncidentMaxAffectedResources)
				setOptional(payload, "incident_allowed_environments", s.IncidentAllowedEnvironments)
				setOptional(payload, "incident_allowed_tools", s.IncidentAllowedTools)
				setOptional(payload, "incident_allowed_connector_types", s.IncidentAllowedConnectorTypes)
				setOptional(payload, "incident_require_compensation", s.IncidentRequireCompensation)
				setOptional(payload, "incident_require_verification", s.IncidentRequireVerification)
				setOptional(payload, "incident_stop_on_ambiguous_input", s.IncidentStopOnAmbiguousInput)
				setOptional(payload, "incident_timeout_minutes", s.IncidentTimeoutMinutes)
				return payload
			},
			GetResourceID: func(obj *v1alpha1.Workflow) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Workflow, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Workflow, data map[string]interface{}) {
				obj.Status.ScheduleLastRunAt = getString(data, "schedule_last_run_at")
				obj.Status.ScheduleNextRunAt = getString(data, "schedule_next_run_at")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Workflow) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.Workflow, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Workflow{}).
		Named("workflow").
		Complete(r)
}
