# TASK-044 — FRAUD DETECTION & RISK ENGINE PLATFORM

## Goal

Build a REAL production-grade Fraud Detection & Risk Engine Platform.

This platform is responsible for:
- real-time fraud detection
- transaction risk scoring
- behavioral anomaly detection
- device fingerprint risk analysis
- account takeover detection
- payment fraud prevention signals
- bot detection signals
- ML-based risk scoring pipeline
- rule-based + ML hybrid decision engine
- risk-based access control decisions

This is NOT a simple rules engine.

The Fraud Platform is:
## REAL-TIME RISK DECISION ENGINE

The system MUST support:
- ultra-low latency risk scoring (<50ms target)
- high-throughput event processing
- ML inference integration
- streaming anomaly detection
- replay-safe risk evaluation

The architecture MUST prioritize:
- decision correctness under adversarial behavior
- latency efficiency
- model accuracy + fallback safety
- resilience under traffic spikes
- explainability of risk decisions

---

## Tech Stack

Use:
- Golang (real-time scoring engine)
- Python (ML training + anomaly models)
- Kafka (event streaming backbone)
- Redis Cluster (feature cache + risk cache)
- ClickHouse (behavior analytics)
- PostgreSQL (rules + config store)
- Vector DB (similar fraud pattern detection)
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- PyTorch / TensorFlow (fraud models)
- XGBoost / LightGBM (scoring models)
- Feast (feature store)
- Kafka Streams / Flink (stream processing)

---

## Core Responsibilities

The Fraud Platform MUST support:

### Real-Time Fraud Scoring
- transaction risk scoring (<50ms)
- request-level risk evaluation
- adaptive scoring models

### Behavioral Anomaly Detection
- login anomaly detection
- purchase pattern anomaly detection
- velocity checks (frequency-based fraud detection)
- geographic inconsistency detection

### Account Takeover Detection
- device fingerprint mismatch detection
- login pattern deviation detection
- session hijacking detection signals

### Bot Detection
- traffic pattern anomaly detection
- behavioral mouse/click pattern analysis (signal-based)
- rate pattern fingerprinting

### Payment Fraud Prevention Signals
- payment risk scoring
- chargeback risk prediction
- suspicious transaction clustering

---

## Architecture Requirements

The Fraud Platform MUST:
- be real-time streaming first
- separate scoring engine from ML pipeline
- support hybrid rules + ML decisioning
- support explainable risk output

The system MUST:
- handle replay-safe risk evaluation
- support idempotent scoring requests
- maintain deterministic scoring for same inputs where required

Use:
- streaming pipeline architecture
- feature store integration
- stateless scoring services (where possible)
- cached inference results for performance

The system MUST tolerate:
- adversarial attack patterns
- replayed fraud attempts
- burst traffic spikes (flash sales)
- ML model degradation or failure
- Kafka lag / delay in feature updates
- Redis cache failure fallback mode

---

## Folder Structure

Generate:

platforms/fraud-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── scoring/
│   ├── rules/
│   ├── ml/
│   ├── features/
│   ├── streaming/
│   ├── detection/
│   ├── anomaly/
│   ├── bot/
│   ├── payment/
│   ├── identity/
│   ├── device/
│   ├── velocity/
│   ├── riskengine/
│   ├── explainability/
│   ├── cache/
│   ├── vector/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── ml/
│   ├── training/
│   ├── inference/
│   ├── models/
│   ├── feature_engineering/
│   └── pipelines/
│
├── rules/
├── features/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Risk Scoring Engine (CRITICAL)

The scoring engine MUST:

- combine rule-based + ML-based scoring
- produce final risk score (0–1 or 0–100)
- return explainability metadata

Rules:
- velocity rules (time-based thresholds)
- device mismatch rules
- geo anomaly rules
- behavioral thresholds

ML:
- classification models
- anomaly detection models
- embedding-based similarity fraud detection

---

