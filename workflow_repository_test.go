package enrollment

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

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


func TestInMemoryWorkflowRepositoryConditionalClaim(t *testing.T) {
    repository := NewInMemoryWorkflowRepository()

    execution := NewWorkflowExecution("workflow-001")
    repository.Save(execution)

    firstClaim := repository.TryClaim("workflow-001", "worker-A")

    if !firstClaim {
        t.Fatal("expected worker-A to claim READY workflow")
    }

    secondClaim := repository.TryClaim("workflow-001", "worker-B")

    if secondClaim {
        t.Fatal("expected worker-B claim to be rejected")
    }

    stored, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to exist")
    }

    if stored.Status != "RUNNING" {
        t.Fatalf(
            "expected RUNNING, got %s",
            stored.Status,
        )
    }

    if stored.ClaimOwner != "worker-A" {
        t.Fatalf(
            "expected worker-A as claim owner, got %s",
            stored.ClaimOwner,
        )
    }
}


func TestInMemoryWorkflowRepositoryConcurrentCompetingClaims(t *testing.T) {
    repository := NewInMemoryWorkflowRepository()

    execution := NewWorkflowExecution("workflow-001")
    repository.Save(execution)

    const workers = 10

    var wg sync.WaitGroup
    var successfulClaims atomic.Int32

    winners := make(chan string, workers)
    start := make(chan struct{})

    wg.Add(workers)

    for i := 0; i < workers; i++ {
        workerID := fmt.Sprintf("worker-%d", i)

        go func() {
            defer wg.Done()

            <-start

            if repository.TryClaim("workflow-001", workerID) {
                successfulClaims.Add(1)
                winners <- workerID
            }
        }()
    }

    close(start)

    wg.Wait()
    close(winners)

    if actual := successfulClaims.Load(); actual != 1 {
        t.Fatalf(
            "expected exactly 1 successful repository claim, got %d",
            actual,
        )
    }

    winningWorker := <-winners

    stored, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to exist")
    }

    if stored.Status != "RUNNING" {
        t.Fatalf(
            "expected RUNNING, got %s",
            stored.Status,
        )
    }

    if stored.ClaimOwner != winningWorker {
        t.Fatalf(
            "expected owner %s, got %s",
            winningWorker,
            stored.ClaimOwner,
        )
    }
}


func TestInMemoryWorkflowRepositoryStaleSaveCannotOverwriteClaim(t *testing.T) {
    repository := NewInMemoryWorkflowRepository()

    execution := NewWorkflowExecution("workflow-001")
    repository.Save(execution)

    // Worker A retrieves a snapshot while the workflow is READY.
    staleSnapshot, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to exist")
    }

    // Worker B claims the stored workflow.
    claimed := repository.TryClaim("workflow-001", "worker-B")

    if !claimed {
        t.Fatal("expected worker-B to claim READY workflow")
    }

    // Worker A attempts to save its outdated READY snapshot.
    repository.Save(staleSnapshot)

    // Retrieve the authoritative stored state.
    stored, exists := repository.Get("workflow-001")

    if !exists {
        t.Fatal("expected workflow execution to still exist")
    }

    if stored.Status != "RUNNING" {
        t.Fatalf(
            "expected RUNNING after stale save, got %s",
            stored.Status,
        )
    }

    if stored.ClaimOwner != "worker-B" {
        t.Fatalf(
            "expected worker-B to remain claim owner, got %s",
            stored.ClaimOwner,
        )
    }
}


