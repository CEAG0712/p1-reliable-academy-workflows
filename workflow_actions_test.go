package enrollment

import (
	"errors"
	"maps"
	"reflect"
	"strings"
	"testing"
	"time"
)

func actionChainDefinition() WorkflowDefinition {
	return WorkflowDefinition{
		ID: "document-processing", Trigger: TriggerDefinition{EventType: "document.received"}, StartStep: "z-first",
		Steps: map[string]StepDefinition{
			"z-first":   {Type: StepAction, Action: "load", Next: "a-middle"},
			"a-middle":  {Type: StepAction, Action: "transform", Next: "m-last"},
			"m-last":    {Type: StepAction, Action: "store"},
			"unreached": {Type: StepAction, Action: "must-not-run"},
		},
	}
}

func TestHandlerRegistryExecuteActionsOrderAndReadOnly(t *testing.T) {
	for rotation := 0; rotation < 4; rotation++ {
		d := actionChainDefinition()
		original := d
		original.Steps = maps.Clone(d.Steps)
		// Vary insertion order; lexical order also differs from the configured path.
		d.Steps = make(map[string]StepDefinition)
		ids := []string{"a-middle", "m-last", "unreached", "z-first"}
		for j := range ids {
			id := ids[(rotation+j)%len(ids)]
			d.Steps[id] = original.Steps[id]
		}
		r := NewHandlerRegistry()
		var calls []string
		for _, name := range []string{"load", "transform", "store", "must-not-run"} {
			if err := r.RegisterAction(name, func() error { calls = append(calls, name); return nil }); err != nil {
				t.Fatal(err)
			}
		}
		conditionCalls := 0
		if err := r.RegisterCondition("unused", func() (bool, error) { conditionCalls++; return true, nil }); err != nil {
			t.Fatal(err)
		}
		actions, conditions := maps.Clone(r.actions), maps.Clone(r.conditions)
		for attempt := 0; attempt < 2; attempt++ {
			calls = nil
			cursor, err := r.ExecuteActions(d, d.StartStep)
			if cursor != "" || err != nil {
				t.Fatalf("got (%q, %v), want terminal/nil", cursor, err)
			}
			if want := []string{"load", "transform", "store"}; !reflect.DeepEqual(calls, want) {
				t.Fatalf("calls = %v, want %v", calls, want)
			}
			if !reflect.DeepEqual(d, original) {
				t.Fatal("executor changed definition")
			}
			if len(r.actions) != len(actions) || len(r.conditions) != len(conditions) {
				t.Fatal("executor changed registry sizes")
			}
			for name, handler := range actions {
				if reflect.ValueOf(r.actions[name]).Pointer() != reflect.ValueOf(handler).Pointer() {
					t.Fatalf("action mapping changed: %s", name)
				}
			}
			for name, handler := range conditions {
				if reflect.ValueOf(r.conditions[name]).Pointer() != reflect.ValueOf(handler).Pointer() {
					t.Fatalf("condition mapping changed: %s", name)
				}
			}
		}
		if conditionCalls != 0 {
			t.Fatal("executor evaluated an unused condition")
		}
		if err := r.InvokeAction("must-not-run"); err != nil {
			t.Fatal(err)
		}
		if got, err := r.EvaluateCondition("unused"); !got || err != nil || conditionCalls != 1 {
			t.Fatal("original condition mapping lost")
		}
	}
}

