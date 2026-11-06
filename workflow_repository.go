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

    return true
}



