# Stack Reference

## Platform Stack

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

## Per-App Stack

Every app uses:

| Component | Technology |
|-----------|------------|
| Frontend | Next.js 14.2.28 |
| Backend | FastAPI + Pydantic |
| Database | PostgreSQL + asyncpg |
| Packaging | Helm 3 |
| Desktop | Tauri v2 (optional) |

## Port Registry

| App | Frontend | Backend |
|-----|----------|---------|
| DPanel | 3000 | — |
| DClaw Chat | 30002 | 8000 |
| DClaw Flow | 30003 | 8001 |
| DClaw Med | 30004 | 8002 |
| DClaw Learn | 30005 | 8003 |
| DClaw SEO | 30006 | 8004 |
| DClaw Create | 30007 | 8005 |
| DClaw RAG | 30008 | 8006 |
| DClaw Agent | 30091 | 8091 |
| DClaw Code | 30012 | 8012 |
| ... | ... | ... |

See `PORT_REGISTRY.md` in `dclaw-platform` for the complete list.
