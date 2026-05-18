# TASK-012 — SHIPMENT SERVICE

## Goal

Build a REAL production-grade Shipment Service.

This service is responsible for:
- shipment orchestration
- carrier integration
- shipping label orchestration
- tracking synchronization
- delivery lifecycle
- warehouse fulfillment coordination
- shipment reconciliation
- logistics retry handling

This is NOT a toy CRUD shipment service.

The Shipment Service must support:
- millions of shipments
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- replay-safe workflows
- logistics integration resiliency

The architecture MUST prioritize:
- shipment correctness
- delivery tracking integrity
- replay safety
- carrier resiliency
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
- Carrier SDK integrations

---

## Core Responsibilities

The Shipment Service MUST support:

### Shipment Orchestration
- shipment creation
- shipment assignment
- shipment coordination
- shipment cancellation
- shipment retry handling

### Carrier Integration
- multiple carrier providers
- carrier abstraction
- label generation
- shipment booking
- carrier failover hooks

### Tracking Synchronization
- tracking updates
- webhook synchronization
- tracking reconciliation
- delayed event handling

### Delivery Lifecycle
- awaiting pickup
- picked up
- in transit
- out for delivery
- delivered
- failed delivery
- returned

### Warehouse Coordination
- fulfillment coordination hooks
- warehouse shipment preparation
- packing workflow hooks

### Reconciliation
- carrier reconciliation
- tracking reconciliation
- delivery reconciliation
- async reconciliation jobs

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Shipment Service MUST:
- support replay-safe workflows
- support distributed retries
- support failure recovery
- support carrier abstraction

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The shipment system MUST tolerate:
- webhook replay
- carrier instability
- retry storms
- delayed tracking events
- duplicate delivery updates
- distributed deployments

---

## Folder Structure

Generate:

services/shipment/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── shipment/
│   ├── carriers/
│   ├── tracking/
│   ├── labels/
│   ├── reconciliation/
│   ├── fulfillment/
│   ├── delivery/
│   ├── webhooks/
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
- shipment records
- tracking history
- carrier mappings
- delivery lifecycle state
- webhook history
- reconciliation state

Generate:
- optimized schemas
- indexes
- immutable tracking history
- audit tables

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- tracking consistency

Never:
- mutate immutable tracking history
- trust carrier callbacks blindly
- ignore reconciliation
- tightly couple warehouse logic

---

## Redis Requirements

Use Redis for:
- idempotency cache
- webhook replay protection
- tracking synchronization cache
- temporary orchestration state

Generate:
- TTL strategy
- replay protection
- distributed coordination
- retry handling

Support:
- distributed deployments
- duplicate prevention
- webhook replay protection

The coordination layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- shipment created
- shipment booked
- shipment picked up
- shipment in transit
- shipment delivered
- shipment failed
- shipment returned
- tracking updated

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

## Carrier Integration Requirements

Support:
- multiple carrier providers
- carrier abstraction
- provider failover hooks
- webhook verification
- retry-safe carrier handling

Generate:
- carrier adapter interfaces
- provider abstraction layer
- webhook signature validation
- carrier retry workflows

The carrier integration MUST be production-grade.

No fake carrier abstraction.

---

## Tracking Requirements

Implement:
- tracking synchronization
- replay-safe tracking updates
- duplicate update handling
- delayed event handling
- immutable tracking history

Support:
- webhook replay
- out-of-order tracking events
- delayed delivery updates

Generate:
- tracking reconciliation jobs
- replay-safe consumers
- tracking integrity validation

---

## Delivery Lifecycle Requirements

Support:
- delivery lifecycle state machine
- valid transition enforcement
- failed delivery handling
- return workflows

Generate:
- lifecycle transition rules
- reconciliation workflows
- replay-safe lifecycle handling

The lifecycle system MUST be production-grade.

---

## Reconciliation Requirements

Support:
- carrier reconciliation
- delivery reconciliation
- tracking reconciliation
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
- /shipments
- /shipments/{id}
- /tracking
- /tracking/history
- /labels
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
- trust carrier callbacks blindly
- expose internal logistics metadata
- expose secret carrier credentials

Generate:
- authorization middleware
- webhook verification middleware
- replay validation
- logistics integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- shipment creation latency
- carrier API latency
- tracking sync latency
- webhook replay count
- reconciliation failures
- delivery transition latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive carrier credentials.

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
- carrier failover tests
- concurrency tests

Test:
- webhook replay
- duplicate tracking updates
- delayed delivery updates
- retry storms
- carrier instability
- reconciliation correctness
- high concurrency

---

## Output Requirements

Explain:
- shipment architecture
- carrier abstraction strategy
- tracking synchronization flow
- replay protection strategy
- reconciliation strategy
- delivery lifecycle strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy shipment CRUD service.
No fake carrier abstraction.
No naive tracking synchronization.

---

## Acceptance Criteria

The Shipment Service must support future integration with:
- Order Service
- Warehouse Service
- Notification Service
- Refund Service
- Analytics Platform

without major future refactors.

The service MUST realistically tolerate:
- carrier instability
- webhook replay
- retry storms
- delayed tracking updates
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