# P1 — Durable Session Handoff

## 1. Project and business purpose

Build a reusable Go workflow execution engine. Student enrollment and
course access are the first use case; keep business logic outside the
generic engine. This is an in-memory engineering portfolio prototype.

## 2. Repository identity and checkpoint revision

Current checkpoint: section 13 (M3.8). Earlier revisions/status below are historical.

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
- Lesson 11 started at 01ef1e37f288459276d3025335b7cdcde669f733; its five
  harness/documentation files were subsequently committed in a5b345d.
- Earlier M3.1 documentation checkpoint recorded a5b345d99075debfe652d0182f0374466b7ff8bc
  plus M3.1 working changes. Those files were subsequently committed as
  93c9861b93091f270d4d9a78a0edc8e19376fd65 (m3.1 implemented).
- M3.2 implementation commit: 4fb635168f61125842be7cf88241c25738c6f265.
- Previous documentation checkpoint revision: 74ab066fc22554384a9337588513f1298a7fe0e9,
  which reconciled M3.2 at 4fb6351. Earlier text describing those documentation
  edits as uncommitted is historical; they are committed at 74ab066.
- M3.3 implementation/acceptance used 74ab066 plus then-uncommitted M3.3 work.
- Historical observed HEAD (2026-11-09): 1a5ad1413c80b7024cb82291bdb7843ff23e857a,
  on main, the confirmed local M3.3 implementation commit. Startup tree was clean.
  Commit inspection confirmed exactly feature_list.json, workflow_trigger.go and
  workflow_trigger_test.go, matching the reported revision and file set.
- That M3.3 reconciliation was committed as d0784d8; the subsequent defense
  protocol upgrade is committed as 701c3e748c4f61093175b171d5038199f4326921.
- At that historical checkpoint, GitHub publication of 1a5ad14 was unconfirmed.
  Initialization warned HEAD differed from cached origin/main; no fetch or live
  remote check was performed.
  Cached refs do not prove publication or absence of a push.

- Historical observed HEAD (2026-11-10): 1e34b5d0108949be651e9e3b96ac4cef4d3c5100
  on main, the expected M3.4 implementation commit; startup tree was clean.
  It changes exactly feature_list.json, workflow_handlers.go and workflow_handlers_test.go.
- At that checkpoint, the previous documentation revision was d0784d8.
  Those documentation edits were subsequently committed as dd35ce0.
- At that historical checkpoint, cached origin/main was 701c3e748c4f61093175b171d5038199f4326921.
  No live remote check occurred then; publication was unverified at that time.

## 3. Historical confirmed milestones at the M3.5 checkpoint

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
  verification scripts exist. M3.1 models, M3.2 validation and M3.3 event eligibility
  matching are verified and committed locally. M3.4 named handler registration is
  verified. M3.5 single-condition branch selection is implemented and committed
  locally; final tracker acceptance remains pending.
  The runner is not implemented.

Evidence: workflow.go, workflow_repository.go, their tests, enrollment.go,
and enrollment_test.go. Preserve ownership, versions, diagnostics, and history.

## 4. Historical M3.5 tracker state (superseded by section 12)

M3.1-M3.4 are verified. M3.4 acceptance is committed as
88fe6ed924efd351cb2a10d5cd8f11f326b81006, superseding the historical pending
statements below. M3.5 is in_progress pending final tracker acceptance;
active_task_id is M3.5, WIP is 1, and M3.6-M3.10 remain pending.
The tracker is unchanged here; M3.5 evidence fields remain empty.
Historical human reviews accepted M3.2 and M3.3 at 10/10 against their bounded
contracts. M3.3 acceptance establishes eligibility, not execution.

Historical observed HEAD (2026-11-11): e0e73da376d55cd22c86fdda8d6011af022a9b0a
on main, matching the expected implementation commit; startup tree was clean.
Commit inspection confirms exactly feature_list.json, workflow_branch.go and
workflow_branch_test.go. Previous documentation checkpoint is
dd35ce009d4bfc36a183671c6fbef7f48e4ed3a8, separate from current HEAD.
Only PROGRESS.md and docs/HANDOFF.md are now edited and uncommitted.
Initialization exited 0, warning HEAD differs from cached origin/main.
No live remote check occurred; M3.5 GitHub publication is unverified.

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

- M3.1 (2026-11-08), HEAD a5b345d plus uncommitted changes: initialization,
  tracker validation and M3.1 preflight exited 0 before activation. Human approval
  came from the explicit user message, not preflight success. Focused command
  `go test -count=1 -run '^TestWorkflowDefinition' ./...` exited 1 (undefined model
  types), then exited 0 after adding the models. Initial `bash scripts/verify.sh`
  exited 1 with four validator fixture failures. After fixture correction it exited
  1 at Go unit tests from sandbox Go-cache denial; permission-enabled rerun exited 0:
  tracker, all 44 Python tests, Go unit/race tests, go vet and both whitespace checks.
  Commands and acceptance review are recorded in PROGRESS.md and feature_list.json.

