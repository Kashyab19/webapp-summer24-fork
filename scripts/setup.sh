#!/bin/bash
set -e

# Create user and group
groupadd -f csye6225
useradd --system --home /sbin/nologin --shell /usr/sbin/nologin -g csye6225 csye6225

# Update and install necessary packages
sudo yum update -y
