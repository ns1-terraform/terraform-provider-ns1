package ns1

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccDataSourceZone_basic(t *testing.T) {
	var zone dns.Zone
	dataSourceName := "data.ns1_zone.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceZoneBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists(dataSourceName, &zone),
					resource.TestCheckResourceAttr(dataSourceName, "zone", "terraform-testda-zone.io"),
					resource.TestCheckResourceAttr(dataSourceName, "ttl", "3600"),
					resource.TestCheckResourceAttr(dataSourceName, "refresh", "43200"),
					resource.TestCheckResourceAttr(dataSourceName, "retry", "7200"),
					resource.TestCheckResourceAttr(dataSourceName, "expiry", "1209600"),
					resource.TestCheckResourceAttr(dataSourceName, "nx_ttl", "3600"),
				),
			},
		},
	})
}

const testAccDataSourceZoneBasic = `
resource "ns1_zone" "it" {
  zone = "terraform-testda-zone.io"
}

data "ns1_zone" "test" {
  zone = "${ns1_zone.it.zone}"
}
`
