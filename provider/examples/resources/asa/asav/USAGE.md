## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_null"></a> [null](#provider\_null) | n/a |
| <a name="provider_template"></a> [template](#provider\_template) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_eip.asa_mgmt_eip](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip) | resource |
| [aws_eip_association.asa_mgmt_ip_assocation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip_association) | resource |
| [aws_instance.asav](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance) | resource |
| [aws_network_interface.asa_inside](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/network_interface) | resource |
| [aws_network_interface.asa_mgmt](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/network_interface) | resource |
| [aws_network_interface.asa_outside](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/network_interface) | resource |
| [aws_route53_record.asa_mgmt_if_dns](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_security_group.inside_interface_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group.mgmt_interface_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group.outside_interface_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [null_resource.wait_for_asa_to_be_ready](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) | resource |
| [aws_ami.asav](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ami) | data source |
| [aws_availability_zones.available](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/availability_zones) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [template_file.asav_startup_config](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_asa_hostname"></a> [asa\_hostname](#input\_asa\_hostname) | The hostname of the ASA. | `string` | `"asa-01"` | no |
| <a name="input_asa_password"></a> [asa\_password](#input\_asa\_password) | The password for the admin user (with privilege level 15) that will be used to connect to the ASAv. | `any` | n/a | yes |
| <a name="input_asa_username"></a> [asa\_username](#input\_asa\_username) | The name of the admin user (with privilege level 15) that will be used to connect to the ASAv. | `any` | n/a | yes |
| <a name="input_asa_version"></a> [asa\_version](#input\_asa\_version) | The version of the ASAv to deploy. | `string` | `"9-19"` | no |
| <a name="input_asav_instance_size"></a> [asav\_instance\_size](#input\_asav\_instance\_size) | Size of the EC2 instance used to run the SDC). Allowed values: See the [Cisco Secure Firewall ASA Virtual](https://aws.amazon.com/marketplace/pp/prodview-sltshxd3bzqbg) page on the AWS Marketplace. | `string` | `"c5a.large"` | no |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | The AWS region the ASAv is being deployed in. This is required in order to wait for the ASAv to become available. | `any` | n/a | yes |
| <a name="input_base_name"></a> [base\_name](#input\_base\_name) | The base name is a prefix assigned to all of the resources created for the ASAv. It typically corresponds to the name of the CDO environment that the ASAv will be onboarded to. | `string` | n/a | yes |
| <a name="input_bastion_sg"></a> [bastion\_sg](#input\_bastion\_sg) | The ASAv is configured, for security reasons, so that SSH access to its management interface is restricted to instances using a specified security group. Specify the security group here. | `any` | n/a | yes |
| <a name="input_enable_password"></a> [enable\_password](#input\_enable\_password) | The enable password lets you enter privileged EXEC mode. Specify the enable password here. | `any` | n/a | yes |
| <a name="input_hosted_zone_id"></a> [hosted\_zone\_id](#input\_hosted\_zone\_id) | The ID of a AWS Route53 zone where a DNS entry to the management interface of the ASA will be created. | `any` | n/a | yes |
| <a name="input_private_subnets"></a> [private\_subnets](#input\_private\_subnets) | The ASAv is deployed with three interfaces (management, outside, and inside) to a AWS VPC with public subnets accessible from the internet, and private subnets inaccessible from the internet. The management and outside interfaces of the ASAv are deployed in the public subnet. Specify the private subnets in the VPC here. | `list(any)` | n/a | yes |
| <a name="input_public_subnets"></a> [public\_subnets](#input\_public\_subnets) | The ASAv is deployed with three interfaces (management, outside, and inside) to a AWS VPC with public subnets accessible from the internet, and private subnets inaccessible from the internet. The management and outside interfaces of the ASAv are deployed in the public subnet. Specify the public subnets in the VPC here. | `list(any)` | n/a | yes |
| <a name="input_ssh_key_name"></a> [ssh\_key\_name](#input\_ssh\_key\_name) | The key to use on the ASAv for SSH. | `string` | `"lh-localdev-server"` | no |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | The AWS VPC to deploy the ASAv in. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_mgmt_interface_public_dns"></a> [mgmt\_interface\_public\_dns](#output\_mgmt\_interface\_public\_dns) | The FQDN of the management interface of the ASAv. This FQDN can be used to connect to the management interface of the ASAv on the port the HTTP server is enabled. |
