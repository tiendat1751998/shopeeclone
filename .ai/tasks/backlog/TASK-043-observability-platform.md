# TASK-043 — OBSERVABILITY PLATFORM

## Goal

Build a REAL production-grade Observability Platform.

This platform is responsible for:
- centralized logs ingestion & processing
- distributed tracing storage & query
- metrics aggregation & time-series system
- audit event data lake
- alerting pipeline (integration layer)
- anomaly detection hooks
- correlation across logs/metrics/traces
- retention & compliance policies
- high-throughput telemetry ingestion

This is NOT a monitoring dashboard setup.

The Observability Platform is:
## SYSTEM-WIDE TELEMETRY BACKBONE

The system MUST support:
- ultra-high throughput ingestion pipelines
- distributed telemetry correlation
- scalable storage backend
- Kubernetes-native deployment
- observability-first architecture

The architecture MUST prioritize:
- ingestion reliability
- query performance
- cost efficiency
- trace completeness
- system-wide correlation correctness

---

## Tech Stack

Use:
- Golang (ingestion + query services)
- Kafka (telemetry streaming backbone)
- ClickHouse (primary analytics store)
- Elasticsearch / OpenSearch (log indexing)
- Prometheus (metrics ingestion/compat layer)
- Grafana (visualization layer integration)
- OpenTelemetry (standard instrumentation format)
- S3-compatible storage (long-term archive)

Optional:
- Loki (log aggregation alternative)
- Tempo (trace backend alternative)
- Thanos (Prometheus scaling layer)
- Parquet + Iceberg (data lake storage layer)

---

## Core Responsibilities

The Observability Platform MUST support:

### Logs Platform
- structured log ingestion (JSON)
- log parsing & normalization
- indexing & search
- log retention policies

### Metrics Platform
- time-series ingestion
- high-cardinality handling
- aggregation pipelines
- long-term storage scaling

### Distributed Tracing Platform
- trace ingestion (OpenTelemetry)
- span correlation
- trace reconstruction
- cross-service request mapping

### Audit Data Lake
- immutable audit event storage
- compliance-grade retention
- replayable event history

### Alerting Pipeline (Integration Layer)
- alert rule ingestion
- anomaly detection hooks
- integration with notification system

---

## Architecture Requirements

The Observability Platform MUST:
- decouple ingestion, processing, and storage layers
- support backpressure handling
- support multi-region ingestion
- support eventual consistency for analytics

The system MUST:
- be horizontally scalable
- handle burst telemetry spikes
- degrade gracefully under storage overload

Use:
- streaming pipelines (Kafka-based)
- batch + streaming hybrid processing
- schema versioning for telemetry

The system MUST tolerate:
- log storms (incident spikes)
- metric cardinality explosion
- trace ingestion overload
- storage backend slowdowns
- Kafka lag under load

---

## Folder Structure

Generate:

platforms/observability-platform/
├── ingestion/
│   ├── logs/
│   ├── metrics/
│   ├── traces/
│   ├── audit/
│   └── normalization/
│
├── processing/
│   ├── pipeline/
│   ├── aggregation/
│   ├── enrichment/
│   ├── correlation/
│   └── sampling/
│
├── storage/
│   ├── clickhouse/
│   ├── elastic/
│   ├── s3/
│   ├── schema/
│   └── retention/
│
├── query/
│   ├── logs/
│   ├── metrics/
│   ├── traces/
│   ├── api/
│   └── caching/
│
├── alerting/
├── rules/
├── dashboards/
├── kafka/
├── configs/
├── deployments/
├── charts/
├── tests/
└── Dockerfile

---

## Logs Requirements (CRITICAL)

Support:
- structured JSON logs
- log normalization
- log enrichment (trace_id, service_id)
- high-throughput ingestion

Generate:
- log ingestion pipeline
- indexing strategy
- query API

Must support:
- millisecond-level search latency (indexed paths)

---

## Metrics Requirements (CRITICAL)

Support:
- Prometheus-compatible ingestion
- time-series aggregation
- histogram + summary metrics
- high-cardinality optimization

