data "aws_ami" "bastion_host" {
  owners      = ["099720109477"]
  most_recent = true
  filter {
    name   = "name"
    values = ["*ubuntu-jammy*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}
