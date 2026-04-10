# Getting Started

## Prerequisites

### To use an already running instance

- A modern web browser such as Chrome, Firefox, Edge, or Safari

### To run the application locally

- Docker
- k3d
- kubectl

## Starting the Application Locally

### Windows (PowerShell)

```powershell
.\scripts\install-local-k3d.ps1
```

### Linux / macOS

```bash
./scripts/install-local-k3d.sh
```

After the setup finishes, open:

```text
http://localhost
```

### Optional Demo Credentials

For local deployments, the setup script creates a default admin account:

- Email: `admin@diabrisk.local`
- Password: `default_admin_password`

## Accessing the Application

1. Open `http://localhost` in your browser
2. Create an account using the registration form, or log in with existing credentials
3. After logging in, wait until the health assessment form is displayed

## Using the Risk Assessment

1. After logging in, fill out the health information form with your data
2. The form includes 21 health indicators covering health conditions, lifestyle, healthcare access, general health, and demographic data
3. Click **Estimate Risk** to get your diabetes risk assessment
4. The result shows:
   - a risk percentage
   - a risk category: Low, Medium, or High
   - a short explanatory message

## What the Application Needs From You

- An email address and password to create an account
- Values for all 21 health input fields
- A browser session with cookies enabled, because the app uses a session cookie after login

## Important Notes

- DiabRisk is an educational screening tool
- It does not provide a medical diagnosis
- The result should be treated as an informational estimate, not a clinical decision
- If you are concerned about your health, consult a qualified healthcare professional
