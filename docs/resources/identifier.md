# Identifier Resource

Provides [Mule B2B Identifier][1] resource. Allows one to add identifiers to a partner, including the host partner. 

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "duns" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "DUNS"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

resource "muleb2b_identifier" "host-duns" {
  partner_id = muleb2b_partner.host.id
  environment_id = data.muleb2b_environment.sbx.id
  identifier_type_id = data.muleb2b_identifier_type.duns.id
  value = "987654321"
}
```

## Argument Reference

* `partner_id` - (Required) ID of partner to add identifier to
* `environment_id` - (Required) ID of environment to add identifier to
* `identifier_type_id` - (Required) ID of the identifier type
* `value` - (Required) Identifier value

## Attribute Reference

* `id` - Identifier's ID
* `status` - Status of the identifier


[1]: https://docs.mulesoft.com/partner-manager/2.0/x12-identity-settings
