# TASK-018 — USER BEHAVIOR PLATFORM

## Goal

Build a REAL production-grade User Behavior Platform.

This platform is responsible for:
- clickstream ingestion
- event tracking
- session analytics
- realtime analytics
- behavioral aggregation
- event pipelines
- user activity streams
- replay-safe ingestion
- downstream synchronization

This is NOT a toy analytics service.

The User Behavior Platform must support:
- billions of events
- ultra-high ingestion throughput
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe ingestion

The architecture MUST prioritize:
- ingestion throughput
- replay correctness
- aggregation correctness
- analytics freshness
- operational stability

---

## Tech Stack

Use:
- Golang
- Kafka
- ClickHouse
- Redis Cluster
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Apache Flink
- Apache Spark
- Apache Pinot
- Kafka Streams

---

## Core Responsibilities

The User Behavior Platform MUST support:

### Clickstream Ingestion
- page views
- product views
- clicks
- impressions
- add-to-cart events
- checkout events

### Session Analytics
- session tracking
- session aggregation
- user activity timelines
- engagement tracking

### Realtime Analytics
- realtime counters
- trending calculations
- live dashboards
- realtime aggregation

### Event Pipelines
- distributed ingestion
- replay-safe processing
- event enrichment
- event partitioning

### Behavioral Aggregation
- user behavior aggregation
- product analytics
- seller analytics
- category analytics

### Downstream Synchronization
- recommendation pipelines
- search analytics
- fraud analytics
- BI pipelines

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate ingestion/aggregation/storage
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The User Behavior Platform MUST:
- support massive ingestion throughput
- support replay-safe ingestion
- support distributed aggregation
- support late-event handling

Use:
- partitioned ingestion
- distributed stream processing
- scalable aggregation pipelines
- resilience patterns

The analytics system MUST tolerate:
- duplicate events
- retry storms
- late-arriving events
- partial failures
- distributed deployments
- out-of-order events

---

## Folder Structure

Generate:

platforms/user-behavior/
├── cmd/
├── internal/
│   ├── config/
│   ├── ingestion/
│   ├── aggregation/
│   ├── analytics/
│   ├── sessions/
│   ├── enrichment/
│   ├── partitioning/
│   ├── synchronization/
│   ├── replay/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── stream/
│   ├── consumers/
│   ├── processors/
│   ├── aggregators/
│   ├── replay/
│   └── pipelines/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Kafka Requirements

Use Kafka for:
- clickstream ingestion
- event buffering
- distributed processing
- replay-safe pipelines

Generate:
- topic strategy
- partitioning strategy
- retention policies
- DLQ topics
- replay pipelines

Requirements:
- replay-safe ingestion
- idempotent producers
- consumer groups
- event versioning

Support:
- billions of events
- high fan-in ingestion
- distributed consumers

Never:
- use naive event ingestion
- ignore partitioning
- ignore replay correctness

---

## ClickHouse Requirements

Use ClickHouse for:
- realtime analytics
- session analytics
- behavioral aggregation
- trending analytics
- low-cost OLAP queries

Generate:
- partitioning strategy
- aggregation strategy
- TTL policies
- materialized views

Requirements:
- high ingestion throughput
- distributed aggregation
- replay-safe aggregation

Never:
- use naive analytics schemas
- ignore partition pruning
- ignore aggregation optimization

---

## Session Analytics Requirements

Support:
- session tracking
- activity timelines
- engagement scoring
- realtime session updates

Generate:
- session aggregation pipelines
- session timeout handling
- replay-safe session updates

The session system MUST tolerate:
- duplicate events
- out-of-order events
- late-arriving events

---

## Realtime Aggregation Requirements

Support:
- realtime counters
- trending analytics
- engagement aggregation
- product popularity tracking

Generate:
- distributed aggregation pipelines
- replay-safe aggregation
- cache synchronization

The aggregation system MUST support:
- high ingestion throughput
- distributed workers
- replay-safe recalculation

---

## Event Enrichment Requirements

Support:
- user enrichment hooks
- product enrichment hooks
- geo enrichment hooks
- device enrichment hooks

Generate:
- enrichment pipelines
- replay-safe enrichment
- fallback enrichment handling

---

## Replay Requirements

Support:
- replay-safe ingestion
- event reprocessing
- backfill pipelines
- replay isolation

Generate:
- replay orchestration
- replay-safe aggregation
- replay reconciliation workflows

The replay system MUST be production-grade.

No fake replay architecture.

---

## Downstream Synchronization Requirements

Support:
- recommendation synchronization
- search analytics synchronization
- fraud analytics synchronization
- BI synchronization

Generate:
- distributed synchronization workflows
- retry-safe synchronization
- replay-safe downstream delivery

---

## Redis Requirements

Use Redis for:
- hot counters
- trending cache
- realtime aggregation cache
- distributed coordination

Generate:
- TTL strategy
- replay protection
- distributed coordination
- cache invalidation

Support:
- realtime reads
- distributed deployments
- high concurrency

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /events
- /analytics/trending
- /analytics/realtime
- /sessions
- /analytics/products

Support:
- pagination
- filtering
- aggregation queries
- realtime queries

---

## Security Requirements

The platform MUST:
- validate ingestion requests
- enforce RBAC
- sanitize input
- isolate analytics access
- protect event integrity

Never:
- expose raw behavioral streams
- expose internal aggregation metadata
- trust external event timestamps blindly

Generate:
- authorization middleware
- replay validation
- ingestion integrity validation
- analytics isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- ingestion throughput
- aggregation latency
- replay latency
- duplicate event count
- late-event count
- ClickHouse query latency

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
- aggregation tests
- stream processing tests
- concurrency tests
- late-event tests

Test:
- duplicate events
- retry storms
- replay correctness
- out-of-order events
- late-arriving events
- high ingestion throughput
- distributed aggregation

---

## Output Requirements

Explain:
- ingestion architecture
- replay strategy
- aggregation strategy
- ClickHouse strategy
- session analytics strategy
- synchronization strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy analytics service.
No fake replay architecture.
No naive aggregation pipelines.

---

## Acceptance Criteria

The User Behavior Platform must support future integration with:
- Recommendation Platform
- Search Platform
- Fraud Platform
- BI Platform
- Marketing Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate events
- late-arriving events
- distributed deployments
- ultra-high ingestion throughput

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