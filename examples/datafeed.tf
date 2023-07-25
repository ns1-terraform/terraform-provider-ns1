# DNS
resource "ns1_datafeed" "example_com_dns" {
  name      = "[DNS] example.com"
  source_id = ns1_datasource.my_monitoring_data_source.id
  config = {
    jobid = ns1_monitoringjob.example_com_dns.id
  }
}

# HTTP
resource "ns1_datafeed" "example_com_http" {
  name      = "[HTTP] example.com"
  source_id = ns1_datasource.my_monitoring_data_source.id
  config = {
    jobid = ns1_monitoringjob.example_com_http.id
  }
}

# PING
resource "ns1_datafeed" "example_com_ping" {
  name      = "[PING] example.com"
  source_id = ns1_datasource.my_monitoring_data_source.id
  config = {
    jobid = ns1_monitoringjob.example_com_ping.id
  }
}

# TCP
resource "ns1_datafeed" "example_com_tcp" {
  name      = "[TCP] example.com"
  source_id = ns1_datasource.my_monitoring_data_source.id
  config = {
    jobid = ns1_monitoringjob.example_com_tcp.id
  }
}
