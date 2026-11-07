package enrollment

import "sync"

type WorkflowTransition struct {
	From     string
	To       string
	WorkerID string
	Version  uint64
	Reason   string
}

type WorkflowExecution struct {
	mu         sync.Mutex
	WorkflowID string
	Status     string
	ClaimOwner string
	Version    uint64
	LastError  string
	Attempts   int
	History    []WorkflowTransition
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

	previousStatus := w.Status

	w.Status = "RUNNING"
	w.Version++
	w.Attempts++

	w.History = append(
		w.History,
		WorkflowTransition{
			From:    previousStatus,
			To:      w.Status,
			Version: w.Version,
		},
	)

	return true
}

func (w *WorkflowExecution) TryClaim(workerID string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Status != "READY" || workerID == "" {
		return false
	}

	previousStatus := w.Status

	w.Status = "RUNNING"
	w.ClaimOwner = workerID
	w.Version++
	w.Attempts++

	w.History = append(
		w.History,
		WorkflowTransition{
			From:     previousStatus,
			To:       w.Status,
			WorkerID: workerID,
			Version:  w.Version,
		},
	)

	return true
}
