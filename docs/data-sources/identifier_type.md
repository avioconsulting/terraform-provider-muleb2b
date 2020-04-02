# Identifier Type Data Source

Provides data for [Mule B2B Identifier Types][1]


## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}

data "muleb2b_identifier_type" "test" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "X12-ISA"
  qualifier_code = "12"
}
```

## Argument Reference

* `name` - (Optional) Exact name for the identifier.
* `label` - (Optional) Exact label for the identifier
* `qualifier_code` - (Optional) Exact code for the identifier qualifier
* `qualifier_label` - (Optional) Exact label for the identifier qualifier
* `environment_id` - (Required) ID of the environment in which to lookup identifier type

## Attribute Reference

* `id` - ID of the identifier type

[1]: https://docs.mulesoft.com/partner-manager/2.0/x12-identity-settings