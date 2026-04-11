from pathlib import Path
from typing import Any, Dict, Optional

import os
import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from joblib import load


# --- scaler data ---
SCALER_MEAN = {
    "BMI": 28.391265570797856,
    "MentHlth": 3.177388836329232,
    "PhysHlth": 4.232379375591297,
}

SCALER_SCALE = {
    "BMI": 6.616426071801123,
    "MentHlth": 7.401493591748837,
    "PhysHlth": 8.709293119641702,
}

# --- paths & cache ---
BASE_DIR = Path(__file__).resolve().parent
REPO_DIR = BASE_DIR.parents[1]

CANDIDATES = [
    os.environ.get("MODEL_PATH"),  # 1) explicit override (Docker/Prod)
    str(REPO_DIR / "models" / "model.joblib"),  # 2) local dev (new name)
    str(REPO_DIR / "models" / "diabrisk_screening.joblib"),  # 3) local dev (old name)
    str(BASE_DIR / "models" / "model.joblib"),  # 4) package-local fallback
    str(BASE_DIR / "models" / "diabrisk_screening.joblib"),  # 5) package-local fallback
    "/opt/models/model.joblib",  # 6) server default
]

MODEL_PATH: Optional[Path] = None
_artifact: Optional[dict] = None


def candidate_model_paths() -> list[Path]:
    return [Path(p).expanduser() for p in CANDIDATES if p]


def resolve_model_path() -> Path:
    for path in candidate_model_paths():
        if path.exists():
            return path.resolve()

    raise RuntimeError(
        "Model file not found. Tried: "
        + ", ".join(str(path) for path in candidate_model_paths())
    )


def apply_training_scaling(X, feature_names):
    X_scaled = X.copy()

    for feat in ("BMI", "MentHlth", "PhysHlth"):
        idx = feature_names.index(feat)
        X_scaled[0, idx] = (X_scaled[0, idx] - SCALER_MEAN[feat]) / SCALER_SCALE[feat]

    return X_scaled


def prepare_X_for_model(features: Dict[str, Any], feature_names: list[str]):
    X = build_X(features, feature_names)
    return apply_training_scaling(X, feature_names)


def load_artifact() -> dict:
    global MODEL_PATH, _artifact
    if _artifact is None:
        if MODEL_PATH is None:
            MODEL_PATH = resolve_model_path()
        if not MODEL_PATH.exists():
            raise RuntimeError(f"Model file not found: {MODEL_PATH}")
        _artifact = load(MODEL_PATH)
    return _artifact


def build_X(features: Dict[str, Any], feature_names: list[str]) -> np.ndarray:
    missing = [f for f in feature_names if f not in features]
    if missing:
        raise HTTPException(
            status_code=400,
            detail={"error": "Missing features", "missing": missing},
        )

    try:
        return np.array([[float(features[f]) for f in feature_names]], dtype=float)
    except (TypeError, ValueError):
        raise HTTPException(
            status_code=400,
            detail={"error": "Invalid feature value type"},
        )


# --- FastAPI ---
app = FastAPI(title="DiabRisk ML Service", version="1.0")


class PredictRequest(BaseModel):
    features: Dict[str, Any] = Field(
        ..., description="Dict: feature_name -> numeric value"
    )


class PredictResponse(BaseModel):
    RiskPercent: float
    Category: str
    Message: str


@app.get("/healthz")
def healthz():
    return {"status": "ok"}


@app.get("/features")
def features():
    art = load_artifact()
    return {
        "feature_names": art["feature_names"],
        "count": len(art["feature_names"]),
    }


@app.post("/predict", response_model=PredictResponse)
def predict(req: PredictRequest):
    art = load_artifact()

    model1 = art["model1"]
    feature_names = art["feature_names"]

    X = prepare_X_for_model(req.features, feature_names)

    risk_score = model1.predict_proba(X)[0, 1]

    # --- decision logic ---
    if risk_score > 0.80:
        category = "high"
        message = "High risk detected. Immediate medical consultation recommended."
    elif risk_score > 0.50:
        category = "medium"
        message = "Moderate risk detected. Consider scheduling a medical checkup."
    else:
        category = "low"
        message = "Low risk detected. No immediate action required."

    return PredictResponse(
        RiskPercent=risk_score,
        Category=category,
        Message=message,
    )
