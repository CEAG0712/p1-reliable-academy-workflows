package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

type httpStatusReader struct {
	reader   application.ExecutionReader
	calls    int
	snapshot *workflow.WorkflowExecution
	exists   bool
}

func (p *httpStatusReader) Get(id string) (*workflow.WorkflowExecution, bool) {
	p.calls++
	if p.reader != nil {
		return p.reader.Get(id)
	}
	return p.snapshot, p.exists
}

type httpStatusProbe struct {
	delegate application.LocalExecutionStatusQuerier
	calls    int
	ctx      context.Context
	command  application.TrustedLocalStatusQuery
	result   application.ExecutionStatusResult
	err      error
	strip    bool
}

func (p *httpStatusProbe) QueryLocal(ctx context.Context, c application.TrustedLocalStatusQuery) (application.ExecutionStatusResult, error) {
	p.calls++
	p.ctx = ctx
	p.command = c
	if p.strip {
		c = application.TrustedLocalStatusQuery{}
	}
	if p.delegate != nil {
		return p.delegate.QueryLocal(ctx, c)
	}
	return p.result, p.err
}
func httpStatusScope(t *testing.T, r application.ExecutionReader, ids ...string) *application.LocalStatusScope {
	t.Helper()
	s, e := application.NewLocalStatusScope(r, ids)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func statusRouter(h gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET("/v1/workflow-executions/:execution_id", h)
	return r
}
func sendStatus(r *gin.Engine, id string, ctx context.Context) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/workflow-executions/"+url.PathEscape(id), nil)
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	r.ServeHTTP(w, req)
	return w
}
func assertStatusJSON(t *testing.T, w *httptest.ResponseRecorder, status int, want any) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status=%d body=%s", w.Code, w.Body)
	}
	if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("headers=%v", w.Header())
	}
	expected, e := json.Marshal(want)
	if e != nil {
		t.Fatal(e)
	}
	if w.Body.String() != string(expected) {
		t.Fatalf("body=%s want=%s", w.Body, expected)
	}
}
func TestLocalStatusHandlerConstruction(t *testing.T) {
	for _, s := range []*application.LocalStatusScope{nil, {}} {
		if h, e := NewLocalStatusHandler(s); h != nil || e == nil {
			t.Fatal("zero scope accepted")
		}
	}
	r := &httpStatusReader{}
	s := httpStatusScope(t, r, "run")
	var typedNil *httpStatusProbe
	for _, q := range []application.LocalExecutionStatusQuerier{nil, typedNil} {
		if h, e := newLocalStatusHandler(s, q); h != nil || e == nil {
			t.Fatal("nil querier")
		}
	}
	p := &httpStatusProbe{}
	if _, e := newLocalStatusHandler(s, p); e != nil || p.calls != 0 || r.calls != 0 {
		t.Fatal("constructor queried")
	}
}
func TestLocalStatusHandlerPathValidation(t *testing.T) {
	invalid := []string{"", strings.Repeat("a", 129), "a b", " padded", "é", "a/b", "a%2fb", "a\x00", "a\n", "a:b", "a@b", "a?b", "a#b", "a\\b", "a+b"}
	for _, id := range invalid {
		t.Run(fmt.Sprintf("%q", id), func(t *testing.T) {
			reader := &httpStatusReader{}
			s := httpStatusScope(t, reader)
			p := &httpStatusProbe{}
			h, e := newLocalStatusHandler(s, p)
			if e != nil {
				t.Fatal(e)
			}
			// Invoke a matched Gin context directly for values a router would not match.
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Params = gin.Params{{Key: "execution_id", Value: id}}
			h(c)
			assertStatusJSON(t, w, 400, ErrorResponse{APIError{"INVALID_REQUEST", "invalid execution ID"}})
			if p.calls != 0 || reader.calls != 0 {
				t.Fatal("invalid path queried")
			}
		})
	}
	for _, id := range []string{"workflow-001", "exec-00000000000000000000000000000001", "AZaz09_.-", strings.Repeat("x", 128), "Case.Sensitive"} {
		reader := &httpStatusReader{snapshot: workflow.NewWorkflowExecution(id), exists: true}
		s := httpStatusScope(t, reader, id)
		h, e := NewLocalStatusHandler(s)
		if e != nil {
			t.Fatal(e)
		}
		assertStatusJSON(t, sendStatus(statusRouter(h), id, nil), 200, ExecutionStatusResponse{id, "READY", 1})
		if reader.calls != 1 {
			t.Fatal("valid path reads")
		}
	}
	s := httpStatusScope(t, &httpStatusReader{})
	h, e := NewLocalStatusHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	r := statusRouter(h)
	w := sendStatus(r, "", nil)
	if w.Code != 404 || w.Body.String() != "404 page not found" {
		t.Fatal("missing segment reached handler")
	}
	// Gin decodes the URL path once. The handler does not decode the residual %.
	for _, path := range []string{"/v1/workflow-executions/a%252Fb", "/v1/workflow-executions/a%20b"} {
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		assertStatusJSON(t, w, 400, ErrorResponse{APIError{"INVALID_REQUEST", "invalid execution ID"}})
	}
}
func TestLocalStatusHandlerExactDelegation(t *testing.T) {
	reader := &httpStatusReader{snapshot: workflow.NewWorkflowExecution("run"), exists: true}
	s := httpStatusScope(t, reader, "run")
	p := &httpStatusProbe{delegate: s.Querier()}
	h, e := newLocalStatusHandler(s, p)
	if e != nil {
		t.Fatal(e)
	}
	ctx := context.WithValue(context.Background(), struct{}{}, "exact")
	w := sendStatus(statusRouter(h), "run", ctx)
	expected, e := s.Command("run")
	if e != nil {
		t.Fatal(e)
	}
	assertStatusJSON(t, w, 200, ExecutionStatusResponse{"run", "READY", 1})
	if p.calls != 1 || p.ctx != ctx || p.command != expected || reader.calls != 1 {
		t.Fatal("delegation altered")
	}
}
func TestLocalStatusHandlerMissingConcealed(t *testing.T) {
	repo := workflow.NewInMemoryWorkflowRepository()
	repo.Save(workflow.NewWorkflowExecution("hidden"))
	reader := &httpStatusReader{reader: repo}
	s := httpStatusScope(t, reader, "missing")
	h, e := NewLocalStatusHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	r := statusRouter(h)
	missing := sendStatus(r, "missing", nil)
	if reader.calls != 1 {
		t.Fatal("miss must read once")
	}
	for _, id := range []string{"hidden", "unknown"} {
		concealed := sendStatus(r, id, nil)
		assertStatusJSON(t, concealed, 404, ErrorResponse{APIError{"NOT_FOUND", "execution not found"}})
		if missing.Code != concealed.Code || missing.Body.String() != concealed.Body.String() || !reflect.DeepEqual(missing.Header(), concealed.Header()) || reader.calls != 1 {
			t.Fatal("concealment distinguishable or read")
		}
	}
	assertStatusJSON(t, missing, 404, ErrorResponse{APIError{"NOT_FOUND", "execution not found"}})
}
func TestLocalStatusHandlerSafeErrors(t *testing.T) {
	for _, tc := range []struct {
		kind          application.StatusQueryErrorKind
		status        int
		code, message string
	}{
		{application.StatusQueryInvalidCommand, 400, "INVALID_REQUEST", "invalid execution ID"},
		{application.StatusQueryMissingAuthority, 401, "UNAUTHENTICATED", "local capability required"},
		{application.StatusQueryForbidden, 403, "FORBIDDEN", "local status query not permitted"},
		{application.StatusQueryConcealed, 404, "NOT_FOUND", "execution not found"},
		{application.StatusQueryNotFound, 404, "NOT_FOUND", "execution not found"},
		{application.StatusQueryUnavailable, 503, "UNAVAILABLE", "local status query unavailable"},
		{application.StatusQueryInvalidSnapshot, 500, "INTERNAL_ERROR", "internal error"},
		{"unknown", 500, "INTERNAL_ERROR", "internal error"},
	} {
		t.Run(string(tc.kind), func(t *testing.T) {
			s := httpStatusScope(t, &httpStatusReader{}, "run")
			p := &httpStatusProbe{err: fmt.Errorf("secret tenant wrapper: %w", &application.StatusQueryError{Kind: tc.kind, Cause: errors.New("private payment failure")})}
			h, e := newLocalStatusHandler(s, p)
			if e != nil {
				t.Fatal(e)
			}
			assertStatusJSON(t, sendStatus(statusRouter(h), "run", nil), tc.status, ErrorResponse{APIError{tc.code, tc.message}})
		})
	}
	s := httpStatusScope(t, &httpStatusReader{}, "run")
	p := &httpStatusProbe{err: errors.New("secret unexpected")}
	h, e := newLocalStatusHandler(s, p)
	if e != nil {
		t.Fatal(e)
	}
	assertStatusJSON(t, sendStatus(statusRouter(h), "run", nil), 500, ErrorResponse{APIError{"INTERNAL_ERROR", "internal error"}})
}
func TestLocalStatusHandlerAuthorityFailsClosed(t *testing.T) {
	reader := &httpStatusReader{}
	s := httpStatusScope(t, reader, "run")
	other := httpStatusScope(t, reader, "run")
	for _, tc := range []struct {
		q             application.LocalExecutionStatusQuerier
		status        int
		code, message string
	}{
		{&httpStatusProbe{delegate: s.Querier(), strip: true}, 401, "UNAUTHENTICATED", "local capability required"},
		{other.Querier(), 403, "FORBIDDEN", "local status query not permitted"},
	} {
		h, e := newLocalStatusHandler(s, tc.q)
		if e != nil {
			t.Fatal(e)
		}
		w := sendStatus(statusRouter(h), "run", context.WithValue(context.Background(), "tenant_id", "authorized"))
		assertStatusJSON(t, w, tc.status, ErrorResponse{APIError{tc.code, tc.message}})
	}
	if reader.calls != 0 {
		t.Fatal("authority rejection read")
	}
}

