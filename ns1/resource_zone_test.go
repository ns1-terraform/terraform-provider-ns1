package ns1

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccZone_basic(t *testing.T) {
	var zone dns.Zone
	defaultHostmaster := "hostmaster@nsone.net"
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
					testAccCheckNSRecord("ns1_zone.it", true),
					testAccCheckZoneHostmaster(&zone, defaultHostmaster),
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
		{
			NetworkIDs: []int{},
			IP:         "2.2.2.2",
			Port:       53,
			Notify:     false,
		},
		{
			NetworkIDs: []int{},
			IP:         "3.3.3.3",
			Port:       5353,
			Notify:     true,
		},
	}
	expectedOtherPorts := []int{53, 5353}
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
					resource.TestCheckResourceAttr("ns1_zone.it", "primary_port", "54"),
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
		{
			NetworkIDs: []int{},
			IP:         "2.2.2.2",
			Port:       53,
			Notify:     false,
		},
		{
			NetworkIDs: []int{},
			IP:         "3.3.3.3",
			Port:       5353,
			Notify:     true,
		},
	}
	expectedOtherPorts := []int{53, 5353}
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
					resource.TestCheckResourceAttr("ns1_zone.it", "primary_port", "54"),
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

func TestAccZone_hostmaster(t *testing.T) {
	var zone dns.Zone
	defaultHostmaster := "hostmaster@nsone.net"
	zoneHostmaster := "hostmaster@rname.test"
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
					testAccCheckZoneHostmaster(&zone, defaultHostmaster),
				),
			},
			{
				Config:             testAccZoneHostmaster(zoneName, zoneHostmaster),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccZoneHostmaster(zoneName, zoneHostmaster),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneHostmaster(&zone, zoneHostmaster),
				),
			},
			//detect loop conditions
			{
				Config:             testAccZoneHostmaster(zoneName, zoneHostmaster),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccZone_TSIG(t *testing.T) {
	var zone dns.Zone
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	tsig := &dns.TSIG{
		Enabled: true,
		Name:    fmt.Sprintf("terraform-test-%s.", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)),
		Hash:    "hmac-sha256",
		Key:     "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA==",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneSecondaryTSIG(zoneName, tsig.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneTsigEnabled(&zone, tsig.Enabled),
					testAccCheckZoneTsigName(&zone, tsig.Name),
					testAccCheckZoneTsigHash(&zone, tsig.Hash),
					testAccCheckZoneTsigKey(&zone, tsig.Key),
					resource.TestCheckResourceAttr("ns1_zone.it", "primary_port", "53"),
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
					testAccCheckNSRecord("ns1_zone.it", false),
				),
			},
			{
				Config: testAccZoneDisableAutoGenerateNSRecordLinkedZone(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("ns1_zone.it", &zone),
					testAccCheckZoneName(&zone, zoneName),
					testAccCheckZoneExists("ns1_zone.linked_zone", &zone),
					testAccCheckZoneName(&zone, "linkedzone_" + zoneName),
				),
			},
			{
				ResourceName:      "ns1_zone.linked_zone",
				ImportState:       true,
				ImportStateId:     zoneName,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccZone_ManualDelete(t *testing.T) {
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
				Check:  testAccCheckZoneExists("ns1_zone.it", &zone),
			},
			// Simulate a manual deletion of the zone and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteZone(zoneName),
				Config:             testAccZoneBasic(zoneName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccZoneBasic(zoneName),
				Check:  testAccCheckZoneExists("ns1_zone.it", &zone),
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundZone, _, err := client.Zones.Get(rs.Primary.Attributes["zone"])

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundZone.ID != p.Attributes["id"] {
			return fmt.Errorf("zone not found")
		}

		*zone = *foundZone

		return nil
	}
}

func testAccCheckNSRecord(n string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource.Primary: no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		p := rs.Primary

		shouldAutogenerate, err := strconv.ParseBool(
			p.Attributes["autogenerate_ns_record"],
		)
		if err != nil {
			return err
		}

		if expected != shouldAutogenerate {
			return fmt.Errorf(
				"autogenerate_ns_record: want %t, got %t",
				expected,
				shouldAutogenerate,
			)
		}

		foundRecord, _, err := client.Records.Get(
			p.Attributes["zone"], p.Attributes["zone"], "NS",
		)
		if shouldAutogenerate {
			if err != nil {
				return fmt.Errorf(
					"NS Record not found (autogenerate_ns_record set to true)",
				)
			}

			if foundRecord.Domain != p.Attributes["zone"] {
				return fmt.Errorf("an NS Record found, but domain does not match")
			}
		} else if err == nil {
			return fmt.Errorf("an NS Record found (autogenerate_ns_record set to false)")
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
			return fmt.Errorf("zone still exists: %#v: %#v", err, zone)
		}
	}

	return nil
}

