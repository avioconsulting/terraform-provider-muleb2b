package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccMuleB2bDocument(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceDocument_InitialConfig(name),
				Check:  testResourceDocument_InitialCheck(),
			},
			{
				Config: testResourceDocument_UpdateConfig(name),
				Check:  testResourceDocument_UpdateCheck(),
			},
		},
	})
}

func testResourceDocument_InitialConfig(name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "Sandbox"
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
  edi_document_type_id = "6117a01c-d661-4517-80a5-5a6fe08d833c"
}`, name, name, name)
}

func testResourceDocument_InitialCheck() resource.TestCheckFunc {
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

		return nil
	}
}

func testResourceDocument_UpdateConfig(name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "Sandbox"
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
  edi_document_type_id = "f4190af3-a1ee-4aa7-8a1c-67d385b57c90"
}`, name, name, name)
}

func testResourceDocument_UpdateCheck() resource.TestCheckFunc {
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
		doc, err := client.GetDocumentById(partner, id)
		if err != nil {
			return err
		}

		if doc == nil || *doc.EdiDocumentTypeId != "f4190af3-a1ee-4aa7-8a1c-67d385b57c90" {
			return fmt.Errorf("server_address did not update")
		}

		return nil
	}
}
