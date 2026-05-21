# TASK-039 — OMS & FULFILLMENT PLATFORM

## Goal

Build a REAL production-grade Order Management System (OMS) & Fulfillment Platform.

This platform is responsible for:
- order lifecycle management
- inventory reservation & allocation
- fulfillment orchestration
- warehouse coordination
- shipment tracking
- delivery state synchronization
- order state machine consistency
- partial fulfillment handling
- cancellation & returns coordination

This is NOT a simple order CRUD system.

The OMS is the:
## EXECUTION ENGINE OF E-COMMERCE ORDERS

The system MUST support:
- ultra-high concurrency order placement
- distributed inventory allocation
- event-driven fulfillment pipeline
- Kubernetes-native deployment
- observability-first architecture
- replay-safe order processing

The architecture MUST prioritize:
- order correctness (NO DUPLICATES / NO LOSS)
- inventory consistency
- fulfillment reliability
- state machine correctness
- scalability under flash-sale traffic

---

## Tech Stack

Use:
- Golang
- Kafka
- PostgreSQL
- Redis Cluster
- ClickHouse
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm
- gRPC

Optional:
- Temporal (workflow orchestration)
- Debezium (CDC sync)
- RabbitMQ (legacy integration simulation)
- ElasticSearch (order search)
- Cassandra (high-scale order events)

---

## Core Responsibilities

The OMS MUST support:

### Order Lifecycle Management
- order creation
- order confirmation
- order cancellation
- order completion
- order state machine enforcement

### Inventory Reservation
- real-time stock reservation
- distributed inventory locking
- oversell prevention
- reservation expiry handling

### Fulfillment Orchestration
- warehouse selection (multi-warehouse routing)
- picking & packing coordination
- shipment creation
- carrier integration

### Shipment Tracking
- delivery status updates
- tracking number management
- multi-carrier support
- delayed event reconciliation

### Returns & Cancellations
- return requests
- refund coordination (integrates Payment Platform)
- partial returns handling
- cancellation windows

---

## Architecture Requirements

The OMS MUST:
- follow clean architecture
- enforce strict state machine design
- separate order / inventory / fulfillment domains
- support event-driven architecture

The OMS MUST:
- support replay-safe order processing
- support idempotent order creation
- support distributed inventory consistency
- support eventual consistency across warehouses

Use:
- Saga pattern (MANDATORY)
- CQRS where appropriate
- state machine enforcement layer
- dependency injection
- modular domain boundaries

The system MUST tolerate:
- retry storms
- duplicate order requests
- delayed inventory sync
- partial warehouse failures
- Kafka replay duplication
- network partitions
- service degradation

---

## Folder Structure

Generate:

platforms/oms-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── orders/
│   ├── state_machine/
│   ├── inventory/
│   ├── allocation/
│   ├── fulfillment/
│   ├── warehouse/
│   ├── shipping/
│   ├── returns/
│   ├── cancellation/
│   ├── orchestration/
│   ├── saga/
│   ├── synchronization/
│   ├── replay/
│   ├── events/
│   ├── cache/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── workflows/
├── db/
│   ├── migrations/
│   ├── schema.sql
│   └── state_models.sql
│
├── integrations/
│   ├── payment/
│   ├── logistics/
│   └── warehouse/
│
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Order State Machine Requirements (CRITICAL)

The order system MUST enforce strict state transitions:

Example states:
- CREATED
- PAYMENT_PENDING
- PAID
- INVENTORY_RESERVED
- ALLOCATED
- PACKING
- SHIPPED
- DELIVERED
- CANCELLED
- RETURNED
- FAILED

Rules:
- NO illegal transitions
- NO state skipping
- NO concurrent state mutation
- state transitions MUST be deterministic

Generate:
- state machine engine
- transition validator
- replay-safe state recovery

---

## Inventory Management Requirements

Support:
- real-time stock deduction
- reservation-based inventory
- rollback on failure
- distributed stock consistency

Generate:
- inventory allocation engine
- reservation service
- expiry-based release system
- replay-safe inventory sync

Critical:
## NO overselling allowed under any scenario

---

## Fulfillment Requirements

