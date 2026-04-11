# Testing

## Running Tests

```bash
# Run all tests
make test

# Run only auth-svc tests
make test-auth

# Run only api-gateway tests
make test-gateway

# Run only ML API tests
make test-ml
```

The ML API tests use the Python dependencies from `requirements.txt`.
If your default `python` does not point at the project environment, pass it explicitly:

```powershell
make test-ml PYTHON=src\Ml\env\Scripts\python.exe
```

## Test Structure

| Service | Test file | Coverage |
|---------|-----------|----------|
| auth-svc | `services/auth-svc/main_test.go` | Token hashing, random string generation, HTTP handlers |
| api-gateway | `services/api-gateway/main_test.go` | Env helpers, health endpoint |
| ml-api | `tests/test_ml_api.py`, `tests/test_ml_training.py` | Feature validation, raw request scaling for processed-model input, model path resolution, prediction response, processed training data loading |

## CI/CD

Tests run automatically on every pull request via GitHub Actions (`.github/workflows/test.yml`).

## Writing New Tests

- Place test files next to the code they test (`*_test.go`)
- Use `gin.SetMode(gin.TestMode)` and `net/http/httptest` for HTTP handler tests
- Keep unit tests independent of the database
- Keep ML API unit tests independent of the real joblib model by using fake model artifacts
