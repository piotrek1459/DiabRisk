# Frontend

The frontend is a Svelte 5 single-page application built with Vite and
served in production from Nginx.

## Current Responsibilities

- render the login and registration views
- restore the current session through `GET /auth/session`
- collect the 21-feature diabetes-risk form
- submit the form to `POST /api/risk`
- display `RiskPercent`, `Category`, and `Message`

## Main Files

| Path | Role |
|------|------|
| `src/App.svelte` | top-level application flow, form state, result rendering |
| `src/lib/Login.svelte` | login form |
| `src/lib/Register.svelte` | registration form |
| `vite.config.js` | Vite config and local dev proxy |

## Local Development

```bash
cd frontend
npm install
npm run dev
```

The Vite dev server runs on `http://localhost:5173`.

## API Proxy

During local frontend development, Vite proxies:

- `/auth` -> `http://localhost`
- `/api` -> `http://localhost`

This expects the backend stack to be available through the k3d ingress on
`http://localhost`.

## Production Build

```bash
cd frontend
npm install
npm run build
```

`Dockerfile.frontend` builds the static assets and copies `dist/` into an
Nginx image that serves the compiled app on port `80`.
