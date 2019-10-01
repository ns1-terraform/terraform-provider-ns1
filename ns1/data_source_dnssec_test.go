package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccDataSourceDNSSEC(t *testing.T) {
	var ds dns.ZoneDNSSEC
	var resourceName = "ns1_zone.it"
	var dataDNSSECName = "data.ns1_dnssec.test"
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
				Config: testAccDataSourceDNSSEC(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDNSSECExists(
						resourceName, dataDNSSECName, zoneName, &ds,
					),
					testAccCheckDataSourceDNSSECKeys(&ds),
					testAccCheckDataSourceDNSSECDelegation(&ds),
				),
			},
		},
	})
}

func testAccCheckDataSourceDNSSECExists(
	resourceName, dataDNSSECName, n string, ds *dns.ZoneDNSSEC,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("test zone not found: %s", resourceName)
		}

		d, ok := s.RootModule().Resources[dataDNSSECName]
		if !ok {
			return fmt.Errorf("data dnssec not found: %s", dataDNSSECName)
		}

		zoneName := d.Primary.Attributes["zone"]
		if zoneName != n {
			return fmt.Errorf(
				"data dnssec zone mismatch: want %s, got %s", n, zoneName,
			)
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.DNSSEC.Get(zoneName)
		if err != nil {
			return fmt.Errorf("API lookup failed: %s", err)
		}

		*ds = *found

		return nil
	}
}

func testAccCheckDataSourceDNSSECKeys(
	d *dns.ZoneDNSSEC,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if d.Keys.TTL != 3600 {
			return fmt.Errorf("Keys.TTL want 3600, got %d", d.Keys.TTL)
		}

		if len(d.Keys.DNSKey) != 2 {
			return fmt.Errorf(
				"Keys.DNSKey length: want 2, got %d", len(d.Keys.DNSKey),
			)
		}
		for i := range d.Keys.DNSKey {
			if err := testAccCheckDNSKey(d.Keys.DNSKey[i]); err != nil {
				return fmt.Errorf("Keys.DNSKey[%d]: %s", i, err)
			}
		}

		return nil
	}
}

func testAccCheckDataSourceDNSSECDelegation(
	d *dns.ZoneDNSSEC,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if d.Delegation.TTL != 3600 {
			return fmt.Errorf("Delegation.TTL want 3600, got %d", d.Delegation.TTL)
		}

		if len(d.Delegation.DNSKey) != 1 {
			return fmt.Errorf(
				"Delgation.DNSKey length: want 1, got %d", len(d.Delegation.DNSKey),
			)
		}
		for i := range d.Delegation.DNSKey {
			if err := testAccCheckDNSKey(d.Delegation.DNSKey[i]); err != nil {
				return fmt.Errorf("Delegation.DNSKey[%d]: %s", i, err)
			}
		}

		if len(d.Delegation.DS) != 1 {
			return fmt.Errorf(
				"Delegation.DS length: want 1, got %d", len(d.Delegation.DS),
			)
		}
		for i := range d.Delegation.DS {
			if err := testAccCheckDNSKey(d.Delegation.DS[i]); err != nil {
				return fmt.Errorf("Delegation.DS[%d]: %s", i, err)
			}
		}

		return nil
	}
}

func testAccCheckDNSKey(key *dns.Key) error {
	if key.Flags == "" {
		return fmt.Errorf("Flags is empty")
	}
	if key.Protocol == "" {
		return fmt.Errorf("Protocol is empty")
	}
	if key.Algorithm == "" {
		return fmt.Errorf("Algorithm is empty")
	}
	if key.PublicKey == "" {
		return fmt.Errorf("PublicKey is empty")
	}
	return nil
}

func testAccDataSourceDNSSEC(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone   = "%s"
  dnssec = true
}

data "ns1_dnssec" "test" {
  zone = "${ns1_zone.it.zone}"
}
`, zoneName)
}
