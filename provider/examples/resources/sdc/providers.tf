terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwid3JpdGUiLCI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciXSwicm9sZXMiOlsiUk9MRV9TVVBFUl9BRE1JTiJdLCJhbXIiOiJzYW1sIiwiaXNzIjoiaXRkIiwiY2x1c3RlcklkIjoiMSIsImlkIjoiMjZkNWE2YjQtMTBhMS00NGY0LTgwZTUtMmYzMTI3ZmRmNmYzIiwic3ViamVjdFR5cGUiOiJ1c2VyIiwianRpIjoiMzZjOTg2NGMtZjg4Ni00NjI2LWIwYWQtODYwYmY0ODk1NzY2IiwicGFyZW50SWQiOiI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciLCJjbGllbnRfaWQiOiJhcGktY2xpZW50In0.g3wH-IdIG7u-boB5Ln1Ft2txoq6aPwqHCyOpSJ2deYf3vJkqNuO_P_DI87iayH3F3fQ845YplTDrkPOMjupVGAruNq8L52MLCfQZX9FjarUT9Gre47d6SxWiUA8GLS-Bumh1TuzjFGSLxYFWvQAi-X36X25ezY3vb3-qtlPzyasuxc3TYhBXm5q8FGQCIqodJ6024O7iN7hk8cchvayVrIuurxT2tZRY6g93wFwuhZRrnH6Ds6UB_83fFmaTQoYbV6NX4ny2cD9WuOSF0nf5TYZkQob4PfirtsGuRMJyDqBQovHN2D0UJ417QRYRsi2OZ8vlrvEIgUv-kdkiHcG0ag"
}


provider "aws" {
  region = "us-east-1"
}

data "aws_region" "current" {}

output  "aws_region" {
  value = data.aws_region.current
}