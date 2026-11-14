package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

const httpInvocationID = "exec-00000000000000000000000000000001"

type httpInvocationIDs struct {
	calls int
	ctx   context.Context
	after func()
	err   error
}

func (p *httpInvocationIDs) NextExecutionID(ctx context.Context) (string, error) {
	p.calls++
	p.ctx = ctx
	if p.after != nil {
		p.after()
	}
	return httpInvocationID, p.err
}

type httpInvocationProbe struct {
	delegate application.LocalWorkflowInvoker
	calls    int
	ctx      context.Context
	command  application.TrustedLocalInvocationCommand
	outcome  application.InvocationOutcome
	err      error
}

func (p *httpInvocationProbe) InvokeLocal(ctx context.Context, c application.TrustedLocalInvocationCommand) (application.InvocationOutcome, error) {
	p.calls++
	p.ctx = ctx
	p.command = c
	if p.delegate != nil {
		p.outcome, p.err = p.delegate.InvokeLocal(ctx, c)
	}
	return p.outcome, p.err
}
func httpInvocationFixture(t *testing.T) (*application.LocalInvocationScope, *workflow.InMemoryWorkflowRepository, *httpInvocationIDs) {
	t.Helper()
	r := workflow.NewInMemoryWorkflowRepository()
	ids := &httpInvocationIDs{}
	s, e := application.NewLocalInvocationScope(r, ids, "http-local-worker")
	if e != nil {
		t.Fatal(e)
	}
	return s, r, ids
}
func invocationRouter(t *testing.T, h gin.HandlerFunc) *gin.Engine {
	t.Helper()
	r := gin.New()
	r.POST("/test-only/invoke", h)
	return r
}
func sendInvocation(r *gin.Engine, body string, ctx context.Context) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test-only/invoke", strings.NewReader(body))
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	r.ServeHTTP(w, req)
	return w
}
func TestLocalInvocationHandlerConstruction(t *testing.T) {
	for _, s := range []*application.LocalInvocationScope{nil, {}} {
		if h, e := NewLocalInvocationHandler(s); h != nil || e == nil {
			t.Fatal("zero scope")
		}
	}
	s, _, ids := httpInvocationFixture(t)
	var typedNil *httpInvocationProbe
	for _, invoker := range []application.LocalWorkflowInvoker{nil, typedNil} {
		if h, e := newLocalInvocationHandler(s, invoker, s.ApprovedDefinitions()); h != nil || e == nil {
			t.Fatal("nil invoker")
		}
	}
	p := &httpInvocationProbe{}
	approved := s.ApprovedDefinitions()
	h, e := newLocalInvocationHandler(s, p, approved)
	if e != nil {
		t.Fatal(e)
	}
	delete(approved, "local.success")
	if p.calls != 0 || ids.calls != 0 {
		t.Fatal("constructor called collaborator")
	}
	sendInvocation(invocationRouter(t, h), validationBody("local.success", "event"), nil)
	if p.calls != 1 {
		t.Fatal("support set not snapshotted")
	}
}
func TestLocalInvocationHandlerRejectedRequests(t *testing.T) {
	for _, tc := range invalidValidationCases() {
		t.Run(tc.name, func(t *testing.T) {
			s, repo, ids := httpInvocationFixture(t)
			p := &httpInvocationProbe{delegate: s.Invoker()}
			h, e := newLocalInvocationHandler(s, p, tc.approved)
			if e != nil {
				t.Fatal(e)
			}
			w := sendInvocation(invocationRouter(t, h), tc.body, nil)
			if w.Code != 400 || p.calls != 0 || ids.calls != 0 {
				t.Fatalf("rejection status=%d calls=%d", w.Code, p.calls)
			}
			assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: tc.message}})
			if _, ok := repo.Get(httpInvocationID); ok {
				t.Fatal("rejection created")
			}
		})
	}
	s, _, ids := httpInvocationFixture(t)
	h, e := NewLocalInvocationHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	r := invocationRouter(t, h)
	for _, field := range []string{"tenant_id", "caller_id", "payment_status", "authorized", "definition", "execution_id", "worker_id"} {
		w := sendInvocation(r, `{"workflow_definition_id":"local.success","event_id":"CONFIRMED","`+field+`":"CONFIRMED"}`, nil)
		if w.Code != 400 || ids.calls != 0 {
			t.Fatal("untrusted field authorized")
		}
	}
}
func TestLocalInvocationHandlerBodyBounds(t *testing.T) {
	for _, fail := range []bool{false, true} {
		s, _, ids := httpInvocationFixture(t)
		h, e := NewLocalInvocationHandler(s)
		if e != nil {
			t.Fatal(e)
		}
		body := &preparationBody{fail: fail}
		req := httptest.NewRequest("POST", "/test-only/invoke", nil)
		req.Body = body
		w := httptest.NewRecorder()
		invocationRouter(t, h).ServeHTTP(w, req)
		if w.Code != 400 || ids.calls != 0 {
			t.Fatal("body rejection entered engine")
		}
		message := "body must not exceed 4096 bytes"
		if fail {
			message = "could not read request body"
		} else if body.read != 4097 {
			t.Fatal("unbounded stream")
		}
		assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: message}})
	}
	s, _, _ := httpInvocationFixture(t)
	h, e := NewLocalInvocationHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	body := validationBody("local.success", "event")
	w := sendInvocation(invocationRouter(t, h), body+strings.Repeat(" ", 4096-len(body)), nil)
	if w.Code != 201 {
		t.Fatalf("4096 rejected %s", w.Body)
	}
}
func TestLocalInvocationHandlerRealEngine(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			s, repo, ids := httpInvocationFixture(t)
			p := &httpInvocationProbe{delegate: s.Invoker()}
			h, e := newLocalInvocationHandler(s, p, s.ApprovedDefinitions())
			if e != nil {
				t.Fatal(e)
			}
			ctx := context.WithValue(context.Background(), struct{}{}, "exact-context")
			w := sendInvocation(invocationRouter(t, h), validationBody(definition, "CONFIRMED"), ctx)
			expected, e := s.Command(application.SubmissionCommand{WorkflowDefinitionID: definition, EventID: "CONFIRMED"})
			if e != nil {
				t.Fatal(e)
			}
			if w.Code != 201 || p.calls != 1 || p.ctx != ctx || p.command != expected || ids.ctx != ctx || ids.calls != 1 {
				t.Fatalf("delegation %d %+v", w.Code, p)
			}
			snapshot, ok := repo.Get(httpInvocationID)
			if !ok {
				t.Fatal("real repository missing")
			}
			status := "COMPLETED"
			if definition == "local.failure" {
				status = "FAILED"
			}
			if definition == "local.wait" || definition == "local.condition" {
				status = "RUNNING"
			}
			version := uint64(3)
			entries := 2
			if status == "RUNNING" {
				version = 2
				entries = 1
			}
			if snapshot.Status != status || snapshot.Version != version || snapshot.ClaimOwner != "http-local-worker" || snapshot.Attempts != 1 || len(snapshot.History) != entries || snapshot.History[0].To != "RUNNING" {
				t.Fatalf("actual runner state %+v", snapshot)
			}
			if status == "FAILED" && (snapshot.LastError == "" || snapshot.History[1].Reason != snapshot.LastError) {
				t.Fatal("missing failure history")
			}
			assertPreparationJSON(t, w, CreateExecutionResponse{ExecutionID: httpInvocationID, Status: status})
			if strings.Contains(w.Body.String(), "controlled") || strings.Contains(w.Body.String(), "owner") || strings.Contains(w.Body.String(), "cursor") {
				t.Fatal("diagnostic leak")
			}
		})
	}
}
func TestLocalInvocationHandlerDuplicateAndCancellation(t *testing.T) {
	s, repo, ids := httpInvocationFixture(t)
	h, e := NewLocalInvocationHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	r := invocationRouter(t, h)
	if w := sendInvocation(r, validationBody("local.success", "event"), nil); w.Code != 201 {
		t.Fatal(w.Body)
	}
	before, _ := repo.Get(httpInvocationID)
	w := sendInvocation(r, validationBody("local.success", "other-event"), nil)
	if w.Code != 409 || ids.calls != 2 {
		t.Fatalf("duplicate %d %s", w.Code, w.Body)
	}
	assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: "CONFLICT", Message: "execution conflict"}})
	after, _ := repo.Get(httpInvocationID)
	if after.Version != before.Version || len(after.History) != len(before.History) {
		t.Fatal("duplicate rerun")
	}
	s, repo, ids = httpInvocationFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ids.after = cancel
	h, e = NewLocalInvocationHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	w = sendInvocation(invocationRouter(t, h), validationBody("local.success", "event"), ctx)
	if w.Code != 503 {
		t.Fatal(w.Code)
	}
	if _, ok := repo.Get(httpInvocationID); ok {
		t.Fatal("cancel created")
	}
}
func TestLocalInvocationHandlerSafeErrors(t *testing.T) {
	for _, tc := range []struct {
		kind   application.InvocationErrorKind
		status int
		code   string
	}{
		{application.InvocationMissingAuthority, 401, "UNAUTHENTICATED"}, {application.InvocationForbidden, 403, "FORBIDDEN"}, {application.InvocationInvalidCommand, 400, "INVALID_REQUEST"}, {application.InvocationDefinitionNotFound, 404, "NOT_FOUND"}, {application.InvocationDuplicate, 409, "CONFLICT"}, {application.InvocationClaimRejected, 409, "CONFLICT"}, {application.InvocationUnavailable, 503, "UNAVAILABLE"}, {application.InvocationInvalidDefinition, 500, "INTERNAL_ERROR"}, {application.InvocationReconciliation, 500, "INTERNAL_ERROR"},
	} {
		t.Run(string(tc.kind), func(t *testing.T) {
			s, _, _ := httpInvocationFixture(t)
			p := &httpInvocationProbe{err: fmt.Errorf("private wrapper: %w", &application.InvocationError{Kind: tc.kind, Cause: errors.New("secret payment tenant diagnostic")})}
			h, e := newLocalInvocationHandler(s, p, s.ApprovedDefinitions())
			if e != nil {
				t.Fatal(e)
			}
			w := sendInvocation(invocationRouter(t, h), validationBody("local.success", "event"), nil)
			if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.code) || strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret") {
				t.Fatalf("unsafe error %d %s", w.Code, w.Body)
			}
		})
	}
}

