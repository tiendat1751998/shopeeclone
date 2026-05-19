# TASK-031 — GLOBAL INFRASTRUCTURE & MULTI-REGION PLATFORM

## Goal

Build a REAL production-grade Global Infrastructure & Multi-Region Platform.

This platform is responsible for:
- multi-region deployment
- global traffic routing
- DR orchestration
- geo-replication
- edge routing
- regional failover
- global observability
- cross-region synchronization
- global service discovery
- infrastructure governance

This is NOT a toy Kubernetes setup.

The Global Infrastructure Platform must support:
- multi-region deployments
- global failover
- distributed deployments
- Kubernetes-native deployment
- observability-first architecture
- replay-safe synchronization

The architecture MUST prioritize:
- regional resiliency
- failover correctness
- operational stability
- disaster recovery readiness
- infrastructure consistency

---

## Tech Stack

Use:
- Kubernetes
- Istio or Linkerd
- ArgoCD
- Crossplane
- Terraform
- Kafka
- Redis Enterprise/Cluster
- PostgreSQL
- Vitess or CockroachDB
- OpenTelemetry
- Prometheus
- Thanos
- Loki
- Tempo
- Helm

Optional:
- Consul
- Envoy
- Cloudflare
- Gateway API
- eBPF observability

---

## Core Responsibilities

The Global Infrastructure Platform MUST support:

### Multi-Region Deployment
- active-active regions
- active-passive regions
- regional workload placement
- global workload orchestration

### Global Traffic Routing
- geo-routing
- latency-aware routing
- failover routing
- edge routing

### Disaster Recovery
- automated failover
- DR orchestration
- recovery validation
- replay-safe recovery

### Geo-Replication
- database replication
- cache replication
- event replication
- object storage replication

### Global Observability
- cross-region tracing
- centralized metrics
- global alerting
- distributed diagnostics

### Infrastructure Governance
- policy enforcement
- GitOps orchestration
- global RBAC
- infrastructure auditing

---

## Architecture Requirements

The platform MUST:
- follow GitOps principles
- separate control-plane/data-plane
- support distributed deployments
- support eventual consistency
- support infrastructure-as-code

The Global Infrastructure Platform MUST:
- support replay-safe synchronization
- support regional isolation
- support global failover
- support degraded regional operation

Use:
- declarative infrastructure
- dependency injection where applicable
- modular infrastructure architecture
- resilience patterns

The infrastructure MUST tolerate:
- retry storms
- duplicate synchronization events
- regional outages
- split-brain risks
- distributed deployments
- replay storms

---

## Folder Structure

Generate:

infrastructure/
├── terraform/
│   ├── global/
│   ├── regions/
│   ├── networking/
│   ├── observability/
│   ├── databases/
│   └── security/
│
├── kubernetes/
│   ├── clusters/
│   ├── namespaces/
│   ├── policies/
│   ├── gateways/
│   ├── mesh/
│   └── observability/
│
├── argocd/
├── crossplane/
├── helm/
├── scripts/
├── runbooks/
├── disaster-recovery/
├── policies/
├── tests/
└── docs/

---

## Kubernetes Requirements

Support:
- multi-cluster orchestration
- regional isolation
- autoscaling
- workload spreading

Generate:
- cluster bootstrap
- node pools
- autoscaling policies
- workload affinity/anti-affinity
- PodDisruptionBudgets
- topology spread constraints

Requirements:
- Kubernetes-native
- production-grade HA
- rolling upgrades
- zero-downtime deployments

Never:
- deploy single-region only
- ignore zone failures
- ignore autoscaling

---

## Service Mesh Requirements

Use Istio or Linkerd for:
- mTLS
- traffic routing
- retries
- circuit breaking
- traffic shadowing

Generate:
- mesh policies
- retry policies
- circuit breakers
- canary routing
- blue/green deployments

Support:
- multi-region mesh
- cross-cluster communication
- failover routing

No fake service mesh setup.

---

## Global Traffic Routing Requirements

Support:
- geo-routing
- latency-aware routing
- regional failover
- edge routing

Generate:
- Gateway API configs
- ingress orchestration
- DNS failover strategies
- traffic balancing policies

