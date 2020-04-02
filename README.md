## Setting Up Your Development Environment
This code should be in the following directory ${GOPATH}/src/github.com/avioconsulting/terraform-provider-muleb2b
Also, you will need the muleb2b-client source as well. It is available at https://bitbucket.org/avioconsulting/muleb2b-client

## Building the provider
Execute `make build VERSION=${TAG}` and it will output `terraform-provider-muleb2b_${TAG}`

## Getting the provider binary
You may download the binary from the [Mule B2B Provider's releases page](https://github.com/avioconsulting/terraform-provider-muleb2b/releases), or you may use the local install script to build (or download) and install the provider as well.

Using the script (Linux/Unix/OS X):
```shell script
$ ./local_install.sh v1.0.0
``` 

Using the script (Windows):
```shell script
PS C:\path\to\this\repository\local_install.ps1 v1.0.0
```

## Using the provider
There are two options for using the provider on your system:
1. Place the terraform-provider-muleb2b_vX.Y.Z file in the same directory as your terraform code before `terraform init` is executed.
2. Move the terraform-provider-muleb2b_VX.Y.Z file into the Terraform plugins directory for your platform
    * Linux/Unix/OS X: ~/.terraform.d/plugins
    * Windows: %APPDATA%\terraform.d\plugins

## Integration Tests
Execute `make testacc` to run the integration tests. These environment variables must be set first: `MULEB2B_BASE_URL`, `MULEB2B_ORG`, `MULEB2B_USERNAME`, and `MULEB2B_PASSWORD`

## Terraform Cloud
Using this provider in Terraform Cloud requires that the linux amd64 version of the terraform-provider-muleb2b_vX.Y.Z binary be added to the Terraform project's directory structure in `${TERRAFORM_PROJECT_BASE_DIR}/terraform.d/plugins/linux_amd64`.

This can be also be done with a Git submodule. Here are the steps:
1. Create a repository for the provider binaries. This repository should contain all the providers that you will use that are not in the public registry.
2. Compile a linux_amd64 binary for the provider and make it executable using `chmod +x`
3. Upload the binary to the root of the provider repository
4. Add the submodule to the Terraform project by running this from the Terraform project's base directory: 
```bash
git submodule add \
  https://github.com/${ORGANIZATION}/${REPOSITORY}.git \
  terraform.d/plugins/linux_amd64
```

In your Terraform Cloud workspace, be sure you have [include submodules on clone](https://www.terraform.io/docs/cloud/workspaces/vcs.html#include-submodules-on-clone) enabled
HashiCorp has a [support article](https://support.hashicorp.com/hc/en-us/articles/360016992613-Using-custom-and-community-providers-in-Terraform-Cloud-and-Enterprise) for this use case as well.