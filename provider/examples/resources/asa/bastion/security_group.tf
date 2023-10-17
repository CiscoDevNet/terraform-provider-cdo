resource "aws_security_group" "bastion-sg" {
  vpc_id      = var.vpc_id
  name        = "cdo-provider-example-bastion-sg"
  description = "Security Group that allows all egress and ingress only on port 22"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "cdo-provider-example-bastion-sg"
  }
}
