# TASK-022 — BILLING & FINANCE PLATFORM

## Goal

Build a REAL production-grade Billing & Finance Platform.

This platform is responsible for:
- double-entry ledger
- wallet systems
- settlements
- payouts
- reconciliation
- financial audit trails
- merchant balances
- refunds
- fee accounting
- accounting synchronization

This is NOT a toy payment service.

The Billing & Finance Platform must support:
- strict financial correctness
- replay-safe transactions
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- audit-grade integrity

The architecture MUST prioritize:
- financial correctness
- replay safety
- reconciliation correctness
- auditability
- operational stability

---

## Tech Stack

Use:
- Golang
- PostgreSQL
- Kafka
- Redis Cluster
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Apache Flink
- Temporal
- ClickHouse for analytics

---

## Core Responsibilities

The Billing & Finance Platform MUST support:

### Double-Entry Ledger
- immutable ledger entries
- debit/credit accounting
- financial balancing
- audit-safe bookkeeping

### Wallet Systems
- buyer wallets
- seller wallets
- platform wallets
- frozen balances
- pending balances

### Settlements
- merchant settlements
- fee accounting
- payout scheduling
- settlement reconciliation

### Refunds
- refund orchestration
- partial refunds
- replay-safe refunds
- refund reconciliation

### Reconciliation
- ledger reconciliation
- payout reconciliation
- external payment reconciliation
- replay-safe reconciliation

### Financial Audit Trails
- immutable audit logs
- compliance hooks
- evidence retention
- transaction traceability

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate ledger/wallets/reconciliation
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Billing & Finance Platform MUST:
- support replay-safe transactions
- support distributed reconciliation
- support audit-grade integrity
- support idempotent money movement

Use:
- Saga orchestration where appropriate
- CQRS where appropriate
- dependency injection
- resilience patterns

The finance system MUST tolerate:
- retry storms
- duplicate events
- delayed settlement events
- partial failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/billing-finance/
├── cmd/
├── internal/
│   ├── config/
│   ├── ledger/
│   ├── wallets/
│   ├── settlements/
│   ├── payouts/
│   ├── refunds/
│   ├── reconciliation/
│   ├── accounting/
│   ├── audits/
│   ├── synchronization/
│   ├── replay/
│   ├── idempotency/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## PostgreSQL Requirements

Use PostgreSQL for:
- ledger entries
- wallet balances
- settlement state
- payout state
- reconciliation state
- audit trails

Generate:
- optimized schemas
- indexes
- immutable ledger tables
- append-only accounting history

Requirements:
- ACID correctness
- replay safety
- strict transactional guarantees
- audit-grade traceability

Never:
- mutate immutable ledger history
- use floating-point money types
- bypass ledger balancing
- tightly couple external payments to ledger writes

---

## Ledger Requirements

Support:
- double-entry accounting
- immutable entries
- balancing validation
- replay-safe posting

Generate:
- posting engine
- ledger balancing checks
- replay-safe transaction orchestration
- reconciliation hooks

The ledger system MUST tolerate:
- duplicate events
- replay storms
- partial posting failures

No fake accounting architecture.

---

## Wallet Requirements

Support:
- buyer wallets
- seller wallets
- frozen balances
- pending balances
- withdrawal holds

Generate:
- wallet orchestration
- balance reservation flows
- replay-safe wallet updates
- reconciliation workflows

The wallet system MUST be production-grade.

---

## Settlement Requirements

Support:
- merchant settlements
- fee accounting
- payout scheduling
- delayed settlements

Generate:
- settlement pipelines
- replay-safe settlement handling
- payout orchestration
- settlement reconciliation

The settlement system MUST tolerate:
- delayed payout events
- duplicate settlement events
- retry storms

---

## Refund Requirements

Support:
- full refunds
- partial refunds
- replay-safe refunds
- asynchronous refund workflows

Generate:
- refund orchestration
- reconciliation hooks
- refund audit trails
- compensation workflows

