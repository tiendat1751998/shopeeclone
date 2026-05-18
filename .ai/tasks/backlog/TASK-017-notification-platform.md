# TASK-017 — NOTIFICATION PLATFORM

## Goal

Build a REAL production-grade Notification Platform.

This platform is responsible for:
- push notifications
- email delivery
- SMS delivery
- in-app notifications
- distributed fanout
- retry orchestration
- notification preferences
- template systems
- delivery analytics
- provider failover

This is NOT a toy notification service.

The Notification Platform must support:
- millions of notifications
- burst traffic fanout
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe delivery

The architecture MUST prioritize:
- delivery reliability
- retry safety
- fanout scalability
- provider resiliency
- operational stability

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- Redis Cluster
- Kafka or NATS JetStream
- ClickHouse
- PostgreSQL or MySQL
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Firebase Cloud Messaging
- APNs
- SES
- SendGrid
- Twilio

---

## Core Responsibilities

The Notification Platform MUST support:

### Push Notifications
- mobile push delivery
- device token management
- notification fanout
- push retry handling

### Email Delivery
- transactional emails
- bulk emails
- provider failover
- bounce handling

### SMS Delivery
- OTP delivery hooks
- transactional SMS
- provider abstraction
- retry-safe delivery

### In-App Notifications
- notification inbox
- unread counters
- real-time delivery hooks
- notification synchronization

### Template Systems
- dynamic templates
- localization hooks
- template versioning
- variable rendering

### Delivery Analytics
- delivery tracking
- open tracking hooks
- click tracking hooks
- retry analytics

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Notification Platform MUST:
- support massive fanout
- support replay-safe delivery
- support provider failover
- support distributed retry orchestration

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The notification system MUST tolerate:
- retry storms
- duplicate delivery events
- provider instability
- delayed delivery
- partial failures
- distributed deployments

---

## Folder Structure

Generate:

platforms/notification/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── push/
│   ├── email/
│   ├── sms/
│   ├── inbox/
│   ├── templates/
│   ├── providers/
│   ├── fanout/
│   ├── retries/
│   ├── analytics/
│   ├── preferences/
│   ├── synchronization/
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

## Database Requirements

Use PostgreSQL or MySQL for:
- notification metadata
- inbox state
- template metadata
- delivery state
- retry state
- provider mappings

Use ClickHouse for:
- delivery analytics
- retry analytics
- engagement analytics

Generate:
- optimized schemas
- indexes
- immutable delivery history
- audit tables

Requirements:
- transactional correctness
- replay safety
- pagination everywhere
- inbox consistency

Never:
- mutate immutable delivery history
- tightly couple provider logic
- trust provider callbacks blindly

---

## Redis Requirements

Use Redis for:
- hot inbox cache
- unread counters
- fanout coordination
- retry coordination
- rate limiting

Generate:
- TTL strategy
- replay protection
- distributed coordination
- retry handling

Support:
- massive fanout
- high concurrency
- distributed deployments

The coordination layer MUST be production-grade.

---

## Event-Driven Requirements

Generate events for:
- notification requested
- notification delivered
- notification failed
- notification retried
- notification opened
- notification clicked

Use:
- Kafka or NATS JetStream

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
- async delivery orchestration

No fake async architecture.

---

## Provider Integration Requirements

Support:
- multiple push providers
- multiple email providers
- multiple SMS providers
- provider failover
- provider abstraction

Generate:
- provider adapter interfaces
- provider retry workflows
- provider health tracking
- failover orchestration

The provider integration MUST be production-grade.

No fake provider abstraction.

---

## Fanout Requirements

Support:
- massive notification fanout
- segmented fanout
- delayed delivery
- distributed scheduling

Generate:
- fanout orchestration
- scheduling workflows
- queue partitioning strategy
- replay-safe fanout handling

The fanout system MUST tolerate:
- burst traffic
- retry storms
- provider throttling
- delayed processing

---

## Template Requirements

Support:
- localized templates
- template versioning
- variable rendering
- dynamic personalization hooks

Generate:
- template rendering engine
- localization strategy
- rendering validation

The template system MUST be realistic.

---

## Inbox Requirements

Support:
- in-app notification inbox
- unread counters
- notification synchronization
- read/unread workflows

Generate:
- inbox caching strategy
- synchronization workflows
- replay-safe inbox updates

---

## Analytics Requirements

Use ClickHouse for:
- delivery analytics
- retry analytics
- engagement analytics
- provider performance analytics

Generate:
- analytics ingestion flows
- replay-safe ingestion
- aggregation strategy

Support:
- high ingestion throughput
- distributed aggregation
- low-cost analytics queries

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /notifications
- /notifications/inbox
- /notifications/preferences
- /notifications/templates
- /notifications/send

Support:
- pagination
- filtering
- localization
- idempotency keys

---

## Security Requirements

The platform MUST:
- validate requests
- enforce RBAC
- sanitize input
- isolate delivery analytics
- protect provider credentials

Never:
- expose raw provider credentials
- expose internal retry orchestration
- expose internal delivery pipelines

Generate:
- authorization middleware
- replay validation
- provider credential isolation
- notification integrity validation

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
- provider failure rate
- retry count
- fanout latency
- inbox synchronization latency
- template rendering latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive provider credentials.

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
- fallback delivery workflows

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
- provider tests
- replay tests
- fanout tests
- inbox tests
- analytics ingestion tests
- concurrency tests

Test:
- duplicate delivery events
- retry storms
- provider instability
- delayed delivery
- fanout spikes
- distributed deployments
- high concurrency

---

## Output Requirements

Explain:
- notification architecture
- provider abstraction strategy
- fanout orchestration
- retry strategy
- inbox synchronization strategy
- template strategy
- analytics strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy notification service.
No fake provider abstraction.
No naive fanout architecture.

---

## Acceptance Criteria

The Notification Platform must support future integration with:
- Order Service
- Payment Service
- Recommendation Platform
- User Platform
- Marketing Platform

without major future refactors.

The platform MUST realistically tolerate:
- retry storms
- provider instability
- duplicate delivery events
- delayed delivery
- distributed deployments
- massive fanout traffic

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