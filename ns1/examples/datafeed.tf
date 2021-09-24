resource "ns1_datasource" "api" {
  name       = "terraform test"
  sourcetype = "nsone_v1"
}

resource "ns1_datafeed" "foobar" {
  name      = "terraform test"
  source_id = ns1_datasource.api.id

  #optional
  config = {
    label = "exampledc2"
  }
}
