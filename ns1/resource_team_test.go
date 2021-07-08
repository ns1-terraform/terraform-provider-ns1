package ns1

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestAccTeam_basic(t *testing.T) {
	var team account.Team
	n := fmt.Sprintf("terraform test team %s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTeamBasic, n),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamExists("ns1_team.foobar", &team),
					testAccCheckTeamName(&team, n),
					testAccCheckTeamDNSPermission(&team, "view_zones", true),
					testAccCheckTeamDNSPermission(&team, "zones_allow_by_default", true),
					testAccCheckTeamDNSPermissionZones(&team, "zones_allow", []string{"mytest.zone"}),
					testAccCheckTeamDNSPermissionZones(&team, "zones_deny", []string{"myother.zone"}),
					testAccCheckTeamDNSPermissionRecords(&team, "dns_records_allow", []account.Record{
						{Domain: "my.ns1.com", Subdomains: false, Zone: "ns1.com", RecordType: "A"}}),
					testAccCheckTeamDNSPermissionRecords(&team, "dns_records_deny", []account.Record{
						{Domain: "my.test.com", Subdomains: true, Zone: "test.com", RecordType: "A"}}),
					testAccCheckTeamDataPermission(&team, "manage_datasources", true),
					testAccCheckTeamIPWhitelists(&team, []account.IPWhitelist{
						{Name: "whitelist-1", Values: []string{"1.1.1.1", "2.2.2.2"}},
						{Name: "whitelist-2", Values: []string{"3.3.3.3", "4.4.4.4"}},
					}),
				),
			},
		},
	})
}

func TestAccTeam_updated(t *testing.T) {
	var team account.Team
	n := fmt.Sprintf("terraform test team %s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTeamBasic, n),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamExists("ns1_team.foobar", &team),
					testAccCheckTeamName(&team, n),
				),
			},
			{
				Config: fmt.Sprintf(testAccTeamUpdated, n),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamExists("ns1_team.foobar", &team),
					testAccCheckTeamName(&team, fmt.Sprintf("%s updated", n)),
					testAccCheckTeamDNSPermission(&team, "view_zones", true),
					testAccCheckTeamDNSPermission(&team, "zones_allow_by_default", true),
					testAccCheckTeamDNSPermissionZones(&team, "zones_allow", []string{}),
					testAccCheckTeamDNSPermissionZones(&team, "zones_deny", []string{}),
					testAccCheckTeamDNSPermissionRecords(&team, "dns_records_allow", []account.Record{}),
					testAccCheckTeamDNSPermissionRecords(&team, "dns_records_deny", []account.Record{}),
					testAccCheckTeamDataPermission(&team, "manage_datasources", false),
					testAccCheckTeamIPWhitelists(&team, []account.IPWhitelist{}),
				),
			},
		},
	})
}

// Verifies that a team is re-created correctly if it is manually deleted.
func TestAccTeam_ManualDelete(t *testing.T) {
	var team account.Team

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamBasic,
				Check:  testAccCheckTeamExists("ns1_team.foobar", &team),
			},
			// Simulate a manual deletion of the team and verify that the plan has a diff.
			{
				PreConfig:          testAccManualDeleteTeam(&team),
				Config:             testAccTeamBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Then re-create and make sure it is there again.
			{
				Config: testAccTeamBasic,
				Check:  testAccCheckTeamExists("ns1_team.foobar", &team),
			},
		},
	})
}

func testAccCheckTeamExists(n string, team *account.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NoID is set")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundTeam, _, err := client.Teams.Get(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		if foundTeam.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Team not found")
		}

		*team = *foundTeam

		return nil
	}
}

func testAccCheckTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ns1_team" {
			continue
		}

		team, _, err := client.Teams.Get(rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("Team still exists: %#v: %#v", err, team.Name)
		}
	}

	return nil
}

func testAccCheckTeamName(team *account.Team, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if team.Name != expected {
			return fmt.Errorf("Name: got: %s want: %s", team.Name, expected)
		}
		return nil
	}
}

func testAccCheckTeamDNSPermission(team *account.Team, perm string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dns := team.Permissions.DNS

		switch perm {
		case "view_zones":
			if dns.ViewZones != expected {
				return fmt.Errorf("DNS.ViewZones: got: %t want: %t", dns.ViewZones, expected)
			}
		case "manage_zones":
			if dns.ManageZones != expected {
				return fmt.Errorf("DNS.ManageZones: got: %t want: %t", dns.ManageZones, expected)
			}
		case "zones_allow_by_default":
			if dns.ZonesAllowByDefault != expected {
				return fmt.Errorf("DNS.ZonesAllowByDefault: got: %t want: %t", dns.ZonesAllowByDefault, expected)
			}
		}

		return nil
	}
}

