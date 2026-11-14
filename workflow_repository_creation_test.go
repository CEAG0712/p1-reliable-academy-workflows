package enrollment

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRepositoryCreateIfAbsentInvalid(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*WorkflowExecution)
	}{
		{"blank", func(w *WorkflowExecution) { w.WorkflowID = " \t" }},
		{"empty", func(w *WorkflowExecution) { w.WorkflowID = "" }},
		{"status", func(w *WorkflowExecution) { w.Status = "RUNNING" }},
		{"version", func(w *WorkflowExecution) { w.Version = 2 }},
		{"owner", func(w *WorkflowExecution) { w.ClaimOwner = "owner" }},
		{"error", func(w *WorkflowExecution) { w.LastError = "diagnostic" }},
		{"attempts", func(w *WorkflowExecution) { w.Attempts = 1 }},
		{"history", func(w *WorkflowExecution) { w.History = []WorkflowTransition{{To: "READY"}} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewInMemoryWorkflowRepository()
			w := NewWorkflowExecution("run")
			tc.change(w)
			if ok, err := r.CreateIfAbsent(w); ok || err == nil {
				t.Fatalf("invalid input: %v %v", ok, err)
			}
			if len(r.executions) != 0 {
				t.Fatal("invalid insertion")
			}
		})
	}
	for _, r := range []*InMemoryWorkflowRepository{nil, {}} {
		if ok, err := r.CreateIfAbsent(NewWorkflowExecution("run")); ok || err == nil {
			t.Fatal("unusable repository accepted")
		}
	}
	if ok, err := NewInMemoryWorkflowRepository().CreateIfAbsent(nil); ok || err == nil {
		t.Fatal("nil input accepted")
	}
}

func TestRepositoryCreateIfAbsentIsolationAndDuplicates(t *testing.T) {
	for _, status := range []string{"READY", "RUNNING", "COMPLETED", "FAILED"} {
		t.Run(status, func(t *testing.T) {
			r := NewInMemoryWorkflowRepository()
			w := NewWorkflowExecution("run")
			if ok, err := r.CreateIfAbsent(w); !ok || err != nil {
				t.Fatalf("insert: %v %v", ok, err)
			}
			w.Status = "tampered"
			stored, exists := r.Get("run")
			if !exists || stored.Status != "READY" || stored.Version != 1 || stored.Attempts != 0 || stored.ClaimOwner != "" || stored.LastError != "" || len(stored.History) != 0 {
				t.Fatalf("initial snapshot: %+v", stored)
			}
			if status != "READY" {
				r.TryClaim("run", "owner")
			}
			if status == "COMPLETED" {
				r.Complete("run", "owner", 2)
			}
			if status == "FAILED" {
				r.Fail("run", "owner", 2, "preserve reason")
			}
			before, _ := r.Get("run")
			if ok, err := r.CreateIfAbsent(NewWorkflowExecution("run")); ok || err != nil {
				t.Fatalf("duplicate: %v %v", ok, err)
			}
			invalid := NewWorkflowExecution("run")
			invalid.Status = "FAILED"
			if ok, err := r.CreateIfAbsent(invalid); ok || err == nil {
				t.Fatal("invalid duplicate hidden")
			}
			after, _ := r.Get("run")
			if !reflect.DeepEqual(before, after) {
				t.Fatal("duplicate changed metadata")
			}
			after.Status = "tampered"
			if len(after.History) != 0 {
				after.History[0].Reason = "tampered"
			}
			again, _ := r.Get("run")
			if !reflect.DeepEqual(before, again) {
				t.Fatal("snapshot escaped")
			}
		})
	}
}

func TestRepositoryCreateIfAbsentConcurrency(t *testing.T) {
	for _, mode := range []string{"same", "different", "with-save"} {
		t.Run(mode, func(t *testing.T) {
			r := NewInMemoryWorkflowRepository()
			start := make(chan struct{})
			var wg sync.WaitGroup
			var inserted atomic.Int32
			for i := 0; i < 24; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					id := "run"
					if mode == "different" {
						id = fmt.Sprintf("run-%d", i)
					}
					w := NewWorkflowExecution(id)
					if mode == "with-save" && i%2 == 0 {
						// Save permits legacy snapshots; this marker identifies its winner.
						w.LastError = "legacy-save"
						r.Save(w)
						return
					}
					ok, err := r.CreateIfAbsent(w)
					if err != nil {
						t.Errorf("create: %v", err)
					}
					if ok {
						inserted.Add(1)
					}
				}(i)
			}
			close(start)
			wg.Wait()
			if mode == "different" {
				if inserted.Load() != 24 || len(r.executions) != 24 {
					t.Fatal("unrelated IDs lost")
				}
				return
			}
			s, _ := r.Get("run")
			if len(r.executions) != 1 || s.Status != "READY" || s.Version != 1 {
				t.Fatal("conflicting authoritative records")
			}
			want := int32(1)
			if s.LastError == "legacy-save" {
				want = 0
			}
			if inserted.Load() != want {
				t.Fatalf("insertion outcomes inconsistent with winner: %d %+v", inserted.Load(), s)
			}
			before, _ := r.Get("run")
			r.Save(NewWorkflowExecution("run"))
			after, _ := r.Get("run")
			if !reflect.DeepEqual(before, after) {
				t.Fatal("Save overwrote winner")
			}
		})
	}
}

func TestRepositoryCreateIfAbsentSaveOrdering(t *testing.T) {
	for _, first := range []string{"save", "create"} {
		t.Run(first, func(t *testing.T) {
			repo := NewInMemoryWorkflowRepository()
			legacy := NewWorkflowExecution("run")
			legacy.LastError = "legacy marker"
			if first == "save" {
				repo.Save(legacy)
				before, _ := repo.Get("run")
				ok, e := repo.CreateIfAbsent(NewWorkflowExecution("run"))
				after, _ := repo.Get("run")
				if ok || e != nil || !reflect.DeepEqual(before, after) {
					t.Fatal("create replaced Save winner")
				}
			} else {
				ok, e := repo.CreateIfAbsent(NewWorkflowExecution("run"))
				if !ok || e != nil {
					t.Fatal("create failed")
				}
				before, _ := repo.Get("run")
				repo.Save(legacy)
				after, _ := repo.Get("run")
				if !reflect.DeepEqual(before, after) {
					t.Fatal("Save replaced create winner")
				}
			}
		})
	}
}
