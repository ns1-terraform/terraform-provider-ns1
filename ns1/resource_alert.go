package ns1

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
)

func alertResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc:     validatePath,
				// DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"subtype": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc:     validateURL,
				// DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			// Read-only
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Optional
			"notification_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"record_ids"},
			},
			"record_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"zone_names"},
			},
		},
		Create:   AlertConfigCreate,
		Read:     AlertConfigRead,
		Update:   AlertConfigUpdate,
		Delete:   AlertConfigDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func alertToResourceData(d *schema.ResourceData, alert *alerting.Alert) error {
	d.SetId(*alert.ID)
	d.Set("name", alert.Name)
	d.Set("type", alert.Type)
	d.Set("subtype", alert.Subtype)
	d.Set("created_at", alert.CreatedAt)
	d.Set("updated_at", alert.UpdatedAt)
	d.Set("created_by", alert.CreatedBy)
	d.Set("updated_by", alert.UpdatedBy)
	d.Set("notification_lists", alert.NotifierListIds)
	d.Set("zone_names", alert.ZoneNames)
	d.Set("record_ids", alert.RecordIds)
	return nil
}

func strPtr(str string) *string {
	return &str
}

func int64Ptr(n int) *int64 {
	n64 := int64(n)
	return &n64
}

func resourceDataToAlert(d *schema.ResourceData) (*alerting.Alert, error) {
	alert := alerting.Alert{
		ID: strPtr(d.Id()),
	}
	if v, ok := d.GetOk("name"); ok {
		alert.Name = strPtr(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		alert.Type = strPtr(v.(string))
	}
	if v, ok := d.GetOk("subtype"); ok {
		alert.Subtype = strPtr(v.(string))
	}
	if v, ok := d.GetOk("created_at"); ok {
		alert.CreatedAt = int64Ptr(v.(int))
	}
	if v, ok := d.GetOk("updated_at"); ok {
		alert.UpdatedAt = int64Ptr(v.(int))
	}
	if v, ok := d.GetOk("created_by"); ok {
		alert.CreatedBy = strPtr(v.(string))
	}
	if v, ok := d.GetOk("updated_by"); ok {
		alert.UpdatedBy = strPtr(v.(string))
	}
	if v, ok := d.GetOk("notification_lists"); ok {
		listIds := v.(*schema.Set)
		alert.NotifierListIds = make([]string, 0, listIds.Len())
		for _, id := range listIds.List() {
			alert.NotifierListIds = append(alert.NotifierListIds, id.(string))
		}
	} else {
		alert.NotifierListIds = []string{}
	}
	if v, ok := d.GetOk("zone_names"); ok {
		zoneNames := v.(*schema.Set)
		alert.ZoneNames = make([]string, 0, zoneNames.Len())
		for _, zone := range zoneNames.List() {
			alert.ZoneNames = append(alert.ZoneNames, zone.(string))
		}
	} else {
		alert.ZoneNames = []string{}
	}
	if v, ok := d.GetOk("record_ids"); ok {
		recordIds := v.(*schema.Set)
		alert.RecordIds = make([]string, 0, recordIds.Len())
		for _, id := range recordIds.List() {
			alert.RecordIds = append(alert.RecordIds, id.(string))
		}
	} else {
		alert.RecordIds = []string{}
	}
	return &alert, nil
}

func AlertConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	var alert *alerting.Alert = nil
	alert, err := resourceDataToAlert(d)
	if err != nil {
		return err
	}

	if resp, err := client.Alerts.Create(alert); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return alertToResourceData(d, alert)
}

func AlertConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	alert, resp, err := client.Alerts.Get(d.Id())
	if err != nil {
		if err == ns1.ErrAlertMissing {
			log.Printf("[DEBUG] NS1 alert (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}

	return alertToResourceData(d, alert)
}

func AlertConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	alert, err := resourceDataToAlert(d)
	if err != nil {
		return err
	}

	if resp, err := client.Alerts.Update(alert); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return alertToResourceData(d, alert)
}

func AlertConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	resp, err := client.Alerts.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}
