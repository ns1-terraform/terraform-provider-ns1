package ns1

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dataset"
)

var DatatypeTypeStringEnum = NewStringEnum([]string{
	"num_queries",
	"num_ebot_response",
	"num_nxd_response",
	"zero_queries",
})

var DatatypeScopeStringEnum = NewStringEnum([]string{
	"account",
	"network_single",
	"record_single",
	"zone_single",
	"network_each",
	"record_each",
	"zone_each",
	"top_n_zones",
	"top_n_records",
})

var RepeatsEveryStringEnum = NewStringEnum([]string{
	"week",
	"month",
	"year",
})

var TimeframeAggregationStringEnum = NewStringEnum([]string{
	"daily",
	"monthly",
	"billing_period",
})

var ExportTypeStringEnum = NewStringEnum([]string{
	"csv",
	"json",
	"xlsx",
})

func datasetResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"datatype": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: DatatypeTypeStringEnum.ValidateFunc,
						},
						"scope": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: DatatypeScopeStringEnum.ValidateFunc,
						},
						"data": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
			"repeat": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"repeats_every": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: RepeatsEveryStringEnum.ValidateFunc,
						},
						"end_after_n": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"timeframe": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregation": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: TimeframeAggregationStringEnum.ValidateFunc,
						},
						"cycles": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"from": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"export_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ExportTypeStringEnum.ValidateFunc,
			},
			"recipient_emails": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Read-only
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
		Create: DatasetCreate,
		Read:   DatasetRead,
		Update: DatasetUpdate,
		Delete: DatasetDelete,
	}
}

// DatasetCreate creates a dataset
func DatasetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	resourceDatatype := d.Get("datatype").([]interface{})[0].(map[string]interface{})
	resourceDatatypeData := resourceDatatype["data"].(map[string]interface{})
	datatypeDate := make(map[string]string)
	for k, v := range resourceDatatypeData {
		datatypeDate[k] = v.(string)
	}
	datatype := &dataset.Datatype{
		Type:  dataset.DatatypeType(resourceDatatype["type"].(string)),
		Scope: dataset.DatatypeScope(resourceDatatype["scope"].(string)),
		Data:  datatypeDate,
	}

	resourceTimeframe := d.Get("timeframe").([]interface{})[0].(map[string]interface{})
	resourceTimeframeCycles := int32(resourceTimeframe["cycles"].(int))
	resourceTimeframeFrom := newUnixTimestamp(int64(resourceTimeframe["from"].(int)))
	resourceTimeframeTo := newUnixTimestamp(int64(resourceTimeframe["to"].(int)))
	timeframe := &dataset.Timeframe{
		Aggregation: dataset.TimeframeAggregation(resourceTimeframe["aggregation"].(string)),
		Cycles:      &resourceTimeframeCycles,
		From:        resourceTimeframeFrom,
		To:          resourceTimeframeTo,
	}

	resourceListRepeat := d.Get("repeat").([]interface{})
	var repeat *dataset.Repeat
	if len(resourceListRepeat) > 0 {
		resourceRepeat := resourceListRepeat[0].(map[string]interface{})
		repeat = &dataset.Repeat{
			Start:        dataset.UnixTimestamp(time.Unix(int64(resourceRepeat["start"].(int)), 0)),
			RepeatsEvery: dataset.RepeatsEvery(resourceRepeat["repeats_every"].(string)),
			EndAfterN:    int32(resourceRepeat["end_after_n"].(int)),
		}
	}

	resourceRecipientEmails := d.Get("recipient_emails").(*schema.Set).List()
	recipientEmails := make([]string, 0)
	for _, email := range resourceRecipientEmails {
		recipientEmails = append(recipientEmails, email.(string))
	}

	r := dataset.NewDataset(
		"",
		d.Get("name").(string),
		datatype,
		repeat,
		timeframe,
		dataset.ExportType(d.Get("export_type").(string)),
		nil,
		recipientEmails,
		dataset.UnixTimestamp{},
		dataset.UnixTimestamp{},
	)

	dt, resp, err := client.Datasets.Create(r)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return datasetToResourceData(d, dt)
}

// DatasetRead reads the dataset from ns1
func DatasetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	cfg, resp, err := client.Datasets.Get(d.Get("id").(string))
	if err != nil {
		if errors.Is(err, ns1.ErrDatasetNotFound) {
			log.Printf("[DEBUG] NS1 redirect config (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}

	return datasetToResourceData(d, cfg)
}

// DatasetDelete deletes the dataset from ns1
func DatasetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Datasets.Delete(d.Get("id").(string))
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// DatasetUpdate updates the dataset from ns1
func DatasetUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("datasets cannot be updated")
}

func newUnixTimestamp(sec int64) *dataset.UnixTimestamp {
	if sec == 0 {
		return nil
	}
	u := dataset.UnixTimestamp(time.Unix(sec, 0))
	return &u
}

func datasetToResourceData(rd *schema.ResourceData, dt *dataset.Dataset) error {
	if len(rd.Id()) == 0 {
		rd.SetId(dt.ID)
	}

	rd.Set("name", dt.Name)
	rd.Set("export_type", dt.ExportType)
	rd.Set("recipient_emails", dt.RecipientEmails)

	rd.Set("datatype", []map[string]interface{}{
		{
			"type":  dt.Datatype.Type,
			"scope": dt.Datatype.Scope,
			"data":  dt.Datatype.Data,
		},
	})

	var timeframeFrom, timeframeTo, timeframeCycles = 0, 0, 0
	if dt.Timeframe.From != nil {
		timeframeFrom = int(time.Time(*dt.Timeframe.From).Unix())
	}
	if dt.Timeframe.To != nil {
		timeframeFrom = int(time.Time(*dt.Timeframe.To).Unix())
	}
	if dt.Timeframe.Cycles != nil {
		timeframeCycles = int(*dt.Timeframe.Cycles)
	}

	rd.Set("timeframe", []map[string]interface{}{
		{
			"aggregation": dt.Timeframe.Aggregation,
			"cycles":      timeframeCycles,
			"from":        timeframeFrom,
			"to":          timeframeTo,
		},
	})

	if dt.Repeat != nil {
		rd.Set("repeat", []map[string]interface{}{
			{
				"start":         time.Time(dt.Repeat.Start).Unix(),
				"repeats_every": dt.Repeat.RepeatsEvery,
				"end_after_n":   dt.Repeat.EndAfterN,
			},
		})
	}

	rd.Set("reports", nil)
	if len(dt.Reports) > 0 {
		reports := make([]map[string]interface{}, 0)
		for _, report := range dt.Reports {
			reports = append(reports, map[string]interface{}{
				"id":         report.ID,
				"status":     report.Status,
				"start":      time.Time(report.Start).Unix(),
				"end":        time.Time(report.End).Unix(),
				"created_at": time.Time(report.CreatedAt).Unix(),
			})
		}
		rd.Set("reports", reports)
	}

	return nil
}
