# TASK-029 — ANALYTICS & BI PLATFORM

## Goal

Build a REAL production-grade Analytics & BI Platform.

This platform is responsible for:
- realtime analytics
- OLAP pipelines
- dashboards
- data lake ingestion
- stream aggregation
- business intelligence
- metrics computation
- behavioral analytics
- operational analytics
- executive reporting

This is NOT a toy analytics service.

The Analytics Platform must support:
- ultra-high event ingestion
- distributed stream processing
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe aggregation

The architecture MUST prioritize:
- aggregation correctness
- replay safety
- query performance
- operational stability
- analytics freshness

---

## Tech Stack

Use:
- Golang
- Kafka
- ClickHouse
- PostgreSQL
- Redis Cluster
- Apache Flink
- Apache Spark
- Apache Iceberg
- MinIO or S3
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Trino
- Presto
- dbt
- Airflow

---

## Core Responsibilities

The Analytics Platform MUST support:

### Realtime Analytics
- realtime dashboards
- operational metrics
- behavioral metrics
- business KPIs

### OLAP Processing
- large-scale aggregations
- analytical queries
- distributed OLAP
- replay-safe computation

### Stream Aggregation
- event aggregation
- windowed aggregation
- replay-safe stream processing
- out-of-order event handling

### Data Lake Ingestion
- raw event ingestion
- immutable event storage
- partitioned storage
- historical replay support

### Business Intelligence
- executive dashboards
- product analytics
- seller analytics
- monetization analytics

### Behavioral Analytics
- user behavior aggregation
- conversion analytics
- retention analytics
- engagement analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate ingestion/aggregation/querying
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Analytics Platform MUST:
- support replay-safe aggregation
- support delayed event handling
- support out-of-order events
- support distributed analytics querying

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The analytics system MUST tolerate:
- retry storms
- duplicate ingestion events
- delayed events
- partial aggregation failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/analytics/
├── cmd/
├── internal/
│   ├── config/
│   ├── ingestion/
│   ├── aggregation/
│   ├── stream/
│   ├── olap/
│   ├── dashboards/
│   ├── reporting/
│   ├── behavioral/
│   ├── metrics/
│   ├── synchronization/
│   ├── replay/
│   ├── lakehouse/
│   ├── cache/
│   ├── events/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── flink/
├── spark/
├── iceberg/
├── sql/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Kafka Requirements

Use Kafka for:
- event ingestion
- stream synchronization
- replay-safe ingestion
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
- distributed ingestion
- high throughput streaming
- async aggregation

---

## ClickHouse Requirements

Use ClickHouse for:
- OLAP queries
- realtime dashboards
- analytical aggregations
- behavioral analytics
- operational analytics

Generate:
- partitioning strategy
- sharding strategy
- materialized views
- TTL policies
- aggregation pipelines

Requirements:
- replay-safe aggregation
- high ingestion throughput
- low-latency OLAP querying

Never:
- use naive aggregation tables
- ignore partition pruning
- ignore retention strategies

---

## Apache Flink Requirements

Use Flink for:
- stream aggregation
- windowed processing
- out-of-order event handling
- replay-safe stream processing

Generate:
- stream pipelines
- watermark strategies
- checkpointing
- replay-safe state management

Requirements:
- exactly-once semantics
- distributed processing
- replay-safe computation

The stream system MUST tolerate:
- delayed events
- replay storms
- partial failures

No fake streaming architecture.

---

## Apache Spark Requirements

Use Spark for:
- batch analytics
- historical aggregation
- replay processing
- ML-oriented analytics hooks

Generate:
- batch pipelines
- replay-safe batch jobs
- distributed historical aggregation
- partition-aware processing

---

## Iceberg/Data Lake Requirements

Use Iceberg + MinIO/S3 for:
- immutable event storage
- historical replay
- partitioned lakehouse storage
- analytics backfills

Generate:
- partition strategies
- retention policies
- replay orchestration
- schema evolution workflows

Requirements:
- immutable raw events
- replay-safe ingestion
- schema evolution support

The lakehouse MUST be production-grade.

---

## Dashboard Requirements

Support:
- realtime dashboards
- operational dashboards
- executive dashboards
- replay-safe metrics refresh

Generate:
- dashboard APIs
- aggregation orchestration
- caching workflows
- distributed refresh pipelines

The dashboard system MUST tolerate:
- delayed aggregations
- replay storms
- high query concurrency

---

## Reporting Requirements

Support:
- scheduled reports
- ad-hoc reports
- seller reports
- executive reports

Generate:
- reporting pipelines
- distributed report generation
- replay-safe reporting
- export orchestration

---

## Behavioral Analytics Requirements

Support:
- retention analytics
- funnel analytics
- conversion analytics
- engagement analytics

Generate:
- behavioral aggregation pipelines
- replay-safe behavioral metrics
- distributed enrichment workflows
- analytics reconciliation

The behavioral analytics system MUST be realistic.

---

## Query Layer Requirements

Support:
- distributed OLAP querying
- replay-safe querying
- dashboard queries
- large-scale analytical queries

Generate:
- query orchestration
- query caching
- fallback querying
- query isolation

The query layer MUST tolerate:
- high concurrency
- partial failures
- delayed partitions

---

## Event-Driven Requirements

Generate events for:
- aggregation completed
- replay triggered
- dashboard refreshed
- ingestion delayed
- analytics anomaly detected
- report generated

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
- distributed aggregation
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
- /analytics/dashboards
- /analytics/reports
- /analytics/metrics
- /analytics/query
- /analytics/behavioral

Support:
- realtime dashboards
- analytical querying
- replay-safe aggregation
- pagination/filtering

---

## Security Requirements

The platform MUST:
- validate analytics access
- enforce RBAC
- isolate tenant analytics
- sanitize query input
- protect data integrity

Never:
- expose raw internal events
- expose sensitive PII
- trust external aggregation blindly
- allow unrestricted OLAP scans

Generate:
- authorization middleware
- replay validation
- query isolation
- analytics auditing

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
- dashboard latency
- query concurrency
- delayed-event count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive analytics payloads.

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
- StatefulSets
- Services
- ConfigMap
- Secret integration
- HPA
- PodDisruptionBudget
- ServiceMonitor
- NetworkPolicy
- Helm charts

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
- streaming tests
- concurrency tests
- delayed-event tests

Test:
- duplicate ingestion events
- replay storms
- out-of-order events
- delayed aggregations
- ultra-high ingestion throughput
- distributed deployments
- analytical query correctness

---

## Output Requirements

Explain:
- analytics architecture
- ingestion strategy
- aggregation strategy
- replay-safe aggregation strategy
- lakehouse strategy
- stream-processing strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy analytics service.
No fake streaming architecture.
No naive OLAP logic.

---

## Acceptance Criteria

The Analytics Platform must support future integration with:
- Search Platform
- Recommendation Platform
- Advertising Platform
- Fraud Platform
- Live Commerce Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate ingestion events
- delayed events
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