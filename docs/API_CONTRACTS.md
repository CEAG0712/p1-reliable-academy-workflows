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
