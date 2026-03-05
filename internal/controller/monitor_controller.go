package controller

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
)

func tagsMapToJSON(tags map[string]string) string {
	if tags == nil {
		return "{}"
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func SetupMonitorReconciler(mgr ctrl.Manager, ac *argonixclient.Client) error {
	r := &ResourceReconciler[*v1alpha1.Monitor]{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ArgonixClient: ac,
		NewObject:     func() *v1alpha1.Monitor { return &v1alpha1.Monitor{} },
		Adapter: ResourceAdapter[*v1alpha1.Monitor]{
			APIEndpoint: "/monitors/",
			BuildPayload: func(obj *v1alpha1.Monitor) map[string]interface{} {
				s := obj.Spec
				p := map[string]interface{}{
					"name":                     s.Name,
					"monitor_type":             s.MonitorType,
					"is_active":                s.IsActive,
					"url":                      s.URL,
					"hostname":                 s.Hostname,
					"port":                     s.Port,
					"dns_record_type":          s.DNSRecordType,
					"dns_expected":             s.DNSExpected,
					"http_method":              s.HTTPMethod,
					"http_headers":             s.HTTPHeaders,
					"http_body":                s.HTTPBody,
					"http_body_content_type":   s.HTTPBodyContentType,
					"follow_redirects":         s.FollowRedirects,
					"verify_ssl":               s.VerifySSL,
					"http_auth_user":           s.HTTPAuthUser,
					"http_auth_pass":           s.HTTPAuthPass,
					"keyword":                  s.Keyword,
					"keyword_exists":           s.KeywordExists,
					"check_interval":           s.CheckInterval,
					"timeout":                  s.Timeout,
					"retries":                  s.Retries,
					"remediation_enabled":      s.RemediationEnabled,
					"remediation_script":       s.RemediationScript,
					"remediation_timeout":      s.RemediationTimeout,
					"remediation_wait_seconds": s.RemediationWaitSeconds,
					"auto_investigate":         s.AutoInvestigate,
					"auto_remediate":           s.AutoRemediate,
					"remediation_strategy":     s.RemediationStrategy,
					"heartbeat_grace_seconds":  s.HeartbeatGraceSeconds,
					"multi_step_config":        s.MultiStepConfig,
					"grpc_service":             s.GRPCService,
					"grpc_method":              s.GRPCMethod,
					"grpc_proto":               s.GRPCProto,
					"grpc_metadata":            s.GRPCMetadata,
					"grpc_tls":                 s.GRPCTLS,
					"assertions":               s.Assertions,
					"ssl_expiry_warn_days":     s.SSLExpiryWarnDays,
					"location":                 s.Location,
					"regions":                  s.Regions,
					"tags":                     s.Tags,
				}
				if s.GroupID != "" {
					p["group"] = s.GroupID
				}
				return p
			},
			GetResourceID: func(obj *v1alpha1.Monitor) string { return obj.Status.ID },
			SetResourceID: func(obj *v1alpha1.Monitor, id string) { obj.Status.ID = id },
			SetStatusFromResponse: func(obj *v1alpha1.Monitor, data map[string]interface{}) {
				obj.Status.CurrentStatus = getString(data, "current_status")
				obj.Status.HeartbeatToken = getString(data, "heartbeat_token")
				obj.Status.DateCreated = getString(data, "date_created")
				obj.Status.DateModified = getString(data, "date_modified")
			},
			GetConditions: func(obj *v1alpha1.Monitor) []metav1.Condition { return obj.Status.Conditions },
			SetConditions: func(obj *v1alpha1.Monitor, c []metav1.Condition) { obj.Status.Conditions = c },
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Monitor{}).
		Named("monitor").
		Complete(r)
}
