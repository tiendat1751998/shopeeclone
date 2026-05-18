# Observability Standards & Telemetry Conventions

Ensuring a fully monitorable, auditable, and alert-ready system.

## 1. Golden Signals Monitoring (RED & USE Methods)
Every service must expose Prometheus metrics covering the **RED Method**:
- **Rate**: Request throughput (requests per second).
- **Errors**: Failures rate (error HTTP 5xx / gRPC status errors).
- **Duration**: Latency distributions (p50, p90, p95, p99 percentiles).

Every resource cluster (Postgres, Redis, Kafka) must monitor the **USE Method**:
- **Utilization**: Percentage of resource actively in use (e.g. CPU, Disk space).
- **Saturation**: Volume of queued work waiting for execution (e.g. Thread pool queues, DB connection wait queue).
- **Errors**: Error count metrics.

## 2. Distributed Tracing Propagation
- **Trace Propagation Standards**: All HTTP headers must propagate tracing context using the W3C Trace Context standard (`traceparent`, `tracestate`).
- **gRPC metadata**: propagate using standard OpenTelemetry metadata injection.
- **Correlation Logs**: Every JSON log entry must contain the `trace_id` and `span_id` fields.
