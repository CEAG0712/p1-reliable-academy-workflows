P1 - Interview Defense Protocol
Purpose and scope
Use verified P1 evidence to build genuine, independently demonstrable engineering competence. Reuse this protocol for P2-P6 with project-specific evidence. Interview practice never substitutes for engineering verification or authorizes scope changes.
Three assessment tracks

1. Backend Engineering: business-to-technical requirements, Go code and data structures, APIs, state machines, concurrency, optimistic concurrency, idempotency, queues/databases, scaling, security, testing, observability, recovery, and architectural tradeoffs.
2. AI Engineering: AI integration boundaries, context and tool use, orchestration when applicable, structured output validation, evaluation, safety, latency, cost, fallback behavior, and agent/harness engineering.
3. Forward Deployed Engineering (FDE) and technical sales: customer discovery, stakeholder alignment, requirements, technical feasibility, solution design, presales demos, POC/pilot, ROI hypotheses, objections, delivery, adoption, executive communication, and customer success.
   Do not claim a capability is implemented if it is only planned. Simulated customers, savings, throughput and business results must be explicitly labeled simulated or proposed.
   Two modes
   CHEAT MODE (default during development): Present a realistic interview question, a concise high-quality model answer grounded in the current repository, a plain-language explanation, code/test evidence, alternative designs and tradeoffs, likely follow-ups, and a brief recall prompt. A supplied answer is study material, not proof of competence.

NO-CHEAT MODE (on request or at formal assessment): Ask one question at a time; do not reveal hints, scripts or model answers before the candidate responds. Probe assumptions, present counterexamples and novel scenarios, and evaluate only the answer actually given. Provide reference answers and corrections only after the round or when the candidate ends it.
Cadence and speed limits

- Bounded engineering task: 2-4 relevant CHEAT MODE questions, usually 5-10 minutes; keep development moving.
- Milestone: structured defense of architecture, selected code paths, correctness, failure modes and tradeoffs; optionally a NO-CHEAT checkpoint.
- Completed project: thorough multi-round NO-CHEAT defense across all three tracks, including a fresh design change, debugging scenario, business case, and explanation of limitations.
- After P1-P6: cross-project comparative architecture, unfamiliar systems, FDE discovery, and executive presentation.
- Defense readiness must not be equated with engineering acceptance. Weak answers create targeted study actions but do not silently revise accepted code or stop unrelated tasks.
  Evidence-first question generation

1. Read AGENTS.md, PROGRESS.md, docs/ROADMAP.md, relevant source files and tests.
2. Identify the exact change/milestone and its verified acceptance evidence.
3. Label each question IMPLEMENTED, PLANNED, or HYPOTHETICAL / FDE SIMULATION.
4. Cover why, how, what can fail, alternatives, what was measured, and what was personally contributed.
5. Do not use generated answers to fabricate professional employment, customer deployments, measured outcomes or individual contribution.
   Project-defense domains

- Business problem, requirements, KPIs, measurable outcomes and how evidence is obtained
- Architecture, boundaries, technology selection, alternatives and tradeoffs
- Code walkthroughs, data structures, Go concurrency, race and lifecycle correctness
- Optimistic concurrency, idempotency, ownership, leases and fencing
- API contracts, data storage, event delivery and distributed-systems guarantees
- Failure recovery, reconciliation, monitoring, testing and security
- Performance, capacity, bottlenecks, cost and production limitations
- Live debugging, incident response and unfamiliar architecture changes
- AI input/output validation, evaluation, reliability, cost, latency and safety
- Personal decisions, specific contributions, lessons learned and credible limitations
- FDE discovery, solution selling, POC/pilot, ROI, stakeholder alignment, objections, deployment and adoption
  Assessment rubric for NO-CHEAT MODE
  For each answer, score 0-4:
- 0 Unable to explain or materially incorrect
- 1 Recognizes concepts but cannot explain mechanism or evidence
- 2 Explains the basic mechanism accurately with substantial gaps
- 3 Defends the design using correct code/evidence, limitations and tradeoffs
- 4 Handles challenging follow-ups, unfamiliar variants and production/business consequences independently
  Track results separately for Backend, AI and FDE. Report question, score, demonstrated evidence, gap, corrective practice and suggested review point. Overall project-defense readiness requires strong independent performance across relevant areas, not a high memorized-answer count. Not every project implements every domain: distinguish implemented defense from hypothetical design questions.
  Personal study records
  Keep personal answers, rehearsal scripts, private performance scores and self-assessments outside the public repository by default. A public repository may retain the protocol and curated engineering evidence. A private defense log may use: date, project, milestone, track, question, mode, score (only for NO-CHEAT), code/test reference, gap and next review.
  M3 example
  For M3.1, ask: "Why must workflow definition models be independent of enrollment?" In CHEAT MODE, teach the answer and show the relevant model/test evidence once implemented. In NO-CHEAT MODE, ask for an independent explanation and challenge the candidate with a second business workflow. Never assert M3.1 is implemented before it is verified.
  Session handoff
  At each formal defense checkpoint, note covered areas, evidence inspected, unresolved knowledge gaps and next practice topic. Separate the defense assessment from source-code acceptance and PROGRESS.md engineering status.

