package application

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	workflow "example.com/p1-reliable-academy-workflows"
)

type statusReaderProbe struct {
	calls    int
	id       string
	snapshot *workflow.WorkflowExecution
	exists   bool
	after    func()
}

func (p *statusReaderProbe) Get(id string) (*workflow.WorkflowExecution, bool) {
	p.calls++
	p.id = id
	if p.after != nil {
		p.after()
	}
	return p.snapshot, p.exists
}

type statusPolicyProbe struct {
	calls    int
	ctx      context.Context
	id       string
	decision statusVisibility
	err      error
	after    func()
}

func (p *statusPolicyProbe) evaluate(ctx context.Context, id string) (statusVisibility, error) {
	p.calls++
	p.ctx = ctx
	p.id = id
	if p.after != nil {
		p.after()
	}
	return p.decision, p.err
}
func statusKind(t *testing.T, err error, want StatusQueryErrorKind) {
	t.Helper()
	var e *StatusQueryError
	if !errors.As(err, &e) || e.Kind != want {
		t.Fatalf("error=%v, want %s", err, want)
	}
}
func statusScope(t *testing.T, r ExecutionReader, ids ...string) *LocalStatusScope {
	t.Helper()
	s, e := NewLocalStatusScope(r, ids)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func statusCommand(t *testing.T, s *LocalStatusScope, id string) TrustedLocalStatusQuery {
	t.Helper()
	c, e := s.Command(id)
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func TestLocalStatusConstruction(t *testing.T) {
	var typedNil *statusReaderProbe
	for _, r := range []ExecutionReader{nil, typedNil} {
		if s, e := NewLocalStatusScope(r, nil); s != nil || e == nil {
			t.Fatal("nil reader accepted")
		}
	}
	p := &statusReaderProbe{}
	for _, id := range []string{"", strings.Repeat("a", 129), "a b", "é", "a/b", "a%2fb", "a\x00", "a:b"} {
		if s, e := NewLocalStatusScope(p, []string{id}); s != nil || e == nil {
			t.Fatalf("bad allowlist %q", id)
		}
	}
	s := statusScope(t, p, "run", "run", strings.Repeat("a", 128))
	if p.calls != 0 {
		t.Fatal("constructor read")
	}
	var nilPolicy *statusPolicyProbe
	for _, policy := range []statusVisibilityPolicy{nil, nilPolicy} {
		if s, e := newStatusService(s.authority, p, policy); s != nil || e == nil {
			t.Fatal("nil policy accepted")
		}
	}
	policy := &statusPolicyProbe{decision: statusVisible}
	if _, e := newStatusService(s.authority, p, policy); e != nil || p.calls != 0 || policy.calls != 0 {
		t.Fatal("constructor evaluated")
	}
	for _, scope := range []*LocalStatusScope{nil, {}} {
		if scope.Querier() != nil {
			t.Fatal("zero scope service")
		}
		if _, e := scope.Command("run"); e == nil {
			t.Fatal("zero scope command")
		}
	}
	for _, service := range []*statusService{nil, {}} {
		_, e := service.QueryLocal(context.Background(), TrustedLocalStatusQuery{})
		statusKind(t, e, StatusQueryUnavailable)
	}
	if _, e := newStatusService(nil, p, policy); e == nil {
		t.Fatal("missing service authority")
	}
}
func TestLocalStatusAuthority(t *testing.T) {
	p := &statusReaderProbe{}
	s := statusScope(t, p, "run")
	other := statusScope(t, p, "run")
	var decoded TrustedLocalStatusQuery
	if e := json.Unmarshal([]byte(`{"authority":{"marker":1},"executionID":"run","execution_id":"run","authorized":true}`), &decoded); e != nil {
		t.Fatal(e)
	}
	for _, c := range []TrustedLocalStatusQuery{{}, {executionID: "run"}, decoded} {
		_, e := s.Querier().QueryLocal(context.WithValue(context.Background(), "authorized", true), c)
		statusKind(t, e, StatusQueryMissingAuthority)
	}
	_, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, other, "run"))
	statusKind(t, e, StatusQueryForbidden)
	if p.calls != 0 {
		t.Fatal("authority bypass read")
	}
	for _, id := range []string{"", "a/b", " padded"} {
		if _, e := s.Command(id); e == nil {
			t.Fatal("invalid scope command")
		}
		_, e := s.Querier().QueryLocal(context.Background(), TrustedLocalStatusQuery{s.authority, id})
		statusKind(t, e, StatusQueryInvalidCommand)
	}
}
func TestLocalStatusVisibility(t *testing.T) {
	p := &statusReaderProbe{snapshot: workflow.NewWorkflowExecution("visible"), exists: true}
	ids := []string{"visible", "missing", "visible"}
	s, e := NewLocalStatusScope(p, ids)
	if e != nil {
		t.Fatal(e)
	}
	ids[0] = "concealed"
	ids[1] = "concealed"
	out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "visible"))
	if e != nil || out != (ExecutionStatusResult{"visible", "READY", 1}) || p.calls != 1 || p.id != "visible" {
		t.Fatalf("projection %+v %v", out, e)
	}
	for _, id := range []string{"concealed", "unknown"} {
		out, e = s.Querier().QueryLocal(context.Background(), statusCommand(t, s, id))
		statusKind(t, e, StatusQueryConcealed)
		if out != (ExecutionStatusResult{}) || p.calls != 1 {
			t.Fatal("concealed disclosure/read")
		}
	}
	p.snapshot = nil
	p.exists = false
	_, e = s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "missing"))
	statusKind(t, e, StatusQueryNotFound)
	if p.calls != 2 || p.id != "missing" {
		t.Fatal("missing reads")
	}
	deny := statusScope(t, p)
	_, e = deny.Querier().QueryLocal(context.Background(), statusCommand(t, deny, "visible"))
	statusKind(t, e, StatusQueryConcealed)
	if p.calls != 2 {
		t.Fatal("empty allowlist permitted read")
	}
}
func TestLocalStatusPolicyFailures(t *testing.T) {
	for _, decision := range []statusVisibility{statusVisible, statusConcealed, 99} {
		for _, failure := range []error{nil, errors.New("private tenant policy")} {
			p := &statusReaderProbe{snapshot: workflow.NewWorkflowExecution("run"), exists: true}
			scope := statusScope(t, p, "run")
			policy := &statusPolicyProbe{decision: decision, err: failure}
			s, e := newStatusService(scope.authority, p, policy)
			if e != nil {
				t.Fatal(e)
			}
			ctx := context.WithValue(context.Background(), struct{}{}, "exact")
			out, e := s.QueryLocal(ctx, statusCommand(t, scope, "run"))
			if policy.calls != 1 || policy.ctx != ctx || policy.id != "run" {
				t.Fatal("policy delegation")
			}
			if failure != nil || decision == 99 {
				statusKind(t, e, StatusQueryUnavailable)
				if failure != nil && !errors.Is(e, failure) {
					t.Fatal("lost cause")
				}
			} else if decision == statusConcealed {
				statusKind(t, e, StatusQueryConcealed)
			} else if e != nil || out.ExecutionID != "run" {
				t.Fatal(out, e)
			}
			want := 0
			if decision == statusVisible && failure == nil {
				want = 1
			}
			if p.calls != want {
				t.Fatal("policy ordering")
			}
		}
	}
}
func TestLocalStatusBasicSnapshots(t *testing.T) {
	for _, state := range []string{"READY", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(state, func(t *testing.T) {
			r := workflow.NewInMemoryWorkflowRepository()
			r.Save(workflow.NewWorkflowExecution("run"))
			if state != "READY" && !r.TryClaim("run", "owner") {
				t.Fatal("claim")
			}
			if state == "COMPLETED" && !r.Complete("run", "owner", 2) {
				t.Fatal("complete")
			}
			if state == "FAILED" && !r.Fail("run", "owner", 2, "private diagnostic") {
				t.Fatal("fail")
			}
			before, _ := r.Get("run")
			s := statusScope(t, r, "run")
			for i := 0; i < 3; i++ {
				out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "run"))
				if e != nil || out != (ExecutionStatusResult{"run", state, before.Version}) {
					t.Fatal(out, e)
				}
			}
			after, _ := r.Get("run")
			if !reflect.DeepEqual(before, after) {
				t.Fatal("query changed authoritative fields")
			}
		})
	}
}

