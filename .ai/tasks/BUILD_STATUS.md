# TASK-006 to TASK-009 Build Status

## COMPLETED: TASK-006 — Inventory Service

Location: /home/datdt/shopeeclone/services/inventory/

### Files Created (35+ files):

**Core:**
- go.mod
- Dockerfile
- cmd/server/main.go

**Config:**
- internal/config/config.go

**Domain:**
- internal/domain/stock.go (Stock, Reservation, Warehouse, StockMovement, FlashSaleInventory)
- internal/domain/events.go (Event types and payloads)
- internal/domain/repository.go (Repository interfaces + StockFilter)

**Infrastructure:**
- internal/infrastructure/mysql/db.go
- internal/infrastructure/mysql/stock_repo.go
- internal/infrastructure/mysql/reservation_repo.go
- internal/infrastructure/mysql/other_repos.go (Warehouse, StockMovement, FlashSale repos)
- internal/infrastructure/redis/store.go (distributed locking, Lua scripts, flash sale atomic ops, idempotency)

**Application:**
- internal/application/service.go (ReserveStock, ConfirmReservation, ReleaseReservation, AdjustStock, GetStock, GetStockAvailability, FlashSaleReserve, CleanupExpiredReservations)

**Transport:**
- internal/transport/http/handler.go
- internal/transport/http/router.go
- internal/transport/kafka/producer.go

**Observability:**
- internal/metrics/metrics.go
- internal/health/health.go
- internal/tracing/tracing.go
- internal/logging/logging.go
- internal/validation/validator.go

**Migrations:**
- migrations/001_initial.sql (warehouses, stocks, reservations, stock_movements, flash_sale_inventory, outbox_events)

**Kubernetes:**
- deployments/deployment.yaml
- deployments/service.yaml
- deployments/configmap.yaml
- deployments/secrets.yaml
- deployments/hpa.yaml
- deployments/pdb.yaml
- deployments/service-monitor.yaml
- deployments/network-policy.yaml

**Helm:**
- charts/Chart.yaml
- charts/values.yaml
- charts/templates/_helpers.tpl

**Tests:**
- tests/unit/domain_test.go

### Key Features Implemented:
- Anti-oversell protection with distributed locking (Redis + optimistic locking in MySQL)
- Reservation lifecycle (pending -> confirmed/released/expired)
- Flash sale stock handling with Redis Lua atomic operations
- Idempotency support
- Event-driven architecture (Kafka)
- Cache-aside pattern with Redis
- Expired reservation cleanup worker
- Full observability (metrics, tracing, structured logging)
- Graceful shutdown
- Kubernetes-native deployment

---

## IN PROGRESS: TASK-007 — Cart Service

Location: /home/datdt/shopeeclone/services/cart/

### Files Created:
- go.mod

### Remaining Files Needed:
- cmd/server/main.go
- internal/config/config.go
- internal/domain/ (cart.go, cart_item.go, events.go, repository.go)
- internal/infrastructure/mysql/ (db.go, cart_repo.go, snapshot_repo.go)
- internal/infrastructure/redis/ (store.go - cart cache, session carts)
- internal/application/ (service.go - add/remove/update/clear/merge/checkout-preview)
- internal/transport/http/ (handler.go, router.go)
- internal/transport/kafka/ (producer.go)
- internal/metrics/metrics.go
- internal/health/health.go
- internal/tracing/tracing.go
- internal/logging/logging.go
- internal/validation/validator.go
- migrations/001_initial.sql
- deployments/ (8 K8s manifest files)
- charts/ (Helm chart files)
- tests/unit/ (domain_test.go, service_test.go)

---

## PENDING: TASK-008 — Promotion Service

Location: /home/datdt/shopeeclone/services/promotion/

### Structure: Same as above with promotion-specific domain:
- Voucher engine (create, redeem, usage limits, expiration)
- Campaign engine (flash-sale, scheduled, seasonal)
- Pricing rules (percentage, fixed, shipping discounts)
- Eligibility engine (user, region, payment, seller, product targeting)
- Stacking rules (mutually exclusive, priority, conflict resolution)
- Abuse prevention (duplicate redemption, rate limiting, fraud detection)

---

## PENDING: TASK-009 — Checkout Service

Location: /home/datdt/shopeeclone/services/checkout/

### Structure: Same as above with checkout-specific domain:
- Checkout orchestration (validation, seller grouping, pricing finalization)
- Pricing freeze (immutable snapshots, anti-price-drift)
- Inventory reservation orchestration (request, rollback, timeout, reconciliation)
- Checkout snapshot (immutable, replay-safe)
- Idempotency protection (anti-double-submit, request fingerprinting)
- Validation engine (cart, promotion, inventory, seller, shipping)
- Saga-like coordination pattern

---

## Architecture Pattern (all services follow):

```
services/<name>/
├── cmd/server/main.go          # Entry point, wire everything
├── internal/
│   ├── config/config.go        # Env-based configuration
│   ├── domain/                 # Business entities, events, repository interfaces
│   ├── application/service.go  # Business logic, use cases
│   ├── infrastructure/
│   │   ├── mysql/              # DB connection + repository implementations
│   │   └── redis/              # Redis store (cache, locks, atomic ops)
│   ├── transport/
│   │   ├── http/               # Gin handler + router
│   │   └── kafka/              # Event producer
│   ├── metrics/metrics.go      # Prometheus metrics
│   ├── health/health.go        # Health check handlers
│   ├── tracing/tracing.go      # OpenTelemetry init
│   ├── logging/logging.go      # Structured logging helpers
│   └── validation/validator.go # Input validation
├── migrations/                  # SQL migrations
├── deployments/                 # K8s manifests
├── charts/                      # Helm chart
├── tests/unit/                  # Unit tests
└── Dockerfile
```
