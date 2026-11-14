package httpapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

// NewLocalRequestIdentityHandler constructs an unregistered fixture-only adapter.
// It does not add a listener, authentication, tenant authority or status access.
// Trusted local composition owns scope and requests; never expose it operationally.
func NewLocalRequestIdentityHandler(scope *application.LocalRequestIdentityScope) (gin.HandlerFunc, error) {
	if scope == nil || scope.Invoker() == nil {
		return nil, fmt.Errorf("configured local request identity scope is required")
	}
	invoker, support := scope.Invoker(), scope.ApprovedDefinitions()
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		body, err := io.ReadAll(io.LimitReader(c.Request.Body, 4097))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: "could not read request body"}})
			return
		}
		request, failure := ParseAndValidateCreateExecutionRequest(body, support)
		if failure != nil {
			c.JSON(failure.StatusCode, failure.Response)
			return
		}
		// Count all case-insensitive map entries, including noncanonical headers
		// supplied by trusted in-process harnesses. Never trim or comma-split keys.
		var values []string
		for name, entries := range c.Request.Header {
			if strings.EqualFold(name, "Idempotency-Key") {
				values = append(values, entries...)
			}
		}
		if len(values) != 1 || !validRequestIdentifier(values[0]) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: APIError{Code: "INVALID_REQUEST", Message: "invalid idempotency key"}})
			return
		}
		command, err := scope.Command(values[0], application.SubmissionCommand{
			WorkflowDefinitionID: request.WorkflowDefinitionID, EventID: request.EventID,
		})
		if err != nil {
			writeInvocationError(c, err)
			return
		}
		out, err := invoker.InvokeLocalRequest(c.Request.Context(), command)
		// Only the sealed scope's concrete invocation service supplies outcomes.
		// It proves accepted lifecycle operations and exact authoritative history;
		// FAILED alone, without that M4.5 evidence, cannot produce 201 FAILED.
		if localCreationResponseValid(out, err) {
			c.JSON(http.StatusCreated, CreateExecutionResponse{ExecutionID: out.ExecutionID, Status: out.Snapshot.Status})
			return
		}
		if err == nil {
			err = &application.InvocationError{Kind: application.InvocationReconciliation}
		}
		writeInvocationError(c, err)
	}, nil
}
