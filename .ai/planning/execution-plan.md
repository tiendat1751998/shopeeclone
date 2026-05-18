# Local Infrastructure Execution Plan

Initialize all local resources using Docker Compose to mimic the cloud environment.

## Docker Compose Setup (`docker-compose.yml` guidelines)
Create `docker-compose.yml` in project root with these exact resources:
- **PostgreSQL**: Port `5432` (Auth, Order, Payment DBs).
- **MongoDB**: Port `27017` (Product Catalog DB).
- **Redis (Cluster mode)**: Port `6379` (Session, Stock reservation, Cart cache).
- **Apache Kafka + Zookeeper**: Ports `9092` (Message Bus).
- **Elasticsearch + Kibana**: Ports `9200` (Search Index Engine).

### Bootstrapping Command
```powershell
docker-compose up -d
```
