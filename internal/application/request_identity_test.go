package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	workflow "example.com/p1-reliable-academy-workflows"
)

func identityFixture(t *testing.T, limit int) (*LocalRequestIdentityScope, *workflow.InMemoryWorkflowRepository, *atomic.Int32) {
	t.Helper()
	repo := workflow.NewInMemoryWorkflowRepository()
	calls := new(atomic.Int32)
	s, err := NewLocalRequestIdentityScope(repo, invocationIDFunc(func(context.Context) (string, error) {
		return fmt.Sprintf("exec-%032x", calls.Add(1)), nil
	}), "identity-worker", limit)
	if err != nil {
		t.Fatal(err)
	}
	return s, repo, calls
}
func identityCommand(t *testing.T, s *LocalRequestIdentityScope, key, definition, event string) TrustedLocalRequestCommand {
	t.Helper()
	c, err := s.Command(key, SubmissionCommand{definition, event})
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func identityInvoke(t *testing.T, s *LocalRequestIdentityScope, key, definition, event string) (InvocationOutcome, error) {
	t.Helper()
	return s.Invoker().InvokeLocalRequest(context.Background(), identityCommand(t, s, key, definition, event))
}
func identityState(s *LocalRequestIdentityScope, key string) requestReservationState {
	s.service.mu.Lock()
	defer s.service.mu.Unlock()
	return s.service.reservations[key].state
}
func TestLocalRequestIdentityConstruction(t *testing.T) {
	for _, service := range []*requestIdentityService{nil, {}} {
		if _, e := service.InvokeLocalRequest(context.Background(), TrustedLocalRequestCommand{}); e == nil {
			t.Fatal("unconfigured service accepted")
		}
	}
	var typedNil *invocationIDs
	for _, tc := range []struct {
		repo  *workflow.InMemoryWorkflowRepository
		ids   ExecutionIDSource
		owner string
		limit int
	}{
		{nil, &invocationIDs{}, "owner", 1}, {workflow.NewInMemoryWorkflowRepository(), nil, "owner", 1},
		{workflow.NewInMemoryWorkflowRepository(), typedNil, "owner", 1},
		{workflow.NewInMemoryWorkflowRepository(), &invocationIDs{}, " ", 1},
		{workflow.NewInMemoryWorkflowRepository(), &invocationIDs{}, "owner", 0},
		{workflow.NewInMemoryWorkflowRepository(), &invocationIDs{}, "owner", -1},
	} {
		if s, err := NewLocalRequestIdentityScope(tc.repo, tc.ids, tc.owner, tc.limit); s != nil || err == nil {
			t.Fatal("bad construction accepted")
		}
	}
	s, repo, calls := identityFixture(t, 2)
	if calls.Load() != 0 {
		t.Fatal("constructor called ID source")
	}
	if _, ok := repo.Get(invocationID); ok {
		t.Fatal("constructor wrote")
	}
	for _, zero := range []*LocalRequestIdentityScope{nil, {}} {
		if zero.Invoker() != nil || zero.ApprovedDefinitions() != nil {
			t.Fatal("zero scope")
		}
		if _, e := zero.Command("key", SubmissionCommand{"local.success", "event"}); e == nil {
			t.Fatal("zero command")
		}
	}
	approved := s.ApprovedDefinitions()
	delete(approved, "local.success")
	if _, e := s.Command("key", SubmissionCommand{"local.success", "event"}); e != nil {
		t.Fatal("support map escaped")
	}
}
func TestLocalRequestIdentityValidationBeforeReservation(t *testing.T) {
	s, _, calls := identityFixture(t, 1)
	good := identityCommand(t, s, "key", "local.success", "event")
	foreign, _, _ := identityFixture(t, 1)
	var decoded TrustedLocalRequestCommand
	if e := json.Unmarshal([]byte(`{"key":"key","authority":{},"selection":{"WorkflowDefinitionID":"local.success","EventID":"event"}}`), &decoded); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		c    TrustedLocalRequestCommand
		ctx  context.Context
		kind InvocationErrorKind
	}{
		{TrustedLocalRequestCommand{}, context.Background(), InvocationMissingAuthority},
		{decoded, context.Background(), InvocationMissingAuthority},
		{identityCommand(t, foreign, "key", "local.success", "event"), context.Background(), InvocationForbidden},
		{good, nil, InvocationInvalidCommand},
	} {
		_, e := s.Invoker().InvokeLocalRequest(tc.ctx, tc.c)
		invocationKind(t, e, tc.kind)
	}
	for _, key := range []string{"", "a,b", " a", "a ", "a/b", "é", strings.Repeat("a", 129)} {
		if _, e := s.Command(key, SubmissionCommand{"local.success", "event"}); e == nil {
			t.Fatal("bad key accepted", key)
		}
		c := good
		c.key = key
		_, e := s.Invoker().InvokeLocalRequest(context.Background(), c)
		invocationKind(t, e, InvocationInvalidCommand)
	}
	for _, selection := range []SubmissionCommand{{"local.success", ""}, {"local.success", "a/b"}, {"other", "event"}} {
		if _, e := s.Command("key", selection); e == nil {
			t.Fatal("bad selection accepted")
		}
		c := good
		c.selection = selection
		if _, e := s.Invoker().InvokeLocalRequest(context.Background(), c); e == nil {
			t.Fatal("direct service gate bypass")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e := s.Invoker().InvokeLocalRequest(ctx, good)
	invocationKind(t, e, InvocationUnavailable)
	if !errors.Is(e, context.Canceled) || calls.Load() != 0 || len(s.service.reservations) != 0 {
		t.Fatal("rejection reserved or operated")
	}
	if _, e := s.Invoker().InvokeLocalRequest(context.Background(), good); e != nil {
		t.Fatal("rejection consumed key", e)
	}
}
func TestLocalRequestIdentityOutcomesAndRepeats(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			s, repo, calls := identityFixture(t, 1)
			out, e := identityInvoke(t, s, "key", definition, "event")
			if definition == "local.failure" {
				invocationKind(t, e, InvocationExecutionFailed)
			} else if e != nil {
				t.Fatal(e)
			}
			if out.Creation != CreationInserted || identityState(s, "key") != requestKnown {
				t.Fatalf("bad known outcome %+v %v", out, e)
			}
			before, _ := repo.Get(out.ExecutionID)
			for i := 0; i < 3; i++ {
				dup, e := identityInvoke(t, s, "key", definition, "event")
				invocationKind(t, e, InvocationDuplicate)
				if dup.RunAttempted || dup.Snapshot != nil || dup.Creation != CreationNotAttempted {
					t.Fatal("duplicate adopted result")
				}
			}
			_, e = identityInvoke(t, s, "key", definition, "changed")
			invocationKind(t, e, InvocationDuplicate)
			after, _ := repo.Get(out.ExecutionID)
			if !reflect.DeepEqual(before, after) || calls.Load() != 1 {
				t.Fatal("repeat mutated/replayed")
			}
			if _, ok := repo.Get(fmt.Sprintf("exec-%032x", 2)); ok {
				t.Fatal("extra record")
			}
		})
	}
}
func TestLocalRequestIdentityCapacityAndScope(t *testing.T) {
	s, repo, calls := identityFixture(t, 2)
	for _, key := range []string{"Key", "key"} {
		if _, e := identityInvoke(t, s, key, "local.success", "same-event"); e != nil {
			t.Fatal(e)
		}
	}
	_, e := identityInvoke(t, s, "third", "local.success", "same-event")
	invocationKind(t, e, InvocationUnavailable)
	_, e = identityInvoke(t, s, "Key", "local.success", "same-event")
	invocationKind(t, e, InvocationDuplicate)
	_, e = identityInvoke(t, s, "Key", "local.failure", "same-event")
	invocationKind(t, e, InvocationDuplicate)
	if calls.Load() != 2 || len(s.service.reservations) != 2 {
		t.Fatal("capacity changed existing identity")
	}
	other, e := NewLocalRequestIdentityScope(repo, invocationIDFunc(func(context.Context) (string, error) { return fmt.Sprintf("exec-%032x", 3), nil }), "other", 1)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = identityInvoke(t, other, "Key", "local.success", "same-event"); e != nil {
		t.Fatal("scopes unexpectedly coordinate", e)
	}
}
func TestLocalRequestIdentityCollision(t *testing.T) {
	for _, state := range []string{"READY", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(state, func(t *testing.T) {
			s, repo, calls := identityFixture(t, 2)
			repo.Save(workflow.NewWorkflowExecution(invocationID))
			if state != "READY" {
				repo.TryClaim(invocationID, "prior")
			}
			if state == "COMPLETED" {
				repo.Complete(invocationID, "prior", 2)
			}
			if state == "FAILED" {
				repo.Fail(invocationID, "prior", 2, "preserved")
			}
			before, _ := repo.Get(invocationID)
			out, e := identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, InvocationDuplicate)
			if out.Creation != CreationDuplicate || out.RunAttempted {
				t.Fatal("collision adopted")
			}
			_, e = identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, InvocationDuplicate)
			after, _ := repo.Get(invocationID)
			if !reflect.DeepEqual(before, after) || calls.Load() != 1 {
				t.Fatal("collision replaced/replayed")
			}
		})
	}
}

