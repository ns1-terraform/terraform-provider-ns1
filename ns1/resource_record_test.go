package ns1

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccRecord_basic(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 10),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/CNAME", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_updated(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 10),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				Config: testAccRecordUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 120),
					testAccCheckRecordUseClientSubnet(&record, false),
					testAccCheckRecordRegionName(&record, []string{"ny", "wa"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 5),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test2.%s", zoneName)},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/CNAME", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_meta(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)
	sub := make(map[string]interface{}, 2)
	sub["BR"] = []interface{}{"SP", "SC"}
	sub["DZ"] = []interface{}{"01", "02", "03"}
	sub["NO"] = []interface{}{"NO-01", "NO-11"}
	sub["SG"] = []interface{}{"SG-03"}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordAnswerMeta(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordAnswerMetaSubdivisions(&record, sub),
					testAccCheckRecordAnswerMetaIPPrefixes(&record, []string{
						"3.248.0.0/13",
						"13.248.96.0/24",
						"13.248.113.0/24",
						"13.248.118.0/24",
						"13.248.119.0/24",
						"13.248.121.0/24",
					}),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/A", zoneName, domainName),
				ImportStateVerify: true,
			},
			{
				Config: testAccRecordAnswerMetaDataFeed(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordAnswerMetaUp("ns1_datafeed.test", &record),
				),
			},
		},
	})
}

func TestAccRecord_meta_with_json(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)
	sub := make(map[string]interface{}, 2)
	sub["BR"] = []interface{}{"SP", "SC"}
	sub["DZ"] = []interface{}{"01", "02", "03"}
	sub["NO"] = []interface{}{"NO-01", "NO-11"}
	sub["SG"] = []interface{}{"SG-03"}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordAnswerMetaWithJson(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordAnswerMetaSubdivisions(&record, sub),
					testAccCheckRecordAnswerMetaIPPrefixes(&record, []string{
						"3.248.0.0/13",
						"13.248.96.0/24",
						"13.248.113.0/24",
						"13.248.118.0/24",
						"13.248.119.0/24",
						"13.248.121.0/24",
					}),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/A", zoneName, domainName),
				ImportStateVerify: true,
			},
			{
				Config: testAccRecordAnswerMetaDataFeed(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordAnswerMetaUp("ns1_datafeed.test", &record),
				),
			},
		},
	})
}

func TestAccRecord_CAA(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordCAA(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.caa", &record),
					testAccCheckRecordDomain(&record, zoneName),
					testAccCheckRecordTTL(&record, 3600),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{"0", "issue", "letsencrypt.org"},
					),
					testAccCheckRecordAnswerRdata(
						t, &record, 1, []string{"0", "issuewild", ";"},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%[1]s/%[1]s/CAA", zoneName),
				ImportStateVerify: true,
			},
			{
				Config: testAccRecordCAAWithSpace(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.caa", &record),
					testAccCheckRecordDomain(&record, zoneName),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{"0", "issue", "inbox2221.ticket; account=xyz"},
					),
				),
			},
		},
	})
}

