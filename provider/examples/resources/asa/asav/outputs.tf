output "mgmt_interface_ip" {
  value       = aws_network_interface.asa_mgmt.private_ip
  description = "The IP Address of the Management Interface of the ASA."
}