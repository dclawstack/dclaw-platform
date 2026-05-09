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
| **3010** | DClaw Design frontend dev | Local dev | Shell | ✅ Assigned |
| **3011** | DClaw Video frontend dev | Local dev | Shell | ✅ Assigned |
| **3012** | DClaw Research frontend dev | Local dev | Shell | ✅ Assigned |
| **3013** | DClaw Legal frontend dev | Local dev | Shell | ✅ Assigned |
| **3014** | DClaw Finance frontend dev | Local dev | Shell | ✅ Assigned |
| **3015** | DClaw Sales frontend dev | Local dev | Shell | ✅ Assigned |
| **3016** | DClaw Support frontend dev | Local dev | Shell | ✅ Assigned |
| **3017** | DClaw HR frontend dev | Local dev | Shell | ✅ Assigned |
| **3018** | DClaw Translate frontend dev | Local dev | Shell | ✅ Assigned |
| **3019** | DClaw Write frontend dev | Local dev | Shell | ✅ Assigned |
| **3020** | DClaw Meet frontend dev | Local dev | Shell | ✅ Assigned |
| **3021** | DClaw Doc frontend dev | Local dev | Shell | ✅ Assigned |
| **3022** | DClaw Sheet frontend dev | Local dev | Shell | ✅ Assigned |
| **3023** | DClaw Slide frontend dev | Local dev | Shell | ✅ Assigned |
| **3024** | DClaw Email frontend dev | Local dev | Shell | ✅ Assigned |
| **3025** | DClaw Calendar frontend dev | Local dev | Shell | ✅ Assigned |
| **3026** | DClaw Task frontend dev | Local dev | Shell | ✅ Assigned |
| **3027** | DClaw Wiki frontend dev | Local dev | Shell | ✅ Assigned |
| **3028** | DClaw Data frontend dev | Local dev | Shell | ✅ Assigned |
| **3029** | DClaw API frontend dev | Local dev | Shell | ✅ Assigned |
| **3030** | DClaw Test frontend dev | Local dev | Shell | ✅ Assigned |
| **8096** | dclaw-design backend | Local + K8s | Shell | ✅ Assigned |
| **8097** | dclaw-video backend | Local + K8s | Shell | ✅ Assigned |
| **8098** | dclaw-research backend | Local + K8s | Shell | ✅ Assigned |
| **8099** | dclaw-legal backend | Local + K8s | Shell | ✅ Assigned |
| **8100** | dclaw-finance backend | Local + K8s | Shell | ✅ Assigned |
| **8101** | dclaw-sales backend | Local + K8s | Shell | ✅ Assigned |
| **8102** | dclaw-support backend | Local + K8s | Shell | ✅ Assigned |
| **8103** | dclaw-hr backend | Local + K8s | Shell | ✅ Assigned |
| **8104** | dclaw-translate backend | Local + K8s | Shell | ✅ Assigned |
| **8105** | dclaw-write backend | Local + K8s | Shell | ✅ Assigned |
| **8106** | dclaw-meet backend | Local + K8s | Shell | ✅ Assigned |
| **8107** | dclaw-doc backend | Local + K8s | Shell | ✅ Assigned |
| **8108** | dclaw-sheet backend | Local + K8s | Shell | ✅ Assigned |
| **8109** | dclaw-slide backend | Local + K8s | Shell | ✅ Assigned |
| **8110** | dclaw-email backend | Local + K8s | Shell | ✅ Assigned |
| **8111** | dclaw-calendar backend | Local + K8s | Shell | ✅ Assigned |
| **8112** | dclaw-task backend | Local + K8s | Shell | ✅ Assigned |
| **8113** | dclaw-wiki backend | Local + K8s | Shell | ✅ Assigned |
| **8114** | dclaw-data backend | Local + K8s | Shell | ✅ Assigned |
| **8115** | dclaw-api backend | Local + K8s | Shell | ✅ Assigned |
| **8116** | dclaw-test backend | Local + K8s | Shell | ✅ Assigned |
| **8443** | *Reserved: DClaw HTTPS dev* | Future | — | ✅ Free |

### Port Ranges by Purpose

- **3000–3009:** DClaw frontend dev servers (Next.js apps)
- **3010–3019:** New DClaw frontend dev servers (batch 2)
- **8008–8010:** DClaw backend dev servers (FastAPI / Go)
- **8088–8095:** DClaw backend services (batch 1)
- **8096–8105:** DClaw backend services (batch 2)
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

