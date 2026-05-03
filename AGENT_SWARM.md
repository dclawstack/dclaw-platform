# DClaw Agent Swarm — Master Manifest

> This document defines the topology, protocols, and operating procedures for the DClaw Stack multi-agent swarm. It is the source of truth for how autonomous agents coordinate across Kimi web chat sessions, GitHub repositories, and Telegram notifications.

---

## 1. Swarm Topology

```
┌─────────────────────────────────────────────────────────────┐
│                     USER / PRODUCT OWNER                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┴─────────────┐
        │      GENERAL AGENT        │  ← User-facing router
        └─────────────┬─────────────┘
                      │
    ┌─────────────────┼─────────────────┐
    │                 │                 │
┌───▼────┐    ┌──────▼──────┐   ┌──────▼─────┐
│ VAULT  │◄──►│   SHELL     │   │  SHIELD    │
│COORD.  │    │   AGENT     │   │  AGENT     │
└───┬────┘    └──────┬──────┘   └────────────┘
    │                │
    │    ┌───────────┼───────────┐
    │    │           │           │
┌───▼────▼───┐  ┌────▼────┐ ┌───▼────────┐
│   CODE     │  │ RESEARCH│ │   MEMORY   │
│   AGENT    │  │  AGENT  │ │   AGENT    │
└────────────┘  └─────────┘ └────────────┘
```

**Coordination principle:** Agents do not talk to each other in real-time. They coordinate asynchronously through **GitHub commits, PRs, and issues**. The `AGENT_SWARM.md` and agent prompts are the shared memory that keeps all nodes synchronized.

---

## 2. Agent Role Registry

| Agent | Primary Owner | Repos | Responsibility |
|-------|--------------|-------|----------------|
| **Vault Coordinator** | Architecture & PRDs | `dclaw-prd` | Sacred context, specs, decision logs, roadmaps, architecture reviews |
| **Shell Agent** | Implementation & Infra | `dclaw-platform`, `dclaw-chat`, product repos | Code, builds, CI/CD, git ops, deployments, Helm charts |
| **Shield Agent** | Security & Compliance | `dclaw-prd` (SECURITY.md), all repos (reviews) | Threat modeling, audits, compliance checks, security reviews on all PRs |
| **Code Agent** | Deep Engineering | Any repo (assigned by Shell/Vault) | Complex refactoring, algorithm design, code reviews, performance optimization |
| **Research Agent** | R&D & Prototyping | `dclaw-prd` (research docs), temp branches | Technology evaluation, prototype building, benchmarking, PoCs |
| **Memory Agent** | Documentation & KG | `dclaw-prd`, `dclaw-obsidian` | Knowledge graphs, docs, weekly synthesis, context summarization |
| **General Agent** | Routing & User Face | None (orchestrates others) | User intake, routes tasks to specialists, coordinates multi-agent workflows |

---

## 3. Repository Ownership Matrix

| Repository | Visibility | Primary Owner | Secondary Owners | Purpose |
|------------|-----------|---------------|------------------|---------|
| `dclawstack/dclaw-prd` | **Public** | Vault Coordinator | Memory Agent | Sacred context, architecture, product specs |
| `dclawstack/dclaw-platform` | **Public** | Shell Agent | Code Agent, Shield Agent | Operator, DPanel, dpanel-api, platform infra |
| `dclawstack/dclaw-chat` | **Public** | Shell Agent | Code Agent | Chat frontend (Next.js) + backend (FastAPI) |
| `dclawstack/dclaw-obsidian` | **Private** | Memory Agent | Vault Coordinator | Obsidian vault, knowledge graphs, notes |
| `dclawstack/dclaw-enterprise` | **Private** | Shell Agent | Shield Agent | White-label builds, enterprise configs |
| `dclawstack/.github` | **Public** | Shell Agent | — | Org workflows, profile README |

**Rule:** An agent must NOT push to a repo owned by another agent without explicit handoff documented in a GitHub issue or PR description.

