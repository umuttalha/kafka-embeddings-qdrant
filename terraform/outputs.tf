output "python_server_server_ip" {
  value = hcloud_server.python_server.ipv4_address
}


output "golang_server_server_ip" {
  value = hcloud_server.golang_server.ipv4_address
}


output "qdrant_server_server_ip" {
  value = hcloud_server.qdrant_server.ipv4_address
}


output "kafka_server_server_ip" {
  value = hcloud_server.kafka_server.ipv4_address
}