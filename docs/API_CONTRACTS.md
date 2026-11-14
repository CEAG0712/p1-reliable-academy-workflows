# Workflow API contracts — M4.2

## Implemented functionality

`GET /health` returns HTTP 200, `application/json; charset=utf-8`, and
`{"status":"ok"}`. It indicates HTTP responsiveness only. Gin's default
unknown-path and unsupported-method handling remains HTTP 404 with body
`404 page not found`. No HTTP authentication mechanism exists in M4.1.

M4.2 adds transport DTO definitions and deterministic contract tests, plus these
decisions. It registers no new routes, middleware or business operations. The
router and server entry point remain unchanged. DTOs alone neither validate
requests nor enforce lifecycle or security rules.

## Approved future HTTP contracts

| Method | Resource path | Availability |
|---|---|---|
| GET | `/health` | Implemented; unchanged |
| POST | `/v1/workflow-executions` | Future contract only; currently default 404 |
| GET | `/v1/workflow-executions/{execution_id}` | Future contract only; currently default 404 |

### Creation

`CreateExecutionRequest` has exactly two string fields:

```json
{"workflow_definition_id":"definition-1","event_id":"event-1"}
```

`WorkflowDefinitionID` identifies the intended definition. `EventID` is event
correlation data, not payment proof, authentication or an implemented idempotency
guarantee. The request contains no payment-status, tenant, worker or claim-owner
selection field. Ordinary malformed or invalid fields use the proposed HTTP 400
contract; no validation handler exists yet.

`CreateExecutionResponse` contains string `ExecutionID` and `Status` fields:

```json
{"execution_id":"execution-1","status":"READY"}
```

Prospective HTTP 201 Created applies only after a resource actually exists and
can be retrieved under the specified local contract. This does not claim durable
persistence. HTTP 202 is reserved for separately justified asynchronous
acceptance; it is not promised now and cannot guarantee scheduling or completion.
No creation handler is registered in M4.2.

### Status retrieval

`ExecutionStatusResponse` contains string `ExecutionID`, string `Status` and
uint64 `Version`, serialized as a JSON number:

```json
{"execution_id":"execution-1","status":"FAILED","version":3}
```

The only currently supported persisted lifecycle states are READY, RUNNING,
COMPLETED and FAILED. WAITING is not a stored execution status. A successfully
authorized retrieval of a FAILED execution can return HTTP 200; business failure
does not automatically mean HTTP 500.

This is an intentional projection of internal state. It excludes ClaimOwner,
transition history, attempts, LastError, mutex fields, tenant data and provider
information. There is no automatic serialization of the engine's execution type
and no projection handler yet.

### Errors and HTTP status matrix

`ErrorResponse` wraps `APIError`, whose Code and Message are strings:

```json
{"error":{"code":"INVALID_REQUEST","message":"workflow_definition_id is required"}}
```

Public messages must exclude internal errors and sensitive diagnostic details.
These documented machine-readable codes and statuses are future contracts, not
evidence of runtime enforcement. The envelope does not replace current health
or default Gin 404 handling.

| HTTP status | Meaning | Public error code |
|---|---|---|
| 200 | Successful authorized status retrieval | — |
| 201 | Execution resource actually created | — |
| 202 | Reserved for separately justified asynchronous acceptance; not promised now | — |
| 400 | Malformed request or invalid request fields | INVALID_REQUEST |
| 401 | Authentication missing or invalid | UNAUTHENTICATED |
| 403 | Authenticated caller lacks permission | FORBIDDEN |
| 404 | Resource missing or intentionally concealed | NOT_FOUND |
| 409 | Defined duplicate or conflicting-resource condition | CONFLICT |
| 500 | Unexpected internal error | INTERNAL_ERROR |
| 503 | Required service unavailable | UNAVAILABLE |

Ordinary validation failures use 400. Introducing 422 requires a separate approved
contract decision. Specific duplicate conditions still require later definition.

