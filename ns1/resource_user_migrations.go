package ns1

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func userResourceV0() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
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
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}

	s = addPermsSchemaV0(s)

	return &schema.Resource{
		Schema: s,
		Create: UserCreate,
		Read:   UserRead,
		Update: UserUpdate,
		Delete: UserDelete,
	}
}
