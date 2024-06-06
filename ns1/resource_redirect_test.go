package ns1

import (
	"fmt"
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
		PreCheck:     func() { testAccPreCheck(t) },
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
				),
			},
			{
				Config: testAccRedirectUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedirectConfigExists("ns1_redirect.it", &redirect),
					testAccCheckRedirectConfigDomain(&redirect, "test."+domainName),
					testAccCheckRedirectConfigFwType(&redirect, "permanent"),
					testAccCheckRedirectConfigTags(&redirect, []string{"test"}),
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
  https_enabled    = true
  https_forced     = true
  query_forwarding = true
  tags             = [ "test", "it" ]
}

resource "ns1_redirect_certificate" "example" {
  domain       = "*.${ns1_zone.test.zone}"
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
  target           = "https://url.com/target/path"
  forwarding_mode  = "capture"
  forwarding_type  = "permanent"
  https_enabled    = true
  https_forced     = true
  query_forwarding = true
  tags             = [ "test" ]
}

resource "ns1_redirect_certificate" "example" {
  domain       = "*.${ns1_zone.test.zone}"
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
			return fmt.Errorf("Name: got: %s want: %s", cfg.Domain, expected)
		}
		return nil
	}
}

func testAccCheckRedirectConfigFwType(cfg *redirect.Configuration, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cfg.ForwardingType.String() != expected {
			return fmt.Errorf("Name: got: %s want: %s", cfg.ForwardingType.String(), expected)
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
			for i, _ := range expected {
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
