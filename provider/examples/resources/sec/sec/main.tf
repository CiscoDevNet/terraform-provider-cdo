data "aws_route53_zone" "selected" {
  name = var.hosted_zone_name
}

# Create SEC. This creates an SEC entry in CDO, and does not bootstrap the SEC. This SEC is configured to be created in AWS; disable this by setting `var.create_resources_in_aws` to false.
resource "cdo_sec" "example" {
}

# Create SDC instance in the private subnet of the AWS VPC. Disable this by setting `var.create_resources_in_aws` to false.
module "sec-instance-in-aws-example" {
  source                     = "CiscoDevNet/cdo-sec/aws"
  version                    = "0.1.0"
  cdo_bootstrap_data         = cdo_sec.example.cdo_bootstrap_data
  # Get the bootstrap data from CDO and pass it to the AWS instance.
  sec_bootstrap_data         = cdo_sec.example.sec_bootstrap_data
  instance_name              = "${var.resource_prefix}-sec"
  # Deploy the instance in the private subnet of the VPC you created.
  vpc_id                     = var.vpc_id
  subnet_id                  = var.subnet_id
  public_subnet_id           = var.lb_public_subnet_id
  secondary_public_subnet_id = var.lb_secondary_public_subnet_id
  hosted_zone_id             = data.aws_route53_zone.selected.id
  dns_name                   = "${var.dns_prefix}.${data.aws_route53_zone.selected.name}"
  env                        = "example"
}

resource "cdo_sec_onboarding" "example" {
  name = cdo_sec.example.name
}