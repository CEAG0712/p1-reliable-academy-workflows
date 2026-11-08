package enrollment

import (
	"reflect"
	"testing"
	"time"
)

func validationDefinition(steps map[string]StepDefinition) WorkflowDefinition {
	return WorkflowDefinition{ID: "workflow", Trigger: TriggerDefinition{EventType: "event"}, StartStep: "start", Steps: steps}
}

func TestWorkflowDefinitionValidateAcceptsValidGraphs(t *testing.T) {
	cases := []struct {
		name  string
		steps map[string]StepDefinition
	}{
		{"action", map[string]StepDefinition{
			"start": {Type: StepAction, Action: "send", Next: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
		{"branch and merge", map[string]StepDefinition{
			"start": {Type: StepCondition, Condition: "eligible", OnTrue: "yes", OnFalse: "no"},
			"yes":   {Type: StepAction, Action: "approve", Next: "end"},
			"no":    {Type: StepAction, Action: "reject", Next: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
		{"wait", map[string]StepDefinition{
			"start": {Type: StepWait, Delay: time.Second, Next: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
		{"terminal wait", map[string]StepDefinition{"start": {Type: StepWait, Delay: time.Second}}},
		{"terminal true branch", map[string]StepDefinition{
			"start": {Type: StepCondition, Condition: "eligible", OnFalse: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
		{"terminal false branch", map[string]StepDefinition{
			"start": {Type: StepCondition, Condition: "eligible", OnTrue: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
		{"both branches terminal", map[string]StepDefinition{"start": {Type: StepCondition, Condition: "eligible"}}},
		{"same branch target", map[string]StepDefinition{
			"start": {Type: StepCondition, Condition: "eligible", OnTrue: "end", OnFalse: "end"},
			"end":   {Type: StepAction, Action: "finish"},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := validationDefinition(tc.steps)
			if err := d.Validate(); err != nil {
				t.Fatalf("valid definition rejected: %v", err)
			}
		})
	}
}

func TestWorkflowDefinitionValidateRejectsMalformedDefinitions(t *testing.T) {
	cases := []struct {
		name   string
		change func(*WorkflowDefinition)
		want   string
	}{
		{"missing ID", func(d *WorkflowDefinition) { d.ID = "" }, "workflow ID is required"},
		{"blank ID", func(d *WorkflowDefinition) { d.ID = " \t" }, "workflow ID is required"},
		{"missing trigger", func(d *WorkflowDefinition) { d.Trigger.EventType = "" }, "trigger event type is required"},
		{"blank trigger", func(d *WorkflowDefinition) { d.Trigger.EventType = " " }, "trigger event type is required"},
		{"missing start", func(d *WorkflowDefinition) { d.StartStep = "" }, "start step is required"},
		{"unknown start", func(d *WorkflowDefinition) { d.StartStep = "missing" }, `start step "missing" does not exist`},
		{"nil steps", func(d *WorkflowDefinition) { d.Steps = nil }, `start step "start" does not exist`},
		{"empty step ID", func(d *WorkflowDefinition) { d.Steps[""] = StepDefinition{Type: StepAction, Action: "send"} }, "step ID is required"},
		{"blank step ID", func(d *WorkflowDefinition) { d.Steps[" "] = StepDefinition{Type: StepAction, Action: "send"} }, "step ID is required"},
		{"unknown type", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: "unknown"} }, `step "start": unsupported type "unknown"`},
		{"missing action", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepAction} }, `step "start": action handler name is required`},
		{"blank action", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepAction, Action: " "} }, `step "start": action handler name is required`},
		{"missing condition", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepCondition} }, `step "start": condition handler name is required`},
		{"blank condition", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: " "} }, `step "start": condition handler name is required`},
		{"broken Next", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepAction, Action: "send", Next: "missing"}
		}, `step "start": Next references missing step "missing"`},
		{"broken wait Next", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepWait, Delay: time.Second, Next: "missing"}
		}, `step "start": Next references missing step "missing"`},
		{"broken OnTrue", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: "check", OnTrue: "missing"}
		}, `step "start": OnTrue references missing step "missing"`},
		{"broken OnFalse", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: "check", OnFalse: "missing"}
		}, `step "start": OnFalse references missing step "missing"`},
		{"condition Next", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: "check", Next: "end"}
		}, `step "start": condition cannot configure Action, Next or Delay`},
		{"action OnTrue", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepAction, Action: "send", OnTrue: "end"}
		}, `step "start": action cannot configure Condition, OnTrue, OnFalse or Delay`},
		{"action OnFalse", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepAction, Action: "send", OnFalse: "end"}
		}, `step "start": action cannot configure Condition, OnTrue, OnFalse or Delay`},
		{"action condition", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepAction, Action: "send", Condition: "check"}
		}, `step "start": action cannot configure Condition, OnTrue, OnFalse or Delay`},
		{"action delay", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepAction, Action: "send", Delay: time.Second}
		}, `step "start": action cannot configure Condition, OnTrue, OnFalse or Delay`},
		{"condition action", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: "check", Action: "send"}
		}, `step "start": condition cannot configure Action, Next or Delay`},
		{"condition delay", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepCondition, Condition: "check", Delay: time.Second}
		}, `step "start": condition cannot configure Action, Next or Delay`},
		{"zero wait", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepWait} }, `step "start": wait delay must be positive`},
		{"negative wait", func(d *WorkflowDefinition) { d.Steps["start"] = StepDefinition{Type: StepWait, Delay: -time.Second} }, `step "start": wait delay must be positive`},
		{"wait branching", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepWait, Delay: time.Second, OnTrue: "end"}
		}, `step "start": wait cannot configure Action, Condition, OnTrue or OnFalse`},
		{"wait false branch", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepWait, Delay: time.Second, OnFalse: "end"}
		}, `step "start": wait cannot configure Action, Condition, OnTrue or OnFalse`},
		{"wait action", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepWait, Delay: time.Second, Action: "send"}
		}, `step "start": wait cannot configure Action, Condition, OnTrue or OnFalse`},
		{"wait condition", func(d *WorkflowDefinition) {
			d.Steps["start"] = StepDefinition{Type: StepWait, Delay: time.Second, Condition: "check"}
		}, `step "start": wait cannot configure Action, Condition, OnTrue or OnFalse`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := validationDefinition(map[string]StepDefinition{
				"start": {Type: StepAction, Action: "send"},
				"end":   {Type: StepAction, Action: "finish"},
			})
			tc.change(&d)
			err := d.Validate()
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestWorkflowDefinitionValidateRejectsCycles(t *testing.T) {
	cases := []struct {
		name  string
		steps map[string]StepDefinition
		want  string
	}{
		{"self", map[string]StepDefinition{"start": {Type: StepAction, Action: "send", Next: "start"}}, `unsupported cycle at step "start"`},
		{"multiple steps", map[string]StepDefinition{
			"start": {Type: StepAction, Action: "send", Next: "wait"},
			"wait":  {Type: StepWait, Delay: time.Second, Next: "start"},
		}, `unsupported cycle at step "start"`},
		{"true branch", map[string]StepDefinition{"start": {Type: StepCondition, Condition: "check", OnTrue: "start"}}, `unsupported cycle at step "start"`},
		{"false branch", map[string]StepDefinition{"start": {Type: StepCondition, Condition: "check", OnFalse: "start"}}, `unsupported cycle at step "start"`},
		{"unreachable", map[string]StepDefinition{
			"start": {Type: StepAction, Action: "send"},
			"other": {Type: StepAction, Action: "send", Next: "other"},
		}, `unsupported cycle at step "other"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validationDefinition(tc.steps).Validate()
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestWorkflowDefinitionValidateIsDeterministicAndReadOnly(t *testing.T) {
	for _, cyclic := range []bool{false, true} {
		for i := 0; i < 100; i++ {
			ids := []string{"start", "alpha", "zeta"}
			steps := make(map[string]StepDefinition)
			for j := range ids {
				id := ids[(i+j)%len(ids)]
				step := StepDefinition{Type: StepAction, Action: "send"}
				if id != "start" {
					step.Next = "missing"
					if cyclic {
						step.Next = id
					}
				}
				steps[id] = step
			}
			d := validationDefinition(steps)
			before := make(map[string]StepDefinition)
			for id, step := range steps {
				before[id] = step
			}
			want := `step "alpha": Next references missing step "missing"`
			if cyclic {
				want = `unsupported cycle at step "alpha"`
			}
			err := d.Validate()
			if err == nil || err.Error() != want {
				t.Fatalf("iteration %d: got %v, want %q", i, err, want)
			}
			if !reflect.DeepEqual(before, d.Steps) {
				t.Fatal("validation mutated definition steps")
			}
		}
	}
}
