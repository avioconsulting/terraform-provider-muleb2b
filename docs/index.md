# Mule B2B Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.

## Example Usage

```hcl
provider "muleb2b" {
  base_url = "https://anypoint.mulesoft.com/"
  organization_id = "be4f0fba-541b-5f82-b51d-f047b6569645"
  username = "test"
  password = "test_password"
}
```
## Authentication
The Mule B2B provider only offers username/password authentication.

### Static Credentials
!> !> Warning: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided by adding a `username` and `password` in-line in the Mule B2B provider block.

Usage:
```hcl-terraform
provider "muleb2b" {
  organization_id = "be4f0fba-541b-5f82-b51d-f047b6569645"
  username = "my-test-user"
  password = "my-test-password"
}
```

### Environment Variables
You can provide your credentials with the `MULEB2B_USERNAME` and `MULEB2B_PASSWORD` environment variables. You may also use the `MULEB2B_ORG` environment variable rather than the  `organization_id` variable.
```hcl-terraform
provider "muleb2b" {}
``` 

Usage:
```shell script
$ export MULEB2B_ORG="be4f0fba-541b-5f82-b51d-f047b6569645"
$ export MULEB2B_USERNAME="my-test-user"
$ export MULEB2B_PASSWORD="my-test-password"
$ terraform plan
```

## Argument Reference

The following arguments are supported in the Mule B2B provider block.

* `base_url` - (Optional) The base URL for the Mule B2B API. Typically it is `https://anypoint.mulesoft.com/`
* `organization_id` - (Optional) Either this or the `MULEB2B_ORG` environment variable are required. This is the organization all the resources will be created under. This is the Business Group Id from your organization on Anypoint.
* `username` - (Optional) Either this or the `MULEB2B_USERNAME` environment variable are required.
* `password` - (Optional) Either this or the `MULEB2B_PASSWORD` environment variable are required.
