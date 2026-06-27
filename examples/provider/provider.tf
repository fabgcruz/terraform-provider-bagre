terraform {
  required_providers {
    bagre = {
      source = "fabgcruz/bagre"
    }
  }
}

# Endpoint and token can also come from the BAGRE_ENDPOINT / BAGRE_TOKEN
# environment variables — recommended for CI so the token never lands in code.
provider "bagre" {
  endpoint  = "https://ipam.example.com"
  api_token = var.bagre_token # generate under Bagre → Tokens de API
}

variable "bagre_token" {
  type      = string
  sensitive = true
}
