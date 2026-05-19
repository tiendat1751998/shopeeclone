# TASK-037 — RECOMMENDATION & PERSONALIZATION PLATFORM

## Goal

Build a REAL production-grade Recommendation & Personalization Platform.

This platform is responsible for:
- real-time recommendation generation
- personalization ranking
- user embedding systems
- session-based ranking
- feed generation pipelines
- candidate generation
- ranking orchestration
- exploration/exploitation balancing
- behavioral personalization
- multi-surface recommendation (home, search, cart, checkout)

This is NOT a toy recommendation service.

The Recommendation Platform must support:
- ultra-high QPS personalization requests
- distributed candidate generation
- distributed ranking pipelines
- Kubernetes-native deployment
- observability-first architecture
- replay-safe personalization updates

The architecture MUST prioritize:
- latency (very strict)
- ranking quality
- freshness of signals
- scalability
- resiliency under load

---

## Tech Stack

Use:
- Golang (serving layer)
- Python (ML/ranking models)
- Kafka (event streaming)
- Redis Cluster (low-latency caching)
- PostgreSQL (metadata, configs)
- ClickHouse (behavior analytics)
- Elasticsearch/OpenSearch (candidate retrieval)
- Vector DB (Milvus / Faiss / Weaviate)
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Ray (distributed embedding computation)
- TensorFlow / PyTorch
- Triton Inference Server
- Feast (feature store integration)

---

## Core Responsibilities

The Recommendation Platform MUST support:

### Candidate Generation
- item retrieval from multiple sources:
  - collaborative filtering
  - content-based retrieval
  - trending items
  - search-derived candidates
  - vector similarity search
- multi-stage candidate fusion

### Ranking Engine
- ML ranking models
- business rules overlay
- personalization scoring
- session-aware ranking
- replay-safe ranking computation

### Real-Time Personalization
- live user behavior updates
- session event ingestion
- immediate feed adaptation
- streaming feature updates

### Feed Generation
- home feed
- product feed
- category feed
- search personalization feed
- cart/checkout upsell feed

### Exploration vs Exploitation
- bandit-based strategies
- diversity injection
- novelty balancing
- fatigue prevention

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate candidate generation / ranking / serving
- support distributed deployments
- support event-driven pipelines
- support eventual consistency

The Recommendation Platform MUST:
- support replay-safe personalization updates
- support degraded ranking fallback
- support cached feed serving under load spikes
- support multi-region traffic distribution

Use:
- CQRS where appropriate
- dependency injection
- modular pipeline architecture
- resilience patterns

The system MUST tolerate:
- retry storms
- duplicate user events
- delayed feature ingestion
- partial ranking model failure
- distributed deployment inconsistency
- replay storms

---

## Folder Structure

Generate:

platforms/recommendation-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── candidates/
│   ├── ranking/
│   ├── features/
│   ├── personalization/
│   ├── sessions/
│   ├── feeds/
│   ├── retrieval/
│   ├── exploration/
│   ├── embeddings/
│   ├── orchestration/
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
│   ├── ranking/
│   ├── embeddings/
│   ├── training/
│   └── serving/
│
├── pipelines/
├── vectorstore/
├── retrieval/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Candidate Generation Requirements

Support:
- collaborative filtering
- content-based retrieval
- trending pipeline
- vector similarity search
- search-derived candidates

Generate:
- multi-source candidate pipelines
- caching layers
- replay-safe candidate ingestion
- distributed retrieval orchestration

Must support:
- high recall under latency constraints
- fallback candidate sources

---

## Ranking Requirements

Support:
- ML-based ranking models
- real-time ranking adjustment
- business rule overlays
- session-based ranking
- personalization scoring

Generate:
- ranking pipeline
- feature retrieval layer
- replay-safe ranking computation
- multi-stage ranking architecture

Ranking MUST:
- be latency optimized (<100ms target)
- support model fallback

---

## Real-Time Personalization Requirements

Support:
- live event ingestion
- session updates
- click/view/cart signals
- immediate ranking adaptation

