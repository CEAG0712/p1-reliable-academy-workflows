P1 Reliable Academy Workflows - Agent Instructions
Purpose
Build a reusable Go workflow execution engine. Student enrollment is the first use case, not a reason to hardcode the engine to enrollment.
Confirm the current state first

- Read PROGRESS.md before starting work. Treat its status
  as a checkpoint that must be reconciled with the
  current source code, Git history, and test results.
- Read README.md, go.mod, workflow.go, workflow_repository.go, and relevant tests before changing code.
- Check the working tree (git status) and the current implementation. README milestone notes may lag local work.
- Ask for clarification or stop with a blocker rather than inventing missing M3 requirements.
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

- feature_list.json: approved M3 task sequence,
  dependencies, acceptance criteria, and evidence.
  Read before selecting implementation work.
- Run `python3 scripts/validate_features.py feature_list.json`
  before task activation and after tracker changes; resolve any
  validation failures before proceeding.
- Work on at most one in_progress task at a time.
- Do not mark a task verified without reviewed
  acceptance evidence and actual test results.
- Do not change task scope or activate a new task
  without authorization.

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
- Standard tests: go test -count=1 ./...
- Race checks: go test -race -count=1 ./...
- Report which commands were actually executed and their outcomes.
- If a required check cannot run, report the blocker; do not describe the task as verified.
- A passing unit test is not proof of cloud integration or end-to-end correctness.
  Handoff
- Report the files changed, behavior implemented, commands run, failures, and next action.
- Do not claim the feature is complete without evidence satisfying its acceptance criteria.

Next harness increment

- scripts/verify.sh is planned.
- Reconcile the tracker, PROGRESS.md, source code,
  and current Git state before implementation.
