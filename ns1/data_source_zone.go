package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"refresh": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"retry": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expiry": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nx_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_primaries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"additional_ports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"primary_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dns_servers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostmaster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"secondaries": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"networks": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"dnssec": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Read: zoneRead,
	}
}
