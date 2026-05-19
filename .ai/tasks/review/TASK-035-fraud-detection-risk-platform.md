# TASK-035 — FRAUD DETECTION & RISK PLATFORM

## Goal

Build a REAL production-grade Fraud Detection & Risk Platform.

This platform is responsible for:
- fraud scoring
- realtime risk evaluation
- anti-abuse pipelines
- payment risk analysis
- account trust systems
- behavioral anomaly detection
- device fingerprinting
- transaction risk scoring
- bot mitigation
- abuse intelligence

This is NOT a toy fraud detection service.

The Fraud Platform must support:
- ultra-low latency risk scoring
- distributed abuse detection
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe risk synchronization

The architecture MUST prioritize:
- fraud-detection correctness
- scoring latency
- abuse resiliency
- operational stability
- replay safety

---

## Tech Stack

Use:
- Golang
- Python
- Kafka
- Redis Cluster
- PostgreSQL
- ClickHouse
- Elasticsearch/OpenSearch
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Flink
- Ray
- MLflow
- Triton
- Vector DB

---

## Core Responsibilities

The Fraud Platform MUST support:

### Realtime Risk Evaluation
- transaction risk scoring
- login risk scoring
- account trust evaluation
- replay-safe scoring orchestration

### Anti-Abuse Pipelines
- bot detection
- abuse throttling
- replay-safe abuse synchronization
- coordinated mitigation

### Behavioral Anomaly Detection
- account anomaly detection
- transaction anomaly detection
- behavioral scoring
- distributed anomaly orchestration

### Device Fingerprinting
- device identity evaluation
- replay-safe fingerprint synchronization
- suspicious-device detection
- distributed device reputation

### Payment Risk Analysis
- payment fraud scoring
- chargeback risk analysis
- replay-safe transaction validation
- realtime risk mitigation

### Trust Systems
- trust-score orchestration
- seller trust evaluation
- buyer trust evaluation
- distributed trust synchronization

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate scoring/mitigation/analytics
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Fraud Platform MUST:
- support replay-safe synchronization
- support distributed scoring
- support realtime mitigation
- support degraded fallback scoring

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The fraud system MUST tolerate:
- retry storms
- duplicate fraud events
- delayed scoring synchronization
- partial mitigation failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/fraud-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── scoring/
│   ├── mitigation/
│   ├── trust/
│   ├── anomalies/
│   ├── fingerprints/
│   ├── payments/
│   ├── accounts/
│   ├── synchronization/
│   ├── replay/
│   ├── orchestration/
│   ├── intelligence/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── scoring/
│   ├── anomaly/
│   ├── training/
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
- trust metadata
- fraud-case metadata
- replay metadata
- mitigation metadata
- operational state

Generate:
- optimized schemas
- indexes
- immutable fraud audit tables
- replay-safe synchronization tables

Requirements:
- replay safety
- operational correctness
- audit consistency

Never:
- mutate immutable fraud history
- tightly couple scoring to OLTP writes

---

## Redis Requirements

Use Redis for:
- hot trust cache
- replay coordination
- distributed throttling
- abuse coordination

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- ultra-low latency scoring
- high concurrency
- distributed deployments

The cache layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- fraud events
- scoring synchronization
- replay-safe ingestion
- mitigation events

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
- distributed scoring
- async synchronization
- realtime mitigation

---

## ClickHouse Requirements

Use ClickHouse for:
- fraud analytics
- anomaly analytics
- abuse analytics
- trust analytics
- scoring analytics

Generate:
- aggregation strategy
- partitioning strategy
- TTL policies
- materialized views

Requirements:
- replay-safe aggregation
- high-ingestion throughput
- realtime analytics

---

## Elasticsearch/OpenSearch Requirements

Use Elasticsearch/OpenSearch for:
- threat intelligence lookup
- device fingerprint indexing
- abuse search
- fraud investigation

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

## Realtime Scoring Requirements

Support:
- transaction scoring
- login scoring
- behavioral scoring
- replay-safe risk evaluation

Generate:
- scoring orchestration
- distributed evaluation pipelines
- replay-safe synchronization
- fallback scoring workflows

The scoring system MUST tolerate:
- delayed signals
- replay storms
- partial inference failures

No fake scoring architecture.

---

## Mitigation Requirements

Support:
- abuse throttling
- account freezing
- transaction blocking
- replay-safe mitigation execution

Generate:
- mitigation orchestration
- distributed enforcement workflows
- replay-safe synchronization
- degraded mitigation fallback

No naive mitigation logic.

---

## Device Fingerprinting Requirements

Support:
- device reputation
- suspicious-device detection
- replay-safe fingerprint synchronization
- distributed fingerprint evaluation

Generate:
- fingerprint pipelines
- replay-safe coordination
- distributed enrichment workflows
- fingerprint reconciliation

The fingerprinting system MUST be realistic.

---

## Behavioral Anomaly Requirements

Support:
- anomaly scoring
- behavioral baselines
- replay-safe anomaly detection
- distributed anomaly orchestration

Generate:
- anomaly pipelines
- distributed evaluation
- replay-safe synchronization
- realtime enrichment workflows

---

## ML Risk Evaluation Requirements

Support:
- fraud scoring models
- anomaly detection models
- behavioral risk models
- realtime inference

Generate:
- inference orchestration
- replay-safe feature retrieval
- distributed inference pipelines
- fallback inference workflows

The ML system MUST tolerate:
- feature outages
- replay storms
- inference degradation
- partial node failures

No fake ML fraud architecture.

---

## Event-Driven Requirements

Generate events for:
- fraud detected
- mitigation triggered
- anomaly detected
- trust score updated
- suspicious device detected
- transaction blocked

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
- /fraud/score
- /fraud/trust
- /fraud/mitigation
- /fraud/anomalies
- /fraud/fingerprints

Support:
- realtime scoring
- replay-safe synchronization
- distributed mitigation
- operational auditing

---

## Security Requirements

The platform MUST:
- validate operational ownership
- enforce RBAC
- isolate mitigation workflows
- sanitize scoring input
- protect trust integrity

Never:
- expose fraud internals publicly
- expose raw scoring signals
- trust external scoring blindly
- allow unrestricted mitigation execution

Generate:
- authorization middleware
- replay validation
- mitigation isolation
- operational audit pipelines

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- fraud-detection latency
- mitigation latency
- anomaly throughput
- replay latency
- scoring throughput
- abuse-block rate

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive fraud signals.

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
- Deployments
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
- GitHub Actions or Drone pipelines
- model validation
- Helm validation
- Kubernetes policy validation
- vulnerability scanning
- GitOps workflows

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- replay tests
- mitigation tests
- anomaly tests
- concurrency tests
- scoring tests

Test:
- duplicate fraud events
- replay storms
- inference degradation
- delayed scoring signals
- ultra-high scoring QPS
- distributed deployments
- mitigation correctness

---

## Output Requirements

Explain:
- fraud-platform architecture
- scoring strategy
- mitigation strategy
- replay-safe synchronization strategy
- anomaly-detection strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy fraud service.
No fake mitigation architecture.
No naive fraud logic.

---

## Acceptance Criteria

The Fraud Platform must support future integration with:
- Payment Platform
- AI/ML Platform
- Analytics Platform
- Advertising Platform
- Live Commerce Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate fraud events
- inference degradation
- distributed deployments
- ultra-high scoring throughput

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.