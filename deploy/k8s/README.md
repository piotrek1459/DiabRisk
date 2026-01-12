# Kubernetes Manifests

## Setting Up OAuth Credentials

The `secrets.yaml` file contains placeholder values. To deploy with working Google OAuth:

### Option 1: Environment Variables (Recommended)

```bash
# Set your OAuth credentials
export GOOGLE_CLIENT_ID='your-client-id.apps.googleusercontent.com'
export GOOGLE_CLIENT_SECRET='GOCSPX-your-client-secret'

# Run the install script (will use environment variables)
./scripts/install-local-k3d.sh
```

### Option 2: Local Secrets File (Not Committed)

Create a `secrets.local.yaml` file (gitignored):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: oauth-secret
type: Opaque
stringData:
  google-client-id: "your-actual-client-id"
  google-client-secret: "your-actual-secret"
  redirect-url: "http://localhost/auth/google/callback"
```

Then apply it manually:
```bash
kubectl apply -f deploy/k8s/secrets.local.yaml
```

## Getting Google OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Create a new OAuth 2.0 Client ID
3. Set authorized redirect URI: `http://localhost/auth/google/callback`
4. Copy the Client ID and Client Secret

## Security Note

**Never commit real credentials to git!** Always use:
- Environment variables
- `.env.local` files (gitignored)
- Kubernetes secrets created from external sources
- Secret management tools (Vault, SOPS, sealed-secrets)
