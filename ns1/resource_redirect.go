package ns1

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

var forwardingTypeStringEnum = NewStringEnum([]string{
	"permanent",
	"temporary",
	"masking",
})

var forwardingModeStringEnum = NewStringEnum([]string{
	"all",
	"capture",
	"none",
})

func redirectConfigResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"domain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateDomain,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"path": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validatePath,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			"target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateURL,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			// Read-only
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"https_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Optional
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"forwarding_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: forwardingModeStringEnum.ValidateFunc,
			},
			"forwarding_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "permanent",
				ValidateFunc: forwardingTypeStringEnum.ValidateFunc,
			},
			"https_forced": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"query_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Create:   RedirectConfigCreate,
		Read:     RedirectConfigRead,
		Update:   RedirectConfigUpdate,
		Delete:   RedirectConfigDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func redirectCertificateResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"domain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateDomain,
				DiffSuppressFunc: caseSensitivityDiffSuppress,
			},
			// Read-only
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"valid_from": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"valid_until": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"errors": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"processing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Create:   RedirectCertCreate,
		Read:     RedirectCertRead,
		Update:   RedirectCertUpdate,
		Delete:   RedirectCertDelete,
		Importer: &schema.ResourceImporter{},
	}
}

// RedirectConfigCreate creates a redirect configuration
func RedirectConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	var tags []string
	terraformTags := d.Get("tags").([]interface{})
	for _, t := range terraformTags {
		tags = append(tags, t.(string))
	}

	r := redirect.NewConfiguration(
		d.Get("domain").(string),
		d.Get("path").(string),
		d.Get("target").(string),
		tags,
		getFwModep(d, "forwarding_mode"),
		getFwTypep(d, "forwarding_type"),
		getBoolp(d, "https_enabled"),
		getBoolp(d, "https_forced"),
		getBoolp(d, "query_forwarding"),
	)

	cert := getStringp(d, "certificate_id")
	if cert != nil {
		_, _, err := client.RedirectCertificates.Get(*cert)
		if err == ns1.ErrRedirectCertificateNotFound {
			cert = nil
		}
	}
	if cert != nil {
		r.CertificateID = cert
		t := true
		r.HttpsEnabled = &t
	} else {
		f := false
		r.HttpsEnabled = &f
	}

	cfg, resp, err := client.Redirects.Create(r)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return redirectConfigToResourceData(d, cfg)
}

// RedirectConfigRead reads the redirect config from ns1
func RedirectConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	cfg, resp, err := client.Redirects.Get(d.Get("id").(string))
	if err != nil {
		if err == ns1.ErrRedirectNotFound {
			log.Printf("[DEBUG] NS1 redirect config (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return ConvertToNs1Error(resp, err)
	}

	return redirectConfigToResourceData(d, cfg)
}

// RedirectConfigDelete deletes the redirect config from ns1
func RedirectConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.Redirects.Delete(d.Get("id").(string))
	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// RedirectConfigUpdate updates the given redirect config in ns1
func RedirectConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	var tags []string
	terraformTags := d.Get("tags").([]interface{})
	if terraformTags != nil {
		tags = []string{}
	}
	for _, t := range terraformTags {
		tags = append(tags, t.(string))
	}

	r := redirect.NewConfiguration(
		d.Get("domain").(string),
		d.Get("path").(string),
		d.Get("target").(string),
		tags,
		getFwModep(d, "forwarding_mode"),
		getFwTypep(d, "forwarding_type"),
		getBoolp(d, "https_enabled"),
		getBoolp(d, "https_forced"),
		getBoolp(d, "query_forwarding"),
	)
	id := d.Id()
	r.ID = &id

	cert := getStringp(d, "certificate_id")
	if cert != nil {
		_, _, err := client.RedirectCertificates.Get(*cert)
		if err == ns1.ErrRedirectCertificateNotFound {
			cert = nil
		}
	}
	if cert != nil {
		r.CertificateID = cert
		t := true
		r.HttpsEnabled = &t
	} else {
		f := false
		r.HttpsEnabled = &f
	}

	cfg, resp, err := client.Redirects.Update(r)
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}
	return redirectConfigToResourceData(d, cfg)
}

// RedirectCertCreate creates a redirect certificate
func RedirectCertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)

	cert, resp, err := client.RedirectCertificates.Create(d.Get("domain").(string))
	if err != nil {
		return ConvertToNs1Error(resp, err)
	}

	return redirectCertToResourceData(d, cert)
}

