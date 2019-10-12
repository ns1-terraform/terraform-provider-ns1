package ns1

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccZone_basic(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneTTL(&zone, 3600),
					testAccCheckZoneRefresh(&zone, 43200),
					testAccCheckZoneRetry(&zone, 7200),
					testAccCheckZoneExpiry(&zone, 1209600),
					testAccCheckZoneNxTTL(&zone, 3600),
					testAccCheckZoneNotPrimary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
					testAccCheckNSRecord(&zone, true),
				),
			},
		},
	})
}

func TestAccZone_updated(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneTTL(&zone, 3600),
					testAccCheckZoneRefresh(&zone, 43200),
					testAccCheckZoneRetry(&zone, 7200),
					testAccCheckZoneExpiry(&zone, 1209600),
					testAccCheckZoneNxTTL(&zone, 3600),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			{
				Config: testAccZoneUpdated(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneTTL(&zone, 10800),
					testAccCheckZoneRefresh(&zone, 3600),
					testAccCheckZoneRetry(&zone, 300),
					testAccCheckZoneExpiry(&zone, 2592000),
					testAccCheckZoneNxTTL(&zone, 3601),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			{
				ResourceName:      "ns1_zone.it",
				ImportState:       true,
				ImportStateId:     zoneName,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccZone_primary_to_secondary_to_normal(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	// sorted by IP please
	expected := []*dns.ZoneSecondaryServer{
		&dns.ZoneSecondaryServer{
			NetworkIDs: []int{0},
			IP:         "2.2.2.2",
			Port:       53,
			Notify:     false,
		},
		&dns.ZoneSecondaryServer{
			NetworkIDs: []int{0},
			IP:         "3.3.3.3",
			Port:       5353,
			Notify:     true,
		},
	}
	expectedOtherPorts := []int{53, 53}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// Start with a Primary Zone
			{
				Config: testAccZonePrimary(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneSecondaries(t, &zone, 0, expected[0]),
					testAccCheckZoneSecondaries(t, &zone, 1, expected[1]),
					testAccCheckZoneNotSecondary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			// Check import
			{
				ResourceName:      "ns1_zone.it",
				ImportState:       true,
				ImportStateId:     zoneName,
				ImportStateVerify: true,
			},
			// should make the zone secondary
			{
				Config: testAccZoneSecondary(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					resource.TestCheckResourceAttr("ns1_zone.it", "primary", "1.1.1.1"),
					resource.TestCheckResourceAttr(
						"ns1_zone.it", "additional_primaries.0", "2.2.2.2",
					),
					resource.TestCheckResourceAttr(
						"ns1_zone.it", "additional_primaries.1", "3.3.3.3",
					),
					testAccCheckOtherPorts(&zone, expectedOtherPorts),
					testAccCheckZoneNotPrimary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			// should correctly clear zone.Primary
			{
				Config: testAccZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneNotPrimary(&zone),
					testAccCheckZoneNotSecondary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
		},
	})
}

func TestAccZone_secondary_to_primary_to_normal(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	// sorted by IP please
	expected := []*dns.ZoneSecondaryServer{
		&dns.ZoneSecondaryServer{
			NetworkIDs: []int{0},
			IP:         "2.2.2.2",
			Port:       53,
			Notify:     false,
		},
		&dns.ZoneSecondaryServer{
			NetworkIDs: []int{0},
			IP:         "3.3.3.3",
			Port:       5353,
			Notify:     true,
		},
	}
	expectedOtherPorts := []int{53, 53}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// Start with a secondary zone
			{
				Config: testAccZoneSecondary(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					resource.TestCheckResourceAttr("ns1_zone.it", "primary", "1.1.1.1"),
					resource.TestCheckResourceAttr("ns1_zone.it", "additional_primaries.0", "2.2.2.2"),
					resource.TestCheckResourceAttr("ns1_zone.it", "additional_primaries.1", "3.3.3.3"),
					testAccCheckOtherPorts(&zone, expectedOtherPorts),
					testAccCheckZoneNotPrimary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			// Check import
			{
				ResourceName:      "ns1_zone.it",
				ImportState:       true,
				ImportStateId:     zoneName,
				ImportStateVerify: true,
			},
			// should make the zone primary
			{
				Config: testAccZonePrimary(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneSecondaries(t, &zone, 0, expected[0]),
					testAccCheckZoneSecondaries(t, &zone, 1, expected[1]),
					testAccCheckZoneNotSecondary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
			// should correctly set zone.Primary disabled
			{
				Config: testAccZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneNotPrimary(&zone),
					testAccCheckZoneNotSecondary(&zone),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
		},
	})
}

func TestAccZone_dnssec(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccDNSSECPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneDNSSEC(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneDNSSEC(&zone, true),
				),
			},
			{
				ResourceName:      "ns1_zone.it",
				ImportState:       true,
				ImportStateId:     zoneName,
				ImportStateVerify: true,
			},
			{
				Config: testAccZoneDNSSECUpdated(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneDNSSEC(&zone, false),
				),
			},
		},
	})
}

func TestAccZone_disable_autogenerate_ns_record(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneDisableAutoGenerateNSRecord(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckNSRecord(&zone, false),
				),
			},
		},
	})
}

