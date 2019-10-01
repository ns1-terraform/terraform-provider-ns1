package ns1

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

var dnsKeySchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"flags": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"algorithm": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_key": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func dataSourceDNSSEC() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dnskey": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dnsKeySchema,
						},
						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"delegation": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dnskey": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dnsKeySchema,
						},
						"ds": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dnsKeySchema,
						},
						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
		Read: dnssecRead,
	}
}

func dnssecToResourceData(d *schema.ResourceData, z *dns.ZoneDNSSEC) error {
	d.SetId(fmt.Sprintf("%sdnssec", z.Zone))
	// remove trailing "." for consistency with resources
	d.Set("zone", strings.TrimSuffix(z.Zone, "."))
	d.Set("keys", flattenKeys(z.Keys))
	d.Set("delegation", flattenDelegation(z.Delegation))
	return nil
}

func dnssecRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	z, _, err := client.DNSSEC.Get(d.Get("zone").(string))
	if err != nil {
		return err
	}
	if err := dnssecToResourceData(d, z); err != nil {
		return err
	}
	return nil
}

func flattenKeys(keys *dns.Keys) []interface{} {
	m := make(map[string]interface{})
	m["dnskey"] = flattenDNSKeys(keys.DNSKey)
	m["ttl"] = keys.TTL
	return []interface{}{m}
}

func flattenDelegation(delegation *dns.Delegation) []interface{} {
	m := make(map[string]interface{})
	m["dnskey"] = flattenDNSKeys(delegation.DNSKey)
	m["ds"] = flattenDNSKeys(delegation.DS)
	m["ttl"] = delegation.TTL
	return []interface{}{m}
}

func flattenDNSKeys(keys []*dns.Key) []interface{} {
	out := make([]interface{}, 0, 0)
	for _, v := range keys {
		m := make(map[string]interface{})
		m["flags"] = v.Flags
		m["protocol"] = v.Protocol
		m["algorithm"] = v.Algorithm
		m["public_key"] = v.PublicKey
		out = append(out, m)
	}
	return out
}
