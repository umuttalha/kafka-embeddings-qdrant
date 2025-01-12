variable "hcloud_token" {
  description = "Hetzner Cloud API Token"
  type        = string
  sensitive   = true
}



variable "ssh_key_name" {
  description = "ssh key name"
  type        = string
  sensitive   = true
}

