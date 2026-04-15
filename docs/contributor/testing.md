# Testing

## Running Tests

```bash
# Run all tests
make test

# Run only auth-svc tests
make test-auth

# Run only api-gateway tests
make test-gateway

# Run only data-svc tests
make test-data

# Run Go integration tests
make test-go-integration

# Run only ML API tests
make test-ml

# Run integration tests
make test-integration

# Recreate the local k3d cluster and run full-stack smoke tests
make test-k3d
```

The ML API tests use the Python dependencies from `requirements.txt`.
If your default `python` does not point at the project environment, pass it explicitly:

```powershell
make test-ml PYTHON=src\Ml\env\Scripts\python.exe
```

For the detailed Python unit test scope and conventions, see `tests/README.md`.

## Test Structure

| Service | Test file | Coverage |
|---------|-----------|----------|
| auth-svc | `services/auth-svc/main_test.go` | Token hashing, random string generation, HTTP handlers |
| auth-svc integration | `services/auth-svc/auth_integration_test.go` | Register, login, session validation, logout against a real PostgreSQL database |
| api-gateway | `services/api-gateway/main_test.go` | Env helpers, health endpoint |
| api-gateway integration | `services/api-gateway/gateway_integration_test.go` | Protected gateway routes with fake auth and ML upstream services |
| data-svc integration | `services/data-svc/database_integration_test.go` | Database readiness, migrations, schema verification |
| ml-api | `tests/test_ml_api.py`, `tests/test_ml_training.py` | Feature validation, raw request scaling for processed-model input, model path resolution, prediction response, processed training data loading |
| integration | `tests/integration/test_ml_api_http.py` | ML API HTTP request handling, joblib artifact loading, prediction response contract |
| smoke | `tests/smoke/test_k3d_stack.py` | Local k3d cluster setup, frontend ingress, auth flow, protected gateway routes, ML prediction flow |

## CI/CD

Tests run automatically on every pull request via GitHub Actions (`.github/workflows/test.yml`).

## Writing New Tests

- Place test files next to the code they test (`*_test.go`)
- Use `gin.SetMode(gin.TestMode)` and `net/http/httptest` for HTTP handler tests
- Keep unit tests independent of the database
- Keep gateway integration tests local by using `httptest` upstream services for auth and ML
- Put DB-backed Go integration tests behind the `integration` build tag
- Keep ML API unit tests independent of the real joblib model by using fake model artifacts
- Put slower cross-boundary checks in `tests/integration` and run them with `make test-integration`
- Run the full local cluster smoke test explicitly with `make test-k3d`; it recreates the k3d cluster by default

## Go Microservice Integration Tests

`make test-go-integration` runs the gateway integration tests plus Go tests
compiled with the `integration` build tag.

The gateway tests do not need Docker or PostgreSQL. They start fake auth and ML
HTTP servers and verify that protected routes validate the session cookie and
proxy requests correctly.

The auth and data-svc integration tests need a PostgreSQL URL. Set one of these
before running them against a reachable test database:

```powershell
$env:DATABASE_URL="postgres://diabrisk:diabrisk_dev_password@localhost:5432/diabrisk?sslmode=disable"
make test-go-integration
```

You can also scope the URLs:

- `AUTH_INTEGRATION_DATABASE_URL` for `auth-svc`
- `DATA_INTEGRATION_DATABASE_URL` for `data-svc`

## Full-Stack k3d Smoke Tests

`make test-k3d` runs `tests/smoke/test_k3d_stack.py`. By default it calls the
local install script, recreates the `diabrisk` k3d cluster, waits for the stack,
and checks the frontend plus the register/session/features/risk path through
the ingress.

Useful environment variables:

- `DIABRISK_SKIP_K3D_SETUP=1` skips cluster setup and tests the currently running stack.
- `DIABRISK_DELETE_K3D_AFTER=1` deletes the cluster after the smoke tests finish.
- `DIABRISK_K3D_CLUSTER=diabrisk-test` changes the cluster name passed to the install script.
- `DIABRISK_BASE_URL=http://localhost` changes the URL used by HTTP checks.
