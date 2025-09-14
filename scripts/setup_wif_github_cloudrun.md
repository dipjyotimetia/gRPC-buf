# Setup Workload Identity Federation for GitHub Actions → Cloud Run

This script configures Google Cloud Workload Identity Federation (WIF) so GitHub Actions can deploy to Cloud Run without long‑lived JSON keys. It also prepares IAM roles, creates service accounts, optionally creates an Artifact Registry repository, and configures Secret Manager access for your Cloud Run runtime service account.

The script is idempotent and safe to re‑run.

## What it sets up

- Workload Identity Pool and OIDC Provider for GitHub Actions
- A deployer Service Account (used by GitHub Actions via WIF)
- A runtime Service Account (attached to your Cloud Run service)
- Required IAM roles for deploying to Cloud Run
- Optional: Artifact Registry repository and writer roles for the deployer SA
- Optional: Secret Manager accessor role for the runtime SA (per‑secret)
- Optional: Update an existing Cloud Run service to use the runtime SA

## Prerequisites

- gcloud CLI installed and authenticated (user or admin service account)
- The caller has permissions to:
  - Enable services (Service Usage)
  - Create and bind IAM roles
  - Create Workload Identity Pools/Providers
  - Create Service Accounts
  - Manage Artifact Registry (if enabled)
  - Manage Secret Manager bindings (if enabled)

The script enables these APIs if needed: IAM, IAM Credentials, Cloud Run, Secret Manager, and Artifact Registry.

## Required environment variables

- PROJECT_ID — e.g. `my-gcp-project`
- REGION — e.g. `us-central1` or a regional location you deploy to
- GITHUB_OWNER — GitHub org or user, e.g. `my-org` or `my-user`
- GITHUB_REPO — Repository name only, e.g. `my-repo`

## Optional variables (with defaults)

- SERVICE_NAME — Cloud Run service to update with the runtime SA (only if UPDATE_CLOUD_RUN=true)
- SECRET_IDS — Comma-separated Secret Manager IDs to grant to the runtime SA, e.g. `DATABASE_URL,JWT_SECRET`
- AR_REPO — Artifact Registry repo name (default: `apps`)
- AR_LOCATION — Artifact Registry location (default: `REGION`, e.g. `us` or `us-central1`)
- POOL_ID — WIF pool ID (default: `github-pool`)
- PROVIDER_ID — WIF provider ID (default: `github-provider`)
- DEPLOYER_SA_NAME — Deployer SA name (default: `github-deployer`)
- RUNTIME_SA_NAME — Runtime SA name (default: `cloudrun-runtime`)
- ALLOWED_REF — Git ref allowed to assume WIF (default: `refs/heads/main`)
- GRANT_AR_WRITER — Grant project‑level Artifact Registry writer to deployer SA (default: `true`)
- GRANT_AR_REPO_WRITER — Grant repo‑level AR writer to deployer SA (default: `true`)
- CREATE_AR_REPO — Create the AR repo if missing (default: `true`)
- PIPELINE_READ_SECRETS — Allow deployer SA to read secrets (default: `false`)
- UPDATE_CLOUD_RUN — Update the Cloud Run service’s service account (default: `false`)

Note on ALLOWED_REF: the provider condition uses equality. To change the allowed branch/tag, set `ALLOWED_REF` to an exact value like `refs/heads/main` or `refs/tags/v1.2.3`. To allow patterns, you would need to adjust the provider condition in GCP manually.

## Usage

1) Export variables (zsh/bash):

```zsh
export PROJECT_ID=dev-aileron-214211
export REGION=australia-southeast2
export GITHUB_OWNER=dipjyotimetia
export GITHUB_REPO=gRPC-buf

# Optional
export SERVICE_NAME=grpc-buf
export SECRET_IDS="DATABASE_URL,JWT_SECRET"
export AR_REPO=grpc-buf
export AR_LOCATION=australia-southeast2

# Toggles (optional, shown with defaults)
export ALLOWED_REF="refs/heads/main"
export GRANT_AR_WRITER=true
export GRANT_AR_REPO_WRITER=true
export CREATE_AR_REPO=true
export PIPELINE_READ_SECRETS=false
export UPDATE_CLOUD_RUN=false
```

2) Run the script:

```zsh
bash scripts/setup_wif_github_cloudrun.sh
```

## Expected output

At the end the script prints the key values to paste into your GitHub Actions workflow, for example:

- workload_identity_provider: projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github-pool/providers/github-provider
- service_account: github-deployer@PROJECT_ID.iam.gserviceaccount.com
- project_id, region, allowed_repo, allowed_ref

Keep these values handy for the workflow configuration.

## Example GitHub Actions job

Use the values printed by the script in place of the placeholders below.

