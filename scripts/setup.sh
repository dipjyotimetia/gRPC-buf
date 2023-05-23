#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
SERVICE_ACCOUNT_NAME="<your-service-account-name>"

# Create the service account
gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
    --project $PROJECT_ID \
    --display-name "$SERVICE_ACCOUNT_NAME"


# Check if service account creation was successful
if [ $? -ne 0 ]; then
    echo "Error creating service account '$SERVICE_ACCOUNT_NAME'"
    exit 1
fi

# Define the roles as an array
roles=("roles/run.admin" "roles/run.serviceAgent" "roles/iam.serviceAccountUser" "roles/run.developer" "roles/storage.admin")

# Iterate over the roles array
for role in "${roles[@]}"
do
    echo "Adding role $role"

    # Add the role to the service account
    gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
        --role="$role"

    # Check if role addition was successful
    if [ $? -ne 0 ]; then
        echo "Error adding role $role to service account"
        exit 1
    fi
done

echo "Service account '$SERVICE_ACCOUNT_NAME' created and configured successfully."
