# Create SDC. This creates an SDC entry in CDO, and does not bootstrap the SDC. This SDC is configured to be created in AWS; disable this by setting `var.create_resources_in_aws` to false.
resource "cdo_sdc" "sdc-in-aws" {
  name = "weilluo-sdc-in-aws"
}

# Create SDC instance in the private subnet of the AWS VPC. Disable this by setting `var.create_resources_in_aws` to false.
module "sdc-instance-in-aws" {
#  source             = "CiscoDevNet/cdo-sdc/aws"
#  version            = "0.0.4"
  source             = "./terraform-aws-cdo-sdc"
  cdo_bootstrap_data = cdo_sdc.sdc-in-aws.bootstrap_data # Get the bootstrap data from CDO and pass it to the AWS instance.
  instance_name      = "cdo-provider-example-sdc-in-aws"
  # Deploy the instance in the private subnet of the VPC you created.
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
  region  = data.aws_region.current
}