# TASK-014 — PRODUCT CATALOG SERVICE

## Goal

Build a REAL production-grade Product Catalog Service.

This service is responsible for:
- SKU management
- product aggregation
- seller product lifecycle
- category systems
- attribute systems
- media orchestration
- search indexing hooks
- product moderation hooks
- catalog synchronization

This is NOT a toy CRUD product service.

The Product Catalog Service must support:
- millions of products
- massive read traffic
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- fault tolerance
- replay-safe workflows

The architecture MUST prioritize:
- read scalability
- catalog consistency
- indexing correctness
- seller isolation
- operational stability

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- MySQL
- Redis Cluster
- Elasticsearch/OpenSearch
- Kafka or NATS JetStream
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- media optimization pipeline hooks

---

## Core Responsibilities

The Product Catalog Service MUST support:

### Product Management
- product creation
- product updates
- product archival
- seller product lifecycle
- product visibility management

### SKU Management
- SKU variants
- pricing metadata
- seller SKU mapping
- SKU lifecycle

### Category Systems
- category hierarchy
- category mapping
- category validation
- category metadata

### Attribute Systems
- dynamic attributes
- attribute validation
- category-specific attributes
- attribute indexing hooks

### Media Orchestration
- image metadata
- video metadata
- media optimization hooks
- CDN integration hooks

### Search Indexing Hooks
- search synchronization
- indexing events
- partial reindexing
- search consistency hooks

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Product Catalog Service MUST:
- support massive read traffic
- support distributed indexing
- support cache invalidation
- support replay-safe workflows

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The catalog system MUST tolerate:
- retry storms
- duplicate events
- indexing lag
- partial failures
- distributed deployments
- stale cache

---

## Folder Structure

Generate:

services/product-catalog/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── products/
│   ├── skus/
│   ├── categories/
│   ├── attributes/
│   ├── media/
│   ├── indexing/
│   ├── moderation/
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
- product metadata
- SKU metadata
- category hierarchy
- attribute metadata
- moderation state
- synchronization state

Generate:
- optimized schemas
- indexes
- immutable audit history
- seller isolation support

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- read optimization

Never:
- use naive product schemas
- ignore indexing strategy
- ignore category hierarchy performance
- tightly couple search engine logic

---

## Redis Requirements

Use Redis for:
- hot product cache
- category cache
- attribute cache
- search synchronization cache
- distributed invalidation coordination

Generate:
- TTL strategy
- cache invalidation strategy
- replay protection
- distributed coordination
- retry handling

Support:
- distributed deployments
- massive read traffic
- high concurrency

The cache layer MUST be production-grade.

---

## Search Indexing Requirements

Support:
- Elasticsearch/OpenSearch indexing
- partial reindexing
- distributed indexing
- eventual consistency
- replay-safe indexing

Generate:
- indexing workflows
- indexing reconciliation jobs
- retry-safe indexing consumers
- indexing versioning strategy

The indexing system MUST tolerate:
- indexing lag
- duplicate events
- retry storms
- partial indexing failures

No fake indexing architecture.

---

## Event-Driven Requirements

Generate events for:
- product created
- product updated
- product archived
- SKU updated
- category updated
- media updated
- indexing triggered

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
- distributed synchronization
- async indexing

No fake async architecture.

---

## Category Requirements

Support:
- hierarchical categories
- category validation
- category metadata
- category synchronization

Generate:
- category tree strategy
- category caching strategy
- hierarchy traversal optimization

The category system MUST be production-grade.

---

## Attribute Requirements

Support:
- dynamic attributes
- category-specific attributes
- validation rules
- indexing metadata

Generate:
- attribute validation engine
- indexing strategy
- distributed synchronization

---

## Media Requirements

Support:
- image metadata
- video metadata
- CDN integration hooks
- media optimization hooks

Generate:
- media synchronization workflows
- CDN invalidation hooks
- retry-safe media handling

The service MUST NOT store raw media directly.

Media truth belongs elsewhere.

---

## Moderation Requirements

Support:
- product moderation hooks
- seller violation hooks
- visibility control
- moderation state management

Generate:
- moderation workflows
- async moderation events
- replay-safe moderation handling

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /products
- /products/{id}
- /skus
- /categories
- /attributes
- /media

Support:
- pagination
- filtering
- sorting
- validation
- idempotency keys

---

## Security Requirements

The service MUST:
- validate ownership
- enforce RBAC
- sanitize input
- validate seller permissions
- isolate seller catalog access

Never:
- trust client indexing state
- expose internal indexing metadata
- expose internal moderation workflows

Generate:
- authorization middleware
- seller isolation validation
- replay validation
- catalog integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- product read latency
- indexing latency
- cache hit ratio
- indexing failure count
- synchronization lag
- category traversal latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive seller credentials.

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
- replay tests
- cache tests
- synchronization tests
- concurrency tests

Test:
- indexing lag
- duplicate events
- retry storms
- stale cache
- distributed synchronization
- category hierarchy performance
- high concurrency

---

## Output Requirements

Explain:
- catalog architecture
- SKU management strategy
- indexing strategy
- cache invalidation strategy
- category hierarchy strategy
- attribute validation strategy
- moderation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy product CRUD service.
No fake indexing system.
No naive category hierarchy.

---

## Acceptance Criteria

The Product Catalog Service must support future integration with:
- Search Platform
- Recommendation Platform
- Inventory Service
- Promotion Service
- Analytics Platform

without major future refactors.

The service MUST realistically tolerate:
- indexing lag
- duplicate events
- retry storms
- stale cache
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