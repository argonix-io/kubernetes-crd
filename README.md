# Argonix Kubernetes CRD Controller

Kubernetes operator for managing [Argonix](https://argonix.io) resources as native Kubernetes Custom Resources.

## Resources

This controller manages the following Argonix resources:

| Kind | API Endpoint | Description |
|------|-------------|-------------|
| `Monitor` | `/monitors/` | HTTP, ping, TCP, DNS, SSL, keyword, gRPC, heartbeat, multi-step monitors |
| `SyntheticTest` | `/synthetic-tests/` | API and browser synthetic tests |
| `Group` | `/groups/` | Monitor groups |
| `AlertChannel` | `/alert-channels/` | Email, Slack, webhook, PagerDuty, OpsGenie, Telegram, Discord, Teams, Jira |
| `NotificationRule` | `/notification-rules/` | Alert routing rules |
| `StatusPage` | `/status-pages/` | Public/private status pages |
| `TestSuite` | `/test-suites/` | Test suite groupings |
| `ManualTestCase` | `/manual-test-cases/` | Manual test case definitions |
| `TestPlan` | `/test-plans/` | Test plan orchestration |

All resources belong to the `argonix.io/v1alpha1` API group.

## Prerequisites

- Go 1.22+
- Kubernetes cluster 1.28+
- `kubectl` configured to access your cluster
- An Argonix API key

## Quick Start

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

The controller uses a **generic reconciler** pattern — a single `ResourceReconciler[T]` handles CRUD for all 9 resource types. Each resource provides a `ResourceAdapter[T]` that maps between the CRD spec and the Argonix API.

### Reconciliation flow

1. **Create**: When a new CR is applied, the controller calls the Argonix API to create the resource and stores the API ID in `.status.id`.
2. **Update**: When the CR spec changes, the controller sends a PUT to the API with the full payload.
3. **Delete**: A finalizer ensures the API resource is deleted before the CR is removed.
4. **Drift detection**: The controller re-reconciles every 5 minutes to detect out-of-band changes.

### Status conditions

Each resource exposes two conditions:
- `Ready` — Whether the resource is provisioned in the Argonix API
- `Synced` — Whether the last sync operation succeeded

## Sample manifests

See [`config/samples/`](config/samples/) for example resource definitions.
