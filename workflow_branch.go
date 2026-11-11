package enrollment

import "fmt"

// SelectConditionSuccessor evaluates one condition and returns its configured
// successor. An empty successor with nil error ends the path; any error returns
// no successor. Definition validation remains the caller's responsibility.
// Selection does not change configuration, registrations or execution state;
// callbacks retain responsibility for their own side effects and concurrency.
func (r *HandlerRegistry) SelectConditionSuccessor(step StepDefinition) (string, error) {
	if step.Type != StepCondition {
		return "", fmt.Errorf("cannot select condition successor for step type %q", step.Type)
	}
	result, err := r.EvaluateCondition(step.Condition)
	if err != nil {
		return "", err
	}
	if result {
		return step.OnTrue, nil
	}
	return step.OnFalse, nil
}
