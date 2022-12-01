package ns1

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9@$&'*+\-=? ^_.{|}~]{3,320}$`)

func userResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUsername,
		},
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"notify": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     schema.TypeBool,
		},
		"teams": {
			Type:     schema.TypeSet,
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
		Create:        UserCreate,
		Read:          UserRead,
		Update:        UserUpdate,
		Delete:        UserDelete,
		Importer:      &schema.ResourceImporter{State: userImportStateFunc},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    userResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: permissionInstanceStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func userToResourceData(d *schema.ResourceData, u *account.User) error {
	d.SetId(u.Username)
	d.Set("username", u.Username)
	d.Set("name", u.Name)
	d.Set("email", u.Email)
	d.Set("teams", u.TeamIDs)
	notify := make(map[string]bool)
	notify["billing"] = u.Notify.Billing
	d.Set("notify", notify)
	d.Set("ip_whitelist", u.IPWhitelist)
	d.Set("ip_whitelist_strict", u.IPWhitelistStrict)
	permissionsToResourceData(d, u.Permissions)
	return nil
}

func resourceDataToUser(u *account.User, d *schema.ResourceData) error {
	u.Name = d.Get("name").(string)
	u.Username = d.Get("username").(string)
	u.Email = d.Get("email").(string)
	if v, ok := d.GetOk("teams"); ok {
		teamsRaw := v.(*schema.Set).List()
		u.TeamIDs = make([]string, len(teamsRaw))
		for i, team := range teamsRaw {
			u.TeamIDs[i] = team.(string)
		}
	} else {
		u.TeamIDs = make([]string, 0)
	}
	if v, ok := d.GetOk("notify"); ok {
		notifyRaw := v.(map[string]interface{})
		u.Notify.Billing = notifyRaw["billing"].(bool)
	}

	if v, ok := d.GetOk("ip_whitelist"); ok {
		ipWhitelistRaw := v.(*schema.Set)
		u.IPWhitelist = make([]string, ipWhitelistRaw.Len())
		for i, ip := range ipWhitelistRaw.List() {
			u.IPWhitelist[i] = ip.(string)
		}
	} else {
		// This still needs to be initialized to a zero value,
		// otherwise it can't be removed.
		u.IPWhitelist = make([]string, 0)
	}

	u.IPWhitelistStrict = d.Get("ip_whitelist_strict").(bool)

	u.Permissions = resourceDataToPermissions(d)
	return nil
}

// UserCreate creates the given user in ns1
func UserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	u := account.User{}
	if err := resourceDataToUser(&u, d); err != nil {
		return err
	}
	if resp, err := client.Users.Create(&u); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	// If a user is assigned to at least one team, then it's permissions need to be refreshed
	// because the current user permissions in Terraform state will be out of date.
	if len(u.TeamIDs) > 0 {
		updatedUser, resp, err := client.Users.Get(u.Username)
		if err != nil {
			return ConvertToNs1Error(resp, err)
		}

		return userToResourceData(d, updatedUser)
	}

	return userToResourceData(d, &u)
}

// UserRead reads the given users data from ns1
func UserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	u, resp, err := client.Users.Get(d.Id())
	if err != nil {
		// No custom error type is currently defined in the SDK for a non-existent user.
		if strings.Contains(err.Error(), "User not found") {
			log.Printf("[DEBUG] NS1 user (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	return userToResourceData(d, u)
}

// UserDelete deletes the given user from ns1
func UserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Users.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// UserUpdate updates the user with given parameters in ns1
func UserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	u := account.User{
		Username: d.Id(),
	}
	if err := resourceDataToUser(&u, d); err != nil {
		return err
	}

	if resp, err := client.Users.Update(&u); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	// If a user's teams has changed then the permissions on the user need to be refreshed
	// because the current user permissions in Terraform state will be out of date.
	if d.HasChange("teams") {
		updatedUser, resp, err := client.Users.Get(d.Id())
		if err != nil {
			return ConvertToNs1Error(resp, err)
		}

		return userToResourceData(d, updatedUser)
	}

	return userToResourceData(d, &u)
}

func validateUsername(
	val interface{}, key string,
) (warns []string, errs []error) {
	v := []byte(val.(string))
	if !usernameRegex.Match(v) {
		errs = append(
			errs,
			fmt.Errorf(
				"username '%s' does not match regular expression `%s`",
				v,
				usernameRegex,
			),
		)
	}
	return warns, errs
}

func userImportStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
