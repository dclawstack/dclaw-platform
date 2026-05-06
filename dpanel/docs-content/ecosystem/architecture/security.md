# Security

## Authentication

DClaw uses **Logto** as the identity provider:

- **Local mode:** Built-in username/password (no external dependency)
- **Cloud mode:** Logto Cloud with SSO (Google, Microsoft, SAML)
- **Enterprise mode:** Self-hosted Logto with custom branding

## Authorization

Role-based access control (RBAC):

| Role | Permissions |
|------|-------------|
| **Owner** | Full platform access |
| **Admin** | Install/uninstall apps, manage users |
| **Developer** | Deploy apps, view logs, read docs |
| **User** | Use installed apps, read docs |
| **Guest** | Browse app store, read public docs |

## PII Shield (ClawShield)

ClawShield runs locally and processes all outgoing data:

1. **Detection** — Scan text for PII (names, emails, SSNs, credit cards)
2. **Redaction** — Replace PII with tokens (`[NAME_1]`, `[EMAIL_1]`)
3. **Reconstruction** — On response, restore original values from token map
4. **Audit** — Log all PII touches for compliance

```python
# Example: Shield middleware
from clawshield import Shield

shield = Shield()

@app.middleware("http")
async def shield_middleware(request, call_next):
    if request.body:
        request.body = shield.redact(await request.body())
    response = await call_next(request)
    if response.body:
        response.body = shield.restore(await response.body())
    return response
```

## Compliance

| Framework | Status | Notes |
|-----------|--------|-------|
| GDPR | ✅ | Local-first, data never leaves EU if configured |
| HIPAA | ✅ | PII Shield + audit logs + BAA available |
| SOC 2 | ⏳ | Type II audit scheduled for Q3 2026 |
| ISO 27001 | ⏳ | Planned for 2027 |

## Secrets Management

- Kubernetes Secrets for runtime credentials
- Sealed Secrets for GitOps workflows
- Vault integration for enterprise (future)

## Vulnerability Scanning

- `dclaw-secure` app scans all container images
- Trivy integration in CI/CD
- CVE database auto-updates
