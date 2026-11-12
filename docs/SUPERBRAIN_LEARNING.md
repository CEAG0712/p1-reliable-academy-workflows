# P1 — Superbrain 2.2 Learning Contract

## Purpose and authority

**SUPERBRAIN 2.2 — DEEP + GENTLE + MOMENTUM:** full technical rigor with low cognitive friction. This is the canonical educational presentation specification for P1 Mission Stations, not a lesson, interactive application, or engineering authorization source.

Follow [Interview Defense](INTERVIEW_DEFENSE.md) for assessment tracks, modes, grading, evidence classifications, and engineering authority. Preserve its hierarchy: AGENTS.md and explicit human authorization; approved roadmap, milestone boundaries and acceptance criteria; actual source, tests and verified evidence; then checkpoint documents and tracker reconciled against Git and source.

This contract cannot override those rules, activate tasks, change scope or trackers, authorize implementation, or bypass correctness and verification gates. Instructional completion is neither engineering acceptance nor independent mastery. Personal answers and assessment scores stay outside the public repository by default.

## Navigation

- [Default learning flow](#default-learning-flow)
- [Mission Station structure](#mission-station-structure)
- [Mission content](#mission-content)
- [Optional Deep Implementation](#optional-deep-implementation)
- [Ask About This Code / Context Rescue](#ask-about-this-code--context-rescue)
- [Final mental model and defense](#final-mental-model-and-defense)
- [Evidence discipline](#evidence-discipline)

## Default learning flow

**SCENARIO → SPOT THE BUG → REVEAL → CODE → UNDERSTAND → REMEMBER → DEFEND.**

CHEAT MODE is the default study experience. Teach with short, connected explanations, realistic decisions, actual code, and optional depth. Explain unfamiliar syntax before depending on it; preserve momentum without dropping technical rigor. A multiple-choice selection is a study interaction, not a graded mastery assessment or engineering gate.

Preserve STRICT/NO-CHEAT rules: use independent assessment only when requested or at an agreed formal checkpoint, ask one question at a time, and withhold hints and model defenses until after the response under the existing protocol. This contract changes no grading rules.

## Mission Station structure

Provide one complete, navigable Mission Station per milestone, with approximately ten concise missions where appropriate. Choose mission count for coherent coverage, not a quota. State project, milestone, source revision, learning goals, prerequisites, engineering status, and known boundaries. A station may teach accepted subtask evidence without claiming the whole milestone is accepted.

Where the teaching interface supports interactivity, preserve this preferred M3.7 interaction design:

- Previous/Next navigation, a mission map, and progress indication.
- Selectable multiple-choice answers, followed by answer and code reveal after selection.
- Collapsible Deep Implementation and optional model-defense reveal.
- Context Rescue questions carrying current milestone and mission metadata.
- A complete final visual mental model and concise final interview defense.

Prefer this interactive station over a long static Markdown lesson. If interactive controls are unavailable, provide a clearly labeled Markdown fallback with a linked contents map, numbered mission anchors, Previous/Next links, and equivalent technical content. Use disclosure sections where supported, but disclose that Markdown cannot enforce selection-triggered reveals or track interactive progress. This specification itself implements no interactive controls.

## Mission content

Each mission must include:

1. **Scenario:** a realistic business situation, affected stakeholder, and observable consequence. Label simulated customers and outcomes.
2. **Spot the bug / decision:** a multiple-choice engineering question with plausible alternatives grounded in the current mechanism or an explicitly hypothetical variation.
3. **Reveal:** after selection, explain the answer and why each alternative succeeds or fails. Apply the assessment-mode restrictions above.
4. **Code:** give the verified source filename, exact function or method, purpose, relevant real code excerpt, and applicable revision. Link supporting tests and distinguish inspected behavior from freshly executed checks.
5. **Understand:** explain the execution path and Go syntax in plain English, including relevant receivers, pointers, maps, returns, errors, and control flow. Connect each explanation to the shown code.
6. **Remember:** provide a short memory association and state where the analogy stops being accurate.
7. **Defend:** offer a natural, concise interview-ready statement grounded in evidence, with an honest limitation.
8. **Deep Implementation:** offer the optional engineering-depth layer defined below.
9. **Staff Engineer challenge:** introduce a meaningful failure, changed requirement, or alternative design; offer an optional model defense in CHEAT MODE, subject to STRICT/NO-CHEAT rules.
10. **Ask About This Code / Context Rescue:** preserve the current context and provide a route back to this mission.

## Optional Deep Implementation

Every mission must offer a substantial optional depth layer for the current code, not merely another one-paragraph summary. Keep the default path concise while making the full reasoning accessible.

Cover relevant aspects of the actual execution mechanism and source evidence; correctness invariants; design rationale; alternatives and tradeoffs; failure and recovery cases; concurrency and security implications; known limitations; and relevant tests, verification commands, results, revisions, and their limits.

Explain how the mechanism preserves each claimed invariant and what would break it. Separate tested recovery from proposed recovery. Discuss Backend, AI, and FDE implications where applicable; label unrelated dimensions NOT APPLICABLE rather than inventing integrations. Model defenses must distinguish implemented behavior from hypothetical extensions.

## Ask About This Code / Context Rescue

For an unfamiliar function or earlier milestone:

1. Preserve project, milestone, mission ID/title, current function, source revision, question, and assessment mode.
2. Locate and inspect verified source where available; identify unavailable evidence instead of inventing code.
3. Explain the business purpose.
4. Show the relevant real code with its file and function reference.
5. Explain the syntax simply.
6. Connect the prerequisite back to the current mission and execution path.
7. End with a short **Return to Mission** cue and navigation target.

Carry that metadata into rescue questions where the interface permits. Otherwise provide a copyable context prompt. Voice questions may use the surrounding ChatGPT interface where supported; do not promise a built-in microphone control. Rescue must respect the active assessment mode and must not silently disclose answers during STRICT/NO-CHEAT assessment.

## Final mental model and defense

End each station with one complete visual mental model connecting the business trigger, relevant components, execution path, state changes, failure paths, and system boundaries. Use a readable diagram with a legend; distinguish implemented, planned, and hypothetical elements. Explain the diagram in plain English and link its important nodes to source or tests.

Finish with a concise interview defense covering the business problem, actual mechanism, key correctness invariant, evidence, main tradeoff, and strongest limitation. Supply study scripts in CHEAT MODE; preserve independent-assessment disclosure rules otherwise. Scripts do not prove competence or personal contribution.

## Evidence discipline

Reuse the existing VERIFIED, PLANNED, HYPOTHETICAL, NOT APPLICABLE, and INSUFFICIENT EVIDENCE classifications. Label FDE simulations explicitly. Recheck source references and excerpts against the stated revision when preparing or refreshing a station; do not present stale excerpts as current code.

Separate repository facts and actual verification results from instructional examples, analogies, model answers, external guidance, and proposed designs. Identify historical versus fresh commands and results. Never invent unavailable source, tests, measurements, deployments, contributions, or external reference contents.

Passing local tests establishes only the exercised local behavior. It does not establish durability, distributed exactly-once effects, production operation, business outcomes, or mastery. Learning progress never updates engineering status or grants permission to implement the next milestone.
