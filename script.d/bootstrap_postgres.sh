#!/bin/bash

set -EeufCo pipefail
IFS=$'\t\n'

echo "Creating bootstrapped model"
psql -h "postgres" -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -f db.d/postgres/0.0.1-bootstrap.sql

echo "Applying fingerprint model update"
psql -h "postgres" -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -f db.d/postgres/0.1.0-fingerprint.sql

echo "Done creating model"