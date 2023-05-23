#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
SECRET_NAME="<your-secret-name>"
SECRET_VALUE="<your-secret-value>"
SERVICE_ACCOUNT_EMAIL="<your-service-account-email>"

# Create the secret
gcloud secrets create $SECRET_NAME --project $PROJECT_ID

# Check if secret creation was successful
if [ $? -ne 0 ]; then
    echo "Error creating secret '$SECRET_NAME'"
    exit 1
fi

# Add the secret payload
echo -n $SECRET_VALUE | gcloud secrets versions add $SECRET_NAME --data-file=-

# Check if secret payload addition was successful
if [ $? -ne 0 ]; then
    echo "Error adding secret payload to '$SECRET_NAME'"
    exit 1
fi

# Grant access to the secret
gcloud secrets add-iam-policy-binding $SECRET_NAME \
    --member=serviceAccount:$SERVICE_ACCOUNT_EMAIL \
    --role=roles/secretmanager.secretAccessor \
    --project $PROJECT_ID

# Check if IAM policy binding was successful
if [ $? -ne 0 ]; then
    echo "Error adding IAM policy binding to secret '$SECRET_NAME'"
    exit 1
fi

echo "Secret '$SECRET_NAME' created and configured successfully."