func TestLocalStatusSnapshotValidation(t *testing.T) {
	cases := []struct {
		name string
		edit func(*workflow.WorkflowExecution)
	}{
		{"wrong-id", func(s *workflow.WorkflowExecution) { s.WorkflowID = "wrong" }},
		{"unknown-state", func(s *workflow.WorkflowExecution) { s.Status = "WAITING" }},
		{"zero-version", func(s *workflow.WorkflowExecution) { s.Version = 0 }},
		{"negative-attempts", func(s *workflow.WorkflowExecution) { s.Attempts = -1 }},
		{"zero-attempts", func(s *workflow.WorkflowExecution) { s.Attempts = 0 }},
		{"missing-owner", func(s *workflow.WorkflowExecution) { s.ClaimOwner = "" }},
		{"foreign-history-owner", func(s *workflow.WorkflowExecution) { s.History[0].WorkerID = "foreign" }},
		{"disconnected-history", func(s *workflow.WorkflowExecution) { s.History[1].From = "READY" }},
		{"zero-history-version", func(s *workflow.WorkflowExecution) { s.History[0].Version = 0 }},
		{"equal-history-version", func(s *workflow.WorkflowExecution) { s.History[0].Version = 3 }},
		{"decreasing-history-version", func(s *workflow.WorkflowExecution) { s.History[0].Version = 4 }},
		{"final-version-mismatch", func(s *workflow.WorkflowExecution) { s.Version = 4 }},
		{"final-state-mismatch", func(s *workflow.WorkflowExecution) { s.Status = "RUNNING" }},
		{"missing-failure-reason", func(s *workflow.WorkflowExecution) { s.LastError = ""; s.History[1].Reason = "" }},
		{"different-failure-reason", func(s *workflow.WorkflowExecution) { s.LastError = "different" }},
		{"nonfailure-reason", func(s *workflow.WorkflowExecution) { s.History[0].Reason = "private" }},
		{"terminal-without-history", func(s *workflow.WorkflowExecution) { s.History = nil }},
		{"running-without-history", func(s *workflow.WorkflowExecution) { s.Status = "RUNNING"; s.History = nil }},
		{"ready-owner", func(s *workflow.WorkflowExecution) { s.Status = "READY"; s.Attempts = 0; s.History = nil }},
		{"ready-attempts", func(s *workflow.WorkflowExecution) { s.Status = "READY"; s.ClaimOwner = ""; s.History = nil }},
		{"ready-history", func(s *workflow.WorkflowExecution) { s.Status = "READY"; s.ClaimOwner = ""; s.Attempts = 0 }},
		{"completed-error", func(s *workflow.WorkflowExecution) {
			s.Status = "COMPLETED"
			s.History[1].To = "COMPLETED"
			s.History[1].Reason = ""
		}},
		{"impossible-terminal-exit", func(s *workflow.WorkflowExecution) {
			s.History = append(s.History, workflow.WorkflowTransition{From: "FAILED", To: "RUNNING", WorkerID: "owner", Version: 4})
			s.Status = "RUNNING"
			s.Version = 4
		}},
		{"unsupported-transition", func(s *workflow.WorkflowExecution) { s.History[0].From = "NEW" }},
		{"ownerless-terminal", func(s *workflow.WorkflowExecution) {
			s.ClaimOwner = ""
			for i := range s.History {
				s.History[i].WorkerID = ""
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := workflow.NewInMemoryWorkflowRepository()
			r.Save(workflow.NewWorkflowExecution("run"))
			r.TryClaim("run", "owner")
			r.Fail("run", "owner", 2, "secret")
			snapshot, _ := r.Get("run")
			tc.edit(snapshot)
			// This private reader returns malformed data; the real repository is not repaired.
			p := &statusReaderProbe{snapshot: snapshot, exists: true}
			scope := statusScope(t, p, "run")
			before, _ := json.Marshal(snapshot)
			out, e := scope.Querier().QueryLocal(context.Background(), statusCommand(t, scope, "run"))
			statusKind(t, e, StatusQueryInvalidSnapshot)
			after, _ := json.Marshal(snapshot)
			if out != (ExecutionStatusResult{}) || p.calls != 1 || string(before) != string(after) {
				t.Fatal("corrupt projection/normalization")
			}
		})
	}
	for _, tc := range []struct {
		name   string
		s      *workflow.WorkflowExecution
		exists bool
		want   StatusQueryErrorKind
	}{
		{"missing", nil, false, StatusQueryNotFound}, {"contradictory", workflow.NewWorkflowExecution("run"), false, StatusQueryInvalidSnapshot}, {"nil-present", nil, true, StatusQueryInvalidSnapshot},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := &statusReaderProbe{snapshot: tc.s, exists: tc.exists}
			s := statusScope(t, p, "run")
			out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "run"))
			statusKind(t, e, tc.want)
			if out != (ExecutionStatusResult{}) || p.calls != 1 {
				t.Fatal("bad result")
			}
		})
	}
}
func TestLocalStatusLegacyCompatibility(t *testing.T) {
	for _, mode := range []string{"ready-marker", "ownerless-start", "claimed-marker", "complete", "save-if-version", "fail", "higher-version", "terminal-suffix"} {
		t.Run(mode, func(t *testing.T) {
			r := workflow.NewInMemoryWorkflowRepository()
			w := workflow.NewWorkflowExecution("workflow-001")
			switch mode {
			case "ready-marker", "claimed-marker":
				w.LastError = "legacy diagnostic"
			case "ownerless-start":
				w.LastError = "legacy start marker"
				if !w.TryStart() {
					t.Fatal("legacy start")
				}
			case "higher-version":
				w.Version = 18
			case "terminal-suffix":
				w.Status = "RUNNING"
				w.ClaimOwner = "owner"
				w.Attempts = 4
				w.Version = 29 // Legacy Save can retain a historical suffix.
			}
			r.Save(w)
			switch mode {
			case "claimed-marker", "complete", "save-if-version", "fail", "higher-version":
				if !r.TryClaim("workflow-001", "owner") {
					t.Fatal("claim")
				}
			}
			switch mode {
			case "complete":
				if !r.Complete("workflow-001", "owner", 2) {
					t.Fatal("complete")
				}
			case "save-if-version":
				s, _ := r.Get("workflow-001")
				s.Status = "COMPLETED"
				if !r.SaveIfVersion(s, s.Version) {
					t.Fatal("save version")
				}
			case "fail":
				if !r.Fail("workflow-001", "owner", 2, "legacy failure") {
					t.Fatal("fail")
				}
			case "higher-version":
				if !r.Complete("workflow-001", "owner", 19) {
					t.Fatal("higher complete")
				}
			case "terminal-suffix":
				if !r.Complete("workflow-001", "owner", 29) {
					t.Fatal("suffix complete")
				}
			}
			before, _ := r.Get("workflow-001")
			s := statusScope(t, r, "workflow-001")
			out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "workflow-001"))
			if e != nil || out != (ExecutionStatusResult{"workflow-001", before.Status, before.Version}) {
				t.Fatalf("legacy rejected %+v %v history=%+v", out, e, before.History)
			}
			after, _ := r.Get("workflow-001")
			if !reflect.DeepEqual(before, after) {
				t.Fatal("legacy mutated")
			}
		})
	}
}
func TestLocalStatusCancellation(t *testing.T) {
	for _, phase := range []string{"nil-context", "before-policy", "after-policy", "during-get", "deadline"} {
		t.Run(phase, func(t *testing.T) {
			p := &statusReaderProbe{snapshot: workflow.NewWorkflowExecution("run"), exists: true}
			scope := statusScope(t, p, "run")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			policy := &statusPolicyProbe{decision: statusVisible}
			s, e := newStatusService(scope.authority, p, policy)
			if e != nil {
				t.Fatal(e)
			}
			var queryCtx context.Context = ctx
			cause := error(context.Canceled)
			switch phase {
			case "nil-context":
				queryCtx = nil
			case "before-policy":
				cancel()
			case "after-policy":
				policy.after = cancel
			case "during-get":
				p.after = cancel
			case "deadline":
				var stop context.CancelFunc
				queryCtx, stop = context.WithDeadline(context.Background(), time.Unix(1, 0))
				defer stop()
				cause = context.DeadlineExceeded
			}
			out, e := s.QueryLocal(queryCtx, statusCommand(t, scope, "run"))
			if phase == "nil-context" {
				statusKind(t, e, StatusQueryInvalidCommand)
			} else {
				statusKind(t, e, StatusQueryUnavailable)
				if !errors.Is(e, cause) {
					t.Fatal("lost cancellation cause")
				}
			}
			wantReads, wantPolicy := 0, 0
			if phase == "after-policy" || phase == "during-get" {
				wantPolicy = 1
			}
			if phase == "during-get" {
				wantReads = 1
			}
			if out != (ExecutionStatusResult{}) || p.calls != wantReads || policy.calls != wantPolicy {
				t.Fatal("cancellation ordering")
			}
		})
	}
}

