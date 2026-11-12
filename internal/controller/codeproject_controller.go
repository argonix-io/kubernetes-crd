package controller

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func buildCodeProjectPayload(obj *v1alpha1.CodeProject) map[string]interface{} {
	s := obj.Spec
	payload := map[string]interface{}{}
	setOptional(payload, "name", s.Name)
	setOptional(payload, "key", s.Key)
	setNullableID(payload, "connector", s.Connector)
	setOptional(payload, "provider_project_id", s.ProviderProjectID)
	setOptional(payload, "default_branch", s.DefaultBranch)
	setOptional(payload, "visibility", s.Visibility)
	setNullableID(payload, "quality_gate", s.QualityGate)
	setNullableID(payload, "long_lived_quality_gate", s.LongLivedQualityGate)
	setOptional(payload, "enabled_engines", s.EnabledEngines)
	setOptional(payload, "publish_quality_gate", s.PublishQualityGate)
	setOptional(payload, "publish_review_feedback", s.PublishReviewFeedback)
	setOptional(payload, "excluded_paths", s.ExcludedPaths)
	setOptional(payload, "webhook_analysis_mode", s.WebhookAnalysisMode)
	setOptional(payload, "webhook_standard_scope", s.WebhookStandardScope)
	setOptional(payload, "webhook_standard_non_default_scope", s.WebhookStandardNonDefaultScope)
	setOptional(payload, "webhook_standard_long_lived_scope", s.WebhookStandardLongLivedScope)
	setOptional(payload, "webhook_ai_scope", s.WebhookAIScope)
	setOptional(payload, "webhook_ai_non_default_scope", s.WebhookAINonDefaultScope)
	setOptional(payload, "webhook_ai_long_lived_scope", s.WebhookAILongLivedScope)
	setOptional(payload, "webhook_long_lived_patterns", s.WebhookLongLivedPatterns)
	setOptional(payload, "webhook_default_gate_scope", s.WebhookDefaultGateScope)
	setOptional(payload, "webhook_non_default_gate_scope", s.WebhookNonDefaultGateScope)
	setOptional(payload, "webhook_long_lived_gate_scope", s.WebhookLongLivedGateScope)
	setOptional(payload, "sql_dialect", s.SQLDialect)
	setOptional(payload, "auto_assign", s.AutoAssign)
	setOptional(payload, "is_active", s.IsActive)
	repository := strings.TrimRight(strings.TrimSpace(s.RepositoryURL), "/")
	payload["repository_url"] = strings.TrimSuffix(repository, ".git")
	return payload
}

func SetupCodeProjectReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.CodeProject]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.CodeProject { return &v1alpha1.CodeProject{} },
		Adapter: ResourceAdapter[*v1alpha1.CodeProject]{
			APIEndpoint:    "/security/code-projects/",
			PatchDriftOnly: true,

			BuildPayload:  buildCodeProjectPayload,
			GetResourceID: func(obj *v1alpha1.CodeProject) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.CodeProject, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.CodeProject, data map[string]interface{}) {
				obj.Status.LatestAnalysis = getString(data, "latest_analysis")
				obj.Status.Key = getString(data, "key")
				obj.Status.Name = getString(data, "name")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.CodeProject) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.CodeProject, conditions []metav1.Condition) { obj.Status.Conditions = conditions },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.CodeProject{}).
		Named("codeproject").
		Complete(r)
}
