package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ValidationFailure carries the approved response for a future HTTP handler.
// No handler is registered by this validation boundary.
type ValidationFailure struct {
	StatusCode int
	Response   ErrorResponse
}

func invalidRequest(message string) (CreateExecutionRequest, *ValidationFailure) {
	return CreateExecutionRequest{}, &ValidationFailure{
		StatusCode: http.StatusBadRequest,
		Response:   ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: message}},
	}
}

// ParseAndValidateCreateExecutionRequest validates supplied bytes without invoking
// business operations. approvedDefinitions is a trusted server-supplied set of
// supported definitions, not caller authorization. Callers must not mutate it
// concurrently. This function does not enforce a streaming HTTP body-read limit.
func ParseAndValidateCreateExecutionRequest(body []byte, approvedDefinitions map[string]struct{}) (CreateExecutionRequest, *ValidationFailure) {
	if len(body) > 4096 {
		return invalidRequest("body must not exceed 4096 bytes")
	}
	const invalidJSON = "body must be one JSON object"
	decoder := json.NewDecoder(bytes.NewReader(body))
	opening, err := decoder.Token()
	if err != nil || opening != json.Delim('{') {
		return invalidRequest(invalidJSON)
	}
	var request CreateExecutionRequest
	seen := make(map[string]bool, 2)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return invalidRequest(invalidJSON)
		}
		key, ok := token.(string)
		if !ok {
			return invalidRequest(invalidJSON)
		}
		if key != "workflow_definition_id" && key != "event_id" {
			return invalidRequest("unsupported request field")
		}
		if seen[key] {
			return invalidRequest("duplicate request field")
		}
		seen[key] = true
		token, err = decoder.Token()
		if err != nil {
			return invalidRequest(invalidJSON)
		}
		value, ok := token.(string)
		if !ok {
			return invalidRequest(key + " must be a string")
		}
		if key == "workflow_definition_id" {
			request.WorkflowDefinitionID = value
		} else {
			request.EventID = value
		}
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return invalidRequest(invalidJSON)
	}
	if _, err := decoder.Token(); err != io.EOF {
		return invalidRequest(invalidJSON)
	}
	// Fixed order makes missing/invalid identifier diagnostics deterministic.
	for _, field := range []struct{ key, value string }{
		{"workflow_definition_id", request.WorkflowDefinitionID},
		{"event_id", request.EventID},
	} {
		if !seen[field.key] {
			return invalidRequest(field.key + " is required")
		}
		if !validRequestIdentifier(field.value) {
			return invalidRequest(field.key + " must contain 1 to 128 ASCII letters, digits, underscores, hyphens or periods")
		}
	}
	if _, approved := approvedDefinitions[request.WorkflowDefinitionID]; !approved {
		return invalidRequest("unsupported workflow definition")
	}
	return request, nil
}

func validRequestIdentifier(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-' || c == '.') {
			return false
		}
	}
	return true
}
