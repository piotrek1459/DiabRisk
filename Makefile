PYTHON ?= python

.PHONY: test test-go test-auth test-gateway test-ml

test: test-go test-ml  ## Run all tests

test-go: test-auth test-gateway  ## Run all Go tests

test-auth:  ## Run auth-svc tests
	cd services/auth-svc && go test -v ./...

test-gateway:  ## Run api-gateway tests
	cd services/api-gateway && go test -v ./...

test-ml:  ## Run ML API unit tests
	$(PYTHON) -m unittest discover -s tests -p "test_*.py"
