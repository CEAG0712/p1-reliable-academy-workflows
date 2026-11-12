package enrollment

import "fmt"

// WorkflowLifecycleRepository is the owner/version-checked M2 boundary used by
// RunActions. Get must return an independent snapshot. TryClaim must atomically
// acquire READY executions for a nonempty owner, and terminal operations must
// reject a different owner, stale version or non-RUNNING state without mutation.
// InMemoryWorkflowRepository implements this contract; SaveIfVersion is omitted.
type WorkflowLifecycleRepository interface {
	TryClaim(workflowID, workerID string) bool
	Get(workflowID string) (*WorkflowExecution, bool)
	Complete(workflowID, workerID string, expectedVersion uint64) bool
	Fail(workflowID, workerID string, expectedVersion uint64, reason string) bool
}

// RunActions claims an existing execution, then executes from d.StartStep using
// ExecuteActions. executionID identifies a run independently of the definition ID;
// associating that run with the intended definition is the caller's responsibility.
// The registry and repository must be nonnil, initialized dependencies.
// Definition mutation and registry concurrency constraints match ExecuteActions.
//
// Check the error first. Empty cursor/nil means repository completion succeeded.
// An execution error uses Fail with a nonempty diagnostic and wraps the cause.
// Rejected terminal updates return errors: prior side effects require reconciliation,
// not blind replay. No refreshed-version retry or fallback terminal write occurs.
//
// Nonempty cursor/nil hands off an UNEXECUTED condition or WAIT, leaving the claim
// RUNNING with the same owner/version. Retain the cursor and definition locally;
// this adapter neither evaluates/prepares the boundary nor resumes/finalizes it.
// Existing condition/wait helpers remain separate; their nil errors alone do not
// finalize M2 state. Calling RunActions again cannot resume a RUNNING execution.
// No durable waiting, scheduler, lease, automatic retry or panic recovery is added.
// A lost handoff or callback panic can leave RUNNING state requiring reconciliation.
func (r *HandlerRegistry) RunActions(repo WorkflowLifecycleRepository, executionID, workerID string, d WorkflowDefinition) (string, error) {
	if !repo.TryClaim(executionID, workerID) {
		return d.StartStep, fmt.Errorf("workflow execution %q: claim rejected for worker %q", executionID, workerID)
	}
	claimed, exists := repo.Get(executionID)
	if !exists || claimed == nil || claimed.WorkflowID != executionID || claimed.Status != "RUNNING" || claimed.ClaimOwner != workerID {
		return d.StartStep, fmt.Errorf("workflow execution %q: claimed snapshot unavailable or invalid; reconcile repository state", executionID)
	}
	// Capture before callbacks: never refresh a stale version to force a write.
	version := claimed.Version
	cursor, err := r.ExecuteActions(d, d.StartStep)
	if err != nil {
		executionErr := fmt.Errorf("workflow execution %q: %w", executionID, err)
		if !repo.Fail(executionID, workerID, version, executionErr.Error()) {
			return cursor, fmt.Errorf("%w; failure transition rejected for owner %q at version %d; reconcile state and side effects", executionErr, workerID, version)
		}
		return cursor, executionErr
	}
	if cursor != "" {
		return cursor, nil
	}
	if !repo.Complete(executionID, workerID, version) {
		return cursor, fmt.Errorf("workflow execution %q: completion transition rejected for owner %q at version %d; actions succeeded, reconcile state and side effects before replay", executionID, workerID, version)
	}
	return "", nil
}
