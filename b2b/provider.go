package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://devx.anypoint.mulesoft.com/",
				Description: "The base URL for the server running the Mule B2B API",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the organization that will be used for API operations",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user for the API operations",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for the user for the API operations",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"muleb2b_partner":       resourcePartner(),
			"muleb2b_endpoint":      resourceEndpoint(),
			"muleb2b_document":      resourceDocument(),
			"muleb2b_document_flow": resourceDocumentFlow(),
			"muleb2b_identifier":    resourceIdentifier(),
			"muleb2b_certificate":   resourceCertificate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"muleb2b_environment":     dataSourceEnvironment(),
			"muleb2b_ediDocumentType": dataSourceEdiDocumentType(),
			"muleb2b_partner":         dataSourcePartner(),
			"muleb2b_identifier_type": dataSourceIdentifierType(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	baseUrl, baseOk := d.GetOk("base_url")
	orgId, orgOk := d.GetOk("organization_id")
	user, userOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")

	if !baseOk {
		if val := os.Getenv("MULEB2B_BASE_URL"); val != "" {
			baseUrl = val
		}
	}

	if !orgOk {
		if val := os.Getenv("MULEB2B_ORG"); val != "" {
			orgId = val
		} else {
			return nil, fmt.Errorf("organization_id needs to be set in the muleb2b provider configuration or MULEB2B_ORG environment variable must be set")
		}
	}

	if !userOk {
		if val := os.Getenv("MULEB2B_USERNAME"); val != "" {
			user = val
		} else {
			return nil, fmt.Errorf("username needs to be set in the muleb2b provider configuration or MULEB2B_USERNAME environment variable must be set")
		}
	}

	if !passwordOk {
		if val := os.Getenv("MULEB2B_PASSWORD"); val != "" {
			password = val
		} else {
			return nil, fmt.Errorf("password needs to be set in the muleb2b provider configuration or MULEB2B_PASSWORD environment variable must be set")
		}
	}

	client, err := muleb2b.NewClient(muleb2b.String(baseUrl.(string)), muleb2b.String(orgId.(string)), nil)
	if err != nil {
		return nil, err
	}
	err = client.Login(user.(string), password.(string))
	if err != nil {
		return nil, err
	}

	return client, nil
}
