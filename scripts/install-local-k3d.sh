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


echo "=== Building Docker images ==="

# Backend
docker build \
  -f Dockerfile.api-gateway \
  -t diabrisk-api:dev \
  services/api-gateway

# Frontend
docker build \
  -f Dockerfile.frontend \
  -t diabrisk-frontend:dev \
  frontend

echo "=== Importing images into k3d cluster '$CLUSTER_NAME' ==="

k3d image import \
  --cluster "$CLUSTER_NAME" \
  diabrisk-api:dev \
  diabrisk-frontend:dev

echo "=== Applying Kubernetes manifests ==="

kubectl apply -f deploy/k8s/api-gateway.yaml
kubectl apply -f deploy/k8s/frontend.yaml
kubectl apply -f deploy/k8s/ingress.yaml

echo "=== Waiting for deployments to be ready ==="

kubectl rollout status deployment/api-gateway --timeout=60s
kubectl rollout status deployment/frontend --timeout=60s

echo "=== Current resources ==="
kubectl get pods
kubectl get svc
kubectl get ingress

echo
echo "If not yet done, add to /etc/hosts:"
echo "  127.0.0.1   diabrisk.local"
echo
echo "Then open: http://diabrisk.local"
