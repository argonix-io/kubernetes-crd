package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupTestPlanReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.TestPlan]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.TestPlan { return &v1alpha1.TestPlan{} },
		Adapter: ResourceAdapter[*v1alpha1.TestPlan]{
			APIEndpoint: "/test-plans/",
			BuildPayload: func(obj *v1alpha1.TestPlan) map[string]interface{} {
				s := obj.Spec
				p := map[string]interface{}{
					"name":        s.Name,
					"description": s.Description,
					"suites":      s.Suites,
					"tags":        s.Tags,
				}
				if s.EndDate != "" {
					p["end_date"] = s.EndDate
				}
				return p
			},
			GetResourceID: func(obj *v1alpha1.TestPlan) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.TestPlan, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.TestPlan, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.TestPlan) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.TestPlan, c []metav1.Condition) { obj.Status.Conditions = c },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.TestPlan{}).
		Named("testplan").
		Complete(r)
}