Historical M3.1 final verification after tracker/checkpoint updates: sandbox run exited 1 at
go vet after tracker, all 44 Python tests, Go unit and race tests passed.
Permission-enabled final rerun of `bash scripts/verify.sh` exited 0 with
all six checks PASS. Final scope assertions confirmed M3.2-M3.10 unchanged,
WIP 1, active_task_id null and existing lifecycle/business Go files unchanged.
Only this evidence text was added afterward; tracker validation and both
Git whitespace checks were rerun successfully.

M3.2 implementation and acceptance evidence (2026-11-08):

- Human engineering review accepted M3.2 at 10/10, with no blocking functional
  defect, against both approved criteria:
  1. Valid definitions are accepted.
  2. Malformed definitions produce deterministic validation errors.
- Historical implementation evidence was recorded at
  93c9861b93091f270d4d9a78a0edc8e19376fd65 plus then-uncommitted M3.2 changes.
  The subsequent implementation commit is now confirmed locally as
  4fb635168f61125842be7cf88241c25738c6f265. No future revision is inferred.
- RED: `go test -count=1 -run '^TestWorkflowDefinitionValidate' ./...` exited 1:
  WorkflowDefinition had no Validate method.
- GREEN: that focused command and its `-v` rerun exited 0. Four test functions
  cover eight valid graphs, 31 malformed cases, five cycle cases and 200
  deterministic/read-only iterations. The unchanged-data assertion covers the
  Steps map in broken-reference and cycle scenarios, not every definition field.
- Implementation-session `bash scripts/verify.sh` initially exited 1 at race
  detection from sandbox Go-cache denial, after tracker validation, 44 Python
  tests and Go unit tests passed. Permission-enabled rerun exited 0, including
  race detection, go vet and both Git whitespace checks.
- Acceptance-finalization tracker validation exited 0. Its unified verifier
  initially exited 1 at go vet from sandbox Go-cache denial after tracker,
  44 Python tests and Go unit/race tests passed. Permission-enabled final rerun
  exited 0 with all six checks PASS. The evidence-only tracker append was
  followed by successful tracker validation and both Git whitespace checks.
- These are preserved pre-commit results, not newly executed checks of this
  documentation session. Fresh post-commit results are recorded separately.

M3.3 implementation and acceptance evidence (2026-11-08/09):

- Human engineering acceptance on 2026-11-09: **10/10** against both approved
  bounded acceptance criteria. CHEAT MODE was chosen during implementation;
  no NO-CHEAT assessment is required for this checkpoint.
- Actual local implementation commit:
  `1a5ad1413c80b7024cb82291bdb7843ff23e857a`. Commit inspection confirms exactly
  feature_list.json, workflow_trigger.go and workflow_trigger_test.go.
- Acceptance criterion 1: "Matching events activate eligible workflows."
  `WorkflowDefinition.MatchesEvent(Event) (bool, error)` returns true/nil for a
  valid definition and an exact event-type match. TestWorkflowDefinitionMatchesEventTypes
  demonstrates document archive and sensor alert workflows using the same API.
  Activation means observable eligibility for a future runner, not actual execution.
- Acceptance criterion 2: "Nonmatching events cause no execution side effects."
  Valid mismatches return false/nil. The API has no repository, handlers, callbacks
  or execution-state access; source inspection establishes this boundary.
  TestWorkflowDefinitionMatchesEventIsDeterministicAndReadOnly checks the entire
  definition against an independent copy, and the event, over 600 evaluations.
- Blank incoming event types return an event error first; invalid definitions
  return existing Validate errors even on a mismatch. Tests cover missing IDs,
  blank triggers and broken graphs. Nonblank names are opaque; comparison is
  exact and case-sensitive without trimming or case folding.
- Historical RED, at 74ab066 plus then-uncommitted tests:
  `go test -count=1 -run '^TestWorkflowDefinitionMatchesEvent' ./...` exited 1:
  MatchesEvent was undefined and Event was undefined. No implementation existed yet.
- Historical GREEN: `go test -count=1 -run '^TestWorkflowDefinitionMatchesEvent' -v ./...`
  exited 0: four test functions, three invalid-definition subtests, and 600 repeated
  deterministic/read-only evaluations passed after the minimal implementation.
- Historical implementation verification: baseline verifier initially exited 1 at
  unit tests from sandbox Go-cache denial; permission-enabled rerun passed all six
  checks. Post-change verifier exited 1 at go vet from sandbox cache denial after
  tracker, 44 Python tests and Go unit/race tests passed; permission-enabled rerun
  exited 0 with all six checks PASS.
- Historical acceptance-finalization verification (2026-11-09), at 74ab066 plus
  the three then-uncommitted M3.3 files: tracker validation exited 0. Initial
  `bash scripts/verify.sh` exited 1 at Go unit tests because sandbox cache trimming
  (trim.txt) was denied after tests reported ok; tracker and 44 Python tests passed,
  later checks were not reached. Permission-enabled rerun exited 0: tracker,
  44 Python tests, Go unit tests, race detection, go vet and both Git whitespace
  checks PASS. Separate whitespace/status checks passed; both then-untracked Go
  files were fully inspected, passed explicit no-index whitespace checks and
  required no gofmt changes. Final evidence-only updates passed tracker/whitespace
  and scope checks. These are preserved pre-commit results, not fresh checks here.
