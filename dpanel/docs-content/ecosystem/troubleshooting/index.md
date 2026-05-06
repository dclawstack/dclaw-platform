# Troubleshooting

Common problems and solutions for the DClaw Platform.

## Quick Diagnostics

```bash
# Check operator health
kubectl get pods -n dclaw-system

# Check app status
kubectl get dclawapps

# View operator logs
kubectl logs -n dclaw-system deployment/dclaw-operator

# Check ingress
kubectl get ingress -A

# Verify databases
kubectl get clusters -A
```

## Sections

- [Common Issues](./common-issues) — Frequently encountered problems
- [FAQ](./faq) — Frequently asked questions
