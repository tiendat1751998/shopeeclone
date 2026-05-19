# TASK-028 — ADVERTISING PLATFORM

## Goal

Build a REAL production-grade Advertising Platform.

This platform is responsible for:
- ad serving
- bidding pipelines
- sponsored ranking
- campaign orchestration
- realtime targeting
- ad analytics
- budget pacing
- auction orchestration
- advertiser billing hooks
- attribution pipelines

This is NOT a toy advertising service.

The Advertising Platform must support:
- ultra-low latency ad serving
- realtime auctions
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe ad tracking

The architecture MUST prioritize:
- ad-serving latency
- auction correctness
- budget correctness
- targeting freshness
- operational stability

---

## Tech Stack

Use:
- Golang
- Kafka
- Redis Cluster
- ClickHouse
- PostgreSQL
- Elasticsearch/OpenSearch
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- vector targeting hooks
- ML ranking hooks
- feature stores
- Flink

---

## Core Responsibilities

The Advertising Platform MUST support:

### Ad Serving
- sponsored product serving
- sponsored seller serving
- realtime ad retrieval
- low-latency ad ranking

### Auction Orchestration
- realtime bidding
- CPC/CPM auctions
- sponsored ranking auctions
- replay-safe auction execution

### Campaign Management
- campaign orchestration
- ad-group management
- targeting configuration
- pacing orchestration

### Realtime Targeting
- behavioral targeting
- keyword targeting
- contextual targeting
- audience segmentation

### Budget Pacing
- spend pacing
- throttling
- budget exhaustion prevention
- replay-safe spend accounting

### Attribution & Analytics
- impression tracking
- click tracking
- conversion attribution
- advertiser analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate serving/bidding/analytics
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Advertising Platform MUST:
- support replay-safe tracking
- support distributed auctions
- support realtime targeting
- support degraded serving fallback

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The advertising system MUST tolerate:
- retry storms
- duplicate tracking events
- delayed attribution updates
- partial auction failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/advertising/
├── cmd/
├── internal/
│   ├── config/
│   ├── serving/
│   ├── bidding/
│   ├── auctions/
│   ├── campaigns/
│   ├── targeting/
│   ├── pacing/
│   ├── analytics/
│   ├── attribution/
│   ├── sponsored/
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
│   ├── targeting/
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
- campaign metadata
- ad-group metadata
- pacing metadata
- attribution metadata
- replay metadata

Generate:
- optimized schemas
- indexes
- immutable tracking audit tables
- replay-safe synchronization tables

Requirements:
- replay safety
- budget correctness
- campaign consistency

Never:
- tightly couple ad serving to OLTP writes
- mutate immutable impression history

---

## Redis Requirements

Use Redis for:
- hot ad cache
- pacing counters
- realtime targeting cache
- distributed coordination

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- ultra-low latency serving
- high concurrency
- distributed deployments

The cache layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- impression events
- click events
- auction events
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
- distributed auctions
- async synchronization
- attribution pipelines

---

## ClickHouse Requirements

Use ClickHouse for:
- impression analytics
- CTR analytics
- conversion analytics
- pacing analytics
- advertiser analytics

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

## Elasticsearch/OpenSearch Requirements

Use Elasticsearch/OpenSearch for:
- keyword targeting
- ad retrieval
- campaign search
- sponsored retrieval

Generate:
- index templates
- shard strategies
- alias strategies
- ILM policies

Requirements:
- replay-safe indexing
- low-latency retrieval
- distributed querying

---

## Ad Serving Requirements

Support:
- sponsored product serving
- low-latency ranking
- contextual retrieval
- replay-safe tracking

Generate:
- ad-serving pipelines
- distributed serving orchestration
- fallback serving strategies
- replay-safe impression handling

The serving system MUST tolerate:
- delayed targeting updates
- replay storms
- partial serving failures

No fake ad-serving architecture.

---

## Auction Requirements

Support:
- realtime auctions
- CPC auctions
- CPM auctions
- sponsored ranking auctions

Generate:
- auction engine
- replay-safe bidding
- distributed auction orchestration
- fallback auction logic

No naive bidding logic.

---

## Targeting Requirements

Support:
- behavioral targeting
- contextual targeting
- keyword targeting
- audience segmentation

Generate:
- targeting pipelines
- replay-safe targeting synchronization
- distributed enrichment workflows
- cache-aware targeting

The targeting system MUST be realistic.

---

## Budget Pacing Requirements

Support:
- spend pacing
- budget throttling
- replay-safe spend accounting
- pacing recalculation

Generate:
- pacing orchestration
- replay-safe counters
- distributed pacing workers
- fallback throttling logic

The pacing system MUST tolerate:
- duplicate impression events
- replay storms
- delayed attribution

---

## Attribution Requirements

Support:
- impression attribution
- click attribution
- conversion attribution
- replay-safe attribution

Generate:
- attribution pipelines
- distributed attribution aggregation
- replay-safe tracking
- delayed conversion handling

---

## ML Inference Requirements

Support:
- CTR prediction hooks
- relevance scoring
- targeting inference
- degraded inference fallback

Generate:
- inference orchestration
- replay-safe feature retrieval
- distributed inference pipelines
- inference fallback workflows

The inference system MUST tolerate:
- feature outages
- model lag
- inference timeouts
- retry storms

No fake ML architecture.

---

## Event-Driven Requirements

Generate events for:
- impression served
- click recorded
- auction completed
- campaign exhausted
- pacing recalculated
- attribution completed

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
- distributed serving
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
- /ads/serve
- /ads/campaigns
- /ads/analytics
- /ads/targeting
- /ads/auctions

Support:
- realtime serving
- replay-safe tracking
- low-latency bidding
- campaign orchestration

---

## Security Requirements

The platform MUST:
- validate advertiser ownership
- enforce RBAC
- sanitize targeting input
- isolate campaign execution
- protect auction integrity

Never:
- expose raw bidding internals
- expose targeting signals
- expose advertiser secrets
- trust external attribution blindly

Generate:
- authorization middleware
- replay validation
- auction integrity validation
- targeting isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- ad-serving latency
- auction latency
- pacing latency
- CTR
- replay latency
- attribution delay

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive advertiser secrets.

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
- auction tests
- targeting tests
- pacing tests
- concurrency tests

Test:
- duplicate impression events
- replay storms
- delayed attribution
- auction correctness
- ultra-high QPS
- distributed deployments
- partial serving failures

---

## Output Requirements

Explain:
- advertising architecture
- serving strategy
- auction strategy
- targeting strategy
- replay-safe tracking strategy
- pacing strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy ad service.
No fake bidding architecture.
No naive targeting logic.

---

## Acceptance Criteria

The Advertising Platform must support future integration with:
- Search Platform
- Recommendation Platform
- Billing Platform
- Analytics Platform
- Fraud Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate tracking events
- delayed attribution
- distributed deployments
- ultra-high ad-serving QPS

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.