- The tracker's historical revision `74ab066fc22554384a9337588513f1298a7fe0e9`
  plus uncommitted M3.3 work now maps to implementation commit 1a5ad14 above.
  Historical "pending review" and "no implementation commit" statements are
  superseded by human acceptance and the observed local commit, respectively.
  The tracker is preserved unchanged by this documentation session.

M3.4 historical implementation evidence (2026-11-10; 701c3e7 plus then-uncommitted work):

- RED: go test -count=1 -run '^TestHandlerRegistry' ./... exited 1 with undefined
  NewHandlerRegistry, ActionHandler and ConditionHandler.
- GREEN: go test -count=1 -run '^TestHandlerRegistry' -v ./... exited 0;
  all six registry test functions passed after implementation and gofmt.
- Baseline bash scripts/verify.sh passed. Post-change sandbox run exited 1 at
  Go unit-test setup from denied build-cache access after tracker and 44 Python
  tests passed. Permission-enabled rerun exited 0: tracker, 44 Python tests,
  Go unit/race tests, go vet and both Git whitespace checks passed. Earlier
  M1/M2, Exercise 82 and M3.1-M3.3 regressions passed. These are preserved
  session-report/tracker results, not freshly executed RED/GREEN checks here.

These local checks establish no cloud integration or production outcome.

## 6. Historical M3.5 blockers and limitations (current limits in section 12)

M3.1-M3.4 are verified; M3.5 implementation is committed with final acceptance pending.
No blocking source mismatch was found. M3.5 live publication is unverified here.
Registry registration must not overlap registry operations across goroutines;
callback safety belongs to callers. Populated registries must not be copied.
Predicate errors must be checked before using results. No panic recovery, retries,
cancellation, idempotency, runner or lifecycle integration exists
in this registry. M3.5 adds single-condition branch selection only.
No cloud deployment or production integration is established.
Historical roadmap and README next-task descriptions remain outside this reconciliation.
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

## 7. Historical M3.5 result and next bounded task

M3.1 generic definition, trigger and step models describe enrollment branching and
customer follow-up. M3.2 adds WorkflowDefinition.Validate() error: required fields,
supported step types, handler names, successor references, compatible fields,
positive WAIT durations and cycles (including disconnected components). Sorted
step IDs and fixed traversal order produce deterministic errors. Empty successors
remain path endings. M3.3 adds Event and MatchesEvent for generic, read-only
eligibility decisions. M3.4 adds ActionHandler, ConditionHandler and HandlerRegistry,
with a constructor and usable zero value. RegisterAction/RegisterCondition reject
blank names, nil callbacks and duplicates without replacing originals. InvokeAction
and EvaluateCondition dispatch exact, case-sensitive names through separate maps;
callback outputs/errors pass through unchanged. No business integration is imported.

Approved criterion 1: "Multiple implementations can be registered and invoked by name."
TestHandlerRegistryNamedInvocation checks two actions, two predicates and the exact
call sequence. Approved criterion 2: "Unknown handlers are rejected."
TestHandlerRegistryRejectsUnknownNames checks missing/blank/mismatched names and
wrong namespaces without callback side effects. Additional tests cover invalid
registration, duplicate preservation, error identity, zero-value use and isolation.
The committed files match the implementation shown in the previous session.

M3.5 adds HandlerRegistry.SelectConditionSuccessor(StepDefinition) (string, error).
TRUE selects OnTrue; FALSE selects OnFalse. Predicate errors return the original
error and no successor. Unknown handler errors propagate; wrong step types are
rejected before callback invocation. Empty successors with nil error terminate paths.
Approved criteria: "Both branches are tested" and "Condition errors are handled
and tested." TestHandlerRegistrySelectConditionSuccessor covers both branches,
both Boolean/error combinations, four terminal cases and one evaluation per attempt.
TestHandlerRegistrySelectConditionSuccessorUnknown covers unknown names and zero
value; TestHandlerRegistrySelectConditionSuccessorRejectsWrongType covers invalid types.
Selection preserves configuration and registrations; callbacks own their side effects.
Callers must validate definitions/references and check errors before successors.
Selection does not advance execution or call repository Complete/Fail. Nil registry
pointers are unsupported; no deduplication, panic recovery or durable guarantee exists.

Historical M3.5 RED/GREEN and verification evidence was recovered from local Codex
session 01a1cba2-65a3-7882-a7fd-6e562fa76160 (2026-11-11), at 88fe6ed plus
then-uncommitted M3.5 changes; full details are in PROGRESS.md's M3.5 checkpoint.
`go test -count=1 -run '^TestHandlerRegistrySelectConditionSuccessor' ./...`
exited 1 with the method undefined at test lines 41, 73, 78 and 93. Its `-v`
GREEN rerun exited 0: three test functions and eight subtests passed.
Baseline verifier passed. Post-change sandbox verifier exited 1 at race-test setup
from Go-cache denial after tracker, 44 Python tests and Go unit tests passed;
permission-enabled rerun exited 0 with all six checks PASS. Historical remote
lookup failed with sandbox DNS exit 128, then permission-enabled lookup confirmed
88fe6ed; it proves no publication of the later M3.5 commit. These are preserved
implementation-session observations, separate from fresh checks below.

