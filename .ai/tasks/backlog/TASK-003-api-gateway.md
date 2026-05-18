# TASK-003 — API GATEWAY

## Goal

Build a REAL production-grade API Gateway for the platform.

The gateway will act as:
- single external entrypoint
- authentication boundary
- rate limiting layer
- observability boundary
- API routing layer
- security enforcement layer

This gateway must support:
- high concurrency
- horizontal scaling
- Kubernetes-native deployment
- distributed tracing
- traffic control
- future microservices integration

This is NOT a toy API gateway.

---

## Tech Stack

Use:
- Golang
- Gin or Fiber
- gRPC
- Redis
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm
- Istio-compatible architecture

---

## Core Responsibilities

The API Gateway MUST handle:

### Authentication
- JWT validation
- access token validation
- refresh token flow integration
- RBAC support
- session validation
- device metadata propagation

### Security
- rate limiting
- IP throttling
- security headers
- request validation
- body size limits
- anti-abuse middleware
- request sanitization
- correlation IDs

### Routing
- REST routing
- gRPC routing
- upstream service discovery
- load balancing
- retry handling
- timeout handling

### Observability
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- request tracing
- distributed correlation IDs

### Reliability
- retries
- circuit breakers
- graceful shutdown
- timeout propagation
- context propagation

---

## Architecture Requirements

The gateway MUST:
- be stateless
- support horizontal scaling
- support Kubernetes autoscaling
- support rolling deployment
- support canary deployment

Use:
- middleware architecture
- dependency injection
- modular route registration

---

## Folder Structure

Generate a production-grade structure:

services/gateway/
├── cmd/
├── internal/
│   ├── config/
│   ├── middleware/
│   ├── transport/
│   ├── routing/
│   ├── auth/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   ├── discovery/
│   ├── ratelimit/
│   ├── resilience/
│   └── health/
│
├── deployments/
├── charts/
├── scripts/
├── configs/
├── tests/
└── Dockerfile

---

## Middleware Requirements

Generate middleware for:

- request logging
- request tracing
- metrics
- JWT validation
- RBAC
- panic recovery
- rate limiting
- request ID generation
- correlation IDs
- timeout handling
- body validation
- CORS
- security headers

---

## Redis Requirements

Use Redis for:
- distributed rate limiting
- session validation
- token blacklist
- cache
- gateway metadata

Generate:
- Redis integration layer
- retry handling
- connection pooling
- timeout handling

---

## Routing Requirements

The gateway must support:

### REST
- REST API proxying
- request transformation
- response transformation

### gRPC
- gRPC upstream integration
- context propagation
- metadata forwarding

---

## Observability Requirements

Generate:
- OpenTelemetry integration
- Prometheus metrics
- structured logging
- trace propagation
- request latency metrics
- error metrics
- upstream metrics

Metrics must include:
- request count
- request duration
- error rate
- upstream latency
- retry count
- rate limit hits

---

## Kubernetes Requirements

Generate:
- Deployment
- Service
- Ingress
- ConfigMap
- Secret integration
- HPA
- PodDisruptionBudget
- NetworkPolicy
- ServiceMonitor
- Helm chart

Support:
- rolling deployment
- canary deployment
- autoscaling

---

## Security Requirements

The gateway MUST:
- never trust client input
- sanitize requests
- validate JWTs
- enforce RBAC
- limit request size
- prevent abuse
- support TLS
- support mTLS internally

Generate:
- security headers
- anti-abuse middleware
- IP filtering support
- rate limiting strategy

---

## Reliability Requirements

Generate:
- retry policies
- timeout policies
- circuit breaker support
- graceful shutdown
- panic recovery
- upstream failover strategy

---

## Configuration Requirements

Generate:
- environment-based configs
- production configs
- Kubernetes configs
- local development configs

Never hardcode:
- secrets
- URLs
- credentials
- tokens

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
- middleware tests
- integration tests
- routing tests
- auth tests
- rate limiting tests

---

## Output Requirements

Explain:
- gateway architecture
- request lifecycle
- middleware flow
- tracing flow
- scaling strategy
- rate limiting strategy
- Redis strategy
- resilience strategy
- Kubernetes topology

Generate production-grade code only.

No toy implementations.
No fake middleware.
No simplified architecture.

---

## Acceptance Criteria

The gateway must be capable of serving as the entrypoint for:
- auth service
- user service
- product service
- inventory service
- order service
- payment service

without requiring major future refactors.

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