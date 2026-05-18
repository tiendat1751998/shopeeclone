# Event-Driven Message Payloads & Schemas

We strictly enforce JSON Schema standards for messages flowing through Apache Kafka topics.

## 1. Topic: `order.created` Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "OrderCreatedEvent",
  "type": "object",
  "required": ["event_id", "order_id", "buyer_id", "seller_id", "total_amount", "items"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "order_id": { "type": "string", "format": "uuid" },
    "buyer_id": { "type": "string", "format": "uuid" },
    "seller_id": { "type": "string", "format": "uuid" },
    "total_amount": { "type": "number", "minimum": 0 },
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["sku_id", "quantity", "price"],
        "properties": {
          "sku_id": { "type": "string" },
          "quantity": { "type": "integer", "minimum": 1 },
          "price": { "type": "number" }
        }
      }
    },
    "timestamp": { "type": "string", "format": "date-time" }
  }
}
```

## 2. Topic: `payment.completed` Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "PaymentCompletedEvent",
  "type": "object",
  "required": ["event_id", "order_id", "payment_id", "amount", "transaction_ref"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "order_id": { "type": "string", "format": "uuid" },
    "payment_id": { "type": "string", "format": "uuid" },
    "amount": { "type": "number" },
    "transaction_ref": { "type": "string" },
    "timestamp": { "type": "string", "format": "date-time" }
  }
}
```