// Fault collaborators are private test composition. Public scope construction
// still accepts only the concrete repository and fixed harmless resolver.
func identityPrivate(t *testing.T, plan ApprovedLocalPlan) (*LocalRequestIdentityScope, *invocationRepo, *invocationIDs) {
	t.Helper()
	core, repo, _, ids, _ := privateInvocation(t, plan)
	s, e := newLocalRequestIdentityScope(core, 4)
	if e != nil {
		t.Fatal(e)
	}
	return s, repo, ids
}

type identityCreateFault struct {
	InvocationRepository
	inserted bool
	cause    error
	insert   bool
}

type identityFailFault struct{ InvocationRepository }

func (r identityFailFault) Fail(id, owner string, version uint64, reason string) bool {
	r.InvocationRepository.Fail(id, owner, version, reason)
	return false // Stored FAILED is not proof this call observed accepted Fail.
}

func TestLocalRequestIdentityFailedSnapshotWithoutAcceptedFail(t *testing.T) {
	plan, e := (localDefinitionResolver{}).ResolveApprovedDefinition(context.Background(), "local.failure")
	if e != nil {
		t.Fatal(e)
	}
	core, repo, _, ids, _ := privateInvocation(t, plan)
	core.deps.repo = identityFailFault{repo}
	s, e := newLocalRequestIdentityScope(core, 1)
	if e != nil {
		t.Fatal(e)
	}
	out, e := identityInvoke(t, s, "key", "local.failure", "event")
	invocationKind(t, e, InvocationReconciliation)
	if out.Kind != OutcomeReconciliation || out.Snapshot == nil || out.Snapshot.Status != "FAILED" || identityState(s, "key") != requestUncertain {
		t.Fatal("FAILED snapshot mistaken for authoritative acceptance")
	}
	_, e = identityInvoke(t, s, "key", "local.failure", "event")
	invocationKind(t, e, InvocationReconciliation)
	if ids.calls != 1 {
		t.Fatal("failed uncertainty replayed")
	}
}