type httpInvocationStripCapability struct {
	delegate application.LocalWorkflowInvoker
}

func (p httpInvocationStripCapability) InvokeLocal(ctx context.Context, _ application.TrustedLocalInvocationCommand) (application.InvocationOutcome, error) {
	return p.delegate.InvokeLocal(ctx, application.TrustedLocalInvocationCommand{})
}
func TestLocalInvocationHandlerAuthorityFailsClosed(t *testing.T) {
	scope, repo, ids := httpInvocationFixture(t)
	other, otherRepo, otherIDs := httpInvocationFixture(t)
	for _, invoker := range []application.LocalWorkflowInvoker{httpInvocationStripCapability{scope.Invoker()}, other.Invoker()} {
		h, err := newLocalInvocationHandler(scope, invoker, scope.ApprovedDefinitions())
		if err != nil {
			t.Fatal(err)
		}
		w := sendInvocation(invocationRouter(t, h), validationBody("local.success", "CONFIRMED"), context.WithValue(context.Background(), "authorized", true))
		if w.Code != 401 && w.Code != 403 {
			t.Fatalf("authority bypass %d %s", w.Code, w.Body)
		}
	}
	if ids.calls != 0 || otherIDs.calls != 0 {
		t.Fatal("unauthorized ID generation")
	}
	for _, r := range []*workflow.InMemoryWorkflowRepository{repo, otherRepo} {
		if _, ok := r.Get(httpInvocationID); ok {
			t.Fatal("unauthorized creation")
		}
	}
}
func TestLocalInvocationHandlerCreatedNotStartedProjection(t *testing.T) {
	scope, repo, _ := httpInvocationFixture(t)
	if ok, e := repo.CreateIfAbsent(workflow.NewWorkflowExecution(httpInvocationID)); !ok || e != nil {
		t.Fatal(ok, e)
	}
	snapshot, ok := repo.Get(httpInvocationID)
	if !ok {
		t.Fatal("missing resource")
	}
	p := &httpInvocationProbe{outcome: application.InvocationOutcome{ExecutionID: httpInvocationID, Creation: application.CreationInserted, Kind: application.OutcomeCreatedNotStarted, Claim: application.ClaimNotObserved, Snapshot: snapshot}, err: &application.InvocationError{Kind: application.InvocationUnavailable, Cause: context.Canceled}}
	h, e := newLocalInvocationHandler(scope, p, scope.ApprovedDefinitions())
	if e != nil {
		t.Fatal(e)
	}
	w := sendInvocation(invocationRouter(t, h), validationBody("local.success", "event"), nil)
	if w.Code != 201 {
		t.Fatalf("READY projection %d %s", w.Code, w.Body)
	}
	assertPreparationJSON(t, w, CreateExecutionResponse{ExecutionID: httpInvocationID, Status: "READY"})
}
func TestLocalInvocationHandlerInconsistentOutcome(t *testing.T) {
	for _, mode := range []string{"zero", "missing-snapshot", "wrong-identity", "waiting", "reconciliation", "error-on-completed", "invalid-version"} {
		t.Run(mode, func(t *testing.T) {
			scope, _, _ := httpInvocationFixture(t)
			p := &httpInvocationProbe{}
			switch mode {
			case "zero":
			default:
				p.outcome = application.InvocationOutcome{ExecutionID: httpInvocationID, Creation: application.CreationInserted, Kind: application.OutcomeCompleted, Snapshot: &workflow.WorkflowExecution{WorkflowID: httpInvocationID, Status: "COMPLETED"}}
			}
			switch mode {
			case "missing-snapshot":
				p.outcome.Snapshot = nil
			case "wrong-identity":
				p.outcome.Snapshot.WorkflowID = "other"
			case "waiting":
				p.outcome.Snapshot.Status = "WAITING"
			case "reconciliation":
				p.outcome.Kind = application.OutcomeReconciliation
			case "error-on-completed":
				p.err = errors.New("private unexpected diagnostic")
			}
			h, e := newLocalInvocationHandler(scope, p, scope.ApprovedDefinitions())
			if e != nil {
				t.Fatal(e)
			}
			w := sendInvocation(invocationRouter(t, h), validationBody("local.success", "event"), nil)
			if w.Code != 500 {
				t.Fatalf("inconsistent success %d %s", w.Code, w.Body)
			}
			assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: "INTERNAL_ERROR", Message: "internal error"}})
		})
	}
}
