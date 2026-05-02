#!/usr/bin/env bash
# Sync the production SQLite database.
#
# Usage:
#   ./scripts/prod-db.sh pull [local_path]   — download prod DB locally
#   ./scripts/prod-db.sh push [local_path]   — upload local DB to prod (stops/starts app)
#
# Requires: ssh access to lounge@192.168.1.197

set -euo pipefail

REMOTE="lounge@192.168.1.197"
REMOTE_DB="/home/lounge/eros/data/db/eros.sqlite"
LOCAL_DB="${2:-./eros-prod.sqlite}"

case "${1:-}" in
  pull)
    echo "Checkpointing WAL on remote..."
    ssh -t "$REMOTE" "sudo sqlite3 '$REMOTE_DB' 'PRAGMA wal_checkpoint(TRUNCATE);'"
    echo "Copying database to $LOCAL_DB..."
    scp "$REMOTE:$REMOTE_DB" "$LOCAL_DB"
    echo "Done. Open with: sqlite3 $LOCAL_DB"
    ;;

  push)
    if [[ ! -f "$LOCAL_DB" ]]; then
      echo "Error: $LOCAL_DB not found" >&2
      exit 1
    fi

    echo "WARNING: This will replace the production database."
    echo "Any data written to prod since your last pull will be LOST."
    read -rp "Continue? [y/N] " confirm
    [[ "$confirm" =~ ^[Yy]$ ]] || { echo "Aborted."; exit 1; }

    echo "Checkpointing local WAL..."
    sqlite3 "$LOCAL_DB" 'PRAGMA wal_checkpoint(TRUNCATE);'

    echo "Stopping backend on remote..."
    ssh "$REMOTE" "docker compose -f /home/lounge/eros/docker-compose.yml stop backend"

    echo "Uploading $LOCAL_DB to remote..."
    scp "$LOCAL_DB" "$REMOTE:/tmp/eros-upload.sqlite"

    echo "Replacing remote DB and cleaning WAL/SHM..."
    ssh "$REMOTE" "rm -f '$REMOTE_DB' '${REMOTE_DB}-wal' '${REMOTE_DB}-shm' && mv /tmp/eros-upload.sqlite '$REMOTE_DB'"

    echo "Starting backend on remote..."
    ssh "$REMOTE" "docker compose -f /home/lounge/eros/docker-compose.yml start backend"

    echo "Done. Production database replaced."
    ;;

  *)
    echo "Usage: $0 {pull|push} [local_path]" >&2
    exit 1
    ;;
esac
