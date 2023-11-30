# variables

variable "cdo_bootstrap_data" {}
variable "sec_bootstrap_data" {}

variable "vpc_id" {
  description = "Specify the VPC to deploy the SEC in"
  type        = string
}

variable "subnet_id" {
  description = "Specify the subnet to deploy the SEC in."
  type        = string
}

variable "lb_public_subnet_id" {
  description = "Specify the public subnet to deploy the Load balancer in front of SEC."
  type        = string
}

variable "lb_secondary_public_subnet_id" {
  description = "Specify the secondary public subnet to deploy the Load balancer in front of SEC."
  type        = string
}

variable "resource_prefix" {
  description = "Prefix applied to name of the resources created."
  type        = string
}

variable "hosted_zone_name" {
  description = "Hosted zone to connect the load balancer."
  type        = string
}

variable "dns_prefix" {
  description = "The DNS name in the hosted zone to connect the load balancer. The DNS name will be: {dns_prefix}.{hosted_zone_name}."
  type        = string
}

# data & resources

data "aws_route53_zone" "selected" {
  name = var.hosted_zone_name
}

# Create SDC instance in the private subnet of the AWS VPC. Disable this by setting `var.create_resources_in_aws` to false.
module "sec-instance-in-aws-example" {
  source                     = "CiscoDevNet/cdo-sec/aws"
  version                    = "0.1.0"
  cdo_bootstrap_data         = var.cdo_bootstrap_data
  # Get the bootstrap data from CDO and pass it to the AWS instance.
  sec_bootstrap_data         = var.sec_bootstrap_data
  instance_name              = "${var.resource_prefix}-sec"
  # Deploy the instance in the private subnet of the VPC you created.
  vpc_id                     = var.vpc_id
  subnet_id                  = var.subnet_id
  public_subnet_id           = var.lb_public_subnet_id
  secondary_public_subnet_id = var.lb_secondary_public_subnet_id
  hosted_zone_id             = data.aws_route53_zone.selected.id
  dns_name                   = "${var.dns_prefix}.${data.aws_route53_zone.selected.name}"
  env                        = var.resource_prefix
  tags                       = {
    ApplicationName = "Terraform Provider CDO"
    ServiceName     = "SEC"
    ResourcePrefix  = var.resource_prefix
  }
}

# outputs

output "sec_fqdn" {
  value = module.sec-instance-in-aws-example.sec_fqdn
}

output "sec_instance_id" {
  value = module.sec-instance-in-aws-example.instance_id
}