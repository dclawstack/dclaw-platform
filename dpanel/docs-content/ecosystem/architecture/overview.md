# Architecture Overview

## Design Principles

1. **Local-first AI** — Ollama runs on your hardware. No data leaves unless you explicitly enable cloud fallback.
2. **BYOK** — Bring your own API keys. We don't resell LLM access.
3. **K8s-native** — Every app is a container. Scale, isolate, migrate freely.
4. **PII Shield** — Automatic anonymization before any cloud call. HIPAA/GDPR ready.
5. **White-label** — Enterprise customers get branded desktop apps with their own certs.

## Per-App Architecture

Every app follows the same pattern:

```
App (e.g., dclaw-chat)
├── Next.js Frontend        → Serves UI, talks to backend API
├── FastAPI Backend         → Business logic, DB access, AI proxies
├── PostgreSQL              → Per-app database (CloudNativePG cluster)
├── Tauri Desktop           → WebView wrapper with native features
└── Helm Chart              → K8s manifests for deployment
```

## Data Flow

```
User Input
    │
    ▼
┌─────────────────┐
│  ClawShield     │ ← PII detection/redaction
│  (local, always)│
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐  ┌──────────┐
│ Ollama │  │ OpenRouter│ ← Local LLM vs Cloud fallback
│ (local)│  │ (cloud)   │
└───┬────┘  └─────┬────┘
    │             │
    └──────┬──────┘
           ▼
    ┌─────────────┐
    │  Response   │
    │  Stream     │
    └─────────────┘
```

## Tech Stack

| Layer | Technology | Version |
|-------|------------|---------|
| K8s Operator | Go + controller-runtime | 0.18 |
| DPanel + Frontends | Next.js 16 + Tailwind + shadcn/ui | latest |
| App Backends | FastAPI + uvicorn + SQLAlchemy 2.0 | latest |
| Databases | PostgreSQL 16 + CloudNativePG | latest |
| Desktop | Tauri v2 | 2.0.0-beta |
| Local LLMs | Ollama | latest |
| Cloud LLMs | OpenRouter + Kimi K2.5 | API |
| Message Bus | Redis | 7.x |
| Object Storage | MinIO | latest |
| Monitoring | Prometheus + Grafana | latest |
| Auth | Logto | latest |
| Billing | Stripe | API |
