
# Create SDC. This creates an SDC entry in CDO, and does not bootstrap the SDC. This SDC is configured to be created in AWS; disable this by setting `var.create_resources_in_aws` to false.
resource "cdo_sdc" "example" {
  name = "${var.resource_prefix}-sdc"
}

# Create SDC instance in the private subnet of the AWS VPC. Disable this by setting `var.create_resources_in_aws` to false.
module "sdc-instance-in-aws-example" {
  source             = "CiscoDevNet/cdo-sdc/aws"
  version            = "0.0.6"
  cdo_bootstrap_data = cdo_sdc.example.bootstrap_data # Get the bootstrap data from CDO and pass it to the AWS instance.
  instance_name      = "${var.resource_prefix}-sdc"
  # Deploy the instance in the private subnet of the VPC you created.
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}

resource "cdo_sdc_onboarding" "sdc-in-aws" {
  name = cdo_sdc.example.name
}