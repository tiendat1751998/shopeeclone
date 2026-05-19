#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Setup script variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKIP_JAVA=false
SKIP_GO=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --skip-java) SKIP_JAVA=true ;;
        --skip-go) SKIP_GO=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

echo -e "\033[0;36m=========================================\033[0m"
echo -e "\033[0;36mStarting Docker Image Build Process\033[0m"
echo -e "\033[0;36m=========================================\033[0m"

# 1. Build Java Image
if [ "$SKIP_JAVA" = false ]; then
    JAVA_DIR="${SCRIPT_DIR}/services/identity-auth"
    DOCKERFILE_PATH="${JAVA_DIR}/Dockerfile"
    if [ -f "$DOCKERFILE_PATH" ]; then
        echo -e "\n\033[0;33m[Java] Building docker image for identity-auth...\033[0m"
        IMAGE_NAME="ghcr.io/shopee-clone/identity-auth:latest"
        docker build -t "$IMAGE_NAME" -f "$DOCKERFILE_PATH" "$JAVA_DIR"
        echo -e "\033[0;32m[Java] Success: Built ${IMAGE_NAME}\033[0m"
    fi
fi

# 2. Build Go Images
if [ "$SKIP_GO" = false ]; then
    SHARED_PACKAGE_PATH="${SCRIPT_DIR}/packages/go-shared"
    GO_MODULES="services/auth services/cart services/catalog-product services/checkout services/gateway services/inventory services/order services/payment services/product services/product-catalog services/promotion services/shipment platforms/advertising platforms/aiml platforms/analytics platforms/api-gateway platforms/billing platforms/developer platforms/fraud platforms/fraud-risk platforms/global-infra platforms/live-commerce platforms/live-scale platforms/logistics-delivery platforms/notification platforms/notification-campaign platforms/oms-fulfillment platforms/payment-ledger platforms/rec-vector platforms/recommendation platforms/search platforms/search-indexing platforms/service-mesh platforms/sre platforms/user-behavior"

    for MODULE in $GO_MODULES; do
        FULL_PATH="${SCRIPT_DIR}/${MODULE}"
        if [ -d "$FULL_PATH" ]; then
            MODULE_NAME=$(basename "${MODULE}")
            
            DOCKERFILE_PATH=""
            if [ -f "${FULL_PATH}/Dockerfile" ]; then
                DOCKERFILE_PATH="${FULL_PATH}/Dockerfile"
            elif [ -f "${FULL_PATH}/deployments/Dockerfile" ]; then
                DOCKERFILE_PATH="${FULL_PATH}/deployments/Dockerfile"
            fi
            
            if [ -n "$DOCKERFILE_PATH" ]; then
                echo -e "\n\033[0;33m[Go] Building docker image for ${MODULE}...\033[0m"
                IMAGE_NAME="ghcr.io/shopee-clone/${MODULE_NAME}:latest"
                
                # Check if Dockerfile requires packages/go-shared
                if grep -q "packages/go-shared" "$DOCKERFILE_PATH"; then
                    echo "-> Module requires packages/go-shared. Creating temporary build context..."
                    
                    TEMP_CONTEXT_DIR="${SCRIPT_DIR}/temp_build_${MODULE_NAME}"
                    rm -rf "$TEMP_CONTEXT_DIR"
                    mkdir -p "$TEMP_CONTEXT_DIR"
                    
                    # Copy module and shared pkg into temp context
                    cp -r "${FULL_PATH}/"* "$TEMP_CONTEXT_DIR/"
                    mkdir -p "${TEMP_CONTEXT_DIR}/packages/go-shared"
                    cp -r "${SHARED_PACKAGE_PATH}/"* "${TEMP_CONTEXT_DIR}/packages/go-shared/"
                    
                    # Create adjusted temp Dockerfile
                    TEMP_DOCKERFILE="${TEMP_CONTEXT_DIR}/Dockerfile.build"
                    sed 's|COPY \.\./\.\./packages/go-shared /app/packages/go-shared|COPY packages/go-shared /packages/go-shared|g' "$DOCKERFILE_PATH" > "$TEMP_DOCKERFILE"
                    
                    docker build -t "$IMAGE_NAME" -f "$TEMP_DOCKERFILE" "$TEMP_CONTEXT_DIR"
                    
                    # Clean up
                    rm -rf "$TEMP_CONTEXT_DIR"
                else
                    docker build -t "$IMAGE_NAME" -f "$DOCKERFILE_PATH" "$FULL_PATH"
                fi
                echo -e "\033[0;32m[Go] Success: Built ${IMAGE_NAME}\033[0m"
            fi
        fi
    done
fi

echo -e "\n\033[0;36m=========================================\033[0m"
echo -e "\033[0;36mAll Docker image builds completed!\033[0m"
echo -e "\033[0;36m=========================================\033[0m"
