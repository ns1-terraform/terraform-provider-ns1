---
layout: "ns1"
page_title: "NS1: ns1_account_whitelist"
sidebar_current: "docs-ns1-resource-account-whitelist"
description: |-
  Provides a NS1 Global IP Whitelist resource.
---

# ns1\_account\_whitelist

Provides a NS1 Global IP Whitelist resource.

This can be used to create, modify, and delete Global IP Whitelists.

## Example Usage

```hcl
resource "ns1_account_whitelist" "example" {
  name  = "Example Whitelist"
  values = ["1.1.1.1","2.2.2.2"]
}
```

~> You current source IP must be present in one of the whitelists to prevent locking yourself out.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the whitelist.
* `values` - (Required) Array of IP addresses/networks from which to allow access.

## Import

`terraform import ns1_account_whitelist.example <ID>`

## NS1 Documentation

[Global IP Whitelist Doc](https://ns1.com/api?docId=2282)
