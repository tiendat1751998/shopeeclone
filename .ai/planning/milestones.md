# Release Milestones & Target KPIs

## Milestone 1: Baseline Checkout MVP (Exit Criteria)
- JWT user registration, access tokens refresh.
- Product catalog allows search filters and pagination.
- Order is submitted, payment webhook processes VNPay simulator, inventory drops.
- Load performance: Handles stable 200 QPS with latencies < 200ms.

## Milestone 2: Flash Sale Ready (Exit Criteria)
- Lua stock checks fully loaded into Redis.
- Concurrency simulation showing 20k checkout requests over 5 seconds results in exact 0 stock discrepancies.
- Zero database locks or pool crashes during peak simulation.

## Milestone 3: Scale & Launch (Exit Criteria)
- Kubernetes auto-scalers configure memory boundaries.
- Continuous Elasticsearch synchronizations complete via Kafka CDC logs under 1 second lag.
