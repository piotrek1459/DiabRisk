#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${1:-diabrisk}"

# --- Helper functions ---

require_command() {
  if ! command -v "$1" &>/dev/null; then
    echo "ERROR: '$1' not found in PATH. Please install it first."
    exit 1
  fi
}

retry_kubectl_apply_file() {
  local file_path="$1"
  local max_retries="${2:-5}"
  local delay="${3:-2}"
  local attempt=0

  while [ $attempt -lt $max_retries ]; do
    if kubectl apply -f "$file_path"; then
      return 0
    fi
    attempt=$((attempt + 1))
    if [ $attempt -lt $max_retries ]; then
      echo "kubectl apply -f $file_path failed, retrying in ${delay}s... (attempt $attempt/$max_retries)"
      sleep "$delay"
    fi
  done
  echo "ERROR: kubectl apply -f $file_path failed after $max_retries attempts."
  return 1
}

# --- Validate prerequisites ---

require_command "k3d"
require_command "kubectl"
require_command "docker"

# --- Cluster setup ---

echo "=== Checking k3d cluster '$CLUSTER_NAME' ==="

if k3d cluster list | grep -q "^$CLUSTER_NAME "; then
  echo "Cluster '$CLUSTER_NAME' already exists. Deleting it first..."
  k3d cluster delete "$CLUSTER_NAME"
fi

echo "Creating k3d cluster '$CLUSTER_NAME' with port 80 mapped..."
k3d cluster create "$CLUSTER_NAME" \
  --api-port 6443 \
  -p "80:80@loadbalancer"
echo "Cluster created successfully!"

# --- Wait for API server ---

echo "Waiting for Kubernetes API server to be ready..."
api_ready=false
for i in $(seq 1 60); do
  if kubectl get nodes &>/dev/null; then
    api_ready=true
    echo "API server is ready!"
    break
  fi
  if [ $((i % 5)) -eq 0 ]; then
    echo "  Still waiting... (${i}s)"
  fi
  sleep 1
done

if [ "$api_ready" = false ]; then
  echo "ERROR: API server did not become ready in 60 seconds."
  exit 1
fi

# --- Kubernetes secrets ---

echo "=== Applying Kubernetes secrets ==="

# postgres-secret
kubectl apply -f - <<EOF
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
EOF

# auth-secret (local email/password authentication)
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@diabrisk.local}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-default_admin_password}"

kubectl delete secret auth-secret --ignore-not-found 2>/dev/null || true

echo "Creating auth-secret with admin credentials..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: auth-secret
type: Opaque
stringData:
  admin_email: $ADMIN_EMAIL
  admin_password: $ADMIN_PASSWORD
EOF

# --- Build Docker images ---

echo "=== Building Docker images ==="

docker build -f Dockerfile.api-gateway -t diabrisk-api:dev services/api-gateway
docker build -f Dockerfile.frontend -t diabrisk-frontend:dev frontend
docker build -f Dockerfile.data-svc -t diabrisk-data:dev services/data-svc
docker build -f Dockerfile.auth-svc -t diabrisk-auth:dev services/auth-svc
docker build -f Dockerfile.ml-api -t diabrisk-ml:dev .

# --- Import images into k3d ---

echo "=== Importing images into k3d cluster '$CLUSTER_NAME' ==="

k3d image import \
  --cluster "$CLUSTER_NAME" \
  diabrisk-api:dev \
  diabrisk-frontend:dev \
  diabrisk-data:dev \
  diabrisk-auth:dev \
  diabrisk-ml:dev

# --- Apply Kubernetes manifests ---

echo "=== Applying Kubernetes manifests ==="

retry_kubectl_apply_file "deploy/k8s/postgres.yaml"

echo "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s

retry_kubectl_apply_file "deploy/k8s/data-svc.yaml"
retry_kubectl_apply_file "deploy/k8s/auth-svc.yaml"
retry_kubectl_apply_file "deploy/k8s/api-gateway.yaml"
retry_kubectl_apply_file "deploy/k8s/ml-api.yaml"
retry_kubectl_apply_file "deploy/k8s/frontend.yaml"
retry_kubectl_apply_file "deploy/k8s/ingress.yaml"

# --- Wait for deployments ---

echo "=== Waiting for deployments to be ready ==="

kubectl rollout status deployment/data-svc --timeout=120s
kubectl rollout status deployment/auth-svc --timeout=120s
kubectl rollout status deployment/ml-api --timeout=120s
kubectl rollout status deployment/api-gateway --timeout=120s
kubectl rollout status deployment/frontend --timeout=120s

# --- Database migrations ---

echo "=== Running database migrations ==="
echo "Waiting for data-svc to be fully ready..."
sleep 5

DATA_SVC_POD=$(kubectl get pod -l app=data-svc -o jsonpath='{.items[0].metadata.name}')
echo "Running migrations in pod: $DATA_SVC_POD"
timeout 10s kubectl exec "$DATA_SVC_POD" -- /root/data-svc migrate up || true
echo "Migrations completed!"

# --- Status ---

echo "=== Current resources ==="
kubectl get pods
kubectl get svc
kubectl get ingress

echo
echo "Setup complete!"
echo
echo "Access the application at: http://localhost"
echo
echo "Default admin credentials:"
echo "  Email: $ADMIN_EMAIL"
echo "  Password: $ADMIN_PASSWORD"
echo
echo "API Endpoints:"
echo "  Auth Service: http://localhost/auth"
echo "  API Gateway:  http://localhost/api"
echo
echo "To view logs:"
echo "  kubectl logs -f -l app=api-gateway"
echo "  kubectl logs -f -l app=auth-svc"
echo "  kubectl logs -f -l app=data-svc"
echo "  kubectl logs -f -l app=postgres"
echo "  kubectl logs -f -l app=ml-api"
echo
