package ns1

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ns1": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NS1_APIKEY"); v == "" {
		t.Fatal("NS1_APIKEY must be set for acceptance tests")
	}
}

func testAccPreCheckSso(t *testing.T) {
	testAccPreCheck(t)

	client, err := sharedClient()
	if err != nil {
		log.Fatalf("failed to get shared client: %s", err)
	}
	kl, _, err := client.APIKeys.List()
	if err != nil {
		t.Skipf("account not authorized for redirects, skipping test")
	}
	if len(kl) == 0 {
		t.Skipf("no api keys found, skipping test")
	}
	if kl[0].Permissions.Account.ManageUsers == false {
		t.Skipf("account not authorized to manage users, skipping test")
	}

}
