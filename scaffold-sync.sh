#!/usr/bin/env bash
# scaffold-sync.sh — propagate hub-root scaffold artifacts into each dclaw-* repo.
#
# Hub-master location: ~/DClaw-Stack/  (loose files; control plane)
# Per-repo destination: <hub>/<app>/   (cp + git commit + optional push)
#
# v1 scope: PATCH-*.md files and PATCHES.md index.
# NOT in v1: AGENTS.md or PLAN-vN.md — those are per-app templates with
# placeholders ({APP_NAME}, {BACKEND_PORT}, ...) and need a renderer.
#
# Usage:
#   ./scaffold-sync.sh                       # dry-run: show what would change
#   ./scaffold-sync.sh --apply               # commit locally (no push)
#   ./scaffold-sync.sh --apply --push        # commit + push to GitHub main
#   ./scaffold-sync.sh --apps dclaw-agent,dclaw-chat   # limit target apps
#   ./scaffold-sync.sh --files 'PATCH-2026-05-*'        # limit source files (glob)
#
# Exit codes:
#   0  ok (or dry-run completed)
#   1  any per-repo failure
#   2  bad invocation

set -euo pipefail

HUB_ROOT="${HUB_ROOT:-$HOME/DClaw-Stack}"
DEFAULT_FILES_GLOB="PATCH-*.md PATCHES.md"

DRY_RUN=true
PUSH=false
APPS_FILTER=""
FILES_GLOB="$DEFAULT_FILES_GLOB"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply)  DRY_RUN=false; shift ;;
    --push)   PUSH=true; shift ;;
    --apps)   APPS_FILTER="$2"; shift 2 ;;
    --files)  FILES_GLOB="$2"; shift 2 ;;
    -h|--help)
      sed -n '2,/^set -/p' "$0" | sed -E 's/^# ?//;/^set -/d'
      exit 0 ;;
    *)
      echo "unknown arg: $1" >&2; exit 2 ;;
  esac
done

# ---------- Discover source files at hub root ----------
declare -a FILES=()
for pat in $FILES_GLOB; do
  while IFS= read -r -d '' f; do
    FILES+=("$f")
  done < <(find "$HUB_ROOT" -maxdepth 1 -type f -name "$pat" -print0 2>/dev/null | sort -z)
done

if [ "${#FILES[@]}" -eq 0 ]; then
  echo "no scaffold files match: $FILES_GLOB (looked in $HUB_ROOT)"
  exit 0
fi

# ---------- Discover target repos ----------
declare -a REPOS=()
if [ -n "$APPS_FILTER" ]; then
  IFS=',' read -ra REQUESTED <<<"$APPS_FILTER"
  for app in "${REQUESTED[@]}"; do
    if [ -d "$HUB_ROOT/$app/.git" ]; then
      REPOS+=("$HUB_ROOT/$app")
    else
      echo "warn: $app not a git repo, skipping" >&2
    fi
  done
else
  while IFS= read -r -d '' d; do
    [ -d "$d/.git" ] && REPOS+=("$d")
  done < <(find "$HUB_ROOT" -maxdepth 1 -type d -name 'dclaw-*' -print0 | sort -z)
fi

# ---------- Report banner ----------
MODE="DRY-RUN"; $DRY_RUN || MODE="APPLY"
$PUSH && MODE="$MODE + PUSH"
printf '\n══ scaffold-sync · %s ══\n' "$MODE"
printf '   hub-root : %s\n' "$HUB_ROOT"
printf '   files    : %d\n' "${#FILES[@]}"
for f in "${FILES[@]}"; do printf '              · %s\n' "$(basename "$f")"; done
printf '   repos    : %d\n\n' "${#REPOS[@]}"

# ---------- Per-repo sync ----------
fail_count=0
update_count=0
skip_count=0
nochange_count=0

for repo in "${REPOS[@]}"; do
  app=$(basename "$repo")

  # Don't blast over uncommitted work
  if ! git -C "$repo" diff --quiet 2>/dev/null || \
     ! git -C "$repo" diff --cached --quiet 2>/dev/null; then
    printf '  ⚠  %s — skipped (uncommitted changes)\n' "$app"
    ((skip_count++))
    continue
  fi

  # Diff source vs destination
  declare -a TO_COPY=()
  for src in "${FILES[@]}"; do
    base=$(basename "$src")
    dst="$repo/$base"
    if [ ! -f "$dst" ] || ! cmp -s "$src" "$dst"; then
      TO_COPY+=("$base")
    fi
  done

  if [ "${#TO_COPY[@]}" -eq 0 ]; then
    printf '  ·  %s — up to date\n' "$app"
    ((nochange_count++))
    continue
  fi

  if $DRY_RUN; then
    printf '  →  %s — would sync: %s\n' "$app" "${TO_COPY[*]}"
    ((update_count++))
    continue
  fi

  # Apply
  for base in "${TO_COPY[@]}"; do
    cp "$HUB_ROOT/$base" "$repo/$base"
    git -C "$repo" add -- "$base"
  done

  msg="chore(scaffold): sync ${#TO_COPY[@]} hub-master patch$([ "${#TO_COPY[@]}" -eq 1 ] || echo es)

Synced from ~/DClaw-Stack/ (hub control plane):
$(printf -- '- %s\n' "${TO_COPY[@]}")"

  if ! git -C "$repo" commit -m "$msg" >/dev/null 2>&1; then
    printf '  ✗  %s — commit failed\n' "$app"
    ((fail_count++))
    continue
  fi

  if $PUSH; then
    branch=$(git -C "$repo" rev-parse --abbrev-ref HEAD)
    if git -C "$repo" push origin "$branch" >/dev/null 2>&1; then
      printf '  ✓  %s — synced + pushed: %s\n' "$app" "${TO_COPY[*]}"
    else
      printf '  ⚠  %s — synced locally; push FAILED\n' "$app"
      ((fail_count++))
      continue
    fi
  else
    printf '  ✓  %s — synced (local commit): %s\n' "$app" "${TO_COPY[*]}"
  fi
  ((update_count++))
done

# ---------- Summary ----------
printf '\n── summary ──\n'
printf '   updated   : %d\n' "$update_count"
printf '   unchanged : %d\n' "$nochange_count"
printf '   skipped   : %d\n' "$skip_count"
printf '   failed    : %d\n' "$fail_count"
$DRY_RUN && printf '\n(dry-run — re-run with --apply to commit)\n'

[ "$fail_count" -eq 0 ] || exit 1
