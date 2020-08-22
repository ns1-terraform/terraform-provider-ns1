package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccDataSourceRecord(t *testing.T) {
	var record dns.Record
	dataSourceName := "data.ns1_record.test"
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)
	rrType := "CNAME"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResource(zoneName, domainName, rrType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					resource.TestCheckResourceAttr(dataSourceName, "zone", zoneName),
					resource.TestCheckResourceAttr(dataSourceName, "domain", domainName),
					resource.TestCheckResourceAttr(dataSourceName, "type", rrType),
				),
			},
		},
	})
}

func testAccDataSourceResource(zoneName string, domainName string, rrType string) string {
	return fmt.Sprintf(`resource "ns1_record" "it" {
  zone = "%s"
  domain = "%s"
  type = "%s"
}

data "ns1_record" "test" {
  zone = "${ns1_record.it.zone}"
  domain = "${ns1_record.it.domain}"
  type = "${ns1_record.it.type}"
}
`, zoneName, domainName, rrType)
}