## Feature Store Requirements

Support:
- real-time feature retrieval
- offline + online feature parity
- feature versioning

Generate:
- feature pipeline engine
- caching layer
- feature freshness control

---

## Streaming Requirements

Use Kafka for:
- user events
- transaction events
- login events
- behavioral signals

Generate:
- stream processing pipeline
- event enrichment layer
- replay-safe consumers

---

## ML Model Integration

Support:
- real-time inference (<50ms)
- fallback scoring models
- versioned model deployment

Generate:
- model serving layer
- A/B model switching
- shadow model evaluation

---

## Rules Engine Requirements

Support:
- configurable fraud rules
- dynamic rule updates
- rule priority system

Generate:
- DSL-based rule engine
- rule evaluation pipeline
- hot-reload rule updates

---

## Explainability Requirements (CRITICAL)

Every decision MUST include:
- risk score
- reason codes
- contributing signals
- rule triggers

Generate:
- explanation builder engine
- audit-safe decision logs

---

## Redis Requirements

Use Redis for:
- real-time risk cache
- velocity tracking
- session-level risk state

Must support:
- atomic updates
- TTL-based risk decay
- failover fallback scoring

---

## ClickHouse Requirements

Use ClickHouse for:
- fraud analytics
- pattern detection
- historical risk trends
- model evaluation

---

## PostgreSQL Requirements

Use PostgreSQL for:
- rules configuration
- model metadata
- risk policies
- audit logs

---

## Vector Search Requirements

Support:
- fraud pattern similarity search
- embedding-based anomaly detection

Generate:
- vector indexing pipeline
- similarity scoring engine

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto definitions

Endpoints:
- /risk/score
- /risk/evaluate
- /risk/explain
- /risk/decision
- /rules/update

Must support:
- ultra-low latency scoring (<50ms)
- idempotent evaluation
- deterministic fallback mode

---

## Security Requirements (CRITICAL)

The Fraud Platform MUST:
- never be bypassable by client input
- validate all upstream signals
- isolate internal risk logic

Never:
- trust frontend risk signals
- allow direct override of risk decisions
- expose raw ML model internals externally

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logs
- risk decision audit logs

Metrics:
- risk score distribution
- fraud detection rate
- false positive rate
- scoring latency
- model drift indicators

---

## Reliability Requirements

Implement:
- retries (only safe operations)
- circuit breakers
- fallback scoring mode
- degraded mode rules engine

Critical:
## FRAUD DECISION MUST NEVER FAIL SILENTLY

---

## Kubernetes Requirements

Generate:
- Deployments
- HPA
- PDB
- ConfigMaps
- Secrets
- Helm charts

Must support:
- high throughput inference scaling
- low latency autoscaling
- burst traffic resilience

---

## CI/CD Requirements

Generate:
- model validation pipeline
- rule validation tests
- load testing
- adversarial testing simulation
- GitOps deployment

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- scoring consistency tests
- adversarial behavior tests
- latency stress tests

Test scenarios:
- fraud burst attack simulation
- replayed transaction attacks
- model drift failure
- Redis cache failure fallback
- Kafka lag in feature pipeline
- flash sale fraud spike

---

## Output Requirements

Explain:
- fraud architecture
- hybrid scoring strategy
- ML + rules integration model
- feature pipeline design
- explainability system
- real-time decision pipeline
- scaling strategy

Generate production-grade code only.

No toy fraud detection system.
No fake scoring logic.
No simplified ML pipeline.

---

## Acceptance Criteria

The Fraud Platform must support integration with:
- Payment Platform
- OMS Platform
- API Gateway
- Recommendation Platform (signals only)
- Notification Platform (risk alerts)
- Observability Platform

without redesign.

The system MUST survive:
- fraud attack bursts
- replay attacks
- model degradation
- streaming delays
- Redis/Kafka outages

WITHOUT compromising financial safety decisions.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
Strict real-time risk enforcement required.