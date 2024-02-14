resource "ns1_dataset" "my_dataset" {
  name     = "%s"
  datatype {
    type  = "num_queries"
    scope = "account"
    data  = {}
  }
  timeframe {
    aggregation = "monthly"
    cycles      = 1
  }
  export_type = "csv"
}