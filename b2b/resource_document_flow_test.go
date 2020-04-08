package b2b

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccMuleB2bResourceDocumentFlow(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceDocumentFlow_InitialConfig(envName, name),
				Check:  testResourceDocumentFlow_InitialCheck(),
			},
			{
				Config: testResourceDocumentFlow_UpdateConfig(envName, name),
				Check:  testResourceDocumentFlow_UpdateCheck(),
			},
		},
	})
}

func testResourceDocumentFlow_InitialConfig(envName, name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
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

resource "muleb2b_endpoint" "test" {
  name = "%s"
  role = "receive"
  type = "http"
  partner_id = data.muleb2b_partner.host.id
  environment_id = data.muleb2b_environment.sbx.id
  http_config {
    server_address = "accTest.mytest.com"
    server_port = 80
    path = "/"
    protocol = "http"
    auth_mode  {
      type = "none"
    }
  }
}

resource "muleb2b_document_flow" "test" {
  name = "%s"
  direction = "inbound"
  environment_id = data.muleb2b_environment.sbx.id
  partner_from_id = muleb2b_partner.test.id
  partner_to_id = data.muleb2b_partner.host.id
  config {
    receiving_endpoint_id = muleb2b_endpoint.test.id
  }
}
`, envName, name, name, name, name)
}

func testResourceDocumentFlow_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dfId := s.Modules[0].Resources["muleb2b_document_flow.test"].Primary.ID

		if dfId == "" {
			return fmt.Errorf("document_flow.id is empty")
		}

		return nil
	}
}

func testResourceDocumentFlow_UpdateConfig(envName, name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
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

resource "muleb2b_endpoint" "test" {
  name = "%s"
  role = "receive"
  type = "http"
  partner_id = data.muleb2b_partner.host.id
  environment_id = data.muleb2b_environment.sbx.id
  http_config {
    server_address = "accTest.mytest.com"
    server_port = 80
    path = "/"
    protocol = "http"
    auth_mode  {
      type = "none"
    }
  }
}

data "muleb2b_ediDocumentType" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  format_type = "X12"
  format_version = "3010"
  document_name = "819"
}

data "muleb2b_ediDocumentType" "json" {
  environment_id = data.muleb2b_environment.sbx.id
  format_type = "JSON"
  format_version = "V1"
  document_name = "JSON"
}

resource "muleb2b_document" "source" {
  name = "%s-source"
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  edi_document_type_id = data.muleb2b_ediDocumentType.test.id
}

resource "muleb2b_document" "target" {
  name = "%s-target"
  partner_id = data.muleb2b_partner.host.id
  environment_id = data.muleb2b_environment.sbx.id
  edi_document_type_id = data.muleb2b_ediDocumentType.json.id
  schema_file = "ewogICIkaWQiOiAiaHR0cHM6Ly9leGFtcGxlLmNvbS9wZXJzb24uc2NoZW1hLmpzb24iLAogICIkc2NoZW1hIjogImh0dHA6Ly9qc29uLXNjaGVtYS5vcmcvZHJhZnQtMDcvc2NoZW1hIyIsCiAgInRpdGxlIjogIlBlcnNvbiIsCiAgInR5cGUiOiAib2JqZWN0IiwKICAicHJvcGVydGllcyI6IHsKICAgICJmaXJzdE5hbWUiOiB7CiAgICAgICJ0eXBlIjogInN0cmluZyIsCiAgICAgICJkZXNjcmlwdGlvbiI6ICJUaGUgcGVyc29uJ3MgZmlyc3QgbmFtZS4iCiAgICB9LAogICAgImxhc3ROYW1lIjogewogICAgICAidHlwZSI6ICJzdHJpbmciLAogICAgICAiZGVzY3JpcHRpb24iOiAiVGhlIHBlcnNvbidzIGxhc3QgbmFtZS4iCiAgICB9LAogICAgImFnZSI6IHsKICAgICAgImRlc2NyaXB0aW9uIjogIkFnZSBpbiB5ZWFycyB3aGljaCBtdXN0IGJlIGVxdWFsIHRvIG9yIGdyZWF0ZXIgdGhhbiB6ZXJvLiIsCiAgICAgICJ0eXBlIjogImludGVnZXIiLAogICAgICAibWluaW11bSI6IDAKICAgIH0KICB9Cn0="
}

resource "muleb2b_document_flow" "test" {
  name = "%s"
  direction = "inbound"
  environment_id = data.muleb2b_environment.sbx.id
  partner_from_id = muleb2b_partner.test.id
  partner_to_id = data.muleb2b_partner.host.id
  config {
    receiving_endpoint_id = muleb2b_endpoint.test.id
    source_doc_type_id = muleb2b_document.source.id
    target_doc_type_id = muleb2b_document.target.id
  }
}`, envName, name, name, name, name, name, name)
}

func testResourceDocumentFlow_UpdateCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		prim := s.Modules[0].Resources["muleb2b_document_flow.test"].Primary.ID
		if prim == "" {
			return fmt.Errorf("muleb2b_document_flow.test.config.source_doc_type_id is empty")
		}
		return nil
	}
}
