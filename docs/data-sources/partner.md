# Partner Data Source

Provides data for a [Mule B2B Partner][1] 

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

data "muleb2b_partner" "test" {
  name           = "Test-2"
  environment_id = data.muleb2b_environment.sbx.id
}
```

## Argument Reference

* `name` - (Optional) Exact name of the partner
* `host` - (Optional) `true` if the host provider should be retrieved, name will be ignored
* `environment_id` - (Required) ID of the environment in which to perform the lookup

## Attribute Reference

* `id` - ID of the retrieved Partner

[1]: https://docs.mulesoft.com/partner-manager/2.0/configure-partner