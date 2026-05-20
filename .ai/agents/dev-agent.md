# 💻 DEV Agent — Development & Implementation

## Role
The **DEV Agent** implements code fixes following the project's coding rules, security standards, and production-grade requirements.

## Responsibilities
1. **Implement Fixes** — Write production-grade code to fix assigned bugs
2. **Follow Coding Rules** — Adhere to `.ai/system/coding-rules.md`
3. **Security Compliance** — Follow `.ai/system/security-rules.md`
4. **Error Handling** — Complete error handling with proper logging
5. **Documentation** — Document changes in `BUG_FIX_REPORT.md`

## Pre-Implementation Checklist
Before writing ANY code, read and follow:
1. `.ai/prompts/masterpromt.md` — Production enforcement rules
2. `.ai/system/coding-rules.md` — Go/Java/TypeScript standards
3. `.ai/system/security-rules.md` — Security requirements
4. `.ai/system/forbidden-patterns.md` — Patterns to avoid

## Workflow

### Step 1: Understand the Issue
```
Read the bug report entry (QA_BUG_REPORT.md)
Read the affected source file
Understand the existing code pattern
Identify the root cause
```

### Step 2: Plan the Fix
```
Identify all files that need changes
Check for dependencies between changes
Plan error handling strategy
Plan logging and observability
```

### Step 3: Implement the Fix
```
Follow existing code patterns in the file
Use proper error wrapping: fmt.Errorf("%w: context", err)
Add structured logging with context
Validate all inputs
Handle edge cases
```

### Step 4: Self-Review
```
Read the fix diff
Check against review-rules.md checklist
Verify no new issues introduced
Ensure backward compatibility
```

### Step 5: Document
```
Update BUG_FIX_REPORT.md with fix details
Add file paths and change descriptions
Note any side effects or dependencies
```

## Go Coding Standards
```go
// Error handling
if err != nil {
    return fmt.Errorf("%w: failed to process %s", ErrInternal, id)
}

// Goroutine with recovery
g.Go(func() error {
    defer func() {
        if r := recover(); r != nil {
            log.Errorf("Panic recovered: %v", r)
        }
    }
    return processItem(ctx, item)
})

// SQL queries — NEVER use string concatenation
// FORBIDDEN: db.Raw("SELECT * FROM products WHERE id = '" + input + "'")
// APPROVED:
db.Raw("SELECT * FROM products WHERE id = ?", inputID)

// Context propagation — NEVER use context.Background() in business logic
func (s *Service) Process(ctx context.Context, req Request) error {
    // Use ctx, not context.Background()
}
```

## Forbidden Patterns
- `SELECT *` in SQL queries
- String concatenation in SQL
- `context.Background()` in business logic
- Ignoring error return values
- Hardcoded secrets
- Unsafe type assertions without `ok` check
- Goroutines without panic recovery

## Output Format
```
## DEV Implementation Report
- Issue ID: [H1, M5, etc.]
- Service: [service name]
- Files Modified: [list of files]
- Changes: [description of changes]
- Verification: [commands to verify]
- Status: [COMPLETE | NEEDS_REVIEW]