#!/usr/bin/env bash

# Configure Workload Identity Federation (GitHub Actions → Google Cloud)
# and IAM for deploying to Cloud Run, with Secret Manager access.
#
# This script is idempotent and safe to re-run.
#
# Prereqs:
# - gcloud CLI installed and authenticated with a user or admin SA
# - Caller has permissions: IAM admin, Service Usage admin, and can create IAM SA
#
# Required env vars:
#   PROJECT_ID            # e.g. my-gcp-project
#   REGION                # e.g. us-central1
#   GITHUB_OWNER          # e.g. my-org or my-user
#   GITHUB_REPO           # e.g. my-repo (name only)
#
# Recommended (optional):
#   SERVICE_NAME          # Cloud Run service name to set runtime SA on (optional)
#   SECRET_IDS            # Comma-separated secret IDs to grant runtime SA access (optional)
#   AR_REPO               # Artifact Registry repo name to push images (default: apps)
#   AR_LOCATION           # Artifact Registry location (default: REGION). e.g. us or us-central1
#
# Tunables (have sensible defaults):
#   POOL_ID               # default: github-pool
#   PROVIDER_ID           # default: github-provider
#   DEPLOYER_SA_NAME      # default: github-deployer
#   RUNTIME_SA_NAME       # default: cloudrun-runtime
#   ALLOWED_REF           # default: refs/heads/main (branch restriction for WIF)
#   GRANT_AR_WRITER       # default: true (project-level AR writer to deployer SA)
#   GRANT_AR_REPO_WRITER  # default: true (repo-level AR writer to deployer SA)
#   CREATE_AR_REPO        # default: true (create AR repo if missing)
#   PIPELINE_READ_SECRETS # default: false (Secret Manager accessor to deployer SA)
#   UPDATE_CLOUD_RUN      # default: false (update Cloud Run service account)
#
# Usage:
export PROJECT_ID=
export REGION=australia-southeast2
export GITHUB_OWNER=dipjyotimetia
export GITHUB_REPO=grpc-buf
export SERVICE_NAME=grpc-buf
export SECRET_IDS="DATABASE_URL,JWT_SECRET"
export AR_REPO=grpc-buf
export AR_LOCATION=australia-southeast2
# bash scripts/setup_wif_github_cloudrun.sh

set -euo pipefail

echo "==> Validating environment..."
command -v gcloud >/dev/null 2>&1 || { echo "gcloud not found in PATH" >&2; exit 1; }

: "${PROJECT_ID:?PROJECT_ID is required}"
: "${REGION:?REGION is required}"
: "${GITHUB_OWNER:?GITHUB_OWNER is required}"
: "${GITHUB_REPO:?GITHUB_REPO is required}"

# Defaults
POOL_ID=${POOL_ID:-github-pool}
PROVIDER_ID=${PROVIDER_ID:-github-provider}
DEPLOYER_SA_NAME=${DEPLOYER_SA_NAME:-github-deployer}
RUNTIME_SA_NAME=${RUNTIME_SA_NAME:-cloudrun-runtime}
ALLOWED_REF=${ALLOWED_REF:-refs/heads/main}
GRANT_AR_WRITER=${GRANT_AR_WRITER:-true}
GRANT_AR_REPO_WRITER=${GRANT_AR_REPO_WRITER:-true}
CREATE_AR_REPO=${CREATE_AR_REPO:-true}
PIPELINE_READ_SECRETS=${PIPELINE_READ_SECRETS:-false}
UPDATE_CLOUD_RUN=${UPDATE_CLOUD_RUN:-false}
AR_REPO=${AR_REPO:-apps}
AR_LOCATION=${AR_LOCATION:-${REGION}}

GITHUB_REPO_FULL="${GITHUB_OWNER}/${GITHUB_REPO}"

echo "==> Getting project number for ${PROJECT_ID}..."
PROJECT_NUMBER=$(gcloud projects describe "${PROJECT_ID}" --format='value(projectNumber)')
if [[ -z "${PROJECT_NUMBER}" ]]; then
  echo "Failed to resolve project number for ${PROJECT_ID}" >&2
  exit 1
fi

POOL_FULL="projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}"
PROVIDER_FULL="${POOL_FULL}/providers/${PROVIDER_ID}"

DEPLOYER_SA_EMAIL="${DEPLOYER_SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
RUNTIME_SA_EMAIL="${RUNTIME_SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "==> Enabling required APIs (idempotent)..."
gcloud services enable \
  iam.googleapis.com \
  iamcredentials.googleapis.com \
  run.googleapis.com \
  secretmanager.googleapis.com \
  artifactregistry.googleapis.com \
  --project "${PROJECT_ID}" >/dev/null

