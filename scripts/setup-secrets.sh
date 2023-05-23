#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
SECRET_NAME="<your-secret-name>"
SECRET_VALUE="<your-secret-value>"

# Create the secret
gcloud secrets create $SECRET_NAME --project $PROJECT_ID

# Add the secret payload
echo -n $SECRET_VALUE | gcloud secrets versions add $SECRET_NAME --data-file=-

# Grant access to the secret
gcloud secrets add-iam-policy-binding $SECRET_NAME \
  --member=serviceAccount:<your-service-account-email> \
  --role=roles/secretmanager.secretAccessor \
  --project $PROJECT_ID

echo "Secret '$SECRET_NAME' created and configured successfully."
