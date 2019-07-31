---
layout: "ns1"
page_title: "NS1: ns1_zone"
sidebar_current: "docs-ns1-resource-zone"
description: |-
  Provides a NS1 Zone resource.
---

# ns1\_zone

Provides a NS1 DNS Zone resource. This can be used to create, modify, and delete zones.

## Example Usage

```hcl
# Create a new DNS zone
resource "ns1_zone" "example" {
  zone = "terraform.example.io"
  ttl  = 600
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The domain name of the zone.
* `link` - (Optional) The target zone(domain name) to link to.
* `primary` - (Optional) The primary zones' IP. This makes the zone a secondary.
* `additional_primaries` - (Optional) List of additional IPs for the primary zone.
* `ttl` - (Optional/Computed) The SOA TTL.
* `refresh` - (Optional/Computed) The SOA Refresh.
* `retry` - (Optional/Computed) The SOA Retry.
* `expiry` - (Optional/Computed) The SOA Expiry.
* `nx_ttl` - (Optional/Computed) The SOA NX TTL.
* `networks` - (Optional/Computed) List of network IDs for which the zone is available.  If no network is provided, the zone will be created in network 0, the primary NS1 Global Network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dns_servers` - (Computed) Authoritative Name Servers.
* `hostmaster` - (Computed) The SOA Hostmaster.