func testAccCheckTeamDataPermission(team *account.Team, perm string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		data := team.Permissions.Data

		switch perm {
		case "push_to_datafeeds":
			if data.PushToDatafeeds != expected {
				return fmt.Errorf("Data.PushToDatafeeds: got: %t want: %t", data.PushToDatafeeds, expected)
			}
		case "manage_datasources":
			if data.ManageDatasources != expected {
				return fmt.Errorf("Data.ManageDatasources: got: %t want: %t", data.ManageDatasources, expected)
			}
		case "manage_datafeeds":
			if data.ManageDatafeeds != expected {
				return fmt.Errorf("Data.ManageDatafeeds: got: %t want: %t", data.ManageDatafeeds, expected)
			}
		}

		return nil
	}
}

func testAccCheckTeamDNSPermissionZones(team *account.Team, perm string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		dns := team.Permissions.DNS

		switch perm {
		case "zones_allow":
			if !reflect.DeepEqual(dns.ZonesAllow, expected) {
				return fmt.Errorf("DNS.ZonesAllow: got: %v want: %v", dns.ZonesAllow, expected)
			}
		case "zones_deny":
			if !reflect.DeepEqual(dns.ZonesDeny, expected) {
				return fmt.Errorf("DNS.ZonesDeny: got: %v want: %v", dns.ZonesDeny, expected)
			}
		}

		return nil
	}
}

func testAccCheckTeamDNSPermissionRecords(team *account.Team, perm string, expected []account.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dns := team.Permissions.DNS

		switch perm {
		case "dns_records_allow":
			if !reflect.DeepEqual(dns.RecordsAllow, expected) {
				return fmt.Errorf("DNS.RecordAllow: got: %v want: %v", dns.RecordsAllow, expected)
			}
		case "dns_records_deny":
			if !reflect.DeepEqual(dns.RecordsDeny, expected) {
				return fmt.Errorf("DNS.RecordDeny: got: %v want: %v", dns.RecordsDeny, expected)
			}
		}

		return nil
	}
}

func testAccCheckTeamIPWhitelists(team *account.Team, expected []account.IPWhitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(team.IPWhitelist) != len(expected) {
			return fmt.Errorf("IPWhitelist: got length: %v want: %v", len(team.IPWhitelist), len(expected))
		}

		for i, l := range expected {
			if l.Name != team.IPWhitelist[i].Name {
				return fmt.Errorf("IPWhitelist: got name: %v want: %v", team.IPWhitelist[i].Name, l.Name)
			}

			if !reflect.DeepEqual(l.Values, team.IPWhitelist[i].Values) {
				return fmt.Errorf("IPWhitelist: got values: %v want: %v", team.IPWhitelist[i].Values, l.Values)
			}
		}
		return nil
	}
}

// Simulate a manual deletion of a team.
func testAccManualDeleteTeam(team *account.Team) func() {
	return func() {
		client := testAccProvider.Meta().(*ns1.Client)
		_, err := client.Teams.Delete(team.ID)
		// Not a big deal if this fails, it will get caught in the test conditions and fail the test.
		if err != nil {
			log.Printf("failed to delete team: %v", err)
		}
	}
}

const testAccTeamBasic = `
resource "ns1_team" "foobar" {
  name = "%s"

  dns_view_zones = true
  dns_zones_allow_by_default = true
  dns_zones_allow = ["mytest.zone"]
  dns_zones_deny = ["myother.zone"]

  data_manage_datasources = true

  ip_whitelist {
	name = "whitelist-1"
	values = ["1.1.1.1", "2.2.2.2"]
  }

  ip_whitelist {
	name = "whitelist-2"
	values = ["3.3.3.3", "4.4.4.4"]
  }

  dns_records_allow {
	domain = "my.ns1.com"
	include_subdomains = false
	zone = "ns1.com"
	type = "A"
}

  dns_records_deny {
	domain = "my.test.com"
	include_subdomains = true
	zone = "test.com"
	type = "A"
 }

}`

const testAccTeamUpdated = `
resource "ns1_team" "foobar" {
  name = "%s updated"

  dns_view_zones = true
  dns_zones_allow_by_default = true

  data_manage_datasources = false
}`
