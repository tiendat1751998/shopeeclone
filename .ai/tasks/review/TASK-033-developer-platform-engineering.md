# TASK-033 — DEVELOPER PLATFORM & INTERNAL ENGINEERING PLATFORM

## Goal

Build a REAL production-grade Developer Platform & Internal Engineering Platform.

This platform is responsible for:
- internal developer platform
- self-service infrastructure
- golden paths
- scaffolding systems
- developer automation
- platform APIs
- service provisioning
- deployment orchestration
- engineering governance
- developer observability

This is NOT a toy admin dashboard.

The Developer Platform must support:
- self-service provisioning
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe provisioning

The architecture MUST prioritize:
- developer productivity
- platform consistency
- operational stability
- governance correctness
- secure automation

---

## Tech Stack

Use:
- Golang
- Kubernetes
- Backstage
- ArgoCD
- Terraform
- Crossplane
- PostgreSQL
- Redis Cluster
- Kafka
- OpenTelemetry
- Prometheus
- Helm

Optional:
- Crossplane compositions
- OPA/Gatekeeper
- Kyverno
- Temporal
- Tekton

---

## Core Responsibilities

The Developer Platform MUST support:

### Self-Service Infrastructure
- environment provisioning
- database provisioning
- Kafka topic provisioning
- Redis provisioning

### Golden Paths
- production service templates
- standardized deployment flows
- production-ready scaffolding
- operational defaults

### Scaffolding Systems
- service generators
- infrastructure generators
- replay-safe provisioning workflows
- dependency orchestration

### Developer Automation
- deployment automation
- environment lifecycle automation
- rollback automation
- replay-safe workflow execution

### Platform APIs
- internal provisioning APIs
- deployment APIs
- platform orchestration APIs
- governance APIs

### Engineering Governance
- policy enforcement
- production guardrails
- operational standards
- audit orchestration

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate provisioning/deployment/governance
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The Developer Platform MUST:
- support replay-safe provisioning
- support distributed orchestration
- support realtime deployment visibility
- support degraded automation fallback

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The platform MUST tolerate:
- retry storms
- duplicate provisioning events
- delayed infrastructure reconciliation
- partial deployment failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/developer-platform/
├── cmd/
├── internal/
│   ├── config/
│   ├── provisioning/
│   ├── deployments/
│   ├── templates/
│   ├── governance/
│   ├── automation/
│   ├── orchestration/
│   ├── environments/
│   ├── synchronization/
│   ├── replay/
│   ├── backstage/
│   ├── workflows/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── templates/
│   ├── microservices/
│   ├── kafka/
│   ├── redis/
│   ├── databases/
│   └── kubernetes/
│
├── crossplane/
├── terraform/
├── argocd/
├── backstage/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Backstage Requirements

Use Backstage for:
- service catalog
- golden paths
- developer portal
- platform orchestration

Generate:
- software templates
- platform plugins
- deployment visibility
- replay-safe provisioning workflows

Requirements:
- production-grade RBAC
- platform governance
- distributed visibility

Never:
- expose infrastructure secrets
- allow unrestricted provisioning
- bypass governance

---

## Provisioning Requirements

Support:
- Kubernetes namespace provisioning
- PostgreSQL provisioning
- Redis provisioning
- Kafka provisioning

Generate:
- provisioning orchestration
- replay-safe provisioning workflows
- distributed reconciliation
- dependency-aware provisioning

The provisioning system MUST tolerate:
- duplicate provisioning requests
- replay storms
- partial infrastructure failures

No fake provisioning architecture.

---

## Deployment Orchestration Requirements

Support:
- GitOps deployment
- canary deployments
- rollback orchestration
- replay-safe deployment execution

Generate:
- deployment orchestration
- rollout coordination
- replay-safe synchronization
- distributed deployment visibility

No naive deployment orchestration.

---

## Governance Requirements

Support:
- policy enforcement
- security standards
- production guardrails
- replay-safe governance synchronization

Generate:
- OPA/Kyverno policies
- governance workflows
- policy validation pipelines
- audit orchestration

The governance system MUST be realistic.

---

## Template/Golden Path Requirements

Support:
- production-ready templates
- service scaffolding
- infrastructure scaffolding
- replay-safe template orchestration

Generate:
- service templates
- Kubernetes templates
- CI/CD templates
- observability templates

The generated templates MUST:
- be production-grade
- include observability
- include resilience
- include security defaults

No toy templates.

---

## Workflow Automation Requirements

Support:
- deployment workflows
- rollback workflows
- provisioning workflows
- replay-safe workflow execution

Generate:
- workflow orchestration
- distributed execution
- replay-safe coordination
- failure recovery workflows

Use:
- Temporal or equivalent if appropriate

---

## Kafka Requirements

Use Kafka for:
- provisioning events
- deployment events
- governance synchronization
- replay-safe orchestration

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
- distributed automation
- async synchronization
- operational orchestration

---

## PostgreSQL Requirements

Use PostgreSQL for:
- provisioning metadata
- deployment metadata
- governance metadata
- replay metadata

Generate:
- optimized schemas
- indexes
- immutable audit tables
- replay-safe synchronization tables

Requirements:
- replay safety
- operational correctness
- audit consistency

---

## Redis Requirements

Use Redis for:
- workflow coordination
- deployment caching
- replay coordination
- distributed throttling

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- high automation throughput
- distributed workflows
- low-latency coordination

---

## Event-Driven Requirements

Generate events for:
- provisioning started
- deployment completed
- rollback triggered
- governance violation detected
- template generated
- environment provisioned

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
- distributed automation
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
- /platform/provision
- /platform/deployments
- /platform/templates
- /platform/governance
- /platform/environments

Support:
- self-service provisioning
- replay-safe synchronization
- deployment visibility
- operational auditing

---

## Security Requirements

The platform MUST:
- validate developer ownership
- enforce RBAC
- isolate provisioning workflows
- sanitize infrastructure input
- protect automation integrity

Never:
- expose cluster-admin blindly
- expose infrastructure secrets
- trust provisioning input blindly
- allow unrestricted deployments

Generate:
- authorization middleware
- replay validation
- provisioning isolation
- audit pipelines

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- provisioning latency
- deployment latency
- rollback count
- governance violations
- replay latency
- workflow concurrency

Logs:
- JSON structured logs
- trace IDs
- correlation IDs

Never log infrastructure secrets.

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
- Deployments
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
- GitHub Actions or Drone pipelines
- Helm validation
- Kubernetes policy validation
- template validation
- vulnerability scanning
- GitOps workflows

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- replay tests
- provisioning tests
- governance tests
- workflow tests
- concurrency tests

Test:
- duplicate provisioning events
- replay storms
- partial deployment failures
- governance violations
- distributed deployments
- infrastructure reconciliation correctness

---

## Output Requirements

Explain:
- developer-platform architecture
- provisioning strategy
- governance strategy
- replay-safe automation strategy
- template strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy internal platform.
No fake provisioning architecture.
No naive governance logic.

---

## Acceptance Criteria

The Developer Platform must support future integration with:
- PRE/SRE Platform
- Global Infrastructure Platform
- Analytics Platform
- Fraud Platform
- Live Commerce Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- duplicate provisioning events
- partial deployment failures
- distributed deployments
- high internal automation throughput

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.