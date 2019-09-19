package ns1

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestAccRecord_basic(t *testing.T) {
	var record dns.Record
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, "test.terraform-record-test.io"),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 10),
					testAccCheckRecordAnswerRdata(&record, 0, "test1.terraform-record-test.io"),
				),
			},
		},
	})
}

func TestAccRecord_updated(t *testing.T) {
	var record dns.Record
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, "test.terraform-record-test.io"),
					testAccCheckRecordTTL(&record, 60),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordRegionName(&record, []string{"cal"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 10),
					testAccCheckRecordAnswerRdata(&record, 0, "test1.terraform-record-test.io"),
				),
			},
			{
				Config: testAccRecordUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.it", &record),
					testAccCheckRecordDomain(&record, "test.terraform-record-test.io"),
					testAccCheckRecordTTL(&record, 120),
					testAccCheckRecordUseClientSubnet(&record, false),
					testAccCheckRecordRegionName(&record, []string{"ny", "wa"}),
					// testAccCheckRecordAnswerMetaWeight(&record, 5),
					testAccCheckRecordAnswerRdata(&record, 0, "test2.terraform-record-test.io"),
				),
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
		},
	})
}

func TestAccRecord_SPF(t *testing.T) {
	var record dns.Record
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordSPF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.spf", &record),
					testAccCheckRecordDomain(&record, "terraform-record-test.io"),
					testAccCheckRecordTTL(&record, 86400),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(&record, 0, "v=DKIM1; k=rsa; p=XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRecord_SRV(t *testing.T) {
	var record dns.Record
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRecordSRV,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("ns1_record.srv", &record),
					testAccCheckRecordDomain(&record, "_some-server._tcp.terraform-record-test.io"),
					testAccCheckRecordTTL(&record, 86400),
					testAccCheckRecordUseClientSubnet(&record, true),
					testAccCheckRecordAnswerRdata(&record, 0, "10"),
					testAccCheckRecordAnswerRdata(&record, 1, "0"),
					testAccCheckRecordAnswerRdata(&record, 2, "2380"),
					testAccCheckRecordAnswerRdata(&record, 3, "node-1.terraform-record-test.io"),
				),
			},
		},
	})
}

func testAccCheckRecordExists(n string, record *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %v", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NoID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		p := rs.Primary

		foundRecord, _, err := client.Records.Get(p.Attributes["zone"], p.Attributes["domain"], p.Attributes["type"])
		if err != nil {
			return fmt.Errorf("Record not found")
		}

		if foundRecord.Domain != p.Attributes["domain"] {
			return fmt.Errorf("Record not found")
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
		return fmt.Errorf("Record still exists: %#v %#v", foundRecord, err)
	}

	return nil
}

func testAccCheckRecordDomain(r *dns.Record, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if r.Domain != expected {
			return fmt.Errorf("Domain: got: %#v want: %#v", r.Domain, expected)
		}
		return nil
	}
}

func testAccCheckRecordTTL(r *dns.Record, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if r.TTL != expected {
			return fmt.Errorf("TTL: got: %#v want: %#v", r.TTL, expected)
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
			return fmt.Errorf("Regions: got: %#v want: %#v", regions, expected)
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
			return fmt.Errorf("Answers[0].Meta.Weight: got: %#v want: %#v", weight, expected)
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

func testAccCheckRecordAnswerRdata(r *dns.Record, idx int, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		recordAnswer := r.Answers[0]
		recordAnswerString := recordAnswer.Rdata[idx]
		if recordAnswerString != expected {
			return fmt.Errorf("Answers[0].Rdata[%d]: got: %#v want: %#v", idx, recordAnswerString, expected)
		}
		return nil
	}
}

const testAccRecordBasic = `
resource "ns1_record" "it" {
  zone              = ns1_zone.test.zone
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"
  ttl               = 60

  // meta {
  //   weight = 5
  //   connections = 3
  //   // up = false // Ignored by d.GetOk("meta.0.up") due to known issue
  // }

  answers {
    answer = "test1.terraform-record-test.io"
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
  zone = "terraform-record-test.io"
}
`

const testAccRecordUpdated = `
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
    answer = "test2.terraform-record-test.io"
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
  zone = "terraform-record-test.io"
}
`

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

const testAccRecordSPF = `
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
  zone = "terraform-record-test.io"
}
`

const testAccRecordSRV = `
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
  zone = "terraform-record-test.io"
}
`

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
