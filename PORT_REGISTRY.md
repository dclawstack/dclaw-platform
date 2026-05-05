# DClaw Stack — Port Registry

> Source of truth for all ports used by DClaw Stack services and agents.  
> **Rule:** Check this file before assigning a new port. Update this file when you change a port.

---

## System Ports (DO NOT USE)

These ports are permanently occupied by system services, PM2, Docker, or K8s tunnels on the dev machine:

| Port | Process | Owner | Notes |
|------|---------|-------|-------|
| 80 | com.docker | Docker/Colima | HTTP ingress |
| 443 | com.docker | Docker/Colima | HTTPS ingress |
| 4040 | ngrok | System | Tunnel dashboard |
| 5000 | ControlCenter | macOS | AirPlay / Control Center |
| 5002 | Python (glm-ocr) | PM2 | OCR service |
| 5200 | Python (seaclip-lite) | PM2 | SeaClip lite |
| 5300 | Python (mortgage-lite) | PM2 | Mortgage lite |
| 5310 | Python (law-lite) | PM2 | Law lite |
| 5432 | com.docker | Docker/Colima | PostgreSQL |
| 6443 | ssh | K8s tunnel | Kubernetes API |
| 6444 | ssh | K8s tunnel | Kubernetes API (alt) |
| 7000 | ControlCenter | macOS | Control Center alt |
| 8000 | com.docker | Docker/Colima | Some Docker service |
| 8080 | kubectl | PM2 (dkubex-portforward) | K8s port-forward |
| 9090 | Python (hub-dashboard) | PM2 | Hub dashboard |
| 9093 | node (app-store) | PM2 | App store |
| 10010 | ssh | K8s tunnel | K8s port-forward |
| 10248-10259 | ssh | K8s tunnel | Kubelet ports |
| 11434 | ollama | System | Local LLM inference |
| 3001 | node | PM2 | Unknown Node service |
| 30080 | ssh | K8s tunnel | K8s port-forward |
| 4321 | node | PM2 | Unknown Node service |
| 4322 | node (hub-docs) | PM2 | Hub docs |
| 49152-49169 | system | macOS | mDNS / limactl |
| 5180 | node | PM2 | Unknown Node service |
| 8174 | Python (aina) | PM2 | Aina service |
| 8175 | node | PM2 | Unknown Node service |

---

## DClaw Stack Ports (ASSIGNED)

| Port | Service | Environment | Owner Agent | Status |
|------|---------|-------------|-------------|--------|
| **3000** | DPanel dev server | Local dev | Shell | ✅ Free |
| **3002** | DClaw Chat frontend dev | Local dev | Shell | ✅ Free |
| **8008** | DClaw Chat backend dev | Local dev | Shell | ✅ Free |
| **8088** | dpanel-api | Local + K8s | Shell | ✅ Assigned |
| **8089** | dclaw-operator metrics | K8s cluster | Shell | ✅ Assigned |
| **8090** | dclaw-rag backend | Local + K8s | Shell | ✅ Assigned |
| **8091** | dclaw-agent backend | Local + K8s | Shell | ✅ Assigned |
| **8092** | dclaw-med backend | Local + K8s | Shell | ✅ Assigned |
| **8093** | dclaw-learn backend | Local + K8s | Shell | ✅ Reserved |
| **8094** | dclaw-code backend | Local + K8s | Shell | ✅ Assigned |
| **8095** | dclaw-seo backend | Local + K8s | Shell | ✅ Reserved |
| **18080** | dclaw-operator metrics (local fallback) | Local dev | Shell | ✅ Assigned |
| **3003** | DClaw Flow frontend dev | Local dev | Shell | ✅ Assigned |
| **3004** | DClaw Med frontend dev | Local dev | Shell | ✅ Assigned |
| **3005** | DClaw Code frontend dev | Local dev | Shell | ✅ Assigned |
| **3008** | DClaw Learn frontend dev | Local dev | Shell | ✅ Assigned |
| **3006** | DClaw SEO frontend dev | Local dev | Shell | ✅ Reserved |
| **3007** | DClaw Create frontend dev | Local dev | Shell | ✅ Reserved |
| **8443** | *Reserved: DClaw HTTPS dev* | Future | — | ✅ Free |

### Port Ranges by Purpose

- **3000–3009:** DClaw frontend dev servers (Next.js apps)
- **8008–8010:** DClaw backend dev servers (FastAPI / Go)
- **8088–8090:** DClaw platform services (dpanel-api, operator metrics)
- **18080–18090:** DClaw platform local fallbacks (when 808x is taken)

---

## Using Ports in Code

### dpanel-api (Go)
```go
// main.go — default port
port := os.Getenv("PORT")
if port == "" {
    port = "8088"  // NOT 8080 (taken by kubectl)
}
```

### DClaw Chat backend (FastAPI / Uvicorn)
```bash
# Local dev — use 8008 to avoid Docker 8000 conflict
uvicorn main:app --host 0.0.0.0 --port 8008

# Production — container uses 8000 (isolated, no conflict)
```

### DClaw Chat frontend (Next.js)
```bash
# Local dev
npm run dev -- --port 3002
```

### dclaw-operator (Go)
```go
// cmd/main.go — metrics port
flag.StringVar(&metricsAddr, "metrics-bind-address", ":8089", ...)
```

---

## Agent Instructions

**Shell Agent:** Before starting any dev server, check this file. If a port is listed as "System Ports", do not use it. If you need a new port, pick from the "DClaw Stack Ports" table or use the 3000–3009 / 8008–8010 / 8088–8090 ranges. Update this file.

**Shield Agent:** Review PRs for hardcoded ports. Ensure new services reference this registry.

**Vault Coordinator:** When specs define new services, assign ports from this registry and update this file.

---

*Last updated: 2026-05-05 by Shell Agent*  
*Next review: When adding a new service or dev environment*
