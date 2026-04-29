param(
    [string]$ClusterName = "diabrisk"
)

$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "=== $Message ==="
}

function Require-Command {
    param([string]$Name)
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Nie znaleziono polecenia '$Name' w PATH."
    }
}

function Wait-DeploymentRollout {
    param(
        [string]$DeploymentName,
        [int]$TimeoutSeconds = 120
    )
    & kubectl rollout status "deployment/$DeploymentName" "--timeout=${TimeoutSeconds}s"
    if ($LASTEXITCODE -ne 0) {
        throw "Rollout deployment/$DeploymentName nie powiódł się."
    }
}

function Retry-KubectlApply {
    param(
        [string]$Input,
        [int]$MaxRetries = 5,
        [int]$DelaySeconds = 2
    )
    
    $retries = 0
    while ($retries -lt $MaxRetries) {
        $Input | & kubectl apply -f -
        if ($LASTEXITCODE -eq 0) {
            return $true
        }
        $retries++
        if ($retries -lt $MaxRetries) {
            Write-Host "kubectl apply failed, retrying in ${DelaySeconds}s... (attempt $retries/$MaxRetries)"
            Start-Sleep -Seconds $DelaySeconds
        }
    }
    return $false
}

function Retry-KubectlApplyFile {
    param(
        [string]$FilePath,
        [int]$MaxRetries = 5,
        [int]$DelaySeconds = 2
    )
    
    $retries = 0
    while ($retries -lt $MaxRetries) {
        & kubectl apply -f $FilePath
        if ($LASTEXITCODE -eq 0) {
            return $true
        }
        $retries++
        if ($retries -lt $MaxRetries) {
            Write-Host "kubectl apply -f $FilePath failed, retrying in ${DelaySeconds}s... (attempt $retries/$MaxRetries)"
            Start-Sleep -Seconds $DelaySeconds
        }
    }
    return $false
}

Require-Command "k3d"
Require-Command "kubectl"
Require-Command "docker"

Write-Step "Checking k3d cluster '$ClusterName'"

$clusterList = & k3d cluster list
if ($clusterList -match "^$ClusterName\s") {
    Write-Host "Cluster '$ClusterName' already exists. Deleting it first..."
    & k3d cluster delete $ClusterName
    if ($LASTEXITCODE -ne 0) {
        throw "Nie udało się usunąć istniejącego klastra '$ClusterName'."
    }
}

Write-Host "Creating k3d cluster '$ClusterName' with ports 80, 3001 and 9091 mapped..."
& k3d cluster create $ClusterName --api-port 6443 -p "80:80@loadbalancer" -p "3001:3000@loadbalancer" -p "9091:9090@loadbalancer"
if ($LASTEXITCODE -ne 0) {
    throw "Nie udało się utworzyć klastra k3d."
}
Write-Host "Cluster created successfully!"

Write-Host "Setting kubectl context and waiting for API server to be ready..."
& kubectl config use-context "k3d-$ClusterName"
if ($LASTEXITCODE -ne 0) {
    throw "Nie udało się ustawić kontekstu kubectl."
}

# Get and update kubeconfig to use localhost instead of host.docker.internal
Write-Host "Updating kubeconfig for Windows compatibility..."
try {
    $kubeConfig = & kubectl config view --raw --flatten
    # Replace host.docker.internal with 127.0.0.1
    $kubeConfig = $kubeConfig -replace 'host\.docker\.internal', '127.0.0.1'
    $kubeConfig | Set-Content "$env:USERPROFILE\.kube\config" -Force
} catch {
    Write-Warning "Could not update kubeconfig, continuing anyway: $_"
}

# Wait for API server to be ready with longer timeout
Write-Host "Waiting for Kubernetes API server to be ready..."
$apiReady = $false
for ($i = 0; $i -lt 60; $i++) {
    $output = & kubectl get nodes 2>&1
    if ($LASTEXITCODE -eq 0) {
        $apiReady = $true
        Write-Host "✅ API server is ready!"
        break
    }
    if ($i % 5 -eq 0) {
        Write-Host "  Still waiting... (${i}s)"
    }
    Start-Sleep -Seconds 1
}

if (-not $apiReady) {
    throw "API server did not become ready in time."
}


Write-Step "Applying Kubernetes secrets"

# postgres-secret
@"
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
type: Opaque
stringData:
  username: diabrisk
  password: diabrisk_dev_password
  database: diabrisk
  url: postgres://diabrisk:diabrisk_dev_password@postgres:5432/diabrisk?sslmode=disable
"@ | Retry-KubectlApply
if (-not $?) {
    throw "Nie udało się zastosować postgres-secret."
}

# auth-secret (for local authentication)
$adminEmail = $env:ADMIN_EMAIL
if ([string]::IsNullOrWhiteSpace($adminEmail)) {
    $adminEmail = "admin@diabrisk.local"
}

$adminPassword = $env:ADMIN_PASSWORD
if ([string]::IsNullOrWhiteSpace($adminPassword)) {
    $adminPassword = "default_admin_password"
}

# Delete old auth-secret if exists - retry up to 3 times
$deleteRetries = 0
while ($deleteRetries -lt 3) {
    & kubectl delete secret auth-secret --ignore-not-found 2>$null
    if ($LASTEXITCODE -eq 0) {
        break
    }
    $deleteRetries++
    Start-Sleep -Seconds 1
}

Write-Host "Creating auth-secret with admin credentials..."
@"
apiVersion: v1
kind: Secret
metadata:
  name: auth-secret
