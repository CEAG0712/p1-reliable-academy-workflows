# P1 — Durable Session Handoff

## 1. Project and business purpose

Build a reusable Go workflow execution engine. Student enrollment and
course access are the first use case; keep business logic outside the
generic engine. This is an in-memory engineering portfolio prototype.

## 2. Repository identity and checkpoint revision

Current checkpoint: section 16 (Stage 4 tracker preparation only).
Earlier revisions/status below are historical; M3.1-M3.10 are now verified.

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

## 13. Historical M3.8 implementation checkpoint (2026-11-12)

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


## 14. Historical M3.9 verification checkpoint (2026-11-12)

Previous implementation/documentation checkpoint revision:
`8a049e88b1d63db73e15e120b17c750179f51868`, which recorded its own start at
`842000fabc0b203a75481b1dcc42d6cb1a9b4eb8`. M3.8 acceptance commit and this
session's observed starting/ending HEAD on main:
`7811e9d08f8a10a4d18da6a74d4209570971daa9`. Startup tree/index were clean.
Initialization passed: cached origin/main matched; no fetch/live remote check.
Local history, source/tests and tracker supersede section 13's pending M3.8 status.

The human explicitly authorized M3.9 activation and implementation only. Pre-activation
tracker validation and M3.9 preflight exited 0; eligibility was not authorization.
M3.9 status is now in_progress and active_task_id M3.9, WIP 1. Acceptance criteria
and evidence fields are unchanged. M3.1-M3.8 remain verified. M3.10 remains pending
and unauthorized. No verified status, staging, commit, push or deployment occurred.

Uncommitted/unstaged files: workflow_runner_concurrency_test.go (new),
feature_list.json (activation only), PROGRESS.md and docs/HANDOFF.md. No production
source, existing tests, scripts, roadmap, README or AGENTS.md changed. No production
defect was found; no fix or architecture expansion was needed.

### Test mapping and evidence

- A/B: TestRunnerConcurrencySuccessfulExecutions and TestRunnerConcurrencyFailedExecutions
  launch ten distinct workers against one repository/initialized registry/immutable
  definition using a start barrier. Exactly one action path, nine claim rejections,
  exact callback order without replay, winning owner and terminal lifecycle are
  asserted. Failure wraps a sentinel recognized by errors.Is and stops subsequent
  actions. Final version 3, attempts 1, exact claim/terminal history and matching
  diagnostic reason are checked; no completion is recorded on action failure.
- C: TestRunnerConcurrencyDuplicateAfterCompletion rejects same-owner and other-owner
  duplicates without new callbacks or any snapshot changes.
- D: TestRunnerConcurrencyInvalidOwnership rejects non-owner Complete and Fail
  against the actual lifecycle contract without mutation, then permits owner completion.
- E: reuse unchanged TestLifecycleRunnerTerminalRejection: four rejected/stale terminal
  scenarios check cause preservation, visible rejection, original owner/version,
  exact operation sequence without refresh/retry/fallback and accepted stored state.
- F: TestRunnerConcurrencyBlockedCallback checks a rejected contender while the winner
  is blocked RUNNING, then normal completion. Channels, atomic counters, five-second
  deadlines and cleanup release/join avoid sleeps and unbounded test synchronization.
- Together these map to M3.9's no-duplicate-action and correct failure/lifecycle criteria.
  The no-distributed-exactly-once criterion is preserved in the limitations below.
  Implementation evidence is separate from human acceptance.

Fresh commands at HEAD plus these uncommitted changes:

- `bash init.sh`, branch/HEAD/status/history checks, tracker validation before/after
  activation and `python3 scripts/task_preflight.py --task M3.9`: exit 0.
- First `go test -count=1 -run '^TestRunnerConcurrency' -v ./...`: exit 0,
  five test functions and two subtests. Baseline GREEN; existing M3.8 behavior already
  satisfies the added verification. No behavioral RED or production fix occurred.
- Initial `go test -count=1 -run '^TestLifecycleRunner' -v ./...` and
  `go test -race -count=1 -run '^TestRunnerConcurrency' ./...`: exit 1 at setup,
  sandbox Go-cache access denied. Permission-enabled reruns exited 0: lifecycle's
  six functions/24 subtests and new concurrency tests passed. Environmental failures
  are distinct from test failures.
- `go test -race -count=20 -run '^TestRunnerConcurrency' ./...`: exit 0.
- `bash scripts/verify.sh`: first invocation exited 0; tracker, 44 Python tests,
  Go unit tests, full race checks, go vet and both Git whitespace checks PASS.
  Only documentation changed afterward; final standalone tracker, formatting,
  whitespace (including the untracked Go file), scope/preservation and status checks
  passed. Exactly the four authorized files changed; the index is empty and HEAD
  unchanged. Production source and existing tests remain byte-for-byte equal to HEAD.

Limits: one shared process-local in-memory repository and only exercised scenarios.
No exactly-once distributed execution, exactly-once external provider effects, safe
retry after every timeout, durable recovery, cross-process exclusivity or production
readiness. Definitions must remain immutable, registration must precede concurrent
registry use, and captured callback data needs its own synchronization. No scheduler,
lease, retries, persistence, cloud/provider integration or panic recovery added.
A lost handoff or panic may leave RUNNING. Ambiguous side effects before a rejected
terminal update require reconciliation rather than blind replay.

Next action: stop for human engineering review of the complete uncommitted diff.
Keep M3.9 in_progress and M3.10 pending. Do not stage, commit, push, deploy, mark
verified or begin M3.10. Detailed acceptance mapping and command outcomes are also
recorded in PROGRESS.md's M3.9 checkpoint. No current behavioral blocker found.


## 15. Historical M3.10 implementation checkpoint (2026-11-12)

Previous implementation/documentation checkpoint revision:
`5909abf099c888536f1bda88ca2d4614a45a4508`; it recorded its own start at
`7811e9d08f8a10a4d18da6a74d4209570971daa9`. Current observed starting/ending
HEAD on main: `e15efdf2da602bd67de54adb5bd8ddf7f8dc3f2f`, the required baseline
and locally inspected M3.9 acceptance commit. Startup tree/index were clean.
M3.9 implementation is committed at 5909abf; its acceptance at e15efdf supersedes
section 14's pending status. Initialization exited 0, cached origin/main matched;
no fetch/live remote verification occurred. README/roadmap status notes lag source,
tests, history and the accepted M3.1-M3.9 tracker; section 13 remains approved scope.

The human explicitly authorized M3.10 activation and implementation ONLY. Tracker
validation and pending-task preflight exited 0 before activation. Preflight establishes
eligibility, not permission. Tracker changes only active_task_id null -> M3.10 and
M3.10 pending -> in_progress. WIP 1, all earlier records/evidence, M3.10 criteria,
objective/dependencies/evidence and unrelated metadata are preserved.

### Shared engine and executable evidence

New `workflow_runner_reusability_test.go` configures two valid WorkflowDefinitions
in one initialized HandlerRegistry and runs independent identities through one
actual InMemoryWorkflowRepository using unchanged RunActions/ExecuteActions code.
All handlers and simulated business effects are deterministic and local to tests.

