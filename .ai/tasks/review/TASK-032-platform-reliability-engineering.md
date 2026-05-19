# TASK-032 — PLATFORM RELIABILITY ENGINEERING (PRE/SRE PLATFORM)

## Goal

Build a REAL production-grade Platform Reliability Engineering Platform.

This platform is responsible for:
- incident management
- SLO/SLI systems
- auto-remediation
- chaos engineering
- reliability orchestration
- operational automation
- alert routing
- runbook automation
- deployment health analysis
- platform diagnostics

This is NOT a toy monitoring dashboard.

The PRE/SRE Platform must support:
- distributed observability
- automated remediation
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe incident ingestion

The architecture MUST prioritize:
- operational resiliency
- incident correctness
- alert quality
- recovery automation
- platform stability

---

## Tech Stack

Use:
- Golang
- Kubernetes
- Prometheus
- Alertmanager
- Thanos
- Loki
- Tempo
- Kafka
- PostgreSQL
- Redis Cluster
- OpenTelemetry
- Argo Rollouts
- Helm

Optional:
- Grafana OnCall
- Cortex/Mimir
- Chaos Mesh
- LitmusChaos
- eBPF observability

---

## Core Responsibilities

The PRE/SRE Platform MUST support:

### Incident Management
- incident orchestration
- escalation workflows
- incident timelines
- replay-safe incident ingestion

### SLO/SLI Management
- SLI aggregation
- SLO evaluation
- error budget tracking
- reliability scoring

### Auto-Remediation
- automated rollback
- pod remediation
- service restart orchestration
- replay-safe remediation execution

### Chaos Engineering
- fault injection
- resilience validation
- replay-safe chaos orchestration
- blast-radius isolation

### Alert Routing
- distributed alert routing
- deduplication
- alert suppression
- escalation policies

### Deployment Health Analysis
- canary analysis
- rollout validation
- deployment anomaly detection
- rollback orchestration

---

## Architecture Requirements

The platform MUST:
- follow clean architecture
- separate observability/remediation/analysis
- support distributed deployments
- support eventual consistency
- support event-driven workflows

The PRE/SRE Platform MUST:
- support replay-safe incident ingestion
- support distributed remediation
- support realtime alerting
- support degraded observability fallback

Use:
- CQRS where appropriate
- dependency injection
- modular architecture
- resilience patterns

The reliability system MUST tolerate:
- retry storms
- duplicate incident events
- delayed observability ingestion
- partial remediation failures
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

platforms/pre-sre/
├── cmd/
├── internal/
│   ├── config/
│   ├── incidents/
│   ├── slo/
│   ├── sli/
│   ├── remediation/
│   ├── chaos/
│   ├── alerts/
│   ├── deployments/
│   ├── diagnostics/
│   ├── analysis/
│   ├── synchronization/
│   ├── replay/
│   ├── observability/
│   ├── cache/
│   ├── events/
│   ├── metrics/
│   ├── tracing/
│   ├── logging/
│   └── health/
│
├── chaos/
├── runbooks/
├── deployments/
├── charts/
├── tests/
├── configs/
└── Dockerfile

---

## Prometheus/Thanos Requirements

Use Prometheus + Thanos for:
- global metrics aggregation
- SLI computation
- distributed alerting
- replay-safe metrics ingestion

Generate:
- federation topology
- recording rules
- alert rules
- replay-safe aggregation

Requirements:
- high availability
- global metric aggregation
- replay-safe observability

Never:
- use single Prometheus instance
- ignore long-term retention
- ignore regional observability

---

## Loki/Tempo Requirements

Use Loki + Tempo for:
- distributed logs
- distributed tracing
- incident diagnostics
- replay-safe trace correlation

Generate:
- retention strategies
- correlation workflows
- replay-safe trace ingestion
- distributed diagnostics

Requirements:
- trace-log correlation
- cross-region diagnostics
- high-ingestion tolerance

---

## Incident Management Requirements

Support:
- incident creation
- escalation workflows
- replay-safe incident synchronization
- operational auditing

Generate:
- incident orchestration
- escalation pipelines
- distributed coordination
- replay-safe recovery workflows

The incident system MUST tolerate:
- duplicate alerts
- replay storms
- delayed ingestion

