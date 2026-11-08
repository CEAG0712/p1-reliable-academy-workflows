P1 — Reliable Academy Workflows
Engineering Handoff and Implementation Roadmap
Date: November 8, 2026
Repository: https://github.com/CEAG0712/p1-reliable-academy-workflows
Default branch: main
Go module: example.com/p1-reliable-academy-workflows
Go version declared: 1.24.5
Current checkpoint: Stage 3 — M2 accepted, Exercise 82 complete; M3 not verified
Purpose: Enable a new AI coding agent to continue implementation without relying on prior conversation history.
1. Executive Summary
P1 is a Go-based, reusable Workflow Automation Platform.

Its purpose is to automate business workflows reliably through configurable triggers, conditions, branching, actions, waits and integrations.

Student enrollment and course-access provisioning is the first business use case, not the architecture of the core engine.

The engine must remain sufficiently generic to support other processes such as CRM updates, customer communications, content delivery and operational workflows.

The project is an independent engineering portfolio and simulated engagement. It is not a verified production deployment.
Current engineering state
Completed:
- Enrollment eligibility and payment validation.
- Tenant-aware entitlement identity.
- In-memory entitlement provider.
- Locally idempotent entitlement grants.
- Concurrency-safe workflow claims.
- Repository snapshot isolation.
- Create-only execution storage.
- Optimistic concurrency control.
- Owner-checked execution completion and failure.
- Execution versioning, attempt counts and transition history.
- Failure-reason consistency correction.
Not completed:
- Reusable workflow-definition and execution engine.
- Full enrollment workflow orchestration.
- Gin REST API.
- DynamoDB-backed execution state.
- AWS Lambda execution.
- SQS/SNS event integration.
- Terraform deployment.
- Distributed recovery and reconciliation.
- LangChain integration.
- Production observability, CI/CD and load testing.
- End-to-end deployment and portfolio handoff.
Do not report planned capabilities as implemented.
2. Sources of Truth and Agent Operating Instructions
Read these repository files before modifying the code:
1. AGENTS.md
2. PROGRESS.md
3. README.md
4. go.mod
5. workflow.go
6. workflow_repository.go
7. workflow_test.go
8. workflow_repository_test.go
9. enrollment.go
10. enrollment_test.go
The source code and tests determine implemented behavior. Documentation and prior assistant instructions are supporting evidence, not substitutes for code inspection.
Existing Harness Engineering files
AGENTS.md is implemented and committed. It instructs agents to inspect current code before changes, work on bounded features, preserve concurrency and lifecycle invariants, avoid unapproved deployment or repository mutations, execute relevant verification commands and report actual outcomes.

PROGRESS.md is implemented and committed. It records milestone status, verification evidence, blockers and next actions.

It currently identifies M3 as the next milestone but explicitly states that local M3 implementation has not been reconciled with GitHub.
Planned Harness Engineering files
The repository identifies these as future harness additions:
- init.sh
- feature_list.json
- scripts/verify.sh
These are planned developer-harness artifacts, not currently verified implementations.

Their purpose is to make agent initialization, feature tracking and verification repeatable. Their exact interfaces and schemas have not yet been established; do not invent them before defining and accepting their requirements.
Agent operating constraints
- Do not silently redesign the architecture.
- Preserve all existing tests unless a contract is deliberately revised.
- Make one bounded change with explicit acceptance criteria at a time.
- Run RED → GREEN → regression verification.
- Never claim tests passed unless actually executed.
- Do not push, deploy or modify production resources without explicit approval.
- Do not expose secrets or credentials.
- Record changed files, verification evidence, blockers and next actions after each implementation increment.
- Reconcile local and remote Git states before editing.
3. Current Implementation — M1
M1 — Execution Versioning & Optimistic Concurrency
Status: Completed and accepted.

M1 established repository-controlled state changes and rejection of stale workflow updates.
3.1 Implemented architecture
The central implementation consists of two components:

WorkflowExecution

Represents an execution and its current state.

Important fields:
- WorkflowID
- Status
- ClaimOwner
- Version
- LastError
- Attempts
- History
- Internal sync.Mutex
InMemoryWorkflowRepository

Stores workflow executions in a map protected by sync.RWMutex.

Important methods:
- NewInMemoryWorkflowRepository()
- Save()
- Get()
- TryClaim()
- SaveIfVersion()
- Complete() — introduced in M2
- Fail() — introduced in M2
These are concrete in-memory operations, not yet an abstract repository interface or DynamoDB implementation.
3.2 Create-only execution storage
Save() creates a workflow execution only if its ID does not already exist.

An existing execution cannot be replaced through Save().

This protects a successfully claimed execution from being overwritten by a stale READY snapshot.

