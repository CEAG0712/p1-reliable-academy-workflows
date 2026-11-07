package enrollment

import "sync"

type WorkflowExecution struct {
	mu         sync.Mutex
	WorkflowID string
	Status     string
	ClaimOwner string
	Version    uint64
}

func NewWorkflowExecution(workflowID string) *WorkflowExecution {
	return &WorkflowExecution{
		WorkflowID: workflowID,
		Status:     "READY",
		ClaimOwner: "",
		Version:    1,
	}
}

func (w *WorkflowExecution) TryStart() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Status != "READY" {
		return false
	}

	w.Status = "RUNNING"
	w.Version++

	return true
}

func (w *WorkflowExecution) TryClaim(workerID string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Status != "READY" || workerID == "" {
		return false
	}

	w.Status = "RUNNING"
	w.ClaimOwner = workerID
	w.Version++

	return true
}
