# Environment Data Source

Provides data about [Mule Environments][1]

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}
```

## Argument Reference

* `name` - (Required) Exact name of the environment

## Attribute Reference

* `id` - ID of the environment

[1]: https://docs.mulesoft.com/access-management/environments