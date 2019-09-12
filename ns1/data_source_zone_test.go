package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccDataSourceZone_basic(t *testing.T) {
	var zone dns.Zone
	dataSourceName := "data.ns1_zone.test"
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
				Config: testAccDataSourceZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists(dataSourceName, &zone),
					resource.TestCheckResourceAttr(dataSourceName, "zone", zoneName),
					resource.TestCheckResourceAttr(dataSourceName, "ttl", "3600"),
					resource.TestCheckResourceAttr(dataSourceName, "refresh", "43200"),
					resource.TestCheckResourceAttr(dataSourceName, "retry", "7200"),
					resource.TestCheckResourceAttr(dataSourceName, "expiry", "1209600"),
					resource.TestCheckResourceAttr(dataSourceName, "nx_ttl", "3600"),
					resource.TestCheckResourceAttr(dataSourceName, "primary", "1.1.1.1"),
					resource.TestCheckResourceAttr(dataSourceName, "additional_primaries.0", "2.2.2.2"),
					resource.TestCheckResourceAttr(dataSourceName, "additional_primaries.1", "3.3.3.3"),
				),
			},
		},
	})
}

func TestAccDataSourceZone_primary(t *testing.T) {
	var zone dns.Zone
	dataSourceName := "data.ns1_zone.test"
	zoneName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	// sorted by IP please
	expected := []*dns.ZoneSecondaryServer{
		&dns.ZoneSecondaryServer{NetworkIDs: []int{0}, IP: "2.2.2.2", Port: 53, Notify: false},
		&dns.ZoneSecondaryServer{NetworkIDs: []int{0}, IP: "3.3.3.3", Port: 5353, Notify: true},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceZonePrimary(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists(dataSourceName, &zone),
					resource.TestCheckResourceAttr(dataSourceName, "zone", zoneName),
					testAccCheckZoneSecondary(&zone, 0, expected[0]),
					testAccCheckZoneSecondary(&zone, 1, expected[1]),
				),
			},
		},
	})
}

func testAccDataSourceZoneBasic(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
  primary = "1.1.1.1"
  additional_primaries = ["2.2.2.2", "3.3.3.3"]
}

data "ns1_zone" "test" {
  zone = "${ns1_zone.it.zone}"
}
`, zoneName)
}

func testAccDataSourceZonePrimary(zoneName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
  zone = "%s"
  secondaries {
    ip     = "2.2.2.2"
    notify = false
    port   = 53
  }
  secondaries {
    ip     = "3.3.3.3"
    notify = true
    port   = 5353
  }
}

data "ns1_zone" "test" {
  zone = "${ns1_zone.it.zone}"
}
`, zoneName)
}
