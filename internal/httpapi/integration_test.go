package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"example.com/p1-reliable-academy-workflows/internal/httpapi"
	"github.com/gin-gonic/gin"
)

const (
	m48ExecutionID = "exec-00000000000000000000000000000001"
	m48MissingID   = "exec-00000000000000000000000000000002"
	m48Worker      = "m48-local-worker"
	m48Key         = "m48-request"
	m48Event       = "correlation-only"
)

// Only identity generation is controlled. The public scopes supply the real
// fixture catalog, invocation service and runner over the real repository.
type m48IDs struct{ calls atomic.Uint64 }

func (s *m48IDs) NextExecutionID(context.Context) (string, error) {
	return fmt.Sprintf("exec-%032x", s.calls.Add(1)), nil
}

type m48Fixture struct {
	router *gin.Engine
	repo   *workflow.InMemoryWorkflowRepository
	ids    *m48IDs
}

func newM48Fixture(t *testing.T) *m48Fixture {
	t.Helper()
	f := &m48Fixture{gin.New(), workflow.NewInMemoryWorkflowRepository(), &m48IDs{}}
	scope, err := application.NewLocalRequestIdentityScope(f.repo, f.ids, m48Worker, 1)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := httpapi.NewLocalRequestIdentityHandler(scope)
	if err != nil {
		t.Fatal(err)
	}
	f.router.POST("/test-only/invoke", handler)
	return f
}

// This observer has only the existing read boundary, not lifecycle capability.
type m48Reader struct {
	repo  *workflow.InMemoryWorkflowRepository
	calls atomic.Uint64
}

func (r *m48Reader) Get(id string) (*workflow.WorkflowExecution, bool) {
	r.calls.Add(1)
	return r.repo.Get(id)
}

func (f *m48Fixture) statusRoute(t *testing.T, path string, visible []string) *m48Reader {
	t.Helper()
	reader := &m48Reader{repo: f.repo}
	scope, err := application.NewLocalStatusScope(reader, visible)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := httpapi.NewLocalStatusHandler(scope)
	if err != nil {
		t.Fatal(err)
	}
	f.router.GET(path+"/:execution_id", handler)
	if reader.calls.Load() != 0 {
		t.Fatal("status composition read repository")
	}
	return reader
}

func m48Body(t *testing.T, definition, event string) string {
	t.Helper()
	body, err := json.Marshal(httpapi.CreateExecutionRequest{WorkflowDefinitionID: definition, EventID: event})
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func m48Header() http.Header { return http.Header{"Idempotency-Key": []string{m48Key}} }

func m48Request(router *gin.Engine, method, path, body string, header http.Header) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if header != nil {
		request.Header = header.Clone()
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	return w
}

// Public response expectations are owned by this test, independently of DTOs.
func m48Created(id, status string) map[string]any {
	return map[string]any{"execution_id": id, "status": status}
}

func m48Status(id, status string, version uint64) map[string]any {
	return map[string]any{"execution_id": id, "status": status, "version": json.Number(strconv.FormatUint(version, 10))}
}

func m48ErrorValue(code, message string) map[string]any {
	return map[string]any{"error": map[string]any{"code": code, "message": message}}
}

func m48ParseJSON(body []byte) (map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	value, err := m48JSONValue(decoder)
	if err != nil {
		return nil, err
	}
	object, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("response must be a JSON object")
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("response must end after one JSON object: %v", err)
	}
	return object, nil
}

// Token returns decoded key names, so escaped aliases share an object's set.
// Each recursive object owns a separate set; arrays recurse into every element.
func m48JSONValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delimiter, compound := token.(json.Delim)
	if !compound {
		return token, nil // UseNumber retains exact integer version values.
	}
	switch delimiter {
	case '{':
		object := make(map[string]any)
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, fmt.Errorf("object key must be a string")
			}
			if _, duplicate := object[key]; duplicate {
				return nil, fmt.Errorf("duplicate JSON object key %q", key)
			}
			value, err := m48JSONValue(decoder)
			if err != nil {
				return nil, err
			}
			object[key] = value
		}
		if closing, err := decoder.Token(); err != nil || closing != json.Delim('}') {
			return nil, fmt.Errorf("invalid object closing token %v: %v", closing, err)
		}
		return object, nil
	case '[':
		array := make([]any, 0)
		for decoder.More() {
			value, err := m48JSONValue(decoder)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}
		if closing, err := decoder.Token(); err != nil || closing != json.Delim(']') {
			return nil, fmt.Errorf("invalid array closing token %v: %v", closing, err)
		}
		return array, nil
	default:
		return nil, fmt.Errorf("unexpected JSON delimiter %q", delimiter)
	}
}

