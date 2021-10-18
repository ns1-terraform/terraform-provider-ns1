---
layout: "ns1"
page_title: "NS1: ns1_tsigkey"
sidebar_current: "docs-ns1-resource-tsigket"
description: |-
  Provides a NS1 TSIG Key resource.
---

# ns1\_tsigkey

Only supported in S4 and DDI. \
Provides a NS1 TSIG Key resource. This can be used to create, modify, and delete TSIG keys.

## Example Usage

```hcl
resource "ns1_tsigkey" "example" {
  name = "ExampleTsigKey"
  algorithm = "hmac-sha256"
  secret = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA=="
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the tsigkey.
* `algorithm` - (Required) The algorithm used to hash the TSIG key's secret.
* `secret` - (Required) The key's secret to be hashed.

## Import

`terraform import ns1_pulsarjob.importTest <name>`

## NS1 Documentation

[TSIG Keys Api Doc](https://ns1.com/api/#tsig)
