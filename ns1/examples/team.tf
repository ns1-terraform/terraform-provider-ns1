resource "ns1_team" "foobar" {
  name = "terraform test"

  dns_view_zones             = true
  dns_zones_allow_by_default = true
  dns_zones_allow            = ["mytest.zone"]
  dns_zones_deny             = ["myother.zone"]

  dns_records_allow {
      domain = "a.example.com"
      include_subdomains = false
      zone = "example.com"
      type = "A"
  }
  dns_records_allow {
      domain = "my.ns1.com"
      include_subdomains = true
      zone = "ns1.com"
      type = "A"
  }
  dns_records_deny {
      domain = "evil-user.com"
      include_subdomains = false
      zone = "evil-user.com"
      type = "A"
  }

  data_manage_datasources = true
}
