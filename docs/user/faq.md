# FAQ

## What data does DiabRisk use?

The ML model is trained on the CDC BRFSS 2015 dataset (Behavioral Risk Factor Surveillance System) containing 253,680 survey responses.

## How accurate is the prediction?

The screening model achieves ~94% accuracy and the severity model ~98% accuracy on the test set. However, this is a screening tool and should not replace professional medical diagnosis.

## Is my data stored?

Assessment results are stored in the database and linked to your user account. You can contact an administrator to request data deletion.

## What do the risk categories mean?

- **Low** (0-50%): Lower probability of diabetes risk based on provided indicators
- **Medium** (50-80%): Moderate probability; consider consulting a healthcare provider
- **High** (80-100%): Higher probability; strongly recommended to consult a healthcare provider
