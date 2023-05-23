#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
TOPIC_NAME="<your-topic-name>"
CLOUD_RUN_SERVICE="<your-cloud-run-service>"
CLOUD_RUN_SERVICE_ACCOUNT="<your-cloud-run-service-account-email>"

# Create the Pub/Sub topic
gcloud pubsub topics create $TOPIC_NAME --project $PROJECT_ID

# Grant permission to the Cloud Run service to publish messages to the topic
gcloud run services add-iam-policy-binding $CLOUD_RUN_SERVICE \
  --member "serviceAccount:$CLOUD_RUN_SERVICE_ACCOUNT" \
  --role roles/pubsub.publisher \
  --project $PROJECT_ID

echo "Pub/Sub topic '$TOPIC_NAME' created successfully."
echo "Permission granted to the Cloud Run service '$CLOUD_RUN_SERVICE' to publish messages to the topic."
