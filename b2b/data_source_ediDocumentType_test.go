package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccMuleB2bEdiDocumentType(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEdiDocumentType_InitialConfig(name, envName),
				Check:  testDataSourceEdiDocumentType_InitialCheck(),
			},
		},
	})
}

func testDataSourceEdiDocumentType_InitialConfig(name, envName string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_ediDocumentType" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  format_type = "X12"
  format_version = "3010"
  document_name = "819"
}

data "muleb2b_identifier_type" "as2" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}

resource "muleb2b_partner" "test" {
  name           = "%s"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "%s-id1"
  }
}

resource "muleb2b_document" "test" {
  name = "%s"
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  edi_document_type_id = data.muleb2b_ediDocumentType.test.id
}`, envName, name, name, name)
}

func testDataSourceEdiDocumentType_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.Modules[0].Resources["muleb2b_document.test"]
		if resourceState == nil {
			return fmt.Errorf("resource not found in state")
		}

		instanceState := resourceState.Primary
		if instanceState == nil {
			return fmt.Errorf("resource has no primary instance")
		}

		id := instanceState.ID
		partner := instanceState.Attributes["partner_id"]

		if partner == "" {
			return fmt.Errorf("partner_id is empty")
		}

		if id == "" {
			return fmt.Errorf("id is not set")
		}

		client := testAccProvider.Meta().(*muleb2b.Client)
		_, err := client.GetDocumentById(partner, id)
		if err != nil {
			return err
		}

		docId := s.Modules[0].Resources["data.muleb2b_ediDocumentType.test"].Primary.ID

		if docId != "f4190af3-a1ee-4aa7-8a1c-67d385b57c90" {
			return fmt.Errorf("EDI Document Type ID (%s) does not match expected (f4190af3-a1ee-4aa7-8a1c-67d385b57c90)", docId)
		}

		return nil
	}
}
