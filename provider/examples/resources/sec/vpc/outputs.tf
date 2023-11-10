output "vpc_id" {
  description = "ID of the cdo-provider example VPC created by this module."
  value       = aws_vpc.vpc.id
}

output "public_subnet_1_id" {
  description = "ID of the cdo-provider example public subnet 1 created by this module."
  value       = aws_subnet.public_subnet_1.id
}

output "public_subnet_2_id" {
  description = "ID of the cdo-provider example public subnet 2 created by this module."
  value       = aws_subnet.public_subnet_2.id
}

output "private_subnet_id" {
  description = "ID of the cdo-provider example private subnet created by this module."
  value       = aws_subnet.private_subnet.id
}