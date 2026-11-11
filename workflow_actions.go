package enrollment

import "fmt"

// ExecuteActions validates the entire definition, then follows Next from startID.
// An empty cursor with nil error means the action path ended. A nonempty cursor
// with nil error hands off an unexecuted condition or WAIT step. On action error,
// the cursor identifies the failed step and the error wraps the underlying cause.
// Invalid input returns startID and an error; an empty starting ID is invalid.
// Progress is local to this call: no lifecycle, persistence or deduplication.
// Callers and callbacks must not mutate the definition during traversal and must
// respect HandlerRegistry's registration and callback concurrency constraints.
func (r *HandlerRegistry) ExecuteActions(d WorkflowDefinition, startID string) (string, error) {
	if err := d.Validate(); err != nil {
		return startID, fmt.Errorf("cannot execute actions: %w", err)
	}
	if _, exists := d.Steps[startID]; !exists {
		return startID, fmt.Errorf("starting step %q does not exist", startID)
	}
	for cursor := startID; ; {
		step := d.Steps[cursor]
		if step.Type == StepCondition || step.Type == StepWait {
			return cursor, nil
		}
		// Validation guarantees that every other step is an action and all
		// nonempty successors exist, with no cycles anywhere in the graph.
		if err := r.InvokeAction(step.Action); err != nil {
			return cursor, fmt.Errorf("action step %q: %w", cursor, err)
		}
		cursor = step.Next
		if cursor == "" {
			return "", nil
		}
	}
}
