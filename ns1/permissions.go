package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ns1/ns1-go.v2/common/conv"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func addPermsSchema(s map[string]*schema.Schema) map[string]*schema.Schema {
	dnsRecords := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Required: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain": {
					Type:     schema.TypeString,
					Required: true,
				},
				"include_subdomains": {
					Type:     schema.TypeBool,
					Required: true,
				},
				"zone": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
	s["dns_records_allow"] = dnsRecords
	s["dns_records_deny"] = dnsRecords
	s["dns_view_zones"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dns_manage_zones"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dns_zones_allow_by_default"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dns_zones_deny"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	s["dns_zones_allow"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	s["data_push_to_datafeeds"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["data_manage_datasources"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["data_manage_datafeeds"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_users"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_payment_methods"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_plan"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
		Deprecated:       "obsolete, should no longer be used",
	}
	s["account_manage_teams"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_apikeys"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_account_settings"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_view_activity_log"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_view_invoices"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_ip_whitelist"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_manage_lists"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_manage_jobs"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_view_jobs"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["security_manage_global_2fa"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["security_manage_active_directory"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dhcp_manage_dhcp"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dhcp_view_dhcp"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["ipam_manage_ipam"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["ipam_view_ipam"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	return s
}

// If a user or API key is part of a team then this suppresses the diff on the permissions,
// since it will inherit the permissions of the teams it is on.
func suppressPermissionDiff(k, old, new string, d *schema.ResourceData) bool {
	// Don't want to suppress the diff if the key has no value -- e.g. the first time this is ran
	// otherwise Terraform complains about nil values.
	if old == "" {
		return false
	}

	oldTeams, newTeams := d.GetChange("teams")

	// Check for both old and new team values - if either of them is set,
	// (e.g. if a user is either being added to a team or removed from one),
	// then ignore diffs on the permission keys.
	if teamsRaw, ok := oldTeams.([]interface{}); ok {
		if len(teamsRaw) > 0 {
			return true
		}
	}

	// For some reason, removing a user from a team by completely
	// deleting the teams block from the config will not show up here,
	// the old value will still be in newTeams, so there is no way to know
	// that a user isn't part of any teams anymore. So `terraform apply`
	// has to be ran again to update the users permissions.
	if teamsRaw, ok := newTeams.([]interface{}); ok {
		if len(teamsRaw) > 0 {
			return true
		}
	}

	return false
}

func permissionsToResourceData(d *schema.ResourceData, permissions *account.PermissionsMapV2) {
	if permissions == nil {
		permissions = &account.PermissionsMapV2{}
	}
	if permissions.Account != nil {
		if permissions.Account.ManageUsers != nil {
			d.Set("account_manage_users", conv.BoolFromPtr(permissions.Account.ManageUsers))
		}
		if permissions.Account.ManagePaymentMethods != nil {
			d.Set("account_manage_payment_methods", conv.BoolFromPtr(permissions.Account.ManagePaymentMethods))
		}
		if permissions.Account.ManagePlan != nil {
			d.Set("account_manage_plan", conv.BoolFromPtr(permissions.Account.ManagePlan))
		}
		if permissions.Account.ManageTeams != nil {
			d.Set("account_manage_teams", conv.BoolFromPtr(permissions.Account.ManageTeams))
		}
		if permissions.Account.ManageApikeys != nil {
			d.Set("account_manage_apikeys", conv.BoolFromPtr(permissions.Account.ManageApikeys))
		}
		if permissions.Account.ManageAccountSettings != nil {
			d.Set("account_manage_account_settings", conv.BoolFromPtr(permissions.Account.ManageAccountSettings))
		}
		if permissions.Account.ViewActivityLog != nil {
			d.Set("account_view_activity_log", conv.BoolFromPtr(permissions.Account.ViewActivityLog))
		}
		if permissions.Account.ViewInvoices != nil {
			d.Set("account_view_invoices", conv.BoolFromPtr(permissions.Account.ViewInvoices))
		}
		if permissions.Account.ManageIPWhitelist != nil {
			d.Set("account_manage_ip_whitelist", conv.BoolFromPtr(permissions.Account.ManageIPWhitelist))
		}
	}
	if permissions.DNS != nil {
		if permissions.DNS.ViewZones != nil {
			d.Set("dns_view_zones", conv.BoolFromPtr(permissions.DNS.ViewZones))
		}
		if permissions.DNS.ManageZones != nil {
			d.Set("dns_manage_zones", conv.BoolFromPtr(permissions.DNS.ManageZones))
		}
		if permissions.DNS.ZonesAllowByDefault != nil {
			d.Set("dns_zones_allow_by_default", conv.BoolFromPtr(permissions.DNS.ZonesAllowByDefault))
		}

		d.Set("dns_zones_deny", permissions.DNS.ZonesDeny)
		d.Set("dns_zones_allow", permissions.DNS.ZonesAllow)
	}
	if permissions.Data != nil {
		if permissions.Data.PushToDatafeeds != nil {
			d.Set("data_push_to_datafeeds", conv.BoolFromPtr(permissions.Data.PushToDatafeeds))
		}
		if permissions.Data.ManageDatasources != nil {
			d.Set("data_manage_datasources", conv.BoolFromPtr(permissions.Data.ManageDatasources))
		}
		if permissions.Data.ManageDatafeeds != nil {
			d.Set("data_manage_datafeeds", conv.BoolFromPtr(permissions.Data.ManageDatafeeds))
		}
	}
	if permissions.Monitoring != nil {
		if permissions.Monitoring.ManageLists != nil {
			d.Set("monitoring_manage_lists", conv.BoolFromPtr(permissions.Monitoring.ManageLists))
		}
		if permissions.Monitoring.ManageJobs != nil {
			d.Set("monitoring_manage_jobs", conv.BoolFromPtr(permissions.Monitoring.ManageJobs))
		}
		if permissions.Monitoring.ViewJobs != nil {
			d.Set("monitoring_view_jobs", conv.BoolFromPtr(permissions.Monitoring.ViewJobs))
		}
	}
	if permissions.Security != nil {
		if permissions.Security.ManageGlobal2FA != nil {
			d.Set("security_manage_global_2fa", conv.BoolFromPtr(permissions.Security.ManageGlobal2FA))
		}
		if permissions.Security.ManageActiveDirectory != nil {
			d.Set("security_manage_active_directory", conv.BoolFromPtr(permissions.Security.ManageActiveDirectory))
		}
	}
	if permissions.DHCP != nil {
		if permissions.DHCP.ManageDHCP != nil {
			d.Set("dhcp_manage_dhcp", conv.BoolFromPtr(permissions.DHCP.ManageDHCP))
		}
		if permissions.DHCP.ViewDHCP != nil {
			d.Set("dhcp_view_dhcp", conv.BoolFromPtr(permissions.DHCP.ViewDHCP))
		}
	}
	if permissions.IPAM != nil {
		if permissions.IPAM.ManageIPAM != nil {
			d.Set("ipam_manage_ipam", conv.BoolFromPtr(permissions.IPAM.ManageIPAM))
		}
		if permissions.IPAM.ViewIPAM != nil {
			d.Set("ipam_view_ipam", conv.BoolFromPtr(permissions.IPAM.ViewIPAM))
		}
	}
}

func resourceDataToPermissions(d *schema.ResourceData) account.PermissionsMapV2 {
	p := account.PermissionsMapV2{}
	if p.DNS == nil {
		p.DNS = &account.PermissionsDNSV2{}
	}

	if v, ok := d.GetOkExists("dns_records_allow"); ok {
		p.DNS.RecordsAllow = SchemaToRecordArray(v)
	} else {
		p.DNS.RecordsAllow = []account.PermissionsRecord{}
	}
	if v, ok := d.GetOkExists("dns_records_deny"); ok {
		p.DNS.RecordsDeny = SchemaToRecordArray(v)
	} else {
		p.DNS.RecordsAllow = []account.PermissionsRecord{}
	}
	if v, ok := d.GetOkExists("dns_view_zones"); ok {
		p.DNS.ViewZones = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("dns_manage_zones"); ok {
		p.DNS.ManageZones = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("dns_zones_allow_by_default"); ok {
		p.DNS.ZonesAllowByDefault = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("dns_zones_deny"); ok {
		denyRaw := v.([]interface{})
		p.DNS.ZonesDeny = make([]string, len(denyRaw))
		for i, deny := range denyRaw {
			p.DNS.ZonesDeny[i] = deny.(string)
		}
	} else {
		p.DNS.ZonesDeny = []string{}
	}
	if v, ok := d.GetOkExists("dns_zones_allow"); ok {
		allowRaw := v.([]interface{})
		p.DNS.ZonesAllow = make([]string, len(allowRaw))
		for i, allow := range allowRaw {
			p.DNS.ZonesAllow[i] = allow.(string)
		}
	} else {
		p.DNS.ZonesAllow = []string{}
	}
	if v, ok := d.GetOkExists("data_push_to_datafeeds"); ok {
		if p.Data == nil {
			p.Data = &account.PermissionsDataV2{}
		}
		p.Data.PushToDatafeeds = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("data_manage_datasources"); ok {
		if p.Data == nil {
			p.Data = &account.PermissionsDataV2{}
		}
		p.Data.ManageDatasources = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("data_manage_datafeeds"); ok {
		if p.Data == nil {
			p.Data = &account.PermissionsDataV2{}
		}
		p.Data.ManageDatafeeds = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_users"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManageUsers = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_payment_methods"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManagePaymentMethods = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_plan"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManagePlan = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_teams"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManageTeams = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_apikeys"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManageApikeys = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_account_settings"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManageAccountSettings = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_view_activity_log"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ViewActivityLog = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_view_invoices"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ViewInvoices = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("account_manage_ip_whitelist"); ok {
		if p.Account == nil {
			p.Account = &account.PermissionsAccountV2{}
		}
		p.Account.ManageIPWhitelist = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("monitoring_manage_lists"); ok {
		if p.Monitoring == nil {
			p.Monitoring = &account.PermissionsMonitoringV2{}
		}
		p.Monitoring.ManageLists = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("monitoring_manage_jobs"); ok {
		if p.Monitoring == nil {
			p.Monitoring = &account.PermissionsMonitoringV2{}
		}
		p.Monitoring.ManageJobs = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("monitoring_view_jobs"); ok {
		if p.Monitoring == nil {
			p.Monitoring = &account.PermissionsMonitoringV2{}
		}
		p.Monitoring.ViewJobs = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("security_manage_global_2fa"); ok {
		if p.Security == nil {
			p.Security = &account.PermissionsSecurityV2{}
		}
		p.Security.ManageGlobal2FA = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("security_manage_active_directory"); ok {
		if p.Security == nil {
			p.Security = &account.PermissionsSecurityV2{}
		}
		p.Security.ManageActiveDirectory = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("dhcp_manage_dhcp"); ok {
		if p.DHCP == nil {
			p.DHCP = &account.PermissionsDHCPV2{}
		}
		p.DHCP.ManageDHCP = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("dhcp_view_dhcp"); ok {
		if p.DHCP == nil {
			p.DHCP = &account.PermissionsDHCPV2{}
		}
		p.DHCP.ViewDHCP = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("ipam_manage_ipam"); ok {
		if p.IPAM == nil {
			p.IPAM = &account.PermissionsIPAMV2{}
		}
		p.IPAM.ManageIPAM = conv.BoolPtrFrom(v.(bool))
	}
	if v, ok := d.GetOkExists("ipam_view_ipam"); ok {
		if p.IPAM == nil {
			p.IPAM = &account.PermissionsIPAMV2{}
		}
		p.IPAM.ViewIPAM = conv.BoolPtrFrom(v.(bool))
	}

	return p
}

func SchemaToRecordArray(v interface{}) []account.PermissionsRecord {
	if schemaRecord, ok := v.([]interface{}); ok {
		var records []account.PermissionsRecord
		for _, sr := range schemaRecord {
			mapRecord := sr.(map[string]interface{})
			record := account.PermissionsRecord{
				Domain:     mapRecord["domain"].(string),
				Subdomains: mapRecord["include_subdomains"].(bool),
				Zone:       mapRecord["zone"].(string),
				RecordType: mapRecord["type"].(string)}
			records = append(records, record)
		}
		return records
	}
	return []account.PermissionsRecord{}
}
