package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Compare complete JSON structures, including numeric types and extra fields.
func assertContractJSON(t *testing.T, value any, want string) {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	decode := func(input string) any {
		t.Helper()
		decoder := json.NewDecoder(strings.NewReader(input))
		decoder.UseNumber()
		var result any
		if err := decoder.Decode(&result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	if !reflect.DeepEqual(decode(string(encoded)), decode(want)) {
		t.Fatalf("JSON = %s, want %s", encoded, want)
	}
}

func TestCreateExecutionRequestJSON(t *testing.T) {
	const input = `{"workflow_definition_id":"definition-1","event_id":"event-1"}`
	var request CreateExecutionRequest
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		t.Fatal(err)
	}
	if request.WorkflowDefinitionID != "definition-1" || request.EventID != "event-1" {
		t.Fatalf("decoded request = %#v", request)
	}
	assertContractJSON(t, request, input)
}

func TestCreateExecutionResponseJSON(t *testing.T) {
	response := CreateExecutionResponse{ExecutionID: "execution-1", Status: "READY"}
	assertContractJSON(t, response, `{"execution_id":"execution-1","status":"READY"}`)
}

func TestExecutionStatusResponseJSON(t *testing.T) {
	for _, status := range []string{"READY", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(status, func(t *testing.T) {
			// The full uint64 range must remain a JSON number without float rounding.
			response := ExecutionStatusResponse{ExecutionID: "execution-1", Status: status, Version: ^uint64(0)}
			assertContractJSON(t, response, `{"execution_id":"execution-1","status":"`+status+`","version":18446744073709551615}`)
		})
	}
}

func TestErrorResponseJSON(t *testing.T) {
	response := ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: "workflow_definition_id is required"}}
	assertContractJSON(t, response, `{"error":{"code":"INVALID_REQUEST","message":"workflow_definition_id is required"}}`)
}

func TestNoSensitiveFieldsExposed(t *testing.T) {
	// Exact allowlists reject any added claim owner, history, attempts, LastError,
	// mutex, tenant, provider or payment fields, including zero-valued fields.
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"request", CreateExecutionRequest{WorkflowDefinitionID: "definition-1", EventID: "event-1"}, `{"workflow_definition_id":"definition-1","event_id":"event-1"}`},
		{"created", CreateExecutionResponse{ExecutionID: "execution-1", Status: "READY"}, `{"execution_id":"execution-1","status":"READY"}`},
		{"status", ExecutionStatusResponse{ExecutionID: "execution-1", Status: "FAILED", Version: 3}, `{"execution_id":"execution-1","status":"FAILED","version":3}`},
		{"error", ErrorResponse{Error: APIError{Code: "INTERNAL_ERROR", Message: "internal error"}}, `{"error":{"code":"INTERNAL_ERROR","message":"internal error"}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertContractJSON(t, test.value, test.want)
		})
	}
}

func TestHealthContractUnchanged(t *testing.T) {
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))
	assertHealth(t, response)
}
