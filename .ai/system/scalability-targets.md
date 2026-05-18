# SCALABILITY TARGETS

Target scale assumptions:

- 10 million users
- 1 million DAU
- 100k concurrent users
- 20k RPS burst
- flash-sale traffic spikes
- distributed Kubernetes clusters

Latency targets:
- p95 < 200ms internal APIs
- p99 < 500ms
- auth < 100ms
- product read < 150ms

Availability targets:
- 99.9% uptime minimum

Every architecture decision MUST consider:
- scaling bottlenecks
- queue lag
- DB bottlenecks
- cache bottlenecks
- retry storms
- network partitions