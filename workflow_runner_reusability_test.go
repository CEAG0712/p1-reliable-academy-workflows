package enrollment

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// Business inputs and effects belong to test-local callbacks, not engine fields.
type reusabilityFixture struct {
	registry                           *HandlerRegistry
	definitions                        []WorkflowDefinition
	calls                              []string
	accessRequest, fulfillment, report string
	metrics                            []int
	total                              int
}

func newReusabilityFixture(t *testing.T, beforeAction func(string), aggregateError error) *reusabilityFixture {
	t.Helper()
	f := &reusabilityFixture{registry: NewHandlerRegistry()}
	f.definitions = []WorkflowDefinition{
		{ID: "enrollment-fulfillment", Trigger: TriggerDefinition{EventType: "purchase.confirmed"}, StartStep: "validate-purchase", Steps: map[string]StepDefinition{
			"validate-purchase":  {Type: StepAction, Action: "enrollment.validate", Next: "prepare-access"},
			"prepare-access":     {Type: StepAction, Action: "enrollment.prepare", Next: "record-fulfillment"},
			"record-fulfillment": {Type: StepAction, Action: "enrollment.record"},
		}},
		{ID: "learning-analytics", Trigger: TriggerDefinition{EventType: "report.requested"}, StartStep: "collect-activity", Steps: map[string]StepDefinition{
			"collect-activity":   {Type: StepAction, Action: "analytics.collect", Next: "aggregate-activity"},
			"aggregate-activity": {Type: StepAction, Action: "analytics.aggregate", Next: "render-report"},
			"render-report":      {Type: StepAction, Action: "analytics.render", Next: "record-publication"},
			"record-publication": {Type: StepAction, Action: "analytics.publish"},
		}},
	}
	student, course, paid := "student-17", "go-basics", true
	validated := false
	rendered := ""
	handlers := []struct {
		name   string
		action ActionHandler
	}{
		{"enrollment.validate", func() error {
			if student == "" || course == "" || !paid {
				return errors.New("ineligible purchase")
			}
			validated = true
			return nil
		}},
		{"enrollment.prepare", func() error {
			if !validated {
				return errors.New("purchase was not validated")
			}
			f.accessRequest = student + ":" + course
			return nil
		}},
		{"enrollment.record", func() error {
			if f.accessRequest == "" {
				return errors.New("missing access request")
			}
			f.fulfillment = "prepared:" + f.accessRequest
			return nil
		}},
		{"analytics.collect", func() error { f.metrics = []int{2, 3, 5}; return nil }},
		{"analytics.aggregate", func() error {
			if aggregateError != nil {
				return aggregateError
			}
			if len(f.metrics) != 3 {
				return errors.New("missing activity metrics")
			}
			for _, value := range f.metrics {
				f.total += value
			}
			return nil
		}},
		{"analytics.render", func() error {
			if f.total != 10 {
				return errors.New("unexpected activity total")
			}
			rendered = "learning activity: 10"
			return nil
		}},
		{"analytics.publish", func() error {
			if rendered == "" {
				return errors.New("missing rendered report")
			}
			f.report = "published:" + rendered
			return nil
		}},
	}
	for _, handler := range handlers {
		if err := f.registry.RegisterAction(handler.name, func() error {
			beforeAction(handler.name)
			f.calls = append(f.calls, handler.name)
			return handler.action()
		}); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range f.definitions {
		if err := d.Validate(); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func reusabilitySnapshot(t *testing.T, repo *InMemoryWorkflowRepository, id string) *WorkflowExecution {
	t.Helper()
	s, ok := repo.Get(id)
	if !ok {
		t.Fatalf("execution %q missing", id)
	}
	return s
}

func assertReusabilityLifecycle(t *testing.T, repo *InMemoryWorkflowRepository, id, worker, status, reason string) {
	t.Helper()
	s := reusabilitySnapshot(t, repo, id)
	version, attempts := uint64(1), 0
	var history []WorkflowTransition
	if status != "READY" {
		version, attempts = 2, 1
		history = append(history, WorkflowTransition{From: "READY", To: "RUNNING", WorkerID: worker, Version: 2})
	}
	if status == "COMPLETED" || status == "FAILED" {
		version = 3
		history = append(history, WorkflowTransition{From: "RUNNING", To: status, WorkerID: worker, Version: 3, Reason: reason})
	}
	if s.WorkflowID != id || s.Status != status || s.ClaimOwner != worker || s.Version != version || s.Attempts != attempts || s.LastError != reason || !reflect.DeepEqual(s.History, history) {
		t.Fatalf("unexpected execution: %+v history=%+v; want %s owner=%q version=%d attempts=%d reason=%q history=%+v", s, s.History, status, worker, version, attempts, reason, history)
	}
}

func reusabilityCalls(index int) []string {
	if index == 0 {
		return []string{"enrollment.validate", "enrollment.prepare", "enrollment.record"}
	}
	return []string{"analytics.collect", "analytics.aggregate", "analytics.render", "analytics.publish"}
}

func TestRunnerReusabilityWorkflows(t *testing.T) {
	for _, first := range []int{0, 1} {
		t.Run([]string{"enrollment-first", "analytics-first"}[first], func(t *testing.T) {
			repo := NewInMemoryWorkflowRepository()
			ids := []string{"purchase-42", "report-2026-11"}
			workers := []string{"fulfillment-worker", "report-worker"}
			current := -1
			f := newReusabilityFixture(t, func(name string) {
				assertReusabilityLifecycle(t, repo, ids[current], workers[current], "RUNNING", "")
				if !strings.HasPrefix(name, []string{"enrollment.", "analytics."}[current]) {
					t.Fatalf("cross-workflow callback %q", name)
				}
			}, nil)
			for _, id := range ids {
				repo.Save(NewWorkflowExecution(id))
				assertReusabilityLifecycle(t, repo, id, "", "READY", "")
			}
			var wantCalls []string
			for _, index := range []int{first, 1 - first} {
				current = index
				otherBefore := reusabilitySnapshot(t, repo, ids[1-index])
				cursor, err := f.registry.RunActions(repo, ids[index], workers[index], f.definitions[index])
				if err != nil || cursor != "" {
					t.Fatalf("workflow %d: (%q, %v)", index, cursor, err)
				}
				wantCalls = append(wantCalls, reusabilityCalls(index)...)
				if !reflect.DeepEqual(f.calls, wantCalls) {
					t.Fatalf("calls=%v want=%v", f.calls, wantCalls)
				}
				assertReusabilityLifecycle(t, repo, ids[index], workers[index], "COMPLETED", "")
				if !reflect.DeepEqual(otherBefore, reusabilitySnapshot(t, repo, ids[1-index])) {
					t.Fatal("other execution changed")
				}
				if index == first && ((index == 0 && f.report != "") || (index == 1 && (f.accessRequest != "" || f.fulfillment != ""))) {
					t.Fatal("other business output changed")
				}
			}
			if f.accessRequest != "student-17:go-basics" || f.fulfillment != "prepared:student-17:go-basics" || f.total != 10 || f.report != "published:learning activity: 10" {
				t.Fatalf("wrong business outputs: %+v", f)
			}
		})
	}
}

func TestRunnerReusabilityActionShapes(t *testing.T) {
	// Exercise the underlying action API independently, without lifecycle claims.
	for index, name := range []string{"three-action-enrollment", "four-action-analytics"} {
		t.Run(name, func(t *testing.T) {
			f := newReusabilityFixture(t, func(string) {}, nil)
			d := f.definitions[index]
			if len(d.Steps) != 3+index {
				t.Fatal("fixture lost its distinct chain length")
			}
			cursor, err := f.registry.ExecuteActions(d, d.StartStep)
			if err != nil || cursor != "" || !reflect.DeepEqual(f.calls, reusabilityCalls(index)) {
				t.Fatalf("cursor=%q err=%v calls=%v", cursor, err, f.calls)
			}
		})
	}
}

func TestRunnerReusabilityIndependentFailure(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	ids := []string{"purchase-independent", "report-failed"}
	workers := []string{"access-worker", "analytics-worker"}
	for _, id := range ids {
		repo.Save(NewWorkflowExecution(id))
	}
	cause := errors.New("invalid activity batch")
	current := 1
	f := newReusabilityFixture(t, func(string) {
		assertReusabilityLifecycle(t, repo, ids[current], workers[current], "RUNNING", "")
	}, cause)
	ready := reusabilitySnapshot(t, repo, ids[0])
	cursor, err := f.registry.RunActions(repo, ids[1], workers[1], f.definitions[1])
	if cursor != "aggregate-activity" || !errors.Is(err, cause) {
		t.Fatalf("failure=(%q, %v)", cursor, err)
	}
	reason := err.Error()
	if !strings.Contains(reason, cause.Error()) || !strings.Contains(reason, cursor) || !strings.Contains(reason, ids[1]) {
		t.Fatalf("incomplete diagnostic: %q", reason)
	}
	assertReusabilityLifecycle(t, repo, ids[1], workers[1], "FAILED", reason)
	if !reflect.DeepEqual(f.calls, []string{"analytics.collect", "analytics.aggregate"}) || f.report != "" || f.total != 0 || f.accessRequest != "" || f.fulfillment != "" {
		t.Fatal("failure continued or contaminated enrollment")
	}
	if !reflect.DeepEqual(ready, reusabilitySnapshot(t, repo, ids[0])) {
		t.Fatal("failure changed other execution")
	}
	failed := reusabilitySnapshot(t, repo, ids[1])
	current = 0
	cursor, err = f.registry.RunActions(repo, ids[0], workers[0], f.definitions[0])
	if err != nil || cursor != "" {
		t.Fatalf("independent enrollment=(%q, %v)", cursor, err)
	}
	assertReusabilityLifecycle(t, repo, ids[0], workers[0], "COMPLETED", "")
	wantCalls := append([]string{"analytics.collect", "analytics.aggregate"}, reusabilityCalls(0)...)
	if !reflect.DeepEqual(f.calls, wantCalls) || f.fulfillment != "prepared:student-17:go-basics" || f.report != "" || !reflect.DeepEqual(failed, reusabilitySnapshot(t, repo, ids[1])) {
		t.Fatal("independent completion corrupted failure or outputs")
	}
}
