#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Setup script variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="${SCRIPT_DIR}/bin"
RUN_TESTS=false
SKIP_JAVA=false
SKIP_GO=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --run-tests) RUN_TESTS=true ;;
        --skip-java) SKIP_JAVA=true ;;
        --skip-go) SKIP_GO=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

# Initialize output folder
rm -rf "${BIN_DIR}"
mkdir -p "${BIN_DIR}"

echo -e "\033[0;36m=========================================\033[0m"
echo -e "\033[0;36mStarting Build Process for Shopee Clone\033[0m"
echo -e "\033[0;36m=========================================\033[0m"

# 1. Build Java Service
if [ "$SKIP_JAVA" = false ]; then
    JAVA_DIR="${SCRIPT_DIR}/services/identity-auth"
    if [ -d "$JAVA_DIR" ]; then
        echo -e "\n\033[0;33m[Java] Building identity-auth service...\033[0m"
        cd "$JAVA_DIR"
        
        MVN_ARGS=("clean" "package")
        if [ "$RUN_TESTS" = false ]; then
            MVN_ARGS+=("-DskipTests=true" "-Dmaven.javadoc.skip=true")
            echo "-> Skipping Java tests (auto-bypass enabled)"
        fi
        
        mvn "${MVN_ARGS[@]}"
        
        # Copy jar to bin
        JAR_FILE=$(find target -maxdepth 1 -name "*.jar" ! -name "*original*" ! -name "*sources*" ! -name "*javadoc*" | head -n 1)
        if [ -n "$JAR_FILE" ]; then
            cp "$JAR_FILE" "${BIN_DIR}/identity-auth.jar"
            echo -e "\033[0;32m[Java] Success: Copied identity-auth.jar to bin/\033[0m"
        else
            echo "Warning: Built JAR file not found"
        fi
        cd "$SCRIPT_DIR"
    fi
fi

# 2. Build Go Services and Platforms
if [ "$SKIP_GO" = false ]; then
    echo -e "\n\033[0;33m[Go] Scanning workspace modules...\033[0m"

    GO_MODULES=(
        # Services
        "services/auth" "services/cart" "services/catalog-product" "services/checkout"
        "services/gateway" "services/inventory" "services/order" "services/payment"
        "services/product" "services/product-catalog" "services/promotion" "services/shipment"
        # Platforms
        "platforms/advertising" "platforms/aiml" "platforms/analytics" "platforms/api-gateway"
        "platforms/billing" "platforms/developer" "platforms/fraud" "platforms/fraud-risk"
        "platforms/global-infra" "platforms/live-commerce" "platforms/live-scale" "platforms/logistics-delivery"
        "platforms/notification" "platforms/notification-campaign" "platforms/oms-fulfillment"
        "platforms/payment-ledger" "platforms/rec-vector" "platforms/recommendation" "platforms/search"
        "platforms/search-indexing" "platforms/service-mesh" "platforms/sre"
    )

    for MODULE in "${GO_MODULES[@]}"; do
        FULL_PATH="${SCRIPT_DIR}/${MODULE}"
        if [ -d "$FULL_PATH" ]; then
            MODULE_NAME=$(basename "${MODULE}")
            
            MAIN_PATH=""
            if [ -f "${FULL_PATH}/cmd/server/main.go" ]; then
                MAIN_PATH="./cmd/server/main.go"
            elif [ -f "${FULL_PATH}/cmd/main.go" ]; then
                MAIN_PATH="./cmd/main.go"
            fi
            
            if [ -n "$MAIN_PATH" ]; then
                echo -e "[Go] Building module: ${MODULE}..."
                cd "$FULL_PATH"
                go build -o "${BIN_DIR}/${MODULE_NAME}" "${MAIN_PATH}"
                echo -e "\033[0;32m[Go] Success: Built ${MODULE_NAME}\033[0m"
                cd "$SCRIPT_DIR"
            fi
        fi
    done
fi

echo -e "\n\033[0;36m=========================================\033[0m"
echo -e "\033[0;36mAll builds completed! Binaries are in: ${BIN_DIR}\033[0m"
echo -e "\033[0;36m=========================================\033[0m"
