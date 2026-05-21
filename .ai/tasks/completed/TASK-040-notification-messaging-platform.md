# TASK-040 — NOTIFICATION & MESSAGING PLATFORM

## Goal

Build a REAL production-grade Notification & Messaging Platform.

This platform is responsible for:
- push notifications (mobile/web)
- email delivery pipeline
- SMS messaging orchestration
- realtime websocket communication
- in-app messaging system
- event-driven notification routing
- multi-channel delivery guarantees
- notification deduplication & idempotency
- user preference management
- delivery tracking & retry orchestration

This is NOT a toy notification service.

The Notification Platform is a:
## UNIFIED COMMUNICATION INFRASTRUCTURE LAYER

The system MUST support:
- ultra-high throughput messaging bursts
- multi-channel fanout processing
- real-time + async hybrid delivery
- Kubernetes-native deployment
- observability-first architecture
- replay-safe message delivery

The architecture MUST prioritize:
- delivery reliability
- idempotency
- scalability
- provider resilience
- low-latency realtime delivery

---

## Tech Stack

Use:
- Golang (core delivery engine)
- Kafka (event streaming backbone)
- Redis Cluster (dedup + rate limiting + cache)
- PostgreSQL (notification state + preferences)
- ClickHouse (message analytics)
- WebSocket Gateway (realtime delivery)
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

External Providers:
- Firebase Cloud Messaging (push)
- APNs (iOS push)
- SMTP providers (email)
- SMS gateways (Twilio-like abstraction)

Optional:
- NATS (low-latency messaging)
- Temporal (workflow orchestration)

---

## Core Responsibilities

The Notification Platform MUST support:

### Multi-Channel Delivery
- push notification (mobile/web)
- email delivery
- SMS delivery
- in-app notifications
- websocket realtime messaging

### Event-Driven Notification Routing
- event ingestion from Kafka
- routing rules engine
- notification template rendering
- priority-based delivery

### Delivery Orchestration
- retry strategies per channel
- fallback channel switching
- provider failover handling
- delayed delivery scheduling

### User Preferences System
- opt-in/out management
- channel preferences
- quiet hours handling
- localization preferences

### Delivery Tracking
- message state tracking
- delivery confirmation
- failure classification
- retry scheduling

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate ingestion / routing / delivery / tracking layers
- support distributed deployment
- support event-driven architecture

The Notification Platform MUST:
- support replay-safe delivery
- ensure idempotent message dispatch
- support provider failover
- handle partial delivery success scenarios

Use:
- CQRS where appropriate
- dependency injection
- modular delivery pipelines
- retry + circuit breaker patterns

The system MUST tolerate:
- duplicate Kafka events
- provider outages (email/SMS/push)
- retry storms
- websocket reconnections
- regional delivery delays
- partial fanout failures

---

## Folder Structure

Generate:

platforms/notification-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── ingestion/
│   ├── routing/
│   ├── templates/
│   ├── channels/
│   │   ├── push/
│   │   ├── email/
│   │   ├── sms/
│   │   ├── websocket/
│   │   └── inapp/
│   ├── delivery/
│   ├── fanout/
│   ├── preferences/
│   ├── scheduling/
│   ├── retry/
│   ├── failover/
│   ├── deduplication/
│   ├── orchestration/
│   ├── events/
│   ├── websocket/
│   ├── cache/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── db/
│   ├── migrations/
│   ├── schema.sql
│   └── notification_models.sql
│
├── providers/
│   ├── fcm/
│   ├── apns/
│   ├── smtp/
│   ├── sms/
│   └── webhook/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Event Ingestion Requirements

Support ingestion from:
- Kafka topics
- internal platform events
- external webhook triggers

Generate:
- event normalization layer
- schema versioning
- replay-safe ingestion
- idempotent event handling

---

## Routing Engine Requirements

Support:
- rule-based routing
- priority-based routing
- user segmentation routing
- event-type routing

Generate:
- routing DSL engine
- rule evaluation system
- fallback routing strategy

---

## Multi-Channel Delivery Requirements

### Push Notifications
- mobile (FCM/APNs)
- web push

### Email
- template rendering
- SMTP provider abstraction
- retry-safe sending

### SMS
- gateway abstraction
- rate-limited sending
- failover provider switching

### WebSocket Realtime
- low-latency delivery
- reconnect-safe sessions
- fanout optimization

### In-App Messaging
- persistent notification storage
- user inbox system

---

## Fanout Requirements (CRITICAL)

