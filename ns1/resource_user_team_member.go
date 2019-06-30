package ns1

import (
	"github.com/hashicorp/terraform/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func userTeamMemberResource() *schema.Resource {
	//s := userTeamMemberSchema()
	//s = addComputedPermsSchema(s)
	return &schema.Resource{
		Schema: userTeamMemberSchema(),
		Create: UserTeamMemberCreate,
		Read:   UserRead,
		Update: UserUpdate,
		Delete: UserDelete,
	}
}

func userTeamMemberSchema() map[string]*schema.Schema {
	s := userResource().Schema
	for k := range s {
		if k != "name" && k != "username" && k != "email" && k != "notify" && k != "teams" {
			s[k].Computed = true
			// Optional must be false to ensure attribute is not settable in configuration
			s[k].Optional = false
		}
	}

	return s
}

// UserTeamMemberCreate creates the given user in ns1
func UserTeamMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	u := account.User{}
	if err := resourceDataToUser(&u, d); err != nil {
		return err
	}

	// Ensure we are not setting permissions inherited from team
	removePermissionsMap(&u)

	if _, err := client.Users.Create(&u); err != nil {
		return err
	}

	return userToResourceData(d, &u)
}
