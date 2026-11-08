package enrollment

import (
	"testing"
	"time"
)

func TestWorkflowDefinitionDescribesEnrollment(t *testing.T) {
	definition := WorkflowDefinition{
		ID:        "enrollment",
		Trigger:   TriggerDefinition{EventType: "payment.confirmed"},
		StartStep: "eligible",
		Steps: map[string]StepDefinition{
			"eligible": {Type: StepCondition, Condition: "enrollment.eligible", OnTrue: "grant", OnFalse: "notify"},
			"grant":    {Type: StepAction, Action: "enrollment.grant", Next: "notify"},
			"notify":   {Type: StepAction, Action: "notification.send"},
		},
	}

	if definition.ID != "enrollment" || definition.Trigger.EventType != "payment.confirmed" {
		t.Fatal("definition must retain its identity and configured trigger")
	}
	condition := definition.Steps[definition.StartStep]
	if condition.Type != StepCondition || condition.Condition != "enrollment.eligible" {
		t.Fatal("start step must describe the named eligibility condition")
	}
	if condition.OnTrue != "grant" || condition.OnFalse != "notify" {
		t.Fatal("condition must retain both configured branch targets")
	}
	grant := definition.Steps[condition.OnTrue]
	if grant.Type != StepAction || grant.Action != "enrollment.grant" || grant.Next != "notify" {
		t.Fatal("true branch must describe the grant action and its successor")
	}
	notify := definition.Steps[condition.OnFalse]
	if notify.Type != StepAction || notify.Action != "notification.send" || notify.Next != "" {
		t.Fatal("notification must describe a terminal action")
	}
}

func TestWorkflowDefinitionDescribesCustomerFollowUp(t *testing.T) {
	definition := WorkflowDefinition{
		ID:        "customer-follow-up",
		Trigger:   TriggerDefinition{EventType: "customer.created"},
		StartStep: "welcome",
		Steps: map[string]StepDefinition{
			"welcome":   {Type: StepAction, Action: "customer.welcome", Next: "delay"},
			"delay":     {Type: StepWait, Delay: 24 * time.Hour, Next: "follow-up"},
			"follow-up": {Type: StepAction, Action: "customer.follow-up"},
		},
	}

	if definition.ID != "customer-follow-up" || definition.Trigger.EventType != "customer.created" {
		t.Fatal("a different workflow must use the same generic definition and trigger models")
	}
	welcome := definition.Steps[definition.StartStep]
	if welcome.Type != StepAction || welcome.Action != "customer.welcome" || welcome.Next != "delay" {
		t.Fatal("start step must describe the welcome action and its successor")
	}
	delay := definition.Steps[welcome.Next]
	if delay.Type != StepWait || delay.Delay != 24*time.Hour || delay.Next != "follow-up" {
		t.Fatal("wait step must retain a duration and successor without executing a wait")
	}
	followUp := definition.Steps[delay.Next]
	if followUp.Type != StepAction || followUp.Action != "customer.follow-up" || followUp.Next != "" {
		t.Fatal("follow-up must describe a terminal action")
	}
}
