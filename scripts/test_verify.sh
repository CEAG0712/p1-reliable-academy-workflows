#!/bin/sh
# Isolated command fixtures test orchestration, not application correctness.
set -eu

script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd -P)
fixture=$(mktemp -d "${TMPDIR:-/tmp}/p1-verify.XXXXXX")
trap 'rm -rf "$fixture"' 0
trap 'exit 1' HUP INT TERM
repository="$fixture/project with spaces"
bin="$fixture/bin"
mkdir -p "$repository/scripts" "$repository/nested/deeper" "$bin"
cp "$script_dir/verify.sh" "$repository/scripts/verify.sh"

cat > "$fixture/command" <<'COMMAND'
#!/bin/sh
set -eu
[ "$PWD" = "$EXPECTED_ROOT" ] || exit 98
name=${0##*/}
if [ "$name $*" = 'python3 -m unittest discover -s scripts -p test_*.py -v' ]; then
    [ "${PYTHONDONTWRITEBYTECODE:-}" = 1 ] || exit 99
fi
printf '%s\n' "$name $*" >> "$COMMAND_LOG"
if [ "$name $*" = "$FAIL_COMMAND" ]; then
    exit "$FAIL_STATUS"
fi
COMMAND
chmod +x "$fixture/command"
for tool in python3 go git; do
    cp "$fixture/command" "$bin/$tool"
done
# dirname is the only external utility used by the verifier itself.
ln -s "$(command -v dirname)" "$bin/dirname"

cat > "$fixture/expected" <<'EXPECTED'
python3 scripts/validate_features.py feature_list.json
python3 -m unittest discover -s scripts -p test_*.py -v
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
git diff --check
git diff --cached --check
EXPECTED

export EXPECTED_ROOT COMMAND_LOG FAIL_COMMAND FAIL_STATUS
EXPECTED_ROOT=$(CDPATH= cd "$repository" && pwd -P)
COMMAND_LOG="$fixture/commands.log"
FAIL_COMMAND=
FAIL_STATUS=0
tests=0

fail() {
    printf 'FAIL: %s\n' "$*" >&2
    cat "$fixture/output" >&2
    exit 1
}

invoke() {
    : > "$COMMAND_LOG"
    if (cd "$invocation_dir" && PATH="$bin" /bin/sh "$repository/scripts/verify.sh") \
        > "$fixture/output" 2>&1; then
        actual_status=0
    else
        actual_status=$?
    fi
}

passed() {
    tests=$((tests + 1))
    printf 'ok %s - %s\n' "$tests" "$*"
}

for invocation_dir in "$repository" "$repository/nested/deeper"; do
    invoke
    [ "$actual_status" -eq 0 ] || fail "successful run exited $actual_status"
    cmp -s "$fixture/expected" "$COMMAND_LOG" || fail 'command order or arguments differ'
    grep -Fq 'PASS: all six verification checks succeeded.' "$fixture/output" || fail 'missing final PASS'
    for number in 1 2 3 4 5 6; do
        grep -Fq "=== Check $number:" "$fixture/output" || fail "missing heading $number"
    done
    passed "success, ordered commands, headings and root resolution from $invocation_dir"
done

index=0
while IFS= read -r FAIL_COMMAND; do
    index=$((index + 1))
    FAIL_STATUS=$((40 + index))
    case "$index" in
        1) label='Check 1: Feature tracker validation' ;;
        2) label='Check 2: Python harness tests' ;;
        3) label='Check 3: Go unit tests' ;;
        4) label='Check 4: Go race detector' ;;
        5) label='Check 5: Go static analysis' ;;
        6) label='Check 6: Git whitespace validation (working tree)' ;;
        7) label='Check 6: Git whitespace validation (index)' ;;
    esac
    invoke
    [ "$actual_status" -eq "$FAIL_STATUS" ] || fail "lost exit code $FAIL_STATUS: got $actual_status"
    head -n "$index" "$fixture/expected" > "$fixture/prefix"
    cmp -s "$fixture/prefix" "$COMMAND_LOG" || fail "did not stop at $FAIL_COMMAND"
    grep -Fq "FAIL: $label (exit $FAIL_STATUS)" "$fixture/output" || fail 'incorrect failure diagnostic'
    if grep -Fq 'PASS' "$fixture/output"; then
        fail 'success message after failure'
    fi
    passed "$FAIL_COMMAND: exit $FAIL_STATUS preserved; stopped immediately; no PASS"
done < "$fixture/expected"

FAIL_COMMAND=
for tool in python3 go git; do
    mv "$bin/$tool" "$fixture/disabled-tool"
    invoke
    mv "$fixture/disabled-tool" "$bin/$tool"
    [ "$actual_status" -eq 127 ] || fail "missing $tool did not exit 127"
    [ ! -s "$COMMAND_LOG" ] || fail "checks ran without $tool"
    grep -Fq "FAIL: required tool unavailable: $tool" "$fixture/output" || fail 'missing tool diagnostic'
    if grep -Fq 'PASS' "$fixture/output"; then
        fail 'success message with missing tool'
    fi
    passed "missing $tool: exit 127 before any checks; no PASS"
done

printf 'PASS: %s verifier shell tests.\n' "$tests"
