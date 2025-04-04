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

# Example of using the data in other resources
output "total_queries" {
  value = data.ns1_billing_usage.queries.clean_queries
}

output "total_ddos_queries" {
  value = data.ns1_billing_usage.queries.ddos_queries
}

output "total_nxd_responses" {
  value = data.ns1_billing_usage.queries.nxd_responses
}

output "queries_limit" {
  value = data.ns1_billing_usage.limits.queries_limit
}

output "total_decisions" {
  value = data.ns1_billing_usage.decisions.total_usage
}

output "decisions_limit" {
  value = data.ns1_billing_usage.limits.decisions_limit
}

output "total_filter_chains" {
  value = data.ns1_billing_usage.filter_chains.total_usage
}

output "filter_chains_limit" {
  value = data.ns1_billing_usage.limits.filter_chains_limit
}

output "total_monitors" {
  value = data.ns1_billing_usage.monitors.total_usage
}

output "monitors_limit" {
  value = data.ns1_billing_usage.limits.monitors_limit
}

output "total_records" {
  value = data.ns1_billing_usage.records.total_usage
}

output "records_limit" {
  value = data.ns1_billing_usage.limits.records_limit
}