package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func TestAccDataSource_basic(t *testing.T) {
	var dataSource data.Source
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceExists("ns1_datasource.foobar", &dataSource),
					testAccCheckDataSourceName(&dataSource, "terraform test"),
					testAccCheckDataSourceType(&dataSource, "nsone_v1"),
				),
			},
		},
	})
}

func TestAccDataSource_updated(t *testing.T) {
	var dataSource data.Source
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceExists("ns1_datasource.foobar", &dataSource),
					testAccCheckDataSourceName(&dataSource, "terraform test"),
					testAccCheckDataSourceType(&dataSource, "nsone_v1"),
				),
			},
			{
				Config: testAccDataSourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceExists("ns1_datasource.foobar", &dataSource),
					testAccCheckDataSourceName(&dataSource, "terraform test"),
					testAccCheckDataSourceType(&dataSource, "nsone_monitoring"),
				),
			},
		},
	})
}

func TestAccDataSource_ManualDelete(t *testing.T) {
	var dataSource data.Source
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBasic,
				Check:  testAccCheckDataSourceExists("ns1_datasource.foobar", &dataSource),
			},
			// Simulate a manual deletion of the data source and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteDataSource(&dataSource),
				Config:             testAccDataSourceBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccDataSourceBasic,
				Check:  testAccCheckDataSourceExists("ns1_datasource.foobar", &dataSource),
			},
		},
	})
}

func testAccCheckDataSourceExists(n string, dataSource *data.Source) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("noID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundSource, _, err := client.DataSources.Get(rs.Primary.Attributes["id"])

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundSource.Name != p.Attributes["name"] {
			return fmt.Errorf("datasource not found")
		}

		*dataSource = *foundSource

		return nil
	}
}

func testAccCheckDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_datasource" {
			continue
		}

		_, _, err := client.DataSources.Get(rs.Primary.Attributes["id"])

		if err == nil {
			return fmt.Errorf("datasource still exists")
		}
	}

	return nil
}

func testAccCheckDataSourceName(dataSource *data.Source, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dataSource.Name != expected {
			return fmt.Errorf("dataSource.Name: got: %#v want: %#v", dataSource.Name, expected)
		}

		return nil
	}
}

func testAccCheckDataSourceType(dataSource *data.Source, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dataSource.Type != expected {
			return fmt.Errorf("dataSource.Type: got: %#v want: %#v", dataSource.Type, expected)
		}

		return nil
	}
}

// Simulate a manual deletion of a data source.
func testAccManualDeleteDataSource(dataSource *data.Source) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.DataSources.Delete(dataSource.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete data source: %v", err)
		}
	}
}

const testAccDataSourceBasic = `
resource "ns1_datasource" "foobar" {
	name = "terraform test"
	sourcetype = "nsone_v1"
}`

const testAccDataSourceUpdated = `
resource "ns1_datasource" "foobar" {
	name = "terraform test"
	sourcetype = "nsone_monitoring"
}`
