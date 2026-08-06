#!/bin/bash
set -e

psql \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  <<-'SQL'
    CREATE ROLE replicator
      WITH REPLICATION
      LOGIN
      PASSWORD 'replica_password';
SQL

# Permit the replica container to make a replication connection.
echo \
  "host replication replicator 0.0.0.0/0 scram-sha-256" \
  >> "$PGDATA/pg_hba.conf"