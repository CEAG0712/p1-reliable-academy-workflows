package enrollment

import (
	"errors"
	"fmt"
	"maps"
	"time"
)

var (
	ErrWaitNotReady = errors.New("wait deadline has not arrived")
	ErrWaitConsumed = errors.New("wait has already been consumed")
)

// WaitSnapshot is a read-only value copy of local execution progress.
// Cursor is the WAIT ID until eligible resume, then the failed action, an
// unexecuted condition/WAIT boundary, or empty for local path completion.
// Check ResumeWait's error first; an empty cursor alone does not prove success.
type WaitSnapshot struct {
	WaitID   string
	Cursor   string
	ResumeAt time.Time
	Consumed bool
}

// WaitState retains one wait in this process only; it is not durable state.
// Retain the returned pointer: do not copy it or use it concurrently (including
// Snapshot). Definitions must not be mutated concurrently with PrepareWait.
// The private definition snapshot prevents later caller mutations or unrelated
// definitions from changing the saved continuation. Handler mappings/captured
// callback inputs are not snapshotted; registry constraints still apply.
type WaitState struct {
	definition WorkflowDefinition
	progress   WaitSnapshot
}

// PrepareWait validates configuration and records the deadline once, without
// invoking handlers or sleeping. The caller supplies a reached WAIT ID and time;
// this helper validates the ID/type, not the provenance of earlier execution.
// Each call creates an independent wait, not a deduplicated execution identity.
func PrepareWait(d WorkflowDefinition, waitID string, now time.Time) (*WaitState, error) {
	if err := d.Validate(); err != nil {
		return nil, fmt.Errorf("cannot prepare wait: %w", err)
	}
	step, exists := d.Steps[waitID]
	if !exists || step.Type != StepWait {
		return nil, fmt.Errorf("step %q is not a configured WAIT", waitID)
	}
	d.Steps = maps.Clone(d.Steps)
	return &WaitState{
		definition: d,
		progress: WaitSnapshot{
			WaitID: waitID, Cursor: waitID, ResumeAt: now.Add(step.Delay),
		},
	}, nil
}

// Snapshot returns progress by value, so editing it cannot reset consumption.
// Like other methods on WaitState, it requires a nonnil pointer from PrepareWait.
func (w *WaitState) Snapshot() WaitSnapshot {
	return w.progress
}

// ResumeWait refuses early attempts without mutation. Equality is eligible.
// It consumes the wait BEFORE callbacks; errors (including unknown handlers)
// and panics never make this wait retryable. On execution error the cursor
// identifies the failed action and the error preserves the callback cause.
// Earlier actions or even the failed callback may already have side effects:
// callers must reconcile those externally, not reset/reprepare for blind retry.
// No automatic retry, compensation, panic recovery or M2 transition occurs.
// On nil error, empty cursor completes the local path; nonempty cursor hands
// off an unexecuted condition or WAIT. A later WAIT requires separate preparation.
func (r *HandlerRegistry) ResumeWait(w *WaitState, now time.Time) (string, error) {
	if w == nil || w.progress.WaitID == "" {
		return "", fmt.Errorf("wait state must come from PrepareWait")
	}
	if w.progress.Consumed {
		return w.progress.Cursor, ErrWaitConsumed
	}
	if now.Before(w.progress.ResumeAt) {
		return w.progress.Cursor, ErrWaitNotReady
	}
	w.progress.Consumed = true
	w.progress.Cursor = w.definition.Steps[w.progress.WaitID].Next
	if w.progress.Cursor == "" {
		return "", nil
	}
	cursor, err := r.ExecuteActions(w.definition, w.progress.Cursor)
	w.progress.Cursor = cursor
	return cursor, err
}