type: Opaque
stringData:
  admin_email: $adminEmail
  admin_password: $adminPassword
"@ | Retry-KubectlApply
if (-not $?) {
    throw "Nie udało się utworzyć auth-secret."
}

Write-Step "Building Docker images"

& docker build -f Dockerfile.api-gateway -t diabrisk-api:dev services/api-gateway
if ($LASTEXITCODE -ne 0) { throw "Build diabrisk-api:dev nie powiódł się." }

& docker build -f Dockerfile.frontend -t diabrisk-frontend:dev frontend
if ($LASTEXITCODE -ne 0) { throw "Build diabrisk-frontend:dev nie powiódł się." }

& docker build -f Dockerfile.data-svc -t diabrisk-data:dev services/data-svc
if ($LASTEXITCODE -ne 0) { throw "Build diabrisk-data:dev nie powiódł się." }

& docker build -f Dockerfile.auth-svc -t diabrisk-auth:dev services/auth-svc
if ($LASTEXITCODE -ne 0) { throw "Build diabrisk-auth:dev nie powiódł się." }

& docker build -f Dockerfile.ml-api -t diabrisk-ml:dev .
if ($LASTEXITCODE -ne 0) { throw "Build diabrisk-ml:dev nie powiódł się." }

Write-Step "Importing images into k3d cluster '$ClusterName'"

& k3d image import `
    --cluster $ClusterName `
    diabrisk-api:dev `
    diabrisk-frontend:dev `
    diabrisk-data:dev `
    diabrisk-auth:dev `
    diabrisk-ml:dev
if ($LASTEXITCODE -ne 0) {
    throw "Import obrazów do klastra nie powiódł się."
}

Write-Step "Applying Kubernetes manifests"

if (-not (Retry-KubectlApplyFile "deploy/k8s/postgres.yaml")) {
    throw "Nie udało się wdrożyć postgres.yaml."
}

Write-Host "Waiting for PostgreSQL to be ready..."
& kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s
if ($LASTEXITCODE -ne 0) { throw "PostgreSQL nie osiągnął stanu ready." }

if (-not (Retry-KubectlApplyFile "deploy/k8s/data-svc.yaml")) {
    throw "Nie udało się wdrożyć data-svc.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/auth-svc.yaml")) {
    throw "Nie udało się wdrożyć auth-svc.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/api-gateway.yaml")) {
    throw "Nie udało się wdrożyć api-gateway.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/ml-api.yaml")) {
    throw "Nie udało się wdrożyć ml-api.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/frontend.yaml")) {
    throw "Nie udało się wdrożyć frontend.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/ingress.yaml")) {
    throw "Nie udało się wdrożyć ingress.yaml."
}

if (-not (Retry-KubectlApplyFile "deploy/k8s/monitoring.yaml")) {
    throw "Failed to deploy monitoring.yaml."
}

Write-Step "Waiting for deployments to be ready"

Wait-DeploymentRollout -DeploymentName "data-svc" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "auth-svc" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "ml-api" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "api-gateway" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "frontend" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "prometheus" -TimeoutSeconds 120
Wait-DeploymentRollout -DeploymentName "grafana" -TimeoutSeconds 180

Write-Step "Running database migrations"
Write-Host "Waiting for data-svc to be fully ready..."
Start-Sleep -Seconds 5

$dataSvcPod = & kubectl get pod -l app=data-svc -o jsonpath="{.items[0].metadata.name}"
if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($dataSvcPod)) {
    throw "Nie udało się pobrać nazwy poda data-svc."
}

Write-Host "Running migrations in pod: $dataSvcPod"

# Uruchom migracje, ale ubij proces po 10 sekundach,
# bo według oryginalnego skryptu binary nie kończy się sam po migracjach.
$proc = Start-Process -FilePath "kubectl" `
    -ArgumentList @("exec", $dataSvcPod, "--", "/root/data-svc", "migrate", "up") `
    -NoNewWindow `
    -PassThru

$finished = $true
try {
    Wait-Process -Id $proc.Id -Timeout 10
} catch {
    $finished = $false
}

if (-not $finished -and -not $proc.HasExited) {
    Write-Host "Migration process still running after 10s, stopping it..."
    Stop-Process -Id $proc.Id -Force
}

Write-Host "Migrations completed!"

Write-Step "Current resources"
& kubectl get pods
& kubectl get svc
& kubectl get ingress

Write-Host ""
Write-Host "Setup complete!"
Write-Host ""
Write-Host "Access the application at: http://localhost"
Write-Host "Prometheus: http://localhost:9091"
Write-Host "Grafana: http://localhost:3001"
Write-Host "Grafana login: admin / admin"
Write-Host ""
Write-Host "Default admin credentials:"
Write-Host "  Email: $adminEmail"
Write-Host "  Password: $adminPassword"
Write-Host ""


Write-Host "  Prometheus: http://localhost:9091"
Write-Host "  Grafana: http://localhost:3001"
Write-Host ""
Write-Host "To view logs:"
Write-Host "  kubectl logs -f -l app=api-gateway"
Write-Host "  kubectl logs -f -l app=auth-svc"
Write-Host "  kubectl logs -f -l app=data-svc"
Write-Host "  kubectl logs -f -l app=postgres"
Write-Host "  kubectl logs -f -l app=prometheus"
Write-Host "  kubectl logs -f -l app=grafana"
Write-Host ""
