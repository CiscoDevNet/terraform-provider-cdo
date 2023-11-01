variable "vpc_id" {
  description = "The AWS VPC to deploy the ASAv in."
  type        = string
}

variable "base_name" {
  description = "The base name is a prefix assigned to all of the resources created for the ASAv. It typically corresponds to the name of the CDO environment that the ASAv will be onboarded to."
  type        = string
}

variable "public_subnets" {
  description = "The ASAv is deployed with three interfaces (management, outside, and inside) to a AWS VPC with public subnets accessible from the internet, and private subnets inaccessible from the internet. The management and outside interfaces of the ASAv are deployed in the public subnet. Specify the public subnets in the VPC here."
  type        = list(any)
  validation {
    condition     = length(var.public_subnets) >= 1
    error_message = "The ASAv requires at least one subnet that the management and outside interfaces can be deployed to."
  }
}

variable "private_subnets" {
  description = "The ASAv is deployed with three interfaces (management, outside, and inside) to a AWS VPC with public subnets accessible from the internet, and private subnets inaccessible from the internet. The management and outside interfaces of the ASAv are deployed in the public subnet. Specify the private subnets in the VPC here."
  type        = list(any)
  validation {
    condition     = length(var.private_subnets) >= 1
    error_message = "The ASAv requires at least one subnet that the management and outside interfaces can be deployed to."
  }
}

variable "asav_instance_size" {
  default     = "c5a.large"
  type        = string
  description = "Size of the EC2 instance used to run the SDC). Allowed values: See the [Cisco Secure Firewall ASA Virtual](https://aws.amazon.com/marketplace/pp/prodview-sltshxd3bzqbg) page on the AWS Marketplace."
}

variable "asa_hostname" {
  description = "The hostname of the ASA."
  type        = string
  default     = "asa-01"
}

variable "ssh_key_name" {
  description = "The key to use on the ASAv for SSH."
  default     = "lh-localdev-server"
  type        = string
}

variable "bastion_sg" {
  description = "The ASAv is configured, for security reasons, so that SSH access to its management interface is restricted to instances using a specified security group. Specify the security group here."
  type        = string
}

variable "asa_password" {
  description = "The password for the admin user (with privilege level 15) that will be used to connect to the ASAv."
  type        = string
}

variable "asa_username" {
  description = "The name of the admin user (with privilege level 15) that will be used to connect to the ASAv."
  type        = string
}

variable "asa_version" {
  description = "The version of the ASAv to deploy."
  default     = "9-19"
  type        = string
}

variable "enable_password" {
  description = "The enable password lets you enter privileged EXEC mode. Specify the enable password here."
  type        = string
}

variable "aws_region" {
  description = "The AWS region the ASAv is being deployed in. This is required in order to wait for the ASAv to become available."
  type        = string
}

variable "bastion_public_ip" {
  description = "We expect the ASA to be deployed in a private subnet accessible only through a Bastion host. Provide the public IP of the bastion host so that the module can SSH into it and wait until the ASA is ready."
  type        = string
}

variable "bastion_private_key" {
  description = "We expect the ASA to be deployed in a private subnet accessible only through a Bastion host. Provide the private SSH key of the bastion host so that the module can SSH into it and wait until the ASA is ready."
  type        = string
}