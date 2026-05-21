# TASK-019 — FRAUD DETECTION PLATFORM

## Goal

Build a REAL production-grade Fraud Detection Platform.

This platform is responsible for:
- fraud scoring
- anomaly detection
- behavioral risk analysis
- realtime fraud pipelines
- payment abuse detection
- account abuse detection
- device fingerprinting hooks
- rule engines
- ML fraud scoring hooks
- investigation workflows

This is NOT a toy fraud service.

The Fraud Detection Platform must support:
- realtime scoring
- ultra-high event throughput
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe evaluation

The architecture MUST prioritize:
- fraud detection accuracy
- low-latency scoring
- replay correctness
- operational stability
- false-positive minimization

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
- Flink
- Feast
- MLflow
- vector similarity hooks

---

## Core Responsibilities

The Fraud Detection Platform MUST support:

### Fraud Scoring
- realtime risk scoring
- rule-based scoring
- behavioral scoring
- payment risk scoring
- account risk scoring

### Anomaly Detection
- abnormal transaction detection
- login anomaly detection
- device anomaly detection
- geographic anomaly detection

### Behavioral Risk Analysis
- session analysis
- behavioral fingerprinting
- velocity checks
- abuse pattern detection

### Rule Engine
- dynamic fraud rules
- configurable thresholds
- replay-safe evaluation
- distributed rule execution

### Device Fingerprinting Hooks
- device metadata hooks
- browser fingerprint hooks
- IP reputation hooks
- geo risk hooks

### Investigation Workflows
- fraud investigation queues
- manual review hooks
- evidence aggregation
- fraud audit trails

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate ingestion/scoring/storage
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Fraud Detection Platform MUST:
- support realtime scoring
- support replay-safe evaluation
- support distributed rule execution
- support low-latency blocking

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The fraud system MUST tolerate:
- retry storms
- duplicate events
- delayed events
- partial failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/fraud/
├── cmd/
├── internal/
│   ├── config/
│   ├── ingestion/
│   ├── scoring/
│   ├── rules/
│   ├── anomaly/
│   ├── behavioral/
│   ├── fingerprinting/
│   ├── investigations/
│   ├── evidence/
│   ├── synchronization/
│   ├── replay/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── scoring/
│   ├── models/
│   ├── training/
│   ├── datasets/
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
- fraud event ingestion
- realtime scoring pipelines
- replay-safe evaluation
- distributed fraud processing

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
- realtime scoring
- distributed workers
- high event throughput

Never:
- use naive scoring pipelines
- ignore replay correctness
- ignore partition strategy

---

## PostgreSQL Requirements

Use PostgreSQL for:
- fraud cases
- investigation metadata
- fraud rules
- audit trails
- evidence metadata

Generate:
- optimized schemas
- indexes
- immutable audit history
- evidence integrity support

Requirements:
- transactional correctness
- replay safety
- investigation consistency

Never:
- mutate immutable fraud audit history
- tightly couple rule execution to storage

---

## Redis Requirements

Use Redis for:
- hot fraud scores
- realtime blocking cache
- replay coordination
- distributed throttling

Generate:
- TTL strategy
- distributed coordination
- replay protection
- rate limiting

Support:
- realtime scoring
- high concurrency
- distributed deployments

The coordination layer MUST be production-grade.

---

## ClickHouse Requirements

Use ClickHouse for:
- fraud analytics
- anomaly analytics
- risk aggregation
- behavioral analytics

Generate:
- aggregation strategy
- partitioning strategy
- TTL policies
- materialized views

Requirements:
- replay-safe aggregation
- high ingestion throughput
- distributed analytics

---

## Rule Engine Requirements

Support:
- dynamic rules
- configurable thresholds
- realtime rule execution
- replay-safe evaluation

Generate:
- rule execution engine
- distributed evaluation pipelines
- rule versioning
- rollback-safe rule deployment

The rule engine MUST tolerate:
- duplicate events
- replay storms
- partial evaluation failures

No fake rule engine.

---

## ML Scoring Requirements

Support:
- ML fraud scoring hooks
- behavioral models
- anomaly models
- realtime inference hooks

Generate:
- scoring pipelines
- feature retrieval hooks
- fallback scoring
- inference orchestration

Support:
- degraded inference mode
- partial feature outages
- low-latency scoring

No naive ML integration.

---

## Device Fingerprinting Requirements

Support:
- browser fingerprint hooks
- IP reputation hooks
- geo risk hooks
- device correlation hooks

Generate:
- fingerprint aggregation
- replay-safe fingerprint updates
- risk enrichment pipelines

The fingerprinting system MUST be realistic.

---

## Investigation Requirements

Support:
- fraud review queues
- manual review workflows
- evidence aggregation
- fraud audit trails

Generate:
- investigation orchestration
- replay-safe evidence handling
- audit integrity workflows

---

## Event-Driven Requirements

Generate events for:
- fraud scored
- anomaly detected
- account flagged
- payment blocked
- investigation opened
- rule updated

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
- distributed scoring
- async investigation workflows

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /fraud/score
- /fraud/anomalies
- /fraud/investigations
- /fraud/rules
- /fraud/decisions

Support:
- pagination
- filtering
- replay-safe scoring
- realtime decisions

---

## Security Requirements

The platform MUST:
- validate ingestion requests
- enforce RBAC
- sanitize input
- isolate fraud investigations
- protect scoring integrity

Never:
- expose raw ML models
- expose internal scoring rules
- expose investigation evidence
- trust external fraud signals blindly

Generate:
- authorization middleware
- replay validation
- fraud integrity validation
- investigation isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- fraud scoring latency
- anomaly detection latency
- rule execution latency
- false-positive rate
- replay latency
- fraud blocking count

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive fraud evidence.

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
- scoring tests
- anomaly tests
- rule-engine tests
- concurrency tests

Test:
- duplicate events
- replay storms
- delayed events
- realtime blocking
- rule rollback
- high event throughput
- distributed scoring

---

## Output Requirements

Explain:
- fraud architecture
- scoring strategy
- rule engine strategy
- anomaly strategy
- ML scoring strategy
- replay strategy
- investigation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy fraud service.
No fake rule engine.
No naive realtime scoring architecture.

---

## Acceptance Criteria

The Fraud Detection Platform must support future integration with:
- Payment Platform
- User Behavior Platform
- Recommendation Platform
- Order Service
- Risk Operations Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate events
- delayed events
- distributed deployments
- realtime fraud scoring traffic

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