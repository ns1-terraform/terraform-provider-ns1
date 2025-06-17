## Teams

resource "ns1_team" "example_team" {
  name = "Example Team"
  dns_manage_zones = false
  dns_view_zones = true
  dns_zones_allow_by_default = false
  dns_zones_allow = ["example.com"]
# ....... all permissions listed on G
  account_manage_ip_whitelist = true
  monitoring_manage_lists     = true
  monitoring_manage_jobs      = true
  monitoring_view_jobs        = true
  security_manage_global_2fa  = true
}


## Users

# Team User
resource "ns1_user" "example_team_user" {
  name      = "Example Team User from Terraform"
  username  = "example-team-user"
  email     = "testuser@example.com"
  teams     = [ns1_team.example_team.id]
  notify = {
    "billing" = false
  }
}

# Admin User
resource "ns1_user" "example_admin_user" {
  name      = "Example Admin User from Terraform"
  username  = "example-admin-user"
  email     = "testuser@example.com"
  notify = {
    "billing" = true
  }
  dns_view_zones                  = true
  dns_manage_zones                = true
  dns_zones_allow_by_default      = true
  data_push_to_datafeeds          = true
  data_manage_datasources         = true
  data_manage_datafeeds           = true
  account_manage_users            = true
  account_manage_payment_methods  = true
  account_manage_teams            = true
  account_manage_apikeys          = true
  account_manage_account_settings = true
  account_view_activity_log       = true
  account_view_invoices           = true
  account_manage_ip_whitelist     = true
  monitoring_manage_lists         = true
  monitoring_manage_jobs          = true
  monitoring_view_jobs            = true
  security_manage_global_2fa      = true
}

# Read only user with IP whitelist
resource "ns1_user" "example_whitelist_user" {
  name      = "Example User with Whitelist from Terraform"
  username  = "example-admin-user"
  email     = "testuser@example.com"
  notify = {
    "billing" = true
  }
  ip_whitelist                    = ["1.1.1.1","2.2.2.2"]
  dns_view_zones                  = true
  dns_manage_zones                = false
  dns_zones_allow_by_default      = false
  data_push_to_datafeeds          = false
  data_manage_datasources         = false
  data_manage_datafeeds           = false
  account_manage_users            = false
  account_manage_payment_methods  = false
  account_manage_teams            = false
  account_manage_apikeys          = false
  account_manage_account_settings = false
  account_view_activity_log       = false
  account_view_invoices           = false
  account_manage_ip_whitelist     = false
  monitoring_manage_lists         = false
  monitoring_manage_jobs          = false
  monitoring_view_jobs            = true
  security_manage_global_2fa      = false
}


## API keys
resource "ns1_apikey" "example" {
  name      = "Example API Key from Terraform"
  monitoring_manage_lists         = true
  monitoring_manage_jobs          = true
  monitoring_view_jobs            = true
  security_manage_global_2fa      = true
}


## Permissions Arguments
#name - (Required) The free form name of the user.
#username - (Required) The users login name.
#email - (Required) The email address of the user.
#notify - (Required) Whether or not to notify the user of specified events. Only billing is available currently.
#teams - (Required) The teams that the user belongs to.
#ip_whitelist - (Optional) Array of IP addresses/networks to which to grant the user access.
#ip_whitelist_strict - (Optional) Set to true to restrict access to only those IP addresses and networks listed in the ip_whitelist field.
#dns_view_zones - (Optional) Whether the user can view the accounts zones.
#dns_manage_zones - (Optional) Whether the user can modify the accounts zones.
#dns_zones_allow_by_default - (Optional) If true, enable the dns_zones_allow list, otherwise enable the dns_zones_deny list.
#dns_zones_allow - (Optional) List of zones that the user may access.
#dns_zones_deny - (Optional) List of zones that the user may not access.
#data_push_to_datafeeds - (Optional) Whether the user can publish to data feeds.
#data_manage_datasources - (Optional) Whether the user can modify data sources.
#data_manage_datafeeds - (Optional) Whether the user can modify data feeds.
#account_manage_users - (Optional) Whether the user can modify account users.
#account_manage_payment_methods - (Optional) Whether the user can modify account payment methods.
#account_manage_plan - (Deprecated) No longer in use.
#account_manage_teams - (Optional) Whether the user can modify other teams in the account.
#account_manage_apikeys - (Optional) Whether the user can modify account apikeys.
#account_manage_account_settings - (Optional) Whether the user can modify account settings.
#account_view_activity_log - (Optional) Whether the user can view activity logs.
#account_view_invoices - (Optional) Whether the user can view invoices.
#account_manage_ip_whitelist - (Optional) Whether the user can manage ip whitelist.
#monitoring_manage_lists - (Optional) Whether the user can modify notification lists.
#monitoring_manage_jobs - (Optional) Whether the user can create, update, and delete monitoring jobs.
#monitoring_create_jobs - (Optional) Whether the user can create monitoring jobs when manage_jobs is not set to true.
#monitoring_update_jobs - (Optional) Whether the user can update monitoring jobs when manage_jobs is not set to true.
#monitoring_delete_jobs - (Optional) Whether the user can delete monitoring jobs when manage_jobs is not set to true.
#monitoring_view_jobs - (Optional) Whether the user can view monitoring jobs.
#security_manage_global_2fa - (Optional) Whether the user can manage global two factor authentication.
#security_manage_active_directory - (Optional) Whether the user can manage global active directory. Only relevant for the DDI product.
#redirects_manage_redirects - (Optional) Whether the user can manager redirects.