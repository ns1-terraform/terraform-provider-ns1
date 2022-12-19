package ns1

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
	"gopkg.in/ns1/ns1-go.v2/rest/model/filter"
)

var recordTypeStringEnum = NewStringEnum([]string{
	"A",
	"AAAA",
	"ALIAS",
	"AFSDB",
	"CAA",
	"CNAME",
	"DNAME",
	"DS",
	"HINFO",
	"MX",
	"NAPTR",
	"NS",
	"PTR",
	"RP",
	"SPF",
	"SRV",
	"TXT",
	"URLFWD",
	"strings",
})

func recordResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateFQDN,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"domain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateFQDN,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: recordTypeStringEnum.ValidateFunc,
			},
			// Optional
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"override_ttl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"meta": {
				Type:             schema.TypeMap,
				Optional:         true,
				DiffSuppressFunc: metaDiffSuppress,
			},
			"link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_client_subnet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"short_answers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Deprecated: `short_answers will be deprecated in a future release.
It is suggested to migrate to a regular "answers" block. Using Terraform 0.12+, a similar convenience to "short_answers" can be achieved with dynamic blocks:
  dynamic "answers" {
    for_each = ["4.4.4.4", "5.5.5.5"]
    content {
      answer  = answers.value
    }
  }`,
			},
			"answers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"answer": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"meta": {
							Type:             schema.TypeMap,
							Optional:         true,
							DiffSuppressFunc: metaDiffSuppress,
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"meta": {
							Type:             schema.TypeMap,
							Optional:         true,
							DiffSuppressFunc: metaDiffSuppress,
						},
					},
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
							Default:  filter.Config{},
						},
					},
				},
			},
		},
		Create:   RecordCreate,
		Read:     RecordRead,
		Update:   RecordUpdate,
		Delete:   RecordDelete,
		Importer: &schema.ResourceImporter{State: recordStateFunc},
	}
}

// errJoin joins errors into a single error
func errJoin(errs []error, sep string) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return errors.New(errs[0].Error() + sep + errs[1].Error())
	case 3:
		// Same special case
		return errors.New(errs[0].Error() + sep + errs[1].Error() + sep + errs[2].Error())
	}

	n := len(sep) * (len(errs) - 1)
	for i := 0; i < len(errs); i++ {
		n += len(errs[i].Error())
	}

	b := make([]byte, n)
	bp := copy(b, errs[0].Error())
	for _, err := range errs[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], err.Error())
	}
	return errors.New(string(b))
}

func recordToResourceData(d *schema.ResourceData, r *dns.Record) error {
	d.SetId(r.ID)
	d.Set("domain", r.Domain)
	d.Set("zone", r.Zone)
	d.Set("type", r.Type)
	d.Set("ttl", r.TTL)

	d.Set("override_ttl", nil)
	if r.Type == "ALIAS" && r.Override_TTL != nil {
		err := d.Set("override_ttl", *r.Override_TTL)

		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting override_ttl for: %s, error: %#v", r.Domain, err)
		}
	}

	if r.Link != "" {
		err := d.Set("link", r.Link)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting link for: %s, error: %#v", r.Domain, err)
		}
	}

	// top level meta works but nested meta doesn't
	if r.Meta != nil {
		err := d.Set("meta", r.Meta.StringMap())
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting meta for: %s, error: %#v", r.Domain, err)
		}
	}
	if r.UseClientSubnet != nil {
		err := d.Set("use_client_subnet", *r.UseClientSubnet)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting use_client_subnet for: %s, error: %#v", r.Domain, err)
		}
	}
	if len(r.Filters) > 0 {
		filters := make([]map[string]interface{}, len(r.Filters))
		for i, f := range r.Filters {
			m := make(map[string]interface{})
			m["filter"] = f.Type
			if f.Disabled {
				m["disabled"] = true
			}
			if f.Config != nil {
				m["config"] = recordMapValueToString(f.Config)
			}
			filters[i] = m
		}
		err := d.Set("filters", filters)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting filters for: %s, error: %#v", r.Domain, err)
		}
	}
	if len(r.Answers) > 0 {
		ans := make([]map[string]interface{}, 0)
		log.Printf("Got back from ns1 answers: %+v", r.Answers)
		for _, answer := range r.Answers {
			ans = append(ans, answerToMap(*answer))
		}
		log.Printf("Setting answers %+v", ans)
		err := d.Set("answers", ans)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting answers for: %s, error: %#v", r.Domain, err)
		}
	}
	if len(r.Regions) > 0 {
		keys := make([]string, 0, len(r.Regions))
		for regionName := range r.Regions {
			keys = append(keys, regionName)
		}
		sort.Strings(keys)
		regions := make([]map[string]interface{}, 0, len(r.Regions))
		for _, k := range keys {
			newRegion := make(map[string]interface{})
			region := r.Regions[k]
			newRegion["name"] = k
			newRegion["meta"] = region.Meta.StringMap()
			regions = append(regions, newRegion)
		}
		log.Printf("Setting regions %+v", regions)
		err := d.Set("regions", regions)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting regions for: %s, error: %#v", r.Domain, err)
		}
	}
	return nil
}