- Enrollment fulfillment (`enrollment-fulfillment`, purchase.confirmed):
  validate-purchase/enrollment.validate -> prepare-access/enrollment.prepare ->
  record-fulfillment/enrollment.record -> empty terminal cursor. Three actions
  validate student/course/paid inputs, prepare student-17:go-basics and record
  prepared:student-17:go-basics. No provider access is granted.
- Learning analytics (`learning-analytics`, report.requested):
  collect-activity/analytics.collect -> aggregate-activity/analytics.aggregate ->
  render-report/analytics.render -> record-publication/analytics.publish -> empty
  terminal cursor. Four actions collect [2,3,5], sum 10, render and record
  published:learning activity: 10 in memory. No report service is contacted.

Different lengths, handler structures, step names, inputs, outputs, execution IDs
and worker IDs demonstrate more than renamed copies. Trigger fields are configuration;
the new tests call the runner directly, without claiming event-dispatch integration.

Acceptance mapping:

1. Both workflows run without engine-core changes: TestRunnerReusabilityWorkflows
   (enrollment-first/analytics-first) checks both orders, exact callbacks, business
   outputs, isolated snapshots, READY -> RUNNING -> COMPLETED, versions 1/2/3,
   owner, attempts and exact history. TestRunnerReusabilityActionShapes checks the
   three-/four-action paths independently via ExecuteActions and empty terminal
   cursors. TestRunnerReusabilityIndependentFailure checks analytics errors.Is,
   failed cursor/cause/diagnostics, FAILED history, no later actions, untouched
   enrollment, then enrollment success without changing the failed execution.
   Source and diff review show no production Go or existing test modifications.
2. M1/M2 tests/race checks: verbose full Go regression was inspected for actual
   enrollment, execution, repository (including Exercise 82) and all M3 test groups;
   all passed. Full race and canonical verification passed.
3. Limitations: documented below and in PROGRESS.md's M3.10 section.

### Fresh commands and outcomes

All startup branch/HEAD/status/index/history checks, bash init.sh, tracker validation
before/after activation and `python3 scripts/task_preflight.py --task M3.10` exited 0.
`gofmt -w workflow_runner_reusability_test.go` exited 0. First focused
`go test -count=1 -run '^TestRunnerReusability' -v ./...` exited 0: three functions,
four subtests. **BASELINE GREEN** against unchanged production code; no artificial
RED or production fix was necessary.

Initial standalone commands below each exited 1 from sandbox Go-cache denial
(test setup failures); permission-enabled reruns each exited 0:

- `go test -count=1 -run '^TestLifecycleRunner' -v ./...`: six functions/24 subtests.
- `go test -count=1 -run '^TestRunnerConcurrency' -v ./...`: five functions/two subtests.
- `go test -race -count=1 -run '^TestRunnerConcurrency' ./...`.
- `go test -count=1 -v ./...` (all M1/M2/M3 suites actually selected).
- `go test -race -count=1 ./...`.
- `go vet ./...`.

`bash scripts/verify.sh` exited 0 on its first invocation: tracker, all 44 Python
harness tests, Go unit tests, full race checks, vet and both Git whitespace checks
PASS. Its Go unit command is `go test -count=1 ./...`; Python runs with bytecode
writes disabled. Environment failures are separate from behavioral test outcomes.
Only checkpoint docs changed after the verifier. Final audit evidence and complete
patches are available at `/private/tmp/p1-m3.10-review/` outside the repository.

### Stage 3 limitations and review boundary

This is an in-memory engineering prototype. Reusability of tested configurations
and callbacks through unchanged engine code is demonstrated; production readiness
is not. No cross-process coordination, durable execution persistence, durable WAIT/
condition cursor recovery, automatic retries or compensation, crash-safe external
side effects, exactly-once provider effects, guaranteed callback-panic recovery,
production-grade scheduler, deployment or measured business results are claimed.

Definitions must remain immutable during execution. Registration must finish before
concurrent registry use; do not copy populated registries. Mutable callback state
needs synchronization. These new domain fixtures use sequential execution. Verified
concurrency guarantees apply only to the tested shared in-memory repository behavior
within one process. Ambiguous external effects need reconciliation. A lost handoff
or panic can leave RUNNING; condition/WAIT helpers do not automatically resume and
finalize lifecycle state. Callers associate execution identities with definitions.
No distributed infrastructure, external integrations or Stage 4 work was introduced.

Final repository changes, all unstaged/uncommitted: workflow_runner_reusability_test.go
(new), feature_list.json (activation only), PROGRESS.md and docs/HANDOFF.md.
M3.10 remains in_progress, active_task_id M3.10; M3.1-M3.9 remain verified; WIP 1.
No M3.10 acceptance commit exists. Nothing was staged, committed, pushed or deployed.
Next action: STOP for human engineering review of the complete diff. No final verified
status, commit, publication or Stage 3 acceptance without separate human authorization.
At next startup, reread AGENTS.md/PROGRESS.md/this checkpoint, run init.sh, inspect
actual HEAD/tree/index/history, validate tracker and reconcile code/tests/evidence.
Do not mistake the previous checkpoint revision or these uncommitted changes for
an implementation commit. Do not rerun pending-task preflight as if M3.10 were
still pending; it is now active and no other task is authorized.


## 16. Current Stage 4 tracker preparation checkpoint (2026-11-13)

Previous documentation checkpoint revision: `ae17f66b6c91a394050eb4fc60c33087e064e6a3`.
Current observed HEAD on main: `5457b77e7835cf951b1c71ab332080d7b81fe735`,
`chore: finalize M3.10 and Stage 3 acceptance`. The startup tree and index were
clean. Live remote main matched via permission-enabled `git ls-remote origin
refs/heads/main` (exit 0); the sandbox attempt exited 128 from DNS failure.
`bash init.sh` and tracker validation exited 0. The initialization cached-ref
message alone is not live evidence. Source/tests, history and all ten verified
M3 records supersede section 15's pre-acceptance status. README and roadmap
historical status notes lag; their future API/trust decisions remain unresolved.

The human authorized Stage 4 TRACKER PREPARATION ONLY. The unified schema-v1
tracker contains exactly 20 ordered tasks: ten structurally identical accepted
M3 records, then pending M4.1-M4.10 contracts from the supplied request. All M4
evidence is null/empty, blockers null, global WIP 1, active_task_id null. M4.1
depends on verified M3.10; every later M4 task depends on the preceding M4 task.
Status meanings/rules are unchanged. The API task IDs follow the supplied
breakdown; the roadmap's older Stage 3 M4 enrollment label is separate. Gin
v1.11.0 remains an unapproved compatibility candidate, not an installed package.
No HTTP status/body/method or authentication/payment-trust decisions were invented.

