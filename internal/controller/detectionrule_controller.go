package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupDetectionRuleReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.DetectionRule]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.DetectionRule { return &v1alpha1.DetectionRule{} },
		Adapter: ResourceAdapter[*v1alpha1.DetectionRule]{
			APIEndpoint: "/security/detection-rules/",
			BuildPayload: func(obj *v1alpha1.DetectionRule) map[string]interface{} {
				s := obj.Spec
				payload := map[string]interface{}{
					"name":            s.Name,
					"description":     s.Description,
					"rule_type":       s.RuleType,
					"severity":        s.Severity,
					"is_active":       s.IsActive,
					"mitre_tactic":    s.MitreTactic,
					"mitre_technique": s.MitreTechnique,
				}
				// Parse config JSON string into actual JSON object.
				parseJSONStringField(s.Config, "config", payload)
				return payload
			},
			GetResourceID: func(obj *v1alpha1.DetectionRule) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.DetectionRule, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.DetectionRule, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.DetectionRule) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.DetectionRule, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DetectionRule{}).
		Named("detectionrule").
		Complete(r)
}
