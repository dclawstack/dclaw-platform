# DClaw Stack — Deployment Status & Action Items

> Last updated: 2026-05-05

---

## 🚀 Current Infrastructure

### Kubernetes Cluster (Colima k3s)
| Component | Status | Details |
|-----------|--------|---------|
| Cluster | ✅ Running | `colima` context, k3s v1.35.0 |
| Ingress Controller | ✅ Installed | nginx-ingress, NodePort 30081 (HTTP) / 30443 (HTTPS) |
| Namespace `dclaw-core` | ✅ Exists | dpanel-api running (old port 8085) |
| Namespace `dclaw` | ✅ Exists | dclaw-learn deployed and tested |
| Namespace `ingress-nginx` | ✅ Exists | nginx controller pods |

### Running Pods (dclaw namespace)
```
dclaw-learn-backend    2/2 Running   NodePort 30093
dclaw-learn-frontend   2/2 Running   NodePort 30005
postgres               1/1 Running   ClusterIP 5432
```

### End-to-End Test Results ✅
| Test | URL | Result |
|------|-----|--------|
| Backend /health | http://localhost:30093/health | ✅ `{"status":"ok","version":"0.1.0"}` |
| Frontend | http://localhost:30005/ | ✅ HTTP 200, DClaw Learn landing page |

---

## 🐧 Linux Installer

### One-Line Install
```bash
curl -fsSL https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/install.sh | bash
```

### What It Installs
- **Directory:** `~/dclaw-stack/`
- **Compose file:** `docker-compose.yml` (shared Postgres + all apps)
- **Default Postgres password:** `dclaw` (edit `.env` to change)

### Exposed Ports
| App | URL |
|-----|-----|
| DPanel | http://localhost:3000 |
| DClaw Flow | http://localhost:3003 |
| DClaw RAG | http://localhost:3009 |
| DClaw Med | http://localhost:3004 |
| DClaw Code | http://localhost:3005 |
| DClaw Learn | http://localhost:3006 |
| DClaw SEO | http://localhost:3007 |
| DClaw Create | http://localhost:3008 |

⚠️ **Port conflict in compose:** dpanel-api and flow-backend both bind to host port 8088. Need to resolve.

---

## ☸️ Kubernetes Deployment

### Deploy Script
```bash
./k8s-deploy.sh [namespace] [tag]
# Default: ./k8s-deploy.sh dclaw latest
```

### Helm Charts Available
| App | Chart Path | Backend Port | Frontend Port | Status |
|-----|-----------|--------------|---------------|--------|
| dclaw-flow | `helm/dclaw-flow/` | 8088 | 3003 | ⏳ Needs ARM64 image |
| dclaw-rag | `helm/dclaw-rag/` | 8090 | 3009 | ⏳ Needs ARM64 image |
| dclaw-agent | `helm/dclaw-agent/` | 8091 | — | ⏳ Needs ARM64 image |
| dclaw-med | `helm/dclaw-med/` | 8092 | 3004 | ⏳ Needs ARM64 image |
| dclaw-code | `helm/dclaw-code/` | 8094 | 3005 | ⏳ Needs ARM64 image |
| **dclaw-learn** | **`helm/dclaw-learn/`** | **8093** | **3005** | **✅ Tested & Working** |
| dclaw-seo | `helm/dclaw-seo/` | 8095 | 3006 | ⏳ Needs ARM64 image |
| dclaw-create | `helm/dclaw-create/` | — | 3007 | ⏳ Needs ARM64 image |

### K8s Deployment Issues Found During Testing
1. **Missing ServiceAccount template** — Helm charts reference `serviceAccountName` but no `serviceaccount.yaml` template exists. Workaround: `--set serviceAccount.create=false`
2. **CloudNativePG dependency** — Charts depend on `cloudnative-pg` Helm chart but CRDs aren't installed in cluster. Workaround: `--set database.enabled=false` + deploy separate Postgres
3. **Frontend targetPort mismatch** — Helm values had `frontend.port: 3000` but Dockerfile exposes `3005`. Fix: Align ports in values.yaml
4. **AMD64-only GHCR images** — Can't run on ARM64 Mac/Colima. Fix: Build locally or wait for multi-arch CI rebuilds

---

## 📦 Container Registry Status

### Current Registry: GHCR (`ghcr.io/dclawstack/*`)
| Image | Architecture | Status |
|-------|-------------|--------|
| `dpanel-api` | Multi-arch? | ✅ Exists locally + GHCR |
| `dclaw-*-backend` | AMD64 only | ⚠️ Built on GitHub Actions (x86 runners) |
| `dclaw-*-frontend` | AMD64 only | ⚠️ Built on GitHub Actions (x86 runners) |

### Problem
All CI-built images are **AMD64 only**. They fail on ARM64 (Apple Silicon / Colima):
```
no matching manifest for linux/arm64/v8
```

