#!/bin/bash
source .env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "host=pg-prod-chat port=5433 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable" up -v