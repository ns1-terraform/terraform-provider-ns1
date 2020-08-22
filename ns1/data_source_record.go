package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"meta": {
				Type:             schema.TypeMap,
				Optional:         true,
				DiffSuppressFunc: metaDiffSuppressUp,
			},
			"link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_client_subnet": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"short_answers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"answers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"answer": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"meta": {
							Type:             schema.TypeMap,
							Optional:         true,
							DiffSuppressFunc: metaDiffSuppress,
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"meta": {
							Type:             schema.TypeMap,
							Optional:         true,
							DiffSuppressFunc: metaDiffSuppress,
						},
					},
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
		Read: dataSourceRecordRead,
	}
}

func dataSourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	r, _, err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	return recordToResourceData(d, r)
}
