package httpapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

type preparationProbe struct {
	calls   int
	command application.SubmissionCommand
	ctx     context.Context
	err     error
}

func (p *preparationProbe) PrepareSubmission(ctx context.Context, command application.SubmissionCommand) (application.PreparationResult, error) {
	p.calls++
	p.command = command
	p.ctx = ctx
	return application.PreparationResult{WorkflowDefinitionID: command.WorkflowDefinitionID, EventID: command.EventID, StepCount: 7}, p.err
}

func preparationRouter(t *testing.T, service application.WorkflowService, approved map[string]struct{}) *gin.Engine {
	t.Helper()
	handler, err := NewPreparationHandler(service, approved)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/test-only/prepare", handler)
	return router
}

func TestPreparationHandlerConstruction(t *testing.T) {
	var typedNil *preparationProbe
	for _, service := range []application.WorkflowService{nil, typedNil} {
		if handler, err := NewPreparationHandler(service, nil); handler != nil || err == nil {
			t.Fatal("missing service accepted")
		}
	}
	p := &preparationProbe{}
	approved := map[string]struct{}{"definition-1": {}}
	router := preparationRouter(t, p, approved)
	if p.calls != 0 {
		t.Fatal("construction called service")
	}
	delete(approved, "definition-1") // The handler owns its support-set snapshot.
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("POST", "/test-only/prepare", strings.NewReader(validationBody("definition-1", "event-1"))))
	if response.Code != 200 || p.calls != 1 {
		t.Fatalf("support snapshot failed: %d %d", response.Code, p.calls)
	}
}

func TestPreparationHandlerRejectsBeforeDelegation(t *testing.T) {
	for _, tc := range invalidValidationCases() {
		t.Run(tc.name, func(t *testing.T) {
			p := &preparationProbe{}
			router := preparationRouter(t, p, tc.approved)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest("POST", "/test-only/prepare", strings.NewReader(tc.body)))
			if response.Code != 400 || p.calls != 0 {
				t.Fatalf("status=%d calls=%d", response.Code, p.calls)
			}
			assertPreparationJSON(t, response, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: tc.message}})
		})
	}
	for _, field := range []string{"tenant_id", "caller_id", "payment_status"} {
		p := &preparationProbe{}
		r := preparationRouter(t, p, map[string]struct{}{"definition-1": {}})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/test-only/prepare", strings.NewReader(`{"workflow_definition_id":"definition-1","event_id":"event-1","`+field+`":"CONFIRMED"}`)))
		if w.Code != 400 || p.calls != 0 {
			t.Fatalf("untrusted %s accepted", field)
		}
	}
}

type preparationBody struct {
	read int
	fail bool
}

func (b *preparationBody) Read(p []byte) (int, error) {
	if b.fail {
		return 0, errors.New("private body diagnostic")
	}
	for i := range p {
		p[i] = ' '
	}
	b.read += len(p)
	return len(p), nil // A never-ending stream must still be bounded.
}
func (*preparationBody) Close() error { return nil }

func TestPreparationHandlerBodyReads(t *testing.T) {
	for _, fail := range []bool{false, true} {
		p := &preparationProbe{}
		r := preparationRouter(t, p, map[string]struct{}{"definition-1": {}})
		body := &preparationBody{fail: fail}
		req := httptest.NewRequest("POST", "/test-only/prepare", nil)
		req.Body = body
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 400 || p.calls != 0 {
			t.Fatalf("read failure/size reached service: %d %d", w.Code, p.calls)
		}
		message := "body must not exceed 4096 bytes"
		if fail {
			message = "could not read request body"
		} else if body.read != 4097 {
			t.Fatalf("read %d bytes", body.read)
		}
		assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: message}})
	}
	// Exactly the documented byte limit remains accepted through the real adapter.
	p := &preparationProbe{}
	r := preparationRouter(t, p, map[string]struct{}{"definition-1": {}})
	valid := validationBody("definition-1", "event-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/test-only/prepare", io.NopCloser(strings.NewReader(valid+strings.Repeat(" ", 4096-len(valid))))))
	if w.Code != 200 || p.calls != 1 {
		t.Fatalf("4096-byte request rejected: %d %d", w.Code, p.calls)
	}
}

