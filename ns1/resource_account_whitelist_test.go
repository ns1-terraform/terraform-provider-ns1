package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestAccAccountWhitelist_basic(t *testing.T) {
	var wl account.IPWhitelist
	var name = fmt.Sprintf("it-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAccountWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountWhitelistBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountWhitelistExists(fmt.Sprintf("ns1_account_whitelist.%s", name), &wl),
				),
			},
		},
	})
}

func TestAccAccountWhitelist_ManualDelete(t *testing.T) {
	var wl account.IPWhitelist
	var name = fmt.Sprintf("it-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAccountWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountWhitelistBasic(name),
				Check:  testAccCheckAccountWhitelistExists(fmt.Sprintf("ns1_account_whitelist.%s", name), &wl),
			},
			// Simulate a manual deletion of the whitelist and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteAccountWhitelist(t, &wl),
				Config:             testAccAccountWhitelistBasic(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccAccountWhitelistBasic(name),
				Check:  testAccCheckAccountWhitelistExists(fmt.Sprintf("ns1_account_whitelist.%s", name), &wl),
			},
		},
	})
}

func TestAccAccountWhitelist_updated(t *testing.T) {
	var wl account.IPWhitelist
	var name = fmt.Sprintf("it-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAccountWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountWhitelistBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountWhitelistExists(fmt.Sprintf("ns1_account_whitelist.%s", name), &wl),
					testAccCheckAccountWhitelistValues(&wl, 5),
				),
			},
			{
				Config: testAccAccountWhitelistUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountWhitelistExists(fmt.Sprintf("ns1_account_whitelist.%s", name), &wl),
					testAccCheckAccountWhitelistValues(&wl, 1),
				),
			},
		},
	})
}

// Other stuff passed here

func testAccCheckAccountWhitelistExists(n string, accountWhitelist *account.IPWhitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundAccountwhitelist, _, err := client.GlobalIPWhitelist.Get(rs.Primary.Attributes["id"])

		p := rs.Primary

		if err != nil {
			return err
		}

		if foundAccountwhitelist.ID != p.Attributes["id"] {
			return fmt.Errorf("Global allow list not found")
		}

		*accountWhitelist = *foundAccountwhitelist

		return nil
	}
}

func testAccManualDeleteAccountWhitelist(t *testing.T, wl *account.IPWhitelist) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.GlobalIPWhitelist.Delete(wl.ID)

		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			t.Logf("failed to delete global white list: %v", err)
		}
	}
}

func testAccAccountWhitelistDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_account_whitelist" {
			continue
		}

		wl, _, err := client.GlobalIPWhitelist.Get(rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("global ip allow list still exists %#v: %#v", err, wl)
		}
	}

	return nil
}

func testAccCheckAccountWhitelistValues(wl *account.IPWhitelist, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		values := len(wl.Values)

		if values != expected {
			return fmt.Errorf("IPWhitelist.Values: got: %v want: %v", values, expected)
		}

		return nil
	}
}

func testAccAccountWhitelistBasic(name string) string {
	return fmt.Sprintf(`resource "ns1_account_whitelist" "%s" {
		name = "%s"
		values			= [
			"0.0.0.0/0",
			"192.0.0.0/8",
			"192.168.0.0/16",
			"192.168.1.0/24",
			"192.168.1.1"
			]
	  }`, name, name)
}

func testAccAccountWhitelistUpdated(name string) string {
	return fmt.Sprintf(`resource "ns1_account_whitelist" "%s" {
  		name = "%s"
		values			= [
			"0.0.0.0/0",
			]
	}`, name, name)
}
