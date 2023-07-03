package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
)

// Creating basic Pulsar jobs
func TestAccPulsarJob_basic(t *testing.T) {
	var (
		job     = pulsar.Job{}
		jobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		appName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)
	// Basic test for JavaScript jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJSPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccURLPath"),
				),
			},
		},
	})

	// Basic test for Bulk Beacon jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBBPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
				),
			},
		},
	})
}

// Updating pulsar jobs without changing its type (JavaScript or Bulk Beacon)
func TestAccPulsarJob_updated_same_type(t *testing.T) {
	var (
		app     = pulsar.Application{}
		job     = pulsar.Job{}
		jobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		appName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		updatedJobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)
	// Update test for JavaScript jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJSPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccURLPath"),
				),
			},
			{
				Config: testAccJSPulsarJobUpdated(appName, updatedJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, updatedJobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, false),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccUpdatedHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccUpdatedURLPath"),
				),
			},
		},
	})

	// update test for Bulk Beacon jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBBPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
				),
			},
			{
				Config: testAccBBPulsarJobUpdated(appName, updatedJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, updatedJobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, false),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
				),
			},
		},
	})
}

// Updating pulsar jobs changing its type (JavaScript <-> Bulk Beacon)
func TestAccPulsarJob_updated_different_type(t *testing.T) {
	var (
		job     = pulsar.Job{}
		jobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		app     = pulsar.Application{}
		appName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		updatedJobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)

	// Updating JavaScript job to Bulk Beacon Job
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJSPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccURLPath"),
				),
			},
			{
				Config: testAccBBPulsarJobConverted(appName, updatedJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, updatedJobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
				),
			},
		},
	})

	// Updating Bulk Beacon Job to JavaScript job
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBBPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
				),
			},
			{
				Config: testAccJSPulsarJobUpdated(appName, updatedJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, updatedJobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, false),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccUpdatedHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccUpdatedURLPath"),
				),
			},
		},
	})
}

// Creating Pulsar jobs with Blend Metric Weights
func TestAccPulsarJob_BlendMetricWeights(t *testing.T) {
	var (
		job     = pulsar.Job{}
		jobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		app     = pulsar.Application{}
		appName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		weights = []pulsar.Weights{
			{
				Name:         "testAccWeight1",
				Weight:       123,
				DefaultValue: 12.3,
				Maximize:     false,
			},
			{
				Name:         "testAccWeight2",
				Weight:       321,
				DefaultValue: 32.1,
				Maximize:     true,
			},
		}
	)

	// Blend Metric Weights test for JavaScript jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJSPulsarJobBlendMetricWeights(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "latency"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCheckPulsarJobSHost(&job, "testAccCompleteHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccCompleteURLPath"),
					testAccCHeckPulsarJobBlendMetricWeights_timestamp(&job, 123),
					testAccCHeckPulsarJobBlendMetricWeights_weights(&job, weights),
				),
			},
		},
	})

	// Blend Metric Weights test for Bulk Beacon jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBBPulsarJobBlendMetricWeights(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists("ns1_application.app", &app),
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
					testAccCheckPulsarJobName(&job, jobName),
					testAccCheckPulsarJobTypeID(&job, "custom"),
					testAccCheckPulsarJobSHost(&job, "testAccHost"),
					testAccCheckPulsarJobSUrlPath(&job, "/testAccUrlPath"),
					testAccCheckPulsarJobAppID(&job, &app),
					testAccCheckPulsarJobActive(&job, true),
					testAccCheckPulsarJobShared(&job, false),
					testAccCheckPulsarJobSCommunity(&job, false),
					testAccCHeckPulsarJobBlendMetricWeights_timestamp(&job, 123),
					testAccCHeckPulsarJobBlendMetricWeights_weights(&job, weights),
				),
			},
		},
	})
}

// Manually deleting Pulsar Jobs
func TestAccPulsarJob_ManualDelete(t *testing.T) {
	var (
		job     = pulsar.Job{}
		jobName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

		appName = fmt.Sprintf("terraform-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
	)
	// Manual deletion test for JavaScript jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJSPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
				),
			},
			// Simulate a manual deletion of the pulsar job and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeletePulsarJob(&job),
				Config:             testAccJSPulsarJobBasic(appName, jobName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccJSPulsarJobBasic(appName, jobName),
				Check:  testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
			},
		},
	})

	// Manual deletion test for Bulk Beacon jobs
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPulsarJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBBPulsarJobBasic(appName, jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
				),
			},
			// Simulate a manual deletion of the pulsar job and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeletePulsarJob(&job),
				Config:             testAccBBPulsarJobBasic(appName, jobName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccBBPulsarJobBasic(appName, jobName),
				Check:  testAccCheckPulsarJobExists("ns1_pulsarjob.it", &job),
			},
		},
	})
}

func testAccJSPulsarJobBasic(appName, jobName string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
	name = "%s"
}

	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "latency"
		app_id = "${ns1_application.app.id}"
		config {
			host = "testAccHost"
			url_path = "/testAccURLPath"
		}
}
`, appName, jobName)
}

func testAccJSPulsarJobUpdated(appName, jobName string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
		name = "%s"
	}
	
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "latency"
		app_id = "${ns1_application.app.id}"
		active = false
		shared = false
		config {
			host = "testAccUpdatedHost"
			url_path = "/testAccUpdatedURLPath"
		}
}
`, appName, jobName)
}

