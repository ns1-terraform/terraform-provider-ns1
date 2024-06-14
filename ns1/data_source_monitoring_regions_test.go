package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMonitoringRegions_basic(t *testing.T) {
	name := "foobar"
	resourceName := fmt.Sprintf("data.ns1_monitoring_regions.%s", name)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMonitoringRegionsBasic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitoringRegionsLength(resourceName),
				),
			},
		},
	})
}

func testAccCheckMonitoringRegionsLength(
	n string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		regions, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// make sure we get some monitoring regions
		if len(regions.Primary.Attributes["regions.#"]) == 0 {
			return fmt.Errorf("no monitoring regions found")
		}

		return nil
	}
}

const testAccMonitoringRegionsBasic = `
data "ns1_monitoring_regions" "%s" {
}
`
