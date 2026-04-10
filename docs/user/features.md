# Features

## User Authentication

- Email/password registration and login
- Secure session management with HttpOnly cookies
- Session expiration after 7 days

## Diabetes Risk Assessment

- 21 health indicators from the BRFSS 2015 survey
- Two-stage cascade Random Forest ML model
- Risk categories: Low, Medium, High
- Instant results with probability percentage

## Health Indicators

The assessment uses the following input parameters:

| Category | Indicators |
|----------|-----------|
| Health conditions | HighBP, HighChol, CholCheck, BMI |
| Lifestyle | Smoker, Stroke, HeartDiseaseorAttack, PhysActivity, Fruits, Veggies, HvyAlcoholConsump |
| Healthcare access | AnyHealthcare, NoDocbcCost |
| Health measures | GenHlth, MentHlth, PhysHlth, DiffWalk |
| Demographics | Sex, Age, Education, Income |
