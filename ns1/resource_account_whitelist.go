package ns1

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func accountWhitelistResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"values": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}

	return &schema.Resource{
		Schema:        s,
		Create:        accountWhitelistCreate,
		Read:          accountWhitelistRead,
		Update:        accountWhitelistUpdate,
		Delete:        accountWhitelistDelete,
		Importer:      &schema.ResourceImporter{},
		SchemaVersion: 1,
	}
}

func accountWhitelistToResourceData(d *schema.ResourceData, wl *account.IPWhitelist) error {
	d.SetId(wl.ID)
	d.Set("name", wl.Name)
	d.Set("values", wl.Values)
	return nil
}

func resourceDataToWhitelist(wl *account.IPWhitelist, d *schema.ResourceData) error {
	wl.ID = d.Id()
	wl.Name = d.Get("name").(string)

	if v, ok := d.GetOk("values"); ok {
		ipWhitelistRaw := v.([]interface{})
		wl.Values = make([]string, len(ipWhitelistRaw))
		for i, ip := range ipWhitelistRaw {
			wl.Values[i] = ip.(string)
		}
	}

	return nil
}

func accountWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	wl := account.IPWhitelist{}
	if err := resourceDataToWhitelist(&wl, d); err != nil {
		return err
	}
	if resp, err := client.GlobalIPWhitelist.Create(&wl); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return accountWhitelistToResourceData(d, &wl)
}

func accountWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	wl, resp, err := client.GlobalIPWhitelist.Get(d.Id())
	if err != nil {
		if err == ns1.ErrIPWhitelistMissing {
			log.Printf("[DEBUG] NS1 global whitelist (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return ConvertToNs1Error(resp, err)
	}

	return accountWhitelistToResourceData(d, wl)
}

func accountWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	wl := account.IPWhitelist{
		ID: d.Id(),
	}

	if err := resourceDataToWhitelist(&wl, d); err != nil {
		return err
	}

	if resp, err := client.GlobalIPWhitelist.Update(&wl); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return accountWhitelistToResourceData(d, &wl)
}

func accountWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.GlobalIPWhitelist.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}
