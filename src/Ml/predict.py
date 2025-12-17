import numpy as np
from joblib import load
from pathlib import Path

# Root projektu: DiabRisk/
BASE_DIR = Path(__file__).resolve().parents[2]
MODEL_PATH = BASE_DIR / "models" / "diabrisk_screening.joblib"

_artifact = None


def load_model():
    global _artifact
    if _artifact is None:
        _artifact = load(MODEL_PATH)
    return _artifact


def predict_risk(features: dict) -> dict:
    artifact = load_model()

    model1 = artifact["model1"]
    model2 = artifact["model2"]
    feature_names = artifact["feature_names"]

    X = np.array([[features[f] for f in feature_names]])

    p_at_risk = model1.predict_proba(X)[0, 1]
    p_prediabetes = model2.predict_proba(X)[0, 1]

    return {
        "risk_at_risk": float(p_at_risk),
        "risk_prediabetes": float(p_at_risk * p_prediabetes),
        "risk_diabetes": float(p_at_risk * (1 - p_prediabetes)),
        "risk_healthy": float(1 - p_at_risk),
    }
