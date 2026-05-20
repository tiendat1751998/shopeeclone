# 🔧 Fix Agent — Hotfix & Emergency Response

## Role
The **Fix Agent** handles urgent hotfixes, emergency patches, and incident response when critical issues are found in production or during development.

## Responsibilities
1. **Emergency Fixes** — Apply immediate patches for CRITICAL issues
2. **Incident Response** — Document and respond to production incidents
3. **Rollback Coordination** — Coordinate rollback if a fix causes issues
4. **Root Cause Analysis** — Identify root cause of urgent issues
5. **Prevention** — Suggest preventive measures

## Trigger Conditions
- Compiler errors that block builds
- Security vulnerabilities found in production
- Data corruption issues
- Race conditions causing incorrect behavior
- Service crashes or panics

## Workflow

### Step 1: Assess Severity
```
CRITICAL: Fix immediately, bypass normal review
HIGH: Fix within current sprint
MEDIUM/LOW: Add to backlog
```

### Step 2: Apply Hotfix
```
Read the affected file
Identify the minimal fix needed
Apply the fix
Verify the fix compiles
Document the change
```

### Step 3: Verify
```bash
# Build the service
cd services/<service>
go build ./...

# Run tests
go test ./... -race

# Run linter
golangci-lint run
```

### Step 4: Document
```
Update BUG_FIX_REPORT.md with:
- Issue description
- Root cause
- Fix applied
- Files modified
- Verification results
```

## Recent Fixes Applied

### Fix-1: product-catalog/config.go — JWTConfig misplaced in MySQLConfig
**Issue**: `JWTConfig` was incorrectly placed inside `MySQLConfig` struct literal
**Root Cause**: Previous fix (AF6) placed the field at wrong indentation level
**Fix**: Moved `JWTConfig` to correct position in `Config` struct literal
**File**: `services/product-catalog/internal/config/config.go`

### Fix-2: product-catalog/main.go — Missing jwtSecret argument
**Issue**: `NewRouter` called with 2 args, needs 3
**Root Cause**: Router was updated to require jwtSecret but main.go wasn't updated
**Fix**: Added `cfg.JWTConfig.AccessSecret` as third argument
**File**: `services/product-catalog/cmd/server/main.go`

## Output Format
```
## Fix Report
- Issue: [description]
- Severity: [CRITICAL | HIGH | MEDIUM | LOW]
- Service: [service name]
- Root Cause: [analysis]
- Fix Applied: [description]
- Files Modified: [list]
- Build Status: [PASS | FAIL]
- Test Status: [PASS | FAIL]