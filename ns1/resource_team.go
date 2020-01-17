package ns1

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"values": {
						Type:     schema.TypeList,
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

	ipWhitelistsRaw := d.Get("ip_whitelist").([]interface{})
	t.IPWhitelist = make([]account.IPWhitelist, 0, len(ipWhitelistsRaw))

	for _, v := range ipWhitelistsRaw {
		ipWhitelistRaw := v.(map[string]interface{})

		valsRaw := ipWhitelistRaw["values"].([]interface{})
		vals := make([]string, 0, len(valsRaw))
		for _, vv := range valsRaw {
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
	if _, err := client.Teams.Create(&t); err != nil {
		return err
	}
	return teamToResourceData(d, &t)
}

// TeamRead reads the team data from ns1
func TeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	t, _, err := client.Teams.Get(d.Id())
	if err != nil {
		if err == ns1.ErrTeamMissing {
			log.Printf("[DEBUG] NS1 team (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}
	return teamToResourceData(d, t)
}

// TeamDelete deletes the given team from ns1
func TeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	_, err := client.Teams.Delete(d.Id())
	d.SetId("")
	return err
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
	if _, err := client.Teams.Update(&t); err != nil {
		return err
	}

	// @TODO - when a teams permissions are updated, all users and keys assigned to that team
	// should have their Terraform state refreshed, there is not a particularly nice way to implement this
	// because teams don't have a concept of what users and keys are assigned to them, only the other way around.
	return teamToResourceData(d, &t)
}
