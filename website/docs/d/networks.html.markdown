---
layout: "ns1"
page_title: "NS1: ns1_networks"
sidebar_current: "docs-ns1-datasource-networks"
description: |-
  Provides details about NS1 Networks.
---

# Data Source: ns1_networks

Provides details about NS1 Networks. Use this if you would simply like to read
information from NS1 into your configurations. For read/write operations, you
should use a resource.

## Example Usage

```hcl
# Get details about NS1 Networks.
data "ns1_networks" "example" {
}
```

## Argument Reference

There are no required arguments.

## Attributes Reference

The following are attributes exported:

* `networks` - A set of the available networks. [Networks](#networks) is
  documented below.

#### Networks

A network has the following fields:

* `label` - Label associated with the network.
* `name` - Name of the network.
* `network_id` - network ID (`int`). Default is network 0, the primary NS1 Managed DNS Network.
