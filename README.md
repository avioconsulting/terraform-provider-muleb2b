## Setting Up Your Development Environment
This code should be in the following directory ${GOPATH}/src/github.com/avioconsulting/terraform-provider-muleb2b
Also, you will need the muleb2b-client source as well. It is available at https://bitbucket.org/avioconsulting/muleb2b-client

## Building the code
Execute `make build` and it will output terraform-provider-muleb2b

## Executing the code 
Place the terraform-provider-muleb2b file in the same directory as your terraform code before `terraform init` is executed.

## Integration Tests
Execute `make testacc` to run the integration tests. These environment variables must be set first: `MULEB2B_BASE_URL`, `MULEB2B_ORG`, `MULEB2B_USERNAME`, and `MULEB2B_PASSWORD`