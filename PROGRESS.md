# P1 - Development Progress

## Source of Truth

- Repository: CEAG0712/p1-reliable-academy-workflows
- Branch: main
- Source files and tests determine implemented behavior.
- This document tracks work status, not proof of correctness by itself.

## Historical Repository Baseline

- Historical user-run baseline commit: 2538574 (test date not recorded)
- Go version: 1.24.5
- Go unit tests: PASS (user-run baseline)
- Go race checks: PASS (user-run baseline)
- go vet: No diagnostics reported in user output
- Results apply to the confirmed baseline, not future changes.

Historical Lesson 10 starting checkpoint (2026-11-08): main at
b0e4db4ec4d20a7e3d150de0342e9f80821669cb, the confirmed Lesson 9
GitHub commit. Starting tree was clean; HEAD matched cached origin/main.
No fetch performed; live synchronization remains unverified. Source/tests
support implemented M1, M2, and Exercise 82. No Go source/test differences
exist between 2538574 and this HEAD; fresh checks are recorded below.

Historical Lesson 11 starting checkpoint (2026-11-08): clean main at
01ef1e37f288459276d3025335b7cdcde669f733, matching the expected Lesson 10
commit and cached origin/main. Local Git history confirms the Lesson 10
documentation commit. No fetch or live remote verification was performed.

## Implemented in the Confirmed Repository

- Enrollment eligibility validation
- In-memory entitlement provider
- Workflow execution state and transition history
- Single-winner workflow claims
- Owner-aware workflow completion and failure
- Version-checked repository updates
- Failure reasons preserved by Fail()
- In-memory repository with independent snapshots

## Harness Status

- AGENTS.md: implemented and committed
- AGENTS.md context-routing rules: implemented
- PROGRESS.md: implemented and committed
- init.sh: implemented; local execution verified
- feature_list.json: implemented; initial JSON validation passed
- scripts/validate_features.py: implemented; G1-G7 tracker validation PASS
  with `python3 scripts/validate_features.py feature_list.json`
- scripts/test_validate_features.py: 27 unittest tests PASS with
  `python3 -m unittest discover -s scripts -p "test_*.py" -v`
  (bytecode writes disabled); original feature_list.json remains unchanged
- scripts/verify.sh: implemented; resolves the root from its own location,
  checks tools, and runs tracker validation, Python harness tests, Go unit
  tests, race detection, go vet, and both Git whitespace checks in order.
  Stops on the first failure, identifies the check, and preserves its exit code.
- scripts/test_verify.sh: 12 isolated shell tests PASS; root/nested invocation
  (including spaces), ordered commands/headings, final PASS, seven command
  failures preserving exit codes 41-47 with no later commands or PASS,
  and missing python3/go/git exiting 127 before any checks.
- Lesson 9 verification (2026-11-08): clean starting tree on main at
  0a6fbc7794fee2cd4d2137631e8805796a3a9336; local HEAD matched cached
  origin/main. Remote was not fetched. No Go source changes since the
  documented 2538574 baseline; no local M3 engine implementation found.
- Commands: `bash init.sh`, `python3 scripts/validate_features.py feature_list.json`,
  `bash -n scripts/verify.sh`, `bash -n scripts/test_verify.sh`, and
  `bash scripts/test_verify.sh` all exited 0.
- Actual `bash scripts/verify.sh`: first exited 1 at Go unit tests because
  the sandbox denied Go build-cache access; no final PASS was printed.
  Retried with permission for cache access and exited 0: tracker PASS,
  27 Python tests PASS, Go unit tests PASS, Go race tests PASS, go vet
  without diagnostics, and both Git whitespace checks PASS.
- Lesson 9 changes to scripts/verify.sh, scripts/test_verify.sh, AGENTS.md,
  and PROGRESS.md are recorded in confirmed GitHub commit
  b0e4db4ec4d20a7e3d150de0342e9f80821669cb (harness - Unified Verification).
  The prior pending-commit/push statement is superseded. This session
  inspected the local commit; it did not fetch live remote state.
- AGENTS.md now requires unified verification before engineering acceptance
  and retains individual commands for targeted debugging.
- Historical Lesson 9 evidence before its commit: `git diff --check` and
  `git diff --cached --check` exited 0. `git status --short` exited 0:
  AGENTS.md and PROGRESS.md modified; both new shell scripts untracked.
- Codex local execution: verified at commit a20b8d1
- Fresh-agent context reconstruction: passed
- docs/ROADMAP.md: implemented and committed
- docs/INTERVIEW_DEFENSE.md: implemented and committed
- Interview-defense protocol: defined; practice during M3 pending
- docs/HANDOFF.md: created in Lesson 10 and committed with AGENTS.md and
  PROGRESS.md as 01ef1e37f288459276d3025335b7cdcde669f733, confirmed in
  local Git history. The user reports the Lesson 10 fresh-agent acceptance
  test succeeded; that historical report is not a newly executed check here.
- Fresh Lesson 10 verification (2026-11-08), starting at b0e4db4:
  `bash init.sh` and `python3 scripts/validate_features.py feature_list.json`
  exited 0. Initial `bash scripts/verify.sh` exited 1 at Go unit tests
  because sandbox access to the Go build cache was denied; tracker and
  27 Python tests passed, with no final PASS. Permission-enabled rerun
  exited 0: tracker, 27 Python tests, Go unit tests, race checks, go vet,
  and working-tree/index whitespace checks passed. This evidence is
  separate from the historical Lesson 9 checks above.
- Historical M3.5 checkpoint: implementation committed; final tracker acceptance pending
  then, with M3.1-M3.4 verified and M3.6 not yet authorized. Superseded below.
- Current M3 active task: M3.10, test implementation uncommitted; human review pending.
  M3.1-M3.9 are verified; Stage 3 final acceptance remains pending.
  See the M3.10 checkpoint below for fresh evidence and boundaries.

## Lesson 11 — Controlled Task Preflight (2026-11-08)

- Historical pre-commit harness-only changes: scripts/task_preflight.py,
  scripts/test_task_preflight.py, AGENTS.md, PROGRESS.md, docs/HANDOFF.md.
  Lesson 11 starting HEAD was 01ef1e37f288459276d3025335b7cdcde669f733.
  These five files were subsequently committed in a5b345d99075debfe652d0182f0374466b7ff8bc.
- Standard-library CLI reuses validate_tracker contract checks; rejects
  unreadable/malformed/invalid trackers, unknown IDs, non-pending tasks,
  unverified dependencies, and another in_progress task. Default tracker
  resolves from the script location, independently of the caller's directory.
- Success explicitly requires human approval, does not authorize implementation,
  and activates nothing. This is a snapshot eligibility check, not an enforceable
  human-authorization system; it neither authenticates approval nor reserves
  eligibility against later tracker changes. Chat/handoff/environment/flags
  are not treated as proof of approval.
- Actually executed: `bash init.sh` and
  `python3 scripts/validate_features.py feature_list.json` exited 0.
  Initialization warned that cached origin/main is not live remote evidence.
- `PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p "test_*.py" -v`
  exited 0: 44 tests passed (27 existing and 17 new preflight tests).
- `python3 scripts/task_preflight.py --task M3.1` exited 0: ELIGIBLE,
  HUMAN APPROVAL REQUIRED, no task activated, no implementation authority.