Support:
- warehouse selection strategy
- multi-warehouse split orders
- picking & packing workflows
- shipment generation

Generate:
- warehouse routing engine
- fulfillment orchestration
- shipment pipeline
- carrier abstraction layer

---

## Saga Orchestration Requirements (MANDATORY)

Use Saga pattern for:

Order flow:
1. Create Order
2. Reserve Inventory
3. Process Payment (integration with Payment Platform)
4. Allocate Warehouse
5. Ship Order

If any step fails:
- compensate previous steps

Generate:
- saga coordinator
- compensation handlers
- replay-safe execution engine

---

## Kafka Requirements

Use Kafka for:
- order events
- inventory updates
- fulfillment events
- shipment updates
- saga orchestration events

Generate:
- topic design
- partitioning strategy
- DLQ topics
- replay pipelines

Requirements:
- idempotent consumers
- event versioning
- replay-safe processing

---

## PostgreSQL Requirements

Use PostgreSQL for:
- order state storage
- inventory state
- saga state
- fulfillment tracking

Generate:
- normalized schemas
- indexing strategy
- transactional integrity rules
- audit logs

Critical:
- state MUST be consistent
- no partial writes without saga tracking

---

## Redis Requirements

Use Redis for:
- inventory locks
- idempotency keys
- hot order cache
- distributed coordination

Must support:
- atomic reservation locking
- safe expiry handling
- failover correctness

---

## ClickHouse Requirements

Use ClickHouse for:
- order analytics
- fulfillment performance
- warehouse efficiency
- delivery time analytics

---

## Event-Driven Requirements

Generate events for:
- order created
- inventory reserved
- order shipped
- order delivered
- order cancelled
- fulfillment failed

Rules:
- retries
- DLQ
- replay-safe consumers
- idempotent processing
- event versioning

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /orders/create
- /orders/status
- /orders/cancel
- /inventory/reserve
- /fulfillment/ship
- /tracking/update

Must support:
- high concurrency
- idempotent order creation
- consistent state machine enforcement

---

## Security Requirements

The OMS MUST:
- validate order ownership
- enforce RBAC for warehouse actions
- isolate inventory operations
- protect order state integrity

Never:
- allow direct state mutation bypassing state machine
- trust external warehouse updates blindly
- expose internal saga state publicly

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logs
- distributed tracing
- correlation IDs

Metrics:
- order throughput
- fulfillment latency
- inventory reservation success rate
- saga failure rate
- state transition latency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log sensitive payment or inventory secrets.

---

## Reliability Requirements

Implement:
- retries (idempotent only)
- timeout handling
- circuit breakers
- graceful shutdown
- saga recovery logic

Critical:
## ORDER STATE MUST NEVER CORRUPT

---

## Kubernetes Requirements

Generate:
- Deployments
- StatefulSets
- Services
- HPA
- PDB
- ConfigMaps
- Secrets
- Helm charts

Must support:
- flash-sale traffic bursts
- horizontal scaling
- safe rolling upgrades

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone pipelines
- migration validation
- saga workflow tests
- chaos testing
- Kubernetes validation
- GitOps deployment

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- saga tests
- inventory race tests
- state machine tests
- concurrency tests

Test scenarios:
- flash sale oversell attempt
- duplicate order requests
- warehouse failure mid-fulfillment
- payment success but inventory failure
- Kafka replay duplication
- partial shipment updates

---

## Output Requirements

Explain:
- OMS architecture
- state machine design
- saga orchestration strategy
- inventory consistency strategy
- fulfillment orchestration strategy
- failure recovery strategy
- scaling strategy

Generate production-grade code only.

No toy order system.
No unsafe inventory logic.
No naive state machine.

---

## Acceptance Criteria

The OMS must support integration with:
- Payment Platform
- Fraud Platform
- Recommendation Platform (signals only)
- Search Platform
- AI/ML Platform
- Logistics Platform

without major redesign.

The system MUST survive:
- flash-sale spikes
- retry storms
- duplicate events
- partial warehouse failures
- distributed inconsistencies

WITHOUT overselling or corrupting orders.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
Strict consistency required for orders + inventory.