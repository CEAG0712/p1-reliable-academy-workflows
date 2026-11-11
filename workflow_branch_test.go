package enrollment

import (
	"errors"
	"testing"
)

func TestHandlerRegistrySelectConditionSuccessor(t *testing.T) {
	wantErr := errors.New("predicate failed")
	for _, tc := range []struct {
		name            string
		result          bool
		err             error
		onTrue, onFalse string
		want            string
	}{
		{"true", true, nil, "yes", "no", "yes"},
		{"false", false, nil, "yes", "no", "no"},
		{"true error", true, wantErr, "yes", "no", ""},
		{"false error", false, wantErr, "yes", "no", ""},
		{"true terminal", true, nil, "", "no", ""},
		{"false terminal", false, nil, "yes", "", ""},
		{"both terminal true", true, nil, "", "", ""},
		{"both terminal false", false, nil, "", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewHandlerRegistry()
			calls := 0
			if err := r.RegisterCondition("check", func() (bool, error) {
				calls++
				return tc.result, tc.err
			}); err != nil {
				t.Fatal(err)
			}
			step := StepDefinition{
				Type: StepCondition, Condition: "check",
				OnTrue: tc.onTrue, OnFalse: tc.onFalse,
			}
			before := step
			for attempt := 1; attempt <= 2; attempt++ {
				got, err := r.SelectConditionSuccessor(step)
				if got != tc.want || err != tc.err {
					t.Fatalf("got (%q, %v), want (%q, %v)", got, err, tc.want, tc.err)
				}
				if calls != attempt {
					t.Fatalf("calls = %d after %d attempts", calls, attempt)
				}
			}
			if step != before || len(r.conditions) != 1 || len(r.actions) != 0 {
				t.Fatal("selection changed step or registry registration counts")
			}
			if err := r.RegisterCondition("check", func() (bool, error) {
				t.Fatal("replacement invoked")
				return false, nil
			}); err == nil {
				t.Fatal("selection removed the original registration")
			}
		})
	}
}

func TestHandlerRegistrySelectConditionSuccessorUnknown(t *testing.T) {
	r := NewHandlerRegistry()
	if err := r.RegisterCondition("check", func() (bool, error) {
		t.Fatal("unknown lookup invoked registered condition")
		return true, nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"missing", "", "Check", " check "} {
		step := StepDefinition{Type: StepCondition, Condition: name, OnTrue: "yes", OnFalse: "no"}
		_, wantErr := r.EvaluateCondition(name)
		if got, err := r.SelectConditionSuccessor(step); got != "" || err == nil || err.Error() != wantErr.Error() {
			t.Fatalf("%q: got (%q, %v), want empty successor and %v", name, got, err, wantErr)
		}
	}
	var zero HandlerRegistry
	if got, err := zero.SelectConditionSuccessor(StepDefinition{Type: StepCondition, Condition: "missing"}); got != "" || err == nil {
		t.Fatalf("zero registry: got (%q, %v), want empty successor/error", got, err)
	}
}

func TestHandlerRegistrySelectConditionSuccessorRejectsWrongType(t *testing.T) {
	r := NewHandlerRegistry()
	if err := r.RegisterCondition("check", func() (bool, error) {
		t.Fatal("wrong step type invoked condition")
		return true, nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, kind := range []StepType{StepAction, StepWait, "", "unsupported"} {
		step := StepDefinition{Type: kind, Condition: "check", OnTrue: "yes", OnFalse: "no"}
		if got, err := r.SelectConditionSuccessor(step); got != "" || err == nil {
			t.Fatalf("%q: got (%q, %v), want empty successor/error", kind, got, err)
		}
	}
}
