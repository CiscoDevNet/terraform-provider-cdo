data "aws_availability_zones" "available" {}

data "aws_ami" "asav" {
  most_recent = true

  filter {
    name   = "name"
    values = ["*${var.asa_version}*"]
  }

  filter {
    name   = "product-code"
    values = ["663uv4erlxz65quhgaz9cida0"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "template_file" "asav_startup_config" {
  template = file("${path.module}/asav_startup_config.txt")
  vars = {
    asa_hostname    = var.asa_hostname
    enable_password = var.enable_password
    asa_username    = var.asa_username
    asa_password    = var.asa_password
  }
}