Only these six allowed files are changed, unstaged/uncommitted:
feature_list.json (metadata and appended contracts), AGENTS.md (M3/M4 preflight,
selection and global WIP), scripts/test_validate_features.py (seven tests),
scripts/test_task_preflight.py (four tests), PROGRESS.md and docs/HANDOFF.md
(truthful current planning checkpoint). Validator/preflight/verifier implementations,
Go source/tests, go.mod/go.sum, API, infrastructure and roadmap are unchanged.

Baseline verifier: initial sandbox exit 1 at Go unit setup after 44 Python tests
passed, from denied Go-cache access; permission-enabled rerun exit 0, all six checks
PASS. Fresh post-extension validation exit 0; Python unittest discovery with
PYTHONDONTWRITEBYTECODE=1 exit 0, 55 tests (all 44 unchanged fixtures plus 11 new).
M4.1 preflight exit 0: ELIGIBLE, HUMAN APPROVAL REQUIRED, no activation or authority.
M4.2 preflight exit 1: expected negative, depends on unverified M4.1. Standalone
Go unit/race/vet commands initially exited 1 from sandbox Go-cache denial; their
permission-enabled reruns each exited 0. Environmental failures are not behavioral
RED. Exact commands are recorded in PROGRESS.md's current checkpoint and the
external review report. Nonverbose Go output does not provide test counts.

Preservation regression reads the immutable local published tracker via git show
5457b77e7835cf951b1c71ab332080d7b81fe735:feature_list.json without network access.
Preparation invariants require empty evidence for the proposed pending M4 tasks;
schema v1 generally permits populated pending evidence, explicitly tested here.
It enforces null nonblocked blockers, one global in_progress slot and matching
active_task_id. Preflight remains read-only and cannot establish human authorization.
Original M3 fixture methods/order were not edited. Existing engine limitations in
section 15 remain applicable; harness success establishes no API/cloud capability.

Read-only baseline snapshot and complete review material:
`/private/tmp/p1-stage4-tracker-review/`. Recheck actual HEAD/tree/index, tracker,
source/tests and verification at startup; do not treat uncommitted edits as a commit.
Next action is human engineering review. No M4 task activated or accepted. Nothing
staged, committed, pushed or deployed; no Gin installed and no API code implemented.
Stop after tracker preparation; no M4.1 implementation or later task is authorized.

Final preparation verification and preservation audit:

- Post-change `bash scripts/verify.sh`: initial sandbox exit 1 at go vet cache
  access after tracker, 55 Python tests and Go unit/race checks passed;
  permission-enabled rerun exit 0, all six checks PASS including vet/whitespace.
- Separate structural audit against the complete published tracker and the
  read-only pre-edit snapshot passed all 14 required preservation checks:
  M3 objects/evidence/order/status, M4 pending/empty evidence/null blockers,
  schema/global WIP/no active task, exact chain/count and permitted file scope.
  Additional checks passed for unchanged HEAD, empty index, no untracked files,
  unchanged status meanings/rules, original fixture method ASTs and every supplied
  M4 title/objective/dependency/acceptance criterion.
- Only these checkpoint verification records were added after the passing verifier.
  Final tracker validation, both whitespace checks, status/stat and complete patch
  export follow those records; no behavior or tracker contract changed afterward.
  Review patch: `/private/tmp/p1-stage4-tracker-review/complete-review.patch`.
  Full report: `/private/tmp/p1-stage4-tracker-review/review-report.txt`.
  Separate audit: `/private/tmp/p1-stage4-tracker-review/preservation-audit.txt`.


## M4.1 Gin Foundation implementation checkpoint (2026-11-13)

Previous documentation checkpoint revision: `ed161d7ced6061dac105f67c3fbe10137030a44c`
(the tracker-preparation publication). That checkpoint recorded its own starting
HEAD as `5457b77e7835cf951b1c71ab332080d7b81fe735`; it is historical, not the
current HEAD. Current observed starting/ending HEAD on main is `ed161d7ced6061dac105f67c3fbe10137030a44c`.
Live remote main matched that required baseline using permission-enabled
`git ls-remote origin refs/heads/main` (exit 0). Startup tree and index were clean.
Published tracker, Git history and fresh tests supersede the older pending M3.10
and uncommitted tracker-preparation statements. README/roadmap historical status
still lags the accepted M3 tracker. No README or roadmap change was authorized.

Human explicitly authorized M4.1 activation and bounded implementation ONLY.
`bash init.sh`, tracker validation, M4.1 preflight and baseline canonical verifier
exited 0. Preflight printed ELIGIBLE and HUMAN APPROVAL REQUIRED; the supplied
request, not eligibility, supplied authorization. Baseline Python count is 60,
not the historical checkpoint's 55. Activation changed only active_task_id to
M4.1 and its status to in_progress, preserving empty evidence until actual runs.
Current M4.1 evidence records observed results; reviewed_by remains null.
M3.1-M3.10 objects/evidence remain unchanged and verified, M4.2-M4.10 unchanged
and pending, schema 1 and WIP 1. Human acceptance remains pending.

### Implementation and acceptance evidence

- `internal/httpapi/router.go`: independent NewRouter, gin.New plus Recovery,
  GET /health only. Returns 200, application/json; charset=utf-8, semantic JSON
  exactly {"status":"ok"}. No listener, engine imports or business dependencies.
- `internal/httpapi/router_test.go`: TestHealth, TestUnknownRoute,
  TestHealthUnsupportedMethod, TestHealthRepeatedCalls. Default Gin 404 for
  GET /does-not-exist and POST /health; exact rejection body excludes health
  success. Ten calls reuse one router. JSON decoding rejects extra fields,
  malformed/trailing JSON and incorrect types. Tests use httptest, no sockets.
- `cmd/p1-api/main.go`: tcp 127.0.0.1:8080 only, standard http.Server with
  five-second ReadHeaderTimeout. Injectable net.Listen seam; listen/serve errors
  wrap the cause and include operation/address. main logs and exits on failure.
- `cmd/p1-api/main_test.go`: TestRunListenFailure and TestRunServeFailure inject
  deterministic errors without public/local sockets, sleeps or port selection.
  Assert loopback target, errors.Is, diagnostic context and listener closure.
- `go.mod` and new `go.sum`: direct Gin v1.11.0 plus required transitive modules.
  Actual toolchain go1.24.5 darwin/arm64; Gin's own go.mod requires 1.23.0.
  Resolved metadata covers 55 modules; all reported Go requirements <= 1.24.5.
  Project go 1.24.5 retained, no toolchain directive, no proxy/system changes.
  Only Gin is a direct third-party requirement. go mod verify exits 0.
- HTTP import graph excludes the root engine package, so neither health nor
  unknown/method rejection can invoke workflow repositories/actions, executions,
  enrollment providers or payment state. This is structural evidence combined
  with handler tests, not a claim of end-to-end workflow integration.

### Actual command results

- RED `go test -count=1 -v ./internal/httpapi/... ./cmd/p1-api/...`: exit 1
  before implementations existed. Permission-enabled run showed undefined
  NewRouter at four sites and run at two sites. Earlier sandbox attempt also
  showed net/http/internal/testcert setup failure; it is environmental, not RED.
