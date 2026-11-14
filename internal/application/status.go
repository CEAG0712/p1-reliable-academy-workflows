package application

import (
	"context"
	"reflect"

	workflow "example.com/p1-reliable-academy-workflows"
)

// ExecutionReader supplies independent, coherent snapshots stable for the call.
// The query boundary has no lifecycle or action capability. Get is synchronous
// and cannot be interrupted by context cancellation.
type ExecutionReader interface {
	Get(string) (*workflow.WorkflowExecution, bool)
}

// Status authority is independent of invocation authority. Nonzero size makes
// distinct live scopes distinguishable by pointer identity.
type localStatusAuthority struct{ marker byte }

// TrustedLocalStatusQuery can only be minted by trusted local composition.
// JSON and ordinary external callers cannot populate its private fields.
type TrustedLocalStatusQuery struct {
	authority   *localStatusAuthority
	executionID string
}
type ExecutionStatusResult struct {
	ExecutionID string
	Status      string
	Version     uint64
}
type LocalExecutionStatusQuerier interface {
	QueryLocal(context.Context, TrustedLocalStatusQuery) (ExecutionStatusResult, error)
}

type StatusQueryErrorKind string

const (
	StatusQueryInvalidCommand   StatusQueryErrorKind = "invalid_command"
	StatusQueryMissingAuthority StatusQueryErrorKind = "missing_authority"
	StatusQueryForbidden        StatusQueryErrorKind = "forbidden"
	StatusQueryConcealed        StatusQueryErrorKind = "concealed"
	StatusQueryNotFound         StatusQueryErrorKind = "not_found"
	StatusQueryUnavailable      StatusQueryErrorKind = "unavailable"
	StatusQueryInvalidSnapshot  StatusQueryErrorKind = "invalid_snapshot"
)

// Cause is internal evidence; adapters must use fixed public messages.
type StatusQueryError struct {
	Kind  StatusQueryErrorKind
	Cause error
}

func (e *StatusQueryError) Error() string { return "local status query: " + string(e.Kind) }
func (e *StatusQueryError) Unwrap() error { return e.Cause }
func statusQueryError(kind StatusQueryErrorKind, cause error) error {
	return &StatusQueryError{Kind: kind, Cause: cause}
}

type statusVisibility uint8

const (
	statusVisible statusVisibility = iota + 1
	statusConcealed
)

type statusVisibilityPolicy interface {
	evaluate(context.Context, string) (statusVisibility, error)
}
type statusAllowlist map[string]struct{}

func (p statusAllowlist) evaluate(_ context.Context, id string) (statusVisibility, error) {
	if _, ok := p[id]; ok {
		return statusVisible, nil
	}
	return statusConcealed, nil
}

// LocalStatusScope is local fixture visibility, NOT production authentication or
// tenant ownership. Trusted host code owns the reader and controls composition.
// Never attach its capability to operational requests. Its allowlist is copied
// once and never mutated; concurrent queries require a concurrency-safe reader.
type LocalStatusScope struct {
	authority *localStatusAuthority
	service   *statusService
}

func NewLocalStatusScope(reader ExecutionReader, visibleExecutionIDs []string) (*LocalStatusScope, error) {
	allowed := make(statusAllowlist, len(visibleExecutionIDs))
	for _, id := range visibleExecutionIDs {
		if !validLocalIdentifier(id) {
			return nil, statusQueryError(StatusQueryInvalidCommand, nil)
		}
		allowed[id] = struct{}{}
	}
	authority := &localStatusAuthority{marker: 1}
	service, err := newStatusService(authority, reader, allowed)
	if err != nil {
		return nil, err
	}
	return &LocalStatusScope{authority: authority, service: service}, nil
}
func (s *LocalStatusScope) configured() bool {
	return s != nil && s.authority != nil && s.service != nil && s.service.authority == s.authority && s.service.configured()
}
func (s *LocalStatusScope) Querier() LocalExecutionStatusQuerier {
	if !s.configured() {
		return nil
	}
	return s.service
}

// Command validates syntax only. The service independently enforces visibility.
// Possession of an ID alone grants no visibility.
func (s *LocalStatusScope) Command(id string) (TrustedLocalStatusQuery, error) {
	if !s.configured() {
		return TrustedLocalStatusQuery{}, statusQueryError(StatusQueryMissingAuthority, nil)
	}
	if !validLocalIdentifier(id) {
		return TrustedLocalStatusQuery{}, statusQueryError(StatusQueryInvalidCommand, nil)
	}
	return TrustedLocalStatusQuery{authority: s.authority, executionID: id}, nil
}

