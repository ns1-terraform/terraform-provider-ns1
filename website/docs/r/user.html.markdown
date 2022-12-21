---
layout: "ns1"
page_title: "NS1: ns1_user"
sidebar_current: "docs-ns1-resource-user"
description: |-
  Provides a NS1 User resource.
---

# ns1\_user

Provides a NS1 User resource. Creating a user sends an invitation email to the
user's email address. This can be used to create, modify, and delete users.
The credentials used must have the `manage_users` permission set.

## Example Usage

```hcl
resource "ns1_team" "example" {
  name = "Example team"

  # Optional IP whitelist
  ip_whitelist = ["1.1.1.1","2.2.2.2"]

  dns_view_zones       = false
  account_manage_users = false
}

resource "ns1_user" "example" {
  name      = "Example User"
  username  = "example_user"
  email     = "user@example.com"
  teams     = [ns1_team.example.id]
  notify = {
    billing = false
  }
}
```

## Permissions
A user will inherit permissions from the teams they are assigned to.
If a user is assigned to a team and also has individual permissions set on the user, the individual permissions
will be overridden by the inherited team permissions.
In a future release, setting permissions on a user that is part of a team will be explicitly disabled.

When a user is removed from all teams completely, they will inherit whatever permissions they had previously.
If a user is removed from all their teams, it will probably be necessary to run `terraform apply` a second time
to update the users permissions from their old team permissions to new user-specific permissions.

See [this NS1 Help Center article](https://help.ns1.com/hc/en-us/articles/360024409034-Managing-user-permissions) for an overview of user permission settings.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the user.
* `username` - (Required) The users login name.
* `email` - (Required) The email address of the user.
* `notify` - (Required) Whether or not to notify the user of specified events. Only `billing` is available currently.
* `teams` - (Required) The teams that the user belongs to.
* `ip_whitelist` - (Optional) Array of IP addresses/networks to which to grant the user access. 
* `ip_whitelist_strict` - (Optional) Set to true to restrict access to only those IP addresses and networks listed in the **ip_whitelist** field.
* `dns_view_zones` - (Optional) Whether the user can view the accounts zones.
* `dns_manage_zones` - (Optional) Whether the user can modify the accounts zones.
* `dns_zones_allow_by_default` - (Optional) If true, enable the `dns_zones_allow` list, otherwise enable the `dns_zones_deny` list.
* `dns_zones_allow` - (Optional) List of zones that the user may access.
* `dns_zones_deny` - (Optional) List of zones that the user may not access.
* `data_push_to_datafeeds` - (Optional) Whether the user can publish to data feeds.
* `data_manage_datasources` - (Optional) Whether the user can modify data sources.
* `data_manage_datafeeds` - (Optional) Whether the user can modify data feeds.
* `account_manage_users` - (Optional) Whether the user can modify account users.
* `account_manage_payment_methods` - (Optional) Whether the user can modify account payment methods.
* `account_manage_plan` - (Deprecated) No longer in use.
* `account_manage_teams` - (Optional) Whether the user can modify other teams in the account.
* `account_manage_apikeys` - (Optional) Whether the user can modify account apikeys.
* `account_manage_account_settings` - (Optional) Whether the user can modify account settings.
* `account_view_activity_log` - (Optional) Whether the user can view activity logs.
* `account_view_invoices` - (Optional) Whether the user can view invoices.
* `account_manage_ip_whitelist` - (Optional) Whether the user can manage ip whitelist.
* `monitoring_manage_lists` - (Optional) Whether the user can modify notification lists.
* `monitoring_manage_jobs` - (Optional) Whether the user can modify monitoring jobs.
* `monitoring_view_jobs` - (Optional) Whether the user can view monitoring jobs.
* `security_manage_global_2fa` - (Optional) Whether the user can manage global two factor authentication.
* `security_manage_active_directory` - (Optional) Whether the user can manage global active directory.
Only relevant for the DDI product.
* `dhcp_manage_dhcp` - (Optional) Whether the user can manage DHCP.
Only relevant for the DDI product.
* `dhcp_view_dhcp` - (Optional) Whether the user can view DHCP.
Only relevant for the DDI product.
* `ipam_manage_ipam` - (Optional) Whether the user can manage IPAM.
Only relevant for the DDI product.

## Import

`terraform import ns1_user.<name> <username>`

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[User Api Docs](https://ns1.com/api#user)

[Managing user permissions](https://help.ns1.com/hc/en-us/articles/360024409034-Managing-user-permissions)

