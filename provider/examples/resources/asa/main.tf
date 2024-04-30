variable "asa_username" {}
variable "asa_hostname" {}

module "vpc" {
  source = "github.com/CiscoDevNet/cdo-automation//modules/aws_vpc"
}

module "bastion" {
  source    = "github.com/CiscoDevNet/cdo-automation//modules/bastion"
  vpc_id    = module.vpc.vpc_id
  subnet_id = module.vpc.public_subnet_id
}

module "sdc_in_aws" {
  source    = "github.com/CiscoDevNet/cdo-automation//modules/cdo/sdc-aws"
  vpc_id    = module.vpc.vpc_id
  subnet_id = module.vpc.private_subnet_id
}

resource "cdo_sdc_onboarding" "sdc" {
  name = module.sdc_in_aws.sdc_name
}

resource "random_password" "asa_password" {
  length = 16
  override_special = "@!"
}

resource "random_password" "asa_enable_password" {
  length = 16
  override_special = "@!"
}

module "asav" {
  source              = "github.com/CiscoDevNet/cdo-automation//modules/cdo/asa"
  bastion_ip          = module.bastion.bastion_ip
  bastion_private_key = module.bastion.bastion_private_key
  bastion_sg          = module.bastion.bastion_sg
  vpc_id              = module.vpc.vpc_id
  public_subnet_id    = module.vpc.public_subnet_id
  private_subnet_id   = module.vpc.private_subnet_id
  sdc_name            = module.sdc_in_aws.sdc_name
  asa_username        = var.asa_username
  asa_hostname        = var.asa_hostname
  asa_password        = random_password.asa_password.result
  asa_enable_password = random_password.asa_enable_password.result

  depends_on = [cdo_sdc_onboarding.sdc]
}

resource "cdo_asa_device" "my-asa" {
  name               = "my-asa"
  socket_address     = module.asav.mgmt_interface_ip
  username           = var.asa_username
  password           = random_password.asa_password.result
  connector_type     = "CDG"
  ignore_certificate = false
 }