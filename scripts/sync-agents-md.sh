#!/usr/bin/env bash
# Render per-app AGENTS.md from the canonical template + commit + push.
#
# This is the bulk-apply pass for unrendered AGENTS.md files across the
# dclaw-* app repos. Renders from
#   dclaw-platform/scaffold/AGENTS.template.md (template)
#   dclaw-platform/agents/swarm-dispatch/apps-manifest.json (per-app vars)
# Substitutes {APP_NAME}, {BACKEND_PORT}, {FRONTEND_PORT}, {DB_NAME}.
#
# Skips repos whose AGENTS.md is already personalized (no placeholders).
#
# Usage:
#   ./sync-agents-md.sh                  # dry-run
#   ./sync-agents-md.sh --apply          # render + commit + push
#   ./sync-agents-md.sh --apply --apps chat,flow

set -euo pipefail

HUB_ROOT="${HUB_ROOT:-$HOME/DClaw-Stack}"
PLATFORM_DIR="$HUB_ROOT/dclaw-platform"
RENDERER="$PLATFORM_DIR/scripts/render-agents-md.py"

DRY_RUN=true
APPS=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) DRY_RUN=false; shift ;;
    --apps)  APPS="$2"; shift 2 ;;
    -h|--help)
      sed -n '2,/^set -/p' "$0" | sed -E 's/^# ?//;/^set -/d'
      exit 0 ;;
    *) echo "unknown arg: $1" >&2; exit 2 ;;
  esac
done

# Step 1: render in-place
APPS_ARG=""
[ -n "$APPS" ] && APPS_ARG="--apps $APPS"
if $DRY_RUN; then
  python3 "$RENDERER" $APPS_ARG
  echo
  echo "(re-run with --apply to write files, commit, and push)"
  exit 0
fi

python3 "$RENDERER" --apply $APPS_ARG

# Step 2: commit + push for each repo whose AGENTS.md was modified
echo
echo "── committing + pushing ──"
ok=0; fail=0; nochange=0
for repo_path in "$HUB_ROOT"/dclaw-*/; do
  repo=$(basename "$repo_path")
  [ -d "$repo_path/.git" ] || continue

  # Filter by --apps if provided
  if [ -n "$APPS" ]; then
    app_id="${repo#dclaw-}"
    if ! echo ",$APPS," | grep -q ",$app_id,"; then continue; fi
  fi

  # Only AGENTS.md should be dirty if anything is
  status=$(/usr/bin/git -C "$repo_path" status --porcelain 2>/dev/null)
  if [ -z "$status" ]; then
    ((nochange++))
    continue
  fi

  # Refuse to touch a repo with unrelated dirty files
  unrelated=$(echo "$status" | grep -v "AGENTS.md" || true)
  if [ -n "$unrelated" ]; then
    printf '  ⚠  %s — has unrelated changes, skipping\n' "$repo"
    ((fail++))
    continue
  fi

  /usr/bin/git -C "$repo_path" add AGENTS.md
  msg="chore(scaffold): render AGENTS.md from template

Personalized from dclaw-platform/scaffold/AGENTS.template.md using
apps-manifest.json vars (APP_NAME / BACKEND_PORT / FRONTEND_PORT / DB_NAME).
No placeholders remain."

  if ! /usr/bin/git -C "$repo_path" commit -m "$msg" >/dev/null 2>&1; then
    printf '  ✗  %s — commit failed\n' "$repo"
    ((fail++))
    continue
  fi

  branch=$(/usr/bin/git -C "$repo_path" rev-parse --abbrev-ref HEAD)
  if ! /usr/bin/git -C "$repo_path" pull --rebase origin "$branch" >/dev/null 2>&1; then
    printf '  ⚠  %s — committed locally; rebase failed\n' "$repo"
    ((fail++))
    continue
  fi
  if /usr/bin/git -C "$repo_path" push origin "$branch" >/dev/null 2>&1; then
    printf '  ✓  %s — rendered + pushed\n' "$repo"
    ((ok++))
  else
    printf '  ⚠  %s — push failed\n' "$repo"
    ((fail++))
  fi
done

echo
echo "── summary ──"
printf '   ok       : %d\n' "$ok"
printf '   unchanged: %d\n' "$nochange"
printf '   failed   : %d\n' "$fail"
