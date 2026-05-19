.PHONY: all build clean build-java build-go test

# Default build directory
BIN_DIR = $(CURDIR)/bin

# Flags
SKIP_TESTS ?= true

# Go executable extension for Windows
ifeq ($(OS),Windows_NT)
    EXE_EXT = .exe
else
    EXE_EXT =
endif

# List of Go modules to build
GO_MODULES = \
    services/auth \
    services/cart \
    services/catalog-product \
    services/checkout \
    services/gateway \
    services/inventory \
    services/order \
    services/payment \
    services/product \
    services/product-catalog \
    services/promotion \
    services/shipment \
    platforms/advertising \
    platforms/aiml \
    platforms/analytics \
    platforms/api-gateway \
    platforms/billing \
    platforms/developer \
    platforms/fraud \
    platforms/fraud-risk \
    platforms/global-infra \
    platforms/live-commerce \
    platforms/live-scale \
    platforms/logistics-delivery \
    platforms/notification \
    platforms/notification-campaign \
    platforms/oms-fulfillment \
    platforms/payment-ledger \
    platforms/rec-vector \
    platforms/recommendation \
    platforms/search \
    platforms/search-indexing \
    platforms/service-mesh \
    platforms/sre

all: clean build

build: build-java build-go

clean:
	@echo "Cleaning output directory..."
	@rm -rf $(BIN_DIR)
	@mkdir -p $(BIN_DIR)

images:
	@chmod +x build-images.sh 2>/dev/null || true
	@./build-images.sh

build-java:
	@echo "========================================="
	@echo "[Java] Building identity-auth service..."
	@echo "========================================="
	@cd services/identity-auth && \
	if [ "$(SKIP_TESTS)" = "true" ]; then \
		mvn clean package -DskipTests=true -Dmaven.javadoc.skip=true; \
	else \
		mvn clean package; \
	fi
	@# Copy the built jar to bin
	@cp services/identity-auth/target/*.jar $(BIN_DIR)/identity-auth.jar 2>/dev/null || true
	@echo "[Java] Copied identity-auth.jar to $(BIN_DIR)"

build-go:
	@echo "========================================="
	@echo "[Go] Compiling Go modules..."
	@echo "========================================="
	@for mod in $(GO_MODULES); do \
		if [ -d "$$mod" ]; then \
			name=$$(basename $$mod); \
			main_path=""; \
			if [ -f "$$mod/cmd/server/main.go" ]; then \
				main_path="cmd/server/main.go"; \
			elif [ -f "$$mod/cmd/main.go" ]; then \
				main_path="cmd/main.go"; \
			fi; \
			if [ -n "$$main_path" ]; then \
				echo "Building $$mod..."; \
				cd $$mod && go build -o $(BIN_DIR)/$$name$(EXE_EXT) $$main_path && cd $(CURDIR); \
			fi; \
		fi; \
	done
	@echo "[Go] Compilation completed successfully."
