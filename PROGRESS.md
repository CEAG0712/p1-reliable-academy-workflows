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
- Current planning checkpoint (2026-11-13): M3.1-M3.10 are verified in
  published acceptance commit 5457b77. The unified tracker adds M4.1-M4.10
  as pending only; active_task_id null, global WIP 1. See the Stage 4
  tracker-preparation checkpoint below; no API implementation is authorized.

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

M3 accepted; Stage 4 Gin API tracker preparation only (M4.1-M4.10 pending).
No task is active. The M3 capability/checkpoint history follows; the current
planning checkpoint and fresh evidence are at the end of this document.

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


## Historical M3.10 — Reusability implementation checkpoint (2026-11-12)

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


## Stage 4 unified tracker preparation checkpoint (2026-11-13)

Previous documentation checkpoint revision: `ae17f66b6c91a394050eb4fc60c33087e064e6a3`.
Current observed HEAD on main: `5457b77e7835cf951b1c71ab332080d7b81fe735`,
`chore: finalize M3.10 and Stage 3 acceptance`. Live remote main matched that
revision via permission-enabled `git ls-remote origin refs/heads/main` (exit 0).
Its first sandbox attempt exited 128: github.com could not be resolved.
Starting tree/index were clean. `bash init.sh` and standalone tracker validation
exited 0. Initialization's cached-ref warning is separate from the live check.
Local history, accepted tracker and passing source/tests supersede the older
pending M3.10 checkpoint. README/roadmap historical status text still lags;
the roadmap's Stage 4 scope remains planning, with HTTP/trust decisions open.
The supplied M4.1-M4.10 breakdown is the authorized tracker scope; roadmap
section 6's older M4 enrollment-integration label is separate from these API IDs.

Human authorization covers local tracker, instructions, regression tests and
truthful checkpoint updates only. `feature_list.json` now has 20 records:
all ten original M3 objects/evidence/order retained, followed by ten pending
M4 contracts copied from the supplied request. The chain is M3.10 -> M4.1 ->
M4.2 -> M4.3 -> M4.4 -> M4.5 -> M4.6 -> M4.7 -> M4.8 -> M4.9 -> M4.10.
Schema version 1, global WIP 1 and active_task_id null are preserved. Every M4
blocker is null and evidence is null/empty. Status meanings and rules are unchanged.
Gin v1.11.0 is only a compatibility candidate, with no dependency approval or
installation. Health status/body/method behavior and API trust contracts still
require decisions during separately authorized work; none are invented here.

AGENTS.md now applies tracker selection, preflight and global WIP to both stages.
Seven validator/preservation tests and four preflight tests were added; the
original 44 fixture tests remain unchanged. The validator, preflight and verifier
implementations already support the unified tracker and were not modified.
Preservation tests compare complete M3 objects against the immutable local Git
baseline. Preparation tests require empty M4 evidence; schema v1 itself permits
populated pending evidence. A separate regression demonstrates that distinction.
Schema v1 enforces null blockers for nonblocked tasks and global active/WIP rules.
Preflight tests check input bytes and temporary-directory contents on both eligible
and rejected paths, and explicitly preserve the human-approval boundary.

Fresh checks before final canonical verification:

- Baseline `bash scripts/verify.sh`: sandbox exit 1 at Go unit setup after 44
  Python tests passed (Go build-cache access denied); permission-enabled rerun
  exit 0, all six checks PASS. This is environmental failure, not behavioral RED.
- Post-extension `python3 scripts/validate_features.py feature_list.json`: exit 0.
- `PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v`:
  exit 0, 55 tests (44 original plus 11 new). Bytecode writes disabled.
- `python3 scripts/task_preflight.py --task M4.1`: exit 0, ELIGIBLE, HUMAN APPROVAL
  REQUIRED, no task activated, no implementation authority.
- `python3 scripts/task_preflight.py --task M4.2`: exit 1, expected negative test:
  depends on unverified M4.1. No mutation occurred.
- `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`:
  initial sandbox invocations each exited 1 from denied Go build-cache access;
  permission-enabled reruns each exited 0. No behavioral test failure occurred.
  Nonverbose Go output reports package success without test counts.

Exactly six permitted files are edited, unstaged/uncommitted: feature_list.json,
AGENTS.md, scripts/test_validate_features.py, scripts/test_task_preflight.py,
PROGRESS.md and docs/HANDOFF.md. Original published tracker snapshot and complete
review material are outside the repository at `/private/tmp/p1-stage4-tracker-review/`.
No Go source/tests, dependencies, API, infrastructure or roadmap changes. Existing
process-local engine limitations remain authoritative; these harness results do
not prove API, cloud, end-to-end or production behavior.

Next action: human engineering review of the complete local preparation diff.
No task activation, staging, commit, push, deployment, Gin installation or API
implementation is authorized. No future revision is claimed. At next startup,
recheck HEAD/tree/index/live state and reconcile this uncommitted checkpoint.

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


## M4.4 workflow-service implementation checkpoint (2026-11-14)

Previous checkpoint's recorded revision: `89c1317b2785df994036bb6ea7ca087599ca828b`
(M4.3 implementation/acceptance started there). Current observed starting HEAD:
`8a447b5a5bce896cf929e90da5e5ffaf9890e05e`, main; live remote main matched
independently (git ls-remote exit 0 after sandbox DNS exit 128). Git history now
confirms published, verified M4.3; older local/uncommitted statements are historical.
Starting tree clean, index empty. init.sh, tracker validation, M4.4 preflight and
baseline canonical verifier all exited 0. Initialization's cached-remote caveat
was reviewed and resolved by the live lookup. README/roadmap status notes lag code.

The human's bounded M4.4 request authorizes activation and implementation after
baseline/design reconciliation. Before activation, the A-G criterion-2 design
was presented in chat: real unregistered Gin handler -> transport-independent
WorkflowService.PrepareSubmission -> read-only DefinitionReader.LookupDefinition,
identity match, existing WorkflowDefinition.Validate and configuration summary.
This checks the actual server graph, unlike M4.3 request syntax/support membership.
No execution is created, claimed, dispatched, scheduled, resumed or authorized.
No criterion text was amended. Only M4.4 was activated; schema/WIP remain 1.
M4.4 remains in_progress, active_task_id M4.4, revision/reviewed_by null, pending
human source review. M4.5-M4.10 remain pending and all other task objects unchanged.

Exactly seven files change:

- internal/application/submission.go: explicit WorkflowService command/result and
  read-only DefinitionReader contract, constructor rejecting nil/typed nil with
  zero dependency calls. Lookup receives context and definition ID exactly once;
  identity mismatch or invalid graph returns ErrInvalidDefinition. Missing,
  unavailable and unexpected dependency errors retain errors.Is identity through
  wrapping; failed operations return zero result. Nil/zero Service fails closed;
  empty identifiers/nil context reject, cancelled context avoids lookup. Transport
  syntax validation remains the adapter's responsibility. No HTTP/Gin types.
