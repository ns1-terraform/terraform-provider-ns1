package ns1

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
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

func dataSourceZoneToResourceData(d *schema.ResourceData, z *dns.Zone) error {
	d.SetId(z.ID)
	d.Set("hostmaster", z.Hostmaster)
	d.Set("ttl", z.TTL)
	d.Set("nx_ttl", z.NxTTL)
	d.Set("refresh", z.Refresh)
	d.Set("retry", z.Retry)
	d.Set("expiry", z.Expiry)
	d.Set("networks", z.NetworkIDs)
	d.Set("dnssec", z.DNSSEC)
	d.Set("dns_servers", strings.Join(z.DNSServers[:], ","))
	d.Set("link", z.Link)
	if z.Secondary != nil && z.Secondary.Enabled {
		d.Set("primary", z.Secondary.PrimaryIP)
		d.Set("additional_primaries", z.Secondary.OtherIPs)
	}
	if z.Primary != nil && z.Primary.Enabled {
		secondaries := make([]map[string]interface{}, 0)
		for _, secondary := range z.Primary.Secondaries {
			secondaries = append(secondaries, secondaryToMap(&secondary))
		}
		err := d.Set("secondaries", secondaries)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting secondaries for: %s, error: %#v", z.Zone, err)
		}
	}
	return nil
}
