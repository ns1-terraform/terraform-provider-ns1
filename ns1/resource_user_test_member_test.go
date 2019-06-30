package ns1

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestAccUserTeamMember_basic(t *testing.T) {
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
				Config: testAccUserTeamMemberBasic(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ns1_user_team_member.u", &user),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "email", "tf_acc_test_ns1@hashicorp.com"),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "name", name),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "teams.#", "1"),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "notify.%", "1"),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "notify.billing", "true"),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "username", username),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "dns_manage_zones", "true"),
					// Test for permissions inherited from team
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "account_manage_plan", "true"),
					resource.TestCheckResourceAttr("ns1_user_team_member.u", "dns_manage_zones", "true"),
				),
			},
		},
	})
}

func testAccCheckUserTeamMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_user_team_member" {
			continue
		}

		user, _, err := client.Users.Get(rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("User still exists: %#v: %#v", err, user.Name)
		}
	}

	return nil
}

func testAccUserTeamMemberBasic(rString string) string {
	return fmt.Sprintf(`
resource "ns1_team" "t" {
  name = "terraform acc test team %s"
  dns_manage_zones = true
  account_manage_plan = true
}

resource "ns1_user_team_member" "u" {
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
