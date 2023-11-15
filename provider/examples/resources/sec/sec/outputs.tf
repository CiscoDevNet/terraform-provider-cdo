output "sec_name" {
  value       = cdo_sec.example.name
  description = "The name of the SDC spun up in AWS."
}

output "sec_bootstrap_data" {
  value       = cdo_sec.example.sec_bootstrap_data
  description = "The bootstrap data of the SEC spun up in AWS."
  sensitive = true
}

output "cdo_bootstrap_data" {
  value       = cdo_sec.example.cdo_bootstrap_data
  description = "The bootstrap data of the SEC in CDO."
  sensitive = true
}

output "sec_fqdn" {
  value = module.sec-instance-in-aws-example.sec_fqdn
}

output "sec_instance_id" {
  value = module.sec-instance-in-aws-example.instance_id
}