## Supplementary Knowledge and Authority

P1 engineering authority remains:

1. AGENTS.md and explicit human authorization.
2. Approved roadmap, milestone boundaries and acceptance criteria.
3. Actual source code, tests and verified repository evidence.
4. PROGRESS.md, docs/HANDOFF.md and feature_list.json, reconciled against Git history and source.

EngineerGame Harness v2.0 is supplementary educational material, not engineering authority.

For the EngineerGame Superbrain 2.2 Mission Station teaching contract, see [Superbrain Learning](SUPERBRAIN_LEARNING.md). It supplements learning only; existing tracks, modes, grading, evidence rules, and engineering authority remain unchanged. It grants no implementation or task-activation authority.

Use relevant coach courses, DDIA, Clean Architecture, Clean Code and curated technical sources to deepen explanations and compare alternatives, only when those references are actually accessible and their contents have been checked.

Distinguish verified repository facts, coach/course guidance, reference-book principles, current external documentation, and proposed engineering recommendations. Do not silently convert educational recommendations into P1 implementation requirements. Preserve source-access and course-excerpt restrictions.

## Defense Evidence Classifications

- **VERIFIED:** Implemented behavior supported by inspected code and actual verification evidence, identified with the applicable revision and scope.
- **PLANNED:** Approved future work not yet implemented.
- **HYPOTHETICAL:** Proposed design, alternative or interview scenario. Label FDE simulations explicitly.
- **NOT APPLICABLE:** Topic outside the relevant task or milestone.
- **INSUFFICIENT EVIDENCE:** Claim that cannot currently be established; identify what evidence would resolve it.

Classify each defense question or claim and cite the applicable file, test, revision or external reference where practical. Existing historical labels such as IMPLEMENTED mean VERIFIED only when the supporting evidence is established; do not silently relabel unverified work.

Never infer deployment, production operation, business outcomes, personal engineering contributions or professional mastery from a passing local test.

## Milestone Defense Dimensions

For each relevant implemented milestone, assess:

1. Business problem and measurable requirements.
2. Actual implementation and code paths.
3. Engineering decisions and rationale.
4. Tradeoffs and reasonable alternatives.
5. Failure scenarios, recovery and known limitations.
6. Security, identity and authorization boundaries.
7. Tests, verification commands and evidence quality.
8. Operational consequences, customer impact and ROI assumptions or measured results.

Use the Backend Engineering, AI Engineering, and FDE tracks already defined above. Evaluate only dimensions relevant to implemented or explicitly planned work; mark others NOT APPLICABLE or PLANNED as appropriate. Do not imply that every task requires cloud deployment, LLM integration, production operations or newly implemented infrastructure.

## Enhanced Defense Modes

**CHEAT MODE (default during implementation):**

- Identify exact repository filenames and functions to study, checking that they actually exist.
- Explain prerequisites in accessible technical language and show actual code before interpreting it.
- Explain implementation choices, invariants, failure mechanisms, authorization boundaries and limitations.
- Provide complete model interview answers, relevant alternatives and tradeoffs.
- Identify tests, verification evidence and relevant Harness v2.0 reference concepts.
- Provide worked practice and spaced review prompts, without turning learning checks into engineering gates.

**STRICT MODE:** The challenging independent-assessment form of the existing NO-CHEAT MODE, not a separate engineering workflow.

- Ask one progressively difficult question at a time without showing model answers before the response.
- Introduce unfamiliar failures, modified requirements and counterexamples.
- Require code, test and architectural justification.
- Assess independently demonstrated answers with the existing 0-4 rubric.
- Record evidence, gaps and corrective practice separately for Backend, AI and FDE.

CHEAT MODE answers support study but do not establish independent mastery. Conduct STRICT/NO-CHEAT MODE only when requested or at the agreed formal defense checkpoint; during P1 development, preserve the user's CHEAT MODE preference. Neither mode changes acceptance or authorizes source edits.

## Project Correctness vs Career Mastery

Track these dimensions separately:

A. **Engineering correctness:** Verified implementation, tests, invariants, regression results and limitations.

B. **Technical understanding:** Ability to explain, reconstruct, debug and justify implemented mechanisms.

C. **Interview performance:** Independently defend unfamiliar scenarios, failure cases, alternatives and tradeoffs.

D. **Professional evidence:** Actual deployments, operations, measurements, customer work and personally performed contributions; distinguish demonstrations and simulations from real-world experience.

The approximately USD 177,000 annual compensation objective is a career target, not a P1 acceptance criterion or guaranteed outcome.

No interview, reference material, coaching exercise or documentation proposal may change P1 code, milestones, tracker status or acceptance criteria without separate authorization. Keep personal defense scores outside the public repository by default.