```yaml
permissions:
  id-token: write
  contents: read

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - id: auth
        name: Authenticate to Google Cloud (WIF)
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: <FROM_SCRIPT_OUTPUT>
          service_account: <FROM_SCRIPT_OUTPUT>

      - name: Set up gcloud
        uses: google-github-actions/setup-gcloud@v2
        with:
          project_id: <YOUR_PROJECT_ID>

      - name: Configure Docker auth for Artifact Registry
        run: |
          gcloud auth configure-docker <AR_LOCATION>-docker.pkg.dev --quiet

      - name: Build and push image
        run: |
          IMAGE=<AR_LOCATION>-docker.pkg.dev/<PROJECT_ID>/<AR_REPO>/<REPO_NAME>:${GITHUB_SHA}
          docker build -t "$IMAGE" .
          docker push "$IMAGE"

      - name: Deploy to Cloud Run
        run: |
          IMAGE=<AR_LOCATION>-docker.pkg.dev/<PROJECT_ID>/<AR_REPO>/<REPO_NAME>:${GITHUB_SHA}
          gcloud run deploy <SERVICE_NAME> \
            --image="$IMAGE" \
            --region=<REGION> \
            --service-account=<RUNTIME_SA_EMAIL> \
            --allow-unauthenticated
```

Tips:
- `<AR_LOCATION>` is often a multi‑region like `us`, or a region like `us-central1`.
- `<REPO_NAME>` is typically your GitHub repo name.
- `<RUNTIME_SA_EMAIL>` will be `cloudrun-runtime@PROJECT_ID.iam.gserviceaccount.com` unless you changed `RUNTIME_SA_NAME`.

## What gets created/updated

- Workload Identity Pool: `github-pool` (by default)
- OIDC Provider: `github-provider` (issuer: `https://token.actions.githubusercontent.com`)
- Deployer SA: `github-deployer@PROJECT_ID.iam.gserviceaccount.com`
  - Roles: `roles/run.admin`, `roles/iam.serviceAccountUser`
  - Optional: `roles/artifactregistry.writer` (project and/or repo level)
  - Optional: `roles/secretmanager.secretAccessor` (if PIPELINE_READ_SECRETS=true)
  - WIF binding with a conditional constraint on `repository_owner` and `ref`
- Runtime SA: `cloudrun-runtime@PROJECT_ID.iam.gserviceaccount.com`
  - Granted `roles/secretmanager.secretAccessor` on specified `SECRET_IDS`
- Optional: Artifact Registry repo `${AR_REPO}` in `${AR_LOCATION}` (Docker format)
- Optional: Update Cloud Run service `${SERVICE_NAME}` to use the runtime SA

## Troubleshooting

- Permission denied / missing APIs: Ensure your caller identity has IAM admin/service usage admin permissions. The script enables required services but cannot escalate your identity.
- Pool or provider already exists: The script is idempotent and will update certain fields when possible.
- Artifact Registry errors: Verify `AR_LOCATION` is a valid AR location and that your project has Artifact Registry enabled.
- Secret binding failures: Ensure the secrets named in `SECRET_IDS` exist in the same project and that you have permission to bind IAM policies on them.
- Cloud Run update fails: Set `UPDATE_CLOUD_RUN=true` and verify `SERVICE_NAME` exists in the specified `REGION`.
- Allowing a different branch or tag: Re‑run with `ALLOWED_REF` set to the exact ref (e.g., `refs/heads/release`). For broader patterns, update the provider condition in GCP to a CEL expression like `startsWith(assertion.ref, 'refs/heads/')`.

## Clean up (optional)

These commands are destructive. Only run them if you want to remove what the script created. Replace placeholders with your actual values.

```zsh
# Remove WIF binding from deployer SA
gcloud iam service-accounts get-iam-policy github-deployer@${PROJECT_ID}.iam.gserviceaccount.com \
  --project=${PROJECT_ID} >/dev/null # inspect before changing

# Delete provider and pool (order matters)
gcloud iam workload-identity-pools providers delete github-provider \
  --workload-identity-pool=github-pool --location=global --project=${PROJECT_ID}

gcloud iam workload-identity-pools delete github-pool \
  --location=global --project=${PROJECT_ID}

# Optionally delete SAs (if not used elsewhere)
gcloud iam service-accounts delete github-deployer@${PROJECT_ID}.iam.gserviceaccount.com --project=${PROJECT_ID}
gcloud iam service-accounts delete cloudrun-runtime@${PROJECT_ID}.iam.gserviceaccount.com --project=${PROJECT_ID}

# Optionally delete the Artifact Registry repo
# WARNING: This deletes container images in the repo.
gcloud artifacts repositories delete <AR_REPO> \
  --location=<AR_LOCATION> --project=${PROJECT_ID}
```

## Notes

- This repository already includes the same usage example in the script header for quick copy‑paste.
- No GitHub secrets are needed for GCP authentication with WIF; the trust is established via OIDC and the conditional binding created by this script.
