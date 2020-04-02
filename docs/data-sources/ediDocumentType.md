# EDI Document Type Data Source

Provides data for [Mule B2B Document Types][1]

## Example Usage

```hcl
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
```

## Argument Reference

* `environment_id` - (Required) The ID of the environment to lookup document type in
* `format_type` - (Required) The name of the EDI format type
* `format_version` - (Required) The version of the EDI format type
* `document_name` - (Required) The exact name of the document

## Attribute Reference

* `id` - ID of the EDI Document Type
* `edi_format_id` - The ID of the EDI format
* `edi_format_version_id` - The ID of the EDI format version

[1]: https://docs.mulesoft.com/partner-manager/2.0/document-types