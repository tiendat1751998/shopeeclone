# Architecture Compliance Scorecard

## Architecture Standards Checklist
- **Database Partition**: Does the system follow Database-per-Service? (Yes/No)
- **Communications**: Are cross-service synchronizations strictly over gRPC? (Yes/No)
- **Data Delivery**: Does the state change update via transactional Outbox? (Yes/No)

## Scale Audit Rubric
- Target QPS capability:
- CPU bottlenecks identified:
- Redis memory sizing:
