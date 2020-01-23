package ns1

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func teamResourceV0() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	s = addPermsSchemaV0(s)

	return &schema.Resource{
		Schema: s,
		Create: TeamCreate,
		Read:   TeamRead,
		Update: TeamUpdate,
		Delete: TeamDelete,
	}
}
