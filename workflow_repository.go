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

	stored.Status = "RUNNING"
	stored.ClaimOwner = workerID
	stored.Version++

	return true
}

func (r *InMemoryWorkflowRepository) SaveIfVersion(
	execution *WorkflowExecution,
	expectedVersion uint64,
) bool {
	if execution == nil {
		return false
	}

	// Copy the caller's data into an independent snapshot.
	snapshot := cloneWorkflowExecution(execution)

	r.mu.Lock()
	defer r.mu.Unlock()

	// The workflow must already exist.
	stored, exists := r.executions[snapshot.WorkflowID]

	if !exists {
		return false
	}

	// Both the expected and submitted versions must
	// match the authoritative stored version.
	if stored.Version != expectedVersion {
		return false
	}

	if snapshot.Version != expectedVersion {
		return false
	}

	// Do not allow this update to change claim ownership.
	if snapshot.ClaimOwner != stored.ClaimOwner {
		return false
	}

	// Only previously claimed executions may be updated.
	if stored.Status != "RUNNING" || stored.ClaimOwner == "" {
		return false
	}

	// M1 permits only these terminal transitions.
	if snapshot.Status != "COMPLETED" &&
		snapshot.Status != "FAILED" {
		return false
	}

	// The update is accepted. Advance the revision.
	snapshot.Version = stored.Version + 1

	// Atomically replace the authoritative stored record.
	r.executions[snapshot.WorkflowID] = snapshot

	return true
}
