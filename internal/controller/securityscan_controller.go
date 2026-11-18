package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupSecurityScanReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.SecurityScan]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.SecurityScan { return &v1alpha1.SecurityScan{} },
		Adapter: ResourceAdapter[*v1alpha1.SecurityScan]{
			APIEndpoint: "/security/scans/",
			BuildPayload: func(obj *v1alpha1.SecurityScan) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"scan_type": s.ScanType,
				}
				if s.Connector != "" {
					payload["connector"] = s.Connector
				}
				if s.TriggeredBy != "" {
					payload["triggered_by"] = s.TriggeredBy
				}
				return payload
			},
			GetResourceID: func(obj *v1alpha1.SecurityScan) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.SecurityScan, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.SecurityScan, data map[string]interface{}) {
				obj.Status.Status = getString(data, "status")
				if v, ok := data["findings_count"].(float64); ok {
					obj.Status.FindingsCount = int64(v)
				}
				obj.Status.StartedAt = getString(data, "started_at")
				obj.Status.CompletedAt = getString(data, "completed_at")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.SecurityScan) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.SecurityScan, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SecurityScan{}).
		Named("securityscan").
		Complete(r)
}