Current limitation: Save() returns no success/failure result. Callers cannot distinguish a successful insertion from a duplicate that was ignored.
3.3 Repository snapshot isolation
Get() returns a cloned execution rather than exposing the repository's stored pointer.

cloneWorkflowExecution() copies data fields and independently copies the history slice.

The mutex itself is not copied.

Invariant: Mutating a retrieved snapshot must not modify the authoritative stored execution.
3.4 Atomic claims
TryClaim(workflowID, workerID):
- Requires an existing READY execution.
- Requires a nonempty worker ID.
- Acquires the repository's exclusive mutex.
- Updates status to RUNNING.
- Records the claim owner.
- Increments the execution version.
- Increments the attempt count.
- Appends a transition-history entry.
The check and mutation occur under one lock.

Invariant: Exactly one competing repository claim may succeed for a given READY execution.

This guarantee applies within a single repository instance and process.
3.5 Versioning
Initial execution version:

Version = 1

Successful claim:

READY v1 → RUNNING v2

Successful terminal transition:

RUNNING v2 → COMPLETED v3

or:

RUNNING v2 → FAILED v3

Every accepted state mutation increments the version once.

Rejected operations must leave the existing version unchanged.
3.6 Optimistic concurrency
SaveIfVersion(execution, expectedVersion) compares the submitted and expected versions against the authoritative stored version while holding the repository lock.

A stale version is rejected.

The implementation also requires:
- Existing execution.
- RUNNING source status.
- Nonempty recorded claim owner.
- Unchanged claim-owner string.
- Permitted destination status.
After Exercise 82, SaveIfVersion accepts COMPLETED only. It cannot transition an execution to FAILED.

This restriction must be preserved.
3.7 Concurrency guarantees and limitations
Implemented:
- Single-winner in-memory claims.
- Atomic expected-version comparison and update.
- Rejection of stale versions.
- Protection from duplicate create-only saves.
- Local synchronization using Go mutexes.
Not implemented:
- Cross-process locking.
- Database transactions.
- Durable conditional writes.
- Lease expiration.
- Fencing tokens.
- Crash-safe external side effects.
- Authenticated service identities.
A race-detector PASS does not establish these missing distributed guarantees.
4. Current Implementation — M2
M2 — Workflow Lifecycle & State Machine
Status: Completed and accepted, including corrective Exercise 82.

M2 introduced owner-checked terminal transitions, failure diagnostics, attempts and execution history.
4.1 Valid execution lifecycle
The supported repository-controlled lifecycle is:

READY → RUNNING → COMPLETED

or:

READY → RUNNING → FAILED

COMPLETED and FAILED are terminal in the current implementation.

There is no implemented transition from FAILED back to READY and no implemented automatic retry cycle.
4.2 Authoritative completion
Repository method:

Complete(workflowID, workerID, expectedVersion) bool

The operation succeeds only when:
- The execution exists.
- The worker ID is nonempty.
- The execution is RUNNING.
- The supplied worker ID equals the stored claim owner.
- The expected version matches the stored version.
Successful completion:
- Sets status to COMPLETED.
- Clears LastError.
- Increments version.
- Appends a completion-history entry.
Rejected completion must leave state unchanged.
4.3 Authoritative failure transition
Repository method:

Fail(workflowID, workerID, expectedVersion, reason) bool

The operation succeeds only when:
- The execution exists.
- The worker ID is nonempty.
- The reason is nonempty.
- The execution is RUNNING.
- The worker ID matches the recorded owner.
- The expected version matches the stored version.
Successful failure:
- Sets status to FAILED.
- Stores LastError = reason.
- Increments version.
- Appends history containing the same failure reason.
Critical invariant: Every accepted transition to FAILED must preserve a nonempty failure reason in both current state and transition history.

This invariant was explicitly added and verified in Exercise 82.
4.4 Exercise 82 correction
The original M2 SaveIfVersion() accepted FAILED but cleared LastError.

That allowed FAILED executions without meaningful failure diagnostics.

The corrective decision was:
- Fail() is the authoritative method for FAILED transitions.
- SaveIfVersion() must reject FAILED.
- SaveIfVersion() may still produce COMPLETED for M1 compatibility.
- The failure reason must be recorded by Fail() in both state and history.
Relevant regression tests:

TestRepositorySaveIfVersionCannotBypassFail

TestRepositoryFailPreservesReasonInStateAndHistory

Both were reported passing and are now present on GitHub.

Do not remove or weaken these tests.
4.5 Attempt tracking
Attempts increases when an execution is successfully claimed for processing.

Rejected claims do not increment attempts.

