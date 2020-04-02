# Document Resource

Provides [Mule B2B Document][1] resource.

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
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
  name           = "my-partner"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "my-id"
  }
}

resource "muleb2b_document" "test" {
  name = "my-document"
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  edi_document_type_id = muleb2b_ediDocumentType.test.id
}
```

## Argument Reference

* `environment_id` - (Required) ID of environment in which to add document
* `partner_id` - (Required) ID of partner in which to add document
* `name` - (Required) Name for the document
* `edi_document_type_id` - (Required) ID of the document's type. See [EDI Document Type Data Source](../data-sources/ediDocumentType.md)
* `schema_file` - (Optional) Base64 encoded contents of the custom schema file.

## Attribute Reference

* `id` - ID of the document
* `custom_schema_id` - ID of the custom schema if one was created

[1]: https://docs.mulesoft.com/partner-manager/2.0/document-types