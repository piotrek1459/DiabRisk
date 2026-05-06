# Kubernetes Manifests

This folder contains the manifests used by the local k3d deployment
scripts.

## Active Deployment Path

The install scripts apply these manifests during local setup:

- `postgres.yaml`
- `data-svc.yaml`
- `auth-svc.yaml`
- `api-gateway.yaml`
- `ml-api.yaml`
- `frontend.yaml`
- `ingress.yaml`
- `monitoring.yaml`

The supported entrypoints for local deployment are:

- `scripts/install-local-k3d.ps1`
- `scripts/install-local-k3d.sh`

## Secrets Used by the Current Runtime

The current stack requires two Kubernetes secrets:

### `postgres-secret`

- `username`
- `password`
- `database`
- `url`

### `auth-secret`

- `admin_email`
- `admin_password`

`secrets.yaml` contains local-development defaults, but the install
scripts also create these secret objects directly during deployment.

## Manifest Notes

| File | Current Role |
|------|--------------|
| `postgres.yaml` | PostgreSQL deployment and service |
| `data-svc.yaml` | migration runner deployment |
| `auth-svc.yaml` | auth service deployment and service |
| `api-gateway.yaml` | gateway deployment and service |
| `ml-api.yaml` | ML service deployment, service, and model-path config |
| `frontend.yaml` | frontend deployment and service |
| `ingress.yaml` | Traefik ingress for `localhost` |
| `monitoring.yaml` | Prometheus and Grafana deployments, services, and provisioning |
| `secrets.yaml` | local default secrets, useful when applying resources manually |
| `configmap.yaml` | legacy file, not applied by the current install scripts |

## Access Pattern

With the default manifests applied:

- `http://localhost/` -> frontend
- `http://localhost/auth/*` -> `api-gateway` -> `auth-svc`
- `http://localhost/api/*` -> `api-gateway` -> `ml-api`
- `http://localhost:9091/` -> Prometheus through the k3d load balancer
- `http://localhost:3001/` -> Grafana through the k3d load balancer

`/healthz` is implemented inside services, but it is not exposed through
the current ingress manifest.

Grafana is provisioned with Prometheus as the default data source.
Default local credentials are `admin` / `admin`.

## Useful Commands

```bash
kubectl get pods
kubectl get svc
kubectl get ingress
kubectl logs -f -l app=api-gateway
kubectl logs -f -l app=auth-svc
kubectl logs -f -l app=data-svc
kubectl logs -f -l app=ml-api
kubectl logs -f -l app=postgres
kubectl logs -f -l app=prometheus
kubectl logs -f -l app=grafana
```

## Rebuilding a Single Service

Example for `auth-svc`:

```bash
docker build -f Dockerfile.auth-svc -t diabrisk-auth:dev services/auth-svc
k3d image import --cluster diabrisk diabrisk-auth:dev
kubectl rollout restart deployment/auth-svc
kubectl rollout status deployment/auth-svc
```
