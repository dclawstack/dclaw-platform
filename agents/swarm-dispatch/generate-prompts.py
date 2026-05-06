#!/usr/bin/env python3
"""Generate per-app v1.0 swarm prompts (spec + implementer) for all 67 DClaw apps."""

import json
import csv
from pathlib import Path
from colorsys import rgb_to_hls, hls_to_rgb

APPS = [
    {"id": "chat", "name": "DClaw Chat", "category": "Communication", "tagline": "AI conversations that remember", "color": "#3B82F6", "status": "P0 Active", "priority": 0, "icon": "💬"},
    {"id": "flow", "name": "DClaw Flow", "category": "Automation", "tagline": "Connect anything, automate everything", "color": "#10B981", "status": "P1 Queued", "priority": 1, "icon": "🌊"},
    {"id": "agent", "name": "DClaw Agent", "category": "Platform", "tagline": "Build, share, and sell AI agents", "color": "#8B5CF6", "status": "P1 Queued", "priority": 1, "icon": "🤖"},
    {"id": "rag", "name": "DClaw RAG", "category": "Platform", "tagline": "Universal knowledge retrieval", "color": "#F59E0B", "status": "P1 Queued", "priority": 1, "icon": "🔍"},
    {"id": "med", "name": "DClaw Med", "category": "Healthcare", "tagline": "Clinical intelligence at your fingertips", "color": "#EF4444", "status": "P2 Queued", "priority": 2, "icon": "🏥"},
    {"id": "learn", "name": "DClaw Learn", "category": "Education", "tagline": "Adaptive learning that works", "color": "#6366F1", "status": "P2 Queued", "priority": 2, "icon": "📚"},
    {"id": "code", "name": "DClaw Code", "category": "Development", "tagline": "AI-native IDE inside your desktop", "color": "#1F2937", "status": "P2 Queued", "priority": 2, "icon": "💻"},
    {"id": "seo", "name": "DClaw SEO", "category": "Marketing", "tagline": "Rank higher with AI", "color": "#F97316", "status": "P3 Queued", "priority": 3, "icon": "📈"},
    {"id": "create", "name": "DClaw Create", "category": "Media", "tagline": "Generate anything", "color": "#EC4899", "status": "P3 Queued", "priority": 3, "icon": "🎨"},
    {"id": "legal", "name": "DClaw Legal", "category": "Legal", "tagline": "Contract review and case law research", "color": "#475569", "status": "P4 Planned", "priority": 4, "icon": "⚖️"},
    {"id": "finance", "name": "DClaw Finance", "category": "Finance", "tagline": "Financial modeling and risk analysis", "color": "#14B8A6", "status": "P4 Planned", "priority": 4, "icon": "💰"},
    {"id": "sales", "name": "DClaw Sales", "category": "Sales", "tagline": "CRM AI, email sequences, and forecasting", "color": "#F97316", "status": "P4 Planned", "priority": 4, "icon": "🎯"},
    {"id": "support", "name": "DClaw Support", "category": "Support", "tagline": "Ticket resolution and knowledge base AI", "color": "#06B6D4", "status": "P4 Planned", "priority": 4, "icon": "🎧"},
    {"id": "hr", "name": "DClaw HR", "category": "HR", "tagline": "Resume screening, interview prep, and onboarding", "color": "#D946EF", "status": "P4 Planned", "priority": 4, "icon": "👥"},
    {"id": "design", "name": "DClaw Design", "category": "Design", "tagline": "UI/UX generation and design system management", "color": "#A855F7", "status": "P4 Planned", "priority": 4, "icon": "🖌️"},
    {"id": "translate", "name": "DClaw Translate", "category": "Language", "tagline": "Real-time translation and localization", "color": "#84CC16", "status": "P4 Planned", "priority": 4, "icon": "🌐"},
    {"id": "write", "name": "DClaw Write", "category": "Content", "tagline": "Long-form writing, blogging, and books", "color": "#F43F5E", "status": "P4 Planned", "priority": 4, "icon": "✍️"},
    {"id": "meet", "name": "DClaw Meet", "category": "Meetings", "tagline": "Transcription, action items, and summaries", "color": "#64748B", "status": "P4 Planned", "priority": 4, "icon": "🎥"},
    {"id": "doc", "name": "DClaw Doc", "category": "Documents", "tagline": "Smart documents and collaborative editing", "color": "#0EA5E9", "status": "P4 Planned", "priority": 4, "icon": "📄"},
    {"id": "sheet", "name": "DClaw Sheet", "category": "Spreadsheets", "tagline": "AI-powered spreadsheet alternative", "color": "#22C55E", "status": "P4 Planned", "priority": 4, "icon": "📊"},
    {"id": "slide", "name": "DClaw Slide", "category": "Presentations", "tagline": "AI-generated decks", "color": "#E11D48", "status": "P4 Planned", "priority": 4, "icon": "🖼️"},
    {"id": "email", "name": "DClaw Mail", "category": "Email", "tagline": "Smart inbox, auto-reply, and scheduling", "color": "#3B82F6", "status": "P4 Planned", "priority": 4, "icon": "📧"},
    {"id": "calendar", "name": "DClaw Calendar", "category": "Scheduling", "tagline": "AI scheduling and conflict resolution", "color": "#F59E0B", "status": "P4 Planned", "priority": 4, "icon": "📅"},
    {"id": "task", "name": "DClaw Task", "category": "Productivity", "tagline": "Smart to-do and project management", "color": "#8B5CF6", "status": "P4 Planned", "priority": 4, "icon": "✅"},
    {"id": "wiki", "name": "DClaw Wiki", "category": "Knowledge", "tagline": "Internal Wikipedia with AI search", "color": "#A855F7", "status": "P4 Planned", "priority": 4, "icon": "📖"},
    {"id": "data", "name": "DClaw Data", "category": "Analytics", "tagline": "Natural language data analysis", "color": "#14B8A6", "status": "P4 Planned", "priority": 4, "icon": "📉"},
    {"id": "api", "name": "DClaw API", "category": "Developer", "tagline": "API design, testing, and documentation", "color": "#6366F1", "status": "P4 Planned", "priority": 4, "icon": "🔌"},
    {"id": "test", "name": "DClaw Test", "category": "QA", "tagline": "Automated testing and bug prediction", "color": "#EF4444", "status": "P4 Planned", "priority": 4, "icon": "🧪"},
    {"id": "deploy", "name": "DClaw Deploy", "category": "DevOps", "tagline": "CI/CD pipeline builder", "color": "#1F2937", "status": "P4 Planned", "priority": 4, "icon": "🚀"},
    {"id": "monitor", "name": "DClaw Monitor", "category": "Observability", "tagline": "AI-powered alerting and root cause analysis", "color": "#F97316", "status": "P4 Planned", "priority": 4, "icon": "📡"},
    {"id": "secure", "name": "DClaw Secure", "category": "Security", "tagline": "Threat detection and vulnerability scanning", "color": "#DC2626", "status": "P4 Planned", "priority": 4, "icon": "🔒"},
    {"id": "backup", "name": "DClaw Backup", "category": "Infrastructure", "tagline": "Intelligent backup and recovery", "color": "#475569", "status": "P4 Planned", "priority": 4, "icon": "💾"},
    {"id": "migrate", "name": "DClaw Migrate", "category": "Infrastructure", "tagline": "Cloud migration assistant", "color": "#64748B", "status": "P4 Planned", "priority": 4, "icon": "🚚"},
    {"id": "cost", "name": "DClaw Cost", "category": "FinOps", "tagline": "Cloud cost optimization", "color": "#14B8A6", "status": "P4 Planned", "priority": 4, "icon": "💵"},
    {"id": "carbon", "name": "DClaw Carbon", "category": "Sustainability", "tagline": "Carbon footprint tracking", "color": "#22C55E", "status": "P4 Planned", "priority": 4, "icon": "🌱"},
    {"id": "compliance", "name": "DClaw Compliance", "category": "Governance", "tagline": "GDPR, SOC2, and HIPAA automation", "color": "#6366F1", "status": "P4 Planned", "priority": 4, "icon": "📋"},
    {"id": "audit", "name": "DClaw Audit", "category": "Governance", "tagline": "Automated audit trails", "color": "#475569", "status": "P4 Planned", "priority": 4, "icon": "🔍"},
    {"id": "policy", "name": "DClaw Policy", "category": "Governance", "tagline": "Policy management and enforcement", "color": "#8B5CF6", "status": "P4 Planned", "priority": 4, "icon": "📜"},
    {"id": "train", "name": "DClaw Train", "category": "L&D", "tagline": "Employee training and certification tracking", "color": "#D946EF", "status": "P4 Planned", "priority": 4, "icon": "🎓"},
    {"id": "recruit", "name": "DClaw Recruit", "category": "Talent", "tagline": "Job posting and candidate ranking", "color": "#EC4899", "status": "P4 Planned", "priority": 4, "icon": "🤝"},
    {"id": "onboard", "name": "DClaw Onboard", "category": "HR", "tagline": "Automated employee onboarding", "color": "#10B981", "status": "P4 Planned", "priority": 4, "icon": "🚀"},
    {"id": "offboard", "name": "DClaw Offboard", "category": "HR", "tagline": "Secure offboarding automation", "color": "#EF4444", "status": "P4 Planned", "priority": 4, "icon": "🚪"},
    {"id": "assets", "name": "DClaw Assets", "category": "IT", "tagline": "Hardware and software asset management", "color": "#64748B", "status": "P4 Planned", "priority": 4, "icon": "🖥️"},
    {"id": "network", "name": "DClaw Network", "category": "IT", "tagline": "Network monitoring and optimization", "color": "#3B82F6", "status": "P4 Planned", "priority": 4, "icon": "🌐"},
    {"id": "inventory", "name": "DClaw Inventory", "category": "Operations", "tagline": "Supply chain intelligence", "color": "#F59E0B", "status": "P4 Planned", "priority": 4, "icon": "📦"},
    {"id": "forecast", "name": "DClaw Forecast", "category": "Operations", "tagline": "Demand forecasting", "color": "#06B6D4", "status": "P4 Planned", "priority": 4, "icon": "🔮"},
    {"id": "quality", "name": "DClaw Quality", "category": "Manufacturing", "tagline": "Quality control AI", "color": "#22C55E", "status": "P4 Planned", "priority": 4, "icon": "✨"},
    {"id": "maintenance", "name": "DClaw Maintenance", "category": "Operations", "tagline": "Predictive maintenance", "color": "#F97316", "status": "P4 Planned", "priority": 4, "icon": "🔧"},
    {"id": "route", "name": "DClaw Route", "category": "Logistics", "tagline": "Route optimization", "color": "#0EA5E9", "status": "P4 Planned", "priority": 4, "icon": "🗺️"},
    {"id": "warehouse", "name": "DClaw Warehouse", "category": "Logistics", "tagline": "Warehouse automation", "color": "#8B5CF6", "status": "P4 Planned", "priority": 4, "icon": "🏭"},
    {"id": "fleet", "name": "DClaw Fleet", "category": "Logistics", "tagline": "Fleet management", "color": "#6366F1", "status": "P4 Planned", "priority": 4, "icon": "🚛"},
    {"id": "energy", "name": "DClaw Energy", "category": "Utilities", "tagline": "Energy consumption optimization", "color": "#F59E0B", "status": "P4 Planned", "priority": 4, "icon": "⚡"},
    {"id": "water", "name": "DClaw Water", "category": "Utilities", "tagline": "Water management", "color": "#0EA5E9", "status": "P4 Planned", "priority": 4, "icon": "💧"},
    {"id": "waste", "name": "DClaw Waste", "category": "Sustainability", "tagline": "Waste reduction AI", "color": "#22C55E", "status": "P4 Planned", "priority": 4, "icon": "♻️"},
    {"id": "building", "name": "DClaw Building", "category": "Real Estate", "tagline": "Smart building management", "color": "#8B5CF6", "status": "P4 Planned", "priority": 4, "icon": "🏢"},
    {"id": "space", "name": "DClaw Space", "category": "Real Estate", "tagline": "Office space optimization", "color": "#A855F7", "status": "P4 Planned", "priority": 4, "icon": "🪑"},
    {"id": "lease", "name": "DClaw Lease", "category": "Real Estate", "tagline": "Lease management", "color": "#14B8A6", "status": "P4 Planned", "priority": 4, "icon": "📜"},
    {"id": "vendor", "name": "DClaw Vendor", "category": "Procurement", "tagline": "Vendor evaluation and management", "color": "#F59E0B", "status": "P4 Planned", "priority": 4, "icon": "🤝"},
    {"id": "contract", "name": "DClaw Contract", "category": "Legal", "tagline": "Contract lifecycle management", "color": "#475569", "status": "P4 Planned", "priority": 4, "icon": "📑"},
    {"id": "risk", "name": "DClaw Risk", "category": "Governance", "tagline": "Enterprise risk management", "color": "#DC2626", "status": "P4 Planned", "priority": 4, "icon": "⚠️"},
    {"id": "crisis", "name": "DClaw Crisis", "category": "Operations", "tagline": "Crisis response planning", "color": "#EF4444", "status": "P4 Planned", "priority": 4, "icon": "🚨"},
    {"id": "continuity", "name": "DClaw Continuity", "category": "Operations", "tagline": "Business continuity planning", "color": "#10B981", "status": "P4 Planned", "priority": 4, "icon": "🛡️"},
    {"id": "knowledge", "name": "DClaw Knowledge", "category": "Platform", "tagline": "Enterprise knowledge graph", "color": "#6366F1", "status": "P4 Planned", "priority": 4, "icon": "🧠"},
    {"id": "research", "name": "DClaw Research", "category": "R&D", "tagline": "Research paper analysis and synthesis", "color": "#8B5CF6", "status": "P4 Planned", "priority": 4, "icon": "🔬"},
    {"id": "patent", "name": "DClaw Patent", "category": "Legal", "tagline": "Patent search and analysis", "color": "#3B82F6", "status": "P4 Planned", "priority": 4, "icon": "📋"},
    {"id": "trademark", "name": "DClaw Trademark", "category": "Legal", "tagline": "Trademark monitoring", "color": "#F59E0B", "status": "P4 Planned", "priority": 4, "icon": "™️"},
    {"id": "video", "name": "DClaw Video", "category": "Media", "tagline": "AI-powered video editing and generation", "color": "#EC4899", "status": "P4 Planned", "priority": 4, "icon": "🎬"},
]


