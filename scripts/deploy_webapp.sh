#!/bin/bash
set -e

# Deploy the web application
sudo mkdir -p /opt/${APPLICATION_NAME}
sudo mv /tmp/webapp /opt/${APPLICATION_NAME}/
sudo chmod +x /opt/${APPLICATION_NAME}/webapp
sudo chown -R csye6225:csye6225 /opt/${APPLICATION_NAME}
