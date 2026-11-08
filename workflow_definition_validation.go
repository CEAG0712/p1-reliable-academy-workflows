package enrollment

import (
	"fmt"
	"sort"
	"strings"
)

// Validate checks the entire configured graph without changing it. It returns
// the first error: definition fields first, then steps in sorted ID order, then
// cycles. Empty successors end a path, including either condition branch.
// Handler names are required references; registration is not checked here.
func (d WorkflowDefinition) Validate() error {
	if strings.TrimSpace(d.ID) == "" {
		return fmt.Errorf("workflow ID is required")
	}
	if strings.TrimSpace(d.Trigger.EventType) == "" {
		return fmt.Errorf("trigger event type is required")
	}
	if strings.TrimSpace(d.StartStep) == "" {
		return fmt.Errorf("start step is required")
	}
	if _, exists := d.Steps[d.StartStep]; !exists {
		return fmt.Errorf("start step %q does not exist", d.StartStep)
	}

	ids := make([]string, 0, len(d.Steps))
	for id := range d.Steps {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("step ID is required")
		}
		step := d.Steps[id]
		switch step.Type {
		case StepAction:
			if strings.TrimSpace(step.Action) == "" {
				return fmt.Errorf("step %q: action handler name is required", id)
			}
			if step.Condition != "" || step.OnTrue != "" || step.OnFalse != "" || step.Delay != 0 {
				return fmt.Errorf("step %q: action cannot configure Condition, OnTrue, OnFalse or Delay", id)
			}
		case StepCondition:
			if strings.TrimSpace(step.Condition) == "" {
				return fmt.Errorf("step %q: condition handler name is required", id)
			}
			if step.Action != "" || step.Next != "" || step.Delay != 0 {
				return fmt.Errorf("step %q: condition cannot configure Action, Next or Delay", id)
			}
		case StepWait:
			if step.Delay <= 0 {
				return fmt.Errorf("step %q: wait delay must be positive", id)
			}
			if step.Action != "" || step.Condition != "" || step.OnTrue != "" || step.OnFalse != "" {
				return fmt.Errorf("step %q: wait cannot configure Action, Condition, OnTrue or OnFalse", id)
			}
		default:
			return fmt.Errorf("step %q: unsupported type %q", id, step.Type)
		}

		// Fixed field order makes diagnostics independent of map iteration.
		for _, link := range []struct{ field, target string }{
			{"Next", step.Next}, {"OnTrue", step.OnTrue}, {"OnFalse", step.OnFalse},
		} {
			if link.target != "" {
				if _, exists := d.Steps[link.target]; !exists {
					return fmt.Errorf("step %q: %s references missing step %q", id, link.field, link.target)
				}
			}
		}
	}

	// A node on the active DFS path is a cycle; a previously finished node is
	// a valid shared successor. Check disconnected components as well.
	const (
		visiting = 1
		finished = 2
	)
	state := make(map[string]int, len(d.Steps))
	var visit func(string) error
	visit = func(id string) error {
		if id == "" || state[id] == finished {
			return nil
		}
		if state[id] == visiting {
			return fmt.Errorf("unsupported cycle at step %q", id)
		}
		state[id] = visiting
		step := d.Steps[id]
		for _, next := range []string{step.Next, step.OnTrue, step.OnFalse} {
			if err := visit(next); err != nil {
				return err
			}
		}
		state[id] = finished
		return nil
	}
	for _, id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
