package application

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	workflow "example.com/p1-reliable-academy-workflows"
)

type catalogProbe struct {
	definition workflow.WorkflowDefinition
	err        error
	calls      int
	id         string
	ctx        context.Context
}

func (p *catalogProbe) LookupDefinition(ctx context.Context, id string) (workflow.WorkflowDefinition, error) {
	p.calls++
	p.id = id
	p.ctx = ctx
	return p.definition, p.err
}

func preparationDefinition() workflow.WorkflowDefinition {
	return workflow.WorkflowDefinition{ID: "definition-1", Trigger: workflow.TriggerDefinition{EventType: "event"}, StartStep: "grant", Steps: map[string]workflow.StepDefinition{"grant": {Type: workflow.StepAction, Action: "entitlement.grant"}}}
}

func TestServiceConstruction(t *testing.T) {
	var typedNil *catalogProbe
	for _, dependency := range []DefinitionReader{nil, typedNil} {
		if service, err := NewService(dependency); service != nil || err == nil {
			t.Fatalf("missing dependency accepted: %v %v", service, err)
		}
	}
	p := &catalogProbe{}
	if service, err := NewService(p); service == nil || err != nil || p.calls != 0 {
		t.Fatalf("construction performed lookup or failed: %v %v calls=%d", service, err, p.calls)
	}
}

func TestPreparationReadsConfigurationOnly(t *testing.T) {
	d := preparationDefinition()
	beforeSteps := map[string]workflow.StepDefinition{"grant": d.Steps["grant"]}
	p := &catalogProbe{definition: d}
	service, err := NewService(p)
	if err != nil {
		t.Fatal(err)
	}
	// Real state with diagnostic/history must survive preparation unchanged.
	repo := workflow.NewInMemoryWorkflowRepository()
	repo.Save(workflow.NewWorkflowExecution("existing"))
	repo.TryClaim("existing", "owner")
	repo.Fail("existing", "owner", 2, "retained diagnostic")
	before, _ := repo.Get("existing")
	ctx := context.WithValue(context.Background(), struct{}{}, "controlled-context")
	command := SubmissionCommand{WorkflowDefinitionID: "definition-1", EventID: "CONFIRMED"}
	for i := 1; i <= 2; i++ {
		result, err := service.PrepareSubmission(ctx, command)
		if err != nil || result != (PreparationResult{WorkflowDefinitionID: "definition-1", EventID: "CONFIRMED", StepCount: 1}) {
			t.Fatalf("preparation: %+v %v", result, err)
		}
		if p.calls != i || p.id != command.WorkflowDefinitionID || p.ctx != ctx {
			t.Fatalf("lookup arguments/calls: %+v", p)
		}
	}
	after, _ := repo.Get("existing")
	if !reflect.DeepEqual(before, after) || !reflect.DeepEqual(d.Steps, beforeSteps) {
		t.Fatal("preparation mutated state or configuration")
	}
	if _, exists := repo.Get(command.EventID); exists {
		t.Fatal("correlation became execution identity")
	}
	// The service receives no repository, registry or provider capability. A graph
	// naming a grant is only inspected; no registered callback can be invoked.
}

func TestPreparationErrors(t *testing.T) {
	cause := errors.New("private dependency diagnostic")
	for _, tc := range []struct {
		name                  string
		definition            workflow.WorkflowDefinition
		dependencyError, want error
	}{
		{"not found", preparationDefinition(), fmt.Errorf("lookup: %w", ErrDefinitionNotFound), ErrDefinitionNotFound},
		{"unavailable", preparationDefinition(), fmt.Errorf("catalog: %w", ErrUnavailable), ErrUnavailable},
		{"unexpected", preparationDefinition(), cause, cause},
		{"invalid graph", workflow.WorkflowDefinition{ID: "definition-1"}, nil, ErrInvalidDefinition},
		{"wrong identity", workflow.WorkflowDefinition{ID: "other"}, nil, ErrInvalidDefinition},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := &catalogProbe{definition: tc.definition, err: tc.dependencyError}
			s, err := NewService(p)
			if err != nil {
				t.Fatal(err)
			}
			result, err := s.PrepareSubmission(context.Background(), SubmissionCommand{"definition-1", "event-1"})
			if !errors.Is(err, tc.want) || result != (PreparationResult{}) || p.calls != 1 {
				t.Fatalf("result=%+v err=%v calls=%d", result, err, p.calls)
			}
		})
	}
}

func TestPreparationInvalidUse(t *testing.T) {
	p := &catalogProbe{definition: preparationDefinition()}
	s, err := NewService(p)
	if err != nil {
		t.Fatal(err)
	}
	for _, command := range []SubmissionCommand{{}, {WorkflowDefinitionID: "definition-1"}, {EventID: "event-1"}} {
		result, err := s.PrepareSubmission(context.Background(), command)
		if !errors.Is(err, ErrInvalidCommand) || result != (PreparationResult{}) || p.calls != 0 {
			t.Fatalf("invalid command: %+v %v", result, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.PrepareSubmission(ctx, SubmissionCommand{"definition-1", "event-1"}); !errors.Is(err, context.Canceled) || p.calls != 0 {
		t.Fatalf("cancellation: %v", err)
	}
	for _, invalidService := range []*Service{nil, {}} {
		if _, err := invalidService.PrepareSubmission(context.Background(), SubmissionCommand{"definition-1", "event-1"}); !errors.Is(err, ErrUnavailable) {
			t.Fatalf("unconstructed service: %v", err)
		}
	}
}
