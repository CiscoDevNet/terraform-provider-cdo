output "asa-username" {
  value = local.asa_username
}

output "asa-password" {
  value     = random_password.asa_password
  sensitive = true
}

output "asa-enable-password" {
  value     = random_password.asa_enable_password
  sensitive = true
}

output "asa-hostname" {
  value = local.asa_hostname
}

output "asa-ip" {
  value = module.terraform-managed-asav-01.mgmt_interface_ip
}