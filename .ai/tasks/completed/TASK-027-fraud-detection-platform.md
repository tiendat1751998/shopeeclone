# TASK-027 — FRAUD DETECTION PLATFORM

## Goal

Build a REAL production-grade Fraud Detection Platform.

This platform is responsible for:
- risk scoring
- fraud analytics
- anomaly detection
- behavioral risk analysis
- payment fraud prevention
- account abuse detection
- device fingerprinting
- realtime rule evaluation
- trust orchestration
- abuse mitigation

This is NOT a toy fraud service.

The Fraud Detection Platform must support:
- realtime risk scoring
- ultra-low latency decisioning
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe behavioral ingestion

The architecture MUST prioritize:
- scoring correctness
- behavioral freshness
- adversarial resiliency
- operational stability
- replay safety

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
- Ray
- feature stores
- graph analysis engines

---

## Core Responsibilities

The Fraud Detection Platform MUST support:

### Risk Scoring
- realtime risk scoring
- behavioral scoring
- account trust scoring
- transaction risk evaluation

### Fraud Analytics
- abuse analytics
- anomaly aggregation
- fraud trend analysis
- attack pattern analysis

### Behavioral Analysis
- behavioral profiling
- activity anomaly detection
- login behavior analysis
- transaction behavior analysis

### Payment Fraud Prevention
- suspicious payment detection
- chargeback risk analysis
- settlement abuse detection
- payout fraud detection

### Account Abuse Detection
- fake account detection
- spam account detection
- account farming detection
- bot activity detection

### Device Fingerprinting
- device correlation
- session fingerprinting
- multi-account linkage
- suspicious device detection

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate scoring/rules/analytics
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Fraud Detection Platform MUST:
- support replay-safe behavioral ingestion
- support realtime scoring
- support distributed rule evaluation
- support degraded scoring fallback

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The fraud system MUST tolerate:
- retry storms
- duplicate behavioral events
- delayed analytics updates
- partial scoring failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/fraud-detection/
├── cmd/
├── internal/
│   ├── config/
│   ├── scoring/
│   ├── rules/
│   ├── analytics/
│   ├── behavioral/
│   ├── payments/
│   ├── accounts/
│   ├── fingerprinting/
│   ├── mitigation/
│   ├── trust/
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
│   ├── scoring/
│   ├── anomaly/
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
- rule metadata
- trust metadata
- mitigation states
- replay metadata
- investigation metadata

Generate:
- optimized schemas
- indexes
- immutable investigation history
- replay-safe synchronization tables

Requirements:
- replay safety
- trust consistency
- investigation correctness

Never:
- tightly couple realtime scoring to OLTP writes
- mutate immutable investigation history

---

## Redis Requirements

Use Redis for:
- hot risk cache
- realtime behavioral cache
- replay coordination
- distributed throttling

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
- behavioral events
- scoring updates
- fraud alerts
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
- distributed behavioral ingestion
- async synchronization
- realtime scoring

---

## ClickHouse Requirements

Use ClickHouse for:
- fraud analytics
- anomaly analytics
- behavioral analytics
- payment-risk analytics
- abuse analytics

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

## Risk Scoring Requirements

Support:
- realtime scoring
- behavioral scoring
- trust scoring
- degraded scoring fallback

Generate:
- scoring pipelines
- replay-safe score updates
- distributed scoring orchestration
- fallback scoring logic

The scoring system MUST tolerate:
- delayed behavioral events
- replay storms
- partial scoring failures

No fake scoring architecture.

---

## Rule Engine Requirements

Support:
- realtime rule evaluation
- dynamic rule updates
- replay-safe rule execution
- distributed rule orchestration

Generate:
- rule engine
- rule versioning
- rule synchronization
- fallback evaluation workflows

No naive rule logic.

---

## Behavioral Analysis Requirements

Support:
- login behavior analysis
- purchase behavior analysis
- device behavior analysis
- anomaly detection hooks

Generate:
- behavioral pipelines
- replay-safe aggregation
- feature freshness orchestration
- anomaly detection workflows

The behavioral system MUST be realistic.

---

## Device Fingerprinting Requirements

Support:
- session fingerprinting
- device linkage
- suspicious device correlation
- replay-safe fingerprint updates

Generate:
- fingerprint orchestration
- device correlation pipelines
- replay-safe synchronization
- mitigation hooks

---

## Mitigation Requirements

Support:
- rate limiting hooks
- transaction holds
- account freezing hooks
- step-up verification hooks

Generate:
- mitigation orchestration
- replay-safe mitigation workflows
- escalation pipelines
- trust enforcement hooks

The mitigation system MUST tolerate:
- replay storms
- adversarial traffic
- distributed attacks

---

## ML Inference Requirements

Support:
- anomaly scoring
- fraud prediction
- risk classification
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
- fraud detected
- risk score updated
- mitigation triggered
- anomaly detected
- account flagged
- payment blocked

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
- /risk/score
- /fraud/anomalies
- /fraud/mitigations
- /trust/accounts
- /fraud/analytics

Support:
- realtime scoring
- replay-safe evaluation
- low-latency risk lookups
- investigation workflows

---

## Security Requirements

The platform MUST:
- validate investigation access
- enforce RBAC
- sanitize behavioral input
- isolate scoring pipelines
- protect trust integrity

Never:
- expose raw fraud rules
- expose internal trust scores
- expose model internals
- trust external scoring blindly

Generate:
- authorization middleware
- replay validation
- scoring integrity validation
- mitigation isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- scoring latency
- inference latency
- anomaly detection latency
- mitigation count
- replay latency
- fraud detection rate

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive investigation evidence.

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
- inference tests
- mitigation tests
- concurrency tests

Test:
- duplicate behavioral events
- replay storms
- delayed analytics updates
- scoring degradation
- adversarial traffic
- ultra-high concurrency
- distributed deployments

---

## Output Requirements

Explain:
- fraud architecture
- scoring strategy
- rule-engine strategy
- behavioral-analysis strategy
- replay-safe ingestion strategy
- mitigation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy fraud service.
No fake rule-engine architecture.
No naive scoring logic.

---

## Acceptance Criteria

The Fraud Detection Platform must support future integration with:
- Billing Platform
- Recommendation Platform
- Notification Platform
- User Platform
- Analytics Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate behavioral events
- delayed analytics updates
- distributed deployments
- adversarial traffic

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.