# Disaster Recovery & System Resiliency Specifications

This document outlines the backup, failover, and restoration protocols for critical business continuity.

## 1. RTO & RPO Targets
| System Subsystem | RTO (Recovery Time Objective) | RPO (Recovery Point Objective) | Strategy |
| :--- | :--- | :--- | :--- |
| **User & Checkout** | `< 2 minutes` | `< 5 seconds` | Multi-AZ active-passive failover |
| **Payment Ledger** | `< 5 minutes` | `0 (Strict ACID)` | Write-Ahead Logs Sync + Daily Audit reconciliations |
| **Product Search** | `< 1 hour` | `< 5 minutes` | CDC Re-indexing from MongoDB source of truth |

## 2. Database Failover & Replication
- **PostgreSQL**: Configure Streaming Replication with 1 primary node (read/write) and 2 read replicas. Use **PgBouncer** or **Patroni** for automated health checks and primary node promotion inside Kubernetes.
- **MongoDB**: Setup 3-node Replica Sets. If the primary node goes offline, the set automatically holds elections to promote a new primary in < 12 seconds.
