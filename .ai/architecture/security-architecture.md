# Network Isolation & API Gateway Access Policies

To prevent security breaches, we enforce strict network segmentation.

## Network Segment Access Matrix
1. **Public Web Entry**: Acknowledges requests from internet -> Port `443` (Cloudflare WAF).
2. **Gateway Segment**: Gateway handles TLS termination, JWT extraction, and routes requests internally through gRPC mTLS.
3. **Database Segment**: PostgreSQL, MongoDB, and Redis clusters listen strictly on internal VPN virtual subnets (`10.0.0.0/16`), with all public inbound ports disabled.

## Kong Routing YAML Definition (Example)
```yaml
_format_version: "2.1"
services:
- name: order-service
  url: http://order-processing.internal:8080
  routes:
  - name: order-route
    paths:
    - /api/v1/orders
    - /api/v1/checkout
  plugins:
  - name: rate-limiting
    config:
      second: 5
      hour: 3600
      policy: redis
      redis_host: redis.internal
```
