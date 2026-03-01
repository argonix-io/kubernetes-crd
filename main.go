package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
	argonixclient "github.com/argonix-io/kubernetes-crd/internal/client"
	"github.com/argonix-io/kubernetes-crd/internal/controller"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var probeAddr string
	var enableLeaderElection bool
	var argonixURL string
	var argonixAPIKey string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager.")
	flag.StringVar(&argonixURL, "argonix-url", "", "Argonix API base URL (default: https://api.argonix.io)")
	flag.StringVar(&argonixAPIKey, "argonix-api-key", "", "Argonix API key")

	opts := zap.Options{Development: true}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Fall back to environment variables.
	if argonixURL == "" {
		argonixURL = os.Getenv("ARGONIX_URL")
	}
	if argonixURL == "" {
		argonixURL = "https://api.argonix.io"
	}
	if argonixAPIKey == "" {
		argonixAPIKey = os.Getenv("ARGONIX_API_KEY")
	}
	if argonixAPIKey == "" {
		setupLog.Error(fmt.Errorf("API key is required"), "Please provide --argonix-api-key flag or ARGONIX_API_KEY env var")
		os.Exit(1)
	}

	// Initialize Argonix API client.
	ac, err := argonixclient.NewClient(context.Background(), argonixURL, argonixAPIKey)
	if err != nil {
		setupLog.Error(err, "Failed to initialize Argonix API client")
		os.Exit(1)
	}
	setupLog.Info("Argonix API client initialized", "url", argonixURL, "organization", ac.OrganizationID)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "argonix-crd-controller.argonix.io",
	})
	if err != nil {
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}

	// Setup all controllers.
	setupFuncs := []struct {
		name  string
		setup func(ctrl.Manager, *argonixclient.Client) error
	}{
		{"Monitor", controller.SetupMonitorReconciler},
		{"SyntheticTest", controller.SetupSyntheticTestReconciler},
		{"Group", controller.SetupGroupReconciler},
		{"AlertChannel", controller.SetupAlertChannelReconciler},
		{"NotificationRule", controller.SetupNotificationRuleReconciler},
		{"StatusPage", controller.SetupStatusPageReconciler},
		{"TestSuite", controller.SetupTestSuiteReconciler},
		{"ManualTestCase", controller.SetupManualTestCaseReconciler},
		{"TestPlan", controller.SetupTestPlanReconciler},
	}

	for _, sf := range setupFuncs {
		if err := sf.setup(mgr, ac); err != nil {
			setupLog.Error(err, "Unable to create controller", "controller", sf.name)
			os.Exit(1)
		}
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "Unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "Unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "Problem running manager")
		os.Exit(1)
	}
}
