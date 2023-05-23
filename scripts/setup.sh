#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
SERVICE_ACCOUNT_NAME="<your-service-account-name>"
ROLE="<your-role>"
PERMISSIONS="<your-permissions>"

# Create the service account
gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
    --project $PROJECT_ID \
    --display-name "$SERVICE_ACCOUNT_NAME"

# Add the role to the service account
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="$ROLE"

# Grant permissions to the service account
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="$PERMISSIONS"

echo "Service account '$SERVICE_ACCOUNT_NAME' created and configured successfully."
