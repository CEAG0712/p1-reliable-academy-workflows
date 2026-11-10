package enrollment

import (
	"fmt"
	"strings"
)

// ActionHandler performs one action; closures may capture adapter inputs.
type ActionHandler func() error

// ConditionHandler returns a predicate result and any evaluation error.
// Callers must check the error before using the result.
type ConditionHandler func() (bool, error)

// HandlerRegistry dispatches callbacks by exact, case-sensitive names. Action
// and condition namespaces are independent. Nonblank names are stored verbatim.
// Its zero value is usable. It must not be copied after registration, and
// registration must not overlap other registry operations across goroutines.
// Callback concurrency safety is the caller's responsibility.
type HandlerRegistry struct {
	actions    map[string]ActionHandler
	conditions map[string]ConditionHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{}
}

// RegisterAction rejects blank names, nil callbacks and duplicates without
// replacing an existing callback. Registration never invokes the callback.
func (r *HandlerRegistry) RegisterAction(name string, handler ActionHandler) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("action handler name is required")
	}
	if handler == nil {
		return fmt.Errorf("action handler %q is nil", name)
	}
	if _, exists := r.actions[name]; exists {
		return fmt.Errorf("action handler %q is already registered", name)
	}
	if r.actions == nil {
		r.actions = make(map[string]ActionHandler)
	}
	r.actions[name] = handler
	return nil
}

// RegisterCondition applies the same registration rules in its own namespace.
func (r *HandlerRegistry) RegisterCondition(name string, handler ConditionHandler) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("condition handler name is required")
	}
	if handler == nil {
		return fmt.Errorf("condition handler %q is nil", name)
	}
	if _, exists := r.conditions[name]; exists {
		return fmt.Errorf("condition handler %q is already registered", name)
	}
	if r.conditions == nil {
		r.conditions = make(map[string]ConditionHandler)
	}
	r.conditions[name] = handler
	return nil
}

// InvokeAction calls only the named action once and returns its error unchanged.
func (r *HandlerRegistry) InvokeAction(name string) error {
	handler, exists := r.actions[name]
	if !exists {
		return fmt.Errorf("unknown action handler %q", name)
	}
	return handler()
}

// EvaluateCondition returns the named callback's result and error unchanged.
// Unknown names return false/error. No workflow branch is selected here.
func (r *HandlerRegistry) EvaluateCondition(name string) (bool, error) {
	handler, exists := r.conditions[name]
	if !exists {
		return false, fmt.Errorf("unknown condition handler %q", name)
	}
	return handler()
}