func recordMapValueToString(configMap map[string]interface{}) map[string]interface{} {
	config := make(map[string]interface{})
	for configKey, configValue := range configMap {
		switch t := configValue.(type) {
		case bool:
			if t {
				config[configKey] = "1"
			} else {
				config[configKey] = "0"
			}
		case float64:
			config[configKey] = strconv.FormatFloat(t, 'f', -1, 64)
		default:
			config[configKey] = t
		}
	}
	return config
}

func answerToMap(a dns.Answer) map[string]interface{} {
	m := make(map[string]interface{})
	m["answer"] = strings.Join(a.Rdata, " ")
	if a.RegionName != "" {
		m["region"] = a.RegionName
	}
	if a.Meta != nil {
		log.Println("got meta: ", a.Meta)
		m["meta"] = metaToMapString(a.Meta)
		log.Println("converted meta to resource string: ", m["meta"])
	}
	return m
}

func resourceDataToRecord(r *dns.Record, d *schema.ResourceData) error {
	r.ID = d.Id()
	log.Printf("answers from template: %+v, %T\n", d.Get("answers"), d.Get("answers"))

	if shortAnswers := d.Get("short_answers").([]interface{}); len(shortAnswers) > 0 {
		for _, answerRaw := range shortAnswers {
			answer := answerRaw.(string)
			switch d.Get("type") {
			case "TXT", "SPF":
				r.AddAnswer(dns.NewTXTAnswer(answer))
			case "CAA":
				r.AddAnswer(dns.NewAnswer(strings.SplitN(answer, " ", 3)))
			default:
				r.AddAnswer(dns.NewAnswer(strings.Split(answer, " ")))
			}
		}
	}
	if answers := d.Get("answers").([]interface{}); len(answers) > 0 {
		for _, answerRaw := range answers {
			answer := answerRaw.(map[string]interface{})
			var a *dns.Answer

			v := answer["answer"].(string)
			switch d.Get("type") {
			case "TXT", "SPF":
				a = dns.NewTXTAnswer(v)
			case "CAA":
				a = dns.NewAnswer(strings.SplitN(v, " ", 3))
			default:
				a = dns.NewAnswer(strings.Split(v, " "))
			}

			if v, ok := answer["region"]; ok {
				a.RegionName = v.(string)
			}

			if v, ok := answer["meta"]; ok {
				log.Println("answer meta", v)
				meta, err := metaHandler(v)
				if err != nil {
					return err
				}

				a.Meta = meta
			}

			r.AddAnswer(a)
		}
	}
	log.Println("number of answers found:", len(r.Answers))

	if v, ok := d.GetOk("ttl"); ok {
		r.TTL = v.(int)
	}

	if r.Type == "ALIAS" {
		Override_TTL := d.Get("override_ttl").(bool)
		r.Override_TTL = &Override_TTL
	}

	if v, ok := d.GetOk("link"); ok {
		if len(r.Answers) > 0 {
			return errors.New("cannot have both link and answers in a record")
		}
		r.LinkTo(v.(string))
	}

	if v, ok := d.GetOk("meta"); ok {
		log.Println("record meta", v)
		meta, err := metaHandler(v)
		if err != nil {
			return err
		}
		r.Meta = meta
		log.Println(r.Meta)
	}
	useClientSubnet := d.Get("use_client_subnet").(bool)
	r.UseClientSubnet = &useClientSubnet

	if rawFilters := d.Get("filters").([]interface{}); len(rawFilters) > 0 {
		filters := make([]*filter.Filter, len(rawFilters))
		for i, filterRaw := range rawFilters {
			fi := filterRaw.(map[string]interface{})
			config := make(map[string]interface{})
			f := filter.Filter{
				Type:   fi["filter"].(string),
				Config: config,
			}
			if disabled, ok := fi["disabled"]; ok {
				f.Disabled = disabled.(bool)
			}
			if rawConfig, ok := fi["config"]; ok {
				f.Config = rawConfig.(map[string]interface{})
			}
			filters[i] = &f
		}
		r.Filters = filters
	}
	if regions := d.Get("regions").([]interface{}); len(regions) > 0 {
		for _, regionRaw := range regions {
			region := regionRaw.(map[string]interface{})
			ns1R := data.Region{
				Meta: data.Meta{},
			}

			if v, ok := region["meta"]; ok {
				meta, err := metaHandler(v)
				if err != nil {
					return err
				}
				log.Println("region meta object", meta)

				ns1R.Meta = *meta
			}
			r.Regions[region["name"].(string)] = ns1R
		}
	}
	return nil
}

