package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	workflow "example.com/p1-reliable-academy-workflows"
)

const invocationID = "exec-00000000000000000000000000000001"

type invocationIDs struct {
	calls int
	ctx   context.Context
	id    string
	err   error
	after func()
}

func (s *invocationIDs) NextExecutionID(ctx context.Context) (string, error) {
	s.calls++
	s.ctx = ctx
	if s.after != nil {
		s.after()
	}
	return s.id, s.err
}
func invocationFixture(t *testing.T) (*LocalInvocationScope, *workflow.InMemoryWorkflowRepository, *invocationIDs) {
	t.Helper()
	r := workflow.NewInMemoryWorkflowRepository()
	ids := &invocationIDs{id: invocationID}
	s, err := NewLocalInvocationScope(r, ids, "local-worker")
	if err != nil {
		t.Fatal(err)
	}
	return s, r, ids
}
func invokeFixture(t *testing.T, s *LocalInvocationScope, ctx context.Context, definition string) (InvocationOutcome, error) {
	t.Helper()
	c, e := s.Command(SubmissionCommand{definition, "CONFIRMED"})
	if e != nil {
		t.Fatal(e)
	}
	return s.Invoker().InvokeLocal(ctx, c)
}
func invocationKind(t *testing.T, err error, want InvocationErrorKind) {
	t.Helper()
	var e *InvocationError
	if !errors.As(err, &e) || e.Kind != want {
		t.Fatalf("error=%v, want %s", err, want)
	}
}
func assertInvocationState(t *testing.T, s *workflow.WorkflowExecution, id, owner, status, reason string) {
	t.Helper()
	if s == nil {
		t.Fatal("missing snapshot")
	}
	history := []workflow.WorkflowTransition{{From: "READY", To: "RUNNING", WorkerID: owner, Version: 2}}
	version := uint64(2)
	if status != "RUNNING" {
		version = 3
		history = append(history, workflow.WorkflowTransition{From: "RUNNING", To: status, WorkerID: owner, Version: 3, Reason: reason})
	}
	if s.WorkflowID != id || s.Status != status || s.ClaimOwner != owner || s.Version != version || s.Attempts != 1 || s.LastError != reason || !reflect.DeepEqual(s.History, history) {
		t.Fatalf("bad state %+v history=%+v", s, s.History)
	}
}
func TestLocalInvocationConstruction(t *testing.T) {
	var ids *invocationIDs
	for _, source := range []ExecutionIDSource{nil, ids} {
		if s, e := NewLocalInvocationScope(workflow.NewInMemoryWorkflowRepository(), source, "worker"); s != nil || e == nil {
			t.Fatal("nil IDs accepted")
		}
	}
	if s, e := NewLocalInvocationScope(nil, &invocationIDs{}, "worker"); s != nil || e == nil {
		t.Fatal("nil repo")
	}
	if s, e := NewLocalInvocationScope(workflow.NewInMemoryWorkflowRepository(), &invocationIDs{}, " \t"); s != nil || e == nil {
		t.Fatal("blank owner")
	}
	scope, repo, source := invocationFixture(t)
	if source.calls != 0 {
		t.Fatal("constructor generated ID")
	}
	if _, ok := repo.Get(invocationID); ok {
		t.Fatal("constructor wrote")
	}
	for _, s := range []*invocationService{nil, {}} {
		if _, e := s.InvokeLocal(context.Background(), TrustedLocalInvocationCommand{}); e == nil {
			t.Fatal("zero service")
		}
	}
	for _, s := range []*LocalInvocationScope{nil, {}} {
		if _, e := s.Command(SubmissionCommand{"local.success", "event"}); e == nil {
			t.Fatal("zero scope")
		}
	}
	approved := scope.ApprovedDefinitions()
	delete(approved, "local.success")
	if _, e := scope.Command(SubmissionCommand{"local.success", "event"}); e != nil {
		t.Fatal("catalog escaped")
	}
}
func TestLocalInvocationRealOutcomes(t *testing.T) {
	for _, definition := range []string{"local.success", "local.failure", "local.wait", "local.condition"} {
		t.Run(definition, func(t *testing.T) {
			scope, repo, ids := invocationFixture(t)
			ctx := context.WithValue(context.Background(), struct{}{}, "correlation")
			out, err := invokeFixture(t, scope, ctx, definition)
			if out.Creation != CreationInserted || !out.RunAttempted || out.Claim != ClaimAccepted || out.ExecutionID != invocationID || out.WorkflowDefinitionID != definition || out.EventID != "CONFIRMED" || ids.calls != 1 || ids.ctx != ctx {
				t.Fatalf("outcome=%+v ids=%+v", out, ids)
			}
			status, reason := "COMPLETED", ""
			switch definition {
			case "local.failure":
				status = "FAILED"
				invocationKind(t, err, InvocationExecutionFailed)
				reason = out.Snapshot.LastError
				if reason == "" {
					t.Fatal("missing diagnostic")
				}
			case "local.wait", "local.condition":
				status = "RUNNING"
				if out.Kind != OutcomeHandoff || out.Cursor == "" || err != nil {
					t.Fatalf("handoff %+v %v", out, err)
				}
			default:
				if out.Kind != OutcomeCompleted || err != nil {
					t.Fatalf("completion %+v %v", out, err)
				}
			}
			assertInvocationState(t, out.Snapshot, invocationID, "local-worker", status, reason)
			stored, _ := repo.Get(invocationID)
			if !reflect.DeepEqual(out.Snapshot, stored) {
				t.Fatal("outcome != repository")
			}
			out.Snapshot.History[0].Reason = "tamper"
			again, _ := repo.Get(invocationID)
			if again.History[0].Reason != "" {
				t.Fatal("outcome escaped")
			}
		})
	}
}
func TestLocalInvocationTrustGate(t *testing.T) {
	scope, repo, ids := invocationFixture(t)
	other, _, _ := invocationFixture(t)
	for _, cmd := range []TrustedLocalInvocationCommand{
		{}, {selection: SubmissionCommand{"local.success", "CONFIRMED"}},
		{authority: other.authority, selection: SubmissionCommand{"local.success", "CONFIRMED"}},
		{authority: scope.authority, selection: SubmissionCommand{"entitlement.grant", "CONFIRMED"}},
	} {
		out, err := scope.Invoker().InvokeLocal(context.WithValue(context.Background(), "authorized", true), cmd)
		if err == nil || out.Creation != CreationNotAttempted || out.RunAttempted || ids.calls != 0 {
			t.Fatalf("gate bypass %+v %v", out, err)
		}
	}
	if _, ok := repo.Get(invocationID); ok {
		t.Fatal("untrusted creation")
	}
	for _, selection := range []SubmissionCommand{{"local.success", ""}, {"local.success", "a/b"}, {"payment.confirmed", "event"}} {
		if _, e := scope.Command(selection); e == nil {
			t.Fatal("bad selection")
		}
	}
}
func TestLocalInvocationDuplicateAndCorrelation(t *testing.T) {
	for _, status := range []string{"READY", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(status, func(t *testing.T) {
			scope, repo, ids := invocationFixture(t)
			repo.Save(workflow.NewWorkflowExecution(invocationID))
			if status != "READY" {
				repo.TryClaim(invocationID, "prior")
			}
			if status == "COMPLETED" {
				repo.Complete(invocationID, "prior", 2)
			}
			if status == "FAILED" {
				repo.Fail(invocationID, "prior", 2, "preserved")
			}
			before, _ := repo.Get(invocationID)
			out, err := invokeFixture(t, scope, context.Background(), "local.success")
			invocationKind(t, err, InvocationDuplicate)
			after, _ := repo.Get(invocationID)
			if out.Creation != CreationDuplicate || out.RunAttempted || out.Snapshot != nil || ids.calls != 1 || !reflect.DeepEqual(before, after) {
				t.Fatal("duplicate adopted/replayed/overwritten")
			}
		})
	}
	scope, repo, ids := invocationFixture(t)
	for i := 1; i <= 2; i++ {
		ids.id = fmt.Sprintf("exec-%032x", i)
		out, err := invokeFixture(t, scope, context.Background(), "local.success")
		if err != nil || out.EventID != "CONFIRMED" {
			t.Fatal(err)
		}
		if _, ok := repo.Get(ids.id); !ok {
			t.Fatal("missing separate run")
		}
	}
}
func TestLocalInvocationIDFailures(t *testing.T) {
	for _, id := range []string{"", "event", "exec-ABCDEF00000000000000000000000000", "exec-0000000000000000000000000000000g"} {
		scope, repo, ids := invocationFixture(t)
		ids.id = id
		out, err := invokeFixture(t, scope, context.Background(), "local.success")
		if err == nil || out.RunAttempted || out.Creation != CreationNotAttempted {
			t.Fatal("bad generated ID accepted")
		}
		if _, ok := repo.Get(id); ok {
			t.Fatal("bad ID inserted")
		}
	}
	scope, _, ids := invocationFixture(t)
	cause := errors.New("private ID source")
	ids.err = cause
	if _, err := invokeFixture(t, scope, context.Background(), "local.success"); !errors.Is(err, cause) {
		t.Fatal("lost ID cause")
	}
	id, err := (CryptoExecutionIDSource{}).NextExecutionID(context.Background())
	if err != nil || !validExecutionID(id) {
		t.Fatalf("crypto ID %s %v", id, err)
	}
}

