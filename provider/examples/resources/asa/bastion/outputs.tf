output "bastion_instance_id" {
  value = aws_instance.bastion.id
}

output "bastion_ip" {
  value = aws_instance.bastion.public_ip
}

output "bastion_sg" {
  value = aws_security_group.bastion-sg.id
}

output "bastion_private_key" {
  value = tls_private_key.rsa_keypair.private_key_pem
}