Support:
- 1 event → millions of users
- batch fanout processing
- segmented fanout
- async delivery pipelines

Generate:
- scalable fanout workers
- partitioned processing
- queue backpressure handling

---

## Deduplication & Idempotency

All messages MUST be idempotent.

Generate:
- dedup keys (event_id + user_id + channel)
- Redis-based dedup store
- Kafka replay safety handling
- retry-safe delivery guarantees

Prevent:
- duplicate notifications
- duplicate SMS/email sends
- duplicate websocket pushes

---

## Scheduling Requirements

Support:
- delayed notifications
- scheduled campaigns
- retry scheduling
- quiet hours enforcement

Generate:
- scheduling queue
- delayed job system
- retry backoff engine

---

## Kafka Requirements

Use Kafka for:
- notification events
- delivery pipelines
- retry queues
- fanout processing
- analytics events

Generate:
- topic strategy
- partitioning strategy
- DLQ topics
- replay-safe consumers

---

## Redis Requirements

Use Redis for:
- deduplication
- rate limiting
- session cache (websocket)
- delivery state cache

Must support:
- high throughput writes
- TTL-based cleanup
- distributed coordination

---

## PostgreSQL Requirements

Use PostgreSQL for:
- notification state
- user preferences
- delivery history
- template metadata

Generate:
- normalized schema
- indexing strategy
- audit logs

---

## ClickHouse Requirements

Use ClickHouse for:
- delivery analytics
- open/click tracking
- channel performance
- failure rate analysis

---

## WebSocket Gateway Requirements

Support:
- realtime push delivery
- session management
- reconnection handling
- horizontal scaling

Generate:
- gateway cluster design
- session routing strategy
- fanout optimization

---

## Template System Requirements

Support:
- dynamic templates
- localization (i18n)
- variable injection
- versioned templates

Generate:
- template engine
- safe rendering system
- fallback templates

---

## Event-Driven Requirements

Generate events for:
- notification created
- notification delivered
- notification failed
- notification retried
- notification opened/clicked

Rules:
- retries
- DLQ handling
- replay-safe processing
- idempotent execution
- versioned events

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto definitions

Endpoints:
- /notify/send
- /notify/bulk
- /notify/status
- /notify/preferences
- /notify/templates

Must support:
- high throughput ingestion
- idempotent requests
- real-time delivery fallback

---

## Security Requirements

The platform MUST:
- validate notification ownership
- enforce RBAC for campaigns
- isolate user data
- prevent spam abuse

Never:
- allow unrestricted mass messaging
- expose internal delivery pipelines
- trust external event payloads blindly

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logs
- distributed tracing
- correlation IDs

Metrics:
- delivery latency
- success rate per channel
- retry rate
- fanout throughput
- websocket connection health

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive user content.

---

## Reliability Requirements

Implement:
- retries (channel-aware)
- circuit breakers
- backpressure control
- graceful shutdown
- provider fallback

Critical:
## DELIVERY FAILURE MUST NOT CAUSE MESSAGE LOSS

---

## Kubernetes Requirements

Generate:
- Deployments
- StatefulSets (WebSocket layer)
- Services
- HPA
- PDB
- ConfigMaps
- Secrets
- Helm charts

Must support:
- burst traffic fanout
- autoscaling delivery workers
- resilient websocket scaling

---

## CI/CD Requirements

Generate:
- GitHub Actions / Drone pipelines
- provider integration tests
- load tests for fanout
- Kubernetes validation
- chaos testing
- GitOps deployment

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- fanout tests
- deduplication tests
- delivery retry tests
- websocket stress tests

Test scenarios:
- provider outage (SMS/email)
- Kafka replay duplication
- fanout to millions of users
- websocket reconnect storms
- retry explosion scenarios

---

## Output Requirements

Explain:
- notification architecture
- routing strategy
- fanout strategy
- delivery reliability model
- deduplication strategy
- websocket scaling strategy
- failure recovery strategy

Generate production-grade code only.

No toy notification system.
No fake delivery logic.
No unsafe fanout shortcuts.

---

## Acceptance Criteria

The Notification Platform must support integration with:
- OMS Platform
- Payment Platform
- Fraud Platform
- Recommendation Platform
- Search Platform
- Analytics Platform

without redesign.

The system MUST survive:
- massive fanout spikes
- provider outages
- replay storms
- websocket reconnection storms
- retry explosions

WITHOUT duplicate or lost notifications.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.