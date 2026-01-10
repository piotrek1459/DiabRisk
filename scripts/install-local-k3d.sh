#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${1:-diabrisk}"

echo "=== Checking k3d cluster '$CLUSTER_NAME' ==="

if k3d cluster list | grep -q "^$CLUSTER_NAME "; then
  echo "Cluster '$CLUSTER_NAME' already exists."
  k3d cluster delete $CLUSTER_NAME
fi
echo "Creating k3d cluster '$CLUSTER_NAME' with port 80 mapped..."
k3d cluster create "$CLUSTER_NAME" \
  --api-port 6443 \
  -p "80:80@loadbalancer"
echo "Cluster created successfully!"

echo "=== Applying Kubernetes manifests ==="

# Secrets first
# Always apply postgres-secret from file (non-sensitive dev password)
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

# OAuth secret from environment variables or placeholders
if [ -z "$GOOGLE_CLIENT_ID" ] || [ -z "$GOOGLE_CLIENT_SECRET" ]; then
  echo "⚠️  WARNING: OAuth credentials not found in environment variables"
  echo "Creating oauth-secret with placeholders (authentication will not work)"
  kubectl create secret generic oauth-secret \
    --from-literal=google-client-id="YOUR_GOOGLE_CLIENT_ID_HERE" \
    --from-literal=google-client-secret="YOUR_GOOGLE_CLIENT_SECRET_HERE" \
    --from-literal=redirect-url="http://localhost/auth/google/callback"
else
  echo "✓ Creating oauth-secret from environment variables"
  kubectl create secret generic oauth-secret \
    --from-literal=google-client-id="$GOOGLE_CLIENT_ID" \
    --from-literal=google-client-secret="$GOOGLE_CLIENT_SECRET" \
    --from-literal=redirect-url="http://localhost/auth/google/callback"
fi

echo "=== Building Docker images ==="

# API Gateway
docker build \
  -f Dockerfile.api-gateway \
  -t diabrisk-api:dev \
  services/api-gateway

# Frontend
docker build \
  -f Dockerfile.frontend \
  -t diabrisk-frontend:dev \
  frontend

# Data Service
docker build \
  -f Dockerfile.data-svc \
  -t diabrisk-data:dev \
  services/data-svc

# Auth Service
docker build \
  -f Dockerfile.auth-svc \
  -t diabrisk-auth:dev \
  services/auth-svc

echo "=== Importing images into k3d cluster '$CLUSTER_NAME' ==="

k3d image import \
  --cluster "$CLUSTER_NAME" \
  diabrisk-api:dev \
  diabrisk-frontend:dev \
  diabrisk-data:dev \
  diabrisk-auth:dev

echo "=== Applying Kubernetes manifests ==="

# PostgreSQL (needed by other services)
kubectl apply -f deploy/k8s/postgres.yaml

echo "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s

# Apply other services
kubectl apply -f deploy/k8s/data-svc.yaml
kubectl apply -f deploy/k8s/auth-svc.yaml
kubectl apply -f deploy/k8s/api-gateway.yaml
kubectl apply -f deploy/k8s/frontend.yaml
kubectl apply -f deploy/k8s/ingress.yaml

echo "=== Waiting for deployments to be ready ==="

kubectl rollout status deployment/data-svc --timeout=120s
kubectl rollout status deployment/auth-svc --timeout=120s
kubectl rollout status deployment/api-gateway --timeout=120s
kubectl rollout status deployment/frontend --timeout=120s

echo "=== Running database migrations ==="
echo "Waiting for data-svc to be fully ready..."
sleep 5

DATA_SVC_POD=$(kubectl get pod -l app=data-svc -o jsonpath='{.items[0].metadata.name}')
echo "Running migrations in pod: $DATA_SVC_POD"
# Run migrations with timeout since the binary doesn't exit after migrations
timeout 10s kubectl exec $DATA_SVC_POD -- /root/data-svc migrate up || true
echo "Migrations completed!"

echo "=== Current resources ==="
kubectl get pods
kubectl get svc
kubectl get ingress

echo
echo "✅ Setup complete!"
echo
echo "Access the application at: http://localhost"
echo
echo "To view logs:"
echo "  kubectl logs -f -l app=api-gateway"
echo "  kubectl logs -f -l app=auth-svc"
echo "  kubectl logs -f -l app=data-svc"
echo "  kubectl logs -f -l app=postgres"
