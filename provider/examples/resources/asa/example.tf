terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "4.66.1"
    }
  }
}

provider "cdo" {
  base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
  api_token = "<replace-with-api-token-generated-from-cdo>"
}

module "bastion" {
  source    = "./bastion"
  subnet_id = module.vpc.public_subnet_id
  vpc_id    = module.vpc.vpc_id
  base_name = local.base_name
}

module "vpc" {
  source = "./aws_vpc"
}

resource "random_password" "asa_password" {
  length           = 16
  override_special = "@!"
}

resource "random_password" "asa_enable_password" {
  length           = 16
  override_special = "@!"
}

locals {
  asa_version        = "9-19"
  asav_instance_size = "c5a.large"
  asa_hostname       = "example-asa-hostname"
  asa_username       = "example-asa-username"
  asa_port           = 443
  connector_name     = "sdc-in-aws"
  base_name          = "example"
}

data "aws_region" "current" {}

module "terraform-managed-asav-01" {
  source              = "./asav"
  base_name           = local.base_name
  vpc_id              = module.vpc.vpc_id
  public_subnets      = [module.vpc.public_subnet_id]
  private_subnets     = [module.vpc.private_subnet_id]
  asa_hostname        = local.asa_hostname
  bastion_sg          = module.bastion.bastion_sg
  asa_username        = local.asa_username
  asa_password        = random_password.asa_password.result
  enable_password     = random_password.asa_enable_password.result
  asa_version         = local.asa_version
  asav_instance_size  = local.asav_instance_size
  aws_region          = data.aws_region.current.id
  bastion_public_ip   = module.bastion.bastion_ip
  bastion_private_key = module.bastion.bastion_private_key
}

# Create SDC. This creates an SDC entry in CDO, and does not bootstrap the SDC. This SDC is configured to be created in AWS; disable this by setting `var.create_resources_in_aws` to false.
resource "cdo_sdc" "sdc-in-aws" {
  name = local.connector_name
}

resource "cdo_sdc_onboarding" "sdc-in-aws" {
  name = cdo_sdc.sdc-in-aws.name
}

# Create SDC instance in the private subnet of the AWS VPC. Disable this by setting `var.create_resources_in_aws` to false.
module "sdc-instance-in-aws" {
  source             = "CiscoDevNet/cdo-sdc/aws"
  version            = "0.0.6"
  cdo_bootstrap_data = cdo_sdc.sdc-in-aws.bootstrap_data
  # Get the bootstrap data from CDO and pass it to the AWS instance.
  instance_name      = "cdo-provider-example-sdc-in-aws"
  # Deploy the instance in the private subnet of the VPC you created.
  vpc_id             = module.vpc.vpc_id
  subnet_id          = module.vpc.private_subnet_id
  env                = "example-local"
}

# onboard ASAv in private subnet using SDC deployed to the same subnet in 04-sdc-aws.tf
resource "cdo_asa_device" "example" {
  name           = local.asa_hostname
  username       = local.asa_username
  password       = random_password.asa_password.result
  socket_address = "${module.terraform-managed-asav-01.mgmt_interface_ip}:${local.asa_port}"

  connector_type = "SDC"
  connector_name = cdo_sdc.sdc-in-aws.name

  ignore_certificate = true

  labels = ["333", "444", "111", "222"]
}

