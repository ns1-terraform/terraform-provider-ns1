package ns1

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

// Creating basic DNS alert
func TestAccAlert_basic(t *testing.T) {
	var (
		alert     = alerting.Alert{}
		alertName = fmt.Sprintf("terraform-test-alert-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertBasic(alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &alert),
					testAccCheckAlertName(&alert, alertName),
					// testAccCheckAlertPreference(&alert, alertPreference),
				),
			},
			{
				ResourceName:      "ns1_alert.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlert_links(t *testing.T) {
	var (
		alert     = alerting.Alert{}
		alertName = fmt.Sprintf("terraform-test-alert-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)

	zoneNames := []string{}
	for i := 0; i < 3; i++ {
		zoneNames = append(zoneNames, fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)))
	}

	// Support objects
	notfierLists := []*monitor.NotifyList{{}, {}, {}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertZones(alertName, zoneNames),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &alert),
					testAccCheckAlertName(&alert, alertName),
					testAccCheckAlertType(&alert, "zone"),
					testAccCheckAlertSubtype(&alert, "transfer_failed"),
					testAccCheckAlertZoneNames(&alert, zoneNames),
					testAccCheckAlertNotifierLists(&alert, []*monitor.NotifyList{}),
				),
			},
			{
				Config: testAccAlertNotifierLists(alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &alert),
					testAccCheckNotifyListExists("ns1_notifylist.email_list_0", notfierLists[0]),
					testAccCheckNotifyListExists("ns1_notifylist.email_list_1", notfierLists[1]),
					testAccCheckNotifyListExists("ns1_notifylist.email_list_2", notfierLists[2]),
					testAccCheckAlertName(&alert, alertName),
					testAccCheckAlertType(&alert, "zone"),
					testAccCheckAlertSubtype(&alert, "transfer_failed"),
					testAccCheckAlertZoneNames(&alert, []string{}),
					testAccCheckAlertNotifierLists(&alert, notfierLists),
				),
			},
			{
				ResourceName:      "ns1_alert.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Update DNS alert
func TestAccAlert_update(t *testing.T) {
	var (
		alert     = alerting.Alert{}
		alertName = fmt.Sprintf("terraform-test-alert-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		updatedAlert     = alerting.Alert{}
		updatedAlertName = fmt.Sprintf("terraform-test-alert-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		zoneName         = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		// Support objects
		nl = monitor.NotifyList{}
		// zone = dns.Zone{}
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertBasic(alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &alert),
					testAccCheckAlertName(&alert, alertName),
					testAccCheckAlertType(&alert, "zone"),
					testAccCheckAlertSubtype(&alert, "transfer_failed"),
					testAccCheckAlertZoneNames(&alert, []string{}),
					testAccCheckAlertNotifierLists(&alert, []*monitor.NotifyList{}),
				),
			},
			{
				Config: testAccAlertUpdated(zoneName, updatedAlertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &updatedAlert),
					// Have to retrieve the notifier list to get random ID.
					testAccCheckNotifyListExists("ns1_notifylist.it", &nl),
					testAccCheckAlertName(&updatedAlert, updatedAlertName),
					testAccCheckAlertType(&updatedAlert, "zone"),
					testAccCheckAlertSubtype(&updatedAlert, "transfer_failed"),
					testAccCheckAlertZoneNames(&updatedAlert, []string{zoneName}),
					testAccCheckAlertNotifierLists(&updatedAlert, []*monitor.NotifyList{&nl}),
				),
			},
			{
				ResourceName:      "ns1_alert.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Manually deleting DNS Alert
func TestAccAlert_ManualDelete(t *testing.T) {
	var (
		alert     = alerting.Alert{}
		alertName = fmt.Sprintf("terraform-test-alert-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)
	// Manual deletion test for DNS Alert
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertBasic(alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertExists("ns1_alert.it", &alert),
				),
			},
			// Simulate a manual deletion of the DNS Alert and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteAlert(&alert),
				Config:             testAccAlertBasic(alertName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccAlertBasic(alertName),
				Check:  testAccCheckAlertExists("ns1_alert.it", &alert),
			},
		},
	})
}

func testAccAlertBasic(alertName string) string {
	return fmt.Sprintf(`resource "ns1_alert" "it" {
	name               = "%s"
	type               = "zone"
	subtype            = "transfer_failed"
	notification_lists = []
	zone_names = []
}`, alertName)
}

// Updates alert "it" created above.
func testAccAlertUpdated(zoneName, alertName string) string {
	return fmt.Sprintf(`resource "ns1_zone" "it" {
		zone = "%s"
	}

	resource "ns1_notifylist" "it" {
		name = "email list"
		notifications {
			type = "email"
			config = {
			email = "jdoe@example.com"
			}
		}
	}

	resource "ns1_alert" "it" {
  		name = "%s"
		type = "zone"
		subtype = "transfer_failed"
		zone_names = ["${ns1_zone.it.zone}"]
		notification_lists = ["${ns1_notifylist.it.id}"]
}
`, zoneName, alertName)
}

func testAccAlertNotifierLists(alertName string) string {
	config := ""
	listResourceNames := []string{}
	for i := range []int{1, 2, 3} {
		listName := fmt.Sprintf("terraform-test-list-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		listResourceNames = append(listResourceNames, fmt.Sprintf("ns1_notifylist.email_list_%d.id", i))
		config += fmt.Sprintf(`
resource "ns1_notifylist" "email_list_%d" {
  name = "%s"
  notifications {
    type = "email"
    config = {
      email = "jdoe@example.com"
    }
  }
}`,
			i, listName)
	}

	config += fmt.Sprintf(`
resource "ns1_alert" "it" {
  name               = "%s"
  type               = "zone"
  subtype            = "transfer_failed"
  notification_lists = [%s]
  zone_names = []
}`,
		alertName, strings.Join(listResourceNames, ","))
	return config
}

func testAccAlertZones(alertName string, zoneNames []string) string {
	config := ""

	zoneResourceNames := make([]string, 0, len(zoneNames))
	for i := range zoneNames {
		zoneResourceNames = append(zoneResourceNames, fmt.Sprintf("ns1_zone.alert_zone_%d.zone", i))
		config += fmt.Sprintf(`
resource "ns1_zone" "alert_zone_%d" {
  zone                 = "%s"
  primary              = "192.0.2.1"
  additional_primaries = ["192.0.2.2"]
  additional_ports = [53]
}`,
			i, zoneNames[i])
	}

	config += fmt.Sprintf(`
resource "ns1_alert" "it" {
  name               = "%s"
  type               = "zone"
  subtype            = "transfer_failed"
  notification_lists = []
  zone_names = [%s]
}`,
		alertName, strings.Join(zoneResourceNames, ","))
	return config
}

func testAccCheckAlertName(alert *alerting.Alert, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *alert.Name != expected {
			return fmt.Errorf("alert.Name: got: %s want: %s", *alert.Name, expected)
		}
		return nil
	}
}

func testAccCheckAlertType(alert *alerting.Alert, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *alert.Type != expected {
			return fmt.Errorf("alert.Type: got: %s want: %s", *alert.Type, expected)
		}
		return nil
	}
}

