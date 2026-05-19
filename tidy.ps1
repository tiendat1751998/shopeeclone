$GoModules = @(
    "packages/go-shared",
    "services/auth", "services/cart", "services/catalog-product", "services/checkout",
    "services/gateway", "services/inventory", "services/order", "services/payment",
    "services/product", "services/product-catalog", "services/promotion", "services/shipment",
    "platforms/advertising", "platforms/aiml", "platforms/analytics", "platforms/api-gateway",
    "platforms/billing", "platforms/developer", "platforms/fraud", "platforms/fraud-risk",
    "platforms/global-infra", "platforms/live-commerce", "platforms/live-scale", "platforms/logistics-delivery",
    "platforms/notification", "platforms/notification-campaign", "platforms/oms-fulfillment",
    "platforms/payment-ledger", "platforms/rec-vector", "platforms/recommendation", "platforms/search",
    "platforms/search-indexing", "platforms/service-mesh", "platforms/sre", "platforms/user-behavior"
)

foreach ($Module in $GoModules) {
    $Path = Join-Path $PSScriptRoot $Module
    if (Test-Path $Path) {
        Write-Host "Tidying $Module..."
        Push-Location $Path
        go mod tidy
        Pop-Location
    }
}
