# Getting Started

Welcome to the DClaw Platform. This guide will help you install, configure, and deploy your first app.

## What is DClaw?

DClaw is a unified AI platform that delivers 65 AI-native apps under one subscription, running on your own infrastructure. Think of it as "Adobe Creative Cloud for AI" — one platform, many specialized tools, complete data ownership.

## Key Concepts

| Concept | Description |
|---------|-------------|
| **DPanel** | The web-based command center and app launcher |
| **DClaw Operator** | Kubernetes operator that provisions and manages apps |
| **DClawApp CRD** | Custom resource that declares an app's desired state |
| **App** | A deployable unit: Next.js frontend + FastAPI backend + PostgreSQL DB |
| **Shield** | Local PII detection and anonymization layer |

## Quick Links

- [Installation](./installation)
- [Quickstart](./quickstart)
- [Configuration](./configuration)

## Architecture at a Glance

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   DPanel    │────▶│  DClaw App  │────▶│  Backend    │
│  (Next.js)  │     │  (Next.js)  │     │  (FastAPI)  │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                                        ┌──────┴──────┐
                                        │  PostgreSQL  │
                                        │ (CloudNative)│
                                        └─────────────┘
```

## Support

If you run into issues, check the [Troubleshooting](../troubleshooting) section or review [Common Issues](../troubleshooting/common-issues).
