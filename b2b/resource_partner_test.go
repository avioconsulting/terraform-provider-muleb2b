package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccMuleB2bPartner(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourcePartner_InitialConfig(name),
				Check:  testResourcePartner_InitialCheck(),
			},
			{
				Config: testResourcePartner_UpdateConfig(name),
				Check:  testResourcePartner_UpdateCheck(),
			},
			{
				Config: testResourcePartner_UpdateConfig2(name),
				Check:  testResourcePartner_UpdateCheck2(),
			},
		},
	})
}

func testResourcePartner_InitialConfig(name string) string {
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
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "remove-me-%s"
  }
  address {
    address_line_1 = "123 Main Street"
    city = "Anytown"
    state = "NY"
    country = "US"
    postal_code = "12345"
  }
  contact {
    name = "John Doe"
    email = "test@test.com"
    phone = "2511231234"
    type = "business"
  }
  contact {
    name = "John Legend"
    email = "john.legend@test.com"
    type = "technical"
  }
  x12_inbound_config {
    character_encoding = "UTF8"
    acknowledgements {}
    validations {}
    control_numbers {}
  }
}
`, name, name, name)
}

func testResourcePartner_InitialCheck() resource.TestCheckFunc {
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
		client.SetEnvironment("3a4d3936-22d9-4d87-a3c0-a8d424bcc032")
		partner, err := client.GetPartner(id)
		if err != nil {
			return err
		}

		identifiers, err := client.ListPartnerIdentifiers(*partner.Id)
		if err != nil {
			return err
		} else if len(identifiers) < 1 {
			return fmt.Errorf("partner identifier did not get created")
		}

		address, err := client.GetPartnerAddress(*partner.Id)
		if err != nil {
			return err
		} else if address == nil {
			return fmt.Errorf("address was not set")
		} else if address.City == nil || *address.City != "Anytown" {
			return fmt.Errorf("city not set")
		}

		contacts, err := client.GetPartnerContacts(*partner.Id)
		if err != nil {
			return err
		} else if len(contacts) != 2 {
			return fmt.Errorf("contact identifiers did not get created")
		}

		x12, err := client.GetPartnerInboundX12Configuration(*partner.Id)
		if err != nil {
			return err
		} else if x12.CharacterSetAndEncoding.CharacterEncoding == nil {
			return fmt.Errorf("nil character encoding")
		} else if *x12.CharacterSetAndEncoding.CharacterEncoding != "UTF8" {
			return fmt.Errorf("incorrect character encoding (%s)", *x12.CharacterSetAndEncoding.CharacterEncoding)
		}
		return nil
	}
}

func testResourcePartner_UpdateConfig(name string) string {
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
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "remove-me-%s"
  }
  address {
    address_line_1 = "123 Main Street"
    address_line_2 = "Ste 1"
    city = "Anytown"
    state = "NY"
    country = "US"
    postal_code = "12345"
  }
  contact {
    name = "John Doe"
    email = "test@test.com"
    phone = "2511231234"
    type = "business"
  }
  contact {
    name = "Bon Jovi"
    email = "bon.jovi@test.com"
    type = "technical"
  }
  contact {
    name = "Cyndi Lauper"
    email = "c.lauper@test.com"
    phone = "1231231234"
    type = "other"
  }
  x12_inbound_config {
    character_encoding = "UTF8"
    acknowledgements {}
    validations {
      fail_when_unused_segments_included = true
      fail_when_too_many_repeats_of_segment = false
    }
    control_numbers {}
  }
}
`, name, name, name)
}