type invocationPlans struct {
	plan  ApprovedLocalPlan
	err   error
	calls int
	ctx   context.Context
	id    string
}

func (p *invocationPlans) ResolveApprovedDefinition(ctx context.Context, id string) (ApprovedLocalPlan, error) {
	p.calls++
	p.ctx = ctx
	p.id = id
	return p.plan, p.err
}

type invocationRunner struct {
	delegate   LifecycleRunner
	calls      int
	repo       workflow.WorkflowLifecycleRepository
	id, worker string
	definition workflow.WorkflowDefinition
}

func (r *invocationRunner) RunActions(repo workflow.WorkflowLifecycleRepository, id, worker string, d workflow.WorkflowDefinition) (string, error) {
	r.calls++
	r.repo = repo
	r.id = id
	r.worker = worker
	r.definition = d
	return r.delegate.RunActions(repo, id, worker, d)
}

type invocationRepo struct {
	*workflow.InMemoryWorkflowRepository
	creates, gets, claims, completes, fails      int
	createError                                  error
	afterCreate                                  func()
	rejectClaim, rejectTerminal, advanceTerminal bool
	alterGet                                     func(int, *workflow.WorkflowExecution, bool) (*workflow.WorkflowExecution, bool)
	terminalID, terminalWorker, reason           string
	terminalVersion                              uint64
}