// Every mutable operation is an observable sentinel, but the service receives
// only ExecutionReader. The callbacks cannot be reached through that interface.
type statusMutationSentinel struct {
	reader                   ExecutionReader
	reads, writes, callbacks atomic.Int32
}

func (p *statusMutationSentinel) Get(id string) (*workflow.WorkflowExecution, bool) {
	p.reads.Add(1)
	return p.reader.Get(id)
}
func (p *statusMutationSentinel) Save(*workflow.WorkflowExecution) { p.writes.Add(1) }
func (p *statusMutationSentinel) CreateIfAbsent(*workflow.WorkflowExecution) (bool, error) {
	p.writes.Add(1)
	return false, nil
}
func (p *statusMutationSentinel) TryClaim(string, string) bool         { p.writes.Add(1); return false }
func (p *statusMutationSentinel) Complete(string, string, uint64) bool { p.writes.Add(1); return false }
func (p *statusMutationSentinel) Fail(string, string, uint64, string) bool {
	p.writes.Add(1)
	return false
}
func (p *statusMutationSentinel) SaveIfVersion(*workflow.WorkflowExecution, uint64) bool {
	p.writes.Add(1)
	return false
}
func (p *statusMutationSentinel) RunActions(workflow.WorkflowLifecycleRepository, string, string, workflow.WorkflowDefinition) (string, error) {
	p.callbacks.Add(1)
	return "", nil
}
func TestLocalStatusReadOnly(t *testing.T) {
	r := workflow.NewInMemoryWorkflowRepository()
	r.Save(workflow.NewWorkflowExecution("concealed"))
	r.Save(workflow.NewWorkflowExecution("visible"))
	r.TryClaim("visible", "owner")
	r.Fail("visible", "owner", 2, "secret tenant diagnostic")
	p := &statusMutationSentinel{reader: r}
	s := statusScope(t, p, "visible", "missing")
	before, _ := r.Get("visible")
	hidden, _ := r.Get("concealed")
	for _, id := range []string{"concealed", "unknown"} {
		out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, id))
		statusKind(t, e, StatusQueryConcealed)
		if out != (ExecutionStatusResult{}) {
			t.Fatal("disclosure")
		}
	}
	if p.reads.Load() != 0 || p.writes.Load() != 0 || p.callbacks.Load() != 0 {
		t.Fatal("concealed called collaborator")
	}
	for i := 0; i < 5; i++ {
		out, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "visible"))
		if e != nil || out.Status != "FAILED" {
			t.Fatal(out, e)
		}
	}
	_, e := s.Querier().QueryLocal(context.Background(), statusCommand(t, s, "missing"))
	statusKind(t, e, StatusQueryNotFound)
	after, _ := r.Get("visible")
	hiddenAfter, _ := r.Get("concealed")
	if p.reads.Load() != 6 || p.writes.Load() != 0 || p.callbacks.Load() != 0 || !reflect.DeepEqual(before, after) || !reflect.DeepEqual(hidden, hiddenAfter) {
		t.Fatal("query mutation or extra reads")
	}
}

