terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
#   base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
#   aPi_token = "<replace-with-api-token-generated-from-cdo>"
      # base_url = "https://staging.dev.lockhart.io"
    # api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwid3JpdGUiLCI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciXSwiYW1yIjoic2FtbCIsInJvbGVzIjpbIlJPTEVfU1VQRVJfQURNSU4iXSwiaXNzIjoiaXRkIiwiY2x1c3RlcklkIjoiMSIsImlkIjoiNjAyMjYxNTUtOTRjYi00YWY5LWIzYTQtZDk0ZTcxYjhmOThkIiwic3ViamVjdFR5cGUiOiJ1c2VyIiwianRpIjoiMTVmZmU0NWYtODJiZC00MmYzLTgwYWQtYzgwMzg0MTFiM2ExIiwicGFyZW50SWQiOiI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciLCJjbGllbnRfaWQiOiJhcGktY2xpZW50In0.vaNSnJTsFe-dr79xXs0cD3Jz1ayHhWCmFTgXs-nX-pwU8BbyytTRaUH0m-s-hDpzxdPuuv3J_nyuWHDIE6cd1QB0SV_WiFkYSNkOBX2JTu_oMeBso1ffYxsNc2gCg0-NK7A67RJjgD2GZJnpseOkYU8H7vWPEHCpxGS9JQpOWvdTCEgbyBtm5h3bmMkQNUx82obqODi-ZlQK1Qro2Q8c050MAf64As4-TlriIgwoqaJxb6W9REAbrL2fOD9p0IXiqTaY0OAuo5IcxL794GX63ey8kOHoxYnnDA2stIQQiqMLO7kLeYrPAgZvMOixm__O8-7k-azQuLuvKAUMxam_OQ"
    base_url = "http://localhost:9000"
    api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwiMTExMTExMTEtMTExMS0xMTExLTExMTEtMTExMTExMTExMTExIiwid3JpdGUiXSwiYW1yIjoicHdkIiwicm9sZXMiOlsiUk9MRV9TVVBFUl9BRE1JTiJdLCJpc3MiOiJpdGQiLCJjbHVzdGVySWQiOiIxIiwiaWQiOiI0OWQxYWVjOC00ZmUwLTQyNGYtYjNkYy1kMWU0YmIyYjZmMjMiLCJzdWJqZWN0VHlwZSI6InVzZXIiLCJqdGkiOiJmNGE3OThhYi1lNDFjLTQxNjQtYTYxMC01MjViMWU4MmZjOWIiLCJwYXJlbnRJZCI6IjExMTExMTExLTExMTEtMTExMS0xMTExLTExMTExMTExMTExMSIsImNsaWVudF9pZCI6ImFwaS1jbGllbnQifQ.eHAHgmW6aZxSmqyZEaKvqLPoMFETmdu5PVEXXmR1FOc3LY5msbW-2SENcuagRVtiMomHDyHGulEgxpRdS-lkE2mNKYPSKPfacdOmOy-MR4iS_ShDQFWjCDe4xx3HiogT9MOhGAqMEB4AYHNDjRVgT5OCC_MyUUZr1J2Ja4f1s7ZsdNVljl7MZnRHEXome51Twil6SqBtiukbMpq5JaQ3khW9hHziJ3rgCA5jnXFhZXWNWKmIVPd8yoSAy1d2ubArzTkFJvAm251VIs5GvHepGdVg1IcnifCYKsOylB0_BoCI7_7xqKAKvhNOE3e5lzybqwwLQCOxD1QSy2mWQ7tk3w"
}


resource "cdo_user" "new_api_only_user" {
  name             = "api_user@kunji.com"
  is_api_only_user = false
  role             = "ROLE_ADMIN"
}

resource "cdo_api_token" "new_api_only_user_api_token" {
    username = cdo_user.new_api_only_user.generated_username
}

output "api_only_user_api_token_value" {
    value = cdo_api_token.new_api_only_user_api_token.api_token
    sensitive = true
}