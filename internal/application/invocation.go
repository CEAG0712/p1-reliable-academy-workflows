package application

import (
	"context"
	"errors"
	"maps"
	"reflect"
	"strings"

	workflow "example.com/p1-reliable-academy-workflows"
)

// These contracts describe controlled local invocation only. They do not add
// caller authentication, tenant ownership, payment proof or event idempotency.
type TrustedLocalInvocationCommand struct {
	authority *localAuthority
	selection SubmissionCommand
}

type ExecutionIDSource interface {
	NextExecutionID(context.Context) (string, error)
}
type LifecycleRunner interface {
	RunActions(workflow.WorkflowLifecycleRepository, string, string, workflow.WorkflowDefinition) (string, error)
}

// Private plan fields prevent callers outside this package injecting callbacks.
// Resolvers must return stable configuration while the service snapshots it.
type ApprovedLocalPlan struct {
	definition workflow.WorkflowDefinition
	runner     LifecycleRunner
	actions    map[string]struct{}
	conditions map[string]struct{}
}
type ApprovedDefinitionResolver interface {
	ResolveApprovedDefinition(context.Context, string) (ApprovedLocalPlan, error)
}
type InvocationRepository interface {
	workflow.WorkflowLifecycleRepository
	CreateIfAbsent(*workflow.WorkflowExecution) (bool, error)
}

type CreationDisposition string

const (
	CreationNotAttempted CreationDisposition = "not_attempted"
	CreationInserted     CreationDisposition = "inserted"
	CreationDuplicate    CreationDisposition = "duplicate"
	CreationRejected     CreationDisposition = "rejected"
	CreationUncertain    CreationDisposition = "uncertain"
)

type ClaimDisposition string

const (
	ClaimNotObserved ClaimDisposition = "not_observed"
	ClaimAccepted    ClaimDisposition = "accepted"
	ClaimRejected    ClaimDisposition = "rejected"
)

type InvocationOutcomeKind string

const (
	OutcomeRejected          InvocationOutcomeKind = "rejected"
	OutcomeCreatedNotStarted InvocationOutcomeKind = "created_not_started"
	OutcomeCompleted         InvocationOutcomeKind = "completed"
	OutcomeFailed            InvocationOutcomeKind = "failed"
	OutcomeHandoff           InvocationOutcomeKind = "handoff"
	OutcomeClaimRejected     InvocationOutcomeKind = "claim_rejected"
	OutcomeReconciliation    InvocationOutcomeKind = "reconciliation_required"
)

// InvocationOutcome is internal evidence, never an HTTP DTO. A nonzero result
// with an error retains creation identity and the actual independently read state.
// Definition/event/cursor association exists only during this call/result.
type InvocationOutcome struct {
	ExecutionID          string
	WorkflowDefinitionID string
	EventID              string
	Creation             CreationDisposition
	RunAttempted         bool
	Claim                ClaimDisposition
	Kind                 InvocationOutcomeKind
	Cursor               string
	Snapshot             *workflow.WorkflowExecution
}
type InvocationErrorKind string

const (
	InvocationInvalidCommand     InvocationErrorKind = "invalid_command"
	InvocationMissingAuthority   InvocationErrorKind = "missing_authority"
	InvocationForbidden          InvocationErrorKind = "forbidden"
	InvocationDefinitionNotFound InvocationErrorKind = "definition_not_found"
	InvocationInvalidDefinition  InvocationErrorKind = "invalid_definition"
	InvocationUnavailable        InvocationErrorKind = "unavailable"
	InvocationDuplicate          InvocationErrorKind = "duplicate"
	InvocationExecutionFailed    InvocationErrorKind = "execution_failed"
	InvocationClaimRejected      InvocationErrorKind = "claim_rejected"
	InvocationReconciliation     InvocationErrorKind = "reconciliation_required"
)

type InvocationError struct {
	Kind  InvocationErrorKind
	Cause error
}

func (e *InvocationError) Error() string { return "local invocation: " + string(e.Kind) }
func (e *InvocationError) Unwrap() error { return e.Cause }
func invocationError(kind InvocationErrorKind, cause error) error {
	return &InvocationError{Kind: kind, Cause: cause}
}

type LocalWorkflowInvoker interface {
	InvokeLocal(context.Context, TrustedLocalInvocationCommand) (InvocationOutcome, error)
}
type invocationDependencies struct {
	authority   *localAuthority
	definitions ApprovedDefinitionResolver
	ids         ExecutionIDSource
	repo        InvocationRepository
	workerID    string
}
type invocationService struct{ deps invocationDependencies }

func nilInvocationDependency(v any) bool {
	if v == nil {
		return true
	}
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return r.IsNil()
	}
	return false
}
func validInvocationDependencies(d invocationDependencies) bool {
	return d.authority != nil && !nilInvocationDependency(d.definitions) && !nilInvocationDependency(d.ids) && !nilInvocationDependency(d.repo) && strings.TrimSpace(d.workerID) != ""
}

