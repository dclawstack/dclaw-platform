# DClaw Onboard v1.0 — Swarm Agent Prompt

> **Agent Role:** Product Architect + Full-Stack Engineer + Design Systems Specialist
> **Mission:** Take DClaw Onboard from scaffold (v0.1.0) to production-ready v1.0 by researching the market, designing future-proof AI-native features, and producing a complete build specification.
> **Output:** A single `v1.0-spec.md` document that serves as the blueprint for implementation.
> **Priority:** P4 Planned

---

## Context You Are Given

### App Identity
```yaml
app_id: onboard
name: DClaw Onboard
category: HR
tagline: Automated employee onboarding
current_version: 0.1.0
status: P4 Planned
repo: https://github.com/dclawstack/dclaw-onboard
docs: https://docs.dclawstack.io/apps/onboard
frontend_port: 3041
backend_port: 8041
db_name: dclaw_onboard
brand_color: #10B981
icon: 🚀
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
Research the HR category and identify the 5 most relevant competitors to DClaw Onboard. For each:
- Name + URL
- Core value proposition
- Pricing model
- Key features (top 10)
- Tech stack (if public)
- User sentiment (G2, Capterra, Reddit, HN)

### 1.2 Feature Matrix
Build a comparative feature matrix:

| Feature | DClaw Onboard (v0.1) | Competitor A | Competitor B | Competitor C | Competitor D | Competitor E |
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
Research 2026-2027 trends in the HR category:
- Emerging AI capabilities (agents, multi-modal, reasoning)
- UX patterns (conversational UI, inline AI, predictive actions)
- Architecture trends (edge computing, local LLMs, federated learning)
- Business model trends (usage-based, seat-based, outcome-based)

---

## Phase 2 — Feature Design (v1.0)

### 2.1 Feature Tiers

#### P0 — Must Have (Launch Blockers)
Features that are table stakes. Without these, DClaw Onboard is not viable.

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
### Feature: {Feature Name}

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
Request: {...}
Response: {...}
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
User: {context}
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
  request: {...}
  response: {...}
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

Apply the **DKube Design System** to DClaw Onboard:

### 4.1 Color Adaptation
Map the app's brand color to the DKube palette:

| Token | Value | Usage |
|-------|-------|-------|
| `--dk-purple` | #10B981 | Primary actions, links |
| `--dk-purple-light` | #29ECAB | Hover states, accents |
| `--dk-purple-deep` | #09724F | Pressed states |
| `--dk-purple-wash` | #F2F7F5 | Section backgrounds |

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
DClaw Onboard supports both light and dark modes. Define the dark mode color mapping:

```css
.dark {
  --dk-surface: #111013;
  --dk-body: #F4F2F8;
  --dk-muted: #9E9AAB;
  ...
}
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
# DClaw Onboard v1.0 Specification

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
