package ns1

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestAccUser_basic(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "teams.#", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "notify.%", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					testAccCheckUserIPWhitelists(&user, []string{"1.1.1.1", "2.2.2.2"}),
					resource.TestCheckResourceAttr("ns1_user.u", "ip_whitelist_strict", "true"),
				),
			},
			{
				Config: testAccUserUpdated(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "teams.#", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "notify.%", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					testAccCheckUserIPWhitelists(&user, []string{}),
					resource.TestCheckResourceAttr("ns1_user.u", "ip_whitelist_strict", "false"),
				),
			},
		},
	})
}

func TestAccUser_ManualDelete(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasic(rString),
				Check:  testAccCheckUserExists("ns1_user.u", &user),
			},
			// Simulate a manual deletion of the user and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteUser(username),
				Config:             testAccUserBasic(rString),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUser_permissions(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
			{
				Config: testAccUserSecurityPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "security_manage_global_2fa", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "security_manage_active_directory", "true"),
				),
			},
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// The user should still have this permission, it would have inherited it from the team.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
		},
	})
}

func TestAccUser_permissions_empty_team(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			// Strange Terraform behavior causes explicitly settings a users team to []
			// to behave differently than removing the block entirely, so test for this as well.
			{
				Config: testAccUserPermissionsEmptyTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
		},
	})
}

// Edge cases exist with starting a user on a team vs. on no team, so test for this as well.
func TestAccUser_permissions_start_no_team(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// The user should still have this permission, it would have inherited it from the team.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
		},
	})
}

// Case when a user starts on a single team and is added to another team.
func TestAccUser_permissions_multiple_teams(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			{
				Config: testAccUserPermissionsOnTwoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			{
				Config: testAccUserPermissionsEmptyTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// The user should still have this permission, it would have inherited it from the team.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
		},
	})
}

// Case when a user starts on no teams and is added to multiple teams at once.
func TestAccUser_permissions_multiple_teams_start_no_team(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
			{
				Config: testAccUserPermissionsOnTwoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
			},
			{
				Config: testAccUserPermissionsEmptyTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// The user should still have this permission, it would have inherited it from the team.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccUserPermissionsNoTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					// But if an apply is ran again, the permission will be removed.
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "false"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_ip_whitelist", "true"),
				),
			},
		},
	})
}

// Import user test
func TestAccUser_import_test(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)
	ignored_fields := []string{""}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "teams.#", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "notify.%", "1"),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					testAccCheckUserIPWhitelists(&user, []string{"1.1.1.1", "2.2.2.2"}),
					resource.TestCheckResourceAttr("ns1_user.u", "ip_whitelist_strict", "true"),
				),
			},
			{
				ResourceName:      "ns1_user.u",
				ImportState:       true,
				ImportStateId:     username,
				ImportStateVerify: true,
				// Ignoring some fields because of how the permissions work right now
				ImportStateVerifyIgnore: ignored_fields,
			},
		},
	})
}

func TestAccUser_emptyIPWhitelist(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserEmptyIPWhitelist(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "ip_whitelist.#", "0"),
					resource.TestCheckResourceAttr("ns1_user.u", "ip_whitelist_strict", "false"),
					testAccCheckUserIPWhitelists(&user, []string{}),
				),
			},
		},
	})
}

func testAccUserEmptyIPWhitelist(rString string) string {
	return fmt.Sprintf(`
resource "ns1_team" "test" {
  name = "terraform test team %s"
}

resource "ns1_user" "u" {
  name     = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email    = "tf_acc_test_ns1@hashicorp.com"
  teams    = [ns1_team.test.id]
  
  ip_whitelist_strict = false
  ip_whitelist = []

  notify = {
    billing = false
  }
}
`, rString, rString, rString)
}

