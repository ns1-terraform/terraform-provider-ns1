---
layout: "ns1"
page_title: "NS1: ns1_dnsview"
sidebar_current: "docs-ns1-resource-dns-view"
description: |-
  Provides a NS1 DNS View resource.
---

# dns\_view

Provides a NS1 DNS View resource. This can be used to create, modify, delete and import DNS views.

## Example Usage

```hcl
# Create a new DNS View
resource "ns1_dnsview" "it" {
  name = "terraform_example"
  preference = 1
  read_acls = ["acl_example_1", "acl_example_2"]
  update_acls = ["acl_example_1", "acl_example_2"]
  zones = ["terraform.example.io"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS view.
* `preference` - (Optional/Computed) A unique value that indicates the priority of this view. This can be any value grather than 0 where a value of 1 indicates top priority.
* `read_acls` - (Optional) List of ACL names used with "read" permissions. The order of the ACLs determine how they are processed.
* `update_acls` - (Optional) List of ACL names used with "update" permissions. The order of the ACLs determine how they are processed.
* `zones` - (Optional) List of zone names.
* `networks` - (Optional) Networks is an array of positive integers corresponding to the service definition(s) to which the zone is published. 

## Import

`terraform import ns1_dnsview.<name> <viewName>`

So for the example above:

`terraform import ns1_dnsview.example terraform_example`

## NS1 Documentation

[DNS View Api Docs](https://ns1.com/api/#dns-views-ddi-only)
