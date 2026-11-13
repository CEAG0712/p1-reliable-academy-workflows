package httpapi

import (
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func assertHealth(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	mediaType, _, err := mime.ParseMediaType(response.Header().Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		t.Fatalf("content type = %q, error = %v", response.Header().Get("Content-Type"), err)
	}
	var body any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body, map[string]any{"status": "ok"}) {
		t.Fatalf("body = %#v, want exactly {status: ok}", body)
	}
}

func TestHealth(t *testing.T) {
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))
	assertHealth(t, response)
}

func TestUnknownRoute(t *testing.T) {
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
	if response.Body.String() != "404 page not found" {
		t.Fatalf("unexpected not-found body: %q", response.Body.String())
	}
}

func TestHealthUnsupportedMethod(t *testing.T) {
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/health", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want Gin default 404", response.Code)
	}
	if response.Body.String() != "404 page not found" {
		t.Fatalf("unexpected method rejection body: %q", response.Body.String())
	}
}

func TestHealthRepeatedCalls(t *testing.T) {
	router := NewRouter()
	for i := 0; i < 10; i++ {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))
		assertHealth(t, response)
	}
}
