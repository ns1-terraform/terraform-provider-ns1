package ns1

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

func TestAccNetworks_basic(t *testing.T) {
	name := "foobar"
	resourceName := fmt.Sprintf("data.ns1_networks.%s", name)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccNetworksBasic, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNumberOfNetworks(resourceName),
					testAccCheckNetworks(resourceName),
				),
			},
		},
	})
}
func testAccCheckNumberOfNetworks(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ns1.Client)
		foundNetworks, _, err := client.Network.Get()
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(n, "networks.#", strconv.Itoa(len(foundNetworks)))(s)
	}
}

func testAccCheckNetworks(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ns1.Client)
		foundNetworks, _, err := client.Network.Get()
		if err != nil {
			return err
		}

		for _, network := range foundNetworks {
			err = resource.TestCheckTypeSetElemNestedAttrs(
				n,
				"networks.*",
				map[string]string{
					"name":       network.Name,
					"label":      network.Label,
					"network_id": strconv.Itoa(network.NetworkID),
				},
			)(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

const testAccNetworksBasic = `
data "ns1_networks" "%s" {
}
`
