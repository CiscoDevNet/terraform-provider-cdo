terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwid3JpdGUiLCI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciXSwiYW1yIjoic2FtbCIsInJvbGVzIjpbIlJPTEVfU1VQRVJfQURNSU4iXSwiaXNzIjoiaXRkIiwiY2x1c3RlcklkIjoiMSIsImlkIjoiNjAyMjYxNTUtOTRjYi00YWY5LWIzYTQtZDk0ZTcxYjhmOThkIiwic3ViamVjdFR5cGUiOiJ1c2VyIiwianRpIjoiYzdlZTFjM2QtMGQ5ZC00MTk2LThkYTEtNDc1ZTE0MDk4ZWYxIiwicGFyZW50SWQiOiI1Nzc4YzRiNC02N2EzLTQ4OTUtOTliZC1hMmY2YjMwY2Q2YjciLCJjbGllbnRfaWQiOiJhcGktY2xpZW50In0.OghPYXbC5Fom2js_sZaRRhaqgU9ytHn6QQLoYLOUBJ1vD66GOUDGhDf9dwU4Q7C5WPL_YTHY0a9A83xWyEWfL1sbgwZTKtBP1zpb0mBd5CvJm8GV8gwjXKRphTQjxYydCtwQzXJIWuMyTtafYHa314a0sNOqufQM79W1Cr0Vxtd6z9ZPrczyaIS7sWs6PxTnSBA_ZiQiNeqqCS0YO-PPVU9_qC-4nz7fXa5l9XpyfL_URMSpVi-p8K9S0DiFJl5_PgYJIM7r6Hm3Mwx4bSqtqf0GynfdzSzgkNvZzQN8owNOTEsF4LNC7S2azmR2rIGMCWozF7LIAUoVHXcr9HqqOw"
}

resource "cdo_asa_device" "my_asa" {
  name               = "burak-crush-mango-lassi"
  connector_type     = "CDG"
  socket_address     = "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443"
  username           = "lockhart"
  password           = "BlueSkittles123!!"
  ignore_certificate = false
}
