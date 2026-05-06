# Installation

## Prerequisites

Before installing DClaw, ensure you have:

- **Kubernetes cluster** (k3s, minikube, or managed EKS/GKE/AKS)
- **kubectl** configured to access your cluster
- **Helm 3** installed
- **Docker** (for local development)
- **Node.js 20+** and **Python 3.11+** (for app development)

## Platform Installation

### 1. Install the DClaw Operator

```bash
helm repo add dclaw https://charts.dclawstack.io
helm repo update
helm install dclaw-operator dclaw/dclaw-operator \
  --namespace dclaw-system \
  --create-namespace
```

### 2. Install DPanel

```bash
helm install dpanel dclaw/dpanel \
  --namespace dclaw-system \
  --set ingress.host=panel.yourdomain.com
```

### 3. Install Core Infrastructure

```bash
# PostgreSQL operator (CloudNativePG)
kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/main/releases/cnpg-1.24.0.yaml

# Redis
helm install redis bitnami/redis \
  --namespace dclaw-system \
  --set auth.enabled=false

# Ingress controller (if not present)
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace
```

### 4. Verify Installation

```bash
kubectl get pods -n dclaw-system
```

You should see the operator and DPanel pods running.

## Local Development Setup

For local development without Kubernetes:

```bash
git clone https://github.com/dclawstack/dclaw-platform.git
cd dclaw-platform/dpanel
npm install
npm run dev
```

DPanel will be available at `http://localhost:3000`.

## Next Steps

- [Quickstart](./quickstart) — Deploy your first app
- [Configuration](./configuration) — Customize the platform
