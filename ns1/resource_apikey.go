package ns1

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func apikeyResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"teams": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"ip_whitelist": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"ip_whitelist_strict": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	s = addPermsSchema(s)

	return &schema.Resource{
		Schema:        s,
		Create:        ApikeyCreate,
		Read:          ApikeyRead,
		Update:        ApikeyUpdate,
		Delete:        ApikeyDelete,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    apikeyResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: permissionInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func apikeyToResourceData(d *schema.ResourceData, k *account.APIKey) error {
	d.SetId(k.ID)
	d.Set("name", k.Name)
	d.Set("key", k.Key)
	d.Set("teams", k.TeamIDs)
	d.Set("ip_whitelist", k.IPWhitelist)
	d.Set("ip_whitelist_strict", k.IPWhitelistStrict)
	permissionsToResourceData(d, k.Permissions)
	return nil
}

func resourceDataToApikey(k *account.APIKey, d *schema.ResourceData) error {
	k.ID = d.Id()
	k.Name = d.Get("name").(string)
	if v, ok := d.GetOk("teams"); ok {
		teamsRaw := v.([]interface{})
		k.TeamIDs = make([]string, len(teamsRaw))
		for i, team := range teamsRaw {
			k.TeamIDs[i] = team.(string)
		}
	} else {
		k.TeamIDs = make([]string, 0)
	}
	k.Permissions = resourceDataToPermissions(d)

	if v, ok := d.GetOk("ip_whitelist"); ok {

		ipWhitelistRaw := v.(*schema.Set)
		k.IPWhitelist = make([]string, ipWhitelistRaw.Len())
		for i, ip := range ipWhitelistRaw.List() {
			k.IPWhitelist[i] = ip.(string)
		}
	} else {
		// This still needs to be initialized to a zero value,
		// otherwise it can't be removed.
		k.IPWhitelist = make([]string, 0)
	}

	k.IPWhitelistStrict = d.Get("ip_whitelist_strict").(bool)

	return nil
}

// ApikeyCreate creates ns1 API key
func ApikeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k := account.APIKey{}
	if err := resourceDataToApikey(&k, d); err != nil {
		return err
	}
	if resp, err := client.APIKeys.Create(&k); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	// If a key is assigned to at least one team, then it's permissions need to be refreshed
	// because the current key permissions in Terraform state will be out of date.
	if len(k.TeamIDs) > 0 {
		updatedKey, resp, err := client.APIKeys.Get(k.ID)
		if err != nil {
			return ConvertToNs1Error(resp, err)
		}
		// Key attribute only avail on initial GET
		updatedKey.Key = k.Key

		return apikeyToResourceData(d, updatedKey)
	}

	return apikeyToResourceData(d, &k)
}

// ApikeyRead reads API key from ns1
func ApikeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k, resp, err := client.APIKeys.Get(d.Id())
	if err != nil {
		if err == ns1.ErrKeyMissing {
			log.Printf("[DEBUG] NS1 API key (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	return apikeyToResourceData(d, k)
}

//ApikeyDelete deletes the given ns1 api key
func ApikeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.APIKeys.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

//ApikeyUpdate updates the given api key in ns1
func ApikeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k := account.APIKey{
		ID: d.Id(),
	}

	if err := resourceDataToApikey(&k, d); err != nil {
		return err
	}

	if resp, err := client.APIKeys.Update(&k); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	// If a key's teams have changed then the permissions on the key need to be refreshed
	// because the current key permissions in Terraform state will be out of date.
	if d.HasChange("teams") {
		updatedKey, resp, err := client.APIKeys.Get(d.Id())
		if err != nil {
			return ConvertToNs1Error(resp, err)
		}

		return apikeyToResourceData(d, updatedKey)
	}

	return apikeyToResourceData(d, &k)
}
