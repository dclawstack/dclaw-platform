#!/usr/bin/env bash
# Propagate claude.yml + claude-code-review.yml from dclaw-chat (canonical
# source) into every dclaw-* app repo missing them.
#
# Note: workflow files alone are not sufficient. Each target repo also
# needs the CLAUDE_CODE_OAUTH_TOKEN secret set. Repos without the secret
# will install the workflow cleanly but the action will fail at runtime
# until the secret is added (gh secret set CLAUDE_CODE_OAUTH_TOKEN -b "$TOKEN").
#
# Usage:
#   ./propagate-claude-workflows.sh                # dry-run
#   ./propagate-claude-workflows.sh --apply        # commit + push
#   ./propagate-claude-workflows.sh --apply --apps dclaw-api,dclaw-assets
#
# This is a one-shot tool. The cleaner long-term home is scaffold-sync v2
# with `--target-dir .github/workflows/` support.

set -euo pipefail

HUB_ROOT="${HUB_ROOT:-$HOME/DClaw-Stack}"
SOURCE_REPO="$HUB_ROOT/dclaw-chat"
WORKFLOWS=("claude.yml" "claude-code-review.yml")

DRY_RUN=true
APPS_FILTER=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) DRY_RUN=false; shift ;;
    --apps)  APPS_FILTER="$2"; shift 2 ;;
    -h|--help)
      sed -n '2,/^set -/p' "$0" | sed -E 's/^# ?//;/^set -/d'
      exit 0 ;;
    *) echo "unknown arg: $1" >&2; exit 2 ;;
  esac
done

# Verify source files
for wf in "${WORKFLOWS[@]}"; do
  src="$SOURCE_REPO/.github/workflows/$wf"
  if [ ! -f "$src" ]; then
    echo "FAIL: source workflow not found: $src" >&2
    exit 1
  fi
done

# Discover targets — every dclaw-* repo missing at least one of the workflows
declare -a TARGETS=()
if [ -n "$APPS_FILTER" ]; then
  IFS=',' read -ra REQUESTED <<<"$APPS_FILTER"
  for app in "${REQUESTED[@]}"; do
    [ -d "$HUB_ROOT/$app/.git" ] && TARGETS+=("$app")
  done
else
  for repo_path in "$HUB_ROOT"/dclaw-*/; do
    repo=$(basename "$repo_path")
    [ -d "$HUB_ROOT/$repo/.git" ] || continue
    # Skip if both workflows already present
    missing=false
    for wf in "${WORKFLOWS[@]}"; do
      [ ! -f "$HUB_ROOT/$repo/.github/workflows/$wf" ] && missing=true
    done
    $missing && TARGETS+=("$repo")
  done
fi

MODE="DRY-RUN"; $DRY_RUN || MODE="APPLY + PUSH"
printf '\n══ propagate-claude-workflows · %s ══\n' "$MODE"
printf '   source repo: %s\n' "$SOURCE_REPO"
printf '   targets    : %d repos\n\n' "${#TARGETS[@]}"

ok=0; fail=0; skip=0
for repo in "${TARGETS[@]}"; do
  path="$HUB_ROOT/$repo"

  # Skip if dirty
  if ! /usr/bin/git -C "$path" diff --quiet 2>/dev/null || \
     ! /usr/bin/git -C "$path" diff --cached --quiet 2>/dev/null; then
    printf '  ⚠  %s — uncommitted changes, skipping\n' "$repo"
    ((skip++))
    continue
  fi

  mkdir -p "$path/.github/workflows"
  added=()
  for wf in "${WORKFLOWS[@]}"; do
    dst="$path/.github/workflows/$wf"
    if [ ! -f "$dst" ]; then
      added+=("$wf")
    fi
  done

  if [ "${#added[@]}" -eq 0 ]; then
    printf '  ·  %s — both workflows already present\n' "$repo"
    continue
  fi

  if $DRY_RUN; then
    printf '  →  %s — would add: %s\n' "$repo" "${added[*]}"
    continue
  fi

  for wf in "${added[@]}"; do
    cp "$SOURCE_REPO/.github/workflows/$wf" "$path/.github/workflows/$wf"
    /usr/bin/git -C "$path" add ".github/workflows/$wf"
  done

  msg="chore(ci): add Claude Code Action workflow(s)

Propagated from dclaw-chat (canonical source):
$(printf -- '- %s\n' "${added[@]}")

NOTE: requires CLAUDE_CODE_OAUTH_TOKEN secret to be set on this repo.
Set with: gh secret set CLAUDE_CODE_OAUTH_TOKEN --body \"\$TOKEN\""

  if ! /usr/bin/git -C "$path" commit -m "$msg" >/dev/null 2>&1; then
    printf '  ✗  %s — commit failed\n' "$repo"
    ((fail++))
    continue
  fi

  branch=$(/usr/bin/git -C "$path" rev-parse --abbrev-ref HEAD)
  if ! /usr/bin/git -C "$path" pull --rebase origin "$branch" >/dev/null 2>&1; then
    printf '  ⚠  %s — committed locally; rebase failed\n' "$repo"
    ((fail++))
    continue
  fi
  if /usr/bin/git -C "$path" push origin "$branch" >/dev/null 2>&1; then
    printf '  ✓  %s — added: %s\n' "$repo" "${added[*]}"
    ((ok++))
  else
    printf '  ⚠  %s — push failed\n' "$repo"
    ((fail++))
  fi
done

printf '\n── summary ──\n'
printf '   ok      : %d\n' "$ok"
printf '   skipped : %d\n' "$skip"
printf '   failed  : %d\n' "$fail"
$DRY_RUN && printf '\n(dry-run — re-run with --apply to commit + push)\n'