The routing system MUST tolerate:
- regional outages
- traffic spikes
- edge failures

---

## Database Replication Requirements

Support:
- cross-region replication
- replay-safe synchronization
- failover orchestration
- consistency management

Generate:
- replication topology
- failover workflows
- replay-safe recovery
- split-brain mitigation

Use:
- Vitess or CockroachDB where appropriate

Never:
- ignore replication lag
- ignore recovery orchestration
- assume perfect consistency

---

## Kafka Replication Requirements

Support:
- multi-region Kafka replication
- replay-safe event synchronization
- disaster recovery
- failover consumers

Generate:
- MirrorMaker2 or equivalent
- replication topology
- replay-safe synchronization
- failover orchestration

The event system MUST tolerate:
- regional outages
- replay storms
- duplicate replication events

---

## Redis Requirements

Support:
- global cache synchronization
- regional cache isolation
- replay-safe cache invalidation
- failover coordination

Generate:
- replication topology
- cache invalidation strategy
- replay-safe synchronization
- regional isolation workflows

---

## Observability Requirements

Generate:
- global Prometheus federation
- Thanos aggregation
- Loki centralized logs
- Tempo distributed tracing
- cross-region correlation IDs

Metrics:
- regional latency
- failover duration
- replication lag
- cross-region throughput
- recovery duration

Logs:
- structured JSON logs
- trace IDs
- correlation IDs

Never:
- isolate observability per-region only
- lose replay traceability

---

## Disaster Recovery Requirements

Support:
- regional failover
- replay-safe recovery
- backup restoration
- traffic rerouting

Generate:
- DR runbooks
- automated failover workflows
- replay-safe restoration
- recovery validation tests

The DR system MUST tolerate:
- full regional outage
- partial regional outage
- replay storms
- split-brain risks

No fake DR architecture.

---

## GitOps Requirements

Use:
- ArgoCD

Generate:
- app-of-apps structure
- environment overlays
- progressive delivery
- rollback orchestration

Support:
- multi-region deployments
- staged rollouts
- canary deployments
- drift detection

---

## Security Requirements

The platform MUST:
- enforce mTLS
- isolate regions
- enforce RBAC
- rotate secrets
- audit infrastructure changes

Never:
- expose cluster internals publicly
- trust cross-region traffic blindly
- hardcode infrastructure secrets
- disable mesh security

Generate:
- NetworkPolicies
- RBAC policies
- secret rotation workflows
- infrastructure audit pipelines

---

## Reliability Requirements

Implement:
- retries
- timeout handling
- graceful failover
- panic recovery
- circuit breakers
- backoff strategies

Support:
- rolling deployments
- autoscaling
- distributed deployments

Generate:
- resilience policies
- retry orchestration
- regional isolation
- recovery workflows

---

## CI/CD Requirements

Generate:
- GitHub Actions or Drone pipelines
- Terraform validation
- Helm validation
- Kubernetes policy validation
- vulnerability scanning
- GitOps deployment workflows

---

## Testing Requirements

Generate:
- infrastructure tests
- failover tests
- replay tests
- chaos tests
- latency tests
- disaster recovery tests

Test:
- regional outages
- replay storms
- replication lag
- split-brain risks
- traffic spikes
- distributed deployments
- failover correctness

---

## Output Requirements

Explain:
- global infrastructure architecture
- failover strategy
- geo-routing strategy
- replay-safe synchronization strategy
- DR strategy
- scaling strategy
- resilience strategy

Generate production-grade infrastructure only.

No toy Kubernetes setup.
No fake DR architecture.
No naive replication logic.

---

## Acceptance Criteria

The Global Infrastructure Platform must support future integration with:
- Search Platform
- Recommendation Platform
- Analytics Platform
- Live Commerce Platform
- Fraud Platform

without major future refactors.

The platform MUST realistically tolerate:
- replay storms
- regional outages
- split-brain risks
- distributed deployments
- global traffic spikes

---

## Constraints

Follow ALL:
- .ai/system/*
- .ai/architecture/*
- .ai/planning/*
- .ai/context/*
- .ai/prompts/*

Production-grade only.