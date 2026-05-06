# DClaw Operator

The DClaw Operator is a Kubernetes controller written in Go that watches `DClawApp` custom resources and automatically provisions all required infrastructure.

## Reconciliation Loop

When a `DClawApp` CR is created or updated, the operator performs these steps:

1. **Validate** — Check CRD spec for required fields
2. **Namespace** — Create a dedicated namespace (`dclaw-{appId}`)
3. **Database** — Provision a CloudNativePG cluster
4. **Backend** — Deploy FastAPI container with env vars
5. **Frontend** — Deploy Next.js container
6. **Ingress** — Configure nginx ingress + cert-manager TLS
7. **NetworkPolicy** — Isolate app namespace
8. **ResourceQuota** — Enforce CPU/memory limits
9. **Monitoring** — Create ServiceMonitor for Prometheus
10. **Status Update** — Write status back to CRD

## CRD: DClawApp

```yaml
apiVersion: platform.dclaw.io/v1
kind: DClawApp
metadata:
  name: chat
spec:
  appId: chat
  appName: DClaw Chat
  version: 0.1.0
  category: communication
  enabled: true
  frontend:
    image: ghcr.io/dclawstack/dclaw-chat:latest
    replicas: 2
  backend:
    image: ghcr.io/dclawstack/dclaw-chat-backend:latest
    replicas: 2
  database:
    enabled: true
    storage: 10Gi
  ingress:
    enabled: true
    host: chat.dclawstack.io
    tls: true
  resources:
    limits: { cpu: 1000m, memory: 2Gi }
    requests: { cpu: 250m, memory: 512Mi }
  branding:
    primaryColor: "#3B82F6"
  billing:
    tier: pro
```

## Operator Code Structure

```
dclaw-operator/
├── main.go                  # Entry point
├── controllers/
│   └── dclawapp_controller.go
├── reconcilers/
│   ├── namespace.go
│   ├── database.go
│   ├── backend.go
│   ├── frontend.go
│   ├── ingress.go
│   ├── network.go
│   └── monitoring.go
└── api/
    └── v1/
        ├── dclawapp_types.go
        └── groupversion_info.go
```

## Installation

```bash
helm install dclaw-operator dclaw/dclaw-operator \
  --namespace dclaw-system \
  --create-namespace
```
