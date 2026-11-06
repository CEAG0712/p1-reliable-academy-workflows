package enrollment

import "testing"

func TestInMemoryWorkflowRepositoryStoresAndRetrievesExecution(t *testing.T) {
    repository := NewInMemoryWorkflowRepository()

    execution := NewWorkflowExecution("workflow-001")

    repository.Save(execution)

    stored, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to exist")
    }

    if stored.WorkflowID != "workflow-001" {
        t.Fatalf(
            "expected workflow-001, got %s",
            stored.WorkflowID,
        )
    }

    if stored.Status != "READY" {
        t.Fatalf(
            "expected READY, got %s",
            stored.Status,
        )
    }
}


func TestInMemoryWorkflowRepositoryGetReturnsIndependentSnapshot(t *testing.T) {
    repository := NewInMemoryWorkflowRepository()

    original := NewWorkflowExecution("workflow-001")

    repository.Save(original)

    retrieved, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to exist")
    }

    retrieved.Status = "COMPLETED"

    storedAgain, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to still exist")
    }

    if storedAgain.Status != "READY" {
        t.Fatalf(
            "expected stored workflow to remain READY, got %s",
            storedAgain.Status,
        )
    }
}