- GREEN same focused command: exit 0, all six test functions passed.
- `go test -race -count=1 ./internal/httpapi/... ./cmd/p1-api/...`: exit 0.
- `go mod download -json github.com/gin-gonic/gin@v1.11.0`,
  `go get github.com/gin-gonic/gin@v1.11.0`, `go mod tidy`, and
  `go list -m -json all`: permission-enabled exits 0. Original download/list
  attempts hit sandbox DNS, original tidy hit module/build-cache permissions.
  Original live-remote attempt also hit DNS (exit 128). Environmental failures
  were resolved through permission-enabled reruns, not classified as incompatibility.
- `python3 scripts/validate_features.py feature_list.json`: exit 0 before/after
  activation; Python discovery with PYTHONDONTWRITEBYTECODE=1: exit 0, 60 tests.
- Standalone `go test -count=1 ./...`, `go test -race -count=1 ./...`,
  `go vet ./...`: each exit 0 with normal cache access, all three packages pass.
- Post-implementation `bash scripts/verify.sh`: exit 0, all six checks PASS,
  including tracker, 60 Python tests, Go unit/race/vet and both whitespace checks.
- `go mod verify`, `go list -deps ./internal/httpapi ./cmd/p1-api`,
  `git diff --check` and `git diff --cached --check`: each exit 0.
- Source formatted with gofmt. Separate untracked-file whitespace, formatting,
  tracker-preservation and permitted-file audits accompany the external patch.
  Only tracker evidence/checkpoint documentation changed after regression;
  final validation and whitespace/audit results are in the review materials.

### Review boundary and limitations

Modified: feature_list.json, go.mod, PROGRESS.md, docs/HANDOFF.md.
New/untracked: go.sum, internal/httpapi/router.go, internal/httpapi/router_test.go,
cmd/p1-api/main.go, cmd/p1-api/main_test.go. Existing engine files/tests, harness
scripts, README, roadmap and AGENTS.md are unchanged. No enrollment endpoint,
workflow invocation, payment API, production/cloud access or deployment config.
Nothing staged, committed, pushed or deployed. HEAD remains the starting SHA.

This establishes only a locally tested HTTP foundation. It does not establish
enrollment HTTP invocation, real payment authorization, real provider integration,
durable execution state, reliable asynchronous orchestration, cloud deployment,
production security or distributed exactly-once processing. Engine limitations
in earlier checkpoints remain applicable. Health indicates HTTP responsiveness
only. HTTP recovery cannot undo side effects or repair persisted workflow state.
No production lifecycle/shutdown system or end-to-end enrollment test is claimed.

Complete review patch (including every untracked file), report, module metadata,
command logs and preservation audit: `/private/tmp/p1-m4.1-review/`.
Next action: STOP for human engineering review. Keep M4.1 in_progress and M4.2
pending; do not stage, commit, push, deploy or mark verified. At next startup,
reconcile actual HEAD/tree/index/tracker/source/tests and this uncommitted
checkpoint. Do not treat the review patch or current edits as a published commit.


## M4.1 acceptance reconciliation checkpoint (2026-11-13)

Published baseline and current observed HEAD on main:
`ed161d7ced6061dac105f67c3fbe10137030a44c`.
Previous committed documentation checkpoint revision: the same `ed161d7`.
The immediately preceding M4.1 implementation checkpoint above remains
uncommitted; its complete reviewed patch is `/private/tmp/p1-m4.1-review/complete-review.patch`.
Neither that checkpoint nor this acceptance update has an implementation or
acceptance commit. Live remote main independently matched the full baseline
through permission-enabled `git ls-remote origin refs/heads/main` (exit 0).
Initial sandbox lookup exited 128 from DNS restriction. Initialization exited 0;
its warning about existing changes was reviewed. Those changes are expected:
all nine modified/new files exactly matched the previously reviewed complete
patch before this reconciliation. Index empty, no unrelated changes. No pending-task
preflight was rerun against the already-active M4.1.

Human engineering owner explicitly authorized acceptance reconciliation following
ChatGPT source-level review of the supplied patch against all seven M4.1 criteria.
Codex independently inspected source/tests and ran the fresh local commands below.
ChatGPT is not claimed to have executed local commands; the human reviewer is not
claimed to have personally run tests. Human source-review acceptance is distinct
from Codex command evidence. The older pending-review/in_progress statements above
are preserved as history and superseded by this checkpoint.

Acceptance mapping, without implementation changes:

1. go.mod pins Gin v1.11.0 and retains Go 1.24.5 without a toolchain directive.
   Fresh go mod verify passed. Resolved metadata contains 55 modules, all reported
   Go requirements <= 1.24.5; actual installed toolchain is go1.24.5 darwin/arm64.
2. internal/httpapi/router.go NewRouter constructs routes without listening;
   cmd/p1-api/main.go owns listener startup. Existing engine files remain unchanged.
3. TestHealth asserts 200, parsed application/json and exact semantic {"status":"ok"}.
4. TestRunListenFailure/TestRunServeFailure assert wrapped causes and diagnostic
   operation/address; tcp 127.0.0.1:8080, five-second ReadHeaderTimeout, listener closure.
5. TestUnknownRoute and TestHealthUnsupportedMethod assert default 404;
   TestHealthRepeatedCalls checks ten calls through one router. httptest and
   injected listener failures require no listening ports or sleeps.
6. Fresh Go/Python/race/vet and canonical checks passed as recorded below.
7. Source and import graph show no engine dependency, enrollment/payment endpoint,
   provider integration, lifecycle modification, production exposure or M4.2 work.

Fresh pre-acceptance commands, each exit 0 on the successful run:

- go mod verify
- go test -count=1 -v ./internal/httpapi/... ./cmd/p1-api/... (six test functions)
- go test -race -count=1 ./internal/httpapi/... ./cmd/p1-api/...
- go test -count=1 ./... (all three packages)
- go test -race -count=1 ./... (all three packages)
- go vet ./... (no diagnostics)
- PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
  (60 tests)
- bash scripts/verify.sh (all six canonical checks PASS)
- go list -m -json all
- go list -deps ./internal/httpapi ./cmd/p1-api

Initial sandbox focused unit/race and module verification passed, but full Go
unit setup exited 1 from denied root-package build-cache access. Permission-enabled
rerun of the required checks passed. This environmental error is not a behavioral
failure or new RED evidence. No source fix, dependency install or toolchain change.

Only after the required checks passed, M4.1 transitioned in_progress -> verified
and active_task_id became null. Its original evidence/result/commands were preserved
and fresh results appended with explicit attribution. All ten M3 task objects and
evidence are unchanged/verified; M4.2-M4.10 unchanged/pending. Schema 1, WIP 1,
all acceptance criteria/dependencies and unrelated metadata preserved; no task
in progress. M4.2 structural eligibility grants no activation authority.

