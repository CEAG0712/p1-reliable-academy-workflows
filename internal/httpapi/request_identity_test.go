package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

type identityHTTPIDs struct {
	calls atomic.Int32
	after func()
	err   error
}

func (s *identityHTTPIDs) NextExecutionID(context.Context) (string, error) {
	n := s.calls.Add(1)
	if s.after != nil {
		s.after()
	}
	return fmt.Sprintf("exec-%032x", n), s.err
}
func identityHTTPFixture(t *testing.T, repo *workflow.InMemoryWorkflowRepository, limit int) (*gin.Engine, *identityHTTPIDs) {
	t.Helper()
	ids := &identityHTTPIDs{}
	s, e := application.NewLocalRequestIdentityScope(repo, ids, "http-identity-worker", limit)
	if e != nil {
		t.Fatal(e)
	}
	h, e := NewLocalRequestIdentityHandler(s)
	if e != nil {
		t.Fatal(e)
	}
	return invocationRouter(t, h), ids
}
func identityHTTPSend(r *gin.Engine, body string, headers http.Header, ctx context.Context) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/test-only/invoke", strings.NewReader(body))
	req.Header = headers
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func identityHeader(key string) http.Header { return http.Header{"Idempotency-Key": []string{key}} }
func identityHTTPError(t *testing.T, w *httptest.ResponseRecorder, status int, code, message string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status %d body %s", w.Code, w.Body)
	}
	assertPreparationJSON(t, w, ErrorResponse{Error: APIError{Code: code, Message: message}})
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("cacheable error")
	}
}
func TestLocalRequestIdentityHTTPConstruction(t *testing.T) {
	for _, s := range []*application.LocalRequestIdentityScope{nil, {}} {
		if h, e := NewLocalRequestIdentityHandler(s); h != nil || e == nil {
			t.Fatal("bad scope accepted")
		}
	}
}
func TestLocalRequestIdentityHTTPHeadersBeforeReservation(t *testing.T) {
	for _, headers := range []http.Header{
		nil, {}, identityHeader(""), {"Idempotency-Key": []string{"key", "key"}}, identityHeader("key,key"), identityHeader("key, other"),
		identityHeader(" key"), identityHeader("key "), identityHeader("a/b"), identityHeader("é"), identityHeader(strings.Repeat("a", 129)),
		{"Idempotency-Key": []string{"key"}, "idempotency-key": []string{"other"}},
	} {
		r, ids := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
		w := identityHTTPSend(r, validationBody("local.success", "event"), headers, nil)
		identityHTTPError(t, w, 400, "INVALID_REQUEST", "invalid idempotency key")
		if ids.calls.Load() != 0 {
			t.Fatal("invalid header invoked")
		}
		if w = identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil); w.Code != 201 {
			t.Fatal("invalid header reserved", w.Code, w.Body)
		}
	}
}
func TestLocalRequestIdentityHTTPBodyBeforeReservation(t *testing.T) {
	for _, body := range []string{"", `{}`, `null`, validationBody("other", "event"), validationBody("local.success", "bad/event"), `{"workflow_definition_id":"local.success","event_id":"event","tenant_id":"tenant"}`, strings.Repeat(" ", 4097)} {
		r, ids := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
		w := identityHTTPSend(r, body, identityHeader("key"), nil)
		if w.Code != 400 || w.Header().Get("Cache-Control") != "no-store" || ids.calls.Load() != 0 {
			t.Fatal("invalid body invoked", w.Code, w.Body)
		}
		w = identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil)
		if w.Code != 201 {
			t.Fatal("invalid body reserved", w.Code, w.Body)
		}
	}
}
func TestLocalRequestIdentityHTTPKnownOutcomes(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			repo := workflow.NewInMemoryWorkflowRepository()
			r, ids := identityHTTPFixture(t, repo, 1)
			w := identityHTTPSend(r, validationBody(definition, "event"), identityHeader("key"), nil)
			if w.Code != 201 || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal(w.Code, w.Body)
			}
			var response CreateExecutionResponse
			if e := json.Unmarshal(w.Body.Bytes(), &response); e != nil {
				t.Fatal(e)
			}
			before, ok := repo.Get(response.ExecutionID)
			if !ok || before.Status != response.Status {
				t.Fatal("not authoritative")
			}
			if definition == "local.failure" && (before.Status != "FAILED" || before.LastError == "" || len(before.History) != 2 || before.History[1].Reason != before.LastError) {
				t.Fatal("false FAILED success")
			}
			for _, body := range []string{validationBody(definition, "event"), `{ "event_id":"event", "workflow_definition_id":"` + definition + `" }`, validationBody(definition, "changed"), validationBody("local.success", "changed")} {
				identityHTTPError(t, identityHTTPSend(r, body, identityHeader("key"), nil), 409, "CONFLICT", "execution conflict")
			}
			after, _ := repo.Get(response.ExecutionID)
			if !reflect.DeepEqual(before, after) || ids.calls.Load() != 1 {
				t.Fatal("repeat mutated/replayed")
			}
			var keys map[string]any
			json.Unmarshal(w.Body.Bytes(), &keys)
			if len(keys) != 2 {
				t.Fatal("projection leak")
			}
			if strings.Contains(w.Body.String(), "worker") || strings.Contains(w.Body.String(), "controlled") {
				t.Fatal("diagnostic leak")
			}
		})
	}
}
func TestLocalRequestIdentityHTTPUncertainCreation(t *testing.T) {
	// The concrete zero-value repository returns false/error without callbacks.
	// The opt-in boundary deliberately retains even this rejected creation key.
	r, ids := identityHTTPFixture(t, &workflow.InMemoryWorkflowRepository{}, 2)
	for i := 0; i < 3; i++ {
		identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil), 500, "INTERNAL_ERROR", "internal error")
	}
	identityHTTPError(t, identityHTTPSend(r, validationBody("local.failure", "event"), identityHeader("key"), nil), 409, "CONFLICT", "execution conflict")
	if ids.calls.Load() != 1 {
		t.Fatal("uncertain creation replayed")
	}
}
func TestLocalRequestIdentityHTTPCancellationAndCapacity(t *testing.T) {
	for _, phase := range []string{"before-reservation", "before-create"} {
		t.Run(phase, func(t *testing.T) {
			repo := workflow.NewInMemoryWorkflowRepository()
			r, ids := identityHTTPFixture(t, repo, 1)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if phase == "before-reservation" {
				cancel()
			} else {
				ids.after = cancel
			}
			identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), ctx), 503, "UNAVAILABLE", "local invocation unavailable")
			if _, ok := repo.Get(httpInvocationID); ok {
				t.Fatal("cancel created")
			}
			ids.after = nil
			w := identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil)
			if phase == "before-reservation" {
				if w.Code != 201 {
					t.Fatal("cancel reserved")
				}
			} else {
				identityHTTPError(t, w, 409, "CONFLICT", "execution conflict")
			}
		})
	}
	r, ids := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
	ids.err = errors.New("private dependency error")
	identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil), 503, "UNAVAILABLE", "local invocation unavailable")
	identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("other"), nil), 503, "UNAVAILABLE", "local invocation unavailable")
	identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil), 409, "CONFLICT", "execution conflict")
	if ids.calls.Load() != 1 {
		t.Fatal("capacity invoked")
	}
}
func TestLocalRequestIdentityHTTPPanicRetention(t *testing.T) {
	r, ids := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
	ids.after = func() { panic("private panic") }
	func() {
		defer func() {
			if recover() != "private panic" {
				t.Error("panic swallowed")
			}
		}()
		identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil)
	}()
	identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader("key"), nil), 500, "INTERNAL_ERROR", "internal error")
	if ids.calls.Load() != 1 {
		t.Fatal("panic replayed")
	}
}

