# Get query usage data for the given timeframe
data "ns1_billing_usage" "queries" {
  metric_type = "queries"
  from = 1738368000  # 2025-02-01 00:00:00 UTC
  to = 1740787199    # 2025-02-28 23:59:59 UTC
}

# Get account limits data for the given timeframe
data "ns1_billing_usage" "limits" {
  metric_type = "limits"
  from = 1738368000  # 2025-02-01 00:00:00 UTC
  to = 1740787199    # 2025-02-28 23:59:59 UTC
}

# Get RUM decisions usage data for the given timeframe
data "ns1_billing_usage" "decisions" {
  metric_type = "decisions"
  from = 1738368000  # 2025-02-01 00:00:00 UTC
  to = 1740787199    # 2025-02-28 23:59:59 UTC
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