The current lifecycle provides no automatic reactivation after terminal failure, so the attempt counter does not yet constitute a full retry-attempt model.
4.6 Execution history
WorkflowTransition contains:
- From
- To
- WorkerID
- Version
- Reason
A successful claim records READY → RUNNING.

Successful completion records RUNNING → COMPLETED.

Failure records RUNNING → FAILED, including the reason.

History is stored in memory and copied when returning snapshots.

The repository protects recorded history rather than trusting arbitrary caller-supplied history.

Current history does not contain timestamps or correlation IDs. Those are future observability considerations, not implemented fields.
4.7 M2 design limitations
Worker ownership is not authentication. A string comparison establishes consistency with the recorded owner, not proof of caller identity.

TryStart() remains a legacy local method capable of starting an execution without a claim owner. It exists for compatibility with earlier tests and is not intended as the M3 runner's repository lifecycle path.

SaveIfVersion() still permits successful completion without accepting an explicit caller worker ID. M3 should prefer Complete().

There is no durable history store or immutable audit log.
5. Non-Negotiable Invariants
The next agent must preserve these requirements across every milestone.
ID	Invariant
INV-01	An execution has an initial version of 1
INV-02	Exactly one concurrent claim may succeed for a READY execution
INV-03	Successful state changes advance the version
INV-04	Rejected transitions leave state and version unchanged
INV-05	Stale snapshots cannot overwrite current repository state through Save()
INV-06	Get() returns independent execution and history snapshots
INV-07	Terminal transitions validate the authoritative stored state
INV-08	Complete() and Fail() require matching owner and expected version
INV-09	SaveIfVersion() cannot transition an execution to FAILED
INV-10	Fail() requires and preserves a nonempty failure reason
INV-11	Completed and failed executions cannot be reopened through current lifecycle methods
INV-12	Successful claims and terminal transitions are represented in execution history
INV-13	Enrollment entitlement identity includes TenantID, StudentID and CourseOfferingID
INV-14	Duplicate local entitlement grants must not create duplicate active entitlements
INV-15	The generic engine must not hardcode enrollment business logic
INV-16	AI output must not independently authorize payment or course access


Additional operational requirements from project discovery:
- Correct tenant/student/course access takes priority over throughput.
- The source-of-truth provider must be consulted before assuming external side-effect success where reconciliation is needed.
- Unsafe or ambiguous recovery requires containment and human escalation.
- External side effects must not be retried blindly after uncertain completion.
These last requirements are established design constraints, not claims that their production implementations already exist.
6. Stage 3 — Remaining Milestones
M3 — Reusable Workflow Definition & Execution Engine
Status: Pending implementation verification.
Objective
Build a generic Go workflow engine that executes configurable workflow definitions rather than enrollment-specific procedural code.
Scope
- Workflow definitions.
- Event triggers.
- Conditions.
- True/false branches.
- Registered action handlers.
- Wait/delay steps.
- Execution runner.
- Integration with M2 lifecycle operations.
- Two distinct example workflows.
Technology
Go, standard-library collections, interfaces/function handlers, context, time, synchronization and the existing in-memory repository.
Dependencies
M1 and M2 repository behavior and Exercise 82's failure-transition contract.
Expected deliverables
- Definition schema or Go models.
- Definition validation.
- Generic runner.
- Action and condition registration.
- Branch traversal.
- Nonblocking wait representation.
- In-memory execution progress required for resumption.
- Tests demonstrating reuse across multiple workflow definitions.
Acceptance criteria
- Malformed definitions are rejected.
- Unknown steps and handlers are detected.
- Triggers activate only matching workflows.
- Conditions choose the correct branch.
- Two distinct workflows execute using the same engine.
- Registered actions operate without changing engine core.
- Waits do not require sleeping workers.
- Resume does not repeat earlier completed steps in the tested in-memory scenario.
- Business action failures use M2 Fail() and retain diagnostic reasons.
- Successful execution uses M2 Complete().
- Concurrent runner invocation is tested.
- Existing M1/M2 tests remain passing.
Explicit limitation
M3's initial wait/resume capability is in-memory. Durable scheduling and crash recovery are not yet implemented.