type statusService struct {
	authority *localStatusAuthority
	reader    ExecutionReader
	policy    statusVisibilityPolicy
}

func nilStatusDependency(v any) bool {
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
func (s *statusService) configured() bool {
	return s != nil && s.authority != nil && !nilStatusDependency(s.reader) && !nilStatusDependency(s.policy)
}

// Private fault-injection seam. Construction never evaluates policy or reads.
func newStatusService(authority *localStatusAuthority, reader ExecutionReader, policy statusVisibilityPolicy) (*statusService, error) {
	s := &statusService{authority: authority, reader: reader, policy: policy}
	if !s.configured() {
		return nil, statusQueryError(StatusQueryUnavailable, nil)
	}
	return s, nil
}
func (s *statusService) QueryLocal(ctx context.Context, c TrustedLocalStatusQuery) (ExecutionStatusResult, error) {
	reject := func(kind StatusQueryErrorKind, cause error) (ExecutionStatusResult, error) {
		return ExecutionStatusResult{}, statusQueryError(kind, cause)
	}
	if !s.configured() {
		return reject(StatusQueryUnavailable, nil)
	}
	if c.authority == nil {
		return reject(StatusQueryMissingAuthority, nil)
	}
	if c.authority != s.authority {
		return reject(StatusQueryForbidden, nil)
	}
	if ctx == nil {
		return reject(StatusQueryInvalidCommand, nil)
	}
	if err := ctx.Err(); err != nil {
		return reject(StatusQueryUnavailable, err)
	}
	if !validLocalIdentifier(c.executionID) {
		return reject(StatusQueryInvalidCommand, nil)
	}
	decision, err := s.policy.evaluate(ctx, c.executionID)
	if cause := ctx.Err(); cause != nil {
		return reject(StatusQueryUnavailable, cause)
	}
	if err != nil {
		return reject(StatusQueryUnavailable, err)
	}
	switch decision {
	case statusConcealed:
		return reject(StatusQueryConcealed, nil)
	case statusVisible:
	default:
		return reject(StatusQueryUnavailable, nil)
	}
	snapshot, exists := s.reader.Get(c.executionID)
	if !exists {
		if snapshot != nil {
			return reject(StatusQueryInvalidSnapshot, nil)
		}
		return reject(StatusQueryNotFound, nil)
	}
	if !validStatusSnapshot(snapshot, c.executionID) {
		return reject(StatusQueryInvalidSnapshot, nil)
	}
	out := ExecutionStatusResult{ExecutionID: snapshot.WorkflowID, Status: snapshot.Status, Version: snapshot.Version}
	if cause := ctx.Err(); cause != nil {
		return reject(StatusQueryUnavailable, cause)
	}
	return out, nil
}

// Validate the stored record, not the invocation fixture's v1/2/3 timeline.
// Coherent historical suffixes are accepted; absent active/terminal history is
// ambiguous. No repairs, normalization, rereads or writes occur.
func validStatusSnapshot(s *workflow.WorkflowExecution, id string) bool {
	if s == nil || s.WorkflowID != id || s.Version == 0 || s.Attempts < 0 {
		return false
	}
	switch s.Status {
	case "READY":
		return s.ClaimOwner == "" && s.Attempts == 0 && len(s.History) == 0 // Legacy LastError is allowed.
	case "RUNNING": // Legacy TryStart can leave ClaimOwner empty and preserve a diagnostic.
	case "COMPLETED":
		if s.ClaimOwner == "" || s.LastError != "" {
			return false
		}
	case "FAILED":
		if s.ClaimOwner == "" || s.LastError == "" {
			return false
		}
	default:
		return false
	}
	if s.Attempts == 0 || len(s.History) == 0 {
		return false
	}
	claims := 0
	for i, h := range s.History {
		if h.Version == 0 || h.WorkerID != s.ClaimOwner {
			return false
		}
		if i > 0 && (s.History[i-1].To != h.From || s.History[i-1].Version >= h.Version) {
			return false
		}
		switch {
		case h.From == "READY" && h.To == "RUNNING":
			claims++
		case h.From == "RUNNING" && (h.To == "COMPLETED" || h.To == "FAILED"):
			if h.WorkerID == "" {
				return false
			}
		default:
			return false
		}
		if h.To == "FAILED" {
			if h.Reason == "" || h.Reason != s.LastError {
				return false
			}
		} else if h.Reason != "" {
			return false
		}
	}
	last := s.History[len(s.History)-1]
	return claims <= s.Attempts && last.To == s.Status && last.Version == s.Version
}
