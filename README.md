# Argonix Kubernetes CRD Controller

Kubernetes operator for managing [Argonix](https://argonix.io) resources as native Kubernetes Custom Resources.

## Resources

This controller manages the following Argonix resources:

| Kind | API Endpoint | Description |
|------|-------------|-------------|
| `Monitor` | `/monitors/` | HTTP, ping, TCP, DNS, SSL, keyword, gRPC, heartbeat, multi-step monitors |
| `SyntheticTest` | `/synthetic-tests/` | API, browser, and Sanity scenarios |
| `Group` | `/groups/` | Organization groups |
| `MonitorGroup` | `/monitor-groups/` | Monitor folders |
| `AlertChannel` | `/alert-channels/` | Email, Slack, webhook, PagerDuty, OpsGenie, Telegram, Discord, Teams, Jira |
| `NotificationRule` | `/notification-rules/` | Alert routing rules |
| `StatusPage` | `/status-pages/` | Public/private status pages |
| `TestSuite` | `/test-suites/` | Test suite groupings |
| `ManualTestCase` | `/manual-test-cases/` | Manual test case definitions |
| `TestPlan` | `/test-plans/` | Test plan orchestration |
| `Connector` | `/argos/connectors/` | Integrations (cloud, K8s, Git, observability, ticketing…) |
| `KnowledgeBase` | `/argos/knowledge-bases/` | Runbooks and docs the AI agent can consult |
| `Workflow` | `/argos/workflows/` | Workflow definitions, schedules, and incident policies |
| `PeriodicJob` | `/argos/jobs/` | Scheduled Argos jobs |
| `SREPatrol` | `/argos/patrols/` | Scheduled sweeps for operational signals |
| `Persona` | `/argos/personas/` | AI agent personas (scoped tools, prompts, approvals) |
| `ChatChannel` | `/argos/chat-channels/` | Chat integrations for Argos (Slack, Teams…) |
| `MaintenanceWindow` | `/maintenance-windows/` | Scheduled maintenance windows |
| `Environment` | `/environments/` | Deployment environments |
| `ArgosSettings` | `/argos/settings/` | Org-level Argos configuration |
| `AlertSource` | `/argos/alert-sources/` | Inbound alert sources for incidents |
| `SecurityPolicy` | `/security/policies/` | Deployment gate / compliance policies |
| `SecurityScan` | `/security/scans/` | Security scan requests |
| `ScanTarget` | `/security/scan-targets/` | Static analysis targets for Git repositories and OCI images |
| `CodeProject` | `/security/code-projects/` | SAST and code quality configuration for Git repositories |
| `DetectionRule` | `/security/detection-rules/` | Threat detection rules |
| `HealthNotebook` | `/argos/health-notebooks/` | Scheduled infrastructure health reports and security briefs |
| `Tenant` | `/finops/tenants/` | FinOps cost owners (team, product, client, environment) |
| `CostBudget` | `/finops/budgets/` | FinOps spend budgets with threshold alerts |

All resources belong to the `argonix.io/v1alpha1` API group.

## Prerequisites

- Go 1.24.3+
- Kubernetes cluster 1.28+
- `kubectl` configured to access your cluster
- An Argonix API key

## Helm installation and upgrade

Current release: chart `0.1.22`, operator image `ghcr.io/argonix-io/kubernetes-crd:0.1.22`.
Use an Argonix API key with `iac:read` and `iac:write` and the organization role required by the resources you manage. Deleting a `PeriodicJob` or `CodeProject` requires an organization admin.

```bash
kubectl create namespace argonix-system
kubectl create secret generic argonix-credentials \
  --namespace argonix-system \
  --from-literal=api-key=YOUR_API_KEY
helm repo add argonix https://charts.argonix.io
helm repo update
helm upgrade --install argonix-crd argonix/argonix-crd \
  --namespace argonix-system --version 0.1.22 \
  --set argonix.existingSecret=argonix-credentials
```

Helm applies the CRD schemas on install and upgrade. If your previous Helm values
pin an image tag, update the override to `0.1.22` or remove it to use the chart's
`appVersion`. For an installation managed directly with `kubectl`, apply the CRD
and RBAC manifests from the same release before upgrading the controller.

## Build from source

### 1. Build

```bash
make build
```

### 2. Install CRDs

```bash
# Generate CRD manifests (requires controller-gen)
make manifests

# Apply CRDs to cluster
make install
```

### 3. Configure credentials

```bash
kubectl create namespace argonix-system

kubectl create secret generic argonix-credentials \
  --namespace argonix-system \
  --from-literal=api-key=YOUR_API_KEY \
  --from-literal=url=https://api.argonix.io
```

### 4. Deploy the controller

```bash
# Build and push Docker image
make docker-build docker-push IMG=your-registry/argonix-controller:latest

# Deploy to cluster
make deploy
```

### 5. Create resources

```yaml
apiVersion: argonix.io/v1alpha1
kind: Monitor
metadata:
  name: my-website
spec:
  name: "My Website"
  monitorType: http
  url: "https://example.com"
  checkInterval: 60
  tags:
    - production
```

```bash
kubectl apply -f monitor.yaml
kubectl get monitors
kubectl describe monitor my-website
```

## Configuration

The controller accepts configuration via flags or environment variables:

| Flag | Env Var | Required | Default | Description |
|------|---------|----------|---------|-------------|
| `--argonix-api-key` | `ARGONIX_API_KEY` | **Yes** | — | Argonix API key |
| `--argonix-url` | `ARGONIX_URL` | No | `https://api.argonix.io` | API base URL |
| `--leader-elect` | — | No | `false` | Enable leader election |
| `--metrics-bind-address` | — | No | `:8080` | Metrics endpoint |
| `--health-probe-bind-address` | — | No | `:8081` | Health probes endpoint |

The organization is automatically discovered from the API key.

### ArgosSettings Secret reference

When `ArgosSettings.spec.llmApiKeySecretRef` is set, create the Secret in the
same namespace as the `ArgosSettings` object. Grant the controller's service
account `get` on that named Secret. For the Helm installation above, this
example grants access to `argos-llm-secret` in the `platform` namespace:

```bash
kubectl create role argos-llm-secret-reader \
  --namespace platform --verb=get --resource=secrets \
  --resource-name=argos-llm-secret
kubectl create rolebinding argos-llm-secret-reader \
  --namespace platform --role=argos-llm-secret-reader \
  --serviceaccount=argonix-system:argonix-crd
```

Use the actual service account name if you override it. The controller reads
the named Secret directly and picks up rotations on its next reconciliation.

## Development

```bash
# Run locally (against current kubeconfig cluster)
make run

# Run tests
make test

# Regenerate deepcopy after changing types
make generate

# Regenerate CRD manifests after changing types
make manifests
```

## Architecture

The controller uses a **generic reconciler** pattern — a single `ResourceReconciler[T]` handles CRUD for all resource types. Each resource provides a `ResourceAdapter[T]` that maps between the CRD spec and the Argonix API.

### Reconciliation flow

1. **Create**: When a new CR is applied, the controller calls the Argonix API to create the resource and stores the API ID in `.status.id`.
2. **Update**: When the CR spec changes, the controller sends a PATCH to the API.
3. **Delete**: A finalizer ensures the API resource is deleted before the CR is removed.
4. **Drift detection**: The controller reads the API resource every 5 minutes, refreshes observed status, and corrects drift in managed fields.

### Status conditions

Each resource exposes two conditions:
- `Ready` — Whether the resource is provisioned in the Argonix API
- `Synced` — Whether the last sync operation succeeded

## Sample manifests

See [`config/samples/`](config/samples/) for example resource definitions.
Sanity checks use `SyntheticTest.spec.testType: scenario`. The `steps` field
keeps its existing JSON string format. Read-only run results and manual run
actions are available through Argonix. Code quality gate UUIDs can be referenced
from a `CodeProject`. Versioned gate definitions remain managed through Argonix.
