package ns1

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func teamResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ip_whitelist": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"values": {
						Type:     schema.TypeSet,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Required: true,
					},
				},
			},
		},
	}

	s = addPermsSchema(s)

	return &schema.Resource{
		Schema:        s,
		Create:        TeamCreate,
		Read:          TeamRead,
		Update:        TeamUpdate,
		Delete:        TeamDelete,
		Importer:      &schema.ResourceImporter{State: teamImportStateFunc},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    teamResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: permissionInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func teamToResourceData(d *schema.ResourceData, t *account.Team) error {
	d.SetId(t.ID)
	d.Set("name", t.Name)

	wl := make([]interface{}, 0)
	for _, l := range t.IPWhitelist {
		wlm := make(map[string]interface{})
		wlm["name"] = l.Name
		wlm["values"] = l.Values
		wl = append(wl, wlm)
	}
	if err := d.Set("ip_whitelist", wl); err != nil {
		return fmt.Errorf("[DEBUG] Error setting IP Whitelist for: %s, error: %#v", t.Name, err)
	}

	permissionsToResourceData(d, t.Permissions)

	return nil
}

func resourceDataToTeam(t *account.Team, d *schema.ResourceData) error {
	t.ID = d.Id()
	t.Name = d.Get("name").(string)

	ipWhitelistsRaw := d.Get("ip_whitelist").(*schema.Set)
	t.IPWhitelist = make([]account.IPWhitelist, 0, ipWhitelistsRaw.Len())

	for _, v := range ipWhitelistsRaw.List() {
		ipWhitelistRaw := v.(map[string]interface{})

		valsRaw := ipWhitelistRaw["values"].(*schema.Set)
		vals := make([]string, 0, valsRaw.Len())
		for _, vv := range valsRaw.List() {
			vals = append(vals, vv.(string))
		}

		t.IPWhitelist = append(t.IPWhitelist, account.IPWhitelist{
			Name:   ipWhitelistRaw["name"].(string),
			Values: vals,
		})
	}

	t.Permissions = resourceDataToPermissions(d)
	return nil
}

// TeamCreate creates the given team in ns1
func TeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	t := account.Team{}
	if err := resourceDataToTeam(&t, d); err != nil {
		return err
	}
	if resp, err := client.Teams.Create(&t); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	// workaround INBOX-2226 - send a GET to refresh object
	_ = teamToResourceData(d, &t)
	return TeamRead(d, meta)
}

// TeamRead reads the team data from ns1
func TeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	t, resp, err := client.Teams.Get(d.Id())
	if err != nil {
		if err == ns1.ErrTeamMissing {
			log.Printf("[DEBUG] NS1 team (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	return teamToResourceData(d, t)
}

// TeamDelete deletes the given team from ns1
func TeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Teams.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// TeamUpdate updates the given team in ns1
func TeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	t := account.Team{
		ID: d.Id(),
	}
	if err := resourceDataToTeam(&t, d); err != nil {
		return err
	}
	if resp, err := client.Teams.Update(&t); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	// @TODO - when a teams permissions are updated, all users and keys assigned to that team
	// should have their Terraform state refreshed, there is not a particularly nice way to implement this
	// because teams don't have a concept of what users and keys are assigned to them, only the other way around.
	return teamToResourceData(d, &t)
}

func teamImportStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
