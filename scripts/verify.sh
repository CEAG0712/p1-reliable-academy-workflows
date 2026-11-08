#!/bin/sh
# Read-only engineering verification; stop at the first failed command.
set -eu

for tool in dirname python3 go git; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        printf 'FAIL: required tool unavailable: %s\n' "$tool" >&2
        exit 127
    fi
done

repository_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd -P)
cd "$repository_root"

run_check() {
    label=$1
    shift
    printf '\n=== %s ===\n' "$label"
    if "$@"; then
        :
    else
        status=$?
        printf 'FAIL: %s (exit %s)\n' "$label" "$status" >&2
        exit "$status"
    fi
}

run_check 'Check 1: Feature tracker validation' \
    python3 scripts/validate_features.py feature_list.json

export PYTHONDONTWRITEBYTECODE=1
run_check 'Check 2: Python harness tests' \
    python3 -m unittest discover -s scripts -p 'test_*.py' -v
run_check 'Check 3: Go unit tests' go test -count=1 ./...
run_check 'Check 4: Go race detector' go test -race -count=1 ./...
run_check 'Check 5: Go static analysis' go vet ./...
run_check 'Check 6: Git whitespace validation (working tree)' git diff --check
run_check 'Check 6: Git whitespace validation (index)' git diff --cached --check

printf '\nPASS: all six verification checks succeeded.\n'
