package ns1

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
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
		Create:   ApplicationCreate,
		Read:     ApplicationRead,
		Update:   ApplicationUpdate,
		Delete:   ApplicationDelete,
		Importer: &schema.ResourceImporter{State: ApplicationStateFunc},
	}
}

func resourceApplicationToResourceData(d *schema.ResourceData, a *pulsar.Application) error {
	d.SetId(a.ID)
	d.Set("name", a.Name)
	d.Set("active", a.Active)
	d.Set("browser_wait_millis", a.BrowserWaitMillis)
	d.Set("jobs_per_transaction", a.JobsPerTransaction)

	d.Set("default_config", []map[string]interface{}{defaultConfigToMap(&a.DefaultConfig)})
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
		a.DefaultConfig = setDefaultConfig(v)
	}
}

func setDefaultConfig(ds interface{}) (d pulsar.DefaultConfig) {
	defaultConfig := ds.([]interface{})[0].(map[string]interface{})
	d = pulsar.DefaultConfig{}
	httpConf := defaultConfig["http"]
	if httpConf != nil {
		d.Http = httpConf.(bool)
	}
	httpsConf := defaultConfig["https"]
	if httpsConf != nil {
		d.Https = httpsConf.(bool)
	}
	xhrConf := defaultConfig["use_xhr"]
	if xhrConf != nil {
		d.UseXhr = xhrConf.(bool)
	}
	staticConf := defaultConfig["static_values"]
	if staticConf != nil {
		d.StaticValues = staticConf.(bool)
	}
	jobConf := defaultConfig["job_timeout_millis"]
	if jobConf != nil {
		d.JobTimeoutMillis = jobConf.(int)
	}
	reqConf := defaultConfig["request_timeout_millis"]
	if reqConf != nil {
		d.RequestTimeoutMillis = reqConf.(int)
	}
	return
}

// ApplicationCreate creates the given application in ns1
func ApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	app := pulsar.NewApplication(d.Get("name").(string))
	resourceDataToApplication(app, d)
	if resp, err := client.Applications.Create(app); err != nil {
		return ConvertToNs1Error(resp, err)
	}
	if err := resourceApplicationToResourceData(d, app); err != nil {
		return err
	}
	return nil
}

// ApplicationRead reads the given application data from ns1
func ApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	app, resp, err := client.Applications.Get(d.Id())
	if err != nil {
		if errors.Is(err, ns1.ErrApplicationMissing) {
			log.Printf("[DEBUG] NS1 application (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return ConvertToNs1Error(resp, err)
	}

	if err := resourceApplicationToResourceData(d, app); err != nil {
		return err
	}
	return nil
}

// ApplicationDelete deletes the given application from ns1
func ApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Applications.Delete(d.Id())
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// ApplicationUpdate updates the application with given params in ns1
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

func ApplicationStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