This reconciliation edits only feature_list.json, PROGRESS.md and docs/HANDOFF.md.
Entire working tree versus HEAD still has exactly these modified paths:
PROGRESS.md, docs/HANDOFF.md, feature_list.json, go.mod; and these untracked paths:
go.sum, internal/httpapi/router.go, internal/httpapi/router_test.go,
cmd/p1-api/main.go, cmd/p1-api/main_test.go. The six implementation/dependency
files are unchanged from the reviewed patch. All other repository files are
unchanged from the published baseline. No staging, commit, push, deployment,
implementation modification or M4.2 activation performed or authorized.

Local HTTP foundation only: no enrollment HTTP invocation, real payment/provider
integration, durable execution state, reliable asynchronous orchestration, cloud
deployment, production security or distributed exactly-once guarantee. Health
means HTTP responsiveness; engine limitations and previous RED/GREEN evidence
remain applicable. No unresolved behavioral defect was found during reconciliation.

Post-update verification results and exact exit codes, independent preservation
checks, the three-file acceptance-only diff and full nine-file review patch
(including all untracked contents) are retained outside the repository under
`/private/tmp/p1-m4.1-acceptance-review/`. Recheck those artifacts for final results.
Next action: STOP for human publication review; do not stage, commit, push,
deploy, activate M4.2 or implement more API functionality. At next startup,
reconcile actual HEAD/tree/index/tracker and fresh evidence, preserving the
distinction between accepted local source and published commits.


## M4.2 API Contracts implementation checkpoint (2026-11-13)

Previous committed documentation checkpoint revision: `eac1f8bf7a14f6853036d91071c6b5881cd373ca`
(the M4.1 publication; its embedded acceptance narrative began at `ed161d7`).
Current observed starting HEAD on main: `eac1f8bf7a14f6853036d91071c6b5881cd373ca`.
Startup tree/index clean. Live remote main independently matched that full SHA:
initial sandbox `git ls-remote origin refs/heads/main` exit 128 (DNS),
permission-enabled retry exit 0. Source, Git history, accepted tracker and fresh
checks supersede older uncommitted M4.1/pending-review text. README/roadmap
historical implementation notes lag; the human supplied the bounded M4.2 decisions.

The supplied human request explicitly authorizes activation and implementation of
M4.2 ONLY. AGENTS.md requires authorization, not a second approval after an
explicit request. Preflight eligibility was not treated as authority.
`bash init.sh`, tracker validation, M4.2 preflight and baseline
`bash scripts/verify.sh` each exited 0. Preflight reported ELIGIBLE / HUMAN
APPROVAL REQUIRED. Baseline verifier passed all six checks and 60 Python tests.
Activation changed only active_task_id and M4.2 status, retaining empty evidence
until actual implementation/test results existed. Human source review is pending;
reviewed_by remains null, M4.2 remains in_progress and active_task_id is M4.2.

Implementation, exactly six changed/new repository files:

- `internal/httpapi/contracts.go` (new): CreateExecutionRequest,
  CreateExecutionResponse, ExecutionStatusResponse, APIError and ErrorResponse;
  exact approved JSON fields, string statuses, uint64 numeric version.
- `internal/httpapi/contracts_test.go` (new): seven required test functions;
  semantic exact JSON, full uint64 precision, four persisted status examples,
  public-field allowlists, health and only-health route-table/404 checks.
- `docs/API_CONTRACTS.md` (new): implemented health versus future POST creation
  and GET status resources; explicit HTTP/error matrix, ordinary validation 400,
  conditional 201, reserved 202, FAILED retrieval may be 200. Documents identity,
  authorization, tenant isolation, payment trust and insertion/duplicate blockers.
- `feature_list.json`: M4.2 activation and actual evidence only.
- `PROGRESS.md` and `docs/HANDOFF.md`: this factual implementation checkpoint.

RED -> GREEN evidence:

- Before contracts.go, direct `go test -count=1 -v ./internal/httpapi/...`
  exited 1 with undefined CreateExecutionRequest, CreateExecutionResponse,
  ExecutionStatusResponse, ErrorResponse and APIError. Original output is in the
  session; a labelled transcription is `/private/tmp/p1-m4.2-review/red.txt`.
- A separate attempt to capture RED output via Python subprocess hit denied
  Go-cache access (exit 1); its assertion failed because missing symbols were not
  in that environmental output. That failure is not RED evidence. The log is
  preserved as red-capture-environmental-failure.txt. The same Go command was
  rerun permission-enabled after implementation, exit 0.
- GREEN same focused Go command exit 0, seven new and four existing test
  functions pass. `gofmt -w` on both new Go files exited 0.

Post-implementation commands, each exit 0:

- `python3 scripts/validate_features.py feature_list.json`
- `PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v`
  (60 tests)
- `go test -count=1 ./...` (all three packages)
- `go test -race -count=1 ./...` (all three packages)
- `go vet ./...` (no diagnostics)
- `bash scripts/verify.sh` (all six checks PASS)
- `git diff --check` and `git diff --cached --check`
- `go list -deps ./internal/httpapi ./cmd/p1-api` (engine package absent)

Preservation audit and final post-checkpoint checks are retained with complete
source, command logs, review report and full tracked-plus-untracked patch under
`/private/tmp/p1-m4.2-review/`. No source/dependency changes beyond the DTO files;
router, server, engine, harness, go.mod/go.sum, README and roadmap are unchanged.
M3.1–M3.10 and M4.1 accepted objects/evidence remain unchanged, M4.3–M4.10 remain
unchanged/pending, original scope/dependencies unchanged, schema 1 and WIP 1.
No business operations or operational execution endpoints are exposed. Structural
import/route checks support absence of invocation; no end-to-end integration is
claimed. No HTTP authentication, authorization, tenant ownership enforcement,
provider authenticity/replay defense, durable processing, scheduling or distributed
exactly-once behavior was added. EventID is correlation only; in-memory Save()
has no insertion/duplicate result. Payment string checks do not authenticate origin.

All six files remain unstaged/uncommitted. No stage, commit, push, deployment,
M4.2 acceptance or M4.3 activation occurred. No future commit SHA is asserted.
Next action: STOP for human source review of the complete patch. At next startup,
reconcile HEAD/tree/index, tracker and source/tests with this checkpoint, and run
init.sh. M4.2 is already active: do not use pending-task preflight to reactivate
it, and do not begin later work without separate authorization.

Final checkpoint checks: preservation/formatting audit exit 0. Its first attempt
incorrectly required exit 0 from `git diff --no-index --check` on added files;
Git returned 1 with no diagnostics for added content. The external audit was
corrected to accept 0/1 only with empty diagnostics, alongside explicit whitespace
and gofmt checks; no repository fix was needed. Final post-checkpoint verifier
first exited 1 at root Go unit setup from denied cache access after 60 Python
tests passed. Permission-enabled `bash scripts/verify.sh` retry exited 0, all six
checks PASS. Logs are final-canonical.txt and
final-canonical-permission-enabled.txt in the review directory. These environment
and audit-tool issues are not behavioral RED. Only evidence text changed afterward;
final tracker validation, whitespace, preservation and patch checks accompany the
review report. Observed HEAD remains eac1f8bf7a14f6853036d91071c6b5881cd373ca.


