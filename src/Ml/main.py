import numpy as np
import pandas as pd
from pathlib import Path
from joblib import dump


RF_N_ESTIMATORS = 300
MODEL2_CLASS_WEIGHT = {0: 1, 1: 2}
MODEL2_SMOTE_SAMPLING_STRATEGY = 0.5


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
    from imblearn.over_sampling import SMOTE
    from sklearn.ensemble import RandomForestClassifier
    from sklearn.metrics import classification_report, confusion_matrix

    y_train_m1 = np.where(y_train == 0, 0, 1)
    y_test_m1 = np.where(y_test == 0, 0, 1)
    smote_m1 = SMOTE(random_state=random_state)
    X_train_m1_res, y_train_m1_res = smote_m1.fit_resample(X_train, y_train_m1)

    model1 = RandomForestClassifier(
        n_estimators=RF_N_ESTIMATORS,
        max_depth=12,
        min_samples_leaf=10,
        n_jobs=-1,
        random_state=random_state,
        class_weight="balanced",
    )
    model1.fit(X_train_m1_res, y_train_m1_res)

    print("\n=== MODEL 1 ===")
    print(f"train distribution [healthy, at-risk]: {np.bincount(y_train_m1)}")
    print(f"SMOTE distribution [healthy, at-risk]: {np.bincount(y_train_m1_res)}")
    print(classification_report(y_test_m1, model1.predict(X_test), digits=3))

    train_mask = y_train != 0
    test_mask = y_test != 0
    y_train_m2 = np.where(y_train[train_mask] == 1, 1, 0)
    y_test_m2 = np.where(y_test[test_mask] == 1, 1, 0)
    smote_m2 = SMOTE(
        random_state=random_state,
        sampling_strategy=MODEL2_SMOTE_SAMPLING_STRATEGY,
    )
    X_train_m2_res, y_train_m2_res = smote_m2.fit_resample(
        X_train[train_mask],
        y_train_m2,
    )

    model2 = RandomForestClassifier(
        n_estimators=RF_N_ESTIMATORS,
        max_depth=12,
        min_samples_leaf=10,
        n_jobs=-1,
        random_state=random_state,
        class_weight=MODEL2_CLASS_WEIGHT,
    )
    model2.fit(X_train_m2_res, y_train_m2_res)

    print("\n=== MODEL 2 ===")
    print(f"class_weight={MODEL2_CLASS_WEIGHT}")
    print(f"train distribution [diabetes, prediabetes]: {np.bincount(y_train_m2)}")
    print(f"SMOTE distribution [diabetes, prediabetes]: {np.bincount(y_train_m2_res)}")
    print(f"test distribution  [diabetes, prediabetes]: {np.bincount(y_test_m2)}")

    y_pred_m2 = model2.predict(X_test[test_mask])
    print(f"pred distribution  [diabetes, prediabetes]: {np.bincount(y_pred_m2)}")
    print("confusion matrix rows=true cols=pred [diabetes, prediabetes]:")
    print(confusion_matrix(y_test_m2, y_pred_m2, labels=[0, 1]))
    print(classification_report(y_test_m2, y_pred_m2, digits=3))

    return CascadeDiabetesModel(model1, model2)


def load_processed_splits(base_dir: Path):
    processed_dir = base_dir / "data" / "processed"

    X_train_df = pd.read_csv(processed_dir / "X_train_processed.csv")
    X_test_df = pd.read_csv(processed_dir / "X_test_processed.csv")
    y_train = pd.read_csv(processed_dir / "y_train.csv").iloc[:, 0].astype(int).to_numpy()
    y_test = pd.read_csv(processed_dir / "y_test.csv").iloc[:, 0].astype(int).to_numpy()

    if list(X_train_df.columns) != list(X_test_df.columns):
        raise ValueError("Processed train and test feature columns do not match")

    return (
        X_train_df.to_numpy(dtype=float),
        X_test_df.to_numpy(dtype=float),
        y_train,
        y_test,
        list(X_train_df.columns),
    )


def main():
    random_state = 42

    base_dir = Path(__file__).resolve().parents[2]
    X_train, X_test, y_train, y_test, feature_names = load_processed_splits(base_dir)

    print("\n>>> Training cascade model...")
    cascade = train_cascade_with_splits(
        X_train,
        X_test,
        y_train,
        y_test,
        random_state=random_state,
    )

    # ===== SAVE MODEL  =====
    model_dir = base_dir / "models"
    model_dir.mkdir(exist_ok=True)

    dump(
        {
            "model1": cascade.model1,
            "model2": cascade.model2,
            "feature_names": feature_names,
            "type": "screening",
            "training_data": "data/processed",
            "model_features": "processed",
            "api_input_features": "raw",
        },
        model_dir / "diabrisk_screening.joblib",
    )

    print(">>> Model saved to models/diabrisk_screening.joblib")


if __name__ == "__main__":
    main()
