package ns1

import (
	"fmt"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccApplication_basic(t *testing.T) {
	var application pulsar.Application
	applicationName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationBasic(applicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.it", &application),
				),
			},
		},
	})
}

func testAccCheckApplicationExists(n string, application *pulsar.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundApplication, _, err := client.Applications.Get(rs.Primary.ID)

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundApplication.ID != p.Attributes["id"] {
			return fmt.Errorf("zone not found")
		}

		*application = *foundApplication

		return nil
	}
}

func testAccCheckApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_application" {
			continue
		}

		application, _, err := client.Applications.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("application still exists: %#v: %#v", err, application)
		}
	}

	return nil
}

func testAccApplicationBasic(appName string) string {
	return fmt.Sprintf(`resource "ns1_application" "it" {
 name = "%s"
}
`, appName)

}
