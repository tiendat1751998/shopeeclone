# CHECKOUT SERVICE - BUG REPORT
**Date:** May 21, 2026  
**Service:** services/checkout  
**Scope:** Domain, Application, Infrastructure, Transport layers

---

## Summary

| Severity | Count |
|----------|-------|
| **CRITICAL** | 0 |
| **HIGH** | 5 |
| **MEDIUM** | 8 |
| **LOW** | 4 |
| **TOTAL** | **17** |

---

## 🔴 HIGH Severity Issues (5)

### H1: `RetryCheckout` doesn't validate user ownership properly
**File:** `internal/application/service.go:358-381`  
**Issue:** The `RetryCheckout` method validates `checkout.UserID != requestingUserID` but the handler doesn't extract `user_id` from JWT token properly - it uses `c.Get("user_id")` which may not exist if auth middleware is not configured correctly.

```go
// service.go:366-368
if checkout.UserID != requestingUserID {
    return domain.ErrUnauthorized
}
```

**Impact:** Unauthorized users could retry other users' checkouts if auth middleware is misconfigured.

**Fix:** Ensure consistent user ID extraction from JWT claims.

---

### H2: `handleError` uses fragile string prefix matching
**File:** `internal/transport/http/handler.go:87-96`  
**Issue:** The error matching uses string prefix comparison which can match wrong errors:

```go
if err.Error() == domainErr.Error() || (len(err.Error()) >= len(domainErr.Error()) && err.Error()[:len(domainErr.Error())] == domainErr.Error()) {
```

**Impact:** Wrapped errors (e.g., `fmt.Errorf("wrap: %w", ErrCheckoutNotFound)`) could match incorrectly or not at all.

**Fix:** Use `errors.Is()` or `errors.As()` for proper error unwrapping:

```go
if errors.Is(err, domainErr) {
    c.AbortWithStatusJSON(status, gin.H{"error_code": domainErr.Error(), "message": err.Error()})
    return
}
```

---

### H3: `KafkaConfig` loaded but never used
**File:** `internal/config/config.go:108-110`  
**Issue:** Kafka configuration is loaded but no Kafka producer/consumer is initialized or used in the service.

**Impact:** Misleading configuration - developers may think events are being published to Kafka.

**Fix:** Either implement Kafka event publishing or remove the unused config.

---

### H4: `GRPCPort` configured but no gRPC server
**File:** `internal/config/config.go:15,82`  
**Issue:** gRPC port is configured but no gRPC server is implemented.

**Impact:** Port conflict if another service tries to use the same port.

**Fix:** Implement gRPC server or remove the config.

---

### H5: No authentication on HTTP endpoints (conditionally applied)
**File:** `internal/transport/http/router.go:38-40`  
**Issue:** Authentication is only applied if `jwtSecret != ""`:

```go
if r.jwtSecret != "" {
    api.Use(auth.GinJWTAuth(r.jwtSecret))
}
```

**Impact:** If JWT_ACCESS_SECRET env var is not set, all endpoints are unprotected.

**Fix:** Make authentication mandatory in production:

```go
if r.jwtSecret == "" {
    log.Fatal("JWT_ACCESS_SECRET is required in production")
}
api.Use(auth.GinJWTAuth(r.jwtSecret))
```

---

## 🟡 MEDIUM Severity Issues (8)

### M1: Dead code - unused packages
**Files:** `internal/tracing/`, `internal/logging/`, `internal/health/`, `internal/validation/`  
**Issue:** These packages exist but are not used in the main service.

**Impact:** Confusing codebase, wasted maintenance effort.

**Fix:** Remove unused packages.

---

### M2: `CheckoutLatency` metric defined but never recorded
**File:** `internal/metrics/metrics.go:34-38`  
**Issue:** The `CheckoutLatency` histogram is defined but never used in the service.

**Impact:** Missing observability for checkout step performance.

**Fix:** Record latency in each step:

```go
start := time.Now()
// ... step execution ...
metrics.CheckoutLatency.WithLabelValues(step).Observe(time.Since(start).Seconds())
```

---

### M3: `gin.SetMode` called after `gin.New()`
**File:** Likely in `cmd/server/main.go`  
**Issue:** Gin mode should be set before creating the engine.

**Impact:** Debug mode may be active in production.

**Fix:** Set mode before creating engine:

```go
gin.SetMode(gin.ReleaseMode)
engine := gin.New()
```

---

### M4: `UpdateStatus` parameter named `id` but it's a reservation key
**File:** `internal/infrastructure/mysql/repos.go:92-94`  
**Issue:** The parameter name is misleading:

```go
func (r *ReservationOrchestrationRepository) UpdateStatus(ctx context.Context, key, status string) error {
```

**Impact:** Confusing code, potential misuse.

