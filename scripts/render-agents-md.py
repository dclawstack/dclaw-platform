#!/usr/bin/env python3
"""Render per-app AGENTS.md from the canonical template + apps-manifest.json.

The template at `dclaw-platform/scaffold/AGENTS.template.md` contains
placeholders `{APP_NAME}`, `{BACKEND_PORT}`, `{FRONTEND_PORT}`, `{DB_NAME}`.
This script substitutes those per-app from the swarm-dispatch manifest.

By default it operates as dry-run — prints which repos would change.
Use --apply to write rendered files (does not commit; pair with the
sibling sync-agents-md.sh for git workflow).

Usage:
    render-agents-md.py                  # dry-run, list changes
    render-agents-md.py --apply          # write rendered AGENTS.md
    render-agents-md.py --apps chat,flow # limit to specific apps
    render-agents-md.py --force          # render even if no placeholders found
"""

import argparse
import json
import os
import re
import sys
from pathlib import Path

HUB_ROOT = Path(os.path.expanduser("~/DClaw-Stack"))
TEMPLATE = HUB_ROOT / "dclaw-platform/scaffold/AGENTS.template.md"
MANIFEST = HUB_ROOT / "dclaw-platform/agents/swarm-dispatch/apps-manifest.json"
PLACEHOLDER_RE = re.compile(r"\{[A-Z_]+\}")


def short_name(full: str) -> str:
    """'DClaw Chat' -> 'Chat'; fallback to the original string."""
    return full.removeprefix("DClaw ").strip() or full


def render(template: str, entry: dict) -> str:
    """Substitute {APP_NAME}/{BACKEND_PORT}/{FRONTEND_PORT}/{DB_NAME}."""
    out = template
    subs = {
        "{APP_NAME}": short_name(entry["name"]),
        "{BACKEND_PORT}": str(entry["backend_port"]),
        "{FRONTEND_PORT}": str(entry["frontend_port"]),
        "{DB_NAME}": entry["db_name"],
    }
    for k, v in subs.items():
        out = out.replace(k, v)
    return out


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--apply", action="store_true", help="actually write files")
    p.add_argument("--apps", help="comma-separated app_ids to limit to")
    p.add_argument("--force", action="store_true",
                   help="render even if AGENTS.md is already personalized")
    args = p.parse_args()

    if not TEMPLATE.exists():
        print(f"FAIL: template missing at {TEMPLATE}", file=sys.stderr)
        sys.exit(1)
    if not MANIFEST.exists():
        print(f"FAIL: manifest missing at {MANIFEST}", file=sys.stderr)
        sys.exit(1)

    template = TEMPLATE.read_text()
    entries = json.loads(MANIFEST.read_text())

    if args.apps:
        wanted = {a.strip() for a in args.apps.split(",")}
        entries = [e for e in entries if e["app_id"] in wanted]

    mode = "APPLY" if args.apply else "DRY-RUN"
    print(f"\n══ render-agents-md · {mode} ══")
    print(f"   template: {TEMPLATE.relative_to(HUB_ROOT)}")
    print(f"   manifest: {len(entries)} entries")
    print()

    rendered = 0
    skipped = 0
    missing = 0

    for entry in entries:
        app_id = entry["app_id"]
        repo = HUB_ROOT / f"dclaw-{app_id}"
        agents_path = repo / "AGENTS.md"

        if not repo.is_dir() or not (repo / ".git").exists():
            print(f"  -  dclaw-{app_id} — repo not present, skipping")
            missing += 1
            continue

        existing = agents_path.read_text() if agents_path.exists() else ""
        already_personalized = (
            existing and not PLACEHOLDER_RE.search(existing)
        )

        if already_personalized and not args.force:
            print(f"  ·  dclaw-{app_id} — already personalized, skipping")
            skipped += 1
            continue

        new_content = render(template, entry)
        if existing == new_content:
            print(f"  ·  dclaw-{app_id} — already at rendered content")
            skipped += 1
            continue

        if args.apply:
            agents_path.write_text(new_content)
            print(f"  ✓  dclaw-{app_id} — rendered ({len(new_content):,} bytes)")
        else:
            print(f"  →  dclaw-{app_id} — would render: "
                  f"name='{short_name(entry['name'])}' "
                  f"be_port={entry['backend_port']} "
                  f"fe_port={entry['frontend_port']} "
                  f"db={entry['db_name']}")
        rendered += 1

    print()
    print("── summary ──")
    print(f"   rendered : {rendered}")
    print(f"   skipped  : {skipped}")
    print(f"   missing  : {missing}")
    if not args.apply:
        print("\n(dry-run — re-run with --apply to write files; "
              "use sync-agents-md.sh to commit + push)")


if __name__ == "__main__":
    main()