Next: human engineering review of these two uncommitted documentation files.
M3.5 remains in_progress pending final tracker acceptance. M3.6 — Sequential Action
Execution is next planned, NOT authorized. Runner execution, wait/resume and M2
lifecycle integration remain future work requiring separate authorization.

## 8. Authorization boundaries and reconciliation checks

Historical M3.5 session authorization (superseded for this session by section 12):
M3.5 documentation reconciliation only: PROGRESS.md and
docs/HANDOFF.md. No application code, tests, feature_list.json or roadmap edits.
Do not commit, push, deploy, start M3.6 or mark M3.5 verified.
Implementation completion does not imply acceptance or authorize a later task.

Historical startup (2026-11-09): git rev-parse HEAD, git branch --show-current,
git status --short (empty), git show --stat --oneline HEAD and
git diff-tree --no-commit-id --name-only -r HEAD all exited 0. HEAD matched
1a5ad1413c80b7024cb82291bdb7843ff23e857a on main with exactly the three approved
M3.3 files. bash init.sh exited 0, warning HEAD differs from cached origin/main;
remote state was not verified. Startup tracker validation exited 0.

The following M3.2 documentation-session evidence is historical:

Startup commands exited 0: git rev-parse HEAD, git status --short (empty),
git show --stat --oneline HEAD and git diff-tree --no-commit-id --name-only -r HEAD.
HEAD matched the reported 4fb635168f61125842be7cf88241c25738c6f265 and only the four
approved implementation files were present in that commit. bash init.sh exited 0
and warned that cached origin/main differs and remote state was not verified.
Fresh documentation-session verification results are recorded below separately
from preserved pre-commit implementation and acceptance results.

Fresh documentation-session verification (2026-11-08), at HEAD
4fb635168f61125842be7cf88241c25738c6f265 plus these two documentation edits:

- `bash scripts/verify.sh` initially exited 1 at go vet because sandbox access to
  a Go build-cache file was denied. Tracker validation, all 44 Python tests,
  `go test -count=1 ./...` and `go test -race -count=1 ./...` had passed.
- Permission-enabled rerun of `bash scripts/verify.sh` exited 0: tracker,
  44 Python tests, Go unit/race tests, go vet and both working-tree/index
  whitespace checks passed. Existing M1/M2, Exercise 82, M3.1 and M3.2 tests pass.
- Separate `git diff --check`, `git status --short` and
  `git diff -- PROGRESS.md docs/HANDOFF.md` exited 0. Status listed only the two
  modified documentation files, with no untracked files. The diff was reviewed.
- Only this verification record was added after the successful unified check;
  final whitespace and scope checks cover that evidence-only addition.
- No application code, tests, feature_list.json or roadmap files changed; no
  commit, push, deployment or M3.3 work was performed. These documentation
  changes remain uncommitted and await human review.

Fresh M3.3 post-commit documentation-session verification (2026-11-09), at
HEAD 1a5ad1413c80b7024cb82291bdb7843ff23e857a plus these two documentation edits:

- Initial `bash scripts/verify.sh` exited 1 at go vet because sandbox access to
  a Go build-cache file was denied. Tracker validation, all 44 Python tests,
  `go test -count=1 ./...` and `go test -race -count=1 ./...` passed first.
  The verifier did not reach whitespace checks or print a final PASS.
- Permission-enabled rerun of `bash scripts/verify.sh` exited 0: tracker
  validation, 44 Python tests, Go unit tests, race detection, go vet and both
  working-tree/index whitespace checks passed. M1, M2, Exercise 82, M3.1,
  M3.2 and M3.3 regression tests remain passing. These are fresh local results;
  they establish no cloud integration or production guarantees.
- Separate `git diff --check`, `git status --short` and
  `git diff -- PROGRESS.md docs/HANDOFF.md` exited 0; the diff was reviewed.
  Only PROGRESS.md and docs/HANDOFF.md were modified, with no untracked files.
- Only this verification record and documentation wording clarifications were
  added after the successful unified check. Final whitespace and scope checks
  cover these documentation-only additions. feature_list.json, Go source/tests,
  roadmap and scripts remain unchanged. No commit, push, deployment or M3.4 work
  occurred. These two documentation edits await human review; HEAD is unchanged.

Historical documentation-session startup (2026-11-10): branch main, clean tree and
HEAD 1e34b5d0108949be651e9e3b96ac4cef4d3c5100 confirmed with Git. Commit inspection
found only the three M3.4 files. bash init.sh exited 0 and warned of cached-ref
mismatch; no live remote verification was performed. Tracker validation exited 0.
Previous checkpoint d0784d8 and defense upgrade 701c3e7 are committed, superseding
historical uncommitted-documentation statements. Only the two documentation files
are being edited now; no new documentation commit exists.

