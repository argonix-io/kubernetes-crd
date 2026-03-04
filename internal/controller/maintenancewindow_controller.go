package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupMaintenanceWindowReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.MaintenanceWindow]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.MaintenanceWindow { return &v1alpha1.MaintenanceWindow{} },
		Adapter: ResourceAdapter[*v1alpha1.MaintenanceWindow]{
			APIEndpoint: "/maintenance-windows/",
			BuildPayload: func(obj *v1alpha1.MaintenanceWindow) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":      s.Name,
					"group_id":  s.GroupID,
					"repeat":    s.Repeat,
					"is_active": s.IsActive,
				}
				if s.StartsAt != "" {
					payload["starts_at"] = s.StartsAt
				}
				if s.EndsAt != "" {
					payload["ends_at"] = s.EndsAt
				}
				if s.TimeFrom != "" {
					payload["time_from"] = s.TimeFrom
				}
				if s.TimeTo != "" {
					payload["time_to"] = s.TimeTo
				}
				if s.Weekdays != "" {
					payload["weekdays"] = s.Weekdays
				}
				if s.DayOfMonth != 0 {
					payload["day_of_month"] = s.DayOfMonth
				}
				if s.CronExpression != "" {
					payload["cron_expression"] = s.CronExpression
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.MaintenanceWindow) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.MaintenanceWindow, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.MaintenanceWindow, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
				obj.Status.ScheduleSummary = getString(data, "schedule_summary")
			},
			GetConditions: func(obj *v1alpha1.MaintenanceWindow) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.MaintenanceWindow, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MaintenanceWindow{}).
		Named("maintenancewindow").
		Complete(r)
}
