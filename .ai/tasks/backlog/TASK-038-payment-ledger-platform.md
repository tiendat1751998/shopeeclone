# TASK-038 — PAYMENT PLATFORM & LEDGER SYSTEM

## Goal

Build a REAL production-grade Payment Platform & Ledger System.

This platform is responsible for:
- wallet system
- transaction ledger (double-entry)
- payment processing orchestration
- settlement engine
- reconciliation system
- refunds / chargebacks
- financial consistency guarantees
- idempotent transaction processing
- distributed financial correctness

This is NOT a toy payment service.

The Payment Platform must be treated as:
## SYSTEM OF FINANCIAL TRUTH

The system MUST support:
- ultra-reliable transaction processing
- distributed consistency guarantees
- replay-safe financial events
- Kubernetes-native deployment
- observability-first architecture

The architecture MUST prioritize:
- financial correctness (HIGHEST PRIORITY)
- idempotency
- auditability
- consistency
- failure safety

---

## Tech Stack

Use:
- Golang (core ledger engine)
- PostgreSQL (strong consistency ledger store)
- Kafka (event-driven transactions)
- Redis Cluster (idempotency + locks)
- ClickHouse (financial analytics)
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Debezium (CDC)
- Temporal (workflow orchestration)
- CockroachDB (if distributed ledger extension needed)

---

## Core Responsibilities

The Payment Platform MUST support:

### Wallet System
- user wallet balances
- frozen / available balances
- multi-currency support (optional extension)
- balance locking mechanisms

### Double-Entry Ledger
- immutable ledger entries
- debit/credit accounting model
- audit-safe financial history
- append-only transaction log

### Payment Processing
- payment authorization
- payment capture
- payment reversal
- idempotent payment execution

### Settlement Engine
- batch settlement processing
- merchant payout calculation
- settlement reconciliation
- delayed settlement windows

### Reconciliation System
- detect inconsistencies
- fix mismatched ledger entries
- audit trails
- replay-safe reconciliation

### Refunds & Chargebacks
- refund lifecycle management
- partial/full refunds
- chargeback tracking
- dispute handling workflows

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- enforce strict separation: API / domain / ledger / persistence
- support distributed deployments
- support event-driven financial workflows

The Payment Platform MUST:
- guarantee idempotent financial operations
- support replay-safe ledger updates
- ensure immutable financial records
- prevent double spending under ALL conditions

Use:
- CQRS where appropriate
- strong domain-driven design (DDD)
- dependency injection
- strict transactional boundaries

The system MUST tolerate:
- retry storms (WITHOUT double charging)
- duplicate payment requests
- Kafka replays
- partial database failures
- network partitions
- service restarts mid-transaction

---

## Folder Structure

Generate:

platforms/payment-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── ledger/
│   ├── wallet/
│   ├── payments/
│   ├── settlement/
│   ├── reconciliation/
│   ├── refunds/
│   ├── accounting/
│   ├── idempotency/
│   ├── orchestration/
│   ├── events/
│   ├── cache/
│   ├── locks/
│   ├── persistence/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── db/
│   ├── migrations/
│   ├── schema.sql
│   └── ledger_models.sql
│
├── workflows/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Ledger System Requirements (CRITICAL)

The ledger MUST be:

- append-only
- immutable
- double-entry accounting based
- auditable at all times

Rules:
- every transaction has debit AND credit
- balance is derived, not overwritten
- NO UPDATE on ledger rows
- NO DELETE on ledger rows

Generate:
- ledger schema
- posting engine
- balance calculation strategy
- audit trail system

---

## Idempotency Requirements (CRITICAL)

All payment operations MUST be idempotent.

Generate:
- idempotency keys
- request deduplication
- retry-safe transaction execution
- Redis-based idempotency guard
- database uniqueness constraints

System MUST prevent:
- double charge
- double refund
- duplicate settlement

---

## Transaction Safety Requirements

The system MUST guarantee:
- atomic ledger writes
- safe rollback strategies
- compensation transactions (NOT naive rollback)
- eventual consistency handling for external systems

