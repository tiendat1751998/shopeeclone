# TASK-030 — LIVE COMMERCE PLATFORM

## Goal

Build a REAL production-grade Live Commerce Platform.

This platform is responsible for:
- livestream commerce
- realtime chat
- live ordering
- realtime engagement
- interactive events
- livestream recommendations
- realtime reactions
- live inventory synchronization
- stream moderation
- realtime fanout

This is NOT a toy livestream service.

The Live Commerce Platform must support:
- ultra-high websocket concurrency
- low-latency interaction
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe event synchronization

The architecture MUST prioritize:
- realtime responsiveness
- fanout scalability
- stream resiliency
- moderation correctness
- operational stability

---

## Tech Stack

Use:
- Golang
- Kafka
- Redis Cluster
- PostgreSQL
- ClickHouse
- WebSocket
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- WebRTC
- SFU infrastructure
- media orchestration
- realtime ML moderation hooks

---

## Core Responsibilities

The Live Commerce Platform MUST support:

### Livestream Commerce
- livestream orchestration
- stream lifecycle management
- seller livestreams
- viewer participation

### Realtime Chat
- websocket chat
- realtime messaging
- distributed fanout
- replay-safe chat delivery

### Live Ordering
- realtime flash purchases
- live product pinning
- checkout hooks
- replay-safe order synchronization

### Realtime Engagement
- reactions
- likes
- polls
- interactive participation

### Livestream Recommendations
- realtime stream recommendations
- trending livestreams
- behavioral livestream ranking
- personalization hooks

### Moderation
- realtime moderation
- spam mitigation
- abuse detection hooks
- stream enforcement workflows

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate fanout/chat/engagement
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Live Commerce Platform MUST:
- support replay-safe synchronization
- support distributed websocket fanout
- support realtime moderation
- support degraded stream fallback

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The live-commerce system MUST tolerate:
- retry storms
- duplicate websocket events
- delayed synchronization
- partial fanout failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/live-commerce/
├── cmd/
├── internal/
│   ├── config/
│   ├── streams/
│   ├── chat/
│   ├── fanout/
│   ├── engagement/
│   ├── orders/
│   ├── moderation/
│   ├── recommendations/
│   ├── synchronization/
│   ├── replay/
│   ├── websocket/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── media/
│   ├── webrtc/
│   ├── sfu/
│   ├── transcoding/
│   └── orchestration/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## PostgreSQL Requirements

Use PostgreSQL for:
- stream metadata
- moderation states
- replay metadata
- engagement metadata
- orchestration state

Generate:
- optimized schemas
- indexes
- immutable engagement history
- replay-safe synchronization tables

Requirements:
- replay safety
- stream consistency
- orchestration correctness

Never:
- tightly couple websocket fanout to OLTP writes
- mutate immutable engagement history

---

## Redis Requirements

Use Redis for:
- websocket session state
- hot engagement cache
- fanout coordination
- distributed throttling

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- ultra-high websocket concurrency
- low-latency fanout
- distributed deployments

The cache layer MUST be production-grade.

---

## Kafka Requirements

Use Kafka for:
- engagement events
- websocket synchronization
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
- distributed fanout
- async synchronization
- realtime interaction

---

## ClickHouse Requirements

Use ClickHouse for:
- livestream analytics
- engagement analytics
- viewer analytics
- moderation analytics
- realtime interaction analytics

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
- distributed websocket sessions
- replay-safe delivery
- realtime fanout
- connection recovery

Generate:
- websocket gateway
- distributed session coordination
- replay-safe messaging
- reconnection workflows

The websocket layer MUST tolerate:
- reconnect storms
- replay storms
- partial node failures

No fake websocket architecture.

---

## Fanout Requirements

Support:
- distributed fanout
- realtime engagement fanout
- live event broadcasting
- replay-safe synchronization

Generate:
- fanout orchestration
- replay-safe event distribution
- distributed broadcast pipelines
- fallback fanout workflows

No naive fanout logic.

---

## Moderation Requirements

Support:
- spam detection
- abuse mitigation
- realtime moderation
- replay-safe moderation events

Generate:
- moderation pipelines
- distributed moderation orchestration
- replay-safe moderation synchronization
- escalation workflows

The moderation system MUST be realistic.

---

## Live Ordering Requirements

Support:
- flash-sale ordering hooks
- realtime inventory synchronization
- replay-safe checkout coordination
- distributed purchase orchestration

Generate:
- live ordering workflows
- inventory synchronization hooks
- replay-safe ordering pipelines
- degraded purchase fallback logic

---

## Livestream Recommendation Requirements

Support:
- trending livestreams
- realtime personalization
- behavioral ranking
- recommendation hooks

Generate:
- livestream recommendation pipelines
- replay-safe ranking synchronization
- distributed recommendation orchestration
- cache-aware recommendations

---

## Event-Driven Requirements

Generate events for:
- livestream started
- livestream ended
- engagement recorded
- moderation triggered
- live order placed
- websocket degraded

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
- /live/streams
- /live/chat
- /live/engagement
- /live/orders
- /live/moderation

Support:
- websocket upgrades
- realtime fanout
- replay-safe synchronization
- low-latency interactions

---

## Security Requirements

The platform MUST:
- validate livestream ownership
- enforce RBAC
- sanitize realtime messages
- isolate websocket sessions
- protect fanout integrity

Never:
- expose websocket internals
- expose moderation internals
- trust client-side engagement blindly
- allow unauthorized stream controls

Generate:
- authorization middleware
- replay validation
- websocket isolation
- moderation enforcement

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
- engagement throughput
- moderation latency
- replay latency
- reconnect rate

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive session tokens.

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
- websocket tests
- fanout tests
- moderation tests
- concurrency tests

Test:
- duplicate websocket events
- replay storms
- reconnect storms
- delayed synchronization
- ultra-high websocket concurrency
- distributed deployments
- partial node failures

---

## Output Requirements

Explain:
- live-commerce architecture
- websocket strategy
- fanout strategy
- replay-safe synchronization strategy
- moderation strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy livestream service.
No fake websocket architecture.
No naive fanout logic.

---

## Acceptance Criteria

The Live Commerce Platform must support future integration with:
- Recommendation Platform
- Notification Platform
- Billing Platform
- Fraud Platform
- Analytics Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate websocket events
- reconnect storms
- distributed deployments
- ultra-high websocket concurrency

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.