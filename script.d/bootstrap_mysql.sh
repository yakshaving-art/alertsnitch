#!/bin/bash

set -EeufCo pipefail
IFS=$'\t\n'

echo "Creating DB"
mysql --user=root --password="${MYSQL_ROOT_PASSWORD}" --host=mysql -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE};"

echo "Creating bootstrapped model"
mysql --user=root --password="${MYSQL_ROOT_PASSWORD}" --host=mysql "${MYSQL_DATABASE}" < db.d/mysql/0.0.1-bootstrap.sql

echo "Applying fingerprint model update"
mysql --user=root --password="${MYSQL_ROOT_PASSWORD}" --host=mysql "${MYSQL_DATABASE}" < db.d/mysql/0.1.0-fingerprint.sql

echo "Done creating model"