---

## 4. Sacred Context Ingestion

Every Kimi web chat agent session must ingest the following public raw URLs on spawn to bootstrap context:

```
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/VISION.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/ARCHITECTURE.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/PRODUCTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/STATUS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/AGENTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/CONVENTIONS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SETUP.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SECURITY.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/README.md
```

After ingestion, read the agent-specific prompt from:
```
https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/agents/[ROLE].md
```

---

## 5. Communication Protocol

### 5.1 GitHub as Shared Memory

Agents communicate exclusively through GitHub artifacts:

| Artifact | Use Case | Required Tags |
|----------|----------|---------------|
| **Commit message** | Work completion signal | `[agent:shell]`, `[agent:vault]`, etc. |
| **PR description** | Handoff between agents | `Handoff: @shield-agent`, `Ready for: @code-agent` |
| **GitHub Issue** | Task assignment, bug reports, decisions | `agent-assigned: shell`, `needs-review: shield` |
| **Repo file** | Persistent shared state | `AGENT_SWARM.md`, `STATUS.md` |

### 5.2 Commit Message Convention

Every commit must be prefixed with the agent role:

```
[agent:vault] docs(arch): add dpanel-api sequence diagram
[agent:shell] feat(dpanel): add install/uninstall UI
[agent:shield] audit(chat): review PII handling in chat.py
[agent:code] refactor(operator): extract reconciler helpers
[agent:research] spike(rag): evaluate pgvector vs milvus
[agent:memory] docs(kg): update component dependency graph
```

### 5.3 Handoff Signals

When an agent completes work that another agent must continue, the commit or PR must include a handoff block:

```markdown
## Handoff

- **From:** @shell-agent
- **To:** @shield-agent
- **Task:** Security review of dpanel-api Dockerfile and RBAC manifests
- **Context:** dpanel-api service reads ConfigMap `dclaw-core/dclaw-apps-registry`. 
  Needs RBAC least-privilege review.
- **Branch:** `feat/dpanel-api-rbac`
- **Files:** `dpanel-api/Dockerfile`, `dpanel-api/main.go`
```

### 5.4 Decision Log

Architecture decisions must be logged in `dclaw-prd/DECISIONS.md` (or created if not existing) by the Vault Coordinator:

```markdown
## 2026-05-03: DPanel API Read Strategy

- **Context:** DPanel needs live app data from the K8s cluster
- **Decision:** Use a dedicated `dpanel-api` Go service reading ConfigMap registry
- **Alternatives considered:** Direct CRD reads from frontend (rejected: CORS/auth complexity)
- **Owner:** @shell-agent
- **Status:** Implemented
```

---

## 6. Current Project State (Snapshot)

