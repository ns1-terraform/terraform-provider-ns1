package ns1

import (
	"errors"
	"os"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_APIKEY", nil),
				Description: descriptions["api_key"],
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"ignore_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_IGNORE_SSL", nil),
				Description: descriptions["ignore_ssl"],
			},
			"enable_ddi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_ENABLE_DDI", nil),
				Description: descriptions["enable_ddi"],
			},
			"rate_limit_parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_RATE_LIMIT_PARALLELISM", nil),
				Description: descriptions["rate_limit_parallelism"],
			},
			"retry_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_RETRY_MAX", nil),
				Description: descriptions["retry_max"],
			},
			"user_agent": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NS1_TF_USER_AGENT", nil),
				Description: descriptions["user_agent"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ns1_zone":     dataSourceZone(),
			"ns1_dnssec":   dataSourceDNSSEC(),
			"ns1_record":   dataSourceRecord(),
			"ns1_networks": dataSourceNetworks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ns1_zone":                 resourceZone(),
			"ns1_record":               recordResource(),
			"ns1_datasource":           dataSourceResource(),
			"ns1_datafeed":             dataFeedResource(),
			"ns1_monitoringjob":        monitoringJobResource(),
			"ns1_notifylist":           notifyListResource(),
			"ns1_user":                 userResource(),
			"ns1_apikey":               apikeyResource(),
			"ns1_team":                 teamResource(),
			"ns1_application":          resourceApplication(),
			"ns1_pulsarjob":            pulsarJobResource(),
			"ns1_tsigkey":              tsigKeyResource(),
			"ns1_subnet":               resourceSubnet(),
			"ns1_dnsview":              dnsView(),
			"ns1_account_whitelist":    accountWhitelistResource(),
			"ns1_redirect":             redirectConfigResource(),
			"ns1_redirect_certificate": redirectCertificateResource(),
		},
		ConfigureFunc: ns1Configure,
	}
}

var errNoAPIKey = errors.New("ns1: could not find api key")

func ns1Configure(d *schema.ResourceData) (interface{}, error) {
	config := Config{}
	key := ""
	if k, ok := d.GetOk("apikey"); ok {
		key = k.(string)
	} else {
		key = os.Getenv("NS1_APIKEY")
	}

	if key == "" {
		return nil, errNoAPIKey
	}

	config.Key = key

	if v, ok := d.GetOk("endpoint"); ok {
		config.Endpoint = v.(string)
	}
	if v, ok := d.GetOk("ignore_ssl"); ok {
		config.IgnoreSSL = v.(bool)
	}
	if v, ok := d.GetOk("enable_ddi"); ok {
		config.EnableDDI = v.(bool)
	}
	if v, ok := d.GetOk("rate_limit_parallelism"); ok {
		config.RateLimitParallelism = v.(int)
	}
	if v, ok := d.GetOk("retry_max"); ok {
		config.RetryMax = v.(int)
	}
	if v, ok := d.GetOk("user_agent"); ok {
		config.UserAgent = v.(string)
	}

	return config.Client()
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key":                "The ns1 API key (required)",
		"endpoint":               "URL prefix (including version) for API calls",
		"ignore_ssl":             "Don't validate server SSL/TLS certificate",
		"rate_limit_parallelism": "Tune response to rate limits, see docs",
		"retry_max":              "Maximum retries for 50x errors (-1 to disable)",
		"user_agent":             "User-Agent string to use in NS1 API requests",
		"enable_ddi":             "Deprecated, no longer in use",
	}

	structs.DefaultTagName = "json"
}
