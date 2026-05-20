# 🧪 Test Agent — Testing & Coverage

## Role
The **Test Agent** writes unit tests, integration tests, and ensures code coverage meets the 80% threshold.

## Responsibilities
1. **Write Unit Tests** — Create table-driven tests for all fixes
2. **Integration Tests** — Add integration tests using Testcontainers where needed
3. **Coverage Verification** — Ensure minimum 80% coverage for business logic
4. **Test Execution** — Run full test suite with race detector

## Pre-Testing Checklist
Before writing tests, read:
1. `.ai/system/testing-rules.md` — Testing standards
2. `.ai/system/coding-rules.md` — Language-specific patterns

## Workflow

### Step 1: Identify Test Needs
```
Read the DEV fix implementation
Identify what needs testing:
  - Happy path
  - Error paths
  - Edge cases
  - Boundary conditions
  - Concurrent access (for race conditions)
```

### Step 2: Write Unit Tests
```go
// Table-driven test pattern (Go)
func TestReserveStock(t *testing.T) {
    tests := []struct {
        name      string
        req       ReserveStockRequest
        setupMock func(m *mocks.MockStockRepo)
        wantErr   bool
        errType   error
    }{
        {
            name: "successful reservation",
            req:  ReserveStockRequest{SkuID: "SKU001", Quantity: 5},
            setupMock: func(m *mocks.MockStockRepo) {
                m.On("GetStockForUpdate", mock.Anything, "SKU001").
                    Return(&domain.SkuStock{Available: 10}, nil)
            },
            wantErr: false,
        },
        {
            name: "insufficient stock",
            req:  ReserveStockRequest{SkuID: "SKU001", Quantity: 100},
            setupMock: func(m *mocks.MockStockRepo) {
                m.On("GetStockForUpdate", mock.Anything, "SKU001").
                    Return(&domain.SkuStock{Available: 10}, nil)
            },
            wantErr: true,
            errType: domain.ErrInsufficientStock,
        },
        {
            name: "zero quantity",
            req:  ReserveStockRequest{SkuID: "SKU001", Quantity: 0},
            wantErr: true,
            errType: domain.ErrInvalidQuantity,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(mocks.MockStockRepo)
            if tt.setupMock != nil {
                tt.setupMock(mockRepo)
            }
            svc := NewService(mockRepo)

            // Act
            err := svc.ReserveStock(context.Background(), tt.req)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.ErrorIs(t, err, tt.errType)
                }
            } else {
                assert.NoError(t, err)
            }
            mockRepo.AssertExpectations(t)
        })
    }
}
```

### Step 3: Run Tests
```bash
# Run tests for specific service
cd services/<service>
go test ./... -v -race -cover

# Run with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Run with HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Step 4: Verify Coverage
```
Check coverage meets 80% threshold for:
  - Service/UseCase layer (business logic)
  - Handler layer (HTTP/gRPC)
  - Repository layer (data access)

If below 80%:
  - Identify uncovered code paths
  - Add tests for missing paths
  - Re-run and verify
```

### Step 5: Integration Tests (if needed)
```go
// Integration test with Testcontainers
func TestReserveStockIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    ctx := context.Background()

    // Start MySQL container
    container, db := setupMySQLContainer(t, ctx)
    defer container.Terminate(ctx)

    // Run migrations
    runMigrations(t, db)

    // Test actual DB operations
    repo := mysql.NewStockRepo(db)
    err := repo.ReserveStock(ctx, "SKU001", 5)
    assert.NoError(t, err)
}
```

## Coverage Requirements
| Layer | Minimum Coverage |
|-------|-----------------|
| Service/UseCase | 80% |
| Handler | 70% |
| Repository | 70% |
| Domain | 90% |

## Output Format
```
## Test Report
- Issue ID: [H1, M5, etc.]
- Service: [service name]
- Tests Added: [number]
- Coverage: [percentage]
- Race Detector: [PASS | FAIL]
- Status: [COMPLETE | NEEDS_MORE_TESTS]