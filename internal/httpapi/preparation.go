package httpapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

// preparationResponse is used only by the unregistered preparation adapter.
// It deliberately cannot represent execution creation, scheduling or permission.
type preparationResponse struct {
	WorkflowDefinitionID string `json:"workflow_definition_id"`
	EventID              string `json:"event_id"`
	StepCount            int    `json:"step_count"`
}

// NewPreparationHandler constructs an unregistered read-only adapter. It must
// not be exposed operationally before identity/authorization/visibility policy
// is approved. NewRouter does not call this constructor. The support set is
// copied at construction; membership grants no caller or payment authority.
func NewPreparationHandler(service application.WorkflowService, approvedDefinitions map[string]struct{}) (gin.HandlerFunc, error) {
	if service == nil {
		return nil, fmt.Errorf("workflow service is required")
	}
	v := reflect.ValueOf(service)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if v.IsNil() {
			return nil, fmt.Errorf("workflow service is required")
		}
	}
	approved := make(map[string]struct{}, len(approvedDefinitions))
	for id := range approvedDefinitions {
		approved[id] = struct{}{}
	}
	return func(c *gin.Context) {
		// Read one extra byte so the existing validator detects an oversized body
		// without allocating memory proportional to untrusted input size.
		body, err := io.ReadAll(io.LimitReader(c.Request.Body, 4097))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: "could not read request body"}})
			return
		}
		request, failure := ParseAndValidateCreateExecutionRequest(body, approved)
		if failure != nil {
			c.JSON(failure.StatusCode, failure.Response)
			return
		}
		result, err := service.PrepareSubmission(c.Request.Context(), application.SubmissionCommand{WorkflowDefinitionID: request.WorkflowDefinitionID, EventID: request.EventID})
		if err != nil {
			status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "internal error"
			switch {
			case errors.Is(err, application.ErrInvalidCommand):
				status, code, message = http.StatusBadRequest, "INVALID_REQUEST", "invalid preparation command"
			case errors.Is(err, application.ErrDefinitionNotFound):
				status, code, message = http.StatusNotFound, "NOT_FOUND", "workflow definition not found"
			case errors.Is(err, application.ErrUnavailable):
				status, code, message = http.StatusServiceUnavailable, "UNAVAILABLE", "preparation service unavailable"
			}
			c.JSON(status, ErrorResponse{Error: APIError{Code: code, Message: message}})
			return
		}
		c.JSON(http.StatusOK, preparationResponse{WorkflowDefinitionID: result.WorkflowDefinitionID, EventID: result.EventID, StepCount: result.StepCount})
	}, nil
}
