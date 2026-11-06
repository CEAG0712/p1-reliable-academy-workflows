
# Reliable Academy Workflows
### Reusable Workflow Automation Platform — P1

A Go-based backend engineering portfolio project focused on
reliable workflow execution, concurrency control, execution
ownership, and distributed systems architecture.

## Project Overview

This project explores how to build a reusable workflow automation
engine capable of executing business workflows reliably.

The first business use case is student enrollment and
course-access provisioning.

The long-term goal is a generic platform supporting multiple
triggers, conditions, actions, delays, retries, and integrations.

## Current Implementation

- Enrollment eligibility validation
- Tenant-aware entitlement identity
- In-memory enrollment provider
- Idempotent local entitlement grants
- Concurrency-safe workflow state transitions
- Owner-aware workflow claiming
- Competing-worker concurrency tests
- In-memory workflow execution repository
- Independent repository snapshots
- Repository synchronization with RWMutex
- Go unit tests and race-detector validation

## Technology Stack

### Implemented
- Go
- Go standard library testing
- Goroutines and channels
- sync.Mutex and sync.RWMutex
- sync/atomic
- In-memory repositories

### Planned
- Gin
- AWS Lambda
- DynamoDB
- Amazon SQS and SNS
- Terraform
- LangChain
- Structured logging, metrics and traces
- CI/CD and load testing

## Testing

Run the test suite:

```bash
go test -count=1 -v ./...
```

Run the race detector:

```bash
go test -race -count=1 -v ./...
```

## Engineering Focus

1. Correct workflow state transitions
2. Single-winner worker claims
3. Tenant-aware business operations
4. Safe handling of duplicate events
5. Isolation of persisted execution state
6. Durable conditional updates
7. Recovery from partial failures

## Current Development Status

Stage 3 — Exercise 60 completed.

The current implementation is a local engineering prototype.

Cloud infrastructure, complete orchestration, and distributed
recovery are not yet implemented.

## Next Milestone

Implement conditional claims at the repository level,
then progress toward DynamoDB-backed execution ownership.

## Business Objectives

The original project proposes the following targets:

- Automate 10,000 workflow tasks per month
- Prevent 2,500 manual errors annually
- Reduce operational costs by $200,000 annually

These targets are not measured production results.

## Project Purpose

Independent backend engineering portfolio project
demonstrating progressive implementation of reliable
workflow automation concepts.