- internal/application/submission_test.go: four deterministic top-level tests,
  construction, exact lookup/context/result, repeated correlation without execution,
  failure/zero-result behavior and cancellation. A real pre-existing FAILED
  repository snapshot retains owner/version/attempts/history/diagnostic unchanged.
  This repository is deliberately not injected; capability/source inspection is
  the primary proof that this boundary cannot perform lifecycle writes or grants.
- internal/httpapi/preparation.go: NewPreparationHandler returns a real Gin
  handler but is never composed by NewRouter/server. It snapshots the server
  support set, reads at most 4097 bytes, reuses unchanged M4.3 validation, maps the
  two fields to the application command and passes request context. The isolated
  adapter's HTTP 200 JSON summary has exactly workflow_definition_id, event_id,
  step_count. It is not a published endpoint or CreateExecutionResponse; no
  execution ID, lifecycle state or acceptance receipt is fabricated. No new
  status number or public error code is introduced. ErrorResponse/APIError remain
  unchanged: invalid command 400 INVALID_REQUEST, missing definition 404 NOT_FOUND,
  unavailable dependency 503 UNAVAILABLE, bad server graph/unexpected errors
  500 INTERNAL_ERROR. Internal messages are never serialized.
- internal/httpapi/preparation_test.go: six new top-level tests (19 HTTP functions
  total), nil/typed nil service, no constructor call, support-set copy, all 74
  existing rejection inputs through real handler with zero service calls, additional
  tenant/caller/payment-field rejection, bounded infinite stream/read-error and
  exact 4096-byte control, exact command/context/result, wrapped safe error matrix.
  Real handler+real service+controlled definition reader succeeds on a valid graph
  then rejects a broken graph for identical valid JSON. This proves meaningful
  delegation beyond a fake call counter. An action name entitlement.grant is
  configuration only: no registry/provider exists in this composition.
- feature_list.json: M4.4 activation and actual evidence only.
- PROGRESS.md and docs/HANDOFF.md: appended factual checkpoint, histories retained.

RED: focused application and HTTP sandbox attempts exited 1 at denied Go-cache
access (environmental only). Approved permission-enabled retries exited 1 from
undefined DefinitionReader/NewService/command/result/errors and no non-test
application files respectively: genuine pre-implementation RED. GREEN focused
application and HTTP commands exited 0 after implementation. gofmt on four new
Go files exited 0. Additional read-boundary coverage was added before full checks.

Fresh verification commands (successful runs each exit 0):

- python3 scripts/validate_features.py feature_list.json (before/after activation)
- PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
  (60 tests)
- go test -count=1 -v ./internal/application/... (four functions)
- go test -count=1 -v ./internal/httpapi/... (19 functions)
- go test -count=1 ./... (four packages)
- go test -race -count=1 ./... (four packages)
- go vet ./... (no diagnostics)
- bash scripts/verify.sh (all six checks)
- git diff --check; git diff --cached --check
- gofmt -l on the four new Go files (empty output)

Full unit/race/vet/canonical sandbox attempts exited 1 at Go-cache permissions;
approved permission-enabled retries all exited 0, no behavioral failure/source fix.
Exact command arguments, exits and logs: /private/tmp/p1-m4.4-review/commands.jsonl
and named .txt logs, with reproducible run_checks.py. Initial baseline commands
are recorded in the session transcript and baseline.json; baseline verifier passed.

One-time route audit: python3 /private/tmp/p1-m4.4-review/route_audit.py, exit 0.
External Go test overlay calls actual NewRouter, observes only GET /health with
unchanged JSON/content type, and default 404 for execution creation/status,
/test-only/prepare, /enrollment and /payment. Exact route enumeration excludes
all other business paths. Test-only handler route is not operationally reachable.
Go import inspection shows Gin absent in engine/application; application imports
the existing engine model for Validate. HTTP's dependency graph now includes that
model package, which is explicitly not an invocation path. Source audit finds no
runner/repository/provider invocation in new production files. Router/server bytes
unchanged. Audit source, overlay, output and command exits retained externally;
no permanent test freezes future authorized routing work.

Preservation audit compares every pre-existing file with the full starting SHA,
all 20 task records and original M4.4 objective/dependency/six criteria. New code
cannot bypass TryClaim/expectedVersion/Fail, because no execution repository or
runner is injected or called. Existing engine/repository/tests, DTOs, validator,
API contract document, dependency files and harness remain byte-identical.
Final documentation checks, preservation audit and complete tracked-plus-untracked
patch are recorded in the external review directory. Tests establish local
behavior only; they do not constitute human acceptance or cloud integration.

Trust limitations remain: no caller authentication, tenant ownership, workflow
permission, payment/provider authenticity, replay defense, durable storage/queue,
production provider, distributed coordination or scheduling. EventID is correlation
only; even CONFIRMED is no payment proof. Support membership is no authorization.
The reader must be trusted, read-only and return stable configuration during a call;
this interface cannot police arbitrary malicious dependency implementations.
Preparation checks graph structure only, not registered handlers, event contents
or execution readiness; its result must never be reused as authorization or a
validated executable snapshot. No production definition reader is wired.

No staging, commit, push, deployment, dependency installation or M4.5 activation
occurred. Current HEAD remains 8a447b5; tree is unstaged/uncommitted and index empty.
STOP for human source review; do not mark verified or clear active_task_id. At
next startup reconcile actual HEAD/tree/tracker/tests and run init.sh; temporary
review artifacts are outside the repository and not publication evidence.

Final post-checkpoint verification (2026-11-14): tracker/Python (60 tests), focused
application (4 functions)/HTTP (19 functions), full Go unit/race (4 packages), vet,
canonical verifier, formatting and working-tree/index whitespace passed on their
successful runs (all exit 0). Final sandbox Go checks initially exited 1 from
cache permissions; approved retries all exited 0. No behavioral failure or source
fix. Corrected external preservation audit exits 0; its first audit-only failure
misread git diff --no-index exit 1 for clean new-file differences as a failure,
then was corrected to check empty diagnostics and whitespace directly. No project
code changed. Exact final logs and manifest are retained externally. Only factual
evidence text changed after regressions; final tracker/whitespace/preservation
and complete patch reverse-applicability checks are recorded in final-artifacts.txt.
No human acceptance/publication occurred; STOP for human source review.


## M4.4 acceptance reconciliation checkpoint (2026-11-14)

