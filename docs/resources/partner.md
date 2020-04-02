# Partner Resource

Provides a [Mule B2b Partner][1] resource. This allows partners to be created, update, and deleted.  

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "duns" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "DUNS"
}

resource "muleb2b_partner" "test" {
  name           = "Test"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.duns.id
    value = "123456789"
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
    name = "Jane Doe"
    email = "jane.doe@test.com"
    type = "technical"
  }
  x12_inbound_config {
    character_encoding = "UTF8"
    acknowledgements {}
    validations {}
    control_numbers {
      require_unique_group_number = true
    }
  }
}

```

## Argument Reference

* `address` - (Optional) Address block representing partner's address
* `contact` - (Optional) Contact block for contacts associated with the partner
* `description` - (Optional) Brief description of the partner's business, and the trading relationship
* `environment_id` - (Required) Environment the partner will be created in
* `identifier` - (Required) Identifier block that uniquely identifies the partner. May specify multiple.
* `name` - (Required) Identifier for the partner
* `website_url` - (Optional) Trading partner's website
* `x12_inbound_config` - (Optional) X12 block for partner's x12 configuration

#### Address
The `address` block allows one to specify the partner's corporate address:
* `address_line_1` - (Required) First line of the partner's address
* `address_line_2` - (Optional) Second line of the partner's address
* `city` - (Required) Partner's city 
* `country` - (Required) Partner's country
* `postal_code` - (Required) Partner's postal code
* `state` - (Required) Partner's state

#### Contact
The `contact` block allows one to specify the partner's contacts
* `email` - (Required) Contact's email address
* `name` - (Required) Contact's full name
* `phone` - (Optional) Contact's phone number
* `type` - (Required) The type of the contact. Can be `"business"`, `"technical"`, or `"other"`

#### Identifier
The `identifier` block specifies the identifier value for the partner
* `identifier_type_id` - (Required) ID of the identifier type. Use the Identifier Type data source to look this up.
* `value` - (Required) Identifier Value. See the [Partner Manager Identifier documentation][3] for value rules

#### X12 Inbound Config
The `x12_inbound_config` block allows one to specify the partner's [X12 configuration][2]
* `character_encoding` - (Optional) Character encoding for messages from provider. Can be`"ASCII"`, `"ISO-8859-1"`, or `"UTF8"`
* `character_set` - (Optional) Characters allowed in string data. Can be `"BASIC"`, `"UNRESTRICTED"`, or `"EXTENDED"`. Defaults to `EXTENDED`
* `acknowledgements` - (Required) See [Acknowledgements](#acknowledgements)
* `control_numbers` - (Required) See [Control Numbers](#control numbers)
* `validations` - (Required) See [Validations](#validations)

##### Acknowledgements
The `acknowledgements` block, part of the `x12_inbound_config` block, allows one to specify acknowledgement details for the partner
* `endpoint_id` - (Optional) ID of the endpoint to be used for sending acknowledgements. Required if `generate_ta1` is true or `failure_acknowledgement_type` is non-zero.
* `generate_ta1` - (Optional) `true` if a technical acknowledgement should be sent to the partner. Defaults `false`.
* `failure_acknowledgement_type` - (Optional) The type of failure acknowledgement to be sent. In this case `0` means no failure acknowledgements will be sent. Can be `0`, `997`, or `999`. Defaults to `0`

##### Control Numbers
The `control_numbers` block, part of the `x12_inbound_config` block, allows one to specify the uniqueness of control numbers
* `require_unique_interchange_number` - (Optional) `true` if interchange control numbers (ISA13) must be unique. Defaults `true`
* `require_unique_group_number` - (Optional) `true` if group numbers (GS06) must be unique. Defaults `false`
* `require_unique_transaction_set_number` - (Optional) `true` if transaction set numbers (ST02) must be unique. Defaults `false`

##### Validations
The `validations` block, part of the `x12_inbound_config` block, allows one to specify the validations run on inbound messages
* `fail_when_value_length_outside_allowed_range` - (Optional) `true` to fail the transaction if values are too long or too short. Defaults `true`
* `fail_when_unused_segments_included` - (Optional) `true` to fail the transaction when segments marked as unused in the schema are included. Defaults `false`
* `fail_when_too_many_repeats_of_segment` - (Optional) `true` to fail the transaction when a segment is repeated too many times. Defaults `true`
* `fail_when_segments_out_of_order` - (Optional) `true` to fail the transaction when segments are out of order. Defaults `true`
* `fail_when_invalid_character_in_value` - (Optional) `true` to fail the transaction when invalid characters are present. Defaults `true`
* `fail_if_value_repeated_too_many_times` - (Optional) `true` to fail the transaction when values are repeated too man or too few times. Defaults `true`
* `fail_if_unknown_segments_used` - (Optional) `true` to fail the transaction when an unknown segment is used. Defaults `false`

## Attribute Reference

* `id` - The ID of the partner

[1]: https://docs.mulesoft.com/partner-manager/2.0/configure-partner
[2]: https://docs.mulesoft.com/partner-manager/2.0/x12-receive-read-settings
[3]: https://docs.mulesoft.com/partner-manager/2.0/x12-identity-settings