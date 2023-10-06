variable "vpc_id" {
  description = "Specify the VPC to deploy the SDC in"
  type        = string
}

variable "subnet_id" {
  description = "Specify the subnet to deploy the SDC in."
  type        = string
}

variable "resource_prefix" {
  description = "Prefix applied to name of the resources created."
  type        = string
}