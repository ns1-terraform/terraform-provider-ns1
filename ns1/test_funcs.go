package ns1

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// Note: sorts the Set of secondaries by IP
func testAccCheckZoneSecondaries(
	t *testing.T, z *dns.Zone, idx int, expected *dns.ZoneSecondaryServer,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Truef(t, z.Primary.Enabled, "primary.enabled: got %v, wanted %v", z.Primary.Enabled, true)
		secondaries := z.Primary.Secondaries
		sort.SliceStable(secondaries, func(i, j int) bool {
			return secondaries[i].IP < secondaries[j].IP
		})
		secondary := secondaries[idx]

		assert.Equalf(t, expected.IP, secondary.IP, "secondary IP: got %v, wanted %v", secondary.IP, expected.IP)
		assert.Equalf(t, expected.Port, secondary.Port, "secondary port: got %v, wanted %v", secondary.Port, expected.Port)
		assert.Equalf(t, expected.Notify, secondary.Notify, "secondary notify: got %v, wanted %v", secondary.Notify, expected.Notify)
		assert.ElementsMatchf(t, expected.NetworkIDs, secondary.NetworkIDs, "secondary network ID mismatch: got %v, wanted %v", z.Primary.Secondaries[idx], expected)

		return nil
	}
}
