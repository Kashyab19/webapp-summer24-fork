#!/bin/bash
set -e

# Debugging: Print the current directory
echo "Current Directory: $(pwd)"

# Debugging: List the contents of the parent directory
echo "Contents of Parent Directory:"
ls -la ..

# Load environment variables from .env file in the parent directory
if [ -f .env ]; then
  echo ".env file found"
  set -o allexport
  source .env
  set +o allexport
else
  echo ".env file not found in the root directory"
  exit 1
fi

# Debugging: Print loaded environment variables
echo "Loaded Environment Variables:"
echo "GCP_PROJECT_ID=$GCP_PROJECT_ID"
echo "GOOGLE_APPLICATION_CREDENTIALS=$GOOGLE_APPLICATION_CREDENTIALS"

# Validate required variables are set
if [ -z "$GCP_PROJECT_ID" ] || [ -z "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
  echo "Required environment variables are missing"
  exit 1
fi

# Create a temporary credentials file from the contents
echo "$GOOGLE_APPLICATION_CREDENTIALS" > /tmp/packer-service-account.json

# Run Packer
packer validate -var "gcp_project_id=$GCP_PROJECT_ID" -var "account_file=/tmp/packer-service-account.json" ./template.pkr.hcl
packer build -var "gcp_project_id=$GCP_PROJECT_ID" -var "account_file=/tmp/packer-service-account.json" ./template.pkr.hcl

# Clean up the temporary credentials file
rm /tmp/packer-service-account.json
