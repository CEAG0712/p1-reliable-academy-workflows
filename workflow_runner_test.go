package enrollment

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Deliberately exposes no SaveIfVersion: the runner needs only the M2 lifecycle.
type lifecycleProbe struct {
	repo                              *InMemoryWorkflowRepository
	get                               func(string) (*WorkflowExecution, bool)
	rejectTerminal                    bool
	operations                        []string
	terminalID, terminalOwner, reason string
	terminalVersion                   uint64
}

func (p *lifecycleProbe) TryClaim(id, owner string) bool {
	p.operations = append(p.operations, "claim")
	return p.repo.TryClaim(id, owner)
}

func (p *lifecycleProbe) Get(id string) (*WorkflowExecution, bool) {
	p.operations = append(p.operations, "get")
	if p.get != nil {
		return p.get(id)
	}
	return p.repo.Get(id)
}

func (p *lifecycleProbe) Complete(id, owner string, version uint64) bool {
	p.operations = append(p.operations, "complete")
	p.terminalID, p.terminalOwner, p.terminalVersion = id, owner, version
	return !p.rejectTerminal && p.repo.Complete(id, owner, version)
}

func (p *lifecycleProbe) Fail(id, owner string, version uint64, reason string) bool {
	p.operations = append(p.operations, "fail")
	p.terminalID, p.terminalOwner, p.terminalVersion, p.reason = id, owner, version, reason
	return !p.rejectTerminal && p.repo.Fail(id, owner, version, reason)
}

func runnerFixture() (*InMemoryWorkflowRepository, *HandlerRegistry) {
	repo := NewInMemoryWorkflowRepository()
	repo.Save(NewWorkflowExecution("execution-42"))
	return repo, NewHandlerRegistry()
}

func assertRunnerState(t *testing.T, repo *InMemoryWorkflowRepository, status, reason string) {
	t.Helper()
	s, exists := repo.Get("execution-42")
	if !exists {
		t.Fatal("missing execution")
	}
	wantHistory := []WorkflowTransition{{From: "READY", To: "RUNNING", WorkerID: "worker-A", Version: 2}}
	wantVersion := uint64(2)
	if status != "RUNNING" {
		wantVersion = 3
		wantHistory = append(wantHistory, WorkflowTransition{From: "RUNNING", To: status, WorkerID: "worker-A", Version: 3, Reason: reason})
	}
	if s.WorkflowID != "execution-42" || s.Status != status || s.ClaimOwner != "worker-A" || s.Version != wantVersion || s.Attempts != 1 || s.LastError != reason || !reflect.DeepEqual(s.History, wantHistory) {
		t.Fatalf("unexpected stored lifecycle: %+v, history=%+v", s, s.History)
	}
}

