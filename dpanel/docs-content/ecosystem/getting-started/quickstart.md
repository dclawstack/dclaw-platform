# Quickstart

This guide walks you through deploying your first DClaw app in under 5 minutes.

## Step 1: Access DPanel

Open DPanel in your browser:

```
https://panel.yourdomain.com
```

Or locally:

```
http://localhost:3000
```

## Step 2: Browse the App Store

The DPanel home screen shows all 65 DClaw apps in a grid. Apps marked **Live** are ready to install. Apps marked **Soon** are in development.

## Step 3: Install an App

1. Click on **DClaw Chat**
2. Click the **Install** button
3. The DClaw Operator will:
   - Create a namespace for the app
   - Deploy the frontend and backend pods
   - Provision a PostgreSQL database
   - Configure ingress and TLS
   - Apply network policies

## Step 4: Verify Deployment

```bash
kubectl get pods -n dclaw-chat
```

You should see:
- `dclaw-chat-frontend-*` (Next.js)
- `dclaw-chat-backend-*` (FastAPI)
- `dclaw-chat-db-*` (PostgreSQL)

## Step 5: Open the App

Click **Open App** in DPanel, or navigate directly to:

```
https://chat.yourdomain.com
```

## Step 6: Configure the App

Each app has a configuration page in DPanel:

1. Go to the app's detail page
2. Click the **Settings** tab
3. Adjust environment variables, resource limits, and scaling

## What's Next?

- Read the [Configuration](./configuration) guide for advanced setup
- Explore the [Architecture](../architecture) docs to understand how apps are provisioned
- Check [Troubleshooting](../troubleshooting) if you encounter issues