func m48MatchJSON(body []byte, want map[string]any) error {
	actual, err := m48ParseJSON(body)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(actual, want) {
		return fmt.Errorf("public JSON = %s, want exactly %v", body, want)
	}
	return nil
}

func m48JSON(t *testing.T, w *httptest.ResponseRecorder, status int, want map[string]any) {
	t.Helper()
	if w.Code != status || w.Header().Get("Content-Type") != "application/json; charset=utf-8" || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("HTTP %d headers=%v body=%s; want %d JSON/no-store", w.Code, w.Header(), w.Body, status)
	}
	// Compare the parsed response with fixed test-owned keys and values,
	// including nested error fields; duplicate members must be rejected first.
	if err := m48MatchJSON(w.Body.Bytes(), want); err != nil {
		t.Fatal(err)
	}
}

func m48Error(t *testing.T, w *httptest.ResponseRecorder, status int, code, message string) {
	t.Helper()
	m48JSON(t, w, status, m48ErrorValue(code, message))
}

func m48Snapshot(t *testing.T, repo *workflow.InMemoryWorkflowRepository) *workflow.WorkflowExecution {
	t.Helper()
	snapshot, exists := repo.Get(m48ExecutionID)
	if !exists || snapshot == nil {
		t.Fatal("expected authoritative execution")
	}
	return snapshot
}

func m48State(t *testing.T, repo *workflow.InMemoryWorkflowRepository, status string, version uint64, reason string) *workflow.WorkflowExecution {
	t.Helper()
	snapshot := m48Snapshot(t, repo)
	history := []workflow.WorkflowTransition{{From: "READY", To: "RUNNING", WorkerID: m48Worker, Version: 2}}
	if status == "COMPLETED" || status == "FAILED" {
		history = append(history, workflow.WorkflowTransition{From: "RUNNING", To: status, WorkerID: m48Worker, Version: 3, Reason: reason})
	}
	if snapshot.WorkflowID != m48ExecutionID || snapshot.Status != status || snapshot.Version != version || snapshot.ClaimOwner != m48Worker || snapshot.Attempts != 1 || snapshot.LastError != reason || !reflect.DeepEqual(snapshot.History, history) {
		t.Fatalf("stored execution=%+v history=%+v; want %s/v%d owner=%s attempts=1 reason=%q history=%+v", snapshot, snapshot.History, status, version, m48Worker, reason, history)
	}
	return snapshot
}

func m48Unchanged(t *testing.T, f *m48Fixture, before *workflow.WorkflowExecution) {
	t.Helper()
	if after := m48Snapshot(t, f.repo); !reflect.DeepEqual(before, after) {
		t.Fatalf("execution changed: before=%+v after=%+v", before, after)
	}
	if f.ids.calls.Load() != 1 {
		t.Fatal("unexpected execution-ID allocation", f.ids.calls.Load())
	}
	if _, exists := f.repo.Get(m48MissingID); exists {
		t.Fatal("unexpected second execution")
	}
}