- `python3 scripts/task_preflight.py --task M3.2` exited 1: depends on
  unverified M3.1. This is the expected successful negative test, run separately.
- `bash scripts/verify.sh` initially exited 1 at go vet because sandbox
  access to a Go build-cache file was denied; tracker, 44 Python tests,
  Go unit and race tests had passed. Permission-enabled rerun exited 0:
  tracker, 44 Python tests, `go test -count=1 ./...`,
  `go test -race -count=1 ./...`, `go vet ./...`, and both Git whitespace
  checks passed. No existing initialization or verification scripts changed.
- `git diff --check`, `git diff --cached --check`, and `git status --short`
  exited 0. Status: three documentation files modified; two Python files
  untracked. No commit or push was made.
- Read-only Python assertions exited 0: feature_list.json equals HEAD bytes
  and its starting SHA-256
  `5fb4828b8fbec11444e4dbfe396e8cbcd882ebaad5e6c3de2554dbef084bc2c3`;
  all ten tasks pending, active_task_id null, WIP 1; exactly five authorized
  files changed. No Go application/test files changed; no M3 task activated
  or implemented. M3.1 remains pending and unauthorized.
- Historical next action was human review of Lesson 11; its commit is now observed in local history.

## Current Engineering Milestone

M3 - Reusable Workflow Definition and Execution Engine

Milestone capabilities (implementation status below):

- Reusable workflow definitions
- Triggers
- Conditions
- Branching
- Action handlers
- Wait/delay steps

Historical status at the M3.5 documentation checkpoint (superseded by M3.6 below):

- M3.1 definition, trigger and step models are verified, committed as
  93c9861b93091f270d4d9a78a0edc8e19376fd65.
- M3.2 Definition Validation is verified and human-accepted at 10/10, committed as
  4fb635168f61125842be7cf88241c25738c6f265.
- M3.3 Trigger Matching is verified and human-accepted at 10/10, committed as
  1a5ad1413c80b7024cb82291bdb7843ff23e857a.
- M3.4 Handler Registration is verified; implementation commit
  1e34b5d0108949be651e9e3b96ac4cef4d3c5100, acceptance commit
  88fe6ed924efd351cb2a10d5cd8f11f326b81006.
- M3.5 Conditions and Branching implementation is completed and committed locally
  as e0e73da376d55cd22c86fdda8d6011af022a9b0a, observed HEAD on main (2026-11-11).
  M3.5 remains in_progress pending final tracker acceptance; active_task_id is M3.5,
  WIP limit is 1. M3.6-M3.10 remain pending. No tracker data changes here.
- Models, validation, event eligibility, named callback dispatch and single-condition
  branch selection are implemented. Runner and lifecycle integration remain planned.
- M3.6 is next planned and NOT authorized. Only PROGRESS.md and docs/HANDOFF.md
  are being reconciled; these edits remain uncommitted.
- Startup tree was clean. Initialization warned HEAD differs from cached origin/main;
  no fetch or live remote check occurred. M3.5 GitHub publication is unverified.

## Historical M3.1 — Workflow Definition Models (2026-11-08)

The following evidence describes the original implementation session. Its seven
then-uncommitted files were subsequently committed as
93c9861b93091f270d4d9a78a0edc8e19376fd65; historical revisions below are retained.

- User explicitly authorized activation and implementation of the bounded proposal,
  plus tracker/checkpoint updates; no commit, push, deployment or M3.2 authorization.
- Previous checkpoint: Lesson 11 starting revision 01ef1e37f288459276d3025335b7cdcde669f733.
  Observed HEAD during M3.1: a5b345d99075debfe652d0182f0374466b7ff8bc. Started clean on main;
  the Lesson 11 preflight files are committed at HEAD, not still uncommitted.
- Acceptance criteria: models compile in Go; tests demonstrate two distinct workflow
  definitions without business-specific engine fields.
- Created workflow_definition.go and workflow_definition_test.go. Definitions describe
  IDs, an event-type trigger, a starting step, and steps keyed by ID. Action/condition
  names, next/true/false links, and time.Duration delays are configuration only.
- Tests describe enrollment branching and customer follow-up with an action/wait/action
  graph using the same types. No handlers, matching, validation, execution or resume logic.
- Updated scripts/test_task_preflight.py and scripts/test_validate_features.py because
  their fixtures assumed the published tracker remained entirely pending. Behavioral
  fixtures now establish their own states; default-path CLI checks compare against the
  explicit current tracker. Production validator/preflight scripts are unchanged.
- Startup: bash init.sh, python3 -B scripts/validate_features.py feature_list.json,
  and python3 -B scripts/task_preflight.py --task M3.1 all exited 0. Initialization
  reported cached origin/main only; preflight required human approval. The user's
  explicit authorization, not eligibility, permitted activation. Activation validation passed.
- RED: go test -count=1 -run '^TestWorkflowDefinition' ./... exited 1 with undefined
  WorkflowDefinition, TriggerDefinition, StepDefinition and step constants.
- GREEN: the same focused command exited 0 after adding models and running gofmt.
- bash scripts/verify.sh initially exited 1 at Python tests: four validator fixture
  failures after the preflight fixtures were corrected. After correcting validator
  fixtures, rerun passed all 44 Python tests but exited 1 at Go unit tests because
  sandbox access to a Go build-cache file was denied. Permission-enabled rerun
  exited 0: tracker, 44 Python tests, go test -count=1 ./..., go test -race -count=1
  ./..., go vet ./..., git diff --check, and git diff --cached --check passed.
- Codex reviewed the models and test evidence against both acceptance criteria;
  M3.1 is marked verified with the uncommitted revision explicitly recorded. This
  records agent review, not human acceptance or completion of the whole M3 milestone.
- Seven changed/untracked files: workflow_definition.go, workflow_definition_test.go,
  scripts/test_task_preflight.py, scripts/test_validate_features.py, feature_list.json,
  PROGRESS.md, docs/HANDOFF.md. Existing lifecycle/business Go files are unchanged.

Historical M3.1 final verification after tracker/checkpoint updates: sandbox run exited 1 at
go vet after tracker, all 44 Python tests, Go unit and race tests passed.
Permission-enabled final rerun of `bash scripts/verify.sh` exited 0 with
all six checks PASS. Final scope assertions confirmed M3.2-M3.10 unchanged,
WIP 1, active_task_id null and existing lifecycle/business Go files unchanged.
Only this evidence text was added afterward; tracker validation and both
Git whitespace checks were rerun successfully.

## M3.2 — Accepted Implementation Checkpoint (2026-11-08)

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

- `WorkflowDefinition.Validate() error` checks required fields, supported types,
  handler names, successor references, type-specific field compatibility, positive
  WAIT durations and cycles across all configured components. Sorted step IDs and
  fixed check order make errors deterministic. Empty successors terminate paths.
- The implementation commit contains exactly feature_list.json,
  workflow_definition.go (comments only), workflow_definition_validation.go and
  workflow_definition_validation_test.go. Models were not redesigned; M1/M2 and
  Exercise 82 source and regression tests are unchanged and remain passing.

## Historical M3.2 Post-Commit Documentation Reconciliation (2026-11-08)

- `git rev-parse HEAD` returned 4fb635168f61125842be7cf88241c25738c6f265;
  `git status --short` was empty; `git show --stat --oneline HEAD` and
  `git diff-tree --no-commit-id --name-only -r HEAD` confirmed exactly the four
  approved implementation files. All commands exited 0.
