<#
.SYNOPSIS
    Builds all Go and Java services in the Shopee Clone monorepo.
.DESCRIPTION
    Loops through all services and platforms listed in the workspace, compiles the Go entrypoints
    to a centralized 'bin/' folder, and packages the Java 'identity-auth' service.
    By default, all tests (Maven tests) are skipped for speed.
.PARAMETER SkipJava
    Switch to skip building Java services.
.PARAMETER SkipGo
    Switch to skip building Go services.
.PARAMETER RunTests
    Switch to enable running tests (disabled by default).
.EXAMPLE
    .\build.ps1
#>

param (
    [switch]$SkipJava = $false,
    [switch]$SkipGo = $false,
    [switch]$RunTests = $false
)

$ErrorActionPreference = "Stop"

# Clear or create output directory
$BinDir = Join-Path $PSScriptRoot "bin"
if (Test-Path $BinDir) {
    Remove-Item -Recurse -Force $BinDir
}
New-Item -ItemType Directory -Path $BinDir | Out-Null

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Starting Build Process for Shopee Clone" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 1. Build Java Service
if (-not $SkipJava) {
    $JavaDir = Join-Path $PSScriptRoot "services/identity-auth"
    if (Test-Path $JavaDir) {
        Write-Host "`n[Java] Building identity-auth service..." -ForegroundColor Yellow
        Push-Location $JavaDir
        try {
            $MvnCmd = "mvn"
            $MvnArgs = @("clean", "package")
            if (-not $RunTests) {
                $MvnArgs += "-DskipTests=true"
                $MvnArgs += "-Dmaven.javadoc.skip=true"
                Write-Host "-> Skipping Java tests (auto-bypass enabled)" -ForegroundColor Gray
            }
            
            # Execute maven directly to preserve console output & handle errors properly
            & $MvnCmd $MvnArgs
            if ($LASTEXITCODE -ne 0) {
                throw "Maven build failed with exit code $LASTEXITCODE"
            }
            
            # Copy built jar to bin
            $JarFile = Get-ChildItem -Path "target" -Filter "*.jar" | Where-Object { $_.Name -notmatch "original|sources|javadoc" } | Select-Object -First 1
            if ($JarFile) {
                Copy-Item $JarFile.FullName (Join-Path $BinDir "identity-auth.jar")
                Write-Host "[Java] Success: Copied identity-auth.jar to bin/" -ForegroundColor Green
            } else {
                Write-Warning "Could not find built JAR file in target/ directory."
            }
        }
        catch {
            Write-Host "[Java] Build FAILED: $_" -ForegroundColor Red
            Pop-Location
            exit 1
        }
        Pop-Location
    }
}

# 2. Build Go Services and Platforms
if (-not $SkipGo) {
    Write-Host "`n[Go] Scanning workspace modules..." -ForegroundColor Yellow
    
    # Target directories based on go.work use list
    $GoModules = @(
        # Services
        "services/auth", "services/cart", "services/catalog-product", "services/checkout",
        "services/gateway", "services/inventory", "services/order", "services/payment",
        "services/product", "services/product-catalog", "services/promotion", "services/shipment",
        # Platforms
        "platforms/advertising", "platforms/aiml", "platforms/analytics", "platforms/api-gateway",
        "platforms/billing", "platforms/developer", "platforms/fraud", "platforms/fraud-risk",
        "platforms/global-infra", "platforms/live-commerce", "platforms/live-scale", "platforms/logistics-delivery",
        "platforms/notification", "platforms/notification-campaign", "platforms/oms-fulfillment",
        "platforms/payment-ledger", "platforms/rec-vector", "platforms/recommendation", "platforms/search",
        "platforms/search-indexing", "platforms/service-mesh", "platforms/sre"
    )

    foreach ($Module in $GoModules) {
        $FullModulePath = Join-Path $PSScriptRoot $Module
        if (Test-Path $FullModulePath) {
            $ModuleName = Split-Path $Module -Leaf
            
            # Find main entry point
            $MainPath = ""
            if (Test-Path (Join-Path $FullModulePath "cmd/server/main.go")) {
                $MainPath = "./cmd/server/main.go"
            } elseif (Test-Path (Join-Path $FullModulePath "cmd/main.go")) {
                $MainPath = "./cmd/main.go"
            }

            if ($MainPath -ne "") {
                Write-Host "[Go] Building module: $Module..." -ForegroundColor DarkYellow
                Push-Location $FullModulePath
                try {
                    $OutName = $ModuleName
                    if ($IsWindows) { $OutName += ".exe" }
                    $OutPath = Join-Path $BinDir $OutName
                    
                    # Compile directly
                    & go build -o $OutPath $MainPath
                    if ($LASTEXITCODE -ne 0) {
                        throw "Go build failed with exit code $LASTEXITCODE"
                    }
                    Write-Host "[Go] Success: Built $OutName" -ForegroundColor Green
                }
                catch {
                    Write-Host "[Go] Build FAILED for ${Module}: $_" -ForegroundColor Red
                    Pop-Location
                    exit 1
                }
                Pop-Location
            }
        }
    }
}

Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "All builds completed! Binaries are in: $BinDir" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
