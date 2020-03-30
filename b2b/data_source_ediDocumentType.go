package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceEdiDocumentType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEdiDocumentTypeRead,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the environment to lookup EDI Document Types in",
			},
			"format_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the EDI format type",
			},
			"format_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"document_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Exact name of the environment",
			},
			"edi_format_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the EDI format",
			},
			"edi_format_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the EDI format version",
			},
		},
	}
}

func dataSourceEdiDocumentTypeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	formatType := d.Get("format_type").(string)
	formatVersion := d.Get("format_version").(string)
	documentName := d.Get("document_name").(string)

	ediFormat, err := client.GetEdiFormatByFormat(formatType)
	if err != nil {
		return err
	}
	if ediFormat == nil {
		return fmt.Errorf("format type (%s) is not valid in this environment (%s)", formatType, envId)
	}

	version, err := client.GetEdiFormatVersionByLabel(formatType, formatVersion)
	if err != nil {
		return err
	}
	if version == nil {
		return fmt.Errorf("format version (%s) is not valid for this type (%s) and/or environment (%s)", formatVersion, formatType, envId)
	}

	docType, err := client.GetEdiDocumentTypeByName(formatType, *version.Id, documentName)
	if err != nil {
		return err
	}
	if docType == nil {
		return fmt.Errorf("EDI Document Type (%s) is not valid for this type (%s), version (%s), and/or environment (%s)", documentName, formatType, formatVersion, envId)
	}

	d.SetId(*docType.Id)
	d.Set("edi_format_version_id", *version.Id)
	d.Set("edi_format_id", *ediFormat.Id)

	return nil
}