**Fix:** Rename parameter to `reservationKey`.

---

### M5: `FindExpired` uses string parameter for time
**File:** `internal/infrastructure/mysql/repos.go:45-48`  
**Issue:** Time is passed as string instead of `time.Time`:

```go
func (r *CheckoutRepository) FindExpired(ctx context.Context, before string, limit int) ([]*domain.Checkout, error) {
```

**Impact:** Time format issues, potential SQL injection risk.

**Fix:** Use `time.Time` parameter:

```go
func (r *CheckoutRepository) FindExpired(ctx context.Context, before time.Time, limit int) ([]*domain.Checkout, error) {
```

---

### M6: `SELECT *` in multiple queries
**File:** `internal/infrastructure/mysql/repos.go`  
**Issue:** Multiple queries use `SELECT *` which is fragile against schema changes.

**Impact:** Breaking changes when schema evolves.

**Fix:** Use explicit column names.

---

### M7: `mustJSON` swallows marshal errors (FIXED but verify)
**File:** `internal/application/service.go:416-422`  
**Issue:** Previously `mustJSON` would panic on error. Now it returns `(string, error)` but callers don't always check the error.

**Impact:** Empty data in DB if marshaling fails.

**Fix:** Ensure all callers check the error return value.

---

### M8: Context not passed to goroutine properly
**File:** `internal/application/service.go:86-95`  
**Issue:** The saga goroutine uses `context.Background()` instead of the request context:

```go
sagaCtx, sagaCancel := context.WithTimeout(context.Background(), 5*time.Minute)
go func() {
    defer sagaCancel()
    s.executeSaga(sagaCtx, checkout.ID, req)
}()
```

**Impact:** Tracing context is lost, no parent span.

**Fix:** Pass request context with timeout:

```go
sagaCtx, sagaCancel := context.WithTimeout(ctx, 5*time.Minute)
```

---

## 🟢 LOW Severity Issues (4)

### L1: Missing error check in `stepComplete`
**File:** `internal/application/service.go:265-292`  
**Issue:** `s.reconcileRepo.Create(ctx, job)` error is not checked.

**Impact:** Reconciliation jobs may be silently lost.

**Fix:** Add error logging:

```go
if err := s.reconcileRepo.Create(ctx, job); err != nil {
    logger.Error("failed to create reconciliation job", zap.Error(err))
}
```

---

### L2: Missing error check in `rollbackReservations`
**File:** `internal/application/service.go:328-338`  
**Issue:** `s.reservationRepo.UpdateStatus` errors are not checked.

**Impact:** Reservations may not be properly released.

**Fix:** Add error logging.

---

### L3: Missing error check in `logStep`
**File:** `internal/application/service.go:340-343`  
**Issue:** `s.stepLogRepo.Create` error is not checked.

**Impact:** Step logs may be silently lost.

**Fix:** Add error logging.

---

### L4: Missing error check in `handleFailure`
**File:** `internal/application/service.go:294-326`  
**Issue:** `s.reconcileRepo.Create` error is not checked.

**Impact:** Reconciliation jobs may be silently lost on failure.

**Fix:** Add error logging.

---

## 🔧 Recommended Fixes Priority

1. **HIGH:** Fix error handling in `handleError` - use `errors.Is()`
2. **HIGH:** Make authentication mandatory in production
3. **HIGH:** Fix context propagation to saga goroutine
4. **MEDIUM:** Record `CheckoutLatency` metrics
5. **MEDIUM:** Use `time.Time` instead of string for `FindExpired`
6. **LOW:** Add error logging for all repository calls

---

## 📊 Code Quality Metrics

| Metric | Value |
|--------|-------|
| Total Lines of Code | ~1,200 |
| Unused Packages | 4 |
| Missing Error Checks | 6 |
| Hardcoded Defaults | 1 (MySQL password) |
| Missing Metrics | 1 |

---

## ✅ Previously Fixed Issues (from QA_BUG_REPORT.md)

1. ✅ `mustJSON` changed to return `(string, error)` 
2. ✅ Added error logging in `handleFailure`, `stepComplete`, `rollbackReservations`, `logStep`
3. ✅ Added auth middleware to routes
4. ✅ Fixed hardcoded secrets with `requireEnv()`

---

## 🔍 Additional Notes

### Race Condition Risk
The `executeSaga` runs in a goroutine and updates the checkout state. If the main request returns immediately and the client polls for status, there's a race condition where the status might not be updated yet.

### Idempotency Check Logic
The idempotency check at line 69-76 returns the existing checkout if it's not failed or rolled back. This means a completed checkout will be returned, which may confuse clients expecting a new checkout.

### Missing Distributed Lock
The inventory reservation doesn't use distributed locks, which could lead to over-selling in high-concurrency scenarios.