P1 Reliable Academy Workflows - Agent Instructions
Purpose
Build a reusable Go workflow execution engine. Student enrollment is the first use case, not a reason to hardcode the engine to enrollment.
Confirm the current state first

- Read PROGRESS.md before starting work. Treat its status
  as a checkpoint that must be reconciled with the
  current source code, Git history, and test results.
- Read docs/HANDOFF.md during session initialization. Cross-check it
  against PROGRESS.md, feature_list.json, Git history, and actual
  source/tests; do not treat a handoff as proof of current state.
- Read README.md, go.mod, workflow.go, workflow_repository.go, and relevant tests before changing code.
- Check the working tree (git status) and the current implementation. README milestone notes may lag local work.
- Ask for clarification or stop with a blocker rather than inventing missing M3 or M4 requirements.
- At session start, run `bash init.sh` when terminal
  access is available.
- Review initialization warnings before editing.
- Initialization does not authorize source changes,
  Git pushes, deployment, or external service access.

Information map - Load context when relevant

- README.md: project purpose, implemented stack,
  and documented commands.
- PROGRESS.md: current milestone, blockers,
  and next actions. Read at session start.
- docs/HANDOFF.md: durable session checkpoint, evidence, boundaries,
  and fresh-agent startup procedure. Read at session start.
- workflow.go and workflow_repository.go:
  read when changing workflow state, ownership,
  version checks, or execution persistence.
- workflow_test.go and workflow_repository_test.go:
  read when changing or verifying workflow behavior.
- enrollment.go and enrollment_test.go:
  read when changing enrollment business rules.
- Do not load unrelated future infrastructure
  documentation unless the task requires it.
- docs/ROADMAP.md: planned milestones, architecture
  boundaries, requirements, and acceptance criteria.
  Consult for scope; verify implementation in code.
- docs/INTERVIEW_DEFENSE.md: project-defense
  procedure for Backend, AI, and FDE interviews.
  Use CHEAT MODE for learning during implementation;
  conduct NO-CHEAT MODE when requested.
- Do not treat interview answers or planned designs
  as proof of implemented engineering behavior.

- feature_list.json: approved unified M3/M4 task sequence,
  dependencies, acceptance criteria, and evidence.
  Read before selecting implementation work.
- Run `python3 scripts/validate_features.py feature_list.json`
  before task activation and after tracker changes; resolve any
  validation failures before proceeding.
- Work on at most one in_progress task at a time across M3 and M4;
  both stages share one active_task_id and WIP limit 1.
- Do not mark a task verified without reviewed
  acceptance evidence and actual test results.
- Do not change task scope or activate a new task
  without authorization.
- Never begin implementation solely because a task is described as next.
- Before requesting activation of an M3 or M4 task, run
  `python3 scripts/task_preflight.py --task TASK_ID` and resolve diagnostics.
  Preflight checks structural eligibility only; eligibility is not authorization.
  Human approval is still required. Preflight success alone must never trigger
  automatic activation, implementation, commits, or pushes. It does not enforce
  human authorization or prove approval from chat, handoff text, environment
  variables, or CLI flags. Preserve WIP = 1 and all verification requirements.

Existing correctness constraints

- Preserve single-winner workflow claims and owner-aware transitions.
- Preserve version checks; never silently overwrite stale workflow state.
- SaveIfVersion currently accepts completion only; failures must use Fail with a nonempty reason.
- Preserve failure diagnostics and the stored transition history.
- Do not silently expand scope or refactor unrelated functionality.
  Execution boundaries
- Work on one bounded feature at a time and identify its acceptance criteria before editing.
- Do not access production resources, add cloud credentials, deploy, or push commits without explicit approval.
- Do not expose credentials or secrets in output or commits.
- Keep the implementation aligned with the code actually present, not only with planned stack items.
  Verification
- Run `bash scripts/verify.sh` before requesting acceptance of
  engineering changes. It checks required tools and runs tracker
  validation, Python harness tests, Go unit tests, race detection,
  static analysis, and working-tree/index whitespace checks in order.
- When changing the verifier, also run `bash -n scripts/verify.sh`,
  `bash -n scripts/test_verify.sh`, and `bash scripts/test_verify.sh`.
- Preserve individual commands for targeted debugging:
- Tracker: python3 scripts/validate_features.py feature_list.json
- Python harness: PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p "test_*.py" -v
- Standard tests: go test -count=1 ./...
- Race checks: go test -race -count=1 ./...
- Static analysis: go vet ./...
- Whitespace: git diff --check and git diff --cached --check
- Report which commands were actually executed and their outcomes.
- If a required check cannot run, report the blocker; do not describe the task as verified.
- A passing unit test is not proof of cloud integration or end-to-end correctness.
  Handoff
- Report the files changed, behavior implemented, commands run, failures, and next action.
- Do not claim the feature is complete without evidence satisfying its acceptance criteria.
- After an authorized engineering session, update docs/HANDOFF.md when
  relevant state or evidence changes. Record the previous checkpoint
  revision separately from the current observed HEAD and working-tree
  changes. Never invent a future commit hash or imply uncommitted work
  is committed; recheck the revision at the next startup.

Harness checkpoint

- scripts/verify.sh and scripts/test_verify.sh are implemented.
- Reconcile the tracker, PROGRESS.md, source code,
  and current Git state before selecting further implementation.
