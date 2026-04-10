# Development Setup

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [k3d](https://k3d.io/) (lightweight Kubernetes in Docker)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Quick Start

### Linux / macOS

```bash
./scripts/install-local-k3d.sh
```

### Windows (PowerShell)

```powershell
.\scripts\install-local-k3d.ps1
```

The script will:
1. Create a k3d cluster
2. Build all Docker images
3. Deploy all services to Kubernetes
4. Run database migrations
5. Print access URLs and admin credentials

## Default Credentials

| | Value |
|-|-------|
| Admin email | `admin@diabrisk.local` |
| Admin password | `default_admin_password` |

Override via environment variables `ADMIN_EMAIL` and `ADMIN_PASSWORD` before running the setup script.

## Frontend Development

```bash
cd frontend
npm install
npm run dev   # runs on http://localhost:5173
```

Vite proxies `/auth` and `/api` requests to the backend.

## Environment Variables

See `.env.example` for all available configuration options.
