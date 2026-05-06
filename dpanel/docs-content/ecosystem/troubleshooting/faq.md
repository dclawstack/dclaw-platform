# Frequently Asked Questions

## General

### What is DClaw?

DClaw is a unified AI platform that delivers 65 AI-native apps under one subscription, running on your own infrastructure.

### How is DClaw different from other AI platforms?

- **Local-first:** Your data stays on your hardware
- **Unified:** One subscription, 65 apps, single auth
- **Open:** BYOK (bring your own API keys), open-source core
- **Enterprise-ready:** On-prem deployment, audit logs, SSO

### What does "local-first" mean?

All AI inference runs locally via Ollama. Cloud LLMs are only used as an opt-in fallback.

## Pricing

### Is DClaw free?

We offer a **Free** tier with 5 apps, local-only, community support.

### Do I need to pay for LLM API calls?

No. DClaw uses local models by default. If you enable cloud fallback, you bring your own API keys (OpenRouter, OpenAI, Anthropic).

### What is the Team tier?

$199/month for all 65 apps, SSO, audit logs, and shared workspaces.

## Technical

### What hardware do I need?

**Minimum:**
- 4 CPU cores
- 8 GB RAM
- 50 GB storage

**Recommended (for local LLMs):**
- 8+ CPU cores
- 32 GB RAM
- Apple Silicon M4 or NVIDIA GPU
- 200 GB SSD storage

### Can I run DClaw without Kubernetes?

Yes, for local development. Use `docker-compose` in `dclaw-platform`.

### How do I update an app?

Update the `version` field in the app's `DClawApp` CRD. The operator will perform a rolling update.

### Can I develop my own apps?

Yes! Follow the [app development guide](../getting-started/configuration) to scaffold a new app.

### How do I backup my data?

Use the `dclaw-backup` app or configure CloudNativePG scheduled backups:

```yaml
spec:
  database:
    backups:
      enabled: true
      schedule: "0 2 * * *"
      retention: 7d
```

## Security

### Is my data secure?

Yes. DClaw is local-first by design. PII Shield anonymizes data before any cloud call.

### Is DClaw HIPAA compliant?

Yes, with the Enterprise tier. We provide a BAA and audit logs.

### Can I use my own identity provider?

Yes. DClaw supports SAML 2.0, OIDC, and LDAP via Logto.