func TestHandlerRegistryExecuteActionsFailure(t *testing.T) {
	ids := []string{"z-first", "a-middle", "m-last"}
	names := []string{"load", "transform", "store"}
	for failed := range ids {
		t.Run(ids[failed], func(t *testing.T) {
			r := NewHandlerRegistry()
			cause := errors.New("adapter failed")
			var calls []string
			for i, name := range names {
				if err := r.RegisterAction(name, func() error {
					calls = append(calls, name)
					if i == failed {
						return cause
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			d := actionChainDefinition()
			for attempt := 0; attempt < 2; attempt++ {
				calls = nil
				cursor, err := r.ExecuteActions(d, d.StartStep)
				if cursor != ids[failed] || !errors.Is(err, cause) || !strings.Contains(err.Error(), ids[failed]) {
					t.Fatalf("got (%q, %v), want failed cursor/wrapped cause", cursor, err)
				}
				if !reflect.DeepEqual(calls, names[:failed+1]) {
					t.Fatalf("calls = %v", calls)
				}
			}
		})
	}
}

func TestHandlerRegistryExecuteActionsUnknownHandler(t *testing.T) {
	for _, name := range []string{"missing", "Transform", " transform "} {
		t.Run(name, func(t *testing.T) {
			r := NewHandlerRegistry()
			calls := 0
			if err := r.RegisterAction("load", func() error { calls++; return nil }); err != nil {
				t.Fatal(err)
			}
			if err := r.RegisterAction("store", func() error { t.Fatal("ran after unknown handler"); return nil }); err != nil {
				t.Fatal(err)
			}
			d := actionChainDefinition()
			step := d.Steps["a-middle"]
			step.Action = name
			d.Steps["a-middle"] = step
			cursor, err := r.ExecuteActions(d, d.StartStep)
			if cursor != "a-middle" || err == nil || !strings.Contains(err.Error(), name) || calls != 1 {
				t.Fatalf("got (%q, %v), calls=%d", cursor, err, calls)
			}
		})
	}
	var zero HandlerRegistry
	if cursor, err := zero.ExecuteActions(actionChainDefinition(), "z-first"); cursor != "z-first" || err == nil {
		t.Fatalf("zero registry: (%q, %v)", cursor, err)
	}
}

func TestHandlerRegistryExecuteActionsRejectsInvalidInput(t *testing.T) {
	for _, tc := range []struct {
		name, start string
		change      func(*WorkflowDefinition)
	}{
		{"unknown cursor", "missing", func(*WorkflowDefinition) {}},
		{"empty cursor", "", func(*WorkflowDefinition) {}},
		{"blank cursor", " ", func(*WorkflowDefinition) {}},
		{"broken successor", "z-first", func(d *WorkflowDefinition) { s := d.Steps["a-middle"]; s.Next = "missing"; d.Steps["a-middle"] = s }},
		{"invalid configured start", "z-first", func(d *WorkflowDefinition) { d.StartStep = "missing" }},
		{"unsupported type", "z-first", func(d *WorkflowDefinition) { d.Steps["unreached"] = StepDefinition{Type: "unsupported"} }},
		{"cycle", "z-first", func(d *WorkflowDefinition) { s := d.Steps["m-last"]; s.Next = "z-first"; d.Steps["m-last"] = s }},
		{"missing identity", "z-first", func(d *WorkflowDefinition) { d.ID = "" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewHandlerRegistry()
			if err := r.RegisterAction("load", func() error { t.Fatal("invalid input invoked action"); return nil }); err != nil {
				t.Fatal(err)
			}
			d := actionChainDefinition()
			tc.change(&d)
			if _, err := r.ExecuteActions(d, tc.start); err == nil {
				t.Fatal("invalid input accepted")
			}
		})
	}
}

func TestHandlerRegistryExecuteActionsBoundary(t *testing.T) {
	for _, kind := range []StepType{StepCondition, StepWait} {
		for _, start := range []string{"before", "boundary"} {
			t.Run(string(kind)+"/"+start, func(t *testing.T) {
				r := NewHandlerRegistry()
				calls := 0
				if err := r.RegisterAction("before", func() error { calls++; return nil }); err != nil {
					t.Fatal(err)
				}
				if err := r.RegisterAction("after", func() error { t.Fatal("crossed boundary"); return nil }); err != nil {
					t.Fatal(err)
				}
				if err := r.RegisterCondition("check", func() (bool, error) { t.Fatal("evaluated condition"); return false, nil }); err != nil {
					t.Fatal(err)
				}
				boundary := StepDefinition{Type: StepWait, Delay: 24 * time.Hour, Next: "after"}
				if kind == StepCondition {
					boundary = StepDefinition{Type: StepCondition, Condition: "check", OnTrue: "after"}
				}
				d := WorkflowDefinition{ID: "boundary", Trigger: TriggerDefinition{EventType: "event"}, StartStep: "before", Steps: map[string]StepDefinition{
					"before": {Type: StepAction, Action: "before", Next: "boundary"}, "boundary": boundary, "after": {Type: StepAction, Action: "after"},
				}}
				cursor, err := r.ExecuteActions(d, start)
				wantCalls := 0
				if start == "before" {
					wantCalls = 1
				}
				if cursor != "boundary" || err != nil || calls != wantCalls {
					t.Fatalf("got (%q, %v), calls=%d", cursor, err, calls)
				}
			})
		}
	}
}

func TestHandlerRegistryExecuteActionsStartsAtSuppliedStep(t *testing.T) {
	r := NewHandlerRegistry()
	if err := r.RegisterAction("store", func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	// Earlier handlers are intentionally absent: this attempt starts at the last action.
	if cursor, err := r.ExecuteActions(actionChainDefinition(), "m-last"); cursor != "" || err != nil {
		t.Fatalf("got (%q, %v)", cursor, err)
	}
}
