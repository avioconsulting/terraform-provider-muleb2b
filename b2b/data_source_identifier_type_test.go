package b2b

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccMuleB2bIdentifierType(t *testing.T) {
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIdentifierType_InitialConfig(envName),
				Check:  testDataSourceIdentifierType_InitialCheck(),
			},
		},
	})
}

func testDataSourceIdentifierType_InitialConfig(envName string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_identifier_type" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}`, envName)
}

func testDataSourceIdentifierType_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		docId := s.Modules[0].Resources["data.muleb2b_identifier_type.test"].Primary.ID

		if docId != "25c1bc8a-801f-4947-a2a6-7721ef971460" {
			return fmt.Errorf("EDI Document Type ID (%s) does not match expected (25c1bc8a-801f-4947-a2a6-7721ef971460)", docId)
		}

		return nil
	}
}

func TestAccMuleB2bIdentifierTypeWithMultipleQualifiers(t *testing.T) {
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIdentifierType_InitialConfigMultipleQualifiers(envName),
				Check:  testDataSourceIdentifierType_InitialCheckMultipleQualifiers(),
			},
		},
	})
}

func testDataSourceIdentifierType_InitialConfigMultipleQualifiers(envName string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_identifier_type" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "X12-ISA"
  qualifier_code = "12"
}`, envName)
}

func testDataSourceIdentifierType_InitialCheckMultipleQualifiers() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		docId := s.Modules[0].Resources["data.muleb2b_identifier_type.test"].Primary.ID

		if docId != "26432f80-b58b-4c96-aac6-58af5d5580fc" {
			return fmt.Errorf("EDI Document Type ID (%s) does not match expected (26432f80-b58b-4c96-aac6-58af5d5580fc)", docId)
		}

		return nil
	}
}
