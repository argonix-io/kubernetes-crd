package controller

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func buildScanTargetPayload(obj *v1alpha1.ScanTarget) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{
		"kind":             s.Kind,
		"enabled_scanners": s.EnabledScanners,
	}
	setOptional(payload, "name", s.Name)
	setNullableID(payload, "connector", s.Connector)
	setNullableID(payload, "project", s.Project)
	setOptional(payload, "branch", s.Branch)
	setOptional(payload, "platform", s.Platform)
	setOptional(payload, "frequency", s.Frequency)
	setOptional(payload, "schedule_hour", s.ScheduleHour)
	setOptional(payload, "schedule_minute", s.ScheduleMinute)
	setOptional(payload, "schedule_days", s.ScheduleDays)
	setOptional(payload, "is_active", s.IsActive)
	setOptional(payload, "webhook_enabled", s.WebhookEnabled)
	setOptional(payload, "enforce_merge_requests", s.EnforceMergeRequests)
	identifier := strings.TrimSpace(s.Identifier)
	payload["identifier"] = identifier
	if s.GateEnvironment != nil {
		payload["gate_environment"] = strings.TrimSpace(*s.GateEnvironment)
	}
	if s.DefaultRef != nil {
		ref := strings.TrimSpace(*s.DefaultRef)
		if s.Kind == "image" {
			suffix := identifier[strings.LastIndex(identifier, "/")+1:]
			switch {
			case strings.Contains(identifier, "@") || strings.Contains(suffix, ":"):
				ref = ""
			case ref == "":
				ref = "latest"
			}
		}
		payload["default_ref"] = ref
	}
	if s.Frequency != nil && *s.Frequency == "weekly" && s.ScheduleDays != nil && len(*s.ScheduleDays) == 0 {
		payload["schedule_days"] = []int{0}
	}
	return payload
}

func SetupScanTargetReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.ScanTarget]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.ScanTarget { return &v1alpha1.ScanTarget{} },
		Adapter: ResourceAdapter[*v1alpha1.ScanTarget]{
			APIEndpoint:     "/security/scan-targets/",
			PatchDriftOnly:  true,
			NormalizeRemote: normalizeScanTargetSchedule,

			BuildPayload:  buildScanTargetPayload,
			GetResourceID: func(obj *v1alpha1.ScanTarget) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.ScanTarget, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.ScanTarget, data map[string]interface{}) {
				obj.Status.LastDispatchedAt = getString(data, "last_dispatched_at")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.ScanTarget) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.ScanTarget, conditions []metav1.Condition) { obj.Status.Conditions = conditions },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ScanTarget{}).
		Named("scantarget").
		Complete(r)
}

func normalizeScanTargetSchedule(desired, remote map[string]interface{}) {
	// The API defaults an empty weekly day selection to Monday, including
	// when the frequency was set outside this manifest.
	if _, managed := desired["frequency"]; !managed && remote["frequency"] == "weekly" {
		if days, ok := desired["schedule_days"].([]int); ok && len(days) == 0 {
			desired["schedule_days"] = []int{0}
		}
	}
}