- `bash init.sh` exited 0. It warned HEAD differs from cached origin/main and
  reported remote state was not fetched or verified. This does not prove either
  a successful push or that no push occurred.
- Previous documentation checkpoint described a5b345d plus M3.1 working changes;
  M3.1 is now committed at 93c9861 and M3.2 at 4fb6351. Earlier test evidence and
  its original revisions above remain historical records.

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

## M3.3 — Accepted Implementation Checkpoint (2026-11-09)

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

## Historical M3.3 Post-Commit Documentation Reconciliation (2026-11-09)

- Previous documentation checkpoint revision: 74ab066fc22554384a9337588513f1298a7fe0e9,
  which reconciled M3.2 at 4fb6351; its claims about uncommitted documentation are
  historical. M3.3 acceptance testing also used 74ab066 plus uncommitted M3.3 files.
- Current observed HEAD: 1a5ad1413c80b7024cb82291bdb7843ff23e857a on main.
  `git rev-parse HEAD`, `git branch --show-current`, `git status --short`,
  `git show --stat --oneline HEAD` and `git diff-tree --no-commit-id --name-only -r HEAD`
  exited 0. Startup tree was clean and the three approved commit files matched.
- `bash init.sh` exited 0, warning that HEAD differs from cached origin/main;
  no remote was fetched or verified. This does not prove whether a push occurred.
- Startup tracker validation exited 0. Only PROGRESS.md and docs/HANDOFF.md are
  being reconciled; these documentation edits are uncommitted. No documentation
  commit SHA exists for this work. Fresh verification results are recorded below.

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

## Historical M3.4 — Committed Implementation Checkpoint (2026-11-10)

- Previous documentation checkpoint: d0784d8, which reconciled M3.3; the defense
  protocol upgrade is committed at 701c3e748c4f61093175b171d5038199f4326921.
  Earlier statements describing the M3.3 documentation edits as uncommitted are historical.
- Actual M3.4 implementation commit: 1e34b5d0108949be651e9e3b96ac4cef4d3c5100.
  Inspection confirms exactly feature_list.json, workflow_handlers.go and
  workflow_handlers_test.go; no unrelated files changed. The committed registry
  and six tests match the implementation presented in the prior session.
- Acceptance criterion 1: "Multiple implementations can be registered and invoked by name."
  TestHandlerRegistryNamedInvocation registers two actions and two predicates,
  checks true/false results and the exact callback sequence. Registration invokes nothing.
- Acceptance criterion 2: "Unknown handlers are rejected."
  TestHandlerRegistryRejectsUnknownNames covers missing and blank names,
  case/padding mismatches and wrong namespaces without invoking callbacks.
- ActionHandler is func() error; ConditionHandler is func() (bool, error).
  NewHandlerRegistry and the usable zero value support RegisterAction,
  RegisterCondition, InvokeAction and EvaluateCondition. Separate maps store
  exact, case-sensitive nonblank names without trimming them.
- TestHandlerRegistryRejectsInvalidRegistration covers blank names and nil functions.
  TestHandlerRegistryDuplicatesPreserveOriginal checks rejection without replacement,
  including the same name in separate action/condition namespaces.
  TestHandlerRegistryPropagatesErrors checks original error identity and unchanged
  predicate results with errors. TestHandlerRegistryZeroValueAndIsolation checks
  zero-value use, separate instances and distinct opaque names.
- Historical RED at 701c3e7 plus then-uncommitted tests:
  go test -count=1 -run '^TestHandlerRegistry' ./... exited 1 with undefined
  NewHandlerRegistry, ActionHandler and ConditionHandler.
- Historical GREEN: go test -count=1 -run '^TestHandlerRegistry' -v ./...
  exited 0 with all six test functions passing after implementation and gofmt.
- Historical implementation verification, preserved from the session report and tracker:
  baseline bash scripts/verify.sh passed all six checks. Post-change sandbox run
  exited 1 at Go unit-test setup from denied build-cache access after tracker and
  44 Python tests passed. Permission-enabled rerun exited 0: tracker, 44 Python
  tests, Go unit/race tests, go vet and working-tree/index whitespace checks passed.
  Existing M1/M2, Exercise 82 and M3.1-M3.3 regressions passed. These results apply
  to 701c3e7 plus the then-uncommitted M3.4 files, not a fresh run of the later commit.
- Implementation completion is distinct from task acceptance: feature_list.json
  still records M3.4 in_progress and human engineering acceptance pending.
  Its historical revision now maps to the actual commit above; no tracker edit is required here.
- Limitations: no concurrent registration guarantee; registration must not overlap
  registry operations across goroutines. Callback safety belongs to callers.
  Populated registries must not be copied. Errors/results pass through unchanged;
  callers must check predicate errors before using booleans. No panic recovery,
  retries, cancellation, idempotency, branch selection, runner or lifecycle integration.
  No cloud deployment, production integration or measured business outcome is established.

## Historical M3.4 Documentation Reconciliation (2026-11-10)

Starting HEAD is 1e34b5d0108949be651e9e3b96ac4cef4d3c5100 on main with a clean tree.
Git branch/HEAD/status/history and commit file-set inspection passed; bash init.sh
exited 0 and warned HEAD differs from cached origin/main. Live remote publication
was not checked. Startup tracker validation exited 0. Fresh verification is
recorded below separately from historical implementation results.

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

## Historical M3.5 — Local Implementation and Documentation Checkpoint (2026-11-11)

- Previous documentation checkpoint: dd35ce009d4bfc36a183671c6fbef7f48e4ed3a8,
  reconciling M3.4 at 1e34b5d. Its pending-acceptance statements are historical:
  M3.4 was subsequently accepted in 88fe6ed924efd351cb2a10d5cd8f11f326b81006.
- Starting/current observed HEAD: e0e73da376d55cd22c86fdda8d6011af022a9b0a
  on main, matching the expected implementation revision. Startup tree was clean.
  Git commit inspection confirms exactly feature_list.json, workflow_branch.go
  and workflow_branch_test.go. Committed source/tests were inspected.
- Approved criterion 1: "Both branches are tested."
  TestHandlerRegistrySelectConditionSuccessor covers true selecting OnTrue and
  false selecting OnFalse, plus four empty-successor terminal cases.
- Approved criterion 2: "Condition errors are handled and tested."
  The same test covers true/error and false/error: both return no successor and
  the original error. Each of eight table cases runs twice and checks exactly
  one predicate invocation per attempt, unchanged step and registration counts,
  and duplicate rejection preserving the original registration.
- HandlerRegistry.SelectConditionSuccessor(StepDefinition) (string, error) checks
  StepCondition before dispatching through EvaluateCondition. Unknown handler
  errors propagate with no successor; TestHandlerRegistrySelectConditionSuccessorUnknown
  covers missing, blank, case/padding mismatches and a zero-value registry.
  TestHandlerRegistrySelectConditionSuccessorRejectsWrongType rejects ACTION,
  WAIT, empty and unsupported types without invoking the predicate.
- Empty successor with nil error is a valid terminal path. An empty successor
  with a non-nil error is failed evaluation; callers must check errors first.
