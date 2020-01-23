package ns1

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func permissionInstanceStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["security_manage_global_2fa"] = false
	rawState["security_manage_active_directory"] = false
	rawState["dhcp_manage_dhcp"] = false
	rawState["dhcp_view_dhcp"] = false
	rawState["ipam_manage_ipam"] = false
	rawState["ipam_view_ipam"] = false

	return rawState, nil
}

func addPermsSchemaV0(s map[string]*schema.Schema) map[string]*schema.Schema {
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