## M4.2 source-review correction checkpoint (2026-11-13)

Human authorized CORRECTION ONLY: source review found the permanent
TestOperationalRoutesNotRegistered test would freeze later legitimate API work.
Removed that test only. All other test code, DTOs/JSON contracts, sensitive-field
allowlists, health tests, docs/API_CONTRACTS.md, router/server, engine and
dependencies remain unchanged from the reviewed implementation.

Previous committed checkpoint revision and current observed HEAD on main:
`eac1f8bf7a14f6853036d91071c6b5881cd373ca`. The six existing modified/new
M4.2 files matched the previous review snapshot before correction. init.sh exit 0;
its existing-changes notice was reviewed. Index empty. No fresh remote lookup
was needed or performed for this local correction; prior live observation is
historical. M4.2 was already active; no pending-task preflight or activation.

Current focused `go test -count=1 -v ./internal/httpapi/...` exit 0: observed
six M4.2 test functions plus four existing router tests, all PASS. Earlier seven
M4.2 test counts above describe the pre-correction run only.

One-time read-only acceptance audit, outside the permanent suite:
`python3 /private/tmp/p1-m4.2-correction/route_acceptance_audit.py`.
It supplies external audit code through Go's test-only overlay and calls actual
NewRouter(), without replacing production source or writing an audit into the
repository. Observed only GET /health; POST /v1/workflow-executions and
GET /v1/workflow-executions/execution-1 each returned 404 "404 page not found".
The actual HTTP/server dependency graph excludes the workflow/business package;
production router/server bytes match publication and correction start. Together
with the inspected health-only handler, this establishes current nonexposure,
not future routing constraints or end-to-end integration. First sandbox audit
exited 1 at go list from denied Go-cache access; permission-enabled retry exited
0 with all assertions PASS. Reproducible code, overlay and actual output are in
`/private/tmp/p1-m4.2-correction/` (route-audit-permission-enabled.txt).

Fresh requested regressions each exited 0: tracker validation; Python unittest
discovery with PYTHONDONTWRITEBYTECODE=1 (60 tests); go test -count=1 ./...;
go test -race -count=1 ./...; go vet ./...; bash scripts/verify.sh (all six checks);
git diff --check; git diff --cached --check. All three Go packages pass.
Full correction preservation audit and final evidence/whitespace checks are in
the correction review materials. Only this test removal and factual updates to
PROGRESS.md, docs/HANDOFF.md and M4.2 tracker evidence changed during correction.
M3/M4.1 accepted objects/evidence unchanged; later tasks unchanged/pending;
M4.2 remains in_progress, active_task_id M4.2, reviewed_by null, schema/WIP 1.
Human source review of the correction remains pending; no verification status.
Exact correction-only and updated complete patches are in the same external
review directory. No staging, commit, push, deployment or M4.3 activation.
Next action: STOP for human review. Uncommitted work is not a published revision.

Final correction preservation audit exited 0, comparing all repository files to
the pre-correction snapshot and all M3/M4.1/later task objects to publication.
Final post-checkpoint canonical verifier first exited 1 at root-package setup
from Go-cache denial; permission-enabled retry exited 0, all six checks PASS
(60 Python tests, all three Go packages unit/race, vet and whitespace). Only this
factual evidence text changed afterward. Both logs, reproducible audits and
exact correction patch are retained in /private/tmp/p1-m4.2-correction/.


## M4.2 acceptance reconciliation checkpoint (2026-11-13)

Published starting baseline and current observed HEAD on main:
`eac1f8bf7a14f6853036d91071c6b5881cd373ca`.
Previous committed documentation checkpoint revision: the same `eac1f8b`;
M4.2 implementation/correction checkpoints above remain uncommitted local work.
Fresh live remote main matched via permission-enabled git ls-remote (exit 0);
initial sandbox lookup exited 128 from DNS restriction. init.sh and tracker
validation exited 0; expected existing-changes notice reviewed. All six local
M4.2 files matched the accepted correction snapshot exactly; index empty.
No pending-task preflight or later activation was run.

Human explicitly authorized acceptance reconciliation following source review
of the original implementation and correction-only patches against all six
criteria. Codex independently inspected source and ran local checks. ChatGPT
is not claimed to have executed local tests; the human is not claimed to have
personally run them. This acceptance supersedes historical pending-review and
in_progress statements without rewriting their evidence.

Acceptance mapping: API_CONTRACTS.md documents health and two future execution
paths; DTOs and six permanent tests cover exact JSON, public allowlists, uint64
version and four persisted states; nested errors and the explicit status matrix
are documented/tested; identity, authorization, tenancy and payment-origin trust
are requirements/blockers, not claimed controls. The historical one-time route
audit observed health only, proposed paths default 404 and no engine imports.
Its code/output remain available under /private/tmp/p1-m4.2-correction/ and are
copied into /private/tmp/p1-m4.2-acceptance-review/historical-route-audit/.
Identical current router/server source ties that evidence to this acceptance;
it is not a permanent restriction on later authorized routing work.

Fresh pre-status checks each exited 0 on successful runs: focused Go unit/race;
full Go unit/race (three packages); go vet; Python unittest discovery with
PYTHONDONTWRITEBYTECODE=1 (60 tests); tracker validation; canonical verify.sh
(all six checks); working-tree/index whitespace. Focused verbose output observed
six M4.2 plus four existing router test functions PASS. gofmt -l was clean and
untracked whitespace checked independently, with no formatting/source changes.
Initial sandbox full Go unit/race/vet and canonical runs exited 1 from denied
root Go-cache access; permission-enabled reruns passed. No behavioral failure
or implementation change. Exact commands/exits are retained in review logs.

Only after passing checks, M4.2 transitioned in_progress -> verified and
active_task_id became null. Historical evidence preserved; fresh results and
truthful human review attribution appended. All M3/M4.1 objects and evidence
unchanged; M4.3–M4.10 unchanged/pending, no active task, schema/WIP 1. Original
requirements/dependencies intact. HTTP identity/authorization/tenant ownership,
provider authenticity/replay defense, durable processing/scheduling and distributed
exactly-once guarantees remain unresolved. No execution/business endpoint added.

This acceptance edits only feature_list.json, PROGRESS.md and docs/HANDOFF.md.
Complete six-file M4.2 working tree additionally includes the unchanged accepted
internal/httpapi/contracts.go, internal/httpapi/contracts_test.go and
docs/API_CONTRACTS.md. Router/server, engine, dependencies and existing tests
remain unchanged from publication. No implementation or acceptance commit exists;
no staging, commit, push, deployment or M4.3 work occurred. Final post-update
checks, preservation audit, acceptance-only patch and complete tracked/untracked
patch are retained in /private/tmp/p1-m4.2-acceptance-review/.
Next action: STOP for human publication review. At next startup reconcile actual
HEAD/tree/index/tracker and source/tests; do not treat local acceptance as already
published or infer permission to activate M4.3.

