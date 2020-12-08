#!/bin/bash

set -EeufCo pipefail
IFS=$'\t\n'

go mod download

go test -v ./... -tags "integration" -coverprofile=coverage.out $(go list ./... | grep -v '/vendor/')

go tool cover -func=coverage.out