// RedirectCertRead reads the redirect certificate from ns1
func RedirectCertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	id := d.Get("id").(string)

	cert, resp, err := client.RedirectCertificates.Get(id)
	if err == nil && cert.Errors != nil && *cert.Errors == "Revoking" {
		// wait for delete
		for i := 0; i < 20; i++ {
			// 20 x 500 milliseconds = max 10 seconds plus network delay
			time.Sleep(500 * time.Millisecond)
			_, _, err = client.RedirectCertificates.Get(id)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		if err == ns1.ErrRedirectCertificateNotFound {
			log.Printf("[DEBUG] NS1 redirect certificate (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return ConvertToNs1Error(resp, err)
	}

	return redirectCertToResourceData(d, cert)
}

// RedirectCertDelete deletes the redirect certificate from ns1
func RedirectCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	id := d.Get("id").(string)
	resp, err := client.RedirectCertificates.Delete(id)
	if err == nil {
		for i := 0; i < 20; i++ {
			// 20 x 500 milliseconds = max 10 seconds plus network delay
			time.Sleep(500 * time.Millisecond)
			_, _, err := client.RedirectCertificates.Get(id)
			if err != nil {
				break
			}
		}
	}

	d.SetId("")
	return ConvertToNs1Error(resp, err)
}

// RedirectCertUpdate updates the given redirect certificate in ns1
func RedirectCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ns1.Client)
	resp, err := client.RedirectCertificates.Update(d.Id())
	return ConvertToNs1Error(resp, err)
}

// validateDomain verifies that the string matches a valid FQDN.
func validateDomain(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	match, err := regexp.MatchString("^(\\*\\.)?([\\w-]+\\.)*[\\w-]+$", v)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s is invalid, got: %s, error: %e", key, v, err))
	}

	if !match {
		errs = append(errs, fmt.Errorf("%s is not a valid FQDN, got: %s", key, v))
	}

	return warns, errs
}

// validatePath verifies that the path matches a valid URL path.
func validatePath(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	match, err := regexp.MatchString("^[*]?[a-zA-Z0-9\\.\\-/$!+(_)' ]+[*]?$", v)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s is invalid, got: %s, error: %e", key, v, err))
	}

	if !match {
		errs = append(errs, fmt.Errorf("%s is not a valid FQDN, got: %s", key, v))
	}

	return warns, errs
}

// validateURL verifies that the string is a valid URL.
func validateURL(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	match, err := regexp.MatchString("^(http://|https://)?[a-zA-Z0-9\\.\\-/$!+(_)' ]+$", v)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s is invalid, got: %s, error: %e", key, v, err))
	}

	if !match {
		errs = append(errs, fmt.Errorf("%s is not a valid FQDN, got: %s", key, v))
	}

	return warns, errs
}

// return nil if the value is not set, a valid pointer if it is
func getBoolp(d *schema.ResourceData, key string) *bool {
	val := d.Get(key)
	if val != nil {
		ret := val.(bool)
		return &ret
	} else {
		return nil
	}
}

// return nil if the value is not set, a valid pointer if it is
func getStringp(d *schema.ResourceData, key string) *string {
	val := d.Get(key)
	if val != nil {
		ret := val.(string)
		if ret != "" {
			return &ret
		}
	}
	return nil
}

// return nil if the value is not set, a valid pointer if it is
func getFwTypep(d *schema.ResourceData, key string) *redirect.ForwardingType {
	val := d.Get(key).(string)
	ret, found := redirect.ParseForwardingType(val)
	if found {
		return &ret
	} else {
		return nil
	}
}

// return nil if the value is not set, a valid pointer if it is
func getFwModep(d *schema.ResourceData, key string) *redirect.ForwardingMode {
	val := d.Get(key).(string)
	ret, found := redirect.ParseForwardingMode(val)
	if found {
		return &ret
	} else {
		return nil
	}
}

func redirectConfigToResourceData(d *schema.ResourceData, r *redirect.Configuration) error {
	d.Set("domain", r.Domain)
	d.Set("path", r.Path)
	d.Set("target", r.Target)
	if r.ID != nil {
		d.SetId(*r.ID)
	}
	if r.CertificateID != nil {
		d.Set("certificate_id", *r.CertificateID)
	}
	if r.ForwardingMode != nil {
		d.Set("forwarding_mode", r.ForwardingMode.String())
	}
	if r.ForwardingType != nil {
		d.Set("forwarding_type", r.ForwardingType.String())
	}
	if r.HttpsEnabled != nil {
		d.Set("https_enabled", *r.HttpsEnabled)
	}
	if r.HttpsForced != nil {
		d.Set("https_forced", *r.HttpsForced)
	}
	if r.QueryForwarding != nil {
		d.Set("query_forwarding", *r.QueryForwarding)
	}
	if r.Tags != nil {
		d.Set("tags", r.Tags)
	}
	if r.LastUpdated != nil {
		d.Set("last_updated", *r.LastUpdated)
	}
	return nil
}

func redirectCertToResourceData(d *schema.ResourceData, r *redirect.Certificate) error {
	d.Set("domain", r.Domain)
	if r.ID != nil {
		d.SetId(*r.ID)
	}
	if r.Certificate != nil {
		d.Set("certificate", *r.Certificate)
	}
	if r.ValidFrom != nil {
		d.Set("valid_from", *r.ValidFrom)
	}
	if r.ValidUntil != nil {
		d.Set("valid_until", *r.ValidUntil)
	}
	if r.Errors != nil {
		d.Set("errors", *r.Errors)
	}
	if r.Processing != nil {
		d.Set("processing", *r.Processing)
	}
	if r.LastUpdated != nil {
		d.Set("last_updated", *r.LastUpdated)
	}
	return nil
}