// No dependency operations are performed during construction. Fault injection
// remains private; the public scope always supplies the fixed harmless resolver.
func newInvocationService(d invocationDependencies) (*invocationService, error) {
	if !validInvocationDependencies(d) {
		return nil, invocationError(InvocationUnavailable, nil)
	}
	return &invocationService{deps: d}, nil
}

func (s *invocationService) InvokeLocal(ctx context.Context, c TrustedLocalInvocationCommand) (InvocationOutcome, error) {
	out := InvocationOutcome{Creation: CreationNotAttempted, Claim: ClaimNotObserved, Kind: OutcomeRejected}
	if s == nil || !validInvocationDependencies(s.deps) {
		return out, invocationError(InvocationUnavailable, nil)
	}
	if ctx == nil {
		return out, invocationError(InvocationInvalidCommand, nil)
	}
	if c.authority == nil {
		return out, invocationError(InvocationMissingAuthority, nil)
	}
	if c.authority != s.deps.authority {
		return out, invocationError(InvocationForbidden, nil)
	}
	if !validLocalIdentifier(c.selection.WorkflowDefinitionID) || !validLocalIdentifier(c.selection.EventID) {
		return out, invocationError(InvocationInvalidCommand, nil)
	}
	if err := ctx.Err(); err != nil {
		return out, invocationError(InvocationUnavailable, err)
	}
	if !approvedLocalDefinition(c.selection.WorkflowDefinitionID) {
		return out, invocationError(InvocationForbidden, nil)
	}
	plan, err := s.deps.definitions.ResolveApprovedDefinition(ctx, c.selection.WorkflowDefinitionID)
	if err != nil {
		kind := InvocationUnavailable
		if errors.Is(err, ErrDefinitionNotFound) {
			kind = InvocationDefinitionNotFound
		}
		return out, invocationError(kind, err)
	}
	// StepDefinition contains value fields; cloning the map gives a full private
	// graph snapshot. Trusted resolvers must not mutate it concurrently with copy.
	plan.definition.Steps = maps.Clone(plan.definition.Steps)
	if plan.definition.ID != c.selection.WorkflowDefinitionID {
		return out, invocationError(InvocationInvalidDefinition, nil)
	}
	if err := validateLocalPlan(plan); err != nil {
		return out, invocationError(InvocationInvalidDefinition, err)
	}
	id, err := s.deps.ids.NextExecutionID(ctx)
	if err != nil {
		return out, invocationError(InvocationUnavailable, err)
	}
	if !validExecutionID(id) {
		return out, invocationError(InvocationInvalidCommand, nil)
	}
	if err := ctx.Err(); err != nil {
		return out, invocationError(InvocationUnavailable, err)
	}
	out.ExecutionID, out.WorkflowDefinitionID, out.EventID = id, plan.definition.ID, c.selection.EventID
	inserted, err := s.deps.repo.CreateIfAbsent(workflow.NewWorkflowExecution(id))
	if err != nil {
		out.Creation = CreationRejected
		if inserted {
			out.Creation = CreationUncertain
			out.Kind = OutcomeReconciliation
			out.Snapshot, _ = s.deps.repo.Get(id)
			return out, invocationError(InvocationReconciliation, err)
		}
		return out, invocationError(InvocationUnavailable, err)
	}
	if !inserted {
		out.Creation = CreationDuplicate
		return out, invocationError(InvocationDuplicate, nil)
	}
	out.Creation = CreationInserted
	var exists bool
	out.Snapshot, exists = s.deps.repo.Get(id)
	if !exists || !freshInvocationSnapshot(out.Snapshot, id) {
		out.Kind = OutcomeReconciliation
		return out, invocationError(InvocationReconciliation, nil)
	}
	if err := ctx.Err(); err != nil {
		out.Kind = OutcomeCreatedNotStarted
		return out, invocationError(InvocationUnavailable, err)
	}
	observer := &lifecycleObserver{repo: s.deps.repo}
	out.RunAttempted = true
	out.Cursor, err = plan.runner.RunActions(observer, id, s.deps.workerID, plan.definition)
	out.Snapshot, exists = s.deps.repo.Get(id)
	if !exists {
		out.Snapshot = nil
	}
	return classifyInvocation(out, err, observer, plan.definition, s.deps.workerID)
}

func freshInvocationSnapshot(s *workflow.WorkflowExecution, id string) bool {
	return s != nil && s.WorkflowID == id && s.Status == "READY" && s.Version == 1 && s.ClaimOwner == "" && s.LastError == "" && s.Attempts == 0 && len(s.History) == 0
}
func invocationSnapshotMatches(s *workflow.WorkflowExecution, id, owner, status, reason string) bool {
	if s == nil || s.WorkflowID != id || s.Status != status || s.ClaimOwner != owner || s.LastError != reason || s.Attempts != 1 {
		return false
	}
	history := []workflow.WorkflowTransition{{From: "READY", To: "RUNNING", WorkerID: owner, Version: 2}}
	version := uint64(2)
	if status != "RUNNING" {
		version = 3
		history = append(history, workflow.WorkflowTransition{From: "RUNNING", To: status, WorkerID: owner, Version: 3, Reason: reason})
	}
	return s.Version == version && reflect.DeepEqual(s.History, history)
}

