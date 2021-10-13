package ns1

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func tsigKeyResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"algorithm": {
			Type:     schema.TypeString,
			Required: true,
		},
		"secret": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	return &schema.Resource{
		Schema:        s,
		Create:        tsigKeyCreate,
		Read:          tsigKeyRead,
		Update:        tsigKeyUpdate,
		Delete:        tsigKeyDelete,
		Importer:      &schema.ResourceImporter{State: tsigKeyImportStateFunc},
		SchemaVersion: 1,
	}
}

func tsigKeyToResourceData(d *schema.ResourceData, k *dns.TSIGKey) error {
	d.SetId(k.Name)
	d.Set("name", k.Name)
	d.Set("algorithm", k.Algorithm)
	d.Set("secret", k.Secret)

	return nil
}

func resourceDataToTsigKey(k *dns.TSIGKey, d *schema.ResourceData) error {
	k.Name = d.Get("name").(string)
	k.Algorithm = d.Get("algorithm").(string)
	k.Secret = d.Get("secret").(string)

	return nil
}

// TsigKeyCreate creates the given TSIG key in ns1
func tsigKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k := dns.TSIGKey{}
	if err := resourceDataToTsigKey(&k, d); err != nil {
		return err
	}
	if resp, err := client.TSIG.Create(&k); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return TsigKeyToResourceData(d, &k)
}

// TsigKeyRead reads the given TSIG key from ns1
func tsigKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k, resp, err := client.TSIG.Get(d.Id())
	if err != nil {
		if errors.Is(err, ns1.ErrTsigKeyMissing) {
			log.Printf("[DEBUG] NS1 TSIG Key (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	// Set Terraform resource data from the tsig key data we just downloaded
	if err := TsigKeyToResourceData(d, k); err != nil {
		return err
	}
	return nil
}

// TsigKeyUpdate updates the TSIG Key with given parameters in ns1
func tsigKeyUpdate(key_schema *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k := dns.TSIGKey{}
	if err := resourceDataToTsigKey(&k, key_schema); err != nil {
		return err
	}

	if resp, err := client.TSIG.Update(&k); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return tsigKeyToResourceData(key_schema, &k)
}

// TsigKeyDelete deletes the given TSIG Key from ns1
func tsigKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.TSIG.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

func tsigKeyImportStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