type statusBarrierReader struct {
	reader   ExecutionReader
	captured chan struct{}
	release  chan struct{}
	calls    atomic.Int32
}

func (r *statusBarrierReader) Get(id string) (*workflow.WorkflowExecution, bool) {
	r.calls.Add(1)
	s, ok := r.reader.Get(id)
	close(r.captured)
	<-r.release
	return s, ok
}
func statusAwait(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("barrier timeout")
	}
}
func TestLocalStatusConcurrentSnapshots(t *testing.T) {
	for _, operation := range []string{"claim", "complete", "fail"} {
		t.Run(operation, func(t *testing.T) {
			r := workflow.NewInMemoryWorkflowRepository()
			r.Save(workflow.NewWorkflowExecution("run"))
			if operation != "claim" && !r.TryClaim("run", "owner") {
				t.Fatal("claim")
			}
			before, _ := r.Get("run")
			barrier := &statusBarrierReader{reader: r, captured: make(chan struct{}), release: make(chan struct{})}
			s := statusScope(t, barrier, "run")
			cmd := statusCommand(t, s, "run")
			done := make(chan struct{})
			var out ExecutionStatusResult
			var err error
			var once sync.Once
			defer once.Do(func() { close(barrier.release) })
			go func() { defer close(done); out, err = s.Querier().QueryLocal(context.Background(), cmd) }()
			statusAwait(t, barrier.captured)
			var ok bool
			switch operation {
			case "claim":
				ok = r.TryClaim("run", "owner")
			case "complete":
				ok = r.Complete("run", "owner", 2)
			case "fail":
				ok = r.Fail("run", "owner", 2, "new diagnostic")
			}
			if !ok {
				t.Fatal("concurrent lifecycle write")
			}
			advanced, _ := r.Get("run")
			once.Do(func() { close(barrier.release) })
			statusAwait(t, done)
			if err != nil || out != (ExecutionStatusResult{"run", before.Status, before.Version}) || barrier.calls.Load() != 1 {
				t.Fatal("snapshot reread or incoherent result", out, err)
			}
			after, _ := r.Get("run")
			if !reflect.DeepEqual(advanced, after) {
				t.Fatal("query overwrote concurrent writer")
			}
		})
	}
	t.Run("simultaneous-readers-and-writer", func(t *testing.T) {
		r := workflow.NewInMemoryWorkflowRepository()
		r.Save(workflow.NewWorkflowExecution("run"))
		p := &statusMutationSentinel{reader: r}
		s := statusScope(t, p, "run")
		cmd := statusCommand(t, s, "run")
		start := make(chan struct{})
		done := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < 24; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for j := 0; j < 20; j++ {
					out, e := s.Querier().QueryLocal(context.Background(), cmd)
					valid := out.Status == "READY" && out.Version == 1 || out.Status == "RUNNING" && out.Version == 2 || out.Status == "COMPLETED" && out.Version == 3
					if e != nil || out.ExecutionID != "run" || !valid {
						t.Errorf("incoherent %+v %v", out, e)
						return
					}
				}
			}()
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if !r.TryClaim("run", "owner") || !r.Complete("run", "owner", 2) {
				t.Error("writer")
			}
		}()
		close(start)
		go func() { wg.Wait(); close(done) }()
		statusAwait(t, done)
		if p.reads.Load() != 480 || p.writes.Load() != 0 || p.callbacks.Load() != 0 {
			t.Fatal("query wrote")
		}
	})
	t.Run("cancel-blocked-get", func(t *testing.T) {
		r := workflow.NewInMemoryWorkflowRepository()
		r.Save(workflow.NewWorkflowExecution("run"))
		p := &statusBarrierReader{reader: r, captured: make(chan struct{}), release: make(chan struct{})}
		s := statusScope(t, p, "run")
		cmd := statusCommand(t, s, "run")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan struct{})
		var out ExecutionStatusResult
		var err error
		var once sync.Once
		defer once.Do(func() { close(p.release) })
		go func() { defer close(done); out, err = s.Querier().QueryLocal(ctx, cmd) }()
		statusAwait(t, p.captured)
		cancel()
		select {
		case <-done:
			t.Fatal("Get was falsely interrupted")
		default:
		}
		once.Do(func() { close(p.release) })
		statusAwait(t, done)
		statusKind(t, err, StatusQueryUnavailable)
		if !errors.Is(err, context.Canceled) || out != (ExecutionStatusResult{}) || p.calls.Load() != 1 {
			t.Fatal("cancelled success")
		}
	})
}

