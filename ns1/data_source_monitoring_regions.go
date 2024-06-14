package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

func dataSourceMonitoringRegions() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnets": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Read: MonitoringingRegionsRead,
	}
}

// MonitoringRegionsRead reads the available Monitoring Regions from ns1.
func MonitoringingRegionsRead(d *schema.ResourceData, meta any) error {
	client := meta.(*ns1.Client)

	regions, resp, err := client.MonitorRegions.List()
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	out := []map[string]any{}

	for _, region := range regions {
		out = append(out, map[string]any{
			"code":    region.Code,
			"name":    region.Name,
			"subnets": region.Subnets,
		})
	}

	d.SetId("1")
	return d.Set("regions", out)
}
