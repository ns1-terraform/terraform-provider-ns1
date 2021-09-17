package ns1

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
)

func pulsarJobResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"customer": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateTypeId,
		},
		"community": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"job_id": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"active": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"shared": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"config": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"url_path": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"http": {
						Type:     schema.TypeBool,
						Optional: true,
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
						Default:  false,
					},
					"static_values": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"blend_metric_weights": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"timestamp": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"weights": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"weight": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"default_value": {
						Type:     schema.TypeFloat,
						Required: true,
					},
					"maximize": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        s,
		Create:        PulsarJobCreate,
		Read:          pulsarJobRead,
		Update:        PulsarJobUpdate,
		Delete:        pulsarJobDelete,
		Importer:      &schema.ResourceImporter{State: pulsarJobImportStateFunc},
		SchemaVersion: 1,
	}
}

func pulsarJobToResourceData(d *schema.ResourceData, j *pulsar.PulsarJob) error {
	d.SetId(j.JobID)
	d.Set("customer", j.Customer)
	d.Set("name", j.Name)
	d.Set("type_id", j.TypeID)
	d.Set("community", j.Community)
	d.Set("app_id", j.AppID)
	d.Set("active", j.Active)
	d.Set("shared", j.Shared)

	if j.Config != nil {
		if err := jobConfigToResourceData(d, j); err != nil {
			return err
		}

	}

	return nil
}

func jobConfigToResourceData(d *schema.ResourceData, j *pulsar.PulsarJob) error {
	config := make(map[string]interface{})
	if j.Config.Host != nil {
		config["host"] = *j.Config.Host
	}
	if j.Config.URL_Path != nil {
		config["url_path"] = *j.Config.URL_Path
	}
	if j.Config.Http != nil {
		config["http"] = strconv.FormatBool(*j.Config.Http)
	}
	if j.Config.Https != nil {
		config["https"] = strconv.FormatBool(*j.Config.Https)
	}
	if j.Config.RequestTimeoutMillis != nil {
		config["request_timeout_millis"] = strconv.Itoa(*j.Config.RequestTimeoutMillis)
	}
	if j.Config.RequestTimeoutMillis != nil {
		config["job_timeout_millis"] = strconv.Itoa(*j.Config.JobTimeoutMillis)
	}
	if j.Config.UseXHR != nil {
		config["use_xhr"] = strconv.FormatBool(*j.Config.UseXHR)
	}
	if j.Config.StaticValues != nil {
		config["static_values"] = strconv.FormatBool(*j.Config.StaticValues)
	}

	if j.Config.BlendMetricWeights != nil {
		if err := d.Set("blend_metric_weights", blendMetricWeightsToMap(j.Config.BlendMetricWeights)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Blend Metric Weights for: %s, error: %#v", j.Name, err)
		}

		if w := j.Config.BlendMetricWeights.Weights; w != nil {
			weights := make([]map[string]interface{}, 0)
			for _, weight := range w {
				weights = append(weights, weightsToResourceData(weight))
			}

			if err := d.Set("weights", weights); err != nil {
				return fmt.Errorf("[DEBUG] Error setting weights for: %s, error: %#v", j.Name, err)
			}
		}
	}

	d.Set("config", config)
	return nil
}

func blendMetricWeightsToMap(b *pulsar.BlendMetricWeights) map[string]interface{} {
	blendMetric := make(map[string]interface{})
	blendMetric["timestamp"] = strconv.Itoa(b.Timestamp)

	return blendMetric
}

func weightsToResourceData(w *pulsar.Weights) map[string]interface{} {
	weight := make(map[string]interface{})
	weight["name"] = w.Name
	weight["weight"] = w.Weight
	weight["default_value"] = w.DefaultValue
	weight["maximize"] = w.Maximize

	return weight
}

func resourceDataToPulsarJob(j *pulsar.PulsarJob, d *schema.ResourceData) error {
	j.Name = d.Get("name").(string)
	j.TypeID = d.Get("type_id").(string)
	j.JobID = d.Id()
	j.AppID = d.Get("app_id").(string)
	j.Active = d.Get("active").(bool)
	j.Shared = d.Get("shared").(bool)

	if v, ok := d.GetOk("config"); ok {
		if config, err := resourceDataToJobConfig(v); err != nil {
			return err
		} else {
			j.Config = config
		}
	}

	if v, ok := d.GetOk("blend_metric_weights"); ok {
		if j.Config == nil {
			j.Config = &pulsar.JobConfig{}
		}

		if m, err := resourceDataToBlendMetric(v); err != nil {
			return err
		} else {
			j.Config.BlendMetricWeights = m
		}

		if v, ok := d.GetOk("weights"); ok {
			j.Config.BlendMetricWeights.Weights = resourceDataToWeights(v)
		}
	}

	return nil
}

