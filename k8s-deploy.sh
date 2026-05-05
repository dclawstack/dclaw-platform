#!/bin/bash
set -euo pipefail

# DClaw Stack — Kubernetes Deployment Script
# Deploys all DClaw app Helm charts to a K8s cluster.
#
# Prerequisites:
#   - kubectl configured and pointing to your cluster
#   - helm installed
#   - Docker images pushed to GHCR (or another registry)
#
# Usage:
#   ./k8s-deploy.sh [namespace] [image_tag]
#   ./k8s-deploy.sh dclaw latest

NAMESPACE="${1:-dclaw}"
TAG="${2:-latest}"
REGISTRY="${REGISTRY:-ghcr.io/dclawstack}"

APPS=(
  "flow:8088"
  "rag:8090"
  "agent:8091"
  "med:8092"
  "code:8094"
  "learn:8093"
  "seo:8095"
)

color() { printf '\033[%sm' "$1"; }
nc() { color '0'; }
green() { color '0;32'; }
yellow() { color '1;33'; }
blue() { color '0;34'; }

echo ""
blue
printf '╔══════════════════════════════════════════════════════════╗\n'
printf '║         DClaw Stack — K8s Deployment                     ║\n'
printf '╚══════════════════════════════════════════════════════════╝\n'
nc
echo ""
echo "Namespace:  $NAMESPACE"
echo "Image tag:  $TAG"
echo "Registry:   $REGISTRY"
echo ""

# ─── Verify cluster access ──────────────────────────────────────────
if ! kubectl cluster-info &>/dev/null; then
  red
  echo "❌ Cannot connect to Kubernetes cluster."
  echo "   Run: kubectl config current-context"
  nc
  exit 1
fi

# ─── Create namespace ───────────────────────────────────────────────
echo "🔧 Creating namespace: $NAMESPACE"
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# ─── Deploy each app ────────────────────────────────────────────────
for spec in "${APPS[@]}"; do
  app_id=$(echo "$spec" | cut -d: -f1)
  port=$(echo "$spec" | cut -d: -f2)
  chart_dir="helm/dclaw-$app_id"

  echo ""
  echo "📦 Deploying dclaw-$app_id (port $port)..."

  if [ ! -d "$chart_dir" ]; then
    yellow
    echo "⚠️  Helm chart not found: $chart_dir — skipping"
    nc
    continue
  fi

  helm upgrade --install "dclaw-$app_id" "$chart_dir" \
    --namespace "$NAMESPACE" \
    --set image.backend.tag="$TAG" \
    --set image.frontend.tag="$TAG" \
    --set image.backend.repository="$REGISTRY/dclaw-$app_id-backend" \
    --set image.frontend.repository="$REGISTRY/dclaw-$app_id-frontend" \
    --wait --timeout 5m \
    2>&1 | tail -5

green
echo "✅ dclaw-$app_id deployed"
nc
done

# ─── Deploy DPanel ──────────────────────────────────────────────────
echo ""
echo "📦 Deploying DPanel..."
if [ -d "helm/dpanel" ]; then
  helm upgrade --install dpanel helm/dpanel \
    --namespace "$NAMESPACE" \
    --set image.tag="$TAG" \
    --wait --timeout 5m \
    2>&1 | tail -5
  green
  echo "✅ DPanel deployed"
  nc
else
  yellow
  echo "⚠️  DPanel helm chart not found — skipping"
  nc
fi

# ─── Summary ────────────────────────────────────────────────────────
echo ""
green
printf '╔══════════════════════════════════════════════════════════╗\n'
printf '║              ✅ Deployment Complete!                     ║\n'
printf '╚══════════════════════════════════════════════════════════╝\n'
nc
echo ""
echo "Get status:   kubectl get pods -n $NAMESPACE"
echo "Get ingress:  kubectl get ingress -n $NAMESPACE"
echo "View logs:    kubectl logs -n $NAMESPACE -l app=dclaw-flow-backend"
echo ""