func assertPreparationJSON(t *testing.T, w *httptest.ResponseRecorder, want any) {
	t.Helper()
	expected := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(expected)
	c.JSON(w.Code, want)
	if w.Body.String() != expected.Body.String() || w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("response=%s expected=%s", w.Body, expected.Body)
	}
}

func TestPreparationHandlerDelegation(t *testing.T) {
	p := &preparationProbe{}
	r := preparationRouter(t, p, map[string]struct{}{"definition-1": {}})
	ctx := context.WithValue(context.Background(), struct{}{}, "test-context")
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test-only/prepare", strings.NewReader(validationBody("definition-1", "event-2"))).WithContext(ctx)
	r.ServeHTTP(w, req)
	if w.Code != 200 || p.calls != 1 || p.ctx != ctx || p.command != (application.SubmissionCommand{WorkflowDefinitionID: "definition-1", EventID: "event-2"}) {
		t.Fatalf("delegation: status=%d probe=%+v", w.Code, p)
	}
	if w.Body.String() != `{"workflow_definition_id":"definition-1","event_id":"event-2","step_count":7}` {
		t.Fatalf("result=%s", w.Body)
	}
}

func TestPreparationHandlerErrors(t *testing.T) {
	for _, tc := range []struct {
		err           error
		status        int
		code, message string
	}{
		{application.ErrDefinitionNotFound, 404, "NOT_FOUND", "workflow definition not found"},
		{application.ErrUnavailable, 503, "UNAVAILABLE", "preparation service unavailable"},
		{application.ErrInvalidCommand, 400, "INVALID_REQUEST", "invalid preparation command"},
		{application.ErrInvalidDefinition, 500, "INTERNAL_ERROR", "internal error"},
		{errors.New("tenant=private payment=CONFIRMED owner=secret"), 500, "INTERNAL_ERROR", "internal error"},
	} {
		t.Run(tc.code+tc.err.Error(), func(t *testing.T) {
			p := &preparationProbe{err: fmt.Errorf("private provider detail: %w", tc.err)}
			r := preparationRouter(t, p, map[string]struct{}{"definition-1": {}})
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest("POST", "/test-only/prepare", strings.NewReader(validationBody("definition-1", "event-1"))))
			if w.Code != tc.status || p.calls != 1 {
				t.Fatalf("status=%d calls=%d", w.Code, p.calls)
			}
			assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: tc.code, Message: tc.message}})
		})
	}
}

type controlledDefinitions struct {
	definition workflow.WorkflowDefinition
	calls      int
}

func (d *controlledDefinitions) LookupDefinition(context.Context, string) (workflow.WorkflowDefinition, error) {
	d.calls++
	return d.definition, nil
}

func TestPreparationHandlerRealService(t *testing.T) {
	// Graph validity changes the real HTTP outcome despite identical valid JSON.
	// There is no registry/provider/repository capability in this composition.
	definitions := &controlledDefinitions{definition: workflow.WorkflowDefinition{ID: "definition-1", Trigger: workflow.TriggerDefinition{EventType: "purchase"}, StartStep: "grant", Steps: map[string]workflow.StepDefinition{"grant": {Type: workflow.StepAction, Action: "entitlement.grant"}}}}
	service, err := application.NewService(definitions)
	if err != nil {
		t.Fatal(err)
	}
	r := preparationRouter(t, service, map[string]struct{}{"definition-1": {}})
	for _, valid := range []bool{true, false} {
		if !valid {
			definitions.definition.StartStep = "missing"
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/test-only/prepare", strings.NewReader(validationBody("definition-1", "CONFIRMED"))))
		if valid {
			if w.Code != 200 || w.Body.String() != `{"workflow_definition_id":"definition-1","event_id":"CONFIRMED","step_count":1}` {
				t.Fatalf("preparation=%d %s", w.Code, w.Body)
			}
		} else {
			if w.Code != 500 || strings.Contains(w.Body.String(), "missing") {
				t.Fatalf("invalid graph=%d %s", w.Code, w.Body)
			}
		}
	}
	if definitions.calls != 2 {
		t.Fatalf("lookup calls=%d", definitions.calls)
	}
}
