# DevOps & Infrastructure Rules

## 1. Multi-Stage Production Containerization
All services must build using secure, lightweight Docker images.

### Go Production Dockerfile
```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# Stage 2: Distroless Run
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/main /main
USER nonroot:nonroot
ENTRYPOINT ["/main"]
```

### Java production Dockerfile
```dockerfile
# Stage 1: Build
FROM maven:3.9-eclipse-temurin-17 AS builder
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Minimal JRE Run
FROM eclipse-temurin:17-jre-alpine
RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring
COPY --from=builder /app/target/*.jar app.jar
ENTRYPOINT ["java", "-XX:+UseG1GC", "-jar", "/app.jar"]
```

## 2. Kubernetes Deployment Configuration Guidelines
- Every Deployment manifest must contain explicit resource configurations:
  ```yaml
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "500m"
  ```
- Define `livenessProbe` and `readinessProbe` targeting Actuator (`/actuator/health`) or standard check paths.
