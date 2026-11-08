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

func TestInMemoryWorkflowRepositoryClaimIncrementsVersion(t *testing.T) {
	repository := NewInMemoryWorkflowRepository()

	execution := NewWorkflowExecution("workflow-001")
	repository.Save(execution)

	claimed := repository.TryClaim("workflow-001", "worker-A")

	if !claimed {
		t.Fatal("expected worker-A to claim READY workflow")
	}

	stored, exists := repository.Get("workflow-001")

	if !exists {
		t.Fatal("expected workflow execution to exist")
	}

	if stored.Version != 2 {
		t.Fatalf(
			"expected stored version 2 after claim, got %d",
			stored.Version,
		)
	}

	rejected := repository.TryClaim("workflow-001", "worker-B")

	if rejected {
		t.Fatal("expected second claim to be rejected")
	}

	storedAgain, exists := repository.Get("workflow-001")

	if !exists {
		t.Fatal("expected workflow execution to still exist")
	}

	if storedAgain.Version != 2 {
		t.Fatalf(
			"expected rejected claim to preserve version 2, got %d",
			storedAgain.Version,
		)
	}
}

func TestInMemoryWorkflowRepositoryRejectsStaleVersionedUpdate(t *testing.T) {
	repository := NewInMemoryWorkflowRepository()

	execution := NewWorkflowExecution("workflow-001")
	repository.Save(execution)

	if !repository.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected worker-A to claim workflow")
	}

	// Both workers read version 2.
	firstSnapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected first snapshot")
	}

	secondSnapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected second snapshot")
	}

	// First update succeeds: version 2 -> version 3.
	firstSnapshot.Status = "COMPLETED"

	firstUpdated := repository.SaveIfVersion(
		firstSnapshot,
		firstSnapshot.Version,
	)

	if !firstUpdated {
		t.Fatal("expected first versioned update to succeed")
	}

	// Second update uses stale version 2.
	secondSnapshot.Status = "FAILED"

	secondUpdated := repository.SaveIfVersion(
		secondSnapshot,
		secondSnapshot.Version,
	)

	if secondUpdated {
		t.Fatal("expected stale versioned update to be rejected")
	}

	stored, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected stored execution")
	}

	if stored.Status != "COMPLETED" {
		t.Fatalf(
			"expected COMPLETED, got %s",
			stored.Status,
		)
	}

	if stored.Version != 3 {
		t.Fatalf(
			"expected version 3, got %d",
			stored.Version,
		)
	}

	if stored.ClaimOwner != "worker-A" {
		t.Fatalf(
			"expected worker-A to remain owner, got %s",
			stored.ClaimOwner,
		)
	}
}

func TestInMemoryWorkflowRepositoryVersionedUpdateCannotChangeOwner(t *testing.T) {
	repository := NewInMemoryWorkflowRepository()

	execution := NewWorkflowExecution("workflow-001")
	repository.Save(execution)

	if !repository.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected worker-A to claim workflow")
	}

	snapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected stored execution")
	}

	// Attempt to replace the original claim owner.
	snapshot.ClaimOwner = "worker-B"
	snapshot.Status = "COMPLETED"

	updated := repository.SaveIfVersion(
		snapshot,
		snapshot.Version,
	)

	if updated {
		t.Fatal("expected ownership change to be rejected")
	}

	stored, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected stored execution")
	}

	if stored.ClaimOwner != "worker-A" {
		t.Fatalf(
			"expected worker-A to remain owner, got %s",
			stored.ClaimOwner,
		)
	}

	if stored.Status != "RUNNING" {
		t.Fatalf(
			"expected RUNNING, got %s",
			stored.Status,
		)
	}

	if stored.Version != 2 {
		t.Fatalf(
			"expected version 2, got %d",
			stored.Version,
		)
	}
}

