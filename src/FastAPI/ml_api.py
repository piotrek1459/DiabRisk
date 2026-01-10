from typing import Any, Dict, List, Optional
from pathlib import Path

import os
import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from joblib import load


BASE_DIR = Path(__file__).resolve().parent  # directory containing ml_api.py

CANDIDATES = [
    os.environ.get("MODEL_PATH"),                     # 1) explicit override (Docker/Prod)
    str(BASE_DIR / "models" / "model.joblib"),        # 2) local dev (new name)
    str(BASE_DIR / "models" / "diabrisk_screening.joblib"),  # 3) local dev (old name)
    "/opt/models/model.joblib",                       # 4) server default
]

MODEL_PATH = next(
    (Path(p).expanduser().resolve() for p in CANDIDATES if p and Path(p).expanduser().exists()),
    None
)

if MODEL_PATH is None:
    raise FileNotFoundError(
        "Model file not found. Tried: " + ", ".join([p for p in CANDIDATES if p])
    )

model = load(MODEL_PATH)
print(f"Loaded model from: {MODEL_PATH}")


_artifact: Optional[dict] = None


def load_artifact() -> dict:
    global _artifact
    if _artifact is None:
        if not MODEL_PATH.exists():
            raise RuntimeError(f"Model file not found: {MODEL_PATH}")
        _artifact = load(MODEL_PATH)
    return _artifact


def build_X(features: Dict[str, Any], feature_names: List[str]) -> np.ndarray:
    missing = [f for f in feature_names if f not in features]
    if missing:
        raise HTTPException(
            status_code=400,
            detail={"error": "Missing features", "missing": missing},
        )
    # Extra pola ignorujemy (ale możesz też zwracać warning)
    return np.array([[float(features[f]) for f in feature_names]], dtype=float)


def predict_class_and_probs(X: np.ndarray, model1, model2):
    """
    Zgodnie z Twoim treningiem kaskadowym:
    - model1: 0 vs (1+2)
    - model2: rozróżnia 1 vs 2 na próbkach "at risk"
    """
    # klasy końcowe: 0=healthy, 1=prediabetes, 2=diabetes
    m1_pred = model1.predict(X)  # 0/1
    y_pred = np.zeros((X.shape[0],), dtype=int)

    at_risk_idx = np.where(m1_pred == 1)[0]
    if len(at_risk_idx) > 0:
        m2_pred = model2.predict(X[at_risk_idx])  # 1 => prediabetes, 0 => diabetes (jak w treningu)
        y_pred[at_risk_idx[m2_pred == 1]] = 1
        y_pred[at_risk_idx[m2_pred == 0]] = 2

    # prawdopodobieństwa “screeningowe” tak jak w Twoim predict.py
    p_at_risk = model1.predict_proba(X)[:, 1]
    p_prediabetes = model2.predict_proba(X)[:, 1]

    probs = np.vstack([
        (1 - p_at_risk),               # healthy
        (p_at_risk * p_prediabetes),   # prediabetes
        (p_at_risk * (1 - p_prediabetes)),  # diabetes
    ]).T

    return y_pred, probs


def confusion_matrix_3(y_true: np.ndarray, y_pred: np.ndarray) -> List[List[int]]:
    cm = np.zeros((3, 3), dtype=int)
    for t, p in zip(y_true, y_pred):
        if t in (0, 1, 2) and p in (0, 1, 2):
            cm[t, p] += 1
    return cm.tolist()


# --- FastAPI ---
app = FastAPI(title="DiabRisk ML Service", version="1.0")


class PredictRequest(BaseModel):
    features: Dict[str, Any] = Field(..., description="Słownik: nazwa_cechy -> wartość (0/1/liczba)")


class PredictResponse(BaseModel):
    RiskPercent: float
    Category: str
    Message: str
   


class AccuracySample(BaseModel):
    features: Dict[str, Any]
    y_true: int = Field(..., description="0=healthy, 1=prediabetes, 2=diabetes")


class AccuracyRequest(BaseModel):
    samples: List[AccuracySample]


@app.get("/healthz")
def healthz():
    return {"status": "ok"}


@app.get("/features")
def features():
    art = load_artifact()
    return {"feature_names": art["feature_names"], "count": len(art["feature_names"])}


@app.post("/predict", response_model=PredictResponse)
def predict(req: PredictRequest):
    art = load_artifact()
    """
    model1 = art["model1"]
    model2 = art["model2"]

    X = build_X(req.features, feature_names)
    y_pred, probs = predict_class_and_probs(X, model1, model2)

    label_map = {0: "healthy", 1: "prediabetes", 2: "diabetes"}
    p = probs[0]
    artifact = load_model()

    model1 = artifact["model1"]
    model2 = artifact["model2"]
    feature_names = artifact["feature_names"]

    X = np.array([[features[f] for f in feature_names]])"""
    model1 = art["model1"]
    model2 = art["model2"]
    feature_names = art["feature_names"]
    X = build_X(req.features, feature_names)    
    p_at_risk = 1-model1.predict_proba(X)[0, 1]
    p_prediabetes = model2.predict_proba(X)[0, 1]
    category="low"
    message="Jesteś zdrowy jak koń"
    if p_at_risk>0.50:
        category="medium"
        message="visit  doctor house"
    if p_at_risk>0.80:
        category="high"
        message="you are so ill "

    return PredictResponse(
        RiskPercent=p_at_risk,        
        Category=category,
        Message=message
    )

    