Previous committed documentation checkpoint / current observed starting HEAD:
`8a447b5a5bce896cf929e90da5e5ffaf9890e05e` (published M4.3), main. The prior
implementation checkpoint recorded this same baseline plus seven uncommitted
M4.4 files; its earlier historical revision 89c1317 remains historical. Live
remote main independently matched the full current SHA (approved lookup exit 0
after sandbox DNS exit 128). Starting tree contained exactly the expected seven
M4.4 files, index empty. init.sh and tracker validation exited 0; expected dirty
tree notice reviewed. No activation preflight rerun: M4.4 was already active.

Human engineering owner explicitly reports source-level acceptance of the
seven-file implementation against all six original criteria with no correction,
and authorizes acceptance reconciliation only. Codex separately inspected source
and performed the local commands below. ChatGPT is not claimed to have executed
local tests; the human is not claimed to have personally run them. This checkpoint
supersedes earlier pending-review status without rewriting its history/evidence.

Complete tracked-plus-untracked starting patch matched the human-reviewed
/private/tmp/p1-m4.4-review/Complete-review.patch byte-for-byte, SHA256
`ec1e71d769909747202e103f5a0cbf7d6bebd09d5ecd3e80cf1ffe6000f963ef`.
Accepted seven files: internal/application/submission.go and submission_test.go;
internal/httpapi/preparation.go and preparation_test.go; feature_list.json;
PROGRESS.md; docs/HANDOFF.md. The four Go files remain byte-identical to that
reviewed source; exactly the three checkpoint/tracker files change this session.

Six-criterion acceptance mapping (original tracker text remains unchanged):

1. Gin stays in HTTP transport; application command/result/interface use context
   and model values. Generic engine imports/source are unchanged.
2. Real NewPreparationHandler invokes WorkflowService.PrepareSubmission with exact
   command/context; TestPreparationHandlerDelegation observes one call and JSON.
   TestPreparationHandlerRealService uses the real service and existing graph
   Validate: identical valid JSON succeeds for valid server configuration and
   fails for invalid configuration. The accepted criterion-2 interpretation is
   meaningful read-only preparation; the handler remains unregistered.
3. TestServiceConstruction/PreparationErrors/PreparationInvalidUse and handler
   tests cover nil/typed nil, no construction calls, nil/zero service, cancellation,
   lookup arguments, identity mismatch, bad graph, errors.Is propagation, zero
   failure result, support-set snapshot and bounded body reads. Nil context rejects
   in inspected source; no separate nil-context test is claimed. Safe errors map
   to unchanged 400 INVALID_REQUEST / 404 NOT_FOUND / 503 UNAVAILABLE / 500
   INTERNAL_ERROR envelopes; wrapped private diagnostics are excluded.
4. The read-only service receives no execution repository, registry or provider;
   its only dependency call is DefinitionReader.LookupDefinition and its model
   operation is Validate. State-preservation test keeps an existing FAILED snapshot
   intact. Existing owner/version/history/diagnostic/single-winner/WAIT/ambiguous-
   effect regressions pass unchanged. Preparation does not itself exercise those
   lifecycle transitions. SaveIfVersion remains completion-only; Fail authoritative.
5. All existing source/tests, runner/repository contracts, API_CONTRACTS.md,
   DTOs/validator, router/server, dependencies and harness remain byte-identical
   to publication. No engine redesign or invocation machinery added.
6. Fresh pre-status regression results passed as recorded below.

Fresh BEFORE-status checks: each successful run exit 0:

- python3 scripts/validate_features.py feature_list.json
- PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
  (60 tests)
- go test -count=1 -v ./internal/application/... (four test functions)
- go test -count=1 -v ./internal/httpapi/... (19 test functions)
- go test -race -count=1 ./internal/application/...
- go test -race -count=1 ./internal/httpapi/...
- go test -count=1 ./...; go test -race -count=1 ./... (four Go packages)
- go vet ./... (no diagnostics); bash scripts/verify.sh (all six checks)
- git diff --check; git diff --cached --check
- gofmt -l on four accepted Go files (empty output, no writes)
- python3 /private/tmp/p1-m4.4-acceptance-review/preservation_audit.py pre
  (exact accepted starting checkpoint, source bytes, task objects, empty index,
  untracked-file whitespace and seven-file scope PASS)

Initial sandbox focused unit/race, full unit/race, vet and canonical attempts
all exited 1 at denied Go-cache access. Each approved permission-enabled retry
exited 0; no source defect, source fix or new behavioral RED. Logs and exact argv/
exit manifest are external under /private/tmp/p1-m4.4-acceptance-review/.

Historical one-time route audit source, overlay and raw output were read and
copied into historical-implementation/ there. It observed actual NewRouter with
only GET /health, unchanged 200 JSON/content type, default 404 for execution
creation/status, test preparation, enrollment/payment. This audit was NOT freshly
executed this acceptance session. Current router/server source is byte-identical
to the audited published baseline. Inspected new source has no runner/repository/
provider calls; model import for read-only Validate is permitted, not invocation.
No permanent future-route prohibition test was introduced.

Only after pre-status checks passed, M4.4 changed in_progress -> verified and
active_task_id became null. Schema/WIP remain 1, all 20 task records retained;
M3.1-M3.10/M4.1-M4.3 complete verified objects unchanged; M4.5-M4.10 complete
pending objects unchanged. Original M4.4 objective, dependency and six criteria
intact, earlier implementation/design/preflight/activation/RED/GREEN/error/audit
history retained. Revision field explicitly distinguishes published baseline,
accepted uncommitted source and local reconciliation; no M4.4 SHA invented.

Trust limitations persist: no caller authentication/authorization, tenant ownership
or isolation, payment-provider authenticity/replay protection, production definition
reader, registered-action availability guarantee, execution identity/creation,
durable persistence/scheduling, real provider integration, cross-process concurrency,
distributed idempotency, rate limiting or operational HTTP execution endpoints.
The injected reader must be trusted, read-only and stable for the call. EventID
is correlation, not payment or idempotency proof. Support membership is no permission;
a preparation summary proves neither caller authority nor execution readiness.

Acceptance is local/uncommitted; observed HEAD remains 8a447b5, index empty.
No implementation changes, new routes, workflow invocation, dependency install,
M4.5 work, staging, commit, push or deployment occurred. GitHub still contains the
published M4.3 baseline. Fresh AFTER-status verification and final independent
preservation checks are recorded separately below and in external artifacts.
STOP for human publication review. Do not infer permission to publish or activate
M4.5. Next startup: run init.sh and reconcile actual HEAD/tree/tracker/tests.

