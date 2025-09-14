#!/usr/bin/env bash
set -euo pipefail

# Local Cloud Run deploy helper for gRPC-buf
# Mirrors .github/workflows/gcpdeploy.yaml flags
#
# Usage:
#   export GCP_PROJECT_ID="your-project-id"
#   # Optional overrides
#   # export SERVICE="grpc-buf"
#   # export REGION="australia-southeast2"
#   # export RUNTIME_SA="runtime-sa@${GCP_PROJECT_ID}.iam.gserviceaccount.com"
#   # export DB_SECRET_NAME="DATABASE_URL"
#   # export JWT_SECRET_NAME="JWT_SECRET"
#   # export SECRET_PROJECT_NUMBER="<if secrets live in another project>"
#   ./scripts/deploy_local.sh

PROJECT_ID=${GCP_PROJECT_ID:?"GCP_PROJECT_ID is required (gcloud config project)"}
SERVICE=${SERVICE:-grpc-buf}
REGION=${REGION:-australia-southeast2}
DB_SECRET_NAME=${DB_SECRET_NAME:-DATABASE_URL}
JWT_SECRET_NAME=${JWT_SECRET_NAME:-JWT_SECRET}
# If secrets are in another project, set SECRET_PROJECT_NUMBER
SECRET_PROJECT_NUMBER=${SECRET_PROJECT_NUMBER:-}

# Compose --set-secrets value
if [[ -n "${SECRET_PROJECT_NUMBER}" ]]; then
  SECRETS_ARG="${DB_SECRET_NAME}=projects/${SECRET_PROJECT_NUMBER}/secrets/${DB_SECRET_NAME}:latest,${JWT_SECRET_NAME}=projects/${SECRET_PROJECT_NUMBER}/secrets/${JWT_SECRET_NAME}:latest"
else
  SECRETS_ARG="${DB_SECRET_NAME}=${DB_SECRET_NAME}:latest,${JWT_SECRET_NAME}=${JWT_SECRET_NAME}:latest"
fi

# Basic sanity check: ensure DATABASE_URL secret won't use the docker-compose hostname "postgres"
# Cloud Run cannot resolve the bare hostname "postgres"; it must be a reachable host (Cloud SQL, public IP, DNS name)
echo "[deploy] Validating secrets..."
DB_URL=""
if DB_URL=$(gcloud secrets versions access latest --secret "${DB_SECRET_NAME}" --project "${PROJECT_ID}" 2>/dev/null); then
  if [[ -z "${DB_URL}" ]]; then
    echo "[deploy] WARNING: Secret ${DB_SECRET_NAME} has empty value. Deployment will likely fail."
  else
    # Match common DSN forms that include the docker-compose hostname 'postgres'
    if [[ ${DB_URL} =~ @postgres(:|/|\?|$) || ${DB_URL} =~ (host=|HOST=)postgres([^a-zA-Z0-9_\-\.]|$) ]]; then
      echo "[deploy] ERROR: ${DB_SECRET_NAME} points at host 'postgres', which is only valid inside docker-compose."
      echo "        Update the secret to a reachable host for Cloud Run (e.g., Cloud SQL or a DNS/IP), then rerun this script."
      echo "        Example: postgres://USER:PASS@<CLOUD_SQL_IP_OR_DNS>:5432/DB?sslmode=require"
      exit 1
    fi
  fi
else
  echo "[deploy] NOTE: Could not read secret ${DB_SECRET_NAME} (it may not exist or you lack access)."
  echo "       Ensure the secret exists and is accessible to the deploying identity and the Cloud Run service."
fi

echo "[deploy] Ensuring required APIs are enabled..."
gcloud services enable \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  --project "${PROJECT_ID}" 1>/dev/null

echo "[deploy] Using project: ${PROJECT_ID}, service: ${SERVICE}, region: ${REGION}"

deploy_flags=(
  --project "${PROJECT_ID}"
  --region "${REGION}"
  --source .
  --allow-unauthenticated
  --use-http2
  --concurrency=80
  --cpu=1
  --memory=256Mi
  --max-instances=1
  --cpu-throttling
  --set-env-vars "ENVIRONMENT=prod"
  --set-secrets "${SECRETS_ARG}"
)

# Optional runtime service account
if [[ -n "${RUNTIME_SA:-}" ]]; then
  deploy_flags+=(--service-account "${RUNTIME_SA}")
fi

# Run deploy (Cloud Build will build from the Dockerfile)
echo "[deploy] Deploying to Cloud Run..."
gcloud run deploy "${SERVICE}" "${deploy_flags[@]}"

echo "[deploy] Done. Fetching service URL..."
gcloud run services describe "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --format='value(status.url)'
