package ns1

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// Creating basic DNS view
func TestAccDNSView_basic(t *testing.T) {
	var (
		view           = dns.DNSView{}
		viewName       = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		viewPreference = 10
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSViewBasic(viewName, viewPreference),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSViewExists("ns1_dnsview.it", &view),
					testAccCheckDNSViewName(&view, viewName),
					testAccCheckDNSViewPreference(&view, viewPreference),
				),
			},
			{
				ResourceName:      "ns1_dnsview.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Update DNS view
func TestAccDNSView_update(t *testing.T) {
	var (
		view           = dns.DNSView{}
		viewName       = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		viewPreference = 10

		updatedView           = dns.DNSView{}
		updatedViewPreference = 5

		zoneName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSViewBasic(viewName, viewPreference),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSViewExists("ns1_dnsview.it", &view),
					testAccCheckDNSViewName(&view, viewName),
					testAccCheckDNSViewPreference(&view, viewPreference),
				),
			},
			{
				Config: testAccDNSViewUpdated(zoneName, viewName, updatedViewPreference),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSViewExists("ns1_dnsview.it", &updatedView),
					testAccCheckDNSViewName(&updatedView, viewName),
					testAccCheckDNSViewPreference(&updatedView, updatedViewPreference),
					testAccCheckDNSViewZones(&updatedView, []string{zoneName}),
				),
			},
			{
				ResourceName:      "ns1_dnsview.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Manually deleting DNS View
func TestAccDNSView_ManualDelete(t *testing.T) {
	var (
		view           = dns.DNSView{}
		viewName       = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		viewPreference = 10
	)
	// Manual deletion test for DNS View
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSViewBasic(viewName, viewPreference),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSViewExists("ns1_dnsview.it", &view),
				),
			},
			// Simulate a manual deletion of the DNS View and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteDNSView(&view),
				Config:             testAccDNSViewBasic(viewName, viewPreference),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccDNSViewBasic(viewName, viewPreference),
				Check:  testAccCheckDNSViewExists("ns1_dnsview.it", &view),
			},
		},
	})
}

func testAccDNSViewBasic(viewName string, viewPreference int) string {
	return fmt.Sprintf(`resource "ns1_dnsview" "it" {
  		name = "%s"
		preference = %d
}
`, viewName, viewPreference)
}

func testAccDNSViewUpdated(zoneName, viewName string, viewPreference int) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
		zone = "%s"
	  }
	
	resource "ns1_dnsview" "it" {
  		name = "%s"
		preference = %d
		zones = ["${ns1_zone.it.zone}"]
}
`, zoneName, viewName, viewPreference)
}

func testAccCheckDNSViewName(view *dns.DNSView, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if view.Name != expected {
			return fmt.Errorf("view.Name: got: %s want: %s", view.Name, expected)
		}
		return nil
	}
}

func testAccCheckDNSViewPreference(view *dns.DNSView, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if view.Preference != expected {
			return fmt.Errorf("view.Preference: got: %d want: %d", view.Preference, expected)
		}
		return nil
	}
}

func testAccCheckDNSViewZones(view *dns.DNSView, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !reflect.DeepEqual(view.Zones, expected) {
			return fmt.Errorf("view.Zones: got: %v want: %v", view.Zones, expected)
		}
		return nil
	}
}

func testAccCheckDNSViewDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_dnsview" {
			continue
		}

		dnsView, _, err := client.View.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("DNS View still exists: %#v: %#v", err, dnsView)
		}

	}

	return nil
}

func testAccCheckDNSViewExists(n string, view *dns.DNSView) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundDNSView, _, err := client.View.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		*view = *foundDNSView

		return nil
	}
}

func testAccManualDeleteDNSView(view *dns.DNSView) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.View.Delete(view.Name)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete DNS view: %v", err)
		}
	}
}