echo "==> Creating Workload Identity Pool if missing..."
if ! gcloud iam workload-identity-pools describe "${POOL_ID}" --location=global --project="${PROJECT_ID}" >/dev/null 2>&1; then
  gcloud iam workload-identity-pools create "${POOL_ID}" \
    --location=global \
    --project="${PROJECT_ID}" \
    --display-name="GitHub Actions Pool"
else
  echo "    Pool exists: ${POOL_ID}"
fi

echo "==> Creating OIDC Provider for GitHub if missing..."
if ! gcloud iam workload-identity-pools providers describe "${PROVIDER_ID}" \
  --location=global \
  --workload-identity-pool="${POOL_ID}" \
  --project="${PROJECT_ID}" >/dev/null 2>&1; then
  gcloud iam workload-identity-pools providers create-oidc "${PROVIDER_ID}" \
    --location=global \
    --workload-identity-pool="${POOL_ID}" \
    --project="${PROJECT_ID}" \
    --display-name="GitHub Actions OIDC" \
    --issuer-uri="https://token.actions.githubusercontent.com" \
    --allowed-audiences="${PROVIDER_FULL}" \
    --attribute-condition="assertion.repository_owner=='${GITHUB_OWNER}' && assertion.ref=='refs/heads/main'" \
    --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner,attribute.ref=assertion.ref,attribute.actor=assertion.actor,attribute.aud=assertion.aud,attribute.job_workflow_ref=assertion.job_workflow_ref"
else
  echo "    Provider exists: ${PROVIDER_ID}"
fi

echo "==> Creating service accounts if missing..."
if ! gcloud iam service-accounts describe "${DEPLOYER_SA_EMAIL}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
  gcloud iam service-accounts create "${DEPLOYER_SA_NAME}" \
    --display-name="GitHub Deployer (WIF)" \
    --project="${PROJECT_ID}"
else
  echo "    Deployer SA exists: ${DEPLOYER_SA_EMAIL}"
fi

if ! gcloud iam service-accounts describe "${RUNTIME_SA_EMAIL}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
  gcloud iam service-accounts create "${RUNTIME_SA_NAME}" \
    --display-name="Cloud Run Runtime" \
    --project="${PROJECT_ID}"
else
  echo "    Runtime SA exists: ${RUNTIME_SA_EMAIL}"
fi

echo "==> Ensuring Artifact Registry repository exists (optional)..."
if [[ "${CREATE_AR_REPO}" == "true" ]]; then
  if ! gcloud artifacts repositories describe "${AR_REPO}" --location="${AR_LOCATION}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
    gcloud artifacts repositories create "${AR_REPO}" \
      --location="${AR_LOCATION}" \
      --repository-format=docker \
      --description="Container images for ${PROJECT_ID} (${AR_LOCATION})" \
      --project="${PROJECT_ID}"
  else
    echo "    Artifact Registry repo exists: ${AR_REPO} (${AR_LOCATION})"
  fi
else
  echo "    Skipping repo creation (CREATE_AR_REPO=false)"
fi

echo "==> Granting project roles to deployer SA (idempotent)..."
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${DEPLOYER_SA_EMAIL}" \
  --role="roles/run.admin" \
  --quiet >/dev/null

# Needed so the deployer can set the service account of the Cloud Run service
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${DEPLOYER_SA_EMAIL}" \
  --role="roles/iam.serviceAccountUser" \
  --quiet >/dev/null

if [[ "${GRANT_AR_WRITER}" == "true" ]]; then
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${DEPLOYER_SA_EMAIL}" \
    --role="roles/artifactregistry.writer" \
    --quiet >/dev/null
fi

if [[ "${GRANT_AR_REPO_WRITER}" == "true" ]]; then
  gcloud artifacts repositories add-iam-policy-binding "${AR_REPO}" \
    --location="${AR_LOCATION}" \
    --project="${PROJECT_ID}" \
    --member="serviceAccount:${DEPLOYER_SA_EMAIL}" \
    --role="roles/artifactregistry.writer" \
    --quiet >/dev/null || true
fi

if [[ "${PIPELINE_READ_SECRETS}" == "true" ]]; then
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${DEPLOYER_SA_EMAIL}" \
    --role="roles/secretmanager.secretAccessor" \
    --quiet >/dev/null
fi

echo "==> Granting Secret Manager access to runtime SA (optional, per-secret)..."
if [[ -n "${SECRET_IDS:-}" ]]; then
  IFS=',' read -r -a _secrets <<< "${SECRET_IDS}"
  for sid in "${_secrets[@]}"; do
    sid_trimmed=$(echo "$sid" | xargs)
    if [[ -n "${sid_trimmed}" ]]; then
      gcloud secrets add-iam-policy-binding "${sid_trimmed}" \
        --project="${PROJECT_ID}" \
        --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
        --role="roles/secretmanager.secretAccessor" \
        --quiet >/dev/null || true
    fi
  done
