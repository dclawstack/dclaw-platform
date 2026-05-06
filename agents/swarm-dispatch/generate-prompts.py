#!/usr/bin/env python3
"""Generate per-app v1.0 swarm prompts for all 67 DClaw apps."""

import json
import csv
from pathlib import Path

# All 67 DClaw apps with metadata
APPS = [
    # P0 — Active
    {"id": "chat", "name": "DClaw Chat", "category": "Communication", "tagline": "AI conversations that remember", "color": "#3B82F6", "status": "P0 Active", "priority": 0, "icon": "💬"},
    # P1 — Queued
    {"id": "flow", "name": "DClaw Flow", "category": "Automation", "tagline": "Connect anything, automate everything", "color": "#10B981", "status": "P1 Queued", "priority": 1, "icon": "🌊"},
    {"id": "agent", "name": "DClaw Agent", "category": "Platform", "tagline": "Build, share, and sell AI agents", "color": "#8B5CF6", "status": "P1 Queued", "priority": 1, "icon": "🤖"},
    {"id": "rag", "name": "DClaw RAG", "category": "Platform", "tagline": "Universal knowledge retrieval", "color": "#F59E0B", "status": "P1 Queued", "priority": 1, "icon": "🔍"},
    # P2 — Queued
    {"id": "med", "name": "DClaw Med", "category": "Healthcare", "tagline": "Clinical intelligence at your fingertips", "color": "#EF4444", "status": "P2 Queued", "priority": 2, "icon": "🏥"},
    {"id": "learn", "name": "DClaw Learn", "category": "Education", "tagline": "Adaptive learning that works", "color": "#6366F1", "status": "P2 Queued", "priority": 2, "icon": "📚"},
    {"id": "code", "name": "DClaw Code", "category": "Development", "tagline": "AI-native IDE inside your desktop", "color": "#1F2937", "status": "P2 Queued", "priority": 2, "icon": "💻"},
    # P3 — Queued
    {"id": "seo", "name": "DClaw SEO", "category": "Marketing", "tagline": "Rank higher with AI", "color": "#F97316", "status": "P3 Queued", "priority": 3, "icon": "📈"},
    {"id": "create", "name": "DClaw Create", "category": "Media", "tagline": "Generate anything", "color": "#EC4899", "status": "P3 Queued", "priority": 3, "icon": "🎨"},
    # P4+ — Planned
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


def generate_prompt(app: dict, idx: int) -> str:
    """Fill the APP_BUILDER_SWARM_PROMPT template for a single app."""
    app_id = app["id"]
    name = app["name"]
    category = app["category"]
    tagline = app["tagline"]
    color = app["color"]
    status = app["status"]
    icon = app["icon"]
    priority = app["priority"]

    # Derive ports and DB name
    frontend_port = 3000 + idx
    backend_port = 8000 + idx
    db_name = f"dclaw_{app_id}"

    # Color variations
    from colorsys import rgb_to_hls, hls_to_rgb
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

    return f"""# {name} v1.0 — Swarm Agent Prompt

> **Agent Role:** Product Architect + Full-Stack Engineer + Design Systems Specialist
> **Mission:** Take {name} from scaffold (v0.1.0) to production-ready v1.0 by researching the market, designing future-proof AI-native features, and producing a complete build specification.
> **Output:** A single `v1.0-spec.md` document that serves as the blueprint for implementation.
> **Priority:** {status}

---

## Context You Are Given

### App Identity
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

### Existing Scaffold
The app currently has:
- **Frontend:** Next.js 14.2.28, Tailwind CSS, standalone output, dark theme
- **Backend:** FastAPI, pydantic-settings, hatchling, SQLAlchemy 2.0 ready
- **Database:** PostgreSQL via CloudNativePG
- **Helm:** Standard DClaw chart with deployment, service, ingress, HPA
- **Docs:** `docs/` directory with getting-started, guides, reference, troubleshooting, releases
- **Mock API:** 2 endpoints (POST create, GET detail) with fake data

### Stack Constraints (Non-Negotiable)
| Layer | Technology | Constraint |
|-------|------------|------------|
| Frontend | Next.js 14 App Router | Must use App Router, not Pages Router |
| Frontend | Tailwind CSS | Must use Tailwind v3+ with custom design tokens |
| Frontend | React | Functional components + hooks only |
| Backend | FastAPI | pydantic v2, async routes, SQLAlchemy 2.0 |
| Backend | Python | 3.11+, type hints on all public APIs |
| Database | PostgreSQL | CloudNativePG for K8s, asyncpg for backend |
| Desktop | Tauri v2 | Optional — only if desktop adds real value |
| AI | Ollama | Local-first. Cloud fallback via OpenRouter |
| PII | ClawShield | Must call Shield before any external API |
| Auth | Logto | JWT tokens, RBAC: Owner/Admin/Developer/User/Guest |
| Packaging | Helm 3 | Standard DClaw chart patterns |

---

## Phase 1 — Market Intelligence & Gap Analysis

### 1.1 Identify Top 5 Competitors
Research the {category} category and identify the 5 most relevant competitors to {name}. For each:
- Name + URL
- Core value proposition
- Pricing model
- Key features (top 10)
- Tech stack (if public)
- User sentiment (G2, Capterra, Reddit, HN)

### 1.2 Feature Matrix
Build a comparative feature matrix:

| Feature | {name} (v0.1) | Competitor A | Competitor B | Competitor C | Competitor D | Competitor E |
|---------|---------------------------|--------------|--------------|--------------|--------------|--------------|
| Feature 1 | ❌ | ✅ | ✅ | ❌ | ✅ | ✅ |
| Feature 2 | ... | ... | ... | ... | ... | ... |

### 1.3 Gap Analysis
Identify the critical gaps:
1. **Missing table-stakes features** — What do users expect that we don't have?
2. **Differentiation opportunities** — What can we do that no competitor does?
3. **AI-native gaps** — Where are competitors still using traditional workflows instead of AI?
4. **Integration gaps** — What third-party tools do users expect integrations with?
5. **UX friction points** — Where do competitors have clunky UX we can improve?

### 1.4 Trends Analysis
Research 2026-2027 trends in the {category} category:
- Emerging AI capabilities (agents, multi-modal, reasoning)
- UX patterns (conversational UI, inline AI, predictive actions)
- Architecture trends (edge computing, local LLMs, federated learning)
- Business model trends (usage-based, seat-based, outcome-based)

---

## Phase 2 — Feature Design (v1.0)

### 2.1 Feature Tiers

#### P0 — Must Have (Launch Blockers)
Features that are table stakes. Without these, {name} is not viable.

| # | Feature | User Story | AI Integration? | Effort |
|---|---------|-----------|-----------------|--------|
| 1 | ... | As a ..., I want ... so that ... | Yes/No | S/M/L |

#### P1 — Differentiators (Launch with)
Features that set DClaw apart from competitors.

| # | Feature | User Story | AI Integration? | Effort |
|---|---------|-----------|-----------------|--------|
| 1 | ... | ... | ... | ... |

#### P2 — Nice to Have (Post-Launch)
Features for v1.1 or v1.2.

| # | Feature | User Story | AI Integration? | Effort |
|---|---------|-----------|-----------------|--------|
| 1 | ... | ... | ... | ... |

### 2.2 AI-Native Design Principles
Every feature must be evaluated against these principles:

1. **Predictive over Reactive** — Can the AI anticipate what the user needs before they ask?
2. **Generative over Manual** — Can the AI generate content/structure instead of requiring manual input?
3. **Conversational over Form-Based** — Can natural language replace complex forms?
4. **Contextual over Stateless** — Does the AI remember previous interactions and project context?
5. **Local over Cloud** — Can this run on Ollama before falling back to OpenRouter?

### 2.3 Feature Specifications
For each P0 and P1 feature, write:

```markdown
### Feature: {{Feature Name}}

**Priority:** P0/P1
**Effort:** S/M/L
**AI-Native Score:** 1-5 (how deeply AI is integrated)

#### User Flow
1. User opens ...
2. System shows ...
3. User interacts ...
4. AI does ...
5. Result is ...

#### UI/UX Design
- Screen layout (describe key zones)
- Micro-interactions (hover states, loading, transitions)
- Error states and empty states
- Mobile adaptation

#### API Contract
```
POST /api/v1/...
Request: {{...}}
Response: {{...}}
```

#### Data Model
```python
class FeatureModel(BaseModel):
    id: UUID
    ...
```

#### AI Prompt (if applicable)
```
System: You are a ...
User: {{context}}
Assistant: ...
```

#### Edge Cases
- What happens when ...?
- How do we handle ...?
```

---

## Phase 3 — Architecture Design

### 3.1 System Architecture Diagram
Create a Mermaid diagram showing:
- Frontend (Next.js App Router) with key routes
- Backend (FastAPI) with routers, services, repositories
- Database schema (ER diagram)
- AI layer (Ollama / OpenRouter / Shield)
- External integrations
- Message bus / cache (Redis)
- File storage (MinIO) if needed

### 3.2 Database Schema
Write the complete SQLAlchemy models:

```python
# models.py
from sqlalchemy import String, DateTime, ForeignKey, JSON, Integer, Float, Boolean, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship
from datetime import datetime
from uuid import uuid4

class User(Base):
    __tablename__ = "users"
    id: Mapped[str] = mapped_column(String(36), primary_key=True, default=lambda: str(uuid4()))
    ...
```

### 3.3 API Design
Full OpenAPI-style spec for all endpoints:

```yaml
POST /api/v1/resources
  summary: Create resource
  request: {{...}}
  response: {{...}}
  errors: [...]
```

### 3.4 Key Code Snippets
Provide exact, production-ready code for the 5 most critical functions:

1. **Main service function** (the core business logic)
2. **AI integration function** (Ollama/OpenRouter call with Shield)
3. **Real-time update function** (WebSocket/SSE if needed)
4. **Background job function** (if using Celery/ARQ)
5. **Auth middleware function** (Logto JWT validation)

Each snippet must be:
- Fully typed
- Error-handled
- Logged
- Documented with docstrings
- Follow DClaw conventions

---

## Phase 4 — Design System Application

Apply the **DKube Design System** to {name}:

### 4.1 Color Adaptation
Map the app's brand color to the DKube palette:

| Token | Value | Usage |
|-------|-------|-------|
| `--dk-purple` | {color} | Primary actions, links |
| `--dk-purple-light` | {light} | Hover states, accents |
| `--dk-purple-deep` | {deep} | Pressed states |
| `--dk-purple-wash` | {wash} | Section backgrounds |

### 4.2 Typography
- Display font: Manrope (headings)
- Body font: Inter (body text)
- Mono font: JetBrains Mono (code, data)

### 4.3 Component Specifications
For each custom component, specify:
- Props interface
- Default values
- Responsive behavior
- Accessibility (ARIA labels, keyboard navigation, focus rings)
- Animation (enter/exit, hover, loading)

### 4.4 Dark Mode
{name} supports both light and dark modes. Define the dark mode color mapping:

```css
.dark {{
  --dk-surface: #111013;
  --dk-body: #F4F2F8;
  --dk-muted: #9E9AAB;
  ...
}}
```

---

## Phase 5 — Testing Strategy

### 5.1 Test Pyramid

| Layer | Framework | Coverage Target | What to Test |
|-------|-----------|-----------------|--------------|
| Unit | pytest | 80%+ | Services, utilities, models |
| Integration | pytest + httpx | 70%+ | API endpoints, DB transactions |
| E2E | Playwright | Critical paths | User flows, AI interactions |
| Contract | schemathesis | 100% | OpenAPI spec compliance |
| Load | locust | Key endpoints | 100 concurrent users |
| Security | bandit + safety | All code | Vulnerabilities, dependency audits |

### 5.2 Critical Test Cases
Write 10 specific test cases for the most important features:

```python
def test_ai_generates_valid_response():
    '''Given context X, AI should return structured output Y.'''
    ...
```

### 5.3 Test Data Strategy
- Factory Boy / Faker for generating test data
- Snapshot testing for AI outputs (within tolerance)
- Seeded random for reproducible AI tests

---

## Phase 6 — Roadmap & Lifecycle

### 6.1 v1.0 Release Checklist
- [ ] All P0 features implemented
- [ ] Test coverage > 80%
- [ ] Performance benchmarks pass
- [ ] Security audit clean
- [ ] Docs complete
- [ ] Helm chart tested
- [ ] Desktop build signed (if applicable)

### 6.2 v1.1 — Quick Wins (4 weeks post-launch)
- Feature A
- Feature B
- Performance optimization

### 6.3 v1.2 — AI Enhancements (8 weeks post-launch)
- Agent integration
- Multi-modal support
- Advanced reasoning

### 6.4 v2.0 — Platform (Q1 2027)
- Plugin ecosystem
- Third-party integrations
- Enterprise features

### 6.5 Deprecation Policy
- Features marked deprecated receive 6 months warning
- API versions supported for 12 months
- Migration guides provided for breaking changes

---

## Phase 7 — Documentation Update

Update the app's `docs/` directory with v1.0 content:

### Getting Started
- `installation.md` — Step-by-step with screenshots
- `quickstart.md` — 5-minute tutorial
- `configuration.md` — All env vars explained

### Guides
- `use-cases.md` — 5 real-world scenarios with walkthroughs
- `best-practices.md` — Security, performance, cost optimization
- `ai-prompts.md` — How to write effective prompts for this app

### Reference
- `architecture.md` — Updated system diagram
- `stack.md` — Complete dependency list with versions
- `api.md` — Full OpenAPI spec rendered

### Troubleshooting
- `common-issues.md` — 10 most common problems
- `faq.md` — 20 questions

### Releases
- `changelog.md` — v0.1.0 → v1.0.0 changes
- `roadmap.md` — v1.1, v1.2, v2.0

---

## Output Format

Produce a single file: `v1.0-spec.md`

Structure:
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


def main():
    base = Path(__file__).parent
    prompts_dir = base / "prompts"
    prompts_dir.mkdir(exist_ok=True)

    manifest = []
    dispatch = []

    for idx, app in enumerate(APPS, start=1):
        app_id = app["id"]
        prompt_text = generate_prompt(app, idx)
        prompt_path = prompts_dir / f"{app_id}-v1.0-prompt.md"
        prompt_path.write_text(prompt_text, encoding="utf-8")

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
            "prompt_file": str(prompt_path.relative_to(base)),
            "frontend_port": 3000 + idx,
            "backend_port": 8000 + idx,
            "db_name": f"dclaw_{app_id}",
        })

        dispatch.append({
            "app_id": app_id,
            "name": app["name"],
            "priority": app["priority"],
            "repo_url": f"https://github.com/dclawstack/dclaw-{app_id}",
            "prompt_path": str(prompt_path.relative_to(base)),
            "status": app["status"],
        })

        print(f"✅ {app_id}")

    # Write manifest
    manifest_path = base / "apps-manifest.json"
    manifest_path.write_text(json.dumps(manifest, indent=2), encoding="utf-8")

    # Write dispatch CSV
    dispatch_path = base / "dispatch.csv"
    with dispatch_path.open("w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=["app_id", "name", "priority", "repo_url", "prompt_path", "status"])
        writer.writeheader()
        writer.writerows(dispatch)

    print(f"\nDone. {len(APPS)} prompts generated.")
    print(f"Manifest: {manifest_path}")
    print(f"Dispatch: {dispatch_path}")


if __name__ == "__main__":
    main()