func TestM48FixtureLifecycle(t *testing.T) {
	for _, tc := range []struct {
		definition, status string
		version            uint64
	}{
		{"local.success", "COMPLETED", 3},
		{"local.failure", "FAILED", 3},
		{"local.wait", "RUNNING", 2},
		{"local.condition", "RUNNING", 2},
	} {
		t.Run(tc.definition, func(t *testing.T) {
			f := newM48Fixture(t)
			// Visibility is an explicit trusted-host decision made independently
			// before POST; receiving a creation response does not authorize GET.
			reader := f.statusRoute(t, "/test-only/status", []string{m48ExecutionID})
			created := m48Request(f.router, "POST", "/test-only/invoke", m48Body(t, tc.definition, m48Event), m48Header())
			m48JSON(t, created, 201, m48Created(m48ExecutionID, tc.status))
			reason := ""
			if tc.status == "FAILED" {
				// Derived from localDefinitionResolver's local.fail error and the
				// ExecuteActions/RunActions wrappers; not from the returned snapshot.
				reason = fmt.Sprintf("workflow execution %q: action step %q: controlled local fixture failure", m48ExecutionID, "start")
			}
			before := m48State(t, f.repo, tc.status, tc.version, reason)
			for i := 0; i < 3; i++ {
				status := m48Request(f.router, "GET", "/test-only/status/"+m48ExecutionID, "", nil)
				m48JSON(t, status, 200, m48Status(m48ExecutionID, tc.status, tc.version))
				m48Unchanged(t, f, before)
			}
			if reader.calls.Load() != 3 {
				t.Fatal("status should read exactly once per request")
			}
			// RUNNING/v2 with only the claim remains stable: GET does not resume
			// WAIT/condition. Callback non-evaluation is additionally covered by
			// existing application handoff tests; no private seam is added here.
		})
	}
}

func TestM48RequestIdentityConflicts(t *testing.T) {
	f := newM48Fixture(t)
	body := m48Body(t, "local.success", m48Event)
	m48JSON(t, m48Request(f.router, "POST", "/test-only/invoke", body, m48Header()), 201, m48Created(m48ExecutionID, "COMPLETED"))
	before := m48State(t, f.repo, "COMPLETED", 3, "")
	for _, tc := range []struct{ name, body string }{
		{"identical", body},
		{"equivalent-json", `{ "event_id":"correlation-only", "workflow_definition_id":"local.\u0073uccess" }`},
		{"changed-event", m48Body(t, "local.success", "different-event")},
		{"changed-definition", m48Body(t, "local.failure", m48Event)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m48Error(t, m48Request(f.router, "POST", "/test-only/invoke", tc.body, m48Header()), 409, "CONFLICT", "execution conflict")
			m48Unchanged(t, f, before)
		})
	}
}

