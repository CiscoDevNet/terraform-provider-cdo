output "sdc_name" {
  value       = cdo_sdc.sdc-in-aws.name
  description = "The name of the SDC spun up in AWS."
}

output "sdc_bootstrap_data" {
  value       = cdo_sdc.sdc-in-aws.bootstrap_data
  description = "The bootstrap data of the SDC spun up in AWS."
}