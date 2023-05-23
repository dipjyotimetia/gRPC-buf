#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
SERVICE_ACCOUNT_NAME="<your-service-account-name>"
TOPIC_NAME="<your-pubsub-topic-name>"
CLOUD_RUN_SERVICE="<your-cloud-run-service>"
REGION="<your-preferred-region>"

# Create a Pub/Sub topic
gcloud pubsub topics create $TOPIC_NAME --project $PROJECT_ID

# Check if topic creation was successful
if [ $? -ne 0 ]; then
    echo "Error creating Pub/Sub topic '$TOPIC_NAME'"
    exit 1
fi

# Add necessary IAM permissions to the service account
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/eventarc.admin"

# Check if IAM policy binding was successful
if [ $? -ne 0 ]; then
    echo "Error adding IAM policy binding to the service account"
    exit 1
fi

# Set up Eventarc trigger for Cloud Run service
gcloud eventarc triggers create $CLOUD_RUN_SERVICE-trigger \
    --destination-run-service=$CLOUD_RUN_SERVICE \
    --destination-region=$REGION \
    --event-filters="type=google.cloud.pubsub.topic.v1.messagePublished" \
    --service-account=$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com

# Check if Eventarc trigger setup was successful
if [ $? -ne 0 ]; then
    echo "Error setting up Eventarc trigger for Cloud Run service '$CLOUD_RUN_SERVICE'"
    exit 1
fi

echo "Eventarc set up successfully with Cloud Run service '$CLOUD_RUN_SERVICE' and Pub/Sub topic '$TOPIC_NAME'."
