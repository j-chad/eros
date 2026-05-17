#!/bin/bash
set -euo pipefail

# Deploy Eros.
# Run from the repo root on your local machine.
#
# Usage:
#   ./scripts/deploy.sh          - deploy everything (server + client)
#   ./scripts/deploy.sh server   - deploy backend + admin to NAS
#   ./scripts/deploy.sh client   - deploy client PWA to GitHub Pages
#
# Prerequisites:
#   server: NAS reachable via SSH, Docker Desktop running, ~/eros/ on NAS with .env
#   client: Node installed (see .nvmrc), git push access to origin

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

deploy_server() {
    # Verify NAS is reachable
    step "Checking NAS connectivity"
    if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$NAS_HOST" true; then
        echo "Error: cannot reach $NAS_HOST - check SSH key auth (ssh-copy-id $NAS_HOST)" >&2
        exit 1
    fi

    # Build admin frontend
    step "Building admin frontend"
    (cd admin && npm run build)

    # Build Docker images
    step "Building backend Docker image (linux/amd64)"
    docker build --platform linux/amd64 -f Dockerfile.backend -t eros-backend .

    step "Building Caddy Docker image (linux/amd64)"
    docker build --platform linux/amd64 -f Dockerfile.caddy -t eros-caddy .

    # Transfer images to NAS
    step "Transferring Docker images to NAS"
    docker save eros-backend eros-caddy | ssh "$NAS_HOST" "docker load"

    # Sync admin build
    step "Syncing admin build to NAS"
    ssh "$NAS_HOST" "mkdir -p $NAS_DIR/admin-build"
    rsync -a --delete admin/build/ "$NAS_HOST:$NAS_DIR/admin-build/"

    # Copy compose + config files
    step "Copying compose and Caddy config"
    scp docker-compose.yml Caddyfile "$NAS_HOST:$NAS_DIR/"

    # Run database migrations
    step "Running database migrations"
    ssh "$NAS_HOST" "cd $NAS_DIR && docker compose stop backend"
    ssh "$NAS_HOST" "cd $NAS_DIR && docker compose run --rm backend eros-backend migrate"

    # Restart the stack
    step "Restarting Docker Compose stack"
    ssh "$NAS_HOST" "cd $NAS_DIR && docker compose up -d --force-recreate"

    step "Server deploy complete"
}

deploy_client() {
    # Build client frontend
    step "Building client frontend"
    (cd client && npm run build)

    # Deploy to GitHub Pages
    step "Deploying client to GitHub Pages"
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

    step "Client deploy complete"
}

TARGET="${1:-}"

case "$TARGET" in
    server)
        deploy_server
        ;;
    client)
        deploy_client
        ;;
    all)
        deploy_server
        deploy_client
        ;;
    *)
        echo "Usage: $0 {server|client|all}" >&2
        exit 1
        ;;
esac

step "Deploy complete"
