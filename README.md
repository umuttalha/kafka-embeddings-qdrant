# Kafka-Embeddings-Qdrant Deployment with Ansible and Terraform

This guide demonstrates how to deploy a Kafka-Embeddings-Qdrant setup using Ansible for configuration management and Terraform for infrastructure provisioning.

## Ansible Playbook Command

To run the Ansible playbook, execute:

```bash
ansible-playbook -i inventory.yml playbook.yml
```

## Example `inventory.yml`

```yaml
all:
  children:
    kafka_servers:
      hosts:
        kafka1:
          ansible_host: <ip_address>
          ansible_user: root
    qdrant_servers:
      hosts:
        qdrant1:
          ansible_host: <ip_address>
          ansible_user: root
    python_servers:
      hosts:
        python1:
          ansible_host: <ip_address>
          ansible_user: root
    golang_servers:
      hosts:
        golang1:
          ansible_host: <ip_address>
          ansible_user: root
```

Replace `<ip_address>` with the actual IP addresses of your servers.

## Example `terraform.tfvars`

```hcl
hcloud_token = ""

ssh_key_name = ""
```

- **`hcloud_token`**: Your Hetzner Cloud API token.
- **`ssh_key_name`**: The name of the SSH key to use for server access.

## Usage Notes

1. **Ansible**:
   - Ensure that the target servers are reachable via SSH.
   - Update `inventory.yml` with the correct server details.

2. **Terraform**:
   - Populate `terraform.tfvars` with your Hetzner Cloud API token and SSH key name.
   - Use Terraform to provision the infrastructure before running the Ansible playbook.

This setup allows for efficient deployment and management of Kafka, Qdrant, and supporting services.