- Limitations: definition validation and successor-reference checks belong to
  callers. Selection returns an ID; it does not traverse a graph, execute actions,
  advance a cursor, wait/resume, claim ownership or call Complete/Fail. M1/M2
  lifecycle invariants are separate and unchanged. Each attempt reevaluates the
  predicate; no deduplication, retries, persistence or exactly-once guarantee.
  Callbacks can mutate captured state, panic or block; their side effects and
  concurrency safety remain caller responsibilities. Registration must not overlap
  registry operations across goroutines; populated registries must not be copied.
  Nil registry pointers are unsupported. No cloud or production outcome is proven.

Historical M3.5 implementation-session evidence, at
88fe6ed924efd351cb2a10d5cd8f11f326b81006 plus then-uncommitted M3.5 changes:

- Evidence recovered from the local Codex implementation-session transcript
  01a1cba2-65a3-7882-a7fd-6e562fa76160 (2026-11-11), cross-checked against
  committed source/tests. M3.5 tracker evidence fields are still empty; this
  documentation records results without filling them or asserting acceptance.
- RED: `go test -count=1 -run '^TestHandlerRegistrySelectConditionSuccessor' ./...`
  exited 1: SelectConditionSuccessor was undefined at test lines 41, 73, 78 and 93.
- GREEN: `go test -count=1 -run '^TestHandlerRegistrySelectConditionSuccessor' -v ./...`
  exited 0 after implementation/gofmt: three test functions and eight subtests passed.
- Initialization, tracker validation before/after activation and final review,
  M3.5 preflight and baseline `bash scripts/verify.sh` passed. Authorization came
  from the user's explicit instruction, not preflight eligibility.
- Historical `git ls-remote origin refs/heads/main` exited 128 from sandbox DNS
  failure; permission-enabled rerun exited 0 confirming the 88fe6ed baseline.
  This predates e0e73da and does not establish M3.5 publication.
- Post-change sandbox `bash scripts/verify.sh` exited 1 at race-test setup from
  Go build-cache access denial after tracker, 44 Python tests and Go unit tests
  passed. Permission-enabled rerun exited 0: tracker, 44 Python tests, Go unit/race
  tests, go vet and working-tree/index whitespace checks all passed.
- gofmt -l returned no files; both Git whitespace checks passed. Explicit no-index
  checks of the new files had no whitespace diagnostics (exit 1 for differing
  files). Scope assertions passed for tracker activation and the two new files.
  These are historical results, not newly reproduced RED/GREEN checks here.

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

## Historical M3.5 Blocker and Next Action

Human engineering review of these two uncommitted documentation files is next.
M3.5 implementation is committed locally but remains in_progress pending final
tracker acceptance. M3.6 — Sequential Action Execution is next planned and NOT
authorized. No commit, push, deployment or additional functionality is authorized.
README/roadmap historical status snapshots remain outside this edit scope.
The whole M3 milestone is not complete.

## Historical M3.6 — Sequential Action Execution Checkpoint (2026-11-11)

Superseded by the M3.7 checkpoint below. Original evidence and status are retained.

- Current observed HEAD: `f189ce6be261e2269b8bd266584f123c11a4ce06` on main.
  Startup tree and index were clean. The implementation commit contains exactly
  feature_list.json, workflow_actions.go and workflow_actions_test.go.
- Previous documentation checkpoint: `959186046ff3dfd3a86c350027bf7afcd2f509cc`.
  Its documentation is committed; earlier uncommitted-documentation statements
  describe historical sessions. Subsequent M3.5 acceptance baseline:
  `395239d11a17a46eaf6097b8f798585833cdda93`, also the M3.6 implementation parent.
- Tracker: M3.1-M3.5 verified; M3.6 in_progress, active_task_id M3.6, WIP 1;
  M3.7-M3.10 pending. Implementation is committed; M3.6 acceptance is pending.
  Only PROGRESS.md and docs/HANDOFF.md are edited here, still uncommitted.
- `bash init.sh` exited 0; cached origin/main remains at 395239d and differs
  from HEAD. No fetch or live remote check occurred in this documentation session;
  current M3.6 publication is unverified. Cached references do not establish it.

### Committed behavior and acceptance evidence

`HandlerRegistry.ExecuteActions(WorkflowDefinition, string) (string, error)`
validates the complete definition before callbacks, rejects unknown starting IDs,
and follows configured Next links through InvokeAction(). Each reached action is
invoked once per traversal attempt. The local cursor advances only after success.
An action error immediately returns the failed cursor and wraps the original cause
with %w. Empty cursor/nil means an action path ended; nonempty cursor/nil hands off
an unexecuted condition or WAIT. Invalid input returns the starting cursor/error.
The method does not mutate definitions or handler mappings; callbacks own their
side effects and must not mutate the definition during traversal.

All six committed tests are in workflow_actions_test.go:

- TestHandlerRegistryExecuteActionsOrderAndReadOnly: approved criterion 1,
  "Actions execute in the required order." Checks load -> transform -> store
  across four map insertion orders and two attempts each, terminal empty cursor,
  unreached actions, unchanged definition and registered mappings/dispatch.
- TestHandlerRegistryExecuteActionsFailure: approved criterion 2,
  "A failed action stops the normal path and produces an error outcome."
  First, middle and final failures return the failed cursor and preserve the
  underlying cause via errors.Is; exact callback prefixes exclude later actions.
  Each case repeats twice, establishing once per reached action per attempt.
- TestHandlerRegistryExecuteActionsUnknownHandler: missing/case/padding mismatches
  and zero-value registry return errors; an unknown middle handler stops traversal.
- TestHandlerRegistryExecuteActionsRejectsInvalidInput: unknown/empty/blank cursors,
  broken successor, invalid configured start, unsupported type, cycle and missing
  identity are rejected before callbacks.
- TestHandlerRegistryExecuteActionsBoundary: condition and 24-hour WAIT handoffs,
  both entered directly and reached after an action; no predicate evaluation,
  sleeping or successor action execution occurs.
- TestHandlerRegistryExecuteActionsStartsAtSuppliedStep: a valid last-step cursor
  succeeds without invoking or requiring registration of preceding actions.

### Historical implementation-session evidence (not rerun here)

At 395239d plus then-uncommitted M3.6 changes, now committed as f189ce6:

- Tracker validation and `python3 scripts/task_preflight.py --task M3.6` exited 0
  before activation. Explicit user approval authorized work; eligibility did not.
- RED: `go test -count=1 -run '^TestHandlerRegistryExecuteActions' ./...` first
  exited 1 from sandbox Go-cache denial. Permission-enabled rerun exited 1 with
  ExecuteActions undefined at test lines 50, 108, 135, 142, 168, 197 and 216.
- GREEN: `go test -count=1 -run '^TestHandlerRegistryExecuteActions' -v ./...`
  exited 0: six test functions and 18 subtests passed after implementation/gofmt.
- Pre-implementation and post-change `bash scripts/verify.sh` sandbox runs exited 1
  at Go unit-test setup from cache denial after tracker and 44 Python tests passed.
  Their permission-enabled reruns exited 0: tracker, 44 Python tests, Go unit/race
  tests, go vet and both Git whitespace checks passed. These environment failures
  are separate from the intentional missing-method RED failure.
- Final gofmt -l returned no files. Tracker, whitespace and scope checks passed;
  no-index whitespace checks of new files had no diagnostics (exit 1 for differences).
  These results are preserved from the implementation-session tool outputs/report,
  not newly executed RED/GREEN or preflight checks of the committed revision.

