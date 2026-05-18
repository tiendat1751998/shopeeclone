# REST API Standards & Error Formats

All client-facing REST APIs must follow standard RESTful conventions.

## 1. URI Naming & HTTP Methods
- **Plural Nouns**: All resources must be plural nouns (e.g., `/api/v1/products`, `/api/v1/orders`).
- **HTTP Verbs matching**:
  - `GET`: Read resources (idempotent).
  - `POST`: Create new resources.
  - `PUT`: Complete update of an existing resource.
  - `PATCH`: Partial update of an existing resource.
  - `DELETE`: Remove a resource.

## 2. Standard JSON Error Payload
All services returning 4xx or 5xx errors must output this exact payload layout:
```json
{
  "error_code": "ORDER_STOCK_INSUFFICIENT",
  "message": "Chi tiết sản phẩm đã hết hàng trong kho.",
  "timestamp": "2026-05-18T10:16:00Z",
  "details": [
    {
      "field": "sku_id",
      "issue": "SKU sku_12345 only has 0 items remaining."
    }
  ],
  "trace_id": "8ce09b64a0c845668c46d963a436e228"
}
```