func TestInMemoryWorkflowRepositoryVersionedUpdateCannotReopenCompleted(t *testing.T) {
	repository := NewInMemoryWorkflowRepository()

	execution := NewWorkflowExecution("workflow-001")
	repository.Save(execution)

	if !repository.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected worker-A to claim workflow")
	}

	snapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected stored execution")
	}

	snapshot.Status = "COMPLETED"

	if !repository.SaveIfVersion(snapshot, snapshot.Version) {
		t.Fatal("expected completion update to succeed")
	}

	completed, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected completed execution")
	}

	// Attempt to reopen the completed execution.
	completed.Status = "READY"

	updated := repository.SaveIfVersion(
		completed,
		completed.Version,
	)

	if updated {
		t.Fatal("expected completed workflow to reject update")
	}

	stored, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected stored execution")
	}

	if stored.Status != "COMPLETED" {
		t.Fatalf(
			"expected COMPLETED, got %s",
			stored.Status,
		)
	}

	if stored.Version != 3 {
		t.Fatalf(
			"expected version 3, got %d",
			stored.Version,
		)
	}
}

func TestInMemoryWorkflowRepositoryConcurrentVersionedUpdates(t *testing.T) {
	repository := NewInMemoryWorkflowRepository()

	execution := NewWorkflowExecution("workflow-001")
	repository.Save(execution)

	if !repository.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected worker-A to claim workflow")
	}

	// Both snapshots represent version 2.
	firstSnapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected first snapshot")
	}

	secondSnapshot, exists := repository.Get("workflow-001")
	if !exists {
		t.Fatal("expected second snapshot")
	}

	firstSnapshot.Status = "COMPLETED"
	secondSnapshot.Status = "FAILED"

	var wg sync.WaitGroup
	var successfulUpdates atomic.Int32

	start := make(chan struct{})

	wg.Add(2)

	go func() {
		defer wg.Done()

		<-start

		if repository.SaveIfVersion(
			firstSnapshot,
			firstSnapshot.Version,
		) {
			successfulUpdates.Add(1)
		}
	}()

	go func() {
		defer wg.Done()

		<-start

		if repository.SaveIfVersion(
			secondSnapshot,
			secondSnapshot.Version,
		) {
			successfulUpdates.Add(1)
		}
	}()

	// Release both workers.
	close(start)

	wg.Wait()

	if actual := successfulUpdates.Load(); actual != 1 {
		t.Fatalf(
			"expected exactly 1 successful versioned update, got %d",
			actual,
		)
	}

	stored, exists := repository.Get("workflow-001")

	if !exists {
		t.Fatal("expected stored execution")
	}

	if stored.Version != 3 {
		t.Fatalf(
			"expected version 3, got %d",
			stored.Version,
		)
	}

	if stored.Status != "COMPLETED" &&
		stored.Status != "FAILED" {
		t.Fatalf(
			"expected COMPLETED or FAILED, got %s",
			stored.Status,
		)
	}

	if stored.ClaimOwner != "worker-A" {
		t.Fatalf(
			"expected worker-A to remain owner, got %s",
			stored.ClaimOwner,
		)
	}
}

func TestRepositoryOwnerCanCompleteWorkflow(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	if !repo.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected claim to succeed")
	}

	completed := repo.Complete("workflow-001", "worker-A", 2)

	if !completed {
		t.Fatal("expected owner to complete workflow")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED, got %s", stored.Status)
	}

	if stored.Version != 3 {
		t.Fatalf("expected version 3, got %d", stored.Version)
	}
}

