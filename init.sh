
#!/bin/sh

# P1 Reliable Academy Workflows
# Read-only coding-agent initialization.

set -eu

echo "=== P1 Harness Initialization ==="

# 1. Locate the repository.
ROOT=$(git rev-parse --show-toplevel 2>/dev/null) || {
  echo "ERROR: Not inside a Git repository."
  exit 1
}

cd "$ROOT"
echo "Repository: $ROOT"

# 2. Check the Go toolchain.
if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: Go is not installed or is not on PATH."
  exit 1
fi

echo "Go toolchain: $(go version)"

# 3. Check critical project files.
for file in \
  go.mod \
  README.md \
  AGENTS.md \
  PROGRESS.md \
  docs/ROADMAP.md \
  docs/INTERVIEW_DEFENSE.md
do
  if [ ! -f "$file" ]; then
    echo "ERROR: Required file missing: $file"
    exit 1
  fi
done

echo "Required project documents: present"

# 4. Report the Git checkpoint.
echo "Branch: $(git branch --show-current)"
echo "Commit: $(git rev-parse --short HEAD)"

echo "Working tree:"
STATUS=$(git status --short)

if [ -z "$STATUS" ]; then
  echo "Clean"
else
  printf '%s\n' "$STATUS"
  echo "NOTICE: Review existing changes before editing."
fi

# 5. Report local tracking information.
if git rev-parse --verify refs/remotes/origin/main \
  >/dev/null 2>&1; then
  LOCAL=$(git rev-parse HEAD)
  REMOTE=$(git rev-parse refs/remotes/origin/main)

  if [ "$LOCAL" = "$REMOTE" ]; then
    echo "Local HEAD matches cached origin/main."
  else
    echo "NOTICE: Local HEAD differs from cached origin/main."
  fi
else
  echo "NOTICE: No cached origin/main reference found."
fi

echo "Remote state was not fetched or verified."

# 6. Direct the next agent action.
echo ""
echo "Read AGENTS.md and PROGRESS.md."
echo "Consult relevant sections of docs/ROADMAP.md."
echo "Inspect source code and tests before implementation."
echo "Do not assume documentation proves code behavior."
echo ""
echo "Initialization checks completed."