### Fix Applied
Updated all CI workflows to build `linux/amd64,linux/arm64` multi-arch images.
**Waiting for:** Next push to each repo to trigger rebuilds.

### Workaround for Local Testing (ARM64 Mac)
Build images locally and import into k3s containerd:
```bash
# Build
docker build -t dclaw-learn-backend:local ./backend
docker build -t dclaw-learn-frontend:local ./frontend

# Save & import into k3s
docker save dclaw-learn-backend:local dclaw-learn-frontend:local | \
  colima ssh -- sudo ctr -a /run/k3s/containerd/containerd.sock -n k8s.io images import -

# Retag without docker.io/library prefix
colima ssh -- sudo ctr -a /run/k3s/containerd/containerd.sock -n k8s.io images tag \
  docker.io/library/dclaw-learn-backend:local dclaw-learn-backend:local
```

---

## 🎯 Action Items

### 🔴 High Priority

1. **Add Docker Hub mirror pushes to all 8 CI workflows**
   - Add `docker/login-action` for Docker Hub
   - Add `dclawstack/{image}:latest` tags alongside GHCR tags
   - Create `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets in GitHub org
   - Repos to update: flow, rag, agent, med, code, learn, seo, create

2. **Fix docker-compose.yml port conflicts**
   - dpanel-api and flow-backend both bind to host port 8088
   - Resolution: Move dpanel-api to 8088 (correct port), flow-backend to separate port or use internal networking only

3. **Rebuild ARM64 images**
   - Trigger CI rebuilds for all repos (multi-arch now configured)
   - Verify `docker pull ghcr.io/dclawstack/dclaw-learn-backend:latest` works on ARM64

4. **Fix Helm chart serviceAccount template**
   - Add missing `serviceaccount.yaml` template OR set `serviceAccount.create: false` default
   - Applies to: flow, rag, agent, med, code, learn, seo, create

5. **Fix Helm chart frontend port alignment**
   - Ensure `frontend.port` in values.yaml matches Dockerfile `EXPOSE` port
   - Applies to all apps with frontends

### 🟡 Medium Priority

6. **Deploy remaining apps to K8s cluster**
   - dclaw-learn is deployed and tested ✅
   - Need ARM64 images or local builds for: flow, rag, agent, med, code, seo, create

7. **Fix dpanel-api deployed port**
   - Currently deployed on internal port 8085 (old config)
   - Should be 8088 per PORT_REGISTRY.md
   - Redeploy with updated Helm values

8. **Install CloudNativePG operator (optional)**
   - Required if using `database.enabled=true` in Helm charts
   - Alternative: Use shared Postgres deployment for simplicity

9. **Test Linux installer end-to-end**
   - Spin up a Linux VM (or use Docker in Docker)
   - Run `install.sh`
   - Verify all services start and respond on expected ports

### 🟢 Low Priority / Future

10. **DClaw Chat integration**
    - Chat repo exists but not fully DClaw-ified (no Helm chart in standard location)
    - Port 8008 backend, 3002 frontend

11. **Tauri desktop apps**
    - Deferred until Apple Developer cert obtained
    - Can proceed with Linux `.AppImage` builds anytime

12. **Operator hardening**
    - `reconcileDatabase` is placeholder (logs only)
    - CloudNativePG integration incomplete

---

## 🧪 Test Checklist

### Local K8s (Colima)
- [x] nginx ingress installed and running
- [x] dclaw-learn backend responds on NodePort 30093
- [x] dclaw-learn frontend serves on NodePort 30005
- [x] `/health` endpoint returns `{"status":"ok"}`
- [ ] dpanel-api responds on correct port (8088)
- [ ] At least 3 DClaw apps deployed simultaneously
- [ ] Ingress routes traffic by host header

### Docker Compose (Linux)
- [ ] `install.sh` downloads and starts all services
- [ ] Postgres is accessible
- [ ] Each frontend loads in browser
- [ ] Each backend `/health` responds
- [ ] No port conflicts

### Cross-Platform Images
- [ ] `docker pull` works on ARM64 (Mac) — pending CI rebuild
- [ ] `docker pull` works on AMD64 (Linux server)
- [ ] Same tag resolves to correct architecture automatically

---

## 📝 Version Notes

| Component | Version / Commit | Location |
|-----------|-----------------|----------|
| install.sh | `ebca1e0` | `dclaw-platform/install.sh` |
| k8s-deploy.sh | `ebca1e0` | `dclaw-platform/k8s-deploy.sh` |
| docker-compose.yml | `ebca1e0` | `dclaw-platform/docker-compose.yml` |
| DPanel launcher | `e7b4187` | `dclaw-platform/dpanel/` |
| Port Registry | `d0af937` | `dclaw-platform/PORT_REGISTRY.md` |
| Master PRD | `main` | `dclawstack/dclaw-prd` |
| DEPLOYMENT_STATUS | `ebca1e0` | `dclaw-platform/DEPLOYMENT_STATUS.md` |
