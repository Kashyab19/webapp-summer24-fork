packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "project_id" {
  type    = string
  default = "your-id-goes-here"
}


variable "source_image_family" {
  type    = string
  default = "centos-stream-9"
}

variable "repository_url" {
  type    = string
  default = "https://github.com/Kashyab19/webapp-summer24-fork.git"
}

variable "branch" {
  type    = string
  default = "main"
}

variable "region" {
  type    = string
  default = "us-central1"
}

source "googlecompute" "centos" {
  project_id          = var.project_id
  source_image_family = var.source_image_family
  region              = var.region
  ssh_username        = "centos"
  zone                = "us-central1-a"
  credentials_file    = "./account-2.json"
}

build {
  sources = ["source.googlecompute.centos"]

  provisioner "file" {
    source      = "./account-2.json"
    destination = "/tmp/account.json"
  }

  provisioner "shell" {
    inline = [
      

      # Update and install necessary packages
      "sudo dnf -y update",
      "sudo dnf -y install golang postgresql-server postgresql-contrib git",

      # Set up PostgreSQL
      "sudo postgresql-setup --initdb",
      "sudo systemctl enable postgresql",
      "sudo systemctl start postgresql",

      # Create the csye6225 user
      "sudo groupadd -r csye6225",
      "sudo useradd -r -g csye6225 -s /usr/sbin/nologin csye6225",

      # Clone the repository
      "sudo git clone -b ${var.branch} ${var.repository_url} /opt/webapp",
      "sudo chown -R csye6225:csye6225 /opt/webapp",

      # Set up systemd service for your application
      "sudo cp /opt/webapp/config/webapp.service /etc/systemd/system/webapp.service",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable webapp.service",
    ]
  }
}