func resourceDataToJobConfig(v interface{}) (*pulsar.JobConfig, error) {
	rawconfig := v.(map[string]interface{})
	j := &pulsar.JobConfig{}

	if v, ok := rawconfig["host"]; ok {
		host := v.(string)
		j.Host = &host
	}
	if v, ok := rawconfig["url_path"]; ok {
		url_path := v.(string)
		j.URL_Path = &url_path
	}
	if v, ok := rawconfig["http"]; ok {
		if b, err := strconv.ParseBool(v.(string)); err != nil {
			return nil, err
		} else {
			j.Http = &b
		}
	}
	if v, ok := rawconfig["https"]; ok {
		if b, err := strconv.ParseBool(v.(string)); err != nil {
			return nil, err
		} else {
			j.Https = &b
		}
	}
	if v, ok := rawconfig["request_timeout_millis"]; ok {
		if i, err := strconv.Atoi(v.(string)); err != nil {
			return nil, err
		} else {
			j.RequestTimeoutMillis = &i
		}
	}
	if v, ok := rawconfig["job_timeout_millis"]; ok {
		if i, err := strconv.Atoi(v.(string)); err != nil {
			return nil, err
		} else {
			j.JobTimeoutMillis = &i
		}
	}
	if v, ok := rawconfig["use_xhr"]; ok {
		if b, err := strconv.ParseBool(v.(string)); err != nil {
			return nil, err
		} else {
			j.UseXHR = &b
		}
	}
	if v, ok := rawconfig["static_values"]; ok {
		if b, err := strconv.ParseBool(v.(string)); err != nil {
			return nil, err
		} else {
			j.StaticValues = &b
		}
	}
	return j, nil
}

func resourceDataToBlendMetric(v interface{}) (*pulsar.BlendMetricWeights, error) {
	rawMetric := v.(map[string]interface{})
	m := &pulsar.BlendMetricWeights{}

	if v, ok := rawMetric["timestamp"]; ok {
		if i, err := strconv.Atoi(v.(string)); err != nil {
			return nil, err
		} else {
			m.Timestamp = i
		}
	}

	m.Weights = make([]*pulsar.Weights, 0)

	return m, nil
}

func resourceDataToWeights(v interface{}) []*pulsar.Weights {
	weightsSet := v.(*schema.Set)
	weights := make([]*pulsar.Weights, weightsSet.Len())

	for i, weightRaw := range weightsSet.List() {
		weight := weightRaw.(map[string]interface{})
		weights[i] = &pulsar.Weights{
			Name:         weight["name"].(string),
			Weight:       weight["weight"].(int),
			DefaultValue: weight["default_value"].(float64),
			Maximize:     weight["maximize"].(bool),
		}
	}

	return weights
}

// PulsarJobCreate creates the given Pulsar Job in ns1
func PulsarJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	j := pulsar.PulsarJob{}
	if err := resourceDataToPulsarJob(&j, d); err != nil {
		return err
	}
	if resp, err := client.PulsarJobs.Create(&j); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return pulsarJobToResourceData(d, &j)
}

// pulsarJobRead reads the given zone data from ns1
func pulsarJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	j, resp, err := client.PulsarJobs.Get(d.Get("app_id").(string), d.Id())
	if err != nil {
		if errors.Is(err, ns1.ErrAppMissing) {
			log.Printf("[DEBUG] NS1 Pulsar Application (%s) not found", d.Get("app_id"))
			d.SetId("")
			return nil
		}

		if errors.Is(err, ns1.ErrJobMissing) {
			log.Printf("[DEBUG] NS1 Pulsar Job (%s) not found", d.Get("job_id"))
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}
	// Set Terraform resource data from the job data we just downloaded
	if err := pulsarJobToResourceData(d, j); err != nil {
		return err
	}
	return nil
}

// PulsarJobUpdate updates the Pulsar Job with given parameters in ns1
func PulsarJobUpdate(job_schema *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	j := pulsar.PulsarJob{
		JobID: job_schema.Id(),
		AppID: job_schema.Get("app_id").(string),
	}
	if err := resourceDataToPulsarJob(&j, job_schema); err != nil {
		return err
	}

	if resp, err := client.PulsarJobs.Update(&j); err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return pulsarJobToResourceData(job_schema, &j)
}

// pulsarJobDelete deletes the given Pulsar Job from ns1
func pulsarJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	j := pulsar.PulsarJob{}
	resourceDataToPulsarJob(&j, d)
	resp, err := client.PulsarJobs.Delete(&j)
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

func validateTypeId(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v != "latency" && v != "custom" {
		errs = append(errs,
			fmt.Errorf(
				"type_id %s invalid, please select between 'latency' or 'custom'", v,
			),
		)
	}
	return warns, errs
}

func pulsarJobImportStateFunc(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid job specifier. Expected 2 ids (\"app_id\"_\"job_id\", got %d)", len(parts))
	}

	d.Set("app_id", parts[0])
	d.Set("job_id", parts[1])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
