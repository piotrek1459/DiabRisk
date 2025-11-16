# Prepared Diabetes Dataset -- File Descriptions

This package contains machine‑learning--ready data derived from the
original **diabetes_012_health_indicators_BRFSS2015.csv** dataset.\
Each file below serves a specific role in the training and evaluation
pipeline.

------------------------------------------------------------------------

## 1. X_train_processed.csv

**Purpose:** Training dataset (features only).\
Contains: - 80% of the data (stratified split)\
- All 21 input features\
- Scaled numeric columns: **BMI**, **MentHlth**, **PhysHlth**\
- Binary and categorical columns preserved as-is\
Ready for use in model training.

------------------------------------------------------------------------

## 2. X_test_processed.csv

**Purpose:** Test dataset (features only).\
Contains: - 20% of the data\
- Same preprocessing as the training set\
Used strictly for evaluating model performance on unseen data.

------------------------------------------------------------------------

## 3. y_train.csv

**Purpose:** Training labels corresponding to `X_train_processed.csv`.\
Target values: - **0** → No diabetes\
- **1** → Prediabetes\
- **2** → Diabetes

------------------------------------------------------------------------

## 4. y_test.csv

**Purpose:** Test labels corresponding to `X_test_processed.csv`.\
Used for accuracy, precision, recall, F1-score and confusion matrix
evaluation.

------------------------------------------------------------------------

## Summary

This structure ensures a clean, standardized training pipeline: - Data
is split correctly\
- Numeric features are scaled\
- Labels are preserved separately\
- Files are ready for any ML algorithm (RF, XGBoost, Logistic
Regression, NN, etc.)
