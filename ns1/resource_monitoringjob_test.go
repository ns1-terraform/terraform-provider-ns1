package ns1

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

func TestAccMonitoringJob_basic(t *testing.T) {
	var mj monitor.Job
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringJobBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitoringJobExists("ns1_monitoringjob.it", &mj),
					testAccCheckMonitoringJobName(&mj, "terraform test"),
					testAccCheckMonitoringJobActive(&mj, true),
					testAccCheckMonitoringJobRegions(&mj, []string{"lga"}),
					testAccCheckMonitoringJobType(&mj, "tcp"),
					testAccCheckMonitoringJobFrequency(&mj, 60),
					testAccCheckMonitoringJobRapidRecheck(&mj, false),
					testAccCheckMonitoringJobPolicy(&mj, "quorum"),
					testAccCheckMonitoringJobConfigSend(&mj, "HEAD / HTTP/1.0\r\n\r\n"),
					testAccCheckMonitoringJobConfigPort(&mj, 443),
					testAccCheckMonitoringJobConfigHost(&mj, "1.2.3.4"),
					testAccCheckMonitoringJobRuleValue(&mj, "200 OK"),
					testAccCheckMonitoringJobRuleComparison(&mj, "contains"),
					testAccCheckMonitoringJobRuleKey(&mj, "output"),
					testAccCheckMonitoringJobMute(&mj, true),
				),
			},
		},
	})
}

func TestAccMonitoringJob_updated(t *testing.T) {
	var mj monitor.Job
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringJobBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitoringJobExists("ns1_monitoringjob.it", &mj),
					testAccCheckMonitoringJobName(&mj, "terraform test"),
					testAccCheckMonitoringJobActive(&mj, true),
					testAccCheckMonitoringJobRegions(&mj, []string{"lga"}),
					testAccCheckMonitoringJobType(&mj, "tcp"),
					testAccCheckMonitoringJobFrequency(&mj, 60),
					testAccCheckMonitoringJobRapidRecheck(&mj, false),
					testAccCheckMonitoringJobPolicy(&mj, "quorum"),
					testAccCheckMonitoringJobConfigSend(&mj, "HEAD / HTTP/1.0\r\n\r\n"),
					testAccCheckMonitoringJobConfigPort(&mj, 443),
					testAccCheckMonitoringJobConfigHost(&mj, "1.2.3.4"),
					testAccCheckMonitoringJobRuleValue(&mj, "200 OK"),
					testAccCheckMonitoringJobRuleComparison(&mj, "contains"),
					testAccCheckMonitoringJobRuleKey(&mj, "output"),
					testAccCheckMonitoringJobMute(&mj, true),
				),
			},
			{
				Config: testAccMonitoringJobUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitoringJobExists("ns1_monitoringjob.it", &mj),
					testAccCheckMonitoringJobName(&mj, "terraform test"),
					testAccCheckMonitoringJobActive(&mj, true),
					testAccCheckMonitoringJobRegions(&mj, []string{"lga"}),
					testAccCheckMonitoringJobType(&mj, "tcp"),
					testAccCheckMonitoringJobFrequency(&mj, 120),
					testAccCheckMonitoringJobRapidRecheck(&mj, true),
					testAccCheckMonitoringJobPolicy(&mj, "all"),
					testAccCheckMonitoringJobConfigSend(&mj, "HEAD / HTTP/1.0\r\n\r\n"),
					testAccCheckMonitoringJobConfigPort(&mj, 443),
					testAccCheckMonitoringJobConfigHost(&mj, "1.1.1.1"),
					testAccCheckMonitoringJobRuleValue(&mj, "200"),
					testAccCheckMonitoringJobRuleComparison(&mj, "<="),
					testAccCheckMonitoringJobRuleKey(&mj, "connect"),
					testAccCheckMonitoringJobMute(&mj, false),
				),
			},
		},
	})
}

func TestAccMonitoringJob_ManualDelete(t *testing.T) {
	var mj monitor.Job

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringJobBasic,
				Check:  testAccCheckMonitoringJobExists("ns1_monitoringjob.it", &mj),
			},
			// Simulate a manual deletion of the monitoring job and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteMonitoringJob(&mj),
				Config:             testAccMonitoringJobBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccMonitoringJobBasic,
				Check:  testAccCheckMonitoringJobExists("ns1_monitoringjob.it", &mj),
			},
		},
	})
}

