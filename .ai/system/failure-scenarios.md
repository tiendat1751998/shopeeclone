# FAILURE SCENARIOS

Every service MUST consider:

- Redis unavailable
- Kafka lag
- DB failover
- pod crashes
- node crashes
- network latency
- partial outages
- retry storms
- thundering herd
- duplicate events
- replay attacks

Architecture MUST:
- degrade gracefully
- isolate failures
- support retries
- support backoff
- support idempotency