package ns1

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"strconv"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func dataFeedResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Create:   DataFeedCreate,
		Read:     DataFeedRead,
		Update:   DataFeedUpdate,
		Delete:   DataFeedDelete,
		Importer: &schema.ResourceImporter{State: dataFeedStateFunc},
	}
}

func dataFeedToResourceData(d *schema.ResourceData, f *data.Feed) {
	configAdapterOut(f)
	d.SetId(f.ID)
	d.Set("name", f.Name)
	d.Set("config", f.Config)
}

func resourceDataToDataFeed(d *schema.ResourceData) (feed *data.Feed, e error) {
	config := d.Get("config").(map[string]interface{})
	if config != nil {
		if testId := config["test_id"]; testId != nil {
			intTestId, err := strconv.Atoi(testId.(string))
			if err != nil {
				return &data.Feed{}, fmt.Errorf("could not convert %v as int %w", testId, err)
			}
			config["test_id"] = intTestId
		}
		if checkId := config["check_id"]; checkId != nil {
			intCheckId, err := strconv.Atoi(checkId.(string))
			if err != nil {
				return &data.Feed{}, fmt.Errorf("could not convert %v as int %w", checkId, err)
			}
			config["check_id"] = intCheckId
		}
		if failOnWarning := config["fail_on_warning"]; failOnWarning != nil {
			boolFailOnWarning, err := strconv.ParseBool(failOnWarning.(string))
			if err != nil {
				return &data.Feed{}, fmt.Errorf("could not convert %v as bool %w", failOnWarning, err)
			}
			config["fail_on_warning"] = boolFailOnWarning
		}
		if failOnNoData := config["fail_on_no_data"]; failOnNoData != nil {
			boolFailOnNoData, err := strconv.ParseBool(failOnNoData.(string))
			if err != nil {
				return &data.Feed{}, fmt.Errorf("could not convert %v as bool %w", failOnNoData, err)
			}
			config["fail_on_no_data"] = boolFailOnNoData
		}
	}

	return &data.Feed{
		Name:     d.Get("name").(string),
		SourceID: d.Get("source_id").(string),
		Config:   config,
	}, nil
}

// DataFeedCreate creates an ns1 datafeed
func DataFeedCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	f, err := resourceDataToDataFeed(d)
	if err != nil {
		return err
	}
	if resp, err := client.DataFeeds.Create(d.Get("source_id").(string), f); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	dataFeedToResourceData(d, f)
	return nil
}

// DataFeedRead reads the datafeed for the given ID from ns1
func DataFeedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	f, resp, err := client.DataFeeds.Get(d.Get("source_id").(string), d.Id())
	if err != nil {
		// No custom error type is currently defined in the SDK for a data feed.
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[DEBUG] NS1 data source (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	dataFeedToResourceData(d, f)
	return nil
}

// DataFeedDelete delets the given datafeed from ns1
func DataFeedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.DataFeeds.Delete(d.Get("source_id").(string), d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// DataFeedUpdate updates the given datafeed with modified parameters
func DataFeedUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	f, err := resourceDataToDataFeed(d)
	if err != nil {
		return err
	}
	f.ID = d.Id()
	if resp, err := client.DataFeeds.Update(d.Get("source_id").(string), f); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	dataFeedToResourceData(d, f)
	return nil
}

// configAdapterOut back the original configuration types
func configAdapterOut(f *data.Feed) {
	config := f.Config
	if config != nil {
		if testId := config["test_id"]; testId != nil {
			strTestId := strconv.Itoa(int(testId.(float64)))
			config["test_id"] = strTestId
		}
		if checkId := config["check_id"]; checkId != nil {
			strCheckId := strconv.Itoa(int(checkId.(float64)))
			config["check_id"] = strCheckId
		}
		if failOnWarning := config["fail_on_warning"]; failOnWarning != nil {
			strFailOnWarning := strconv.FormatBool(failOnWarning.(bool))
			config["fail_on_warning"] = strFailOnWarning
		}
		if failOnNoData := config["fail_on_no_data"]; failOnNoData != nil {
			strFailOnNoData := strconv.FormatBool(failOnNoData.(bool))
			config["fail_on_no_data"] = strFailOnNoData
		}
		f.Config = config
	}
}

func dataFeedStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid datafeed specifier.  Expecting 1 slashe (\"datasource_id/datafeed_id\"), got %d", len(parts)-1)
	}

	d.SetId(parts[1])
	d.Set("source_id", parts[0])

	return []*schema.ResourceData{d}, nil
}
