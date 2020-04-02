# Document Flow Resource

Provides [Mule B2B Document Flow][1] resource.

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
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
  name           = "my-partner"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "my-id"
  }
}

resource "muleb2b_endpoint" "test" {
  name = "my-endpoint"
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
  name = "my-flow"
  direction = "inbound"
  environment_id = data.muleb2b_environment.sbx.id
  partner_from_id = muleb2b_partner.test.id
  partner_to_id = data.muleb2b_partner.host.id
  config {
    receiving_endpoint_id = muleb2b_endpoint.test.id
  }
}
```

## Argument Reference

* `name` - (Required) List arguments this resource takes.
* `direction` - (Required) - direction of the document flow. Only inbound is supported by Mule B2B at this time.
* `environment_id` - (Required) - ID of the environment in which to create the document flow 
* `partner_from_id` - (Required) - ID of the partner from which the messages will be received
* `partner_to_id` - (Required) - ID of the partner to which messages will be sent
* `config` - (Required) Low level configuration details of the document flow

#### Config
The `config` block allows one to configure the endpoint's SFTP settings

* `preprocessing_endpoint_id` - (Optional)
* `receiving_endpoint_id` - (Optional) ID of the endpoint that will initiate the flow
* `receiving_ack_endpoint_id` - (Optional) ID of the endpoint the acknowledgement should be sent to
* `target_endpoint_id` - (Optional) ID of the endpoint that will receive the document at the end of the flow
* `source_doc_type_id` - (Optional) ID of the document that will be received from the sender
* `target_doc_type_id` - (Optional) ID of the document that will be sent to the target
* `document_mapping` - (Optional) Block describing how to transform from the source document to the target document

##### Document Mapping
The `document_mapping` block, a part of the `config` block, specifies the mapping between the `source_doc_type_id` and the `target_doc_type_id`

* `file_name` - (Required) Name of the mapping file
* `file_content` - (Required) Base64 encoded contents of the mapping file

## Attribute Reference

* `id` - ID of the document flow

[1]: https://docs.mulesoft.com/partner-manager/2.0/message-flows