package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupTestSuiteReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.TestSuite]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.TestSuite { return &v1alpha1.TestSuite{} },
		Adapter: ResourceAdapter[*v1alpha1.TestSuite]{
			APIEndpoint: "/test-suites/",
			BuildPayload: func(obj *v1alpha1.TestSuite) map[string]interface{} {
				s := obj.Spec
				return map[string]interface{}{
					"name":              s.Name,
					"description":       s.Description,
					"tags":              s.Tags,
					"synthetic_tests":   s.SyntheticTests,
					"manual_test_cases": s.ManualTestCases,
				}
			},
			GetResourceID: func(obj *v1alpha1.TestSuite) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.TestSuite, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.TestSuite, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.TestSuite) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.TestSuite, c []metav1.Condition) { obj.Status.Conditions = c },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.TestSuite{}).
		Named("testsuite").
		Complete(r)
}
