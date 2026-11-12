package enrollment

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type runnerResult struct {
	worker, cursor string
	err            error
}

// Channels provide ordering; this deadline only bounds a broken test or runner.
func awaitRunnerValue[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for runner synchronization")
		var zero T
		return zero
	}
}

func runnerSnapshot(t *testing.T, repo *InMemoryWorkflowRepository) *WorkflowExecution {
	t.Helper()
	s, exists := repo.Get("execution-42")
	if !exists {
		t.Fatal("missing execution")
	}
	return s
}

func assertRunnerWinner(t *testing.T, repo *InMemoryWorkflowRepository, worker, status, reason string) {
	t.Helper()
	s := runnerSnapshot(t, repo)
	wantHistory := []WorkflowTransition{
		{From: "READY", To: "RUNNING", WorkerID: worker, Version: 2},
		{From: "RUNNING", To: status, WorkerID: worker, Version: 3, Reason: reason},
	}
	if s.Status != status || s.ClaimOwner != worker || s.Version != 3 || s.Attempts != 1 || s.LastError != reason || !reflect.DeepEqual(s.History, wantHistory) {
		t.Fatalf("unexpected lifecycle: %+v, history=%+v", s, s.History)
	}
}

func assertRunnerClaimRejection(t *testing.T, result runnerResult, start string) {
	t.Helper()
	if result.err == nil || !strings.Contains(result.err.Error(), "claim rejected") || result.cursor != start {
		t.Fatalf("expected rejection before execution, got %+v", result)
	}
}

func TestRunnerConcurrencySuccessfulExecutions(t *testing.T) {
	testRunnerCompetingExecutions(t, false)
}

func TestRunnerConcurrencyFailedExecutions(t *testing.T) {
	testRunnerCompetingExecutions(t, true)
}

