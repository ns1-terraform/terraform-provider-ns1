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

# Example of using the data in other resources
output "total_queries" {
  value = ns1_billing_usage.queries.clean_queries
}

output "total_ddos_queries" {
  value = ns1_billing_usage.queries.ddos_queries
}

output "total_nxd_responses" {
  value = ns1_billing_usage.queries.nxd_responses
}

output "queries_limit" {
  value = ns1_billing_usage.limits.queries_limit
}

output "total_decisions" {
  value = ns1_billing_usage.decisions.total_usage
}

output "decisions_limit" {
  value = ns1_billing_usage.limits.decisions_limit
}

output "total_filter_chains" {
  value = ns1_billing_usage.filter_chains.total_usage
}

output "filter_chains_limit" {
  value = ns1_billing_usage.limits.filter_chains_limit
}

output "total_monitors" {
  value = ns1_billing_usage.monitors.total_usage
}

output "monitors_limit" {
  value = ns1_billing_usage.limits.monitors_limit
}

output "total_records" {
  value = ns1_billing_usage.records.total_usage
}

output "records_limit" {
  value = ns1_billing_usage.limits.records_limit
}