# Architectural Blueprint & Distributed Standards

## 1. Microservices Separation & Boundaries
- **Database Isolation**: No service may access another service's database. Cross-boundary data requests must flow through gRPC (synchronous query) or Kafka (asynchronous state change).
- **Communication Protocol**:
  - Internal Service-to-Service: **gRPC** (Proto3).
  - External-to-Internal: **JSON over REST** or **GraphQL** resolved at the Kong API Gateway.

## 2. The Transactional Outbox Pattern
To prevent dual-write failures (e.g. database commit succeeds but message broker publish fails):
```sql
-- DDL for Outbox Table in each PostgreSQL Service
CREATE TABLE outbox_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_outbox_processed ON outbox_events(processed) WHERE processed = FALSE;
```
- A background worker (Debezium or a polling daemon) reads `processed = FALSE` events, publishes them to Kafka with partition matching, and marks them `processed = TRUE`.

## 3. Distributed Transactions via SAGA
- Prefer **Choreography SAGA** for low-complexity, decoupled pipelines.
- Prefer **Orchestrator-based SAGA** for the checkout-to-payment cycle:
  ```mermaid
  sequenceDiagram
      autonumber
      participant O as Checkout Orchestrator
      participant I as Inventory Service
      participant P as Payment Service
      participant D as Order Service

      O->>I: ReserveStock(sku_id, qty)
      Note over I: Lock stock in Redis/DB
      I-->>O: StockReserved (Success)
      O->>P: ProcessPayment(order_id, amount)
      alt Payment Fails
          P-->>O: PaymentFailed
          O->>I: ReleaseStockCompensate(sku_id, qty)
          O->>D: UpdateOrderStatus(order_id, "FAILED")
      else Payment Succeeds
          P-->>O: PaymentSucceeded
          O->>D: UpdateOrderStatus(order_id, "PAID")
      end
  ```
