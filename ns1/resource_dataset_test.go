package ns1

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dataset"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataset_basic(t *testing.T) {
	var resultDt dataset.Dataset
	expectedDt, tfPlan := testAccDatasetBasic()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatasetDestroy,
		Steps: []resource.TestStep{
			{
				Config: tfPlan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasetExists(&resultDt, t),
					testAccCheckDatasetMatchExpected(expectedDt, &resultDt, t),
				),
			},
		},
	})
}
func testAccCheckDatasetExists(dt *dataset.Dataset, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources["ns1_dataset.my_dataset"]

		if rs == nil || rs.Primary.ID == "" {
			return fmt.Errorf("no id is set for the dataset")
		}

		client := testAccProvider.Meta().(*ns1.Client)

		foundDt, _, err := client.Datasets.Get(rs.Primary.Attributes["id"])
		if err != nil {
			t.Log(err)
			return err
		}

		if foundDt.ID != rs.Primary.Attributes["id"] {
			return fmt.Errorf("dataset mismatch: resource vs api")
		}

		*dt = *foundDt

		return nil
	}
}

func testAccCheckDatasetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ns1.Client)

	rs := s.RootModule().Resources["ns1_dataset.my_dataset"]
	if rs == nil || rs.Primary.ID == "" {
		return fmt.Errorf("no id is set for the dataset")
	}

	dt, _, _ := client.Datasets.Get(rs.Primary.Attributes["id"])
	if dt != nil {
		return fmt.Errorf("dataset still exists: %#v", dt)
	}

	return nil
}

func testAccCheckDatasetMatchExpected(expected, result *dataset.Dataset, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, expected.Datatype, result.Datatype)
		assert.Equal(t, expected.Repeat, result.Repeat)
		assert.Equal(t, expected.Timeframe, result.Timeframe)
		assert.Equal(t, expected.ExportType, result.ExportType)
		assert.Equal(t, expected.RecipientEmails, result.RecipientEmails)
		return nil
	}
}

func testAccDatasetBasic() (*dataset.Dataset, string) {
	var timeframeCycles = int32(1)
	var repeatStart = time.Now().Add(time.Minute).Unix()

	dt := dataset.NewDataset(
		"",
		fmt.Sprintf("tf-test-dataset-%s", acctest.RandStringFromCharSet(15, acctest.CharSetAlphaNum)),
		&dataset.Datatype{
			Type:  dataset.DatatypeTypeNumQueries,
			Scope: dataset.DatatypeScopeAccount,
			Data:  nil,
		},
		&dataset.Repeat{
			Start:        dataset.UnixTimestamp(time.Unix(repeatStart, 0)),
			RepeatsEvery: dataset.RepeatsEveryMonth,
			EndAfterN:    1,
		},
		&dataset.Timeframe{
			Aggregation: dataset.TimeframeAggregationMontly,
			Cycles:      &timeframeCycles,
		},
		dataset.ExportTypeCSV,
		nil,
		nil,
		dataset.UnixTimestamp{},
		dataset.UnixTimestamp{},
	)

	plan := fmt.Sprintf(
		`
			resource "ns1_dataset" "my_dataset" {
				name     = "%s"
				datatype {
					type  = "%s"
					scope = "%s"
					data  = {}
				}
				repeat {
					start = %d
					repeats_every = "%s"
					end_after_n = %d
				}
				timeframe {
					aggregation = "%s"
					cycles      = %d
				}
				export_type = "%s"
			}
		`,
		dt.Name,
		dt.Datatype.Type,
		dt.Datatype.Scope,
		time.Time(dt.Repeat.Start).Unix(),
		dt.Repeat.RepeatsEvery,
		dt.Repeat.EndAfterN,
		dt.Timeframe.Aggregation,
		*dt.Timeframe.Cycles,
		dt.ExportType,
	)

	return dt, plan
}
