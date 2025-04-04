package ns1

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	billingusage "gopkg.in/ns1/ns1-go.v2/rest/model/billingusage"
)

func TestAccBillingUsage_queries(t *testing.T) {
	var queries billingusage.Queries
	from := int32(time.Now().AddDate(0, -1, 0).Unix())
	to := int32(time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageQueries(from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageQueriesExists("data.ns1_billing_usage.test", &queries),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeQueries),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "from", fmt.Sprintf("%d", from)),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "to", fmt.Sprintf("%d", to)),
				),
			},
		},
	})
}

func TestAccBillingUsage_limits(t *testing.T) {
	var limits billingusage.Limits
	from := int32(time.Now().AddDate(0, -1, 0).Unix())
	to := int32(time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageLimits(from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageLimitsExists("data.ns1_billing_usage.test", &limits),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeLimits),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "from", fmt.Sprintf("%d", from)),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "to", fmt.Sprintf("%d", to)),
				),
			},
		},
	})
}

func TestAccBillingUsage_decisions(t *testing.T) {
	var usage billingusage.TotalUsage
	from := int32(time.Now().AddDate(0, -1, 0).Unix())
	to := int32(time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageDecisions(from, to),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageDecisionsExists("data.ns1_billing_usage.test", &usage),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeDecisions),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "from", fmt.Sprintf("%d", from)),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "to", fmt.Sprintf("%d", to)),
				),
			},
		},
	})
}

func TestAccBillingUsage_filter_chains(t *testing.T) {
	var usage billingusage.TotalUsage

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageFilterChains(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageFilterChainsExists("data.ns1_billing_usage.test", &usage),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeFilterChains),
				),
			},
		},
	})
}

func TestAccBillingUsage_monitors(t *testing.T) {
	var usage billingusage.TotalUsage

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageMonitors(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageMonitorsExists("data.ns1_billing_usage.test", &usage),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeMonitors),
				),
			},
		},
	})
}

func TestAccBillingUsage_records(t *testing.T) {
	var usage billingusage.TotalUsage

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingUsageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingUsageRecords(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBillingUsageRecordsExists("data.ns1_billing_usage.test", &usage),
					resource.TestCheckResourceAttr("data.ns1_billing_usage.test", "metric_type", MetricTypeRecords),
				),
			},
		},
	})
}

func testAccCheckBillingUsageQueriesExists(n string, queries *billingusage.Queries) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		from, err := strconv.ParseInt(rs.Primary.Attributes["from"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing from timestamp: %v", err)
		}

		to, err := strconv.ParseInt(rs.Primary.Attributes["to"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing to timestamp: %v", err)
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetQueries(
			int32(from),
			int32(to),
		)
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage queries not found")
		}

		*queries = *found
		return nil
	}
}

func testAccCheckBillingUsageLimitsExists(n string, limits *billingusage.Limits) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		from, err := strconv.ParseInt(rs.Primary.Attributes["from"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing from timestamp: %v", err)
		}

		to, err := strconv.ParseInt(rs.Primary.Attributes["to"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing to timestamp: %v", err)
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetLimits(
			int32(from),
			int32(to),
		)
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage limits not found")
		}

		*limits = *found
		return nil
	}
}

func testAccCheckBillingUsageDecisionsExists(n string, usage *billingusage.TotalUsage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		from, err := strconv.ParseInt(rs.Primary.Attributes["from"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing from timestamp: %v", err)
		}

		to, err := strconv.ParseInt(rs.Primary.Attributes["to"], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing to timestamp: %v", err)
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetDecisions(
			int32(from),
			int32(to),
		)
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage decisions not found")
		}

		*usage = *found
		return nil
	}
}

func testAccCheckBillingUsageFilterChainsExists(n string, usage *billingusage.TotalUsage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetFilterChains()
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage filter chains not found")
		}

		*usage = *found
		return nil
	}
}

func testAccCheckBillingUsageMonitorsExists(n string, usage *billingusage.TotalUsage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetMonitors()
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage monitors not found")
		}

		*usage = *found
		return nil
	}
}

func testAccCheckBillingUsageRecordsExists(n string, usage *billingusage.TotalUsage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)
		found, _, err := client.BillingUsage.GetRecords()
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Billing usage records not found")
		}

		*usage = *found
		return nil
	}
}

func testAccCheckBillingUsageDestroy(s *terraform.State) error {
	// Billing usage is a data source, so it doesn't need a destroy check
	return nil
}

func testAccBillingUsageQueries(from, to int32) string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
  from = %d
  to = %d
}
`, MetricTypeQueries, from, to)
}

func testAccBillingUsageLimits(from, to int32) string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
  from = %d
  to = %d
}
`, MetricTypeLimits, from, to)
}

func testAccBillingUsageDecisions(from, to int32) string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
  from = %d
  to = %d
}
`, MetricTypeDecisions, from, to)
}

func testAccBillingUsageFilterChains() string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
}
`, MetricTypeFilterChains)
}

func testAccBillingUsageMonitors() string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
}
`, MetricTypeMonitors)
}

func testAccBillingUsageRecords() string {
	return fmt.Sprintf(`
data "ns1_billing_usage" "test" {
  metric_type = "%s"
}
`, MetricTypeRecords)
}
