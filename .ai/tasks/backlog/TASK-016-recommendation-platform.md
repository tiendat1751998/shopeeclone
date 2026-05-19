# TASK-016 — RECOMMENDATION PLATFORM

## Goal

Build a REAL production-grade Recommendation Platform.

This platform is responsible for:
- recommendation pipelines
- personalization
- behavioral analytics
- feature stores
- ranking models
- candidate generation
- online/offline inference
- recommendation serving
- feedback learning loops

This is NOT a toy recommendation service.

The Recommendation Platform must support:
- millions of users
- ultra-high read traffic
- low-latency inference
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe ingestion

The architecture MUST prioritize:
- recommendation latency
- personalization quality
- ranking freshness
- feature consistency
- operational stability

---

## Tech Stack

Use:
- Golang
- Python
- Gin/Fiber
- gRPC
- Redis Cluster
- Kafka
- ClickHouse
- Elasticsearch/OpenSearch
- PostgreSQL or MySQL
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Feast feature store
- Ray Serve
- MLflow
- vector database hooks

---

## Core Responsibilities

The Recommendation Platform MUST support:

### Candidate Generation
- collaborative filtering hooks
- content-based candidate generation
- trending candidate generation
- seller-based recommendations
- category-based recommendations

### Personalization
- user behavior tracking
- real-time personalization
- user embeddings hooks
- affinity scoring

### Ranking Models
- recommendation ranking
- relevance scoring
- freshness scoring
- CTR optimization hooks
- engagement scoring

### Online Inference
- low-latency serving
- feature retrieval
- real-time ranking
- fallback recommendations

### Offline Pipelines
- feature generation
- embedding generation
- behavioral aggregation
- training dataset generation

### Feedback Learning
- click feedback
- impression tracking
- conversion feedback
- recommendation analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Recommendation Platform MUST:
- support ultra-low latency inference
- support replay-safe ingestion
- support distributed feature pipelines
- support cache orchestration

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The recommendation system MUST tolerate:
- retry storms
- duplicate events
- feature lag
- stale recommendations
- partial failures
- distributed deployments

---

## Folder Structure

Generate:

platforms/recommendation/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── candidates/
│   ├── ranking/
│   ├── personalization/
│   ├── embeddings/
│   ├── features/
│   ├── inference/
│   ├── feedback/
│   ├── analytics/
│   ├── synchronization/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── training/
│   ├── serving/
│   ├── embeddings/
│   ├── datasets/
│   └── pipelines/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Feature Store Requirements

Support:
- online features
- offline features
- replay-safe ingestion
- feature versioning
- point-in-time correctness

Generate:
- feature ingestion pipelines
- feature synchronization
- feature TTL strategy
- online/offline consistency handling

Support:
- high ingestion throughput
- low-latency retrieval
- distributed synchronization

No fake feature store architecture.

---

## Candidate Generation Requirements

Support:
- collaborative filtering hooks
- content-based filtering
- trending recommendations
- seller affinity
- category affinity

Generate:
- candidate generation pipelines
- candidate caching strategy
- distributed candidate retrieval

The candidate system MUST tolerate:
- stale candidates
- retry storms
- duplicate ingestion events

---

## Ranking Requirements

Support:
- recommendation ranking
- relevance scoring
- freshness scoring
- engagement scoring
- personalization scoring

Generate:
- ranking pipeline
- online ranking inference
- ranking reconciliation
- fallback ranking strategy

The ranking system MUST be production-grade.

---

## Online Inference Requirements

Support:
- low-latency inference
- online feature retrieval
- inference caching
- fallback recommendations

Generate:
- inference serving APIs
- cache orchestration
- timeout handling
- fallback strategies

Support:
- degraded-mode recommendations
- partial feature outages
- distributed deployments

No naive inference architecture.

---

## Offline Pipeline Requirements

Support:
- embedding generation
- behavioral aggregation
- training data generation
- recommendation analytics

Generate:
- distributed pipelines
- replay-safe ingestion
- scheduled jobs
- aggregation workflows

The offline pipeline MUST support:
- massive event ingestion
- replay-safe aggregation
- scalable analytics

---

## Behavioral Analytics Requirements

Use ClickHouse for:
- impressions
- clicks
- conversions
- recommendation analytics
- engagement aggregation

Generate:
- analytics ingestion flows
- replay-safe ingestion
- aggregation strategy

Support:
- high ingestion throughput
- distributed aggregation
- low-cost analytics queries

---

## Event-Driven Requirements

Generate events for:
- recommendation requested
- impression recorded
- click recorded
- conversion recorded
- feature updated
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
- distributed ingestion
- async synchronization

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /recommendations
- /recommendations/home
- /recommendations/product
- /recommendations/similar
- /recommendations/trending

Support:
- pagination
- filtering
- fallback handling
- personalization hooks

---

## Security Requirements

The platform MUST:
- validate requests
- enforce RBAC
- sanitize input
- isolate recommendation analytics
- protect feature integrity

Never:
- expose raw ML internals
- expose private embeddings
- expose raw behavioral analytics
- trust external ranking state blindly

Generate:
- authorization middleware
- replay validation
- analytics isolation
- feature integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- recommendation latency
- inference latency
- feature retrieval latency
- cache hit ratio
- stale recommendation rate
- ranking latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive behavioral data.

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
- fallback recommendation workflows

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
- ranking tests
- inference tests
- replay tests
- cache tests
- analytics ingestion tests
- concurrency tests

Test:
- stale recommendations
- duplicate events
- retry storms
- feature lag
- degraded inference
- high concurrency
- ultra-high read traffic

---

## Output Requirements

Explain:
- recommendation architecture
- feature store strategy
- ranking strategy
- inference strategy
- candidate generation strategy
- personalization strategy
- analytics strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy recommendation service.
No fake ML pipelines.
No naive inference architecture.

---

## Acceptance Criteria

The Recommendation Platform must support future integration with:
- Search Platform
- Product Catalog Service
- Analytics Platform
- Advertising Platform
- User Behavior Platform

without major future refactors.

The platform MUST realistically tolerate:
- stale features
- duplicate events
- retry storms
- degraded inference
- distributed deployments
- ultra-high read traffic

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