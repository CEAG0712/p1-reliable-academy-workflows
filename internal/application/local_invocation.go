package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	workflow "example.com/p1-reliable-academy-workflows"
)

// localAuthority has nonzero size so distinct live scopes have distinct identity.
// It is a local fixture capability, NOT external caller authentication.
type localAuthority struct{ marker byte }

// LocalInvocationScope permits only package-owned harmless fixtures. Possession
// belongs to trusted local/test composition. Never attach it to an operational
// route: doing so would delegate this local capability to incoming requests.
// No arbitrary graph, callback, resolver or runner can be supplied publicly.
// Concurrent use requires a concurrency-safe injected ID source. Do not copy it.
type LocalInvocationScope struct {
	authority *localAuthority
	service   *invocationService
}

// NewLocalInvocationScope checks wiring without calling dependencies. The repo
// must come from NewInMemoryWorkflowRepository; an uninitialized repo fails at
// creation without running callbacks. Worker identity enforces ownership only.
func NewLocalInvocationScope(repo *workflow.InMemoryWorkflowRepository, ids ExecutionIDSource, workerID string) (*LocalInvocationScope, error) {
	authority := &localAuthority{marker: 1}
	service, err := newInvocationService(invocationDependencies{authority: authority, definitions: localDefinitionResolver{}, ids: ids, repo: repo, workerID: workerID})
	if err != nil {
		return nil, err
	}
	return &LocalInvocationScope{authority: authority, service: service}, nil
}
func (s *LocalInvocationScope) configured() bool {
	return s != nil && s.authority != nil && s.service != nil && s.service.deps.authority == s.authority && validInvocationDependencies(s.service.deps)
}

// Invoker returns nil for an unconstructed scope. Its service always checks the
// command capability itself, including calls that bypass an HTTP adapter.
func (s *LocalInvocationScope) Invoker() LocalWorkflowInvoker {
	if !s.configured() {
		return nil
	}
	return s.service
}

// Command grants a trusted host's local fixture selection, never JSON authority.
func (s *LocalInvocationScope) Command(selection SubmissionCommand) (TrustedLocalInvocationCommand, error) {
	if !s.configured() {
		return TrustedLocalInvocationCommand{}, invocationError(InvocationMissingAuthority, nil)
	}
	if !validLocalIdentifier(selection.WorkflowDefinitionID) || !validLocalIdentifier(selection.EventID) {
		return TrustedLocalInvocationCommand{}, invocationError(InvocationInvalidCommand, nil)
	}
	if !approvedLocalDefinition(selection.WorkflowDefinitionID) {
		return TrustedLocalInvocationCommand{}, invocationError(InvocationForbidden, nil)
	}
	return TrustedLocalInvocationCommand{authority: s.authority, selection: selection}, nil
}

// ApprovedDefinitions returns an independent support set, not authorization.
func (s *LocalInvocationScope) ApprovedDefinitions() map[string]struct{} {
	if !s.configured() {
		return nil
	}
	return map[string]struct{}{"local.success": {}, "local.failure": {}, "local.wait": {}, "local.condition": {}}
}
func approvedLocalDefinition(id string) bool {
	switch id {
	case "local.success", "local.failure", "local.wait", "local.condition":
		return true
	}
	return false
}
func validLocalIdentifier(s string) bool {
	if len(s) < 1 || len(s) > 128 {
		return false
	}
	for i := range len(s) {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-' || c == '.') {
			return false
		}
	}
	return true
}
func validExecutionID(id string) bool {
	if len(id) != 37 || !strings.HasPrefix(id, "exec-") {
		return false
	}
	for _, c := range id[5:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

// CryptoExecutionIDSource uses only standard-library local randomness. Collision
// rejection is authoritative; randomness does not imply idempotency or secrecy.
type CryptoExecutionIDSource struct{}

func (CryptoExecutionIDSource) NextExecutionID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", invocationError(InvocationInvalidCommand, nil)
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "exec-" + hex.EncodeToString(bytes[:]), nil
}

type localDefinitionResolver struct{}

func (localDefinitionResolver) ResolveApprovedDefinition(ctx context.Context, id string) (ApprovedLocalPlan, error) {
	if err := ctx.Err(); err != nil {
		return ApprovedLocalPlan{}, err
	}
	if !approvedLocalDefinition(id) {
		return ApprovedLocalPlan{}, ErrDefinitionNotFound
	}
	// Fresh maps and a fresh private registry per call. Registration never
	// overlaps invocation. These callbacks have no external or business effects.
	registry := workflow.NewHandlerRegistry()
	plan := ApprovedLocalPlan{
		definition: workflow.WorkflowDefinition{ID: id, Trigger: workflow.TriggerDefinition{EventType: "local.fixture"}, StartStep: "start", Steps: map[string]workflow.StepDefinition{
			"start": {Type: workflow.StepAction, Action: "local.noop"},
		}},
		runner: registry, actions: make(map[string]struct{}), conditions: make(map[string]struct{}),
	}
	for _, h := range []struct {
		name     string
		callback workflow.ActionHandler
	}{
		{"local.noop", func() error { return nil }},
		{"local.fail", func() error { return errors.New("controlled local fixture failure") }},
	} {
		if err := registry.RegisterAction(h.name, h.callback); err != nil {
			return ApprovedLocalPlan{}, err
		}
		plan.actions[h.name] = struct{}{}
	}
	if err := registry.RegisterCondition("local.check", func() (bool, error) { return true, nil }); err != nil {
		return ApprovedLocalPlan{}, err
	}
	plan.conditions["local.check"] = struct{}{}
	switch id {
	case "local.failure":
		plan.definition.Steps["start"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "local.fail"}
	case "local.wait", "local.condition":
		plan.definition.Steps["start"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "local.noop", Next: "boundary"}
		plan.definition.Steps["after"] = workflow.StepDefinition{Type: workflow.StepAction, Action: "local.noop"}
		if id == "local.wait" {
			plan.definition.Steps["boundary"] = workflow.StepDefinition{Type: workflow.StepWait, Delay: time.Second, Next: "after"}
		} else {
			plan.definition.Steps["boundary"] = workflow.StepDefinition{Type: workflow.StepCondition, Condition: "local.check", OnTrue: "after"}
		}
	}
	return plan, nil
}
func validateLocalPlan(p ApprovedLocalPlan) error {
	if nilInvocationDependency(p.runner) {
		return errors.New("local runner is required")
	}
	if err := p.definition.Validate(); err != nil {
		return err
	}
	// Validate the whole graph, including disconnected and post-handoff paths.
	// Manifests are built only after successful registration by trusted setup.
	for _, step := range p.definition.Steps {
		switch step.Type {
		case workflow.StepAction:
			if step.Action != "local.noop" && step.Action != "local.fail" {
				return errors.New("unsafe local action")
			}
			if _, ok := p.actions[step.Action]; !ok {
				return errors.New("unregistered local action")
			}
		case workflow.StepCondition:
			if step.Condition != "local.check" {
				return errors.New("unsafe local condition")
			}
			if _, ok := p.conditions[step.Condition]; !ok {
				return errors.New("unregistered local condition")
			}
		}
	}
	return nil
}