### Limitations and next action

Local cursor only; no persisted progress, automatic retries, idempotency or rollback
of completed external side effects. Repeated attempts can replay actions. No timeout,
cancellation, panic recovery or concurrent execution ownership guarantee exists in
ExecuteActions. It evaluates no conditions, schedules no WAIT/resume and integrates
with none of M2 TryClaim/Complete/Fail. Registry concurrency constraints still apply;
nil registry pointers are unsupported. No production deployment, provider integration,
measured business results or distributed exactly-once behavior is established.

Next: human review of these two documentation edits and the M3.6 engineering evidence.
M3.6 remains in_progress; M3.7-M3.10 remain pending, and M3.7 is unauthorized.
No source/test/tracker changes, commit, push, deployment or new functionality are
permitted in this documentation session. Fresh verification is recorded separately.

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

## Historical M3.7 — Wait and Resume Checkpoint (2026-11-11)

- Starting observed HEAD: `10a270486f5df4b591852db90b6264b5e6d5ab2c` on main.
  Startup tree/index were clean. Git history and commit inspection confirm exactly
  feature_list.json, workflow_wait.go and workflow_wait_test.go in this commit.
- Previous documentation checkpoint: `020e37acdab851ab4e71230e939ad1f9120569d7`.
  M3.6 acceptance/implementation-parent baseline: `ec0331523dec72ed1c6f8b7a12681458fd3dcacd`,
  identified by the user as the previous published baseline. Local history confirms
  both commits; historical pending-M3.6/uncommitted-documentation statements are superseded.
- Tracker: M3.1-M3.6 verified; M3.7 in_progress, active_task_id M3.7, WIP limit 1;
  M3.8-M3.10 pending and unauthorized. M3.7 evidence fields remain empty and unchanged.
  Committed implementation is not final acceptance or completion of the M3 milestone.
- `bash init.sh` exited 0, warning HEAD differs from cached origin/main.
  No fetch/live remote check occurred; publication of M3.7 is unverified here.
  README/roadmap historical status snapshots remain outside this edit scope.

### Committed interfaces and behavior

workflow_wait.go adds ErrWaitNotReady, ErrWaitConsumed, WaitSnapshot and WaitState,
plus `PrepareWait(WorkflowDefinition, string, time.Time) (*WaitState, error)`,
`(*WaitState).Snapshot() WaitSnapshot` and
`(*HandlerRegistry).ResumeWait(*WaitState, time.Time) (string, error)`.
WaitSnapshot contains WaitID, Cursor, ResumeAt and Consumed.

PrepareWait validates the entire definition and requested WAIT ID/type, copies the
configured Steps map, and calculates an absolute deadline once with now.Add(step.Delay).
Preparation invokes no callbacks and does not sleep. The caller retains the returned
in-memory WaitState pointer. Snapshot returns observable progress by value.

ResumeWait rejects strictly early attempts with ErrWaitNotReady, preserving progress,
cursor, deadline and consumption state. Equality at the deadline and later attempts
are eligible if unconsumed. Eligible resume consumes the wait before downstream actions,
then starts at the saved WAIT successor through ExecuteActions. Previously completed
actions are not replayed in the tested sequential scenario. Consumed attempts return
ErrWaitConsumed without replay. A terminal WAIT returns empty cursor/nil error;
subsequent condition or WAIT boundaries are handed off unexecuted. Downstream action
errors preserve the failed cursor and underlying cause; the wait stays consumed.
Callers must check errors before interpreting the cursor. Nil/zero wait states are rejected.

### Acceptance mapping and committed tests

All three existing M3.7 tracker criteria are mapped without changing the tracker:

1. "Early resume does not advance incorrectly."
   TestWaitPreparationAndEarlyResume checks the initial snapshot and three early
   attempts, including one nanosecond before the deadline, with no progress change.
2. "Eligible resume continues from the correct step without replaying earlier completed actions in the tested in-memory scenario."
   TestWaitResumeWithoutReplay covers exact-deadline and late resume, the exact
   before/after callback sequence, and duplicate rejection without replay.
   TestWaitResumeBoundaries covers terminal, condition and second-WAIT handoffs,
   directly and via an action. TestWaitFailureConsumesWithoutRetry covers first/later
   downstream failures and an unknown handler, cause preservation and consumed retries.
3. "No durable restart guarantee is claimed."
   Process-local WaitState retains a private definition/progress snapshot only;
   inspected source and the limitations below establish this boundary. Passing
   tests do not demonstrate restart recovery.

TestWaitRejectsInvalidInput covers invalid WAIT IDs/types, malformed definitions and
nil/zero state. TestWaitDefinitionSnapshot checks that later caller definition/map
changes cannot redirect the saved continuation. These are six test functions with
11 subtests (two resume, six boundary, three failure cases), matching the supplied
implementation-session report and inspected committed tests.

### Historical implementation-session evidence (not freshly rerun here)

Preserved from the user's supplied implementation-session report, at ec03315 plus
then-uncommitted M3.7 work, now committed as 10a2704:

- M3.7 preflight passed before activation. Explicit human authorization, not
  preflight eligibility, authorized implementation. No preflight rerun for this active task.
- Focused RED tests failed because the new wait/resume API was undefined.
  Focused GREEN tests passed after implementation: six test functions and 11 subtests.
- Permission-enabled full verification reruns passed: tracker validation, 44 Python
  harness tests, Go unit and race tests, go vet, and working-tree/index whitespace.
- Sandbox Go-cache permission failures were environmental, separate from the
  intentional undefined-API RED compilation failure. These historical observations
  are not newly executed failing tests or historical preflight checks in this session.

### Known limitations and next action

- In-memory state only; no durability across process restarts, automatic scheduling,
  background polling, cloud integration, database or external storage.
- WaitState requires single-caller use, retaining its pointer without copying it.
  Concurrent ResumeWait or Snapshot access is unsupported; definitions must not be
  mutated concurrently with PrepareWait. Handler registry constraints still apply.
- PrepareWait validates configuration, not proof that earlier actions executed.
  Repeated PrepareWait calls create independent states; no deduplication.
- Consumption before downstream execution prevents ordinary sequential replay,
  without guaranteeing exactly-once external effects. Downstream failures leave the
  wait consumed and require a separately defined recovery/reconciliation mechanism.
- Callback panics are not recovered; the saved cursor may not identify the exact
  failed action after a panic. Callbacks can block; no cancellation or timeout added.
- Handler registrations and captured callback inputs are not snapshotted.
- No M2 TryClaim/Complete/Fail integration, production deployment, measured business
  outcome or distributed correctness guarantee is established.

Next: human engineering review of these two uncommitted documentation files and
M3.7 evidence. Only PROGRESS.md and docs/HANDOFF.md are authorized for editing.
No source/test/tracker changes, commit, push, deployment or M3.8 work is authorized.

### Fresh M3.7 documentation-session verification (2026-11-11)

At HEAD 10a270486f5df4b591852db90b6264b5e6d5ab2c plus these two documentation edits:

- `python3 scripts/validate_features.py feature_list.json` exited 0 before editing
  and again before the unified verifier. Startup Git/commit-scope checks passed.
