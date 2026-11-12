package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupSyntheticTestReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.SyntheticTest]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.SyntheticTest { return &v1alpha1.SyntheticTest{} },
		Adapter: ResourceAdapter[*v1alpha1.SyntheticTest]{
			APIEndpoint:    "/synthetic-tests/",
			PatchDriftOnly: true,
			NormalizeRemote: func(desired, remote map[string]interface{}) {
				if id, exists := remote["resolved_resource"]; exists {
					remote["resource_id"] = id
				}
			},
			ValidateSpec: func(obj *v1alpha1.SyntheticTest) error {
				return validateJSONField("steps", obj.Spec.Steps, '[')
			},
			BuildPayload:  buildSyntheticTestPayload,
			GetResourceID: func(obj *v1alpha1.SyntheticTest) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.SyntheticTest, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.SyntheticTest, data map[string]interface{}) {
				obj.Status.CurrentStatus = getString(data, "current_status")
				obj.Status.NextRunAt = getString(data, "next_run_at")
				obj.Status.LastRunAt = getString(data, "last_run_at")
				obj.Status.LastWarningCount = getInt(data, "last_warning_count")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.SyntheticTest) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.SyntheticTest, c []metav1.Condition) { obj.Status.Conditions = c },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SyntheticTest{}).
		Named("synthetictest").
		Complete(r)
}

func buildSyntheticTestPayload(obj *v1alpha1.SyntheticTest) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{
		"name":           s.Name,
		"description":    s.Description,
		"is_active":      s.IsActive,
		"test_type":      s.TestType,
		"check_interval": s.CheckInterval,
		"timeout":        s.Timeout,
		"tags":           append([]string{}, s.Tags...),
		"locations":      append([]string{}, s.Locations...),
	}
	parseJSONStringField(s.Steps, "steps", payload)
	setNullableID(payload, "resource_id", s.ResourceID)
	setNullableID(payload, "application_key", s.ApplicationKey)
	setOptional(payload, "input_variables", s.InputVariables)
	setOptional(payload, "allow_mutating_actions", s.AllowMutatingActions)
	setOptional(payload, "run_timeout", s.RunTimeout)
	setOptional(payload, "schedule_mode", s.ScheduleMode)
	setOptional(payload, "schedule", s.Schedule)
	setOptional(payload, "schedule_timezone", s.ScheduleTimezone)
	return payload
}
