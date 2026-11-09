package enrollment

import (
	"reflect"
	"testing"
)

func triggerDefinition() WorkflowDefinition {
	return WorkflowDefinition{
		ID: "archive", Trigger: TriggerDefinition{EventType: "document.approved"},
		StartStep: "start",
		Steps: map[string]StepDefinition{
			"start": {Type: StepAction, Action: "document.archive"},
		},
	}
}

func TestWorkflowDefinitionMatchesEventTypes(t *testing.T) {
	archive := triggerDefinition()
	alert := triggerDefinition()
	alert.ID = "alert"
	alert.Trigger.EventType = "sensor.threshold"
	alert.Steps["start"] = StepDefinition{Type: StepAction, Action: "alert.send"}
	for _, d := range []WorkflowDefinition{archive, alert} {
		for _, eventType := range []string{"document.approved", "sensor.threshold", "unknown", "Document.approved", " document.approved "} {
			got, err := d.MatchesEvent(Event{Type: eventType})
			want := eventType == d.Trigger.EventType
			if err != nil || got != want {
				t.Fatalf("%s / %q: got (%v, %v), want (%v, nil)", d.ID, eventType, got, err, want)
			}
		}
	}
	// Nonblank types are opaque, just as in the existing definition contract.
	archive.Trigger.EventType = " document.approved "
	if got, err := archive.MatchesEvent(Event{Type: " document.approved "}); err != nil || !got {
		t.Fatalf("identical padded types: got (%v, %v), want (true, nil)", got, err)
	}
}

func TestWorkflowDefinitionMatchesEventRejectsBlankTypes(t *testing.T) {
	for _, eventType := range []string{"", " ", "\t\n", "\u2003"} {
		for _, d := range []WorkflowDefinition{triggerDefinition(), {}} {
			got, err := d.MatchesEvent(Event{Type: eventType})
			if got || err == nil || err.Error() != "event type is required" {
				t.Fatalf("%q: got (%v, %v), want event error before definition validation", eventType, got, err)
			}
		}
	}
}

func TestWorkflowDefinitionMatchesEventRejectsInvalidDefinitions(t *testing.T) {
	for _, change := range []struct {
		name string
		edit func(*WorkflowDefinition)
		want string
	}{
		{"missing ID", func(d *WorkflowDefinition) { d.ID = "" }, "workflow ID is required"},
		{"blank trigger", func(d *WorkflowDefinition) { d.Trigger.EventType = " " }, "trigger event type is required"},
		{"broken graph", func(d *WorkflowDefinition) { d.StartStep = "missing" }, `start step "missing" does not exist`},
	} {
		t.Run(change.name, func(t *testing.T) {
			d := triggerDefinition()
			change.edit(&d)
			// Invalid definitions are errors even when their trigger would not match.
			for _, eventType := range []string{"document.approved", "unknown"} {
				got, err := d.MatchesEvent(Event{Type: eventType})
				if got || err == nil || err.Error() != change.want {
					t.Fatalf("%q: got (%v, %v), want (false, %q)", eventType, got, err, change.want)
				}
			}
		})
	}
}

func TestWorkflowDefinitionMatchesEventIsDeterministicAndReadOnly(t *testing.T) {
	for _, invalid := range []bool{false, true} {
		d := triggerDefinition()
		before := triggerDefinition() // Independent Steps map, not a shallow copy.
		if invalid {
			d.StartStep, before.StartStep = "missing", "missing"
		}
		for i := 0; i < 100; i++ {
			for _, eventType := range []string{"document.approved", "unknown", ""} {
				event := Event{Type: eventType}
				got, err := d.MatchesEvent(event)
				wantError := ""
				if eventType == "" {
					wantError = "event type is required"
				} else if invalid {
					wantError = `start step "missing" does not exist`
				}
				if got != (!invalid && eventType == "document.approved") ||
					(err == nil) != (wantError == "") || (err != nil && err.Error() != wantError) {
					t.Fatalf("iteration %d: got (%v, %v), expected error %q", i, got, err, wantError)
				}
				if !reflect.DeepEqual(d, before) || event.Type != eventType {
					t.Fatal("matching changed an input")
				}
			}
		}
	}
}