Final AFTER-status acceptance checks: tracker/Python (60 tests), focused application
(4 functions)/HTTP (19 functions), focused race checks, full Go unit/race (4 packages),
vet, canonical verifier (all six checks), formatting and working-tree/index
whitespace each passed on their successful runs, exit 0. Initial post-status
sandbox Go checks each exited 1 from denied Go-cache access; approved retries
all exited 0. No behavioral source failure or implementation change. Independent
post-preservation audit exited 0: exactly three acceptance-only edits, four accepted
Go files unchanged, other 19 task objects unchanged, original six criteria intact,
M4.4 verified, no active task, empty index. Only factual evidence text was appended
after regressions. Final tracker/whitespace/preservation and both patch reverse-
applicability checks are retained in final-artifacts.txt. Acceptance-only and
complete seven-file patches, fresh verification report, six-criterion mapping,
trust report and historical audit copies are in the external review directory.
Acceptance remains local/uncommitted; STOP for human publication review.


## M4.5 workflow invocation implementation checkpoint (2026-11-14)

Previous checkpoint's recorded starting revision: `8a447b5a5bce896cf929e90da5e5ffaf9890e05e`.
Current observed starting HEAD: `4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2`, main.
Git history and a fresh live origin/main lookup confirm published verified M4.4;
older uncommitted/publication-pending statements remain historical. Initial tree
clean and index empty. init.sh, tracker validation, M4.5 preflight and baseline
canonical verifier exited 0 (60 Python tests, four Go packages). Remote lookup
first exited 128 from sandbox DNS; approved retry exited 0 with the full required
SHA. Initialization's cached-remote warning was resolved by that live lookup.
Baseline raw command output is in the session transcript, separate from later
external command logs. README/roadmap status prose still lags actual source.

The human-submitted bounded implementation handoff explicitly approves activation
and the prior seven-criterion architecture, subject to three guardrails: local scope
is NOT production authentication; creation uses the existing repository without
changing old methods; the observer delegates exactly once and preserves ambiguity.
After baseline/design reconciliation only M4.5 became in_progress and active_task_id
M4.5. Schema/WIP remain 1, all 20 records retained; other 19 complete task objects,
original M4.5 objective/dependency/seven criteria and prior evidence are unchanged.
M4.6-M4.10 remain pending. M4.5 revision/reviewed_by are null pending source review.

Exactly eleven permitted files change:
- workflow_repository.go adds CreateIfAbsent(bool,error), validation and standard
  imports only. It clones fresh READY/v1 data, validates before duplicate handling,
  then checks/inserts under the existing map's mutex. true/nil proves this insertion;
  false/nil means duplicate; invalid input/uninitialized storage is false/error.
  Save, clone/Get, TryClaim, Complete, Fail and SaveIfVersion remain byte-identical.
- workflow_repository_creation_test.go adds four functions covering invalid initial
  fields, nil/zero storage, metadata/snapshot preservation across all lifecycle
  states, same/different-key concurrency, concurrent Save competition and explicit
  Save-first/Create-first ordering. Race checks pass; no alternate map exists.
- internal/application/invocation.go adds separate invocation contracts, typed
  errors, private dependency injection, orchestration and lifecycle observer.
  Exact authority and fixed fixture membership precede resolution/identity/creation.
  Plans are cloned, identity/graph/safe registrations checked, IDs validated and
  explicit insertion checked. Same repository handles creation, runner lifecycle
  and result snapshots. Real HandlerRegistry.RunActions remains unchanged.
  Observer forwards each method once with original arguments/results. Coherent
  snapshots and observed writes distinguish terminal success/failure, handoff,
  claim rejection and reconciliation. No error-text parsing or authoritative
  application state mutation. No retries, replacement execution or replay.
- internal/application/local_invocation.go supplies private nonzero-sized authority,
  scope-bound commands, fixed local.success/failure/wait/condition fixtures and
  standard-library CryptoExecutionIDSource (exec- plus 32 lowercase hex digits).
  Fresh private registries contain only no-op, fixed error and true/nil condition
  callbacks. Successful-registration manifests cover the entire graph, including
  disconnected/post-boundary steps. Public constructors accept no arbitrary graph,
  resolver, callback or runner. Construction calls no operational collaborator.
  Initialized repository storage remains a caller precondition, enforced at creation;
  constructor does not inspect repository state. Trusted injected Go dependencies
  are not sandboxed. Scope/ID-source concurrency constraints are documented.
- internal/application/invocation_test.go adds eighteen functions covering dependency
  nil/typed-nil, no construction calls, exact context/arguments, graph identity and
  handler failures, snapshot isolation, direct/JSON/foreign capability rejection,
  generation failures, duplicates, repeated event correlation, concurrent single
  creation, real runner ownership/history, failure diagnostics, WAIT/conditions,
  terminal rejection despite another terminal state, inconsistent/missing snapshots,
  observer delegation and cancellation before/after creation/during a callback.
- internal/httpapi/invocation.go adds an unregistered bounded handler. It reuses
  unchanged M4.3 validation and M4.2 DTOs; public construction takes only a sealed
  local scope. Test-only private injection records delegation/errors. Real creation
  can return 201 COMPLETED/FAILED/RUNNING or READY created_not_started; errors use
  safe existing envelopes. Uncertain/inconsistent outcomes never project success.
- internal/httpapi/invocation_test.go adds nine functions. The real isolated Gin
  route delegates through scope/service/runner to a real repository, with exact
  context/command propagation. Success/failure/handoff/duplicate, all 74 previous
  rejection inputs, extra trust fields, 4097-byte bounded stream/read error,
  exact 4096-byte control, safe errors and missing/foreign authority are covered.
  READY interruption projection uses an explicitly created real record with a
  controlled result; real cancellation-after-creation is tested in application.
- docs/API_CONTRACTS.md appends the local-only trust/identity/outcome/HTTP contract.
- feature_list.json records activation and genuine results only.
- PROGRESS.md and docs/HANDOFF.md append this checkpoint, retaining history.

RED/GREEN evidence: new creation API tests initially failed from missing method
(structural RED, exit 1), then GREEN/race exit 0. Application/HTTP first attempts
exited 1 from Go-cache denial; approved retries exited 1 for missing new symbols
(structural RED), then implementations passed. Additional genuine behavioral REDs
caught claim rejection hiding an inconsistent final snapshot and invalid HTTP
creation-outcome projection; each was corrected in authorized new source and
returned GREEN. No environmental failure is counted as behavioral RED. Two
behavioral fixes preserve reconciliation rather than inferring success from state.

Fresh successful checks (exit 0): focused creation/application/HTTP unit tests;
creation/application/HTTP race checks; full Go unit and race across four packages;
go vet; 60 Python harness tests; tracker validation; canonical verifier (all six
checks); one-time route overlay; static trust audit; preservation audit. Some Go
checks initially exited 1 from denied cache access; approved retries exited 0.
Exact argv, exits, raw logs and RED distinctions are retained under
/private/tmp/p1-m4.5-review/commands.jsonl and named logs. No dependency install.
The first external preservation audit's fragment-format check miscounted a trailing
separator newline; corrected audit passed without a project-source change.