Fresh documentation-session verification (2026-11-10), at HEAD
1e34b5d0108949be651e9e3b96ac4cef4d3c5100 plus these two documentation edits:

- python3 scripts/validate_features.py feature_list.json exited 0.
- Initial bash scripts/verify.sh exited 1 at go vet because sandbox access to
  a Go build-cache file was denied. Tracker, all 44 Python tests, Go unit tests
  and race detection passed first; whitespace checks were not reached.
- Permission-enabled rerun of bash scripts/verify.sh exited 0: tracker,
  44 Python tests, go test -count=1 ./..., go test -race -count=1 ./...,
  go vet ./... and working-tree/index whitespace checks passed. Existing
  M1/M2, Exercise 82 and M3.1-M3.4 tests pass locally; this proves no cloud outcome.
- Only this verification record and documentation wording were added afterward.
  Final git diff --check, git diff -- PROGRESS.md docs/HANDOFF.md and git status
  --short passed; the diff was reviewed. Scope checks confirmed only the two
  authorized documentation files changed, with no untracked or staged files;
  feature_list.json and Go source/tests remain byte-for-byte equal to HEAD.
- Ending HEAD remains 1e34b5d0108949be651e9e3b96ac4cef4d3c5100. These documentation
  edits are uncommitted and await human review. No commit, push, deployment,
  live remote verification or M3.5 work occurred.

Fresh M3.5 documentation-session verification (2026-11-11), at HEAD
e0e73da376d55cd22c86fdda8d6011af022a9b0a plus these two documentation edits:

- `bash init.sh` and `python3 scripts/validate_features.py feature_list.json`
  exited 0. Branch/HEAD/status/history and commit file-set checks passed.
- Initial `bash scripts/verify.sh` exited 0 with all six checks PASS. The run
  after documentation edits exited 1 at Go unit-test setup from sandbox Go-cache
  denial, after tracker and all 44 Python tests passed; later checks were not run.
- Permission-enabled `bash scripts/verify.sh` rerun exited 0: tracker validation,
  44 Python tests, `go test -count=1 ./...`, `go test -race -count=1 ./...`,
  `go vet ./...`, `git diff --check` and `git diff --cached --check` all passed.
  Existing M1/M2, Exercise 82 and M3.1-M3.5 regression tests pass locally.
- Only documentation wording and this evidence record were added afterward;
  final tracker validation, whitespace, full documentation diff and status/scope
  checks cover these additions. Only PROGRESS.md and docs/HANDOFF.md changed;
  no staged or untracked files. All other tracked files remain equal to HEAD.
- Ending HEAD remains e0e73da376d55cd22c86fdda8d6011af022a9b0a. No commit, push,
  deployment, live remote check, M3.6 work or M3.5 tracker acceptance occurred.
  These documentation edits await human engineering review; no unresolved blocker.

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

## 11. Historical M3.6 documentation reconciliation (2026-11-11)

Superseded by section 12. Original evidence and status are retained.

Previous documentation checkpoint: `959186046ff3dfd3a86c350027bf7afcd2f509cc`;
M3.5 acceptance and M3.6 parent baseline: `395239d11a17a46eaf6097b8f798585833cdda93`.
Current observed HEAD: `f189ce6be261e2269b8bd266584f123c11a4ce06` on main.
Startup tree/index were clean. Git inspection confirms exactly feature_list.json,
workflow_actions.go and workflow_actions_test.go in the implementation commit.
Earlier pending-M3.5/uncommitted-documentation statements are historical; M3.5
acceptance is committed at 395239d. M3.1-M3.5 are verified. M3.6 is in_progress,
active_task_id M3.6, WIP 1; M3.7-M3.10 remain pending. Tracker evidence fields for
M3.6 remain empty and unchanged; implementation completion is not acceptance.

Initialization exited 0, warning HEAD differs from cached origin/main at 395239d.
No fetch or live remote check occurred here; M3.6 publication is unverified in this
session. These documentation changes are uncommitted and have no new commit SHA.
Only PROGRESS.md and docs/HANDOFF.md are authorized; no source, tests, tracker,
roadmap, commit, push, deployment, M3.6 verified status or M3.7 work is authorized.

### Committed capability and acceptance mapping

`HandlerRegistry.ExecuteActions(WorkflowDefinition, string) (string, error)` validates
the entire definition before callbacks and rejects unknown starting IDs. It follows
Next through InvokeAction(), invoking each reached action once per attempt and
advancing its local cursor only after nil error. An error stops immediately and
returns the failed cursor with the underlying cause wrapped by %w. Empty cursor/nil
ends the action path; nonempty cursor/nil hands off an unexecuted condition or WAIT.
Definitions and handler mappings are not mutated by the method. Invalid input
returns the supplied starting cursor/error; callers must check errors first.

The six committed workflow_actions_test.go functions establish:

- TestHandlerRegistryExecuteActionsOrderAndReadOnly: criterion 1, "Actions execute
  in the required order"; four insertion orders, two attempts each, exact callback
  sequence, terminal path, unreached actions and configuration/registry preservation.
