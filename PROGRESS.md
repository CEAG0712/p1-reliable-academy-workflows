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
- Current M3 active task: M3.6, implementation committed; human acceptance pending.
  M3.1-M3.5 are verified; M3.7-M3.10 remain pending and unauthorized.

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

## M3.6 — Sequential Action Execution Checkpoint (2026-11-11)

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
