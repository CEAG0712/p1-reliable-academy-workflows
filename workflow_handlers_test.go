package enrollment

import (
	"errors"
	"reflect"
	"testing"
)

func TestHandlerRegistryNamedInvocation(t *testing.T) {
	r := NewHandlerRegistry()
	var calls []string
	for _, name := range []string{"document.archive", "sensor.alert"} {
		if err := r.RegisterAction(name, func() error {
			calls = append(calls, name)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"document.ready", "sensor.ready"} {
		if err := r.RegisterCondition(name, func() (bool, error) {
			calls = append(calls, name)
			return name == "document.ready", nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	if len(calls) != 0 {
		t.Fatal("registration invoked a callback")
	}
	for _, name := range []string{"sensor.alert", "document.archive", "sensor.alert"} {
		if err := r.InvokeAction(name); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"sensor.ready", "document.ready"} {
		got, err := r.EvaluateCondition(name)
		if err != nil || got != (name == "document.ready") {
			t.Fatalf("%s: got (%v, %v)", name, got, err)
		}
	}
	want := []string{"sensor.alert", "document.archive", "sensor.alert", "sensor.ready", "document.ready"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestHandlerRegistryRejectsUnknownNames(t *testing.T) {
	r := NewHandlerRegistry()
	calls := 0
	if err := r.RegisterAction("action", func() error { calls++; return nil }); err != nil {
		t.Fatal(err)
	}
	if err := r.RegisterCondition("condition", func() (bool, error) { calls++; return true, nil }); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"missing", "", " \t", "\u2003", "Action", " action ", "condition"} {
		if err := r.InvokeAction(name); err == nil {
			t.Fatalf("action %q: expected error", name)
		}
	}
	for _, name := range []string{"missing", "", " \t", "\u2003", "Condition", " condition ", "action"} {
		if got, err := r.EvaluateCondition(name); got || err == nil {
			t.Fatalf("condition %q: got (%v, %v), want false/error", name, got, err)
		}
	}
	if calls != 0 {
		t.Fatal("unknown lookup invoked a callback")
	}
}

func TestHandlerRegistryRejectsInvalidRegistration(t *testing.T) {
	r := NewHandlerRegistry()
	action := ActionHandler(func() error { return nil })
	condition := ConditionHandler(func() (bool, error) { return true, nil })
	for _, name := range []string{"", " ", "\t\n", "\u2003"} {
		if err := r.RegisterAction(name, action); err == nil {
			t.Fatalf("blank action %q accepted", name)
		}
		if err := r.RegisterCondition(name, condition); err == nil {
			t.Fatalf("blank condition %q accepted", name)
		}
	}
	var nilAction ActionHandler
	var nilCondition ConditionHandler
	if err := r.RegisterAction("valid", nilAction); err == nil {
		t.Fatal("nil action accepted")
	}
	if err := r.RegisterCondition("valid", nilCondition); err == nil {
		t.Fatal("nil condition accepted")
	}
	// Rejected registration must not reserve the name.
	if err := r.RegisterAction("valid", action); err != nil {
		t.Fatal(err)
	}
	if err := r.RegisterCondition("valid", condition); err != nil {
		t.Fatal(err)
	}
}

func TestHandlerRegistryDuplicatesPreserveOriginal(t *testing.T) {
	r := NewHandlerRegistry()
	actionCalls, conditionCalls := 0, 0
	if err := r.RegisterAction("shared", func() error { actionCalls++; return nil }); err != nil {
		t.Fatal(err)
	}
	if err := r.RegisterCondition("shared", func() (bool, error) { conditionCalls++; return true, nil }); err != nil {
		t.Fatal(err)
	}
	if err := r.RegisterAction("shared", func() error { t.Fatal("replacement action called"); return nil }); err == nil {
		t.Fatal("duplicate action accepted")
	}
	if err := r.RegisterCondition("shared", func() (bool, error) {
		t.Fatal("replacement condition called")
		return false, nil
	}); err == nil {
		t.Fatal("duplicate condition accepted")
	}
	if actionCalls != 0 || conditionCalls != 0 {
		t.Fatal("registration invoked a callback")
	}
	if err := r.InvokeAction("shared"); err != nil {
		t.Fatal(err)
	}
	if got, err := r.EvaluateCondition("shared"); !got || err != nil {
		t.Fatalf("got (%v, %v), want true/nil", got, err)
	}
	if actionCalls != 1 || conditionCalls != 1 {
		t.Fatalf("original calls = %d/%d, want 1/1", actionCalls, conditionCalls)
	}
}

func TestHandlerRegistryPropagatesErrors(t *testing.T) {
	r := NewHandlerRegistry()
	want := errors.New("adapter failed")
	if err := r.RegisterAction("fail", func() error { return want }); err != nil {
		t.Fatal(err)
	}
	if got := r.InvokeAction("fail"); got != want {
		t.Fatalf("action error = %v, want original %v", got, want)
	}
	for _, result := range []bool{false, true} {
		local := NewHandlerRegistry()
		if err := local.RegisterCondition("fail", func() (bool, error) { return result, want }); err != nil {
			t.Fatal(err)
		}
		if got, err := local.EvaluateCondition("fail"); got != result || err != want {
			t.Fatalf("condition = (%v, %v), want (%v, original error)", got, err, result)
		}
	}
}

func TestHandlerRegistryZeroValueAndIsolation(t *testing.T) {
	var first, second HandlerRegistry
	if err := second.InvokeAction("opaque"); err == nil {
		t.Fatal("empty registry accepted unknown action")
	}
	if got, err := second.EvaluateCondition("opaque"); got || err == nil {
		t.Fatal("empty registry accepted unknown condition")
	}
	for _, r := range []*HandlerRegistry{&first, &second} {
		for _, name := range []string{"opaque", "Opaque", " opaque "} {
			if err := r.RegisterAction(name, func() error { return nil }); err != nil {
				t.Fatal(err)
			}
			if err := r.RegisterCondition(name, func() (bool, error) { return true, nil }); err != nil {
				t.Fatal(err)
			}
			if err := r.InvokeAction(name); err != nil {
				t.Fatal(err)
			}
			if got, err := r.EvaluateCondition(name); !got || err != nil {
				t.Fatalf("%q: got (%v, %v)", name, got, err)
			}
		}
	}
}
