package b2b

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccMuleB2bPartnerDS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePartner_InitialConfig(),
				Check:  testDataSourcePartner_InitialCheck(),
			},
		},
	})
}

func testDataSourcePartner_InitialConfig() string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

data "muleb2b_partner" "test" {
  name           = "MM-Test-2"
  environment_id = data.muleb2b_environment.sbx.id
}`)
}

func testDataSourcePartner_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hostId := s.Modules[0].Resources["data.muleb2b_partner.host"].Primary.ID
		testId := s.Modules[0].Resources["data.muleb2b_partner.test"].Primary.ID

		if hostId != "9e39423e-da65-429a-b735-941e4b3fd350" {
			return fmt.Errorf("partner ID (%s) does not match expected (9e39423e-da65-429a-b735-941e4b3fd350)", hostId)
		}

		if testId != "9e39423e-da65-429a-b735-941e4b3fd350" {
			return fmt.Errorf("partner ID (%s) does not match expected (9e39423e-da65-429a-b735-941e4b3fd350)", testId)
		}

		return nil
	}
}
