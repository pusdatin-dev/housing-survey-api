#!/bin/bash

# Load env file secara manual agar cron bisa pakai
set -a
source /usr/local/bin/.env
set +a

echo "[DEBUG] ENV: DB_NAME=$DB_NAME, DB_HOST=$DB_HOST, DB_USER=$DB_USER"

BACKUP_DIR="/var/backup"
mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.sql"

# Langsung pakai PGPASSWORD di sini
# ALLOWED_TABLES=$(cat /usr/local/bin/allowed_tables.txt | xargs -I {} echo -t {})
EXCLUDES=$(cat /usr/local/bin/forbidden_tables.txt | xargs -n1 -I {} echo --exclude-table={})


echo "[INFO] Running pg_dump with excludes:"
echo "$EXCLUDES"

PGPASSWORD="$DB_PASSWORD" pg_dump \
  --data-only \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d "$DB_NAME" \
  $EXCLUDES \
  > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
  echo "[SUCCESS] Backup successful: $BACKUP_FILE"
else
  echo "[ERROR] Backup failed!"
fi
