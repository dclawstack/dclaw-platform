# Common Issues

## App stuck in "Pending" phase

**Symptoms:** `kubectl get dclawapp chat` shows `phase: Pending` for >5 minutes.

**Causes:**
- Operator not running
- Image pull failures
- Insufficient cluster resources

**Solutions:**

```bash
# Check operator logs
kubectl logs -n dclaw-system deployment/dclaw-operator

# Check for image pull errors
kubectl get events -n dclaw-chat --sort-by='.lastTimestamp'

# Check resource availability
kubectl describe nodes
```

## Database connection refused

**Symptoms:** Backend pods crash with `connection refused` to PostgreSQL.

**Solutions:**

```bash
# Check if database cluster is ready
kubectl get clusters -n dclaw-chat

# Check database pod status
kubectl get pods -n dclaw-chat -l app=dclaw-chat-db

# Verify connection string
kubectl get secret dclaw-chat-db-credentials -n dclaw-chat -o yaml
```

## Ingress not routing traffic

**Symptoms:** App returns 404 or connection timeout.

**Solutions:**

```bash
# Verify ingress exists
kubectl get ingress -n dclaw-chat

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# Test DNS resolution
nslookup chat.yourdomain.com

# Verify TLS certificate
kubectl get certificate -n dclaw-chat
```

## DPanel shows "No apps found"

**Symptoms:** DPanel home screen is empty.

**Causes:**
- dpanel-api not reachable
- GitHub rate limiting
- Network issues

**Solutions:**

```bash
# Check dpanel-api
kubectl get pods -n dclaw-system -l app=dpanel-api

# Check DPanel network connectivity
kubectl exec -n dclaw-system deployment/dpanel -- curl -s dpanel-api:8080/health

# Verify GitHub access
curl -I https://raw.githubusercontent.com/dclawstack/dclaw-chat/main/frontend/public/dclaw-manifest.json
```

## Ollama not responding

**Symptoms:** Chat app shows "Local model unavailable".

**Solutions:**

```bash
# Check Ollama pod
kubectl get pods -n dclaw-system -l app=ollama

# Test Ollama API
kubectl exec -n dclaw-system deployment/ollama -- curl -s localhost:11434/api/tags

# Pull a model
kubectl exec -n dclaw-system deployment/ollama -- ollama pull llama3.2
```

## High memory usage

**Symptoms:** Nodes running out of memory, pods evicted.

**Solutions:**

```bash
# Check resource usage
kubectl top nodes
kubectl top pods -A

# Adjust resource limits
kubectl edit dclawapp chat -n dclaw-chat
# Reduce limits.cpu and limits.memory

# Scale down replicas
kubectl patch dclawapp chat -n dclaw-chat --type merge -p '{"spec":{"frontend":{"replicas":1},"backend":{"replicas":1}}}'
```