### Completed
- [x] Monorepo scaffold with AGENTS.md
- [x] GitHub org `dclawstack` with 6 repos
- [x] Obsidian vault with PARA structure, pushed to `dclaw-obsidian`
- [x] Telegram pipeline (TEST-01 workflow) — operational
- [x] DClaw Operator (Go) — 9-step reconciler, compiles, pushed
- [x] DClaw Chat — Next.js frontend + FastAPI backend + Dockerfile + tests
- [x] Tauri v2 Desktop — unsigned shell with CI
- [x] Agent Swarm Runtime — registry, orchestrator, 5 specialist agents
- [x] Helm charts for `dclaw-chat`
- [x] Sacred Context PRD — 8 docs in `dclaw-prd`
- [x] DPanel — Next.js 16 app on Vercel (https://dpanel-three.vercel.app)
- [x] DPanel App Store UI — install/uninstall with localStorage
- [x] dpanel-api — Go service reading K8s ConfigMap, Dockerfile

### In Progress / Next
- [x] Deploy `dpanel-api` to K8s cluster with RBAC
- [x] Wire DPanel frontend to live `dpanel-api` (set `NEXT_PUBLIC_DPANEL_API_URL`)
- [ ] Operator database reconciliation (CloudNativePG integration)
- [ ] Apple Developer cert for signed Tauri builds
- [ ] Agent swarm web chat backbone (this document)
- [ ] First multi-agent coordinated task

### Active Issues
- `dclawstack/dclaw-panel` repo exists by mistake — needs deletion
- Operator `reconcileDatabase` is placeholder (logs only)
- Tauri macOS unsigned (users must manually allow)

---

## 7. Agent Spawn Procedure (Kimi Web Chat)

### Step 1: Open a new Kimi web chat session
### Step 2: Paste the sacred context ingestion command
```
Please read and internalize the following documents. They define the DClaw Stack sacred context. After reading each, confirm the key takeaways:

https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/VISION.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/ARCHITECTURE.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/PRODUCTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/STATUS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/AGENTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/CONVENTIONS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SETUP.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SECURITY.md
```

### Step 3: Paste the agent-specific role prompt
```
Now read your specific agent prompt and confirm you understand your responsibilities, repo ownership, and handoff protocols:

https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/agents/VAULT_COORDINATOR.md
(or SHELL_AGENT.md, SHIELD_AGENT.md, etc.)
```

### Step 4: Assign an initial task
The agent is now operational. Give it a task using the handoff format.

---

## 8. Security & Trust Boundaries

| Boundary | Rule |
|----------|------|
| **Private repos** | Only Shell Agent and Shield Agent may push to `dclaw-enterprise` and `dclaw-obsidian` |
| **Secrets** | No agent may commit secrets. Use GitHub org/repo secrets only. |
| **Token scope** | Shell Agent uses `VERCEL_TOKEN`, `TELEGRAM_BOT_TOKEN`. Vault Coordinator has no deployment tokens. |
| **Review gate** | All PRs to `main` should ideally be reviewed by Shield Agent for security-critical changes |
| **Escalation** | If an agent is unsure about an architecture decision, it must open an issue tagged `decision-needed: vault` |

---

## 9. State Synchronization

Since Kimi web chat sessions are ephemeral, agents synchronize state through:

1. **GitHub repos** — source of truth for all code and docs
2. **This file (`AGENT_SWARM.md`)** — source of truth for swarm topology and protocols
3. **`dclaw-prd/STATUS.md`** — source of truth for project status
4. **Commit history** — chronological log of all agent actions

Before starting work, an agent should:
1. `git pull` on all repos it owns
2. Read the latest `STATUS.md`
3. Check for open issues tagged with its role

---

## 10. Port Registry

All DClaw Stack services must use assigned ports to avoid conflicts with PM2, Docker, and K8s tunnels on the dev machine.

**Source of truth:** `dclaw-platform/PORT_REGISTRY.md`

**Key rules:**
- NEVER use 8080 (taken by kubectl port-forward)
- NEVER use 8000 (taken by Docker/Colima)
- NEVER use 5000, 7000 (taken by macOS ControlCenter)
- Use DClaw ranges: 3000–3009 (frontend dev), 8008–8010 (backend dev), 8088–8090 (platform services)

## 11. Conventions Summary

- **Language:** TypeScript (frontend), Go (platform/operator), Python (ML/backend), Rust (desktop)
- **Frontend:** Next.js 16 + Tailwind CSS v4 + shadcn/ui (DPanel), Next.js 14 (Chat)
- **Backend:** FastAPI + SQLAlchemy + Pydantic (Chat), Go + client-go (Operator, dpanel-api)
- **Desktop:** Tauri v2 (unsigned macOS for now)
- **K8s:** Helm charts, CloudNativePG (placeholder), ingress-nginx, cert-manager
- **CI/CD:** GitHub Actions
- **Deploy:** Vercel (DPanel), GitHub Container Registry (images), K8s (operator + apps)
- **Visibility:** Public repos = growth/marketing. Private repos = moat (enterprise, vault).

---

*Last updated: 2026-05-03 by Shell Agent*
*Next review: On any agent role change or repo addition*
