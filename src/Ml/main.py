import os
import numpy as np
import pandas as pd

from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, confusion_matrix, roc_auc_score
from sklearn.model_selection import train_test_split
from imblearn.over_sampling import SMOTE


class CascadeDiabetesModel:
    """
    Cascade model logic:
      - model1 predicts: 0 (healthy) vs 1 (at-risk: prediabetes or diabetes)
      - model2 predicts: 0 (diabetes) vs 1 (prediabetes)
    Final output:
      0 → healthy
      1 → prediabetes
      2 → diabetes
    """

    def __init__(self, model1, model2):
        self.model1 = model1
        self.model2 = model2

    def predict(self, X):
        """
        Cascade inference logic.
        X → input feature data (same preprocessing as during training).
        Returns array of final class predictions: 0, 1, or 2.
        """
        m1_pred = self.model1.predict(X)  # 0 vs 1 (at-risk)

        final_pred = np.zeros(len(X), dtype=int)
        at_risk_idx = np.where(m1_pred == 1)[0]

        if len(at_risk_idx) > 0:
            X_at_risk = X[at_risk_idx]
            m2_pred = self.model2.predict(X_at_risk)  # 0=diabetes, 1=prediabetes

            final_pred[at_risk_idx[m2_pred == 1]] = 1  # prediabetes
            final_pred[at_risk_idx[m2_pred == 0]] = 2  # diabetes

        return final_pred


def train_cascade_with_splits(X_train, X_test, y_train, y_test, random_state=42):
    """
    Trains both models in the cascade:
      - Model 1: distinguishes healthy vs. at-risk (prediabetes + diabetes)
      - Model 2: identifies prediabetes vs. diabetes for at-risk samples
    X_train, X_test: feature matrices
    y_train, y_test: original labels (0, 1, 2) after global SMOTE resampling
    """

    # ===== Model 1: 0 vs (1 + 2) =====
    y_train_m1 = np.where(y_train == 0, 0, 1)
    y_test_m1 = np.where(y_test == 0, 0, 1)

    model1 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state,
        class_weight="balanced"
    )
    model1.fit(X_train, y_train_m1)

    y_proba_m1 = model1.predict_proba(X_test)[:, 1]
    y_pred_m1 = (y_proba_m1 >= 0.5).astype(int)

    print("\n=== MODEL 1: 0 (healthy) vs 1 (at-risk) ===")
    print(classification_report(y_test_m1, y_pred_m1, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test_m1, y_pred_m1))
    try:
        print("ROC AUC:", roc_auc_score(y_test_m1, y_proba_m1))
    except ValueError:
        print("ROC AUC could not be calculated.")

    # ===== Model 2: 1 (prediabetes) vs 2 (diabetes) =====
    train_mask_m2 = y_train != 0
    test_mask_m2 = y_test != 0

    X_train_m2 = X_train[train_mask_m2]
    y_train_m2_raw = y_train[train_mask_m2]

    X_test_m2 = X_test[test_mask_m2]
    y_test_m2_raw = y_test[test_mask_m2]

    # Map to binary classes: 1 (prediabetes) → 1, 2 (diabetes) → 0
    y_train_m2 = np.where(y_train_m2_raw == 1, 1, 0)
    y_test_m2 = np.where(y_test_m2_raw == 1, 1, 0)

    model2 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state,
        class_weight="balanced"
    )
    model2.fit(X_train_m2, y_train_m2)

    y_proba_m2 = model2.predict_proba(X_test_m2)[:, 1]
    y_pred_m2 = (y_proba_m2 >= 0.5).astype(int)

    print("\n=== MODEL 2: 0 (diabetes) vs 1 (prediabetes) ===")
    print(classification_report(y_test_m2, y_pred_m2, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test_m2, y_pred_m2))
    try:
        print("ROC AUC:", roc_auc_score(y_test_m2, y_proba_m2))
    except ValueError:
        print("ROC AUC could not be calculated.")

    # ===== Cascade evaluation =====
    cascade = CascadeDiabetesModel(model1, model2)
    y_cascade = cascade.predict(X_test)

    print("\n=== CASCADE RESULT: final labels 0/1/2 vs y_test ===")
    print(classification_report(y_test, y_cascade, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test, y_cascade))

    return cascade


def main():
    random_state = 42
    target_col = "Diabetes_012"

    # Resolve project root path for CSV - independent of current working directory
    base_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    csv_path = os.path.join(base_dir, "data", "raw", "diabetes_012_health_indicators_BRFSS2015.csv")

    print(">>> Loading data from:", csv_path)
    df = pd.read_csv(csv_path)

    if target_col not in df.columns:
        raise ValueError(
            f"Target column '{target_col}' not found in the dataset. "
            f"Available columns: {list(df.columns)}"
        )

    X = df.drop(columns=[target_col])
    y = df[target_col].values

    print("Original dataset size:", X.shape, "| Label size:", y.shape)
    try:
        print("Original class distribution:", np.bincount(y.astype(int)))
    except Exception:
        pass

    # ===== GLOBAL SMOTE on full 3-class dataset BEFORE split =====
    print("\n>>> Applying global SMOTE on 3 classes (0/1/2) prior to train-test split...")
    smote_global = SMOTE(random_state=random_state)
    X_res, y_res = smote_global.fit_resample(X.values, y)

    print("Dataset size after SMOTE:", X_res.shape, "| Label size:", y_res.shape)
    try:
        print("Balanced class distribution after SMOTE:", np.bincount(y_res.astype(int)))
    except Exception:
        pass

    # ===== Split balanced dataset =====
    X_train, X_test, y_train, y_test = train_test_split(
        X_res,
        y_res,
        test_size=0.2,
        random_state=random_state,
        stratify=y_res
    )

    print("\n>>> Training cascade model on SMOTE-balanced dataset...")
    _ = train_cascade_with_splits(X_train, X_test, y_train, y_test, random_state=random_state)
    print(">>> Training complete.")


if __name__ == "__main__":
    main()
