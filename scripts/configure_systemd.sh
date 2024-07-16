#!/bin/bash
set -e

# Move and enable systemd service
sudo mv /tmp/${SERVICE_NAME} /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable ${SERVICE_NAME}
