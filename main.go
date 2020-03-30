package main

import (
	"github.com/avioconsulting/terraform-provider-muleb2b/b2b"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return b2b.Provider()
		},
	})
}
