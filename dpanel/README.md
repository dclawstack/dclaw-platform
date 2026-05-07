# DPanel

> **Your AI app store and command center.**
>
> DPanel is the frontend for the DClaw Platform — a Next.js 16 app that lets you browse, install, and manage all 65+ DClaw apps. It also serves as the documentation hub for the entire ecosystem.

---

## Features

| Feature | Description |
|---------|-------------|
| **App Store** | Browse 65+ apps with search, filters, and categories |
| **Install Manager** | One-click install apps to your K8s cluster |
| **Docs Hub** | Per-app docs + ecosystem docs with custom markdown renderer |
| **Persistent Sidebar** | Collapsible nav with all apps grouped by category |
| **Dark Mode** | DKube design system — deep purples on near-black |
| **Mobile Responsive** | Hamburger menu, touch-friendly cards |

---

## Tech Stack

- [Next.js](https://nextjs.org/) 16.2.4 — App Router, Turbopack
- [Tailwind CSS](https://tailwindcss.com/) v4 — `@theme inline` syntax
- [React](https://react.dev/) 19 — Server Components
- [TypeScript](https://www.typescriptlang.org/) — Full type coverage

---

## Design System

DPanel uses **DKube** — the shared design system for all DClaw apps.

```css
/* globals.css */
@import "tailwindcss";

@theme inline {
  --color-dk-purple: #6B53A3;
  --color-dk-surface: #0E0E10;
  --color-dk-body: #F4F2F8;
  /* ... see ../design-system/dkube.css */
}
```

**Fonts:**
- Manrope (display headings)
- Inter (body text)
- JetBrains Mono (code)

---

## Development

### Prerequisites

- Node.js 20+
- npm or pnpm

### Install

```bash
cd dclaw-platform/dpanel
npm install
```

### Run

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

### Build

```bash
npm run build
```

### Lint

```bash
npm run lint
```

---

## Project Structure

```
dpanel/
├── app/                       # Next.js App Router
│   ├── (docs)/               # Docs layout group
│   │   ├── docs/             # Ecosystem docs
│   │   └── docs/apps/[id]/   # Per-app docs
│   ├── app/[id]/             # App detail pages
│   ├── page.tsx              # Home / app grid
│   └── layout.tsx            # Root layout with fonts
├── components/               # React components
│   ├── sidebar.tsx           # Persistent app nav
│   ├── app-card.tsx          # Store card component
│   ├── install-button.tsx    # K8s install trigger
│   └── ui/                   # shadcn/ui primitives
├── lib/
│   ├── apps.ts               # App registry (65 apps)
│   └── utils.ts              # cn() and helpers
├── docs-content/             # Ecosystem markdown
│   ├── ecosystem/            # Getting started, architecture
│   └── apps/                 # Per-app doc stubs
└── public/                   # Static assets
```

---

## App Registry

Apps are defined in `lib/apps.ts`. Each app has:

```ts
{
  id: "chat",
  name: "DClaw Chat",
  tagline: "AI conversations that remember",
  description: "...",
  features: ["..."],
  version: "1.0.0",
  icon: MessageSquare,
  color: "#3B82F6",
  bgColor: "rgba(59, 130, 246, 0.15)",
  domain: "chat.dclawstack.io",
  category: "Communication",
  status: "live",
}
```

### Categories

- Communication
- Automation
- Platform
- Healthcare
- Education
- Development
- Marketing
- Media
- Legal
- Finance
- Sales
- Support
- HR
- Design
- Productivity
- Knowledge
- Infrastructure
- Security
- Governance
- Operations
- Logistics
- Utilities
- Real Estate
- Sustainability

---

## Docs System

### Ecosystem Docs

Located in `docs-content/ecosystem/`:

| Doc | Path |
|-----|------|
| Getting Started | `/docs/ecosystem/getting-started` |
| Architecture | `/docs/ecosystem/architecture` |
| Stack Reference | `/docs/ecosystem/stack-reference` |
| Troubleshooting | `/docs/ecosystem/troubleshooting` |
| Releases | `/docs/ecosystem/releases` |

### Per-App Docs

Each app repo has a `docs/` directory with:

```
docs/
├── meta.json              # Nav structure
├── getting-started/       # Installation, quickstart, config
├── guides/               # Use cases, best practices
├── reference/            # Architecture, stack, API
├── troubleshooting/      # Common issues, FAQ
└── releases/             # Changelog, roadmap
```

DPanel renders these dynamically via the docs API.

---

## Deployment

### Standalone

```bash
npm run build
npm start
```

### Docker

```bash
docker build -t dpanel .
docker run -p 3000:3000 dpanel
```

### Kubernetes

```bash
kubectl apply -f k8s/
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | DPanel API (Go server) |
| `NEXT_PUBLIC_LOGTO_ENDPOINT` | — | Auth server |

---

## License

MIT
