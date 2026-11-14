// Package application provides read-only submission preparation. Preparation
// validates configuration; it never authorizes, creates or runs an execution.
package application

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	workflow "example.com/p1-reliable-academy-workflows"
)

var (
	ErrDefinitionNotFound = errors.New("workflow definition not found")
	ErrUnavailable        = errors.New("preparation dependency unavailable")
	ErrInvalidDefinition  = errors.New("invalid server workflow definition")
	ErrInvalidCommand     = errors.New("invalid submission command")
)

// SubmissionCommand contains validated identifiers, not trusted caller identity,
// payment evidence or an idempotency key. Transport callers validate syntax.
type SubmissionCommand struct {
	WorkflowDefinitionID string
	EventID              string
}

// PreparationResult is a configuration summary, not an execution receipt or
// readiness promise. Handler registration and event authenticity are not checked.
type PreparationResult struct {
	WorkflowDefinitionID string
	EventID              string
	StepCount            int
}

type WorkflowService interface {
	PrepareSubmission(context.Context, SubmissionCommand) (PreparationResult, error)
}

// DefinitionReader must perform read-only lookup of server-owned definitions.
// Returned configuration must remain stable for the duration of preparation.
// Missing definitions return ErrDefinitionNotFound; availability failures may
// wrap ErrUnavailable. It supplies no execution or business-action capability.
type DefinitionReader interface {
	LookupDefinition(context.Context, string) (workflow.WorkflowDefinition, error)
}

type Service struct{ definitions DefinitionReader }

// NewService validates wiring without calling the reader or any external system.
func NewService(definitions DefinitionReader) (*Service, error) {
	if nilReader(definitions) {
		return nil, fmt.Errorf("definition reader is required")
	}
	return &Service{definitions: definitions}, nil
}

func nilReader(reader DefinitionReader) bool {
	if reader == nil {
		return true
	}
	v := reflect.ValueOf(reader)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func (s *Service) PrepareSubmission(ctx context.Context, command SubmissionCommand) (PreparationResult, error) {
	if s == nil || nilReader(s.definitions) {
		return PreparationResult{}, ErrUnavailable
	}
	if ctx == nil || command.WorkflowDefinitionID == "" || command.EventID == "" {
		return PreparationResult{}, ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return PreparationResult{}, err
	}
	d, err := s.definitions.LookupDefinition(ctx, command.WorkflowDefinitionID)
	if err != nil {
		return PreparationResult{}, fmt.Errorf("read definition: %w", err)
	}
	if d.ID != command.WorkflowDefinitionID {
		return PreparationResult{}, fmt.Errorf("%w: lookup identity mismatch", ErrInvalidDefinition)
	}
	if err := d.Validate(); err != nil {
		return PreparationResult{}, fmt.Errorf("%w: %v", ErrInvalidDefinition, err)
	}
	return PreparationResult{WorkflowDefinitionID: d.ID, EventID: command.EventID, StepCount: len(d.Steps)}, nil
}