func (r *invocationRepo) CreateIfAbsent(w *workflow.WorkflowExecution) (bool, error) {
	r.creates++
	if r.createError != nil {
		return false, r.createError
	}
	ok, e := r.InMemoryWorkflowRepository.CreateIfAbsent(w)
	if ok && r.afterCreate != nil {
		r.afterCreate()
	}
	return ok, e
}
func (r *invocationRepo) Get(id string) (*workflow.WorkflowExecution, bool) {
	r.gets++
	s, ok := r.InMemoryWorkflowRepository.Get(id)
	if r.alterGet != nil {
		return r.alterGet(r.gets, s, ok)
	}
	return s, ok
}
func (r *invocationRepo) TryClaim(id, worker string) bool {
	r.claims++
	if r.rejectClaim {
		return false
	}
	return r.InMemoryWorkflowRepository.TryClaim(id, worker)
}
func (r *invocationRepo) Complete(id, worker string, v uint64) bool {
	r.completes++
	r.terminalID, r.terminalWorker, r.terminalVersion = id, worker, v
	if r.advanceTerminal {
		r.InMemoryWorkflowRepository.Complete(id, worker, v)
	}
	if r.rejectTerminal {
		return false
	}
	return r.InMemoryWorkflowRepository.Complete(id, worker, v)
}
func (r *invocationRepo) Fail(id, worker string, v uint64, reason string) bool {
	r.fails++
	r.terminalID, r.terminalWorker, r.terminalVersion, r.reason = id, worker, v, reason
	if r.advanceTerminal {
		r.InMemoryWorkflowRepository.Complete(id, worker, v)
	}
	if r.rejectTerminal {
		return false
	}
	return r.InMemoryWorkflowRepository.Fail(id, worker, v, reason)
}
func privateInvocation(t *testing.T, plan ApprovedLocalPlan) (*invocationService, *invocationRepo, *invocationPlans, *invocationIDs, TrustedLocalInvocationCommand) {
	t.Helper()
	repo := &invocationRepo{InMemoryWorkflowRepository: workflow.NewInMemoryWorkflowRepository()}
	plans := &invocationPlans{plan: plan}
	ids := &invocationIDs{id: invocationID}
	authority := &localAuthority{marker: 1}
	s, err := newInvocationService(invocationDependencies{authority, plans, ids, repo, "worker"})
	if err != nil {
		t.Fatal(err)
	}
	return s, repo, plans, ids, TrustedLocalInvocationCommand{authority, SubmissionCommand{"local.success", "event"}}
}
func localTestPlan(t *testing.T) ApprovedLocalPlan {
	t.Helper()
	p, e := (localDefinitionResolver{}).ResolveApprovedDefinition(context.Background(), "local.success")
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestLocalInvocationPrivateDependencies(t *testing.T) {
	s, repo, plans, ids, _ := privateInvocation(t, localTestPlan(t))
	good := s.deps
	var nilPlans *invocationPlans
	var nilIDs *invocationIDs
	var nilRepo *invocationRepo
	for _, edit := range []func(*invocationDependencies){
		func(d *invocationDependencies) { d.authority = nil }, func(d *invocationDependencies) { d.definitions = nil }, func(d *invocationDependencies) { d.definitions = nilPlans }, func(d *invocationDependencies) { d.ids = nil }, func(d *invocationDependencies) { d.ids = nilIDs }, func(d *invocationDependencies) { d.repo = nil }, func(d *invocationDependencies) { d.repo = nilRepo }, func(d *invocationDependencies) { d.workerID = " " },
	} {
		d := good
		edit(&d)
		if s, e := newInvocationService(d); s != nil || e == nil {
			t.Fatal("invalid dependency")
		}
	}
	if repo.creates+repo.gets+repo.claims+repo.completes+repo.fails+plans.calls+ids.calls != 0 {
		t.Fatal("constructor called collaborator")
	}
}
func TestLocalInvocationPlans(t *testing.T) {
	for _, name := range []string{"missing", "unavailable", "mismatch", "graph", "nil-runner", "typed-nil-runner", "unregistered", "forbidden-unreachable", "forbidden-after-wait", "unregistered-condition"} {
		t.Run(name, func(t *testing.T) {
			plan := localTestPlan(t)
			s, repo, plans, ids, c := privateInvocation(t, plan)
			want := InvocationInvalidDefinition
			switch name {
			case "missing":
				plans.err = fmt.Errorf("private: %w", ErrDefinitionNotFound)
				want = InvocationDefinitionNotFound
			case "unavailable":
				plans.err = errors.New("private dependency")
				want = InvocationUnavailable
			case "mismatch":
				plans.plan.definition.ID = "wrong"
			case "graph":
				plans.plan.definition.StartStep = "missing"
			case "nil-runner":
				plans.plan.runner = nil
			case "typed-nil-runner":
				var r *invocationRunner
				plans.plan.runner = r
			case "unregistered":
				plans.plan.actions = map[string]struct{}{}
			case "forbidden-unreachable":
				plans.plan.definition.Steps["unused"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "entitlement.grant"}
			case "forbidden-after-wait":
				plans.plan.definition.Steps["start"] = workflow.StepDefinition{Type: workflow.StepWait, Delay: time.Second, Next: "grant"}
				plans.plan.definition.Steps["grant"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "payment.confirm"}
			case "unregistered-condition":
				plans.plan.definition.Steps["start"] = workflow.StepDefinition{Type: workflow.StepCondition, Condition: "local.check"}
				plans.plan.conditions = nil
			}
			out, err := s.InvokeLocal(context.Background(), c)
			invocationKind(t, err, want)
			if ids.calls != 0 || repo.creates != 0 || out.RunAttempted {
				t.Fatal("bad plan entered creation")
			}
			if plans.err != nil && !errors.Is(err, plans.err) {
				t.Fatal("lost resolver cause")
			}
		})
	}
}
func TestLocalInvocationExactDelegationAndSnapshot(t *testing.T) {
	plan := localTestPlan(t)
	spy := &invocationRunner{delegate: plan.runner}
	plan.runner = spy
	s, repo, plans, ids, c := privateInvocation(t, plan)
	original := maps.Clone(plan.definition.Steps)
	ids.after = func() {
		plan.definition.Steps["start"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "entitlement.grant"}
	}
	ctx := context.WithValue(context.Background(), struct{}{}, "exact")
	out, err := s.InvokeLocal(ctx, c)
	if err != nil || plans.ctx != ctx || plans.id != "local.success" || plans.calls != 1 || ids.ctx != ctx || spy.calls != 1 || spy.id != invocationID || spy.worker != "worker" || !reflect.DeepEqual(spy.definition.Steps, original) {
		t.Fatalf("delegation %+v %v", out, err)
	}
	observer, ok := spy.repo.(*lifecycleObserver)
	if !ok || observer.repo != repo {
		t.Fatal("different lifecycle backing store")
	}
	if repo.creates != 1 || repo.claims != 1 || repo.gets != 3 || repo.completes != 1 || repo.fails != 0 || repo.terminalID != invocationID || repo.terminalWorker != "worker" || repo.terminalVersion != 2 {
		t.Fatalf("wrong operations %+v", repo)
	}
}
func TestLocalInvocationPrecreationRejections(t *testing.T) {
	s, repo, plans, ids, c := privateInvocation(t, localTestPlan(t))
	for _, cmd := range []TrustedLocalInvocationCommand{{}, {selection: c.selection}, {authority: &localAuthority{marker: 2}, selection: c.selection}, {authority: c.authority, selection: SubmissionCommand{"unsupported", "event"}}} {
		if _, e := s.InvokeLocal(context.Background(), cmd); e == nil {
			t.Fatal("gate bypass")
		}
	}
	if repo.creates+repo.gets+repo.claims+plans.calls+ids.calls != 0 {
		t.Fatal("gate rejection called dependency")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := s.InvokeLocal(ctx, c); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if plans.calls != 0 {
		t.Fatal("cancelled lookup")
	}
	if _, e := s.InvokeLocal(nil, c); e == nil {
		t.Fatal("nil context")
	}
	cause := errors.New("create unavailable")
	repo.createError = cause
	out, e := s.InvokeLocal(context.Background(), c)
	if !errors.Is(e, cause) || out.Creation != CreationRejected || out.RunAttempted || repo.claims != 0 {
		t.Fatal("creation error executed")
	}
}
func TestLocalInvocationCancellation(t *testing.T) {
	for _, phase := range []string{"before-create", "after-create"} {
		t.Run(phase, func(t *testing.T) {
			s, repo, _, ids, c := privateInvocation(t, localTestPlan(t))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if phase == "before-create" {
				ids.after = cancel
			} else {
				repo.afterCreate = cancel
			}
			out, err := s.InvokeLocal(ctx, c)
			if !errors.Is(err, context.Canceled) || out.RunAttempted || repo.claims != 0 {
				t.Fatalf("cancel %+v %v", out, err)
			}
			if phase == "after-create" {
				if out.Kind != OutcomeCreatedNotStarted || out.Creation != CreationInserted || !freshInvocationSnapshot(out.Snapshot, invocationID) {
					t.Fatal("lost READY run")
				}
			} else if repo.creates != 0 {
				t.Fatal("cancel created")
			}
		})
	}
	t.Run("during-callback", func(t *testing.T) {
		plan := localTestPlan(t)
		registry := workflow.NewHandlerRegistry()
		entered, release := make(chan struct{}), make(chan struct{})
		var calls atomic.Int32
		registry.RegisterAction("local.noop", func() error { calls.Add(1); close(entered); <-release; return nil })
		plan.runner = registry
		s, _, _, _, c := privateInvocation(t, plan)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var once sync.Once
		defer once.Do(func() { close(release) })
		done := make(chan struct {
			out InvocationOutcome
			err error
		}, 1)
		go func() {
			o, e := s.InvokeLocal(ctx, c)
			done <- struct {
				out InvocationOutcome
				err error
			}{o, e}
		}()
		select {
		case <-entered:
		case <-time.After(5 * time.Second):
			t.Fatal("callback did not enter")
		}
		cancel()
		once.Do(func() { close(release) })
		select {
		case r := <-done:
			if r.err != nil || r.out.Kind != OutcomeCompleted || calls.Load() != 1 {
				t.Fatalf("cancelled synchronous result %+v %v", r.out, r.err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("callback did not finish")
		}
	})
}
func TestLocalInvocationTerminalRejections(t *testing.T) {
	for _, failed := range []bool{false, true} {
		for _, advanced := range []bool{false, true} {
			t.Run(fmt.Sprintf("failed=%v/advanced=%v", failed, advanced), func(t *testing.T) {
				plan := localTestPlan(t)
				registry := workflow.NewHandlerRegistry()
				calls := 0
				cause := errors.New("private callback diagnostic")
				registry.RegisterAction("local.noop", func() error {
					calls++
					if failed {
						return cause
					}
					return nil
				})
				plan.runner = registry
				s, repo, _, _, c := privateInvocation(t, plan)
				repo.rejectTerminal = true
				repo.advanceTerminal = advanced
				out, err := s.InvokeLocal(context.Background(), c)
				invocationKind(t, err, InvocationReconciliation)
				if out.Kind != OutcomeReconciliation || calls != 1 || repo.claims != 1 || repo.completes+repo.fails != 1 || repo.terminalVersion != 2 || repo.terminalWorker != "worker" {
					t.Fatalf("uncertain replay %+v %+v", out, repo)
				}
				if failed && (!errors.Is(err, cause) || !strings.Contains(repo.reason, cause.Error())) {
					t.Fatal("lost failure")
				}
				want := "RUNNING"
				if advanced {
					want = "COMPLETED"
				}
				assertInvocationState(t, out.Snapshot, invocationID, "worker", want, "")
			})
		}
	}
}
func TestLocalInvocationClaimAndSnapshotRejections(t *testing.T) {
	for _, mode := range []string{"claim", "initial-missing", "initial-false-with-value", "initial-invalid", "claimed-missing", "claimed-owner", "final-missing", "final-invalid", "claim-invalid-final"} {
		t.Run(mode, func(t *testing.T) {
			plan := localTestPlan(t)
			calls := 0
			registry := workflow.NewHandlerRegistry()
			registry.RegisterAction("local.noop", func() error { calls++; return nil })
			plan.runner = registry
			s, repo, _, _, c := privateInvocation(t, plan)
			repo.rejectClaim = mode == "claim" || mode == "claim-invalid-final"
			repo.alterGet = func(n int, w *workflow.WorkflowExecution, ok bool) (*workflow.WorkflowExecution, bool) {
				if n == 1 {
					switch mode {
					case "initial-missing":
						return nil, false
					case "initial-false-with-value":
						return w, false
					case "initial-invalid":
						w.Version = 99
					}
				}
				if n == 2 && mode == "claim-invalid-final" {
					w.Status = "invalid"
				}
				if n == 2 {
					switch mode {
					case "claimed-missing":
						return nil, false
					case "claimed-owner":
						w.ClaimOwner = "foreign"
					}
				}
				if n == 3 {
					switch mode {
					case "final-missing":
						return nil, false
					case "final-invalid":
						w.Version = 99
					}
				}
				return w, ok
			}
			out, err := s.InvokeLocal(context.Background(), c)
			want := InvocationReconciliation
			if mode == "claim" {
				want = InvocationClaimRejected
			}
			invocationKind(t, err, want)
			if mode != "final-missing" && mode != "final-invalid" && calls != 0 {
				t.Fatal("rejection executed callback")
			}
			if repo.creates != 1 || repo.claims > 1 || calls > 1 {
				t.Fatal("retry")
			}
			if mode == "claim" && out.Claim != ClaimRejected {
				t.Fatal("claim not reported")
			}
		})
	}
}
func TestLocalInvocationHandoffBoundaries(t *testing.T) {
	for _, kind := range []workflow.StepType{workflow.StepWait, workflow.StepCondition} {
		for _, direct := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/direct=%v", kind, direct), func(t *testing.T) {
				id := "local.wait"
				if kind == workflow.StepCondition {
					id = "local.condition"
				}
				plan, e := (localDefinitionResolver{}).ResolveApprovedDefinition(context.Background(), id)
				if e != nil {
					t.Fatal(e)
				}
				plan.definition.ID = "local.success"
				if direct {
					plan.definition.StartStep = "boundary"
				}
				registry := workflow.NewHandlerRegistry()
				calls := 0
				registry.RegisterAction("local.noop", func() error { calls++; return nil })
				registry.RegisterCondition("local.check", func() (bool, error) { t.Error("condition evaluated"); return false, nil })
				plan.runner = registry
				s, repo, _, _, c := privateInvocation(t, plan)
				out, err := s.InvokeLocal(context.Background(), c)
				want := 1
				if direct {
					want = 0
				}
				if err != nil || out.Kind != OutcomeHandoff || out.Cursor != "boundary" || calls != want || repo.completes+repo.fails != 0 {
					t.Fatalf("boundary %+v %v calls=%d", out, err, calls)
				}
				assertInvocationState(t, out.Snapshot, invocationID, "worker", "RUNNING", "")
			})
		}
	}
}
func TestLocalInvocationFailureCauseAndOwner(t *testing.T) {
	plan := localTestPlan(t)
	registry := workflow.NewHandlerRegistry()
	cause := errors.New("sensitive sentinel")
	calls := 0
	registry.RegisterAction("local.noop", func() error { calls++; return cause })
	plan.runner = registry
	s, repo, _, _, c := privateInvocation(t, plan)
	out, err := s.InvokeLocal(context.Background(), c)
	invocationKind(t, err, InvocationExecutionFailed)
	if !errors.Is(err, cause) || strings.Contains(err.Error(), cause.Error()) || repo.fails != 1 || repo.completes != 0 || calls != 1 || repo.terminalVersion != 2 || repo.terminalWorker != "worker" || repo.reason == "" {
		t.Fatal("failure contract")
	}
	assertInvocationState(t, out.Snapshot, invocationID, "worker", "FAILED", repo.reason)
	before, _ := repo.Get(invocationID)
	if repo.InMemoryWorkflowRepository.Complete(invocationID, "foreign", 2) || repo.InMemoryWorkflowRepository.Fail(invocationID, "worker", 2, "stale") {
		t.Fatal("stale terminal accepted")
	}
	after, _ := repo.Get(invocationID)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("rejection mutated")
	}
}

