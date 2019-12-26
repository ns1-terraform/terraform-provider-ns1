package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func addPermsSchema(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["dns_view_zones"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dns_manage_zones"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["dns_zones_allow_by_default"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
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
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["data_manage_datasources"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["data_manage_datafeeds"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_users"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_payment_methods"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_plan"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_teams"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_apikeys"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_manage_account_settings"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_view_activity_log"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["account_view_invoices"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_manage_lists"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_manage_jobs"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
		DiffSuppressFunc: suppressPermissionDiff,
	}
	s["monitoring_view_jobs"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          false,
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

func permissionsToResourceData(d *schema.ResourceData, permissions account.PermissionsMap) {
	d.Set("dns_view_zones", permissions.DNS.ViewZones)
	d.Set("dns_manage_zones", permissions.DNS.ManageZones)
	d.Set("dns_zones_allow_by_default", permissions.DNS.ZonesAllowByDefault)
	d.Set("dns_zones_deny", permissions.DNS.ZonesDeny)
	d.Set("dns_zones_allow", permissions.DNS.ZonesAllow)
	d.Set("data_push_to_datafeeds", permissions.Data.PushToDatafeeds)
	d.Set("data_manage_datasources", permissions.Data.ManageDatasources)
	d.Set("data_manage_datafeeds", permissions.Data.ManageDatafeeds)
	d.Set("account_manage_users", permissions.Account.ManageUsers)
	d.Set("account_manage_payment_methods", permissions.Account.ManagePaymentMethods)
	d.Set("account_manage_plan", permissions.Account.ManagePlan)
	d.Set("account_manage_teams", permissions.Account.ManageTeams)
	d.Set("account_manage_apikeys", permissions.Account.ManageApikeys)
	d.Set("account_manage_account_settings", permissions.Account.ManageAccountSettings)
	d.Set("account_view_activity_log", permissions.Account.ViewActivityLog)
	d.Set("account_view_invoices", permissions.Account.ViewInvoices)
	d.Set("monitoring_manage_lists", permissions.Monitoring.ManageLists)
	d.Set("monitoring_manage_jobs", permissions.Monitoring.ManageJobs)
	d.Set("monitoring_view_jobs", permissions.Monitoring.ViewJobs)
}

func resourceDataToPermissions(d *schema.ResourceData) account.PermissionsMap {
	var p account.PermissionsMap
	if v, ok := d.GetOk("dns_view_zones"); ok {
		p.DNS.ViewZones = v.(bool)
	}
	if v, ok := d.GetOk("dns_manage_zones"); ok {
		p.DNS.ManageZones = v.(bool)
	}
	if v, ok := d.GetOk("dns_zones_allow_by_default"); ok {
		p.DNS.ZonesAllowByDefault = v.(bool)
	}
	if v, ok := d.GetOk("dns_zones_deny"); ok {
		denyRaw := v.([]interface{})
		p.DNS.ZonesDeny = make([]string, len(denyRaw))
		for i, deny := range denyRaw {
			p.DNS.ZonesDeny[i] = deny.(string)
		}
	} else {
		p.DNS.ZonesDeny = make([]string, 0)
	}
	if v, ok := d.GetOk("dns_zones_allow"); ok {
		allowRaw := v.([]interface{})
		p.DNS.ZonesAllow = make([]string, len(allowRaw))
		for i, allow := range allowRaw {
			p.DNS.ZonesAllow[i] = allow.(string)
		}
	} else {
		p.DNS.ZonesAllow = make([]string, 0)
	}
	if v, ok := d.GetOk("data_push_to_datafeeds"); ok {
		p.Data.PushToDatafeeds = v.(bool)
	}
	if v, ok := d.GetOk("data_manage_datasources"); ok {
		p.Data.ManageDatasources = v.(bool)
	}
	if v, ok := d.GetOk("data_manage_datafeeds"); ok {
		p.Data.ManageDatafeeds = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_users"); ok {
		p.Account.ManageUsers = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_payment_methods"); ok {
		p.Account.ManagePaymentMethods = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_plan"); ok {
		p.Account.ManagePlan = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_teams"); ok {
		p.Account.ManageTeams = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_apikeys"); ok {
		p.Account.ManageApikeys = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_account_settings"); ok {
		p.Account.ManageAccountSettings = v.(bool)
	}
	if v, ok := d.GetOk("account_view_activity_log"); ok {
		p.Account.ViewActivityLog = v.(bool)
	}
	if v, ok := d.GetOk("account_view_invoices"); ok {
		p.Account.ViewInvoices = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_manage_lists"); ok {
		p.Monitoring.ManageLists = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_manage_jobs"); ok {
		p.Monitoring.ManageJobs = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_view_jobs"); ok {
		p.Monitoring.ViewJobs = v.(bool)
	}
	return p
}
