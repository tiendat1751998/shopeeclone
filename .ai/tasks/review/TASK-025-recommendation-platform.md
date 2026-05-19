# TASK-025 — RECOMMENDATION PLATFORM

## Goal

Build a REAL production-grade Recommendation Platform.

This platform is responsible for:
- personalized recommendations
- ranking pipelines
- feature stores
- candidate generation
- realtime personalization
- behavioral scoring
- feed generation
- recommendation analytics
- vector retrieval hooks
- model serving orchestration

This is NOT a toy recommendation service.

The Recommendation Platform must support:
- ultra-high feed QPS
- low-latency ranking
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe feature synchronization

The architecture MUST prioritize:
- personalization freshness
- ranking latency
- behavioral correctness
- feed resiliency
- operational stability

---

## Tech Stack

Use:
- Golang
- Python
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
- Feast
- Milvus
- Vespa
- Ray
- TensorFlow Serving
- Triton Inference Server

---

## Core Responsibilities

The Recommendation Platform MUST support:

### Personalized Recommendations
- home feed recommendations
- product recommendations
- seller recommendations
- realtime personalization

### Candidate Generation
- collaborative filtering hooks
- content-based candidates
- trending candidates
- behavioral candidates

### Ranking Pipelines
- personalized ranking
- CTR prediction hooks
- engagement scoring
- sponsored ranking hooks

### Feature Stores
- realtime features
- behavioral features
- product features
- user embeddings

### Feed Generation
- personalized feeds
- realtime feed refresh
- replay-safe feed generation
- fallback feeds

### Recommendation Analytics
- CTR analytics
- engagement analytics
- recommendation quality analytics
- ranking analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate candidate/ranking/feature-serving
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Recommendation Platform MUST:
- support replay-safe feature synchronization
- support low-latency ranking
- support realtime personalization
- support degraded inference modes

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The recommendation system MUST tolerate:
- retry storms
- duplicate feature events
- delayed behavioral updates
- partial inference failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/recommendation/
├── cmd/
├── internal/
│   ├── config/
│   ├── candidates/
│   ├── ranking/
│   ├── features/
│   ├── feeds/
│   ├── embeddings/
│   ├── personalization/
│   ├── analytics/
│   ├── synchronization/
│   ├── replay/
│   ├── inference/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── ranking/
│   ├── embeddings/
│   ├── training/
│   ├── pipelines/
│   └── serving/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## PostgreSQL Requirements

Use PostgreSQL for:
- recommendation metadata
- feature metadata
- feed state
- model metadata
- replay metadata

Generate:
- optimized schemas
- indexes
- immutable recommendation audit tables
- replay-safe synchronization tables

Requirements:
- replay safety
- feed consistency
- metadata correctness

Never:
- tightly couple ranking requests to OLTP writes
- mutate immutable recommendation audit history

---

## Redis Requirements

Use Redis for:
- hot feed cache
- realtime features
- ranking cache
- distributed coordination

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- ultra-high QPS
- low-latency reads
- distributed deployments

The cache layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- behavioral events
- feature updates
- ranking updates
- replay-safe synchronization
- analytics streaming

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
- distributed personalization
- async synchronization
- realtime feature updates

---

## ClickHouse Requirements

Use ClickHouse for:
- recommendation analytics
- CTR analytics
- engagement analytics
- ranking analytics
- feed analytics

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

## Candidate Generation Requirements

Support:
- trending candidates
- collaborative filtering hooks
- content-based candidates
- behavioral candidates

Generate:
- distributed candidate pipelines
- replay-safe candidate generation
- fallback candidate strategies
- candidate cache orchestration

The candidate generation system MUST tolerate:
- delayed behavioral events
- replay storms
- partial failures

No fake candidate generation.

---

## Ranking Requirements

Support:
- personalized ranking
- engagement scoring
- CTR prediction hooks
- sponsored ranking hooks

Generate:
- ranking pipelines
- feature orchestration
- replay-safe ranking updates
- degraded inference fallback logic

No naive ranking logic.

---

## Feature Store Requirements

Support:
- realtime features
- behavioral features
- embeddings
- replay-safe feature updates

Generate:
- feature synchronization pipelines
- replay-safe aggregation
- feature invalidation workflows
- feature freshness orchestration

The feature system MUST be realistic.

---

## Feed Generation Requirements

Support:
- personalized feeds
- realtime refresh
- fallback feeds
- replay-safe feed generation

Generate:
- feed orchestration
- cache-aware feed serving
- distributed feed pipelines
- feed reconciliation workflows

---

## ML Inference Requirements

Support:
- ranking inference
- embedding retrieval
- degraded inference fallback
- realtime feature retrieval

Generate:
- inference orchestration
- feature retrieval hooks
- replay-safe inference pipelines
- inference fallback workflows

The inference system MUST tolerate:
- feature outages
- model lag
- inference timeouts
- retry storms

No fake ML architecture.

---

## Vector Retrieval Requirements

Support:
- vector similarity hooks
- embedding retrieval
- ANN retrieval hooks
- semantic recommendation hooks

Generate:
- embedding synchronization
- replay-safe vector updates
- retrieval orchestration
- fallback retrieval logic

---

## Event-Driven Requirements

Generate events for:
- feed generated
- feature updated
- ranking updated
- candidate generated
- model refreshed
- inference degraded

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
- distributed ranking
- replay-safe synchronization

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /recommendations/feed
- /recommendations/products
- /recommendations/sellers
- /recommendations/trending
- /recommendations/analytics

Support:
- pagination
- realtime personalization
- replay-safe feeds
- low-latency ranking

---

## Security Requirements

The platform MUST:
- validate recommendation ownership
- enforce RBAC
- sanitize feature input
- isolate model serving
- protect ranking integrity

Never:
- expose raw ranking signals
- expose internal embeddings
- expose model internals
- trust external ranking signals blindly

Generate:
- authorization middleware
- replay validation
- ranking integrity validation
- model isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- feed latency
- ranking latency
- inference latency
- feature freshness
- replay latency
- degraded inference count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive embeddings or ranking features.

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
- feature-store tests
- inference tests
- concurrency tests

Test:
- duplicate feature events
- replay storms
- delayed behavioral updates
- feed latency
- inference degradation
- ultra-high QPS
- distributed deployments

---

## Output Requirements

Explain:
- recommendation architecture
- candidate generation strategy
- ranking strategy
- feature-store strategy
- replay-safe feature strategy
- inference strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy recommendation service.
No fake feature-store architecture.
No naive personalization logic.

---

## Acceptance Criteria

The Recommendation Platform must support future integration with:
- Search Platform
- Advertising Platform
- User Behavior Platform
- Analytics Platform
- Live Commerce Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate feature events
- delayed behavioral updates
- distributed deployments
- ultra-high feed QPS

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.