// Case when a user is on a team and that team updates it's permissions.
// This test is currently failing, as this is not implemented yet - this doesn't
// actually cause any issues because it's just Terraforms state that doesn't have the
// new permission values yet, the backend does, and when `terraform refresh` is ran,
// the state will be updated appropriately.
// The test is left here for documentation purposes.
/*
func TestAccUser_permissions_team_update(t *testing.T) {
	var user account.User
	rString := acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("terraform acc test user %s", rString)
	username := fmt.Sprintf("tf_acc_test_user_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsOnTeam(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
				),
			},
			{
				Config: testAccUserPermissionsTeamUpdate(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user.u", &user),
					resource.TestCheckResourceAttr("ns1_user.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_account_settings", "true"),
					resource.TestCheckResourceAttr("ns1_user.u", "account_manage_apikeys", "true"),
				),
			},
		},
	})
}

func testAccUserPermissionsTeamUpdate(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
  account_manage_apikeys = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  teams = ["${ns1_team.t.id}"]

  notify = {
  	billing = false
  }
}
`, rString, rString, rString)
}
*/

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		expErrs int
	}{
		{
			"valid",
			"vAlId_uS3r",
			0,
		},
		{
			"valid - email",
			"valid_us3r@example.com",
			0,
		},
		{
			"invalid - punctuation",
			"inv@l!d_us3r",
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWarns, outErrs := validateUsername(tt.in, "")
			assert.Equal(t, tt.expErrs, len(outErrs))
			assert.Equal(t, 0, len(outWarns))
		})
	}
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_user" {
			continue
		}

		user, _, err := client.Users.Get(rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("user still exists: %#v: %#v", err, user.Name)
		}
	}

	return nil
}

func testAccCheckUserExists(n string, user *account.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundUser, _, err := client.Users.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundUser.Username != rs.Primary.ID {
			return fmt.Errorf("user not found (%#v != %s)", foundUser, rs.Primary.ID)
		}

		*user = *foundUser

		return nil
	}
}

func testAccCheckUserIPWhitelists(user *account.User, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		sort.Strings(user.IPWhitelist)
		sort.Strings(expected)

		if !reflect.DeepEqual(user.IPWhitelist, expected) {
			return fmt.Errorf("IPWhitelist: got values: %v want: %v", user.IPWhitelist, expected)
		}
		return nil
	}
}

// Simulate a manual deletion of a user.
func testAccManualDeleteUser(user string) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Users.Delete(user)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete user: %v", err)
		}
	}
}

func testAccUserBasic(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"

  account_view_invoices = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"
  teams = ["${ns1_team.t.id}"]
  notify = {
    billing = false
  }
  ip_whitelist        = ["1.1.1.1", "2.2.2.2"]
  ip_whitelist_strict = true
}
`, rString, rString, rString)
}

func testAccUserUpdated(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"

  account_view_invoices = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"
  teams = ["${ns1_team.t.id}"]
  notify = {
    billing = true
  }
}
`, rString, rString, rString)
}

func testAccUserPermissionsOnTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
  account_manage_ip_whitelist = false
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  teams = ["${ns1_team.t.id}"]

  notify = {
    billing = false
  }
}
`, rString, rString, rString)
}

func testAccUserPermissionsNoTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  notify = {
    billing = false
  }

  account_manage_ip_whitelist = true
}
`, rString, rString, rString)
}

func testAccUserSecurityPermissionsNoTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  notify = {
    billing = false
  }

  account_manage_ip_whitelist = true
  security_manage_global_2fa = false
}
`, rString, rString, rString)
}

// Explicitly sets the users team to []
func testAccUserPermissionsEmptyTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  teams = []

  notify = {
    billing = false
  }
}
`, rString, rString, rString)
}

func testAccUserPermissionsOnTwoTeam(rString string) string {
	return fmt.Sprintf(`resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  account_manage_account_settings = true
}

resource "ns1_team" "t2" {
  name = "terraform acc test team %s-2"
  account_manage_apikeys = true
  account_manage_ip_whitelist = true
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  teams = ["${ns1_team.t.id}","${ns1_team.t2.id}"]

  notify = {
    billing = false
  }

  account_manage_ip_whitelist = true
}
`, rString, rString, rString, rString)
}
