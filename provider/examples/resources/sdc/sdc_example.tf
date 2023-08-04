terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwid3JpdGUiLCI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciXSwiYW1yIjoic2FtbCIsInJvbGVzIjpbIlJPTEVfU1VQRVJfQURNSU4iXSwiaXNzIjoiaXRkIiwiY2x1c3RlcklkIjoiMSIsImlkIjoiMjZkNWE2YjQtMTBhMS00NGY0LTgwZTUtMmYzMTI3ZmRmNmYzIiwic3ViamVjdFR5cGUiOiJ1c2VyIiwianRpIjoiMzFiMjNjODMtYjE1NC00NzFiLWE5ZjAtMGM3ZjhhNTIzMmRiIiwicGFyZW50SWQiOiI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciLCJjbGllbnRfaWQiOiJhcGktY2xpZW50In0.2Wav57iXmb2v-feNEleqpDk8on790bPOrP3NdC9Hlwp-FYiuPhkWO9Ai5m3z5ikeL025qYGC9EfkIiEEysSfReDbN55FkLl7XRHQT9nCGDIShCq6aHqu2EgV0qAgeaMg6R28IamPEaTJQ44H07mGR4w3h56Le3cFgFumUzN7jScnXGkZ4UWkqtjw2UGNF3525f9WA5dzSMIjk4Gvi8f4Wm954C2PWk6DjjefQIeRyypXW4JflVSZWfh2vxBCxdTKTMr0GTi72jckwrXj7R1MHBebJ26Ohb7rHa5-4VgbCXww4R4WucGyejwsHW6rIaxh8zcjXMcsUjwLUjgF6PSxiQ"
}

resource "cdo_sdc" "example" {
  name = "tf-sdc-1"
}

output "sdc_name" {
    value = cdo_sdc.example.name
}

output "sdc_uid" {
    value = cdo_sdc.example.id
}

output "sdc_bootstrap_data" {
    value = cdo_sdc.example.bootstrap_data
}