package enrollment

import (
	"fmt"
	"strings"
)

// Event is the minimal incoming trigger input. Type is an opaque, nonblank name;
// matching does not trim it, fold case, or interpret business payloads.
type Event struct {
	Type string
}

// MatchesEvent reports eligibility for a future runner, not execution or a claim.
// It checks the event first, then validates the entire definition even on a
// mismatch. Invalid input returns false and an error; a valid mismatch returns
// false and nil. A valid exact match returns true and nil.
// It only reads inputs and has no repository, handler, or dispatch side effects.
func (d WorkflowDefinition) MatchesEvent(event Event) (bool, error) {
	if strings.TrimSpace(event.Type) == "" {
		return false, fmt.Errorf("event type is required")
	}
	if err := d.Validate(); err != nil {
		return false, err
	}
	return event.Type == d.Trigger.EventType, nil
}