Final acceptance results: post-reconciliation tracker/Python (60 tests), full Go
unit/race, vet and canonical verifier all exited 0 on successful runs. Initial
post-update sandbox race/vet/canonical attempts exited 1 from denied root cache
access; full unit passed even in that attempt. Permission-enabled retries all
passed. Independent preservation audit exited 0: exact three-file acceptance-only
scope, historical evidence intact, accepted corrected source and all other task
objects unchanged, no active task. No untracked whitespace or gofmt diagnostics.
Only factual evidence/checkpoint text changed after those regressions; final
validation, whitespace, preservation and patch applicability checks are recorded
in /private/tmp/p1-m4.2-acceptance-review/. Acceptance remains local/uncommitted;
STOP for human publication review, with no authorization for M4.3 or publication.


## M4.3 request-validation implementation checkpoint (2026-11-14)

Previous committed documentation checkpoint revision:
`89c1317b2785df994036bb6ea7ca087599ca828b` (published M4.2).
Current observed starting HEAD on main: the same full SHA. Starting tree clean,
index empty; live remote main independently matched via permission-enabled
`git ls-remote origin refs/heads/main` exit 0 after sandbox DNS failure exit 128.
Git history and the accepted tracker supersede historical M4.2 local/uncommitted
and pending-publication statements above. README/roadmap status notes lag source.
Initialization `bash init.sh`, tracker validation and M4.3 preflight exited 0;
initialization reported clean tree and no warnings beyond cached-remote limits.
Baseline canonical verifier initially exited 1 at Go-cache trim permissions;
permission-enabled retry exited 0, all six checks and 60 Python tests passed.
That environment failure is not behavioral RED evidence.

The human's bounded M4.3 request explicitly authorizes activation/implementation.
AGENTS.md requires human authority, not a second approval after this explicit
request. Preflight's ELIGIBLE / HUMAN APPROVAL REQUIRED result established only
structural prerequisites. Activation changed only active_task_id to M4.3 and its
status to in_progress; empty evidence was retained until actual results existed.
The objective, all six original criteria and dependency remain unchanged.
M4.3 remains in_progress, reviewed_by/revision null; human source review is pending.

Exactly six repository files changed/new:

- internal/httpapi/validation.go: ValidationFailure and the recommended
  ParseAndValidateCreateExecutionRequest API; standard-library-only, no engine
  imports or downstream dependencies. Rejects >4096 raw bytes, non-single-object
  JSON, unknown/duplicate keys, missing/nonstring/invalid identifiers and
  unsupported definitions. Identifiers use 1–128 ASCII bytes from letters,
  digits, underscore, hyphen and period, with no trimming. Trusted server support
  map requires exact case-sensitive membership; no wildcard/default or authority
  grant. All rejections return a zero DTO and safe deterministic 400
  INVALID_REQUEST using unchanged M4.2 envelopes. payment_status is rejected.
- internal/httpapi/validation_test.go: 74 rejection cases, each checked twice for
  stable error/zero-DTO behavior and exact error-envelope JSON; 74 negative seam
  cases asserting zero downstream/repository creation/business calls. Each seam
  case has a valid positive control. Successful cases include 1/128-byte IDs,
  every allowed character class, reordered/escaped fields, ordinary whitespace
  and exactly 4096-byte valid bodies; 129-byte IDs and >4096-byte bodies fail.
  A fake repository sentinel observes control flow; no real repository or
  business operation is constructed/invoked. This establishes the local validation
  boundary, not future handler integration or external end-to-end correctness.
- docs/API_CONTRACTS.md: appended approved M4.3 validation policy, preserving
  historical M4.2 text and its deferred trust/security requirements.
- feature_list.json: sole-task activation and actual evidence, no acceptance.
- PROGRESS.md and docs/HANDOFF.md: this factual checkpoint.

RED: before production validation.go existed,
`go test -count=1 -v ./internal/httpapi/...` exited 1 with undefined
ValidationFailure and ParseAndValidateCreateExecutionRequest. Actual output is
red.txt in the external review directory. GREEN: same command exited 0 after
implementation; both new Go files were formatted with gofmt (exit 0). The focused
suite has 13 passing test functions: ten existing and three new table-driven
functions. Detailed names/outcomes are in green.txt and focused.txt.

One-time audit command:
`python3 /private/tmp/p1-m4.3-review/route_audit.py`.
External Go test-only overlay calls actual NewRouter() without changing repository
files: only GET /health registered, health retains its 200 JSON/content type;
POST /v1/workflow-executions and GET /v1/workflow-executions/execution-1 retain
404 `404 page not found`. Actual HTTP/server dependency graph excludes the engine;
router/server production bytes match HEAD. The inspected sole health handler
exposes no enrollment or payment operation. Initial audit exited 1 at go list
cache denial after route checks passed; permission-enabled full retry exited 0.
Audit source, overlay, both logs and exit records are outside the repository.
No permanent route-freezing test was added.

Fresh post-implementation commands all exited 0:

- python3 scripts/validate_features.py feature_list.json
- PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
  (60 tests)
- go test -count=1 -v ./internal/httpapi/... (13 test functions)
- go test -count=1 ./... (all three packages)
- go test -race -count=1 ./... (all three packages)
- go vet ./... (no diagnostics)
- bash scripts/verify.sh (all six checks PASS)
- git diff --check and git diff --cached --check

Full logs, exact command/exit manifest, reproducible preservation audit, formatting
and untracked-whitespace checks, complete tracked-plus-untracked patch and review
report are retained in `/private/tmp/p1-m4.3-review/`. Final post-document checks
are recorded there separately; do not infer acceptance from passing checks.
All 20 task records retained, schema/WIP 1, M3/M4.1/M4.2 verified objects/evidence
unchanged; M4.4–M4.10 unchanged/pending. No original criterion/dependency edited.
DTOs, existing tests, router/server, engine/enrollment, dependencies and harness
remain byte-identical to publication. No cloud/provider/auth/cache/queue system,
application service, execution endpoint, invocation or streaming HTTP middleware
was introduced. Approval-map membership is support selection, not authorization.

Current work is unstaged and uncommitted; HEAD remains 89c1317. No stage, commit,
push, dependency installation or deployment performed. Next action: STOP for
human source review. Do not mark M4.3 verified, publish, activate M4.4 or implement
later work. At next startup reconcile observed HEAD/tree/index/tracker and fresh
evidence with this checkpoint; init.sh is required. Review files are temporary
external artifacts, not committed repository content.

Final post-document verification: sandbox canonical verifier exited 1 during
root-package unit setup from denied Go-cache access; permission-enabled retry
exited 0, all six checks PASS (60 Python tests, full Go unit/race, vet and both
whitespace checks). Logs: final-canonical.txt and
final-canonical-permission-enabled.txt in the review directory. No behavioral
failure or source fix. Final preservation audit exited 0; formatting and untracked
whitespace clean. Only factual evidence text changed after these regressions;
final tracker/whitespace/preservation and reverse patch-applicability checks are
recorded in the external review report. No acceptance or publication occurred.


## M4.3 acceptance reconciliation checkpoint (2026-11-14)