One-time external route overlay freshly instantiated actual NewRouter: exactly
GET /health, unchanged JSON/content type; execution create/status, test-only invoke,
enrollment, payment and admin paths returned default 404 for GET/POST. Router,
server and middleware source bytes remain unchanged; no permanent future-route
prohibition test was installed. Trust/semantic audit inspects public setup, command
sealing, fixed callback bodies, same-store use and exact observer delegation.
Preservation verifies old source/tests/dependencies/harness, repository method bytes,
other task objects, append-only histories, authorized scope and empty index.
Only new repository code is formatted; old pre-existing repository formatting is
retained byte-for-byte as required. New Go files and additive method formatting
and new-file whitespace are independently checked.

Limits: local capability is not external caller authentication or tenant/payment
proof. No payment/entitlement workflow or real/fake provider grant is invoked.
EventID (including CONFIRMED) is correlation, not authorization or an idempotency key.
Definition/event/cursor association is transient, not stored execution history or
M4.6 visibility. WAIT/condition remains RUNNING/v2; no stored WAITING, condition
execution, scheduler or automatic continuation. Cancel before creation prevents
creation; after insertion/before runner it retains READY/v1. Once synchronous
RunActions begins, cancellation does not stop callbacks or erase observed outcomes.
Rejected terminal writes or missing/inconsistent state require reconciliation even
if the final snapshot is terminal. Crashes/panics/lost responses have no automatic
recovery. No durable store/queue, cloud/provider integration or distributed
exactly-once behavior is established by these tests.

Current observed HEAD remains 4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2. All eleven
files are unstaged/uncommitted, index empty. No staging/commit/push/deployment,
acceptance, M4.5 verification or M4.6 activation. External complete patch, source
report, criterion matrix, creation/trust/lifecycle/route/preservation audits,
verification/RED logs, limitations and review checklist support human source review.
STOP before acceptance/publication; leave M4.5 in_progress and active_task_id M4.5.
At next startup reconcile actual HEAD/tree/index/tracker and run init.sh; do not
confuse uncommitted implementation, test passes or design approval with acceptance.


Final M4.5 post-checkpoint verification: canonical-final-retry exited 0 after
canonical-final exited 1 at sandbox Go-cache access. All six checks passed, including
60 Python tests and full unit/race across four Go packages, vet and whitespace.
Four creation functions (21 RUN entries including parents/subtests), eighteen new
application and nine new HTTP functions are present; focused-final executed 22
application and 28 HTTP functions (426 RUN entries including parents/subtests).
Final preservation confirms exactly eleven allowed paths and original methods/task
objects/history retained; tracker and both whitespace checks passed. Complete patch
reverse-check exited 0. Only this factual evidence and tracker command/result text
changed after those regressions; final tracker/preservation/patch checks and refreshed
artifacts are outside the repository. M4.5 remains in_progress, active_task_id M4.5,
revision/reviewed_by null, index empty, HEAD unchanged. STOP for human source review.


## M4.5 acceptance reconciliation checkpoint (2026-11-14)

Previous implementation checkpoint's recorded revision and current observed starting
HEAD: `4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2`, main (published M4.4).
The preceding implementation checkpoint records this same baseline with eleven
uncommitted M4.5 paths. Its earlier 8a447b5 reference remains historical. Fresh
live origin/main matches the full starting SHA: sandbox DNS exit 128, approved
retry exit 0. init.sh/tracker validation exited 0; expected eleven-file dirty-tree
notice reviewed, index empty. No activation preflight rerun: M4.5 was already active.

Human owner reports ChatGPT source-level review accepted the exact complete
implementation against all seven original criteria without required corrections,
and explicitly authorizes acceptance reconciliation only. Codex separately matched
patch bytes, inspected source and ran fresh local tests. ChatGPT is not claimed to
have independently executed local commands; the human is not claimed to have
manually run tests. This supersedes pending-review status without rewriting history.

Complete reviewed patch matched reconstructed tracked-plus-untracked bytes:
SHA256 `09302f37efbb41c3342da9cf9a525186ed0a3a454c7d0f20d63c3734c46804f2`.
Eight frozen files retain recorded hashes: workflow_repository.go;
workflow_repository_creation_test.go; internal/application/invocation.go;
internal/application/local_invocation.go; internal/application/invocation_test.go;
internal/httpapi/invocation.go; internal/httpapi/invocation_test.go;
docs/API_CONTRACTS.md. The other three reviewed files (feature_list.json, PROGRESS.md,
docs/HANDOFF.md) matched the starting artifact and were saved externally before
acceptance edits. Original supporting artifacts were inspected and copied as
historical evidence. No Go/test/contract formatting or implementation edit occurred.

Seven original criteria, text unchanged:
1. A valid local request selects an approved workflow definition.
2. Execution identity and creation behavior follow a documented contract.
3. Invocation uses the existing lifecycle-controlled runner.
4. Successful and failed paths preserve owner/version/history invariants.
5. Untrusted client inputs cannot independently authorize payment or course access.
6. No unsupported durable queue, real provider, or distributed exactly-once behavior is claimed.
7. HTTP-to-engine tests and prior regressions pass.

Source/evidence mapping: (1) fixed server-owned local fixture resolver, exact ID,
existing graph Validate and full safe-registration manifests; Plans/ExactDelegation/
HandlerRealEngine tests. (2) server-generated exec- plus 32 lowercase hex digits,
explicit atomic insertion/duplicate/error semantics and initial snapshot verification;
four creation tests, IDFailures/DuplicateAndCorrelation/ConcurrentSameID.
(3) real unchanged HandlerRegistry.RunActions, same-store observer and exact arguments;
ExactDelegationAndSnapshot/CallbackOwnership/HandlerRealEngine.
(4) authoritative owner/version/history and nonempty diagnostic, accepted writes plus
coherent snapshots; RealOutcomes/FailureCauseAndOwner/HandoffBoundaries/ObserverDelegates/
TerminalRejections/ClaimAndSnapshotRejections, plus existing repository/runner races.
(5) sealed scope/direct/JSON/foreign-token gates; fixed harmless local.noop/fail/check;
TrustGate/PrecreationRejections/JSONCannotMintAuthority and HTTP authority/rejection tests.
(6) source/contract inspection and cancellation/correlation/handoff tests preserve all
local-only limitations. (7) real isolated Gin-to-service-to-runner-to-repository path,
existing regressions unchanged; fresh full tests and production route audit pass.
Full criterion/source/test/regression/limitation mapping is external.

