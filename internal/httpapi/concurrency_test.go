package httpapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"example.com/p1-reliable-academy-workflows/internal/httpapi"
	"github.com/gin-gonic/gin"
)

// Reuse M4.8's public composition parts and strict test-owned JSON assertions.
// Only capacity varies; the real resolver, runner and repository remain intact.
func newM49Fixture(t *testing.T, capacity int) *m48Fixture {
	t.Helper()
	f := &m48Fixture{gin.New(), workflow.NewInMemoryWorkflowRepository(), &m48IDs{}}
	scope, err := application.NewLocalRequestIdentityScope(f.repo, f.ids, m48Worker, capacity)
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

type m49Request struct {
	method, path, body string
	header             http.Header
}

func m49Post(t *testing.T, key, definition, event string) m49Request {
	t.Helper()
	return m49Request{"POST", "/test-only/invoke", m48Body(t, definition, event), http.Header{"Idempotency-Key": []string{key}}}
}

// All requests and routes are prepared before starting goroutines. Each worker
// owns a recorder and result slot. The WaitGroup joins every worker before any
// assertions run on the main test goroutine. The start channel increases
// contention; it does not prove simultaneous arrival at the reservation mutex.
// These bounded synchronous fixtures have no listeners, external I/O, sleeps,
// injected blocking dependencies or WAIT timers (WAIT returns a handoff).
// The standard go test timeout bounds a regression that deadlocks internally.
func m49Race(router *gin.Engine, requests []m49Request) []*httptest.ResponseRecorder {
	results := make([]*httptest.ResponseRecorder, len(requests))
	start := make(chan struct{})
	var ready, done sync.WaitGroup
	ready.Add(len(requests))
	done.Add(len(requests))
	for i, request := range requests {
		go func() {
			defer done.Done()
			ready.Done()
			<-start
			results[i] = m48Request(router, request.method, request.path, request.body, request.header)
		}()
	}
	ready.Wait()
	close(start)
	done.Wait()
	return results
}

func m49ID(n int) string { return fmt.Sprintf("exec-%032x", n) }

func m49Expected(definition, id string) (string, uint64, string) {
	switch definition {
	case "local.failure":
		return "FAILED", 3, fmt.Sprintf("workflow execution %q: action step %q: controlled local fixture failure", id, "start")
	case "local.wait", "local.condition":
		return "RUNNING", 2, ""
	default:
		return "COMPLETED", 3, ""
	}
}

func m49State(t *testing.T, f *m48Fixture, id, definition string) *workflow.WorkflowExecution {
	t.Helper()
	status, version, reason := m49Expected(definition, id)
	history := []workflow.WorkflowTransition{{From: "READY", To: "RUNNING", WorkerID: m48Worker, Version: 2}}
	if version == 3 {
		history = append(history, workflow.WorkflowTransition{From: "RUNNING", To: status, WorkerID: m48Worker, Version: 3, Reason: reason})
	}
	s, exists := f.repo.Get(id)
	if !exists || s == nil {
		t.Fatalf("missing execution %s", id)
	}
	if s.WorkflowID != id || s.Status != status || s.Version != version || s.ClaimOwner != m48Worker || s.Attempts != 1 || s.LastError != reason || !reflect.DeepEqual(s.History, history) {
		t.Fatalf("execution=%+v history=%+v; want %s/v%d owner=%s attempts=1 reason=%q history=%+v", s, s.History, status, version, m48Worker, reason, history)
	}
	return s
}

func m49Count(t *testing.T, f *m48Fixture, count int) {
	t.Helper()
	if got := f.ids.calls.Load(); got != uint64(count) {
		t.Fatalf("ID allocations=%d want=%d", got, count)
	}
	// The controlled ID source is the only source of IDs on this public path.
	for n := 1; n <= count; n++ {
		if _, exists := f.repo.Get(m49ID(n)); !exists {
			t.Fatalf("lost execution %d", n)
		}
	}
	if _, exists := f.repo.Get(m49ID(count + 1)); exists {
		t.Fatal("unexpected additional execution")
	}
}

func TestM49SameKey(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			f := newM49Fixture(t, 1)
			requests := make([]m49Request, 20)
			for i := range requests {
				requests[i] = m49Post(t, "same-key", definition, "same-event")
			}
			admitted := 0
			status, _, _ := m49Expected(definition, m49ID(1))
			for _, w := range m49Race(f.router, requests) {
				if w.Code == 201 {
					admitted++
					m48JSON(t, w, 201, m48Created(m49ID(1), status))
				} else {
					m48Error(t, w, 409, "CONFLICT", "execution conflict")
				}
			}
			if admitted != 1 {
				t.Fatalf("admitted=%d want=1", admitted)
			}
			m49Count(t, f, 1)
			m49State(t, f, m49ID(1), definition)
		})
	}
}

