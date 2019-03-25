resource "ns1_datasource" "foobar" {
  #required
  name       = "terraform test"
  sourcetype = "nsone_v1"

  #optional
  config = {}
}
