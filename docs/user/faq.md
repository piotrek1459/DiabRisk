# FAQ

## What data is the model based on?

The project uses the BRFSS 2015 health-indicators dataset stored in
`data/raw/`. The current repository also includes precomputed processed
splits in `data/processed/`.

## Do I need an account to use the assessment?

Yes. In the current application flow, the risk form is displayed only after
logging in.

## What kind of account does the current version use?

The current browser flow uses local email/password authentication and a
session cookie.

## What does the application show after I submit the form?

The application shows:

- a risk percentage
- a category: `Low`, `Medium`, or `High`
- a short explanatory message

## What do the risk categories mean in the current UI?

- **Low**: up to 50%
- **Medium**: above 50% up to 80%
- **High**: above 80%

These categories are based on the current ML API response mapping.

## Is my data stored?

The current application stores account data and authentication-session
data. The user-facing browser flow does not provide a visible history of
submitted assessments.

## Can I export or download my results?

No. The current UI does not expose export or report-download features.

## Is this medical advice?

No. DiabRisk is an educational screening tool. It does not replace medical
consultation, diagnosis, or treatment.