The M3 implementation package previously prepared in the conversation has not been accepted as completed work. GitHub does not currently contain its engine files.
M4 — Enrollment Automation & Business Action Integration
Status: Pending.
Objective
Execute the first complete enrollment automation through the generic M3 engine.
Scope
- Connect incoming payment-event data to an enrollment workflow.
- Validate confirmed payment.
- Validate required tenant, student and course-offering fields.
- Evaluate enrollment eligibility through registered conditions.
- Grant course entitlements through a business action adapter.
- Verify correct entitlement identity.
- Preserve duplicate-event safety.
- Return correct completion or failure outcomes.
- Provide extension points for CRM updates, notifications and later AI steps.
Technology
Go, existing PaymentEvent model, EntitlementKey, FakeEnrollmentProvider, generic M3 engine and existing repository.
Dependencies
M3 runner, conditions, branches and action handlers.
Expected deliverables
- Enrollment workflow definition.
- Enrollment-specific condition handlers.
- Entitlement-grant action adapter.
- Adapter-facing integration tests.
- Negative tests for wrong tenant, wrong student and invalid payment.
- Demonstration of payment-to-course-access execution.
Acceptance criteria
- Confirmed, valid eligible payments can drive the intended workflow.
- Pending or invalid payments cannot grant access.
- Tenant/student/offering identity remains correct.
- Duplicate processing does not create duplicate local entitlements.
- A wrong-student entitlement must never be considered successful.
- The workflow records success or failure through M2 methods.
- The generic engine has no embedded enrollment-specific conditionals.
- The implementation does not claim real provider integration or production savings.
The simulated business SLO target is correct course access for 99% of eligible purchases within two minutes. That is a later end-to-end acceptance target, not something established by the existing local tests.
M5 — Integration Testing, Reliability Verification & Refactoring
Status: Pending.
Objective
Finish Stage 3 with a verified, maintainable Go workflow core ready for API and infrastructure integration.
Scope
- Full workflow happy-path tests.
- Negative and malformed-input tests.
- Competing-worker tests.
- Version-conflict tests.
- Action-failure tests.
- Branch and wait/resume tests.
- Enrollment isolation and duplicate-processing tests.
- Regression verification.
- Bounded refactoring.
- Technical documentation.
Technology
Go testing, goroutines, channels, sync primitives, race detector and Go vet.
Dependencies
M1 through M4.
Expected deliverables
- Comprehensive Stage 3 regression suite.
- Documented workflow contracts and limitations.
- Reproducible verification commands.
- Validated integration between workflow engine and enrollment adapters.
- Stage 3 acceptance report.
Acceptance criteria
- All implemented paths satisfy their tests.
- Existing regressions remain passing.
- Concurrency tests pass under repeated execution.
- Race detector reports no races in exercised code.
- Go vet succeeds.
- Git diff integrity checks succeed.
- No known failure-reason regression remains.
- The engine can run more than one workflow type.
- Remaining cloud and durability gaps are explicitly documented.
M5 does not imply that production disaster recovery, distributed exactly-once processing or end-to-end cloud deployment has already been achieved.
7. Stage 4 — Gin API
Status: Planned.
Objective
Expose the workflow platform through a Go HTTP API.
Scope
- Workflow invocation endpoints.
- Request validation.
- Execution creation/triggering.
- Status retrieval.
- Explicit error responses.
- Integration with the generic workflow runner.
Technology
Go and Gin.
Dependencies
Accepted Stage 3 engine and business adapters.
Expected deliverables
- HTTP routing layer.
- Request/response contracts.
- Handlers invoking workflow services.
- HTTP-level tests.
- API documentation.
Acceptance criteria
- Valid requests trigger the intended workflow.
- Invalid requests are rejected before business side effects.
- API handlers do not bypass repository lifecycle rules.
- Status and errors are represented consistently.
- Duplicate or conflicting requests are handled according to defined contracts.
- Integration tests verify handler-to-engine behavior.
Specific endpoint paths, status-code matrices and authentication design have not yet been finalized. The next agent must not fabricate those as previously approved contracts.
8. Stage 5 — AWS Infrastructure
Status: Planned.
Objective
Move from a process-local prototype toward durable, event-driven execution infrastructure.
Scope
- AWS Lambda workers.
- DynamoDB workflow state.
- DynamoDB conditional writes.
- Amazon SQS work queues.
- Amazon SNS event notifications.
- Terraform infrastructure definitions.
- Necessary access permissions and deployment configuration.
Technology
Go, AWS Lambda, DynamoDB, SQS, SNS and Terraform.
Dependencies
Stage 3 engine, Stage 4 API and defined persistence/event interfaces.
Expected deliverables
- Infrastructure-as-code.
- Durable execution repository adapter.
- Queue/worker integration.
- Event publication/consumption paths.
- Deployment configuration.
- Integration tests against appropriate AWS-compatible test infrastructure or isolated test resources.
Acceptance criteria
- Durable repository operations preserve expected-version semantics.
- Competing workers cannot claim the same execution simultaneously through database conditions.
- Queue delivery does not bypass idempotency rules.
- Event processing respects tenant identity.
- Infrastructure can be recreated through Terraform.
- Least-privilege permissions and secret handling are reviewed.
- Failed infrastructure operations produce actionable diagnostics.
Distributed correctness must be verified at the database boundary; Go's in-memory mutex does not provide cross-Lambda synchronization.
9. Stage 6 — Reliability, Recovery and Reconciliation
Status: Planned.
Objective
Handle failure scenarios that cannot be solved through simple local mutexes or retries.
Scope
- Bounded retries.
- Backoff and transient-error handling.
- Lease-based execution ownership.
- Fencing or equivalent stale-worker protection.
- Queue redelivery and duplicate processing.
- Reconciliation after uncertain external side effects.
- Backpressure and provider rate limits.
- Failure containment.
- Human escalation for unsafe recovery cases.
Technology
Go, DynamoDB conditional operations, SQS/Lambda integration, external provider adapters and operational logging.
Dependencies
Durable Stage 5 execution state and event delivery.
Expected deliverables
- Recovery state and transition rules.
- Retry and backoff policies.
- Lease and stale-worker tests.
- External side-effect reconciliation adapters.
- Duplicate-event recovery scenarios.
- Manual escalation path for ambiguous states.
- Fault-injection test scenarios.
Acceptance criteria
- Duplicate queue messages do not cause duplicate unauthorized grants.
- A stale worker cannot overwrite state after ownership has changed.
- A lease expiration alone is not treated as proof that an external side effect never occurred.
- Crash-after-success scenarios reconcile against the authoritative external provider before replaying an unsafe action.
- Retry budgets are bounded.
- Rate limiting and backpressure prevent uncontrolled provider overload.
- Wrong-student and cross-tenant access cases are contained and never silently accepted.
- Uncertain cases requiring human approval remain visible rather than being incorrectly marked successful.
The existing feasibility exercises used an illustrative provider-consistency window of up to 30 seconds and a two-minute course-access SLO. These are design inputs to test against, not validated production performance.
10. Stage 7 — LangChain AI Integration
Status: Planned.
Objective
Integrate AI-assisted actions into the workflow platform without delegating critical authorization decisions to an LLM.
Scope
The specifically defined first AI action is a localized welcome message after confirmed enrollment.

