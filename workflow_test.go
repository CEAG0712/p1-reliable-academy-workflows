package enrollment

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewWorkflowExecutionStartsReady(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	if execution.WorkflowID != "workflow-001" {
		t.Fatalf(
			"expected workflow-001, got %s",
			execution.WorkflowID,
		)
	}

	if execution.Status != "READY" {
		t.Fatalf(
			"expected READY, got %s",
			execution.Status,
		)
	}
}

func TestWorkflowExecutionCanStartOnlyOnce(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	firstClaim := execution.TryStart()

	if !firstClaim {
		t.Fatal("expected READY workflow to start")
	}

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected RUNNING, got %s",
			execution.Status,
		)
	}

	secondClaim := execution.TryStart()

	if secondClaim {
		t.Fatal("expected second start attempt to be rejected")
	}

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected workflow to remain RUNNING, got %s",
			execution.Status,
		)
	}
}

func TestWorkflowExecutionConcurrentClaimsAllowOnlyOneWinner(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	const workers = 10

	var wg sync.WaitGroup
	var successfulClaims atomic.Int32

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			if execution.TryStart() {
				successfulClaims.Add(1)
			}
		}()
	}

	wg.Wait()

	if actual := successfulClaims.Load(); actual != 1 {
		t.Fatalf(
			"expected exactly 1 successful claim, got %d",
			actual,
		)
	}

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected RUNNING, got %s",
			execution.Status,
		)
	}
}

func TestNewWorkflowExecutionStartsWithoutClaimOwner(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	if execution.Status != "READY" {
		t.Fatalf(
			"expected READY, got %s",
			execution.Status,
		)
	}

	if execution.ClaimOwner != "" {
		t.Fatalf(
			"expected empty claim owner, got %s",
			execution.ClaimOwner,
		)
	}
}

func TestWorkflowExecutionClaimPreservesOriginalOwner(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	firstClaim := execution.TryClaim("worker-A")

	if !firstClaim {
		t.Fatal("expected worker-A to claim READY workflow")
	}

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected RUNNING, got %s",
			execution.Status,
		)
	}

	if execution.ClaimOwner != "worker-A" {
		t.Fatalf(
			"expected worker-A as claim owner, got %s",
			execution.ClaimOwner,
		)
	}

	secondClaim := execution.TryClaim("worker-B")

	if secondClaim {
		t.Fatal("expected worker-B claim to be rejected")
	}

	if execution.ClaimOwner != "worker-A" {
		t.Fatalf(
			"expected original owner worker-A, got %s",
			execution.ClaimOwner,
		)
	}

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected workflow to remain RUNNING, got %s",
			execution.Status,
		)
	}
}

func TestWorkflowExecutionConcurrentCompetingOwners(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	const workers = 10

	var wg sync.WaitGroup
	var successfulClaims atomic.Int32

	winners := make(chan string, workers)

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		workerID := fmt.Sprintf("worker-%d", i)

		go func() {
			defer wg.Done()

			if execution.TryClaim(workerID) {
				successfulClaims.Add(1)
				winners <- workerID
			}
		}()
	}

	wg.Wait()
	close(winners)

	if actual := successfulClaims.Load(); actual != 1 {
		t.Fatalf(
			"expected exactly 1 successful claim, got %d",
			actual,
		)
	}

	winningWorker := <-winners

	if execution.Status != "RUNNING" {
		t.Fatalf(
			"expected RUNNING, got %s",
			execution.Status,
		)
	}

	if execution.ClaimOwner != winningWorker {
		t.Fatalf(
			"expected owner %s, got %s",
			winningWorker,
			execution.ClaimOwner,
		)
	}
}

func TestNewWorkflowExecutionStartsAtVersionOne(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	if execution.Version != 1 {
		t.Fatalf(
			"expected new workflow version 1, got %d",
			execution.Version,
		)
	}
}

func TestWorkflowExecutionClaimIncrementsVersion(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	claimed := execution.TryClaim("worker-A")

	if !claimed {
		t.Fatal("expected worker-A to claim READY workflow")
	}

	if execution.Version != 2 {
		t.Fatalf(
			"expected version 2 after claim, got %d",
			execution.Version,
		)
	}
}

func TestWorkflowExecutionStartIncrementsVersion(t *testing.T) {
	execution := NewWorkflowExecution("workflow-001")

	started := execution.TryStart()

	if !started {
		t.Fatal("expected READY workflow to start")
	}

	if execution.Version != 2 {
		t.Fatalf(
			"expected version 2 after start, got %d",
			execution.Version,
		)
	}
}
