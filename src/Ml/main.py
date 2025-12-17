import os
import numpy as np
import pandas as pd
from pathlib import Path
from joblib import dump

from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, confusion_matrix, roc_auc_score
from sklearn.model_selection import train_test_split
from imblearn.over_sampling import SMOTE


class CascadeDiabetesModel:
    def __init__(self, model1, model2):
        self.model1 = model1
        self.model2 = model2

    def predict(self, X):
        m1_pred = self.model1.predict(X)
        final_pred = np.zeros(len(X), dtype=int)

        at_risk_idx = np.where(m1_pred == 1)[0]
        if len(at_risk_idx) > 0:
            m2_pred = self.model2.predict(X[at_risk_idx])
            final_pred[at_risk_idx[m2_pred == 1]] = 1
            final_pred[at_risk_idx[m2_pred == 0]] = 2

        return final_pred


def train_cascade_with_splits(X_train, X_test, y_train, y_test, random_state=42):

    y_train_m1 = np.where(y_train == 0, 0, 1)
    y_test_m1 = np.where(y_test == 0, 0, 1)

    model1 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state,
        class_weight="balanced"
    )
    model1.fit(X_train, y_train_m1)

    print("\n=== MODEL 1 ===")
    print(classification_report(y_test_m1, model1.predict(X_test), digits=3))

    train_mask = y_train != 0
    test_mask = y_test != 0

    model2 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state,
        class_weight="balanced"
    )
    model2.fit(X_train[train_mask], np.where(y_train[train_mask] == 1, 1, 0))

    print("\n=== MODEL 2 ===")
    print(classification_report(
        np.where(y_test[test_mask] == 1, 1, 0),
        model2.predict(X_test[test_mask]),
        digits=3
    ))

    return CascadeDiabetesModel(model1, model2)


def main():
    random_state = 42
    target_col = "Diabetes_012"

    base_dir = Path(__file__).resolve().parents[2]
    csv_path = base_dir / "data" / "raw" / "diabetes_012_health_indicators_BRFSS2015.csv"

    df = pd.read_csv(csv_path)

    X = df.drop(columns=[target_col])
    y = df[target_col].values.astype(int)

    smote = SMOTE(random_state=random_state)
    X_res, y_res = smote.fit_resample(X.values, y)

    X_train, X_test, y_train, y_test = train_test_split(
        X_res, y_res, test_size=0.2, random_state=random_state, stratify=y_res
    )

    print("\n>>> Training cascade model...")
    cascade = train_cascade_with_splits(X_train, X_test, y_train, y_test)

    # ===== SAVE MODEL (ARTEFAKT) =====
    model_dir = base_dir / "models"
    model_dir.mkdir(exist_ok=True)

    dump(
        {
            "model1": cascade.model1,
            "model2": cascade.model2,
            "feature_names": list(X.columns),
            "type": "screening"
        },
        model_dir / "diabrisk_screening.joblib"
    )

    print(">>> Model saved to models/diabrisk_screening.joblib")


if __name__ == "__main__":
    main()
