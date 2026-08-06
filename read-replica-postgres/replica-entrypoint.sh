#!/bin/bash
set -e

# The script initially runs as root so that it can prepare the mounted volume.
mkdir -p "$PGDATA"
chown -R postgres:postgres "$PGDATA"
chmod 700 "$PGDATA"

if [ ! -s "$PGDATA/PG_VERSION" ]; then
  echo "Waiting for primary..."

  until pg_isready -h primary -p 5432 -U postgres; do
    sleep 1
  done

  echo "Taking base backup from primary..."

  rm -rf "${PGDATA:?}"/*

  export PGPASSWORD="$REPLICATION_PASSWORD"

  gosu postgres pg_basebackup \
    --host=primary \
    --port=5432 \
    --username="$REPLICATION_USER" \
    --pgdata="$PGDATA" \
    --wal-method=stream \
    --write-recovery-conf \
    --progress

  unset PGPASSWORD
fi

chown -R postgres:postgres "$PGDATA"
chmod 700 "$PGDATA"

exec gosu postgres postgres