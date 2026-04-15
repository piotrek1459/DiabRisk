# Python Tests

This directory contains the Python test layers for DiabRisk. The root-level
`test_*.py` files are fast ML unit tests. Slower cross-boundary checks live in
subdirectories:

- `integration/` for ML API HTTP integration tests,
- `smoke/` for full local k3d stack smoke tests.

## Scope

The current Python unit test suite covers two areas:

| File | Production code | What it protects |
|------|-----------------|------------------|
| `test_ml_api.py` | `src/FastAPI/ml_api.py` | ML API request validation, model path resolution, raw request preprocessing, prediction response contract |
| `test_ml_training.py` | `src/Ml/main.py` | Loading training and test splits from `data/processed` |
| `integration/test_ml_api_http.py` | `src/FastAPI/ml_api.py` | HTTP request handling, joblib artifact loading, prediction response contract |
| `smoke/test_k3d_stack.py` | local k3d deployment | Frontend ingress, auth flow, protected gateway routes, ML prediction flow |

These tests intentionally do not load the real `models/diabrisk_screening.joblib`
artifact. The real model is large and its predictive quality is not the goal of
this unit test layer.

## Running

Run only the Python unit tests:

```powershell
make test-ml PYTHON=src\Ml\env\Scripts\python.exe
```

Run the Python ML API integration tests:

```powershell
make test-integration PYTHON=src\Ml\env\Scripts\python.exe
```

Run the full local k3d smoke tests:

```powershell
make test-k3d PYTHON=src\Ml\env\Scripts\python.exe
```

Run all unit tests in the repository:

```powershell
make test PYTHON=src\Ml\env\Scripts\python.exe
```

On Linux/macOS, use the Python interpreter from your environment instead:

```bash
make test-ml PYTHON=.venv/bin/python
```

If your default `python` already has the project dependencies installed, this is
enough:

```bash
make test-ml
```

## Dependencies

The tests use the runtime dependencies from `requirements.txt`, including:

- `fastapi`
- `pydantic`
- `numpy`
- `pandas`
- `joblib`
- `scikit-learn`
- `imblearn`

The suite uses Python's built-in `unittest`, so there is no separate pytest
dependency.

## ML API Contract

The HTTP request contract stays raw and frontend-friendly. A request sends raw
feature values such as:

```json
{
  "features": {
    "BMI": 25.5,
    "MentHlth": 2,
    "PhysHlth": 1
  }
}
```

The model is trained on processed features from `data/processed`, where `BMI`,
`MentHlth`, and `PhysHlth` are scaled. Therefore `ml_api.py` owns runtime
preprocessing:

1. `build_X()` validates the request and builds a numeric vector in artifact
   `feature_names` order.
2. `prepare_X_for_model()` applies the training-time scaling.
3. `predict()` sends the processed vector to `model1.predict_proba()`.

The tests use `FakeModel` instead of the real model artifact to verify that the
API sends scaled data into model inference and maps the returned probability to
the response contract.

## Training Data Contract

`src/Ml/main.py` trains from precomputed processed splits:

- `data/processed/X_train_processed.csv`
- `data/processed/X_test_processed.csv`
- `data/processed/y_train.csv`
- `data/processed/y_test.csv`

`test_ml_training.py` verifies that `load_processed_splits()` reads those four
inputs, returns NumPy arrays, preserves feature names from the processed feature
columns, casts labels to integers, and fails if train/test feature columns drift
apart.

The training test mocks `pandas.read_csv`; it does not read the real CSV files
and does not train a model.

## What These Tests Do Not Cover

The root-level tests are unit tests. They deliberately do not cover:

- real model artifact loading,
- model accuracy or prediction quality,
- end-to-end calls through `api-gateway`,
- Kubernetes manifests, rollouts, probes, or resource usage,
- database behavior.

Those belong in smoke or integration tests.

## CI

GitHub Actions runs the Python unit tests in `.github/workflows/test.yml` using
Python 3.12:

```yaml
python -m pip install -r requirements.txt
make test-ml
```

The Go tests run in a separate CI job through `make test-go`.

## Adding Tests

Keep new Python unit tests fast and deterministic:

- Use fake model objects instead of `models/diabrisk_screening.joblib`.
- Do not train models in unit tests.
- Do not depend on Kubernetes, Postgres, or network calls.
- Prefer testing small functions such as request validation, preprocessing, and
  data-loading contracts.
- If a future test needs real artifacts or services, put it in a separate smoke
  or integration test layer instead of this directory's unit suite.
