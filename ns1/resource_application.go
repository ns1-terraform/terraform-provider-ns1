package ns1

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
	"strconv"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"browser_wait_millis": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"jobs_per_transaction": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"https": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"request_timeout_millis": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"job_timeout_millis": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"use_xhr": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"static_values": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		Create: ApplicationCreate,
		Read:   ApplicationRead,
		Update: ApplicationUpdate,
		Delete: ApplicationDelete,
	}
}

func resourceApplicationToResourceData(d *schema.ResourceData, a *pulsar.Application) error {
	d.SetId(a.ID)
	d.Set("name", a.Name)
	d.Set("active", a.Active)
	d.Set("browser_wait_millis", a.BrowserWaitMillis)
	d.Set("jobs_per_transaction", a.JobsPerTransaction)

	d.Set("default_config", defaultConfigToMap(&a.DefaultConfig))
	return nil
}
func defaultConfigToMap(d *pulsar.DefaultConfig) map[string]interface{} {
	dm := make(map[string]interface{})
	dm["http"] = d.Http
	dm["https"] = d.Https
	dm["request_timeout_millis"] = d.RequestTimeoutMillis
	dm["job_timeout_millis"] = d.JobTimeoutMillis
	dm["use_xhr"] = d.UseXhr
	dm["static_values"] = d.StaticValues
	return dm
}

func resourceDataToApplication(a *pulsar.Application, d *schema.ResourceData) {

	a.ID = d.Id()
	if v, ok := d.GetOk("name"); ok {
		a.Name = v.(string)
	}
	if v, ok := d.GetOk("active"); ok {
		a.Active = v.(bool)
	}
	if v, ok := d.GetOk("browser_wait_millis"); ok {
		a.BrowserWaitMillis = v.(int)
	}
	if v, ok := d.GetOk("jobs_per_transaction"); ok {
		a.JobsPerTransaction = v.(int)
	}

	if v, ok := d.GetOk("default_config"); ok {
		defaultConfigSet := v.(map[string]interface{})
		a.DefaultConfig = pulsar.DefaultConfig{}
		a.DefaultConfig.Http, _ = strconv.ParseBool(defaultConfigSet["http"].(string))
		a.DefaultConfig.Https, _ = strconv.ParseBool(defaultConfigSet["https"].(string))
		a.DefaultConfig.RequestTimeoutMillis, _ = strconv.Atoi(defaultConfigSet["request_timeout_millis"].(string))
		a.DefaultConfig.JobTimeoutMillis, _ = strconv.Atoi(defaultConfigSet["job_timeout_millis"].(string))
		a.DefaultConfig.UseXhr, _ = strconv.ParseBool(defaultConfigSet["use_xhr"].(string))
		a.DefaultConfig.StaticValues, _ = strconv.ParseBool(defaultConfigSet["static_values"].(string))
	}
}

// applicationCreate creates the given zone in ns1
func ApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	a := pulsar.NewApplication(d.Get("name").(string))
	resourceDataToApplication(a, d)
	if resp, err := client.Applications.Create(a); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceApplicationToResourceData(d, a); err != nil {
		return err
	}
	return nil
}

// applicationRead reads the given zone data from ns1
func ApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	app, _, _ := client.Applications.Get(d.Id())
	if err := resourceApplicationToResourceData(d, app); err != nil {
		return err
	}
	return nil
}

// applicationDelete deletes the given zone from ns1
func ApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Applications.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// applicationUpdate updates the zone with given params in ns1
func ApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	app := pulsar.NewApplication(d.Get("name").(string))
	resourceDataToApplication(app, d)
	if resp, err := client.Applications.Update(app); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceApplicationToResourceData(d, app); err != nil {
		return err
	}
	return nil
}
