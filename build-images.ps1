<#
.SYNOPSIS
    Builds Docker images for all services and platforms in the Shopee Clone monorepo.
.DESCRIPTION
    Scans the services and platforms directories, locates their Dockerfiles (either at the root 
    or inside the deployments folder), builds them using Docker, and tags them as 
    ghcr.io/shopee-clone/<module-name>:latest.
.PARAMETER SkipJava
    Switch to skip building Java images.
.PARAMETER SkipGo
    Switch to skip building Go images.
.EXAMPLE
    .\build-images.ps1
#>

param (
    [switch]$SkipJava = $false,
    [switch]$SkipGo = $false
)

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Starting Docker Image Build Process" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

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
    "platforms/search-indexing", "platforms/service-mesh", "platforms/sre", "platforms/user-behavior"
)

# 1. Build Java Image
if (-not $SkipJava) {
    $JavaDir = Join-Path $PSScriptRoot "services/identity-auth"
    $DockerfilePath = Join-Path $JavaDir "Dockerfile"
    
    if (Test-Path $DockerfilePath) {
        Write-Host "`n[Java] Building docker image for identity-auth..." -ForegroundColor Yellow
        $ImageName = "ghcr.io/shopee-clone/identity-auth:latest"
        
        & docker build -t $ImageName -f $DockerfilePath $JavaDir
        if ($LASTEXITCODE -ne 0) {
            Write-Host "[Java] Docker build FAILED for identity-auth" -ForegroundColor Red
            exit 1
        }
        Write-Host "[Java] Success: Built $ImageName" -ForegroundColor Green
    }
}

# 2. Build Go Images
if (-not $SkipGo) {
    $SharedPackagePath = Join-Path $PSScriptRoot "packages/go-shared"
    
    foreach ($Module in $GoModules) {
        $FullModulePath = Join-Path $PSScriptRoot $Module
        if (Test-Path $FullModulePath) {
            $ModuleName = Split-Path $Module -Leaf
            
            # Find Dockerfile (root level or deployments folder)
            $DockerfilePath = ""
            if (Test-Path (Join-Path $FullModulePath "Dockerfile")) {
                $DockerfilePath = Join-Path $FullModulePath "Dockerfile"
            } elseif (Test-Path (Join-Path $FullModulePath "deployments/Dockerfile")) {
                $DockerfilePath = Join-Path $FullModulePath "deployments/Dockerfile"
            }
            
            if ($DockerfilePath -ne "") {
                Write-Host "`n[Go] Building docker image for $Module..." -ForegroundColor Yellow
                $ImageName = "ghcr.io/shopee-clone/${ModuleName}:latest"
                
                # Check if Dockerfile contains references to shared package outside the build context
                $DockerfileContent = Get-Content $DockerfilePath -Raw
                $RequiresSharedPackage = $DockerfileContent -match "packages/go-shared"
                
                if ($RequiresSharedPackage) {
                    # Handle outside-context packages safely by creating a temp directory context
                    Write-Host "-> Module requires packages/go-shared. Creating temporary build context..." -ForegroundColor Gray
                    
                    $TempContextDir = Join-Path $PSScriptRoot "temp_build_${ModuleName}"
                    if (Test-Path $TempContextDir) { Remove-Item -Recurse -Force $TempContextDir }
                    New-Item -ItemType Directory -Path $TempContextDir | Out-Null
                    
                    # Copy module contents and go-shared into the temporary context
                    Copy-Item -Recurse -Force "${FullModulePath}\*" $TempContextDir
                    
                    $TempSharedPath = Join-Path $TempContextDir "packages/go-shared"
                    New-Item -ItemType Directory -Path $TempSharedPath | Out-Null
                    Copy-Item -Recurse -Force "${SharedPackagePath}\*" $TempSharedPath
                    
                    # Create a temporary Dockerfile adjusting the COPY path to be local
                    $TempDockerfilePath = Join-Path $TempContextDir "Dockerfile.build"
                    $TempDockerfileContent = $DockerfileContent -replace "COPY \.\./\.\./packages/go-shared /app/packages/go-shared", "COPY packages/go-shared /packages/go-shared"
                    $TempDockerfileContent = $TempDockerfileContent -replace "RUN go mod download", "# RUN go mod download"
                    $TempDockerfileContent = $TempDockerfileContent -replace "COPY \. \.", "COPY . .`nRUN go mod tidy"
                    Set-Content -Path $TempDockerfilePath -Value $TempDockerfileContent
                    
                    try {
                        & docker build -t $ImageName -f $TempDockerfilePath $TempContextDir
                        if ($LASTEXITCODE -ne 0) {
                            throw "Docker build failed"
                        }
                        Write-Host "[Go] Success: Built $ImageName" -ForegroundColor Green
                    }
                    catch {
                        Write-Host "[Go] Docker build FAILED for ${Module}: $_" -ForegroundColor Red
                        Remove-Item -Recurse -Force $TempContextDir
                        exit 1
                    }
                    
                    # Cleanup
                    Remove-Item -Recurse -Force $TempContextDir
                } else {
                    # Build directly using module folder as context
                    & docker build -t $ImageName -f $DockerfilePath $FullModulePath
                    if ($LASTEXITCODE -ne 0) {
                        Write-Host "[Go] Docker build FAILED for ${Module}" -ForegroundColor Red
                        exit 1
                    }
                    Write-Host "[Go] Success: Built $ImageName" -ForegroundColor Green
                }
            }
        }
    }
}

Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "All Docker image builds completed!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
