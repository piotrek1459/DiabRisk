import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report
from imblearn.over_sampling import SMOTE
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix

X_train = pd.read_csv("C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/processed/X_train_processed.csv")
X_test = pd.read_csv("C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/processed/X_test_processed.csv")

y_train = pd.read_csv("C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/processed/y_train.csv")
y_test = pd.read_csv("C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/processed/y_test.csv")

y_train = y_train["Diabetes_012"]
y_test = y_test["Diabetes_012"]

print("Rozkład klas PRZED SMOTE:")
print(y_train.value_counts())


smote = SMOTE(sampling_strategy="not majority", random_state=42)
X_train_res, y_train_res = smote.fit_resample(X_train, y_train)

print("\nRozkład klas PO SMOTE:")
print(y_train_res.value_counts())


model = RandomForestClassifier(
    n_estimators=200,
    n_jobs=-1,
    random_state=42
)

model.fit(X_train_res, y_train_res)


y_pred = model.predict(X_test)
print("\n=== Classification report (po SMOTE) ===")
print(classification_report(y_test, y_pred))
print(confusion_matrix(y_test, y_pred))
