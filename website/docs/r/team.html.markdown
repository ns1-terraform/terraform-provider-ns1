---
layout: "ns1"
page_title: "NS1: ns1_team"
sidebar_current: "docs-ns1-resource-team"
description: |-
  Provides a NS1 Team resource.
---

# ns1\_team

Provides a NS1 Team resource. This can be used to create, modify, and delete
teams. The credentials used must have the `manage_teams` permission set.

## Example Usage

```hcl
# Create a new NS1 Team
resource "ns1_team" "example" {
  name = "Example team"

    
  # Optional IP whitelists
  ip_whitelist {
    name = "whitelist-1"
    values = ["1.1.1.1", "2.2.2.2"]
  }
  ip_whitelist {
    name = "whitelist-2"
    values = ["3.3.3.3", "4.4.4.4"]
  }

  # Configure permissions
  dns_view_zones       = false
  account_manage_users = false
}

# Another team
resource "ns1_team" "example2" {
  name = "another team"

  dns_view_zones = true
  dns_zones_allow_by_default = true
  dns_zones_allow = ["mytest.zone"]
  dns_zones_deny = ["myother.zone"]

  data_manage_datasources = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the team.
* `ip_whitelist` - (Optional) The IP addresses to whitelist for this key.
* `dns_view_zones` - (Optional) Whether the team can view the accounts zones.
* `dns_manage_zones` - (Optional) Whether the team can modify the accounts zones.
* `dns_zones_allow_by_default` - (Optional) If true, enable the `dns_zones_allow` list, otherwise enable the `dns_zones_deny` list.
* `dns_zones_allow` - (Optional) List of zones that the team may access.
* `dns_zones_deny` - (Optional) List of zones that the team may not access.
* `data_push_to_datafeeds` - (Optional) Whether the team can publish to data feeds.
* `data_manage_datasources` - (Optional) Whether the team can modify data sources.
* `data_manage_datafeeds` - (Optional) Whether the team can modify data feeds.
* `account_manage_users` - (Optional) Whether the team can modify account users.
* `account_manage_payment_methods` - (Optional) Whether the team can modify account payment methods.
* `account_manage_plan` - (Optional) Whether the team can modify the account plan.
* `account_manage_teams` - (Optional) Whether the team can modify other teams in the account.
* `account_manage_apikeys` - (Optional) Whether the team can modify account apikeys.
* `account_manage_account_settings` - (Optional) Whether the team can modify account settings.
* `account_view_activity_log` - (Optional) Whether the team can view activity logs.
* `account_view_invoices` - (Optional) Whether the team can view invoices.
* `monitoring_manage_lists` - (Optional) Whether the team can modify notification lists.
* `monitoring_manage_jobs` - (Optional) Whether the team can modify monitoring jobs.
* `monitoring_view_jobs` - (Optional) Whether the team can view monitoring jobs.
* `security_manage_global_2fa` - (Optional) Whether the team can manage global two factor authentication.
* `security_manage_active_directory` - (Optional) Whether the team can manage global active directory.
Only relevant for the DDI product.
* `dhcp_manage_dhcp` - (Optional) Whether the team can manage DHCP.
Only relevant for the DDI product.
* `dhcp_view_dhcp` - (Optional) Whether the team can view DHCP.
Only relevant for the DDI product.
* `ipam_manage_ipam` - (Optional) Whether the team can manage IPAM.
Only relevant for the DDI product.
* `ipam_view_ipam` - (Optional) Whether the team can view IPAM.
Only relevant for the DDI product.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[Team Api Docs](https://ns1.com/api#team)
