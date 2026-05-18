# TASK-011 — PAYMENT SERVICE

## Goal

Build a REAL production-grade Payment Service.

This service is responsible for:
- payment orchestration
- PSP integration
- payment authorization
- payment capture
- webhook validation
- reconciliation
- anti-double-charge protection
- distributed payment consistency
- refund orchestration hooks
- fraud prevention hooks

This is NOT a toy CRUD payment service.

The Payment Service must support:
- massive concurrency
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- replay-safe workflows
- financial auditability

The architecture MUST prioritize:
- financial correctness
- anti-double-charge protection
- replay safety
- resiliency
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
- PSP SDK integrations
- HSM/KMS integration hooks

---

## Core Responsibilities

The Payment Service MUST support:

### Payment Orchestration
- payment authorization
- payment capture
- payment cancellation
- payment expiration
- distributed payment coordination

### PSP Integration
- multiple PSP providers
- PSP failover hooks
- webhook validation
- PSP retry handling
- provider abstraction

### Anti-Double-Charge
- idempotent payment requests
- duplicate prevention
- replay-safe workflows
- distributed coordination

### Reconciliation
- PSP reconciliation
- order reconciliation
- settlement reconciliation
- async reconciliation jobs

### Refund Hooks
- refund orchestration hooks
- partial refund support
- refund reconciliation hooks

### Fraud Prevention Hooks
- suspicious payment detection hooks
- risk scoring hooks
- abuse detection hooks

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Payment Service MUST:
- support replay-safe workflows
- support distributed retries
- support failure recovery
- support PSP abstraction

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The payment system MUST tolerate:
- duplicate requests
- retry storms
- webhook replay
- partial failures
- PSP instability
- distributed deployments

---

## Folder Structure

Generate:

services/payment/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── payment/
│   ├── authorization/
│   ├── capture/
│   ├── reconciliation/
│   ├── webhooks/
│   ├── providers/
│   ├── idempotency/
│   ├── fraud/
│   ├── refunds/
│   ├── settlement/
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
- payment records
- authorization state
- capture state
- reconciliation state
- webhook history
- refund state
- settlement state

Generate:
- optimized schemas
- indexes
- immutable financial records
- audit tables

Requirements:
- transactional correctness
- replay safety
- financial integrity
- pagination everywhere

Never:
- mutate finalized financial records
- ignore reconciliation
- ignore idempotency
- trust PSP webhook blindly

---

## Redis Requirements

Use Redis for:
- idempotency cache
- anti-double-charge coordination
- webhook replay protection
- temporary orchestration state

Generate:
- TTL strategy
- distributed coordination strategy
- replay protection
- retry handling

Support:
- distributed deployments
- duplicate prevention
- webhook replay protection

The coordination layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- payment authorized
- payment captured
- payment failed
- payment expired
- payment refunded
- reconciliation triggered
- webhook received

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
- distributed workflows
- async reconciliation

No fake async architecture.

---

## PSP Integration Requirements

Support:
- multiple PSP providers
- PSP abstraction
- provider failover hooks
- webhook verification
- retry-safe PSP handling

Generate:
- PSP adapter interfaces
- provider abstraction layer
- webhook signature validation
- PSP retry workflow

The PSP architecture MUST be production-grade.

No fake PSP abstraction.

---

## Webhook Requirements

Implement:
- webhook signature validation
- replay protection
- duplicate webhook handling
- retry-safe processing
- webhook audit logs

Support:
- PSP retries
- webhook reordering
- duplicate delivery

Generate:
- webhook validation middleware
- replay-safe consumers
- webhook reconciliation jobs

---

## Anti-Double-Charge Requirements

Implement:
- distributed idempotency
- payment fingerprinting
- replay-safe authorization
- duplicate prevention

Support:
- client retries
- PSP retries
- webhook replay
- network retries

Generate:
- idempotency middleware
- replay validation
- duplicate coordination workflow

---

## Reconciliation Requirements

Support:
- PSP reconciliation
- order reconciliation
- settlement reconciliation
- async reconciliation jobs

Generate:
- reconciliation workflows
- failure recovery strategy
- retry-safe reconciliation

The reconciliation layer MUST be realistic.

---

## Fraud Prevention Requirements

Support:
- suspicious payment hooks
- abnormal retry detection
- replay detection
- risk scoring hooks

Generate:
- fraud event hooks
- async fraud pipeline hooks
- replay detection logic

The service MUST support future fraud systems.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /payments
- /payments/authorize
- /payments/capture
- /payments/refund
- /payments/status
- /webhooks

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
- validate webhook signatures
- sanitize input
- prevent replay attacks

Never:
- trust client payment state
- trust PSP callbacks blindly
- expose internal financial metadata
- expose secret provider credentials

Generate:
- authorization middleware
- webhook verification middleware
- replay validation
- financial integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- payment authorization latency
- payment capture latency
- webhook latency
- reconciliation latency
- duplicate prevention count
- PSP failure rate
- replay attack detection count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log:
- card data
- secret keys
- sensitive PSP credentials

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
- webhook tests
- replay tests
- reconciliation tests
- PSP failover tests
- concurrency tests

Test:
- duplicate requests
- webhook replay
- PSP retries
- retry storms
- partial failures
- reconciliation correctness
- high concurrency

---

## Output Requirements

Explain:
- payment architecture
- PSP abstraction strategy
- webhook validation flow
- anti-double-charge strategy
- replay protection strategy
- reconciliation strategy
- fraud prevention hooks
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy payment CRUD service.
No fake PSP abstraction.
No naive reconciliation.

---

## Acceptance Criteria

The Payment Service must support future integration with:
- Order Service
- Refund Service
- Fraud Service
- Settlement Service
- Notification Service

without major future refactors.

The service MUST realistically tolerate:
- PSP instability
- webhook replay
- duplicate requests
- retry storms
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