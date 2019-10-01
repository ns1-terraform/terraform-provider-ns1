---
layout: "ns1"
page_title: "NS1: ns1_dnssec"
sidebar_current: "docs-ns1-datasource-dnssec"
description: |-
  Provides DNSSEC details about a NS1 Zone.
---

# Data Source: ns1_dnssec

Provides DNSSEC details about a NS1 Zone.

## Example Usage

```hcl
# Get DNSSEC details about a NS1 Zone.
resource "ns1_zone" "example" {
  zone   = "terraform.example.io"
  dnssec = true
}

data "ns1_dnssec" "example" {
  zone = "${ns1_zone.example.zone}"
}
```

## Argument Reference

* `zone` - (Required) The name of the zone to get DNSSEC details for.

## Attributes Reference

In addition to the argument above, the following are exported:

* `keys` - (Computed) - [Keys](#keys-1) field is documented below.
* `delegation` - (Computed) - [Delegation](#delegation-1) field is documented
  below.

#### Keys

`keys` has the following fields:

* `dnskey` - (Computed) List of Keys. [Key](#key) is documented below.
* `ttl` - (Computed) TTL for the Keys (int).

#### Delegation

`delegation` has the following fields:

* `dnskey` - (Computed) List of Keys. [Key](#key) is documented below.
* `ds` - (Computed) List of Keys. [Key](#key) is documented below.
* `ttl` - (Computed) TTL for the Keys (int).

#### Key

A `key` has the following (string) fields:

* `flags` - (Computed) Flags for the key.
* `protocol` - (Computed) Protocol of the key.
* `algorithm` - (Computed) Algorithm of the key.
* `public_key` - (Computed) Public key for the key.
