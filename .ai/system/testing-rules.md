# Testing Rules & Stress Testing Standards

Our codebase requires rigorous quality gates before code reaches staging or production.

## 1. Unit Testing Guidelines
- **Coverage Mandate**: Minimum **80% code coverage** required for all business logic (Service/UseCase layers).
- **Test Isolation**: Database calls must be completely mocked in unit tests. Go must use `go-sqlmock` or interfaces, Java must use `@Mock` and Mockito.
- **Table-Driven Tests (Go)**: Use table-driven testing pattern for rich variation inputs.

## 2. Integration & Stress Tests
- **Testcontainers**: Integration tests requiring actual databases (Postgres, MongoDB, Redis) must boot clean, isolated containers dynamically using the **Testcontainers** library.
- **K6 Stress Testing**: Flash sale endpoints must be tested against simulated peak load profiles using K6 load simulation scripts.
  ```js
  // Example K6 stress script target
  export const options = {
    stages: [
      { duration: '1m', target: 5000 }, // ramp up to 5k parallel virtual users
      { duration: '3m', target: 5000 }, // stay at 5k users
      { duration: '1m', target: 0 },    // ramp down
    ],
  };
  ```
