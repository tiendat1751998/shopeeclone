# TASK-007 — CART SERVICE

## Goal

Build a REAL production-grade Cart Service.

This service is responsible for:
- shopping cart management
- cart aggregation
- pricing preview
- inventory preview
- promotion preview
- seller grouping
- cart synchronization
- checkout preparation
- reservation preview
- multi-device cart consistency

This is NOT a toy CRUD cart service.

The Cart Service must support:
- millions of carts
- high concurrency
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- eventual consistency
- fault tolerance

The architecture MUST prioritize:
- low-latency cart reads
- scalability
- consistency handling
- checkout readiness
- operational resilience

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
- Lua scripting for Redis operations

---

## Core Responsibilities

The Cart Service MUST support:

### Cart Management
- add item
- update quantity
- remove item
- clear cart
- save for later
- cart expiration

### Cart Aggregation
- seller grouping
- SKU aggregation
- pricing aggregation
- shipping preview
- inventory preview

### Checkout Preparation
- reservation preview
- checkout validation
- cart snapshot generation
- checkout payload generation

### Multi-Device Synchronization
- device synchronization
- guest-to-user cart merge
- session cart merge
- concurrent cart updates

### Promotion Preview
- promotion preview
- voucher preview
- shipping discount preview
- eligibility validation

### Inventory Preview
- stock availability preview
- reservation readiness
- unavailable item handling

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Cart Service MUST:
- be read-optimized
- support low-latency reads
- support horizontal scaling
- support distributed cache
- support async recalculation

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The cart system MUST tolerate:
- stale pricing
- stale inventory
- retry storms
- duplicate events
- concurrent updates

---

## Folder Structure

Generate:

services/cart/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── cart/
│   ├── aggregation/
│   ├── pricing/
│   ├── promotion/
│   ├── inventory/
│   ├── checkout/
│   ├── sync/
│   ├── merge/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   ├── validation/
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
- persistent carts
- cart snapshots
- merge history
- checkout preparation state

Use Redis for:
- hot carts
- active sessions
- low-latency reads
- temporary aggregation
- cart TTL handling

Generate:
- optimized schemas
- indexes
- cache strategy
- invalidation strategy

Requirements:
- pagination support
- high read optimization
- eventual consistency support
- stale data handling

Never:
- use naive cart recalculation
- ignore concurrent updates
- tightly couple checkout logic

---

## Redis Requirements

Use Redis for:
- hot cart cache
- cart synchronization
- session carts
- temporary checkout state
- inventory preview cache

Generate:
- TTL strategy
- cache invalidation strategy
- distributed cache consistency
- retry handling
- stale cache recovery

Support:
- high concurrency
- distributed deployments
- multi-device synchronization

The cache layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- cart updated
- item added
- item removed
- cart merged
- checkout prepared
- cart expired

Use:
- Kafka or NATS JetStream

Requirements:
- retries
- DLQ
- replay-safe consumers
- idempotency
- event versioning
- consumer groups

Support:
- eventual consistency
- async recalculation
- distributed synchronization

No fake async architecture.

---

## Pricing Aggregation Requirements

Support:
- pricing preview
- discount preview
- seller grouping
- subtotal calculation
- shipping estimation hooks

Requirements:
- stale pricing tolerance
- async recalculation
- distributed pricing updates

The service MUST NOT own pricing truth.

Pricing truth belongs elsewhere.

---

## Inventory Preview Requirements

Support:
- stock preview
- reservation readiness
- unavailable item detection
- quantity validation

Requirements:
- eventual consistency
- stale inventory tolerance
- async validation

The service MUST NOT own inventory truth.

Inventory truth belongs to Inventory Service.

---

## Checkout Preparation Requirements

Generate:
- cart snapshot
- checkout payload
- seller grouping
- reservation preparation
- validation summary

Support:
- replay safety
- idempotency
- concurrent checkout attempts

The cart service MUST prepare checkout,
NOT execute orders.

---

## Multi-Device Sync Requirements

Support:
- concurrent sessions
- guest cart merge
- user cart merge
- conflict resolution
- stale cart handling

Generate:
- merge strategy
- conflict handling
- reconciliation flow

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /cart
- /cart/items
- /cart/merge
- /cart/checkout-preview
- /cart/promotions

Support:
- pagination
- filtering
- idempotency keys
- validation

---

## Security Requirements

The service MUST:
- validate ownership
- enforce RBAC
- validate cart operations
- sanitize input
- prevent unauthorized access

Never:
- trust client pricing
- trust client inventory
- expose internal aggregation metadata

Generate:
- authorization middleware
- ownership validation
- cart isolation validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- cart latency
- Redis latency
- cache hit ratio
- merge conflicts
- checkout preview latency
- promotion preview latency
- inventory preview latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive user/cart data.

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
- concurrency tests
- merge tests
- cache tests
- replay tests
- checkout preview tests

Test:
- concurrent updates
- stale pricing
- stale inventory
- guest cart merge
- duplicate events
- retry storms
- high concurrency

---

## Output Requirements

Explain:
- cart architecture
- aggregation strategy
- cache strategy
- pricing preview flow
- inventory preview flow
- checkout preparation flow
- merge strategy
- event-driven synchronization
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy cart CRUD service.
No fake aggregation.
No fake async architecture.

---

## Acceptance Criteria

The Cart Service must support future integration with:
- Product Service
- Inventory Service
- Promotion Service
- Checkout Service
- Order Service

without major future refactors.

The service MUST realistically tolerate:
- stale pricing
- stale inventory
- high concurrency
- distributed deployments
- retry storms

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