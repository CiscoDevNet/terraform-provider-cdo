terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://ci.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJhZTk4ZDI1Zi0xMDg5LTQyODYtYTNjNS01MDVkY2I0NDMxYTIiLCJyZWFkIiwid3JpdGUiXSwicm9sZXMiOlsiUk9MRV9BRE1JTiJdLCJhbXIiOiJwd2QiLCJpc3MiOiJpdGQiLCJjbHVzdGVySWQiOiIxIiwiaWQiOiI2NzJjMGE0MS1kMjAzLTQ2YzEtYmE5ZS0wNDVmYWUwYTc5ZGQiLCJzdWJqZWN0VHlwZSI6InVzZXIiLCJqdGkiOiJlZjM3ZWM2Yi0yOGFiLTQxNDktYTczYi1hYzBiOTU4ZTMzNDQiLCJwYXJlbnRJZCI6ImFlOThkMjVmLTEwODktNDI4Ni1hM2M1LTUwNWRjYjQ0MzFhMiIsImNsaWVudF9pZCI6ImFwaS1jbGllbnQifQ.NjtSbsv7RN7SgSc-hPXqbTWx4jFT6f4H1Xg_6IrsbueXk7anrffZeoJtVcJhxxrXU79XxFYb3aP6ycwnvJ_TVh9Byye_GyLJog_--HCl0vzKoKua0DnuA8zU0JzBIMLtNYM2J_LSdQakncBeQm1A0nSRQXpOc8kOpQdbulhdiLr9iUV2ybdc8mw1tB-JpSw3B3oKDKe6Z6Tlq-5ijwYxrSUFpG4NtHW4iulMixSy8udO5BR2OvGbxDH8Mhbbzu8JIxynTt6COMaz-ATx2pKe83Qi7sxagwoDjKtFuGjZDWL9fGqhwHZPCe3ucaJn-Of5Z9DPoKpS9MhtOei8F2V6jg"
}

resource "cdo_ftd_device" "test" {
  name               = "test-wl-ftd"
  access_policy_name = "Default Access Control Policy"
  performance_tier   = "FTDv5"
  virtual            = true
  licenses           = ["BASE"]
}

resource "cdo_ftd_device_onboarding" "test" {
  ftd_id = cdo_ftd_device.test.id
}