Previous committed documentation checkpoint revision:
`89c1317b2785df994036bb6ea7ca087599ca828b` (published M4.2).
Current observed starting HEAD on main is the same full SHA; live remote main
independently matched via permission-enabled git ls-remote (exit 0) after sandbox
DNS failure (exit 128). Starting tree had exactly the expected six M4.3 files,
index empty. Complete tracked-plus-untracked patch matched the human-reviewed
`/private/tmp/p1-m4.3-review/Complete-review.patch` byte-for-byte. Initialization
and tracker validation exited 0; expected existing-changes notice reviewed.
No pending-task preflight rerun: M4.3 was already the sole active task.
Historical pending-review statements above describe implementation-time status;
this explicitly authorized acceptance checkpoint supersedes them.

The human engineering owner reports source-level review by ChatGPT accepted the
six-file M4.3 implementation against all six original criteria with no correction
required, and explicitly authorized acceptance reconciliation only. Codex separately
inspected source/tests and ran the fresh local verification reported below.
ChatGPT is not claimed to have executed local Go/Python commands; the human is not
claimed to have personally run tests. Source-level human acceptance, Codex command
evidence, local tracker reconciliation and future publication are distinct.

Acceptance mapping:

1. ParseAndValidateCreateExecutionRequest / validRequestIdentifier and
   TestParseAndValidateRejectedRequests reject malformed/non-object/trailing JSON,
   missing/nonstring/invalid identifiers, duplicate/unknown fields and oversized
   bodies. TestValidateApprovedDefinition covers 1/128-byte identifier boundaries,
   allowed characters, escapes/order, whitespace and exactly 4096-byte valid input.
2. Unsupported, case-mismatched, nil/empty support sets and wildcard-only maps
   reject deterministically. Exact server support membership is not authorization;
   no dynamic function selection or privileged fallback exists.
3. invalidRequest returns a zero request and 400 INVALID_REQUEST using unchanged
   ErrorResponse/APIError; assertValidationFailure checks exact nested JSON and
   safe fixed messages for every rejection. No operational HTTP error handler.
4. TestRejectedRequestHasNoSideEffects calls the real validator through a test-only
   boundary, returns on failure and asserts zero downstream/simulated repository
   creation/business calls for 74 invalid cases; each has a valid positive control.
   The sentinel is not a real repository. Production validation has no downstream
   dependency/engine import. This proves local validation control flow, not a
   deployed or end-to-end HTTP workflow system.
5. payment_status including CONFIRMED is an unknown field, rejected before the
   callback. Existing DTOs and enrollment rules are unchanged; no payment authority,
   provider authentication or entitlement operation added.
6. Fresh prior regressions passed on successful runs as recorded below.

Fresh BEFORE-status commands (each successful run exit 0):

- go test -count=1 -v ./internal/httpapi/...: 13 test functions (ten existing,
  three M4.3), 74 rejection scenarios each checked twice and 74 negative callback
  cases each with a positive control; all PASS.
- go test -race -count=1 ./internal/httpapi/...
- go test -count=1 ./... and go test -race -count=1 ./... (three packages)
- go vet ./... (no diagnostics)
- PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
  (60 tests)
- python3 scripts/validate_features.py feature_list.json
- bash scripts/verify.sh (all six checks)
- git diff --check and git diff --cached --check
- gofmt -l internal/httpapi/validation.go internal/httpapi/validation_test.go
  (empty output, no file changes); separate untracked-whitespace audit PASS.

Initial sandbox full Go unit/race/vet/canonical commands each exited 1 from denied
root-package Go-cache file access, not behavioral failures. Each permission-enabled
retry exited 0. Focused unit/race, Python, tracker, formatting/whitespace checks
passed initially. Exact commands, output and exits are retained separately under
`/private/tmp/p1-m4.3-acceptance-review/` (pre-*.txt and result manifests).

Historical one-time route audit method and output were inspected and copied to
`historical-route-audit/` in that directory. The overlay calls actual NewRouter();
observed only GET /health, unchanged 200 JSON/content type, default 404 for both
proposed execution paths, and no engine dependency in HTTP/server import graph.
Current router/server bytes remain identical to audited publication, and inspection
confirms the sole health handler exposes no business operation. This historical
audit was not freshly executed in this acceptance session. It establishes current
nonexposure without permanently freezing future approved routes; no permanent
route-nonexposure test was added. Genuine historical RED/GREEN logs preserved.

Only after all required pre-status checks passed, M4.3 changed in_progress ->
verified and active_task_id became null. Its original objective, dependency and
six criteria remain unchanged; earlier result text and command evidence retained,
new review attribution and acceptance results appended. Revision text identifies
published baseline plus accepted uncommitted source and local reconciliation,
not a fabricated M4.3 commit. Schema/WIP remain 1, all 20 task records retained.
M3.1–M3.10/M4.1/M4.2 complete verified objects/evidence unchanged; M4.4–M4.10
complete pending objects unchanged. No task activated.

Exactly three acceptance-only files change: feature_list.json, PROGRESS.md and
docs/HANDOFF.md. The complete working tree still contains six M4.3 files; both
validation source/test files and docs/API_CONTRACTS.md are byte-identical to the
accepted snapshot. Existing DTOs, tests, routes, server, engine/enrollment,
dependencies, harness and infrastructure are unchanged. Final AFTER-status
regressions and independent preservation audit are recorded in the external
verification report; do not substitute historical results for those final checks.

Authentication, authorization, tenant ownership, payment authenticity/replay
protection, durable storage/scheduling, production provider integration and
end-to-end guarantees remain unresolved. No invocation, operational endpoint,
application-service interface, security/provider/cache/queue/rate-limit system,
dependency installation or M4.4 work was introduced.

Acceptance is local and uncommitted. HEAD remains 89c1317; index empty. No stage,
commit, push or deployment performed; GitHub has not been updated by this session.
Next action: STOP for human publication review. Future publication and M4.4
activation remain unauthorized. At next startup run init.sh and reconcile actual
HEAD/tree/index/tracker and current evidence; do not equate verified local status
with publication. External review artifacts include acceptance-only and complete
six-file patches, six-criterion mapping, command logs and preservation audit.

Final acceptance verification: AFTER-status tracker, Python (60 tests), focused
Go (13 functions; 74 rejection and 74 negative seam cases), full Go unit/race,
vet, canonical verifier, formatting and working-tree/index whitespace checks
passed on their successful runs, all exit 0. Initial post-status sandbox vet and
canonical runs exited 1 at root Go-cache permissions; permission-enabled retries
both exited 0. Full unit/race passed on the initial post-status attempt. Exact
pre/post commands, failures and retries are preserved in external manifests/logs.
Independent preservation audit exited 0: exactly three acceptance-only changes,
accepted source/policy unchanged, histories/evidence preserved, no active task,
empty index. Only this factual evidence text and tracker evidence changed after
regressions; final tracker/whitespace/preservation and patch-applicability checks
are recorded in the verification report. Acceptance remains local/uncommitted;
STOP for human publication review, with no publication or M4.4 authorization.
