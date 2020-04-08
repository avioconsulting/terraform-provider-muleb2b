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

func TestAccMuleB2bIdentifier(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	number := acctest.RandIntRange(100, 10000)
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceIdnetifier_InitialConfig(envName, name, number),
				Check:  testResourceIdentifier_InitialCheck(number),
			},
			{
				Config: testResourceIdentifier_UpdateConfig(envName, name, number),
				Check:  testResourceIdentifier_UpdateCheck(number),
			},
		},
	})
}

func testResourceIdnetifier_InitialConfig(envName, name string, number int) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
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
    value = "%d1"
  }
}

resource "muleb2b_identifier" "abc" {
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  identifier_type_id = data.muleb2b_identifier_type.duns.id
  value = "%d2"
}
`, envName, name, number, number)
}

func testResourceIdentifier_InitialCheck(number int) resource.TestCheckFunc {
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

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%d1", number)) {
			return fmt.Errorf("identifier with value of %d1 not created", number)
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%d2", number)) {
			return fmt.Errorf("identifier with value of %d2 not created", number)
		}

		return nil
	}
}

func testResourceIdentifier_UpdateConfig(envName, name string, number int) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
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
    value = "%d1"
  }
}

resource "muleb2b_identifier" "abc" {
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  identifier_type_id = data.muleb2b_identifier_type.duns.id
  value = "%d3"
}
`, envName, name, number, number)
}

func testResourceIdentifier_UpdateCheck(number int) resource.TestCheckFunc {
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

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%d1", number)) {
			return fmt.Errorf("identifier with value of %d1 not created", number)
		}

		if !testResourceIdentifierFindExpected(identifiers, fmt.Sprintf("%d3", number)) {
			return fmt.Errorf("identifier with value of %d3 not created", number)
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