func TestLocalInvocationCallbackOwnership(t *testing.T) {
	plan := localTestPlan(t)
	registry := workflow.NewHandlerRegistry()
	var repo *invocationRepo
	calls := 0
	if err := registry.RegisterAction("local.noop", func() error {
		calls++
		snapshot, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
		assertInvocationState(t, snapshot, invocationID, "worker", "RUNNING", "")
		before, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
		if repo.InMemoryWorkflowRepository.Complete(invocationID, "foreign", 2) || repo.InMemoryWorkflowRepository.Fail(invocationID, "worker", 1, "stale") {
			t.Error("wrong owner/version accepted")
		}
		after, _ := repo.InMemoryWorkflowRepository.Get(invocationID)
		if !reflect.DeepEqual(before, after) {
			t.Error("rejected update mutated")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	plan.runner = registry
	s, r, _, _, c := privateInvocation(t, plan)
	repo = r
	out, err := s.InvokeLocal(context.Background(), c)
	if err != nil || calls != 1 || out.Kind != OutcomeCompleted {
		t.Fatal(out, err)
	}
}

type invocationIDFunc func(context.Context) (string, error)

func (f invocationIDFunc) NextExecutionID(ctx context.Context) (string, error) { return f(ctx) }
func TestLocalInvocationConcurrentSameID(t *testing.T) {
	repo := workflow.NewInMemoryWorkflowRepository()
	var ids atomic.Int32
	scope, err := NewLocalInvocationScope(repo, invocationIDFunc(func(context.Context) (string, error) { ids.Add(1); return invocationID, nil }), "worker")
	if err != nil {
		t.Fatal(err)
	}
	command, err := scope.Command(SubmissionCommand{"local.success", "same-event"})
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	var winners, duplicates atomic.Int32
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			out, e := scope.Invoker().InvokeLocal(context.Background(), command)
			if out.Kind == OutcomeCompleted && e == nil {
				winners.Add(1)
			} else if out.Creation == CreationDuplicate && !out.RunAttempted {
				duplicates.Add(1)
			} else {
				t.Errorf("unexpected %+v %v", out, e)
			}
		}()
	}
	close(start)
	wg.Wait()
	if winners.Load() != 1 || duplicates.Load() != 15 || ids.Load() != 16 {
		t.Fatal("not a single winner")
	}
	snapshot, _ := repo.Get(invocationID)
	assertInvocationState(t, snapshot, invocationID, "worker", "COMPLETED", "")
}

func TestLocalInvocationObserverDelegates(t *testing.T) {
	repo := &invocationRepo{InMemoryWorkflowRepository: workflow.NewInMemoryWorkflowRepository()}
	repo.InMemoryWorkflowRepository.Save(workflow.NewWorkflowExecution(invocationID))
	o := &lifecycleObserver{repo: repo}
	if !o.TryClaim(invocationID, "worker") {
		t.Fatal("claim")
	}
	snapshot, exists := o.Get(invocationID)
	if !exists || snapshot == nil {
		t.Fatal("get")
	}
	if o.Complete(invocationID, "foreign", 2) {
		t.Fatal("bad owner accepted")
	}
	if !o.Fail(invocationID, "worker", 2, "exact reason") {
		t.Fatal("fail")
	}
	if repo.claims != 1 || repo.gets != 1 || repo.completes != 1 || repo.fails != 1 || len(o.operations) != 4 {
		t.Fatal("multiple delegation")
	}
	if o.operations[1].snapshot != snapshot || !o.operations[1].exists || o.operations[2].accepted || o.operations[2].owner != "foreign" || o.operations[3].reason != "exact reason" || o.operations[3].version != 2 || !o.operations[3].accepted {
		t.Fatal("changed arguments/results")
	}
}

func TestLocalInvocationJSONCannotMintAuthority(t *testing.T) {
	scope, repo, ids := invocationFixture(t)
	var command TrustedLocalInvocationCommand
	if err := json.Unmarshal([]byte(`{"authority":{"marker":1},"selection":{"WorkflowDefinitionID":"local.success","EventID":"CONFIRMED"},"authorized":true}`), &command); err != nil {
		t.Fatal(err)
	}
	if _, err := scope.Invoker().InvokeLocal(context.Background(), command); err == nil {
		t.Fatal("JSON minted authority")
	}
	if ids.calls != 0 {
		t.Fatal("JSON authorized ID")
	}
	if _, ok := repo.Get(invocationID); ok {
		t.Fatal("JSON created")
	}
}