func testAccCheckZoneTsigEnabled(zone *dns.Zone, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Secondary.TSIG.Enabled != expected {
			return fmt.Errorf("zone.Secondary.TSIG.Enabled: got: %v %T want: %v %T", zone.Secondary.TSIG.Enabled, zone.Secondary.TSIG.Enabled, expected, expected)
		}
		return nil
	}
}

func testAccCheckZoneTsigName(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Secondary.TSIG.Name != expected {
			return fmt.Errorf("zone.Secondary.TSIG.Name: got: %s want: %s", zone.Secondary.TSIG.Name, expected)
		}
		return nil
	}
}

func testAccCheckZoneTsigHash(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Secondary.TSIG.Hash != expected {
			return fmt.Errorf("zone.Secondary.TSIG.Hash: got: %s want: %s", zone.Secondary.TSIG.Hash, expected)
		}
		return nil
	}
}

func testAccCheckZoneTsigKey(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Secondary.TSIG.Key != expected {
			return fmt.Errorf("zone.Secondary.TSIG.Key: got: %s want: %s .", zone.Secondary.TSIG.Key, expected)
		}
		return nil
	}
}

func testAccCheckZoneName(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Zone != expected {
			return fmt.Errorf("zone: got: %s want: %s", zone.Zone, expected)
		}
		return nil
	}
}

func testAccCheckZoneTTL(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.TTL != expected {
			return fmt.Errorf("zone.TTL: got: %d want: %d", zone.TTL, expected)
		}
		return nil
	}
}
func testAccCheckZoneRefresh(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Refresh != expected {
			return fmt.Errorf("zone.Refresh: got: %d want: %d", zone.Refresh, expected)
		}
		return nil
	}
}
func testAccCheckZoneRetry(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Retry != expected {
			return fmt.Errorf("zone.Retry: got: %d want: %d", zone.Retry, expected)
		}
		return nil
	}
}
func testAccCheckZoneExpiry(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Expiry != expected {
			return fmt.Errorf("zone.Expiry: got: %d want: %d", zone.Expiry, expected)
		}
		return nil
	}
}
func testAccCheckZoneNxTTL(zone *dns.Zone, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.NxTTL != expected {
			return fmt.Errorf("zone.NxTTL: got: %d want: %d", zone.NxTTL, expected)
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
			return fmt.Errorf("z.Primary.Enabled: got: true want: false")
		}
		if len(z.Primary.Secondaries) != 0 {
			return fmt.Errorf("secondaries: got: len(%d) want: len(0)", len(z.Primary.Secondaries))
		}
		return nil
	}
}

func testAccCheckZoneNotSecondary(z *dns.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if z.Secondary != nil {
			// Note that other fields are not cleared. We just toggle "enabled"
			if z.Secondary.Enabled != false {
				return fmt.Errorf("z.Secondary.Enabled: got: true want: false")
			}
		}
		return nil
	}
}

func testAccCheckZoneDNSSEC(zone *dns.Zone, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.DNSSEC == nil {
			return fmt.Errorf("DNSSEC field not defined.")
		}
		if *zone.DNSSEC != expected {
			return fmt.Errorf("DNSSEC: got: %t want: %t", *zone.DNSSEC, expected)
		}
		return nil
	}
}

func testAccCheckZoneHostmaster(zone *dns.Zone, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if zone.Hostmaster != expected {
			return fmt.Errorf("Hostmaster: got: %s want: %s", zone.Hostmaster, expected)
		}
		return nil
	}
}

// Simulate a manual deletion of a zone.
func testAccManualDeleteZone(zone string) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Zones.Delete(zone)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete zone: %v", err)
		}
	}
}

func testAccZoneBasic(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
}
`, zoneName)
}

func testAccZoneHostmaster(zoneName string, hostmaster string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
  hostmaster = "%s"
}
`, zoneName, hostmaster)
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
  primary_port = 54
  additional_primaries = ["2.2.2.2", "3.3.3.3"]
  additional_ports = [53, 5353]
}
`, zoneName)
}

func testAccZoneSecondaryTSIG(zoneName, tsigName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone    = "%s"
  primary = "1.1.1.1"
  # primary_port left unspecified to test default/computed case
  tsig = {
    enabled = true
    name = "%s"
    hash = "hmac-sha256"
    key = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA=="
  }
}
`, zoneName, tsigName)
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

func testAccZoneDisableAutoGenerateNSRecordLinkedZone(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
}

resource "ns1_zone" "linked_zone" {
  zone = "linkedzone_%s"
  link = ns1_zone.it.zone
  autogenerate_ns_record = false
}
`, zoneName, zoneName)
}