Generate:
- streaming personalization engine
- event-driven updates
- replay-safe session state management

---

## Feed Generation Requirements

Support:
- home feed
- category feed
- product feed
- search feed personalization
- upsell/cross-sell feed

Generate:
- feed orchestration layer
- caching strategy
- precomputed feed pipelines
- real-time fallback generation

---

## Vector Search Requirements

Use vector DB for:
- semantic similarity search
- embedding-based retrieval
- hybrid ranking pipelines

Generate:
- embedding indexing pipelines
- similarity search orchestration
- fallback retrieval mechanisms

---

## Kafka Requirements

Use Kafka for:
- user behavior events
- clickstream ingestion
- ranking signals
- replay-safe personalization updates

Generate:
- topic strategy
- partitioning strategy
- DLQ topics
- replay pipelines

Requirements:
- idempotent consumers
- replay-safe ingestion
- event versioning

---

## Redis Requirements

Use Redis for:
- hot feed caching
- session state
- real-time personalization cache
- ranking acceleration

Generate:
- TTL strategy
- cache invalidation policies
- distributed coordination

Must support:
- ultra-low latency feed delivery

---

## ClickHouse Requirements

Use ClickHouse for:
- user behavior analytics
- recommendation effectiveness
- ranking performance
- A/B testing evaluation

Generate:
- aggregation pipelines
- materialized views
- partitioning strategy

---

## ML Integration Requirements

Integrate with ML Platform for:
- ranking models
- embedding generation
- feature retrieval
- inference serving

Must support:
- online inference
- fallback ranking models
- versioned model rollout

---

## Event-Driven Requirements

Generate events for:
- user clicked item
- item viewed
- item purchased
- feed generated
- ranking updated

Requirements:
- retries
- DLQ
- replay-safe processing
- idempotent consumers
- versioned events

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /recommend/home
- /recommend/search
- /recommend/session
- /recommend/feed
- /recommend/rank

Must support:
- low latency (<100ms)
- personalized responses
- fallback responses under load

---

## Security Requirements

The platform MUST:
- validate user identity context
- enforce RBAC for internal APIs
- isolate recommendation domains
- protect behavioral data

Never:
- expose raw user behavioral streams
- allow unrestricted feature injection
- trust external signals blindly

Generate:
- authorization middleware
- replay validation
- data isolation layer

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- feed latency
- ranking latency
- CTR (click-through rate)
- conversion rate
- cache hit ratio
- replay lag

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive user behavior raw.

---

## Reliability Requirements

Implement:
- retries
- timeout handling
- graceful shutdown
- circuit breakers
- backoff strategies

Support:
- autoscaling
- rolling deployment
- distributed execution

Generate:
- resilience middleware
- fallback ranking strategies
- degraded mode feeds

---

## Kubernetes Requirements

Generate:
- Deployments
- StatefulSets
- Services
- ConfigMap
- Secret integration
- HPA
- PDB
- ServiceMonitor
- NetworkPolicy
- Helm charts

Must support:
- burst traffic (flash sale)
- autoscaling
- low latency serving

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone pipelines
- ML model validation
- Helm validation
- Kubernetes policy validation
- load testing
- GitOps workflows

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- ranking tests
- feed generation tests
- concurrency tests
- replay tests

Test:
- traffic spikes
- cache invalidation correctness
- ranking fallback correctness
- distributed consistency issues

---

## Output Requirements

Explain:
- recommendation architecture
- candidate generation strategy
- ranking strategy
- personalization strategy
- replay-safe pipeline design
- scaling strategy
- latency optimization strategy

Generate production-grade code only.

No toy recommender system.
No fake ML pipeline.
No naive ranking logic.

---

## Acceptance Criteria

The Recommendation Platform must support:
- AI/ML Platform integration
- Search Platform integration
- Fraud Platform integration
- Analytics Platform integration
- Live Commerce Platform integration

without major refactors.

The platform MUST realistically tolerate:
- ultra-high QPS traffic spikes
- replay storms
- partial model failures
- cache failures
- distributed inconsistency

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.