## Security requirements and trust blockers

### Authentication and authorization

Execution creation and status access require an approved caller identity model
before operational exposure. Worker claim-owner strings are not authentication.
No JWT validation, authentication middleware or identity-provider integration is
implemented. Authenticated callers require permission for the requested workflow
definition and execution; an execution ID alone does not grant access.
Authorization checks must precede restricted business operations. HTTP
authorization enforcement is not implemented.

### Tenant isolation

A client-submitted tenant ID is untrusted. The caller's authorized tenant context
must be established using trusted authority. The existing workflow repository
does not itself prove HTTP tenant ownership. Cross-tenant execution visibility
must be blocked until ownership checks exist. Concealed cross-tenant resources
may return 404 where required by future security policy. Adding a tenant-ID field
alone would not provide protection; these DTOs do not include one.

### Payment trust

`IsPaymentConfirmed()` compares a string to `CONFIRMED`.
`CanProceedToEnrollment()` evaluates business fields but does not authenticate
payment-provider origin. A client-supplied payment-status value cannot authorize
a grant. Provider authenticity, authorization and replay protection are not
implemented by the current HTTP foundation. No payment-driven or
entitlement-granting HTTP operation is permitted in M4.2.

### Execution acceptance and duplicates

The current repository is in-memory. `Save()` does not return an insertion or
duplicate result; it silently retains an existing execution. EventID alone does
not guarantee idempotency. No durable asynchronous processing or distributed
exactly-once behavior exists. Creation and duplicate semantics require later M4
implementation and verification; the DTOs do not establish those guarantees.

## Blocked or deferred operations

Operational creation/status handlers, caller identity, authorization, tenant
ownership checks, duplicate handling, scheduling, durable storage and engine
integration remain deferred. Enrollment endpoints, action-running endpoints,
payment callbacks, entitlement grants and administrative routes are outside
M4.2. M4.3–M4.10 require separate authorization. No cloud integration, deployment,
provider authentication or end-to-end workflow behavior is established by these
contracts or their local tests.


## M4.3 request-validation policy

`ParseAndValidateCreateExecutionRequest(body []byte, approvedDefinitions map[string]struct{})`
returns the existing `CreateExecutionRequest` and a nil failure on success.
On rejection it returns a zero-valued DTO and `*ValidationFailure`, containing
`StatusCode: 400` and the existing `ErrorResponse` with code `INVALID_REQUEST`.
Messages are fixed public text; parser errors, submitted unknown keys and
internal configuration are not exposed. No HTTP handler serializes this yet.

The supplied raw body must be at most 4096 bytes and contain exactly one JSON
object. Empty/malformed input, null, non-object values, concatenated JSON and
trailing non-whitespace content are rejected. Leading/trailing JSON whitespace
is allowed within the byte limit. This byte-slice check is not a streaming HTTP
body-read limit; a future operational handler must enforce that separately.

Only `workflow_definition_id` and `event_id` are accepted. Both are required JSON
strings. Unknown keys and duplicate keys (including equivalent escaped spellings)
are rejected by a top-level token decoder before a request can succeed. Nested
values are unsupported; this is not a recursive schema framework.

Both decoded identifiers must contain 1–128 ASCII bytes drawn exclusively from
`A-Z`, `a-z`, `0-9`, `_`, `-` and `.`. Empty, padded, whitespace-only, control,
Unicode and other punctuation values are invalid; values are never trimmed.
JSON escapes are decoded normally, then this same identifier policy applies.

The approved-definition map is a trusted, server-supplied support set. It must
contain the requested definition with exact, case-sensitive matching. Nil/empty
sets and unsupported selections fail with `unsupported workflow definition`;
there is no wildcard, default workflow or dynamically selected function. The
validator reads but does not mutate the set; callers must not mutate it
concurrently. Membership establishes local support, not caller authorization.