- Initial `bash scripts/verify.sh` exited 1 at Go unit-test setup because sandbox
  Go build-cache access was denied. Tracker and all 44 Python tests passed first;
  race, vet and whitespace checks were not reached. No final PASS was printed.
- Permission-enabled `bash scripts/verify.sh` rerun exited 0: tracker, 44 Python
  tests, `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `git diff --check` and `git diff --cached --check` passed. These are fresh local
  regression results covering M1/M2, Exercise 82 and M3.1-M3.7, not production evidence.
- Only this verification record and its handoff counterpart were added afterward.
  Final tracker validation, both whitespace checks, complete documentation diff
  review and status/scope assertions passed. Exactly PROGRESS.md and docs/HANDOFF.md
  are modified; no staged or untracked files. All other tracked files, including
  source, tests and feature_list.json, remain byte-for-byte equal to HEAD.
- Ending HEAD remains 10a270486f5df4b591852db90b6264b5e6d5ab2c. M3.7 remains
  in_progress/active with WIP 1; M3.8-M3.10 remain pending. No commit, push,
  deployment, live remote check, tracker acceptance or M3.8 work occurred.
  These uncommitted documentation changes await human engineering review.

## Session Handoff Rules

At the end of each coding session, record:

- Commit or working-tree reference
- Behavior implemented
- Files changed
- Tests actually run and their outcomes
- Unresolved errors or blockers
- Next action

Do not replace a failed check with an unsupported PASS claim.
Do not mark a feature complete without its acceptance evidence.
Update docs/HANDOFF.md when relevant state or evidence changes. Distinguish
the previous checkpoint from current HEAD and uncommitted work. A task
listed as next does not authorize implementation.

## Historical M3.8 — M2 Lifecycle Integration (2026-11-12)

- Previous documentation checkpoint: `9c57598597a1d081cc6fbc60da1c683a976a3478`.
  M3.7 acceptance is committed as `4545eef`; historical pending-M3.7 statements
  above are superseded. Starting/current observed HEAD on main:
  `842000fabc0b203a75481b1dcc42d6cb1a9b4eb8`, the supplied expected baseline.
  Startup tree/index were clean. Initialization passed and reported cached
  origin/main matches HEAD; no fetch or live remote check was performed.
- Explicit human authorization covered M3.8 only. Tracker validation and
  `python3 scripts/task_preflight.py --task M3.8` passed before activation.
  Only M3.8 status and active_task_id changed; WIP remains 1, M3.1-M3.7 verified,
  M3.9-M3.10 pending. M3.8 remains in_progress; evidence fields are preserved.
- Added workflow_runner.go and workflow_runner_test.go; changed feature_list.json,
  PROGRESS.md and docs/HANDOFF.md. All changes are uncommitted and unstaged.
  Existing M2/M3 source/tests, acceptance criteria, roadmap and AGENTS.md are unchanged.

### API and bounded behavior

`(*HandlerRegistry).RunActions(repo WorkflowLifecycleRepository, executionID, workerID string, d WorkflowDefinition) (string, error)`
uses an explicit interface containing only TryClaim, Get, Complete and Fail.
The repository run identity is separate from the definition ID; callers supply the
correct association. Nonnil initialized dependencies are required. A successful
claim and usable owned RUNNING snapshot precede ExecuteActions(d, d.StartStep).
The claimed version is captured before callbacks and never refreshed for a retry.
Execution errors use Fail with nonempty diagnostics and preserve errors.Is causes.
Only empty cursor/nil error from ExecuteActions requests Complete. Both terminal
return values are checked; rejection returns a reconciliation error, with no
fallback write or blind replay. SaveIfVersion is absent from the dependency API.

Nonempty cursor/nil is an unexecuted condition/WAIT handoff: return it and preserve
RUNNING, owner, version 2, attempts 1 and claim-only history. No boundary callback,
wait preparation, deadline evaluation, resume/finalization or scheduling is added.
Callers retain cursor/definition locally. Existing condition and wait helpers remain
separate; their results do not automatically update M2. RunActions cannot resume a
RUNNING execution. There is no new WAITING status, persistence or recovery model.

### Fresh RED, GREEN and regression evidence

- `go test -count=1 -run '^TestLifecycleRunner' ./...` exited 1 before implementation:
  RunActions undefined at test lines 89, 149, 183, 224, 275 and 295. This was a genuine
  missing-API RED compilation failure, with no cache permission failure in that run.
- After implementation/gofmt,
  `go test -count=1 -run '^TestLifecycleRunner' -v ./...` exited 0:
  six test functions and 24 deterministic subtests passed.
- Standalone tracker validation passed before activation, after activation and
  before full verification. Initial `bash scripts/verify.sh` exited 1 at Go unit
  test setup from sandbox Go-cache denial; tracker and all 44 Python tests passed,
  later checks were not reached. This environmental failure is separate from RED.
- Permission-enabled `bash scripts/verify.sh` exited 0: tracker, 44 Python tests,
  `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `git diff --check` and `git diff --cached --check` passed. Existing regression
  behavior remains passing locally. Only checkpoint documentation changed afterward;
  final tracker/formatting/whitespace/scope checks cover those additions.

### Acceptance mapping and limits

1. Success/failure preserve M2 status/version/history: TestLifecycleRunnerTerminalOutcomes
   covers success and first/middle/final action failures, exact callback prefixes,
   READY -> RUNNING -> COMPLETED/FAILED, versions 1 -> 2 -> 3, retained owner,
   attempts 1, exact two-entry history, failure cursor and LastError/history reason.
   TestLifecycleRunnerRejectedClaim confirms zero callbacks and unchanged state on
   rejection. TestLifecycleRunnerTerminalRejection covers rejected/stale completion
   and failure, original-cause preservation and no version refresh/fallback writes.
   TestLifecycleRunnerInvalidClaimSnapshot prevents callbacks on missing/invalid
   post-claim snapshots. TestLifecycleRunnerHandoffs establishes boundary treatment.
2. Fail with nonempty reason; no SaveIfVersion failure: terminal-outcome and
   TestLifecycleRunnerExecutionValidationFailure tests verify Fail invocation,
   diagnostics in state/history, malformed-definition/unknown-handler errors, and
   cause wrapping. The probe/interface expose no SaveIfVersion operation.

This is a synchronous, process-local action-path lifecycle adapter. No integrated
condition/WAIT continuation, durable cursor, scheduler, leases, retry, compensation,
panic recovery, cancellation/timeout, cloud or provider integration is established.
Lost handoffs, unusable post-claim snapshots or callback panics may leave RUNNING.
Callback side effects can succeed before a rejected terminal update; reconciliation
is required. Registry/definition concurrency constraints remain applicable. No
exactly-once external effects, production outcomes or M3.9/M3.10 acceptance claimed.
Next: human engineering review. Do not mark verified, commit, push, deploy or begin
M3.9 without separate explicit authorization.


## Historical M3.9 — Concurrency and Failure Verification (2026-11-12)

Previous implementation/documentation checkpoint revision:
`8a049e88b1d63db73e15e120b17c750179f51868` (its recorded starting HEAD was
`842000fabc0b203a75481b1dcc42d6cb1a9b4eb8`). M3.8 acceptance is committed at
`7811e9d08f8a10a4d18da6a74d4209570971daa9`, the observed starting and ending
HEAD on main for this session. Startup tree/index were clean; cached origin/main
matched HEAD. No fetch or live remote verification occurred. Git/source/tracker
inspection supersedes the historical M3.8 uncommitted/in_progress statements.

