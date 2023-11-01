#resource "null_resource" "wait_for_asa_to_finish_booting" {
#  # If an IP changes, these scripts will run again
#  triggers = {
#    nodes_ips = aws_instance.asav.private_ip
#  }
#  connection {
#    type        = "ssh"
#    user        = "ubuntu"
#    host        = var.bastion_public_ip
#    private_key = var.bastion_private_key
#  }
#
#  provisioner "file" {
#    source      = "${path.module}/wait_for_asa.sh"
#    destination = "/tmp/wait_for_asa.sh"
#  }
#
#  provisioner "remote-exec" {
#    inline = [
#      "chmod u+x /tmp/wait_for_asa.sh",
#      "/tmp/wait_for_asa.sh ${aws_network_interface.asa_mgmt.private_dns_name}"
#    ]
#  }
#}