func TestAccRecord_SPF(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordSPF(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.spf", &record),
					testAccCheckRecordDomain(&record, zoneName),
					testAccCheckRecordTTL(&record, 86400),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{"v=DKIM1; k=rsa; p=XXXXXXXX"},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%[1]s/%[1]s/SPF", zoneName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_SRV(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("_some-server._tcp.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordSRV(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.srv", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 86400),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(
						t,
						&record,
						0,
						[]string{"10", "0", "2380", fmt.Sprintf("node-1.%s", zoneName)},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/SRV", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_DS(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("_some-server._tcp.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordDS(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.ds", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 86400),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(
						t,
						&record,
						0,
						[]string{"262", "13", "2", "287787bd551bcab4f57d0c1dcaf312eebe36cc338bebb90d1402353c7096785d"},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/DS", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_URLFWD(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("fwd.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccURLFWDPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordURLFWD(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.urlfwd", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 3600),
					testAccCheckRecordAnswerRdata(
						t,
						&record,
						0,
						[]string{"/", "https://example.com", "301", "2", "0"},
					),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/URLFWD", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_validationError(t *testing.T) {
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordInvalid(rString),
				/* The error block should look like this:

				config is invalid: 4 problems:

					- zone has an invalid leading ".", got: .terraform-test-e677cntkkar21ak.io.
					- zone has an invalid trailing ".", got: .terraform-test-e677cntkkar21ak.io.
					- domain has an invalid leading ".", got: .test.terraform-test-e677cntkkar21ak.io.
					- domain has an invalid trailing ".", got: .test.terraform-test-e677cntkkar21ak.io.

				*/
				ExpectError: regexp.MustCompile(`config is invalid: 4 problems:\n\n(\s*- (zone|domain) has an invalid (leading|trailing) \"\.\", got: .*){4}`),
			},
		},
	})
}

// Verifies that a record is re-created correctly if it is manually deleted.
func TestAccRecord_ManualDelete(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("test.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordBasic(rString),
				Check:  testAccCheckRecordExists("ns1_record.it", &record),
			},
			// Simulate a manual deletion of the record and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteRecord(zoneName, domainName, "CNAME"),
				Config:             testAccRecordBasic(rString),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccRecordBasic(rString),
				Check:  testAccCheckRecordExists("ns1_record.it", &record),
			},
		},
	})
}

func TestAccRecord_CaseInsensitive(t *testing.T) {
	var CapitalLettersCases = []struct {
		name   string
		domain string
		zone   string
	}{
		{
			"root level record no cap letters",
			"terraform-test-#.io",
			"terraform-test-#.io",
		},
		{
			"root level record domain as title",
			"Terraform-test-#.io",
			"terraform-test-#.io",
		},
		{
			"root level record zone as title",
			"terraform-test-#.io",
			"Terraform-test-#.io",
		},
		{
			"root level record zone and domain as title",
			"Terraform-test-#.io",
			"Terraform-test-#.io",
		},
		{
			"root level record zone and domain random capitalization",
			"TerrAForm-test-#.io",
			"TerrafORm-test-#.io",
		},
		{
			"record no cap letters",
			"test.terraform-test-#.io",
			"terraform-test-#.io",
		},
		{
			"record domain as title",
			"test.Terraform-test-#.io",
			"terraform-test-#.io",
		},
		{
			"record zone as title",
			"test.terraform-test-#.io",
			"Terraform-test-#.io",
		},
		{
			"record zone and domain as title",
			"test.Terraform-test-#.io",
			"Terraform-test-#.io",
		},
	}
	for _, tt := range CapitalLettersCases {
		var record dns.Record

		rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

		zoneName := strings.Replace(tt.zone, "#", rString, 1)
		domainName := strings.Replace(tt.domain, "#", rString, 1)

		tfFile := testAccRecordBasicCaseSensitive(domainName, zoneName)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckRecordDestroy,
			Steps: []resource.TestStep{
				// Simulate an apply
				{
					Config: tfFile,
					Check:  testAccCheckRecordExists("ns1_record.it", &record),
				},
				// Simulate a plan to check if has any diff
				{
					Config:             tfFile,
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
				},
			},
		})

	}

}

