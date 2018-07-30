package ns1

import (
	"github.com/hashicorp/terraform/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func recordMeta() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"meta": {
				Type: schema.TypeMap,
			},
			"zone": {
				Type: schema.TypeString,
			},
			"domain": {
				Type: schema.TypeString,
			},
			"type": {
				Type: schema.TypeString,
			},
		},

		Create:   RecordMetaCreate,
		Read:     RecordMetaRead,
		Update:   RecordMetaUpdate,
		Delete:   RecordMetaDelete,
		Importer: &schema.ResourceImporter{State: RecordStateFunc},
	}
}

func answerMeta() *schema.Resource {
	return &schema.Resource {
		Schema: map[string]*schema.Schema{
			"meta": {
				Type: schema.TypeMap,
			},
			"zone": {
				Type: schema.TypeString,
			},
			"domain": {
				Type: schema.TypeString,
			},
			"type": {
				Type: schema.TypeString,
			},
			"answer": {
				Type: schema.TypeString,
			},
		},
		Create:   AnswerMetaCreate,
		Read:     AnswerMetaRead,
		Update:   AnswerMetaUpdate,
		Delete:   AnswerMetaDelete,
		Importer: &schema.ResourceImporter{State: RecordStateFunc},
	}
}

func regionMeta() *schema.Resource {
	return &schema.Resource {
		Schema: map[string]*schema.Schema{
			"meta": {
				Type: schema.TypeMap,
			},
			"zone": {
				Type: schema.TypeString,
			},
			"domain": {
				Type: schema.TypeString,
			},
			"type": {
				Type: schema.TypeString,
			},
			"region": {
				Type: schema.TypeString,
			},
		},
		Create:   RegionMetaCreate,
		Read:     RegionMetaRead,
		Update:   RegionMetaUpdate,
		Delete:   RegionMetaDelete,
		Importer: &schema.ResourceImporter{State: RecordStateFunc},
	}
}

func metaToResourceData(d *schema.ResourceData, meta *data.Meta) error {
	d.Set("up", meta.Up)
	d.Set("connections", meta.Connections)
	d.Set("requests", meta.Requests)
	d.Set("asn", meta.ASN)
	d.Set("caprovince", meta.CAProvince)
	d.Set("country", meta.Country)
	d.Set("georegion", meta.Georegion)
	d.Set("highwatermark", meta.HighWatermark)
	d.Set("lowwatermark", meta.LowWatermark)
	d.Set("ipprefixes", meta.IPPrefixes)
	d.Set("latitude", meta.Latitude)
	d.Set("longitude", meta.Longitude)
	d.Set("loadavg", meta.LoadAvg)
	d.Set("priority", meta.Priority)
	d.Set("pulsar", meta.Pulsar)
	d.Set("note", meta.Note)
	d.Set("usstate", meta.USState)
	d.Set("weight", meta.Weight)
	return nil
}

func resourceDataToMeta(d *schema.ResourceData) *data.Meta {
	m := &data.Meta{}
	if v, ok := d.GetOk("meta"); ok {
		m = data.MetaFromMap(v.(map[string]interface{}))
	}
	return m
}

func RecordMetaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	r.Meta = resourceDataToMeta(d)
	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	return recordToResourceData(d, r)
}

func RecordMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return RecordMetaCreate(d, meta)
}

func RecordMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	return metaToResourceData(d, r.Meta)
}

func RecordMetaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	r.Meta = &data.Meta{}

	if  _, err := client.Records.Update(r); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func AnswerMetaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	a := d.Get("answer").(string)
	for _, an := range r.Answers {
		if an.String() == a {
			m = data.MetaFromMap(d.Get("meta").(map[string]interface{}))
			an.Meta = m
		}
	}

	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	return metaToResourceData(d, m)
}

func AnswerMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func AnswerMetaRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func AnswerMetaDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func RegionMetaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r := dns.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err := resourceDataToRecord(r, d); err != nil {
		return err
	}
	if _, err := client.Records.Create(r); err != nil {
		return err
	}
	return recordToResourceData(d, r)
}

func RegionMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func RegionMetaRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func RegionMetaDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// RecordRead reads the DNS record from ns1
//func RecordRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*ns1.Client)
//
//	r, _, err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
//	if err != nil {
//		return err
//	}
//
//	return recordToResourceData(d, r)
//}
//
//// RecordDelete deletes the DNS record from ns1
//func RecordDelete(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*ns1.Client)
//	_, err := client.Records.Delete(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
//	d.SetId("")
//	return err
//}
//
//// RecordUpdate updates the given dns record in ns1
//func RecordUpdate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*ns1.Client)
//	r := dns.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
//	if err := resourceDataToRecord(r, d); err != nil {
//		return err
//	}
//	if _, err := client.Records.Update(r); err != nil {
//		return err
//	}
//	return recordToResourceData(d, r)
//}
