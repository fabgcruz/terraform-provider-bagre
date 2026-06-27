# terraform-provider-bagre

Manage **[Bagre](https://bagre.dev) IPAM** (IP Address Management) as code.

Built on the Terraform Plugin Framework. The same binary works with **[OpenTofu](https://opentofu.org)** and **Terraform** — they share the plugin protocol. Bagre is open-source (MIT), so this provider is developed **OpenTofu-first** while staying fully compatible with Terraform.

> ⚠️ **Status: early / work in progress.** The provider plumbing and the `bagre_site` resource are implemented and tested end-to-end. More resources and data sources are on the way (see Roadmap).

## Requirements

- A running Bagre instance (`>= 0.6` — needs API token support).
- A Bagre **API token** with `READ_WRITE` scope. Generate one under **Tokens de API** (`/admin/api-tokens`).
- OpenTofu `>= 1.6` or Terraform `>= 1.0`.

## Usage

```hcl
terraform {
  required_providers {
    bagre = {
      source = "fabgcruz/bagre"
    }
  }
}

provider "bagre" {
  endpoint  = "https://ipam.example.com" # or env BAGRE_ENDPOINT
  api_token = var.bagre_token            # or env BAGRE_TOKEN
}

resource "bagre_site" "dc1" {
  code        = "DC1"
  name        = "Datacenter São Paulo"
  description = "Primary datacenter"
}
```

Set credentials via environment variables in CI so the token never lands in code:

```sh
export BAGRE_ENDPOINT="https://ipam.example.com"
export BAGRE_TOKEN="bagre_…"
tofu apply   # or: terraform apply
```

## Resources & data sources

| Type | Name | Status |
|------|------|--------|
| resource | `bagre_site` | ✅ implemented |
| resource | `bagre_subnet` | 🔜 planned |
| resource | `bagre_ip_reservation` | 🔜 planned |
| data source | `bagre_next_available_ip` | 🔜 planned |
| data source | `bagre_subnet_by_cidr` | 🔜 planned |

## Development

```sh
make build      # builds ./terraform-provider-bagre
make vet
```

To run a local build against a Bagre instance without publishing, use a dev
override in `~/.tofurc` (or `~/.terraformrc`):

```hcl
provider_installation {
  dev_overrides {
    "registry.opentofu.org/fabgcruz/bagre" = "/absolute/path/to/this/repo"
  }
  direct {}
}
```

Then run `tofu plan` / `tofu apply` directly (no `init` needed with dev overrides).

## License

[MPL-2.0](./LICENSE).
