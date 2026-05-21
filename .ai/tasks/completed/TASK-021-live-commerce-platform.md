# TASK-021 — LIVE COMMERCE PLATFORM

## Goal

Build a REAL production-grade Live Commerce Platform.

This platform is responsible for:
- livestream commerce
- realtime interactions
- chat fanout
- live recommendations
- realtime engagement
- live order orchestration
- realtime moderation
- livestream analytics
- gift/reaction systems
- distributed websocket orchestration

This is NOT a toy livestream service.

The Live Commerce Platform must support:
- millions of concurrent viewers
- ultra-low latency fanout
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe event streaming

The architecture MUST prioritize:
- realtime responsiveness
- fanout scalability
- websocket resiliency
- moderation correctness
- operational stability

---

## Tech Stack

Use:
- Golang
- WebSocket
- gRPC
- Redis Cluster
- Kafka
- ClickHouse
- PostgreSQL
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- WebRTC hooks
- SFU integration hooks
- media CDN hooks

---

## Core Responsibilities

The Live Commerce Platform MUST support:

### Livestream Commerce
- livestream sessions
- live product pinning
- realtime product synchronization
- live checkout hooks

### Realtime Interactions
- reactions
- likes
- gifts
- viewer engagement
- realtime counters

### Chat Fanout
- massive websocket fanout
- realtime messaging
- distributed room coordination
- replay-safe delivery

### Live Recommendations
- trending livestreams
- personalized livestream discovery
- realtime recommendation hooks

### Moderation
- realtime moderation
- spam detection hooks
- toxic message filtering hooks
- moderation workflows

### Livestream Analytics
- concurrent viewer analytics
- engagement analytics
- retention analytics
- interaction analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate fanout/chat/analytics
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Live Commerce Platform MUST:
- support ultra-low latency fanout
- support replay-safe event streaming
- support distributed websocket orchestration
- support realtime moderation

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The livestream system MUST tolerate:
- retry storms
- websocket reconnect storms
- duplicate events
- delayed fanout
- partial failures
- distributed deployments

---

## Folder Structure

Generate:

platforms/live-commerce/
├── cmd/
├── internal/
│   ├── config/
│   ├── livestreams/
│   ├── chat/
│   ├── fanout/
│   ├── websocket/
│   ├── moderation/
│   ├── recommendations/
│   ├── engagement/
│   ├── gifts/
│   ├── analytics/
│   ├── synchronization/
│   ├── replay/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## PostgreSQL Requirements

Use PostgreSQL for:
- livestream metadata
- moderation state
- room metadata
- replay metadata
- gift metadata

Generate:
- optimized schemas
- indexes
- immutable interaction history
- audit tables

Requirements:
- transactional correctness
- replay safety
- room consistency

Never:
- mutate immutable interaction history
- tightly couple websocket state to DB

---

## Redis Requirements

Use Redis for:
- websocket coordination
- fanout coordination
- realtime counters
- room presence
- distributed throttling

Generate:
- TTL strategy
- replay protection
- distributed coordination
- rate limiting

Support:
- massive websocket concurrency
- high fanout
- distributed deployments

The coordination layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- chat events
- livestream events
- replay-safe streaming
- distributed analytics
- moderation pipelines

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
- realtime streaming
- massive fanout
- distributed consumers

---

## ClickHouse Requirements

Use ClickHouse for:
- engagement analytics
- viewer analytics
- livestream analytics
- retention analytics
- interaction aggregation

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

## WebSocket Requirements

Support:
- distributed websocket gateways
- room fanout
- reconnect handling
- replay-safe delivery
- heartbeat handling

Generate:
- websocket gateway architecture
- reconnect orchestration
- distributed room coordination
- replay-safe message delivery

The websocket system MUST tolerate:
- reconnect storms
- duplicate deliveries
- partial fanout failures

No fake websocket architecture.

---

## Chat Requirements

Support:
- realtime messaging
- message fanout
- moderation hooks
- replay-safe delivery

Generate:
- chat orchestration
- message persistence hooks
- distributed room routing
- replay-safe chat delivery

---

## Moderation Requirements

Support:
- spam filtering hooks
- toxic content detection hooks
- moderation queues
- user muting/banning

Generate:
- moderation pipelines
- replay-safe moderation workflows
- moderation audit trails

The moderation system MUST be realistic.

---

## Engagement Requirements

Support:
- reactions
- likes
- gifts
- realtime engagement counters

Generate:
- distributed engagement aggregation
- replay-safe counters
- realtime synchronization

---

## Event-Driven Requirements

Generate events for:
- livestream started
- livestream ended
- message sent
- reaction added
- gift sent
- moderation triggered

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
- distributed fanout
- replay-safe streaming

No fake async architecture.

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- WebSocket APIs
- OpenAPI specs
- proto files

Endpoints:
- /live
- /live/chat
- /live/reactions
- /live/gifts
- /live/recommendations

Support:
- pagination
- realtime subscriptions
- websocket authentication
- replay-safe messaging

---

## Security Requirements

The platform MUST:
- validate websocket sessions
- enforce RBAC
- sanitize chat input
- isolate livestream rooms
- protect moderation integrity

Never:
- expose raw websocket topology
- expose moderation internals
- trust client viewer counts
- trust client engagement counters

Generate:
- websocket auth middleware
- replay validation
- room isolation
- moderation integrity validation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- websocket connection count
- fanout latency
- reconnect rate
- moderation latency
- engagement latency
- replay latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive moderation evidence.

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
- reconnect recovery workflows

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
- websocket tests
- replay tests
- fanout tests
- moderation tests
- concurrency tests

Test:
- reconnect storms
- duplicate events
- replay correctness
- websocket scaling
- massive fanout
- realtime moderation
- distributed deployments

---

## Output Requirements

Explain:
- livestream architecture
- websocket architecture
- fanout strategy
- moderation strategy
- replay strategy
- analytics strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy livestream service.
No fake websocket architecture.
No naive fanout system.

---

## Acceptance Criteria

The Live Commerce Platform must support future integration with:
- Recommendation Platform
- Notification Platform
- Advertising Platform
- Order Service
- Analytics Platform

without major future refactors.

The platform MUST realistically tolerate:
- reconnect storms
- duplicate events
- replay storms
- distributed deployments
- millions of concurrent viewers

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