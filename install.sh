#!/bin/bash
set -euo pipefail

# DClaw Stack — One-Line Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/dclawstack/dclaw-platform/main/install.sh | bash
#
# This script installs the entire DClaw Platform on a Linux/macOS machine
# with Docker and Docker Compose.

REPO="dclawstack/dclaw-platform"
BRANCH="main"
INSTALL_DIR="${INSTALL_DIR:-$HOME/dclaw-stack}"

color() {
  printf '\033[%sm' "$1"
}
nc() { color '0'; }
green() { color '0;32'; }
yellow() { color '1;33'; }
red() { color '0;31'; }
blue() { color '0;34'; }

echo ""
blue
printf '╔══════════════════════════════════════════════════════════╗\n'
printf '║           DClaw Stack — One-Line Installer               ║\n'
printf '╚══════════════════════════════════════════════════════════╝\n'
nc
echo ""

# ─── Check prerequisites ────────────────────────────────────────────
echo "🔍 Checking prerequisites..."

MISSING=""
if ! command -v docker &>/dev/null; then
  MISSING="$MISSING docker"
fi
if ! docker compose version &>/dev/null && ! docker-compose version &>/dev/null; then
  MISSING="$MISSING docker-compose"
fi
if ! command -v curl &>/dev/null; then
  MISSING="$MISSING curl"
fi

if [ -n "$MISSING" ]; then
  red
  echo "❌ Missing prerequisites:$MISSING"
  echo ""
  echo "Install them first:"
  echo "  Ubuntu/Debian:  sudo apt-get install docker.io docker-compose-plugin curl"
  echo "  Fedora/RHEL:    sudo dnf install docker docker-compose-plugin curl"
  echo "  macOS:          brew install docker docker-compose curl"
  nc
  exit 1
fi

green
echo "✅ Prerequisites OK"
nc

# ─── Create install directory ───────────────────────────────────────
echo ""
echo "📁 Installing to: $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

# ─── Download compose file ──────────────────────────────────────────
echo ""
echo "⬇️  Downloading docker-compose.yml..."
curl -fsSL "https://raw.githubusercontent.com/$REPO/$BRANCH/docker-compose.yml" -o docker-compose.yml

# ─── Create .env ────────────────────────────────────────────────────
if [ ! -f .env ]; then
  echo ""
  echo "🔐 Creating .env file..."
  cat > .env << 'EOF'
# DClaw Stack Environment
# Change this password before running in production!
POSTGRES_PASSWORD=dclaw
EOF
  yellow
  echo "⚠️  Default PostgreSQL password is 'dclaw'"
  echo "   Edit .env to change it before production use."
  nc
fi

# ─── Pull images ────────────────────────────────────────────────────
echo ""
echo "🐳 Pulling Docker images..."
docker compose pull

# ─── Start services ─────────────────────────────────────────────────
echo ""
echo "🚀 Starting DClaw Stack..."
docker compose up -d

# ─── Wait for health ────────────────────────────────────────────────
echo ""
echo "⏳ Waiting for services to start..."
for i in {1..30}; do
  if docker compose ps | grep -q "healthy\|running"; then
    break
  fi
  sleep 1
done

# ─── Summary ────────────────────────────────────────────────────────
echo ""
green
printf '╔══════════════════════════════════════════════════════════╗\n'
printf '║              ✅ DClaw Stack is running!                  ║\n'
printf '╚══════════════════════════════════════════════════════════╝\n'
nc
echo ""
echo "📍 Access your apps:"
echo "   DPanel:       http://localhost:3000"
echo "   DClaw Flow:   http://localhost:3003"
echo "   DClaw RAG:    http://localhost:3009"
echo "   DClaw Med:    http://localhost:3004"
echo "   DClaw Code:   http://localhost:3005"
echo "   DClaw Learn:  http://localhost:3006"
echo "   DClaw SEO:    http://localhost:3007"
echo "   DClaw Create: http://localhost:3008"
echo ""
echo "🛠️  Management commands:"
echo "   cd $INSTALL_DIR"
echo "   docker compose logs -f        # View logs"
echo "   docker compose down            # Stop all"
echo "   docker compose pull            # Update images"
echo "   docker compose up -d           # Start after update"
echo ""