func testAccCheckMonitoringJobState(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["ns1_monitoringjob.it"]
		if !ok {
			return fmt.Errorf("not found: %s", "ns1_monitoringjob.it")
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

func testAccCheckMonitoringJobExists(n string, monitoringJob *monitor.Job) resource.TestCheckFunc {
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

		foundMj, _, err := client.Jobs.Get(id)

		if err != nil {
			return err
		}

		if foundMj.ID != id {
			return fmt.Errorf("monitoring Job not found want: %#v, got %#v", id, foundMj)
		}

		*monitoringJob = *foundMj

		return nil
	}
}

func testAccCheckMonitoringJobDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_monitoringjob" {
			continue
		}

		mj, _, err := client.Jobs.Get(rs.Primary.Attributes["id"])

		if err == nil {
			return fmt.Errorf("monitoring Job still exists %#v: %#v", err, mj)
		}
	}

	return nil
}

func testAccCheckMonitoringJobName(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Name != expected {
			return fmt.Errorf("mj.Name: got: %#v want: %#v", mj.Name, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobActive(mj *monitor.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Active != expected {
			return fmt.Errorf("mj.Active: got: %#v want: %#v", mj.Active, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobRegions(mj *monitor.Job, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !reflect.DeepEqual(mj.Regions, expected) {
			return fmt.Errorf("mj.Regions: got: %#v want: %#v", mj.Regions, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobType(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Type != expected {
			return fmt.Errorf("mj.Type: got: %#v want: %#v", mj.Type, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobFrequency(mj *monitor.Job, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Frequency != expected {
			return fmt.Errorf("mj.Frequency: got: %#v want: %#v", mj.Frequency, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobRapidRecheck(mj *monitor.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.RapidRecheck != expected {
			return fmt.Errorf("mj.RapidRecheck: got: %#v want: %#v", mj.RapidRecheck, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobPolicy(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Policy != expected {
			return fmt.Errorf("mj.Policy: got: %#v want: %#v", mj.Policy, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobConfigSend(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Config["send"].(string) != expected {
			return fmt.Errorf("mj.Config.send: got: %#v want: %#v", mj.Config["send"].(string), expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobConfigPort(mj *monitor.Job, expected float64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Config["port"].(float64) != expected {
			return fmt.Errorf("mj.Config.port: got: %#v want: %#v", mj.Config["port"].(float64), expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobConfigHost(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Config["host"].(string) != expected {
			return fmt.Errorf("mj.Config.host: got: %#v want: %#v", mj.Config["host"].(string), expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobRuleValue(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Rules[0].Value.(string) != expected {
			return fmt.Errorf("mj.Rules[0].Value: got: %#v want: %#v", mj.Rules[0].Value.(string), expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobRuleComparison(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Rules[0].Comparison != expected {
			return fmt.Errorf("mj.Rules[0].Comparison: got: %#v want: %#v", mj.Rules[0].Comparison, expected)
		}
		return nil
	}
}

func testAccCheckMonitoringJobRuleKey(mj *monitor.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Rules[0].Key != expected {
			return fmt.Errorf("mj.Rules[0].Key: got: %#v want: %#v", mj.Rules[0].Key, expected)
		}
		return nil
	}
}
func testAccCheckMonitoringJobMute(mj *monitor.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mj.Mute != expected {
			return fmt.Errorf("mj.mute: got: %#v want: %#v", mj.Mute, expected)
		}
		return nil
	}
}

// Simulate a manual deletion of a monitoring job.
func testAccManualDeleteMonitoringJob(mj *monitor.Job) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Jobs.Delete(mj.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete monitoring job: %v", err)
		}
	}
}

const testAccMonitoringJobBasic = `
resource "ns1_monitoringjob" "it" {
  job_type = "tcp"
  name     = "terraform test"

  regions   = ["lga"]
  frequency = 60
  mute      = true

  config = {
    ssl = "1",
    send = "HEAD / HTTP/1.0\r\n\r\n"
    port = 443
    host = "1.2.3.4"
  }
  rules {
    value = "200 OK"
    comparison = "contains"
    key = "output"
  }
}
`

const testAccMonitoringJobUpdated = `
resource "ns1_monitoringjob" "it" {
  job_type = "tcp"
  name     = "terraform test"

  active        = true
  regions       = ["lga"]
  frequency     = 120
  rapid_recheck = true
  policy        = "all"
  mute      = false

  config = {
    ssl = "1",
    send = "HEAD / HTTP/1.0\r\n\r\n"
    port = 443
    host = "1.1.1.1"
  }
  rules {
    value = 200
    comparison = "<="
    key = "connect"
  }
}
`
