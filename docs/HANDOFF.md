# P1 — Durable Session Handoff

## 1. Project and business purpose

Build a reusable Go workflow execution engine. Student enrollment and
course access are the first use case; keep business logic outside the
generic engine. This is an in-memory engineering portfolio prototype.

## 2. Repository identity and checkpoint revision

- Repository: CEAG0712/p1-reliable-academy-workflows; branch: main.
- Module: example.com/p1-reliable-academy-workflows; Go: 1.24.5.
- Previous verified checkpoint and Lesson 10 starting HEAD:
  `b0e4db4ec4d20a7e3d150de0342e9f80821669cb`, the confirmed Lesson 9
  GitHub commit (harness - Unified Verification).
- On 2026-11-08, starting tree was clean and HEAD matched cached
  origin/main. No fetch performed; live local-to-remote synchronization
  is unverified.
- Lesson 10 was committed as `01ef1e37f288459276d3025335b7cdcde669f733`
  (chore: add durable coding-agent handoff), confirmed in local Git history.
- Current observed HEAD and Lesson 11 starting revision (2026-11-08):
  `01ef1e37f288459276d3025335b7cdcde669f733`, on main with a clean starting
  tree, matching expected checkpoint and cached origin/main. No fetch was run.
  Lesson 11 changes are uncommitted: scripts/task_preflight.py,
  scripts/test_task_preflight.py, AGENTS.md, PROGRESS.md, and this document.
  No Lesson 11 commit is asserted; recheck HEAD/tree at the next startup.

## 3. Confirmed implemented milestones

Source/tests inspected at the starting revision support:

- M1: create-only repository storage, independent snapshots, versioning,
  stale-version rejection, and single-winner claims within one repository
  instance and process.
- M2: owner-aware completion/failure, attempts, and transition history.
- Exercise 82: SaveIfVersion accepts completion only; failures must use
  Fail with a nonempty reason. LastError and failure-history reasons persist.
- Enrollment eligibility, tenant/student/course entitlement identity, and
  locally idempotent grants are implemented separately.
- Initialization, tracker validation, Python harness tests, and unified
  verification scripts exist. M3 is not implemented.

Evidence: workflow.go, workflow_repository.go, their tests, enrollment.go,
and enrollment_test.go. Preserve ownership, versions, diagnostics, and history.

## 4. Current tracker state

feature_list.json is unchanged: all ten tasks M3.1 through M3.10 are
pending; active_task_id is null; work_in_progress_limit is 1. Preserve
their existing order, dependencies, scope, and acceptance criteria in
that file and docs/ROADMAP.md section 13. No task was activated or verified.

## 5. Verification evidence and commands

- Historical user-run unit/race/vet evidence applies to baseline 2538574;
  its test date was not recorded. It is not evidence of testing later HEADs.
- Historical Lesson 9 evidence (2026-11-08), recorded in PROGRESS.md,
  began at 0a6fbc7794fee2cd4d2137631e8805796a3a9336 with harness edits:
  verifier exit 1 from denied Go cache access, then permission-enabled
  rerun exit 0; syntax checks and 12 isolated shell tests also passed.
  That work was subsequently recorded in b0e4db4; these original checks
  were not run as a clean checkout of that later commit.
