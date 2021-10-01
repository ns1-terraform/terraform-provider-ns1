package ns1

import (
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
		Create:        PulsarJobCreate,
		SchemaVersion: 1,
	}
}

func tsigKeyToResourceData(d *schema.ResourceData, k *dns.Tsig_key) error {
	d.SetId(k.Name)
	d.Set("algorithm", k.Algorithm)
	d.Set("secret", k.Secret)

	return nil
}

func resourceDataToTsigKey(k *dns.Tsig_key, d *schema.ResourceData) error {
	k.Name = d.Id()
	k.Algorithm = d.Get("algorithm").(string)
	k.Secret = d.Get("secret").(string)

	return nil
}

// TsigKeyCreate creates the given TSIG key in ns1
func TsigKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	k := dns.Tsig_key{}
	if err := resourceDataToTsigKey(&k, d); err != nil {
		return err
	}
	if resp, err := client.TSIG.Create(&k); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return tsigKeyToResourceData(d, &k)
}