func TestM48ValidationBeforeEffects(t *testing.T) {
	valid := m48Body(t, "local.success", m48Event)
	for _, tc := range []struct {
		name, body, message string
		header              http.Header
	}{
		{"malformed-json", `{"event_id":`, "body must be one JSON object", m48Header()},
		{"unknown-field", `{"workflow_definition_id":"local.success","event_id":"correlation-only","payment_status":"CONFIRMED"}`, "unsupported request field", m48Header()},
		{"duplicate-field", `{"workflow_definition_id":"local.success","event_id":"one","event_\u0069d":"two"}`, "duplicate request field", m48Header()},
		{"missing-field", `{"workflow_definition_id":"local.success"}`, "event_id is required", m48Header()},
		{"invalid-identifier", m48Body(t, "local.success", "bad/event"), "event_id must contain 1 to 128 ASCII letters, digits, underscores, hyphens or periods", m48Header()},
		{"unsupported-definition", m48Body(t, "not-approved", m48Event), "unsupported workflow definition", m48Header()},
		{"missing-key", valid, "invalid idempotency key", nil},
		{"repeated-key", valid, "invalid idempotency key", http.Header{"Idempotency-Key": []string{m48Key, m48Key}}},
		{"malformed-key", valid, "invalid idempotency key", http.Header{"Idempotency-Key": []string{"bad/key"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newM48Fixture(t)
			m48Error(t, m48Request(f.router, "POST", "/test-only/invoke", tc.body, tc.header), 400, "INVALID_REQUEST", tc.message)
			if f.ids.calls.Load() != 0 {
				t.Fatal("rejected request entered identity generation")
			}
			if _, exists := f.repo.Get(m48ExecutionID); exists {
				t.Fatal("rejected request inserted execution")
			}
			// Capacity one proves rejection did not reserve the key or consume
			// admission capacity. For invalid/missing headers, supply a valid key.
			m48JSON(t, m48Request(f.router, "POST", "/test-only/invoke", valid, m48Header()), 201, m48Created(m48ExecutionID, "COMPLETED"))
			before := m48State(t, f.repo, "COMPLETED", 3, "")
			m48Unchanged(t, f, before)
		})
	}
}

func TestM48IndependentStatusVisibility(t *testing.T) {
	f := newM48Fixture(t)
	denied := f.statusRoute(t, "/test-only/concealed", nil)
	allowed := f.statusRoute(t, "/test-only/status", []string{m48ExecutionID, m48MissingID})
	m48JSON(t, m48Request(f.router, "POST", "/test-only/invoke", m48Body(t, "local.success", m48Event), m48Header()), 201, m48Created(m48ExecutionID, "COMPLETED"))
	before := m48State(t, f.repo, "COMPLETED", 3, "")
	concealed := m48Request(f.router, "GET", "/test-only/concealed/"+m48ExecutionID, "", nil)
	missing := m48Request(f.router, "GET", "/test-only/status/"+m48MissingID, "", nil)
	for _, w := range []*httptest.ResponseRecorder{concealed, missing} {
		m48Error(t, w, 404, "NOT_FOUND", "execution not found")
	}
	if concealed.Body.String() != missing.Body.String() || !reflect.DeepEqual(concealed.Header(), missing.Header()) || denied.calls.Load() != 0 || allowed.calls.Load() != 1 {
		t.Fatal("concealment differed from allowed missing response or read storage")
	}
	m48JSON(t, m48Request(f.router, "GET", "/test-only/status/"+m48ExecutionID, "", nil), 200, m48Status(m48ExecutionID, "COMPLETED", 3))
	if allowed.calls.Load() != 2 {
		t.Fatal("authorized queries should read once each")
	}
	m48Unchanged(t, f, before)
}

func TestM48ProductionRouteIsolation(t *testing.T) {
	router := httpapi.NewRouter()
	routes := router.Routes()
	if len(routes) != 1 || routes[0].Method != "GET" || routes[0].Path != "/health" {
		t.Fatalf("production routes=%+v", routes)
	}
	health := m48Request(router, "GET", "/health", "", nil)
	if health.Code != 200 || health.Body.String() != `{"status":"ok"}` || health.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("health contract changed: %d %v %s", health.Code, health.Header(), health.Body)
	}
	for _, path := range []string{
		"/v1/workflow-executions", "/v1/workflow-executions/" + m48ExecutionID,
		"/test-only/invoke", "/test-only/status/" + m48ExecutionID,
		"/test-only/concealed/" + m48ExecutionID,
		"/payment", "/payment/callback", "/enrollment", "/entitlements", "/admin",
	} {
		for _, method := range []string{"GET", "POST"} {
			w := m48Request(router, method, path, "", nil)
			if w.Code != 404 || w.Body.String() != "404 page not found" {
				t.Fatalf("production %s %s: %d %s", method, path, w.Code, w.Body)
			}
		}
	}
}

