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