// A Client instance we can use outside of a TestStep
func sharedClient() (*ns1.Client, error) {
	var ignoreSSL bool
	v := os.Getenv("NS1_IGNORE_SSL")
	if v == "" {
		ignoreSSL = false
	} else {
		v, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
		ignoreSSL = v
	}
	config := &Config{
		Key:       os.Getenv("NS1_APIKEY"),
		Endpoint:  os.Getenv("NS1_ENDPOINT"),
		IgnoreSSL: ignoreSSL,
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// See if we have DNSSEC permission by trying to create a zone with it
func testAccDNSSECPreCheck(t *testing.T) {
	client, err := sharedClient()
	if err != nil {
		log.Fatalf("failed to get shared client: %s", err)
	}
	zoneName := fmt.Sprintf(
		"terraform-dnssec-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	dnssec := true
	_, err = client.Zones.Create(
		&dns.Zone{Zone: zoneName, DNSSEC: &dnssec},
	)
	if err != nil {
		if strings.Contains(err.Error(), "400 DNSSEC support is not enabled") {
			t.Skipf("account not authorized for DNSSEC changes, skipping test")
			return
		}
		log.Fatalf("failed to create test zone %s: %s", zoneName, err)
	}
	_, err = client.Zones.Delete(zoneName)
	if err != nil {
		log.Fatalf("failed to delete test zone %s", zoneName)
	}
}

func testAccCheckZoneExists(n string, zone *dns.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NoID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundZone, _, err := client.Zones.Get(rs.Primary.Attributes["zone"])

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundZone.ID != p.Attributes["id"] {
			return fmt.Errorf("Zone not found")
		}

		*zone = *foundZone

		return nil
	}
}

func testAccCheckNSRecord(zone *dns.Zone, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[zone.Zone]

		if !ok {
			return fmt.Errorf("Not found: %s", zone.Zone)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NoID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		p := rs.Primary

		if expected != p.Attributes["autogenerate_ns_record"] {
			return fmt.Errorf("Discrepancy between expectation and configuration")
		}

		foundRecord, _, err := client.Records.Get(p.Attributes["zone"], p.Attributes["zone"], "NS")
		if p.Attributes["autogenerate_ns_record"] {
			if err != nil {
				return fmt.Errorf("NS Record not found (autogenerate_ns_record set to true)")
			}

			if foundRecord.Domain != p.Attributes["zone"] {
				return fmt.Errorf("NS Record not found (autogenerate_ns_record set to true)")
			}
		} else if err == nil {
			return fmt.Errorf("NS Record found (autogenerate_ns_record set to false)")
		}

		return nil
	}
}

func testAccCheckZoneDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_zone" {
			continue
		}

		zone, _, err := client.Zones.Get(rs.Primary.Attributes["zone"])

		if err == nil {
			return fmt.Errorf("Zone still exists: %#v: %#v", err, zone)
		}
	}

	return nil
}

