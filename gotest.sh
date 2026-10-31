#!/usr/bin/env bash
# Runs a go command in golang:1.22 on the compose network so tests reach MySQL at host "mysql".
# Usage: ./gotest.sh go test ./...
cd "$(dirname "$0")"
MSYS_NO_PATHCONV=1 exec docker run --rm --network brightloop-subscription-system_default \
  -e GOTOOLCHAIN=local -e TEST_MYSQL_ADDR=mysql:3306 \
  -v "$(pwd -W 2>/dev/null || pwd)/backend:/src" -v gomod:/go/pkg/mod -v gobuild:/root/.cache/go-build \
  -w /src golang:1.22 "$@"
