package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDocument() *schema.Resource {
	return &schema.Resource{
		Create: resourceDocumentCreate,
		Read:   resourceDocumentRead,
		Update: resourceDocumentUpdate,
		Delete: resourceDocumentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the document type",
			},
			"partner_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the partner to create the document under",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the environment to create the document under",
			},
			"edi_document_type_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the document type",
			},
			"schema_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base64 Encoded contents of the schema file. Not needed unless using a custom schema.",
			},
			"custom_schema_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the custom schema",
			},
		},
	}
}

func resourceDocumentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	name := d.Get("name").(string)
	partnerId := d.Get("partner_id").(string)
	docTypeId := d.Get("edi_document_type_id").(string)
	schemaFile, ok := d.GetOk("schema_file")

	doc := muleb2b.Document{
		Name:              muleb2b.String(name),
		EdiDocumentTypeId: muleb2b.String(docTypeId),
	}

	if ok {
		doc.SchemaContent = muleb2b.String(schemaFile.(string))
		doc.SchemaType = muleb2b.String("customSchemaType")
		doc.Standard = muleb2b.Boolean(false)
	} else {
		doc.Standard = muleb2b.Boolean(true)
	}

	id, err := client.CreateDocument(partnerId, &doc)

	if err != nil {
		return err
	} else if id == nil {
		return fmt.Errorf("nil id returned from document service")
	}

	d.SetId(*id)

	return resourceDocumentRead(d, m)
}

func resourceDocumentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	id := d.Id()
	partnerId := d.Get("partner_id").(string)

	doc, err := client.GetDocumentById(partnerId, id)
	if err != nil {
		return err
	}

	if doc != nil {
		d.Set("partner_id", partnerId)
		d.Set("edi_document_type_id", *doc.EdiDocumentTypeId)
		if doc.SchemaContent != nil {
			d.Set("schema_file", *doc.SchemaContent)
			d.Set("custom_schema_id", *doc.CustomSchemaId)
		}
	}
	return nil
}

func resourceDocumentUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDocumentCreate(d, m)
}

func resourceDocumentDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
