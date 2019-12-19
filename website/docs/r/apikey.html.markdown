---
layout: "ns1"
page_title: "NS1: ns1_apikey"
sidebar_current: "docs-ns1-resource-apikey"
description: |-
  Provides a NS1 Api Key resource.
---

# ns1\_apikey

Provides a NS1 Api Key resource. This can be used to create, modify, and delete api keys.

## Example Usage

```hcl
resource "ns1_team" "example" {
  name = "Example team"
}

resource "ns1_apikey" "example" {
  name  = "Example key"
  teams = ["${ns1_team.example.id}"]

  # Configure permissions 
  dns_view_zones       = false
  account_manage_users = false
}
```

## Permissions
An API key will inherit permissions from the teams it is assigned to.
When a key is removed from all teams completely, it will inherit whatever permissions it had previously.
If a key is removed from all it's teams, it will probably be necessary to run `terraform apply` a second time
to update the keys permissions from it's old team permissions to new key-specific permissions.
See [the NS1 API docs](https://ns1.com/api#getget-all-account-users) for an overview of permission semantics.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the apikey.
* `key` - (Required) The apikeys authentication token.
* `teams` - (Optional) The teams that the apikey belongs to.
* `dns_view_zones` - (Optional) Whether the apikey can view the accounts zones.
* `dns_manage_zones` - (Optional) Whether the apikey can modify the accounts zones.
* `dns_zones_allow_by_default` - (Optional) If true, enable the `dns_zones_allow` list, otherwise enable the `dns_zones_deny` list.
* `dns_zones_allow` - (Optional) List of zones that the apikey may access.
* `dns_zones_deny` - (Optional) List of zones that the apikey may not access.
* `data_push_to_datafeeds` - (Optional) Whether the apikey can publish to data feeds.
* `data_manage_datasources` - (Optional) Whether the apikey can modify data sources.
* `data_manage_datafeeds` - (Optional) Whether the apikey can modify data feeds.
* `account_manage_users` - (Optional) Whether the apikey can modify account users.
* `account_manage_payment_methods` - (Optional) Whether the apikey can modify account payment methods.
* `account_manage_plan` - (Optional) Whether the apikey can modify the account plan.
* `account_manage_teams` - (Optional) Whether the apikey can modify other teams in the account.
* `account_manage_apikeys` - (Optional) Whether the apikey can modify account apikeys.
* `account_manage_account_settings` - (Optional) Whether the apikey can modify account settings.
* `account_view_activity_log` - (Optional) Whether the apikey can view activity logs.
* `account_view_invoices` - (Optional) Whether the apikey can view invoices.
* `monitoring_manage_lists` - (Optional) Whether the apikey can modify notification lists.
* `monitoring_manage_jobs` - (Optional) Whether the apikey can modify monitoring jobs.
* `monitoring_view_jobs` - (Optional) Whether the apikey can view monitoring jobs.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.
