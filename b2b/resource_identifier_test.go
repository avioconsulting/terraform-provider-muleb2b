package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccMuleB2bIdentifier(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceIdnetifier_InitialConfig(name),
				Check:  testResourceIdentifier_InitialCheck(name),
			},
			{
				Config: testResourceIdentifier_UpdateConfig(name),
				Check:  testResourceIdentifier_UpdateCheck(name),
			},
		},
	})
}

func testResourceIdnetifier_InitialConfig(name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "duns" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "DUNS"
}

resource "muleb2b_partner" "test" {
  name           = "%s"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "%s-id1"
  }
}

resource "muleb2b_identifier" "abc" {
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  identifier_type_id = data.muleb2b_identifier_type.duns.id
  value = "%s-id2"
}
`, name, name, name)
}

func testResourceIdentifier_InitialCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.Modules[0].Resources["muleb2b_partner.test"]
		if resourceState == nil {
			return fmt.Errorf("resource not found in state")
		}

		instanceState := resourceState.Primary
		if instanceState == nil {
			return fmt.Errorf("resource has no primary instance")
		}

		id := instanceState.ID

		if id == "" {
			return fmt.Errorf("id is not set")
		}

		client := testAccProvider.Meta().(*muleb2b.Client)

		identifiers, err := client.ListPartnerIdentifiers(id)
		if err != nil {
			return err
		} else if len(identifiers) != 2 {
			return fmt.Errorf("partner identifier did not get created")
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%s-id1", name)) {
			return fmt.Errorf("identifier with value of %s-id1 not created", name)
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%s-id2", name)) {
			return fmt.Errorf("identifier with value of %s-id2 not created", name)
		}

		return nil
	}
}

func testResourceIdentifier_UpdateConfig(name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "duns" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "DUNS"
}

resource "muleb2b_partner" "test" {
  name           = "%s"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "%s-id1"
  }
}

resource "muleb2b_identifier" "abc" {
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  identifier_type_id = data.muleb2b_identifier_type.duns.id
  value = "%s-id3"
}
`, name, name, name)
}

func testResourceIdentifier_UpdateCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState := s.Modules[0].Resources["muleb2b_partner.test"]
		if resourceState == nil {
			return fmt.Errorf("resource not found in state")
		}

		instanceState := resourceState.Primary
		if instanceState == nil {
			return fmt.Errorf("resource has no primary instance")
		}

		id := instanceState.ID

		if id == "" {
			return fmt.Errorf("id is not set")
		}

		client := testAccProvider.Meta().(*muleb2b.Client)

		identifiers, err := client.ListPartnerIdentifiers(id)
		if err != nil {
			return err
		} else if len(identifiers) < 2 {
			return fmt.Errorf("additional partner identifiers did not get created")
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%s-id1", name)) {
			return fmt.Errorf("identifier with value of %s-id1 not created", name)
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%s-id3", name)) {
			return fmt.Errorf("identifier with value of %s-id3 not created", name)
		}

		return nil
	}
}

func testResourceIdentifierFindExpected(list []*muleb2b.Identifier, value string) bool {
	for _, i := range list {
		if i.Value != nil && *i.Value == value {
			return true
		}
	}
	return false
}
