package ns1

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

var networkSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	},
}

func dataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     networkSchema,
			},
		},
		Read: networksRead,
	}
}

// networkRead reads the networks from ns1
func networksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	networks, resp, err := client.Network.Get()
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return networksToResourceData(d, networks)
}

func networksToResourceData(d *schema.ResourceData, n []*dns.Network) error {

	d.SetId("123_test")
	networks := make([]interface{}, 0)
	for _, network := range n {
		networkMap := make(map[string]interface{})
		networkMap["label"] = network.Label
		networkMap["name"] = network.Name
		networkMap["network_id"] = network.NetworkID
		networks = append(networks, networkMap)
	}
	if err := d.Set("networks", networks); err != nil {
		return fmt.Errorf("[DEBUG] Error Getting Networks, error: %#v", err)

	}

	return nil
}
