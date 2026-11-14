package application

import (
	"context"
	"sync"

	workflow "example.com/p1-reliable-academy-workflows"
)

// TrustedLocalRequestCommand is minted by a local fixture scope. JSON cannot
// populate its authority, key or selection. It conveys no tenant/payment proof.
type TrustedLocalRequestCommand struct {
	authority *requestIdentityAuthority
	key       string
	selection SubmissionCommand
}

type LocalRequestInvoker interface {
	InvokeLocalRequest(context.Context, TrustedLocalRequestCommand) (InvocationOutcome, error)
}

type requestIdentityAuthority struct{ marker byte }
type requestReservationState uint8

const (
	requestReserved requestReservationState = iota + 1
	requestKnown
	requestUncertain
)

type requestReservation struct {
	selection   SubmissionCommand // Exact value tuple, not concatenation or a lossy hash.
	state       requestReservationState
	executionID string // Internal evidence only; never a replay response or visibility grant.
}

// LocalRequestIdentityScope owns bounded, process-local request reservations.
// Reuse one scope for its request population; separate scopes do not coordinate,
// even with the same repository. Never expose this fixture capability on an
// operational route. It is not authentication or durable request idempotency.
// Concurrent invocation requires a concurrency-safe ID source. Do not copy it.
type LocalRequestIdentityScope struct {
	service *requestIdentityService
}

type requestIdentityService struct {
	authority    *requestIdentityAuthority
	invocation   *invocationService
	limit        int
	mu           sync.Mutex
	reservations map[string]*requestReservation
}

// NewLocalRequestIdentityScope performs no operational dependency calls. The
// public composition accepts only the existing concrete repository and package-
// owned harmless fixtures, with no arbitrary resolver/runner injection path.
func NewLocalRequestIdentityScope(repo *workflow.InMemoryWorkflowRepository, ids ExecutionIDSource, workerID string, maxReservations int) (*LocalRequestIdentityScope, error) {
	core, err := NewLocalInvocationScope(repo, ids, workerID)
	if err != nil {
		return nil, err
	}
	return newLocalRequestIdentityScope(core.service, maxReservations)
}

// Private seam permits controlled faults in package tests only.
func newLocalRequestIdentityScope(core *invocationService, limit int) (*LocalRequestIdentityScope, error) {
	if core == nil || !validInvocationDependencies(core.deps) || limit <= 0 {
		return nil, invocationError(InvocationUnavailable, nil)
	}
	return &LocalRequestIdentityScope{service: &requestIdentityService{
		authority: &requestIdentityAuthority{marker: 1}, invocation: core,
		limit: limit, reservations: make(map[string]*requestReservation),
	}}, nil
}
func (s *requestIdentityService) configured() bool {
	return s != nil && s.authority != nil && s.invocation != nil &&
		validInvocationDependencies(s.invocation.deps) && s.limit > 0 && s.reservations != nil
}
func (s *LocalRequestIdentityScope) Invoker() LocalRequestInvoker {
	if s == nil || !s.service.configured() {
		return nil
	}
	return s.service
}
func (s *LocalRequestIdentityScope) ApprovedDefinitions() map[string]struct{} {
	if s == nil || !s.service.configured() {
		return nil
	}
	return map[string]struct{}{"local.success": {}, "local.failure": {}, "local.wait": {}, "local.condition": {}}
}
func validRequestSelection(key string, selection SubmissionCommand) error {
	if !validLocalIdentifier(key) || !validLocalIdentifier(selection.WorkflowDefinitionID) || !validLocalIdentifier(selection.EventID) {
		return invocationError(InvocationInvalidCommand, nil)
	}
	if !approvedLocalDefinition(selection.WorkflowDefinitionID) {
		return invocationError(InvocationForbidden, nil)
	}
	return nil
}
func (s *LocalRequestIdentityScope) Command(key string, selection SubmissionCommand) (TrustedLocalRequestCommand, error) {
	if s == nil || !s.service.configured() {
		return TrustedLocalRequestCommand{}, invocationError(InvocationMissingAuthority, nil)
	}
	if err := validRequestSelection(key, selection); err != nil {
		return TrustedLocalRequestCommand{}, err
	}
	return TrustedLocalRequestCommand{authority: s.service.authority, key: key, selection: selection}, nil
}

func (s *requestIdentityService) InvokeLocalRequest(ctx context.Context, c TrustedLocalRequestCommand) (InvocationOutcome, error) {
	rejected := InvocationOutcome{Creation: CreationNotAttempted, Claim: ClaimNotObserved, Kind: OutcomeRejected}
	if !s.configured() {
		return rejected, invocationError(InvocationUnavailable, nil)
	}
	if c.authority == nil {
		return rejected, invocationError(InvocationMissingAuthority, nil)
	}
	if c.authority != s.authority {
		return rejected, invocationError(InvocationForbidden, nil)
	}
	if ctx == nil {
		return rejected, invocationError(InvocationInvalidCommand, nil)
	}
	if err := validRequestSelection(c.key, c.selection); err != nil {
		return rejected, err
	}
	if err := ctx.Err(); err != nil {
		return rejected, invocationError(InvocationUnavailable, err)
	}

	s.mu.Lock()
	// Recheck cancellation at the admission boundary. Cancellation racing this
	// check can still win a reservation; once admitted, that key is never released.
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return rejected, invocationError(InvocationUnavailable, err)
	}
	if prior, exists := s.reservations[c.key]; exists {
		kind := InvocationDuplicate
		if prior.selection == c.selection && prior.state == requestUncertain {
			kind = InvocationReconciliation
			rejected.Kind = OutcomeReconciliation
		}
		s.mu.Unlock()
		return rejected, invocationError(kind, nil)
	}
	if len(s.reservations) >= s.limit {
		s.mu.Unlock()
		return rejected, invocationError(InvocationUnavailable, nil)
	}
	reservation := &requestReservation{selection: c.selection, state: requestReserved}
	s.reservations[c.key] = reservation
	s.mu.Unlock()

	// No dependency calls occur under mu. This defer runs on cancellation,
	// returned errors and panic unwinding. A panic propagates, while its key stays
	// protected and becomes uncertain. There is no eviction, retry or result cache.
	state, executionID := requestUncertain, ""
	defer func() {
		s.mu.Lock()
		reservation.state, reservation.executionID = state, executionID
		s.mu.Unlock()
	}()
	out, err := s.invocation.InvokeLocal(ctx, TrustedLocalInvocationCommand{
		authority: s.invocation.deps.authority, selection: c.selection,
	})
	executionID = out.ExecutionID
	// Existing M4.5 classifies false/error as rejected. This opt-in policy is
	// deliberately more conservative: ANY creation-call error blocks replay and
	// requires reconciliation. A FAILED outcome backed by accepted Fail plus the
	// authoritative snapshot/history remains known, retaining its original error.
	if out.Creation == CreationRejected || out.Creation == CreationUncertain || out.Kind == OutcomeReconciliation {
		out.Kind = OutcomeReconciliation
		return out, invocationError(InvocationReconciliation, err)
	}
	state = requestKnown // Classification comes from the sealed concrete M4.5 service.
	return out, err
}
