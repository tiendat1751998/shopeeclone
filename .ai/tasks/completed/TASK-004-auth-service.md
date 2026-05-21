# TASK-004 — AUTHENTICATION SERVICE

## Goal

Build a REAL production-grade Authentication Service.

This service is responsible for:
- user authentication
- JWT lifecycle
- refresh token rotation
- RBAC
- session management
- device session tracking
- audit logging
- account protection
- token validation
- authentication security

This is NOT a toy auth service.

The service must support:
- high concurrency
- horizontal scaling
- Kubernetes-native deployment
- distributed tracing
- observability-first design
- fault tolerance
- secure-by-default architecture

---

## Tech Stack

Use:
- Golang
- Gin/Fiber
- gRPC
- MySQL
- Redis Cluster
- OpenTelemetry
- Prometheus
- Kubernetes
- Helm

Password hashing:
- Argon2id or bcrypt

Token format:
- JWT

---

## Core Responsibilities

The Auth Service MUST support:

### Authentication
- register
- login
- logout
- token refresh
- token revoke
- password reset
- email verification

### Session Management
- device sessions
- refresh token rotation
- session invalidation
- session expiration
- concurrent session handling

### Authorization
- RBAC
- permission validation
- role hierarchy
- claims propagation

### Security
- brute-force protection
- rate limiting
- suspicious login detection
- token blacklist
- replay attack prevention
- audit logging

---

## Architecture Requirements

The service MUST:
- follow clean architecture
- separate domain/application/infrastructure
- support stateless scaling
- support eventual consistency where appropriate
- support distributed deployments

Use:
- middleware architecture
- dependency injection
- repository pattern carefully
- service isolation

---

## Folder Structure

Generate:

services/auth/
├── cmd/
├── internal/
│   ├── config/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── transport/
│   ├── middleware/
│   ├── auth/
│   ├── jwt/
│   ├── session/
│   ├── rbac/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   ├── security/
│   ├── audit/
│   ├── validation/
│   └── health/
│
├── migrations/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## JWT Requirements

Implement:
- access tokens
- refresh tokens
- token rotation
- token revocation
- token expiration
- secure claims
- device-aware tokens

Requirements:
- short-lived access tokens
- rotating refresh tokens
- secure signing
- Redis blacklist support

Never:
- store plaintext tokens
- trust unsigned tokens
- allow refresh token reuse

---

## Session Requirements

Use Redis for:
- session storage
- refresh token metadata
- token blacklist
- login attempt tracking
- rate limiting
- device sessions

Support:
- multiple devices
- session invalidation
- force logout
- session expiration

Generate:
- Redis session layer
- retry handling
- timeout handling
- connection pooling

---

## RBAC Requirements

Implement:
- roles
- permissions
- permission inheritance
- role hierarchy
- claims propagation

Support:
- admin
- seller
- buyer
- internal services

---

## Database Requirements

Use MySQL for:
- users
- credentials
- audit logs
- device metadata
- role mappings

Generate:
- migrations
- indexes
- constraints
- optimized queries

Rules:
- never use SELECT *
- always paginate
- avoid N+1 queries

---

## Security Requirements

MUST include:
- password hashing
- replay prevention
- brute-force protection
- suspicious login detection
- audit logging
- secure headers
- request validation
- sanitization

Implement:
- login throttling
- account lockout
- suspicious IP detection
- correlation IDs

Never:
- log passwords
- expose sensitive errors
- hardcode secrets
- trust client claims

---

## Audit Requirements

Generate audit logs for:
- login
- logout
- token refresh
- failed login
- password reset
- role changes
- session invalidation

Logs must include:
- trace ID
- user ID
- device ID
- IP
- timestamp

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- request correlation

Metrics:
- login attempts
- failed logins
- refresh requests
- session count
- token validation latency
- Redis latency
- DB latency

---

## API Requirements

Generate:
- REST APIs
- gRPC APIs
- OpenAPI specs
- proto files

Endpoints:
- /register
- /login
- /logout
- /refresh
- /verify-email
- /reset-password
- /sessions
- /revoke

---

## Reliability Requirements

Implement:
- retries
- timeout handling
- graceful shutdown
- panic recovery
- circuit breaker support

Support:
- horizontal scaling
- rolling deployment
- Kubernetes autoscaling

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
- rolling deployment
- autoscaling

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
- auth flow tests
- JWT tests
- RBAC tests
- Redis tests
- middleware tests
- integration tests

Test:
- replay attacks
- refresh reuse
- brute force
- concurrent sessions

---

## Output Requirements

Explain:
- auth architecture
- JWT lifecycle
- refresh rotation strategy
- Redis session strategy
- RBAC model
- scaling strategy
- audit strategy
- observability flow
- security hardening

Generate production-grade code only.

No toy auth.
No fake JWT implementation.
No insecure flows.

---

## Acceptance Criteria

The Auth Service must be capable of supporting:
- API Gateway
- User Service
- Seller Service
- Product Service
- Order Service
- Payment Service

without major future refactors.

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