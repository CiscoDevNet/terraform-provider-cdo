terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  # base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
  # api_token = "<replace-with-api-token-generated-from-cdo>"
  base_url = "https://staging.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwid3JpdGUiLCI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciXSwicm9sZXMiOlsiUk9MRV9TVVBFUl9BRE1JTiJdLCJhbXIiOiJzYW1sIiwiaXNzIjoiaXRkIiwiY2x1c3RlcklkIjoiMSIsImlkIjoiNjAyMjYxNTUtOTRjYi00YWY5LWIzYTQtZDk0ZTcxYjhmOThkIiwic3ViamVjdFR5cGUiOiJ1c2VyIiwianRpIjoiYzBiYjVkOGQtNWMzMS00NWEzLTgxMDEtZjhhYjNkNDIxODViIiwicGFyZW50SWQiOiI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciLCJjbGllbnRfaWQiOiJhcGktY2xpZW50In0.pk552gv1J9aUtWKTLUaVqzbREZTO39H1lBY6iAo3g1WD2Ih341pPqLaZjM15v0XsiRET9-NctP1G2ZoMe6W9UmXh7eGp0oSqffZkAUUqJDKXMEXuy2FMV1v3Hsumz2fDE0F7JtaGnhvGpvJVF25kPe0cmKC1lhZ_VOdiOqpApOIrW2chvKW4YWix0pwV7Gz_NqUt2wK9ZbadD5VmZ22loxIKB1zdB5O-IVzyuwkZduXPwknm0BFRXxjRKKJDfLUHsrjYb5AR7Up-XW40P-nyaHrEUE9EJ94DNb9RZiZ3gGjXMfORtKArRyAvXnUxJL0efbHVYIarmV3_1C1K5yticA"
}

data "cdo_cdfmc" "current" {
}

output "cdfmc_hostname" {
    value = data.cdo_cdfmc.current.hostname
}

output "cdfmc_software_version" {
    value = data.cdo_cdfmc.current.software_version
}

output "cdfmc_uid" {
    value = data.cdo_cdfmc.current.id
}

output "cdfmc_domain_uuid" {
    value = data.cdo_cdfmc.current.domain_uuid
}