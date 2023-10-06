output "sdc_bootstrap_data" {
  value = module.example_sdc.sdc_bootstrap_data
  sensitive = true
}

output "sdc_name" {
  value = module.example_sdc.sdc_name
}