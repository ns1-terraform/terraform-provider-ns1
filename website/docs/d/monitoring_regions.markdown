---
layout: "ns1"
page_title: "NS1: ns1_monitoring_regions"
sidebar_current: "docs-ns1-datasource-monitoring-regions"
description: |-
  Provides details of all available monitoring regions.
---

# Data Source: ns1_monitoring_regions

Provides details of all available monitoring regions.

## Example Usage

```hcl
# Get details of all available monitoring regions.
data "ns1_monitoring_regions" "example" {
}
```

## Argument Reference

There are no required arguments.

## Attributes Reference

The following are attributes exported:

* `regions` - A set of the available monitoring regions. [Regions](#regions) is
  documented below.

#### Regions

A region has the following fields:

* `code` - 3-letter city code identifying the location of the monitor.
* `name` - City name identifying the location of the monitor.
* `subnets` - A list of IPv4 and IPv6 subnets the monitor sources requests from.