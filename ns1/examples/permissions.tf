# permissions are a shared collections of values that can be applied to many different resources

# this exists to document those resources and the available permissions

# all permissions on these types are optional

# resources that support permissions:
#   ns1_apikey
#   ns1_team
#   ns1_user

# permissions values and types:
#   dns_view_zones: boolean - allows the requestor to view zones
#   dns_manage_zones: boolean - allows the requestor to edit/manage zones
#   dns_zones_allow_by_default: boolean
#   dns_zones_deny: list of strings - explicitly deny these zones for this user/team/key
#   dns_zones_allow: list of strings - explicitly allow these zones for this user/team/key
#   data_push_to_datafeeds: boolean - allows the requestor to push to datafeeds
#   data_manage_datasources: boolean - allows the requestor to manage datasources
#   data_manage_datafeeds: boolean - allows the requestor to manage datafeeds
#   account_manage_users: boolean - allows the requstor to manage users
#   account_manage_payment_methods: boolean - allows the requestor to manage payment methods
#   account_manage_teams: boolean - allows the requestor to manage teams
#   account_manage_apikeys: boolean - allows the requestor to manage apikeys
#   account_manage_account_settings: boolean - allows the requestor to manage account settings
#   account_view_activity_log: boolean - allows the requestor to view the activity log
#   account_view_invoices: boolean - allows the requestor to view account invoices
#   monitoring_manage_lists: boolean - allows the requestor to manage monitoring lists
#   monitoring_manage_jobs: boolean - allows the requestor to manage monitoring jobs
#   monitoring_view_jobs: boolean - allows the requestor to view monitoring jobs
