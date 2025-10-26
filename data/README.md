# 🩺 Diabetes Health Indicators Dataset (Balanced 50/50 Split)

**Source:** Centers for Disease Control and Prevention (CDC) — Behavioral Risk Factor Surveillance System (BRFSS 2015)  
**Original dataset:** [CDC Diabetes Health Indicators (UCI Repository)](https://archive.ics.uci.edu/dataset/891/cdc+diabetes+health+indicators)  
**Balanced version (50/50):** [Kaggle - Diabetes Health Indicators Dataset](https://www.kaggle.com/datasets/alexteboul/diabetes-health-indicators-dataset)

---

## 📘 Overview
This dataset comes from the **2015 CDC Behavioral Risk Factor Surveillance System (BRFSS)**, an annual health survey conducted across the United States.

It contains **health and demographic indicators** used to predict whether an individual has **diabetes**.

The current version (`diabetes_binary_5050split_health_indicators_BRFSS2015.csv`) is **balanced** — it contains an equal number of records with and without diabetes (50% each).

---

## 📊 Dataset Summary
- **Number of records:** 70,692  
- **Number of features:** 22  
- **Target variable:** `Diabetes_binary`  
  - `0` → No diabetes  
  - `1` → Diabetes (diagnosed)

---

## 🧩 Feature Description

| Column | Description | Type / Values |
|--------|--------------|----------------|
| `Diabetes_binary` | Diabetes status (target variable) | 0 = No, 1 = Yes |
| `HighBP` | High blood pressure | 0 = No, 1 = Yes |
| `HighChol` | High cholesterol | 0 = No, 1 = Yes |
| `CholCheck` | Cholesterol check in the last 5 years | 0 = No, 1 = Yes |
| `BMI` | Body Mass Index | Numeric |
| `Smoker` | Smoked at least 100 cigarettes in lifetime | 0 = No, 1 = Yes |
| `Stroke` | Ever had a stroke | 0 = No, 1 = Yes |
| `HeartDiseaseorAttack` | Coronary heart disease or myocardial infarction | 0 = No, 1 = Yes |
| `PhysActivity` | Any physical activity in past 30 days | 0 = No, 1 = Yes |
| `Fruits` | Consumed fruit ≥ 1 time per day | 0 = No, 1 = Yes |
| `Veggies` | Consumed vegetables ≥ 1 time per day | 0 = No, 1 = Yes |
| `HvyAlcoholConsump` | Heavy alcohol consumption | 0 = No, 1 = Yes |
| `AnyHealthcare` | Have any kind of health care coverage | 0 = No, 1 = Yes |
| `NoDocbcCost` | Could not see a doctor due to cost | 0 = No, 1 = Yes |
| `GenHlth` | General health (1 = excellent → 5 = poor) | Integer [1–5] |
| `MentHlth` | Number of days mental health not good (0–30) | Integer |
| `PhysHlth` | Number of days physical health not good (0–30) | Integer |
| `DiffWalk` | Difficulty walking or climbing stairs | 0 = No, 1 = Yes |
| `Sex` | Gender | 0 = Female, 1 = Male |
| `Age` | Age category (1 = 18–24, 2 = 25–29, …, 13 = 80+) | Integer [1–13] |
| `Education` | Education level (1 = never attended → 6 = college graduate) | Integer [1–6] |
| `Income` | Annual income category (1 = <$10,000 → 8 = ≥$75,000) | Integer [1–8] |

---

## ⚙️ Data Notes
- The dataset is **balanced (50% diabetic / 50% non-diabetic)** — ideal for machine learning classification tasks.  
- All variables are **numerical or binary**, making preprocessing straightforward.  
- No missing values are present in this version.  
- Values are encoded as integers for easier model training.

---

## 🧠 Example Use Cases
This dataset can be used for:
- **Binary classification:** Predicting whether a person has diabetes.  
- **Feature importance analysis:** Determining which health indicators contribute most to diabetes risk.  
- **Health analytics:** Building explainable AI models for early diabetes detection.  
- **Educational projects:** Practicing data cleaning, visualization, and model training.

---

## 🧾 License
This dataset is derived from publicly available **CDC BRFSS 2015 data** and is distributed for **educational and research purposes only** under the U.S. Public Domain license.

---

## 📅 Version Info
- **Original survey year:** 2015  
- **Processed and balanced version:** 2021  
- **File name:** `diabetes_binary_5050split_health_indicators_BRFSS2015.csv`  
- **Maintained by:** CDC National Center for Chronic Disease Prevention and Health Promotion

---

## 🔗 References
- Centers for Disease Control and Prevention (CDC). [Behavioral Risk Factor Surveillance System (BRFSS)](https://www.cdc.gov/brfss/index.html)  
- UCI Machine Learning Repository: [CDC Diabetes Health Indicators Dataset](https://archive.ics.uci.edu/dataset/891/cdc+diabetes+health+indicators)  
- Kaggle Dataset: [Diabetes Health Indicators (Balanced 50/50 Split)](https://www.kaggle.com/datasets/alexteboul/diabetes-health-indicators-dataset)

---

*Prepared for educational and machine learning purposes.*
