package enrollment

import "sync"

type WorkflowExecution struct {
	mu         sync.Mutex
    WorkflowID string
    Status     string
	ClaimOwner string
}

func NewWorkflowExecution(workflowID string) *WorkflowExecution {
	return &WorkflowExecution{
		WorkflowID: workflowID,
		Status: "READY",
	}
}


func (w *WorkflowExecution) TryStart() bool {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.Status != "READY" {
        return false
    }

    w.Status = "RUNNING"
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

    return true
}
