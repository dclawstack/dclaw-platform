# DClaw Agent Swarm — Kimi Web Chat Bootstrap Guide

> How to spawn and operate the multi-agent swarm using Kimi web chat sessions.

---

## Quick Start (30 seconds per agent)

### 1. Open Kimi Web Chat
Go to https://kimi.moonshot.cn (or the Kimi web interface you use).

### 2. Create a new chat session
Each agent = one chat session. Open 2–7 tabs depending on how many agents you need.

### 3. Bootstrap Sacred Context
Paste this into **each** new chat:

```
You are an autonomous agent in the DClaw Stack multi-agent swarm. Before accepting tasks, you must ingest the sacred context. Please read and internalize the following public documents. After each one, summarize the key points in 1–2 sentences:

https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/VISION.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/ARCHITECTURE.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/PRODUCTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/STATUS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/AGENTS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/CONVENTIONS.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SETUP.md
https://raw.githubusercontent.com/dclawstack/dclaw-prd/main/SECURITY.md
https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/AGENT_SWARM.md

Confirm when complete.
```

### 4. Assign Agent Role
After the agent confirms context ingestion, paste the role-specific prompt:

```
Now read your specific agent prompt and confirm your identity:

https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/agents/VAULT_COORDINATOR.md
```

Replace `VAULT_COORDINATOR.md` with the appropriate role for each tab:
- `SHELL_AGENT.md` — implementation, builds, deployment
- `SHIELD_AGENT.md` — security, compliance, audits
- `CODE_AGENT.md` — deep coding, refactoring, reviews
- `RESEARCH_AGENT.md` — R&D, prototyping, benchmarking
- `MEMORY_AGENT.md` — documentation, knowledge graphs
- `GENERAL_AGENT.md` — user-facing router, coordinator

### 5. Give First Task
The agent is now operational. Assign work using the handoff format (see AGENT_SWARM.md §5.3).

---

## Minimal Swarm (2 agents)

For simple tasks, you only need:

1. **Vault Coordinator** — thinks, specifies, decides
2. **Shell Agent** — implements, builds, deploys

## Full Swarm (7 agents)

For complex parallel work:

| Tab | Agent | Prompt URL |
|-----|-------|------------|
| 1 | General | `agents/GENERAL_AGENT.md` |
| 2 | Vault Coordinator | `agents/VAULT_COORDINATOR.md` |
| 3 | Shell Agent | `agents/SHELL_AGENT.md` |
| 4 | Shield Agent | `agents/SHIELD_AGENT.md` |
| 5 | Code Agent | `agents/CODE_AGENT.md` |
| 6 | Research Agent | `agents/RESEARCH_AGENT.md` |
| 7 | Memory Agent | `agents/MEMORY_AGENT.md` |

---

## How Agents Coordinate

Agents do NOT chat with each other. They coordinate through **GitHub**:

1. **Vault Coordinator** writes a spec → commits to `dclaw-prd`
2. **Shell Agent** sees the commit → implements → commits to `dclaw-platform`
3. **Shell Agent** opens PR → tags `shield-agent` for review
4. **Shield Agent** reviews → approves or requests changes
5. **Shell Agent** merges → deploys
6. **Memory Agent** documents the decision → updates `dclaw-prd` and `dclaw-obsidian`

You (the human) monitor the GitHub repos and Telegram notifications to track progress.

---

## Important Notes

- **Ephemeral sessions:** Kimi web chat sessions may reset. If an agent loses context, re-paste its role prompt and say "Resume from where we left off. Check the latest GitHub commits on [repo] to synchronize."
- **No shared memory:** Each chat session is isolated. The ONLY shared memory is GitHub.
- **Commit attribution:** Always use `[agent:role]` prefix in commits so other agents know who did what.
- **Decision authority:** The Vault Coordinator owns architecture. The Shell Agent owns implementation. Disputes are resolved by you (human) or by opening a `decision-needed` issue.

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Agent forgot everything | Re-paste sacred context + role prompt |
| Two agents editing same file | Check AGENT_SWARM.md repo ownership matrix. One agent must yield. |
| Agent unsure what to do next | Tell it to `git pull` and read `STATUS.md` |
| Agent wants to push secrets | Stop it. Secrets go in GitHub Settings → Secrets only. |
| Agent made a bad decision | Open a revert PR, tag the correct agent, update DECISIONS.md |
