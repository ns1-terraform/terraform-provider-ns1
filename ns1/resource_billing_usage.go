package ns1

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

var MetricTypeStringEnum = NewStringEnum([]string{
	"queries",
	"decisions",
	"records",
	"filter-chains",
	"monitors",
	"limits",
})

func billingUsageResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: MetricTypeStringEnum.ValidateFunc,
			},
			"from": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  int(time.Now().AddDate(0, -1, 0).Unix()), // Default to 1 month ago
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  int(time.Now().Unix()),
			},
			// Queries specific fields
			"clean_queries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ddos_queries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nxd_responses": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"by_network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"clean_queries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ddos_queries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nxd_responses": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"billable_queries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"daily": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timestamp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"clean_queries": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ddos_queries": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"nxd_responses": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			// Limits specific fields
			"queries_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"china_queries_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"records_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter_chains_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"monitors_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"decisions_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nxd_protection_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ddos_protection_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"include_dedicated_dns_network_in_managed_dns_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// Other metrics specific fields
			"total_usage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Read: billingUsageRead,
	}
}

// BillingUsageRead reads the billing usage data from NS1
func billingUsageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	metricType := d.Get("metric_type").(string)
	from := int32(d.Get("from").(int))
	to := int32(d.Get("to").(int))

	var err error

	switch metricType {
	case "queries":
		err = readQueriesUsage(d, client, from, to)
	case "limits":
		err = readLimitsUsage(d, client, from, to)
	case "decisions":
		err = readDecisionsUsage(d, client, from, to)
	case "filter-chains":
		err = readFilterChainsUsage(d, client)
	case "monitors":
		err = readMonitorsUsage(d, client)
	case "records":
		err = readRecordsUsage(d, client)
	default:
		return fmt.Errorf("unsupported metric type: %s", metricType)
	}

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s-%d-%d", metricType, from, to))
	return nil
}

// readQueriesUsage reads and sets the queries usage data
func readQueriesUsage(d *schema.ResourceData, client *ns1.Client, from, to int32) error {
	queries, resp, err := client.BillingUsage.GetQueries(from, to)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	d.Set("clean_queries", queries.CleanQueries)
	d.Set("ddos_queries", queries.DdosQueries)
	d.Set("nxd_responses", queries.NxdResponses)

	// Set by_network data
	networks := make([]map[string]interface{}, len(queries.ByNetwork))
	for i, network := range queries.ByNetwork {
		networks[i] = map[string]interface{}{
			"network":          network.Network,
			"clean_queries":    network.CleanQueries,
			"ddos_queries":     network.DdosQueries,
			"nxd_responses":    network.NxdResponses,
			"billable_queries": network.BillableQueries,
		}

		// Set daily data
		daily := make([]map[string]interface{}, len(network.Daily))
		for j, day := range network.Daily {
			daily[j] = map[string]interface{}{
				"timestamp":     day.Timestamp,
				"clean_queries": day.CleanQueries,
				"ddos_queries":  day.DdosQueries,
				"nxd_responses": day.NxdResponses,
			}
		}
		networks[i]["daily"] = daily
	}
	return d.Set("by_network", networks)
}

// readLimitsUsage reads and sets the limits usage data
func readLimitsUsage(d *schema.ResourceData, client *ns1.Client, from, to int32) error {
	limits, resp, err := client.BillingUsage.GetLimits(from, to)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	d.Set("queries_limit", limits.QueriesLimit)
	d.Set("china_queries_limit", limits.ChinaQueriesLimit)
	d.Set("records_limit", limits.RecordsLimit)
	d.Set("filter_chains_limit", limits.FilterChainsLimit)
	d.Set("monitors_limit", limits.MonitorsLimit)
	d.Set("decisions_limit", limits.DecisionsLimit)
	d.Set("nxd_protection_enabled", limits.NxdProtectionEnabled)
	d.Set("ddos_protection_enabled", limits.DdosProtectionEnabled)
	d.Set("include_dedicated_dns_network_in_managed_dns_usage", limits.IncludeDedicatedDnsNetworkInManagedDnsUsage)

	return nil
}

// readDecisionsUsage reads and sets the decisions usage data
func readDecisionsUsage(d *schema.ResourceData, client *ns1.Client, from, to int32) error {
	usage, resp, err := client.BillingUsage.GetDecisions(from, to)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return d.Set("total_usage", usage.TotalUsage)
}

// readFilterChainsUsage reads and sets the filter chains usage data
func readFilterChainsUsage(d *schema.ResourceData, client *ns1.Client) error {
	usage, resp, err := client.BillingUsage.GetFilterChains()
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return d.Set("total_usage", usage.TotalUsage)
}

// readMonitorsUsage reads and sets the monitors usage data
func readMonitorsUsage(d *schema.ResourceData, client *ns1.Client) error {
	usage, resp, err := client.BillingUsage.GetMonitors()
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return d.Set("total_usage", usage.TotalUsage)
}

// readRecordsUsage reads and sets the records usage data
func readRecordsUsage(d *schema.ResourceData, client *ns1.Client) error {
	usage, resp, err := client.BillingUsage.GetRecords()
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return d.Set("total_usage", usage.TotalUsage)
}
