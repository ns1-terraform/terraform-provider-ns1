package ns1

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

func TestAccRedirectConfig_basic(t *testing.T) {
	var redirect redirect.Configuration
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	domainName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccRedirectPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedirectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigTarget(&redirect, "https://url.com/target/path"),
					testAccCheckRedirectConfigFwType(&redirect, "masking"),
					testAccCheckRedirectConfigTags(&redirect, []string{"test", "it"}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
			{
				Config: testAccRedirectUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigTarget(&redirect, "https://url.com/target/path?q=param#frag"),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
			{
				ResourceName:      "ns1_redirect.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "ns1_redirect_certificate.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedirectConfig_http_to_https(t *testing.T) {
	var redirect redirect.Configuration
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	domainName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccRedirectPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedirectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectHTTP(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigTarget(&redirect, "https://url.com/target/path"),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{}),
					testAccCheckRedirectConfigHTTPS(&redirect, false),
					testAccCheckRedirectConfigCertIdPresent(&redirect, false),
				),
			},
			{
				Config: testAccRedirectUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigTarget(&redirect, "https://url.com/target/path?q=param#frag"),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
			{
				Config: testAccRedirectHTTPwithCert(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigTarget(&redirect, "https://url.com/target/path"),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
		},
	})
}

func TestAccRedirectConfig_https_to_http(t *testing.T) {
	var redirect redirect.Configuration
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	domainName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccRedirectPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedirectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigFwType(&redirect, "masking"),
					testAccCheckRedirectConfigTags(&redirect, []string{"test", "it"}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
			{
				Config: testAccRedirectHTTP(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{}),
					testAccCheckRedirectConfigHTTPS(&redirect, false),
					testAccCheckRedirectConfigCertIdPresent(&redirect, false),
				),
			},
		},
	})
}

func TestAccRedirectConfig_remoteChanges(t *testing.T) {
	var redirect redirect.Configuration
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	domainName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccRedirectPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRedirectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigFwType(&redirect, "masking"),
					testAccCheckRedirectConfigTags(&redirect, []string{"test", "it"}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
			{
				PreConfig: eraseAll(t, domainName),
				Config:    testAccRedirectBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigFwType(&redirect, "masking"),
					testAccCheckRedirectConfigTags(&redirect, []string{"test", "it"}),
					testAccCheckRedirectConfigHTTPS(&redirect, true),
					testAccCheckRedirectConfigCertIdPresent(&redirect, true),
				),
			},
		},
	})
}