Use:
- database transactions (STRICT)
- outbox pattern
- saga pattern for distributed workflows

---

## Settlement Requirements

Support:
- merchant payout calculation
- scheduled settlement windows
- batch processing engine
- reconciliation before payout

Generate:
- settlement pipelines
- batch workers
- replay-safe settlement execution

---

## Reconciliation Requirements

Support:
- ledger vs wallet reconciliation
- missing transaction detection
- mismatch correction via compensation entries

Generate:
- reconciliation jobs
- audit comparison tools
- anomaly detection workflows

---

## Event-Driven Requirements

Generate events for:
- payment initiated
- payment completed
- payment failed
- refund processed
- settlement completed

Use Kafka with:

- idempotent producers
- replay-safe consumers
- DLQ handling
- event versioning

---

## PostgreSQL Requirements (CORE OF TRUTH)

Use PostgreSQL as:
## SINGLE SOURCE OF TRUTH

Generate:
- normalized ledger schema
- strict constraints
- indexes for audit queries
- transactional integrity rules

NEVER:
- use eventual consistency for balances
- allow untracked updates
- bypass transactions

---

## Redis Requirements

Use Redis for:
- idempotency keys
- distributed locks
- short-lived transaction state

Must:
- expire safely
- avoid lock race conditions
- support failover correctness

---

## ClickHouse Requirements

Use ClickHouse for:
- financial analytics
- revenue reporting
- fraud correlation (read-only integration)
- settlement analytics

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto definitions

Endpoints:
- /payments/initiate
- /payments/capture
- /payments/refund
- /wallet/balance
- /ledger/transactions
- /settlement/run

Must support:
- strict idempotency
- audit traceability
- deterministic behavior

---

## Security Requirements

The platform MUST:
- enforce strict RBAC for financial APIs
- validate all transaction inputs
- prevent unauthorized balance manipulation
- isolate financial operations

Never:
- expose internal ledger mutations
- trust external payment status blindly
- bypass accounting rules

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logs
- financial audit logs

Metrics:
- transaction success rate
- ledger write latency
- reconciliation mismatch rate
- settlement delay
- idempotency hits

Logs:
- MUST be immutable audit logs
- include correlation IDs

Never log sensitive financial secrets.

---

## Reliability Requirements

Implement:
- retries (SAFE ONLY with idempotency)
- timeout handling
- circuit breakers
- graceful shutdown
- failure recovery workflows

Critical rule:
## NEVER RETRY WITHOUT IDEMPOTENCY GUARANTEE

---

## Kubernetes Requirements

Generate:
- Deployments
- StatefulSets (for ledger workers if needed)
- Services
- HPA
- PDB
- ConfigMaps
- Secrets
- Helm charts

Must support:
- zero-downtime deployments
- safe rolling updates
- strict resource isolation

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone pipelines
- migration validation
- ledger integrity tests
- chaos testing
- security scanning
- GitOps deployment

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- ledger consistency tests
- idempotency tests
- reconciliation tests
- concurrency tests
- failure injection tests

Test scenarios:
- duplicate payment requests
- retry storms
- partial system failure
- DB crash mid-transaction
- Kafka replay duplication
- settlement mismatch

---

## Output Requirements

Explain:
- payment architecture
- ledger consistency model
- idempotency strategy
- settlement design
- reconciliation strategy
- failure recovery strategy
- scaling strategy

Generate production-grade code only.

No toy wallet system.
No fake payment flows.
No unsafe financial shortcuts.

---

## Acceptance Criteria

The Payment Platform must support future integration with:
- Fraud Platform
- Recommendation Platform (signals only)
- Analytics Platform
- Search Platform (analytics only)
- AI/ML Platform (risk signals)

without compromising financial integrity.

The system MUST survive:
- retry storms
- duplicate events
- partial failures
- distributed deployments
- database crashes
- Kafka replays

WITHOUT losing financial correctness.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
Strict financial correctness required.