func TestM49ConflictingSelections(t *testing.T) {
	f := newM49Fixture(t, 1)
	requests := make([]m49Request, 20)
	definitions := []string{"local.success", "local.failure"}
	for i := range requests {
		requests[i] = m49Post(t, "same-key", definitions[i%2], fmt.Sprintf("event-%d", i))
	}
	winner := -1
	for i, w := range m49Race(f.router, requests) {
		if w.Code != 201 {
			m48Error(t, w, 409, "CONFLICT", "execution conflict")
			continue
		}
		if winner != -1 {
			t.Fatal("multiple admissions for one key")
		}
		winner = i
		status, _, _ := m49Expected(definitions[i%2], m49ID(1))
		m48JSON(t, w, 201, m48Created(m49ID(1), status))
	}
	if winner == -1 {
		t.Fatal("no admission")
	}
	m49Count(t, f, 1)
	// The repository does not persist definition/event associations. Distinct
	// success/failure lifecycles prove the observed winner's definition here;
	// exact event tuple conflict checks remain in existing application tests.
	m49State(t, f, m49ID(1), definitions[winner%2])
}

func TestM49IndependentKeys(t *testing.T) {
	m49DistinctKeys(t, 12, 12)
}

func TestM49BoundedCapacity(t *testing.T) {
	m49DistinctKeys(t, 20, 3)
}