type statusFixtureID string

func (id statusFixtureID) NextExecutionID(context.Context) (string, error) { return string(id), nil }
func TestLocalStatusHandlerRealEngine(t *testing.T) {
	for _, fixture := range []string{"ready", "local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(fixture, func(t *testing.T) {
			repo := workflow.NewInMemoryWorkflowRepository()
			id := "exec-123456789012345678901234567890ab"
			if fixture == "ready" {
				if ok, e := repo.CreateIfAbsent(workflow.NewWorkflowExecution(id)); !ok || e != nil {
					t.Fatal(ok, e)
				}
			} else {
				invocation, e := application.NewLocalInvocationScope(repo, statusFixtureID(id), "secret-worker")
				if e != nil {
					t.Fatal(e)
				}
				invokeHandler, e := NewLocalInvocationHandler(invocation)
				if e != nil {
					t.Fatal(e)
				}
				create := sendInvocation(invocationRouter(t, invokeHandler), validationBody(fixture, "CONFIRMED"), nil)
				if create.Code != 201 {
					t.Fatal(create.Code, create.Body)
				}
			}
			before, _ := repo.Get(id)
			reader := &httpStatusReader{reader: repo}
			s := httpStatusScope(t, reader, id)
			h, e := NewLocalStatusHandler(s)
			if e != nil {
				t.Fatal(e)
			}
			r := statusRouter(h)
			state, version := "READY", uint64(1)
			switch fixture {
			case "local.success":
				state, version = "COMPLETED", 3
			case "local.failure":
				state, version = "FAILED", 3
			case "local.wait", "local.condition":
				state, version = "RUNNING", 2
			}
			for i := 0; i < 4; i++ {
				w := sendStatus(r, id, nil)
				assertStatusJSON(t, w, 200, ExecutionStatusResponse{id, state, version})
				var keys map[string]any
				if e := json.Unmarshal(w.Body.Bytes(), &keys); e != nil || len(keys) != 3 {
					t.Fatal("projection keys")
				}
				for _, key := range []string{"execution_id", "status", "version"} {
					if _, ok := keys[key]; !ok {
						t.Fatal("missing key")
					}
				}
				for _, secret := range []string{"secret-worker", "controlled", "LastError", "claim_owner", "history", "attempts", "tenant", "cursor", "CONFIRMED", "authority"} {
					if strings.Contains(w.Body.String(), secret) {
						t.Fatal("disclosure", secret)
					}
				}
			}
			after, _ := repo.Get(id)
			if !reflect.DeepEqual(before, after) || reader.calls != 4 {
				t.Fatal("GET changed metadata or replayed lifecycle")
			}
			if state == "RUNNING" && (len(after.History) != 1 || after.Attempts != 1) {
				t.Fatal("handoff resumed")
			}
		})
	}
}
func TestLocalStatusHandlerCorruptionAndCancellation(t *testing.T) {
	for _, mode := range []string{"nil-present", "contradictory", "wrong-id", "invalid-state", "cancelled"} {
		t.Run(mode, func(t *testing.T) {
			p := &httpStatusReader{snapshot: workflow.NewWorkflowExecution("run"), exists: true}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch mode {
			case "nil-present":
				p.snapshot = nil
			case "contradictory":
				p.exists = false
			case "wrong-id":
				p.snapshot.WorkflowID = "private-other"
			case "invalid-state":
				p.snapshot.Status = "WAITING"
			case "cancelled":
				cancel()
			}
			s := httpStatusScope(t, p, "run")
			h, e := NewLocalStatusHandler(s)
			if e != nil {
				t.Fatal(e)
			}
			w := sendStatus(statusRouter(h), "run", ctx)
			if mode == "cancelled" {
				assertStatusJSON(t, w, 503, ErrorResponse{APIError{"UNAVAILABLE", "local status query unavailable"}})
				if p.calls != 0 {
					t.Fatal("cancel read")
				}
			} else {
				assertStatusJSON(t, w, 500, ErrorResponse{APIError{"INTERNAL_ERROR", "internal error"}})
			}
		})
	}
}