func removeEmptyMeta(v map[string]interface{}) {
	for metaKey, metaValue := range v {
		switch t := metaValue.(type) {
		case string:
			if t == "" {
				delete(v, metaKey)
			}
		}
	}
}

// RecordCreate creates DNS record in ns1
func RecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r := dns.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err := resourceDataToRecord(r, d); err != nil {
		return err
	}
	if resp, err := client.Records.Create(r); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	return recordToResourceData(d, r)
}

// RecordRead reads the DNS record from ns1
func RecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	r, resp, err := client.Records.Get(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		if err == ns1.ErrRecordMissing {
			log.Printf("[DEBUG] NS1 record (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}

	return recordToResourceData(d, r)
}

// RecordDelete deletes the DNS record from ns1
func RecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Records.Delete(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// RecordUpdate updates the given dns record in ns1
func RecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	r := dns.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err := resourceDataToRecord(r, d); err != nil {
		return err
	}
	if resp, err := client.Records.Update(r); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	return recordToResourceData(d, r)
}

func recordStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid record specifier.  Expecting 2 slashes (\"zone/domain/type\"), got %d", len(parts)-1)
	}

	d.Set("zone", parts[0])
	d.Set("domain", parts[1])
	d.Set("type", parts[2])

	return []*schema.ResourceData{d}, nil
}

// ignores case difference between state and resoruce written in terraform file
// so we can keep data consistency between tf state and ddi since ddi is case insensitive
func caseSensitivityDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// metaDiffSuppress evaluates fields in the meta attribute.
// fields that could be []string have diff suppressed if the difference is in ordering of elements,
// since the API often changes the order.
// boolean fields are normalized for string representations of bools.
func metaDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasSuffix(k, ".subdivisions") {
		// Converting and sorting string 'new'
		newSubdivMap, _ := subdivisionConverter(new)
		newJsonStr, _ := json.Marshal(newSubdivMap)

		// Converting and sorting string 'old'
		oldSubdivMap, _ := subdivisionConverter(old)
		oldJsonStr, _ := json.Marshal(oldSubdivMap)
		return string(oldJsonStr) == string(newJsonStr)
	}

	if strings.HasSuffix(k, ".georegion") ||
		strings.HasSuffix(k, ".country") ||
		strings.HasSuffix(k, ".us_state") ||
		strings.HasSuffix(k, ".ca_province") ||
		strings.HasSuffix(k, ".ip_prefixes") ||
		strings.HasSuffix(k, ".asn") {

		compareMap := make(map[string]bool)
		for _, value := range strings.Split(old, ",") {
			compareMap[strings.TrimSpace(value)] = true
		}
		for _, value := range strings.Split(new, ",") {
			value = strings.TrimSpace(value)
			if _, ok := compareMap[value]; ok {
				delete(compareMap, value)
			} else {
				return false
			}
		}

		return len(compareMap) == 0
	}

	if metaDiffSuppressUp(k, old, new, d) {
		return true
	}

	return false
}

