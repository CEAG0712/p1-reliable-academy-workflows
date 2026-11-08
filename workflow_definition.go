package enrollment

import "time"

// WorkflowDefinition describes configuration, separately from execution state.
// Step IDs are the keys in Steps; an empty successor denotes the end of a path.
// Validate checks configuration; definitions do not execute handlers.
type WorkflowDefinition struct {
	ID        string
	Trigger   TriggerDefinition
	StartStep string
	Steps     map[string]StepDefinition
}

// TriggerDefinition names an event type without coupling it to a business payload.
type TriggerDefinition struct {
	EventType string
}

type StepType string

const (
	StepAction    StepType = "action"
	StepCondition StepType = "condition"
	StepWait      StepType = "wait"
)

// StepDefinition describes one configured step. Action and Condition are handler
// names, not implementations. Action and wait steps use Next; conditions use
// OnTrue and OnFalse. Delay is configuration, not a runtime resume deadline.
// Validate checks type-specific constraints.
type StepDefinition struct {
	Type      StepType
	Action    string
	Condition string
	Next      string
	OnTrue    string
	OnFalse   string
	Delay     time.Duration
}
