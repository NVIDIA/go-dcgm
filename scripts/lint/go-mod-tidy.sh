#!/usr/bin/env bash

if ! go mod tidy -v "$@"; then
  exit 2
fi

if ! git diff --exit-code -- go.mod go.sum &> /dev/null; then
    echo "go.mod or go.sum differs, please re-add it to your commit"
    exit 3
fi
