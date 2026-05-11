# DClaw Stack — Agent Swarm Operations Manual

> **For AI Agents Joining the DClaw Swarm**
>
> Read this entire document before writing any code. This is your operating manual.
> If you follow this document, you will build enterprise-grade software at unprecedented velocity.
> If you ignore this document, you will break the scaffold and waste everyone's time.

---

## Table of Contents

1. [What Is the DClaw Stack?](#1-what-is-the-dclaw-stack)
2. [Your First 5 Minutes](#2-your-first-5-minutes)
3. [The Scaffold — Every File Explained](#3-the-scaffold--every-file-explained)
4. [How to Build a Feature — Step by Step](#4-how-to-build-a-feature--step-by-step)
5. [AGENTS.md — The Constitution](#5-agentsmd--the-constitution)
6. [PLAN-v1.2.md — Your Feature Spec](#6-plan-v12md--your-feature-spec)
7. [Backend Patterns](#7-backend-patterns)
8. [Frontend Patterns](#8-frontend-patterns)
9. [Docker & Deployment](#9-docker--deployment)
10. [CI/CD & Git Workflow](#10-cicd--git-workflow)
11. [Common Tasks & Recipes](#11-common-tasks--recipes)
12. [Troubleshooting](#12-troubleshooting)
13. [Full Working Example](#13-full-working-example)

---

## 1. What Is the DClaw Stack?

The DClaw Stack is **72 independent vertical SaaS applications** built from one canonical scaffold. Think of it as:

- **One codebase pattern** replicated 72 times
- **One team** (20 people + AI agents) building in parallel
- **Zero merge conflicts** because every app is its own repo
- **Enterprise quality from day one** because the scaffold includes tests, migrations, Docker, Helm, and CI/CD

You are an AI agent in this swarm. Your job is to read the spec (`PLAN-v1.2.md`), follow the rules (`AGENTS.md`), and write code that fits the scaffold.

### The Philosophy

**Velocity through Constraint.** We don't argue about frameworks. We don't debate folder structures. Every decision is made once, encoded in the scaffold, and replicated everywhere. Your job is NOT to reinvent. Your job is to **fill in the blanks**.

### Your Role

You will be assigned to a specific app repository (e.g., `dclaw-crm`, `dclaw-finance`, `dclaw-inventory`). Within that repo, you have full autonomy. You will:

1. Read `AGENTS.md` to learn the rules
2. Read `PLAN-v1.2.md` to learn what features to build
3. Write backend code (models → schemas → repositories → routers → tests → migrations)
4. Write frontend code (API types → pages → components)
5. Verify with `docker compose up -d`
6. Commit and push

---

## 2. Your First 5 Minutes

When you are assigned to a repo, do this in order:

### Step 1: Read AGENTS.md
```bash
cat AGENTS.md
```
This tells you:
- Your backend port and frontend port
- Your database name
- The architecture locks (what you CANNOT change)
- The anti-patterns (what you MUST NOT do)
- The directory structure

### Step 2: Read PLAN-v1.2.md
```bash
cat PLAN-v1.2.md
```
This tells you:
- What P0 features to build (Must Have)
- What P1 features to build (Should Have)
- What P2 features to build (Could Have)
- The implementation priority (what order to build in)

### Step 3: Check the Scaffold
```bash
ls backend/app/
ls frontend/src/
cat backend/app/models/base.py
cat backend/app/repositories/base_repo.py
cat frontend/src/lib/api.ts
```
Understand what already exists before you add anything.

### Step 4: Verify Docker Works
```bash
docker compose up -d
curl http://localhost:YOUR_BACKEND_PORT/health/
cd frontend && npm run build
```
If this doesn't work, fix it before writing any feature code.

### Step 5: Start Building
Pick the first P0 feature from `PLAN-v1.2.md` and follow the build process in Section 4.

---

## 3. The Scaffold — Every File Explained

The `dclaw-scaffold` is the source of truth. Every app repo is a copy of this scaffold with domain-specific code added. Here is every file and what it does.

### Backend Structure

```
backend/
├── app/
│   ├── __init__.py
│   ├── api/
│   │   ├── __init__.py
│   │   ├── main.py              ← FastAPI app factory with lifespan handler
│   │   ├── routes/
│   │   │   ├── __init__.py
│   │   │   └── health.py        ← GET /health/ endpoint
│   │   └── v1/                  ← YOUR ROUTERS GO HERE
│   │       ├── __init__.py
│   │       └── ...              ← e.g., customers.py, deals.py
│   ├── core/
│   │   ├── __init__.py
│   │   ├── config.py            ← Pydantic Settings (database_url, app_env, etc.)
│   │   └── database.py          ← create_async_engine, AsyncSession, get_db(), init_db()
│   ├── models/
│   │   ├── __init__.py          ← Imports all models for Alembic
│   │   ├── base.py              ← class Base(DeclarativeBase): pass
│   │   └── ...                  ← YOUR MODELS GO HERE
│   ├── repositories/
│   │   ├── __init__.py
│   │   ├── base_repo.py         ← Generic BaseRepository[T]
│   │   └── ...                  ← YOUR REPOSITORIES GO HERE
│   ├── schemas/
│   │   ├── __init__.py
│   │   └── ...                  ← YOUR PYDANTIC SCHEMAS GO HERE
│   ├── services/
│   │   ├── __init__.py
│   │   └── ...                  ← BUSINESS LOGIC / AI SERVICES GO HERE
│   └── utils/
│       ├── __init__.py
│       └── ...                  ← HELPERS (utc_now, etc.)
├── alembic/
│   ├── env.py                   ← Async Alembic environment
│   ├── script.py.mako
│   └── versions/                ← MIGRATION FILES GO HERE
├── tests/
│   ├── __init__.py              ← REQUIRED for pytest discovery
│   ├── conftest.py              ← Test DB override, AsyncClient fixture
│   └── test_*.py                ← YOUR TESTS GO HERE
├── Dockerfile                   ← python:3.11-slim, non-root, stdlib healthcheck
├── requirements.txt             ← Dependencies (pytest-asyncio==0.24.0 PINNED)
└── alembic.ini                  ← Alembic config
```

### Frontend Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx           ← Root layout with providers
│   │   ├── page.tsx             ← Dashboard / home page
│   │   ├── globals.css          ← Tailwind imports + CSS variables
│   │   └── ...                  ← YOUR PAGES GO HERE (e.g., customers/page.tsx)
│   ├── components/
│   │   └── ui/                  ← PRE-BUILT UI COMPONENTS
│   │       ├── button.tsx       ← Variants: default, destructive, outline, secondary, ghost, link
│   │       ├── card.tsx         ← Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter
│   │       ├── input.tsx        ← Standard text input
│   │       ├── label.tsx        ← Form label
│   │       ├── badge.tsx        ← Variants: default, secondary, destructive, outline
│   │       ├── select.tsx       ← Native select with onValueChange
│   │       ├── dialog.tsx       ← Modal with trigger, content, header, title
│   │       ├── table.tsx        ← Table, TableHeader, TableBody, TableRow, TableHead, TableCell
│   │       ├── tabs.tsx         ← Tabs, TabsList, TabsTrigger, TabsContent
│   │       └── avatar.tsx       ← Avatar, AvatarImage, AvatarFallback
│   └── lib/
│       ├── api.ts               ← Typed fetch wrapper + API functions
│       └── utils.ts             ← cn() helper (clsx + tailwind-merge)
├── package.json                 ← Next.js 14 + Tailwind v3 + tailwindcss-animate
├── next.config.js               ← output: 'standalone' for Docker
├── tailwind.config.ts           ← Theme config with CSS variables
└── Dockerfile                   ← node:20-alpine, ARG NEXT_PUBLIC_API_URL before build
```

### Root-Level Files

```
├── docker-compose.yml           ← PostgreSQL + backend + frontend with healthchecks
├── helm/                        ← Kubernetes deployment charts
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── secrets.yaml
├── .github/
│   └── workflows/
│       ├── ci.yml               ← pytest + npm run build on push/PR
│       ├── build-backend.yml    ← Docker build → GHCR (optional)
│       ├── build-frontend.yml   ← Docker build → GHCR (optional)
│       ├── deploy.yml           ← Helm deploy to K8s (optional)
│       └── claude.yml           ← Claude Code Action (optional)
├── AGENTS.md                    ← THE CONSTITUTION — READ THIS FIRST
├── PLAN-v1.2.md                 ← YOUR FEATURE ROADMAP
├── .env.example                 ← Environment variable template
└── .gitignore
```

---

## 4. How to Build a Feature — Step by Step

This is the exact process. Do not skip steps. Do not reorder steps.

### The Golden Path

```
Read PLAN-v1.2.md → Read AGENTS.md → Backend Model → Backend Schema → Backend Repository → Backend Router → Backend Tests → Alembic Migration → Frontend API Types → Frontend Page → Frontend Components → Docker Verify → Git Commit & Push
```

### Example: Building a "Customer" Feature

Let's say `PLAN-v1.2.md` says: **P0.3 — Contact Management: Customer CRUD with search**

#### Step 1: Read the Rules (1 minute)
```bash
cat AGENTS.md | grep -A 20 "Anti-Patterns"
```
Key reminders:
- Use `DeclarativeBase` in `app.models.base`
- Use `Mapped[...]` and `mapped_column()`
- Use `default=` not `default_factory=`
- Use `lazy="selectin"` for relationships
- Never use in-memory mock dicts

#### Step 2: Backend — Add Model (5 minutes)

Create `backend/app/models/customer.py`:

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
```

Update `backend/app/models/__init__.py`:
```python
from app.models.customer import Customer
# Add other models here
```

#### Step 3: Backend — Add Schema (3 minutes)

Create `backend/app/schemas/customer.py`:

```python
from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, ConfigDict

class CustomerBase(BaseModel):
    name: str
    email: str
    phone: str | None = None
    company: str | None = None
    status: str = "lead"
    notes: str | None = None

class CustomerCreate(CustomerBase):
    pass

class CustomerUpdate(BaseModel):
    name: str | None = None
    email: str | None = None
    phone: str | None = None
    company: str | None = None
    status: str | None = None
    notes: str | None = None

class CustomerResponse(CustomerBase):
    model_config = ConfigDict(from_attributes=True)
    id: UUID
    created_at: datetime
    updated_at: datetime

class CustomerListResponse(BaseModel):
    items: list[CustomerResponse]
    total: int
```

#### Step 4: Backend — Add Repository (3 minutes)

Create `backend/app/repositories/customer_repo.py`:

```python
from uuid import UUID
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from app.models.customer import Customer
from app.repositories.base_repo import BaseRepository

class CustomerRepository(BaseRepository[Customer]):
    def __init__(self, db: AsyncSession):
        super().__init__(db, Customer)

    async def update(self, customer: Customer, **kwargs) -> Customer:
        for key, value in kwargs.items():
            if value is not None and hasattr(customer, key):
                setattr(customer, key, value)
        await self.db.commit()
        await self.db.refresh(customer)
        return customer

    async def get_by_email(self, email: str) -> Customer | None:
        result = await self.db.execute(
            select(Customer).where(Customer.email == email)
        )
        return result.scalar_one_or_none()
```

#### Step 5: Backend — Add Router (5 minutes)

Create `backend/app/api/v1/customers.py`:

```python
from uuid import UUID
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from app.repositories.customer_repo import CustomerRepository
from app.schemas.customer import CustomerCreate, CustomerUpdate, CustomerResponse, CustomerListResponse
from app.models.customer import Customer

router = APIRouter()

@router.get("/", response_model=CustomerListResponse)
async def list_customers(
    limit: int = Query(20, ge=1, le=100),
    offset: int = Query(0, ge=0),
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    items, total = await repo.list_all(limit=limit, offset=offset)
    return CustomerListResponse(items=items, total=total)

@router.post("/", response_model=CustomerResponse, status_code=201)
async def create_customer(
    data: CustomerCreate,
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    existing = await repo.get_by_email(data.email)
    if existing:
        raise HTTPException(status_code=409, detail="Email already registered")
    customer = Customer(**data.model_dump())
    created = await repo.create(customer)
    return created

@router.get("/{customer_id}", response_model=CustomerResponse)
async def get_customer(
    customer_id: UUID,
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    customer = await repo.get_by_id(customer_id)
    if not customer:
        raise HTTPException(status_code=404, detail="Customer not found")
    return customer

@router.put("/{customer_id}", response_model=CustomerResponse)
async def update_customer(
    customer_id: UUID,
    data: CustomerUpdate,
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    customer = await repo.get_by_id(customer_id)
    if not customer:
        raise HTTPException(status_code=404, detail="Customer not found")
    update_data = {k: v for k, v in data.model_dump(exclude_unset=True).items() if v is not None}
    updated = await repo.update(customer, **update_data)
    return updated

@router.delete("/{customer_id}", status_code=204)
async def delete_customer(
    customer_id: UUID,
    db: AsyncSession = Depends(get_db),
):
    repo = CustomerRepository(db)
    customer = await repo.get_by_id(customer_id)
    if not customer:
        raise HTTPException(status_code=404, detail="Customer not found")
    await repo.delete(customer)
```

Wire it up in `backend/app/api/main.py`:
```python
from fastapi import FastAPI
from app.api.routes.health import router as health_router
from app.api.v1.customers import router as customers_router

app = FastAPI()
app.include_router(health_router, prefix="/health", tags=["health"])
app.include_router(customers_router, prefix="/api/v1/customers", tags=["customers"])
```

#### Step 6: Backend — Add Tests (5 minutes)

Create `backend/tests/test_customers.py`:

```python
import pytest
from httpx import AsyncClient
from uuid import uuid4

@pytest.mark.asyncio
async def test_list_customers(client: AsyncClient):
    response = await client.get("/api/v1/customers/")
    assert response.status_code == 200
    data = response.json()
    assert "items" in data
    assert "total" in data

@pytest.mark.asyncio
async def test_create_customer(client: AsyncClient):
    payload = {"name": "Alice", "email": "alice@example.com", "status": "lead"}
    response = await client.post("/api/v1/customers/", json=payload)
    assert response.status_code == 201
    data = response.json()
    assert data["name"] == "Alice"
    assert data["email"] == "alice@example.com"

@pytest.mark.asyncio
async def test_create_customer_duplicate_email(client: AsyncClient):
    payload = {"name": "Alice", "email": "alice@example.com"}
    await client.post("/api/v1/customers/", json=payload)
    response = await client.post("/api/v1/customers/", json=payload)
    assert response.status_code == 409

@pytest.mark.asyncio
async def test_get_customer(client: AsyncClient):
    create_resp = await client.post("/api/v1/customers/", json={"name": "Bob", "email": "bob@example.com"})
    customer_id = create_resp.json()["id"]
    response = await client.get(f"/api/v1/customers/{customer_id}")
    assert response.status_code == 200
    assert response.json()["name"] == "Bob"

@pytest.mark.asyncio
async def test_get_customer_not_found(client: AsyncClient):
    response = await client.get(f"/api/v1/customers/{uuid4()}")
    assert response.status_code == 404

@pytest.mark.asyncio
async def test_update_customer(client: AsyncClient):
    create_resp = await client.post("/api/v1/customers/", json={"name": "Charlie", "email": "charlie@example.com"})
    customer_id = create_resp.json()["id"]
    response = await client.put(f"/api/v1/customers/{customer_id}", json={"name": "Charles"})
    assert response.status_code == 200
    assert response.json()["name"] == "Charles"

@pytest.mark.asyncio
async def test_delete_customer(client: AsyncClient):
    create_resp = await client.post("/api/v1/customers/", json={"name": "Dave", "email": "dave@example.com"})
    customer_id = create_resp.json()["id"]
    response = await client.delete(f"/api/v1/customers/{customer_id}")
    assert response.status_code == 204
    get_resp = await client.get(f"/api/v1/customers/{customer_id}")
    assert get_resp.status_code == 404
```

#### Step 7: Backend — Generate Migration (2 minutes)

```bash
cd backend
alembic revision --autogenerate -m "add customers table"
```

#### Step 8: Frontend — Add API Types (3 minutes)

Update `frontend/src/lib/api.ts`:

```typescript
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE}${path}`;
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  });
  if (!response.ok) {
    const error = await response.text();
    throw new ApiError(`API error ${response.status}: ${error}`, response.status);
  }
  return response.json();
}

// === Health ===
export async function getHealth() {
  return fetchJson<{ status: string }>("/health/");
}

// === Customers ===
export interface Customer {
  id: string;
  name: string;
  email: string;
  phone: string | null;
  company: string | null;
  status: string;
  notes: string | null;
  created_at: string;
  updated_at: string;
}

export interface CustomerCreate {
  name: string;
  email: string;
  phone?: string;
  company?: string;
  status?: string;
  notes?: string;
}

export interface CustomerListResponse {
  items: Customer[];
  total: number;
}

export async function listCustomers(limit = 20, offset = 0): Promise<CustomerListResponse> {
  return fetchJson<CustomerListResponse>(`/api/v1/customers/?limit=${limit}&offset=${offset}`);
}

export async function createCustomer(data: CustomerCreate): Promise<Customer> {
  return fetchJson<Customer>("/api/v1/customers/", { method: "POST", body: JSON.stringify(data) });
}

export async function getCustomer(id: string): Promise<Customer> {
  return fetchJson<Customer>(`/api/v1/customers/${id}`);
}

export async function updateCustomer(id: string, data: Partial<CustomerCreate>): Promise<Customer> {
  return fetchJson<Customer>(`/api/v1/customers/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export async function deleteCustomer(id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/v1/customers/${id}`, { method: "DELETE" });
  if (!response.ok) {
    throw new ApiError(`API error ${response.status}`, response.status);
  }
}

export { ApiError };
```

#### Step 9: Frontend — Build Page (10 minutes)

Create `frontend/src/app/customers/page.tsx`:

```tsx
"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { listCustomers, createCustomer, type CustomerCreate } from "@/lib/api";

export default function CustomersPage() {
  const [customers, setCustomers] = useState<{ items: any[]; total: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<CustomerCreate>({ name: "", email: "" });

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

  const filtered = customers?.items.filter((c) => {
    const q = search.toLowerCase();
    return (
      c.name.toLowerCase().includes(q) ||
      c.email.toLowerCase().includes(q) ||
      (c.company ?? "").toLowerCase().includes(q)
    );
  });

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
              <div><Label>Phone</Label><Input value={form.phone ?? ""} onChange={(e) => setForm({ ...form, phone: e.target.value })} /></div>
              <div><Label>Company</Label><Input value={form.company ?? ""} onChange={(e) => setForm({ ...form, company: e.target.value })} /></div>
              <Button onClick={handleCreate} disabled={!form.name || !form.email}>Save</Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      <Input placeholder="Search customers..." value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-sm" />

      <Card><CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead><TableHead>Email</TableHead><TableHead>Company</TableHead><TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading && <TableRow><TableCell colSpan={4} className="text-center text-slate-500">Loading...</TableCell></TableRow>}
            {!loading && filtered?.length === 0 && <TableRow><TableCell colSpan={4} className="text-center text-slate-500">No customers found.</TableCell></TableRow>}
            {filtered?.map((c) => (
              <TableRow key={c.id}>
                <TableCell><Link href={`/customers/${c.id}`} className="font-medium hover:underline">{c.name}</Link></TableCell>
                <TableCell>{c.email}</TableCell>
                <TableCell>{c.company ?? "—"}</TableCell>
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

#### Step 10: Verify Docker (2 minutes)

```bash
docker compose up -d
curl http://localhost:YOUR_BACKEND_PORT/health/
cd frontend && npm run build
cd backend && pytest
```

All must pass.

#### Step 11: Commit & Push (1 minute)

```bash
git add .
git commit -m "feat: add customer CRUD per PLAN-v1.2.md P0.3

- Customer model with status enum
- Full CRUD repository + router
- 6 test cases covering all endpoints
- Frontend page with search, create dialog, table
- Alembic migration"
git push origin main
```

**Total time: ~40 minutes for a complete CRUD feature.**

---

## 5. AGENTS.md — The Constitution

### Architecture Locks (NEVER CHANGE)

| Layer | Rule | Violation |
|-------|------|-----------|
| **Backend Framework** | FastAPI with `lifespan` handler | Using Flask, Django, or sync handlers |
| **ORM Base** | `DeclarativeBase` from `app.models.base` | Using `MappedAsDataclass` or `declarative_base()` |
| **Schema Version** | Pydantic v2 with `ConfigDict(from_attributes=True)` | Using Pydantic v1 `orm_mode` |
| **DB Driver** | Async SQLAlchemy (`create_async_engine` + `AsyncSession`) | Using sync `create_engine` |
| **Data Access** | Repository pattern in `app/repositories/` | Inline DB access in routers |
| **Session Management** | `Depends(get_db)` dependency injection | Manual `AsyncSession` with `__anext__()` |
| **Mock Data** | NEVER use in-memory dicts | Using `MOCK_CUSTOMERS = {}` |
| **Test Framework** | `pytest-asyncio==0.24.0` PINNED | Upgrading to newer versions |
| **Frontend Framework** | Next.js 14+ App Router | Using Pages Router |
| **Styling** | Tailwind CSS v3 + pre-built components | Installing shadcn CLI v4 or `@base-ui/react` |
| **API Client** | Typed `fetchJson<T>()` in `src/lib/api.ts` | Using axios or untyped fetch |
| **Build Var** | `NEXT_PUBLIC_API_URL` baked at Docker build | Hardcoding `localhost:PORT` |
| **Backend Image** | `python:3.11-slim` with non-root `appuser` | Running as root |
| **Healthcheck** | Python stdlib `urllib.request.urlopen()` | Using `curl` (not in slim image) |
| **CI File** | `.github/workflows/ci.yml` must exist | Deleting CI workflow |
| **Utils File** | `src/lib/utils.ts` with `cn()` | Deleting the file |
| **Animate Dep** | `tailwindcss-animate` in `dependencies` | Putting it in `devDependencies` |

### Anti-Patterns (NEVER DO)

| Anti-Pattern | Why It Breaks | Correct Alternative |
|--------------|---------------|---------------------|
| `declarative_base()` in `database.py` | Separate metadata → zero tables | `from app.models.base import Base` |
| `curl` in healthcheck on `python:*-slim` | No `curl` binary → silent failure | `python -c "import urllib.request; urllib.request.urlopen(...)"` |
| In-memory `MOCK_*` dicts | Data lost on restart | Create repository + real DB |
| Missing `ARG NEXT_PUBLIC_API_URL` | Wrong API URL baked in | Add `ARG` before `RUN npm run build` |
| Manual `get_db()` with `__anext__()` | Session leaks | `Depends(get_db)` |
| Hardcoded `localhost:PORT` | Breaks Docker/K8s | Use `process.env.NEXT_PUBLIC_API_URL` |
| No alembic migration for new models | Schema drift | `alembic revision --autogenerate` |
| Installing `shadcn` CLI v4 | Breaks Tailwind v3 | Use pre-built components |
| Using `@base-ui/react` | Incompatible with Tailwind v3 | Use pre-built components |
| Non-standard Postgres port in tests | CI maps 5432 only | Always use `localhost:5432` |
| Upgrading `pytest-asyncio` | v1.3.0 breaks fixtures | Keep `pytest-asyncio==0.24.0` |
| Deleting `.github/workflows/ci.yml` | No CI runs | Leave CI workflow intact |
| Missing `src/lib/utils.ts` | UI components fail | Already in scaffold — don't delete |
| `MappedAsDataclass` in `Base` | Relationship conflicts | Use plain `DeclarativeBase` |
| `default_factory` in `mapped_column()` | `ArgumentError` | Use `default=` with callable |
| Timezone-aware `datetime` | `DataError` with `TIMESTAMP WITHOUT TIME ZONE` | Use naive UTC via `utc_now()` |

### Database Rules

1. All models MUST inherit from `Base` in `app.models.base`
2. All models MUST use `Mapped[...]` and `mapped_column()`
3. **Never use `default_factory=` in `mapped_column()`** — use `default=` instead
4. Relationships MUST specify `lazy="selectin"`
5. All new tables MUST get an alembic migration
6. Use `ondelete="CASCADE"` for child tables
7. Use `ondelete="SET NULL"` for optional references

### How to Add a Feature (From AGENTS.md)

1. **Read this file** and `PLAN-v1.2.md`
2. **Backend:**
   - Add/update model in `app/models/`
   - Add/update schema in `app/schemas/`
   - Add repository in `app/repositories/`
   - Add/update router in `app/api/v1/`
   - Add tests in `tests/`
   - Generate alembic migration
3. **Frontend:**
   - Add API types/functions to `src/lib/api.ts`
   - Add page in `src/app/` or component using pre-built UI components
4. **Docker:** Verify `docker compose config` and `docker compose up -d`
5. **Commit** with conventional commit message

### Testing Requirements

- Every new repository MUST have tests
- Every new router endpoint MUST be covered
- Use `pytest-asyncio` with `async` test functions and `@pytest.mark.asyncio`
- Use `httpx.AsyncClient` with `ASGITransport`
- Override `get_db` dependency with test session in `conftest.py`
- Tests MUST use `localhost:5432` for PostgreSQL (CI requirement)

---

## 6. PLAN-v1.2.md — Your Feature Spec

Every app has a `PLAN-v1.2.md` that serves as your product specification. Read it before writing any code.

### Structure

```markdown
# {AppName} — v1.2 Feature Roadmap

## Pre-Flight Checklist
- [ ] package-lock.json committed
- [ ] next-env.d.ts committed
- [ ] Docker healthchecks correct
- [ ] Dockerfile ARG correct

## v1.0 Feature Inventory (Current)
- [ ] Core CRUD
- [ ] Dashboard
- [ ] Real backend (no mocks)
- [ ] Docker + Helm
- [ ] Alembic migrations
- [ ] Backend tests

## v1.2 Roadmap

### P0 — Must Have (Ship in v1.0, demo-ready)

#### 1. AI {App} Copilot
**Description:** AI assistant that helps with domain tasks
**AI Angle:** What the AI does specifically
**Backend:** Which endpoints to create
**Frontend:** Which components to build
**Files:** Where to put the code

#### 2. Core Feature A
#### 3. Core Feature B
#### 4. Core Feature C

### P1 — Should Have (v1.1–1.2)
#### 5. Advanced Feature D
#### 6. Advanced Feature E
#### 7. Advanced Feature F
#### 8. Advanced Feature G

### P2 — Could Have (v1.3+)
#### 9. Future Feature H
#### 10. Future Feature I
#### 11. Future Feature J
#### 12. Future Feature K

## Implementation Priority
1. Week 1–2: P0.1 + P0.2
2. Week 3–4: P0.3 + P0.4
3. Week 5–6: P1.5 + P1.6
4. Week 7–8: P1.7 + P1.8
```

### How to Use It

1. Start with **P0 features**. These are your demo-ready must-haves.
2. For each feature, the plan tells you:
   - **Description:** What the feature does
   - **AI Angle:** How AI enhances it (every P0 has an AI component)
   - **Backend:** Which API endpoints to create
   - **Frontend:** Which pages/components to build
   - **Files:** Suggested file paths
3. Build in the order specified in **Implementation Priority**.
4. When you finish a feature, check it off in the inventory.

---

## 7. Backend Patterns

### The Model Pattern

Every model follows this exact structure:

```python
from uuid import UUID, uuid4
from datetime import datetime
from sqlalchemy import String, Integer, Boolean, ForeignKey, Enum, Text, func
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.models.base import Base

class Entity(Base):
    __tablename__ = "entities"  # Plural, lowercase, snake_case

    # Primary key — always UUID, always auto-generated
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4, init=False)

    # Required fields
    name: Mapped[str] = mapped_column(String(255), nullable=False)

    # Optional fields — always with default=None
    description: Mapped[str | None] = mapped_column(Text, nullable=True, default=None)

    # Enum fields
    status: Mapped[str] = mapped_column(
        Enum("active", "inactive", name="entity_status"),
        default="active",
    )

    # Foreign key — optional references use nullable=True
    parent_id: Mapped[UUID | None] = mapped_column(
        ForeignKey("parents.id", ondelete="SET NULL"),
        nullable=True,
        default=None,
    )

    # Timestamps — always server-managed
    created_at: Mapped[datetime] = mapped_column(server_default=func.now(), init=False)
    updated_at: Mapped[datetime] = mapped_column(
        server_default=func.now(), onupdate=func.now(), init=False
    )

    # Relationships — always lazy="selectin", cascade for children
    children: Mapped[list["Child"]] = relationship(
        back_populates="parent", lazy="selectin", cascade="all, delete-orphan", init=False
    )
    parent: Mapped["Parent"] = relationship(back_populates="children", init=False)
```

### The Schema Pattern

```python
from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, ConfigDict

# Base — shared fields
class EntityBase(BaseModel):
    name: str
    description: str | None = None
    status: str = "active"

# Create — what the client sends to create
class EntityCreate(EntityBase):
    pass

# Update — what the client sends to update (all optional)
class EntityUpdate(BaseModel):
    name: str | None = None
    description: str | None = None
    status: str | None = None

# Response — what the API returns (from_attributes=True is REQUIRED)
class EntityResponse(EntityBase):
    model_config = ConfigDict(from_attributes=True)
    id: UUID
    created_at: datetime
    updated_at: datetime

# List response — paginated
class EntityListResponse(BaseModel):
    items: list[EntityResponse]
    total: int
```

### The Repository Pattern

```python
from uuid import UUID
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from app.models.entity import Entity
from app.repositories.base_repo import BaseRepository

class EntityRepository(BaseRepository[Entity]):
    def __init__(self, db: AsyncSession):
        super().__init__(db, Entity)

    # Override update to handle partial updates
    async def update(self, entity: Entity, **kwargs) -> Entity:
        for key, value in kwargs.items():
            if value is not None and hasattr(entity, key):
                setattr(entity, key, value)
        await self.db.commit()
        await self.db.refresh(entity)
        return entity

    # Add custom queries here
    async def get_by_name(self, name: str) -> Entity | None:
        result = await self.db.execute(
            select(Entity).where(Entity.name == name)
        )
        return result.scalar_one_or_none()
```

### The Router Pattern

```python
from uuid import UUID
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from app.repositories.entity_repo import EntityRepository
from app.schemas.entity import EntityCreate, EntityUpdate, EntityResponse, EntityListResponse
from app.models.entity import Entity

router = APIRouter()

@router.get("/", response_model=EntityListResponse)
async def list_entities(
    limit: int = Query(20, ge=1, le=100),
    offset: int = Query(0, ge=0),
    db: AsyncSession = Depends(get_db),
):
    repo = EntityRepository(db)
    items, total = await repo.list_all(limit=limit, offset=offset)
    return EntityListResponse(items=items, total=total)

@router.post("/", response_model=EntityResponse, status_code=201)
async def create_entity(data: EntityCreate, db: AsyncSession = Depends(get_db)):
    repo = EntityRepository(db)
    entity = Entity(**data.model_dump())
    return await repo.create(entity)

@router.get("/{entity_id}", response_model=EntityResponse)
async def get_entity(entity_id: UUID, db: AsyncSession = Depends(get_db)):
    repo = EntityRepository(db)
    entity = await repo.get_by_id(entity_id)
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    return entity

@router.put("/{entity_id}", response_model=EntityResponse)
async def update_entity(
    entity_id: UUID,
    data: EntityUpdate,
    db: AsyncSession = Depends(get_db),
):
    repo = EntityRepository(db)
    entity = await repo.get_by_id(entity_id)
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    update_data = {k: v for k, v in data.model_dump(exclude_unset=True).items() if v is not None}
    return await repo.update(entity, **update_data)

@router.delete("/{entity_id}", status_code=204)
async def delete_entity(entity_id: UUID, db: AsyncSession = Depends(get_db)):
    repo = EntityRepository(db)
    entity = await repo.get_by_id(entity_id)
    if not entity:
        raise HTTPException(status_code=404, detail="Entity not found")
    await repo.delete(entity)
```

### The Test Pattern

```python
import pytest
from httpx import AsyncClient
from uuid import uuid4

@pytest.mark.asyncio
async def test_list_entities(client: AsyncClient):
    response = await client.get("/api/v1/entities/")
    assert response.status_code == 200
    data = response.json()
    assert "items" in data
    assert "total" in data

@pytest.mark.asyncio
async def test_create_entity(client: AsyncClient):
    payload = {"name": "Test Entity"}
    response = await client.post("/api/v1/entities/", json=payload)
    assert response.status_code == 201
    assert response.json()["name"] == "Test Entity"

@pytest.mark.asyncio
async def test_get_entity_not_found(client: AsyncClient):
    response = await client.get(f"/api/v1/entities/{uuid4()}")
    assert response.status_code == 404
```

---

## 8. Frontend Patterns

### The Page Pattern

Every CRUD page follows this structure:

```tsx
"use client";

import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { listEntities, createEntity } from "@/lib/api";

export default function EntitiesPage() {
  const [entities, setEntities] = useState<{ items: any[]; total: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ name: "" });

  const load = async () => {
    setLoading(true);
    const data = await listEntities(100, 0);
    setEntities(data);
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    await createEntity(form);
    setOpen(false);
    setForm({ name: "" });
    await load();
  };

  const filtered = entities?.items.filter((e) =>
    e.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-4">
      {/* Header with Create Button */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Entities</h1>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger><Button>Add Entity</Button></DialogTrigger>
          <DialogContent>
            <DialogHeader><DialogTitle>New Entity</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <div><Label>Name</Label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
              <Button onClick={handleCreate} disabled={!form.name}>Save</Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      {/* Search */}
      <Input placeholder="Search..." value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-sm" />

      {/* Table */}
      <Card><CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow><TableHead>Name</TableHead><TableHead>Status</TableHead></TableRow>
          </TableHeader>
          <TableBody>
            {loading && <TableRow><TableCell colSpan={2} className="text-center">Loading...</TableCell></TableRow>}
            {filtered?.map((e) => (
              <TableRow key={e.id}>
                <TableCell className="font-medium">{e.name}</TableCell>
                <TableCell><Badge>{e.status}</Badge></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent></Card>
    </div>
  );
}
```

### Pre-Built UI Components

Use these components. Do NOT install shadcn CLI or `@base-ui/react`.

```tsx
import { Button } from "@/components/ui/button";
// Variants: default, destructive, outline, secondary, ghost, link
// Sizes: default, sm, lg, icon
<Button variant="destructive" size="sm">Delete</Button>

import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card";
<Card><CardHeader><CardTitle>Title</CardTitle></CardHeader><CardContent>Content</CardContent></Card>

import { Input } from "@/components/ui/input";
<Input placeholder="Type here..." value={value} onChange={(e) => setValue(e.target.value)} />

import { Badge } from "@/components/ui/badge";
// Variants: default, secondary, destructive, outline
<Badge variant="secondary">Draft</Badge>

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
<Dialog open={open} onOpenChange={setOpen}>
  <DialogTrigger><Button>Open</Button></DialogTrigger>
  <DialogContent><DialogHeader><DialogTitle>Title</DialogTitle></DialogHeader></DialogContent>
</Dialog>

import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
<Select onValueChange={(v) => setValue(v)}>
  <SelectTrigger><SelectValue placeholder="Select..." /></SelectTrigger>
  <SelectContent>
    <SelectItem value="a">Option A</SelectItem>
    <SelectItem value="b">Option B</SelectItem>
  </SelectContent>
</Select>

import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
<Tabs defaultValue="tab1">
  <TabsList><TabsTrigger value="tab1">Tab 1</TabsTrigger><TabsTrigger value="tab2">Tab 2</TabsTrigger></TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
  <TabsContent value="tab2">Content 2</TabsContent>
</Tabs>
```

### The API Client Pattern

```typescript
// frontend/src/lib/api.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE}${path}`;
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  });
  if (!response.ok) {
    const error = await response.text();
    throw new ApiError(`API error ${response.status}: ${error}`, response.status);
  }
  return response.json();
}

// Add app-specific functions here
export async function listEntities(limit = 20, offset = 0) {
  return fetchJson<{ items: any[]; total: number }>(`/api/v1/entities/?limit=${limit}&offset=${offset}`);
}

export async function createEntity(data: { name: string }) {
  return fetchJson<any>("/api/v1/entities/", { method: "POST", body: JSON.stringify(data) });
}

export { ApiError };
```

---

## 9. Docker & Deployment

### docker-compose.yml

Every app has this structure:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: dclaw_app
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build: ./backend
    ports:
      - "YOUR_BACKEND_PORT:8000"
    environment:
      DATABASE_URL: postgresql+asyncpg://postgres:postgres@postgres:5432/dclaw_app
      APP_ENV: dev
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "python", "-c", "import urllib.request; urllib.request.urlopen('http://localhost:8000/health/')"]
      interval: 10s
      timeout: 5s
      retries: 3

  frontend:
    build:
      context: ./frontend
      args:
        NEXT_PUBLIC_API_URL: http://localhost:YOUR_BACKEND_PORT
    ports:
      - "YOUR_FRONTEND_PORT:3000"
    depends_on:
      - backend
```

### Backend Dockerfile

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY app/ ./app/
COPY alembic/ ./alembic/
COPY alembic.ini .

# Create non-root user
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

# Healthcheck using stdlib (no curl in slim)
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health/')"

EXPOSE 8000

CMD ["uvicorn", "app.api.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### Frontend Dockerfile

```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

# ARG must be declared before RUN npm run build
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
```

### Verify Docker

```bash
# Build and start everything
docker compose up -d --build

# Check health
curl http://localhost:YOUR_BACKEND_PORT/health/

# Check logs
docker compose logs -f backend
docker compose logs -f frontend

# Run tests
cd backend && pytest

# Build frontend
cd frontend && npm run build

# Stop everything
docker compose down
```

---

## 10. CI/CD & Git Workflow

### Standard CI Workflow (ci.yml)

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

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
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with: { python-version: "3.11" }
      - run: pip install uv && uv pip install -r backend/requirements.txt
      - run: cd backend && python -m pytest -v
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

### Git Workflow

```bash
# 1. Pull latest
git pull origin main

# 2. Make changes (follow the build process in Section 4)

# 3. Test locally
docker compose up -d
cd backend && pytest
cd frontend && npm run build

# 4. Commit with conventional messages
# feat: new feature
# fix: bug fix
# docs: documentation
# test: adding tests
# refactor: code refactoring
# chore: maintenance
git add .
git commit -m "feat: add customer CRUD per PLAN-v1.2.md P0.3"

# 5. Push
git push origin main
```

---

## 11. Common Tasks & Recipes

### Add a New Model

```bash
# 1. Create model file
touch backend/app/models/product.py
# 2. Write model (see Section 7 pattern)
# 3. Update backend/app/models/__init__.py
# 4. Create schema in backend/app/schemas/product.py
# 5. Create repository in backend/app/repositories/product_repo.py
# 6. Create router in backend/app/api/v1/products.py
# 7. Wire router in backend/app/api/main.py
# 8. Write tests in backend/tests/test_products.py
# 9. Generate migration
cd backend && alembic revision --autogenerate -m "add products"
```

### Add a New Page

```bash
# 1. Create page directory
mkdir -p frontend/src/app/products
# 2. Create page.tsx (see Section 8 pattern)
touch frontend/src/app/products/page.tsx
# 3. Add API functions to frontend/src/lib/api.ts
# 4. Build and verify
cd frontend && npm run build
```

### Add an AI Service

```python
# backend/app/services/ai_service.py
from openai import AsyncOpenAI
from app.core.config import settings

client = AsyncOpenAI(api_key=settings.openai_api_key)

async def generate_summary(text: str) -> str:
    response = await client.chat.completions.create(
        model="gpt-4",
        messages=[
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": f"Summarize this: {text}"},
        ],
    )
    return response.choices[0].message.content
```

### Run Migrations

```bash
cd backend
# Generate
alembic revision --autogenerate -m "description"
# Upgrade
alembic upgrade head
# Downgrade
alembic downgrade -1
# History
alembic history
```

### Connect to Another App's API

```python
# In your service layer
import httpx

async def fetch_from_inventory(sku: str):
    async with httpx.AsyncClient() as client:
        response = await client.get(f"http://dclaw-inventory:8000/api/v1/items/{sku}")
        return response.json()
```

---

## 12. Troubleshooting

### Backend won't start

```bash
# Check if postgres is running
docker compose ps
# Check logs
docker compose logs backend
# Common issues:
# - Database URL wrong → check .env
# - Port conflict → check PORT_REGISTRY.md
# - Migration not run → cd backend && alembic upgrade head
```

### Tests fail

```bash
# Common issues:
# - pytest-asyncio version wrong → must be 0.24.0
# - Database not accessible → check localhost:5432
# - Missing __init__.py in tests/ → touch backend/tests/__init__.py
# - Model not imported in models/__init__.py → add import
```

### Frontend build fails

```bash
# Common issues:
# - tailwindcss-animate missing from dependencies → npm install tailwindcss-animate
# - next-env.d.ts missing → npx next build will regenerate
# - Missing utils.ts → check src/lib/utils.ts exists
# - Import path wrong → use @/components/ui/ not relative paths
```

### Docker build fails

```bash
# Common issues:
# - Missing ARG NEXT_PUBLIC_API_URL in frontend Dockerfile
# - curl in healthcheck → use python stdlib instead
# - Port mismatch → check EXPOSE and docker-compose ports align
```

### CI fails

```bash
# Common issues:
# - Postgres service not starting → check healthcheck in ci.yml
# - Tests using wrong port → must use localhost:5432
# - Frontend build timeout → increase timeout-minutes
# - Missing package-lock.json → npm install && git add package-lock.json
```

---

## 13. Full Working Example

Here is a complete, minimal working app that demonstrates every pattern. Copy this and adapt it.

### File: `backend/app/models/item.py`

```python
from uuid import UUID, uuid4
from datetime import datetime
from sqlalchemy import String, Boolean, func
from sqlalchemy.orm import Mapped, mapped_column
from app.models.base import Base

class Item(Base):
    __tablename__ = "items"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4, init=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str | None] = mapped_column(String(1000), nullable=True, default=None)
    completed: Mapped[bool] = mapped_column(default=False)
    created_at: Mapped[datetime] = mapped_column(server_default=func.now(), init=False)
    updated_at: Mapped[datetime] = mapped_column(
        server_default=func.now(), onupdate=func.now(), init=False
    )
```

### File: `backend/app/schemas/item.py`

```python
from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, ConfigDict

class ItemBase(BaseModel):
    name: str
    description: str | None = None
    completed: bool = False

class ItemCreate(ItemBase):
    pass

class ItemUpdate(BaseModel):
    name: str | None = None
    description: str | None = None
    completed: bool | None = None

class ItemResponse(ItemBase):
    model_config = ConfigDict(from_attributes=True)
    id: UUID
    created_at: datetime
    updated_at: datetime

class ItemListResponse(BaseModel):
    items: list[ItemResponse]
    total: int
```

### File: `backend/app/repositories/item_repo.py`

```python
from uuid import UUID
from sqlalchemy.ext.asyncio import AsyncSession
from app.models.item import Item
from app.repositories.base_repo import BaseRepository

class ItemRepository(BaseRepository[Item]):
    def __init__(self, db: AsyncSession):
        super().__init__(db, Item)

    async def update(self, item: Item, **kwargs) -> Item:
        for key, value in kwargs.items():
            if value is not None and hasattr(item, key):
                setattr(item, key, value)
        await self.db.commit()
        await self.db.refresh(item)
        return item
```

### File: `backend/app/api/v1/items.py`

```python
from uuid import UUID
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from app.core.database import get_db
from app.repositories.item_repo import ItemRepository
from app.schemas.item import ItemCreate, ItemUpdate, ItemResponse, ItemListResponse
from app.models.item import Item

router = APIRouter()

@router.get("/", response_model=ItemListResponse)
async def list_items(
    limit: int = Query(20, ge=1, le=100),
    offset: int = Query(0, ge=0),
    db: AsyncSession = Depends(get_db),
):
    repo = ItemRepository(db)
    items, total = await repo.list_all(limit=limit, offset=offset)
    return ItemListResponse(items=items, total=total)

@router.post("/", response_model=ItemResponse, status_code=201)
async def create_item(data: ItemCreate, db: AsyncSession = Depends(get_db)):
    repo = ItemRepository(db)
    item = Item(**data.model_dump())
    return await repo.create(item)

@router.get("/{item_id}", response_model=ItemResponse)
async def get_item(item_id: UUID, db: AsyncSession = Depends(get_db)):
    repo = ItemRepository(db)
    item = await repo.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    return item

@router.put("/{item_id}", response_model=ItemResponse)
async def update_item(
    item_id: UUID, data: ItemUpdate, db: AsyncSession = Depends(get_db),
):
    repo = ItemRepository(db)
    item = await repo.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    update_data = {k: v for k, v in data.model_dump(exclude_unset=True).items() if v is not None}
    return await repo.update(item, **update_data)

@router.delete("/{item_id}", status_code=204)
async def delete_item(item_id: UUID, db: AsyncSession = Depends(get_db)):
    repo = ItemRepository(db)
    item = await repo.get_by_id(item_id)
    if not item:
        raise HTTPException(status_code=404, detail="Item not found")
    await repo.delete(item)
```

### File: `backend/tests/test_items.py`

```python
import pytest
from httpx import AsyncClient
from uuid import uuid4

@pytest.mark.asyncio
async def test_list_items(client: AsyncClient):
    response = await client.get("/api/v1/items/")
    assert response.status_code == 200
    data = response.json()
    assert "items" in data
    assert "total" in data

@pytest.mark.asyncio
async def test_create_item(client: AsyncClient):
    payload = {"name": "Test Item", "completed": False}
    response = await client.post("/api/v1/items/", json=payload)
    assert response.status_code == 201
    assert response.json()["name"] == "Test Item"

@pytest.mark.asyncio
async def test_get_item_not_found(client: AsyncClient):
    response = await client.get(f"/api/v1/items/{uuid4()}")
    assert response.status_code == 404

@pytest.mark.asyncio
async def test_update_item(client: AsyncClient):
    create_resp = await client.post("/api/v1/items/", json={"name": "Old Name"})
    item_id = create_resp.json()["id"]
    response = await client.put(f"/api/v1/items/{item_id}", json={"name": "New Name"})
    assert response.status_code == 200
    assert response.json()["name"] == "New Name"

@pytest.mark.asyncio
async def test_delete_item(client: AsyncClient):
    create_resp = await client.post("/api/v1/items/", json={"name": "To Delete"})
    item_id = create_resp.json()["id"]
    response = await client.delete(f"/api/v1/items/{item_id}")
    assert response.status_code == 204
```

### File: `frontend/src/lib/api.ts` (additions)

```typescript
export interface Item {
  id: string;
  name: string;
  description: string | null;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface ItemCreate {
  name: string;
  description?: string;
  completed?: boolean;
}

export async function listItems(limit = 20, offset = 0) {
  return fetchJson<{ items: Item[]; total: number }>(`/api/v1/items/?limit=${limit}&offset=${offset}`);
}

export async function createItem(data: ItemCreate) {
  return fetchJson<Item>("/api/v1/items/", { method: "POST", body: JSON.stringify(data) });
}

export async function getItem(id: string) {
  return fetchJson<Item>(`/api/v1/items/${id}`);
}

export async function updateItem(id: string, data: Partial<ItemCreate>) {
  return fetchJson<Item>(`/api/v1/items/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

export async function deleteItem(id: string) {
  const response = await fetch(`${API_BASE}/api/v1/items/${id}`, { method: "DELETE" });
  if (!response.ok) throw new ApiError(`API error ${response.status}`, response.status);
}
```

### File: `frontend/src/app/items/page.tsx`

```tsx
"use client";
import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { listItems, createItem, updateItem, deleteItem, type ItemCreate } from "@/lib/api";

export default function ItemsPage() {
  const [items, setItems] = useState<{ items: any[]; total: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<ItemCreate>({ name: "" });

  const load = async () => { setLoading(true); setItems(await listItems(100, 0)); setLoading(false); };
  useEffect(() => { load(); }, []);

  const handleCreate = async () => { await createItem(form); setOpen(false); setForm({ name: "" }); await load(); };
  const toggleComplete = async (id: string, completed: boolean) => { await updateItem(id, { completed: !completed }); await load(); };
  const handleDelete = async (id: string) => { await deleteItem(id); await load(); };

  const filtered = items?.items.filter((i) => i.name.toLowerCase().includes(search.toLowerCase()));

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Items</h1>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger><Button>Add Item</Button></DialogTrigger>
          <DialogContent>
            <DialogHeader><DialogTitle>New Item</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <div><Label>Name</Label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
              <div><Label>Description</Label><Input value={form.description ?? ""} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
              <Button onClick={handleCreate} disabled={!form.name}>Save</Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>
      <Input placeholder="Search items..." value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-sm" />
      <Card><CardContent className="p-0">
        <Table>
          <TableHeader><TableRow><TableHead>Name</TableHead><TableHead>Status</TableHead><TableHead>Actions</TableHead></TableRow></TableHeader>
          <TableBody>
            {loading && <TableRow><TableCell colSpan={3} className="text-center">Loading...</TableCell></TableRow>}
            {filtered?.map((i) => (
              <TableRow key={i.id}>
                <TableCell className={i.completed ? "line-through text-slate-400" : "font-medium"}>{i.name}</TableCell>
                <TableCell><Badge variant={i.completed ? "default" : "secondary"}>{i.completed ? "Done" : "Open"}</Badge></TableCell>
                <TableCell>
                  <Button size="sm" variant="ghost" onClick={() => toggleComplete(i.id, i.completed)}>{i.completed ? "Undo" : "Done"}</Button>
                  <Button size="sm" variant="destructive" onClick={() => handleDelete(i.id)}>Delete</Button>
                </TableCell>
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

## Quick Reference Card

| Task | Command | File |
|------|---------|------|
| Start stack | `docker compose up -d` | `./` |
| Run backend tests | `cd backend && pytest` | `backend/tests/` |
| Build frontend | `cd frontend && npm run build` | `frontend/src/` |
| Generate migration | `cd backend && alembic revision --autogenerate -m "desc"` | `backend/alembic/versions/` |
| Upgrade DB | `cd backend && alembic upgrade head` | — |
| Add model | Create in `backend/app/models/` | `backend/app/models/*.py` |
| Add router | Create in `backend/app/api/v1/` | `backend/app/api/v1/*.py` |
| Add page | Create in `frontend/src/app/` | `frontend/src/app/*/page.tsx` |
| Commit | `git commit -m "feat: description"` | — |

---

> **Remember:**
> 1. Read `AGENTS.md` before every coding session
> 2. Read `PLAN-v1.2.md` to know what to build
> 3. Follow the scaffold patterns exactly
> 4. Test before committing
> 5. Never use mock data, never break the architecture locks
>
> **Now go build something amazing.**

---

*Document Version: 1.0*  
*For: DClaw Stack Agent Swarm*  
*Stack: FastAPI + SQLAlchemy 2.0 + Next.js 14 + PostgreSQL + Docker + Helm*  
*Apps: 72 vertical SaaS applications*  
*Philosophy: Velocity through Constraint*
