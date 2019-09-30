---
layout: "ns1"
page_title: "NS1: ns1_zone"
sidebar_current: "docs-ns1-datasource-zone"
description: |-
  Provides details about a NS1 Zone.
---

# Data Source: ns1_zone

Provides details about a NS1 Zone. Use this if you would simply like to read
information from NS1 into your configurations. For read/write operations, you
should use a resource.

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

In addition to the argument above, the following are exported:

* `link` - The linked target zone.
* `primary` - The primary ip.
* `additional_primaries` - List of additional IPs for the primary zone.
* `ttl` - The SOA TTL.
* `refresh` - The SOA Refresh.
* `retry` - The SOA Retry.
* `expiry` - The SOA Expiry.
* `nx_ttl` - The SOA NX TTL.
* `dnssec` - Whether or not DNSSEC is enabled for the zone.
* `networks` - List of network IDs for which the zone is available.
* `dns_servers` - Authoritative Name Servers.
* `hostmaster` - The SOA Hostmaster.
* `secondaries` - List of secondary servers. [Secondaries](#secondaries-1) is
  documented below.

#### Secondaries

A secondary has the following fields:

* `ip` - IPv4 address of the secondary server.
* `port` - Port of the the secondary server. Default `53`.
* `notify` - Whether we send `NOTIFY` messages to the secondary host
  when the zone changes. Default `false`.
* `networks` - List of network IDs (`int`) for which the zone should be made
  available. Default is network 0, the primary NSONE Global Network.