`payment_status`, including `CONFIRMED`, is an unknown field and is rejected.
The DTO and enrollment eligibility rules remain unchanged. No client payment
value is authority to grant an entitlement.

Validation performs no repository access, execution creation, business callback,
workflow invocation, grant, network call, scheduling or provider authentication.
Negative tests place an observable fake downstream callback after successful
validation only, assert zero calls/repository creation/business calls for every
rejection, and use valid-request controls to prove the seam reaches downstream.
The fake repository is a test sentinel, not an engine repository integration;
production validation imports only standard-library decoding/HTTP packages and
has no downstream dependency. These tests prove this local validation boundary,
not future handler-to-engine integration.

No operational route, middleware, authentication or application service is added.
Health and Gin's default 404 behavior remain unchanged. Historical M4.2 statements
above describe that milestone; only this validation function is newly implemented.

## M4.5 — Controlled local fixture invocation (implementation pending source review)

This section describes an unregistered local/test adapter. Production `NewRouter`
continues to expose only GET `/health`. No operational execution endpoint, HTTP
status retrieval, caller authentication or tenant ownership implementation is added.
Historical milestone statements above retain their original scope.

### Trust and composition

`application.NewLocalInvocationScope(repo, ids, workerID)` binds a concrete existing
in-memory repository, trusted execution-ID source and nonblank server worker ID.
It performs structural checks without calling ID, definition, repository or runner
operations. Supply a repository from `NewInMemoryWorkflowRepository`; uninitialized
storage fails at creation without invoking actions. Injected Go dependencies belong
to trusted host composition; this is not a sandbox for malicious process code.

The scope can issue a `TrustedLocalInvocationCommand` containing private authority
and selection fields. `InvokeLocal` independently requires the exact scope token.
Zero, JSON-populated and foreign-scope commands cannot obtain execution authority.
Supported-definition membership, client context values, EventID=CONFIRMED, worker
strings and arbitrary tenant/payment fields are not identity or payment proof.

The fixed catalog is local.success, local.failure, local.wait and local.condition.
Each call builds a fresh private graph and registry. The only actions are local.noop
(return nil) and local.fail (return a fixed fixture error); local.check returns
true/nil and is not evaluated by this action runner. There are no provider, grant,
payment, filesystem or network effects. All graph action/condition references,
including disconnected and post-handoff paths, must be safe and appear in the
successful-registration manifest. Existing graph Validate runs independently.
There is no public callback, resolver, graph or runner injection path.

`httpapi.NewLocalInvocationHandler(scope)` creates a handler but registers nothing.
Only tests compose it with `gin.New()` and `httptest`, without a listener. It reads
at most 4097 bytes and reuses unchanged M4.3 validation. The trusted test harness
controls requests and owns the capability the handler attaches. Registering this
handler operationally would delegate its scope to incoming clients and is prohibited
within M4.5. The sealed local scope is NOT production authentication.

### Identity, creation and association

Execution IDs are server-controlled: `exec-` plus 32 lowercase hexadecimal characters.
`CryptoExecutionIDSource` uses 16 random bytes from the standard library; tests
inject deterministic IDs. The format is checked before repository access. Worker
ID is a separate ownership string. WorkflowDefinition.ID selects configuration;
EventID is correlation; WorkflowExecution.WorkflowID is the execution-record key.
These roles must not be substituted for one another.

The additive repository `CreateIfAbsent(*WorkflowExecution) (bool, error)` clones
and validates fresh READY/version-1 data with empty owner/error/history and zero
attempts. Blank IDs, invalid candidates and uninitialized storage return false/error.
Duplicate detection and insertion use the same existing map and lock as Save:
true/nil proves insertion by this call; false/nil means duplicate with no mutation.
Invalid duplicates remain errors. true/error violates the creation contract and is
classified as uncertain if a faulty injected collaborator ever returns it.
Existing Save and all lifecycle method bodies/semantics are unchanged.
Callers must not directly mutate candidate fields concurrently with snapshotting.

