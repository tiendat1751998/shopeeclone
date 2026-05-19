# TASK-023 — LOGISTICS & DELIVERY PLATFORM

## Goal

Build a REAL production-grade Logistics & Delivery Platform.

This platform is responsible for:
- shipment orchestration
- route optimization
- delivery tracking
- warehouse routing
- courier integrations
- realtime tracking
- delivery estimation
- pickup orchestration
- dispatch coordination
- fulfillment synchronization

This is NOT a toy shipment service.

The Logistics & Delivery Platform must support:
- massive shipment volume
- realtime tracking
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe shipment processing

The architecture MUST prioritize:
- shipment correctness
- tracking freshness
- route optimization
- courier resiliency
- operational stability

---

## Tech Stack

Use:
- Golang
- PostgreSQL
- Kafka
- Redis Cluster
- ClickHouse
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- graph routing engines
- geospatial indexing
- optimization workers

---

## Core Responsibilities

The Logistics & Delivery Platform MUST support:

### Shipment Orchestration
- shipment creation
- shipment lifecycle management
- split shipments
- multi-package orchestration

### Delivery Tracking
- realtime tracking updates
- tracking timelines
- delivery milestones
- courier synchronization

### Route Optimization
- warehouse routing
- dispatch optimization
- courier route optimization
- delivery zone orchestration

### Courier Integrations
- external courier synchronization
- webhook ingestion
- replay-safe courier updates
- courier reconciliation

### Fulfillment Synchronization
- warehouse synchronization
- inventory movement hooks
- pickup orchestration
- dispatch coordination

### Delivery Estimation
- ETA calculation
- dynamic delay estimation
- route-aware estimation
- traffic-aware estimation hooks

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate orchestration/tracking/routing
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Logistics Platform MUST:
- support replay-safe shipment processing
- support distributed courier integrations
- support realtime tracking
- support route optimization

Use:
- Saga orchestration where appropriate
- CQRS where appropriate
- dependency injection
- resilience patterns

The logistics system MUST tolerate:
- retry storms
- duplicate shipment events
- delayed courier updates
- partial failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/logistics-delivery/
├── cmd/
├── internal/
│   ├── config/
│   ├── shipments/
│   ├── tracking/
│   ├── routing/
│   ├── dispatch/
│   ├── couriers/
│   ├── fulfillment/
│   ├── pickups/
│   ├── estimations/
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
- shipment metadata
- routing metadata
- courier mappings
- dispatch state
- fulfillment state

Generate:
- optimized schemas
- indexes
- immutable shipment timelines
- audit tables

Requirements:
- ACID correctness
- replay safety
- shipment consistency

Never:
- mutate immutable tracking history
- tightly couple courier webhooks to shipment writes

---

## Shipment Requirements

Support:
- shipment lifecycle management
- split shipments
- package grouping
- replay-safe shipment transitions

Generate:
- shipment orchestration engine
- replay-safe transition handling
- shipment reconciliation hooks
- distributed shipment workflows

The shipment system MUST tolerate:
- duplicate events
- delayed courier updates
- replay storms

No fake shipment orchestration.

---

## Tracking Requirements

Support:
- realtime tracking
- tracking timelines
- delivery milestones
- courier synchronization

Generate:
- tracking pipelines
- replay-safe tracking aggregation
- distributed tracking synchronization
- tracking reconciliation

The tracking system MUST be production-grade.

---

## Routing Requirements

Support:
- warehouse routing
- dispatch optimization
- delivery zone optimization
- route-aware estimations

Generate:
- routing engine hooks
- distributed optimization workers
- replay-safe routing updates
- fallback routing workflows

No naive routing architecture.

---

## Courier Integration Requirements

Support:
- courier webhooks
- courier polling
- replay-safe courier synchronization
- courier reconciliation

Generate:
- integration adapters
- webhook ingestion pipelines
- replay-safe processing
- courier failure handling

The courier integration layer MUST tolerate:
- webhook duplication
- delayed updates
- external API failures
- retry storms

---

## Fulfillment Requirements

Support:
- warehouse synchronization
- pickup orchestration
- dispatch coordination
- inventory movement hooks

Generate:
- fulfillment orchestration
- replay-safe warehouse synchronization
- dispatch workflows
- warehouse reconciliation

---

## Kafka Requirements

Use Kafka for:
- shipment events
- tracking events
- courier events
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
- shipment streaming
- distributed tracking
- async synchronization

---

## Redis Requirements

Use Redis for:
- hot shipment cache
- tracking cache
- distributed coordination
- replay coordination

Generate:
- TTL strategy
- distributed coordination
- replay protection
- rate limiting

Support:
- high concurrency
- distributed deployments
- realtime reads

---

## ClickHouse Requirements

Use ClickHouse for:
- shipment analytics
- delivery analytics
- courier analytics
- routing analytics
- ETA analytics

Generate:
- aggregation strategy
- partitioning strategy
- TTL policies
- materialized views

Requirements:
- replay-safe aggregation
- high ingestion throughput
- realtime analytics

---

## Estimation Requirements

Support:
- ETA calculations
- delay estimation
- route-aware estimations
- courier-aware estimations

Generate:
- estimation pipelines
- replay-safe estimation recalculation
- distributed estimation workers
- fallback estimation logic

The estimation system MUST be realistic.

---

## Event-Driven Requirements

Generate events for:
- shipment created
- shipment dispatched
- shipment delayed
- shipment delivered
- courier update received
- pickup completed

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
- distributed shipment workflows

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /shipments
- /tracking
- /couriers
- /dispatch
- /estimations

Support:
- pagination
- filtering
- replay-safe updates
- realtime tracking queries

---

## Security Requirements

The platform MUST:
- validate shipment ownership
- enforce RBAC
- sanitize input
- isolate courier integrations
- protect shipment integrity

Never:
- expose internal routing topology
- expose courier credentials
- trust external courier updates blindly
- allow unauthorized shipment mutations

Generate:
- authorization middleware
- replay validation
- shipment integrity validation
- courier isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- shipment orchestration latency
- tracking freshness
- courier API latency
- ETA calculation latency
- replay latency
- delivery delay count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive courier credentials.

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
- tracking tests
- routing tests
- courier integration tests
- concurrency tests

Test:
- duplicate courier events
- replay storms
- delayed updates
- shipment consistency
- ETA recalculation
- high concurrency
- distributed deployments

---

## Output Requirements

Explain:
- logistics architecture
- routing strategy
- tracking strategy
- courier integration strategy
- replay-safe shipment strategy
- estimation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy shipment service.
No fake routing architecture.
No naive tracking logic.

---

## Acceptance Criteria

The Logistics & Delivery Platform must support future integration with:
- Order Service
- Inventory Platform
- Billing Platform
- Analytics Platform
- Notification Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate courier events
- delayed tracking updates
- distributed deployments
- massive shipment volume

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