func testResourcePartner_UpdateCheck() resource.TestCheckFunc {
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
		client.SetEnvironment("3a4d3936-22d9-4d87-a3c0-a8d424bcc032")
		partner, err := client.GetPartner(id)
		if err != nil {
			return err
		}

		identifiers, err := client.ListPartnerIdentifiers(*partner.Id)
		if err != nil {
			return err
		} else if len(identifiers) < 2 {
			return fmt.Errorf("additional partner identifiers did not get created")
		}

		address, err := client.GetPartnerAddress(*partner.Id)
		if err != nil {
			return err
		} else if address == nil {
			return fmt.Errorf("address was not set")
		} else if address.Addr2 == nil || *address.Addr2 != "Ste 1" {
			return fmt.Errorf("address_line_2 not set")
		}

		contacts, err := client.GetPartnerContacts(*partner.Id)
		if err != nil {
			return err
		} else if len(contacts) != 3 {
			return fmt.Errorf("contact identifiers did not get created")
		}

		x12, err := client.GetPartnerInboundX12Configuration(*partner.Id)
		if err != nil {
			return err
		} else if x12.ParserSettings.FailDocumentWhenUnusedSegmentsAreIncluded == nil {
			return fmt.Errorf("nil FailDocumentWhenUnusedSegmentsAreIncluded found")
		} else if *x12.ParserSettings.FailDocumentWhenUnusedSegmentsAreIncluded != true {
			return fmt.Errorf("incorrect value for FailDocumentWhenUnusedSegmentsAreIncluded found, should be true")
		} else if x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment == nil {
			return fmt.Errorf("nil FailDocumentWhenTooManyRepeatsOfSegment found")
		} else if *x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment != false {
			return fmt.Errorf("incorrect value for FailDocumentWhenTooManyRepeatsOfSegment found, should be false")
		}

		return nil
	}
}

func testResourcePartner_UpdateConfig2(name string) string {
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
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "remove-me-%s"
  }
  contact {
    name = "John Doe"
    email = "test@test.com"
    phone = "2511231234"
    type = "business"
  }
  contact {
    name = "Bon Jovi"
    email = "bon.jovi@test.com"
    type = "technical"
  }
  contact {
    name = "Cyndi Lauper"
    email = "c.lauper@test.com"
    phone = "1231231234"
    type = "other"
  }
  x12_inbound_config {
    character_encoding = "UTF8"
    acknowledgements {}
    validations {
      fail_when_unused_segments_included = true
      fail_when_too_many_repeats_of_segment = false
    }
    control_numbers {}
  }
}
`, name, name, name)
}

func testResourcePartner_UpdateCheck2() resource.TestCheckFunc {
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
		client.SetEnvironment("3a4d3936-22d9-4d87-a3c0-a8d424bcc032")
		partner, err := client.GetPartner(id)
		if err != nil {
			return err
		}

		identifiers, err := client.ListPartnerIdentifiers(*partner.Id)
		if err != nil {
			return err
		} else if len(identifiers) < 2 {
			return fmt.Errorf("additional partner identifiers did not get created")
		}

		// Address cannot be fully deleted, the contents can be emptied though
		address, err := client.GetPartnerAddress(*partner.Id)
		if err != nil {
			return err
		} else if address == nil {
			return fmt.Errorf("address was not set")
		} else if address.Addr1 != nil && *address.Addr1 != "" {
			return fmt.Errorf("address_line_1 is set to non-empty value")
		}

		contacts, err := client.GetPartnerContacts(*partner.Id)
		if err != nil {
			return err
		} else if len(contacts) != 3 {
			return fmt.Errorf("contact identifiers did not get created")
		}

		x12, err := client.GetPartnerInboundX12Configuration(*partner.Id)
		if err != nil {
			return err
		} else if x12.ParserSettings.FailDocumentWhenUnusedSegmentsAreIncluded == nil {
			return fmt.Errorf("nil FailDocumentWhenUnusedSegmentsAreIncluded found")
		} else if *x12.ParserSettings.FailDocumentWhenUnusedSegmentsAreIncluded != true {
			return fmt.Errorf("incorrect value for FailDocumentWhenUnusedSegmentsAreIncluded found, should be true")
		} else if x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment == nil {
			return fmt.Errorf("nil FailDocumentWhenTooManyRepeatsOfSegment found")
		} else if *x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment != false {
			return fmt.Errorf("incorrect value for FailDocumentWhenTooManyRepeatsOfSegment found, should be false")
		}

		return nil
	}
}
