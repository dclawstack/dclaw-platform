#!/usr/bin/env bash
# DClaw Swarm Dispatch Runner
# Usage:
#   ./run-swarm.sh              # Print all dispatch commands
#   ./run-swarm.sh chat         # Print command for a specific app
#   ./run-swarm.sh --priority 0 # Print only P0 apps
#   ./run-swarm.sh --batch 5    # Print first 5 apps
#
# Copy-paste the output into Kimi web agent sessions or use with automation.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFEST="$SCRIPT_DIR/apps-manifest.json"

if ! command -v jq &>/dev/null; then
  echo "❌ jq is required. Install it first: brew install jq"
  exit 1
fi

filter_app=""
filter_priority=""
batch_size=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --priority)
      filter_priority="$2"
      shift 2
      ;;
    --batch)
      batch_size="$2"
      shift 2
      ;;
    --help|-h)
      echo "Usage: $0 [app_id] [--priority N] [--batch N]"
      exit 0
      ;;
    *)
      filter_app="$1"
      shift
      ;;
  esac
done

build_command() {
  local app_id="$1"
  local name="$2"
  local repo_url="$3"
  local prompt_path="$4"
  local priority="$5"

  local abs_prompt_path="$SCRIPT_DIR/$prompt_path"

  cat <<EOF

# ============================================================
# Swarm Agent Dispatch: $name ($app_id)
# Priority: P$priority
# Repo: $repo_url
# Prompt: $prompt_path
# ============================================================

## Step 1 — Open Kimi web agent
URL: https://kimi.ai (or your Kimi agent endpoint)

## Step 2 — Paste the system prompt below

SYSTEM PROMPT (copy everything between the dashed lines):

---START SYSTEM PROMPT---
$(cat "$abs_prompt_path")
---END SYSTEM PROMPT---

## Step 3 — Give the agent its mission

Agent, you are building $name v1.0.

1. Clone the scaffold repo:
   git clone $repo_url
   cd dclaw-$app_id

2. Read the existing scaffold (frontend/, backend/, helm/, docs/)

3. Produce a complete v1.0-spec.md following the phases in the system prompt

4. Commit the spec to the repo:
   git checkout -b v1.0-spec
   git add v1.0-spec.md
   git commit -m "docs: add v1.0 specification"
   git push origin v1.0-spec

5. Open a PR on GitHub for review.

## Step 4 — Repeat for next app
Run: $0 $app_id

EOF
}

# Build jq filter
jq_filter=".[]"

if [[ -n "$filter_app" ]]; then
  jq_filter=".[] | select(.app_id == \"$filter_app\")"
fi

if [[ -n "$filter_priority" ]]; then
  jq_filter=".[] | select(.priority == $filter_priority)"
fi

# Read manifest and output commands
apps=$(jq -c "$jq_filter" "$MANIFEST")

count=0
while IFS= read -r app; do
  if [[ -n "$batch_size" && "$count" -ge "$batch_size" ]]; then
    break
  fi

  app_id=$(echo "$app" | jq -r '.app_id')
  name=$(echo "$app" | jq -r '.name')
  repo_url=$(echo "$app" | jq -r '.repo_url')
  prompt_path=$(echo "$app" | jq -r '.prompt_file')
  priority=$(echo "$app" | jq -r '.priority')

  build_command "$app_id" "$name" "$repo_url" "$prompt_path" "$priority"
  ((count++))
done <<< "$apps"

echo ""
echo "# ============================================================"
echo "# Total dispatched: $count app(s)"
echo "# ============================================================"
