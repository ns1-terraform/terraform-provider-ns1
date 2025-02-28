package ns1

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestAccRecord_remove_all_filters(t *testing.T) {
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
					testAccCheckRecordFilterCount(&record, 3),
				),
			},
			{
				Config: testAccRecordUpdatedNoFilters(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordFilterCount(&record, 0),
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

func TestAccRecord_NewTypes(t *testing.T) {
	testCases := []struct {
		recType         string
		domainPrefix    string
		configFuncs     []func(string) string
		preCheck        func(*testing.T)
		expectedAnswers [][]string
		expectedDomains []func(zoneName string) string
		ttl             int
	}{
		{
			recType: "CAA",
			configFuncs: []func(string) string{
				testAccRecordCAA,
				testAccRecordCAAWithSpace,
			},
			// Update these expected answers to include the double quotes
			expectedAnswers: [][]string{
				{"0", "issue", "\"letsencrypt.org\""},
				{"0", "issuewild", "\";\""},
			},
			expectedDomains: []func(string) string{
				func(zone string) string { return fmt.Sprintf("caa.%s", zone) },
				func(zone string) string { return zone },
			},
			ttl: 3600,
		},
		{
			recType:      "APL",
			domainPrefix: "apl",
			configFuncs:  []func(string) string{testAccRecordAPL},
			expectedAnswers: [][]string{
				{"1:127.0.0.0/16"},
			},

			ttl: 3600,
		},
		{
			recType: "SPF",
			// Only 1 config step
			configFuncs: []func(string) string{
				testAccRecordSPF,
			},
			expectedAnswers: [][]string{
				// Just one answer from old test
				{"v=DKIM1; k=rsa; p=XXXXXXXX"},
			},
			expectedDomains: []func(string) string{
				func(zone string) string { return zone },
			},
			ttl: 86400,
		},
		{
			recType:      "SRV",
			domainPrefix: "_some-server._tcp",
			configFuncs:  []func(string) string{testAccRecordSRV},
			expectedAnswers: [][]string{
				{"10", "0", "2380", "node-1.${ns1_zone.test.zone}"},
			},
			ttl: 86400,
		},
		{
			recType:      "DS",
			domainPrefix: "_some-server._tcp",
			configFuncs:  []func(string) string{testAccRecordDS},
			expectedAnswers: [][]string{
				{"262", "13", "2", "287787bd551bcab4f57d0c1dcaf312eebe36cc338bebb90d1402353c7096785d"},
			},
			ttl: 86400,
		},
		{
			recType:      "GPOS",
			domainPrefix: "gpos",
			configFuncs:  []func(string) string{testAccRecordGPOS},
			expectedAnswers: [][]string{
				{"6", "53", "0"},
			},
			ttl: 3600,
		},
		{
			recType:      "IPSECKEY",
			domainPrefix: "ipseckey",
			configFuncs:  []func(string) string{testAccRecordIPSECKEY},
			expectedAnswers: [][]string{
				{"1", "0", "2", ".", "abba"},
			},
			ttl: 3600,
		},
		{
			recType:      "OPENPGPKEY",
			domainPrefix: "openpgpkey",
			configFuncs:  []func(string) string{testAccRecordOPENPGPKEY},
			expectedAnswers: [][]string{
				{"abba"},
			},
			ttl: 3600,
		},
		{
			recType:      "SSHFP",
			domainPrefix: "sshfp",
			configFuncs:  []func(string) string{testAccRecordSSHFP},
			expectedAnswers: [][]string{
				{"1", "1", "abba"},
			},
			ttl: 3600,
		},
		{
			recType:      "URI",
			domainPrefix: "uri",
			configFuncs:  []func(string) string{testAccRecordURI},
			expectedAnswers: [][]string{
				{"1", "2", "\"http://localhost\""},
			},
			ttl: 3600,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.recType, func(t *testing.T) {
			var record dns.Record
			rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
			zoneName := fmt.Sprintf("terraform-test-%s.io", rString)

			var defaultDomain string

			if tc.recType == "SRV" {
				defaultDomain = fmt.Sprintf("_some-server._tcp.%s", zoneName)
				expectedAnswer := []string{"10", "0", "2380", fmt.Sprintf("node-1.%s", zoneName)}
				tc.expectedAnswers = [][]string{expectedAnswer}
			} else if tc.recType == "DS" {
				defaultDomain = fmt.Sprintf("_some-server._tcp.%s", zoneName)
			} else if tc.domainPrefix != "" {
				defaultDomain = fmt.Sprintf("%s.%s", tc.domainPrefix, zoneName)
			} else {
				defaultDomain = zoneName
			}

			if tc.preCheck != nil {
				tc.preCheck(t)
			}

			if tc.recType == "CAA" {
				for i, configFunc := range tc.configFuncs {
					var expectedDomain string
					if len(tc.expectedDomains) > i {
						expectedDomain = tc.expectedDomains[i](zoneName)
					} else {
						expectedDomain = defaultDomain
					}

					var checkFuncs []resource.TestCheckFunc

					checkFuncs = append(checkFuncs,
						testAccCheckRecordExists(fmt.Sprintf("ns1_record.%s", strings.ToLower(tc.recType)), &record),
						testAccCheckRecordDomain(&record, expectedDomain),
						testAccCheckRecordTTL(&record, tc.ttl),
						testAccCheckRecordUseClientSubnet(&record, true),
					)

					if i == 0 {
						checkFuncs = append(checkFuncs,
							testAccCheckRecordAnswerRdata(t, &record, 0, []string{"0", "issue", "\"letsencrypt.org\""}),
							testAccCheckRecordAnswerRdata(t, &record, 1, []string{"0", "issuewild", "\";\""}),
						)
					} else if i == 1 {
						checkFuncs = append(checkFuncs,
							testAccCheckRecordAnswerRdata(t, &record, 0, []string{"0", "issue", "inbox2221.ticket; account=xyz"}),
						)
					}

					resource.Test(t, resource.TestCase{
						PreCheck:     func() { testAccPreCheck(t) },
						Providers:    testAccProviders,
						CheckDestroy: testAccCheckRecordDestroy,
						Steps: []resource.TestStep{
							{
								Config: configFunc(rString),
								Check:  resource.ComposeTestCheckFunc(checkFuncs...),
							},
							{
								ResourceName:      fmt.Sprintf("ns1_record.%s", strings.ToLower(tc.recType)),
								ImportState:       true,
								ImportStateId:     fmt.Sprintf("%s/%s/%s", zoneName, expectedDomain, tc.recType),
								ImportStateVerify: true,
							},
						},
					})
				}
			} else {
				// Standard handling for all other record types
				for i, configFunc := range tc.configFuncs {
					var expectedDomain string
					if len(tc.expectedDomains) > i {
						expectedDomain = tc.expectedDomains[i](zoneName)
					} else {
						expectedDomain = defaultDomain
					}

					checks := []resource.TestCheckFunc{
						testAccCheckRecordExists(fmt.Sprintf("ns1_record.%s", strings.ToLower(tc.recType)), &record),
						testAccCheckRecordDomain(&record, expectedDomain),
						testAccCheckRecordTTL(&record, tc.ttl),
						testAccCheckRecordUseClientSubnet(&record, true),
					}

					answerChecksSlice := createAnswerChecks(t, &record, tc.expectedAnswers)
					allCheckFuncs := append(checks, answerChecksSlice...)

					resource.Test(t, resource.TestCase{
						PreCheck:     func() { testAccPreCheck(t) },
						Providers:    testAccProviders,
						CheckDestroy: testAccCheckRecordDestroy,
						Steps: []resource.TestStep{
							{
								Config: configFunc(rString),
								Check:  resource.ComposeTestCheckFunc(allCheckFuncs...),
							},
							{
								ResourceName:      fmt.Sprintf("ns1_record.%s", strings.ToLower(tc.recType)),
								ImportState:       true,
								ImportStateId:     fmt.Sprintf("%s/%s/%s", zoneName, expectedDomain, tc.recType),
								ImportStateVerify: true,
							},
						},
					})
				}
			}
		})
	}
}

func createAnswerChecks(t *testing.T, record *dns.Record, answers [][]string) []resource.TestCheckFunc {

	checks := make([]resource.TestCheckFunc, len(answers))
	for i, answer := range answers {
		checks[i] = testAccCheckRecordAnswerRdata(t, record, i, answer)
	}
	return checks
}

func TestAccRecord_WithTags(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	zoneName := fmt.Sprintf("terraform-test-%s.io", rString)
	domainName := fmt.Sprintf("tagged.%s", zoneName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordWithTags(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.tagged", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTagData(
						map[string]string{"tag1": "location1", "tag2": "location2"},
						&record,
					),
				),
			},
			{
				ResourceName:      "ns1_record.tagged",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/A", zoneName, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRecord_updatedWithTags(t *testing.T) {
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
				Config: testAccRecordUpdatedWithTags(rString),
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
					testAccCheckRecordTagData(
						map[string]string{"tag1": "location1", "tag2": "location2"},
						&record,
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

func TestAccRecord_updatedWithRegions(t *testing.T) {
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
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				Config: testAccRecordUpdatedWithRegionWeight(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				Config: testAccRecordUpdatedWithRegions(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"ny", "cal"}),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				Config: testAccRecordUpdatedWithRegionWeight(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					testAccCheckRecordAnswerRdata(
						t, &record, 0, []string{fmt.Sprintf("test1.%s", zoneName)},
					),
				),
			},
			{
				Config: testAccRecordUpdatedWithNoRegions(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{}),
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

func TestAccRecord_validationError(t *testing.T) {
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordInvalid(rString),
				/* The error block has lines like this:
				   Error: zone has an invalid leading ".", got: .terraform-test-vmatw2m3iunjw4j.io.
				   Error: zone has an invalid trailing ".", got: .terraform-test-vmatw2m3iunjw4j.io.
				   Error: domain has an invalid leading ".", got: .test.terraform-test-vmatw2m3iunjw4j.io.
				   Error: domain has an invalid trailing ".", got: .test.terraform-test-vmatw2m3iunjw4j.io.
				*/
				ExpectError: regexp.MustCompile(`(?s)(Error: (zone|domain) has an invalid (leading|trailing) "\.", got: .*){4}`),
			},
			{
				Config:      testAccRecordNoAnswers(rString),
				ExpectError: regexp.MustCompile(`Invalid body`),
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

func TestAccRecord_OverrideAddressRecordsNilToTrue(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileBasicAlias := testAccRecordBasicALIAS(rString)
	tfFileOverrideAddressRecordsAliasTrue := testAccRecordBasicALIASOverrideAddressRecords(rString, true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_address_records not set
			{
				Config: tfFileBasicAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, false),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileBasicAlias,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Change override_address_records to true Plan and Apply
			{
				Config:             tfFileOverrideAddressRecordsAliasTrue,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: tfFileOverrideAddressRecordsAliasTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, true),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideAddressRecordsAliasTrue,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRecord_OverrideAddressRecordsTrueToNil(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileBasicAlias := testAccRecordBasicALIAS(rString)
	tfFileOverrideAddressRecordsAliasTrue := testAccRecordBasicALIASOverrideAddressRecords(rString, true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_address_records true
			{
				Config: tfFileOverrideAddressRecordsAliasTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, true),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideAddressRecordsAliasTrue,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:             tfFileBasicAlias,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Change override_address_records from true setting to "null" results in false
			{
				Config: tfFileBasicAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, false),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileBasicAlias,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRecord_OverrideAddressRecordsTrueToFalse(t *testing.T) {
	var record dns.Record
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	tfFileOverrideAddressRecordsAliasTrue := testAccRecordBasicALIASOverrideAddressRecords(rString, true)
	tfFileDoNotOverrideAddressRecordsAliasFalse := testAccRecordBasicALIASOverrideAddressRecords(rString, false)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			// Create an ALIAS record with override_address_records true
			{
				Config: tfFileOverrideAddressRecordsAliasTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, true),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileOverrideAddressRecordsAliasTrue,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Change override_address_records to false setting to false Plan and Apply
			{
				Config:             tfFileDoNotOverrideAddressRecordsAliasFalse,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: tfFileDoNotOverrideAddressRecordsAliasFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordOverrideAddressRecords(&record, false),
				),
			},
			// Plan again to detect "loop" conditions
			{
				Config:             tfFileDoNotOverrideAddressRecordsAliasFalse,
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

func testAccCheckRecordFilterCount(r *dns.Record, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(r.Filters) != expected {
			return fmt.Errorf("r.Filters, got: %d, want: %d", len(r.Filters), expected)
		}

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

func testAccCheckRecordDomain(record *dns.Record, expectedDomain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if record.Domain != expectedDomain {
			return fmt.Errorf("r.Domain: got: %q, want: %q", record.Domain, expectedDomain)
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

func testAccCheckRecordOverrideAddressRecords(r *dns.Record, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *r.OverrideAddressRecords != expected {
			return fmt.Errorf("Override Address Records: got: %#v want: %#v", *r.OverrideAddressRecords, expected)
		}

		return nil
	}
}

func testAccCheckRecordOverrideTTL(r *dns.Record, expectedNil bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if expectedNil {
			if r.OverrideTTL != nil {
				return fmt.Errorf("Override TTL: got: %#v want: null", *r.OverrideTTL)
			}
			return nil
		}
		if r.OverrideTTL == nil {
			return fmt.Errorf("Override TTL: got: %v want: notNil", r.OverrideTTL)
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

func testAccCheckRecordTagData(expected interface{}, r *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordTags := r.Tags
		if !reflect.DeepEqual(recordTags, expected) {
			return fmt.Errorf("tags: got: %#v want: %#v", recordTags, expected)
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

func testAccRecordBasicALIASOverrideAddressRecords(rString string, overrideAddressRecords bool) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              			= "${ns1_zone.test.zone}"
  domain            			= "${ns1_zone.test.zone}"
  type              			= "ALIAS"
  ttl               			= 60
  override_address_records 		= %v
  answers {
    answer = "test.${ns1_zone.test.zone}"
  }
}

resource "ns1_zone" "test" {
	zone = "terraform-test-%s.io"
}
`, overrideAddressRecords, rString)
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

func testAccRecordUpdatedNoFilters(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
	zone              = "${ns1_zone.test.zone}"
	domain            = "test.${ns1_zone.test.zone}"
	type              = "CNAME"
	ttl               = 60

	answers {
		answer = "test1.${ns1_zone.test.zone}"
		region = "cal"
	}

	regions {
		name = "cal"
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
    "nrt"
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
resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}

resource "ns1_record" "caa" {
  zone   = ns1_zone.test.zone
  domain = "caa.${ns1_zone.test.zone}"
  type   = "CAA"
  ttl    = 3600
  
  answers {
    answer = "0 issue \"letsencrypt.org\""
  }
  
  answers {
    answer = "0 issuewild \";\""
  }
}
`, rString)
}

func testAccRecordAPL(rString string) string {
	zone := fmt.Sprintf("terraform-test-%s.io", rString)
	return fmt.Sprintf(`
resource "ns1_zone" "test" {
	zone = "%s"
}

resource "ns1_record" "apl" {
	zone   = ns1_zone.test.zone
	domain = "apl.%s"
	type   = "APL"
	ttl    = 3600

	answers {
		answer = "1:127.0.0.0/16"
	}
	
}
`, zone, zone)
}

func testAccRecordGPOS(rString string) string {
	return fmt.Sprintf(`
resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}

resource "ns1_record" "gpos" {
  zone   = ns1_zone.test.zone
  domain = "gpos.${ns1_zone.test.zone}"
  type   = "GPOS"
  ttl    = 3600

  answers {
    answer = "6 53 0"
  }
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

func testAccRecordIPSECKEY(rString string) string {
	return fmt.Sprintf(`
resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}

resource "ns1_record" "ipseckey" {
  zone   = ns1_zone.test.zone
  domain = "ipseckey.${ns1_zone.test.zone}"
  type   = "IPSECKEY"
  ttl    = 3600

  answers {
    answer = "1 0 2 . abba"
  }
}
`, rString)
}

func testAccRecordWithTags(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "tagged" {
  zone     = "${ns1_zone.test.zone}"
  domain   = "tagged.${ns1_zone.test.zone}"
  type     = "A"
  answers {
    answer = "1.2.3.4"
  }
  tags = {tag1 = "location1", tag2 = "location2"}
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordUpdatedWithTags(rString string) string {
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

  tags = {tag1: "location1", tag2: "location2"}
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
}

func testAccRecordUpdatedWithRegionWeight(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"
  ttl               = 60

  answers {
    answer = "test1.${ns1_zone.test.zone}"
    region = "cal"
  }

  regions {
    name = "cal"
    meta = {
      weight = 100
    }
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

func testAccRecordUpdatedWithRegions(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"
  ttl               = 60

  answers {
    answer = "test1.${ns1_zone.test.zone}"
    region = "cal"
  }

  regions {
    name = "cal"
    meta = {
      weight = 90
    }
  }

  regions {
    name = "ny"
    meta = {
      weight = 10
    }
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

func testAccRecordUpdatedWithNoRegions(rString string) string {
	return fmt.Sprintf(`
	resource "ns1_record" "it" {
	  zone              = "${ns1_zone.test.zone}"
	  domain            = "test.${ns1_zone.test.zone}"
	  type              = "CNAME"
	  ttl               = 60

	  answers {
	    answer = "test1.${ns1_zone.test.zone}"
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

func testAccRecordOPENPGPKEY(rString string) string {
	return fmt.Sprintf(`
	resource "ns1_zone" "test" {
	  zone = "terraform-test-%s.io"
	}
	
	resource "ns1_record" "openpgpkey" {
	  zone   = ns1_zone.test.zone
	  domain = "openpgpkey.${ns1_zone.test.zone}"
	  type   = "OPENPGPKEY"
	  ttl    = 3600
	
	  answers {
		answer = "abba"
	  }
	}
	`, rString)
}

func testAccRecordSSHFP(rString string) string {
	zone := fmt.Sprintf("terraform-test-%s.io", rString)
	return fmt.Sprintf(`
resource "ns1_zone" "test" {
  zone = "%s"
}

resource "ns1_record" "sshfp" {
  zone   = ns1_zone.test.zone
  domain = "sshfp.%s"
  type   = "SSHFP"
  ttl    = 3600

  answers {
    answer = "1 1 abba"
	}
}
`, zone, zone)
}

func testAccRecordURI(rString string) string {
	zone := fmt.Sprintf("terraform-test-%s.io", rString)
	return fmt.Sprintf(`
resource "ns1_zone" "test" {
  zone = "%s"
}

resource "ns1_record" "uri" {
  zone   = ns1_zone.test.zone
  domain = "uri.%s"
  type   = "URI"
  ttl    = 3600

  answers {
    answer = "1 2 \"http://localhost\""
	}
}
`, zone, zone)
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

// there must be at least one answer
func testAccRecordNoAnswers(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "terraform-test-%s.io"
  domain            = "test.terraform-test-%s.io"
  type              = "CNAME"
  ttl               = 60
  answers {}
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
