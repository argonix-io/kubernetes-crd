package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestPeriodicJobReconcilesLifecycleAndRunStatus(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	prompt := "Summarize availability"
	channels := []string{"chat-b", "chat-a"}
	job := &v1alpha1.PeriodicJob{
		ObjectMeta: metav1.ObjectMeta{Name: "availability", Namespace: "default"},
		Spec: v1alpha1.PeriodicJobSpec{
			Name: "Availability", Prompt: &prompt, NotificationChannelIDs: &channels,
		},
	}
	kube := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(job).WithObjects(job).Build()

	var methods []string
	remotePrompt := prompt
	runStatus := "running"
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		methods = append(methods, req.Method)
		path := "/api/0.1/organizations/org/argos/jobs/"
		if req.Method != http.MethodPost {
			path += "job/"
		}
		if req.URL.Path != path {
			t.Fatalf("request path = %q, want %q", req.URL.Path, path)
		}
		if req.Method == http.MethodPost || req.Method == http.MethodPatch {
			var payload map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if req.Method == http.MethodPost && len(payload) != 3 {
				t.Fatalf("create payload = %#v", payload)
			}
			if req.Method == http.MethodPatch && len(payload) != 1 {
				t.Fatalf("patch should contain only the changed prompt: %#v", payload)
			}
			remotePrompt = payload["prompt"].(string)
		}
		body := fmt.Sprintf(`{"id":"job","name":"Availability","prompt":%q,"notification_channels":["chat-a","chat-b"],"latest_run_status":%q}`, remotePrompt, runStatus)
		statusCode := http.StatusOK
		if req.Method == http.MethodDelete {
			statusCode = http.StatusNoContent
			body = ""
		}
		return &http.Response{StatusCode: statusCode, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	reconciler := ResourceReconciler[*v1alpha1.PeriodicJob]{
		Client: kube, Scheme: scheme,
		ArgonixClient: &argonixclient.Client{BaseURL: "https://example.test", OrganizationID: "org", HTTPClient: httpClient},
		Adapter:       periodicJobAdapter(), NewObject: func() *v1alpha1.PeriodicJob { return &v1alpha1.PeriodicJob{} },
	}
	key := client.ObjectKeyFromObject(job)
	reconcile := func() {
		t.Helper()
		if _, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: key}); err != nil {
			t.Fatal(err)
		}
	}
	reconcile() // Add finalizer.
	reconcile() // Create.
	runStatus = "completed"
	reconcile() // Refresh without rewriting differently ordered relations.
	if err := kube.Get(ctx, key, job); err != nil {
		t.Fatal(err)
	}
	if job.Status.ID != "job" || job.Status.LatestRunStatus != "completed" {
		t.Fatalf("status after GET = %#v", job.Status)
	}
	prompt = "Summarize availability and incidents"
	job.Spec.Prompt = &prompt
	if err := kube.Update(ctx, job); err != nil {
		t.Fatal(err)
	}
	reconcile()
	if err := kube.Delete(ctx, job); err != nil {
		t.Fatal(err)
	}
	reconcile()
	if got, want := strings.Join(methods, " "), "POST GET GET PATCH DELETE"; got != want {
		t.Fatalf("requests = %q, want %q", got, want)
	}
}

func TestSanityCheckSendsScenarioStepsAsJSONArray(t *testing.T) {
	manual := "manual"
	readOnly := false
	check := &v1alpha1.SyntheticTest{Spec: v1alpha1.SyntheticTestSpec{
		Name: "Health", TestType: "scenario", Steps: `[{"id":"health","kind":"http"}]`,
		ScheduleMode: &manual, AllowMutatingActions: &readOnly,
	}}
	payload := buildSyntheticTestPayload(check)
	if _, ok := payload["steps"].([]interface{}); !ok {
		t.Fatalf("steps were not decoded: %#v", payload["steps"])
	}
	if payload["schedule_mode"] != "manual" || payload["allow_mutating_actions"] != false {
		t.Fatalf("schedule/action policy was lost: %#v", payload)
	}
	if err := validateJSONField("steps", `{"invalid":true}`, '['); err == nil {
		t.Fatal("object accepted as a steps array")
	}
}
