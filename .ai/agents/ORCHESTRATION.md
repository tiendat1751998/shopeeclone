# 🔄 Orchestration Workflow — How Agents Collaborate

## Overview
This document defines how the Architect Agent coordinates with QA, DEV, Test, Review, and Fix subagents to fix bugs and improve the Shopee Clone monorepo.

## Agent Hierarchy
```
                    ┌─────────────────┐
                    │   ARCHITECT      │
                    │  (Coordinator)   │
                    └────────┬────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
     ┌──────┴──────┐  ┌─────┴──────┐  ┌──────┴──────┐
     │  QA Agent   │  │ DEV Agent  │  │ Review Agent│
     │ (Validate)  │  │ (Implement)│  │ (Approve)   │
     └──────┬──────┘  └─────┬──────┘  └──────┬──────┘
            │                │                │
            │         ┌──────┴──────┐         │
            │         │ Test Agent  │         │
            │         │ (Coverage)  │         │
            │         └──────┬──────┘         │
            │                │                │
            └────────────────┼────────────────┘
                             │
                    ┌────────┴────────┐
                    │   Fix Agent     │
                    │ (Hotfix/Emerg.) │
                    └─────────────────┘
```

## Sprint Workflow

### Phase 1: TRIAGE (Architect)
**Input**: `QA_BUG_REPORT.md`, `SECURITY_AUDIT_REPORT.md`
**Output**: Prioritized sprint backlog

```
1. Read all bug reports
2. Group issues by service
3. Prioritize: CRITICAL > HIGH > MEDIUM > LOW
4. Within each severity: Security > Data Loss > Race Condition > Code Quality
5. Create sprint backlog with assignments
6. Update .ai/planning/sprint-plan.md
```

### Phase 2: DEVELOPMENT (DEV Agent)
**Input**: Sprint backlog, bug report entries
**Output**: Code fixes

```
1. Read the bug report entry
2. Read the affected source file
3. Follow .ai/system/coding-rules.md
4. Follow .ai/system/security-rules.md
5. Implement the fix
6. Self-review against .ai/system/review-rules.md
7. Update BUG_FIX_REPORT.md
```

### Phase 3: TESTING (Test Agent)
**Input**: Code fixes from DEV
**Output**: Test coverage report

```
1. Read the fix implementation
2. Write table-driven unit tests
3. Ensure 80% coverage for business logic
4. Run: go test ./... -race -cover
5. Add integration tests if needed
6. Report coverage percentage
```

### Phase 4: REVIEW (Review Agent)
**Input**: Code fixes + tests
**Output**: Approval or change requests

```
1. Read the fix diff
2. Check security (no SQLi, XSS, auth bypass)
3. Check performance (no N+1, missing indexes)
4. Check reliability (error handling, graceful shutdown)
5. Check code quality (clean code, proper logging)
6. Decision: APPROVE | REQUEST_CHANGES | REJECT
```

### Phase 5: QA VERIFICATION (QA Agent)
**Input**: Approved fixes
**Output**: QA verification report

```
1. Read the approved fix
2. Verify the fix addresses the root cause
3. Check for regressions
4. Update QA_BUG_REPORT.md (move to Fixed)
5. Report verification status
```

### Phase 6: HOTFIX (Fix Agent) — As Needed
**Input**: Urgent issues, compiler errors, production incidents
**Output**: Emergency patches

```
1. Assess severity
2. Apply minimal fix
3. Verify build: go build ./...
4. Document the fix
5. Update BUG_FIX_REPORT.md
```

## Communication Protocol

### Task Assignment Format
```
[TASK ASSIGNMENT]
- Issue ID: H1
- Service: order
- Assigned To: DEV Agent
- Priority: P1
- Deadline: Current Sprint
- Description: Fix unsafe type assertions in handler.go
- File: services/order/internal/transport/http/handler.go
- Verification: go build ./... && go test ./...
```

### Status Update Format
```
[STATUS UPDATE]
- Issue ID: H1
- Phase: DEV
- Status: IN_PROGRESS | COMPLETE | BLOCKED
- Notes: [details]
- Next: [next phase]
```

## Rules All Agents Must Follow

### From .ai/prompts/masterpromt.md
- NEVER use placeholders or fake implementations
- ALWAYS implement production-grade code
- ALWAYS include error handling, logging, metrics
- NEVER hardcode secrets
- ALWAYS validate inputs

### From .ai/system/coding-rules.md
- Go: Use errgroup for goroutines, structured error wrapping
- Java: LAZY fetch, @Transactional boundaries
- TypeScript: Strict types, Server Components for SEO

### From .ai/system/security-rules.md
- Parameterized SQL queries only
- JWT with short lifetimes (15m access, 7d refresh)
- Rate limiting on critical endpoints

### From .ai/system/testing-rules.md
- Minimum 80% coverage for business logic
- Table-driven tests in Go
- Testcontainers for integration tests

### From .ai/system/review-rules.md
- Security checklist before merge
- Performance checklist before merge
- Reliability checklist before merge