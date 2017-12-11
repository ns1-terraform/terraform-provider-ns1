package ns1

import (
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"


	"github.com/hashicorp/terraform/helper/schema"
)

type MetaType string


const (
	// MetaTypeRecord means that this metadata should be attached at the "record" level
	MetaTypeRecord MetaType = "record"
	// MetaTypeRegion means that this metadata should be attached at the "region" level
	MetaTypeRegion = "region"
	// MetaTypeAnswer means that this metadata should be attached at the "answer" level
	MetaTypeAnswer = "answer"
)

// MetaResource is a terraform specific representation of MetaData in the NS1 platform
type MetaResource struct {
	// Data is the actual metadata
	Data map[string]interface{}

	// MetaType is the level this metadata should be associated with on the given record
	// One of either:
	//   MetaTypeRecord - the lowest level priority
	//   MetaTypeRegion - the medium level priority
	//   MetaTypeAnswer - the highest level priority
	MetaType

	// after looking at the python api these are necessary unless we can add a post by record id
	// comments in the python code indicate that id is read only for a reason though

	// Zone is the zone of the record this metadata is attached to in the NS1 API
	Zone string
	// Domain is the domain of the record to which this metadata is attached
	Domain string
	// RecordType is the type of record for this zone and domain
	RecordType string
}

func metaDataTypeValidate(i interface{}, s string) ([]string, []error) {
	return nil, nil
}

func metaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeSet,
				Required: true,

				// BIGGER TODO(cmc) - I do not yet understand why ValidateFunc is unsupported for
				// non-primitive types, but I intend to find out.

				//ValidateFunc: func(i interface{}, s string) ([]string, []error) {
				//	// TODO(cmc) - write validate function that addresses type limitations of specific keywords
				//	// TODO(cmc) - it does not need to be as complex as the original reflection code and should be
				//	// TODO(cmc) - easy enough to do with a type switch and a map of id's to behaviors
				//	return nil, nil
				//},
			},
			"meta_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			// These three should always be associated with a prior resource, as meta has a record
			// dependency.

			// In .tf that should look something like:

			//  resource "ns1_zone" "example" {
			//   zone = "terraform.example.io"
			//   ttl  = 600
		    // }
			// resource "ns1_record" "my_record" {
			//   zone   = "${ns1_zone.tld.zone}"
			//   domain = "www.${ns1_zone.tld.zone}"
			//   type   = "CNAME"
			//   ttl    = 60
			// }
			//
			// resource "ns1_meta" "my_record_answer_meta" {
			//   data = {
			//		"up" = true
			//	  }
			//   meta_type = "answer"
			//
			//   zone = "${my_record.zone}"
			//   domain = "${my_record.domain}"
			//   record_type = "${my_record.type}"
			// }
			//
			//

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
		},
		// explicitly leave out update, as meta isn't an NS1 resource - it's only a TF resource
		Create: MetaCreate,
		Read:   MetaRead,
		Delete: MetaDelete,
	}
}

// MetaCreate updates the associated record with new metadata
func MetaCreate(d *schema.ResourceData, tfMeta interface{}) error {
	return RecordUpdate(d, tfMeta)
}

// MetaRead finds the appropriate metadata for this record and associated terraform resource
func MetaRead(d *schema.ResourceData, tfMeta interface{}) error {
	return RecordRead(d, tfMeta)
}

// MetaDelete removes the specified metadata from the ns1 record
func MetaDelete(d *schema.ResourceData, tfMeta interface{}) error {

	client := tfMeta.(*ns1.Client)

	r, _, err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	if r.Meta != nil && MetaType(d.Get("meta_type").(string)) == MetaTypeRecord {
		r.Meta = nil
		_, err = client.Records.Update(r)
		return err

	}
	return nil
}