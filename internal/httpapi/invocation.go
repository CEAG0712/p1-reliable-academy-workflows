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

// NewLocalInvocationHandler is unregistered and restricted to harmless fixtures.
// The controlled test harness owns scope and requests. Attaching this handler to
// a listener would delegate its scope to clients; this is NOT authentication.
// Production NewRouter does not call this constructor or register an invoke route.
func NewLocalInvocationHandler(scope *application.LocalInvocationScope) (gin.HandlerFunc, error) {
	if scope == nil {
		return nil, fmt.Errorf("local scope is required")
	}
	return newLocalInvocationHandler(scope, scope.Invoker(), scope.ApprovedDefinitions())
}

// The injectable seam is package-private; the public path supplies only the
// sealed scope's own invoker. Construction performs no operational dependency calls.
func newLocalInvocationHandler(scope *application.LocalInvocationScope, invoker application.LocalWorkflowInvoker, approved map[string]struct{}) (gin.HandlerFunc, error) {
	if scope == nil || scope.Invoker() == nil || invoker == nil {
		return nil, fmt.Errorf("configured local scope and invoker are required")
	}
	v := reflect.ValueOf(invoker)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if v.IsNil() {
			return nil, fmt.Errorf("local invoker is required")
		}
	}
	support := make(map[string]struct{}, len(approved))
	for id := range approved {
		support[id] = struct{}{}
	}
	return func(c *gin.Context) {
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
		command, err := scope.Command(application.SubmissionCommand{WorkflowDefinitionID: request.WorkflowDefinitionID, EventID: request.EventID})
		if err != nil {
			writeInvocationError(c, err)
			return
		}
		outcome, err := invoker.InvokeLocal(c.Request.Context(), command)
		if localCreationResponseValid(outcome, err) {
			c.JSON(http.StatusCreated, CreateExecutionResponse{ExecutionID: outcome.ExecutionID, Status: outcome.Snapshot.Status})
			return
		}
		// A malformed result must never become a successful creation projection.
		if err == nil {
			err = &application.InvocationError{Kind: application.InvocationReconciliation}
		}
		writeInvocationError(c, err)
	}, nil
}
func localCreationResponseValid(out application.InvocationOutcome, err error) bool {
	if out.Creation != application.CreationInserted || out.ExecutionID == "" || out.Snapshot == nil || out.Snapshot.WorkflowID != out.ExecutionID {
		return false
	}
	var failure *application.InvocationError
	errors.As(err, &failure)
	switch out.Kind {
	case application.OutcomeCompleted:
		return err == nil && out.RunAttempted && out.Claim == application.ClaimAccepted && out.Snapshot.Status == "COMPLETED" && out.Snapshot.Version == 3
	case application.OutcomeHandoff:
		return err == nil && out.RunAttempted && out.Claim == application.ClaimAccepted && out.Snapshot.Status == "RUNNING" && out.Snapshot.Version == 2 && out.Cursor != ""
	case application.OutcomeFailed:
		return failure != nil && failure.Kind == application.InvocationExecutionFailed && out.RunAttempted && out.Claim == application.ClaimAccepted && out.Snapshot.Status == "FAILED" && out.Snapshot.Version == 3
	case application.OutcomeCreatedNotStarted:
		return failure != nil && failure.Kind == application.InvocationUnavailable && out.Snapshot.Status == "READY" && out.Snapshot.Version == 1 && !out.RunAttempted && out.Claim == application.ClaimNotObserved
	}
	return false
}
func writeInvocationError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "internal error"
	var failure *application.InvocationError
	if errors.As(err, &failure) {
		switch failure.Kind {
		case application.InvocationInvalidCommand:
			status, code, message = http.StatusBadRequest, "INVALID_REQUEST", "invalid invocation command"
		case application.InvocationMissingAuthority:
			status, code, message = http.StatusUnauthorized, "UNAUTHENTICATED", "local capability required"
		case application.InvocationForbidden:
			status, code, message = http.StatusForbidden, "FORBIDDEN", "local invocation not permitted"
		case application.InvocationDefinitionNotFound:
			status, code, message = http.StatusNotFound, "NOT_FOUND", "workflow definition not found"
		case application.InvocationDuplicate, application.InvocationClaimRejected:
			status, code, message = http.StatusConflict, "CONFLICT", "execution conflict"
		case application.InvocationUnavailable:
			status, code, message = http.StatusServiceUnavailable, "UNAVAILABLE", "local invocation unavailable"
		}
	}
	c.JSON(status, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