The workflow engine should call an AI action adapter rather than embedding prompts or model-specific behavior in core orchestration.
Technology
LangChain, the Go workflow platform and a supported service integration boundary.
Dependencies
M3 action registration, M4 enrollment workflow and stable lifecycle/error handling.
Expected deliverables
- AI action adapter contract.
- Localized welcome-message generation.
- Output validation.
- Error and timeout handling.
- Model integration tests using controlled outputs or mocks.
- Failure and fallback behavior.
Acceptance criteria
- AI runs only after the appropriate deterministic business checks.
- AI output does not authorize payment or course access.
- Generated content is validated before downstream use.
- Provider failures have bounded handling.
- Sensitive data is not unnecessarily passed to model services.
- AI action failure does not corrupt workflow or entitlement state.
- The adapter is replaceable without rewriting workflow orchestration.
No specific model provider, prompt template, model version or latency target has been formally selected.
11. Stage 8 — Observability, CI/CD and Performance
Status: Planned.
Objective
Make platform behavior measurable, reproducible and continuously verifiable.
Scope
- Structured execution logs.
- Metrics for claims, latency, failures, retries and queue behavior.
- Execution identifiers and correlation across components.
- Distributed tracing where infrastructure supports it.
- Automated tests in CI.
- Deployment pipeline.
- Load and concurrency testing.
- SLO reporting.
Technology
Go testing, CI/CD tooling, logging/metrics/tracing facilities and AWS telemetry as selected during implementation.
Dependencies
Running API, workflow engine, AWS infrastructure and business integrations.
Expected deliverables
- Structured logging and telemetry integration.
- Automated regression checks.
- CI/CD pipeline configuration.
- Load-test scenarios and results.
- SLO dashboards or documented measurements.
- Operational troubleshooting documentation.
Acceptance criteria
- Each execution can be traced across its relevant components.
- Errors identify the failing operation without exposing secrets.
- CI runs required unit, integration and static-analysis checks.
- Load tests exercise the planned burst and steady-state workloads.
- Workflow latency and correctness are measurable.
- SLO compliance is based on measurements, not proposed business figures.
- Deployment failures are observable and actionable.
The earlier project used 10,000 workflow tasks/month and illustrative throughput/burst scenarios as business and feasibility targets. These must not be represented as achieved results until measured.
12. Stage 9 — Deployment, Demonstration and Handoff
Status: Planned.
Objective
Complete P1 as a deployable, demonstrable engineering portfolio project.
Scope
- Deployment and verification.
- Working demonstration.
- Architectural documentation.
- Operational runbooks.
- Known limitations and tradeoffs.
- Handoff material.
- Evidence-based portfolio claims.
Technology
The completed Go/Gin/AWS/Terraform stack, plus the tooling actually selected and implemented during earlier stages.
Dependencies
Accepted preceding milestones and sufficient integration evidence.
Expected deliverables
- Deployment artifacts.
- Demonstration of complete enrollment automation.
- Repository documentation.
- Architecture diagram.
- Test and verification evidence.
- Incident and recovery runbooks.
- Portfolio-oriented technical case study.
- Engineering handoff.
Acceptance criteria
- The platform can be deployed through documented steps.
- A valid enrollment scenario runs end to end.
- Invalid and duplicate events are handled safely.
- Workflow state and diagnostics can be inspected.
- Recovery scenarios have tested or explicitly documented behavior.
- Operational metrics and known limitations are available.
- Claims of savings, error reductions and throughput are supported by real measurements or clearly labeled as simulated targets.
Proposed business targets include 10,000 tasks/month, 2,500 avoided manual errors/year and approximately $200,000/year savings. These are unverified project targets, not achieved production outcomes.
13. M3 — Sequential Engineering Task Plan
M3 should be implemented as a sequence of bounded, independently verifiable features.

