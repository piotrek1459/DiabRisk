# Testing

## Running Tests

```bash
# Run all tests
make test

# Run only auth-svc tests
make test-auth

# Run only api-gateway tests
make test-gateway
```

## Test Structure

| Service | Test file | Coverage |
|---------|-----------|----------|
| auth-svc | `services/auth-svc/main_test.go` | Token hashing, random string generation, HTTP handlers |
| api-gateway | `services/api-gateway/main_test.go` | Env helpers, health endpoint |

## CI/CD

Tests run automatically on every pull request via GitHub Actions (`.github/workflows/test.yml`).

## Writing New Tests

- Place test files next to the code they test (`*_test.go`)
- Use `gin.SetMode(gin.TestMode)` and `net/http/httptest` for HTTP handler tests
- Keep unit tests independent of the database
