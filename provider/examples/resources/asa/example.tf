terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://ci.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwiMjA5MDc5M2UtNDQ1My00ZjI3LTg2YTYtZDY0YTAxMDQ4N2JjIiwid3JpdGUiXSwicm9sZXMiOlsiUk9MRV9BRE1JTiJdLCJhbXIiOiJwd2QiLCJpc3MiOiJpdGQiLCJjbHVzdGVySWQiOiIxIiwiaWQiOiJmY2NlYTEzMC1jNDQ1LTQ1MzEtOTc1NS1jYmE5NWRlNGYwMTIiLCJzdWJqZWN0VHlwZSI6InVzZXIiLCJqdGkiOiI2MzlhNDgyNy1iYWQyLTQ2ZGUtYTY4ZC0wYTg1OTk1NDdmMGEiLCJwYXJlbnRJZCI6IjIwOTA3OTNlLTQ0NTMtNGYyNy04NmE2LWQ2NGEwMTA0ODdiYyIsImNsaWVudF9pZCI6ImFwaS1jbGllbnQifQ.3gyCSmx6-26dmcTZHFWIVYIDOhDP9cjJhg4wfPqCC-3BlZq6JWdMSNqOP_Q1SXz7GWrqvE4LRqZwj9i0XFWTDWuKvolJY2UcV3k5IVwmDjARn-97pujjjqERCMcm-x5hfJh9SegpVMUYYLqgQE16a-TVGAgzb2qus94qeot7Udga4bdVzwcdZWbxPFGhl8EK4Y9Qx34jHID4QyVW7u4PR2r701c5xh6sBUc4Jy37DU16exeTWOmkxLZdmhwdJv5t2Wm_Sy5BIGhNHoQwPZC_F0-T5i6ockML1iu_7rxpjcY-HxVger91Hliqyzd9gcDr7akUCdvNC0MerBKwZDgJCQ"
}

resource "cdo_asa_device" "my_asa" {
  name               = "wl-asa"
  connector_type     = "CDG"
  socket_address     = "52.53.230.145:443"
  username           = "admin"
  password           = ""
  ignore_certificate = true
}