Its previously planned scope is preserved. The following task breakdown is implementation sequencing, not an additional architectural redesign.
M3.1 — Workflow Definition Models
Implement reusable models representing definitions, triggers and steps.

Verify that the model can describe two different workflows without adding business-specific fields.

Acceptance: Models compile; tests demonstrate distinct workflow definitions.
M3.2 — Definition Validation
Implement checks for missing IDs, invalid start steps, unsupported step types, broken links, missing action/condition references and invalid branching.

For the initial engine, reject unsupported graph cycles rather than introducing unbounded loops.

Acceptance: Valid definitions are accepted; malformed definitions produce deterministic errors.
M3.3 — Trigger Matching
Define how an incoming event is matched to a registered workflow definition.

Verify that nonmatching triggers cannot create or launch unintended executions.

Acceptance: Matching events start eligible workflows; mismatches cause no execution side effects.
M3.4 — Handler Registration
Add registration of action handlers and condition predicates.

Keep the workflow engine independent of enrollment providers, CRM systems and AI vendors.

Acceptance: Multiple handler implementations can be registered and invoked by name; unknown handlers are rejected.
M3.5 — Conditions and Branching
Implement condition evaluation and selection of the configured true/false successor step.

Acceptance: Both branches are tested, including condition-error behavior.
M3.6 — Sequential Action Execution
Execute configured action steps in graph order and advance execution progress when an action succeeds.

Acceptance: Actions execute in the required order; a failed action stops the normal path and produces an error outcome.
M3.7 — Wait and Resume
Represent a delay step with a resume deadline instead of calling time.Sleep() for the full delay.

Retain the execution cursor and necessary context for an in-memory resume.

Acceptance: Early resume attempts do not advance incorrectly; an eligible resume continues from the correct step without replaying earlier actions in the tested scenario.

This does not yet require persistence across restarts.
M3.8 — M2 Lifecycle Integration
Connect runner execution to the repository's owner-checked lifecycle.

Use TryClaim() as the ownership acquisition operation and Complete()/`Fail()` for terminal outcomes.

Do not use SaveIfVersion() to mark executions FAILED.

Acceptance: Successful and failed workflows preserve M2 status, version and history invariants.
M3.9 — Concurrency and Failure Verification
Test competing execution attempts, invalid ownership, action errors and duplicate invocations.

Acceptance: No duplicate action execution in the tested in-process concurrent scenario; failure reasons and lifecycle state remain correct.

Distributed exactly-once side effects are not claimed.
M3.10 — Reusability and Milestone Acceptance
Run at least two genuinely different workflow definitions through the same runner.

Execute the entire Stage 3 regression suite and race checks.

