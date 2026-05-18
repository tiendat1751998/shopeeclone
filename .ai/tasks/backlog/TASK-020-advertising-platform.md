# TASK-020 — ADVERTISING PLATFORM

## Goal

Build a REAL production-grade Advertising Platform.

This platform is responsible for:
- sponsored products
- bidding systems
- ad ranking
- ad analytics
- campaign management
- budget pacing
- realtime bidding hooks
- targeting systems
- impression/click tracking
- advertiser billing hooks

This is NOT a toy ads service.

The Advertising Platform must support:
- ultra-high read traffic
- low-latency ranking
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe billing

The architecture MUST prioritize:
- ranking latency
- billing correctness
- advertiser isolation
- pacing correctness
- operational stability

---

## Tech Stack

Use:
- Golang
- Kafka
- Redis Cluster
- ClickHouse
- PostgreSQL
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- vector targeting hooks
- ML ranking hooks
- Flink

---

## Core Responsibilities

The Advertising Platform MUST support:

### Sponsored Products
- sponsored listings
- promoted products
- sponsored recommendations
- ad placement orchestration

### Campaign Management
- campaign lifecycle
- ad groups
- targeting configuration
- budget management

### Bidding Systems
- CPC bidding
- CPM bidding
- bid adjustments
- realtime bid evaluation

### Ad Ranking
- sponsored ranking
- relevance scoring
- quality scoring
- pacing-aware ranking

### Budget Pacing
- spend pacing
- daily budget limits
- distributed spend tracking
- pacing throttling

### Ad Analytics
- impressions
- clicks
- conversions
- ROAS analytics
- campaign analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate serving/ranking/billing
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Advertising Platform MUST:
- support low-latency ad serving
- support replay-safe billing
- support distributed pacing
- support advertiser isolation

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The advertising system MUST tolerate:
- retry storms
- duplicate events
- delayed billing events
- partial failures
- distributed deployments
- pacing inconsistencies

---

## Folder Structure

Generate:

platforms/advertising/
├── cmd/
├── internal/
│   ├── config/
│   ├── campaigns/
│   ├── bidding/
│   ├── ranking/
│   ├── pacing/
│   ├── targeting/
│   ├── billing/
│   ├── analytics/
│   ├── fraud/
│   ├── serving/
│   ├── synchronization/
│   ├── replay/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── ranking/
│   ├── scoring/
│   ├── targeting/
│   └── pipelines/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## PostgreSQL Requirements

Use PostgreSQL for:
- campaign metadata
- ad groups
- targeting configuration
- advertiser metadata
- billing state

Generate:
- optimized schemas
- indexes
- immutable billing history
- audit tables

Requirements:
- transactional correctness
- replay safety
- advertiser isolation

Never:
- mutate immutable billing history
- tightly couple ranking to billing writes

---

## Redis Requirements

Use Redis for:
- hot campaign cache
- pacing counters
- realtime budget tracking
- distributed coordination

Generate:
- TTL strategy
- replay protection
- distributed pacing coordination
- rate limiting

Support:
- realtime serving
- high concurrency
- distributed deployments

The coordination layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- impression events
- click events
- conversion events
- replay-safe billing
- distributed analytics

Generate:
- topic strategy
- partitioning strategy
- DLQ topics
- replay pipelines

Requirements:
- replay-safe ingestion
- idempotent producers
- consumer groups
- event versioning

Support:
- ultra-high traffic
- realtime analytics
- distributed consumers

---

## ClickHouse Requirements

Use ClickHouse for:
- impression analytics
- click analytics
- conversion analytics
- campaign analytics
- ROAS aggregation

Generate:
- aggregation strategy
- partitioning strategy
- TTL policies
- materialized views

Requirements:
- replay-safe aggregation
- high ingestion throughput
- low-cost OLAP queries

---

## Ad Serving Requirements

Support:
- low-latency ad serving
- sponsored product injection
- ad placement orchestration
- pacing-aware serving

Generate:
- ad serving APIs
- ranking orchestration
- cache orchestration
- fallback serving strategies

The serving system MUST tolerate:
- cache misses
- ranking lag
- stale pacing data
- distributed deployments

No fake ad serving architecture.

---

## Bidding Requirements

Support:
- CPC bidding
- CPM bidding
- realtime bid evaluation
- pacing-aware bidding

Generate:
- bidding engine
- bid evaluation pipelines
- replay-safe bidding
- distributed bid coordination

The bidding system MUST be production-grade.

---

## Ranking Requirements

Support:
- sponsored ranking
- relevance scoring
- quality scoring
- personalization hooks

Generate:
- ranking pipeline
- quality scoring strategy
- pacing-aware ranking
- replay-safe ranking updates

No naive ranking logic.

---

## Budget Pacing Requirements

Support:
- realtime spend tracking
- distributed budget enforcement
- pacing throttling
- budget reconciliation

Generate:
- pacing engine
- distributed spend coordination
- replay-safe spend aggregation
- reconciliation workflows

The pacing system MUST tolerate:
- duplicate events
- delayed billing events
- retry storms

No fake pacing architecture.

---

## Fraud Requirements

Support:
- click fraud hooks
- impression fraud hooks
- abuse scoring hooks

Generate:
- fraud synchronization hooks
- replay-safe fraud aggregation
- advertiser abuse detection hooks

---

## Event-Driven Requirements

Generate events for:
- impression recorded
- click recorded
- conversion recorded
- budget exhausted
- pacing throttled
- ranking updated

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
- distributed analytics
- replay-safe billing

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /ads/serve
- /ads/campaigns
- /ads/bidding
- /ads/analytics
- /ads/budgets

Support:
- pagination
- filtering
- realtime serving
- advertiser isolation

---

## Security Requirements

The platform MUST:
- validate advertiser ownership
- enforce RBAC
- sanitize input
- isolate advertiser data
- protect billing integrity

Never:
- expose internal ranking logic
- expose raw bidding internals
- expose private advertiser analytics
- trust client billing state

Generate:
- authorization middleware
- replay validation
- advertiser isolation
- billing integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- ad serving latency
- bid evaluation latency
- pacing latency
- click-through rate
- replay latency
- budget exhaustion count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive advertiser billing data.

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
- ranking tests
- pacing tests
- billing tests
- concurrency tests

Test:
- duplicate events
- replay storms
- delayed billing events
- budget exhaustion
- pacing correctness
- ultra-high serving traffic
- distributed deployments

---

## Output Requirements

Explain:
- advertising architecture
- bidding strategy
- pacing strategy
- ranking strategy
- replay-safe billing strategy
- analytics strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy advertising service.
No fake pacing engine.
No naive realtime bidding architecture.

---

## Acceptance Criteria

The Advertising Platform must support future integration with:
- Search Platform
- Recommendation Platform
- Billing Platform
- Fraud Platform
- Analytics Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate events
- delayed billing events
- distributed deployments
- ultra-high serving traffic

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