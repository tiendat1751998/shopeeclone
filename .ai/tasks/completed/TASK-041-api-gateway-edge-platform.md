# TASK-041 — API GATEWAY & EDGE PLATFORM

## Goal

Build a REAL production-grade API Gateway & Edge Platform.

This platform is responsible for:
- global request routing
- authentication & authorization enforcement
- rate limiting (global + per-service + per-user)
- WAF-like protection (abuse, bots, attacks)
- request shaping & transformation
- API versioning control
- service discovery routing
- multi-region traffic steering
- circuit breaking at edge
- observability injection (trace, metrics, logs)

This is NOT a simple reverse proxy.

The API Gateway is the:
## SYSTEM-WIDE ENTRY CONTROL PLANE

The system MUST support:
- ultra-high throughput traffic routing
- global traffic governance
- distributed backend routing
- Kubernetes-native deployment
- observability-first architecture
- edge resilience under attack

The architecture MUST prioritize:
- latency (EXTREMELY LOW OVERHEAD)
- security enforcement
- traffic control correctness
- system-wide stability protection
- scalability under attack conditions

---

## Tech Stack

Use:
- Golang (high-performance gateway core)
- Envoy (optional data plane integration)
- Redis Cluster (rate limiting + session + dedup)
- Kafka (audit + traffic events)
- PostgreSQL (routing rules + config)
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Optional:
- Nginx / Envoy / Kong concepts (but custom gateway core required)
- eBPF (advanced traffic inspection)
- Lua/WASM plugins (extensibility layer)

---

## Core Responsibilities

The Gateway MUST support:

### Request Routing
- service discovery routing
- path-based routing
- header-based routing
- weighted routing (canary / A-B testing)
- multi-region routing decisions

### Authentication & Authorization
- JWT validation
- mTLS validation (internal traffic optional extension)
- API key validation
- RBAC enforcement hooks

### Rate Limiting (CRITICAL)
- per-user rate limit
- per-IP rate limit
- per-service rate limit
- global circuit protection

Generate:
- token bucket / leaky bucket algorithms
- distributed rate limit coordination via Redis

### WAF / Abuse Protection
- bot detection hooks
- request anomaly detection
- IP reputation blocking
- payload validation

### Request Transformation
- header injection (trace IDs)
- request normalization
- API version rewriting
- payload validation

### Observability Injection
- trace propagation (OpenTelemetry)
- correlation ID injection
- metrics tagging
- structured logging context

---

## Architecture Requirements

The Gateway MUST:
- be ultra-low latency (<5–20ms overhead target)
- be horizontally scalable
- be stateless where possible
- support distributed configuration updates

The Gateway MUST:
- act as system-wide failure isolation layer
- prevent cascading failures
- degrade gracefully under overload

Use:
- plugin-based architecture
- middleware pipeline design
- config-driven routing engine
- dependency injection

The system MUST tolerate:
- backend service outages
- Redis rate limiter failure (fallback mode)
- Kafka lag
- partial region outages
- traffic spikes (10x–100x)
- malicious request floods

---

## Folder Structure

Generate:

platforms/api-gateway/
├── cmd/
├── internal/
│   ├── config/
│   ├── routing/
│   ├── auth/
│   ├── authz/
│   ├── ratelimit/
│   ├── waf/
│   ├── middleware/
│   ├── proxy/
│   ├── discovery/
│   ├── transform/
│   ├── circuitbreaker/
│   ├── loadbalancer/
│   ├── observability/
│   ├── plugins/
│   ├── resilience/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── plugins/
├── policies/
├── configs/
├── deployments/
├── charts/
├── tests/
└── Dockerfile

---

## Routing Engine Requirements

Support:
- path routing (/api/v1/*)
- service routing (service registry based)
- weighted routing (canary deploys)
- geo-routing (multi-region)

Generate:
- routing DSL
- dynamic config reload system
- fallback routing policies

---

## Rate Limiting Requirements (CRITICAL)

Support:
- per-user limits
- per-IP limits
- per-route limits
- global system protection

Generate:
- distributed token bucket
- Redis-based coordination
- local fallback limiter (fail-open/close strategy)

---

## Authentication & Authorization

Support:
- JWT validation
- API key validation
- RBAC hooks
- service-to-service identity

Generate:
- auth middleware
- token validation pipeline
- identity propagation layer

---

## WAF / Abuse Protection

Support:
- request filtering rules
- IP blocking
- bot heuristics
- payload inspection hooks

Generate:
- rule engine
- blocking pipeline
- anomaly detection hooks (integration-ready)

---

## Load Balancing Requirements

Support:
- round robin
- least latency
- weighted distribution
- health-aware routing

Generate:
- load balancer engine
- health check integration
- dynamic endpoint updates

---

## Circuit Breaker Requirements

Support:
- per-service circuit breaker
- adaptive failure detection
- fallback routing

Generate:
- breaker state machine
- failure threshold engine
- recovery logic

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logs
- distributed tracing
- correlation IDs

Metrics:
- request latency (p50/p95/p99)
- request throughput
- error rate per service
- rate limit hits
- circuit breaker triggers

Logs:
- JSON structured logs
- request lifecycle logs
- correlation IDs

Never log sensitive payload data.

---

## Redis Requirements

Use Redis for:
- rate limiting counters
- distributed coordination
- short-lived auth/session caching

Must support:
- failover fallback mode
- TTL safety
- atomic operations (Lua scripts recommended)

---

## Kafka Requirements

Use Kafka for:
- request audit logs
- traffic analytics
- abuse detection signals

Generate:
- topic design
- partition strategy
- DLQ handling

---

## Security Requirements

The Gateway MUST:
- enforce authentication globally
- validate all incoming requests
- block malformed traffic early
- isolate internal services

Never:
- trust backend responses blindly
- bypass auth for internal routing
- expose internal service topology

---

## Resilience Requirements

Implement:
- retries (careful, idempotent only)
- timeout enforcement
- circuit breakers
- bulkhead isolation
- graceful degradation

Critical:
## GATEWAY MUST PREVENT CASCADING FAILURE ACROSS ENTIRE SYSTEM

---

## Kubernetes Requirements

Generate:
- Deployments
- HPA
- PDB
- ConfigMaps
- Secrets
- Service
- Ingress
- Helm charts

Must support:
- horizontal scaling under traffic spikes
- zero downtime rollout
- edge traffic burst handling

---

## CI/CD Requirements

Generate:
- build pipelines
- load testing gates
- security scanning
- config validation
- canary deployment support

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- load tests
- abuse simulation tests
- failover tests
- rate limit stress tests

Test scenarios:
- flash sale traffic spike
- DDoS-like traffic flood
- backend service failure cascade
- Redis outage fallback
- multi-region routing failure
- circuit breaker activation

---

## Output Requirements

Explain:
- gateway architecture
- routing strategy
- rate limiting strategy
- security enforcement model
- circuit breaker design
- multi-region traffic strategy
- failure isolation strategy

Generate production-grade code only.

No toy gateway.
No naive proxy implementation.
No unsafe routing logic.

---

## Acceptance Criteria

The API Gateway must support integration with:
- OMS Platform
- Payment Platform
- Fraud Platform
- Recommendation Platform
- Search Platform
- Notification Platform
- AI/ML Platform

without redesign.

The system MUST survive:
- massive traffic spikes
- DDoS-like conditions
- backend outages
- partial region failures
- Redis/Kafka degradation

WITHOUT cascading failure propagation.

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.
Strict edge reliability required.