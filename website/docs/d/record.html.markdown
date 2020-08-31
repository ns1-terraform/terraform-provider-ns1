---
layout: "ns1"
page_title: "NS1: ns1_record"
sidebar_current: "docs-ns1-datasource-record"
description: |-
  Provides details about a NS1 Record.
---

# Data Source: ns1_record

Provides details about a NS1 Record. Use this if you would simply like to read
information from NS1 into your configurations. For read/write operations, you
should use a resource.

## Example Usage

```hcl
# Get details about a NS1 Record.
data "ns1_record" "example" {
  zone   = "example.io"
  domain = "terraform.example.io"
  type   = "A"
}
```

## Argument Reference

* `zone` - (Required) The zone the record belongs to.
* `domain` - (Required) The records' domain.
* `type` - (Required) The records' RR type.

## Attributes Reference

In addition to the argument above, the following are exported:

* `ttl` - The records' time to live (in seconds).
* `link` - The target record this links to.
* `use_client_subnet` - Whether to use EDNS client subnet data when available (in filter chain).
* `meta` - Map of metadata
* `regions` - List of regions.
* `answers` - List of NS1 answers.
* `filters` - List of NS1 filters.
