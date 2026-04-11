# Features

## Authentication

- email/password registration
- email/password login
- logout through the current browser session
- session restoration with an HttpOnly cookie

## Diabetes Risk Assessment

- protected form available after login
- 21 input indicators
- immediate result returned after submission
- output shown directly in the browser

## Health Indicators

The current form uses these inputs:

| Category | Indicators |
|----------|-----------|
| Health conditions | `HighBP`, `HighChol`, `CholCheck`, `BMI` |
| Lifestyle | `Smoker`, `Stroke`, `HeartDiseaseorAttack`, `PhysActivity`, `Fruits`, `Veggies`, `HvyAlcoholConsump` |
| Healthcare access | `AnyHealthcare`, `NoDocbcCost` |
| Health measures | `GenHlth`, `MentHlth`, `PhysHlth`, `DiffWalk` |
| Demographics | `Sex`, `Age`, `Education`, `Income` |

## What the Result Contains

After submitting the form, the application displays:

- **Risk percentage** shown in the UI
- **Risk category**: `Low`, `Medium`, or `High`
- **Short explanation** returned by the ML service

## What Is Available in the Current UI

- registration and login
- session-based logout
- the risk-assessment form
- immediate display of the latest result

## What Is Not Available in the Current UI

- assessment history
- report download
- CSV or JSON export
- account deletion screen
- feature-level explanation charts