func (r identityCreateFault) CreateIfAbsent(w *workflow.WorkflowExecution) (bool, error) {
	if r.insert {
		r.InvocationRepository.CreateIfAbsent(w)
	}
	return r.inserted, r.cause
}
func TestLocalRequestIdentityUncertainty(t *testing.T) {
	for _, mode := range []string{"false-error", "hidden-insert-error", "true-error", "missing-after-create", "contradictory-after-create", "terminal-rejected", "terminal-lied"} {
		t.Run(mode, func(t *testing.T) {
			s, repo, ids := identityPrivate(t, localTestPlan(t))
			cause := errors.New("private creation diagnostic")
			switch mode {
			case "false-error", "hidden-insert-error", "true-error":
				s.service.invocation.deps.repo = identityCreateFault{repo, mode == "true-error", cause, mode != "false-error"}
			case "missing-after-create":
				repo.alterGet = func(int, *workflow.WorkflowExecution, bool) (*workflow.WorkflowExecution, bool) { return nil, false }
			case "contradictory-after-create":
				repo.alterGet = func(_ int, w *workflow.WorkflowExecution, _ bool) (*workflow.WorkflowExecution, bool) {
					return w, false
				}
			case "terminal-rejected":
				repo.rejectTerminal = true
			case "terminal-lied":
				repo.rejectTerminal = true
				repo.advanceTerminal = true
			}
			out, e := identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, InvocationReconciliation)
			if out.Kind != OutcomeReconciliation || identityState(s, "key") != requestUncertain {
				t.Fatal("uncertainty classified known")
			}
			if strings.Contains(mode, "error") && !errors.Is(e, cause) {
				t.Fatal("cause lost")
			}
			counts := []int{ids.calls, repo.creates, repo.gets, repo.claims, repo.completes, repo.fails}
			for i := 0; i < 3; i++ {
				_, e = identityInvoke(t, s, "key", "local.success", "event")
				invocationKind(t, e, InvocationReconciliation)
			}
			_, e = identityInvoke(t, s, "key", "local.failure", "event")
			invocationKind(t, e, InvocationDuplicate)
			if !reflect.DeepEqual(counts, []int{ids.calls, repo.creates, repo.gets, repo.claims, repo.completes, repo.fails}) {
				t.Fatal("uncertainty replayed")
			}
		})
	}
}
func TestLocalRequestIdentityCancellation(t *testing.T) {
	for _, phase := range []string{"before-create", "after-create"} {
		t.Run(phase, func(t *testing.T) {
			s, repo, ids := identityPrivate(t, localTestPlan(t))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if phase == "before-create" {
				ids.after = cancel
			} else {
				repo.afterCreate = cancel
			}
			c := identityCommand(t, s, "key", "local.success", "event")
			out, e := s.Invoker().InvokeLocalRequest(ctx, c)
			invocationKind(t, e, InvocationUnavailable)
			if !errors.Is(e, context.Canceled) || out.RunAttempted || repo.claims != 0 || identityState(s, "key") != requestKnown {
				t.Fatal("cancel boundary")
			}
			if phase == "after-create" && !freshInvocationSnapshot(out.Snapshot, invocationID) {
				t.Fatal("lost READY")
			}
			if phase == "before-create" && repo.creates != 0 {
				t.Fatal("cancel created")
			}
			_, e = s.Invoker().InvokeLocalRequest(context.Background(), c)
			invocationKind(t, e, InvocationDuplicate)
			if ids.calls != 1 {
				t.Fatal("cancel key reused")
			}
		})
	}
}
func TestLocalRequestIdentityPrecreationErrorsRetained(t *testing.T) {
	for _, mode := range []string{"ids", "resolver", "plan", "claim"} {
		t.Run(mode, func(t *testing.T) {
			s, repo, ids := identityPrivate(t, localTestPlan(t))
			want := InvocationUnavailable
			switch mode {
			case "ids":
				ids.err = errors.New("ids")
			case "resolver":
				s.service.invocation.deps.definitions = &invocationPlans{err: ErrDefinitionNotFound}
				want = InvocationDefinitionNotFound
			case "plan":
				s.service.invocation.deps.definitions = &invocationPlans{}
				want = InvocationInvalidDefinition
			case "claim":
				repo.rejectClaim = true
				want = InvocationClaimRejected
			}
			_, e := identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, want)
			_, e = identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, InvocationDuplicate)
			if identityState(s, "key") != requestKnown {
				t.Fatal("known refusal uncertain")
			}
		})
	}
}
func TestLocalRequestIdentityPanicRetention(t *testing.T) {
	for _, phase := range []string{"ids", "create", "callback"} {
		t.Run(phase, func(t *testing.T) {
			plan := localTestPlan(t)
			var callbacks atomic.Int32
			if phase == "callback" {
				r := workflow.NewHandlerRegistry()
				r.RegisterAction("local.noop", func() error { callbacks.Add(1); panic("private panic") })
				plan.runner = r
			}
			s, repo, ids := identityPrivate(t, plan)
			if phase == "ids" {
				ids.after = func() { panic("private panic") }
			}
			if phase == "create" {
				repo.afterCreate = func() { panic("private panic") }
			}
			func() {
				defer func() {
					if recover() != "private panic" {
						t.Error("panic swallowed")
					}
				}()
				identityInvoke(t, s, "key", "local.success", "event")
			}()
			if identityState(s, "key") != requestUncertain {
				t.Fatal("panic stranded as known/in-progress")
			}
			_, e := identityInvoke(t, s, "key", "local.success", "event")
			invocationKind(t, e, InvocationReconciliation)
			if ids.calls != 1 || callbacks.Load() > 1 {
				t.Fatal("panic replay")
			}
		})
	}
}
func TestLocalRequestIdentityConcurrent(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		t.Run(fmt.Sprint(conflict), func(t *testing.T) {
			plan := localTestPlan(t)
			registry := workflow.NewHandlerRegistry()
			entered, release := make(chan struct{}), make(chan struct{})
			var once sync.Once
			defer once.Do(func() { close(release) })
			var callbacks atomic.Int32
			registry.RegisterAction("local.noop", func() error { callbacks.Add(1); close(entered); <-release; return nil })
			plan.runner = registry
			s, repo, ids := identityPrivate(t, plan)
			winner := make(chan error, 1)
			go func() { _, e := identityInvoke(t, s, "key", "local.success", "event"); winner <- e }()
			select {
			case <-entered:
			case <-time.After(5 * time.Second):
				t.Fatal("winner did not enter")
			}
			if identityState(s, "key") != requestReserved {
				t.Fatal("in-flight state")
			}
			start := make(chan struct{})
			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					event := "event"
					if conflict {
						event = fmt.Sprint(i)
					}
					_, e := identityInvoke(t, s, "key", "local.success", event)
					invocationKind(t, e, InvocationDuplicate)
				}(i)
			}
			close(start)
			done := make(chan struct{})
			go func() { wg.Wait(); close(done) }()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("duplicates blocked on callback")
			}
			running, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
			assertInvocationState(t, running, invocationID, "worker", "RUNNING", "")
			once.Do(func() { close(release) })
			select {
			case e := <-winner:
				if e != nil {
					t.Fatal(e)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("winner stuck")
			}
			if callbacks.Load() != 1 || ids.calls != 1 || repo.creates != 1 {
				t.Fatal("duplicate invocation")
			}
			stored, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
			assertInvocationState(t, stored, invocationID, "worker", "COMPLETED", "")
		})
	}
}

