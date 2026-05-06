# Architecture

This section describes the system architecture of the DClaw Platform.

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          DClaw Platform                             │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│  │   Users     │    │   Users     │    │  Enterprise │            │
│  │  (Web)      │    │  (Desktop)  │    │  (On-Prem)  │            │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘            │
│         └──────────────────┼──────────────────┘                    │
│                            ▼                                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  DPanel (Next.js + Tailwind)                                │   │
│  │  • 9-dot app launcher                                       │   │
│  │  • App store (install/uninstall)                            │   │
│  │  • Billing dashboard                                        │   │
│  │  • Team management                                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                            │                                        │
│         ┌──────────────────┼──────────────────┐                    │
│         ▼                  ▼                  ▼                    │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│  │ DClaw Chat  │    │ DClaw Flow  │    │ DClaw Med   │            │
│  │ (Product 1) │    │ (Product 2) │    │ (Product 3) │            │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘            │
│         └──────────────────┼──────────────────┘                    │
│                            ▼                                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  DClaw Operator (Go + controller-runtime)                   │   │
│  │  • Watches DClawApp CRDs                                    │   │
│  │  • Auto-provisions: namespace, deployments, DB, ingress     │   │
│  │  • Network policies for isolation                           │   │
│  │  • Resource quotas per app                                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                            │                                        │
│         ┌──────────────────┼──────────────────┐                    │
│         ▼                  ▼                  ▼                    │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│  │  Identity   │    │   Billing   │    │   Shield    │            │
│  │  (Logto)    │    │  (Stripe)   │    │(PII protect)│            │
│  └─────────────┘    └─────────────┘    └─────────────┘            │
│                                                                     │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│  │   Voice     │    │   Updater   │    │  Message Bus │           │
│  │(Wake word)  │    │(Auto-update)│    │   (Redis)   │            │
│  └─────────────┘    └─────────────┘    └─────────────┘            │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  Infrastructure Layer                                       │   │
│  │  • Kubernetes (k3s for local, managed for cloud)            │   │
│  │  • PostgreSQL (CloudNativePG)                               │   │
│  │  • Redis (message bus + cache)                              │   │
│  │  • MinIO (object storage)                                   │   │
│  │  • Prometheus + Grafana (monitoring)                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Sub-Sections

- [Overview](./overview) — High-level design principles
- [Operator](./operator) — How the K8s operator works
- [Networking](./networking) — Ingress, service mesh, and network policies
- [Security](./security) — Authentication, authorization, and Shield PII protection