func testRunnerCompetingExecutions(t *testing.T, fail bool) {
	t.Helper()
	repo, r := runnerFixture()
	d := actionChainDefinition() // Shared read-only after all registrations finish.
	cause := errors.New("provider rejected transform")
	var mu sync.Mutex
	var calls []string
	for _, name := range []string{"load", "transform", "store", "must-not-run"} {
		if err := r.RegisterAction(name, func() error {
			mu.Lock()
			calls = append(calls, name)
			mu.Unlock()
			if fail && name == "transform" {
				return cause
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	const workers = 10
	start := make(chan struct{})
	ready := make(chan struct{}, workers)
	results := make(chan runnerResult, workers)
	for i := 0; i < workers; i++ {
		worker := fmt.Sprintf("worker-%d", i)
		go func() {
			ready <- struct{}{}
			<-start
			cursor, err := r.RunActions(repo, "execution-42", worker, d)
			results <- runnerResult{worker, cursor, err}
		}()
	}
	// Cleanup also releases workers if a readiness assertion fails.
	var release sync.Once
	t.Cleanup(func() { release.Do(func() { close(start) }) })
	for i := 0; i < workers; i++ {
		awaitRunnerValue(t, ready)
	}
	release.Do(func() { close(start) })
	winners, rejected := 0, 0
	var winner runnerResult
	for i := 0; i < workers; i++ {
		result := awaitRunnerValue(t, results)
		if result.err == nil || errors.Is(result.err, cause) {
			winners++
			winner = result
		} else {
			assertRunnerClaimRejection(t, result, d.StartStep)
			rejected++
		}
	}
	if winners != 1 || rejected != workers-1 {
		t.Fatalf("action-path outcomes=%d, rejections=%d", winners, rejected)
	}
	status, reason, cursor := "COMPLETED", "", ""
	wantCalls := []string{"load", "transform", "store"}
	if fail {
		status, cursor = "FAILED", "a-middle"
		wantCalls = []string{"load", "transform"}
		if !errors.Is(winner.err, cause) {
			t.Fatalf("lost action failure: %v", winner.err)
		}
		reason = winner.err.Error()
		if reason == "" || !strings.Contains(reason, cause.Error()) || !strings.Contains(reason, cursor) {
			t.Fatalf("missing failure diagnostic: %q", reason)
		}
	} else if winner.err != nil {
		t.Fatal(winner.err)
	}
	if winner.cursor != cursor {
		t.Fatalf("cursor=%q, want %q", winner.cursor, cursor)
	}
	mu.Lock()
	gotCalls := append([]string(nil), calls...)
	mu.Unlock()
	if !reflect.DeepEqual(gotCalls, wantCalls) {
		t.Fatalf("callbacks=%v, want exactly one ordered path %v", gotCalls, wantCalls)
	}
	assertRunnerWinner(t, repo, winner.worker, status, reason)
}

func TestRunnerConcurrencyDuplicateAfterCompletion(t *testing.T) {
	repo, r := runnerFixture()
	d := actionChainDefinition()
	var calls []string
	for _, name := range []string{"load", "transform", "store", "must-not-run"} {
		if err := r.RegisterAction(name, func() error { calls = append(calls, name); return nil }); err != nil {
			t.Fatal(err)
		}
	}
	if cursor, err := r.RunActions(repo, "execution-42", "worker-A", d); err != nil || cursor != "" {
		t.Fatalf("first invocation=(%q, %v)", cursor, err)
	}
	assertRunnerWinner(t, repo, "worker-A", "COMPLETED", "")
	before := runnerSnapshot(t, repo)
	for _, worker := range []string{"worker-A", "worker-B"} {
		cursor, err := r.RunActions(repo, "execution-42", worker, d)
		assertRunnerClaimRejection(t, runnerResult{worker, cursor, err}, d.StartStep)
		if after := runnerSnapshot(t, repo); !reflect.DeepEqual(before, after) {
			t.Fatalf("duplicate changed lifecycle: before=%+v after=%+v", before, after)
		}
	}
	if !reflect.DeepEqual(calls, []string{"load", "transform", "store"}) {
		t.Fatalf("duplicate callbacks: %v", calls)
	}
}

func TestRunnerConcurrencyInvalidOwnership(t *testing.T) {
	// Exercise the real lifecycle boundary, rather than inventing a runner API.
	for _, operation := range []string{"complete", "fail"} {
		t.Run(operation, func(t *testing.T) {
			repo, _ := runnerFixture()
			if !repo.TryClaim("execution-42", "worker-A") {
				t.Fatal("fixture claim failed")
			}
			before := runnerSnapshot(t, repo)
			var accepted bool
			if operation == "complete" {
				accepted = repo.Complete("execution-42", "worker-B", before.Version)
			} else {
				accepted = repo.Fail("execution-42", "worker-B", before.Version, "intruder failure")
			}
			if accepted || !reflect.DeepEqual(before, runnerSnapshot(t, repo)) {
				t.Fatal("non-owner terminal operation accepted or changed lifecycle")
			}
			if !repo.Complete("execution-42", "worker-A", before.Version) {
				t.Fatal("rejection prevented owner's valid completion")
			}
			assertRunnerWinner(t, repo, "worker-A", "COMPLETED", "")
		})
	}
}

func TestRunnerConcurrencyBlockedCallback(t *testing.T) {
	repo, r := runnerFixture()
	d := actionChainDefinition()
	entered, release := make(chan struct{}, 2), make(chan struct{})
	var unblock sync.Once
	var loadCalls, laterCalls atomic.Int32
	if err := r.RegisterAction("load", func() error {
		loadCalls.Add(1)
		entered <- struct{}{}
		select {
		case <-release:
			return nil
		case <-time.After(5 * time.Second):
			return errors.New("test callback release timed out")
		}
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"transform", "store"} {
		if err := r.RegisterAction(name, func() error { laterCalls.Add(1); return nil }); err != nil {
			t.Fatal(err)
		}
	}
	winner := make(chan runnerResult, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		cursor, err := r.RunActions(repo, "execution-42", "worker-A", d)
		winner <- runnerResult{"worker-A", cursor, err}
	}()
	t.Cleanup(func() {
		unblock.Do(func() { close(release) })
		awaitRunnerValue(t, done)
	})
	awaitRunnerValue(t, entered)
	assertRunnerState(t, repo, "RUNNING", "")
	before := runnerSnapshot(t, repo)
	competitor := make(chan runnerResult, 1)
	go func() {
		cursor, err := r.RunActions(repo, "execution-42", "worker-B", d)
		competitor <- runnerResult{"worker-B", cursor, err}
	}()
	assertRunnerClaimRejection(t, awaitRunnerValue(t, competitor), d.StartStep)
	if loadCalls.Load() != 1 || laterCalls.Load() != 0 || !reflect.DeepEqual(before, runnerSnapshot(t, repo)) {
		t.Fatal("contender executed callbacks or changed the blocked lifecycle")
	}
	unblock.Do(func() { close(release) })
	result := awaitRunnerValue(t, winner)
	if result.err != nil || result.cursor != "" {
		t.Fatalf("winner failed after release: %+v", result)
	}
	if loadCalls.Load() != 1 || laterCalls.Load() != 2 {
		t.Fatal("unexpected callbacks after release")
	}
	assertRunnerWinner(t, repo, "worker-A", "COMPLETED", "")
}
