# Database DDL Schemas & Cache Structures

## 1. PostgreSQL Schema - Order & Checkout Service
```sql
-- DDL representing transactional orders
CREATE TABLE orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    payment_method VARCHAR(20) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    discount_amount DECIMAL(15, 2) DEFAULT 0.00,
    shipping_fee DECIMAL(15, 2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'CREATED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
    sku_id VARCHAR(50) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX idx_orders_buyer ON orders(buyer_id);
CREATE INDEX idx_orders_status ON orders(status);
```

## 2. MongoDB Schema - Product Catalog Service
```json
// Representing hierarchical and dynamic e-commerce catalog structure
{
  "$jsonSchema": {
    "bsonType": "object",
    "required": ["spu_id", "title", "category_id", "skus"],
    "properties": {
      "spu_id": { "bsonType": "string" },
      "title": { "bsonType": "string" },
      "description": { "bsonType": "string" },
      "category_id": { "bsonType": "string" },
      "skus": {
        "bsonType": "array",
        "items": {
          "bsonType": "object",
          "required": ["sku_id", "price", "stock", "variations"],
          "properties": {
            "sku_id": { "bsonType": "string" },
            "price": { "bsonType": "double" },
            "stock": { "bsonType": "int" },
            "variations": {
              "bsonType": "array",
              "items": {
                "bsonType": "object",
                "required": ["name", "value"],
                "properties": {
                  "name": { "bsonType": "string" },
                  "value": { "bsonType": "string" }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

## 3. Redis Cache Keys Directory
- **Active Cart Keys**: `cart:{userId}` -> Hash mapping `skuId` to `quantity` (TTL: 30 days).
- **Inventory Reservation Key**: `stock:{skuId}` -> String mapping quantity remaining (TTL: dynamic/campaign length).
- **Session Tokens cache**: `session:{userId}` -> String containing encrypted payload (TTL: 15 mins).
