# Cisco ASAv module

The Cisco ASA is a unified threat management device, combining several network security functions in one box. The Cisco ASA can be deployed on Amazon Web Services in a virtual form factor, known as an ASAv (where the v stands for virtual).

Use this Terraform module if you wish to deploy an ASAv in your AWS VPC.

## Learn more
The [Cisco documentation](https://www.cisco.com/c/en/us/td/docs/security/asa/roadmap/asaroadmap.html) has more details on the ASA.

# Pre-requisities
- An AWS account
- An AWS VPC with at least one private subnet and one public subnet.
- An AWS Route53 Hosted Zone
- Terraform knowledge

# Supported Regions

This module can be used in all regions that [Cisco Secure Firewall ASA Virtual](https://aws.amazon.com/marketplace/pp/prodview-sltshxd3bzqbg) marketplace offering is available.

# Architecture

The ASAv is deployed within a single availability zone in an AWS VPC with three interfaces:
- a Management interface, which is deployed on the public subnet. This interface is used to manage the ASA using CDO (or ASDM), and will be assigned a Route53 DNS entry and elastic IP.
- an Outside interface, which is deployed on the public subnet.
- an Inside interface, which is deployed on the private subnet.

## Security
- Port 443 on the Management subnet is accessible to the world.
- Port 22 on the Management subnet is only accessible by AWS resources that are associated with a security group specified in the configuration. See [the usage documentation](USAGE.md) fore more.
- The outside and inside interface have no ingress rules.

This is a limited Terraform module at the time of writing, and only has enough code to facilitate building an ASAv that can be managed on CDO, and is probably useless for any real firewalling use case.

# Usage
Please see [the usage documentation](USAGE.md) and the example below.