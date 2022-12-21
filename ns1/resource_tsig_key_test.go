package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// Creating TSIG Key
func TestAccTsigKey_basic(t *testing.T) {
	var (
		key          = dns.TSIGKey{}
		keyName      = fmt.Sprintf("terraform-test-%s.", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		keyAlgorithm = "hmac-sha256"
		keySecret    = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
	)
	// Basic test for TSIG Key
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsigKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsigKeyBasic(keyName, keyAlgorithm, keySecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsigKeyExists("ns1_tsigkey.it", &key),
					testAccCheckTsigKeyName(&key, keyName),
					testAccCheckTsigKeyAlgorithm(&key, keyAlgorithm),
					testAccCheckTsigKeySecret(&key, keySecret),
				),
			},
			{
				ResourceName:      "ns1_tsigkey.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Updating TSIG Keys
func TestAccTsigKey_updated(t *testing.T) {
	var (
		key          = dns.TSIGKey{}
		keyName      = fmt.Sprintf("terraform-test-%s.", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		keyAlgorithm = "hmac-sha256"
		keySecret    = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="

		updatedAlgorithm = "hmac-sha256"
		updatedSecret    = "Mo1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
	)

	// Updating TSIG Key
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsigKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsigKeyBasic(keyName, keyAlgorithm, keySecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsigKeyExists("ns1_tsigkey.it", &key),
					testAccCheckTsigKeyName(&key, keyName),
					testAccCheckTsigKeyAlgorithm(&key, keyAlgorithm),
					testAccCheckTsigKeySecret(&key, keySecret),
				),
			},
			{
				ResourceName:      "ns1_tsigkey.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTsigKeyBasic(keyName, updatedAlgorithm, updatedSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsigKeyExists("ns1_tsigkey.it", &key),
					testAccCheckTsigKeyAlgorithm(&key, updatedAlgorithm),
					testAccCheckTsigKeySecret(&key, updatedSecret),
				),
			},
			{
				ResourceName:      "ns1_tsigkey.it",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Manually deleting TSIG Key
func TestAccTsigKey_ManualDelete(t *testing.T) {
	var (
		key          = dns.TSIGKey{}
		keyName      = fmt.Sprintf("terraform-test-%s.", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))
		keyAlgorithm = "hmac-sha256"
		keySecret    = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsigKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsigKeyBasic(keyName, keyAlgorithm, keySecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsigKeyExists("ns1_tsigkey.it", &key),
					testAccCheckTsigKeyName(&key, keyName),
					testAccCheckTsigKeyAlgorithm(&key, keyAlgorithm),
					testAccCheckTsigKeySecret(&key, keySecret),
				),
			},
			// Simulate a manual deletion of the TSIG key and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteTsigKey(&key),
				Config:             testAccTsigKeyBasic(keyName, keyAlgorithm, keySecret),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccTsigKeyBasic(keyName, keyAlgorithm, keySecret),
				Check:  testAccCheckTsigKeyExists("ns1_tsigkey.it", &key),
			},
		},
	})
}

func testAccTsigKeyBasic(keyName string, keyAlgorithm string, keySecret string) string {
	return fmt.Sprintf(`resource "ns1_tsigkey" "it" {
  		name = "%s"
		algorithm = "%s"
		secret = "%s"
}
`, keyName, keyAlgorithm, keySecret)
}

func testAccCheckTsigKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_tsigkey" {
			continue
		}

		tsigKey, _, err := client.TSIG.Get(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("TSIG key still exists: %#v: %#v", err, tsigKey)
		}

	}

	return nil
}

func testAccCheckTsigKeyExists(n string, key *dns.TSIGKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundTsigKey, _, err := client.TSIG.Get(rs.Primary.ID)

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundTsigKey.Name != p.Attributes["id"] {
			return fmt.Errorf("TSIG key not found")
		}

		*key = *foundTsigKey

		return nil
	}
}

func testAccCheckTsigKeyName(key *dns.TSIGKey, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key.Name != expected {
			return fmt.Errorf("key.Name: got: %s want: %s", key.Name, expected)
		}
		return nil
	}
}

func testAccCheckTsigKeyAlgorithm(key *dns.TSIGKey, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key.Algorithm != expected {
			return fmt.Errorf("key.Algorithm: got: %s want: %s", key.Algorithm, expected)
		}
		return nil
	}
}

func testAccCheckTsigKeySecret(key *dns.TSIGKey, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key.Secret != expected {
			return fmt.Errorf("key.Secret: got: %s want: %s", key.Secret, expected)
		}
		return nil
	}
}

func testAccManualDeleteTsigKey(key *dns.TSIGKey) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.TSIG.Delete(key.Name)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete pulsar job: %v", err)
		}
	}
}
