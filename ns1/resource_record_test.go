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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordMeta(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, domainName),
					testAccCheckRecordAnswerMetaIPPrefixes(&record, []string{"1.1.1.1/32", "2.2.2.2/32"}),
				),
			},
			{
				ResourceName:      "ns1_record.it",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/A", zoneName, domainName),
				ImportStateVerify: true,
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
    filter = "up"
  }

  filters {
    filter = "geotarget_country"
  }

  filters {
    filter = "select_first_n"
    config = {N=1}
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
    filter = "up"
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

func testAccRecordMeta(rString string) string {
	return fmt.Sprintf(`
resource "ns1_record" "it" {
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "A"
  answers {
    answer = "1.2.3.4"

    meta = {
      weight = 5
      ip_prefixes = "1.1.1.1/32,2.2.2.2/32"
    }
  }
}

resource "ns1_zone" "test" {
  zone = "terraform-test-%s.io"
}
`, rString)
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

func TestRegionsMetaDiffSuppress(t *testing.T) {
	metaKeys := []string{"georegion", "country", "us_state", "ca_province"}

	for _, metaKey := range metaKeys {
		key := fmt.Sprintf("somepath.%s", metaKey)

		if regionsMetaDiffSuppress(key, "val1", "val2", nil) {
			t.Errorf("does not return that different strings are different (%s)", metaKey)
		}

		if !regionsMetaDiffSuppress(key, "val1", "val1", nil) {
			t.Errorf("does return that identical strings are different (%s)", metaKey)
		}

		if !regionsMetaDiffSuppress(key, "val1,val2", "val1,val2", nil) {
			t.Errorf("does return that identical strings with multiple elements are different (%s)", metaKey)
		}

		if !regionsMetaDiffSuppress(key, "val2,val1", "val1,val2", nil) {
			t.Errorf("does return that identical values with different orders are different (%s)", metaKey)
		}
	}

	if regionsMetaDiffSuppress("somepath.ignorekey", "val2,val1", "val1,val2", nil) {
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
