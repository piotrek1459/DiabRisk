import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report
from imblearn.over_sampling import SMOTE
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix
import numpy as np

# 1. Wczytanie danych
df = pd.read_csv("C:/Users/Jeremi/Desktop/se_project/DiabRisk/data/raw/diabetes_012_health_indicators_BRFSS2015.csv")

x=df["Diabetes_012"].value_counts()




print(x)

import pandas as pd

# Załóżmy że Twój dataframe to df
# Najpierw znajdź indeksy wierszy z Diabetes_012 == 0.0
# Usuń najpierw z klasy 0.0
# Najpierw wybierz rekordy do usunięcia z każdej klasy
idx0 = df[df['Diabetes_012'] == 0.0].sample(n=206000, random_state=42).index
idx2 = df[df['Diabetes_012'] == 2.0].sample(n=30000, random_state=42).index

# Połącz i WYMIESZAJ losowo
all_indices = np.concatenate([idx0, idx2])
np.random.shuffle(all_indices)

# Usuń wszystkie na raz – kolejność losowa
df_balanced = df.drop(index=all_indices).reset_index(drop=True)

print("\n📊 Liczność klas po usunięciu:")
print(df_balanced['Diabetes_012'].value_counts())


# (opcjonalnie) sprawdź ile zostało wierszy dla 0.0
print(df_balanced['Diabetes_012'].value_counts())

# --- 2. Podział na X i y ---
X = df_balanced.drop(columns=['Diabetes_012'])
y = df_balanced['Diabetes_012']

# --- 3. Train/Test split ---
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42, stratify=y
)

# --- 4. Trenowanie Random Forest ---
model = RandomForestClassifier(
    n_estimators=200,      # więcej drzew
    max_depth=None,        # pełne drzewo
    class_weight='balanced',  # dla pewności jeszcze dopasowujemy wagami
    random_state=42,
    n_jobs=-1              # pełne wykorzystanie CPU
)
model.fit(X_train, y_train)

# --- 5. Predykcja i metryki ---
y_pred = model.predict(X_test)

print(f"\n🎯 Accuracy: {accuracy_score(y_test, y_pred):.4f}")
print("\n=== Classification Report ===")
print(classification_report(y_test, y_pred))

print("\n=== Confusion Matrix ===")
print(confusion_matrix(y_test, y_pred))