func TestAccRecord_OverrideTTLNilToTrue(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileBasicALIAS := testAccRecordBasicALIAS(rString)
	tfFileOverrideTtlALIAS := testAccRecordBasicALIASOverrideTTL(rString, true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_ttl not set
			{
				Config: tfFileBasicALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileBasicALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Change override TTL to true Plan and Apply
			{
				Config:             tfFileOverrideTtlALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: tfFileOverrideTtlALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNotNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideTtlALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRecord_OverrideTTLTrueToNil(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileBasicALIAS := testAccRecordBasicALIAS(rString)
	tfFileOverrideTtlALIAS := testAccRecordBasicALIASOverrideTTL(rString, true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_ttl true
			{
				Config: tfFileOverrideTtlALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNotNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideTtlALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:             tfFileBasicALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Change override TTL to false setting to "null"
			{
				Config: tfFileBasicALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileBasicALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRecord_OverrideTTLTrueToFalse(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileOverrideTtlALIAS := testAccRecordBasicALIASOverrideTTL(rString, true)
	tfFileDoNotOverrideTtlALIAS := testAccRecordBasicALIASOverrideTTL(rString, false)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_ttl true
			{
				Config: tfFileOverrideTtlALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNotNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideTtlALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Change override TTL to false setting to false Plan and Apply
			{
				Config:             tfFileDoNotOverrideTtlALIAS,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: tfFileDoNotOverrideTtlALIAS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideTTL(&record, ExpectOverrideTTLNil()),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileDoNotOverrideTtlALIAS,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRecord_Link(t *testing.T) {
	var record1 dns.Record
	var record2 dns.Record
	var record3 dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	linkDomainName := fmt.Sprintf("link.%s", zoneName)
	targetDomainName := fmt.Sprintf("target.%s", zoneName)
	newTargetDomainName := fmt.Sprintf("newtarget.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordLink(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.link", &record1),
					testAccCheckRecordExists("ns1_record.target", &record2),
					testAccCheckRecordDomain(&record1, linkDomainName),
					testAccCheckRecordDomain(&record2, targetDomainName),
					testAccCheckRecordTTL(&record1, 666),
					testAccCheckRecordTTL(&record2, 777),
					testAccCheckRecordAnswerRdata(
						t, &record2, 0, []string{"99.86.99.86"},
					),
				),
			},
			{
				Config: testAccRecordLinkUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.link", &record1),
					testAccCheckRecordExists("ns1_record.target", &record2),
					testAccCheckRecordExists("ns1_record.newtarget", &record3),
					testAccCheckRecordDomain(&record1, linkDomainName),
					testAccCheckRecordDomain(&record2, targetDomainName),
					testAccCheckRecordDomain(&record3, newTargetDomainName),
					testAccCheckRecordTTL(&record1, 666),
					testAccCheckRecordTTL(&record2, 777),
					testAccCheckRecordTTL(&record3, 888),
					testAccCheckRecordAnswerRdata(
						t, &record2, 0, []string{"99.86.99.86"},
					),
					testAccCheckRecordAnswerRdata(
						t, &record3, 0, []string{"16.19.20.19"},
					),
				),
			},
		},
	})
}

// See if URLFWD permission exist by trying to create a record with it.
func testAccURLFWDPreCheck(t *testing.T) {
	client, err := sharedClient()
	if err != nil {
		log.Fatalf("failed to get shared client: %s", err)
	}

	name := fmt.Sprintf("terraform-urlfwd-test-%s.io", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	_, err = client.Zones.Create(&dns.Zone{Zone: name})
	if err != nil {
		t.Fatalf("failed to create test zone %s: %s", name, err)
	}

	record := dns.NewRecord(name, fmt.Sprintf("domain.%s", name), "URLFWD")
	record.Answers = []*dns.Answer{{Rdata: []string{"/", "https://example.com", "301", "2", "0"}}}

	_, err = client.Records.Create(record)
	if err != nil {
		if strings.Contains(err.Error(), "400 URLFWD records are not enabled") {
			t.Skipf("account not authorized for URLFWD records, skipping test")
			return
		}

		t.Fatalf("failed to create test record %s: %s", name, err)
	}

	_, err = client.Zones.Delete(record.Zone)
	if err != nil {
		t.Fatalf("failed to delete test record %s", name)
	}
}

func testAccCheckRecordExists(n string, record *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %v", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("noID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		p := rs.Primary

		foundRecord, _, err := client.Records.Get(p.Attributes["zone"], p.Attributes["domain"], p.Attributes["type"])
		if err != nil {
			return fmt.Errorf("record not found")
		}

		if foundRecord.Domain != p.Attributes["domain"] {
			return fmt.Errorf("record not found")
		}

		*record = *foundRecord

		return nil
	}
}

func testAccCheckRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	var recordDomain string
	var recordZone string
	var recordType string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_record" {
			continue
		}

		if rs.Type == "ns1_record" {
			recordType = rs.Primary.Attributes["type"]
			recordDomain = rs.Primary.Attributes["domain"]
			recordZone = rs.Primary.Attributes["zone"]
		}
	}

	foundRecord, _, err := client.Records.Get(recordZone, recordDomain, recordType)
	if err != ns1.ErrRecordMissing {
		return fmt.Errorf("record still exists: %#v %#v", foundRecord, err)
	}

	return nil
}

func testAccCheckRecordDomain(r *dns.Record, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if r.Domain != expected {
			return fmt.Errorf("r.Domain: got: %#v want: %#v", r.Domain, expected)
		}
		return nil
	}
}

func testAccCheckRecordTTL(r *dns.Record, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if r.TTL != expected {
			return fmt.Errorf("r.TTL: got: %#v want: %#v", r.TTL, expected)
		}
		return nil
	}
}

