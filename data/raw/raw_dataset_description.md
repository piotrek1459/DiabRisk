
# BRFSS 2015 – Raw Diabetes Dataset Description  
**File:** `diabetes_012_health_indicators_BRFSS2015.csv`

This Markdown document provides a clear and concise description of the **original, unprocessed dataset** before any machine‑learning preprocessing was applied.

---

# 📂 1. Dataset Overview
- **Source:** Behavioral Risk Factor Surveillance System (BRFSS), 2015  
- **Rows:** 253,680  
- **Columns:** 22  
- **Task Type:** Multiclass classification  
- **Target variable:** `Diabetes_012`  
- Values:
  - **0** → No diabetes  
  - **1** → Prediabetes  
  - **2** → Diabetes  

The dataset contains health indicators, lifestyle factors, chronic conditions, and demographic information.

---

# 📊 2. Column-by-Column Description

## 🎯 Target Variable
### **Diabetes_012**
- Classification label  
- Values:
  - 0 – No diabetes  
  - 1 – Prediabetes  
  - 2 – Diabetes  

---

# 🩺 3. Health Condition Indicators

### **HighBP**
- High blood pressure  
- 1 = Yes, 0 = No

### **HighChol**
- High cholesterol  
- 1 = Yes, 0 = No

### **CholCheck**
- Cholesterol checked within the last 5 years  
- 1 = Yes, 0 = No

### **BMI**
- Body Mass Index (raw numeric value)

---

# 🚬 4. Lifestyle & Behavioral Indicators

### **Smoker**
- Smoked at least 100 cigarettes in lifetime  
- 1 = Yes, 0 = No

### **Stroke**
- Ever diagnosed with a stroke  
- 1 = Yes, 0 = No

### **HeartDiseaseorAttack**
- Coronary heart disease or myocardial infarction  
- 1 = Yes, 0 = No

### **PhysActivity**
- Physical activity within the last 30 days  
- 1 = Yes, 0 = No

### **Fruits**
- Eats fruit at least once per day  
- 1 = Yes, 0 = No

### **Veggies**
- Eats vegetables at least once per day  
- 1 = Yes, 0 = No

### **HvyAlcoholConsump**
- Heavy alcohol consumption  
- 1 = Yes, 0 = No

---

# 🧠 5. Mental & Physical Health Measures

### **GenHlth**
- General health  
- Categorical scale:
  - 1 = Excellent  
  - 5 = Poor  

### **MentHlth**
- Number of days in last 30 days with bad mental health  
- Range: 0–30

### **PhysHlth**
- Number of days in last 30 days with bad physical health  
- Range: 0–30

### **DiffWalk**
- Difficulty walking  
- 1 = Yes, 0 = No

---

# 🏥 6. Healthcare Access

### **AnyHealthcare**
- Has any form of healthcare coverage  
- 1 = Yes, 0 = No

### **NoDocbcCost**
- Could not see a doctor due to cost  
- 1 = Yes, 0 = No

---

# 👤 7. Demographic Information

### **Sex**
- 0 = Female  
- 1 = Male  

### **Age**
- Categorical age groups (1–13)  
- Example: 1 = 18–24, 13 = 80+

### **Education**
- Educational level (1–6)  
- 1 = Never attended school  
- 6 = Some post‑secondary education

### **Income**
- Income category (1–8)  
- Higher number = higher income bracket

---

# 📘 Summary
The raw dataset provides:
- A wide range of **behavioral, physiological, and demographic indicators**
- Multiclass labels allowing detection of **prediabetes and diabetes**
- Clean, numeric data with no missing values
- A strong basis for **classification, risk modeling, and health analytics**

This document describes the dataset *as-is*, before any preprocessing such as scaling, splitting, or normalization.

