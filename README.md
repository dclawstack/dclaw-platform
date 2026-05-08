# DClaw Platform

> **The AI app store and command center for the DClaw Stack.**
>
> DPanel is your central hub for discovering, installing, and managing all 65+ DClaw apps. Built as a Next.js 16 app with the DKube Design System.

---

## What is DClaw?

DClaw is an AI-native app platform — think **Adobe Creative Cloud meets AI agents**. Instead of subscribing to 20 different SaaS tools, you get a unified platform of 65+ AI-powered apps that share a common design system, auth layer, and infrastructure.

| Layer | Technology |
|-------|-----------|
| **Frontend** | Next.js 16, Tailwind CSS v4, React Server Components |
| **Design System** | DKube — purple brand, dark neutrals, Manrope/Inter/JetBrains Mono |
| **App Store** | DPanel — browse, install, and launch apps |
| **Docs Hub** | Per-app docs + ecosystem docs with custom markdown renderer |
| **API** | Go HTTP server serving app registry from Kubernetes ConfigMap |
| **Auth** | Logto — JWT tokens, RBAC (Owner/Admin/Developer/User/Guest) |

---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                         DPanel                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Home    │  │ App Store│  │  Docs    │  │ Settings │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                              │
│  DClaw Chat  DClaw Flow  DClaw Agent  DClaw RAG  ...       │
│     🔴Live     ⏳Soon      ⏳Soon      ⏳Soon                │
└──────────────────────────┬───────────────────────────────────┘
                           │
              ┌────────────┼────────────┐
              ↓            ↓            ↓
        ┌─────────┐  ┌─────────┐  ┌─────────┐
        │ dpanel- │  │ dclaw-  │  │ dclaw-  │
        │  api    │  │  chat   │  │  flow   │
        │  (Go)   │  │(Next.js)│  │(Next.js)│
        └─────────┘  └─────────┘  └─────────┘
              │            │            │
              └────────────┼────────────┘
                           ↓
                    ┌─────────────┐
                    │  Kubernetes │
                    │  ConfigMap  │
                    │(app registry)│
                    └─────────────┘
```

---

## Directory Structure

```
dclaw-platform/
├── dpanel/                    # Next.js 16 app store + docs hub
│   ├── app/                   # App Router pages
│   ├── components/            # Sidebar, app cards, docs renderer
│   ├── lib/apps.ts            # App registry (65 apps)
│   └── docs-content/          # Ecosystem markdown docs
├── dpanel-api/                # Go HTTP server
│   ├── main.go                # App registry API
│   └── docs-openapi.yaml      # Docs API spec
├── design-system/             # DKube shared tokens
│   ├── dkube.css              # CSS variables
│   └── tailwind-preset.js     # Tailwind v3 preset
├── agents/
│   ├── swarm-prompts/         # APP_BUILDER + IMPLEMENTER prompts
│   └── swarm-dispatch/        # Per-app prompts + manifest
└── helm/                      # Platform K8s charts
```

---

## DKube Design System

All DClaw apps share the **DKube Design System**.

### Core Tokens

| Token | Value |
|-------|-------|
| Brand Purple | `#6B53A3` |
| Purple Light | `#9985BF` |
| Purple Deep | `#4A3A7A` |
| Surface | `#0E0E10` |
| Surface Raised | `#1F1F23` |
| Body Text | `#F4F2F8` |
| Muted | `#9E9AAB` |

### Fonts
- **Display:** Manrope
- **Body:** Inter
- **Mono:** JetBrains Mono

### Usage in Apps

**Tailwind v3 apps:**
```ts
// tailwind.config.ts
import dkubePreset from "../dclaw-platform/design-system/tailwind-preset";
export default { presets: [dkubePreset] };
```

**Tailwind v4 apps (DPanel):**
```css
@import "tailwindcss";
@theme inline {
  --color-dk-purple: #6B53A3;
  --color-dk-surface: #0E0E10;
  /* ... see design-system/dkube.css */
}
```

---

## Swarm Agent Dispatch

The platform includes a **swarm dispatch system** for AI agents to build apps.

### Prompts Generated

