# Monitoring (Prometheus)

This directory contains a basic Prometheus configuration for running
Prometheus outside Kubernetes.

For the local k3d stack, the install scripts deploy Prometheus and Grafana
from `deploy/k8s/monitoring.yaml`. In that default path, open:

```text
http://localhost:9091
http://localhost:3001
```

Grafana is provisioned with Prometheus as the default data source. The local
credentials are:

```text
admin / admin
```

## Running Prometheus

If you want to run Prometheus manually on the host instead of using the
Kubernetes manifest, run this from the repository root:

```sh
prometheus.exe --config.file=obs/prometheus.yml
```

By default, the Prometheus UI will be available at:

```text
http://localhost:9090
```

## Local Targets

The `prometheus.yml` file assumes the following local ports:

| Service | Port | Metrics endpoint |
| --- | ---: | --- |
| `api-gateway` | `8080` | `http://localhost:8080/metrics` |
| `auth-svc` | `8081` | `http://localhost:8081/metrics` |
| `data-svc` | `8082` | `http://localhost:8082/metrics` |
| `ml-api` | `8000` | `http://localhost:8000/metrics` |

When running the stack in Kubernetes, expose the services locally with
`kubectl port-forward`, for example:

```sh
kubectl port-forward svc/api-gateway 8080:8080
kubectl port-forward svc/auth-svc 8081:8081
kubectl port-forward svc/data-svc 8082:8082
kubectl port-forward svc/ml-api 8000:8000
```

## Useful Queries

Use these PromQL queries in Prometheus or Grafana Explore:

```promql
up
```

```promql
sum by (job) (up)
```

```promql
sum by (job, status) (rate({__name__=~"diabrisk_.*_http_requests_total"}[5m]))
```

```promql
histogram_quantile(
  0.95,
  sum by (job, le) (rate({__name__=~"diabrisk_.*_http_request_duration_seconds_bucket"}[5m]))
)
```

```promql
sum by (category) (diabrisk_ml_api_predictions_total)
```

## Metrics

The services expose Prometheus metrics for:

- request totals by method, route, and status,
- request duration histograms,
- response size histograms,
- ML prediction totals by risk category.
