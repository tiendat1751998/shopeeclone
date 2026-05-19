# TASK-024 — SEARCH PLATFORM

## Goal

Build a REAL production-grade Search Platform.

This platform is responsible for:
- distributed search
- indexing pipelines
- autocomplete
- typo tolerance
- ranking pipelines
- realtime indexing
- query understanding
- faceted filtering
- search analytics
- personalization hooks

This is NOT a toy search service.

The Search Platform must support:
- ultra-high QPS
- low-latency retrieval
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe indexing

The architecture MUST prioritize:
- retrieval latency
- indexing correctness
- ranking freshness
- query resiliency
- operational stability

---

## Tech Stack

Use:
- Golang
- Elasticsearch or OpenSearch
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
- vector search hooks
- ANN indexing hooks
- ML ranking hooks

---

## Core Responsibilities

The Search Platform MUST support:

### Distributed Search
- product search
- seller search
- category search
- distributed query execution

### Indexing Pipelines
- realtime indexing
- incremental indexing
- replay-safe indexing
- bulk reindex orchestration

### Autocomplete
- search suggestions
- trending suggestions
- query completion
- prefix indexing

### Typo Tolerance
- fuzzy matching
- typo correction
- synonym expansion
- token normalization

### Ranking Pipelines
- relevance scoring
- popularity scoring
- personalization hooks
- sponsored ranking hooks

### Search Analytics
- query analytics
- click analytics
- CTR analytics
- no-result analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate retrieval/indexing/ranking
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Search Platform MUST:
- support replay-safe indexing
- support distributed retrieval
- support low-latency queries
- support realtime indexing

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The search system MUST tolerate:
- retry storms
- duplicate indexing events
- delayed indexing updates
- partial failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/search/
├── cmd/
├── internal/
│   ├── config/
│   ├── retrieval/
│   ├── indexing/
│   ├── ranking/
│   ├── autocomplete/
│   ├── typo/
│   ├── analytics/
│   ├── personalization/
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
├── indexing/
│   ├── workers/
│   ├── pipelines/
│   ├── bulk/
│   └── reindex/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Elasticsearch/OpenSearch Requirements

Use Elasticsearch/OpenSearch for:
- distributed retrieval
- inverted indexing
- autocomplete
- typo tolerance
- faceted search

Generate:
- index templates
- shard strategies
- replica strategies
- ILM policies
- alias strategies

Requirements:
- replay-safe indexing
- distributed retrieval
- index versioning
- rolling reindex support

Never:
- use naive indexing
- ignore shard balancing
- ignore index rollover
- tightly couple writes to search serving

---

## Indexing Requirements

Support:
- realtime indexing
- incremental indexing
- bulk indexing
- replay-safe indexing

Generate:
- indexing pipelines
- replay-safe synchronization
- reindex orchestration
- failure recovery workflows

The indexing system MUST tolerate:
- duplicate events
- delayed updates
- replay storms

No fake indexing architecture.

---

## Retrieval Requirements

Support:
- distributed retrieval
- faceted filtering
- category filtering
- seller filtering
- low-latency search

Generate:
- retrieval orchestration
- distributed query pipelines
- fallback retrieval strategies
- cache-aware retrieval

The retrieval system MUST be production-grade.

---

## Ranking Requirements

Support:
- relevance ranking
- popularity ranking
- personalization hooks
- sponsored ranking hooks

Generate:
- ranking pipelines
- scoring orchestration
- replay-safe ranking updates
- ranking feature pipelines

No naive ranking logic.

---

## Autocomplete Requirements

Support:
- realtime suggestions
- trending suggestions
- prefix search
- query completions

Generate:
- autocomplete pipelines
- replay-safe suggestion aggregation
- distributed suggestion updates
- cache synchronization

---

## Typo Tolerance Requirements

Support:
- fuzzy matching
- typo correction
- synonym handling
- token normalization

Generate:
- normalization pipelines
- typo dictionaries
- synonym orchestration
- query rewriting workflows

The typo-tolerance system MUST be realistic.

---

## Kafka Requirements

Use Kafka for:
- indexing events
- ranking updates
- replay-safe synchronization
- search analytics streaming

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
- distributed indexing
- async synchronization
- analytics streaming

---

## Redis Requirements

Use Redis for:
- hot query cache
- autocomplete cache
- trending queries
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

---

## ClickHouse Requirements

Use ClickHouse for:
- query analytics
- click analytics
- CTR aggregation
- no-result analytics
- search performance analytics

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

## Query Understanding Requirements

Support:
- token normalization
- query rewriting
- synonym expansion
- intent classification hooks

Generate:
- normalization pipelines
- replay-safe query processing
- distributed enrichment workflows
- query analytics hooks

---

## Event-Driven Requirements

Generate events for:
- document indexed
- reindex triggered
- query executed
- no-result detected
- ranking updated
- autocomplete updated

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
- replay-safe indexing
- distributed retrieval

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /search
- /search/autocomplete
- /search/trending
- /search/categories
- /search/filters

Support:
- pagination
- faceted filtering
- realtime autocomplete
- typo-tolerant queries

---

## Security Requirements

The platform MUST:
- validate query input
- enforce RBAC
- sanitize queries
- isolate indexing pipelines
- protect ranking integrity

Never:
- expose internal ranking signals
- expose raw index topology
- trust external indexing events blindly
- allow query injection

Generate:
- authorization middleware
- replay validation
- indexing integrity validation
- query sanitization

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- query latency
- indexing latency
- reindex duration
- cache hit ratio
- replay latency
- no-result count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive ranking internals.

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
- retrieval tests
- ranking tests
- indexing tests
- concurrency tests

Test:
- duplicate indexing events
- replay storms
- delayed indexing updates
- search latency
- ranking correctness
- ultra-high QPS
- distributed deployments

---

## Output Requirements

Explain:
- search architecture
- indexing strategy
- retrieval strategy
- ranking strategy
- replay-safe indexing strategy
- autocomplete strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy search service.
No fake indexing architecture.
No naive retrieval logic.

---

## Acceptance Criteria

The Search Platform must support future integration with:
- Recommendation Platform
- Advertising Platform
- Analytics Platform
- Inventory Platform
- User Behavior Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate indexing events
- delayed indexing updates
- distributed deployments
- ultra-high search QPS

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.