func testAccCheckRecordUseClientSubnet(r *dns.Record, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *r.UseClientSubnet != expected {
			return fmt.Errorf("UseClientSubnet: got: %#v want: %#v", *r.UseClientSubnet, expected)
		}
		return nil
	}
}

func ExpectOverrideTTLNil() bool {
	return true
}

func ExpectOverrideTTLNotNil() bool {
	return false
}

func testAccCheckRecordOverrideTTL(r *dns.Record, expectedNil bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if expectedNil {
			if r.Override_TTL != nil {
				return fmt.Errorf("Override TTL: got: %#v want: null", *r.Override_TTL)
			}
			return nil
		}
		if r.Override_TTL == nil {
			return fmt.Errorf("Override TTL: got: %v want: notNil", r.Override_TTL)
		}
		return nil
	}
}

func testAccCheckRecordRegionName(r *dns.Record, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		regions := make([]string, len(r.Regions))

		i := 0
		for k := range r.Regions {
			regions[i] = k
			i++
		}
		sort.Strings(regions)
		sort.Strings(expected)
		if !reflect.DeepEqual(regions, expected) {
			return fmt.Errorf("regions: got: %#v want: %#v", regions, expected)
		}
		return nil
	}
}

func testAccCheckRecordAnswerMetaSubdivisions(r *dns.Record, expected map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswer := r.Answers[0]
		recordMetas := recordAnswer.Meta
		recordSubdiv := recordMetas.Subdivisions.(map[string]interface{})

		recordSlice := mapToSlice(recordSubdiv)
		expectedSlice := mapToSlice(expected)

		if !reflect.DeepEqual(recordSlice, expectedSlice) {
			return fmt.Errorf("r.Answers[0].Meta.Subdivisions: got: %#v want: %#v", recordSubdiv, expected)
		}

		return nil
	}
}

func mapToSlice(m map[string]interface{}) []string {
	sliceString := make([]string, 0)
	for key, slice := range m {
		for _, v := range slice.([]interface{}) {
			sliceString = append(sliceString, fmt.Sprintf("%v-%v", key, v))
		}
	}

	sort.Strings(sliceString)
	return sliceString
}

func testAccCheckRecordAnswerMetaWeight(r *dns.Record, expected float64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswer := r.Answers[0]
		recordMetas := recordAnswer.Meta
		weight := recordMetas.Weight.(float64)
		if weight != expected {
			return fmt.Errorf("r.Answers[0].Meta.Weight: got: %#v want: %#v", weight, expected)
		}
		return nil
	}
}

