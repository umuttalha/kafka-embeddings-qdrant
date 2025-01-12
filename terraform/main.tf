terraform {
  required_providers {
    hcloud = {
      source = "hetznercloud/hcloud"
      version = "~> 1.45"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

data "hcloud_ssh_key" "default" {
  name = var.ssh_key_name
}


# Kafka Sunucusu
resource "hcloud_server" "kafka_server" {
  name        = "kafka-server"
  server_type = "cx22"
  image       = "ubuntu-22.04"
  location    = "nbg1"
  ssh_keys    = [data.hcloud_ssh_key.default.id]

  network {
    network_id = hcloud_network.private_network.id
    ip         = "10.0.1.10"
  }
}

# Python Consumer Sunucusu
resource "hcloud_server" "python_server" {
  name        = "python-server"
  server_type = "cx22"
  image       = "ubuntu-22.04"
  location    = "nbg1"
  ssh_keys    = [data.hcloud_ssh_key.default.id]

  network {
    network_id = hcloud_network.private_network.id
    ip         = "10.0.1.20"
  }
}

# Golang Producer Sunucusu
resource "hcloud_server" "golang_server" {
  name        = "golang-server"
  server_type = "cx22"
  image       = "ubuntu-22.04"
  location    = "nbg1"
  ssh_keys    = [data.hcloud_ssh_key.default.id]

  network {
    network_id = hcloud_network.private_network.id
    ip         = "10.0.1.30"
  }
}

# Qdrant Veritabanı Sunucusu
resource "hcloud_server" "qdrant_server" {
  name        = "qdrant-server"
  server_type = "cx22"
  image       = "ubuntu-22.04"
  location    = "nbg1"
  ssh_keys    = [data.hcloud_ssh_key.default.id]

  network {
    network_id = hcloud_network.private_network.id
    ip         = "10.0.1.40"
  }
}

# Özel Network
resource "hcloud_network" "private_network" {
  name     = "private-network"
  ip_range = "10.0.0.0/16"
}

resource "hcloud_network_subnet" "network_subnet" {
  network_id   = hcloud_network.private_network.id
  type         = "cloud"
  network_zone = "eu-central"
  ip_range     = "10.0.1.0/24"
}

# Firewall Rules
resource "hcloud_firewall" "common_firewall" {
  name = "common-firewall"

  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "22"
    source_ips = ["0.0.0.0/0"]
  }
}

output "kafka_ip" {
  value = hcloud_server.kafka_server.ipv4_address
}

output "python_ip" {
  value = hcloud_server.python_server.ipv4_address
}

output "golang_ip" {
  value = hcloud_server.golang_server.ipv4_address
}

output "qdrant_ip" {
  value = hcloud_server.qdrant_server.ipv4_address
}