| Type | Count | Purpose |
|------|-------|---------|
| Spec Writer | 67 | Research market, write v1.0 spec |
| Implementer | 67 | Build the app from scaffold to working code |

### Key Files

| File | Purpose |
|------|---------|
| `agents/swarm-prompts/APP_BUILDER_SWARM_PROMPT.md` | Master spec-writing template |
| `agents/swarm-prompts/IMPLEMENTER_SWARM_PROMPT.md` | Master implementation template |
| `agents/swarm-dispatch/apps-manifest.json` | All 67 apps with metadata |
| `agents/swarm-dispatch/dispatch.csv` | Machine-readable dispatch map |
| `agents/swarm-dispatch/run-swarm.sh` | CLI to generate agent commands |

### Usage

```bash
cd agents/swarm-dispatch

# Print command for a single app
./run-swarm.sh chat

# Print all P0 apps
./run-swarm.sh --priority 0

# Print first 5 apps
./run-swarm.sh --batch 5
```

---

## App Lifecycle on DPanel

```
Scaffolded → Spec Written → Building → Alpha → Beta → Live → Deprecated
```

### Status Badges

| Status | Badge | Meaning |
|--------|-------|---------|
| Live | 🔴 | Deployed and working |
| Alpha | 🟡 | Feature-complete, testing |
| Beta | 🟠 | Stabilizing, docs ready |
| Soon | ⚪ | Scaffolded, not built |

### Adding an App to DPanel

Edit `dpanel/lib/apps.ts`:

```ts
{
  id: "your-app",
  name: "DClaw YourApp",
  tagline: "What it does",
  description: "Longer description...",
  features: ["Feature 1", "Feature 2"],
  version: "1.0.0",
  icon: YourIcon,
  color: "#3B82F6",
  bgColor: "rgba(59, 130, 246, 0.15)",
  domain: "yourapp.dclawstack.io",
  category: "Category",
  status: "live",
}
```

---

## Development

### DPanel (Next.js 16)

```bash
cd dpanel
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

### DPanel API (Go)

```bash
cd dpanel-api
go run main.go
```

API runs on `http://localhost:8080`.

### Build

```bash
cd dpanel
npm run build
```

---

## Deployment

### Docker Compose (Recommended for Self-Hosted)

```bash
docker compose up -d
```

### Kubernetes

```bash
cd helm/dclaw-platform
helm dependency build
helm upgrade --install dclaw-platform . \
  --namespace dclaw-platform \
  --create-namespace
```

---

## The 65 Apps

### P0 — Active

| App | Status | Description |
|-----|--------|-------------|
| [Chat](https://github.com/dclawstack/dclaw-chat) | 🔴 Live | AI conversations that remember |

### P1 — Platform

| App | Status | Description |
|-----|--------|-------------|
| Flow | ⏳ Soon | Visual workflow builder |
| Agent | ⏳ Soon | Build, share, and sell AI agents |
| RAG | ⏳ Soon | Universal knowledge retrieval |

### P2 — Verticals

| App | Status | Description |
|-----|--------|-------------|
| Med | ⏳ Soon | Clinical intelligence |
| Learn | ⏳ Soon | Adaptive learning platform |
| Code | ⏳ Soon | AI-native IDE |

### P3 — Scale

| App | Status | Description |
|-----|--------|-------------|
| SEO | ⏳ Soon | SEO toolkit |
| Create | ⏳ Soon | Generative media studio |

### P4+ — Planned

Legal, Finance, Sales, Support, HR, Design, Translate, Write, Meet, Doc, Sheet, Slide, Mail, Calendar, Task, Wiki, Data, API, Test, Deploy, Monitor, Secure, Backup, Migrate, Cost, Carbon, Compliance, Audit, Policy, Train, Recruit, Onboard, Offboard, Assets, Network, Inventory, Forecast, Quality, Maintenance, Route, Warehouse, Fleet, Energy, Water, Waste, Building, Space, Lease, Vendor, Contract, Risk, Crisis, Continuity, Knowledge, Research, Patent, Trademark, Video

---

## Contributing

1. Fork the repo
2. Create a branch
3. Make changes
4. Open a PR

---

## License

MIT

## Contributor

Rajendra M (01.r.machani@gmail.com)
