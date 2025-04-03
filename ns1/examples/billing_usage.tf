# Get query usage data for the given timeframe
resource "ns1_billing_usage" "queries" {
  metric_type = "queries"
  from = 1731605824  # 2024-11-15 00:00:00 UTC
  to = 1734197824    # 2024-12-15 00:00:00 UTC
}

# Get account limits data for the given timeframe
resource "ns1_billing_usage" "limits" {
  metric_type = "limits"
  from = 1731605824  # 2024-11-15 00:00:00 UTC
  to = 1734197824    # 2024-12-15 00:00:00 UTC
}

# Get RUM decisions usage data for the given timeframe
resource "ns1_billing_usage" "decisions" {
  metric_type = "decisions"
  from = 1731605824  # 2024-11-15 00:00:00 UTC
  to = 1734197824    # 2024-12-15 00:00:00 UTC
}

# Get filter chains usage data
resource "ns1_billing_usage" "filter_chains" {
  metric_type = "filter-chains"
}

# Get monitoring jobs usage data
resource "ns1_billing_usage" "monitors" {
  metric_type = "monitors"
}

# Get records usage data
resource "ns1_billing_usage" "records" {
  metric_type = "records"
}
