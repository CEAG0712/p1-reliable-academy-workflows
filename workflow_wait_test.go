package enrollment

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func waitDefinition() WorkflowDefinition {
	return WorkflowDefinition{
		ID: "follow-up", Trigger: TriggerDefinition{EventType: "received"}, StartStep: "before",
		Steps: map[string]StepDefinition{
			"before": {Type: StepAction, Action: "before", Next: "wait"},
			"wait":   {Type: StepWait, Delay: 24 * time.Hour, Next: "after"},
			"after":  {Type: StepAction, Action: "after"},
		},
	}
}

func TestWaitPreparationAndEarlyResume(t *testing.T) {
	d := waitDefinition()
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	r := NewHandlerRegistry()
	if err := r.RegisterAction("after", func() error { t.Fatal("early callback"); return nil }); err != nil {
		t.Fatal(err)
	}
	w, err := PrepareWait(d, "wait", now)
	if err != nil {
		t.Fatal(err)
	}
	want := WaitSnapshot{WaitID: "wait", Cursor: "wait", ResumeAt: now.Add(24 * time.Hour)}
	if got := w.Snapshot(); got != want {
		t.Fatalf("prepared = %+v, want %+v", got, want)
	}
	for _, at := range []time.Time{now, now.Add(time.Hour), want.ResumeAt.Add(-time.Nanosecond)} {
		if cursor, err := r.ResumeWait(w, at); cursor != "wait" || !errors.Is(err, ErrWaitNotReady) {
			t.Fatalf("early = (%q, %v)", cursor, err)
		}
		if got := w.Snapshot(); got != want {
			t.Fatalf("early attempt changed state: %+v", got)
		}
	}
}

func TestWaitResumeWithoutReplay(t *testing.T) {
	for _, late := range []time.Duration{0, time.Hour} {
		t.Run(late.String(), func(t *testing.T) {
			d, r := waitDefinition(), NewHandlerRegistry()
			var calls []string
			for _, name := range []string{"before", "after"} {
				if err := r.RegisterAction(name, func() error { calls = append(calls, name); return nil }); err != nil {
					t.Fatal(err)
				}
			}
			cursor, err := r.ExecuteActions(d, d.StartStep)
			if cursor != "wait" || err != nil {
				t.Fatalf("initial traversal = (%q, %v)", cursor, err)
			}
			w, err := PrepareWait(d, cursor, time.Time{})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(calls, []string{"before"}) {
				t.Fatalf("preparation executed successor: %v", calls)
			}
			at := w.Snapshot().ResumeAt.Add(late)
			if cursor, err := r.ResumeWait(w, at); cursor != "" || err != nil {
				t.Fatalf("resume = (%q, %v)", cursor, err)
			}
			saved := w.Snapshot()
			if !saved.Consumed || saved.Cursor != "" {
				t.Fatalf("completed state = %+v", saved)
			}
			if cursor, err := r.ResumeWait(w, at.Add(time.Hour)); cursor != "" || !errors.Is(err, ErrWaitConsumed) {
				t.Fatalf("duplicate = (%q, %v)", cursor, err)
			}
			if w.Snapshot() != saved || !reflect.DeepEqual(calls, []string{"before", "after"}) {
				t.Fatalf("replayed or changed state: %v, %+v", calls, w.Snapshot())
			}
		})
	}
}

func TestWaitResumeBoundaries(t *testing.T) {
	for _, boundary := range []string{"", "condition", "second-wait"} {
		for _, viaAction := range []bool{false, true} {
			t.Run(boundary+"/"+map[bool]string{false: "direct", true: "action"}[viaAction], func(t *testing.T) {
				d, r := waitDefinition(), NewHandlerRegistry()
				d.Steps["condition"] = StepDefinition{Type: StepCondition, Condition: "check"}
				d.Steps["second-wait"] = StepDefinition{Type: StepWait, Delay: time.Hour}
				step := d.Steps["wait"]
				step.Next = boundary
				calls := 0
				if viaAction {
					step.Next = "after"
					d.Steps["after"] = StepDefinition{Type: StepAction, Action: "after", Next: boundary}
				}
				d.Steps["wait"] = step
				if err := r.RegisterAction("after", func() error { calls++; return nil }); err != nil {
					t.Fatal(err)
				}
				if err := r.RegisterCondition("check", func() (bool, error) { t.Fatal("evaluated boundary"); return false, nil }); err != nil {
					t.Fatal(err)
				}
				w, err := PrepareWait(d, "wait", time.Time{})
				if err != nil {
					t.Fatal(err)
				}
				if cursor, err := r.ResumeWait(w, w.Snapshot().ResumeAt); cursor != boundary || err != nil {
					t.Fatalf("boundary = (%q, %v), want %q", cursor, err, boundary)
				}
				wantCalls := 0
				if viaAction {
					wantCalls = 1
				}
				if calls != wantCalls || w.Snapshot().Cursor != boundary || !w.Snapshot().Consumed {
					t.Fatalf("calls = %d, state = %+v", calls, w.Snapshot())
				}
			})
		}
	}
}