No naive refund architecture.

---

## Reconciliation Requirements

Support:
- internal reconciliation
- external reconciliation
- payout reconciliation
- replay-safe reconciliation

Generate:
- reconciliation jobs
- replay-safe aggregation
- reconciliation reports
- mismatch investigation workflows

The reconciliation system MUST be realistic.

---

## Kafka Requirements

Use Kafka for:
- ledger events
- settlement events
- payout events
- replay-safe synchronization

Generate:
- topic strategy
- partitioning strategy
- replay pipelines
- DLQ topics

Requirements:
- replay-safe ingestion
- idempotent producers
- consumer groups
- event versioning

Support:
- financial event streaming
- distributed reconciliation
- async synchronization

---

## Redis Requirements

Use Redis for:
- idempotency coordination
- replay coordination
- hot balance cache
- distributed throttling

Generate:
- TTL strategy
- distributed coordination
- replay protection
- rate limiting

Support:
- high concurrency
- distributed deployments
- low-latency balance reads

---

## Audit Requirements

Support:
- immutable audit logs
- transaction traceability
- evidence retention
- compliance hooks

Generate:
- audit pipelines
- replay-safe audit aggregation
- integrity verification workflows

The audit system MUST be audit-grade.

---

## Event-Driven Requirements

Generate events for:
- ledger posted
- wallet updated
- settlement completed
- payout initiated
- refund processed
- reconciliation mismatch detected

Use:
- Kafka

Requirements:
- retries
- DLQ
- replay-safe consumers
- idempotent processing
- event versioning
- consumer groups

Support:
- eventual consistency
- replay-safe synchronization
- distributed reconciliation

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /wallets
- /ledger
- /settlements
- /refunds
- /payouts

Support:
- pagination
- filtering
- idempotency keys
- replay-safe transaction handling

---

## Security Requirements

The platform MUST:
- validate financial ownership
- enforce RBAC
- sanitize input
- isolate financial data
- protect accounting integrity

Never:
- expose internal ledger balancing internals
- expose sensitive financial audit evidence
- trust external settlement states blindly
- allow direct balance mutation bypassing ledger

Generate:
- authorization middleware
- replay validation
- accounting integrity validation
- financial isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- ledger posting latency
- reconciliation latency
- payout latency
- refund latency
- replay latency
- mismatch detection count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive financial secrets.

---

## Reliability Requirements

Implement:
- retries
- timeout handling
- graceful shutdown
- panic recovery
- circuit breakers
- backoff strategies

Support:
- rolling deployment
- autoscaling
- distributed deployments

Generate:
- resilience middleware
- retry policies
- failure isolation
- replay recovery workflows

---

## Kubernetes Requirements

Generate:
- Deployment
- Service
- ConfigMap
- Secret integration
- HPA
- PodDisruptionBudget
- ServiceMonitor
- NetworkPolicy
- Helm chart

Support:
- readiness/liveness probes
- autoscaling
- rolling deployment
- canary deployment

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone pipeline
- linting
- testing
- vulnerability scanning
- Docker build
- Helm validation

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- replay tests
- reconciliation tests
- ledger balancing tests
- refund tests
- concurrency tests

Test:
- duplicate events
- replay storms
- delayed settlement events
- ledger consistency
- reconciliation correctness
- high concurrency
- distributed deployments

---

## Output Requirements

Explain:
- ledger architecture
- wallet strategy
- reconciliation strategy
- settlement strategy
- replay-safe accounting strategy
- audit strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy payment service.
No fake ledger architecture.
No naive balance mutation logic.

---

## Acceptance Criteria

The Billing & Finance Platform must support future integration with:
- Payment Platform
- Order Service
- Advertising Platform
- Fraud Platform
- BI Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate events
- delayed settlement events
- distributed deployments
- strict financial correctness requirements

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
If a real production system would require operational complexity,
YOU MUST model that complexity realistically instead of simplifying it away.