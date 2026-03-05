package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorSpec defines the desired state of a Monitor.
type MonitorSpec struct {
	// Name is the display name of the monitor.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// MonitorType is the type of monitor.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=http;ping;tcp;dns;ssl;keyword;grpc;heartbeat;multi_step_http
	MonitorType string `json:"monitorType"`

	// IsActive indicates whether the monitor is enabled.
	// +kubebuilder:default=true
	// +optional
	IsActive bool `json:"isActive,omitempty"`

	// URL is the target URL for http/keyword/ssl monitors.
	// +optional
	URL string `json:"url,omitempty"`

	// Hostname is the target hostname for ping/tcp/dns monitors.
	// +optional
	Hostname string `json:"hostname,omitempty"`

	// Port is the target port for tcp monitors.
	// +optional
	Port int64 `json:"port,omitempty"`

	// DNSRecordType is the DNS record type to query.
	// +kubebuilder:default="A"
	// +optional
	DNSRecordType string `json:"dnsRecordType,omitempty"`

	// DNSExpected is the expected DNS response.
	// +optional
	DNSExpected string `json:"dnsExpected,omitempty"`

	// HTTPMethod is the HTTP method to use.
	// +kubebuilder:default="GET"
	// +optional
	HTTPMethod string `json:"httpMethod,omitempty"`

	// HTTPHeaders is a JSON-encoded map of custom HTTP headers.
	// +kubebuilder:default="{}"
	// +optional
	HTTPHeaders string `json:"httpHeaders,omitempty"`

	// HTTPBody is the request body.
	// +optional
	HTTPBody string `json:"httpBody,omitempty"`

	// HTTPBodyContentType is the Content-Type for the request body.
	// +kubebuilder:default="application/json"
	// +optional
	HTTPBodyContentType string `json:"httpBodyContentType,omitempty"`

	// FollowRedirects indicates whether to follow HTTP redirects.
	// +kubebuilder:default=true
	// +optional
	FollowRedirects bool `json:"followRedirects,omitempty"`

	// VerifySSL indicates whether to verify SSL certificates.
	// +kubebuilder:default=true
	// +optional
	VerifySSL bool `json:"verifySSL,omitempty"`

	// HTTPAuthUser is the HTTP Basic Auth username.
	// +optional
	HTTPAuthUser string `json:"httpAuthUser,omitempty"`

	// HTTPAuthPass is the HTTP Basic Auth password.
	// +optional
	HTTPAuthPass string `json:"httpAuthPass,omitempty"`

	// Keyword is the keyword to search for in the response.
	// +optional
	Keyword string `json:"keyword,omitempty"`

	// KeywordExists - alert when keyword is missing (true) or found (false).
	// +kubebuilder:default=true
	// +optional
	KeywordExists bool `json:"keywordExists,omitempty"`

	// CheckInterval is the number of seconds between checks.
	// +kubebuilder:default=300
	// +optional
	CheckInterval int64 `json:"checkInterval,omitempty"`

	// Timeout is the request timeout in seconds.
	// +kubebuilder:default=30
	// +optional
	Timeout int64 `json:"timeout,omitempty"`

	// Retries is the number of retries before marking the monitor as down.
	// +kubebuilder:default=0
	// +optional
	Retries int64 `json:"retries,omitempty"`

	// RemediationEnabled indicates whether auto-remediation is enabled.
	// +kubebuilder:default=false
	// +optional
	RemediationEnabled bool `json:"remediationEnabled,omitempty"`

	// RemediationScript is the shell script to run for remediation.
	// +optional
	RemediationScript string `json:"remediationScript,omitempty"`

	// RemediationTimeout is the script timeout in seconds.
	// +kubebuilder:default=60
	// +optional
	RemediationTimeout int64 `json:"remediationTimeout,omitempty"`

	// RemediationWaitSeconds is the wait time after remediation in seconds.
	// +kubebuilder:default=30
	// +optional
	RemediationWaitSeconds int64 `json:"remediationWaitSeconds,omitempty"`

	// AutoInvestigate enables Argos AI auto-investigation when the monitor goes down.
	// +kubebuilder:default=false
	// +optional
	AutoInvestigate bool `json:"autoInvestigate,omitempty"`

	// AutoRemediate enables Argos AI auto-remediation after investigation.
	// +kubebuilder:default=false
	// +optional
	AutoRemediate bool `json:"autoRemediate,omitempty"`

	// RemediationStrategy is the remediation strategy: auto or approval_required.
	// +kubebuilder:default="approval_required"
	// +kubebuilder:validation:Enum=auto;approval_required
	// +optional
	RemediationStrategy string `json:"remediationStrategy,omitempty"`

	// HeartbeatGraceSeconds is the grace period in seconds for heartbeat monitors.
	// +kubebuilder:default=0
	// +optional
	HeartbeatGraceSeconds int64 `json:"heartbeatGraceSeconds,omitempty"`

	// MultiStepConfig is a JSON-encoded array of multi-step HTTP configurations.
	// +kubebuilder:default="[]"
	// +optional
	MultiStepConfig string `json:"multiStepConfig,omitempty"`

	// GRPCService is the gRPC service name.
	// +optional
	GRPCService string `json:"grpcService,omitempty"`

	// GRPCMethod is the gRPC method name.
	// +optional
	GRPCMethod string `json:"grpcMethod,omitempty"`

	// GRPCProto is the protobuf definition.
	// +optional
	GRPCProto string `json:"grpcProto,omitempty"`

	// GRPCMetadata is a JSON-encoded map of gRPC metadata.
	// +kubebuilder:default="{}"
	// +optional
	GRPCMetadata string `json:"grpcMetadata,omitempty"`

	// GRPCTLS indicates whether to use TLS for gRPC.
	// +kubebuilder:default=true
	// +optional
	GRPCTLS bool `json:"grpcTLS,omitempty"`

	// Assertions is a JSON-encoded array of response assertions.
	// +kubebuilder:default="[]"
	// +optional
	Assertions string `json:"assertions,omitempty"`

	// SSLExpiryWarnDays is the number of days before SSL expiry to warn.
	// +kubebuilder:default=30
	// +optional
	SSLExpiryWarnDays int64 `json:"sslExpiryWarnDays,omitempty"`

	// Location is the primary check location.
	// +kubebuilder:default="eu-france"
	// +optional
	Location string `json:"location,omitempty"`

	// Regions is a list of region codes for multi-region checks.
	// +optional
	Regions []string `json:"regions,omitempty"`

	// Tags is a list of tags.
	// +optional
	Tags []string `json:"tags,omitempty"`

	// GroupID is the UUID of the group this monitor belongs to.
	// +optional
	GroupID string `json:"groupId,omitempty"`
}

// MonitorStatus defines the observed state of a Monitor.
type MonitorStatus struct {
	// ID is the UUID of the monitor in the Argonix API.
	ID string `json:"id,omitempty"`

	// CurrentStatus is the current status (up, down, degraded, maintenance, unknown).
	CurrentStatus string `json:"currentStatus,omitempty"`

	// HeartbeatToken is the auto-generated heartbeat token.
	HeartbeatToken string `json:"heartbeatToken,omitempty"`

	// DateCreated is the creation timestamp.
	DateCreated string `json:"dateCreated,omitempty"`

	// DateModified is the last modification timestamp.
	DateModified string `json:"dateModified,omitempty"`

	// Conditions represent the latest available observations of the resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.monitorType`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.currentStatus`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.spec.isActive`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Monitor is the Schema for the monitors API.
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MonitorSpec   `json:"spec,omitempty"`
	Status            MonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MonitorList contains a list of Monitor.
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}