func testAccCheckRecordAnswerMetaIPPrefixes(r *dns.Record, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswer := r.Answers[0]
		recordMetas := recordAnswer.Meta
		ipPrefixes := make([]string, len(recordMetas.IPPrefixes.([]interface{})))
		for i, v := range recordMetas.IPPrefixes.([]interface{}) {
			ipPrefixes[i] = v.(string)
		}

		sort.Strings(ipPrefixes)
		sort.Strings(expected)
		if !reflect.DeepEqual(ipPrefixes, expected) {
			return fmt.Errorf("ip-prefixes: got: %#v want: %#v", ipPrefixes, expected)
		}
		return nil
	}
}

func testAccCheckRecordAnswerMetaUp(expected interface{}, r *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswer := r.Answers[0]
		recordMetas := recordAnswer.Meta
		up := recordMetas.Up
		switch up.(type) {
		case bool:
			if expected.(bool) != recordMetas.Up {
				return fmt.Errorf("r.Answers[0].Meta.Up: got: %#v want: %#v", up, expected)
			}
		case map[string]interface{}:
			// feed mapping: expected points us to the datafeed, which has the id
			//               that should be in our map
			rs, ok := s.RootModule().Resources[expected.(string)]
			if !ok {
				return fmt.Errorf("resource not found in state: %v", expected)
			}
			ch := map[string]interface{}{"feed": rs.Primary.ID}
			if !reflect.DeepEqual(ch, up) {
				return fmt.Errorf("r.Answers[0].Meta.Up: map: %#v want: %#v", up, ch)
			}
		}
		return nil
	}
}

func testAccCheckRecordAnswerRdata(
	t *testing.T, r *dns.Record, answerIdx int, expected []string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswerRdata := r.Answers[answerIdx].Rdata
		assert.ElementsMatch(t, recordAnswerRdata, expected)
		return nil
	}
}

// Simulate a manual deletion of a record.
func testAccManualDeleteRecord(zone, domain, recordType string) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Records.Delete(zone, domain, recordType)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete record: %v", err)
		}
	}
}

func testAccRecordBasic(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"
  ttl               = 60

  // meta {
  //   weight = 5
  //   connections = 3
  //   // up = false // Ignored by d.GetOk("meta.0.up") due to known issue
  // }

  answers {
    answer = "test1.${ns1_zone.test.zone}"
    region = "cal"

    // meta {
    //   weight = 10
    //   up = true
    // }
  }

  regions {
    name = "cal"
    // meta {
    //   up = true
    //   us_state = ["CA"]
    // }
  }

  filters {
    filter = "geotarget_country"
  }

  filters {
    filter = "select_first_n"
    config = {N=1}
  }

  filters {
    filter = "up"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordBasicCaseSensitive(domain string, zone string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "%s"
  type              = "CNAME"
  ttl               = 60

  answers {
    answer = "test1.${ns1_zone.test.zone}"
    region = "cal"

    // meta {
    //   weight = 10
    //   up = true
    // }
  }

  regions {
    name = "cal"
    // meta {
    //   up = true
    //   us_state = ["CA"]
    // }
  }

}

resource "ns1_zone" "test" {
  zone = "%s"
}
`, domain, zone)
}

func testAccRecordBasicALIASOverrideTTL(rString string, overridettl bool) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "${ns1_zone.test.zone}"
  type              = "ALIAS"
  ttl               = 60
  override_ttl 		= %v
  answers {
    answer = "test.${ns1_zone.test.zone}"
  }
}

resource "ns1_zone" "test" {
	zone = "terraform-test-%s.io"
}
`, overridettl, rString)
}

func testAccRecordBasicALIAS(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "${ns1_zone.test.zone}"
  type              = "ALIAS"
  ttl               = 60
  answers {
    answer = "test.${ns1_zone.test.zone}"
  }
}

