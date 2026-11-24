package investigator

// FixtureScope selects the single synthetic case available to Day 3.
// Its identifiers are demo inputs, with no production authority.
type FixtureScope struct {
	AcademyID  string `json:"academy_id"`
	WorkflowID string `json:"workflow_id"`
}

const (
	FixtureAcademyID  = "academy-demo"
	FixtureWorkflowID = "enrollment-demo-01"
	grantAccess       = "GRANT_COURSE_ACCESS"
	providerTimeout   = "TIMEOUT"
)

type evidence struct {
	Scope      FixtureScope     `json:"scope"`
	Status     string           `json:"workflow_status,omitempty"`
	Attempts   *int             `json:"attempts,omitempty"`
	Transition *transition      `json:"latest_transition,omitempty"`
	Provider   *providerOutcome `json:"provider_outcome,omitempty"`
}

type transition struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Version uint64 `json:"version"`
}

type providerOutcome struct {
	Operation string `json:"operation"`
	Outcome   string `json:"outcome"`
}

func loadFixture() evidence {
	attempts := 1
	return evidence{
		Scope:      FixtureScope{AcademyID: FixtureAcademyID, WorkflowID: FixtureWorkflowID},
		Status:     "FAILED",
		Attempts:   &attempts,
		Transition: &transition{From: "RUNNING", To: "FAILED", Version: 3},
		Provider:   &providerOutcome{Operation: grantAccess, Outcome: providerTimeout},
	}
}

// sanitizedEvidence copies only supported values. Free-text diagnostics,
// names, contact details and worker identifiers are outside the schema.
func sanitizedEvidence(raw evidence) evidence {
	clean := evidence{Scope: raw.Scope} // The caller has matched both fixture identifiers.
	switch raw.Status {
	case "READY", "RUNNING", "FAILED", "COMPLETED":
		clean.Status = raw.Status
	}
	if raw.Attempts != nil && *raw.Attempts >= 0 {
		attempts := *raw.Attempts
		clean.Attempts = &attempts
	}
	if entry := raw.Transition; entry != nil && entry.To == clean.Status && entry.Version >= 2 {
		if (entry.From == "READY" && entry.To == "RUNNING") ||
			(entry.From == "RUNNING" && (entry.To == "FAILED" || entry.To == "COMPLETED")) {
			latest := *entry
			clean.Transition = &latest
		}
	}
	if raw.Provider != nil && raw.Provider.Operation == grantAccess && raw.Provider.Outcome == providerTimeout {
		clean.Provider = &providerOutcome{Operation: grantAccess, Outcome: providerTimeout}
	}
	return clean
}