No fake incident architecture.

---

## SLO/SLI Requirements

Support:
- error budget computation
- SLO evaluation
- service reliability scoring
- burn-rate alerting

Generate:
- SLO computation pipelines
- replay-safe aggregation
- distributed evaluation workflows
- alert orchestration

No naive SLO logic.

---

## Auto-Remediation Requirements

Support:
- automated rollback
- pod remediation
- deployment pausing
- replay-safe remediation

Generate:
- remediation workflows
- rollback orchestration
- replay-safe action pipelines
- distributed remediation coordination

The remediation system MUST tolerate:
- replay storms
- partial remediation failures
- cascading failures

---

## Chaos Engineering Requirements

Support:
- network chaos
- pod chaos
- latency injection
- replay-safe chaos execution

Generate:
- chaos experiments
- blast-radius isolation
- replay-safe orchestration
- resilience validation workflows

The chaos system MUST be production-grade.

---

## Alerting Requirements

Support:
- alert deduplication
- escalation routing
- replay-safe alert ingestion
- alert suppression

Generate:
- routing policies
- deduplication workflows
- escalation orchestration
- replay-safe synchronization

The alerting system MUST tolerate:
- alert floods
- replay storms
- regional outages

---

## Deployment Analysis Requirements

Support:
- canary analysis
- rollout validation
- anomaly detection
- replay-safe deployment analysis

Generate:
- deployment analysis pipelines
- rollback orchestration
- distributed validation workflows
- health scoring

---

## Kafka Requirements

Use Kafka for:
- incident events
- observability synchronization
- remediation events
- replay-safe ingestion

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
- distributed observability
- async synchronization
- operational automation

---

## PostgreSQL Requirements

Use PostgreSQL for:
- incident metadata
- SLO metadata
- remediation metadata
- replay metadata

Generate:
- optimized schemas
- indexes
- immutable incident audit tables
- replay-safe synchronization tables

Requirements:
- replay safety
- operational correctness
- audit consistency

---

## Redis Requirements

Use Redis for:
- hot alert cache
- remediation coordination
- replay coordination
- distributed throttling

Generate:
- TTL strategy
- replay protection
- cache invalidation
- distributed coordination

Support:
- high alert throughput
- distributed remediation
- low-latency coordination

---

## Event-Driven Requirements

Generate events for:
- incident triggered
- SLO violated
- remediation executed
- chaos experiment started
- deployment degraded
- rollback triggered

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
- distributed remediation
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
- /incidents
- /slos
- /alerts
- /remediations
- /deployments/analysis

Support:
- realtime incident updates
- replay-safe synchronization
- distributed diagnostics
- operational auditing

---

## Security Requirements

The platform MUST:
- validate operational access
- enforce RBAC
- isolate remediation workflows
- sanitize operational input
- protect incident integrity

Never:
- expose infrastructure secrets
- expose internal remediation tokens
- trust external observability blindly
- allow unrestricted remediation execution

Generate:
- authorization middleware
- replay validation
- remediation isolation
- operational auditing

---

## Observability Requirements

Generate:
- OpenTelemetry tracing
- Prometheus metrics
- structured logging
- distributed tracing
- correlation IDs

Metrics:
- incident rate
- remediation latency
- alert throughput
- SLO violations
- replay latency
- deployment rollback count

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
- chaos validation
- vulnerability scanning
- GitOps workflows

---

## Testing Requirements

Generate:
- unit tests
- integration tests
- replay tests
- chaos tests
- remediation tests
- alerting tests
- concurrency tests

Test:
- duplicate incident events
- replay storms
- alert floods
- cascading failures
- regional outages
- distributed deployments
- remediation correctness

---

## Output Requirements

Explain:
- reliability architecture
- incident-management strategy
- SLO strategy
- replay-safe remediation strategy
- chaos-engineering strategy
- scaling strategy
- resilience strategy

Generate production-grade code only.

No toy monitoring setup.
No fake remediation architecture.
No naive alerting logic.

---

## Acceptance Criteria

The PRE/SRE Platform must support future integration with:
- Global Infrastructure Platform
- Analytics Platform
- Fraud Platform
- Live Commerce Platform
- Advertising Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- alert floods
- cascading failures
- distributed deployments
- regional outages

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.