- TestHandlerRegistryExecuteActionsFailure: criterion 2, "A failed action stops the
  normal path and produces an error outcome"; first/middle/final errors, exact
  reached prefixes, failed cursor, errors.Is cause preservation and repeated attempts.
- TestHandlerRegistryExecuteActionsUnknownHandler: unknown/case/padding mismatches
  and zero-value registry fail; later actions do not run after the missing handler.
- TestHandlerRegistryExecuteActionsRejectsInvalidInput: invalid starting cursors and
  malformed definitions, including broken links, unsupported types and cycles,
  fail before callbacks.
- TestHandlerRegistryExecuteActionsBoundary: direct and action-to-condition/WAIT
  handoffs, including a 24-hour delay, without evaluation, sleeping or crossing.
- TestHandlerRegistryExecuteActionsStartsAtSuppliedStep: valid supplied last-step
  entry without replaying preceding steps in that explicitly chosen traversal.

### Historical implementation evidence, separate from fresh verification

Preserved from the prior implementation-session tool outputs/report, at 395239d
plus then-uncommitted work, now committed as f189ce6; not rerun as historical checks:

- Pre-activation tracker validation and M3.6 preflight passed. Human authorization,
  not eligibility, permitted activation. Preflight is not rerun for this active task.
- `go test -count=1 -run '^TestHandlerRegistryExecuteActions' ./...`: sandbox exit 1
  from Go-cache denial; permission-enabled RED exit 1 with ExecuteActions undefined
  at lines 50, 108, 135, 142, 168, 197 and 216 of the test file.
- Its `-v` GREEN command exited 0: six test functions and 18 subtests passed.
- Pre-implementation and post-change unified verifiers hit sandbox cache denial at
  Go unit-test setup after tracker and 44 Python tests passed. Permission-enabled
  reruns exited 0: tracker, 44 Python tests, Go unit tests, race detection, go vet
  and both Git whitespace checks passed. Cache failures were environmental;
  the undefined-method failure was intentional RED evidence.
- gofmt, final tracker/whitespace and scope checks passed; new-file no-index checks
  emitted no whitespace diagnostics (exit 1 means files differ from /dev/null).

### Current limitations and next action

Cursor progress is local, with no persistence, automatic retries, idempotency or
rollback of earlier external side effects; repeated attempts may replay actions.
No timeout, cancellation, panic recovery or concurrent execution ownership guarantee.
No condition evaluation inside ExecuteActions, WAIT scheduling/resume or M2
TryClaim/Complete/Fail integration. Callers/callbacks must not mutate definitions
during traversal and must respect registry concurrency constraints; nil registry
pointers are unsupported. No production deployment, provider integration, measured
business results, durable recovery or distributed exactly-once guarantee is shown.

Next: human engineering review of the two uncommitted documentation files and M3.6
evidence. M3.6 remains in_progress. M3.7-M3.10 remain pending; M3.7 is unauthorized.
Fresh documentation-session checks are recorded below separately from historical
implementation evidence. No unresolved source/tracker discrepancy was found.

### Fresh documentation-session verification (2026-11-11)

At HEAD f189ce6be261e2269b8bd266584f123c11a4ce06 plus these two documentation edits:

- `bash init.sh` and `python3 scripts/validate_features.py feature_list.json`
  exited 0; branch/HEAD/clean-start/history/commit-scope checks passed.
- Initial `bash scripts/verify.sh` exited 1 at go vet from denied sandbox Go-cache
  access, after tracker, all 44 Python tests, Go unit tests and race tests passed.
  The verifier did not reach whitespace checks or print final PASS.
- Permission-enabled `bash scripts/verify.sh` rerun exited 0: tracker, 44 Python
  tests, `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`
  and working-tree/index whitespace checks all passed. Existing M1/M2,
  Exercise 82 and M3.1-M3.6 regression tests pass locally.
- Final whitespace review found an extra blank line at HANDOFF.md EOF; removed.
- Only documentation wording and this evidence record were added afterward.
  Final tracker validation, `git diff --check`, `git diff --cached --check`, full
  documentation diff review and status/scope assertions passed. Only PROGRESS.md
  and docs/HANDOFF.md changed; no staged or untracked files. Source, tests,
  feature_list.json and all other tracked files remain byte-for-byte equal to HEAD.
- Ending HEAD remains f189ce6be261e2269b8bd266584f123c11a4ce06. No commit, push,
  deployment, live remote check, tracker acceptance or M3.7 work occurred.
  These two uncommitted documentation files await human engineering review.

## 12. Historical M3.7 documentation reconciliation (2026-11-11)

Previous documentation checkpoint: `020e37acdab851ab4e71230e939ad1f9120569d7`.
M3.6 acceptance and M3.7 parent baseline: `ec0331523dec72ed1c6f8b7a12681458fd3dcacd`,
the previous published baseline supplied by the user. Current observed starting HEAD:
`10a270486f5df4b591852db90b6264b5e6d5ab2c` on main, with clean startup tree/index.
Git inspection confirms exactly feature_list.json, workflow_wait.go and
workflow_wait_test.go in HEAD. M3.6 acceptance is committed at ec03315, superseding
section 11's pending status. M3.1-M3.6 are verified; M3.7 is in_progress,
active_task_id M3.7, WIP limit 1; M3.8-M3.10 remain pending and unauthorized.
M3.7 tracker evidence fields remain empty and unchanged; implementation is not acceptance.