def color_vars(color: str) -> tuple:
    def hex_to_rgb(h):
        h = h.lstrip("#")
        return tuple(int(h[i:i+2], 16) / 255.0 for i in (0, 2, 4))
    def rgb_to_hex(r, g, b):
        return f"#{int(r*255):02X}{int(g*255):02X}{int(b*255):02X}"
    r, g, b = hex_to_rgb(color)
    h, l, s = rgb_to_hls(r, g, b)
    light = rgb_to_hex(*hls_to_rgb(h, min(l + 0.15, 0.95), s))
    deep = rgb_to_hex(*hls_to_rgb(h, max(l - 0.15, 0.05), s))
    wash = rgb_to_hex(*hls_to_rgb(h, 0.96, s * 0.3))
    return light, deep, wash


def generate_spec_prompt(app: dict, idx: int) -> str:
    app_id = app["id"]
    name = app["name"]
    category = app["category"]
    tagline = app["tagline"]
    color = app["color"]
    status = app["status"]
    icon = app["icon"]
    frontend_port = 3000 + idx
    backend_port = 8000 + idx
    db_name = f"dclaw_{app_id}"
    light, deep, wash = color_vars(color)

    return f"""# {name} v1.0 — Spec Writer Prompt

> **Agent Role:** Product Architect + Full-Stack Engineer + Design Systems Specialist
> **Mission:** Research the market and produce a complete v1.0 build specification for {name}.
> **Output:** A single `v1.0-spec.md` document.
> **Priority:** {status}

---

## App Identity
```yaml
app_id: {app_id}
name: {name}
category: {category}
tagline: {tagline}
current_version: 0.1.0
status: {status}
repo: https://github.com/dclawstack/dclaw-{app_id}
docs: https://docs.dclawstack.io/apps/{app_id}
frontend_port: {frontend_port}
backend_port: {backend_port}
db_name: {db_name}
brand_color: {color}
icon: {icon}
```

## Existing Scaffold
- **Frontend:** Next.js 14.2.28, Tailwind CSS, standalone output, dark theme
- **Backend:** FastAPI, pydantic-settings, hatchling, SQLAlchemy 2.0 ready
- **Database:** PostgreSQL via CloudNativePG
- **Helm:** Standard DClaw chart with deployment, service, ingress, HPA
- **Docs:** `docs/` directory with stubs
- **Mock API:** 2 endpoints (POST create, GET detail) with fake data

## Stack Constraints
| Layer | Technology | Constraint |
|-------|------------|------------|
| Frontend | Next.js 14 App Router | App Router only |
| Frontend | Tailwind CSS | Custom design tokens |
| Backend | FastAPI | pydantic v2, async routes, SQLAlchemy 2.0 |
| Backend | Python | 3.11+, type hints on all public APIs |
| Database | PostgreSQL | CloudNativePG for K8s, asyncpg for backend |
| AI | Ollama | Local-first. Cloud fallback via OpenRouter |
| PII | ClawShield | Must call Shield before any external API |
| Auth | Logto | JWT tokens, RBAC: Owner/Admin/Developer/User/Guest |
| Packaging | Helm 3 | Standard DClaw chart patterns |

## Phases

### Phase 1 — Market Intelligence & Gap Analysis
1. Identify top 5 competitors in {category}
2. Build comparative feature matrix
3. Identify gaps: table-stakes, differentiation, AI-native, integrations, UX
4. Research 2026-2027 trends

### Phase 2 — Feature Design (v1.0)
- **P0:** Must-have launch blockers
- **P1:** Differentiators
- **P2:** Nice-to-have post-launch

For each P0/P1 feature specify: user flow, UI/UX design, API contract, data model, AI prompt (if applicable), edge cases.

### Phase 3 — Architecture Design
- Mermaid system architecture diagram
- Complete SQLAlchemy models
- Full OpenAPI-style API spec
- 5 production-ready code snippets (service, AI integration, real-time, background job, auth)

### Phase 4 — Design System
Map brand color to DKube palette:
| Token | Value |
|-------|-------|
| `--dk-purple` | {color} |
| `--dk-purple-light` | {light} |
| `--dk-purple-deep` | {deep} |
| `--dk-purple-wash` | {wash} |

Typography: Manrope (display), Inter (body), JetBrains Mono (mono).

### Phase 5 — Testing Strategy
| Layer | Framework | Coverage |
|-------|-----------|----------|
| Unit | pytest | 80%+ |
| Integration | pytest + httpx | 70%+ |
| E2E | Playwright | Critical paths |
| Contract | schemathesis | 100% |
| Load | locust | 100 concurrent |
| Security | bandit + safety | All code |

### Phase 6 — Roadmap
- v1.0 release checklist
- v1.1 quick wins (4 weeks)
- v1.2 AI enhancements (8 weeks)
- v2.0 platform (Q1 2027)

### Phase 7 — Documentation Update
Fill docs with: installation, quickstart, configuration, use-cases, best-practices, architecture, stack, API, troubleshooting, changelog, roadmap.

## Output Format
Produce `v1.0-spec.md`:
```
# {name} v1.0 Specification
## Executive Summary
## Market Intelligence
## Feature Design
## Architecture
## Design System
## Testing Strategy
## Roadmap
## Appendix: Code Snippets
## Appendix: Database Schema
## Appendix: API Spec
```

This document must be comprehensive enough that a skilled developer can build v1.0 from it without asking questions.
"""


