# Networking

## Network Architecture

```
Internet
    │
    ▼
┌─────────────┐
│  Ingress    │ ← nginx-ingress + cert-manager
│  Controller │
└──────┬──────┘
       │
   ┌───┴───┐
   ▼       ▼
┌─────┐ ┌─────┐
│ /   │ │ /api│
│ Web │ │ API │
└──┬──┘ └──┬──┘
   │       │
   ▼       ▼
┌────────┐ ┌────────┐
│Frontend│ │Backend │
│ Pods   │ │ Pods   │
└───┬────┘ └───┬────┘
    │          │
    │    ┌─────┴─────┐
    │    ▼           ▼
    │ ┌──────┐  ┌────────┐
    │ │Postgre│  │ Ollama │
    │ │SQL    │  │ (sidecar│
    │ └──────┘  │ or svc) │
    │           └────────┘
    │
    └──────► Redis (cache + pub/sub)
```

## Ingress Configuration

Each app gets a dedicated subdomain:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: dclaw-chat
  namespace: dclaw-chat
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  rules:
    - host: chat.dclawstack.io
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: dclaw-chat-frontend
                port:
                  number: 3000
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: dclaw-chat-backend
                port:
                  number: 8000
  tls:
    - hosts:
        - chat.dclawstack.io
      secretName: chat-tls
```

## Network Policies

Apps are isolated by default:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: dclaw-chat-isolation
  namespace: dclaw-chat
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              name: dclaw-chat
    - to:
        - namespaceSelector:
            matchLabels:
              name: dclaw-system
```

## Service Mesh (Future)

For P2 2027, consider adding Linkerd or Istio for:
- mTLS between services
- Traffic splitting for canary deploys
- Circuit breaking
- Distributed tracing
