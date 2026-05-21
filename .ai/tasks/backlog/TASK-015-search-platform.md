# TASK-015 — SEARCH PLATFORM

## Goal

Build a REAL production-grade Search Platform.

This platform is responsible for:
- distributed product search
- autocomplete
- typo tolerance
- ranking systems
- indexing pipelines
- relevance scoring
- search analytics
- query understanding
- search cache orchestration

This is NOT a toy search service.

The Search Platform must support:
- millions of queries per minute
- ultra-low latency reads
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- replay-safe indexing

The architecture MUST prioritize:
- search latency
- ranking correctness
- indexing consistency
- cache efficiency
- operational stability

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- Redis Cluster
- Elasticsearch/OpenSearch
- Kafka or NATS JetStream
- ClickHouse
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- vector search hooks
- personalization hooks

---

## Core Responsibilities

The Search Platform MUST support:

### Product Search
- keyword search
- category search
- SKU search
- seller search
- attribute filtering

### Autocomplete
- search suggestions
- query completion
- trending queries
- typo-tolerant suggestions

### Ranking Systems
- relevance scoring
- popularity scoring
- freshness scoring
- sponsored ranking hooks
- personalization hooks

### Query Understanding
- typo tolerance
- synonym handling
- query normalization
- tokenization
- language-aware parsing

### Indexing Pipelines
- distributed indexing
- partial reindexing
- replay-safe indexing
- indexing reconciliation

### Search Analytics
- query analytics
- click analytics
- ranking analytics
- zero-result analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Search Platform MUST:
- support massive read traffic
- support distributed indexing
- support cache orchestration
- support replay-safe indexing

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The search system MUST tolerate:
- retry storms
- duplicate indexing events
- indexing lag
- partial failures
- distributed deployments
- stale cache

---

## Folder Structure

Generate:

platforms/search/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── indexing/
│   ├── ranking/
│   ├── autocomplete/
│   ├── typo/
│   ├── analytics/
│   ├── filters/
│   ├── cache/
│   ├── synchronization/
│   ├── reconciliation/
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

## Elasticsearch/OpenSearch Requirements

Use Elasticsearch/OpenSearch for:
- product indexing
- autocomplete indexes
- ranking indexes
- search aggregation
- filtering
- typo-tolerant queries

Generate:
- index mappings
- analyzers
- tokenizers
- ranking fields
- index lifecycle policies

Requirements:
- distributed indexing
- replay-safe indexing
- partial reindexing
- index versioning

Never:
- use naive indexing
- ignore shard planning
- ignore query performance
- tightly couple catalog writes to search indexing

---

## Redis Requirements

Use Redis for:
- hot query cache
- autocomplete cache
- ranking cache
- distributed invalidation coordination

Generate:
- TTL strategy
- cache invalidation strategy
- replay protection
- retry handling

Support:
- massive read traffic
- high concurrency
- distributed deployments

The cache layer MUST be production-grade.

---

## Indexing Pipeline Requirements

Support:
- distributed indexing
- partial reindexing
- async synchronization
- replay-safe indexing
- index reconciliation

Generate:
- indexing consumers
- retry-safe indexing workers
- indexing reconciliation jobs
- versioned indexing strategy

The indexing system MUST tolerate:
- duplicate events
- retry storms
- indexing lag
- partial indexing failures

No fake indexing architecture.

---

## Ranking Requirements

Support:
- relevance scoring
- popularity scoring
- freshness scoring
- sponsored ranking hooks
- personalization hooks

Generate:
- ranking pipeline
- scoring strategy
- ranking reconciliation
- ranking cache strategy

The ranking system MUST be production-grade.

---

## Autocomplete Requirements

Support:
- prefix matching
- typo tolerance
- trending suggestions
- query normalization

Generate:
- autocomplete indexing strategy
- cache strategy
- typo handling strategy

---

## Query Understanding Requirements

Support:
- synonym handling
- typo tolerance
- tokenization
- multilingual parsing hooks

Generate:
- normalization pipeline
- parsing strategy
- typo correction strategy

The query pipeline MUST be realistic.

---

## Search Analytics Requirements

Use ClickHouse for:
- query analytics
- click analytics
- ranking analytics
- search performance analytics

Generate:
- analytics ingestion flow
- replay-safe analytics ingestion
- aggregation strategy

The analytics pipeline MUST support:
- high ingestion throughput
- distributed aggregation
- replay-safe ingestion

---

## Event-Driven Requirements

Generate events for:
- indexing triggered
- index updated
- ranking updated
- autocomplete updated
- analytics ingested

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
- distributed indexing
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
- /search
- /autocomplete
- /filters
- /ranking
- /analytics

Support:
- pagination
- filtering
- sorting
- typo tolerance
- query normalization

---

## Security Requirements

The platform MUST:
- validate requests
- enforce RBAC
- sanitize input
- prevent query abuse
- isolate analytics access

Never:
- expose internal ranking metadata
- expose raw indexing internals
- trust external indexing state blindly

Generate:
- authorization middleware
- rate limiting
- replay validation
- search integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- search latency
- indexing latency
- cache hit ratio
- ranking latency
- autocomplete latency
- indexing failure count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive search analytics data.

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
- reconciliation workflows

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
- indexing tests
- ranking tests
- autocomplete tests
- replay tests
- cache tests
- concurrency tests

Test:
- indexing lag
- duplicate events
- retry storms
- stale cache
- ranking correctness
- high concurrency
- massive read traffic

---

## Output Requirements

Explain:
- search architecture
- indexing strategy
- ranking strategy
- cache orchestration
- autocomplete strategy
- analytics strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy search service.
No fake ranking system.
No naive indexing architecture.

---

## Acceptance Criteria

The Search Platform must support future integration with:
- Product Catalog Service
- Recommendation Platform
- Analytics Platform
- Advertising Platform
- Personalization Platform

without major future refactors.

The platform MUST realistically tolerate:
- indexing lag
- duplicate events
- retry storms
- stale cache
- distributed deployments
- massive read traffic

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