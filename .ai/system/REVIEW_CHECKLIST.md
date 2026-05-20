# Service Review Checklist

Per review-gates.md, every service MUST have all 10 gates before merging.

## Review Gates

- [ ] **Tracing** — OpenTelemetry spans propagated via W3C Trace Context
- [ ] **Metrics** — Prometheus RED metrics exposed at /metrics
- [ ] **Structured Logs** — JSON logs with trace_id, span_id, service name
- [ ] **Retries** — Transient failures retried with exponential backoff
- [ ] **Graceful Shutdown** — SIGTERM/SIGINT handled, in-flight requests drained
- [ ] **Tests** — Unit tests (80% coverage) + integration tests (Testcontainers)
- [ ] **Kubernetes Manifests** — Deployment, Service, HPA, PDB, NetworkPolicy
- [ ] **Helm Charts** — Templated manifests with values.yaml
- [ ] **Security Review** — No hardcoded secrets, parameterized queries, auth middleware
- [ ] **Scalability Review** — Connection pooling, cache strategy, rate limiting

## Per-Service Status

| Service | Tracing | Metrics | Logs | Retries | Shutdown | Tests | K8s | Helm | Security | Scale |
|---------|---------|---------|------|---------|----------|-------|-----|------|----------|-------|
| auth | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| cart | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ |
| catalog-product | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | ✅ |
| checkout | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | ✅ |
| gateway | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ |
| inventory | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| order | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| payment | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| product | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ⚠️ | ✅ |
| product-catalog | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| promotion | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ⚠️ | ✅ |
| shipment | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ⚠️ | ✅ |

Legend: ✅ = Complete | ⚠️ = Partial | ❌ = Missing
