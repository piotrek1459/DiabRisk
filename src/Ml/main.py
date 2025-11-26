import numpy as np
import pandas as pd

from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, confusion_matrix, roc_auc_score
from sklearn.model_selection import train_test_split
from imblearn.over_sampling import SMOTE


class CascadeDiabetesModel:
    """
    Kaskadowy model:
      - model1: 0 (zdrowy) vs 1 (zagrożony: prediab+diab)
      - model2: 0 (diabetes) vs 1 (prediabetes)
    Zwraca finalnie: 0, 1, 2 (jak w oryginalnym Diabetes_012).
    """

    def __init__(self, model1, model2):
        self.model1 = model1
        self.model2 = model2

    def predict(self, X):
        """
        X – numpy array / DataFrame po tych samych przekształceniach co przy treningu.
        Zwraca wektor etykiet 0/1/2.
        """
        m1_pred = self.model1.predict(X)  # 0 vs 1

        final_pred = np.zeros(len(X), dtype=int)
        at_risk_idx = np.where(m1_pred == 1)[0]

        if len(at_risk_idx) > 0:
            X_at_risk = X[at_risk_idx]
            m2_pred = self.model2.predict(X_at_risk)  # 0=diab, 1=pre

            final_pred[at_risk_idx[m2_pred == 1]] = 1  # prediabetes
            final_pred[at_risk_idx[m2_pred == 0]] = 2  # diabetes

        return final_pred


def train_cascade_with_splits(X_train, X_test, y_train, y_test, random_state=42):
    """
    X_train, X_test – cechy (np. po skalowaniu/encodowaniu),
    y_train, y_test – etykiety 0/1/2 (oryginalne).
    Zwraca: obiekt CascadeDiabetesModel.
    """

    # ===== MODEL 1: 0 vs (1+2) =====
    y_train_m1 = np.where(y_train == 0, 0, 1)
    y_test_m1 = np.where(y_test == 0, 0, 1)

    smote1 = SMOTE(random_state=random_state)
    X_train_m1_res, y_train_m1_res = smote1.fit_resample(X_train, y_train_m1)

    model1 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state
    )
    model1.fit(X_train_m1_res, y_train_m1_res)

    y_pred_m1 = model1.predict(X_test)
    y_proba_m1 = model1.predict_proba(X_test)[:, 1]

    print("\n=== MODEL 1: 0 (zdrowy) vs 1 (zagrożony) ===")
    print(classification_report(y_test_m1, y_pred_m1, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test_m1, y_pred_m1))
    try:
        print("ROC AUC:", roc_auc_score(y_test_m1, y_proba_m1))
    except ValueError:
        print("Nie dało się policzyć AUC.")

    # ===== MODEL 2: 1 vs 2 (tylko próbki z klas 1 i 2) =====
    train_mask_m2 = y_train != 0
    test_mask_m2 = y_test != 0

    X_train_m2 = X_train[train_mask_m2]
    y_train_m2_raw = y_train[train_mask_m2]

    X_test_m2 = X_test[test_mask_m2]
    y_test_m2_raw = y_test[test_mask_m2]

    # Mapowanie: 1 (prediabetes) -> 1, 2 (diabetes) -> 0
    y_train_m2 = np.where(y_train_m2_raw == 1, 1, 0)
    y_test_m2 = np.where(y_test_m2_raw == 1, 1, 0)

    smote2 = SMOTE(random_state=random_state)
    X_train_m2_res, y_train_m2_res = smote2.fit_resample(X_train_m2, y_train_m2)

    model2 = RandomForestClassifier(
        n_estimators=200,
        n_jobs=-1,
        random_state=random_state
    )
    model2.fit(X_train_m2_res, y_train_m2_res)

    y_pred_m2 = model2.predict(X_test_m2)
    y_proba_m2 = model2.predict_proba(X_test_m2)[:, 1]

    print("\n=== MODEL 2: 0 (diabetes) vs 1 (prediabetes) ===")
    print(classification_report(y_test_m2, y_pred_m2, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test_m2, y_pred_m2))
    try:
        print("ROC AUC:", roc_auc_score(y_test_m2, y_proba_m2))
    except ValueError:
        print("Nie dało się policzyć AUC.")

    # ===== EWALUACJA KASKADY NA X_test, y_test =====
    cascade = CascadeDiabetesModel(model1, model2)
    y_cascade = cascade.predict(X_test)

    print("\n=== KASKADA: finalne etykiety 0/1/2 vs y_test ===")
    print(classification_report(y_test, y_cascade, digits=3))
    print("Confusion matrix:")
    print(confusion_matrix(y_test, y_cascade))

    return cascade


def main():
    # 🔧 KONFIG – DOSTOSUJ DO SIEBIE
    csv_path = "C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/raw/diabetes_012_health_indicators_BRFSS2015.csv"   # ← zmień jeśli masz inną ścieżkę
    target_col = "Diabetes_012"                      # ← zmień jeśli etykieta nazywa się inaczej
    test_size = 0.2
    random_state = 42

    print(">>> Wczytuję dane z:", csv_path)
    df = pd.read_csv(csv_path)

    if target_col not in df.columns:
        raise ValueError(f"Nie znaleziono kolumny docelowej '{target_col}' w DataFrame. "
                         f"Dostępne kolumny: {list(df.columns)}")

    X = df.drop(columns=[target_col])
    y = df[target_col].values

    print("Rozmiar X:", X.shape, "| Rozmiar y:", y.shape)
    try:
        print("Rozkład klas w y:", np.bincount(y.astype(int)))
    except Exception:
        pass

    X_train, X_test, y_train, y_test = train_test_split(
        X.values,
        y,
        test_size=test_size,
        random_state=random_state,
        stratify=y
    )

    print(">>> Zaczynam trenowanie kaskadowego modelu...")
    _ = train_cascade_with_splits(X_train, X_test, y_train, y_test, random_state=random_state)
    print(">>> Trenowanie zakończone.")


if __name__ == "__main__":
    main()
