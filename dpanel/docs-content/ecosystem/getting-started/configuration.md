# Configuration

## Platform-Level Configuration

### DClaw Operator

The operator is configured via a ConfigMap in the `dclaw-system` namespace:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dclaw-operator-config
  namespace: dclaw-system
data:
  defaultReplicas: "2"
  defaultStorage: "10Gi"
  ingressClass: "nginx"
  certManagerIssuer: "letsencrypt-prod"
  enableNetworkPolicies: "true"
  enableResourceQuotas: "true"
```

### DPanel

DPanel reads configuration from environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_DPANEL_API_URL` | dpanel-api endpoint | `""` |
| `DPANEL_AUTH_DOMAIN` | Logto auth domain | `""` |
| `DPANEL_BILLING_ENABLED` | Enable Stripe billing | `"false"` |
| `DPANEL_THEME` | UI theme | `"dark"` |

### Per-App Configuration

Each app is configured via its `DClawApp` CRD:

```yaml
apiVersion: platform.dclaw.io/v1
kind: DClawApp
metadata:
  name: chat
spec:
  appId: chat
  appName: DClaw Chat
  version: 0.2.0
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
```

## Environment Variables

Apps receive configuration through environment variables:

```yaml
# Backend
DATABASE_URL: postgresql+asyncpg://user:pass@db:5432/dclaw_chat
REDIS_URL: redis://redis:6379
OLLAMA_URL: http://ollama:11434
OPENROUTER_API_KEY: ""

# Frontend
NEXT_PUBLIC_API_URL: /api
NEXT_PUBLIC_APP_NAME: DClaw Chat
```

## Secrets Management

Sensitive values are stored in Kubernetes Secrets:

```bash
kubectl create secret generic chat-secrets \
  --from-literal=OPENROUTER_API_KEY=sk-... \
  --namespace dclaw-chat
```

## Local Overrides

For local development, create a `.env.local` file in each app's frontend directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:8000
NEXT_PUBLIC_APP_NAME="DClaw Chat (Local)"
```