An inserted record is immediately retrievable through the same process-local
repository's Get; the application verifies that snapshot before runner entry.
A duplicate does not adopt the existing record, run actions or generate a replacement.
Repeated EventID with distinct generated IDs can create independent runs. Neither
EventID nor random identity establishes request idempotency or exactly-once behavior.

Definition, event, execution identity, plan and handoff cursor are associated only
within the call/result. The execution model stores no definition/event/cursor metadata.
This provides no durable association, M4.6 visibility policy or continuation storage.

### Runner outcomes and interruption

The unchanged real HandlerRegistry.RunActions receives the same authoritative
repository through a private call-local observer. The observer forwards each call
exactly once, unchanged, and records its actual result. It holds no execution map,
performs no independent state write and never refreshes versions or retries.

Normal completion requires accepted claim/Complete plus the coherent COMPLETED/v3
snapshot, retained owner, attempts 1 and claim/completion history. Action failure
requires accepted Fail plus FAILED/v3 and the same nonempty diagnostic in state/history.
The application retains causes internally; public messages are fixed safe text.
A rejected terminal write always requires reconciliation, even if a later snapshot
is COMPLETED. Missing/inconsistent snapshots also require reconciliation. There is
no error-string parsing, fallback terminal write, replacement execution or replay.

WAIT/condition returns a nonempty cursor and RUNNING/v2 with claim history only.
It does not persist WAITING, evaluate a condition, prepare a deadline or resume.
Cancellation before creation prevents creation; cancellation detected after insertion
and before runner entry retains READY/v1 as created_not_started. The synchronous
runner and callbacks accept no context: cancellation during invocation does not
stop them or discard a valid observed result. No background goroutine simulates
cancellation. Crashes/panics/lost responses can leave READY/RUNNING or a terminal
record without a delivered outcome. No queue/scheduler/recovery mechanism resumes it.

### Isolated HTTP response contract

Existing request/response DTOs and ErrorResponse remain unchanged. HTTP 201 is used
only for explicit insertion plus a retrievable coherent application outcome:
COMPLETED/v3, FAILED/v3, RUNNING/v2 handoff, or READY/v1 created_not_started. It reports
resource creation and actual observed status, not universal action success,
asynchronous scheduling or durable persistence. No HTTP 202 or WAITING is introduced.

| Condition | HTTP | Existing code |
|---|---|---|
| Invalid transport/command | 400 | INVALID_REQUEST |
| Missing local capability (not external-user authentication) | 401 | UNAUTHENTICATED |
| Foreign capability/disallowed local operation | 403 | FORBIDDEN |
| Missing selected definition | 404 | NOT_FOUND |
| Duplicate generated key or rejected claim | 409 | CONFLICT |
| Unavailable dependency or cancellation before creation | 503 | UNAVAILABLE |
| Invalid plan, inconsistent result or reconciliation required | 500 | INTERNAL_ERROR |

Cancellation after confirmed creation can still project 201 READY if a response is
writable. A known FAILED execution can project 201 FAILED despite an internal
execution error. Reconciliation never projects success; its safe error envelope
cannot carry the full internal outcome/identity. The controlled test harness can
inspect that outcome and repository. No public recovery/status API is implemented.
Owner, attempts, diagnostics, history, cursor and capabilities are never serialized.

M4.5 establishes no real payment authenticity, tenant authorization, entitlement grant,
durable queue/storage/scheduler, automatic retry/resume, cloud integration or
distributed exactly-once guarantee. Acceptance awaits human source review.

## M4.6 — Local execution status (implementation awaiting human source review)

This additive checkpoint describes a constructible, unregistered status adapter.
Historical statements above describe their respective milestones. Production
`NewRouter` still exposes only GET `/health`; execution status remains default
404 there. A controlled local/test harness may explicitly compose
`GET /v1/workflow-executions/:execution_id` with `NewLocalStatusHandler(scope)`.
No listener, middleware, dependency or operational route is added.