Initialization exited 0 and warned HEAD differs from cached origin/main. No fetch or
live remote verification occurred; M3.7 publication is unverified here. Only PROGRESS.md
and docs/HANDOFF.md are authorized for editing, uncommitted with no new revision.
No functionality, tests, tracker changes, commit, push, deployment or M3.8 work authorized.

### Committed capability and acceptance mapping

workflow_wait.go exposes ErrWaitNotReady, ErrWaitConsumed, WaitSnapshot, WaitState,
PrepareWait(), WaitState.Snapshot() and HandlerRegistry.ResumeWait(). Exact signatures
and detailed test mapping are in PROGRESS.md's M3.7 checkpoint.
PrepareWait validates the definition and WAIT ID/type, copies the configured step map,
and fixes ResumeAt once using now.Add(step.Delay), without callbacks or sleeping.
The caller retains the in-memory WaitState pointer; Snapshot returns a value copy.
Strictly early resume returns ErrWaitNotReady without changing progress, cursor,
deadline or consumption. Exact-deadline and later resume are eligible if unconsumed.
Eligible resume consumes before callbacks and uses ExecuteActions from the saved
WAIT successor. Earlier completed actions are not replayed in the tested sequential
scenario. Consumed attempts return ErrWaitConsumed without replay. Terminal WAIT
returns empty cursor/nil; condition and further WAIT boundaries remain unexecuted.
Downstream action errors preserve the failed cursor/cause and leave the wait consumed.

- Criterion 1, "Early resume does not advance incorrectly.": TestWaitPreparationAndEarlyResume.
- Criterion 2, "Eligible resume continues from the correct step without replaying earlier completed actions in the tested in-memory scenario.":
  TestWaitResumeWithoutReplay, supported by TestWaitResumeBoundaries and
  TestWaitFailureConsumesWithoutRetry for handoffs and consumed failure handling.
- Criterion 3, "No durable restart guarantee is claimed.": process-local WaitState
  source design and explicit limitations; no test demonstrates restart recovery.
- TestWaitRejectsInvalidInput covers invalid configuration/states;
  TestWaitDefinitionSnapshot covers later caller-map mutation isolation.
  Six test functions and 11 subtests match inspected tests and the implementation report.

### Historical implementation evidence (separate from fresh checks)

The user's supplied implementation-session report records, at ec03315 plus work now
committed as 10a2704: M3.7 preflight passed before activation; human authorization,
not eligibility, permitted implementation. Focused RED failed because the wait/resume
API was undefined; focused GREEN passed six functions and 11 subtests after implementation.
Permission-enabled full verification reruns passed tracker, 44 Python tests, Go unit
and race tests, go vet and both Git whitespace checks. Sandbox Go-cache denials were
environmental, distinct from intentional RED compilation failure. These historical
checks/preflight were not freshly rerun in this documentation session.

### Current limitations and next action

State is in memory only, with no restart durability, automatic scheduling/background
polling, cloud integration, database or external storage. Retain one WaitState pointer;
single-caller use only, no copying or concurrent ResumeWait/Snapshot access. Do not
mutate definitions concurrently with PrepareWait. Preparation validates configuration,
not proof of prior action execution; repeated preparation creates independent states
without deduplication. Consumption prevents ordinary sequential replay, not exactly-once
external effects. Downstream failures remain consumed and need separately defined
recovery/reconciliation. Panics are not recovered and can leave an imprecise failed
cursor. Callbacks can block; no cancellation/timeout added. Handler registrations and
captured callback inputs are not snapshotted; registry constraints still apply.
No M2 TryClaim/Complete/Fail integration, production deployment, measured business
outcome or distributed correctness guarantee. M1/M2 invariants remain separate.

Next: human engineering review of these two documentation edits and M3.7 evidence.
M3.7 remains in_progress; M3.8-M3.10 remain pending and unauthorized. No task activation
or tracker acceptance occurs here. Fresh verification is recorded separately below.

### Fresh M3.7 documentation-session verification (2026-11-11)

At HEAD 10a270486f5df4b591852db90b6264b5e6d5ab2c plus the two documentation edits:

- `bash init.sh`, standalone `python3 scripts/validate_features.py feature_list.json`,
  startup branch/HEAD/clean-tree/history/implementation-file-set checks passed.
- Initial `bash scripts/verify.sh` exited 1 at Go unit-test setup from sandbox
  Go-cache denial after tracker and 44 Python tests passed. No race/vet/whitespace
  checks or final PASS followed that failure.
- Permission-enabled rerun exited 0 with all six checks PASS: tracker, 44 Python
  tests, Go unit tests, race detection, go vet and working-tree/index whitespace.
  These are fresh local regression results, separate from historical RED/GREEN evidence.
