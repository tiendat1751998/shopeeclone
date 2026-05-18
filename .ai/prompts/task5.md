# TASK-005 — PRODUCT SERVICE

## Goal

Build a REAL production-grade Product Service.

This service is responsible for:
- product catalog
- SKU management
- product variants
- pricing metadata
- category hierarchy
- seller product management
- product attributes
- product moderation
- media metadata
- inventory integration
- search indexing events

This is NOT a toy CRUD product service.

The service must support:
- millions of products
- high read throughput
- horizontal scaling
- Kubernetes-native deployment
- distributed tracing
- event-driven architecture
- observability-first design
- fault tolerance

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- MySQL
- Redis Cluster
- OpenSearch
- NATS JetStream or Kafka
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

---

## Core Responsibilities

The Product Service MUST support:

### Product Catalog
- create product
- update product
- soft delete product
- product visibility
- product moderation status
- product publishing workflow

### SKU Management
- SKU creation
- SKU variants
- variant combinations
- SKU metadata
- SKU pricing metadata
- SKU dimensions
- SKU media metadata

### Category Management
- hierarchical categories
- category inheritance
- category attributes
- category metadata

### Seller Product Management
- seller-owned products
- seller validation
- ownership validation
- seller product quotas

### Product Attributes
- dynamic attributes
- searchable attributes
- filterable attributes
- category-specific attributes

### Media Metadata
- image metadata
- video metadata
- CDN metadata
- asset references

### Search Integration
- search indexing events
- product update events
- search synchronization
- eventual consistency

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support event-driven workflows
- support eventual consistency
- support distributed deployments

The service MUST:
- be read-optimized
- support high concurrency
- support horizontal scaling
- support caching
- support async indexing

Use:
- dependency injection
- modular architecture
- service isolation
- CQRS where appropriate

---

## Folder Structure

Generate:

services/product/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── product/
│   ├── sku/
│   ├── category/
│   ├── attributes/
│   ├── moderation/
│   ├── pricing/
│   ├── media/
│   ├── search/
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
- products
- SKUs
- category hierarchy
- product attributes
- moderation state
- seller ownership
- media metadata

Generate:
- migrations
- indexes
- constraints
- optimized schemas

Requirements:
- pagination everywhere
- indexed lookups
- optimized joins
- soft delete support
- audit timestamps

Support:
- large product catalogs
- read-heavy workloads
- eventual consistency

NEVER:
- use SELECT *
- ignore indexing
- ignore query optimization
- tightly couple product and inventory schemas

---

## Redis Requirements

Use Redis for:
- hot product cache
- category cache
- attribute cache
- product metadata cache
- seller cache
- query acceleration

Generate:
- cache invalidation strategy
- cache warming strategy
- TTL strategy
- retry handling
- connection pooling

Support:
- distributed deployments
- cache consistency
- stale cache handling

---

## OpenSearch Requirements

Use OpenSearch for:
- product search
- filtering
- faceting
- sorting
- autocomplete

Generate:
- indexing strategy
- sync strategy
- event-driven indexing
- index mappings
- query optimization

Support:
- eventual consistency
- async indexing
- replay-safe indexing

The search architecture MUST be realistic.

No fake search abstraction.

---

## Event-Driven Requirements

Generate events for:
- product created
- product updated
- product deleted
- SKU updated
- category updated
- moderation updated

Use:
- Kafka or NATS JetStream

Requirements:
- retries
- DLQ
- idempotency
- replay-safe consumers
- event versioning
- consumer groups

The async architecture MUST be production-grade.

---

## Moderation Requirements

Support:
- moderation workflow
- pending approval
- rejection reasons
- automated moderation hooks
- manual moderation support

Generate:
- moderation states
- audit trails
- moderation events

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /products
- /skus
- /categories
- /attributes
- /seller/products
- /moderation

Support:
- pagination
- filtering
- sorting
- search
- validation

---

## Security Requirements

The service MUST:
- validate ownership
- enforce RBAC
- sanitize input
- validate seller identity
- prevent unauthorized updates

Never:
- trust seller-provided ownership
- expose internal metadata improperly
- expose moderation internals

Generate:
- authorization middleware
- ownership validation
- moderation validation

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
- cache hit ratio
- DB latency
- indexing latency
- search sync lag
- event retry count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive metadata.

---

## Reliability Requirements

Implement:
- retries
- timeout handling
- graceful shutdown
- panic recovery
- circuit breakers

Support:
- rolling deployment
- autoscaling
- distributed deployments

Generate:
- retry policies
- resilience middleware
- backoff strategies

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
- cache tests
- OpenSearch tests
- event tests
- RBAC tests
- moderation tests
- API tests

Test:
- high concurrency
- cache invalidation
- event replay safety
- indexing consistency
- pagination correctness

---

## Output Requirements

Explain:
- product architecture
- SKU strategy
- category hierarchy strategy
- search indexing flow
- OpenSearch strategy
- cache strategy
- scaling strategy
- event-driven flow
- Kubernetes deployment strategy
- observability flow

Generate production-grade code only.

No toy CRUD service.
No fake search implementation.
No fake async architecture.

---

## Acceptance Criteria

The Product Service must support future integration with:
- Inventory Service
- Cart Service
- Order Service
- Recommendation System
- Search System
- Seller Center

without major future refactors.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.