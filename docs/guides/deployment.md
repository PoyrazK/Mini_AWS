# Deployment Guide

This guide covers how to deploy **The Cloud** in different environments.

## Local Development (Docker Compose)

The easiest way to run the entire stack locally is using Docker Compose.

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Update the `.env` file:**
   Make sure to set a secure `SECRETS_ENCRYPTION_KEY` (at least 16 characters).

3. **Start the services:**
   ```bash
   docker compose up -d
   ```

The API will be available at `http://localhost:8080`.

## Production (Docker Compose)

For a production-like setup on a single server, use the production profile:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

This applies production overrides like:
- Container restart policies (`always`).
- Log rotation.
- Production environment flags.

## Kubernetes

We provide basic Kubernetes manifests in the `k8s/` directory.

1. **Create the namespace:**
   ```bash
   kubectl apply -f k8s/namespace.yaml
   ```

2. **Apply Configs and Secrets:**
   ```bash
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/secrets.yaml
   ```
   *Note: Edit `k8s/secrets.yaml` and `k8s/configmap.yaml` with your actual production values before applying.*

3. **Deploy Database:**
   ```bash
   kubectl apply -f k8s/db-deployment.yaml
   ```

4. **Deploy API:**
   ```bash
   kubectl apply -f k8s/api-deployment.yaml
   ```

### Security Considerations

- **Secrets Manager**: The application requires a `SECRETS_ENCRYPTION_KEY`. Ensure this is stored securely using Kubernetes Secrets or a CSP-specific secrets manager.
- **Docker Socket**: The API requires access to `/var/run/docker.sock` to manage instances. In a production Kubernetes environment, you should consider using a remote Docker host or a Kubernetes-native instance manager (e.g., using CRDs).
- **Persistent Storage**: The examples use `emptyDir` for database storage. For production, use `PersistentVolumeClaims` backed by reliable storage.