def generate_impl_prompt(app: dict, idx: int) -> str:
    app_id = app["id"]
    name = app["name"]
    category = app["category"]
    tagline = app["tagline"]
    color = app["color"]
    status = app["status"]
    icon = app["icon"]
    frontend_port = 3000 + idx
    backend_port = 8000 + idx
    db_name = f"dclaw_{app_id}"
    light, deep, wash = color_vars(color)

    return f"""# {name} v1.0 — Implementer Prompt

> **Agent Role:** Senior Full-Stack Engineer + QA Engineer + Technical Writer
> **Mission:** Build {name} from scaffold to working v1.0. Write code, tests, and docs. You do NOT write specs — you BUILD.
> **Output:** Fully functional app with passing tests and real docs.
> **Priority:** {status}

---

## Golden Rules
1. **Read Before You Write** — Examine existing files before modifying.
2. **Test-Driven** — Write test first, then implementation.
3. **Commit Per Feature** — Format: `feat(scope): description`.
4. **No Stubs** — Every function does real work. No `pass`, no `TODO`.
5. **Type Safety** — Every Python function typed. Every TS prop has interface.
6. **Error Handling** — Every async call has try/except. Proper HTTP codes.
7. **Security First** — PII through ClawShield. Auth validates JWT.
8. **Dark Mode Default** — All UI must look correct in dark mode.

## App Identity
```yaml
app_id: {app_id}
name: {name}
category: {category}
tagline: {tagline}
current_version: 0.1.0
target_version: 1.0.0
status: {status}
repo: https://github.com/dclawstack/dclaw-{app_id}
docs: https://docs.dclawstack.io/apps/{app_id}
frontend_port: {frontend_port}
backend_port: {backend_port}
db_name: {db_name}
brand_color: {color}
icon: {icon}
```

## Implementation Order

### Phase 0 — Setup & Audit
```bash
git clone https://github.com/dclawstack/dclaw-{app_id}.git
cd dclaw-{app_id}
git checkout -b v1.0-implementation
```
Read ALL existing files. Verify builds. Commit: `chore: audit scaffold`.

### Phase 1 — Database Layer
- Design schema for {category} domain
- SQLAlchemy models in `backend/app/models/`
  - Every model: `id` (UUID), `created_at`, `updated_at`, `deleted_at`
  - Use `Mapped[]` syntax, relationships, indexes
- Alembic migration: `alembic revision --autogenerate -m "v1.0 schema"`
- Model unit tests in `backend/tests/unit/test_models.py`
- Commit: `feat(db): add v1.0 schema`

### Phase 2 — Backend Core
- **Config:** `backend/app/core/config.py` with env vars
- **Dependencies:** `backend/app/core/deps.py` with DB session + auth
- **Exceptions:** Custom exceptions in `backend/app/core/exceptions.py`
- **Repositories:** `backend/app/repositories/` — CRUD per model, typed, handle None
- **Services:** `backend/app/services/` — business logic, docstrings, error handling, logging
- Commit: `feat(backend): add repositories and services`

### Phase 3 — Backend API
- **Auth middleware:** Verify JWT from Logto. Extract user + roles. 401/403 correctly.
- **Routers:** Replace mocks in `backend/app/api/v1/`
  - Prefix: `/api/v1/{app_id}`
  - Typed request/response models
  - Proper HTTP codes (201 create, 200/204 update, 404 not found)
  - `Depends()` for auth and DB
  - Try/except with custom handlers
- **AI Integration** (if applicable):
  - `backend/app/services/ai_service.py`
  - ClawShield → Ollama → OpenRouter fallback
  - Log token usage and latency
- **Background jobs** (if needed): Celery/ARQ with Redis
- **WebSocket/SSE** (if real-time): Authenticated, graceful disconnect
- Commit per router: `feat(api): add {{resource}} CRUD`

### Phase 4 — Frontend Foundation
- **Design tokens:** Update `tailwind.config.ts` with brand color `{color}`
- **Global styles:** CSS variables in `globals.css`, dark mode overrides
- **API client:** `frontend/src/lib/api.ts` — Axios with JWT interceptor, 401 redirect
- **Auth context:** `frontend/src/contexts/AuthContext.tsx`
- Commit: `feat(frontend): add tokens, api client, auth`

### Phase 5 — Frontend Pages & Components
- **Layout:** DKube fonts (Manrope, Inter, JetBrains Mono), dark mode, toast provider
- **Dashboard:** App shell with responsive sidebar
- **List views:** Data table with sort/filter/paginate/search, empty state, skeletons, bulk actions
- **Detail views:** Read-only + edit toggle, related data tabs, activity log
- **Forms:** React Hook Form + Zod, AI-assisted fields, auto-save drafts, confirmation modals
- **AI components** (if applicable): Chat bubble, streaming text, code block renderer, feedback buttons
- Commit per page: `feat(frontend): add {{page}}`

### Phase 6 — Testing
- **Unit:** `backend/tests/unit/` — pytest-asyncio, mock externals, 80%+ coverage
- **Integration:** `backend/tests/integration/` — httpx.AsyncClient, test every endpoint, auth, CRUD, errors
- **E2E:** `frontend/e2e/` — Playwright, critical flows, dark mode, responsive
- Run all tests. Fix failures. Do not proceed until green.
- Commit: `test: add unit, integration, and e2e tests`

### Phase 7 — Documentation
Replace all stubs with real content:
- **Getting Started:** installation, quickstart (5-min), configuration
- **Guides:** 3-5 use-cases, best-practices, ai-prompts (if AI app)
- **Reference:** architecture (Mermaid), stack (deps), API (OpenAPI render)
- **Troubleshooting:** 5-10 common issues, 10-15 FAQ
- **Releases:** changelog v0.1.0→v1.0.0, roadmap v1.1/v1.2
- Commit: `docs: write v1.0 documentation`

### Phase 8 — Helm & Deployment
- Update `helm/values.yaml`: image tags, resources (frontend 500m/512Mi, backend 1000m/1Gi), ingress host `{app_id}.dclawstack.io`, HPA min 2 max 10
- Verify: `helm lint helm/`
- Commit: `chore(helm): update for v1.0`

### Phase 9 — Final Validation
```bash
cd frontend && npm run build        # 0 errors
cd backend && ruff check app/       # 0 violations
cd backend && pytest                # 0 failures
cd frontend && npx playwright test  # 0 failures
```
- Bump version to 1.0.0 in package.json, pyproject.toml, docs/meta.json
- Final commit: `chore: bump version to 1.0.0`
- Push: `git push origin v1.0-implementation`
- Open PR on GitHub

## Code Quality Checklist (Before Every Commit)
- [ ] `ruff check app/` passes
- [ ] `mypy app/` passes
- [ ] `pytest` passes
- [ ] `npm run build` passes
- [ ] No `console.log` or `print()` (use logger)
- [ ] No hardcoded secrets
- [ ] No `TODO`/`FIXME` in production code
- [ ] Every route has auth middleware
- [ ] Every DB query uses parameterized statements
- [ ] Every AI call goes through ClawShield

## File Naming Conventions
| Type | Pattern | Example |
|------|---------|---------|
| Models | `app/models/{{entity}}.py` | `app/models/conversation.py` |
| Repositories | `app/repositories/{{entity}}_repo.py` | `app/repositories/conversation_repo.py` |
| Services | `app/services/{{feature}}_service.py` | `app/services/chat_service.py` |
| Routers | `app/api/v1/{{entity}}.py` | `app/api/v1/conversations.py` |
| Schemas | `app/schemas/{{entity}}.py` | `app/schemas/conversation.py` |
| Components | `src/components/{{Name}}.tsx` | `src/components/ChatBubble.tsx` |
| Hooks | `src/hooks/use{{Feature}}.ts` | `src/hooks/useChat.ts` |
| Pages | `src/app/{{route}}/page.tsx` | `src/app/chat/page.tsx` |
| Tests (BE) | `tests/{{layer}}/test_{{module}}.py` | `tests/unit/test_chat_service.py` |
| Tests (FE) | `e2e/{{feature}}.spec.ts` | `e2e/chat.spec.ts` |

## AI Integration Pattern (If Applicable)
```python
# backend/app/services/ai_service.py
import httpx
from app.core.config import settings
from app.core.shield import shield

async def generate_with_ai(prompt: str, user_id: str, model: str = "llama3.2") -> str:
    clean_prompt = await shield.scrub(prompt, user_id=user_id)
    try:
        async with httpx.AsyncClient(timeout=60.0) as client:
            resp = await client.post(
                f"{{settings.OLLAMA_URL}}/api/generate",
                json={{"model": model, "prompt": clean_prompt, "stream": False}},
            )
            resp.raise_for_status()
            return resp.json()["response"]
    except Exception:
        async with httpx.AsyncClient(timeout=60.0) as client:
            resp = await client.post(
                "https://openrouter.ai/api/v1/chat/completions",
                headers={"Authorization": f"Bearer {{settings.OPENROUTER_KEY}}"},
                json={{
                    "model": "anthropic/claude-3.5-sonnet",
                    "messages": [{{"role": "user", "content": clean_prompt}}],
                }},
            )
            resp.raise_for_status()
            return resp.json()["choices"][0]["message"]["content"]
```

## Auth Middleware Pattern
```python
from fastapi import Depends, HTTPException
from fastapi.security import HTTPBearer
import jwt

security = HTTPBearer()

async def get_current_user(credentials=Depends(security)):
    try:
        payload = jwt.decode(
            credentials.credentials,
            settings.LOGTO_JWKS,
            algorithms=["RS256"],
            audience=settings.LOGTO_AUDIENCE,
        )
        user_id = payload.get("sub")
        if not user_id:
            raise HTTPException(status_code=401, detail="Invalid token")
        return user
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")
```

## Final Checklist
- [ ] Real SQLAlchemy models + Alembic migrations
- [ ] Real FastAPI routers with full CRUD + auth
- [ ] Real business logic in services
- [ ] Real Next.js pages with working UI
- [ ] Passing pytest (unit + integration)
- [ ] Passing Playwright E2E
- [ ] Real documentation in `docs/`
- [ ] Updated Helm chart
- [ ] Zero lint errors, zero type errors
- [ ] PR on GitHub ready for review

This is working code, not a spec.
"""


