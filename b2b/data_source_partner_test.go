package b2b

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccMuleB2bPartnerDS(t *testing.T) {
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePartner_InitialConfig(envName),
				Check:  testDataSourcePartner_InitialCheck(),
			},
		},
	})
}

func testDataSourcePartner_InitialConfig(envName string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

data "muleb2b_partner" "test" {
  name           = data.muleb2b_partner.host.name
  environment_id = data.muleb2b_environment.sbx.id
}`, envName)
}

func testDataSourcePartner_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hostId := s.Modules[0].Resources["data.muleb2b_partner.host"].Primary.ID
		testId := s.Modules[0].Resources["data.muleb2b_partner.test"].Primary.ID

		if hostId != testId {
			return fmt.Errorf("partner ID (%s) does not match expected (%s)", hostId, testId)
		}

		return nil
	}
}
