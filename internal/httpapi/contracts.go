package httpapi

// CreateExecutionRequest describes the future creation contract. EventID is
// correlation data, not payment proof, caller identity or an idempotency guarantee.
type CreateExecutionRequest struct {
	WorkflowDefinitionID string `json:"workflow_definition_id"`
	EventID              string `json:"event_id"`
}

// CreateExecutionResponse describes a resource actually created and retrievable
// under the future local contract; it does not promise durable persistence.
type CreateExecutionResponse struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
}

// ExecutionStatusResponse is an intentional public projection, independent of
// internal execution fields. Persisted states are READY, RUNNING, COMPLETED, FAILED.
type ExecutionStatusResponse struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
	Version     uint64 `json:"version"`
}

// APIError contains a public code and safe message, excluding internal diagnostics.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse is the shared future error envelope. Existing health and default
// Gin 404 responses retain their M4.1 behavior.
type ErrorResponse struct {
	Error APIError `json:"error"`
}