Acceptance: Both workflows run without engine-core changes; all prior tests remain passing; limitations are documented.
M3 sequencing rule
For every task:
1. Record the contract and expected behavior.
2. Write a focused failing test.
3. Implement the minimum change.
4. Verify the focused test.
5. Run relevant regression tests.
6. Report files changed and commands run.
7. Update the implementation checkpoint only after verification.
The original roadmap grouped M3 into six work packages. This ten-task breakdown makes those work packages easier for a coding agent to implement incrementally without changing the milestone's scope.
14. Architecture Boundaries
Implemented today
- Go domain and execution code.
- In-memory execution repository.
- In-memory entitlement provider.
- Local synchronization and atomic claims.
- Version-controlled execution updates.
- Explicit completion and failure.
- Execution history.
- Enrollment eligibility helpers.
- Local idempotent entitlement grants.
- Automated Go tests.
- Agent operating instructions in AGENTS.md.
- Development status documentation in PROGRESS.md.
Planned within P1
- Generic workflow runner.
- Action and condition registration.
- Branching and delay processing.
- Enrollment workflow integration.
- Gin API.
- Lambda workers.
- DynamoDB execution state.
- SQS/SNS messaging.
- Terraform.
- Distributed retry and recovery.
- Provider reconciliation.
- LangChain actions.
- Observability.
- CI/CD.
- Load testing.
- Deployment and technical handoff.
Not currently established as P1 deliverables
- A general-purpose visual no-code workflow designer.
- A multi-tenant commercial SaaS administration portal.
- Billing and subscription management.
- A general-purpose autonomous multi-agent platform.
- Training or fine-tuning proprietary foundation models.
- A production CRM replacement.
- A complete learning-management system.
- Exactly-once distributed external side effects as an unconditional guarantee.
These are exclusions from the currently approved scope, not permanent prohibitions. They must not be introduced without explicit scope decisions.
Critical separation of responsibilities
Workflow core: Orchestration, graph traversal, conditions, waits, handler invocation and execution progress.

Execution repository: Authoritative execution state, ownership, versioning, lifecycle and history.

Business adapters: Enrollment, entitlement provision, CRM and notifications.

AI adapter: Validated AI-assisted output; not payment or authorization decisions.

Infrastructure: Persistence, event transport, deployment, identity and distributed coordination.

Business adapters must not bypass repository lifecycle controls, and the workflow engine must not duplicate provider-specific business logic.
15. Verification Strategy
15.1 Test pyramid
Unit tests: Individual state transitions, validators, definitions, conditions and handlers.

Component tests: Runner behavior with in-memory repository and fake adapters.

Integration tests: Gin handlers, DynamoDB conditional updates, queue delivery, Lambda processing and external adapter boundaries.

End-to-end tests: Trigger payment event, execute workflow, verify correct course entitlement, inspect final execution state and diagnostics.

Each layer proves a different property. Passing local tests does not establish cloud deployment correctness.
15.2 Concurrency verification
Preserve and extend:
- Single-winner claims.
- Stale version rejection.
- Concurrent versioned updates.
- Duplicate workflow invocation.
- Competing workers.
- Idempotent entitlement grants.
Run targeted concurrent tests repeatedly with go test -race.

For AWS-backed execution, verify concurrency through real conditional database operations rather than relying on process-local locks.
15.3 Failure recovery verification
Before marking Stage 6 complete, test:
- Provider timeout.
- Queue redelivery.
- Duplicate event.
- Worker crash before side effect.
- Worker crash after side effect but before state persistence.
- Stale ownership or fencing token.
- Expired lease with uncertain provider outcome.
- Provider consistency delay.
- Retry budget exhaustion.
- Manual escalation.
- Incorrect tenant/student access.
A timed-out call or expired lease must not automatically be interpreted as proof that the external action failed.

Reconcile against the authoritative provider before retrying unsafe side effects.
15.4 Business correctness
Use tenant-aware keys and verify that entitlements cannot cross tenant, student or course-offering boundaries.

The most important negative tests involve accidentally granting access to the wrong student.

AI-generated content must never override deterministic eligibility validation.
15.5 Workflow engine verification
M3 must test:
- Definition validation.
- Trigger matching.
- Branch choice.
- Action dispatch.
- Handler errors.
- Wait/resume behavior.
- Duplicate runner invocation.
- Ownership validation.
- Lifecycle completion/failure.
- Reusability across at least two workflows.
For initial in-memory waits, explicitly verify the absence of worker sleeping during the waiting period.
15.6 Operational and performance verification
Later stages should measure:
- Event receipt-to-processing latency.
- Queue wait time.
- Workflow completion latency.
- Correct entitlement provisioning.
- Duplicate-processing suppression.
- Failure rates.
- Retry counts.
- Provider rate-limit behavior.
- Reconciliation outcomes.
- Manual escalations.
The project SLO target is 99% of eligible purchases receiving correct course access within two minutes.

The target is not achieved until end-to-end measurements demonstrate compliance.
15.7 Standard verification commands
gofmt -w *.go
go test -count=1 -v ./...
go test -race -count=1 ./...
go test -race -count=100 -run 'Test(InMemoryWorkflowRepositoryConcurrent|Repository)' ./...
go vet ./...
git diff --check
git status --short
Adjust targeted test patterns when working on M3 or later code.

The final acceptance report must record the actual command, result, relevant failure output and implementation revision.
16. GitHub Repository State — November 8, 2026
Repository inspected:

