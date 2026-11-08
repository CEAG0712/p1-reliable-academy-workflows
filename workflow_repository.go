package enrollment

import "sync"

type InMemoryWorkflowRepository struct {
	mu         sync.RWMutex
	executions map[string]*WorkflowExecution
}

func NewInMemoryWorkflowRepository() *InMemoryWorkflowRepository {
	return &InMemoryWorkflowRepository{
		executions: make(map[string]*WorkflowExecution),
	}
}

// Creates a new execution containing copies of its data fields.
// It does not copy the source execution's mutex.
func cloneWorkflowExecution(source *WorkflowExecution) *WorkflowExecution {
	source.mu.Lock()
	defer source.mu.Unlock()

	return &WorkflowExecution{
		WorkflowID: source.WorkflowID,
		Status:     source.Status,
		ClaimOwner: source.ClaimOwner,
		Version:    source.Version,
		LastError:  source.LastError,
		Attempts:   source.Attempts,
		History: append(
			[]WorkflowTransition(nil),
			source.History...,
		),
	}
}

func (r *InMemoryWorkflowRepository) Save(execution *WorkflowExecution) {
	snapshot := cloneWorkflowExecution(execution)

	r.mu.Lock()
	defer r.mu.Unlock()

	// Do not overwrite an existing workflow execution.
	if _, exists := r.executions[snapshot.WorkflowID]; exists {
		return
	}

	r.executions[snapshot.WorkflowID] = snapshot
}

func (r *InMemoryWorkflowRepository) Get(workflowID string) (*WorkflowExecution, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, exists := r.executions[workflowID]

	if !exists {
		return nil, false
	}

	return cloneWorkflowExecution(stored), true
}

func (r *InMemoryWorkflowRepository) TryClaim(
	workflowID string,
	workerID string,
) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.executions[workflowID]

	if !exists {
		return false
	}

	if stored.Status != "READY" || workerID == "" {
		return false
	}

	previousStatus := stored.Status

	stored.Status = "RUNNING"
	stored.ClaimOwner = workerID
	stored.Version++
	stored.Attempts++

	stored.History = append(
		stored.History,
		WorkflowTransition{
			From:     previousStatus,
			To:       stored.Status,
			WorkerID: workerID,
			Version:  stored.Version,
		},
	)

	return true
}

func (r *InMemoryWorkflowRepository) SaveIfVersion(
	execution *WorkflowExecution,
	expectedVersion uint64,
) bool {
	if execution == nil {
		return false
	}

	snapshot := cloneWorkflowExecution(execution)

	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.executions[snapshot.WorkflowID]
	if !exists {
		return false
	}

	if stored.Version != expectedVersion ||
		snapshot.Version != expectedVersion {
		return false
	}

	if snapshot.ClaimOwner != stored.ClaimOwner {
		return false
	}

	if stored.Status != "RUNNING" ||
		stored.ClaimOwner == "" {
		return false
	}


    // Failure transitions must use Fail(), which requires
    // a nonempty reason and validates worker ownership.
    if snapshot.Status != "COMPLETED" {
        return false
    }


	// Preserve repository-owned metadata, not caller-supplied history.
	snapshot.Attempts = stored.Attempts
	snapshot.History = append(
		[]WorkflowTransition(nil),
		stored.History...,
	)

	previousStatus := stored.Status

	snapshot.Version = stored.Version + 1
	snapshot.LastError = ""

	snapshot.History = append(
		snapshot.History,
		WorkflowTransition{
			From:     previousStatus,
			To:       snapshot.Status,
			WorkerID: stored.ClaimOwner,
			Version:  snapshot.Version,
		},
	)

	r.executions[snapshot.WorkflowID] = snapshot

	return true
}

func (r *InMemoryWorkflowRepository) Complete(
	workflowID string,
	workerID string,
	expectedVersion uint64,
) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.executions[workflowID]

	if !exists {
		return false
	}

	if workerID == "" ||
		stored.Status != "RUNNING" ||
		stored.ClaimOwner != workerID ||
		stored.Version != expectedVersion {
		return false
	}

	previousStatus := stored.Status

	stored.Status = "COMPLETED"
	stored.LastError = ""
	stored.Version++

	stored.History = append(
		stored.History,
		WorkflowTransition{
			From:     previousStatus,
			To:       stored.Status,
			WorkerID: workerID,
			Version:  stored.Version,
		},
	)

	return true
}

func (r *InMemoryWorkflowRepository) Fail(
	workflowID string,
	workerID string,
	expectedVersion uint64,
	reason string,
) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.executions[workflowID]

	if !exists {
		return false
	}

	if workerID == "" ||
		reason == "" ||
		stored.Status != "RUNNING" ||
		stored.ClaimOwner != workerID ||
		stored.Version != expectedVersion {
		return false
	}

	previousStatus := stored.Status

	stored.Status = "FAILED"
	stored.LastError = reason
	stored.Version++

	stored.History = append(
		stored.History,
		WorkflowTransition{
			From:     previousStatus,
			To:       stored.Status,
			WorkerID: workerID,
			Version:  stored.Version,
			Reason:   reason,
		},
	)

	return true
}
