package b2b

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"muleb2b": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MULEB2B_ORG"); v == "" {
		t.Fatal("MULEB2B_ORG must be set for acceptance tests")
	}
	if v := os.Getenv("MULEB2B_USERNAME"); v == "" {
		t.Fatal("MULEB2B_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("MULEB2B_PASSWORD"); v == "" {
		t.Fatal("MULEB2B_PASSWORD must be set for acceptance tests")
	}
}
