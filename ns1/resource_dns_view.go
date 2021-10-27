package ns1

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func dnsView() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"created_at": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"read_acls": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"update_acls": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"zones": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"networks": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Optional: true,
		},
		"preference": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}

	return &schema.Resource{
		Schema:        s,
		Create:        DNSViewCreate,
		Read:          DNSViewRead,
		Update:        DNSViewUpdate,
		Delete:        DNSViewDelete,
		Importer:      &schema.ResourceImporter{State: DNSViewImportStateFunc},
		SchemaVersion: 1,
	}
}

func dnsViewToResourceData(d *schema.ResourceData, v *dns.DNSView) error {
	d.SetId(v.Name)
	d.Set("name", v.Name)
	d.Set("created_at", v.Created_at)
	d.Set("updated_at", v.Updated_at)
	d.Set("read_acls", v.Read_acls)
	d.Set("update_acls", v.Update_acls)
	d.Set("zones", v.Zones)
	d.Set("networks", v.Networks)
	d.Set("preference", v.Preference)

	return nil
}

func resourceDataToDNSView(v *dns.DNSView, d *schema.ResourceData) error {
	v.Name = d.Get("name").(string)
	v.Created_at = d.Get("created_at").(int)
	v.Updated_at = d.Get("updated_at").(int)
	if acls, ok := d.GetOk("read_acls"); ok {
		readAclsRaw := acls.([]interface{})
		v.Read_acls = make([]string, len(readAclsRaw))
		for i, readAcls := range readAclsRaw {
			v.Read_acls[i] = readAcls.(string)
		}
	} else {
		v.Read_acls = []string{}
	}
	if acls, ok := d.GetOk("update_acls"); ok {
		updateAclsRaw := acls.([]interface{})
		v.Update_acls = make([]string, len(updateAclsRaw))
		for i, updateAcls := range updateAclsRaw {
			v.Update_acls[i] = updateAcls.(string)
		}
	} else {
		v.Update_acls = []string{}
	}
	if z, ok := d.GetOk("zones"); ok {
		zonesRaw := z.(*schema.Set)
		v.Zones = make([]string, zonesRaw.Len())
		for i, zone := range zonesRaw.List() {
			v.Zones[i] = zone.(string)
		}
	} else {
		v.Zones = []string{}
	}
	if n, ok := d.GetOk("networks"); ok {
		networksRaw := n.(*schema.Set)
		v.Networks = make([]int, networksRaw.Len())
		for i, network := range networksRaw.List() {
			v.Networks[i] = network.(int)
		}
	} else {
		v.Networks = []int{}
	}
	v.Preference = d.Get("preference").(int)

	return nil
}

// DNSViewCreate creates the given DNS View in ns1
func DNSViewCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	v := dns.DNSView{}
	if err := resourceDataToDNSView(&v, d); err != nil {
		return err
	}
	if resp, err := client.View.Create(&v); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return dnsViewToResourceData(d, &v)
}

// DNSViewRead reads the given DNS view data from ns1
func DNSViewRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	v, resp, err := client.View.Get(d.Id())
	if err != nil {
		if errors.Is(err, ns1.ErrViewMissing) {
			log.Printf("[DEBUG] NS1 DNS View (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	// Set Terraform resource data from the job data we just downloaded
	return dnsViewToResourceData(d, v)
}

// DNSViewUpdate updates the DNS view with given parameters in ns1
func DNSViewUpdate(view_schema *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	v := dns.DNSView{
		Name: view_schema.Id(),
	}
	if err := resourceDataToDNSView(&v, view_schema); err != nil {
		return err
	}
	if resp, err := client.View.Update(&v); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return dnsViewToResourceData(view_schema, &v)
}

// DNSViewDelete deletes the given DNS view from ns1
func DNSViewDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	v := dns.DNSView{}
	resourceDataToDNSView(&v, d)
	resp, err := client.View.Delete(v.Name)
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// DNSViewImportStateFunc import the given DNS view from ns1
func DNSViewImportStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