resource "ns1_zone" "test" {
	zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordUpdated(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"
  ttl               = 120
  use_client_subnet = false

  // meta {
  //   weight = 5
  //   connections = 3
  //   // up = false // Ignored by d.GetOk("meta.0.up") due to known issue
  // }

  answers {
    answer = "test2.${ns1_zone.test.zone}"
    region = "ny"

    // meta {
    //   weight = 5
    //   up = true
    // }
  }

  regions {
    name = "ny"
    // meta {
    //   us_state = ["NY"]
    // }
  }

  regions {
    name = "wa"
    // meta {
    //   us_state = ["WA"]
    // }
  }

  filters {
    filter = "select_first_n"
    config = {N=1}
  }

  filters {
    filter = "geotarget_country"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordAnswerMeta(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "A"
  answers {
    answer = "1.2.3.4"

    meta = {
	  up = true
	  subdivisions = "BR-SP,BR-SC,DZ-01,DZ-02,DZ-03,NO-NO-01,NO-NO-11,SG-SG-03"
      weight = 5
      ip_prefixes = "3.248.0.0/13,13.248.96.0/24,13.248.113.0/24,13.248.118.0/24,13.248.119.0/24,13.248.121.0/24"
      pulsar = jsonencode([{
        "job_id"     = "abcdef",
        "bias"       = "*0.55",
        "a5m_cutoff" = 0.9
      }])
    }
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordAnswerMetaWithJson(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
	zone              = "${ns1_zone.test.zone}"
	domain            = "test.${ns1_zone.test.zone}"
	type              = "A"
	answers {
	  answer = "1.2.3.4"
  
	  meta = {
		up = true
		subdivisions = jsonencode({
			"BR" = ["SP", "SC"],
			"DZ" = ["01", "02", "03"],
			"NO" = ["NO-01", "NO-11"],
			"SG" = ["SG-03"]
		})
		weight = 5
		ip_prefixes = "3.248.0.0/13,13.248.96.0/24,13.248.113.0/24,13.248.118.0/24,13.248.119.0/24,13.248.121.0/24"
		pulsar = jsonencode([{
		  "job_id"     = "abcdef",
		  "bias"       = "*0.55",
		  "a5m_cutoff" = 0.9
		}])
	  }
	}
  }

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordAnswerMetaDataFeed(rString string) string {
	return fmt.Sprintf(`
resource "ns1_monitoringjob" "test" {
  name = "terraform-test-%s"
  active = true
  regions = [
    "ams"
  ]
  job_type = "http"
  frequency = 60
  rapid_recheck = true
  policy = "all"
  config = {
    method = "GET"
    url = "https://www.example.com"
    connect_timeout = "2000"
    idle_timeout = "3"
    ipv6 = false
    follow_redirect = false
    tls_add_verify = false
    user_agent = "just testing"
  }
  rules {
    value = "200"
    comparison = "=="
    key = "status_code"
  }
}
resource "ns1_datasource" "test" {
  name       = "test datasource"
  sourcetype = "nsone_monitoring"
}
resource "ns1_datafeed" "test" {
  name = "monitoring datafeed"
  source_id = ns1_datasource.test.id
  config = {
    jobid = ns1_monitoringjob.test.id
  }
}
resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
resource "ns1_record" "it" {
  zone   = ns1_zone.test.zone
  domain = "datafeed.${ns1_zone.test.zone}"
  type   = "A"
  answers {
    answer = "192.0.2.1"
    meta = {
      up = jsonencode({"feed": "${ns1_datafeed.test.id}"})
    }
  }
  filters {
    filter = "up"
  }
}
`, rString, rString)
}

func testAccRecordCAA(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "caa" {
  zone     = "${ns1_zone.test.zone}"
  domain   = "${ns1_zone.test.zone}"
  type     = "CAA"
  answers {
    answer = "0 issue letsencrypt.org"
  }
  answers {
    answer = "0 issuewild ;"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordCAAWithSpace(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "caa" {
        zone   = ns1_zone.test.zone
        domain = ns1_zone.test.zone
        type   = "CAA"

        answers {
                answer = "0 issue inbox2221.ticket; account=xyz"
        }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordSPF(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "spf" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "${ns1_zone.test.zone}"
  type              = "SPF"
  ttl               = 86400
  use_client_subnet = "true"
  answers {
    answer = "v=DKIM1; k=rsa; p=XXXXXXXX"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordSRV(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "srv" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "_some-server._tcp.${ns1_zone.test.zone}"
  type              = "SRV"
  ttl               = 86400
  use_client_subnet = "true"
  answers {
    answer = "10 0 2380 node-1.${ns1_zone.test.zone}"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordDS(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "ds" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "_some-server._tcp.${ns1_zone.test.zone}"
  type              = "DS"
  ttl               = 86400
  use_client_subnet = "true"
  answers {
    answer = "262 13 2 287787bd551bcab4f57d0c1dcaf312eebe36cc338bebb90d1402353c7096785d"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordURLFWD(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "urlfwd" {
  zone     = "${ns1_zone.test.zone}"
  domain   = "fwd.${ns1_zone.test.zone}"
  type     = "URLFWD"
  answers {
    answer = "/ https://example.com 301 2 0"
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

// zone and domain have leading and trailing dots and should fail validation.
func testAccRecordInvalid(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = ".terraform-test-%s.io."
  domain            = ".test.terraform-test-%s.io."
  type              = "CNAME"
  ttl               = 60
}
`, rString, rString)
}

func testAccRecordLink(rString string) string {
	return fmt.Sprintf(`
# the name being tested that is a link
resource "ns1_record" "link" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "link.${ns1_zone.test.zone}"
  type              = "A"
  link              = "target.${ns1_zone.test.zone}"
  ttl               = 666
}

# the record that is the destination of the link
resource "ns1_record" "target" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "target.${ns1_zone.test.zone}"
  type              = "A"
  ttl               = 777
  answers {
             answer = "99.86.99.86"
	  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordLinkUpdated(rString string) string {
	return fmt.Sprintf(`
# the name being tested that is a link
resource "ns1_record" "link" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "link.${ns1_zone.test.zone}"
  type              = "A"
  link              = "newtarget.${ns1_zone.test.zone}"
  ttl               = 666
}

# the record that is the original destination of the link
resource "ns1_record" "target" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "target.${ns1_zone.test.zone}"
  type              = "A"
  ttl               = 777
  answers {
             answer = "99.86.99.86"
	  }
}

# the new destination for the link
resource "ns1_record" "newtarget" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "newtarget.${ns1_zone.test.zone}"
  type              = "A"
  ttl               = 888
  answers {
             answer = "16.19.20.19"
          }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func TestRegionsMetaDiffSuppress(t *testing.T) {
	metaKeys := []string{"georegion", "country", "us_state", "ca_province"}

	for _, metaKey := range metaKeys {
		key := fmt.Sprintf("somepath.%s", metaKey)

		if metaDiffSuppress(key, "val1", "val2", nil) {
			t.Errorf("does not return that different strings are different (%s)", metaKey)
		}

		if !metaDiffSuppress(key, "val1", "val1", nil) {
			t.Errorf("does return that identical strings are different (%s)", metaKey)
		}

		if !metaDiffSuppress(key, "val1,val2", "val1,val2", nil) {
			t.Errorf("does return that identical strings with multiple elements are different (%s)", metaKey)
		}

		if !metaDiffSuppress(key, "val2,val1", "val1,val2", nil) {
			t.Errorf("does return that identical values with different orders are different (%s)", metaKey)
		}
	}

	if metaDiffSuppress("somepath.ignorekey", "val2,val1", "val1,val2", nil) {
		t.Errorf("is processing non-related meta keys")
	}
}

func TestValidateFQDN(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		expErrs int
	}{
		{
			"valid",
			"terraform-test-zone.io",
			0,
		},
		{
			"invalid - leading .",
			".terraform-test-zone.io",
			1,
		},
		{
			"invalid - trailing .",
			"terraform-test-zone.io.",
			1,
		},
		{
			"invalid - leading and trailing .",
			".terraform-test-zone.io.",
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWarns, outErrs := validateFQDN(tt.in, "")
			assert.Equal(t, tt.expErrs, len(outErrs))
			assert.Equal(t, 0, len(outWarns))
		})
	}
}
