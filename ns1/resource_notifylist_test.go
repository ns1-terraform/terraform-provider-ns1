package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

func TestAccNotifyList_basic(t *testing.T) {
	var nl monitor.NotifyList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotifyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyListBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test", &nl),
					testAccCheckNotifyListName(&nl, "terraform test"),
				),
			},
			{
				ResourceName:      "ns1_notifylist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotifyList_updated(t *testing.T) {
	var nl monitor.NotifyList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotifyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyListBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test", &nl),
					testAccCheckNotifyListName(&nl, "terraform test"),
				),
			},
			{
				Config: testAccNotifyListUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test", &nl),
					testAccCheckNotifyListName(&nl, "terraform test"),
				),
			},
		},
	})
}

func TestAccNotifyList_multiple(t *testing.T) {
	var nl monitor.NotifyList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotifyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyListMultiple,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test_multiple", &nl),
					testAccCheckNotifyListName(&nl, "terraform test multiple"),
				),
			},
// This fails because the schema.TypeList is ordered. We want to switch
// this to schema.TypeSet but cannot due to SDK issue #652 / #895, fix for
// which is waiting for review, see
// https://github.com/hashicorp/terraform-plugin-sdk/pull/1042
//			{
//				Config: testAccNotifyListMultipleDifferentOrder ,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNotifyListExists("ns1_notifylist.test_multiple2", &nl),
//					testAccCheckNotifyListName(&nl, "terraform test multiple2"),
//				),
//			},
			{
				ResourceName:      "ns1_notifylist.test_multiple",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotifyList_types(t *testing.T) {
	var nl monitor.NotifyList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotifyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyListSlack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test_slack", &nl),
					testAccCheckNotifyListName(&nl, "terraform test slack"),
				),
			},
			{
				ResourceName:      "ns1_notifylist.test_slack",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNotifyListPagerDuty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotifyListExists("ns1_notifylist.test_pagerduty", &nl),
					testAccCheckNotifyListName(&nl, "terraform test pagerduty"),
				),
			},
			{
				ResourceName:      "ns1_notifylist.test_pagerduty",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotifyList_ManualDelete(t *testing.T) {
	var nl monitor.NotifyList

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotifyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotifyListBasic,
				Check:  testAccCheckNotifyListExists("ns1_notifylist.test", &nl),
			},
			// Simulate a manual deletion of the notify list and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteNotifyList(&nl),
				Config:             testAccNotifyListBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccNotifyListBasic,
				Check:  testAccCheckNotifyListExists("ns1_notifylist.test", &nl),
			},
		},
	})
}

func testAccCheckNotifyListState(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["ns1_notifylist.test"]
		if !ok {
			return fmt.Errorf("not found: %s", "ns1_notifylist.test")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		p := rs.Primary
		if p.Attributes[key] != value {
			return fmt.Errorf(
				"%s != %s (actual: %s)", key, value, p.Attributes[key])
		}

		return nil
	}
}

func testAccCheckNotifyListExists(n string, nl *monitor.NotifyList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %v", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("ID is not set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundNl, _, err := client.Notifications.Get(id)

		if err != nil {
			return err
		}

		if foundNl.ID != id {
			return fmt.Errorf("notify List not found want: %#v, got %#v", id, foundNl)
		}

		*nl = *foundNl

		return nil
	}
}

func testAccCheckNotifyListDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_notifylist" {
			continue
		}

		nl, _, err := client.Notifications.Get(rs.Primary.Attributes["id"])

		if err == nil {
			return fmt.Errorf("notify List still exists %#v: %#v", err, nl)
		}
	}

	return nil
}

func testAccCheckNotifyListName(nl *monitor.NotifyList, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nl.Name != expected {
			return fmt.Errorf("nl.Name: got: %#v want: %#v", nl.Name, expected)
		}
		return nil
	}
}

// Simulate a manual deletion of a notify list.
func testAccManualDeleteNotifyList(nl *monitor.NotifyList) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Notifications.Delete(nl.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete notify list: %v", err)
		}
	}
}

const testAccNotifyListBasic = `
resource "ns1_notifylist" "test" {
  name = "terraform test"
  notifications {
    type = "webhook"
    config = {
      url = "http://localhost:9090"
    }
  }
}
`

const testAccNotifyListUpdated = `
resource "ns1_notifylist" "test" {
  name = "terraform test"
  notifications {
    type = "webhook"
    config = {
      url = "http://localhost:9091"
    }
  }
}
`
const testAccNotifyListSlack = `
resource "ns1_notifylist" "test_slack" {
  name = "terraform test slack"
  notifications {
    type = "slack"
    config = {
      username = "tf-test"
      url = "http://localhost:9091"
      channel = "TF Test Channel"
    }
  }
}
`
const testAccNotifyListPagerDuty = `
resource "ns1_notifylist" "test_pagerduty" {
  name = "terraform test pagerduty"
  notifications {
    type = "pagerduty"
    config = {
      service_key = "tftestkey"
    }
  }
}
`

const testAccNotifyListMultiple = `
resource "ns1_notifylist" "test_multiple" {
  name = "terraform test multiple"
  notifications {
    type = "pagerduty"
    config = {
      service_key = "tftestkey"
    }
  }
  notifications {
    type = "webhook"
    config = {
      url = "http://localhost:9090"
    }
  }
}
`

const testAccNotifyListMultipleDifferentOrder = `
resource "ns1_notifylist" "test_multiple2" {
  name = "terraform test multiple2"
  notifications {
    type = "webhook"
    config = {
      url = "http://localhost:9090"
    }
  }
  notifications {
    type = "pagerduty"
    config = {
      service_key = "tftestkey"
    }
  }
}
`
