package ns1

import (
	"fmt"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
	"log"
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
	d := pulsar.DefaultConfig{
		Http:                 false,
		Https:                false,
		RequestTimeoutMillis: 0,
		JobTimeoutMillis:     0,
		UseXhr:               false,
		StaticValues:         false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationBasic(applicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.it", &application),
					testAccCheckApplicationName(&application, applicationName),
					testAccCheckApplicationBrowser(&application, 0),
					testAccCheckApplicationJobs(&application, 0),
					testAccCheckApplicationDefaultConfig(&application, d),
				),
			},
		},
	})
}

func TestAccApplication_updated(t *testing.T) {
	var application pulsar.Application
	applicationName := fmt.Sprintf(
		"terraform-test-%s.io",
		acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum),
	)
	basicConfig := pulsar.DefaultConfig{
		Http:                 false,
		Https:                false,
		RequestTimeoutMillis: 0,
		JobTimeoutMillis:     0,
		UseXhr:               false,
		StaticValues:         false,
	}

	updatedConfig := pulsar.DefaultConfig{
		Http:                 true,
		Https:                false,
		RequestTimeoutMillis: 100,
		JobTimeoutMillis:     100,
		UseXhr:               false,
		StaticValues:         true,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationBasic(applicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.it", &application),
					testAccCheckApplicationName(&application, applicationName),
					testAccCheckApplicationBrowser(&application, 0),
					testAccCheckApplicationJobs(&application, 0),
					testAccCheckApplicationDefaultConfig(&application, basicConfig),
				),
			},
			{
				Config: testAccApplicationUpdated(applicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.it", &application),
					testAccCheckApplicationName(&application, applicationName),
					testAccCheckApplicationBrowser(&application, 123),
					testAccCheckApplicationJobs(&application, 100),
					testAccCheckApplicationDefaultConfig(&application, updatedConfig),
				),
			},
		},
	})
}

func TestAccApplication_ManualDelete(t *testing.T) {
	application := pulsar.Application{}
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
				Check:  testAccCheckApplicationExists("ns1_application.it", &application),
			},
			// Simulate a manual deletion of the application and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteApplication(&application),
				Config:             testAccApplicationBasic(applicationName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccApplicationBasic(applicationName),
				Check:  testAccCheckApplicationExists("ns1_application.it", &application),
			},
		},
	})
}

func testAccCheckApplicationDefaultConfig(app *pulsar.Application, expected pulsar.DefaultConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.DefaultConfig != expected {
			return fmt.Errorf("application DefaultConfig: got: %v want: %v", app.DefaultConfig, expected)
		}
		return nil
	}
}

func testAccCheckApplicationJobs(app *pulsar.Application, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.JobsPerTransaction != expected {
			return fmt.Errorf("application JobsPerTransaction: got: %d want: %d", app.JobsPerTransaction, expected)
		}
		return nil
	}
}

func testAccCheckApplicationBrowser(app *pulsar.Application, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.BrowserWaitMillis != expected {
			return fmt.Errorf("application BrowserWaitMillis: got: %d want: %d", app.BrowserWaitMillis, expected)
		}
		return nil
	}
}

func testAccCheckApplicationName(app *pulsar.Application, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if app.Name != expected {
			return fmt.Errorf("application name: got: %s want: %s", app.Name, expected)
		}
		return nil
	}
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
			return fmt.Errorf("application not found")
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

// Simulate a manual deletion of a application.
func testAccManualDeleteApplication(application *pulsar.Application) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Applications.Delete(application.ID)
		if err != nil {
			log.Printf("failed to delete application: %v", err)
		}
	}
}

func testAccApplicationBasic(appName string) string {
	return fmt.Sprintf(`resource "ns1_application" "it" {
 name = "%s"
}
`, appName)

}

func testAccApplicationUpdated(appName string) string {
	return fmt.Sprintf(`resource "ns1_application" "it" {
 name = "%s"
 browser_wait_millis = 123
 jobs_per_transaction = 100
 default_config = {
  http     = true
  https = false
  request_timeout_millis = 100
  job_timeout_millis = 100
  static_values = true
 }
}
`, appName)

}