func TestRepositoryRejectsCompletionByDifferentWorker(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	if repo.Complete("workflow-001", "worker-B", 2) {
		t.Fatal("expected non-owner completion to be rejected")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "RUNNING" || stored.Version != 2 {
		t.Fatal("rejected completion must preserve state and version")
	}
}

func TestRepositoryRejectsStaleCompletion(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	if repo.Complete("workflow-001", "worker-A", 1) {
		t.Fatal("expected stale version to be rejected")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "RUNNING" || stored.Version != 2 {
		t.Fatal("stale completion changed stored execution")
	}
}

func TestRepositoryOwnerCanFailWorkflow(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	failed := repo.Fail(
		"workflow-001",
		"worker-A",
		2,
		"provider timeout",
	)

	if !failed {
		t.Fatal("expected owner to fail workflow")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "FAILED" {
		t.Fatalf("expected FAILED, got %s", stored.Status)
	}

	if stored.LastError != "provider timeout" {
		t.Fatalf("unexpected error: %s", stored.LastError)
	}

	if stored.Version != 3 {
		t.Fatalf("expected version 3, got %d", stored.Version)
	}
}

func TestRepositoryRejectsFailureWithoutReason(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	if repo.Fail("workflow-001", "worker-A", 2, "") {
		t.Fatal("expected empty failure reason to be rejected")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "RUNNING" || stored.Version != 2 {
		t.Fatal("invalid failure changed execution")
	}
}

func TestRepositoryRejectsFailureByDifferentWorker(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	if repo.Fail("workflow-001", "worker-B", 2, "error") {
		t.Fatal("expected non-owner failure to be rejected")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "RUNNING" || stored.Version != 2 {
		t.Fatal("rejected failure changed execution")
	}
}

func TestRepositoryTracksSuccessfulClaimAndCompletion(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	if !repo.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected claim to succeed")
	}

	if repo.TryClaim("workflow-001", "worker-B") {
		t.Fatal("expected second claim to fail")
	}

	if !repo.Complete("workflow-001", "worker-A", 2) {
		t.Fatal("expected completion to succeed")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", stored.Attempts)
	}

	if len(stored.History) != 2 {
		t.Fatalf(
			"expected 2 history entries, got %d",
			len(stored.History),
		)
	}

	claim := stored.History[0]

	if claim.From != "READY" ||
		claim.To != "RUNNING" ||
		claim.WorkerID != "worker-A" ||
		claim.Version != 2 {
		t.Fatal("incorrect claim history entry")
	}

	completion := stored.History[1]

	if completion.From != "RUNNING" ||
		completion.To != "COMPLETED" ||
		completion.WorkerID != "worker-A" ||
		completion.Version != 3 {
		t.Fatal("incorrect completion history entry")
	}
}

func TestRepositoryTracksFailureHistory(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	repo.TryClaim("workflow-001", "worker-A")

	if !repo.Fail(
		"workflow-001",
		"worker-A",
		2,
		"provider timeout",
	) {
		t.Fatal("expected failure transition to succeed")
	}

	stored, _ := repo.Get("workflow-001")

	if len(stored.History) != 2 {
		t.Fatalf(
			"expected 2 transitions, got %d",
			len(stored.History),
		)
	}

	failure := stored.History[1]

	if failure.To != "FAILED" ||
		failure.Reason != "provider timeout" ||
		failure.Version != 3 {
		t.Fatal("incorrect failure history entry")
	}
}

func TestRepositoryHistorySnapshotIsIndependent(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	snapshot, exists := repo.Get("workflow-001")

	if !exists {
		t.Fatal("expected stored workflow")
	}

	if len(snapshot.History) != 1 {
		t.Fatal("expected claim history entry")
	}

	// Attempt to tamper with the returned snapshot.
	snapshot.History[0].To = "COMPLETED"

	storedAgain, _ := repo.Get("workflow-001")

	if storedAgain.History[0].To != "RUNNING" {
		t.Fatal("snapshot mutation changed stored history")
	}
}

func TestRepositoryVersionedUpdateRecordsHistory(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))
	repo.TryClaim("workflow-001", "worker-A")

	snapshot, _ := repo.Get("workflow-001")
	snapshot.Status = "COMPLETED"

	if !repo.SaveIfVersion(snapshot, snapshot.Version) {
		t.Fatal("expected versioned update to succeed")
	}

	stored, _ := repo.Get("workflow-001")

	if len(stored.History) != 2 {
		t.Fatalf(
			"expected 2 history entries, got %d",
			len(stored.History),
		)
	}

	if stored.History[1].To != "COMPLETED" ||
		stored.History[1].Version != 3 {
		t.Fatal("versioned update history was not recorded")
	}
}

func TestRepositoryCannotCompleteReadyWorkflow(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	if repo.Complete("workflow-001", "worker-A", 1) {
		t.Fatal("READY workflow must not complete")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "READY" ||
		stored.Version != 1 ||
		len(stored.History) != 0 {
		t.Fatal("invalid completion changed execution")
	}
}

func TestRepositoryCannotCompleteTwice(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	repo.TryClaim("workflow-001", "worker-A")

	if !repo.Complete("workflow-001", "worker-A", 2) {
		t.Fatal("expected first completion to succeed")
	}

	if repo.Complete("workflow-001", "worker-A", 3) {
		t.Fatal("expected second completion to be rejected")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "COMPLETED" ||
		stored.Version != 3 ||
		len(stored.History) != 2 {
		t.Fatal("duplicate completion changed execution")
	}
}

func TestRepositoryCannotFailCompletedWorkflow(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("workflow-001"))

	repo.TryClaim("workflow-001", "worker-A")
	repo.Complete("workflow-001", "worker-A", 2)

	if repo.Fail(
		"workflow-001",
		"worker-A",
		3,
		"late error",
	) {
		t.Fatal("completed workflow must not become FAILED")
	}

	stored, _ := repo.Get("workflow-001")

	if stored.Status != "COMPLETED" ||
		stored.Version != 3 ||
		len(stored.History) != 2 {
		t.Fatal("terminal-state mutation detected")
	}
}

func TestRepositoryCannotCompleteMissingWorkflow(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()

	if repo.Complete("missing", "worker-A", 1) {
		t.Fatal("missing workflow must not complete")
	}

	if repo.Fail("missing", "worker-A", 1, "error") {
		t.Fatal("missing workflow must not fail")
	}
}

func TestRepositorySaveIfVersionCannotBypassFail(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()

	repo.Save(NewWorkflowExecution("workflow-001"))

	if !repo.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected worker-A to claim workflow")
	}

	snapshot, exists := repo.Get("workflow-001")

	if !exists {
		t.Fatal("expected workflow to exist")
	}

	// Attempt to fail the workflow without using Fail().
	snapshot.Status = "FAILED"

	updated := repo.SaveIfVersion(
		snapshot,
		snapshot.Version,
	)

	if updated {
		t.Fatal("expected SaveIfVersion to reject FAILED transition")
	}

	stored, exists := repo.Get("workflow-001")

	if !exists {
		t.Fatal("expected workflow to remain stored")
	}

	if stored.Status != "RUNNING" {
		t.Fatalf("expected RUNNING, got %s", stored.Status)
	}

	if stored.Version != 2 {
		t.Fatalf("expected version 2, got %d", stored.Version)
	}

	if stored.LastError != "" {
		t.Fatalf("expected no failure reason, got %s", stored.LastError)
	}

	if stored.ClaimOwner != "worker-A" {
		t.Fatalf("expected worker-A, got %s", stored.ClaimOwner)
	}

	if len(stored.History) != 1 {
		t.Fatalf(
			"expected 1 history entry, got %d",
			len(stored.History),
		)
	}
}

func TestRepositoryFailPreservesReasonInStateAndHistory(t *testing.T) {
	repo := NewInMemoryWorkflowRepository()

	repo.Save(NewWorkflowExecution("workflow-001"))

	if !repo.TryClaim("workflow-001", "worker-A") {
		t.Fatal("expected claim to succeed")
	}

	reason := "provider timeout"

	if !repo.Fail("workflow-001", "worker-A", 2, reason) {
		t.Fatal("expected Fail to succeed")
	}

	stored, exists := repo.Get("workflow-001")

	if !exists {
		t.Fatal("expected stored workflow")
	}

	if stored.Status != "FAILED" {
		t.Fatalf("expected FAILED, got %s", stored.Status)
	}

	if stored.LastError != reason {
		t.Fatalf(
			"expected LastError %q, got %q",
			reason,
			stored.LastError,
		)
	}

	if stored.Version != 3 {
		t.Fatalf("expected version 3, got %d", stored.Version)
	}

	if len(stored.History) != 2 {
		t.Fatalf(
			"expected 2 history entries, got %d",
			len(stored.History),
		)
	}

	failure := stored.History[1]

	if failure.To != "FAILED" ||
		failure.Reason != reason ||
		failure.WorkerID != "worker-A" ||
		failure.Version != 3 {
		t.Fatal("failure history does not match stored state")
	}
}
