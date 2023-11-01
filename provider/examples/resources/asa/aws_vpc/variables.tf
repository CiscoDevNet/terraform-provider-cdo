variable "combined_subnet_ranges" {
  type = map(string)

  default = {
    "VPC"    = "172.200.0.0/16"
    "public" = "172.200.1.0/24"
  }
}

variable "vpc_cidr" {
  type        = string
  description = "Specify the CIDR block for the AWS VPC."
  default     = "172.200.0.0/16"
}

variable "public_subnet_cidr" {
  description = "Specify the CIDR block for the public subnet."
  default     = "172.200.16.0/20"
}

variable "private_subnet_cidr" {
  description = "Specify the CIDR block for the private subnet."
  default     = "172.200.32.0/20"
}