func TestLocalRequestIdentityHTTPBodyBounds(t *testing.T) {
	for _, fail := range []bool{false, true} {
		r, ids := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
		body := &preparationBody{fail: fail}
		req := httptest.NewRequest("POST", "/test-only/invoke", nil)
		req.Body = body
		req.Header = identityHeader("key")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		message := "body must not exceed 4096 bytes"
		if fail {
			message = "could not read request body"
		} else if body.read != 4097 {
			t.Fatal("unbounded body read")
		}
		identityHTTPError(t, w, 400, "INVALID_REQUEST", message)
		if ids.calls.Load() != 0 {
			t.Fatal("bad stream invoked")
		}
		valid := validationBody("local.success", "event")
		w = identityHTTPSend(r, valid+strings.Repeat(" ", 4096-len(valid)), identityHeader("key"), nil)
		if w.Code != 201 {
			t.Fatal("bad stream reserved or boundary rejected", w.Code, w.Body)
		}
	}
}
func TestLocalRequestIdentityHTTPKeyBoundaries(t *testing.T) {
	for _, key := range []string{"a", strings.Repeat("a", 128), "Aa09_.-"} {
		r, _ := identityHTTPFixture(t, workflow.NewInMemoryWorkflowRepository(), 1)
		w := identityHTTPSend(r, validationBody("local.success", "event"), http.Header{"idempotency-key": []string{key}}, nil)
		if w.Code != 201 {
			t.Fatal("valid exact key rejected", key, w.Code)
		}
		identityHTTPError(t, identityHTTPSend(r, validationBody("local.success", "event"), identityHeader(key), nil), 409, "CONFLICT", "execution conflict")
	}
}