### Independent visibility and observation

`application.NewLocalStatusScope(reader, visibleExecutionIDs)` accepts only an
`ExecutionReader` exposing `Get(string) (*WorkflowExecution, bool)` and a list of
IDs. Nil/typed-nil readers and invalid members fail construction. Construction
makes no repository read or policy evaluation. IDs use the same 1–128 ASCII-byte
alphabet as M4.3: letters, digits, underscore, hyphen and period. Exact bytes and
case are retained. Duplicates collapse; an empty list denies all. The constructor
copies members into a private map, retains no caller slice and exposes no setter.
Legacy Save can store IDs outside this grammar; M4.6 does not guarantee those IDs
are HTTP-addressable. The handler neither trims nor decodes path values again.
An unmatched missing segment retains Gin's default 404, outside this handler.

A distinct, private, nonzero-sized `localStatusAuthority` binds commands to one
scope. `Command(id)` validates syntax; `Querier().QueryLocal(ctx, command)`
independently checks service construction, exact capability, context, syntax and
visibility. Zero/JSON commands have no authority; a foreign scope is rejected.
M4.5 invocation authority is a different type and grants no status visibility.
Context tenant/caller/worker strings and ID possession cannot bypass the gate.

Policy precedes the sole Get. Denied IDs perform zero repository reads and reveal
no snapshot; internal `concealed` means policy denial, not proof of existence.
An allowed missing ID performs one Get and produces internal `not_found`.
Policy error/unknown decisions fail unavailable. No executable policy is publicly
injectable; fault injection is private to package tests. No service dependency
can call lifecycle mutations or actions. Query projection never normalizes or
repairs data and never returns a WorkflowExecution pointer.

This is local fixture visibility, NOT production authentication, tenant ownership,
payment proof or an in-process sandbox. Trusted host code can construct scopes and
already possesses the repository. Exposing a scoped handler operationally would
delegate that fixture capability to clients. Missing/concealed responses are equal
in status, JSON and headers; timing equality is not promised.

### Stored snapshots and legacy compatibility

A present snapshot must be nonnil, match the exact requested ID, have positive
stored Version and nonnegative Attempts, and use READY/RUNNING/COMPLETED/FAILED.
READY requires empty owner/history and zero attempts; legacy LastError is allowed.
RUNNING needs positive attempts and coherent history, permitting ownerless legacy
TryStart and retained legacy diagnostics. COMPLETED needs an owner, positive
attempts, coherent history and empty LastError. FAILED needs an owner, positive
attempts, coherent history and a nonempty LastError matching the final reason.

History permits READY -> RUNNING and RUNNING -> COMPLETED/FAILED only. Entries
must connect, use strictly increasing positive versions, have consistent owners,
and end at the snapshot's status/version. Ownerless start is supported for
RUNNING. Failure reasons are nonempty; nonfailure reasons are empty. Claim entries
cannot exceed attempts. Coherent historical suffixes are accepted, without a
version-1 origin requirement or a hardcoded version-2/3 classification. Real
Complete, Fail and SaveIfVersion transitions are supported, including higher
versions. History-free active/terminal records are ambiguous and safely rejected;
this is a deliberate compatibility limit even though legacy Save can store them.
Contradictory nonnil/false, nil/true, wrong ID and other invalid records return 500.
No repair or corruption detail is disclosed.

### Projection, errors and headers

Success returns HTTP 200 using the unchanged `ExecutionStatusResponse` with exactly
`execution_id`, `status`, `version` (the actual stored uint64 as a JSON number).
FAILED is a successful retrieval and returns 200 with status FAILED. LastError,
owner, attempts, history, event, cursor, tenant information and capabilities never
appear. Real WAIT/condition handoffs remain RUNNING, without WAITING, evaluation,
callback replay, automatic continuation or another terminal transition.

