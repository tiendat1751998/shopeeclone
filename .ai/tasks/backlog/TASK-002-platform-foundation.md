# TASK-002 — PLATFORM FOUNDATION

## Goal

Build the production-grade Kubernetes and observability foundation for the platform.

This task establishes:
- Kubernetes base infrastructure
- observability stack
- ingress
- service mesh
- monitoring
- logging
- tracing
- GitOps
- secrets management
- local development platform

This is NOT optional.
All future services will depend on this platform foundation.

---

## Requirements

The infrastructure must support:
- Kubernetes-native deployment
- microservices scalability
- observability-first architecture
- GitOps workflow
- production-grade security
- horizontal scaling
- fault tolerance

---

## Infrastructure Components

### Kubernetes
- namespace structure
- resource quotas
- limit ranges
- network policies
- RBAC

### Ingress
- NGINX Ingress or APISIX
- TLS support
- rate limiting
- gzip
- security headers

### Service Mesh
- Istio
- mTLS
- traffic policies
- retries
- circuit breaking

### Observability
- Prometheus
- Grafana
- Loki
- Tempo
- OpenTelemetry Collector

### GitOps
- ArgoCD

### Secrets
- External Secrets or Vault

### Storage
- MinIO

### Messaging
- NATS JetStream or Kafka

### Cache
- Redis Cluster

---

## Deliverables

Generate:
- folder structure
- Helm charts
- Kubernetes manifests
- ArgoCD applications
- observability dashboards
- Prometheus rules
- alerting rules
- Loki configs
- Tempo configs
- OpenTelemetry configs
- ingress configs
- Istio configs
- Redis deployment
- NATS/Kafka deployment
- MinIO deployment

---

## Kubernetes Structure

Expected structure:

platform/
├── deployments/
│   ├── helm/
│   ├── kubernetes/
│   ├── argocd/
│   ├── monitoring/
│   ├── logging/
│   ├── tracing/
│   ├── ingress/
│   ├── istio/
│   ├── redis/
│   ├── nats/
│   ├── kafka/
│   └── minio/
│
├── infrastructure/
│   ├── terraform/
│   ├── scripts/
│   └── environments/
│
└── .ai/

---

## Security Requirements

Must include:
- mTLS
- network policies
- RBAC
- pod security
- non-root containers
- secrets management
- TLS everywhere

---

## Observability Requirements

Every future service MUST support:
- tracing
- metrics
- structured logs
- correlation IDs

Generate shared observability standards.

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone CI
- image scanning
- linting
- testing
- Helm validation
- Kubernetes validation

---

## Scalability Requirements

The platform must support:
- horizontal scaling
- auto-scaling
- rolling deployment
- canary deployment
- blue-green deployment

---

## Reliability Requirements

Must include:
- retries
- circuit breaking
- pod disruption budgets
- anti-affinity
- readiness/liveness probes

---

## Output Requirements

Explain:
- architecture decisions
- Kubernetes topology
- observability flow
- GitOps flow
- networking flow
- scaling strategy
- security hardening

Generate production-grade configs only.

No toy examples.
No fake configs.
No simplified infrastructure.

---

## Acceptance Criteria

The platform foundation must be capable of supporting:
- API Gateway
- Auth Service
- Product Service
- Inventory Service
- Order Service
- Payment Service

without major architectural refactors later.

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