All three architectural guardrails hold: local scope is NOT production authentication;
CreateIfAbsent uses the existing map/mutex, leaving Save/Get/TryClaim/Complete/Fail/
SaveIfVersion bodies unchanged; observer delegates each original operation once and
retains reconciliation after rejected terminal updates, even with final COMPLETED.
No action replay, version-refresh retry, independent state write or alternate runner.
READY/v1 -> RUNNING/v2 -> COMPLETED/v3 or FAILED/v3 retains owner/attempt/history.
WAIT/condition hands off RUNNING/v2 without evaluation/resume or stored WAITING.
Creation interrupted before runner can remain READY/v1; synchronous callbacks have
no context cancellation support. Definition/event/cursor association is transient.

Fresh BEFORE-status checks, successful runs each exit 0:
- go test -count=1 -v -run '^TestRepositoryCreateIfAbsent' . (4 functions, 21 RUN entries)
- go test -race -count=1 -run '^TestRepositoryCreateIfAbsent' .
- go test -count=1 -v -run '^TestLocalInvocation' ./internal/application (18 functions)
- go test -race -count=1 ./internal/application/...
- go test -count=1 -v -run '^TestLocalInvocationHandler' ./internal/httpapi (9 functions)
- go test -race -count=1 ./internal/httpapi/...
- go test -count=1 -v ./internal/application/... (22 functions)
- go test -count=1 -v ./internal/httpapi/... (28 functions)
- go test -count=1 ./...; go test -race -count=1 ./... (four packages)
- go vet ./...; PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v (60 tests)
- python3 scripts/validate_features.py feature_list.json; bash scripts/verify.sh (all six)
- git diff --check; git diff --cached --check
- read-only formatting/new-file whitespace, fresh trust and preservation audits
- fresh external route overlay calling real NewRouter (only GET /health, unchanged
  status/JSON/content type; execution creation/status, test invoke, enrollment,
  payment and admin paths default 404 for GET/POST). No listener or permanent test.

Initial pre-status Go unit/race/vet/canonical/route attempts exited 1 from sandbox
cache access; approved retries exited 0. Python/tracker/trust/preservation/whitespace
passed directly. No behavioral source failure, correction or new RED. Exact argv,
environment, exits, test counts and raw logs: /private/tmp/p1-m4.5-acceptance-review/.
Code/source inspection is separate from the human-reported review; no audit proves
production security or distributed behavior. Eight frozen hashes and prior method/
source/test/dependency/harness bytes are preserved. New Go files/additive method
formatting inspected read-only; legacy repository formatting remains untouched.

Only after those gates passed, M4.5 in_progress -> verified; active_task_id -> null.
Schema/WIP remain 1, all twenty task records retained, other nineteen complete task
objects unchanged. All original M4.5 scope and historical implementation evidence
retained; revision describes published baseline plus accepted uncommitted source
and local reconciliation, without inventing a publication SHA. Exactly the three
evidence/checkpoint files change this session; the whole local M4.5 set remains eleven.

Limitations remain explicit: no production caller authentication, tenant ownership,
payment-provider proof, entitlement authority, operational workflow/status route,
real provider/cloud integration, event idempotency, durable persistence/queue,
scheduler, automatic retry/resume/recovery, distributed exactly-once or production
rate limiting. A trusted local harness owns requests/capability; attaching that
capability to operational traffic would delegate it to clients. Injected Go code
is trusted composition, not sandboxed. EventID=CONFIRMED is correlation only.

Acceptance is local/uncommitted. HEAD stays 4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2;
index empty. No source/test/API edits, dependencies, staging, commit, push, deployment
or M4.6 activation. M4.6-M4.10 pending. STOP for human publication review. Fresh
AFTER-status verification/final preservation and both review patches are recorded
separately below and externally. At next startup reconcile actual HEAD/tree/index,
tracker and evidence; do not equate local acceptance with GitHub publication.


Final AFTER-status M4.5 acceptance evidence (2026-11-14): successful runs each
exit 0: focused creation (4 functions/21 RUN entries), invocation (18/56), HTTP
(9/103); all application (22/65) and HTTP (28/361) functions/RUN entries including
parents/subtests; focused races; full Go unit/race across four packages; go vet;
60 Python harness tests; tracker validation; canonical verifier all six checks;
both whitespace checks; read-only Go formatting/additive-method/new-file whitespace;
fresh route/trust and independent preservation audits. Post-status cache-permission
failures exited 1; approved retries exited 0. No behavioral failure or source fix.
Eight frozen implementation/test/API-contract hashes unchanged; only three acceptance
session files differ; total eleven local paths, other nineteen task objects, seven
original criteria and historical evidence preserved. M4.5 verified, active_task_id
null, index empty, main HEAD unchanged at 4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2.
Only this factual evidence and tracker commands/result were appended after regressions;
final tracker/whitespace/preservation and acceptance-only/complete patch reverse checks
are recorded outside the repository in /private/tmp/p1-m4.5-acceptance-review/.
No source/test/API contract/dependency/formatting edits, stage, commit, push, deployment
or M4.6 activation. Local acceptance remains uncommitted; STOP for human publication
review and separate publication authorization.


## M4.6 bounded implementation checkpoint (2026-11-14)

Previous M4.5 acceptance checkpoint recorded revision:
`4329b18d5dcb4928dc0d7ec72bc2babaa2f0f2a2` (then uncommitted acceptance).
Current observed published starting HEAD and unchanged current HEAD:
`ad792a3618071d6ae6ff7631b9ae622906a8d6de`, main, commit establishing verified M4.5.
Older uncommitted/publication-pending statements above remain historical. Fresh
Git/source/tracker reconciliation confirms M3.1-M3.10 and M4.1-M4.5 verified,
M4.6-M4.10 initially pending, active null, schema/WIP 1, twenty tasks, clean tree
and empty index. Fresh live origin/main matches the full required SHA. Initial
lookup exited 128 (sandbox DNS); approved retry exited 0. Inspected read-only
init.sh before running; initialization, tracker, M4.6 preflight and canonical
baseline verifier passed (60 Python tests, four Go packages). Cached-remote warning
was resolved by the live check. README milestone prose still lags actual source.
Initial gate outputs are in the session transcript; a second baseline canonical
run and subsequent command outputs are saved externally.

The deliberately submitted human-approved handoff authorizes bounded activation
and implementation, not acceptance/publication. All baseline gates passed before
activating M4.6 alone; tracker validation immediately passed. The six original
criteria were frozen unchanged in an external source/implementation/test/check/
limitation/review matrix before coding. Other nineteen full task objects, schema,
WIP, count, objective and dependency remain unchanged. M4.6 is in_progress and
active_task_id M4.6; evidence.revision/reviewed_by remain null. M4.7-M4.10 pending.

