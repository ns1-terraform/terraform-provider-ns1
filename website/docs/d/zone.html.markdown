---
layout: "ns1"
page_title: "NS1: ns1_zone"
sidebar_current: "docs-ns1-datasource-zone"
description: |-
  Provides details about a NS1 Zone.
---

# Data Source: ns1_zone

Provides details about a NS1 Zone.

## Example Usage

```hcl
# Get details about a NS1 Zone.
data "ns1_zone" "example" {
  zone = "terraform.example.io"
}
```

## Argument Reference

* `zone` - (Required) The domain name of the zone.

## Attributes Reference

* `link` - The linked target zone.
* `ttl` - The SOA TTL.
* `refresh` - The SOA Refresh.
* `retry` - The SOA Retry.
* `expiry` - The SOA Expiry.
* `nx_ttl` - The SOA NX TTL.
* `primary` - The primary ip.
* `dns_servers` - Authoritative Name Servers.
