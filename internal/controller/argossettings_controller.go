package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func SetupArgosSettingsReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &argosSettingsReconciler{
		Client:        mgr.GetClient(),
		ArgonixClient: ac,
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ArgosSettings{}).
		Named("argossettings").
		Complete(r)
}

type argosSettingsReconciler struct {
	Client        k8sclient.Client
	ArgonixClient *argonixclient.Client
}

func (r *argosSettingsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the ArgosSettings CR
	var obj v1alpha1.ArgosSettings
	if err := r.Client.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, k8sclient.IgnoreNotFound(err)
	}

	// Build the payload
	payload := map[string]interface{}{
		"llm_provider":        obj.Spec.LLMProvider,
		"llm_model":           obj.Spec.LLMModel,
		"llm_base_url":        obj.Spec.LLMBaseURL,
		"custom_instructions": obj.Spec.CustomInstructions,
		"demo_mode":           obj.Spec.DemoMode,
	}

	// Resolve API key from Secret
	if obj.Spec.LLMApiKeySecretRef != nil {
		key := obj.Spec.LLMApiKeySecretRef.Key
		if key == "" {
			key = "api-key"
		}
		var secret corev1.Secret
		if err := r.Client.Get(ctx, k8stypes.NamespacedName{
			Namespace: obj.Namespace,
			Name:      obj.Spec.LLMApiKeySecretRef.Name,
		}, &secret); err != nil {
			log.Error(err, "Failed to read LLM API key secret")
			return ctrl.Result{}, err
		}
		apiKey, ok := secret.Data[key]
		if !ok {
			err := fmt.Errorf("key %q not found in Secret %s", key, obj.Spec.LLMApiKeySecretRef.Name)
			log.Error(err, "Missing key in Secret")
			return ctrl.Result{}, err
		}
		payload["llm_api_key"] = string(apiKey)
	}

	// Read existing settings (auto-created by API)
	var existing map[string]interface{}
	if err := r.ArgonixClient.Read(ctx, "/argos/settings/", &existing); err != nil {
		log.Error(err, "Failed to read argos settings from API")
		return ctrl.Result{}, err
	}
	settingsID, _ := existing["id"].(string)

	// Update via PATCH
	var result map[string]interface{}
	if err := r.ArgonixClient.Update(ctx, fmt.Sprintf("/argos/settings/%s/", settingsID), payload, &result); err != nil {
		log.Error(err, "Failed to update argos settings")
		return ctrl.Result{}, err
	}

	// Update status
	obj.Status.ID = settingsID
	obj.Status.DateModified = getString(result, "date_modified")
	obj.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "Synced",
			Message:            "Settings synced to Argonix API",
		},
	}
	if err := r.Client.Status().Update(ctx, &obj); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