Explicit human authorization covered activation and implementation of M3.9 ONLY,
with a preference for tests only. Initialization, tracker validation and pre-activation
`python3 scripts/task_preflight.py --task M3.9` exited 0. Preflight established
eligibility, not authorization. Only M3.9 status changed to in_progress and
active_task_id to M3.9; WIP remains 1. Acceptance criteria and all evidence fields
are preserved. M3.1-M3.8 remain verified; M3.10 remains pending and unauthorized.

Added workflow_runner_concurrency_test.go: five test functions, two ownership
subtests. No production changes or defects found. The existing RunActions,
TryClaim, Complete and Fail implementation satisfies the exercised invariants.
Other changed files: feature_list.json, PROGRESS.md and docs/HANDOFF.md only.
All four files are unstaged/uncommitted; no acceptance is recorded.

### Acceptance criteria and requested scenarios

- No duplicate action execution: TestRunnerConcurrencySuccessfulExecutions (A)
  and TestRunnerConcurrencyFailedExecutions (B) each launch ten distinct workers
  through a readiness/start barrier against one repository, fully registered shared
  registry and immutable definition. Buffered result channels and a mutex-protected
  callback trace establish exactly one action path and nine claim rejections.
  Success executes load/transform/store once; failure stops at transform and
  preserves the sentinel via errors.Is. Exact terminal status, version 3, attempts 1,
  winning owner, two history entries and diagnostic reason are asserted without
  choosing a specific winner. No completion transition is accepted on action failure.
- Correct failure reasons and lifecycle state: A/B above plus
  TestRunnerConcurrencyDuplicateAfterCompletion (C) check same-owner and different-owner
  duplicates, unchanged full snapshots and no callbacks. TestRunnerConcurrencyInvalidOwnership
  (D) directly rejects non-owner Complete/Fail at the real lifecycle boundary,
  compares full snapshots, then proves owner completion still works.
- Action error plus terminal rejection (E): reuse unchanged
  TestLifecycleRunnerTerminalRejection, four rejected/stale Complete/Fail subtests.
  It checks error cause, visible rejection, exact claim/get/terminal operation order,
  version captured before callbacks, no refresh/retry/fallback and stored state
  matching accepted transitions. Ambiguous external effects require reconciliation.
- TestRunnerConcurrencyBlockedCallback (F) holds the winner RUNNING on a channel,
  rejects a second runner without callbacks or snapshot changes, then releases
  the winner to one claim and one completion. Atomic callback counters, five-second
  deadlines, buffered results and cleanup release/join bound synchronization.
  No sleeps are used to arrange contention.
- No distributed exactly-once guarantee: the limitations below apply to every test.
  Passing tests provide implementation evidence; human acceptance remains pending.

### Actual commands and outcomes

- `bash init.sh`, branch/HEAD/status/history checks, pre-activation tracker validation
  and preflight exited 0. Initialization warned cached refs are not live remote proof.
  Tracker validation after activation exited 0.
- First `go test -count=1 -run '^TestRunnerConcurrency' -v ./...`: exit 0, five
  functions and two subtests. Baseline GREEN on unchanged production code; no RED
  defect or artificial failing test was introduced.
- `go test -count=1 -run '^TestLifecycleRunner' -v ./...` and
  `go test -race -count=1 -run '^TestRunnerConcurrency' ./...` initially exited 1
  at setup from sandbox Go-cache denial. Permission-enabled reruns exited 0:
  six lifecycle functions/24 subtests and all new concurrency tests passed.
  Cache errors are environmental, not behavioral RED results.
- `go test -race -count=20 -run '^TestRunnerConcurrency' ./...`: exit 0.
- `bash scripts/verify.sh`: exit 0 on its first session invocation. Tracker, all 44
  Python harness tests, `go test -count=1 ./...`, `go test -race -count=1 ./...`,
  `go vet ./...`, `git diff --check` and `git diff --cached --check` passed.
- Only checkpoint documentation changed afterward. Final standalone tracker,
  whitespace (including the untracked Go file), formatting, scope/preservation
  checks and status inspection passed. Only the four authorized files changed;
  production source and existing tests are byte-for-byte equal to HEAD. The index
  is empty and HEAD is unchanged. Complete diffs are supplied for human review.

### Limits and next action

Guarantees apply only to one shared in-memory repository and the actual scenarios
exercised. No exactly-once distributed execution or exactly-once external provider
effects, safe retry after every timeout, durable recovery, cross-process exclusivity
or production readiness is established. Registration must finish before concurrent
registry use; callbacks must synchronize their own mutable data and definitions
must remain immutable. No leases, automatic retry, scheduler, recovery model,
provider integration, cancellation or panic recovery was added. A callback panic
or lost condition/WAIT handoff can leave RUNNING. Side effects preceding a rejected
terminal transition require reconciliation rather than blind replay.

Next action: human engineering review of these uncommitted tests, tracker activation
and checkpoints. M3.9 remains in_progress, not verified; M3.10 remains pending.
Nothing was staged, committed, pushed or deployed. Do not begin M3.10.


## M3.10 — Reusability implementation checkpoint (2026-11-12)

### Revision, authority and reconciliation

Previous implementation/documentation checkpoint revision:
`5909abf099c888536f1bda88ca2d4614a45a4508`. That checkpoint recorded its own
starting HEAD `7811e9d08f8a10a4d18da6a74d4209570971daa9` and pending M3.9 review.
Current observed starting/ending HEAD on main:
`e15efdf2da602bd67de54adb5bd8ddf7f8dc3f2f`, the required user-supplied baseline.
Local history confirms M3.9 implementation at 5909abf and acceptance at e15efdf;
these supersede earlier pending-M3.9 statements without changing historical evidence.
Startup tree/index were clean; initialization exited 0 with cached origin/main
matching HEAD. No fetch or live remote verification occurred in this session.
README and roadmap implementation-status notes lag the actual accepted M3.1-M3.9
source/tests and tracker; roadmap section 13 still defines this task's scope.

The human explicitly authorized activation and implementation of M3.10 ONLY in the
attached request. Pre-activation tracker validation and M3.10 preflight exited 0;
eligibility was not authorization. Tracker activation changes only active_task_id
from null to M3.10 and M3.10 status from pending to in_progress. WIP remains 1.
M3.1-M3.9 records/evidence, M3.10 objective/dependencies/criteria/evidence and all
other metadata are preserved. No final milestone acceptance is claimed.

### Contract, workflow shapes and test evidence

New file: `workflow_runner_reusability_test.go`. All business inputs, deterministic
handlers and simulated outputs are test-local. Both definitions validate through
WorkflowDefinition.Validate. No payment, course-access or analytics service is called.

- Enrollment fulfillment: definition `enrollment-fulfillment`, configured trigger
  `purchase.confirmed`; three actions:
  `validate-purchase` (enrollment.validate) -> `prepare-access` (enrollment.prepare)
  -> `record-fulfillment` (enrollment.record) -> terminal empty cursor.
  Handlers validate a paid purchase with student/course inputs, prepare
  `student-17:go-basics`, then record `prepared:student-17:go-basics`.
  This records a prepared request, not an actual provider entitlement.