Exactly eight authorized paths change:
- internal/application/status.go: independent nonzero-sized private status token,
  sealed commands and local scope; validated copied exact-ID allowlist, duplicates
  collapsed and empty deny-all; narrow public construction, no operational calls
  during construction; private policy injection for faults. Service independently
  checks capability/context/syntax/policy before exactly one allowed Get. Denied
  IDs never read. Get-only interface, no lifecycle/action dependency, no repairs.
  Typed errors retain internal causes; result contains ID/status/actual version.
  Separate snapshot validator accepts ownerless TryStart RUNNING, READY diagnostic
  markers, owner-aware claim/Complete/Fail, SaveIfVersion completion, higher stored
  versions and coherent historical suffixes. History/state/version/owner/reason
  contradictions and history-free active/terminal ambiguity are safely rejected.
- internal/application/status_test.go: twelve functions, 62 RUN entries including
  parents/subtests. Constructor nil/typed-nil and no-call checks, zero/foreign/JSON
  authority, copied/empty/duplicate allowlists, direct-call ordering, private policy
  errors/unknown decisions, four states, legacy compatibility and malformed
  snapshots. Counts prove one allowed read, zero concealed reads/writes/callbacks.
  Repeated queries preserve every authoritative field. Channel barriers cover
  claim/Complete/Fail between snapshot capture and response, simultaneous readers
  and writer, blocked synchronous Get cancellation; no sleeps. A real registry
  action counter proves status does not replay a callback or advance a WAIT.
- internal/httpapi/status.go: unregistered handler, unchanged M4.3 grammar and
  M4.2 DTO. Explicit three-field projection; all emitted responses no-store.
  Fixed safe error envelopes; missing/concealed share exact 404 body/headers.
  FAILED retrieval is 200. Normal public composition accepts only sealed scope.
- internal/httpapi/status_test.go: ten functions, 43 RUN entries. Path boundaries,
  direct matched invalid IDs, no double decoding, default unmatched segment,
  exact context/command delegation, safe errors, real missing/foreign capability,
  exact response keys, no disclosure, equality of missing/concealed including all
  response headers, corruption/cancellation and invalid-projection guards. Real
  isolated M4.5 HTTP invocation -> same repository -> independent M4.6 HTTP status
  covers controlled READY and local.success/failure/wait/condition. Repeated GETs
  preserve owner/version/error/attempts/history; handoffs stay RUNNING/v2.
- docs/API_CONTRACTS.md: additive local-only status contract and limits.
- feature_list.json: M4.6 activation and factual evidence only.
- PROGRESS.md and docs/HANDOFF.md: append this factual checkpoint, preserve history.

Application tests first exited 1 from genuinely missing status symbols (structural
RED), then GREEN. HTTP first attempt exited 1 from Go-cache denial; approved retry
exited 1 from missing handler symbols (structural RED), then GREEN. No behavioral
RED is claimed; environmental failures are distinct. No unrelated source broken.
Successful commands each exited 0: focused application and HTTP verbose tests;
package races; go test -count=1 ./...; go test -race -count=1 ./...; go vet ./...;
PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts -p 'test_*.py' -v
(actual 60); tracker validation; bash scripts/verify.sh (all six); gofmt -l on new
Go files; route, trust and preservation audits. Go cache denials in focused/race/
full/vet/canonical attempts were followed by approved successful retries, without
source workarounds. Exact argv, exit codes and raw outputs are retained under
/private/tmp/p1-m4.6-review/commands.jsonl and named logs. No dependency install.

External route-overlay test freshly instantiates actual NewRouter: exactly GET
/health, unchanged 200/JSON/content type. GET/POST execution creation/status,
test-only invoke/status, enrollment, payment/callback, entitlement and admin paths
retain default 404. No listener or permanent future-route prohibition test added.
Static trust/order/Get-only audit complements manual source inspection and tests;
it is not a production security proof. Preservation checks all protected tracked
bytes, including repository methods, runner/actions, old source/tests, invocation,
DTO/validator/router/cmd, dependencies and harness. Other nineteen task objects and
all six criteria are unchanged; checkpoint/contract history is append-only.

Limits: fixture visibility is NOT production authentication, tenant ownership,
payment/entitlement authority or an in-process sandbox. Injected readers are
trusted and must supply stable independent snapshots; concurrent use needs a safe
reader. ID possession/context identity cannot authorize visibility. Concealed does
not prove existence; no constant-time promise. Legacy IDs outside the grammar
need not be addressable. History-free active/terminal records may be ambiguous.
No database error return is invented for synchronous Get. Cancellation cannot
interrupt it; observed cancellation suppresses success, but the final check can
race response delivery. Earlier coherent snapshots can legitimately be returned
after another goroutine completes/fails; no reread/polling freshness guarantee.
Explicit panic containment omitted: application panics propagate; composed Gin
recovery can return bodyless 500/no-store before success JSON, with no JSON error
or safe framework-log guarantee. Recovery tested without production wiring.
No operational route, new middleware, real provider/cloud, durable storage/queue,
scheduler, automatic retry/resume/recovery or distributed exactly-once behavior.

All work remains unstaged/uncommitted; index empty; HEAD unchanged. No acceptance,
M4.6 verification, active-task clearing, M4.7 activation, staging, commit, push or
deployment. External complete tracked-plus-untracked patch, source report, six-
criterion matrix, RED/GREEN logs, verification and five focused audits, route and
preservation reports, limitations and checklist are prepared for human source
review. Final post-documentation canonical/whitespace/preservation and patch checks
are recorded externally; do not infer results until those logs show success.
STOP before human acceptance/publication. At next startup recheck actual HEAD,
tree/index/tracker and init.sh; this checkpoint does not imply publication.


Final M4.6 evidence: post-documentation canonical retry exited 0, all six checks
including actual 60 Python tests and full Go unit/race across four packages. Its
sandbox attempt exited 1 at race-cache access after unit tests passed; approved
retry passed without source changes. Final trust/preservation and complete-patch
reverse checks exited 0. Only this factual text and tracker command/result evidence
were appended afterward; final tracker/whitespace/untracked-file/preservation and
refreshed patch checks are external. Twelve application/ten HTTP functions (62/43
RUN entries) executed. M4.6 remains in_progress, active_task_id M4.6; acceptance
revision/reviewed_by null. Exactly eight unstaged/uncommitted paths, empty index,
main HEAD ad792a3618071d6ae6ff7631b9ae622906a8d6de unchanged. No acceptance,
publication or M4.7 activation. STOP for human source review.

