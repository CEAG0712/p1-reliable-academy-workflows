package httpapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func validationBody(definition, event string) string {
	encoded, _ := json.Marshal(CreateExecutionRequest{WorkflowDefinitionID: definition, EventID: event})
	return string(encoded)
}

type invalidValidationCase struct {
	name, body, message string
	approved            map[string]struct{}
}

func invalidValidationCases() []invalidValidationCase {
	approved := map[string]struct{}{"definition-1": {}}
	valid := validationBody("definition-1", "event-1")
	cases := []invalidValidationCase{}
	add := func(name, body, message string) {
		cases = append(cases, invalidValidationCase{name, body, message, approved})
	}
	for _, body := range []string{"", " \n\t", "{", "{\"event_id\":", "{\"event_id\":\"x\",}", "null", "[]", "\"text\"", "42", "true", "false"} {
		add("shape/"+body, body, "body must be one JSON object")
	}
	add("missing/both", `{}`, "workflow_definition_id is required")
	add("missing/definition", `{"event_id":"event-1"}`, "workflow_definition_id is required")
	add("missing/event", `{"workflow_definition_id":"definition-1"}`, "event_id is required")
	for _, field := range []string{"workflow_definition_id", "event_id"} {
		for _, value := range []string{"null", "1", "true", "false", "[]", "{}"} {
			add("type/"+field+"/"+value, `{"`+field+`":`+value+`}`, field+" must be a string")
		}
		for _, value := range []string{"", " ", " padded", "padded ", "a b", "a\tb", "a\nb", "a\x00b", "café", "中", "😀", "a/b", "a:b", "a+b", strings.Repeat("a", 129)} {
			definition, event := "definition-1", "event-1"
			if field == "workflow_definition_id" {
				definition = value
			} else {
				event = value
			}
			add("identifier/"+field+"/"+value, validationBody(definition, event), field+" must contain 1 to 128 ASCII letters, digits, underscores, hyphens or periods")
		}
	}
	for _, field := range []string{"unknown", "payment_status", "Event_id", "approvedDefinitions"} {
		add("unknown/"+field, `{"workflow_definition_id":"definition-1","event_id":"event-1","`+field+`":"CONFIRMED"}`, "unsupported request field")
	}
	add("duplicate/event", `{"event_id":"event-1","event_id":"event-2","workflow_definition_id":"definition-1"}`, "duplicate request field")
	add("duplicate/definition", `{"workflow_definition_id":"definition-1","workflow_definition_id":"definition-1","event_id":"event-1"}`, "duplicate request field")
	add("duplicate/escaped", `{"event_id":"event-1","event_\u0069d":"event-2","workflow_definition_id":"definition-1"}`, "duplicate request field")
	for _, suffix := range []string{`{}`, `null`, `[]`, `garbage`} {
		add("trailing/"+suffix, valid+suffix, "body must be one JSON object")
	}
	add("oversize", valid+strings.Repeat(" ", 4097-len(valid)), "body must not exceed 4096 bytes")
	add("oversize/UTF8bytes", valid+strings.Repeat("é", 2100), "body must not exceed 4096 bytes")
	add("unsupported", validationBody("private-definition", "event-1"), "unsupported workflow definition")
	add("case-sensitive", validationBody("Definition-1", "event-1"), "unsupported workflow definition")
	cases = append(cases,
		invalidValidationCase{"nil allowlist", valid, "unsupported workflow definition", nil},
		invalidValidationCase{"empty allowlist", valid, "unsupported workflow definition", map[string]struct{}{}},
		invalidValidationCase{"no wildcard", valid, "unsupported workflow definition", map[string]struct{}{"*": {}}},
	)
	return cases
}

func assertValidationFailure(t *testing.T, request CreateExecutionRequest, failure *ValidationFailure, message string) {
	t.Helper()
	if request != (CreateExecutionRequest{}) {
		t.Fatalf("rejection leaked DTO: %#v", request)
	}
	if failure == nil {
		t.Fatal("expected rejection")
	}
	if failure.StatusCode != 400 || failure.Response.Error.Code != "INVALID_REQUEST" || failure.Response.Error.Message != message {
		t.Fatalf("failure = %#v, want HTTP 400 INVALID_REQUEST %q", failure, message)
	}
	encoded, _ := json.Marshal(message)
	assertContractJSON(t, failure.Response, `{"error":{"code":"INVALID_REQUEST","message":`+string(encoded)+`}}`)
}

func TestParseAndValidateRejectedRequests(t *testing.T) {
	for _, tc := range invalidValidationCases() {
		t.Run(tc.name, func(t *testing.T) {
			for repeat := 0; repeat < 2; repeat++ {
				request, failure := ParseAndValidateCreateExecutionRequest([]byte(tc.body), tc.approved)
				assertValidationFailure(t, request, failure, tc.message)
			}
		})
	}
}

func TestValidateApprovedDefinition(t *testing.T) {
	for _, identifier := range []string{"a", strings.Repeat("a", 128), "AZaz09_-.", "definition-1"} {
		t.Run(identifier, func(t *testing.T) {
			approved := map[string]struct{}{identifier: {}, "another": {}}
			body := validationBody(identifier, identifier)
			for _, input := range []string{body, " \n" + body + "\t", body + strings.Repeat(" ", 4096-len(body))} {
				request, failure := ParseAndValidateCreateExecutionRequest([]byte(input), approved)
				if failure != nil || request != (CreateExecutionRequest{WorkflowDefinitionID: identifier, EventID: identifier}) {
					t.Fatalf("request = %#v, failure = %#v", request, failure)
				}
			}
			if len(approved) != 2 {
				t.Fatal("validator mutated allowlist")
			}
		})
	}
	// Field order and JSON escapes are syntax, not identifier normalization.
	request, failure := ParseAndValidateCreateExecutionRequest([]byte(`{"event_id":"e","workflow_definition_id":"\u0061"}`), map[string]struct{}{"a": {}})
	if failure != nil || request.WorkflowDefinitionID != "a" || request.EventID != "e" {
		t.Fatalf("escaped/reordered request = %#v, failure = %#v", request, failure)
	}
}

// Test-only boundary makes downstream entry observable without introducing an
// application service, execution repository or real business operation.
func validationBoundary(body []byte, approved map[string]struct{}, downstream func(CreateExecutionRequest)) *ValidationFailure {
	request, failure := ParseAndValidateCreateExecutionRequest(body, approved)
	if failure != nil {
		return failure
	}
	downstream(request)
	return nil
}

func TestRejectedRequestHasNoSideEffects(t *testing.T) {
	for _, tc := range invalidValidationCases() {
		t.Run(tc.name, func(t *testing.T) {
			calls, repositoryCreations, businessCalls := 0, 0, 0
			var fakeRepository *struct{}
			downstream := func(CreateExecutionRequest) {
				calls++
				repositoryCreations++
				fakeRepository = &struct{}{}
				businessCalls++
			}
			failure := validationBoundary([]byte(tc.body), tc.approved, downstream)
			assertValidationFailure(t, CreateExecutionRequest{}, failure, tc.message)
			if calls != 0 || repositoryCreations != 0 || businessCalls != 0 || fakeRepository != nil {
				t.Fatal("rejected request entered downstream")
			}
			// Positive control proves this seam actually reaches the callback.
			failure = validationBoundary([]byte(validationBody("definition-1", "event-1")), map[string]struct{}{"definition-1": {}}, downstream)
			if failure != nil || calls != 1 || repositoryCreations != 1 || businessCalls != 1 || fakeRepository == nil {
				t.Fatal("positive control did not enter downstream exactly once")
			}
		})
	}
}
