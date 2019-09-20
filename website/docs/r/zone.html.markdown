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

# Create a new primary zone
resource "ns1_zone" "example_primary" {
  zone     = "terraform-primary.example.io"
  secondaries {
    ip     = "2.2.2.2"
  }
  secondaries {
    ip     = "3.3.3.3"
    port   = 5353
    notify = true
  }
}

# Create a new secondary zone
resource "ns1_zone" "example_primary" {
  zone     = "terraform-primary.example.io"
  primary  = "2.2.2.2"
  additional_primaries = ["3.3.3.3", "4.4.4.4"]
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The domain name of the zone.
* `link` - (Optional) The target zone(domain name) to link to.
* `primary` - (Optional) The primary zones' IP. This makes the zone a
  secondary. Conflicts with `secondaries`.
* `additional_primaries` - (Optional) List of additional IPs for the primary
  zone. Conflicts with `secondaries`.
* `ttl` - (Optional/Computed) The SOA TTL.
* `refresh` - (Optional/Computed) The SOA Refresh.
* `retry` - (Optional/Computed) The SOA Retry.
* `expiry` - (Optional/Computed) The SOA Expiry.
* `nx_ttl` - (Optional/Computed) The SOA NX TTL.
* `networks` - (Optional/Computed) List of network IDs for which the zone is
  available. If no network is provided, the zone will be created in network 0,
  the primary NS1 Global Network.
* `secondaries` - (Optional) List of secondary servers. This makes the zone a
  primary. Conflicts with `primary` and `additional_primaries`.
  [Secondaries](#secondaries-1) is documented below.

#### Secondaries

A zone can have zero or more `secondaries`. Note how this is implemented in the
example above. A secondary has the following fields:

* `ip` - (Required) IPv4 address of the secondary server.
* `port` - (Optional) Port of the the secondary server. Default `53`.
* `notify` - (Optional) Whether we send `NOTIFY` messages to the secondary host
  when the zone changes. Default `false`.
* `networks` - (Computed) - List of network IDs (`int`) for which the zone
  should be made available. Default is network 0, the primary NSONE Global
  Network. Normally, you should not have to worry about this.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dns_servers` - (Computed) Authoritative Name Servers.
* `hostmaster` - (Computed) The SOA Hostmaster.

## Import

`terraform import ns1_zone.<name> <zone>`

So for the example above:

`terraform import ns1_zone.example terraform.example.io`
