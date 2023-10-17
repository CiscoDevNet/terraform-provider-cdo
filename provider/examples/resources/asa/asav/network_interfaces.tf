resource "aws_network_interface" "asa_mgmt" {
  description       = "${var.base_name}-asav-mgmt"
  subnet_id         = var.private_subnets[0]
  source_dest_check = false
  security_groups   = [aws_security_group.mgmt_interface_security_group.id]
}

resource "aws_network_interface" "asa_outside" {
  description       = "${var.base_name}-asav-outside"
  subnet_id         = var.public_subnets[0]
  source_dest_check = false
  security_groups   = [aws_security_group.outside_interface_security_group.id]
}

resource "aws_network_interface" "asa_inside" {
  description       = "${var.base_name}-asav-inside"
  subnet_id         = var.private_subnets[0]
  source_dest_check = false
  security_groups   = [aws_security_group.inside_interface_security_group.id]
}