https://github.com/CEAG0712/p1-reliable-academy-workflows

Default branch inspected: main.

Latest commit observed during inspection:

b9b4a013cd9436e943fedad1989b35469fa8ed43

The latest observed commit includes AGENTS.md and PROGRESS.md.
Confirmed present on GitHub
- README.md
- go.mod
- enrollment.go
- workflow.go
- workflow_repository.go
- workflow_repository_test.go
- AGENTS.md
- PROGRESS.md
The published workflow.go contains the M2 history and attempt fields.

The published workflow_repository.go contains Complete(), Fail(), SaveIfVersion() and the Exercise 82 restriction against FAILED transitions through SaveIfVersion.

The published repository test file contains the two Exercise 82 corrective tests.

Therefore, the repository is not merely at Exercise 60 or M1. It has advanced through the M2 correction.
Confirmed missing from inspected GitHub paths
- workflow_definition.go
- workflow_engine.go
The expected M3 source files were not found on the default branch.

The specific path .github/workflows/go.yml was also not found; this alone does not establish that no other CI workflow configuration exists.
Documentation discrepancy
README.md still states that Stage 3 Exercise 60 is complete and that conditional repository claims are the next milestone.

That section is outdated relative to the actual published implementation.

PROGRESS.md is more recent. It correctly identifies M3 as the current engineering milestone and requires reconciliation of local M3 work.

However, its last confirmed baseline reference, 2538574, is older than the latest observed repository commit.

A new coding agent must reconcile these documents with current source and Git history before updating them.
Local versus remote uncertainty
The user has reported successful local testing through M2 Exercise 82, and the corresponding implementation and tests are now visible on GitHub.

The user has also received an M3 implementation package previously in the conversation.

There is no evidence that this M3 package was implemented, copied into the working tree, tested locally or pushed to GitHub.

The local Git working tree is not directly accessible through the GitHub connector. Therefore, uncommitted local files and unpushed commits cannot be confirmed.

Do not assert that local and remote repositories are identical.
17. Required Local-to-GitHub Reconciliation
Before starting M3, run from the user's local project directory:
git status --short
git branch --show-current
git rev-parse HEAD
git fetch origin
git rev-parse origin/main
git log --oneline --decorate -5
git diff --name-status origin/main
git diff --stat origin/main
Also inspect whether these local files exist:
ls -l workflow_definition.go workflow_engine.go workflow_engine_test.go
A missing-file message is expected if M3 has not been started locally.

Compare:
- Local HEAD versus origin/main.
- Uncommitted changes.
- Untracked files.
- M3 files, if present.
- Local test outcomes.
- PROGRESS.md and README.md status claims.
Do not overwrite local files or reset the repository merely to match GitHub.

If M3 files exist locally, inspect and test them before applying another implementation package.
18. Recommended Next Agent Assignment
Assignment: Reconcile the Stage 3 baseline and implement the first bounded M3 feature.
Step A — Reconcile state
Read AGENTS.md and PROGRESS.md.

Inspect current source code, tests, Git state and any local M3 files.

Report differences without modifying code.
Step B — Confirm M3.1 contract
Define and test the minimal generic workflow-definition models necessary to represent triggers and steps.

Do not introduce enrollment-specific fields into the engine.
Step C — Implement and verify
Write the RED tests, implement the models, run targeted tests and preserve all M1/M2 regressions.

Do not proceed to the next M3 task until the current task has been verified.
Step D — Report
At task completion, provide:
- Objective.
- Files changed.
- Behavior implemented.
- Tests executed.
- Actual results.
- Known limitations.
- Local commit/working-tree state.
- Next bounded engineering task.
Do not push or deploy without explicit permission.
19. Acceptance Status Summary
Component	Status
M1 — Execution Versioning & Optimistic Concurrency	Accepted
M2 — Workflow Lifecycle & State Machine	Accepted
Exercise 82 — Failure-Reason Consistency	Accepted
M3 — Reusable Workflow Definition & Execution Engine	Pending
M4 — Enrollment Integration	Pending
M5 — Stage 3 Integration Verification	Pending
Stage 4 — Gin API	Planned
Stage 5 — AWS Infrastructure	Planned
Stage 6 — Reliability & Recovery	Planned
Stage 7 — LangChain AI	Planned
Stage 8 — Observability & CI/CD	Planned
Stage 9 — Deployment & Handoff	Planned


Final handoff rule: Preserve tested M1/M2 behavior; verify the repository and local working tree before implementing anything; keep the core engine generic; and do not present planned infrastructure or business targets as completed capabilities.

The next authorized engineering increment is M3.1, after baseline reconciliation.