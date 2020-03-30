package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIdentifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentifierCreate,
		Read:   resourceIdentifierRead,
		Update: nil,
		Delete: resourceIdentifierDelete,

		Schema: map[string]*schema.Schema{
			"partner_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the partner to create the identifier under",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the environment the partner is in",
			},
			"identifier_type_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the identifier type to use",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The identifier value",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the identifier",
			},
		},
	}
}

func resourceIdentifierCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partnerId := d.Get("partner_id").(string)

	identifier := muleb2b.Identifier{
		IdentifierTypeQualifierId: muleb2b.String(d.Get("identifier_type_id").(string)),
		Status:                    muleb2b.String("ACTIVE"),
		Value:                     muleb2b.String(d.Get("value").(string)),
	}

	err := client.CreatePartnerIdentifier(partnerId, &identifier)
	if err != nil {
		return err
	}

	newId, err := client.GetPartnerIdentifierByQualifierIdAndValue(partnerId, *identifier.IdentifierTypeQualifierId, *identifier.Value)
	if err != nil {
		return err
	} else if newId == nil {
		return fmt.Errorf("identifier (%s, %s) not created for partner (%s)", *identifier.IdentifierTypeQualifierId, *identifier.Value, partnerId)
	}

	d.SetId(*newId.Id)

	return nil
}

func resourceIdentifierRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partnerId := d.Get("partner_id").(string)

	identifier, err := client.GetPartnerIdentifierById(partnerId, d.Id())
	if err != nil {
		return err
	}

	d.Set("identifier_type_id", *identifier.IdentifierTypeQualifierId)
	d.Set("value", *identifier.Value)
	d.Set("status", *identifier.Status)

	return nil
}

func resourceIdentifierDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partnerId := d.Get("partner_id").(string)

	err := client.DeletePartnerIdentifier(partnerId, d.Id())
	return err
}