def main():
    base = Path(__file__).parent
    prompts_dir = base / "prompts"
    prompts_dir.mkdir(exist_ok=True)

    manifest = []
    dispatch = []

    for idx, app in enumerate(APPS, start=1):
        app_id = app["id"]

        spec_text = generate_spec_prompt(app, idx)
        impl_text = generate_impl_prompt(app, idx)

        spec_path = prompts_dir / f"{app_id}-v1.0-spec-prompt.md"
        impl_path = prompts_dir / f"{app_id}-v1.0-impl-prompt.md"

        spec_path.write_text(spec_text, encoding="utf-8")
        impl_path.write_text(impl_text, encoding="utf-8")

        manifest.append({
            "app_id": app_id,
            "name": app["name"],
            "category": app["category"],
            "tagline": app["tagline"],
            "color": app["color"],
            "status": app["status"],
            "priority": app["priority"],
            "icon": app["icon"],
            "repo_url": f"https://github.com/dclawstack/dclaw-{app_id}",
            "spec_prompt": str(spec_path.relative_to(base)),
            "impl_prompt": str(impl_path.relative_to(base)),
            "frontend_port": 3000 + idx,
            "backend_port": 8000 + idx,
            "db_name": f"dclaw_{app_id}",
        })

        dispatch.append({
            "app_id": app_id,
            "name": app["name"],
            "priority": app["priority"],
            "repo_url": f"https://github.com/dclawstack/dclaw-{app_id}",
            "spec_prompt": str(spec_path.relative_to(base)),
            "impl_prompt": str(impl_path.relative_to(base)),
            "status": app["status"],
        })

        print(f"✅ {app_id}")

    manifest_path = base / "apps-manifest.json"
    manifest_path.write_text(json.dumps(manifest, indent=2), encoding="utf-8")

    dispatch_path = base / "dispatch.csv"
    with dispatch_path.open("w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=["app_id", "name", "priority", "repo_url", "spec_prompt", "impl_prompt", "status"])
        writer.writeheader()
        writer.writerows(dispatch)

    print(f"\nDone. {len(APPS)} apps × 2 prompts = {len(APPS)*2} files.")
    print(f"Manifest: {manifest_path}")
    print(f"Dispatch: {dispatch_path}")


if __name__ == "__main__":
    main()