func m49DistinctKeys(t *testing.T, count, capacity int) {
	t.Helper()
	f := newM49Fixture(t, capacity)
	definitions := []string{"local.success", "local.failure", "local.wait", "local.condition"}
	requests := make([]m49Request, count)
	for i := range requests {
		requests[i] = m49Post(t, fmt.Sprintf("key-%d", i), definitions[i%4], "same-event")
	}
	seen := make(map[string]bool)
	for i, w := range m49Race(f.router, requests) {
		if w.Code != 201 {
			m48Error(t, w, 503, "UNAVAILABLE", "local invocation unavailable")
			continue
		}
		object, err := m48ParseJSON(w.Body.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		id, ok := object["execution_id"].(string)
		if !ok || seen[id] {
			t.Fatalf("missing or duplicate admitted ID: %v", object)
		}
		seen[id] = true
		status, _, _ := m49Expected(definitions[i%4], id)
		m48JSON(t, w, 201, m48Created(id, status))
		m49State(t, f, id, definitions[i%4])
	}
	if len(seen) != capacity {
		t.Fatalf("admitted=%d want=%d", len(seen), capacity)
	}
	m49Count(t, f, capacity)
	for n := 1; n <= capacity; n++ {
		if !seen[m49ID(n)] {
			t.Fatalf("stored ID %d has no admitted response", n)
		}
	}
	// Reservations are not evicted on completion; capacity stays bounded.
	w := m48Request(f.router, "POST", "/test-only/invoke", requests[0].body, http.Header{"Idempotency-Key": []string{"new-key"}})
	m48Error(t, w, 503, "UNAVAILABLE", "local invocation unavailable")
	m49Count(t, f, capacity)
}

func TestM49ValidationBeforeEffects(t *testing.T) {
	f := newM49Fixture(t, 1)
	valid := m49Post(t, "valid-key", "local.success", "valid-event")
	invalid := []struct {
		body, message string
		header        http.Header
	}{
		{`{"workflow_definition_id":`, "body must be one JSON object", valid.header},
		{`{"workflow_definition_id":"local.success","event_id":"one","event_\u0069d":"two"}`, "duplicate request field", valid.header},
		{m48Body(t, "unsupported", "valid-event"), "unsupported workflow definition", valid.header},
		{`{"workflow_definition_id":"local.success","event_id":"valid-event","tenant_id":"foreign"}`, "unsupported request field", valid.header},
		{`{"workflow_definition_id":"local.success","event_id":"valid-event","authority":{}}`, "unsupported request field", valid.header},
		{`{"workflow_definition_id":"local.success","event_id":"valid-event","idempotency_key":"valid-key"}`, "unsupported request field", valid.header},
		{valid.body, "invalid idempotency key", nil},
		{valid.body, "invalid idempotency key", http.Header{"Idempotency-Key": []string{"bad/key"}}},
		{valid.body, "invalid idempotency key", http.Header{"Idempotency-Key": []string{"valid-key", "valid-key"}}},
		{valid.body, "invalid idempotency key", http.Header{"Idempotency-Key": []string{"valid-key"}, "idempotency-key": []string{"valid-key"}}},
	}
	requests := []m49Request{valid}
	for _, tc := range invalid {
		requests = append(requests, m49Request{"POST", valid.path, tc.body, tc.header})
	}
	results := m49Race(f.router, requests)
	m48JSON(t, results[0], 201, m48Created(m49ID(1), "COMPLETED"))
	for i, tc := range invalid {
		m48Error(t, results[i+1], 400, "INVALID_REQUEST", tc.message)
	}
	m49Count(t, f, 1)
	before := m49State(t, f, m49ID(1), "local.success")
	// Even if invalid requests reach validation first, the shared valid key
	// remains available. A repeat is a conflict, never another execution.
	m48Error(t, m48Request(f.router, valid.method, valid.path, valid.body, valid.header), 409, "CONFLICT", "execution conflict")
	if after := m49State(t, f, m49ID(1), "local.success"); !reflect.DeepEqual(before, after) {
		t.Fatal("rejected repeat mutated lifecycle")
	}
	m49Count(t, f, 1)
}

func TestM49ConcurrentStatusVisibility(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			f := newM49Fixture(t, 1)
			allowed := f.statusRoute(t, "/test-only/status", []string{m49ID(1), m49ID(2)})
			denied := f.statusRoute(t, "/test-only/concealed", nil)
			getRequests := make([]m49Request, 0, 40)
			for i := 0; i < 20; i++ {
				getRequests = append(getRequests,
					m49Request{method: "GET", path: "/test-only/status/" + m49ID(1)},
					m49Request{method: "GET", path: "/test-only/concealed/" + m49ID(1)})
			}
			post := m49Post(t, "status-key", definition, "status-event")
			requests := append([]m49Request{post}, getRequests...)
			results := m49Race(f.router, requests)
			status, version, _ := m49Expected(definition, m49ID(1))
			m48JSON(t, results[0], 201, m48Created(m49ID(1), status))
			for i, w := range results[1:] {
				if i%2 == 1 || w.Code == 404 {
					m48Error(t, w, 404, "NOT_FOUND", "execution not found")
					continue
				}
				// Reads may precede creation or observe READY1, RUNNING2 or
				// the terminal state. Never demand the eventual state mid-race.
				var expected map[string]any
				for _, candidate := range []map[string]any{
					m48Status(m49ID(1), "READY", 1),
					m48Status(m49ID(1), "RUNNING", 2),
					m48Status(m49ID(1), status, version),
				} {
					if m48MatchJSON(w.Body.Bytes(), candidate) == nil {
						expected = candidate
						break
					}
				}
				if expected == nil {
					t.Fatalf("incoherent status response %d: %s", w.Code, w.Body)
				}
				m48JSON(t, w, 200, expected)
			}
			if denied.calls.Load() != 0 || allowed.calls.Load() != 20 {
				t.Fatalf("status reads allowed=%d concealed=%d", allowed.calls.Load(), denied.calls.Load())
			}
			before := m49State(t, f, m49ID(1), definition)
			// A second, GET-only phase proves no lifecycle mutation or replay,
			// including the nonterminal WAIT/condition executions.
			for i, w := range m49Race(f.router, getRequests) {
				if i%2 == 0 {
					m48JSON(t, w, 200, m48Status(m49ID(1), status, version))
				} else {
					m48Error(t, w, 404, "NOT_FOUND", "execution not found")
				}
			}
			missing := m48Request(f.router, "GET", "/test-only/status/"+m49ID(2), "", nil)
			concealed := m48Request(f.router, "GET", "/test-only/concealed/"+m49ID(1), "", nil)
			m48Error(t, missing, 404, "NOT_FOUND", "execution not found")
			m48Error(t, concealed, 404, "NOT_FOUND", "execution not found")
			if missing.Body.String() != concealed.Body.String() || !reflect.DeepEqual(missing.Header(), concealed.Header()) || denied.calls.Load() != 0 || allowed.calls.Load() != 41 {
				t.Fatal("concealment read storage or differed from missing; authorized read count changed")
			}
			if after := m49State(t, f, m49ID(1), definition); !reflect.DeepEqual(before, after) {
				t.Fatal("GET mutated stored state")
			}
			m49Count(t, f, 1)
		})
	}
}