func TestLifecycleRunnerTerminalOutcomes(t *testing.T) {
	for _, failed := range []string{"", "load", "transform", "store"} {
		t.Run("failure="+failed, func(t *testing.T) {
			repo, r := runnerFixture()
			p := &lifecycleProbe{repo: repo}
			cause := errors.New("provider unavailable")
			var calls []string
			for _, name := range []string{"load", "transform", "store"} {
				if err := r.RegisterAction(name, func() error {
					// Ownership must be acquired before even the first callback.
					assertRunnerState(t, repo, "RUNNING", "")
					calls = append(calls, name)
					if name == failed {
						return cause
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			cursor, err := r.RunActions(p, "execution-42", "worker-A", actionChainDefinition())
			status, operation, reason := "COMPLETED", "complete", ""
			wantCalls := []string{"load", "transform", "store"}
			wantCursor := ""
			if failed != "" {
				status, operation = "FAILED", "fail"
				if !errors.Is(err, cause) {
					t.Fatalf("lost original error: %v", err)
				}
				for i, name := range wantCalls {
					if name == failed {
						wantCalls = wantCalls[:i+1]
						break
					}
				}
				wantCursor = map[string]string{"load": "z-first", "transform": "a-middle", "store": "m-last"}[failed]
				reason = p.reason
				if reason == "" || !strings.Contains(reason, cause.Error()) || !strings.Contains(reason, wantCursor) {
					t.Fatalf("missing failure diagnostics: %q", reason)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if cursor != wantCursor || !reflect.DeepEqual(calls, wantCalls) {
				t.Fatalf("cursor=%q calls=%v", cursor, calls)
			}
			if p.terminalID != "execution-42" || p.terminalOwner != "worker-A" || p.terminalVersion != 2 || !reflect.DeepEqual(p.operations, []string{"claim", "get", operation}) {
				t.Fatalf("wrong lifecycle calls: %+v", p)
			}
			assertRunnerState(t, repo, status, reason)
		})
	}
}

func TestLifecycleRunnerRejectedClaim(t *testing.T) {
	for _, state := range []string{"missing", "empty owner", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(state, func(t *testing.T) {
			repo, r := runnerFixture()
			id, owner := "execution-42", "worker-A"
			switch state {
			case "missing":
				id = "missing"
			case "empty owner":
				owner = ""
			default:
				if !repo.TryClaim(id, "original-owner") {
					t.Fatal("fixture claim failed")
				}
				if state == "COMPLETED" && !repo.Complete(id, "original-owner", 2) {
					t.Fatal("fixture completion failed")
				}
				if state == "FAILED" && !repo.Fail(id, "original-owner", 2, "prior failure") {
					t.Fatal("fixture failure failed")
				}
			}
			before, _ := repo.Get("execution-42")
			if err := r.RegisterAction("load", func() error { t.Fatal("rejected claim invoked callback"); return nil }); err != nil {
				t.Fatal(err)
			}
			p := &lifecycleProbe{repo: repo}
			if _, err := r.RunActions(p, id, owner, actionChainDefinition()); err == nil || !strings.Contains(err.Error(), "claim") {
				t.Fatalf("claim rejection = %v", err)
			}
			after, _ := repo.Get("execution-42")
			if !reflect.DeepEqual(before, after) || !reflect.DeepEqual(p.operations, []string{"claim"}) {
				t.Fatal("rejected claim changed lifecycle or continued")
			}
		})
	}
}

func TestLifecycleRunnerHandoffs(t *testing.T) {
	for _, kind := range []StepType{StepCondition, StepWait} {
		for _, direct := range []bool{false, true} {
			t.Run(string(kind)+map[bool]string{false: "/after-action", true: "/direct"}[direct], func(t *testing.T) {
				repo, r := runnerFixture()
				p := &lifecycleProbe{repo: repo}
				d := waitDefinition()
				if kind == StepCondition {
					d.Steps["wait"] = StepDefinition{Type: StepCondition, Condition: "check", OnTrue: "after"}
				}
				if direct {
					d.StartStep = "wait"
				}
				calls := 0
				if err := r.RegisterAction("before", func() error { calls++; return nil }); err != nil {
					t.Fatal(err)
				}
				if err := r.RegisterAction("after", func() error { t.Fatal("crossed handoff"); return nil }); err != nil {
					t.Fatal(err)
				}
				if err := r.RegisterCondition("check", func() (bool, error) { t.Fatal("evaluated handoff"); return false, nil }); err != nil {
					t.Fatal(err)
				}
				cursor, err := r.RunActions(p, "execution-42", "worker-A", d)
				wantCalls := 1
				if direct {
					wantCalls = 0
				}
				if cursor != "wait" || err != nil || calls != wantCalls || !reflect.DeepEqual(p.operations, []string{"claim", "get"}) {
					t.Fatalf("handoff = (%q, %v), calls=%d operations=%v", cursor, err, calls, p.operations)
				}
				assertRunnerState(t, repo, "RUNNING", "")
				// Existing wait preparation still works, without a new lifecycle status.
				if kind == StepWait {
					w, err := PrepareWait(d, cursor, time.Time{})
					if err != nil || w.Snapshot().Cursor != cursor {
						t.Fatalf("prepare handoff: %v", err)
					}
				}
			})
		}
	}
}

func TestLifecycleRunnerTerminalRejection(t *testing.T) {
	for _, failAction := range []bool{false, true} {
		for _, stale := range []bool{false, true} {
			t.Run(map[bool]string{false: "complete", true: "fail"}[failAction]+map[bool]string{false: "/rejected", true: "/stale"}[stale], func(t *testing.T) {
				repo, r := runnerFixture()
				p := &lifecycleProbe{repo: repo, rejectTerminal: !stale}
				cause := errors.New("action error")
				d := actionChainDefinition()
				var calls []string
				if err := r.RegisterAction("load", func() error {
					calls = append(calls, "load")
					if stale && !repo.Complete("execution-42", "worker-A", 2) {
						t.Fatal("fixture advance failed")
					}
					if failAction {
						return cause
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				for _, name := range []string{"transform", "store"} {
					if err := r.RegisterAction(name, func() error {
						calls = append(calls, name)
						return nil
					}); err != nil {
						t.Fatal(err)
					}
				}
				cursor, err := r.RunActions(p, "execution-42", "worker-A", d)
				op := "completion"
				wantOperation := "complete"
				wantCalls := []string{"load", "transform", "store"}
				if failAction {
					op = "failure"
					wantOperation = "fail"
					wantCalls = []string{"load"}
				}
				if !reflect.DeepEqual(calls, wantCalls) {
					t.Fatalf("callbacks = %v, want %v", calls, wantCalls)
				}
				if err == nil || !strings.Contains(err.Error(), op) || !strings.Contains(err.Error(), "rejected") {
					t.Fatalf("terminal rejection hidden: %v", err)
				}
				if failAction && (!errors.Is(err, cause) || cursor != "z-first" || !strings.Contains(p.reason, cause.Error()) || !strings.Contains(p.reason, "z-first") || !strings.Contains(err.Error(), p.reason)) {
					t.Fatalf("failure context lost: (%q, %v)", cursor, err)
				}
				if !failAction && (cursor != "" || errors.Is(err, cause) || p.reason != "" || !strings.Contains(err.Error(), "actions succeeded") || !strings.Contains(err.Error(), "side effects")) {
					t.Fatalf("completion ambiguity missing: (%q, %v)", cursor, err)
				}
				if p.terminalID != "execution-42" || p.terminalVersion != 2 || p.terminalOwner != "worker-A" || !reflect.DeepEqual(p.operations, []string{"claim", "get", wantOperation}) {
					t.Fatalf("refreshed stale version or retried transition: %+v", p)
				}
				status := "RUNNING"
				if stale {
					status = "COMPLETED"
				}
				assertRunnerState(t, repo, status, "")
			})
		}
	}
}

func TestLifecycleRunnerInvalidClaimSnapshot(t *testing.T) {
	for _, bad := range []string{"missing", "nil", "owner", "status", "identity"} {
		t.Run(bad, func(t *testing.T) {
			repo, r := runnerFixture()
			p := &lifecycleProbe{repo: repo}
			p.get = func(id string) (*WorkflowExecution, bool) {
				s, _ := repo.Get(id)
				switch bad {
				case "missing":
					return nil, false
				case "nil":
					return nil, true
				case "owner":
					s.ClaimOwner = "other"
				case "status":
					s.Status = "COMPLETED"
				case "identity":
					s.WorkflowID = "other"
				}
				return s, true
			}
			if err := r.RegisterAction("load", func() error { t.Fatal("unusable claim invoked action"); return nil }); err != nil {
				t.Fatal(err)
			}
			if _, err := r.RunActions(p, "execution-42", "worker-A", actionChainDefinition()); err == nil {
				t.Fatal("unusable claim snapshot accepted")
			}
			if !reflect.DeepEqual(p.operations, []string{"claim", "get"}) {
				t.Fatalf("untrusted snapshot used for terminal transition: %v", p.operations)
			}
			assertRunnerState(t, repo, "RUNNING", "")
		})
	}
}

func TestLifecycleRunnerExecutionValidationFailure(t *testing.T) {
	for _, invalid := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown-handler", true: "invalid-definition"}[invalid], func(t *testing.T) {
			repo, r := runnerFixture()
			d := actionChainDefinition()
			if invalid {
				d.StartStep = "missing"
			}
			p := &lifecycleProbe{repo: repo}
			if _, err := r.RunActions(p, "execution-42", "worker-A", d); err == nil {
				t.Fatal("execution error hidden")
			}
			if p.reason == "" || !reflect.DeepEqual(p.operations, []string{"claim", "get", "fail"}) {
				t.Fatalf("error did not use Fail: %+v", p)
			}
			assertRunnerState(t, repo, "FAILED", p.reason)
		})
	}
}
