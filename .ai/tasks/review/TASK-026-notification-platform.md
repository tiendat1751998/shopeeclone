# TASK-026 — NOTIFICATION PLATFORM

## Goal

Build a REAL production-grade Notification Platform.

This platform is responsible for:
- push notifications
- email pipelines
- SMS orchestration
- realtime alerts
- notification preferences
- distributed fanout
- campaign delivery
- retry orchestration
- delivery analytics
- multi-channel routing

This is NOT a toy notification service.

The Notification Platform must support:
- massive fanout
- low-latency delivery
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe delivery

The architecture MUST prioritize:
- delivery reliability
- provider resiliency
- preference correctness
- retry safety
- operational stability

---

## Tech Stack

Use:
- Golang
- Kafka
- Redis Cluster
- PostgreSQL
- ClickHouse
- gRPC
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Provider integrations:
- Firebase Cloud Messaging (FCM)
- Apple Push Notification service (APNs)
- SMTP providers
- SMS gateways

Optional:
- WebPush
- in-app notifications
- WebSocket realtime notifications

---

## Core Responsibilities

The Notification Platform MUST support:

### Push Notifications
- mobile push delivery
- token management
- platform-specific payloads
- provider retries

### Email Pipelines
- transactional emails
- campaign emails
- template rendering
- bounce handling

### SMS Orchestration
- OTP delivery
- transactional SMS
- provider failover
- regional routing

### Notification Preferences
- opt-in/opt-out
- quiet hours
- channel preferences
- frequency limiting

### Distributed Fanout
- bulk notifications
- segmented campaigns
- realtime fanout
- replay-safe delivery

### Delivery Analytics
- delivery rates
- open rates
- click-through rates
- provider performance analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate delivery/providers/preferences
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Notification Platform MUST:
- support replay-safe delivery
- support distributed fanout
- support provider failover
- support preference enforcement

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The notification system MUST tolerate:
- retry storms
- duplicate delivery events
- delayed provider callbacks
- partial failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/notification/
├── cmd/
├── internal/
│   ├── config/
│   ├── push/
│   ├── email/
│   ├── sms/
│   ├── preferences/
│   ├── campaigns/
│   ├── fanout/
│   ├── providers/
│   ├── analytics/
│   ├── synchronization/
│   ├── replay/
│   ├── idempotency/
│   ├── templates/
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
- notification metadata
- delivery states
- preference states
- provider metadata
- replay metadata

Generate:
- optimized schemas
- indexes
- immutable delivery history
- audit tables

Requirements:
- replay safety
- delivery consistency
- preference correctness

Never:
- mutate immutable delivery history
- tightly couple provider callbacks to delivery writes

---

## Push Notification Requirements

Support:
- FCM integration
- APNs integration
- device token lifecycle
- replay-safe delivery

Generate:
- push orchestration
- provider failover
- retry-safe delivery
- token invalidation workflows

The push system MUST tolerate:
- provider outages
- token invalidation storms
- duplicate delivery events

No fake push architecture.

---

## Email Requirements

Support:
- transactional email
- campaign delivery
- bounce handling
- replay-safe sending

Generate:
- email rendering pipelines
- distributed email fanout
- bounce reconciliation
- provider failover

The email system MUST be production-grade.

---

## SMS Requirements

Support:
- OTP delivery
- transactional SMS
- provider routing
- replay-safe sending

Generate:
- SMS orchestration
- provider failover
- retry-safe SMS workflows
- regional routing logic

No naive SMS architecture.

---

## Preference Requirements

Support:
- opt-in/opt-out
- quiet hours
- per-channel preferences
- frequency limiting

Generate:
- preference orchestration
- replay-safe preference updates
- preference cache synchronization
- enforcement middleware

The preference system MUST be realistic.

---

## Campaign Requirements

Support:
- segmented campaigns
- bulk fanout
- scheduled delivery
- replay-safe campaigns

Generate:
- campaign orchestration
- distributed fanout pipelines
- retry-safe bulk delivery
- campaign analytics hooks

---

## Kafka Requirements

Use Kafka for:
- delivery events
- campaign events
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
- analytics streaming

---

## Redis Requirements

Use Redis for:
- rate limiting
- preference cache
- hot delivery cache
- distributed coordination

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- massive fanout
- low-latency reads
- distributed deployments

---

## ClickHouse Requirements

Use ClickHouse for:
- delivery analytics
- open analytics
- CTR analytics
- provider analytics
- campaign analytics

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

## Provider Integration Requirements

Support:
- provider failover
- provider throttling
- webhook callbacks
- replay-safe provider synchronization

Generate:
- provider adapters
- webhook ingestion pipelines
- retry-safe callbacks
- provider reconciliation

The provider integration layer MUST tolerate:
- delayed callbacks
- provider outages
- webhook duplication
- retry storms

---

## Event-Driven Requirements

Generate events for:
- notification queued
- notification delivered
- notification failed
- preference updated
- campaign started
- provider degraded

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
- /notifications
- /preferences
- /campaigns
- /delivery
- /providers

Support:
- pagination
- filtering
- replay-safe delivery
- realtime notification status

---

## Security Requirements

The platform MUST:
- validate notification ownership
- enforce RBAC
- sanitize templates
- isolate providers
- protect delivery integrity

Never:
- expose provider credentials
- expose internal fanout topology
- trust provider callbacks blindly
- allow unauthorized campaign delivery

Generate:
- authorization middleware
- replay validation
- template sanitization
- provider isolation

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- delivery latency
- provider latency
- retry count
- bounce rate
- replay latency
- fanout throughput

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive OTP values or provider secrets.

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
- provider integration tests
- preference tests
- fanout tests
- concurrency tests

Test:
- duplicate delivery events
- replay storms
- delayed provider callbacks
- provider outages
- massive fanout
- high concurrency
- distributed deployments

---

## Output Requirements

Explain:
- notification architecture
- provider failover strategy
- fanout strategy
- replay-safe delivery strategy
- preference enforcement strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy notification service.
No fake fanout architecture.
No naive retry logic.

---

## Acceptance Criteria

The Notification Platform must support future integration with:
- Order Service
- Billing Platform
- Recommendation Platform
- Live Commerce Platform
- Fraud Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate delivery events
- delayed provider callbacks
- distributed deployments
- massive fanout

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.