package ns1

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

const (
	MetricTypeQueries      = "queries"
	MetricTypeLimits       = "limits"
	MetricTypeDecisions    = "decisions"
	MetricTypeRecords      = "records"
	MetricTypeFilterChains = "filter-chains"
	MetricTypeMonitors     = "monitors"
)

var MetricTypeStringEnum = NewStringEnum([]string{
	MetricTypeQueries,
	MetricTypeDecisions,
	MetricTypeRecords,
	MetricTypeFilterChains,
	MetricTypeMonitors,
	MetricTypeLimits,
})

var timeFrameRequiredMetricTypes = []string{
	MetricTypeQueries,
	MetricTypeLimits,
	MetricTypeDecisions,
}

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
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
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

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// BillingUsageRead reads the billing usage data from NS1
func billingUsageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	metricType := d.Get("metric_type").(string)

	// Validate that from and to are provided for metric types that require them
	if contains(timeFrameRequiredMetricTypes, metricType) {
		_, okFrom := d.GetOk("from")
		_, okTo := d.GetOk("to")
		if !okFrom || !okTo {
			return fmt.Errorf("from and to parameters are required for metric_type: %s", metricType)
		}
	}

	var err error

	switch metricType {
	case MetricTypeQueries:
		from := int32(d.Get("from").(int))
		to := int32(d.Get("to").(int))
		err = readQueriesUsage(d, client, from, to)
	case MetricTypeLimits:
		from := int32(d.Get("from").(int))
		to := int32(d.Get("to").(int))
		err = readLimitsUsage(d, client, from, to)
	case MetricTypeDecisions:
		from := int32(d.Get("from").(int))
		to := int32(d.Get("to").(int))
		err = readDecisionsUsage(d, client, from, to)
	case MetricTypeFilterChains:
		err = readFilterChainsUsage(d, client)
	case MetricTypeMonitors:
		err = readMonitorsUsage(d, client)
	case MetricTypeRecords:
		err = readRecordsUsage(d, client)
	default:
		return fmt.Errorf("unsupported metric type: %s", metricType)
	}

	if err != nil {
		return err
	}

	// Set a unique ID for the data source
	if contains(timeFrameRequiredMetricTypes, metricType) {
		from := int32(d.Get("from").(int))
		to := int32(d.Get("to").(int))
		d.SetId(fmt.Sprintf("%s-%d-%d", metricType, from, to))
	} else {
		d.SetId(metricType)
	}

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
