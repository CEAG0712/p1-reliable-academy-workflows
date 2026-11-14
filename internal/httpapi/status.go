package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"example.com/p1-reliable-academy-workflows/internal/application"
	"github.com/gin-gonic/gin"
)

// NewLocalStatusHandler constructs an unregistered local/test adapter. The host
// owns its scope and requests; exposing it operationally would delegate fixture
// visibility to clients. It supplies no production authentication or ownership.
// NewRouter does not register this handler. No dependency is called here.
func NewLocalStatusHandler(scope *application.LocalStatusScope) (gin.HandlerFunc, error) {
	if scope == nil {
		return nil, fmt.Errorf("local status scope is required")
	}
	return newLocalStatusHandler(scope, scope.Querier())
}

// Injection is private; public composition always uses the sealed scope's service.
func newLocalStatusHandler(scope *application.LocalStatusScope, querier application.LocalExecutionStatusQuerier) (gin.HandlerFunc, error) {
	if scope == nil || scope.Querier() == nil || querier == nil {
		return nil, fmt.Errorf("configured local status scope and querier are required")
	}
	v := reflect.ValueOf(querier)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if v.IsNil() {
			return nil, fmt.Errorf("local status querier is required")
		}
	}
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		id := c.Param("execution_id")
		// Gin supplies the path parameter. Never trim, case-fold or decode it again.
		if !validRequestIdentifier(id) {
			writeStatusError(c, &application.StatusQueryError{Kind: application.StatusQueryInvalidCommand})
			return
		}
		command, err := scope.Command(id)
		if err != nil {
			writeStatusError(c, err)
			return
		}
		result, err := querier.QueryLocal(c.Request.Context(), command)
		if err != nil {
			writeStatusError(c, err)
			return
		}
		if result.ExecutionID != id || result.Version == 0 || !publicExecutionStatus(result.Status) {
			writeStatusError(c, &application.StatusQueryError{Kind: application.StatusQueryInvalidSnapshot})
			return
		}
		c.JSON(http.StatusOK, ExecutionStatusResponse{ExecutionID: result.ExecutionID, Status: result.Status, Version: result.Version})
	}, nil
}
func publicExecutionStatus(s string) bool {
	switch s {
	case "READY", "RUNNING", "COMPLETED", "FAILED":
		return true
	}
	return false
}
func writeStatusError(c *gin.Context, err error) {
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "internal error"
	var failure *application.StatusQueryError
	if errors.As(err, &failure) && failure != nil {
		switch failure.Kind {
		case application.StatusQueryInvalidCommand:
			status, code, message = http.StatusBadRequest, "INVALID_REQUEST", "invalid execution ID"
		case application.StatusQueryMissingAuthority:
			status, code, message = http.StatusUnauthorized, "UNAUTHENTICATED", "local capability required"
		case application.StatusQueryForbidden:
			status, code, message = http.StatusForbidden, "FORBIDDEN", "local status query not permitted"
		case application.StatusQueryConcealed, application.StatusQueryNotFound:
			status, code, message = http.StatusNotFound, "NOT_FOUND", "execution not found"
		case application.StatusQueryUnavailable:
			status, code, message = http.StatusServiceUnavailable, "UNAVAILABLE", "local status query unavailable"
		}
	}
	c.JSON(status, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