External untracked-file audit correction: its first run exited 1 because it
incorrectly required git diff --no-index to return zero for differing files.
There were no whitespace diagnostics. Corrected external audit exited 0, with
independent trailing-whitespace and exact gofmt-byte checks for all four new files.
No project source/test change resulted; failure/retry logs are retained externally.


## M4.6 acceptance-reconciliation checkpoint (2026-11-14)

Previous implementation checkpoint recorded revision: ad792a3618071d6ae6ff7631b9ae622906a8d6de plus eight uncommitted paths. Current observed starting HEAD remains that exact published M4.5 SHA; no future revision is invented. Earlier in-progress/source-review-pending statements remain historical and are superseded by this local acceptance checkpoint.

M4.6 acceptance reconciliation (2026-11-14), local and uncommitted. Published starting/current HEAD ad792a3618071d6ae6ff7631b9ae622906a8d6de on main; live origin/main matched after sandbox DNS exit128 and approved retry exit0. Human engineering owner reports ChatGPT source-level review accepted the submitted eight-file patch without corrections and authorized acceptance reconciliation. Codex independently compared bytes, inspected source/tests and executed fresh local checks; ChatGPT is not claimed to have independently run these commands, nor the human to have manually executed them. Exact reviewed patch SHA256 dac02892e77b9f1c4cde09bea6da4e3765e1def3b9b1431c5974c8d82fa0feae matched reconstructed tracked-plus-untracked bytes and all eight reconstructed file outputs. Historical manifest all88 entries matched. Five source/test/API files frozen and unchanged (hashes below and external Starting-hashes.json). Three acceptance files saved byte-for-byte externally before edits; historical evidence preserved.
Fresh pre-acceptance commands all passed on valid runs: focused application12 functions/62 RUN entries, focused HTTP10/43; focused races; full application34/127 and HTTP38/404 functions/RUN entries; package races; full Go units/races across four packages; go vet; actual60 Python harness tests; tracker; canonical all six; worktree/index whitespace; each new Go file exact gofmt bytes and independent whitespace. Route overlay first exited1 from sandbox Go-cache denial, approved retry exited0; only GET /health, unchanged JSON/type/status, all probed execution/test/enrollment/payment/entitlement/admin paths default404. No source workaround or behavioral failure. External historical-manifest parser first exited1 on a heading; corrected parser verified all88 entries, no repository changes.
Fresh source/test and reproducible audits confirm independent private nonzero-sized status authority, copied exact-ID allowlist, zero/foreign/JSON fail-closed, no public arbitrary policy setter, no constructor operations, independent direct-service gates. Policy precedes sole Get; concealed zero reads, visible missing one, identical public404 body/code/message/all headers, no constant-time claim. Strict Get-only query cannot call lifecycle/actions or change owner/version/diagnostic/attempt/history. Four-state explicit three-field projection uses actual stored version; FAILED lookup200 preserves private diagnostic. Legacy TryStart/READY marker/claim/Complete/Fail/SaveIfVersion/high versions/coherent suffixes supported, ambiguous history-free active/terminal rejected. Repeated real M4.5 WAIT/condition status remains RUNNING2 without cursor disclosure, callbacks/resume/terminal writes. Channel-barrier claim/Complete/Fail snapshot interleavings, 24 simultaneous readers and writer, cancellation-before/after-policy and blocked-Get evidence pass. Earlier coherent snapshot may follow concurrent completion; cancellation cannot interrupt synchronous Get or guarantee prevention after final check. Causes retained internally; fixed safe HTTP errors/no-store. Application panics propagate; composed Gin recovery has no JSON-envelope/safe-log guarantee.
All six original criteria satisfied within the approved local-only design by reviewed source, fresh tests and external Six-criterion-acceptance.txt. Protected repository/runner/invocation/DTO/validator/router/server/dependencies/harness and other19 task objects unchanged; original criteria/order/objective/M4.5 dependency preserved,20 tasks/schema1/WIP1. Only after pre-acceptance gates passed: M4.6 in_progress -> verified; active_task_id -> null; M4.7-M4.10 pending. No source correction, formatting, installation, stage/commit/push/deployment or next activation. Local fixture authority is NOT production authentication, tenant ownership, payment/entitlement authority or malicious in-process isolation. No operational status endpoint, durable storage/queue, scheduler, retry/resume/recovery, real provider/cloud or distributed guarantee. External artifacts /private/tmp/p1-m4.6-acceptance-review/; exact argv/exits/raw logs commands.jsonl. Post-acceptance results are recorded separately after execution. STOP for separate human publication review; local acceptance does not publish Git changes.
Frozen SHA256:
3814a7a0f2282b940ac5a04d5017e8d27e0b0cfc1de80baec5c5ea0b87ba745d  internal/application/status.go
67a4fe8e47cf40b6ac41c0288b3b1eb273cf7729e4c308e172c1b3003e402e88  internal/application/status_test.go
c95dd107ef4e9ec81a676e54d46b7a832df67fc5c7fd0bb4a56a220c8c3e395e  internal/httpapi/status.go
9dd2112d90387886b1ee65c91cfebb67920750059da42002a8dd482e03b1a7aa  internal/httpapi/status_test.go
d14fcff10b7fb21b4628e34fc268c572237b54966adf21eabec96dff0ea3f4bb  docs/API_CONTRACTS.md

Complete local eight-file set: internal/application/status.go; internal/application/status_test.go; internal/httpapi/status.go; internal/httpapi/status_test.go; docs/API_CONTRACTS.md; feature_list.json; PROGRESS.md; docs/HANDOFF.md. Only the final three changed during acceptance; five frozen files unchanged. Git index empty. Reconcile actual HEAD/tree/tracker and external post-acceptance evidence on next startup. Separate publication authorization required.


Final post-acceptance M4.6 evidence (2026-11-14): Codex fresh focused application12 functions/62 RUN entries and HTTP10/43 passed; application/HTTP package races, full Go unit/race across four packages, go vet, actual60 Python tests, tracker, canonical all six, both whitespace checks all exit0 on valid post-retry runs. Initial post HTTP package race exited1 from sandbox Go-cache access; approved full post-retry sequence exited0, no source workaround or behavioral failure. Post read-only formatting/new-file whitespace/trust audits passed; all five frozen hashes unchanged. Final preservation confirms only three acceptance-session edits, exact eight-file complete scope, protected baseline bytes and other19 task objects unchanged, M4.6 verified/active null, M4.7-M4.10 pending, empty index and unchanged main HEAD ad792a3618071d6ae6ff7631b9ae622906a8d6de. Acceptance-only and complete tracked-plus-untracked patches reverse-check exit0. This factual results append follows those runs; final validation/canonical/preservation and refreshed patch checks are recorded externally after this append. No source/test/API edit, staging/commit/push/deployment or next activation. STOP for human publication review; acceptance remains local/uncommitted.
