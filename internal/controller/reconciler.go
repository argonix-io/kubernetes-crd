package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
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

	// ValidateSpec rejects malformed fields before sending a request.
	ValidateSpec func(obj T) error

	// NormalizeRemote reconciles API read and write representations before comparison.
	NormalizeRemote func(desired, remote map[string]interface{})

	// PatchDriftOnly avoids rewriting fields with API side effects.
	PatchDriftOnly bool

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
			before := obj.DeepCopyObject().(T)
			controllerutil.RemoveFinalizer(obj, finalizerName)
			if err := r.Patch(ctx, obj, client.MergeFrom(before)); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Ensure finalizer.
	if !controllerutil.ContainsFinalizer(obj, finalizerName) {
		before := obj.DeepCopyObject().(T)
		controllerutil.AddFinalizer(obj, finalizerName)
		if err := r.Patch(ctx, obj, client.MergeFrom(before)); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if r.Adapter.ValidateSpec != nil {
		if err := r.Adapter.ValidateSpec(obj); err != nil {
			r.setCondition(obj, "Synced", metav1.ConditionFalse, "InvalidSpec", err.Error())
			_ = r.Status().Update(ctx, obj)
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

	if r.Adapter.NormalizeRemote != nil {
		r.Adapter.NormalizeRemote(payload, remote)
	}
	// Update if drift detected.
	if changes := driftedFields(payload, remote); len(changes) != 0 {
		if r.Adapter.PatchDriftOnly {
			payload = changes
		}
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
	} else {
		before := obj.DeepCopyObject()
		r.Adapter.SetStatusFromResponse(obj, remote)
		for _, condition := range r.Adapter.GetConditions(obj) {
			if (condition.Type == "Synced" || condition.Type == "Ready") && condition.Status != metav1.ConditionTrue {
				r.setCondition(obj, condition.Type, metav1.ConditionTrue, "Synced", "Resource in sync")
			}
		}
		if !equality.Semantic.DeepEqual(before, obj) {
			if err := r.Status().Update(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: requeueInterval}, nil
}

func driftedFields(payload, remote map[string]interface{}) map[string]interface{} {
	changes := make(map[string]interface{})
	for key, desired := range payload {
		actual, exists := remote[key]
		desiredJSON, desiredErr := json.Marshal(desired)
		actualJSON, actualErr := json.Marshal(actual)
		if !exists || desiredErr != nil || actualErr != nil || !bytes.Equal(desiredJSON, actualJSON) {
			changes[key] = desired
		}
	}
	return changes
}

// normalizeIDList compares unordered API relations that may use nested objects on reads.
func normalizeIDList(desired, remote map[string]interface{}, writeKey, readKey string) {
	if _, managed := desired[writeKey]; !managed {
		return
	}
	rows, ok := remote[readKey].([]interface{})
	if !ok {
		return
	}
	actual := make([]string, 0, len(rows))
	for _, row := range rows {
		if object, nested := row.(map[string]interface{}); nested {
			row = object["id"]
		}
		id, ok := row.(string)
		if !ok {
			return
		}
		actual = append(actual, id)
	}
	var expected []string
	encoded, err := json.Marshal(desired[writeKey])
	if err != nil || json.Unmarshal(encoded, &expected) != nil {
		return
	}
	sort.Strings(actual)
	sort.Strings(expected)
	if len(actual) == len(expected) {
		for i := range actual {
			if actual[i] != expected[i] {
				return
			}
		}
		remote[writeKey] = desired[writeKey]
	}
}

func normalizeArgosRelations(desired, remote map[string]interface{}) {
	normalizeIDList(desired, remote, "connector_ids", "connectors")
	normalizeIDList(desired, remote, "notification_channel_ids", "notification_channels")
	normalizeIDList(desired, remote, "alert_channel_ids", "alert_channels")
	normalizeIDList(desired, remote, "application_root_ids", "application_root_ids")
}

// setOptional leaves API fields unmanaged when the manifest omits them.
func setOptional[V any](payload map[string]interface{}, key string, value *V) {
	if value != nil {
		payload[key] = *value
	}
}

// setNullableID translates an explicit empty ID into an API null.
func setNullableID(payload map[string]interface{}, key string, value *string) {
	if value == nil {
		return
	}
	if *value == "" {
		payload[key] = nil
	} else {
		payload[key] = *value
	}
}

func validateJSONField(field, value string, delimiter byte) error {
	if value == "" {
		return nil
	}
	data := bytes.TrimSpace([]byte(value))
	if !json.Valid(data) || data[0] != delimiter {
		return fmt.Errorf("%s must be a JSON %s", field, map[byte]string{'[': "array", '{': "object"}[delimiter])
	}
	return nil
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
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func parseJSONStringField(jsonStr, key string, payload map[string]interface{}) {
	if jsonStr != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(jsonStr), &v); err == nil {
			payload[key] = v
		}
	}
}

func getInt(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
}
