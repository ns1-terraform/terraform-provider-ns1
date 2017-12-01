package ns1

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type MetaType string


const (
	MetaTypeRecord MetaType = "record"
	MetaTypeRegion = "region"
	MetaTypeAnswer = "answer"
)

type MetaResource struct {
	// definitely needed
	Data map[string]interface{}
	MetaType

	// maybe not needed
	Zone string
	Domain string
	RecordType string

	// possibly just need record_id
	// RecordID string
}

func metaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"meta_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type: schema.TypeString,
				Required: true,
			},
			"domain": {
				Type: schema.TypeString,
				Required: true,
			},
			"record_type": {
				Type: schema.TypeString,
				Required: true,
			},

			//"record_id": {
			//	Type:     schema.TypeString,
			//	Required: true,
			//	Computed: true,
			//},
		},
		Create: MetaCreate,
		Read:   MetaRead,
		Delete: MetaDelete,
	}
}

// no meta update, this is an ephemeral resource to NS1, it's only a resource for terraform

func MetaCreate(d *schema.ResourceData, tfMeta interface{}) error {

	// just call record.Update here and insert metadata wherever it's supposed to be

	return nil
}

func MetaRead(d *schema.ResourceData, tfMeta interface{}) error {

	// get the record and parse it for the correct meta type, and then set that data accurately
	// in the tf representation

	return nil
}

func MetaDelete(d *schema.ResourceData, tfMeta interface{}) error {

	// upload empty meta for the record?

	return nil
}