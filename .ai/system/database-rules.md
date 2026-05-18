# Database Rules & Performance Specifications

To maintain high scale and fast transactional commits, all microservice databases must adhere to these rules.

## 1. Indexing & Query Optimizations
- **Composite Indexes**: Always match the indexing order with the query filtering order.
  ```sql
  -- If querying: WHERE shop_id = ? AND status = ? ORDER BY created_at DESC
  CREATE INDEX idx_shop_status_date ON orders(shop_id, status, created_at DESC);
  ```
- **Unused Index Prevention**: Never write queries targeting columns that are parts of composite indexes without matching the leftmost prefix.
- **Explain Analyze**: Every new SQL query must run through `EXPLAIN ANALYZE` to ensure index scans (no sequential table scans allowed on tables with > 10,000 rows).

## 2. Transactions & Limits
- **Max Transaction Duration**: No transaction may remain open for more than **2 seconds** to prevent pool exhaustion.
- **No Long-Running Processing inside Transactions**: Third-party API calls, gRPC requests, and password hashing must be executed outside the active database transaction scope.
- **Soft Deletes Indexing**: For tables utilizing soft deletes (`deleted_at` timestamp), always append `WHERE deleted_at IS NULL` to index definitions.
