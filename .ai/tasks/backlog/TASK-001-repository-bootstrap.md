# TASK-001 — REPOSITORY BOOTSTRAP + SHARED FOUNDATION

## Goal

Build the REAL production-grade engineering foundation for the platform.

This task is responsible for:
- monorepo structure
- shared libraries
- platform conventions
- observability foundation
- logging foundation
- tracing foundation
- metrics foundation
- middleware foundation
- CI/CD foundation
- Docker foundation
- Kubernetes foundation
- developer tooling foundation

This task does NOT implement business services.

This task builds the platform foundation ALL future services depend on.

This is NOT a toy starter template.

---

## Tech Stack

Use:
- Golang
- gRPC
- Protobuf
- OpenTelemetry
- Prometheus
- Docker
- Kubernetes
- Helm
- Makefile
- GitHub Actions or Drone CI

Optional:
- Buf for proto management
- golangci-lint
- pre-commit hooks

---

## Monorepo Requirements

Generate a REAL production-grade monorepo structure.

The repository MUST support:
- multiple services
- shared libraries
- shared middleware
- shared observability
- shared protobuf contracts
- future scaling

Generate:

platform/
├── .ai/
├── .github/
├── build/
├── deployments/
├── docs/
├── scripts/
├── tools/
├── proto/
├── shared/
├── services/
├── tests/
├── Makefile
├── docker-compose.yml
├── go.work
└── README.md

---

## Shared Libraries Requirements

Generate shared libraries for:

shared/
├── config/
├── logger/
├── tracing/
├── metrics/
├── middleware/
├── errors/
├── validation/
├── response/
├── auth/
├── cache/
├── redis/
├── mysql/
├── grpc/
├── http/
├── resilience/
├── pagination/
├── idempotency/
├── observability/
├── security/
└── testing/

Shared libraries MUST:
- be reusable
- avoid tight coupling
- support future services
- support high observability
- support resilience patterns

This is NOT utility-folder spaghetti.

---

## Configuration System Requirements

Generate:
- environment-based configuration
- local configs
- production configs
- Kubernetes configs
- validation system
- configuration loader
- secret integration support

Requirements:
- no hardcoded secrets
- typed configuration
- startup validation
- environment overrides

Support:
- local development
- Docker
- Kubernetes
- CI/CD pipelines

---

## Logging Foundation Requirements

Generate structured logging system.

Requirements:
- JSON logs
- correlation IDs
- trace IDs
- request IDs
- service metadata
- environment metadata

Support:
- distributed tracing correlation
- Kubernetes logging
- centralized logging systems

Generate:
- logger abstraction
- middleware integration
- context propagation

No fake logging wrapper.

---

## Tracing Foundation Requirements

Generate:
- OpenTelemetry setup
- trace propagation
- gRPC tracing
- HTTP tracing
- middleware integration
- exporter configuration

Requirements:
- distributed tracing
- context propagation
- service correlation
- trace sampling support

Support:
- Jaeger
- Tempo
- OTLP exporters

---

## Metrics Foundation Requirements

Generate:
- Prometheus metrics
- HTTP metrics
- gRPC metrics
- DB metrics
- Redis metrics
- custom business metrics support

Requirements:
- low-overhead metrics
- histogram support
- latency metrics
- error metrics

Generate:
- metrics registry
- middleware integration
- instrumentation helpers

---

## Middleware Foundation Requirements

Generate reusable middleware for:
- logging
- tracing
- metrics
- panic recovery
- timeout handling
- request IDs
- correlation IDs
- RBAC hooks
- validation
- rate limiting hooks

Support:
- HTTP
- gRPC

Middleware MUST:
- be production-grade
- support observability
- support resilience
- avoid tight coupling

---

## Protobuf Structure Requirements

Generate production-grade proto structure:

proto/
├── common/
├── auth/
├── user/
├── product/
├── inventory/
├── order/
├── payment/
└── shared/

Requirements:
- versioning support
- backward compatibility awareness
- shared contracts
- reusable messages
- common error contracts

Support:
- Buf
- code generation
- gRPC services

Generate:
- proto generation scripts
- Makefile integration

---

## Docker Foundation Requirements

Generate:
- production-grade Dockerfiles
- multi-stage builds
- shared base images
- security hardening
- non-root containers

Generate:
- docker-compose local environment
- shared network structure
- local infra setup

Support:
- local development
- CI pipelines
- Kubernetes deployment

---

## Kubernetes Foundation Requirements

Generate:
- shared Helm structure
- base manifests
- common labels
- common annotations
- observability integration
- ingress conventions

Support:
- rolling deployments
- autoscaling
- canary deployment
- service monitoring

Generate:
- namespace strategy
- resource conventions
- deployment conventions

---

## CI/CD Foundation Requirements

Generate:
- GitHub Actions or Drone pipelines
- lint pipelines
- test pipelines
- security scanning
- Docker build pipelines
- Helm validation
- proto validation

Requirements:
- reusable workflows
- caching strategy
- parallel pipelines

Generate:
- branch strategy recommendations
- CI conventions
- artifact conventions

---

## Makefile Requirements

Generate production-grade Makefile targets:

- build
- test
- lint
- proto
- docker
- run
- dev
- migrate
- generate
- benchmark
- security-scan

Requirements:
- developer-friendly
- composable
- scalable

---

## Linting & Quality Requirements

Generate:
- golangci-lint config
- formatting rules
- proto linting
- import organization
- pre-commit hooks

Requirements:
- enterprise readability
- maintainability
- consistency enforcement

Support:
- CI integration
- local validation

---

## Testing Foundation Requirements

Generate:
- shared test helpers
- integration test structure
- mock conventions
- benchmark structure
- concurrency test helpers

Support:
- unit testing
- integration testing
- load testing
- resilience testing

---

## Security Requirements

Generate:
- secret handling conventions
- TLS-ready configs
- secure container defaults
- RBAC-ready middleware hooks
- security headers support

Never:
- hardcode secrets
- expose internal configs
- disable TLS by default

---

## Developer Experience Requirements

Generate:
- onboarding docs
- local development docs
- architecture docs
- Makefile shortcuts
- local infra bootstrap scripts

Requirements:
- fast onboarding
- reproducible environments
- operational clarity

---

## Observability Requirements

The platform foundation MUST support:
- OpenTelemetry
- Prometheus
- structured logging
- distributed tracing
- correlation IDs

Every future service MUST inherit observability automatically.

---

## Reliability Requirements

Generate:
- graceful shutdown helpers
- retry helpers
- timeout helpers
- resilience primitives
- circuit breaker foundation

Support:
- distributed deployments
- Kubernetes restarts
- rolling deployments

---

## Output Requirements

Explain:
- monorepo architecture
- shared library philosophy
- observability architecture
- proto strategy
- CI/CD strategy
- Docker strategy
- Kubernetes conventions
- developer workflow
- resilience foundation

Generate production-grade code only.

No toy starter template.
No fake shared libraries.
No utility-folder spaghetti.
No fake observability.

---

## Acceptance Criteria

The foundation MUST support future implementation of:
- API Gateway
- Auth Service
- User Service
- Product Service
- Inventory Service
- Cart Service
- Order Service
- Payment Service

without major future refactors.

The repository MUST realistically support:
- large engineering teams
- distributed development
- CI/CD scaling
- Kubernetes-native deployment
- production observability

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