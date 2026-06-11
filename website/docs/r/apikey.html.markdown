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

### Legacy API Key (Static Secret)

```hcl
resource "ns1_apikey" "static_key" {
  name = "Static API Key"

  # Configure permissions
  dns_view_zones = true
  dns_manage_zones = true
}

# The static secret is available in the key attribute
output "api_key_secret" {
  value     = ns1_apikey.static_key.key
  sensitive = true
}
```

### API Key with Secret Expiration

```hcl
resource "ns1_apikey" "expiring" {
  name = "Expiring API Key"

  # Set secret expiration period to 30 days
  # Secrets will expire after this period and must be manually rotated
  # Accepts any duration in '<number>d' format (e.g., "10d", "30d", "90d")
  expiry_duration = "30d"

  # Configure permissions
  dns_view_zones  = true
  dns_manage_zones = true
}

# Access secret metadata (not the actual secret key values)
output "secret_info" {
  value = ns1_apikey.expiring.secrets
  sensitive = true
}
```

## Important Notes

~> **Changing expiry_duration forces recreation.** When you modify the `expiry_duration` field of an existing API key, Terraform will destroy the old key and create a new one. This means the API key ID and all secrets will change. Any external references to the old key will break. Plan your migrations carefully and update dependent systems before changing this value.

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
* `ip_whitelist` - (Optional, default: `[]`) Array of IP addresses/networks to which to grant the API key access.
* `ip_whitelist_strict` - (Optional, default: `false`) Set to true to restrict access to only those IP addresses and networks listed in the **ip_whitelist** field.
* `expiry_duration` - (Optional) Duration for secret expiration in `<number>d` format (e.g., `"10d"`, `"30d"`, `"90d"`). When set, API key secrets will expire after the specified period and must be manually rotated using the NS1 API or Portal. The API key can have up to 2 active secrets at a time to allow for graceful rotation without service interruption. If not set, a legacy API key with a permanent secret (stored in the `key` attribute) is created. Changing this value will force recreation of the API key.
* `dns_view_zones` - (Optional, default: `false`) Whether the apikey can view the accounts zones.
* `dns_manage_zones` - (Optional, default: `false`) Whether the apikey can modify the accounts zones.
* `dns_zones_allow_by_default` - (Optional, default: `false`) If true, enable the `dns_zones_allow` list, otherwise enable the `dns_zones_deny` list.
* `dns_zones_allow` - (Optional, default: `[]`) List of zones that the apikey may access.
* `dns_zones_deny` - (Optional, default: `[]`) List of zones that the apikey may not access.
* `dns_records_allow` - (Optional, default: `[]`) List of records that the apikey may access.
* `dns_records_deny` - (Optional, default: `[]`) List of records that the apikey may not access.
* `data_push_to_datafeeds` - (Optional, default: `false`) Whether the apikey can publish to data feeds.
* `data_manage_datasources` - (Optional, default: `false`) Whether the apikey can modify data sources.
* `data_manage_datafeeds` - (Optional, default: `false`) Whether the apikey can modify data feeds.
* `account_manage_users` - (Optional, default: `false`) Whether the apikey can modify account users.
* `account_manage_payment_methods` - (Optional, default: `false`) Whether the apikey can modify account payment methods.
* `account_manage_plan` - (Deprecated) No longer in use.
* `account_manage_teams` - (Optional, default: `false`) Whether the apikey can modify other teams in the account.
* `account_manage_apikeys` - (Optional, default: `false`) Whether the apikey can modify account apikeys.
* `account_manage_account_settings` - (Optional, default: `false`) Whether the apikey can modify account settings.
* `account_view_activity_log` - (Optional, default: `false`) Whether the apikey can view activity logs.
* `account_view_invoices` - (Optional), default: `false` Whether the apikey can view invoices.
* `account_manage_ip_whitelist` - (Optional, default: `false`) Whether the apikey can manage ip whitelist.
* `monitoring_manage_lists` - (Optional, default: `false`) Whether the apikey can modify notification lists.
* `monitoring_manage_jobs` - (Optional, default: `false`) Whether the apikey can create, update, and delete monitoring jobs.
* `monitoring_create_jobs` - (Optional, default: `false`) Whether the apikey can create monitoring jobs when manage_jobs is not set to true.
* `monitoring_update_jobs` - (Optional, default: `false`) Whether the apikey can update monitoring jobs when manage_jobs is not set to true.
* `monitoring_delete_jobs` - (Optional, default: `false`) Whether the apikey can delete monitoring jobs when manage_jobs is not set to true.
* `monitoring_view_jobs` - (Optional, default: `false`) Whether the apikey can view monitoring jobs.
* `security_manage_global_2fa` - (Optional, default: `true`) Whether the apikey can manage global two factor authentication.
* `security_manage_active_directory` - (Optional, default: `true`) Whether the apikey can manage global active directory. Only relevant for the DDI product.
* `redirects_manage_redirects` - (Optional, default: `false`) Whether the apikey can manage redirects.
* `insights_view_insights` - (Optional, default: `false`) Whether the apikey can view DNS insights.
* `insights_manage_insights` - (Optional, default: `false`) Whether the apikey can manage DNS insights.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key` - (Computed) The API key authentication token. Only populated for legacy API keys (when `expiry_duration` is not set). For API keys with expiration, use the secret keys from the `secrets` attribute instead.
* `secrets` - (Computed) List of secrets for this API key. Only populated when `expiry_duration` is set. Each secret contains:
  * `id` - The unique identifier for the secret.
  * `expires_at` - The expiration date/time of the secret in ISO 8601 format.
  * `last_access` - The last time this secret was used for authentication.
  * `enabled` - Whether this secret is currently enabled for authentication.

**Note:** The actual secret key values (starting with `nss_`) are only returned when a secret is first created and are not stored in Terraform state for security reasons. You must save these values when they are first created, as they cannot be retrieved later. To rotate secrets (generate new ones), use the NS1 API or Portal - Terraform does not manage secret rotation.

## Import

-> Imported keys will not have their key stored in the state file.

`terraform import ns1_apikey`

So for the example above:

`terraform import ns1_apikey.example <ID>`


## NS1 Documentation

[ApiKeys Api Doc](https://ns1.com/api#api-key)