func testAccRedirectBasic(rString string) string {
	return fmt.Sprintf(`
resource "ns1_redirect" "it" {
  certificate_id   = "${ns1_redirect_certificate.example.id}"
  domain           = "test.${ns1_zone.test.zone}"
  path             = "/from/path/*"
  target           = "https://url.com/target/path"
  forwarding_mode  = "capture"
  forwarding_type  = "masking"
  https_forced     = true
  query_forwarding = true
  tags             = [ "test", "it" ]
}

resource "ns1_redirect_certificate" "example" {
  domain = "*.${ns1_zone.test.zone}"
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRedirectUpdated(rString string) string {
	return fmt.Sprintf(`
resource "ns1_redirect" "it" {
  certificate_id   = "${ns1_redirect_certificate.example.id}"
  domain           = "test.${ns1_zone.test.zone}"
  path             = "/from/path/*"
  target           = "https://url.com/target/path?q=param#frag"
  forwarding_mode  = "capture"
  forwarding_type  = "permanent"
  https_forced     = true
  query_forwarding = true
  tags             = [ ]
}

resource "ns1_redirect_certificate" "example" {
  domain = "*.${ns1_zone.test.zone}"
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRedirectHTTP(rString string) string {
	return fmt.Sprintf(`
resource "ns1_redirect" "it" {
  domain           = "test.${ns1_zone.test.zone}"
  path             = "/from/path/*"
  target           = "https://url.com/target/path"
  forwarding_mode  = "all"
  forwarding_type  = "permanent"
  https_forced     = false
  tags             = [ ]
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRedirectHTTPwithCert(rString string) string {
	return fmt.Sprintf(`
resource "ns1_redirect" "it" {
  certificate_id   = ""
  domain           = "test.${ns1_zone.test.zone}"
  path             = "/from/path/*"
  target           = "https://url.com/target/path"
  forwarding_mode  = "all"
  forwarding_type  = "permanent"
  tags             = [ ]
}

resource "ns1_redirect_certificate" "example" {
  domain = "*.${ns1_zone.test.zone}"
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccCheckRedirectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	var id string
	var certId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_redirect" && rs.Type != "ns1_redirect_certificate" {
			continue
		}

		if rs.Type == "ns1_redirect" {
			id = rs.Primary.ID
		}

		if rs.Type == "ns1_redirect_certificate" {
			certId = rs.Primary.ID
		}
	}

	if id != "" {
		foundRecord, _, err := client.Redirects.Get(id)
		if err != ns1.ErrRedirectNotFound {
			return fmt.Errorf("redirect still exists: %#v %#v", foundRecord, err)
		}
	}
	if certId != "" {
		foundRecord, _, err := client.RedirectCertificates.Get(certId)
		if err != ns1.ErrRedirectCertificateNotFound {
			return fmt.Errorf("certificate still exists: %#v %#v", foundRecord, err)
		}
	}

	return nil
}

func testAccCheckRedirectConfigExists(n string, cfg *redirect.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %v", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("noID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		p := rs.Primary

		foundCfg, _, err := client.Redirects.Get(p.Attributes["id"])
		if err != nil {
			return fmt.Errorf("redirect not found")
		}

		if foundCfg.ID == nil || *foundCfg.ID != p.Attributes["id"] {
			return fmt.Errorf("redirect not found")
		}

		*cfg = *foundCfg

		return nil
	}
}

func testAccCheckRedirectConfigDomain(cfg *redirect.Configuration, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cfg.Domain != expected {
			return fmt.Errorf("Domain: got: %s want: %s", cfg.Domain, expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigTarget(cfg *redirect.Configuration, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cfg.Target != expected {
			return fmt.Errorf("Target: got: %s want: %s", cfg.Domain, expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigFwType(cfg *redirect.Configuration, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cfg.ForwardingType.String() != expected {
			return fmt.Errorf("ForwardingType: got: %s want: %s", cfg.ForwardingType.String(), expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigHTTPS(cfg *redirect.Configuration, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *cfg.HttpsEnabled != expected {
			return fmt.Errorf("HttpsEnabled: got: %t want: %t", *cfg.HttpsEnabled, expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigCertIdPresent(cfg *redirect.Configuration, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (cfg.CertificateID != nil) != expected {
			return fmt.Errorf("CertificateID present: got: %t want: %t", cfg.CertificateID != nil, expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigTags(cfg *redirect.Configuration, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		diff := false
		if len(cfg.Tags) != len(expected) {
			diff = true
		} else {
			sort.Strings(expected)
			sort.Strings(cfg.Tags)
			for i := range expected {
				if cfg.Tags[i] != expected[i] {
					diff = true
				}
			}
		}
		if diff {
			return fmt.Errorf("Name: got: %s want: %s", strings.Join(cfg.Tags, ","), strings.Join(expected, ","))
		}
		return nil
	}
}

func eraseAll(t *testing.T, domain string) func() {
	return func() {
		client, err := sharedClient()
		if err != nil {
			t.Fatalf("failed to get shared client: %e", err)
		}
		// delete all configs
		redirects := client.Redirects
		cfgs, _, err := redirects.List()
		if err != nil {
			t.Fatalf("failed to get list of redirects: %e", err)
		}
		for _, v := range cfgs {
			if v.ID != nil && strings.HasSuffix(v.Domain, domain) {
				_, err = redirects.Delete(*v.ID)
				if err != nil {
					t.Fatalf("failed to delete redirect %s: %e", v.Domain+"/"+v.Path, err)
				}
			}
		}
		// delete all certs
		certificates := client.RedirectCertificates
		certs, _, err := certificates.List()
		if err != nil {
			t.Fatalf("failed to get list of certificates: %e", err)
		}
		for _, v := range certs {
			if v.ID != nil && strings.HasSuffix(v.Domain, domain) {
				_, err = certificates.Delete(*v.ID)
				if err != nil {
					t.Fatalf("failed to delete cert %s: %e", v.Domain, err)
				}
			}
		}
	}
}

// See if we have redirect permissions by trying to list redirects
func testAccRedirectPreCheck(t *testing.T) {
	client, err := sharedClient()
	if err != nil {
		log.Fatalf("failed to get shared client: %s", err)
	}
	_, _, err = client.Redirects.List()
	if err != nil {
		t.Skipf("account not authorized for redirects, skipping test")
	}
}
