# VPC
resource "aws_vpc" "vpc" {
  cidr_block           = var.combined_subnet_ranges["VPC"]
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name : "${local.vpc_resource_prefix}-vpc"
  }
}