func TestM48StrictJSONAssertions(t *testing.T) {
	created := m48Created("execution-1", "COMPLETED")
	status := m48Status("execution-1", "COMPLETED", 3)
	failure := m48ErrorValue("CONFLICT", "execution conflict")
	for _, tc := range []struct {
		name, body string
		want       map[string]any
	}{
		{"unexpected-creation-field", `{"execution_id":"execution-1","status":"COMPLETED","owner":"private"}`, created},
		{"unexpected-status-field", `{"execution_id":"execution-1","status":"COMPLETED","version":3,"history":[]}`, status},
		{"unexpected-error-field", `{"error":{"code":"CONFLICT","message":"execution conflict"},"debug":true}`, failure},
		{"unexpected-nested-error-field", `{"error":{"code":"CONFLICT","message":"execution conflict","cause":"private"}}`, failure},
		{"missing-creation-field", `{"execution_id":"execution-1"}`, created},
		{"missing-status-field", `{"execution_id":"execution-1","status":"COMPLETED"}`, status},
		{"missing-error-field", `{}`, failure},
		{"missing-nested-error-field", `{"error":{"code":"CONFLICT"}}`, failure},
		{"wrong-value", `{"execution_id":"execution-1","status":"FAILED"}`, created},
		{"wrong-version", `{"execution_id":"execution-1","status":"COMPLETED","version":4}`, status},
		{"wrong-version-type", `{"execution_id":"execution-1","status":"COMPLETED","version":"3"}`, status},
		{"large-version-mismatch", `{"execution_id":"execution-1","status":"COMPLETED","version":18446744073709551614}`, m48Status("execution-1", "COMPLETED", ^uint64(0))},
		{"duplicate-top-level", `{"execution_id":"execution-1","status":"private diagnostic","status":"COMPLETED"}`, created},
		{"duplicate-nested-error", `{"error":{"code":"CONFLICT","message":"private","message":"execution conflict"}}`, failure},
		{"duplicate-escaped-name", `{"execution_id":"execution-1","status":"private","statu\u0073":"COMPLETED"}`, created},
		{"duplicate-nested-escaped-name", `{"error":{"code":"private","co\u0064e":"CONFLICT","message":"execution conflict"}}`, failure},
		{"duplicate-in-array-object", `{"items":[{"code":"private","code":"ok"}]}`, map[string]any{"items": []any{map[string]any{"code": "ok"}}}},
		{"malformed-json", `{"execution_id":`, created},
		{"malformed-array", `{"items":[1,]}`, map[string]any{"items": []any{json.Number("1")}}},
		{"mismatched-closing-delimiter", `{"items":[1}}`, map[string]any{"items": []any{json.Number("1")}}},
		{"trailing-second-value", `{"execution_id":"execution-1","status":"COMPLETED"}{}`, created},
		{"trailing-garbage", `{"execution_id":"execution-1","status":"COMPLETED"}garbage`, created},
		{"non-object-array", `[]`, created},
		{"non-object-string", `"text"`, created},
		{"non-object-null", `null`, created},
	} {
		t.Run("reject/"+tc.name, func(t *testing.T) {
			if err := m48MatchJSON([]byte(tc.body), tc.want); err == nil {
				t.Fatalf("accepted invalid response: %s", tc.body)
			}
		})
	}
	for _, tc := range []struct {
		name, body string
		want       map[string]any
	}{
		{"creation", `{"execution_id":"execution-1","status":"COMPLETED"}`, created},
		{"status", `{"execution_id":"execution-1","status":"COMPLETED","version":3}`, status},
		{"error", `{"error":{"code":"CONFLICT","message":"execution conflict"}}`, failure},
		{"reordered-creation", `{"status":"COMPLETED","execution_id":"execution-1"}`, created},
		{"reordered-error", `{"error":{"message":"execution conflict","code":"CONFLICT"}}`, failure},
		{"escaped-unique-name", `{"execution_id":"execution-1","statu\u0073":"COMPLETED"}`, created},
		{"distinct-nested-objects", `{"left":{"code":"a"},"right":{"code":"b"}}`, map[string]any{"left": map[string]any{"code": "a"}, "right": map[string]any{"code": "b"}}},
		{"array-objects-and-scalars", `{"items":[{"code":"a"},{"code":"b"},[null,true,"text",3]]}`, map[string]any{"items": []any{map[string]any{"code": "a"}, map[string]any{"code": "b"}, []any{nil, true, "text", json.Number("3")}}}},
		{"large-version", `{"execution_id":"execution-1","status":"COMPLETED","version":18446744073709551615}`, m48Status("execution-1", "COMPLETED", ^uint64(0))},
	} {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			if err := m48MatchJSON([]byte(tc.body), tc.want); err != nil {
				t.Fatal(err)
			}
			actual, err := m48ParseJSON([]byte(tc.body))
			if err != nil || !reflect.DeepEqual(actual, tc.want) {
				t.Fatalf("parsed=%v err=%v want=%v", actual, err, tc.want)
			}
		})
	}
}
