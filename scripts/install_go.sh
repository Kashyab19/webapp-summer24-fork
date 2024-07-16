#!/bin/bash
set -e

# Install Go
sudo yum install -y wget
wget https://dl.google.com/go/go1.15.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz

# Set up Go environment
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh

# Install Git
sudo yum install -y git
