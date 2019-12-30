package ns1

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestAccAPIKey_basic(t *testing.T) {
	var apiKey account.APIKey
	name := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
				),
			},
		},
	})
}

func TestAccAPIKey_updated(t *testing.T) {
	var apiKey account.APIKey
	name := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	updatedName := fmt.Sprintf("%s-updated", name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
				),
			},
			{
				Config: testAccAPIKeyBasic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, updatedName),
				),
			},
		},
	})
}

func TestAccAPIKey_ManualDelete(t *testing.T) {
	var apiKey account.APIKey
	name := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyBasic(name),
				Check:  testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
			},
			// Simulate a manual deletion of the API key and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteAPIKey(&apiKey),
				Config:             testAccAPIKeyBasic(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccAPIKeyBasic(name),
				Check:  testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
			},
		},
	})
}

func TestAccAPIKey_permissions(t *testing.T) {
	var apiKey account.APIKey
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test key %s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
				),
			},
			{
				Config: testAccAPIKeyPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
				),
			},
			{
				Config: testAccAPIKeyPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
					// The key should still have this permission, it would have inherited it from the team.
					resource.TestCheckResourceAttr("ns1_apikey.it", "account_manage_account_settings", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAPIKeyPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIKeyExists("ns1_apikey.it", &apiKey),
					testAccCheckAPIKeyName(&apiKey, name),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_apikey.it", "account_manage_account_settings", "false"),
				),
			},
		},
	})
}

func testAccCheckAPIKeyExists(n string, apiKey *account.APIKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundAPIKey, _, err := client.APIKeys.Get(rs.Primary.Attributes["id"])

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundAPIKey.ID != p.Attributes["id"] {
			return fmt.Errorf("API key not found")
		}

		*apiKey = *foundAPIKey

		return nil
	}
}

func testAccCheckAPIKeyName(apiKey *account.APIKey, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if apiKey.Name != expected {
			return fmt.Errorf("apiKey: got: %s want: %s", apiKey.Name, expected)
		}

		return nil
	}
}

func testAccCheckAPIKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	var apiKey string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ns1_apikey" {
			apiKey = rs.Primary.Attributes["id"]
		}
	}

	key, _, _ := client.APIKeys.Get(apiKey)
	if key != nil {
		return fmt.Errorf("apiKey still exists: %#v", key)
	}

	return nil
}

// Simulate a manual deletion of an API key.
func testAccManualDeleteAPIKey(apiKey *account.APIKey) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.APIKeys.Delete(apiKey.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete api key: %v", err)
		}
	}
}

func testAccAPIKeyBasic(apiKeyName string) string {
	return fmt.Sprintf(`resource "ns1_apikey" "it" {
  name = "%s"
}
`, apiKeyName)
}

func testAccAPIKeyPermissionsOnTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_apikey" "it" {
  name = "terraform acc test key %s"

  teams = ["${ns1_team.t.id}"]
}

`, rString, rString)
}

func testAccAPIKeyPermissionsNoTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_apikey" "it" {
  name = "terraform acc test key %s"
}

`, rString, rString)
}