func TestLocalStatusGateBeforePolicy(t *testing.T) {
	p := &statusReaderProbe{}
	scope := statusScope(t, p, "run")
	policy := &statusPolicyProbe{decision: statusVisible}
	s, e := newStatusService(scope.authority, p, policy)
	if e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		command TrustedLocalStatusQuery
		ctx     context.Context
		kind    StatusQueryErrorKind
	}{
		{TrustedLocalStatusQuery{}, context.Background(), StatusQueryMissingAuthority},
		{TrustedLocalStatusQuery{&localStatusAuthority{marker: 1}, "run"}, context.Background(), StatusQueryForbidden},
		{TrustedLocalStatusQuery{scope.authority, "a/b"}, context.Background(), StatusQueryInvalidCommand},
		{TrustedLocalStatusQuery{scope.authority, "run"}, nil, StatusQueryInvalidCommand},
	} {
		_, e := s.QueryLocal(tc.ctx, tc.command)
		statusKind(t, e, tc.kind)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e = s.QueryLocal(ctx, statusCommand(t, scope, "run"))
	statusKind(t, e, StatusQueryUnavailable)
	if policy.calls != 0 || p.calls != 0 {
		t.Fatal("invalid query reached policy/reader")
	}
	policy.decision = statusConcealed
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	policy.after = cancel
	_, e = s.QueryLocal(ctx, statusCommand(t, scope, "run"))
	statusKind(t, e, StatusQueryUnavailable)
	if policy.calls != 1 || p.calls != 0 {
		t.Fatal("denied cancellation read")
	}
}
func TestLocalStatusNoCallbackReplay(t *testing.T) {
	repo := workflow.NewInMemoryWorkflowRepository()
	repo.Save(workflow.NewWorkflowExecution("run"))
	registry := workflow.NewHandlerRegistry()
	var callbacks atomic.Int32
	if e := registry.RegisterAction("count", func() error { callbacks.Add(1); return nil }); e != nil {
		t.Fatal(e)
	}
	definition := workflow.WorkflowDefinition{ID: "fixture", Trigger: workflow.TriggerDefinition{EventType: "local"}, StartStep: "action", Steps: map[string]workflow.StepDefinition{
		"action": {Type: workflow.StepAction, Action: "count", Next: "wait"}, "wait": {Type: workflow.StepWait, Delay: time.Second},
	}}
	if cursor, e := registry.RunActions(repo, "run", "owner", definition); e != nil || cursor != "wait" {
		t.Fatal(cursor, e)
	}
	before, _ := repo.Get("run")
	s := statusScope(t, repo, "run")
	cmd := statusCommand(t, s, "run")
	for i := 0; i < 5; i++ {
		out, e := s.Querier().QueryLocal(context.Background(), cmd)
		if e != nil || out.Status != "RUNNING" || out.Version != 2 {
			t.Fatal(out, e)
		}
	}
	after, _ := repo.Get("run")
	if callbacks.Load() != 1 || !reflect.DeepEqual(before, after) {
		t.Fatal("status replayed callback or advanced handoff")
	}
}
