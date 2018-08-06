package ns1

import (
	"github.com/hashicorp/terraform/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

// recordMeta represents metadata at the record level - in NS1's data model, this is metadata at the root level of the record object
// {
//   "_id" : "random_hash_here",
//   "domain" : "example.domain.com",
//   "use_client_subnet" : true,
//   "answers" : [
//     {
//       "answer" : [
//         "example.domain.com"
//      ],
//       "region" : "ec2-ap-south-1",
//       "meta" : {
//         "up" : {
//           "feed" : "feed_id"
//        }
//      },
//       "_id" : "answer_id"
//    }
//  ],
//   "meta" : { # <---- this is the meta we're talking about
//         "country" : [
//           "GB"
//        ]
//  }, ...
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


// recordMeta represents metadata at the answer level - in NS1's data model, this is metadata at the root level of the answer object
// {
//   "_id" : "random_hash_here",
//   "domain" : "example.domain.com",
//   "use_client_subnet" : true,
//   "answers" : [
//     {
//       "answer" : [
//         "example.domain.com"
//      ],
//       "region" : "ec2-ap-south-1",
//       "meta" : { # <--- this is the meta we're talking about
//         "up" : {
//           "feed" : "feed_id"
//        }
//      },
//       "_id" : "answer_id"
//    }
//  ],
//   "meta" : {
//         "country" : [
//           "GB"
//        ]
//  }, ...
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
			"answer_id": {
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

// regionMeta represents metadata at the region level - in NS1's data model, this is metadata at the root level of the region object
// {
//   "_id" : "random_hash_here",
//   "domain" : "example.domain.com",
//   "use_client_subnet" : true,
//   "answers" : [
//     {
//       "answer" : [
//         "example.domain.com"
//      ],
//       "region" : "ec2-ap-south-1",
//       "meta" : {
//         "up" : {
//           "feed" : "feed_id"
//        }
//      },
//       "_id" : "answer_id"
//    }
//  ],
//   "meta" : { # <---- this is the meta we're talking about
//         "country" : [
//           "GB"
//        ]
//  },
//  "regions" : {
//     "ec2-ap-northeast-1" : {
//       "meta" : { # <--- this is the meta we're talking about
//         "country" : [
//           "JP"
//        ]
//      }
//    }, ...
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

// metaToResourceData converts a meta object to TF resource. This function omits the SetId call, so that call should be made in another function before this one
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

// resourceDataToMeta converts resource data to metadata
func resourceDataToMeta(d *schema.ResourceData) *data.Meta {
	m := &data.Meta{}
	if v, ok := d.GetOk("meta"); ok {
		m = data.MetaFromMap(v.(map[string]interface{}))
	}
	return m
}

// RecordMetaCreate creates a meta object at the top level of an NS1 record
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

// RecordMetaUpdate updates metadata at the top level of an NS1 record
func RecordMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return RecordMetaCreate(d, meta)
}

// RecordMetaRead reads metadata at the top level of an NS1 record
func RecordMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	return metaToResourceData(d, r.Meta)
}

// RecordMetaDelete deletes metadata at the top level of an NS1 record
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

// AnswerMetaCreate creates metadata at the top level of an NS1 answer
func AnswerMetaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	var id string
	a := d.Get("answer").(string)
	for _, an := range r.Answers {
		if an.String() == a {
			m = data.MetaFromMap(d.Get("meta").(map[string]interface{}))
			an.Meta = m
			id = an.ID
		}
	}
	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	d.SetId(id)
	return metaToResourceData(d, m)
}

// AnswerMetaUpdate updates metadata at the top level of an NS1 answer
func AnswerMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return AnswerMetaCreate(d, meta)
}

// AnswerMetaRead reads metadata at the top level of an NS1 answer
func AnswerMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	var id string
	a := d.Get("answer").(string)
	for _, an := range r.Answers {
		if an.String() == a {
			id = an.ID
			m = an.Meta
		}
	}
	d.SetId(id)
	return metaToResourceData(d, m)
}

// AnswerMetaDelete deletes metadata at the top level of an NS1 answer
func AnswerMetaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	a := d.Get("answer").(string)
	for _, an := range r.Answers {
		if an.String() == a {
			an.Meta = nil
		}
	}

	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

// RegionMetaCreate creates metadata at the top level of an NS1 region
func RegionMetaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	var id string
	tfRegion := d.Get("region").(string)
	for regionName, region := range r.Regions {
		if regionName == tfRegion {
			m = data.MetaFromMap(d.Get("meta").(map[string]interface{}))
			region.Meta = *m
			id = regionName
		}
	}
	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	d.SetId(id)
	return metaToResourceData(d, m)
}

// RegionMetaUpdate updates metadata at the top level of an NS1 region
func RegionMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	return RegionMetaCreate(d, meta)
}

// RegionMetaRead reads metadata at the top level of an NS1 region
func RegionMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	var id string
	tfRegion := d.Get("region").(string)
	for regionName, region := range r.Regions {
		if regionName == tfRegion {
			m = &region.Meta
			id = regionName
		}
	}

	d.SetId(id)
	return metaToResourceData(d, m)
}

// RegionMetaDelete deletes metadata at the top level of an NS1 region
func RegionMetaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r, _ , err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	var m *data.Meta
	tfRegion := d.Get("region").(string)
	for regionName, region := range r.Regions {
		if regionName == tfRegion {
			region.Meta = *m
		}
	}
	if _, err := client.Records.Update(r); err != nil {
		return err
	}
	d.SetId("")
	return nil
}