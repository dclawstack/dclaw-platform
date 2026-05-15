#!/usr/bin/env bash
# Regenerate dpanel/public/cluster-status.json from the live dclawstack
# org state. Run this whenever the Cluster Status tab should be refreshed.
#
# Usage:
#   ./scripts/refresh-cluster-status.sh           # regenerate snapshot
#   ./scripts/refresh-cluster-status.sh --redeploy # also rebuild + redeploy dpanel
#
# Deps: gh CLI (authenticated), python3.

set -euo pipefail

PLATFORM_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HUB_ROOT="$(cd "$PLATFORM_DIR/.." && pwd)"
SCRIPT_DIR="$PLATFORM_DIR/scripts"
DPANEL_PUBLIC="$PLATFORM_DIR/dpanel/public"

REDEPLOY=false
[ "${1:-}" = "--redeploy" ] && REDEPLOY=true

echo "[1/3] discovering dclaw-* app repos…"
gh api orgs/dclawstack/repos --paginate -q '.[] | select(.name | startswith("dclaw-")) | .name' 2>/dev/null \
  | sort -u > /tmp/dclaw_app_repos.txt
grep -vE '^(dclaw-platform|dclaw-panel|dclaw-core|dclaw-prd|dclaw-obsidian|dclaw-enterprise|dclaw-rag-scaffold\.md)$' \
  /tmp/dclaw_app_repos.txt > /tmp/dclaw_app_repos_only.txt
echo "    $(wc -l < /tmp/dclaw_app_repos_only.txt | xargs) app repos"

echo "[2/3] probing readiness (this takes ~30s)…"
python3 "$SCRIPT_DIR/_readiness_probe.py" || {
  echo "readiness probe failed" >&2
  exit 1
}

echo "[3/3] writing cluster-status.json into dpanel/public…"
python3 - <<PY
import json, datetime, pathlib
src = json.load(open("/tmp/readiness.json"))
wrapped = {
    "generated_at": datetime.datetime.now(datetime.timezone.utc).isoformat(),
    "repos": src,
}
out = pathlib.Path("$DPANEL_PUBLIC/cluster-status.json")
out.write_text(json.dumps(wrapped, indent=2))
print(f"    wrote {out} ({len(wrapped['repos'])} repos)")
PY

if $REDEPLOY; then
  echo
  echo "[+1] rebuilding + redeploying dpanel (--redeploy)…"
  cd "$PLATFORM_DIR/dpanel"
  docker build -t dpanel-frontend:local . >/dev/null
  docker tag dpanel-frontend:local 192.168.64.2:5000/dpanel-frontend:local
  docker save 192.168.64.2:5000/dpanel-frontend:local | colima ssh -- sudo k3s ctr -n k8s.io images import - >/dev/null
  kubectl rollout restart deploy/dpanel-frontend -n dclaw
  echo "    dpanel-frontend rolling…"
fi

echo
echo "✓ snapshot refreshed. open http://localhost:30010/cluster-status"
