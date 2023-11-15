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
  type = string
}

variable "dns_prefix" {
  description = "The DNS name in the hosted zone to connect the load balancer. The DNS name will be: {dns_prefix}.{hosted_zone_name}."
  type = string
}
