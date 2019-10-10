package ns1

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// Note: sorts the Set of secondaries by IP
func testAccCheckZoneSecondaries(
	t *testing.T, z *dns.Zone, idx int, expected *dns.ZoneSecondaryServer,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.True(t, z.Primary.Enabled)

		secondaries := z.Primary.Secondaries
		sort.SliceStable(secondaries, func(i, j int) bool {
			return secondaries[i].IP < secondaries[j].IP
		})
		secondary := secondaries[idx]

		assert.Equal(t, expected.IP, secondary.IP)
		assert.Equal(t, expected.Port, secondary.Port)
		assert.Equal(t, expected.Notify, secondary.Notify)
		assert.ElementsMatch(t, expected.NetworkIDs, secondary.NetworkIDs)

		return nil
	}
}