func testAccCheckAlertSubtype(alert *alerting.Alert, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *alert.Subtype != expected {
			return fmt.Errorf("alert.Subtype: got: %s want: %s", *alert.Subtype, expected)
		}
		return nil
	}
}

func testAccCheckAlertZoneNames(alert *alerting.Alert, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actualSorted := alert.ZoneNames
		expectedSorted := expected
		sort.Strings(actualSorted)
		sort.Strings(expectedSorted)
		if !reflect.DeepEqual(actualSorted, expectedSorted) {
			return fmt.Errorf("alert.Zones: got: %v want: %v", actualSorted, expectedSorted)
		}
		return nil
	}
}

func testAccCheckAlertNotifierLists(alert *alerting.Alert, expectedLists []*monitor.NotifyList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actualSorted := alert.NotifierListIds
		sort.Strings(actualSorted)
		expected := []string{}
		for _, nl := range expectedLists {
			expected = append(expected, nl.ID)
		}
		sort.Strings(expected)
		if !reflect.DeepEqual(actualSorted, expected) {
			return fmt.Errorf("alert.NotifierListIds: got: %v want: `%v`", actualSorted, expected)
		}
		return nil
	}
}

func testAccCheckAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_alert" {
			continue
		}

		alert, _, err := client.Alerts.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("DNS Alert still exists: %#v: %#v", err, alert)
		}

	}

	return nil
}

func testAccCheckAlertExists(n string, alert *alerting.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundAlert, _, err := client.Alerts.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		*alert = *foundAlert

		return nil
	}
}

func testAccManualDeleteAlert(alert *alerting.Alert) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Alerts.Delete(*alert.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete DNS alert: %v", err)
		}
	}
}
