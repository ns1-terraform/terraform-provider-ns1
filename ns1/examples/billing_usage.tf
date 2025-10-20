locals {
  now           = timestamp()
  now_unix      = provider["time"].rfc3339_parse(local.now).unix
  cur_mon = formatdate("YYYY-MM", local.now)
  begin_mon          = "${local.cur_mon}-01T00:00:00Z"
  bmon_unix     = provider["time"].rfc3339_parse(local.begin_mon).unix
  end_mon          = timeadd(local.begin_mon, "720h")
  emon_unix     = provider["time"].rfc3339_parse(local.end_mon).unix
}

# Get query usage data for the given timeframe
data "ns1_billing_usage" "queries" {
  metric_type = "queries"
  from = local.bmon_unix  # beginning of the month
  to = local.emon_unix    # end of the month
}

# Get account limits data for the given timeframe
data "ns1_billing_usage" "limits" {
  metric_type = "limits"
  from = local.bmon_unix  # beginning of the month
  to = local.emon_unix    # end of the month
}

# Get RUM decisions usage data for the given timeframe
data "ns1_billing_usage" "decisions" {
  metric_type = "decisions"
  from = local.bmon_unix  # beginning of the month
  to = local.emon_unix    # end of the month
}

# Get filter chains usage data
data "ns1_billing_usage" "filter_chains" {
  metric_type = "filter-chains"
}

# Get monitoring jobs usage data
data "ns1_billing_usage" "monitors" {
  metric_type = "monitors"
}

# Get records usage data
data "ns1_billing_usage" "records" {
  metric_type = "records"
}
