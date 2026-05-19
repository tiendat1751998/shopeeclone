#!/usr/bin/env bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

GO_MODULES="packages/go-shared services/auth services/cart services/catalog-product services/checkout services/gateway services/inventory services/order services/payment services/product services/product-catalog services/promotion services/shipment platforms/advertising platforms/aiml platforms/analytics platforms/api-gateway platforms/billing platforms/developer platforms/fraud platforms/fraud-risk platforms/global-infra platforms/live-commerce platforms/live-scale platforms/logistics-delivery platforms/notification platforms/notification-campaign platforms/oms-fulfillment platforms/payment-ledger platforms/rec-vector platforms/recommendation platforms/search platforms/search-indexing platforms/service-mesh platforms/sre platforms/user-behavior"

echo "Tidying all Go modules..."
for MODULE in $GO_MODULES; do
    if [ -d "${SCRIPT_DIR}/${MODULE}" ]; then
        echo "Tidying ${MODULE}..."
        cd "${SCRIPT_DIR}/${MODULE}"
        go mod tidy
    fi
done
cd "${SCRIPT_DIR}"
echo "All modules tidied successfully!"
