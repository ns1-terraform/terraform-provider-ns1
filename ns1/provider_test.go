package ns1

import (
	// "log"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/ns1/ns1-go.v2/rest/model/alerting"
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
		t.Fatalf("failed to get shared client: %s", err)
	}

	// Create a test alert with saml_certificate_expired subtype
	var subtype string = "saml_certificate_expired"
	var typeAlert string = "account"
	testAlertName := fmt.Sprintf("terraform-precheck-sso-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	testAlert := &alerting.Alert{
		Name:            &testAlertName,
		Type:            &typeAlert,
		Subtype:         &subtype,
		NotifierListIds: []string{},
	}

	// Try to create the alert
	resp, err := client.Alerts.Create(testAlert)
	if err != nil {
		t.Skipf("skipping SSO test; unable to create saml_certificate_expired alert: %s", err)
	}

	// If creation succeeded, clean up by deleting the test alert
	if resp.StatusCode == 201 && testAlert.ID != nil {
		_, deleteErr := client.Alerts.Delete(*testAlert.ID)
		if deleteErr != nil {
			log.Printf("warning: failed to clean up precheck alert: %v", deleteErr)
		}
	}
}