// suppress a detected diff if it reflects an unchanged boolean value
func metaDiffSuppressUp(k, old, new string, _ *schema.ResourceData) bool {
	if strings.HasSuffix(k, "up") {
		newB, err := strconv.ParseBool(new)
		if err != nil {
			return false
		}
		oldB, err := strconv.ParseBool(old)
		if err != nil {
			return false
		}
		if newB == oldB {
			return true
		}
	}
	return false
}

// validateFQDN verifies that an FQDN doesn't have leading or trailing dots.
func validateFQDN(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	if strings.HasPrefix(v, ".") {
		errs = append(errs, fmt.Errorf("%s has an invalid leading \".\", got: %s", key, v))
	}

	if strings.HasSuffix(v, ".") {
		errs = append(errs, fmt.Errorf("%s has an invalid trailing \".\", got: %s", key, v))
	}

	return warns, errs
}

func metaToMapString(m *data.Meta) map[string]interface{} {
	stringMap := m.StringMap()
	if stringMap != nil {
		subdivisions := stringMap["subdivisions"]
		if subdivisions != nil {
			subd := m.Subdivisions.(map[string]interface{})
			keys := make([]string, 0, len(subd))
			var array []string
			var subArray interface{}
			for k := range subd {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, subRegion := range keys {
				subArray = subd[subRegion]
				for _, sub := range subArray.([]interface{}) {
					array = append(array, subRegion+"-"+sub.(string))
				}
			}
			stringMap["subdivisions"] = strings.Join(array[:], ",")
		}
	}
	return stringMap
}

func isJson(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func metaHandler(meta interface{}) (*data.Meta, error) {
	metaMap, ok := meta.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error to convert meta to map")
	}

	if v, ok := metaMap["subdivisions"]; ok {
		subdivisionsMap, err := subdivisionConverter(v.(string))
		if err != nil {
			return nil, err
		}

		metaMap["subdivisions"] = subdivisionsMap
	}

	removeEmptyMeta(metaMap)
	metaData := data.MetaFromMap(metaMap)
	log.Println(metaData)
	errs := metaData.Validate()
	if len(errs) > 0 {
		return nil, errJoin(append([]error{errors.New("found error/s in answer metadata")}, errs...), ",")
	}

	return metaData, nil
}

func subdivisionConverter(s string) (map[string]interface{}, error) {
	// Get a string in "Country-Subdivision" or Json format, sort and return a map[string]interface{}
	subdivisionsMap := make(map[string]interface{})
	subdivisions := strings.Split(s, ",")
	if isJson(s) {
		json.Unmarshal([]byte(s), &subdivisionsMap)
		for k, v := range subdivisionsMap {
			vals, ok := v.([]interface{})
			if !ok {
				return nil, fmt.Errorf("error to convert interface{} to []interface{}")
			}
			strSlice := []string{}
			for _, s := range vals {
				str, ok := s.(string)
				if !ok {
					return nil, fmt.Errorf("couldn't get subdivision value correctly")
				}
				strSlice = append(strSlice, str)
			}
			subdivisionsMap[k] = strSlice
		}
	} else {
		for _, sub := range subdivisions {
			sub = strings.Join(strings.Fields(sub), "")
			subp := strings.SplitN(sub, "-", 2)
			if len(subp) != 2 {
				return nil, fmt.Errorf("invalid subdivision format. expecting (\"Country-Subdivision\") got %s", sub)
			}
			if subdivisionsMap[subp[0]] == nil {
				subdivisionsMap[subp[0]] = []string{}
			}
			subdivisionsMap[subp[0]] = append(subdivisionsMap[subp[0]].([]string), subp[1])
		}
	}

	for _, slices := range subdivisionsMap {
		sort.Slice(slices, func(i, j int) bool {
			return slices.([]string)[i] < slices.([]string)[j]
		})
	}

	return subdivisionsMap, nil
}
