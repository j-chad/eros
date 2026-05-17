#!/bin/bash
set -euo pipefail

# Deploy Eros to the NAS.
# Run from the repo root on your local machine.
#
# Prerequisites:
#   - NAS reachable via SSH as "$NAS_HOST"
#   - Docker Desktop running locally
#   - Node installed (see .nvmrc)
#   - ~/eros/ directory exists on the NAS with .env file

NAS_HOST="lounge@192.168.1.197"
NAS_DIR="~/eros"
CLIENT_DOMAIN="eros.jacksonc.dev"
GH_REPO=$(git remote get-url origin)

step() { printf "\n\033[1;34m==> %s\033[0m\n" "$1"; }

# Verify we're in the repo root
if [[ ! -f Dockerfile.backend ]]; then
    echo "Error: run this script from the repo root" >&2
    exit 1
fi

# Verify NAS is reachable
step "Checking NAS connectivity"
if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$NAS_HOST" true; then
    echo "Error: cannot reach $NAS_HOST - check SSH key auth (ssh-copy-id $NAS_HOST)" >&2
    exit 1
fi

# 1. Build frontends
step "Building admin frontend"
(cd admin && npm run build)

step "Building client frontend"
(cd client && npm run build)

# 2. Build Docker images
step "Building backend Docker image (linux/amd64)"
docker build --platform linux/amd64 -f Dockerfile.backend -t eros-backend .

step "Building Caddy Docker image (linux/amd64)"
docker build --platform linux/amd64 -f Dockerfile.caddy -t eros-caddy .

# 3. Transfer images to NAS
step "Transferring Docker images to NAS"
docker save eros-backend eros-caddy | ssh "$NAS_HOST" "docker load"

# 4. Sync admin build
step "Syncing admin build to NAS"
ssh "$NAS_HOST" "mkdir -p $NAS_DIR/admin-build"
rsync -a --delete admin/build/ "$NAS_HOST:$NAS_DIR/admin-build/"

# 5. Deploy client build to GH Pages
cp client/build/200.html client/build/404.html
touch client/build/.nojekyll
echo "$CLIENT_DOMAIN" > client/build/CNAME
(
  cd client/build
  git init
  git checkout -b gh-pages
  git add .
  git commit -m "Deploy client"
  git remote add origin "$GH_REPO"
  git push -f origin gh-pages
)

# 5. Copy compose + config files
step "Copying compose and Caddy config"
scp docker-compose.yml Caddyfile "$NAS_HOST:$NAS_DIR/"

# 6. Run database migrations
step "Running database migrations"
ssh "$NAS_HOST" "cd $NAS_DIR && docker compose stop backend"
ssh "$NAS_HOST" "cd $NAS_DIR && docker compose run --rm backend eros-backend migrate"

# 7. Restart the stack
step "Restarting Docker Compose stack"
ssh "$NAS_HOST" "cd $NAS_DIR && docker compose up -d --force-recreate"

step "Deploy complete"
