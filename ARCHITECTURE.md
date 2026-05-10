# DClaw Stack — Architecture & Operating Model

> **The AI-First Enterprise Suite: How 72 Vertical SaaS Apps Are Built in Parallel by a 20-Person Team**
>
> *"We don't write code. We orchestrate code."*

---

## Table of Contents

1. [Executive Manifesto](#1-executive-manifesto)
2. [The DClaw Stack at a Glance](#2-the-dclaw-stack-at-a-glance)
3. [Core Architecture](#3-core-architecture)
4. [The Scaffold Philosophy](#4-the-scaffold-philosophy)
5. [One Day in the Life](#5-one-day-in-the-life)
6. [Team Isolation Model](#6-team-isolation-model)
7. [The Build Process](#7-the-build-process)
8. [AGENTS.md Governance](#8-agentsmd-governance)
9. [CI/CD & Automation](#9-cicd--automation)
10. [Scaling to Enterprise](#10-scaling-to-enterprise)
11. [Appendix](#11-appendix)

---

## 1. Executive Manifesto

### The Thesis

Traditional enterprise software takes years to build, hundreds of engineers, and billions in capital. The DClaw Stack inverts this. We prove that a small, AI-first team can scaffold, build, and deploy an entire enterprise software suite — 72 vertical SaaS applications spanning CRM, HR, Finance, DevOps, AI/ML, Content, Compliance, Healthcare, and more — in a matter of weeks, not years.

### The Method

**Velocity through Constraint.** We don't argue about frameworks. We don't debate folder structures. We don't bikeshed on naming conventions. Every decision is made once, encoded in a scaffold, and replicated 72 times. The scaffold *is* the architecture. The architecture *is* the contract.

**AI as Team Member.** Every repository has Claude Code Action integrated. Every app has a domain-specific `PLAN-v1.2.md` that serves as the product spec. Developers don't start from blank screens — they start from a scaffold with working tests, a running database, pre-built UI components, and an AI copilot that knows the codebase.

**Parallel Independence.** Each of the 72 apps lives in its own repository with its own CI/CD, its own database, its own deployment pipeline. A developer working on the CRM app never conflicts with a developer working on the Inventory app. Merge conflicts are isolated to single-app boundaries.

**Enterprise from Day One.** Security, observability, testing, migrations, and deployment are not afterthoughts. They are baked into the scaffold. Every app ships with passing tests, Alembic migrations, Docker containers, Helm charts, and GitHub Actions workflows from commit zero.

### The Result

A 20-person team, in a single day, can:
- Assign 72 repositories to individual owners
- Scaffold all 72 apps with working backends, frontends, databases, and CI/CD
- Generate domain-specific AI copilot features for each vertical
- Deploy the entire suite to a Kubernetes cluster
- All while working remotely, asynchronously, and — yes — enjoying their holiday.

This is not science fiction. This is the DClaw Stack.

---

## 2. The DClaw Stack at a Glance

### The Numbers

| Metric | Value |
|--------|-------|
| **Total Apps** | 72 vertical SaaS applications |
| **Vertical Groups** | 15 industry domains |
| **Repositories** | 72 independent Git repos + 4 infrastructure repos |
| **Backend Ports** | 8095–8155 (61 ports assigned) |
| **Frontend Ports** | 3006–3068 (63 ports assigned) |
| **Shared Scaffold** | 1 canonical scaffold (`dclaw-scaffold`) |
| **CI/CD Pipelines** | ~58 with pytest + build, ~15 with Docker + Helm deploy |
| **AI Integrations** | Claude Code Action in ~25 repos, domain-specific AI copilots in all 72 |

### The 15 Vertical Groups

| # | Group | Apps | Representative |
|---|-------|------|--------------|
| 1 | **CRM & Sales** | crm, marketing, sales, email, seo | Customer lifecycle management |
| 2 | **HR & People** | hr, onboard, recruit, offboard, train | Talent pipeline & culture |
| 3 | **Finance & Legal** | contract, finance, legal, patent, trademark, cost | Financial operations & IP |
| 4 | **Operations & Logistics** | inventory, lease, fleet, route, warehouse | Supply chain & assets |
| 5 | **DevOps & IT** | code, maintenance, api, deploy, monitor, secure, test, backup, migrate | Infrastructure lifecycle |
| 6 | **AI & Data Intelligence** | agent, rag, research, data, knowledge, forecast | AI/ML operations |
| 7 | **Content & Media** | doc, video, slide, sheet, wiki, write, create, translate | Content creation suite |
| 8 | **Communication & Collab** | chat, meet, calendar, task | Team collaboration |
| 9 | **Compliance & Risk** | audit, compliance, policy, risk, carbon | Governance & ESG |
| 10 | **Healthcare** | med | Clinical records & scheduling |
| 11 | **Real Estate & Construction** | real-estate, building, space | Property management |
| 12 | **Utilities & Environment** | energy, water, waste | Sustainability operations |
| 13 | **Crisis & Continuity** | crisis, continuity | Business resilience |
| 14 | **Creative & Design** | design, flow | Visual design & workflows |
| 15 | **Project & Support** | project, support | Delivery & customer success |

### The Infrastructure Layer

| Repo | Purpose |
|------|---------|
| `dclaw-core` | Shared libraries, SDKs, client packages |
| `dclaw-platform` | Central coordination, port registry, docs, operator |
| `dclaw-panel` | DPanel — the unified admin dashboard |
| `dclaw-prd-tmp` | Product requirements & design system |
| `dclaw-scaffold` | The canonical scaffold — source of truth for all apps |

---

## 3. Core Architecture

### 3.1 Backend: FastAPI + SQLAlchemy 2.0

Every DClaw app shares an identical backend architecture:

```
backend/
├── app/
│   ├── api/
│   │   ├── main.py              # FastAPI app with lifespan handler
│   │   ├── routes/health.py     # Health check endpoint
│   │   └── v1/                  # App-specific routers
│   ├── core/
│   │   ├── config.py            # Pydantic Settings
│   │   └── database.py          # create_async_engine + get_db
│   ├── models/
│   │   ├── base.py              # DeclarativeBase (NOT MappedAsDataclass)
│   │   └── *.py                 # App-specific models
│   ├── repositories/
│   │   ├── base_repo.py         # Generic BaseRepository[T]
│   │   └── *_repo.py            # Entity-specific repositories
│   ├── schemas/
│   │   └── *.py                 # Pydantic v2 with ConfigDict
│   ├── services/
│   │   └── *.py                 # Business logic / AI integrations
│   └── utils/
│       └── *.py                 # utc_now(), helpers
├── alembic/                     # Async migrations
├── tests/
│   ├── conftest.py              # Test DB override, AsyncClient fixture
│   └── test_*.py                # pytest-asyncio tests
└── Dockerfile                   # python:3.11-slim, non-root, stdlib healthcheck
```

**Key Patterns:**

- **`DeclarativeBase`** — Plain SQLAlchemy 2.0 base class. NO `MappedAsDataclass`. NO `declarative_base()`.
- **Generic Repository** — `BaseRepository[T]` provides `list_all`, `get_by_id`, `create`, `delete`, `count`. Subclass per entity for custom queries.
- **Dependency Injection** — Every endpoint uses `Depends(get_db)`. Never manual `AsyncSession`.
- **Pydantic v2** — All schemas use `ConfigDict(from_attributes=True)`. No `orm_mode`.
- **Async Everything** — `create_async_engine` + `AsyncSession` + `async/await` from HTTP layer to database.

**Example — Model:**

```python
from uuid import UUID, uuid4
from datetime import datetime
from sqlalchemy import String, Enum, Text, func
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.models.base import Base

class Customer(Base):
    __tablename__ = "customers"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4, init=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    email: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    status: Mapped[str] = mapped_column(
        Enum("lead", "active", "churned", name="customer_status"),
        default="lead",
    )
    created_at: Mapped[datetime] = mapped_column(server_default=func.now(), init=False)

    deals: Mapped[list["Deal"]] = relationship(
        back_populates="customer", lazy="selectin", cascade="all, delete-orphan", init=False
    )
```

**Example — Repository:**

```python
from app.models.customer import Customer
from app.repositories.base_repo import BaseRepository

class CustomerRepository(BaseRepository[Customer]):
    def __init__(self, db):
        super().__init__(db, Customer)

    async def get_by_email(self, email: str) -> Customer | None:
        result = await self.db.execute(
            select(Customer).where(Customer.email == email)
        )
        return result.scalar_one_or_none()
```

**Example — API Endpoint:**

```python
@router.get("/", response_model=CustomerListResponse)
async def list_customers(
    limit: int = Query(20, ge=1, le=100),
    offset: int = Query(0, ge=0),
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    items, total = await repo.list_all(limit=limit, offset=offset)
    return CustomerListResponse(items=items, total=total)
```

### 3.2 Frontend: Next.js 14 + Tailwind CSS

Every DClaw app shares an identical frontend architecture:

```
frontend/
├── src/
│   ├── app/                     # Next.js 14 App Router
│   │   ├── layout.tsx           # Root layout with providers
│   │   ├── page.tsx             # Dashboard / home
│   │   └── [entity]/            # CRUD pages
│   ├── components/
│   │   └── ui/                  # Pre-built UI components
│   │       ├── button.tsx
│   │       ├── card.tsx
│   │       ├── input.tsx
│   │       ├── label.tsx
│   │       ├── badge.tsx
│   │       ├── select.tsx
│   │       ├── dialog.tsx
│   │       ├── table.tsx
│   │       ├── tabs.tsx
│   │       └── avatar.tsx
│   └── lib/
│       ├── api.ts               # Typed fetch wrapper
│       └── utils.ts             # cn() helper
├── package.json                 # Next.js 14 + Tailwind v3 + tailwindcss-animate
├── next.config.js               # Standalone output for Docker
└── Dockerfile                   # node:20-alpine, ARG NEXT_PUBLIC_API_URL
```

**Key Patterns:**

- **App Router** — Server components by default. Client components only for interactivity.
- **Pre-Built UI Components** — 10 shadcn-compatible components in every repo. NO shadcn CLI v4. NO `@base-ui/react`.
- **Typed API Client** — `src/lib/api.ts` provides a thin `fetchJson<T>()` wrapper with error handling.
- **Environment Variables** — `NEXT_PUBLIC_API_URL` baked at Docker build time. Dockerfile MUST declare `ARG` before `RUN npm run build`.
- **Tailwind v3** — CSS variables for theming. `tailwindcss-animate` in dependencies (not devDependencies).

**Example — CRUD Page:**

```tsx
"use client";
import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { listCustomers, createCustomer } from "@/lib/api";

export default function CustomersPage() {
  const [customers, setCustomers] = useState([]);
  const [open, setOpen] = useState(false);
  // ... load, create, render
}
```

A full CRUD list page — search, create dialog, sortable table, status badges — is ~150 lines of TypeScript. This is the power of pre-built components and a typed API client.

### 3.3 Database: PostgreSQL + Alembic

- **PostgreSQL 16** with asyncpg driver
- **Every app has its own database** — complete isolation, no shared schemas
- **Alembic migrations** — `alembic revision --autogenerate` for every model change
- **Test database** — `dclaw_app_test` on `localhost:5432`, overridden via `DATABASE_URL` env var

### 3.4 Docker & Deployment

- **Backend Dockerfile:** `python:3.11-slim`, non-root `appuser`, healthcheck via `urllib.request.urlopen()`
- **Frontend Dockerfile:** `node:20-alpine`, multi-stage, `ARG NEXT_PUBLIC_API_URL` before build
- **docker-compose.yml:** Per-app compose with PostgreSQL service + healthchecks
- **Helm Charts:** Kubernetes deployment + service + secret templates for every app
- **Container Registry:** GHCR (`ghcr.io/dclaw-stack/<repo>`)
- **Deploy Target:** Kubernetes via Helm with `KUBECONFIG` secret

---

## 4. The Scaffold Philosophy

### Why One Scaffold for 72 Apps?

The `dclaw-scaffold` repository is the single source of truth for how every DClaw app is structured. It is not a template you copy once and forget. It is a living contract that evolves and propagates.

**The Propagation Model:**

```
dclaw-scaffold/         ← Source of truth
    ├── backend/        ← Standard FastAPI structure
    ├── frontend/       ← Standard Next.js structure
    ├── docker-compose.yml
    ├── helm/
    ├── .github/workflows/ci.yml
    └── AGENTS.md       ← The constitution

        ↓ rsync / surgical merge

dclaw-crm/              ← Full rsync + domain models
    ├── PLAN-v1.2.md    ← "AI CRM Copilot + Pipeline + Contacts + Deals"
    └── ...

dclaw-finance/          ← Full rsync + domain models
    ├── PLAN-v1.2.md    ← "AI Finance Copilot + Ledger + Invoicing + Reporting"
    └── ...
```

### The Three Merge Rules

When propagating scaffold changes to existing apps, we use three strategies:

1. **Skeleton repos** — Full `rsync`. These are empty repos that receive the entire scaffold.
2. **Partial/Substantial repos** — Surgical merge. Preserve `app/models/`, `app/api/`, `app/services/`. Update infrastructure files (Docker, CI, package.json).
3. **Frontend-only / Special repos** — Add missing backend scaffold around existing code. Never move existing code.

### The Velocity Multiplier

| Task | Traditional | DClaw Scaffold |
|------|-------------|----------------|
| New app scaffolding | 2–3 days | 5 minutes |
| CRUD endpoint (backend) | 2 hours | 15 minutes |
| CRUD page (frontend) | 3 hours | 20 minutes |
| Database migration | 30 minutes | 2 minutes |
| CI/CD pipeline | 1 day | Already present |
| Docker setup | 2 hours | Already present |
| Tests | 2 hours | Already present |
| **Total per app** | **~5 days** | **~1 hour** |

With 72 apps, the scaffold saves approximately **350 developer-days** of setup work.

---

## 5. One Day in the Life

### 9:00 AM — The Standup (30 minutes)

Twenty team members gather on a video call. Not in an office. From beaches, mountains, and living rooms. It's a holiday week, but the team treats building as recreation.

**The Lead Architect** shares the screen:

> "Today we scaffold the remaining 44 apps. Each of you gets 2–4 repositories. Your `AGENTS.md` has your ports, your stack, and your anti-patterns. Your `PLAN-v1.2.md` has your feature roadmap. Your job: get the scaffold running, commit, push, and start on P0 features."

The team lead opens the **Port Registry**:

| Developer | Apps | Backend Ports | Frontend Ports |
|-----------|------|---------------|----------------|
| Alice | crm, marketing | 8095, 8102 | 3006, 3015 |
| Bob | finance, legal | 8096, 8099 | 3007, 3013 |
| Charlie | hr, recruit, onboard | 8097, 8105, 8097 | 3008, 3018, 3008 |
| ... | ... | ... | ... |

### 9:30 AM — Repo Assignment (15 minutes)

Each developer:
1. Opens their assigned repos on GitHub
2. Verifies they have write access
3. Checks that `AGENTS.md` has the correct ports
4. Reviews their `PLAN-v1.2.md` to understand the vertical domain

### 9:45 AM — The Great Scaffold (2 hours)

Three parallel agents per app:

**Agent 1: Backend Architect**
```bash
cd dclaw-crm
# Scaffold already present from previous session
# Task: Add domain models
python backend/app/models/customer.py     # 31 lines
python backend/app/models/deal.py         # 28 lines
python backend/app/models/activity.py     # 25 lines
alembic revision --autogenerate -m "add crm models"
```

**Agent 2: Frontend Builder**
```bash
cd dclaw-crm/frontend
# Pre-built UI components already present
# Task: Build CRUD pages
src/app/customers/page.tsx     # 151 lines — full CRUD
src/app/deals/page.tsx         # 148 lines — full CRUD
src/app/activities/page.tsx    # 95 lines — activity feed
npm run build  # Must pass
```

**Agent 3: DevOps Engineer**
```bash
cd dclaw-crm
# CI/CD already present from scaffold
# Task: Verify Docker compose
docker compose up -d
# Test health endpoint
curl http://localhost:8095/health/
# Verify frontend builds
cd frontend && npm run build
```

### 12:00 PM — Lunch & AI Review

Developers take a break. Meanwhile, Claude Code Action is active in every repository:

- **Code Review:** `claude-code-review.yml` runs on every PR, suggesting improvements
- **Feature Implementation:** Developers can `@claude` in issues: "Implement the AI CRM Copilot from PLAN-v1.2.md P0.1"
- **Bug Fixes:** Claude scans for anti-patterns from `AGENTS.md` and auto-suggests fixes

### 1:00 PM — Feature Implementation (3 hours)

Each developer implements P0 features from their `PLAN-v1.2.md`:

**Alice (CRM):**
- P0.1: AI CRM Copilot — `/api/v1/ai/crm-chat` endpoint + chat sidebar
- P0.2: Pipeline Management — deal stage transitions
- P0.3: Contact Management — customer CRUD with search
- P0.4: Activity Timeline — chronological feed

**Bob (Finance):**
- P0.1: AI Finance Copilot — natural language to ledger queries
- P0.2: Double-Entry Ledger — journal entries with validation
- P0.3: Invoice Generation — PDF export + email
- P0.4: Financial Reporting — P&L, balance sheet, cash flow

Every feature follows the same pattern:
1. Read `AGENTS.md` — verify no anti-patterns
2. Read `PLAN-v1.2.md` — understand the feature spec
3. Backend: Model → Schema → Repository → Router → Tests → Migration
4. Frontend: API types → Page → Components
5. Docker: `docker compose up -d` to verify
6. CI: Push → GitHub Actions runs pytest + build

### 4:00 PM — Integration Testing (1 hour)

The team spins up the full stack:

```bash
cd dclaw-platform
docker compose up -d  # Starts PostgreSQL + 8 app pairs
```

Each developer verifies:
- Their backend responds on the assigned port
- Their frontend builds and connects to the backend
- Database migrations ran successfully
- Tests pass in CI

### 5:00 PM — The Push

```bash
# Each developer, in their repos
git add .
git commit -m "feat: implement P0 features per PLAN-v1.2.md

- AI Copilot with domain-specific endpoints
- Core CRUD with repository pattern
- Full test coverage
- Docker + CI verified"
git push origin main
```

GitHub Actions fire across 72 repositories. 58 CI pipelines run pytest. 15 build Docker images. 15 deploy to staging Kubernetes. Slack notifications ping the team channel.

### 6:00 PM — Done

The team logs off. The CI/CD pipelines finish. 72 apps are scaffolded, tested, and deployed. P0 features are implemented. The suite is live.

Tomorrow: P1 features, integration between apps, and the AI copilot training.

**Total time: 6 hours. 72 apps. One team. Zero merge conflicts.**

---

## 6. Team Isolation Model

### One Repo Per App Per Developer

The DClaw Stack's most important architectural decision is **complete repository isolation**. Every app is a fully independent repository with:

- Its own Git history
- Its own CI/CD pipeline
- Its own Docker containers
- Its own PostgreSQL database
- Its own Kubernetes namespace
- Its own feature roadmap (`PLAN-v1.2.md`)
- Its own governance (`AGENTS.md`)

### Why This Matters

| Problem | Monorepo Solution | DClaw Solution |
|---------|-------------------|----------------|
| Merge conflicts | Complex branching strategy | No conflicts — different repos |
| Build times | Incremental builds, caching | Parallel builds — 72 pipelines |
| Deployment risk | Feature flags, canary deploys | Independent deploys — one app at a time |
| Team scaling | Code ownership files | Natural ownership — one repo per team |
| Tech debt spread | Lint rules, code reviews | Isolated — debt in one repo doesn't infect others |
| Rollback | Complex revert strategies | `git revert` + `helm rollback` per app |

### The Cost

The trade-off is **cross-app integration**. When the CRM app needs to reference customer data from the Finance app, we use:

1. **API Contracts** — OpenAPI specs published by each app
2. **Async Events** — Message queue for cross-app notifications (future)
3. **Shared Core Library** — `dclaw-core` for common utilities, auth, and client SDKs
4. **DPanel** — The unified dashboard that aggregates data from all apps

---

## 7. The Build Process

### From PLAN-v1.2.md to Production Code

Every feature in every app follows the same 6-step process:

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  PLAN-v1.2.md   │────▶│  AGENTS.md      │────▶│  Scaffold       │
│  (The Spec)     │     │  (The Rules)    │     │  (The Foundation│
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                                               │
         │                                               ▼
         │                                      ┌─────────────────┐
         │                                      │  Backend        │
         │                                      │  Model          │
         │                                      │  Schema         │
         │                                      │  Repository     │
         │                                      │  Router         │
         │                                      │  Tests          │
         │                                      │  Migration      │
         │                                      └─────────────────┘
         │                                               │
         │                                               ▼
         │                                      ┌─────────────────┐
         │                                      │  Frontend       │
         │                                      │  API Types      │
         │                                      │  Page           │
         │                                      │  Components     │
         │                                      └─────────────────┘
         │                                               │
         │                                               ▼
         │                                      ┌─────────────────┐
         └────────────────────────────────────▶│  CI/CD          │
                                               │  pytest         │
                                               │  npm build      │
                                               │  Docker         │
                                               │  Helm Deploy    │
                                               └─────────────────┘
```

### Step-by-Step

**Step 1: Read the Spec**
```bash
cat PLAN-v1.2.md | grep "P0"  # Find your features
cat AGENTS.md | grep "Anti-Patterns"  # Know what not to do
```

**Step 2: Backend — Add Model**
```python
# backend/app/models/customer.py
class Customer(Base):
    __tablename__ = "customers"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4, init=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
```

**Step 3: Backend — Add Schema**
```python
# backend/app/schemas/customer.py
class CustomerCreate(BaseModel):
    name: str
    email: str
```

**Step 4: Backend — Add Repository**
```python
# backend/app/repositories/customer_repo.py
class CustomerRepository(BaseRepository[Customer]):
    def __init__(self, db):
        super().__init__(db, Customer)
```

**Step 5: Backend — Add Router**
```python
# backend/app/api/v1/customers.py
@router.post("/", response_model=CustomerResponse, status_code=201)
async def create_customer(data: CustomerCreate, db: AsyncSession = Depends(get_db)):
    repo = CustomerRepository(db)
    return await repo.create(Customer(**data.model_dump()))
```

**Step 6: Backend — Add Tests**
```python
# backend/tests/test_customers.py
@pytest.mark.asyncio
async def test_create_customer(client: AsyncClient):
    response = await client.post("/api/v1/customers/", json={"name": "Alice", "email": "alice@example.com"})
    assert response.status_code == 201
```

**Step 7: Backend — Generate Migration**
```bash
cd backend
alembic revision --autogenerate -m "add customers table"
```

**Step 8: Frontend — Add API Types**
```typescript
// frontend/src/lib/api.ts
export async function createCustomer(data: CustomerCreate) {
  return fetchJson<CustomerResponse>("/api/v1/customers/", {
    method: "POST",
    body: JSON.stringify(data),
  });
}
```

**Step 9: Frontend — Build Page**
```tsx
// frontend/src/app/customers/page.tsx
export default function CustomersPage() {
  // Use pre-built components: Card, Table, Dialog, Button, Input
}
```

**Step 10: Verify**
```bash
docker compose up -d
curl http://localhost:8095/health/
cd frontend && npm run build
cd backend && pytest
```

**Step 11: Commit & Push**
```bash
git add .
git commit -m "feat: add customer CRUD per PLAN-v1.2.md P0.3"
git push origin main
```

**Total time per feature: 20–40 minutes.**

---

## 8. AGENTS.md Governance

### The Constitution

Every repository contains an `AGENTS.md` file. It is not documentation. It is **law**. AI agents (Claude Code, GitHub Copilot,Cursor) and human developers must read it before making any code changes.

### Architecture Locks (Non-Negotiable)

These are encoded as absolute rules. If an AI suggests changing them, the developer rejects it:

| Layer | Lock | Violation |
|-------|------|-----------|
| **Backend** | FastAPI with `lifespan` handler | Using Flask, Django, or sync handlers |
| **Backend** | SQLAlchemy 2.0 `DeclarativeBase` | Using `MappedAsDataclass` or `declarative_base()` |
| **Backend** | Pydantic v2 with `ConfigDict` | Using Pydantic v1 `orm_mode` |
| **Backend** | Async SQLAlchemy (`create_async_engine`) | Using sync `create_engine` |
| **Backend** | Repository pattern in `app/repositories/` | Inline DB access in routers |
| **Backend** | `Depends(get_db)` injection | Manual `AsyncSession` with `__anext__()` |
| **Backend** | NO mock data (in-memory dicts) | Using `MOCK_CUSTOMERS = {}` |
| **Frontend** | Next.js 14+ App Router | Using Pages Router |
| **Frontend** | Tailwind CSS v3 + pre-built components | Installing shadcn CLI v4 or `@base-ui/react` |
| **Frontend** | `NEXT_PUBLIC_API_URL` baked at build | Hardcoding `localhost:PORT` |
| **Docker** | `python:3.11-slim` with non-root `appuser` | Running as root |
| **Docker** | Python stdlib healthcheck | Using `curl` (not in slim image) |
| **Tests** | `pytest-asyncio==0.24.0` pinned | Upgrading to newer versions |
| **CI/CD** | `.github/workflows/ci.yml` present | Deleting CI workflow |

### Anti-Patterns Table

The most critical section of `AGENTS.md` is the anti-patterns table. It prevents the most common mistakes that break the scaffold:

| Anti-Pattern | Why It Breaks | Correct Alternative |
|--------------|---------------|---------------------|
| `declarative_base()` in `database.py` | Separate metadata → zero tables | `from app.models.base import Base` |
| `curl` in healthcheck on `python:*-slim` | No `curl` binary → silent failure | `python -c "import urllib.request; urllib.request.urlopen(...)"` |
| In-memory `MOCK_*` dicts | Data lost on restart | Create repository + real DB |
| Missing `ARG NEXT_PUBLIC_API_URL` | Wrong API URL baked into container | Add `ARG` before `RUN npm run build` |
| Manual `get_db()` with `__anext__()` | Session leaks, connection pool exhaustion | `Depends(get_db)` |
| Hardcoded `localhost:PORT` | Breaks Docker/K8s networking | Use `process.env.NEXT_PUBLIC_API_URL` |
| No alembic migration for new models | Schema drift, production crashes | `alembic revision --autogenerate` |
| Installing `shadcn` CLI v4 | Breaks Tailwind v3 build | Use pre-built components in scaffold |
| Using `@base-ui/react` | Incompatible with Tailwind v3 | Use pre-built components in scaffold |
| Non-standard Postgres port in tests | CI service only maps 5432 | Always use `localhost:5432` in `conftest.py` |
| Upgrading `pytest-asyncio` | v1.3.0 breaks fixture scoping | Keep `pytest-asyncio==0.24.0` pinned |
| Deleting `.github/workflows/ci.yml` | No CI runs, no quality gate | Leave CI workflow intact |
| Missing `src/lib/utils.ts` | UI components fail to import `cn()` | Already in scaffold — do NOT delete |
| `MappedAsDataclass` in `Base` | Relationship/foreign-key sync conflicts | Use plain `DeclarativeBase` only |
| `default_factory` in `mapped_column()` | `ArgumentError` on plain `DeclarativeBase` | Use `default=` with callable |
| Timezone-aware `datetime` in models | `DataError` with `TIMESTAMP WITHOUT TIME ZONE` | Use naive UTC via `utc_now()` |

### Database Rules

1. All models MUST inherit from `Base` in `app.models.base`
2. All models MUST use `Mapped[...]` and `mapped_column()`
3. **Never use `default_factory=` in `mapped_column()`** — use `default=` instead
4. Relationships MUST specify `lazy="selectin"`
5. All new tables MUST get an alembic migration
6. Use `ondelete="CASCADE"` for child tables
7. Use `ondelete="SET NULL"` for optional references

---

## 9. CI/CD & Automation

### GitHub Actions Ecosystem

Every app repository has a CI/CD pipeline that runs on every push and PR:

**Standard CI (`ci.yml`):**
```yaml
jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: dclaw_app_test
        ports: ["5432:5432"]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.11" }
      - run: pip install uv && uv pip install -r backend/requirements.txt
      - run: cd backend && pytest
        env:
          DATABASE_URL: postgresql+asyncpg://postgres:postgres@localhost:5432/dclaw_app_test

  frontend-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: "20" }
      - run: cd frontend && npm ci && npm run build
        env:
          NEXT_PUBLIC_API_URL: http://localhost:8000
```

**Docker Build (`build-backend.yml` / `build-frontend.yml`):**
```yaml
- uses: docker/build-push-action@v5
  with:
    context: ./backend
    push: true
    platforms: linux/amd64,linux/arm64
    tags: ghcr.io/dclaw-stack/${{ github.repository }}-backend:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Helm Deploy (`deploy.yml`):**
```yaml
- run: |
    helm upgrade --install ${{ github.repository }} ./helm \
      --namespace dclaw \
      --set backend.image.tag=${{ github.sha }} \
      --set frontend.image.tag=${{ github.sha }} \
      --wait --timeout 5m
```

### Claude Code Action

~25 repositories have Claude Code Action integrated for AI-assisted development:

**OAuth Variant (read-only):**
```yaml
- uses: anthropics/claude-code-action@v1
  with:
    claude_code_oauth_token: ${{ secrets.CLAUDE_CODE_OAUTH_TOKEN }}
    additional_permissions: '{"actions": "read"}'
```

**API Key Variant (write):**
```yaml
- uses: anthropics/claude-code-action@v1
  with:
    anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
    claude_args: --bare --model claude-opus-4-7 --max-turns 10
```

**Auto Code Review (`claude-code-review.yml`):**
```yaml
- uses: anthropics/claude-code-action@v1
  with:
    plugin: code-review@claude-code-plugins
    prompt: /code-review:code-review ${{ github.repository }}/pull/${{ github.event.pull_request.number }}
```

### The Automation Flywheel

1. **Developer pushes code** → CI runs pytest + build
2. **CI passes** → Docker images build and push to GHCR
3. **Manual trigger** → Helm deploys to Kubernetes
4. **Claude reviews** → AI suggests improvements on PRs
5. **Claude implements** → `@claude` in issues triggers feature implementation from `PLAN-v1.2.md`

---

## 10. Scaling to Enterprise

### From 72 Independent Apps to an Integrated Suite

The DClaw Stack is designed to scale from "72 independent apps" to "one unified enterprise platform" through four integration layers:

**Layer 1: Shared Authentication (Now)**
- OAuth2 / OIDC via `dclaw-core` auth library
- JWT tokens with cross-app validation
- Role-based access control (RBAC) per app

**Layer 2: DPanel — The Unified Dashboard (In Progress)**
- Single pane of glass for all 72 apps
- Cross-app analytics and reporting
- User provisioning and app discovery
- Built in `dclaw-panel` with Next.js + FastAPI

**Layer 3: API Mesh & Event Bus (Future)**
- Async message queue (RabbitMQ / NATS) for cross-app events
- GraphQL federation layer for unified queries
- OpenAPI spec registry in `dclaw-platform`

**Layer 4: AI Orchestrator (Future)**
- Central AI coordinator that routes queries to the right app's copilot
- Cross-app RAG — "How much did we spend on the Acme deal?" queries CRM + Finance
- Multi-agent workflow orchestration

### Enterprise Non-Functional Requirements

Every app in the DClaw Stack ships with:

| Requirement | Implementation |
|-------------|---------------|
| **Observability** | Structured logging + OpenTelemetry traces (future) |
| **Security** | Non-root containers, secrets via K8s, no hardcoded credentials |
| **Scalability** | Stateless backends, horizontal pod autoscaling (HPA) |
| **Reliability** | Healthchecks, readiness probes, graceful shutdown |
| **Compliance** | Audit logging, data retention policies, GDPR-ready deletion |
| **Testing** | pytest with >80% coverage target, frontend build gates |
| **Documentation** | `AGENTS.md` + `PLAN-v1.2.md` in every repo |

---

## 11. Appendix

### A. Port Registry (Excerpt)

| App | Backend | Frontend | Postgres DB |
|-----|---------|----------|-------------|
| dclaw-crm | 8095 | 3006 | dclaw_crm |
| dclaw-finance | 8096 | 3007 | dclaw_finance |
| dclaw-hr | 8097 | 3008 | dclaw_hr |
| dclaw-inventory | 8098 | 3009 | dclaw_inventory |
| dclaw-legal | 8099 | 3013 | dclaw_legal |
| dclaw-project | 8100 | 3010 | dclaw_project |
| dclaw-support | 8101 | 3014 | dclaw_support |
| dclaw-marketing | 8102 | 3015 | dclaw_marketing |
| ... | ... | ... | ... |

> **Full registry:** See `dclaw-platform/PORT_REGISTRY.md` (247 lines, 65+ apps)

### B. Scaffold Directory Structure

```
dclaw-scaffold/
├── backend/
│   ├── app/
│   │   ├── api/
│   │   │   ├── main.py
│   │   │   ├── routes/
│   │   │   │   └── health.py
│   │   │   └── v1/
│   │   ├── core/
│   │   │   ├── config.py
│   │   │   └── database.py
│   │   ├── models/
│   │   │   └── base.py
│   │   ├── repositories/
│   │   │   └── base_repo.py
│   │   ├── schemas/
│   │   ├── services/
│   │   └── utils/
│   ├── alembic/
│   ├── tests/
│   │   ├── conftest.py
│   │   └── test_health.py
│   ├── Dockerfile
│   └── requirements.txt
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   └── globals.css
│   │   ├── components/
│   │   │   └── ui/
│   │   │       ├── button.tsx
│   │   │       ├── card.tsx
│   │   │       ├── input.tsx
│   │   │       ├── label.tsx
│   │   │       ├── badge.tsx
│   │   │       ├── select.tsx
│   │   │       ├── dialog.tsx
│   │   │       ├── table.tsx
│   │   │       ├── tabs.tsx
│   │   │       └── avatar.tsx
│   │   └── lib/
│   │       ├── api.ts
│   │       └── utils.ts
│   ├── package.json
│   ├── next.config.js
│   └── Dockerfile
├── docker-compose.yml
├── helm/
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── secrets.yaml
├── .github/
│   └── workflows/
│       └── ci.yml
├── AGENTS.md
├── PLAN-v1.2.md
├── .env.example
└── .gitignore
```

### C. Model Example

```python
# backend/app/models/customer.py
from uuid import UUID, uuid4
from datetime import datetime
from sqlalchemy import String, Enum, Text, func
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.models.base import Base

class Customer(Base):
    __tablename__ = "customers"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4, init=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    email: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    phone: Mapped[str | None] = mapped_column(String(50), nullable=True, default=None)
    company: Mapped[str | None] = mapped_column(String(255), nullable=True, default=None)
    status: Mapped[str] = mapped_column(
        Enum("lead", "active", "churned", name="customer_status"),
        default="lead",
    )
    notes: Mapped[str | None] = mapped_column(Text, nullable=True, default=None)
    created_at: Mapped[datetime] = mapped_column(server_default=func.now(), init=False)
    updated_at: Mapped[datetime] = mapped_column(
        server_default=func.now(), onupdate=func.now(), init=False
    )

    deals: Mapped[list["Deal"]] = relationship(
        back_populates="customer", lazy="selectin", cascade="all, delete-orphan", init=False
    )
    activities: Mapped[list["Activity"]] = relationship(
        back_populates="customer", lazy="selectin", cascade="all, delete-orphan", init=False
    )
```

### D. Repository Example

```python
# backend/app/repositories/base_repo.py
from uuid import UUID
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from typing import TypeVar, Generic
from app.models.base import Base

T = TypeVar("T", bound=Base)

class BaseRepository(Generic[T]):
    def __init__(self, db: AsyncSession, model: type[T]):
        self.db = db
        self.model = model

    async def list_all(self, limit: int = 20, offset: int = 0) -> tuple[list[T], int]:
        result = await self.db.execute(
            select(self.model).limit(limit).offset(offset)
        )
        items = list(result.scalars().all())
        count_result = await self.db.execute(select(func.count()).select_from(self.model))
        total = count_result.scalar() or 0
        return items, total

    async def get_by_id(self, item_id: UUID) -> T | None:
        result = await self.db.execute(
            select(self.model).where(self.model.id == item_id)
        )
        return result.scalar_one_or_none()

    async def create(self, obj: T) -> T:
        self.db.add(obj)
        await self.db.commit()
        await self.db.refresh(obj)
        return obj

    async def delete(self, obj: T) -> None:
        await self.db.delete(obj)
        await self.db.commit()

    async def count(self) -> int:
        result = await self.db.execute(select(func.count()).select_from(self.model))
        return result.scalar() or 0
```

### E. Frontend Page Example

```tsx
// frontend/src/app/customers/page.tsx
"use client";
import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { listCustomers, createCustomer } from "@/lib/api";

export default function CustomersPage() {
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ name: "", email: "" });

  const load = async () => {
    setLoading(true);
    const data = await listCustomers(100, 0);
    setCustomers(data);
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    await createCustomer(form);
    setOpen(false);
    setForm({ name: "", email: "" });
    await load();
  };

  const filtered = customers?.items.filter((c) =>
    c.name.toLowerCase().includes(search.toLowerCase()) ||
    c.email.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Customers</h1>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger><Button>Add Customer</Button></DialogTrigger>
          <DialogContent>
            <DialogHeader><DialogTitle>New Customer</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <div><Label>Name</Label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
              <div><Label>Email</Label><Input value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} /></div>
              <Button onClick={handleCreate} disabled={!form.name || !form.email}>Save</Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>
      <Input placeholder="Search customers..." value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-sm" />
      <Card><CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow><TableHead>Name</TableHead><TableHead>Email</TableHead><TableHead>Status</TableHead></TableRow>
          </TableHeader>
          <TableBody>
            {loading && <TableRow><TableCell colSpan={3} className="text-center">Loading...</TableCell></TableRow>}
            {!loading && filtered?.length === 0 && <TableRow><TableCell colSpan={3} className="text-center">No customers found.</TableCell></TableRow>}
            {filtered?.map((c) => (
              <TableRow key={c.id}>
                <TableCell className="font-medium">{c.name}</TableCell>
                <TableCell>{c.email}</TableCell>
                <TableCell><Badge variant={c.status === "active" ? "default" : "secondary"}>{c.status}</Badge></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent></Card>
    </div>
  );
}
```

---

> **Document Version:** 1.0  
> **Last Updated:** 2026-05-09  
> **Maintainer:** DClaw Platform Team  
> **Source of Truth:** `dclaw-platform/ARCHITECTURE.md`