func TestLocalStatusHandlerInvalidProjection(t *testing.T) {
	for _, result := range []application.ExecutionStatusResult{{}, {ExecutionID: "other", Status: "READY", Version: 1}, {ExecutionID: "run", Status: "WAITING", Version: 2}, {ExecutionID: "run", Status: "RUNNING", Version: 0}} {
		s := httpStatusScope(t, &httpStatusReader{}, "run")
		p := &httpStatusProbe{result: result}
		h, e := newLocalStatusHandler(s, p)
		if e != nil {
			t.Fatal(e)
		}
		assertStatusJSON(t, sendStatus(statusRouter(h), "run", nil), 500, ErrorResponse{APIError{"INTERNAL_ERROR", "internal error"}})
	}
}

type httpStatusPanicReader struct{}

func (httpStatusPanicReader) Get(string) (*workflow.WorkflowExecution, bool) {
	panic("private collaborator panic")
}
func TestLocalStatusHandlerRecovery(t *testing.T) {
	// Explicit containment is omitted. Gin recovery can handle a collaborator
	// panic before JSON is written; it supplies no JSON error-envelope guarantee.
	s := httpStatusScope(t, httpStatusPanicReader{}, "run")
	h, e := NewLocalStatusHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	r := gin.New()
	r.Use(gin.RecoveryWithWriter(io.Discard))
	r.GET("/v1/workflow-executions/:execution_id", h)
	w := sendStatus(r, "run", nil)
	if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || strings.Contains(w.Body.String(), "private") {
		t.Fatal("panic produced success or disclosed payload")
	}
}
