# A site groups subnets by physical location / datacenter.
resource "bagre_site" "dc1" {
  code        = "DC1"
  name        = "Datacenter São Paulo"
  description = "Primary datacenter"
}

output "dc1_id" {
  value = bagre_site.dc1.id
}