func testAccJSPulsarJobBlendMetricWeights(appName, jobName string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
		name = "%s"
	}
		
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "latency"
		app_id = "${ns1_application.app.id}"
		config {
			host = "testAccCompleteHost"
			url_path = "/testAccCompleteURLPath"
		}
		blend_metric_weights {
			timestamp = 123
		}
		weights {
			name = "testAccWeight1"
			weight = 123
			default_value = 12.3
			maximize = false
		}
		weights {
			name = "testAccWeight2"
			weight = 321
			default_value = 32.1
			maximize = true
		}
}
`, appName, jobName)
}

func testAccBBPulsarJobBasic(appName string, jobName string) string {
	return fmt.Sprintf(`
	resource "ns1_application" "app" {
		name = "%s"
	}
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "custom"
		app_id = "${ns1_application.app.id}"
}
`, appName, jobName)
}

func testAccBBPulsarJobUpdated(jobName, appId string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
		name = "%s"
	}
	
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "custom"
		app_id = "${ns1_application.app.id}"
		active = false
		shared = false
}
`, jobName, appId)
}

func testAccBBPulsarJobConverted(appName, jobName string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
		name = "%s"
	}
	
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "custom"
		app_id = "${ns1_application.app.id}"
		config {
			host = ""
			url_path = ""
		}
}
`, appName, jobName)
}

func testAccBBPulsarJobBlendMetricWeights(appName, jobName string) string {
	return fmt.Sprintf(`resource "ns1_application" "app" {
		name = "%s"
	}
	
	resource "ns1_pulsarjob" "it" {
  		name = "%s"
		type_id = "custom"
		app_id = "${ns1_application.app.id}"
		config {
			host = "testAccHost"
			url_path = "/testAccUrlPath"
		}
		blend_metric_weights {
			timestamp = 123
		}
		weights {
			name = "testAccWeight1"
			weight = 123
			default_value = 12.3
			maximize = false
		}
		weights {
			name = "testAccWeight2"
			weight = 321
			default_value = 32.1
			maximize = true
		}
}
`, appName, jobName)
}

func testAccCheckPulsarJobExists(n string, job *pulsar.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundPulsarJob, _, err := client.PulsarJobs.Get(rs.Primary.Attributes["app_id"], rs.Primary.ID)

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundPulsarJob.JobID != p.Attributes["id"] {
			return fmt.Errorf("pulsar job not found")
		}

		*job = *foundPulsarJob

		return nil
	}
}

func testAccCheckPulsarJobDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_pulsarjob" {
			continue
		}

		pulsarJob, _, err := client.PulsarJobs.Get(rs.Primary.Attributes["app_id"], rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("pulsar job still exists: %#v: %#v", err, pulsarJob)
		}

	}

	return nil
}

func testAccCheckPulsarJobName(job *pulsar.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.Name != expected {
			return fmt.Errorf("job.Name: got: %s want: %s", job.Name, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobTypeID(job *pulsar.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.TypeID != expected {
			return fmt.Errorf("job.TypeID: got: %s want: %s", job.TypeID, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobAppID(job *pulsar.Job, app *pulsar.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.AppID != app.ID {
			return fmt.Errorf("job.AppID: got: %s want: %s", job.AppID, app.ID)
		}
		return nil
	}
}

func testAccCheckPulsarJobActive(job *pulsar.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.Active != expected {
			return fmt.Errorf("job.Active: got: %t want: %t", job.Active, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobShared(job *pulsar.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.Shared != expected {
			return fmt.Errorf("job.Shared: got: %t want: %t", job.Shared, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobSCommunity(job *pulsar.Job, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.Community != expected {
			return fmt.Errorf("job.Community: got: %t want: %t", job.Community, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobSHost(job *pulsar.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *job.Config.Host != expected {
			return fmt.Errorf("job.Config.Host: got: %s want: %s", *job.Config.Host, expected)
		}
		return nil
	}
}

func testAccCheckPulsarJobSUrlPath(job *pulsar.Job, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *job.Config.URLPath != expected {
			return fmt.Errorf("job.Config.URL_Path: got: %s want: %s", *job.Config.URLPath, expected)
		}
		return nil
	}
}

func testAccCHeckPulsarJobBlendMetricWeights_timestamp(job *pulsar.Job, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if job.Config.BlendMetricWeights.Timestamp != expected {
			return fmt.Errorf("job.Config.BlendMetricWeights.Timestamp: got: %v want: %v", job.Config.BlendMetricWeights.Timestamp, expected)
		}
		return nil
	}
}

func testAccCHeckPulsarJobBlendMetricWeights_weights(job *pulsar.Job, expected []pulsar.Weights) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		numberWeights := len(job.Config.BlendMetricWeights.Weights)
		numberExpected := len(expected)
		if numberWeights != numberExpected {
			return fmt.Errorf("job.Config.BlendMetricWeights.Weights: got: %v elements want: %v elements", numberWeights, numberExpected)
		}

		for i, weight := range job.Config.BlendMetricWeights.Weights {
			if weight.Name != expected[i].Name {
				return fmt.Errorf("job.Config.BlendMetricWeights.Weights[%v].Name: got: %s want: %s", i, weight.Name, expected[i].Name)
			}
			if weight.Weight != expected[i].Weight {
				return fmt.Errorf("job.Config.BlendMetricWeights.Weights[%v].Weight: got: %v want: %v", i, weight.Weight, expected[i].Weight)
			}
			if weight.DefaultValue != expected[i].DefaultValue {
				return fmt.Errorf("job.Config.BlendMetricWeights.Weights[%v].DefaultValue: got: %v want: %v", i, weight.DefaultValue, expected[i].DefaultValue)
			}
			if weight.Maximize != expected[i].Maximize {
				return fmt.Errorf("job.Config.BlendMetricWeights.Weights[%v].Maximize: got: %v want: %v", i, weight.Maximize, expected[i].Maximize)
			}
		}
		return nil
	}
}

func testAccManualDeletePulsarJob(job *pulsar.Job) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.PulsarJobs.Delete(job)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete pulsar job: %v", err)
		}
	}
}
