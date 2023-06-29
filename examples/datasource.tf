# Monitoring Data Source
resource "ns1_datasource" "my_monitoring_data_source" {
  name       = "My Monitoring Data Source"
  sourcetype = "nsone_monitoring"
}

# Custom Data Source
resource "ns1_datasource" "my_custom_data_source" {
  name       = "My Custom Data Source"
  sourcetype = "nsone_v1"
}