fi

echo "==> Binding WIF principal to deployer SA with branch restriction..."
WIF_MEMBER="principalSet://iam.googleapis.com/${POOL_FULL}/attribute.repository/${GITHUB_REPO_FULL}"
COND_TITLE="github-branch"
COND_DESC="Allow GitHub Actions from ${GITHUB_REPO_FULL} on ${ALLOWED_REF}"

# Expression requires escaped quotes
COND_EXPR="attribute.repository=='${GITHUB_REPO_FULL}' && attribute.ref=='${ALLOWED_REF}' && attribute.aud=='${PROVIDER_FULL}'"

set +e
gcloud iam service-accounts add-iam-policy-binding "${DEPLOYER_SA_EMAIL}" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="${WIF_MEMBER}" \
  --condition="expression=${COND_EXPR},title=${COND_TITLE},description=${COND_DESC}" >/dev/null 2>&1
RC=$?
set -e
if [[ $RC -ne 0 ]]; then
  echo "    Note: binding may already exist or conditions differ; attempting unconditional check..."
  # Try without condition if necessary (won't remove any existing bindings)
  gcloud iam service-accounts add-iam-policy-binding "${DEPLOYER_SA_EMAIL}" \
    --project="${PROJECT_ID}" \
    --role="roles/iam.workloadIdentityUser" \
    --member="${WIF_MEMBER}" >/dev/null || true
fi

if [[ "${UPDATE_CLOUD_RUN}" == "true" && -n "${SERVICE_NAME:-}" ]]; then
  echo "==> Updating Cloud Run service to use runtime SA..."
  set +e
  gcloud run services update "${SERVICE_NAME}" \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --service-account="${RUNTIME_SA_EMAIL}" \
    --quiet >/dev/null
  if [[ $? -ne 0 ]]; then
    echo "    Could not update service '${SERVICE_NAME}'. Ensure it exists, or deploy first." >&2
  fi
  set -e
fi

cat <<EOF

Completed.

Key values for your GitHub Actions workflow:

  workload_identity_provider: ${PROVIDER_FULL}
  service_account:            ${DEPLOYER_SA_EMAIL}
  project_id:                 ${PROJECT_ID}
  region:                     ${REGION}
  allowed_repo:               ${GITHUB_REPO_FULL}
  allowed_ref:                ${ALLOWED_REF}

EOF

# Keep the example literal to avoid expanding $GITHUB_SHA, etc.
cat <<'EOF'
Example GitHub Actions job:

  permissions:
    id-token: write   # to request OIDC token
    contents: read    # to checkout code

  steps:
    - uses: actions/checkout@v4

    - id: auth
      name: Authenticate to Google Cloud (WIF)
      uses: google-github-actions/auth@v2
      with:
        workload_identity_provider: ${PROVIDER_FULL}
        service_account: ${DEPLOYER_SA_EMAIL}

    - name: Set up gcloud
      uses: google-github-actions/setup-gcloud@v2
      with:
        project_id: ${PROJECT_ID}

    # Optional: if you push images to Artifact Registry
    - name: Configure Docker auth for Artifact Registry
      run: |
        gcloud auth configure-docker ${AR_LOCATION}-docker.pkg.dev --quiet

    # Build and push your image (example)
    - name: Build and push image
      run: |
        IMAGE=${AR_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${AR_REPO}/${GITHUB_REPO}:$GITHUB_SHA
        docker build -t "$IMAGE" .
        docker push "$IMAGE"

    # Deploy to Cloud Run
    - name: Deploy to Cloud Run
      run: |
        IMAGE=${AR_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${AR_REPO}/${GITHUB_REPO}:$GITHUB_SHA
        gcloud run deploy ${SERVICE_NAME:-my-service} \
          --image="$IMAGE" \
          --region=${REGION} \
          --service-account=${RUNTIME_SA_EMAIL} \
          --allow-unauthenticated

Notes:
- Runtime SA (${RUNTIME_SA_EMAIL}) has Secret Manager access to: ${SECRET_IDS:-<none configured>}
- Pipeline SA (${DEPLOYER_SA_EMAIL}) Artifact Registry writer (project): ${GRANT_AR_WRITER}
- Pipeline SA (${DEPLOYER_SA_EMAIL}) Artifact Registry writer (repo ${AR_REPO}@${AR_LOCATION}): ${GRANT_AR_REPO_WRITER}
- Pipeline SA Secret Manager accessor: ${PIPELINE_READ_SECRETS}
- WIF restricted to repo ${GITHUB_REPO_FULL} and ref ${ALLOWED_REF}

EOF