package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

const (
	finalizerName   = "argonix.io/finalizer"
	requeueInterval = 5 * time.Minute
)

// ResourceAdapter provides type-specific operations for a CRD resource.
type ResourceAdapter[T client.Object] struct {
	// APIEndpoint is the API path (e.g., "/monitors/").
	APIEndpoint string

	// BuildPayload converts the K8s spec into an API request payload.
	BuildPayload func(obj T) map[string]interface{}

	// GetResourceID returns the remote resource ID from the status.
	GetResourceID func(obj T) string

	// SetResourceID sets the remote resource ID in the status.
	SetResourceID func(obj T, id string)

	// SetStatusFromResponse updates the status fields from an API response.
	SetStatusFromResponse func(obj T, data map[string]interface{})

	// GetConditions returns the status conditions slice.
	GetConditions func(obj T) []metav1.Condition

	// SetConditions sets the status conditions slice.
	SetConditions func(obj T, conditions []metav1.Condition)
}

// ResourceReconciler is a generic reconciler for Argonix CRD resources.
type ResourceReconciler[T client.Object] struct {
	client.Client
	Scheme        *runtime.Scheme
	ArgonixClient *argonixclient.Client
	Adapter       ResourceAdapter[T]
	NewObject     func() T
}

func (r *ResourceReconciler[T]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	obj := r.NewObject()
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle deletion.
	if !obj.GetDeletionTimestamp().IsZero() {
		if controllerutil.ContainsFinalizer(obj, finalizerName) {
			resourceID := r.Adapter.GetResourceID(obj)
			if resourceID != "" {
				if err := r.ArgonixClient.Delete(ctx, r.Adapter.APIEndpoint+resourceID+"/"); err != nil {
					if !argonixclient.IsNotFound(err) {
						logger.Error(err, "Failed to delete remote resource")
						return ctrl.Result{}, err
					}
				}
				logger.Info("Deleted remote resource", "id", resourceID)
			}
			controllerutil.RemoveFinalizer(obj, finalizerName)
			if err := r.Update(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Ensure finalizer.
	if !controllerutil.ContainsFinalizer(obj, finalizerName) {
		controllerutil.AddFinalizer(obj, finalizerName)
		if err := r.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
	}

	resourceID := r.Adapter.GetResourceID(obj)
	payload := r.Adapter.BuildPayload(obj)

	if resourceID == "" {
		// Create.
		var resp map[string]interface{}
		if err := r.ArgonixClient.Create(ctx, r.Adapter.APIEndpoint, payload, &resp); err != nil {
			r.setCondition(obj, "Ready", metav1.ConditionFalse, "CreateFailed", err.Error())
			_ = r.Status().Update(ctx, obj)
			return ctrl.Result{}, err
		}
		if id, ok := resp["id"].(string); ok {
			r.Adapter.SetResourceID(obj, id)
		}
		r.Adapter.SetStatusFromResponse(obj, resp)
		r.setCondition(obj, "Ready", metav1.ConditionTrue, "Created", "Resource created successfully")
		r.setCondition(obj, "Synced", metav1.ConditionTrue, "Synced", "Resource in sync")
		if err := r.Status().Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("Created remote resource", "id", r.Adapter.GetResourceID(obj))
		return ctrl.Result{RequeueAfter: requeueInterval}, nil
	}

	// Read existing to check drift.
	var remote map[string]interface{}
	if err := r.ArgonixClient.Read(ctx, r.Adapter.APIEndpoint+resourceID+"/", &remote); err != nil {
		if argonixclient.IsNotFound(err) {
			// Recreate.
			logger.Info("Remote resource not found, recreating", "id", resourceID)
			r.Adapter.SetResourceID(obj, "")
			var resp map[string]interface{}
			if createErr := r.ArgonixClient.Create(ctx, r.Adapter.APIEndpoint, payload, &resp); createErr != nil {
				r.setCondition(obj, "Ready", metav1.ConditionFalse, "RecreateFailed", createErr.Error())
				_ = r.Status().Update(ctx, obj)
				return ctrl.Result{}, createErr
			}
			if id, ok := resp["id"].(string); ok {
				r.Adapter.SetResourceID(obj, id)
			}
			r.Adapter.SetStatusFromResponse(obj, resp)
			r.setCondition(obj, "Ready", metav1.ConditionTrue, "Recreated", "Resource recreated successfully")
			r.setCondition(obj, "Synced", metav1.ConditionTrue, "Synced", "Resource in sync")
			if err := r.Status().Update(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: requeueInterval}, nil
		}
		return ctrl.Result{}, err
	}

	// Update if drift detected.
	if !r.payloadMatchesRemote(payload, remote) {
		var resp map[string]interface{}
		if err := r.ArgonixClient.Update(ctx, r.Adapter.APIEndpoint+resourceID+"/", payload, &resp); err != nil {
			r.setCondition(obj, "Synced", metav1.ConditionFalse, "UpdateFailed", err.Error())
			_ = r.Status().Update(ctx, obj)
			return ctrl.Result{}, err
		}
		r.Adapter.SetStatusFromResponse(obj, resp)
		r.setCondition(obj, "Ready", metav1.ConditionTrue, "Updated", "Resource updated successfully")
		r.setCondition(obj, "Synced", metav1.ConditionTrue, "Synced", "Resource in sync")
		if err := r.Status().Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("Updated remote resource", "id", resourceID)
	}

	return ctrl.Result{RequeueAfter: requeueInterval}, nil
}

func (r *ResourceReconciler[T]) payloadMatchesRemote(payload map[string]interface{}, remote map[string]interface{}) bool {
	for key, desired := range payload {
		actual, exists := remote[key]
		if !exists {
			return false
		}
		if !equality.Semantic.DeepEqual(desired, actual) {
			// Try JSON comparison for complex types.
			desiredJSON, err1 := json.Marshal(desired)
			actualJSON, err2 := json.Marshal(actual)
			if err1 != nil || err2 != nil || string(desiredJSON) != string(actualJSON) {
				return false
			}
		}
	}
	return true
}

func (r *ResourceReconciler[T]) setCondition(obj T, condType string, status metav1.ConditionStatus, reason, message string) {
	conditions := r.Adapter.GetConditions(obj)
	newCondition := metav1.Condition{
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	for i, c := range conditions {
		if c.Type == condType {
			if c.Status != status {
				conditions[i] = newCondition
			} else {
				conditions[i].Reason = reason
				conditions[i].Message = message
			}
			r.Adapter.SetConditions(obj, conditions)
			return
		}
	}
	conditions = append(conditions, newCondition)
	r.Adapter.SetConditions(obj, conditions)
}

// Helper functions for payload building.

func sliceToJSON(s []string) string {
	if s == nil {
		return "[]"
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