Every response emitted by this handler sets `Cache-Control: no-store`. Normal
JSON responses use `application/json; charset=utf-8`. Error messages are fixed;
causes and panic payloads are never serialized by the adapter.

| Condition | HTTP | Code | Fixed message |
|---|---:|---|---|
| Invalid matched ID/command | 400 | INVALID_REQUEST | invalid execution ID |
| Missing local capability | 401 | UNAUTHENTICATED | local capability required |
| Foreign local capability | 403 | FORBIDDEN | local status query not permitted |
| Missing or concealed | 404 | NOT_FOUND | execution not found |
| Policy error/unknown or observed cancellation | 503 | UNAVAILABLE | local status query unavailable |
| Invalid snapshot/unexpected query error | 500 | INTERNAL_ERROR | internal error |

Explicit collaborator panic containment is omitted. Application panics propagate;
Gin recovery, when composed by a trusted host, can emit a bodyless 500 before any
success JSON is written. That framework response has no JSON-envelope guarantee.
The recovery test checks no false 200, no panic payload and retained no-store.
These tests do not add recovery middleware to production or promise safe logs.

Context is checked before policy, after policy/before Get, and before returning
success. Original cancellation causes are retained internally. Synchronous Get
cannot be interrupted; a cancellation observed afterward suppresses success.
Cancellation can race response delivery after the final check. A coherent earlier
RUNNING snapshot may be returned after a legitimate concurrent Complete/Fail; no
reread, polling-until-terminal, response-lifetime lock or freshness guarantee exists.
Readers must return stable independent snapshots and support concurrency when used
concurrently. No durability, cross-process consistency, database I/O-error contract,
queue, scheduler, retries, recovery, cloud/provider/entitlement integration or
distributed exactly-once behavior is established. Human source review remains
required; this checkpoint does not accept or publish M4.6.

## M4.7 — Opt-in local request identity (awaiting human source review)

The additive `NewLocalRequestIdentityScope(repo, ids, workerID, maxReservations)`
and `NewLocalRequestIdentityHandler(scope)` provide controlled local invocation
with bounded request reservations. They register no routes or listeners.
Production `NewRouter` continues to expose only GET `/health`. M4.5's existing
scope/handler retain their fresh-execution behavior: repeated EventIDs may create
different runs there. M4.6 visibility remains a separate capability and allowlist.

### Identity, validation and admission

The new handler requires exactly one `Idempotency-Key` header value containing
1–128 ASCII letters, digits, underscores, hyphens or periods. Values are exact,
case-sensitive and never trimmed, normalized or comma-split. Missing, empty,
repeated, comma-combined and malformed values return 400 INVALID_REQUEST with
fixed message `invalid idempotency key`. Header names are case-insensitive,
including noncanonical map entries supplied by an in-process harness. A network
HTTP parser may normalize transport whitespace before this handler sees it;
validation applies to the value delivered to the handler, not raw wire bytes.

The unchanged body parser runs after a maximum 4097-byte read; its maximum valid
body remains 4096 bytes. Malformed bodies, unsupported definitions and bad identity
fields are rejected before reservation. The service independently checks its own
construction, sealed scope capability, context, key, selection syntax and approved
fixture membership before admitting a request. JSON/context values cannot mint
authority. Constructors perform no repository, resolver, runner or ID-source calls.

A request identity belongs to exactly one constructed scope. Its fingerprint is
the exact validated `(WorkflowDefinitionID, EventID)` value tuple. Equivalent JSON
formatting/order does not change that fingerprint; changing either field does.
EventID remains correlation data, not payment/tenant proof or the reservation key.

One mutex protects reservation insertion, lookup and disposition changes. At most
one request invokes the existing M4.5 service per scope/key. Repository/ID/resolver/
runner/callback calls occur outside that mutex. Separate keys can progress while
another request is blocked. ID sources must support concurrent use. The existing
repository still atomically creates by execution ID, claims once, checks ownership
and versions, and owns attempts/history/diagnostics; no lifecycle method changes.

