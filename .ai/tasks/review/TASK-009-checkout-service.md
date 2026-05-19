# TASK-009 — CHECKOUT SERVICE

## Goal

Build a REAL production-grade Checkout Service.

This service is responsible for:
- checkout orchestration
- pricing freeze
- inventory reservation orchestration
- promotion finalization
- checkout validation
- checkout snapshot generation
- anti-double-submit protection
- idempotency orchestration
- pre-order workflow

This is NOT a toy CRUD checkout service.

The Checkout Service must support:
- massive concurrency
- flash-sale traffic
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- eventual consistency

The architecture MUST prioritize:
- checkout correctness
- resiliency
- idempotency
- anti-duplicate-order protection
- operational stability

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- Redis Cluster
- MySQL
- Kafka or NATS JetStream
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Saga orchestration helpers

---

## Core Responsibilities

The Checkout Service MUST support:

### Checkout Orchestration
- checkout validation
- seller grouping
- pricing finalization
- promotion finalization
- inventory reservation orchestration

### Pricing Freeze
- immutable pricing snapshot
- immutable promotion snapshot
- checkout pricing integrity
- anti-price-drift handling

### Reservation Orchestration
- inventory reservation requests
- reservation rollback
- reservation timeout handling
- reservation reconciliation

### Checkout Snapshot
- immutable checkout snapshot
- replay-safe checkout state
- retry-safe workflow state

### Idempotency Protection
- anti-double-submit
- idempotent checkout requests
- replay-safe processing
- duplicate request prevention

### Validation Engine
- cart validation
- promotion validation
- inventory validation
- seller validation
- shipping validation hooks

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Checkout Service MUST:
- support orchestration workflows
- support distributed retries
- support failure recovery
- support saga-like coordination

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The checkout system MUST tolerate:
- retry storms
- partial failures
- duplicate requests
- stale cache
- distributed deployments

---

## Folder Structure

Generate:

services/checkout/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── checkout/
│   ├── pricing/
│   ├── reservation/
│   ├── validation/
│   ├── orchestration/
│   ├── snapshot/
│   ├── idempotency/
│   ├── reconciliation/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── migrations/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Database Requirements

Use MySQL for:
- checkout snapshots
- orchestration state
- idempotency state
- retry state
- reconciliation state

Generate:
- optimized schemas
- indexes
- constraints
- immutable snapshot structures

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- audit timestamps

Never:
- mutate finalized checkout snapshots
- ignore replay safety
- ignore idempotency
- tightly couple order creation logic

---

## Redis Requirements

Use Redis for:
- idempotency cache
- checkout session cache
- anti-double-submit protection
- distributed coordination
- temporary workflow state

Generate:
- TTL strategy
- distributed coordination strategy
- retry handling
- stale workflow cleanup

Support:
- distributed deployments
- retry storms
- duplicate submissions

The coordination layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- checkout started
- checkout validated
- pricing frozen
- inventory reserved
- reservation failed
- checkout finalized
- checkout expired

Use:
- Kafka or NATS JetStream

Requirements:
- retries
- DLQ
- idempotent consumers
- replay-safe processing
- event versioning
- consumer groups

Support:
- eventual consistency
- distributed orchestration
- async reconciliation

No fake orchestration.

---

## Pricing Freeze Requirements

Support:
- immutable pricing snapshots
- immutable promotion snapshots
- anti-price-drift handling
- checkout replay safety

Generate:
- pricing freeze workflow
- pricing integrity validation
- stale pricing reconciliation

The service MUST tolerate:
- async promotion updates
- pricing recalculation delays
- distributed cache inconsistency

---

## Inventory Reservation Requirements

Support:
- reservation orchestration
- reservation rollback
- reservation expiration
- distributed retries
- failure recovery

Generate:
- reservation coordination workflow
- timeout handling
- compensation workflow
- reconciliation jobs

The service MUST NOT own inventory truth.

Inventory truth belongs to Inventory Service.

---

## Idempotency Requirements

Implement:
- anti-double-submit protection
- request fingerprinting
- replay-safe workflows
- distributed idempotency validation

Support:
- retries
- duplicate requests
- network retries
- client resubmission

Generate:
- idempotency middleware
- replay validation
- duplicate workflow handling

---

## Validation Requirements

Support:
- cart validation
- promotion validation
- inventory validation
- seller validation
- shipping validation hooks

Generate:
- distributed validation flow
- async reconciliation support
- failure recovery strategy

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /checkout
- /checkout/preview
- /checkout/validate
- /checkout/finalize
- /checkout/status

Support:
- pagination
- filtering
- validation
- idempotency keys

---

## Security Requirements

The service MUST:
- validate ownership
- enforce RBAC
- validate checkout integrity
- sanitize input
- prevent duplicate submissions

Never:
- trust client pricing
- trust client promotion totals
- trust client inventory state
- expose orchestration internals

Generate:
- authorization middleware
- idempotency validation
- workflow integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- checkout latency
- reservation latency
- pricing freeze latency
- retry count
- idempotency rejection count
- orchestration failure count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive checkout/payment metadata.

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
- compensation workflows
- failure isolation

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
- orchestration tests
- replay tests
- idempotency tests
- reservation rollback tests
- compensation tests
- concurrency tests

Test:
- duplicate submissions
- retry storms
- partial failures
- stale pricing
- reservation rollback
- distributed reconciliation
- high concurrency

---

## Output Requirements

Explain:
- checkout architecture
- orchestration strategy
- pricing freeze workflow
- reservation workflow
- idempotency strategy
- replay protection strategy
- compensation strategy
- reconciliation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy checkout CRUD service.
No fake orchestration.
No naive retry handling.

---

## Acceptance Criteria

The Checkout Service must support future integration with:
- Cart Service
- Inventory Service
- Promotion Service
- Order Service
- Payment Service

without major future refactors.

The service MUST realistically tolerate:
- flash-sale traffic
- duplicate submissions
- retry storms
- partial failures
- distributed deployments

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