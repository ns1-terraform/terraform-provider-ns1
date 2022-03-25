package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/ipam"
	"strconv"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"total_addresses": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"free_addresses": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"used_addresses": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_scoped": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
		Create:   subnetCreate,
		Read:     subnetRead,
		Update:   subnetUpdate,
		Delete:   subnetDelete,
	}
}

func resourceSubnetToResourceData(d *schema.ResourceData, a *ipam.Address) error {
	d.SetId(strconv.Itoa(a.ID))
	d.Set("prefix", a.Prefix)
	d.Set("network_id", a.Network)
	d.Set("parent_id", a.Parent)
	d.Set("desc", a.Desc)
	d.Set("name", a.Name)
	d.Set("status", a.Status)
	d.Set("tags", a.Tags)
	d.Set("total_addresses", a.Total)
	d.Set("children", a.Children)
	d.Set("free_addresses", a.Free)
	d.Set("used_addresses", a.Used)
	d.Set("dhcp_scoped", a.DHCPScoped)
	return nil
}

func resourceDataToSubnet(a *ipam.Address, d *schema.ResourceData) {
	a.ID, _ = strconv.Atoi(d.Id())
	if v, ok := d.GetOk("prefix"); ok {
		a.Prefix = v.(string)
	}
	if v, ok := d.GetOk("network_id"); ok {
		a.Network = v.(int)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		a.Parent = v.(int)
	}
	if v, ok := d.GetOk("desc"); ok {
		a.Desc = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		a.Name = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		a.Status = ipam.AddrStatus(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		a.Tags = v.(map[string]interface{})
	}
	if v, ok := d.GetOk("total_addresses"); ok {
		a.Total = v.(string)
	}
	if v, ok := d.GetOk("children"); ok {
		a.Children = v.(int)
	}
	if v, ok := d.GetOk("free_addresses"); ok {
		a.Free = v.(string)
	}
	if v, ok := d.GetOk("used_addresses"); ok {
		a.Used = v.(string)
	}
	if v, ok := d.GetOk("dhcp_scoped"); ok {
		a.DHCPScoped = v.(bool)
	}
}

// subnetCreate creates the given subnet in ns1
func subnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	a := ipam.Address{}
	resourceDataToSubnet(&a, d)
	sub, resp, err := client.IPAM.CreateSubnet(&a)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceSubnetToResourceData(d, sub); err != nil {
		return err
	}
	return nil
}

// subnetRead reads the given subnet data from ns1
func subnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	id, _ := strconv.Atoi(d.Id())
	s,resp, err := client.IPAM.GetSubnet(id)

	if err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceSubnetToResourceData(d, s); err != nil {
		return err
	}
	return nil
}

// subnetDelete deletes the given subnet from ns1
func subnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	id, _ := strconv.Atoi(d.Id())
	subnet, err := client.IPAM.DeleteSubnet(id)
	d.SetId("")
	return ConvertToNs1Error(subnet, err)
}

// subnetUpdate updates the subnet with given params in ns1
func subnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	a := ipam.Address{}
	resourceDataToSubnet(&a, d)
	subnet, _, resp, err := client.IPAM.EditSubnet(&a, true)

	if err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceSubnetToResourceData(d, subnet); err != nil {
		return err
	}
	return nil
}