func coherentInvocationSnapshot(s *workflow.WorkflowExecution, id string) bool {
	if freshInvocationSnapshot(s, id) {
		return true
	}
	if s == nil || s.ClaimOwner == "" {
		return false
	}
	switch s.Status {
	case "RUNNING", "COMPLETED":
		return invocationSnapshotMatches(s, id, s.ClaimOwner, s.Status, "")
	case "FAILED":
		return s.LastError != "" && invocationSnapshotMatches(s, id, s.ClaimOwner, "FAILED", s.LastError)
	}
	return false
}

// The observer stores observations only. Each method forwards exactly once,
// preserves arguments/results and never retries or writes state independently.
// It is used synchronously by the existing runner and is private to one call.
type lifecycleObservation struct {
	operation string
	id, owner string
	version   uint64
	reason    string
	accepted  bool
	snapshot  *workflow.WorkflowExecution
	exists    bool
}
type lifecycleObserver struct {
	repo       workflow.WorkflowLifecycleRepository
	operations []lifecycleObservation
}

func (o *lifecycleObserver) TryClaim(id, owner string) bool {
	ok := o.repo.TryClaim(id, owner)
	o.operations = append(o.operations, lifecycleObservation{operation: "claim", id: id, owner: owner, accepted: ok})
	return ok
}
func (o *lifecycleObserver) Get(id string) (*workflow.WorkflowExecution, bool) {
	s, ok := o.repo.Get(id)
	o.operations = append(o.operations, lifecycleObservation{operation: "get", id: id, snapshot: s, exists: ok})
	return s, ok
}
func (o *lifecycleObserver) Complete(id, owner string, version uint64) bool {
	ok := o.repo.Complete(id, owner, version)
	o.operations = append(o.operations, lifecycleObservation{operation: "complete", id: id, owner: owner, version: version, accepted: ok})
	return ok
}
func (o *lifecycleObserver) Fail(id, owner string, version uint64, reason string) bool {
	ok := o.repo.Fail(id, owner, version, reason)
	o.operations = append(o.operations, lifecycleObservation{operation: "fail", id: id, owner: owner, version: version, reason: reason, accepted: ok})
	return ok
}
func classifyInvocation(out InvocationOutcome, runErr error, o *lifecycleObserver, d workflow.WorkflowDefinition, worker string) (InvocationOutcome, error) {
	uncertain := func() (InvocationOutcome, error) {
		out.Kind = OutcomeReconciliation
		return out, invocationError(InvocationReconciliation, runErr)
	}
	ops := o.operations
	if len(ops) == 0 {
		return uncertain()
	}
	claim := ops[0]
	if claim.operation != "claim" || claim.id != out.ExecutionID || claim.owner != worker {
		return uncertain()
	}
	out.Claim = ClaimRejected
	if claim.accepted {
		out.Claim = ClaimAccepted
	}
	if !coherentInvocationSnapshot(out.Snapshot, out.ExecutionID) {
		return uncertain()
	}
	if !claim.accepted {
		if len(ops) != 1 || runErr == nil || out.Cursor != d.StartStep {
			return uncertain()
		}
		out.Kind = OutcomeClaimRejected
		return out, invocationError(InvocationClaimRejected, runErr)
	}
	if len(ops) < 2 || ops[1].operation != "get" || ops[1].id != out.ExecutionID || !ops[1].exists || !invocationSnapshotMatches(ops[1].snapshot, out.ExecutionID, worker, "RUNNING", "") {
		return uncertain()
	}
	if len(ops) == 2 {
		step, ok := d.Steps[out.Cursor]
		if runErr != nil || !ok || (step.Type != workflow.StepWait && step.Type != workflow.StepCondition) || !invocationSnapshotMatches(out.Snapshot, out.ExecutionID, worker, "RUNNING", "") {
			return uncertain()
		}
		out.Kind = OutcomeHandoff
		return out, nil
	}
	if len(ops) != 3 {
		return uncertain()
	}
	terminal := ops[2]
	if terminal.id != out.ExecutionID || terminal.owner != worker || terminal.version != ops[1].snapshot.Version || !terminal.accepted {
		return uncertain()
	}
	switch terminal.operation {
	case "complete":
		if runErr != nil || out.Cursor != "" || !invocationSnapshotMatches(out.Snapshot, out.ExecutionID, worker, "COMPLETED", "") {
			return uncertain()
		}
		out.Kind = OutcomeCompleted
		return out, nil
	case "fail":
		if runErr == nil || terminal.reason == "" || !invocationSnapshotMatches(out.Snapshot, out.ExecutionID, worker, "FAILED", terminal.reason) {
			return uncertain()
		}
		out.Kind = OutcomeFailed
		return out, invocationError(InvocationExecutionFailed, runErr)
	default:
		return uncertain()
	}
}
