package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupManualTestCaseReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.ManualTestCase]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.ManualTestCase { return &v1alpha1.ManualTestCase{} },
		Adapter: ResourceAdapter[*v1alpha1.ManualTestCase]{
			APIEndpoint: "/manual-test-cases/",
			BuildPayload: func(obj *v1alpha1.ManualTestCase) map[string]interface{} {
				s := obj.Spec
				stepsJSON, _ := json.Marshal(s.Steps)
				return map[string]interface{}{
					"title":         s.Title,
					"description":   s.Description,
					"preconditions": s.Preconditions,
					"steps":         string(stepsJSON),
					"priority":      s.Priority,
					"tags":          s.Tags,
				}
			},
			GetResourceID: func(obj *v1alpha1.ManualTestCase) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.ManualTestCase, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.ManualTestCase, data map[string]interface{}) {
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.ManualTestCase) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.ManualTestCase, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ManualTestCase{}).
		Named("manualtestcase").
		Complete(r)
}
