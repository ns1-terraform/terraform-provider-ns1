resource "ns1_team" "foobar" {
  name = "terraform test"

  dns_view_zones             = true
  dns_zones_allow_by_default = true
  dns_zones_allow            = ["mytest.zone"]
  dns_zones_deny             = ["myother.zone"]

  data_manage_datasources = true
}
