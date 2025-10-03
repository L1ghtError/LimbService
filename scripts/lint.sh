#!/bin/bash
set -e

if ! command -v golangci-lint &> /dev/null; then
  echo "⚠️ Warning: golangci-lint not found. Skipping lint check."
  exit 0
fi

golangci-lint run