func TestWaitRejectsInvalidInput(t *testing.T) {
	for _, id := range []string{"", "missing", "before", "after"} {
		if w, err := PrepareWait(waitDefinition(), id, time.Time{}); w != nil || err == nil {
			t.Fatalf("%q: got (%v, %v)", id, w, err)
		}
	}
	for _, change := range []func(*WorkflowDefinition){
		func(d *WorkflowDefinition) { d.ID = "" },
		func(d *WorkflowDefinition) { d.Steps["wait"] = StepDefinition{Type: StepWait} },
		func(d *WorkflowDefinition) { s := d.Steps["wait"]; s.Next = "missing"; d.Steps["wait"] = s },
		func(d *WorkflowDefinition) { s := d.Steps["wait"]; s.Next = "before"; d.Steps["wait"] = s },
	} {
		d := waitDefinition()
		change(&d)
		if w, err := PrepareWait(d, "wait", time.Time{}); w != nil || err == nil {
			t.Fatalf("invalid definition: (%v, %v)", w, err)
		}
	}
	r := NewHandlerRegistry()
	for _, w := range []*WaitState{nil, {}} {
		if cursor, err := r.ResumeWait(w, time.Time{}); cursor != "" || err == nil {
			t.Fatalf("invalid state: (%q, %v)", cursor, err)
		}
	}
}

func TestWaitDefinitionSnapshot(t *testing.T) {
	d := waitDefinition()
	w, err := PrepareWait(d, "wait", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	// No definition/step argument can rebind this wait. Even mutation of the
	// caller's original map cannot redirect the saved continuation.
	d.ID, d.StartStep = "unrelated", "other"
	d.Steps["wait"] = StepDefinition{Type: StepWait, Delay: time.Hour, Next: "before"}
	d.Steps["after"] = StepDefinition{Type: StepAction, Action: "wrong"}
	delete(d.Steps, "before")
	r := NewHandlerRegistry()
	calls := 0
	if err := r.RegisterAction("after", func() error { calls++; return nil }); err != nil {
		t.Fatal(err)
	}
	if cursor, err := r.ResumeWait(w, w.Snapshot().ResumeAt); cursor != "" || err != nil || calls != 1 {
		t.Fatalf("snapshot resume = (%q, %v), calls %d", cursor, err, calls)
	}
}

func TestWaitFailureConsumesWithoutRetry(t *testing.T) {
	for _, failed := range []string{"after", "last", "missing-handler"} {
		t.Run(failed, func(t *testing.T) {
			d, r := waitDefinition(), NewHandlerRegistry()
			d.Steps["after"] = StepDefinition{Type: StepAction, Action: "after", Next: "last"}
			d.Steps["last"] = StepDefinition{Type: StepAction, Action: "last"}
			cause := errors.New("external result uncertain")
			var calls []string
			for _, name := range []string{"after", "last"} {
				if failed == "missing-handler" && name == "last" {
					continue
				}
				if err := r.RegisterAction(name, func() error {
					calls = append(calls, name)
					if name == failed {
						return cause
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			w, err := PrepareWait(d, "wait", time.Time{})
			if err != nil {
				t.Fatal(err)
			}
			wantCursor, wantCalls := "last", []string{"after", "last"}
			if failed == "after" {
				wantCursor, wantCalls = "after", []string{"after"}
			}
			if failed == "missing-handler" {
				wantCalls = []string{"after"}
			}
			cursor, err := r.ResumeWait(w, w.Snapshot().ResumeAt)
			if cursor != wantCursor || err == nil || (failed != "missing-handler" && !errors.Is(err, cause)) {
				t.Fatalf("failure = (%q, %v)", cursor, err)
			}
			saved := w.Snapshot()
			if !saved.Consumed || saved.Cursor != wantCursor {
				t.Fatalf("failed state = %+v", saved)
			}
			if cursor, err := r.ResumeWait(w, saved.ResumeAt); cursor != wantCursor || !errors.Is(err, ErrWaitConsumed) {
				t.Fatalf("retry = (%q, %v)", cursor, err)
			}
			if w.Snapshot() != saved || !reflect.DeepEqual(calls, wantCalls) {
				t.Fatalf("retry changed execution: %v", calls)
			}
		})
	}
}