func TestLocalRequestIdentityCompetingReservations(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		t.Run(fmt.Sprint(conflict), func(t *testing.T) {
			s, repo, calls := identityFixture(t, 10)
			start := make(chan struct{})
			var wg sync.WaitGroup
			var winners, duplicates atomic.Int32
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					event := "event"
					if conflict {
						event = fmt.Sprint(i)
					}
					c := identityCommand(t, s, "key", "local.success", event)
					out, e := s.Invoker().InvokeLocalRequest(context.Background(), c)
					if e == nil && out.Kind == OutcomeCompleted {
						winners.Add(1)
					} else {
						invocationKind(t, e, InvocationDuplicate)
						duplicates.Add(1)
					}
				}(i)
			}
			close(start)
			done := make(chan struct{})
			go func() { wg.Wait(); close(done) }()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("competing admission stuck")
			}
			if winners.Load() != 1 || duplicates.Load() != 9 || calls.Load() != 1 || len(s.service.reservations) != 1 {
				t.Fatal("multiple reservation winners")
			}
			stored, _ := repo.Get(invocationID)
			assertInvocationState(t, stored, invocationID, "identity-worker", "COMPLETED", "")
		})
	}
}
func TestLocalRequestIdentityIndependentProgress(t *testing.T) {
	repo := workflow.NewInMemoryWorkflowRepository()
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(release) })
	var ids atomic.Int32
	s, e := NewLocalRequestIdentityScope(repo, invocationIDFunc(func(context.Context) (string, error) {
		n := ids.Add(1)
		if n == 1 {
			close(entered)
			<-release
		}
		return fmt.Sprintf("exec-%032x", n), nil
	}), "owner", 2)
	if e != nil {
		t.Fatal(e)
	}
	first := make(chan error, 1)
	go func() { _, e := identityInvoke(t, s, "first", "local.success", "event"); first <- e }()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("first did not block")
	}
	second := make(chan error, 1)
	go func() { _, e := identityInvoke(t, s, "second", "local.success", "event"); second <- e }()
	select {
	case e := <-second:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("scope lock held across dependency")
	}
	once.Do(func() { close(release) })
	select {
	case e := <-first:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("first stuck")
	}
	if ids.Load() != 2 {
		t.Fatal("independent request lost")
	}
}
func TestLocalRequestIdentityCancelledCallback(t *testing.T) {
	plan := localTestPlan(t)
	registry := workflow.NewHandlerRegistry()
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(release) })
	var callbacks atomic.Int32
	registry.RegisterAction("local.noop", func() error { callbacks.Add(1); close(entered); <-release; return nil })
	plan.runner = registry
	s, repo, _ := identityPrivate(t, plan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := identityCommand(t, s, "key", "local.success", "event")
	type result struct {
		out InvocationOutcome
		err error
	}
	done := make(chan result, 1)
	go func() { o, e := s.Invoker().InvokeLocalRequest(ctx, c); done <- result{o, e} }()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("callback not entered")
	}
	cancel()
	_, e := s.Invoker().InvokeLocalRequest(context.Background(), c)
	invocationKind(t, e, InvocationDuplicate)
	once.Do(func() { close(release) })
	select {
	case r := <-done:
		if r.err != nil || r.out.Kind != OutcomeCompleted {
			t.Fatal("synchronous result discarded", r.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("callback stuck")
	}
	_, e = s.Invoker().InvokeLocalRequest(context.Background(), c)
	invocationKind(t, e, InvocationDuplicate)
	stored, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
	assertInvocationState(t, stored, invocationID, "worker", "COMPLETED", "")
	if callbacks.Load() != 1 {
		t.Fatal("cancelled callback replayed")
	}
}
