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
                
                # Handle outside-context packages safely by creating a temp directory context
                Write-Host "Creating temporary build context for $ModuleName..." -ForegroundColor Gray
                
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
                $DockerfileContent = Get-Content $DockerfilePath -Raw
                
                # Copy central proto directory if needed
                $HasProto = $false
                $GoModPath = Join-Path $FullModulePath "go.mod"
                if ((Test-Path $GoModPath) -and ((Get-Content $GoModPath) -match "github.com/shopee-clone/shopee/proto") -or ($ModuleName -eq "catalog-product")) {
                    Write-Host "-> Copying central proto module into build context..." -ForegroundColor Gray
                    $TempProtoPath = Join-Path $TempContextDir "proto"
                    New-Item -ItemType Directory -Path $TempProtoPath | Out-Null
                    Copy-Item -Recurse -Force "${PSScriptRoot}\proto\*" $TempProtoPath
                    $HasProto = $true
                }
                
                # Apply replacements to go.mod inside the temp context
                $TempGoModPath = Join-Path $TempContextDir "go.mod"
                if (Test-Path $TempGoModPath) {
                    $GoModContent = Get-Content $TempGoModPath -Raw
                    $GoModContent = $GoModContent -replace "replace github\.com/shopee-clone/shopee/proto => \.\./\.\./proto", "replace github.com/shopee-clone/shopee/proto => /proto"
                    Set-Content -Path $TempGoModPath -Value $GoModContent
                }
                
                $TempDockerfileContent = $DockerfileContent -replace "; ", "`n"
                
                if ($TempDockerfileContent -match "packages/go-shared") {
                    $TempDockerfileContent = $TempDockerfileContent -replace "COPY \.\./\.\./packages/go-shared /app/packages/go-shared", "COPY packages/go-shared /packages/go-shared"
                } else {
                    $TempDockerfileContent = $TempDockerfileContent -replace "COPY go\.mod go\.sum \./", "COPY go.mod go.sum ./`nCOPY packages/go-shared /packages/go-shared"
                }
                
                # If proto directory was copied, inject COPY proto /proto into the Dockerfile
                if ($HasProto) {
                    $TempDockerfileContent = $TempDockerfileContent -replace "COPY packages/go-shared /packages/go-shared", "COPY packages/go-shared /packages/go-shared`nCOPY proto /proto"
                }
                
                $TempDockerfileContent = $TempDockerfileContent -replace "RUN go mod download", "# RUN go mod download"
                
                $ProtoLogic = "COPY . .`n" +
                              "RUN apt-get update && apt-get install -y protobuf-compiler libprotobuf-dev git wget && \`n" +
                              "    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \`n" +
                              "    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \`n" +
                              "    mkdir -p /tmp/validate/validate && \`n" +
                              "    wget -qO /tmp/validate/validate/validate.proto https://raw.githubusercontent.com/bufbuild/protoc-gen-validate/main/validate/validate.proto && \`n" +
                              "    if [ -d `"/proto/shopee/catalog/v1`" ]; then \`n" +
                              "        echo `"Compiling catalog.proto...`" && \`n" +
                              "        mkdir -p /proto/catalog/v1 && \`n" +
                              "        protoc --proto_path=/proto/shopee --proto_path=/tmp/validate --proto_path=/usr/include --go_out=/proto --go_opt=module=github.com/shopee-clone/shopee/proto --go-grpc_out=/proto --go-grpc_opt=module=github.com/shopee-clone/shopee/proto /proto/shopee/catalog/v1/catalog.proto || exit 1; \`n" +
                              "    fi && \`n" +
                              "    find . -name `"*.proto`" -type f -not -path `"./proto/shopee/*`" | while read f; do \`n" +
                              "        echo `"Compiling `$f...`" && \`n" +
                              "        protoc --proto_path=. --proto_path=/tmp/validate --proto_path=/usr/include --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative `"`$f`" || exit 1; \`n" +
                              "    done`n" +
                              "RUN go mod tidy"
                              
                $TempDockerfileContent = $TempDockerfileContent -replace "COPY \. \.", $ProtoLogic
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
            }
        }
    }
}

Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "All Docker image builds completed!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
