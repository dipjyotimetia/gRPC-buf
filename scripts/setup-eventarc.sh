#!/bin/bash

# Set the variables
PROJECT_ID="<your-project-id>"
TOPIC_NAME="<your-pubsub-topic-name>"
CLOUD_RUN_SERVICE="<your-cloud-run-service>"
REGION="<your-preferred-region>"

# Create a Pub/Sub topic
gcloud pubsub topics create $TOPIC_NAME --project $PROJECT_ID

# Deploy the Cloud Run service
gcloud run deploy $CLOUD_RUN_SERVICE \
  --image gcr.io/$PROJECT_ID/$CLOUD_RUN_SERVICE \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated

# Set up Eventarc trigger for Cloud Run service
gcloud eventarc triggers create $CLOUD_RUN_SERVICE-trigger \
  --destination-run-service=$CLOUD_RUN_SERVICE \
  --destination-region=$REGION \
  --event-filters="type=google.cloud.pubsub.topic.v1.messagePublished" \
  --target-service-account=$(gcloud run services describe $CLOUD_RUN_SERVICE \
    --region $REGION \
    --format='value(status.assignedTraffic[0].url)' \
    --project $PROJECT_ID)

echo "Eventarc set up successfully with Cloud Run service '$CLOUD_RUN_SERVICE' and Pub/Sub topic '$TOPIC_NAME'."
