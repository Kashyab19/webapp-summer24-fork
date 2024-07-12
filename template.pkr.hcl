#Not used
packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "gcp_project_id" {
  type    = string
  default = "numeric-cinema-428517-p4"
}

variable "source_image_family" {
  type    = string
  default = "centos-stream-9"
}

variable "machine_type" {
  type    = string
  default = "e2-medium"
}

variable "application_name" {
  type    = string
  default = "webapp"
}

variable "service_name" {
  type    = string
  default = "webapp.service"
}

variable "account_file" {
  type    = string
  default = "./packer-sa-3.json"
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

source "googlecompute" "example" {
  project_id           = var.gcp_project_id
  source_image_family  = var.source_image_family
  zone                 = "us-central1-a"
  disk_size            = 30
  image_name           = "web-app-img-${local.timestamp}"
  image_family         = "web-app-family"
  machine_type         = var.machine_type
  ssh_username         = "centos"
  wait_to_add_ssh_keys = "20s"
  credentials_file     = var.account_file
  network              = "default"
  metadata = {
    "enable-oslogin" = "FALSE"
  }
  ssh_agent_auth = false
}

build {
  sources = [
    "source.googlecompute.example"
  ]

  provisioner "shell" {
    inline = [
      "sudo dnf update -y || exit 1",
      "sudo dnf install -y postgresql-server postgresql-contrib || exit 1",
      "sudo postgresql-setup --initdb || exit 1",
      "sudo systemctl enable postgresql || exit 1",
      "sudo systemctl start postgresql || exit 1",
      "sudo -u postgres psql -c \"CREATE USER csye6225 WITH ENCRYPTED PASSWORD 'csye6225';\" || exit 1",
      "sudo -u postgres psql -c \"GRANT ALL PRIVILEGES ON DATABASE postgres TO csye6225;\" || exit 1",
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo groupadd -f csye6225 || exit 1",
      "sudo useradd --system --home /sbin/nologin --shell /usr/sbin/nologin -g csye6225 csye6225 || exit 1",
      "sudo dnf install -y wget || exit 1",
      "wget https://dl.google.com/go/go1.15.7.linux-amd64.tar.gz || exit 1",
      "sudo tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz || exit 1",
      "echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh || exit 1",
      "source /etc/profile.d/go.sh || exit 1",
      "sudo dnf install -y git || exit 1"
    ]
  }

  provisioner "file" {
    source      = "bin/webapp"
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    inline = [
      "sudo mkdir -p /opt/${var.application_name} || exit 1",
      "sudo mv /tmp/webapp /opt/${var.application_name}/ || exit 1",
      "sudo chmod +x /opt/${var.application_name}/webapp || exit 1",
      "sudo chown -R csye6225:csye6225 /opt/${var.application_name} || exit 1",
    ]
  }

  provisioner "file" {
    source      = "config/webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    inline = [
      "chmod 644 /etc/systemd/system/${var.service_name} || exit 1",
      "sudo mv /tmp/webapp.service /etc/systemd/system/ || exit 1",
      "sudo systemctl daemon-reload || exit 1",
      "sudo systemctl enable ${var.service_name} || exit 1",
      "sudo systemctl start ${var.service_name} || exit 1"  // Optionally start the service
    ]
  }
}
