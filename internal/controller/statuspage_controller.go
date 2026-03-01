package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupStatusPageReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.StatusPage]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.StatusPage { return &v1alpha1.StatusPage{} },
		Adapter: ResourceAdapter[*v1alpha1.StatusPage]{
			APIEndpoint: "/status-pages/",
			BuildPayload: func(obj *v1alpha1.StatusPage) map[string]interface{} {
				s := obj.Spec
				return map[string]interface{}{
					"name":              s.Name,
					"slug":              s.Slug,
					"custom_domain":     s.CustomDomain,
					"visibility":        s.Visibility,
					"logo_url":          s.LogoURL,
					"accent_color":      s.AccentColor,
					"custom_css":        s.CustomCSS,
					"show_health_graph": s.ShowHealthGraph,
					"is_active":         s.IsActive,
				}
			},
			GetResourceID: func(obj *v1alpha1.StatusPage) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.StatusPage, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.StatusPage, data map[string]interface{}) {
				obj.Status.PublicURL = getString(data, "public_url")
			},
			GetConditions: func(obj *v1alpha1.StatusPage) []metav1.Condition {
				return obj.Status.Conditions
			},
			SetConditions: func(obj *v1alpha1.StatusPage, c []metav1.Condition) {
				obj.Status.Conditions = c
			},
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.StatusPage{}).
		Named("statuspage").
		Complete(r)
}