func testAccCheckZoneName(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Zone != expected {
			return fmt.Errorf("Zone: got: %s want: %s", zone.Zone, expected)
		}
		return nil
	}
}

func testAccCheckZoneTTL(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.TTL != expected {
			return fmt.Errorf("TTL: got: %d want: %d", zone.TTL, expected)
		}
		return nil
	}
}
func testAccCheckZoneRefresh(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Refresh != expected {
			return fmt.Errorf("Refresh: got: %d want: %d", zone.Refresh, expected)
		}
		return nil
	}
}
func testAccCheckZoneRetry(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Retry != expected {
			return fmt.Errorf("Retry: got: %d want: %d", zone.Retry, expected)
		}
		return nil
	}
}
func testAccCheckZoneExpiry(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Expiry != expected {
			return fmt.Errorf("Expiry: got: %d want: %d", zone.Expiry, expected)
		}
		return nil
	}
}
func testAccCheckZoneNxTTL(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.NxTTL != expected {
			return fmt.Errorf("NxTTL: got: %d want: %d", zone.NxTTL, expected)
		}
		return nil
	}
}

func testAccCheckOtherPorts(zone *dns.Zone, expected []int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(zone.Secondary.OtherPorts) != len(expected) {
			return fmt.Errorf("other_ports: got: %d want %d", len(zone.Secondary.OtherPorts), len(expected))
		}
		for i, v := range zone.Secondary.OtherPorts {
			if v != expected[i] {
				return fmt.Errorf("other_ports[%d]: got: %d want %d", i, v, expected[i])
			}
		}
		return nil
	}
}

func testAccCheckZoneNotPrimary(z *dns.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if z.Primary.Enabled != false {
			return fmt.Errorf("Primary.Enabled: got: true want: false")
		}
		if len(z.Primary.Secondaries) != 0 {
			return fmt.Errorf("Secondaries: got: len(%d) want: len(0)", len(z.Primary.Secondaries))
		}
		return nil
	}
}

func testAccCheckZoneNotSecondary(z *dns.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if z.Secondary != nil {
			// Note that other fields are not cleared. We just toggle "enabled"
			if z.Secondary.Enabled != false {
				return fmt.Errorf("Secondary.Enabled: got: true want: false")
			}
		}
		return nil
	}
}

func testAccCheckZoneDNSSEC(zone *dns.Zone, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *zone.DNSSEC != expected {
			return fmt.Errorf("DNSSEC: got: %t want: %t", *zone.DNSSEC, expected)
		}
		return nil
	}
}

func testAccZoneBasic(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
}
`, zoneName)
}

func testAccZoneUpdated(zoneName string) string {
	return fmt.Sprintf(`
resource "ns1_zone" "it" {
  zone    = "%s"
  ttl     = 10800
  refresh = 3600
  retry   = 300
  expiry  = 2592000
  nx_ttl  = 3601
  dnssec  = false
  # link    = "1.2.3.4.in-addr.arpa" # TODO
  # primary = "1.2.3.4.in-addr.arpa" # TODO
}
`, zoneName)
}

func testAccZonePrimary(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone    = "%s"
  secondaries {
    ip       = "2.2.2.2"
  }
  secondaries {
    ip       = "3.3.3.3"
    notify   = true
    port     = 5353
  }
}
`, zoneName)
}

func testAccZoneSecondary(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone    = "%s"
  ttl     = 10800
  primary = "1.1.1.1"
  additional_primaries = ["2.2.2.2", "3.3.3.3"]
}
`, zoneName)
}

func testAccZoneDNSSEC(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone   = "%s"
  dnssec = true
}
`, zoneName)
}

func testAccZoneDNSSECUpdated(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone   = "%s"
  dnssec = false
}
`, zoneName)
}

func testAccZoneDisableAutoGenerateNSRecord(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone                   = "%s"
  autogenerate_ns_record = false
}
`, zoneName)
}
