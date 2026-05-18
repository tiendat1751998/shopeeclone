# TASK-013 — INVENTORY SERVICE

## Goal

Build a REAL production-grade Inventory Service.

This service is responsible for:
- stock management
- distributed reservations
- warehouse stock coordination
- oversell prevention
- reservation expiration
- stock reconciliation
- flash-sale inventory protection
- inventory synchronization

This is NOT a toy CRUD inventory service.

The Inventory Service must support:
- millions of SKUs
- massive flash-sale traffic
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- replay-safe workflows

The architecture MUST prioritize:
- stock correctness
- oversell prevention
- reservation integrity
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
- Redis Lua scripts
- warehouse allocation engine

---

## Core Responsibilities

The Inventory Service MUST support:

### Stock Management
- stock allocation
- stock deduction
- stock replenishment
- warehouse stock tracking
- SKU stock tracking

### Reservation System
- inventory reservation
- reservation expiration
- reservation rollback
- distributed reservation coordination

### Oversell Prevention
- atomic reservation handling
- distributed concurrency control
- replay-safe reservation logic
- flash-sale protection

### Warehouse Coordination
- multi-warehouse stock
- warehouse allocation
- regional inventory
- warehouse synchronization

### Reconciliation
- stock reconciliation
- reservation reconciliation
- warehouse reconciliation
- async reconciliation jobs

### Flash Sale Protection
- burst traffic protection
- stock throttling
- anti-oversell logic
- reservation queue protection

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Inventory Service MUST:
- support replay-safe workflows
- support distributed retries
- support failure recovery
- support reservation coordination

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The inventory system MUST tolerate:
- retry storms
- duplicate events
- reservation races
- partial failures
- distributed deployments
- stale cache

---

## Folder Structure

Generate:

services/inventory/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── stock/
│   ├── reservations/
│   ├── warehouses/
│   ├── allocation/
│   ├── reconciliation/
│   ├── flashsale/
│   ├── synchronization/
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
- stock records
- warehouse stock
- reservation state
- reconciliation state
- allocation mappings
- stock history

Generate:
- optimized schemas
- indexes
- immutable stock history
- audit tables

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- stock consistency

Never:
- mutate immutable stock history
- ignore reservation races
- ignore reconciliation
- use naive decrement logic

---

## Redis Requirements

Use Redis for:
- hot inventory cache
- reservation coordination
- flash-sale protection
- oversell prevention
- distributed locking

Generate:
- TTL strategy
- atomic reservation operations
- replay protection
- distributed coordination
- retry handling

Support:
- distributed deployments
- flash-sale traffic
- high concurrency

The coordination layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- stock reserved
- stock released
- stock deducted
- stock replenished
- reservation expired
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
- async synchronization
- distributed reconciliation

No fake async architecture.

---

## Reservation Requirements

Implement:
- reservation coordination
- reservation expiration
- replay-safe reservations
- distributed reservation locking
- duplicate reservation prevention

Support:
- retry storms
- duplicate requests
- concurrent reservations
- network retries

Generate:
- reservation workflows
- expiration jobs
- reconciliation workflows
- rollback handling

The reservation system MUST be production-grade.

---

## Oversell Prevention Requirements

Support:
- atomic stock reservations
- distributed concurrency control
- flash-sale protection
- anti-oversell guarantees

Generate:
- reservation throttling
- queue protection
- distributed coordination strategy
- reconciliation fallback

No naive stock decrement logic.

---

## Warehouse Coordination Requirements

Support:
- multi-warehouse allocation
- regional stock management
- warehouse synchronization
- warehouse failover hooks

Generate:
- allocation strategy
- synchronization workflows
- reconciliation jobs

---

## Reconciliation Requirements

Support:
- stock reconciliation
- reservation reconciliation
- warehouse reconciliation
- async reconciliation jobs

Generate:
- reconciliation workflows
- failure recovery strategy
- retry-safe reconciliation

The reconciliation layer MUST be realistic.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /inventory
- /inventory/reserve
- /inventory/release
- /inventory/deduct
- /inventory/warehouses
- /inventory/status

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
- sanitize input
- prevent reservation abuse
- validate stock operations

Never:
- trust external stock state blindly
- expose internal coordination metadata
- expose warehouse secrets

Generate:
- authorization middleware
- reservation validation
- replay validation
- stock integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- reservation latency
- stock deduction latency
- reservation failure count
- oversell prevention count
- flash-sale throttling count
- reconciliation failures

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive warehouse credentials.

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
- reservation tests
- replay tests
- reconciliation tests
- flash-sale tests
- concurrency tests

Test:
- reservation races
- duplicate requests
- retry storms
- oversell scenarios
- flash-sale spikes
- reconciliation correctness
- high concurrency

---

## Output Requirements

Explain:
- inventory architecture
- reservation coordination strategy
- oversell prevention strategy
- warehouse allocation strategy
- replay protection strategy
- reconciliation strategy
- flash-sale protection strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy inventory CRUD service.
No fake reservation logic.
No naive stock decrement architecture.

---

## Acceptance Criteria

The Inventory Service must support future integration with:
- Checkout Service
- Order Service
- Warehouse Service
- Analytics Platform
- Flash Sale Platform

without major future refactors.

The service MUST realistically tolerate:
- flash-sale traffic
- reservation races
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