# TASK-010 — ORDER SERVICE

## Goal

Build a REAL production-grade Order Service.

This service is responsible for:
- order creation
- order lifecycle management
- order state machine
- seller order splitting
- post-checkout workflows
- order reconciliation
- cancellation workflows
- shipment orchestration hooks
- distributed order processing

This is NOT a toy CRUD order service.

The Order Service must support:
- millions of orders
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- eventual consistency
- replay-safe workflows

The architecture MUST prioritize:
- order correctness
- lifecycle integrity
- resiliency
- replay safety
- operational stability

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- MySQL
- Redis Cluster
- Kafka or NATS JetStream
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Saga orchestration helpers

---

## Core Responsibilities

The Order Service MUST support:

### Order Creation
- order creation from checkout snapshot
- replay-safe order creation
- seller order splitting
- immutable order snapshot
- order numbering

### Order Lifecycle
- pending
- awaiting payment
- paid
- processing
- packed
- shipped
- delivered
- completed
- cancelled
- refunded

### Order State Machine
- valid state transitions
- transition validation
- transition audit logs
- replay-safe transitions

### Seller Order Splitting
- multi-seller order splitting
- seller-level fulfillment
- seller shipment coordination

### Cancellation Workflow
- user cancellation
- seller cancellation
- timeout cancellation
- compensation workflows

### Reconciliation
- order reconciliation
- payment reconciliation hooks
- inventory reconciliation hooks
- shipment reconciliation hooks

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Order Service MUST:
- support replay-safe workflows
- support distributed retries
- support failure recovery
- support lifecycle orchestration

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The order system MUST tolerate:
- retry storms
- duplicate events
- partial failures
- distributed deployments
- stale cache

---

## Folder Structure

Generate:

services/order/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── order/
│   ├── lifecycle/
│   ├── state_machine/
│   ├── reconciliation/
│   ├── cancellation/
│   ├── shipment/
│   ├── fulfillment/
│   ├── numbering/
│   ├── snapshot/
│   ├── idempotency/
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
- orders
- order items
- order snapshots
- lifecycle history
- transition logs
- reconciliation state
- cancellation state

Generate:
- optimized schemas
- indexes
- constraints
- immutable snapshots
- audit tables

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- lifecycle consistency

Never:
- mutate immutable snapshots
- ignore transition integrity
- ignore reconciliation
- tightly couple payment logic

---

## Redis Requirements

Use Redis for:
- hot order cache
- idempotency cache
- workflow coordination
- temporary orchestration state

Generate:
- TTL strategy
- invalidation strategy
- distributed coordination
- retry handling

Support:
- distributed deployments
- replay protection
- retry storms

The cache layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- order created
- order paid
- order processing
- order shipped
- order delivered
- order cancelled
- order refunded
- reconciliation triggered

Use:
- Kafka or NATS JetStream

Requirements:
- retries
- DLQ
- replay-safe consumers
- idempotent processing
- event versioning
- consumer groups

Support:
- eventual consistency
- async workflows
- distributed reconciliation

No fake async architecture.

---

## Order State Machine Requirements

Implement:
- strict state transitions
- transition validation
- invalid transition prevention
- replay-safe transitions
- transition audit logs

Generate:
- lifecycle transition rules
- transition middleware
- reconciliation handling

The lifecycle system MUST be production-grade.

---

## Seller Splitting Requirements

Support:
- multi-seller orders
- seller-level fulfillment
- shipment grouping
- seller cancellation isolation

Generate:
- splitting strategy
- shipment coordination hooks
- reconciliation flow

---

## Cancellation Requirements

Support:
- user cancellation
- seller cancellation
- payment timeout cancellation
- compensation workflows
- reservation rollback hooks

Generate:
- cancellation workflow
- distributed rollback handling
- reconciliation jobs

The service MUST tolerate:
- partial cancellation
- async rollback
- duplicate cancellation requests

---

## Reconciliation Requirements

Support:
- payment reconciliation hooks
- inventory reconciliation hooks
- shipment reconciliation hooks
- async reconciliation jobs

Generate:
- reconciliation workflows
- failure recovery strategy
- retry-safe reconciliation

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /orders
- /orders/{id}
- /orders/status
- /orders/cancel
- /orders/history

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
- validate lifecycle transitions
- sanitize input
- prevent unauthorized transitions

Never:
- trust client lifecycle state
- expose internal reconciliation metadata
- expose internal workflow state

Generate:
- authorization middleware
- transition validation
- ownership validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- order creation latency
- lifecycle transition latency
- reconciliation latency
- cancellation latency
- retry count
- reconciliation failures

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive payment metadata.

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
- compensation workflows

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
- lifecycle tests
- replay tests
- reconciliation tests
- cancellation tests
- concurrency tests

Test:
- duplicate events
- retry storms
- partial failures
- invalid transitions
- replay safety
- distributed reconciliation
- high concurrency

---

## Output Requirements

Explain:
- order architecture
- lifecycle state machine
- seller splitting strategy
- reconciliation strategy
- cancellation workflow
- replay protection strategy
- event-driven flow
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy order CRUD service.
No fake state machine.
No fake async orchestration.

---

## Acceptance Criteria

The Order Service must support future integration with:
- Checkout Service
- Payment Service
- Shipment Service
- Notification Service
- Refund Service

without major future refactors.

The service MUST realistically tolerate:
- retry storms
- duplicate events
- partial failures
- distributed deployments
- high concurrency

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