# DClaw Stack — Deployment Status & Action Items

> Last updated: 2026-05-05

---

## 🚀 Current Infrastructure

### Kubernetes Cluster (Colima k3s)
| Component | Status | Details |
|-----------|--------|---------|
| Cluster | ✅ Running | `colima` context, k3s v1.35.0 |
| Ingress Controller | ✅ Installed | nginx-ingress, NodePort 30081 (HTTP) / 30443 (HTTPS) |
| Namespace `dclaw-core` | ✅ Exists | dpanel-api running |
| Namespace `dclaw` | ✅ Exists | Empty (ready for apps) |
| Namespace `ingress-nginx` | ✅ Exists | nginx controller pods |

### Running Pods
```
dclaw-core   dpanel-api-57d69c4687-dwfdk   1/1 Running   NodePort 30080
```

### Port Mappings (Local Access)
| Service | NodePort | Internal Port | URL |
|---------|----------|---------------|-----|
| nginx ingress | 30081 | 80 | http://localhost:30081 |
| dpanel-api (old) | 30080 | 8085 | http://localhost:30080 |

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

### Backends
| App | Port |
|-----|------|
| dpanel-api | 8088 |
| flow-backend | 8088 |
| rag-backend | 8090 |
| agent-backend | 8091 |
| med-backend | 8092 |
| learn-backend | 8093 |
| code-backend | 8094 |
| seo-backend | 8095 |

⚠️ **Port conflict:** flow-backend and dpanel-api both claim 8088 in compose. Need to resolve.

---

## ☸️ Kubernetes Deployment

### Deploy Script
```bash
./k8s-deploy.sh [namespace] [tag]
# Default: ./k8s-deploy.sh dclaw latest
```

### Helm Charts Available
| App | Chart Path | Backend Port |
|-----|-----------|--------------|
| dclaw-flow | `helm/dclaw-flow/` | 8088 |
| dclaw-rag | `helm/dclaw-rag/` | 8090 |
| dclaw-agent | `helm/dclaw-agent/` | 8091 |
| dclaw-med | `helm/dclaw-med/` | 8092 |
| dclaw-code | `helm/dclaw-code/` | 8094 |
| dclaw-learn | `helm/dclaw-learn/` | 8093 |
| dclaw-seo | `helm/dclaw-seo/` | 8095 |

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

### 🟡 Medium Priority

4. **Deploy apps to K8s cluster**
   - Currently only dpanel-api (old version) is running
   - Need ARM64 images first, OR build locally for testing
   - Once images available: `./k8s-deploy.sh dclaw latest`

5. **Fix dpanel-api deployed port**
   - Currently deployed on internal port 8085 (old config)
   - Should be 8088 per PORT_REGISTRY.md
   - Redeploy with updated Helm values

6. **Test Linux installer end-to-end**
   - Spin up a Linux VM (or use Docker in Docker)
   - Run `install.sh`
   - Verify all services start and respond on expected ports

### 🟢 Low Priority / Future

7. **DClaw Chat integration**
   - Chat repo exists but not fully DClaw-ified (no Helm chart in standard location)
   - Port 8008 backend, 3002 frontend

8. **Tauri desktop apps**
   - Deferred until Apple Developer cert obtained
   - Can proceed with Linux `.AppImage` builds anytime

9. **Operator hardening**
   - `reconcileDatabase` is placeholder (logs only)
   - CloudNativePG integration incomplete

---

## 🧪 Test Checklist

### Local K8s (Colima)
- [ ] dpanel-api responds on correct port (8088)
- [ ] nginx ingress routes traffic
- [ ] At least one DClaw app deploys and responds
- [ ] `/health` endpoints return `{"status":"ok"}`

### Docker Compose (Linux)
- [ ] `install.sh` downloads and starts all services
- [ ] Postgres is accessible
- [ ] Each frontend loads in browser
- [ ] Each backend `/health` responds
- [ ] No port conflicts

### Cross-Platform Images
- [ ] `docker pull` works on ARM64 (Mac)
- [ ] `docker pull` works on AMD64 (Linux server)
- [ ] Same tag resolves to correct architecture automatically

---

## 📝 Version Notes

| Component | Version / Commit | Location |
|-----------|-----------------|----------|
| install.sh | `b0112c5` | `dclaw-platform/install.sh` |
| k8s-deploy.sh | `b0112c5` | `dclaw-platform/k8s-deploy.sh` |
| docker-compose.yml | `b0112c5` | `dclaw-platform/docker-compose.yml` |
| DPanel launcher | `e7b4187` | `dclaw-platform/dpanel/` |
| Port Registry | `d0af937` | `dclaw-platform/PORT_REGISTRY.md` |
| Master PRD | `main` | `dclawstack/dclaw-prd` |
