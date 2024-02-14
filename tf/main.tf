resource "ns1_dataset" "my_dataset" {
  name     = "my date"
  datatype = {
    type  = "num_queries"
    scope = "account"
    data  = []
  }
  timeframe = {
    aggregation = "monthly"
    cycles      = 1
  }
  repeat      = null
  export_type = "csv"
}