---
layout: "ns1"
page_title: "NS1: ns1_billing_usage"
sidebar_current: "docs-ns1-datasource-billing-usage"
description: |-
  Provides billing usage details about a NS1 account.
---

# Data Source: ns1_billing_usage

Provides billing usage details about a NS1 account.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `metric_type` - (Required) The type of billing metric to retrieve. Must be one of: `queries`, `limits`, `decisions`, `filter-chains`, `monitors`, `records`.
* `from` - (Required for `queries`, `limits`, and `decisions`) The start timestamp for the data range in Unix epoch format.
* `to` - (Required for `queries`, `limits`, and `decisions`) The end timestamp for the data range in Unix epoch format.

## Attribute Reference

The following attributes are exported:

### Common Attributes

* `total_usage` - (Computed) The total usage count for the metric. Available for `decisions`, `filter-chains`, `monitors`, and `records` metrics.

### Queries Metric Attributes

* `clean_queries` - (Computed) The total number of clean queries (excluding DDoS and NXD).
* `ddos_queries` - (Computed) The total number of DDoS queries.
* `nxd_responses` - (Computed) The total number of NXD responses.
* `by_network` - (Computed) A list of network-specific query data containing:
  * `network` - The network ID.
  * `clean_queries` - Clean queries for this network.
  * `ddos_queries` - DDoS queries for this network.
  * `nxd_responses` - NXD responses for this network.
  * `billable_queries` - Total billable queries for this network.
  * `daily` - Daily breakdown containing:
    * `timestamp` - The timestamp for the day.
    * `clean_queries` - Clean queries for this day.
    * `ddos_queries` - DDoS queries for this day.
    * `nxd_responses` - NXD responses for this day.

### Limits Metric Attributes

* `queries_limit` - (Computed) The queries limit for this billing cycle.
* `china_queries_limit` - (Computed) The queries limit for the China network.
* `records_limit` - (Computed) The records limit for this billing cycle.
* `filter_chains_limit` - (Computed) The filter chains limit for this billing cycle.
* `monitors_limit` - (Computed) The monitoring jobs limit for this billing cycle.
* `decisions_limit` - (Computed) The RUM decisions limit for this billing cycle.
* `nxd_protection_enabled` - (Computed) Whether NXD Protection is enabled.
* `ddos_protection_enabled` - (Computed) Whether DDoS Protection is enabled.
* `include_dedicated_dns_network_in_managed_dns_usage` - (Computed) Whether dedicated DNS usage counts towards managed DNS usage.