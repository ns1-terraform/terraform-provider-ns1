package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func TestAccDataFeed_basic(t *testing.T) {
	var dataFeed data.Feed
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFeedBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataFeedExists("ns1_datafeed.foobar", "ns1_datasource.api", &dataFeed, t),
					testAccCheckDataFeedName(&dataFeed, "terraform test"),
					testAccCheckDataFeedConfig(&dataFeed, "label", "exampledc2"),
				),
			},
		},
	})
}

func TestAccDataFeed_updated(t *testing.T) {
	var dataFeed data.Feed
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFeedBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataFeedExists("ns1_datafeed.foobar", "ns1_datasource.api", &dataFeed, t),
					testAccCheckDataFeedName(&dataFeed, "terraform test"),
					testAccCheckDataFeedConfig(&dataFeed, "label", "exampledc2"),
				),
			},
			{
				Config: testAccDataFeedUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataFeedExists("ns1_datafeed.foobar", "ns1_datasource.api", &dataFeed, t),
					testAccCheckDataFeedName(&dataFeed, "terraform test"),
					testAccCheckDataFeedConfig(&dataFeed, "label", "exampledc3"),
				),
			},
		},
	})
}

func TestAccDataFeed_ManualDelete(t *testing.T) {
	var dataFeed data.Feed

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFeedBasic,
				Check:  testAccCheckDataFeedExists("ns1_datafeed.foobar", "ns1_datasource.api", &dataFeed, t),
			},
			// Simulate a manual deletion of the data feed and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteDataFeed(&dataFeed),
				Config:             testAccDataFeedBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccDataFeedBasic,
				Check:  testAccCheckDataFeedExists("ns1_datafeed.foobar", "ns1_datasource.api", &dataFeed, t),
			},
		},
	})
}

func testAccCheckDataFeedExists(n string, dsrc string, dataFeed *data.Feed, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		ds, ok := s.RootModule().Resources[dsrc]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("noID is set")
		}

		if ds.Primary.ID == "" {
			return fmt.Errorf("noID is set for the datasource")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundFeed, _, err := client.DataFeeds.Get(ds.Primary.Attributes["id"], rs.Primary.Attributes["id"])

		p := rs.Primary

		if err != nil {
			t.Log(err)
			return err
		}

		if foundFeed.Name != p.Attributes["name"] {
			return fmt.Errorf("DataFeed not found")
		}

		*dataFeed = *foundFeed

		// SourceID doesn't come back in the data feed.
		dataFeed.SourceID = ds.Primary.Attributes["id"]

		return nil
	}
}

func testAccCheckDataFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	var dataFeedID string
	var dataSourceID string

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "ns1_datasource" {
			dataSourceID = rs.Primary.Attributes["id"]
		}

		if rs.Type == "ns1_datafeed" {
			dataFeedID = rs.Primary.Attributes["id"]
		}
	}

	df, _, _ := client.DataFeeds.Get(dataSourceID, dataFeedID)
	if df != nil {
		return fmt.Errorf("dataFeed still exists: %#v", df)
	}

	return nil
}

func testAccCheckDataFeedName(dataFeed *data.Feed, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dataFeed.Name != expected {
			return fmt.Errorf("dataFeed.Name: got: %#v want: %#v", dataFeed.Name, expected)
		}

		return nil
	}
}

func testAccCheckDataFeedConfig(dataFeed *data.Feed, key, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dataFeed.Config[key] != expected {
			return fmt.Errorf("dataFeed.Config[%s]: got: %#v, want: %s", key, dataFeed.Config[key], expected)
		}

		return nil
	}
}

// Simulate a manual deletion of a data feed.
func testAccManualDeleteDataFeed(dataFeed *data.Feed) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.DataFeeds.Delete(dataFeed.SourceID, dataFeed.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete data feed: %v", err)
		}
	}
}

const testAccDataFeedBasic = `
resource "ns1_datasource" "api" {
  name = "terraform test"
  sourcetype = "nsone_v1"
}

resource "ns1_datafeed" "foobar" {
  name = "terraform test"
  source_id = "${ns1_datasource.api.id}"
  config = {
    label = "exampledc2"
  }
}`

const testAccDataFeedUpdated = `
resource "ns1_datasource" "api" {
  name = "terraform test"
  sourcetype = "nsone_v1"
}

resource "ns1_datafeed" "foobar" {
  name = "terraform test"
  source_id = "${ns1_datasource.api.id}"
  config = {
    label = "exampledc3"
  }
}`