### Reservation dispositions and responses

Reservations distinguish reserved/in-progress, known outcome and uncertain outcome.
They store only the bounded key, exact tuple, disposition and available execution
identity. They do not cache snapshots or grant status visibility.

| Situation | Response and retained behavior |
|---|---|
| First accepted invocation | Existing M4.5 outcome rules apply |
| Same key and payload, in progress or known outcome | 409 CONFLICT, `execution conflict`; no dependency calls or replay |
| Same key, different validated payload | 409 CONFLICT, same fixed message; original reservation unchanged |
| Same key/payload after uncertainty | 500 INTERNAL_ERROR, `internal error`; no dependency calls or replay |
| New key at capacity | 503 UNAVAILABLE, `local invocation unavailable`; no reservation or invocation |
| Execution-ID collision / CreateIfAbsent duplicate | 409; no adoption, overwrite, regenerated-ID retry or callbacks |
| Any creation-call error, including false/error or true/error | 500; uncertain reservation retained |
| Missing/contradictory snapshot or rejected terminal write | 500; uncertain reservation retained |
| Known pre-creation dependency error | Existing safe 503/404/500 as applicable; admitted key remains known and blocked |
| Confirmed COMPLETED, FAILED, RUNNING handoff or created-not-started READY | Existing 201 two-field creation projection |

The new boundary deliberately treats false/error creation more conservatively
than M4.5. It never assumes a creation error proves absence under a faulty
collaborator. A 503 does not authorize automatic replay: even known pre-creation
failures consume an admitted key, and identical repeats return 409.

201 FAILED requires the sealed concrete M4.5 service's accepted claim/Fail
observations and coherent authoritative FAILED/v3 snapshot, attempts/owner/history
and matching nonempty diagnostic. A FAILED snapshot with rejected Fail is uncertain
and returns 500. RUNNING handoffs retain v2; no WAITING, evaluation or continuation
is introduced. Owner, history, attempts, LastError, event, cursor, key and authority
never enter HTTP DTOs. All responses emitted by the new handler set no-store.

### Cancellation, panic and recovery limits

Cancellation observed before admission leaves no reservation. Cancellation racing
admission may consume a key. Once admitted, cancellation/errors/panics never remove
the reservation. Before creation, observed cancellation prevents creation; after
confirmed creation, existing M4.5 evidence may project 201 READY. Cancellation
during synchronous callbacks does not interrupt them or discard a coherent result.

Deferred bookkeeping marks a panic-affected reservation uncertain and lets the
panic propagate. This adapter adds no recovery middleware or safe-log guarantee.
If the underlying invocation panics before returning its outcome, the reservation
may lack the generated execution ID even if storage changed. Trusted host
inspection/manual reconciliation may be necessary. In-progress duplicates can
return 409 before an eventual uncertain outcome is known; later repeats return 500.

The positive host-selected capacity bounds retained request entries. There is no
TTL, eviction, key release, result replay, automatic retry, resume or repair.
Reservations can become stranded; exhaustion can deny new work. Do not replace a
scope or choose a new key as an automatic retry strategy for an ambiguous effect.
Scopes do not coordinate with other scopes, legacy invocations or processes, even
with one repository. Restart/reconstruction loses reservations; execution records
do not durably retain request fingerprints or definition/event associations.

Request acceptance, execution-ID insertion, single-winner claim and idempotent
external effects are separate properties. This design establishes only bounded
scope-local invocation suppression in controlled composition. It supplies no
production authentication, tenant ownership, payment/entitlement proof, operational
invocation/status endpoint, durable recovery, cloud integration or distributed
exactly-once effects. M4.10 must evaluate the roadmap with this controlled-local
delivery boundary explicitly; fixture handlers do not establish operational API
readiness. M4.7 remains in_progress pending human source review and acceptance.