- Only these verification records were added afterward. Final tracker validation,
  `git diff --check`, `git diff --cached --check`, complete two-file diff review,
  `git status --short` and scope/preservation assertions passed. Exactly PROGRESS.md
  and docs/HANDOFF.md are modified; no staged/untracked files. Source, tests,
  feature_list.json and all other tracked files remain byte-for-byte equal to HEAD.
- Ending HEAD is unchanged: 10a270486f5df4b591852db90b6264b5e6d5ab2c. M3.7 remains
  in_progress/active, WIP 1; M3.8-M3.10 pending and unauthorized. No commit, push,
  deployment, live remote check or acceptance occurred. Stop for human engineering
  review of these uncommitted documentation changes; no unresolved blocker found.

## 13. Current M3.8 implementation checkpoint (2026-11-12)

Previous documentation checkpoint: `9c57598597a1d081cc6fbc60da1c683a976a3478`.
M3.7 acceptance commit: `4545eef`, inspected in local history, superseding section 12.
Starting/current observed HEAD on main: `842000fabc0b203a75481b1dcc42d6cb1a9b4eb8`.
This matches the user-supplied baseline; startup tree/index were clean. Initialization
passed and reported cached origin/main matches HEAD. No live remote access occurred.

The user explicitly authorized activation and implementation of M3.8 only. Tracker
validation and M3.8 preflight passed before activation; eligibility was not authority.
M3.1-M3.7 remain verified; M3.8 is in_progress/active; WIP 1; M3.9-M3.10 pending.
Tracker acceptance criteria, objective and evidence fields are unchanged. No acceptance,
commit, push, deployment, M3.9 or M3.10 work occurred in this session.

Uncommitted files: workflow_runner.go, workflow_runner_test.go (new), feature_list.json
(activation only), PROGRESS.md and docs/HANDOFF.md (checkpoint reconciliation).
Existing source/tests, AGENTS.md, roadmap, scripts and README are unchanged.

### Integration and review evidence

New API:
`(*HandlerRegistry).RunActions(repo WorkflowLifecycleRepository, executionID, workerID string, d WorkflowDefinition) (string, error)`.
The lifecycle interface includes TryClaim/Get/Complete/Fail only; the in-memory M2
repository supplies them. Claim precedes callbacks; Get validates the owned RUNNING
snapshot and captures the version before ExecuteActions. Terminal success requires
empty cursor/nil plus accepted Complete. Execution errors call Fail with nonempty
context, retain the failed cursor and wrap the original cause. Rejected terminal
operations return reconciliation errors without version refresh, retry or fallback.

A nonempty cursor/nil hands off an unexecuted condition/WAIT and leaves RUNNING with
unchanged claim metadata. Retain cursor/definition locally. No condition evaluation,
PrepareWait, lifecycle-aware resume/finalization, scheduler, WAITING status or durable
waiting is added. Existing helpers remain separate; calling RunActions again cannot
resume the RUNNING execution. This is the bounded boundary treatment, not completed
end-to-end orchestration. The execution/definition association belongs to the caller.

Acceptance mapping is recorded in PROGRESS.md's M3.8 section. The six new tests cover
success, first/middle/final action errors, claim rejection, terminal owner/version,
exact status/version/attempt/history invariants, nonempty failure diagnostics,
rejected/stale terminal operations, invalid claimed snapshots, validation/handler
errors and direct/action-to-condition/WAIT handoffs (24 subtests).

Fresh commands/results at HEAD plus these uncommitted changes:

- `bash init.sh`, branch/HEAD/status/log checks, tracker validation and pre-activation
  `python3 scripts/task_preflight.py --task M3.8`: exit 0. Tracker validation after
  activation and before full verification also passed.
- RED `go test -count=1 -run '^TestLifecycleRunner' ./...`: exit 1, RunActions undefined
  at test lines 89, 149, 183, 224, 275, 295; no sandbox/cache failure in this run.
- GREEN `go test -count=1 -run '^TestLifecycleRunner' -v ./...`: exit 0; six functions,
  24 subtests. Both new Go files formatted with gofmt.
- Initial `bash scripts/verify.sh`: exit 1 at Go unit-test setup from denied sandbox
  Go-cache access after tracker and 44 Python tests passed; later checks not reached.
- Permission-enabled `bash scripts/verify.sh`: exit 0; tracker, 44 Python tests,
  Go unit tests, race detection, go vet and working-tree/index whitespace all PASS.
  Only documentation changed afterward; final tracker/formatting/whitespace/scope
  checks cover those additions. These are local results, not cloud/E2E evidence.

Limitations: process-local action-path execution only; no durable progress/recovery,
lease, retry, compensation, panic recovery, cancellation/timeout or cloud integration.
A lost handoff, unusable claim snapshot or callback panic can leave RUNNING. Callbacks
may have side effects before a rejected terminal write: reconcile rather than replay.
Nonnil initialized dependencies and existing registry/definition concurrency rules
are required. No exactly-once external-side-effect or production guarantee is claimed.
Next action: human engineering review of the uncommitted implementation and evidence.
M3.8 must remain in_progress; no commit/push/deployment or later task is authorized.