- Learning analytics: definition `learning-analytics`, configured trigger
  `report.requested`; four actions:
  `collect-activity` (analytics.collect) -> `aggregate-activity` (analytics.aggregate)
  -> `render-report` (analytics.render) -> `record-publication` (analytics.publish)
  -> terminal empty cursor. Handlers collect [2,3,5], sum 10, render, then record
  `published:learning activity: 10` in local memory only.

These differ in chain length, step IDs, handler sets and structure (validation/request
preparation versus collection/aggregation/rendering), inputs and business output.
Each has a separate execution ID and worker. A single initialized HandlerRegistry
contains both handler sets, and the lifecycle tests share one actual
InMemoryWorkflowRepository. No registration overlaps execution. Definitions are
immutable while running; callbacks execute sequentially in these new tests.
Both use RunActions -> ExecuteActions -> named action dispatch with the existing
TryClaim/Get/Complete/Fail repository boundary. Trigger fields are configuration;
these tests invoke the runner directly and do not claim event dispatch integration.

Three test functions / four subtests:

- `TestRunnerReusabilityWorkflows/enrollment-first` and `/analytics-first` run both
  definitions in either order through RunActions. Exact callback sequences exclude
  foreign actions. Each execution starts READY/version 1/no owner/zero attempts,
  is checked RUNNING/version 2/expected owner/one attempt inside every callback,
  and ends COMPLETED/version 3/one attempt/empty LastError. Exact claim/terminal
  histories, terminal cursors, business outputs and unchanged other execution
  snapshots are asserted. The first run cannot populate the other's output.
- `TestRunnerReusabilityActionShapes/three-action-enrollment` and
  `/four-action-analytics` independently call the same ExecuteActions API on
  three- and four-step definitions, checking exact order and terminal empty cursor.
  This detects assumptions about fixed step names or action counts.
- `TestRunnerReusabilityIndependentFailure` fails analytics aggregation with an
  injected sentinel, checks errors.Is, failed cursor `aggregate-activity`, diagnostic
  execution/step/cause, matching nonempty LastError/failure-history reason, owner,
  attempts/version/history and no render/publication continuation. Enrollment's
  READY snapshot and business outputs remain unchanged. Enrollment then completes
  successfully without modifying the failed analytics snapshot or publishing a report.

Initial focused execution was **BASELINE GREEN** (exit 0), using unchanged production
code. No artificial RED, production correction or architectural capability was added.
Source inspection and final diff review establish all production Go files and prior
Go tests remain byte-for-byte equal to HEAD.

### Acceptance mapping (evidence prepared, human acceptance pending)

| Unchanged M3.10 acceptance criterion | Supporting evidence |
| --- | --- |
| Both workflows run without modifying the engine core. | TestRunnerReusabilityWorkflows, TestRunnerReusabilityActionShapes and TestRunnerReusabilityIndependentFailure; preservation comparison against HEAD and full diff review. |
| Existing M1/M2 tests and race checks pass. | Verbose full Go regression includes enrollment_test.go, workflow_test.go and workflow_repository_test.go plus every existing M3 suite; full race check and canonical verifier passed. |
| Known limitations are documented. | Stage 3 limitations below and docs/HANDOFF.md section 15. |

### Fresh executed verification and actual exit codes

- `bash init.sh`, `git branch --show-current`, `git rev-parse HEAD`,
  `git status --short`, `git diff --cached --name-only`, local history inspection,
  `python3 scripts/validate_features.py feature_list.json`, and pre-activation
  `python3 scripts/task_preflight.py --task M3.10`: exit 0.
  Tracker validation after activation also exited 0.
- `gofmt -w workflow_runner_reusability_test.go`: exit 0.
- `go test -count=1 -run '^TestRunnerReusability' -v ./...`: exit 0;
  three functions/four subtests, BASELINE GREEN.
- Initial standalone lifecycle, concurrency, concurrency-race, verbose full Go,
  full-race and vet commands below each exited 1 due to denied sandbox Go-cache
  access (test commands failed at setup). These are environmental failures,
  not behavioral RED; they do not count as passing tests.
- Permission-enabled reruns all exited 0:
  `go test -count=1 -run '^TestLifecycleRunner' -v ./...` (six functions/24 subtests),
  `go test -count=1 -run '^TestRunnerConcurrency' -v ./...` (five functions/two subtests),
  `go test -race -count=1 -run '^TestRunnerConcurrency' ./...`,
  `go test -count=1 -v ./...`, `go test -race -count=1 ./...`, `go vet ./...`.
  Verbose output was inspected: all enrollment, repository, execution and M3
  definition/validation/trigger/registry/branch/action/wait/lifecycle/concurrency/
  reusability test groups actually ran and passed, including Exercise 82 tests
  TestRepositorySaveIfVersionCannotBypassFail and
  TestRepositoryFailPreservesReasonInStateAndHistory.
- `bash scripts/verify.sh`: exit 0 on its first invocation in this session.
  All six checks PASS: tracker, 44 Python harness tests, `go test -count=1 ./...`,
  `go test -race -count=1 ./...`, `go vet ./...` and working-tree/index whitespace.
  Python command executed by the verifier:
  `PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p "test_*.py" -v`.
- Only checkpoint documentation changed after the passing verifier. Final tracker,
  formatting, tracked/index/untracked whitespace and preservation/status audits
  are recorded in the review artifacts; no new behavioral changes followed tests.

### Stage 3 limitations

Stage 3 is an in-memory engineering prototype. Demonstrated reusability means
configurations and registered callbacks drive two different tested business paths
through unchanged generic code. This does not establish production readiness.

No cross-process workflow coordination, durable execution persistence, durable
WAIT/condition cursor recovery, automatic retries, automatic compensation,
crash-safe external side-effect execution, exactly-once provider effects, guaranteed
recovery after callback panic, production-grade scheduler, production deployment
or measured business outcomes are demonstrated. No distributed infrastructure or
external provider integration was added.

Definitions must remain immutable during execution. Handler registration must
complete before concurrent registry use; populated registries must not be copied.
Mutable callback state requires appropriate synchronization. The new business
fixtures are sequential examples, not concurrent business-adapter implementations.
Verified concurrency guarantees apply only to tested behavior of one shared
in-memory repository within a single process. Ambiguous external effects require
reconciliation; successful callbacks may precede a rejected terminal write.
A lost handoff or panic can leave an execution RUNNING. RunActions stops at an
unexecuted condition/WAIT; caller-held cursors and separate helpers do not provide
a lifecycle-aware resume/finalization framework. Execution-to-definition association
is the caller's responsibility. No Stage 4 work occurred.

### Review boundary and changed files

Exactly four repository files are changed, all unstaged/uncommitted:
`workflow_runner_reusability_test.go` (new), `feature_list.json` (activation only),
`PROGRESS.md` and `docs/HANDOFF.md` (checkpoint reconciliation).
Complete test, activation and documentation patches are provided outside the
repository under `/private/tmp/p1-m3.10-review/` for human engineering review.
M3.10 remains in_progress, active_task_id M3.10; M3.1-M3.9 remain verified; WIP 1.
No M3.10 acceptance commit exists. Nothing was staged, committed, pushed or deployed.
Stop for human engineering review; separate authorization is required for the
implementation commit, publication, final verification status and Stage 3 acceptance.
