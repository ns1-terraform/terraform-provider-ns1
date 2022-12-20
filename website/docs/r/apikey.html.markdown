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
  teams = [ns1_team.example.id]

  # Optional IP whitelist
  ip_whitelist = ["1.1.1.1","2.2.2.2"]

  # Configure permissions 
  dns_view_zones       = false
  account_manage_users = false
}
```

## Permissions
An API key will inherit permissions from the teams it is assigned to.
If a key is assigned to a team and also has individual permissions set on the key, the individual permissions
will be overridden by the inherited team permissions.
In a future release, setting permissions on a key that is part of a team will be explicitly disabled.

When a key is removed from all teams completely, it will inherit whatever permissions it had previously.
If a key is removed from all it's teams, it will probably be necessary to run `terraform apply` a second time
to update the keys permissions from it's old team permissions to new key-specific permissions.

See [the NS1 API docs](https://ns1.com/api#getget-all-account-users) for an overview of permission semantics or for [more details](https://help.ns1.com/hc/en-us/articles/360024409034-Managing-user-permissions) about the individual permission flags.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the apikey.
* `teams` - (Optional) The teams that the apikey belongs to.
* `ip_whitelist` - (Optional) Array of IP addresses/networks to which to grant the API key access.
* `ip_whitelist_strict` - (Optional) Set to true to restrict access to only those IP addresses and networks listed in the **ip_whitelist** field.
* `dns_view_zones` - (Optional) Whether the apikey can view the accounts zones.
* `dns_manage_zones` - (Optional) Whether the apikey can modify the accounts zones.
* `dns_zones_allow_by_default` - (Optional) If true, enable the `dns_zones_allow` list, otherwise enable the `dns_zones_deny` list.
* `dns_zones_allow` - (Optional) List of zones that the apikey may access.
* `dns_zones_deny` - (Optional) List of zones that the apikey may not access.
* `dns_records_allow` - (Optional) List of records that the apikey may access.
* `dns_records_deny` - (Optional) List of records that the apikey may not access.
* `data_push_to_datafeeds` - (Optional) Whether the apikey can publish to data feeds.
* `data_manage_datasources` - (Optional) Whether the apikey can modify data sources.
* `data_manage_datafeeds` - (Optional) Whether the apikey can modify data feeds.
* `account_manage_users` - (Optional) Whether the apikey can modify account users.
* `account_manage_payment_methods` - (Optional) Whether the apikey can modify account payment methods.
* `account_manage_plan` - (Deprecated) No longer in use.
* `account_manage_teams` - (Optional) Whether the apikey can modify other teams in the account.
* `account_manage_apikeys` - (Optional) Whether the apikey can modify account apikeys.
* `account_manage_account_settings` - (Optional) Whether the apikey can modify account settings.
* `account_view_activity_log` - (Optional) Whether the apikey can view activity logs.
* `account_view_invoices` - (Optional) Whether the apikey can view invoices.
* `account_manage_ip_whitelist` - (Optional) Whether the apikey can manage ip whitelist.
* `monitoring_manage_lists` - (Optional) Whether the apikey can modify notification lists.
* `monitoring_manage_jobs` - (Optional) Whether the apikey can modify monitoring jobs.
* `monitoring_view_jobs` - (Optional) Whether the apikey can view monitoring jobs.
* `security_manage_global_2fa` - (Optional) Whether the apikey can manage global two factor authentication.
* `security_manage_active_directory` - (Optional) Whether the apikey can manage global active directory.
Only relevant for the DDI product.
* `dhcp_manage_dhcp` - (Optional) Whether the apikey can manage DHCP.
Only relevant for the DDI product.
* `dhcp_view_dhcp` - (Optional) Whether the apikey can view DHCP.
Only relevant for the DDI product.
* `ipam_manage_ipam` - (Optional) Whether the apikey can manage IPAM.
Only relevant for the DDI product.
* `ipam_view_ipam` - (Optional) Whether the apikey can view IPAM.
Only relevant for the DDI product.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key` - (Computed) The apikeys authentication token.

## Import

`terraform import ns1_apikey`

So for the example above:

`terraform import ns1_apikey.example <ID>`


## NS1 Documentation

[ApiKeys Api Doc](https://ns1.com/api#api-key)
