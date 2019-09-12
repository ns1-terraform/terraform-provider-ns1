package ns1

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Optional
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// SOA attributes per https://tools.ietf.org/html/rfc1035).
			"refresh": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"retry": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"expiry": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// SOA MINUMUM overloaded as NX TTL per https://tools.ietf.org/html/rfc2308
			"nx_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// TODO: test
			"link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"secondaries"},
			},
			"additional_primaries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"secondaries"},
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
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"secondaries": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"primary", "additional_primaries"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  53,
						},
						"networks": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
		Create:   resourceZoneCreate,
		Read:     resourceZoneRead,
		Update:   resourceZoneUpdate,
		Delete:   resourceZoneDelete,
		Importer: &schema.ResourceImporter{State: resourceZoneStateFunc},
	}
}

func resourceZoneToResourceData(d *schema.ResourceData, z *dns.Zone) error {
	d.SetId(z.ID)
	d.Set("hostmaster", z.Hostmaster)
	d.Set("ttl", z.TTL)
	d.Set("nx_ttl", z.NxTTL)
	d.Set("refresh", z.Refresh)
	d.Set("retry", z.Retry)
	d.Set("expiry", z.Expiry)
	d.Set("networks", z.NetworkIDs)
	d.Set("dns_servers", strings.Join(z.DNSServers[:], ","))
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
	if z.Link != nil && *z.Link != "" {
		d.Set("link", *z.Link)
	}
	return nil
}

func secondaryToMap(s *dns.ZoneSecondaryServer) map[string]interface{} {
	m := make(map[string]interface{})
	m["ip"] = s.IP
	m["port"] = s.Port
	m["notify"] = s.Notify
	m["networks"] = s.NetworkIDs
	return m
}

func resourceToZoneData(z *dns.Zone, d *schema.ResourceData) {
	z.ID = d.Id()
	if v, ok := d.GetOk("hostmaster"); ok {
		z.Hostmaster = v.(string)
	}
	if v, ok := d.GetOk("ttl"); ok {
		z.TTL = v.(int)
	}
	if v, ok := d.GetOk("nx_ttl"); ok {
		z.NxTTL = v.(int)
	}
	if v, ok := d.GetOk("refresh"); ok {
		z.Refresh = v.(int)
	}
	if v, ok := d.GetOk("retry"); ok {
		z.Retry = v.(int)
	}
	if v, ok := d.GetOk("expiry"); ok {
		z.Expiry = v.(int)
	}
	if v, ok := d.GetOk("primary"); ok {
		z.MakeSecondary(v.(string))
	}
	if v, ok := d.GetOk("additional_primaries"); ok {
		otherIPsRaw := v.([]interface{})
		z.Secondary.OtherIPs = make([]string, len(otherIPsRaw))
		for i, otherIP := range otherIPsRaw {
			z.Secondary.OtherIPs[i] = otherIP.(string)
		}
		// Fill a list of matching length with '53' for OtherPorts
		// to match functionality of MakeSecondary for PrimaryPort
		// TODO: Add ability to set custom OtherPorts and PrimaryPort
		//otherPorts := make([]int, len(z.Secondary.OtherIPs))
		z.Secondary.OtherPorts = make([]int, len(z.Secondary.OtherIPs))
		for i := range z.Secondary.OtherPorts {
			z.Secondary.OtherPorts[i] = 53
		}
	}
	if v, ok := d.GetOk("secondaries"); ok {
		secondariesSet := v.(*schema.Set)
		secondaries := make([]dns.ZoneSecondaryServer, secondariesSet.Len())
		for i, secondaryRaw := range secondariesSet.List() {
			secondary := secondaryRaw.(map[string]interface{})
			networkIDSet := secondary["networks"].(*schema.Set)
			secondaries[i] = dns.ZoneSecondaryServer{
				NetworkIDs: setToInts(networkIDSet),
				IP:         secondary["ip"].(string),
				Port:       secondary["port"].(int),
				Notify:     secondary["notify"].(bool),
			}
		}
		z.MakePrimary(secondaries...)
	} else {
		// Ensure Primary is cleared out if we remove all of our secondaries
		if _, ok := d.GetOk("primary"); !ok {
			z.Primary = &dns.ZonePrimary{
				Enabled:     false,
				Secondaries: make([]dns.ZoneSecondaryServer, 0),
			}
		}
	}
	if v, ok := d.GetOk("link"); ok {
		z.LinkTo(v.(string))
	}
	if v, ok := d.GetOk("networks"); ok {
		networkIDSet := v.(*schema.Set)
		z.NetworkIDs = setToInts(networkIDSet)
	}
}

// resourceZoneCreate creates the given zone in ns1
func resourceZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	z := dns.NewZone(d.Get("zone").(string))
	resourceToZoneData(z, d)
	if _, err := client.Zones.Create(z); err != nil {
		return err
	}
	if err := resourceZoneToResourceData(d, z); err != nil {
		return err
	}
	return nil
}

// resourceZoneRead reads the given zone data from ns1
func resourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	z, _, err := client.Zones.Get(d.Get("zone").(string))
	if err != nil {
		return err
	}
	if err := resourceZoneToResourceData(d, z); err != nil {
		return err
	}
	return nil
}

// resourceZoneDelete deletes the given zone from ns1
func resourceZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	_, err := client.Zones.Delete(d.Get("zone").(string))
	d.SetId("")
	return err
}

// resourceZoneUpdate updates the zone with given params in ns1
func resourceZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	z := dns.NewZone(d.Get("zone").(string))
	resourceToZoneData(z, d)
	if _, err := client.Zones.Update(z); err != nil {
		return err
	}
	if err := resourceZoneToResourceData(d, z); err != nil {
		return err
	}
	return nil
}

func resourceZoneStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("zone", d.Id())
	return []*schema.ResourceData{d}, nil
}

// translates *schema.Set to []int
func setToInts(schemaSet *schema.Set) []int {
	ints := make([]int, schemaSet.Len())
	for i, v := range schemaSet.List() {
		ints[i] = v.(int)
	}
	return ints
}