| **3031** | DClaw Deploy frontend dev | Local dev | Shell | ✅ Assigned |
| **3032** | DClaw Monitor frontend dev | Local dev | Shell | ✅ Assigned |
| **3033** | DClaw Secure frontend dev | Local dev | Shell | ✅ Assigned |
| **3034** | DClaw Backup frontend dev | Local dev | Shell | ✅ Assigned |
| **3035** | DClaw Migrate frontend dev | Local dev | Shell | ✅ Assigned |
| **3036** | DClaw Cost frontend dev | Local dev | Shell | ✅ Assigned |
| **3037** | DClaw Carbon frontend dev | Local dev | Shell | ✅ Assigned |
| **3038** | DClaw Compliance frontend dev | Local dev | Shell | ✅ Assigned |
| **3039** | DClaw Audit frontend dev | Local dev | Shell | ✅ Assigned |
| **3040** | DClaw Policy frontend dev | Local dev | Shell | ✅ Assigned |
| **3041** | DClaw Train frontend dev | Local dev | Shell | ✅ Assigned |
| **3042** | DClaw Recruit frontend dev | Local dev | Shell | ✅ Assigned |
| **3043** | DClaw Onboard frontend dev | Local dev | Shell | ✅ Assigned |
| **3044** | DClaw Offboard frontend dev | Local dev | Shell | ✅ Assigned |
| **3045** | DClaw Assets frontend dev | Local dev | Shell | ✅ Assigned |
| **3046** | DClaw Network frontend dev | Local dev | Shell | ✅ Assigned |
| **3047** | DClaw Inventory frontend dev | Local dev | Shell | ✅ Assigned |
| **3048** | DClaw Forecast frontend dev | Local dev | Shell | ✅ Assigned |
| **3049** | DClaw Quality frontend dev | Local dev | Shell | ✅ Assigned |
| **3050** | DClaw Maintenance frontend dev | Local dev | Shell | ✅ Assigned |
| **3051** | DClaw Route frontend dev | Local dev | Shell | ✅ Assigned |
| **3052** | DClaw Warehouse frontend dev | Local dev | Shell | ✅ Assigned |
| **3053** | DClaw Fleet frontend dev | Local dev | Shell | ✅ Assigned |
| **3054** | DClaw Energy frontend dev | Local dev | Shell | ✅ Assigned |
| **3055** | DClaw Water frontend dev | Local dev | Shell | ✅ Assigned |
| **3056** | DClaw Waste frontend dev | Local dev | Shell | ✅ Assigned |
| **3057** | DClaw Building frontend dev | Local dev | Shell | ✅ Assigned |
| **3058** | DClaw Space frontend dev | Local dev | Shell | ✅ Assigned |
| **3059** | DClaw Lease frontend dev | Local dev | Shell | ✅ Assigned |
| **3060** | DClaw Vendor frontend dev | Local dev | Shell | ✅ Assigned |
| **8117** | dclaw-deploy backend | Local + K8s | Shell | ✅ Assigned |
| **8118** | dclaw-monitor backend | Local + K8s | Shell | ✅ Assigned |
| **8119** | dclaw-secure backend | Local + K8s | Shell | ✅ Assigned |
| **8120** | dclaw-backup backend | Local + K8s | Shell | ✅ Assigned |
| **8121** | dclaw-migrate backend | Local + K8s | Shell | ✅ Assigned |
| **8122** | dclaw-cost backend | Local + K8s | Shell | ✅ Assigned |
| **8123** | dclaw-carbon backend | Local + K8s | Shell | ✅ Assigned |
| **8124** | dclaw-compliance backend | Local + K8s | Shell | ✅ Assigned |
| **8125** | dclaw-audit backend | Local + K8s | Shell | ✅ Assigned |
| **8126** | dclaw-policy backend | Local + K8s | Shell | ✅ Assigned |
| **8127** | dclaw-train backend | Local + K8s | Shell | ✅ Assigned |
| **8128** | dclaw-recruit backend | Local + K8s | Shell | ✅ Assigned |
| **8129** | dclaw-onboard backend | Local + K8s | Shell | ✅ Assigned |
| **8130** | dclaw-offboard backend | Local + K8s | Shell | ✅ Assigned |
| **8131** | dclaw-assets backend | Local + K8s | Shell | ✅ Assigned |
| **8132** | dclaw-network backend | Local + K8s | Shell | ✅ Assigned |
| **8133** | dclaw-inventory backend | Local + K8s | Shell | ✅ Assigned |
| **8134** | dclaw-forecast backend | Local + K8s | Shell | ✅ Assigned |
| **8135** | dclaw-quality backend | Local + K8s | Shell | ✅ Assigned |
| **8136** | dclaw-maintenance backend | Local + K8s | Shell | ✅ Assigned |
| **8137** | dclaw-route backend | Local + K8s | Shell | ✅ Assigned |
| **8138** | dclaw-warehouse backend | Local + K8s | Shell | ✅ Assigned |
| **8139** | dclaw-fleet backend | Local + K8s | Shell | ✅ Assigned |
| **8140** | dclaw-energy backend | Local + K8s | Shell | ✅ Assigned |
| **8141** | dclaw-water backend | Local + K8s | Shell | ✅ Assigned |
| **8142** | dclaw-waste backend | Local + K8s | Shell | ✅ Assigned |
| **8143** | dclaw-building backend | Local + K8s | Shell | ✅ Assigned |
| **8144** | dclaw-space backend | Local + K8s | Shell | ✅ Assigned |
| **8145** | dclaw-lease backend | Local + K8s | Shell | ✅ Assigned |
| **8146** | dclaw-vendor backend | Local + K8s | Shell | ✅ Assigned |

| **8147** | dclaw-patent backend | Local + K8s | Shell | ✅ Assigned |
| **8148** | dclaw-trademark backend | Local + K8s | Shell | ✅ Assigned |
| **8149** | dclaw-contract backend | Local + K8s | Shell | ✅ Assigned |
| **8150** | dclaw-continuity backend | Local + K8s | Shell | ✅ Assigned |
| **8151** | dclaw-crisis backend | Local + K8s | Shell | ✅ Assigned |
| **8152** | dclaw-knowledge backend | Local + K8s | Shell | ✅ Assigned |
| **3061** | DClaw Patent frontend dev | Local dev | Shell | ✅ Assigned |
| **3062** | DClaw Trademark frontend dev | Local dev | Shell | ✅ Assigned |
| **3063** | DClaw Contract frontend dev | Local dev | Shell | ✅ Assigned |
| **3064** | DClaw Continuity frontend dev | Local dev | Shell | ✅ Assigned |
| **3065** | DClaw Crisis frontend dev | Local dev | Shell | ✅ Assigned |
| **3066** | DClaw Knowledge frontend dev | Local dev | Shell | ✅ Assigned |
