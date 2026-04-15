PYTHON ?= python

.PHONY: test test-go test-auth test-gateway test-data test-go-integration test-ml test-integration test-k3d

test: test-go test-ml  ## Run all tests

test-go: test-auth test-gateway test-data  ## Run all fast Go tests

test-auth:  ## Run auth-svc tests
	cd services/auth-svc && go test -v ./...

test-gateway:  ## Run api-gateway tests
	cd services/api-gateway && go test -v ./...

test-data:  ## Run data-svc tests
	cd services/data-svc && go test -v ./...

test-go-integration:  ## Run Go integration tests; DATABASE_URL enables DB-backed auth/data checks
	cd services/api-gateway && go test -v ./...
	cd services/data-svc && go test -tags=integration -v ./...
	cd services/auth-svc && go test -tags=integration -v ./...

test-ml:  ## Run ML API unit tests
	$(PYTHON) -m unittest tests.test_ml_api tests.test_ml_training

test-integration:  ## Run integration tests
	$(PYTHON) -m unittest discover -s tests/integration -p "test_*.py"

test-k3d:  ## Recreate local k3d cluster and run full-stack smoke tests
	$(PYTHON) -m unittest discover -s tests/smoke -p "test_*.py"