Generate:
- metrics ingestion pipeline
- aggregation workers
- long-term storage strategy

Must handle:
## HIGH CARDINALITY WITHOUT EXPLOSION FAILURE

---

## Tracing Requirements (CRITICAL)

Support:
- OpenTelemetry ingestion
- span correlation
- distributed trace reconstruction
- service dependency graph

Generate:
- trace collector
- span processor
- trace query API

Must support:
- full request lifecycle reconstruction across all services

---

## Audit Data Lake Requirements

Support:
- immutable event storage
- compliance retention
- replay capability for investigations

Generate:
- append-only storage pipeline
- S3 archival strategy
- query engine integration

Critical:
## AUDIT DATA MUST NEVER BE MUTATED OR DELETED

---

## Kafka Requirements

Use Kafka for:
- logs stream
- metrics stream
- traces stream
- audit stream
- alert stream

Generate:
- topic design
- partition strategy
- replay-safe ingestion
- DLQ handling

---

## ClickHouse Requirements

Use ClickHouse for:
- log analytics
- metrics aggregation
- trace queries
- business observability insights

Generate:
- schema design
- partitioning strategy
- materialized views
- aggregation pipelines

---

## Elasticsearch Requirements

Use Elasticsearch for:
- log search
- structured log filtering
- debugging queries

Generate:
- index templates
- shard strategy
- retention policies

---

## Query API Requirements

Support:
- logs query API
- metrics query API
- trace query API
- correlation search API

Must support:
- low latency debugging queries
- cross-service trace lookup
- time-range filtering

---

## Alerting Requirements

Support:
- threshold-based alerts
- anomaly detection hooks
- multi-channel integration

Generate:
- alert rules engine
- evaluation pipeline
- integration hooks to Notification Platform

---

## Retention & Cost Control

Support:
- tiered storage (hot/warm/cold)
- retention policies per signal type
- sampling for high-volume traces

Generate:
- retention engine
- sampling strategy
- archival pipeline

---

## Security Requirements

The platform MUST:
- protect sensitive logs
- enforce RBAC on telemetry access
- isolate tenant data
- mask sensitive fields

Never:
- expose raw secrets in logs
- allow unrestricted query access
- bypass audit restrictions

---

## Observability of Observability (META REQUIREMENT)

The system MUST monitor itself:
- ingestion lag
- pipeline failures
- storage latency
- query performance

Generate:
- self-monitoring metrics
- health dashboards

---

## Reliability Requirements

Implement:
- retries (safe ingestion only)
- backpressure handling
- circuit breakers
- graceful degradation

Critical:
## TELEMETRY LOSS MUST BE MINIMIZED UNDER ALL CONDITIONS

---

## Kubernetes Requirements

Generate:
- Deployments
- StatefulSets (ClickHouse/Elastic)
- HPA
- PDB
- ConfigMaps
- Secrets
- Helm charts

Must support:
- high ingestion bursts
- storage scaling
- zero downtime upgrades

---

## CI/CD Requirements

Generate:
- ingestion pipeline tests
- schema validation tests
- load tests for telemetry spikes
- retention policy validation
- GitOps deployment

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- ingestion stress tests
- query latency tests
- failure simulation tests

Test scenarios:
- log storm during outage
- trace explosion under flash sale
- metric cardinality spike
- Kafka lag accumulation
- Elasticsearch degradation
- ClickHouse slow queries

---

## Output Requirements

Explain:
- observability architecture
- ingestion pipeline design
- storage strategy
- trace correlation model
- metrics scaling strategy
- log indexing strategy
- cost optimization strategy

Generate production-grade code only.

No toy logging system.
No fake metrics pipeline.
No simplified tracing model.

---

## Acceptance Criteria

The Observability Platform must support:
- API Gateway
- OMS Platform
- Payment Platform
- Recommendation Platform
- Notification Platform
- Service Mesh

without redesign.

The system MUST survive:
- log storms
- metric explosions
- trace overload
- storage degradation
- Kafka lag

WITHOUT losing observability continuity.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
Observability is a system-critical dependency.