- Newly executed Lesson 10 checks (2026-11-08), starting at b0e4db4:
  `bash init.sh` and
  `python3 scripts/validate_features.py feature_list.json` exited 0.
  `bash scripts/verify.sh` initially exited 1 at Go unit tests due to
  sandbox cache denial, after tracker and 27 Python tests passed.
  The permission-enabled rerun exited 0: tracker, 27 Python tests,
  `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `git diff --check`, and `git diff --cached --check` all passed.
- Lesson 11 checks (2026-11-08), starting at 01ef1e3 with the five
  uncommitted harness/documentation files listed above: initialization and
  tracker validation exited 0; Python unittest discovery with bytecode disabled
  exited 0 (44 tests: 27 existing, 17 preflight). M3.1 preflight exited 0 with
  HUMAN APPROVAL REQUIRED and no activation; M3.2 exited 1 for unverified
  M3.1, the expected successful negative test. Exact commands are in PROGRESS.md.
  Unified verifier initially exited 1 at go vet due to sandbox cache denial,
  after Python, Go unit and race checks passed; permission-enabled rerun
  exited 0, including tracker, Python, unit/race/vet and both whitespace checks.
  Separate Git whitespace checks and status exited 0. Read-only assertions
  confirmed tracker equality to HEAD and starting SHA-256, ten pending tasks,
  null active_task_id, WIP 1, and exactly five authorized changed files.
  No Go, tracker, initialization or existing verification scripts changed.
  No M3 work, commit, push, or deployment was performed.

These local checks establish no cloud integration or production outcome.

## 6. Unresolved blockers and limitations

Live GitHub reconciliation and explicit M3 authorization remain required.
The user reports the Lesson 10 fresh-agent acceptance test succeeded.
This is historical user-reported evidence, not a newly executed test here.
docs/ROADMAP.md contains historical repository/harness snapshots; use its
scope, but reconcile status against current files, Git, and tests.
Process-local locks and memory do not establish distributed exactly-once
behavior, durable recovery, complete orchestration, or production readiness.
No deployment or measured business/performance results are established.
Preflight checks a tracker snapshot only: it does not authenticate approval,
reserve eligibility, or enforce authorization. Later state changes require
another check; an eligible result never grants implementation authority.

## 7. Next bounded engineering task

First review the Lesson 11 changes and reconcile the repository. M3.1 — Workflow
Definition Models is next, but is NOT authorized for implementation yet.
Acceptance: models compile and tests demonstrate two distinct workflow
definitions without business-specific engine fields. Preserve the approved
ten-task sequence in feature_list.json and docs/ROADMAP.md section 13.

## 8. Authorization boundaries

Lesson 11 authorizes only the two preflight Python files and updates to
AGENTS.md, PROGRESS.md, and this handoff. No M3,
tracker, Go source/tests, initialization/verifier changes, or dependencies.
Do not discard local changes, commit, push, deploy, use sudo, access
production, or add/expose credentials. A task listed as next grants no
implementation authorization. Await human review before commit/push.

## 9. Fresh-agent startup procedure

1. Read AGENTS.md, PROGRESS.md, and this document; run `bash init.sh`
   and review warnings. Initialization grants no implementation authority.
2. Read README.md, go.mod, feature_list.json, scripts/verify.sh, and
   relevant roadmap sections. Inspect actual source and tests.
3. Record `git status --short`, `git branch --show-current`,
   `git rev-parse HEAD`, and `git log -5 --oneline`; cross-check the
   previous checkpoint and current tree. Fetch/compare origin/main only
   when remote access is authorized; cached refs are not live evidence.
4. Run `python3 scripts/validate_features.py feature_list.json`; identify
   the authorized bounded task and acceptance criteria before activation.
   Before requesting M3 activation, run
   `python3 scripts/task_preflight.py --task TASK_ID`. Eligibility is not
   authorization: obtain explicit human approval before activation or
   implementation. Success must not automatically trigger activation,
   implementation, commits, or pushes. Preserve WIP = 1. Stop for missing
   requirements or authorization. Preflight is not an enforcement system.
5. Run `bash scripts/verify.sh` before requesting engineering acceptance;
   report actual commands, exits, failures, reruns, and revision/tree tested.

## 10. Session handoff update rules

After an authorized engineering session, update this document when state
or evidence changes. Record previous checkpoint, current observed HEAD,
uncommitted files, behavior, dated checks/results, blockers, and next action.
Keep historical evidence tied to its revision and context. Never mark a
task verified without reviewed acceptance evidence and actual tests.
Never invent a future SHA, claim a clean tree after edits without observing
it, or call pending work committed. A handoff update creates a new working
state before its commit exists; verify its actual revision at next startup.
