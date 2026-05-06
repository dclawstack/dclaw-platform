# Changelog

## v0.1.0 — May 2026

### Added
- DClaw Operator with 10-step reconciliation pipeline
- DPanel Next.js app with app launcher grid
- DClaw Chat MVP (multi-model chat, persistent history, voice input)
- Agent Swarm runtime (registry, orchestrator, 5 agents)
- Helm charts for all apps
- CloudNativePG integration for per-app databases
- PII Shield (local anonymization layer)
- Tauri v2 desktop shell (unsigned builds)
- Telegram CI pipeline (TEST-01)

### Changed
- Migrated from monorepo to per-app repositories
- Upgraded to Next.js 16 and Tailwind v4

### Fixed
- Namespace race condition in operator
- Database connection pooling in FastAPI backends

## Upcoming

### v0.2.0 — June 2026 (Planned)
- DClaw Flow visual workflow builder
- DClaw RAG universal knowledge retrieval
- DClaw Agent marketplace
- DPanel docs system
- CloudNativePG cluster CR support

### v0.3.0 — July 2026 (Planned)
- DClaw Med clinical intelligence
- DClaw Learn adaptive learning
- DClaw Code AI-native IDE
- Desktop auto-updater
- Voice wake word prototype

### v1.0.0 — September 2026 (Planned)
- All 65 apps scaffolded or in beta
- Enterprise SSO and audit logs
- White-label builds
- YC Demo Day
