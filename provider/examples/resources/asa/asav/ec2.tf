resource "aws_instance" "asav" {
  ami           = data.aws_ami.asav.id
  instance_type = var.asav_instance_size
  key_name      = aws_key_pair.rsa_keypair.id
  tags = {
    Name = "${var.base_name}-${var.asa_hostname}-cisco-asav"
  }
  network_interface {
    network_interface_id = aws_network_interface.asa_mgmt.id
    device_index         = 0
  }
  network_interface {
    network_interface_id = aws_network_interface.asa_outside.id
    device_index         = 1
  }
  network_interface {
    network_interface_id = aws_network_interface.asa_inside.id
    device_index         = 2
  }
  user_data = data.template_file.asav_startup_config.rendered

  lifecycle {
    ignore_changes = [ami]
  }
}
