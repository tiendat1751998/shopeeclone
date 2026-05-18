# Code Review & PR Checklist Template

Every Pull Request must be reviewed against this strict rubric before merging to `main`.

## Rubric Checklists

### 1. Security & Compliance
- [ ] No plaintext secrets or sensitive tokens committed.
- [ ] All inputs are strictly validated (size, format, XSS escaping).
- [ ] No raw SQL queries containing string concatenation (SQLi check).
- [ ] Cryptographic methods use modern standards (BCrypt, AES-256-GCM).

### 2. Scalability & Performance
- [ ] JPA relationships configured with `LAZY` load; N+1 issues resolved.
- [ ] Redis caching applied to high-read static data (e.g. category tree).
- [ ] Database index matching added for any newly added query columns.
- [ ] Goroutines have clear context limits and recovery functions.

### 3. Reliability & Testing
- [ ] Unit tests written; code coverage is at least 80%.
- [ ] Database migration scripts (`Vxx__schema.sql`) are backward-compatible.
- [ ] Outbox pattern implemented for state modifications needing Kafka events.
