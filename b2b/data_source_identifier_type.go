package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIdentifierType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentifierTypeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exact name of the identifier",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exact label of the identifier",
			},
			"qualifier_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exact code for the identifier qualifier",
			},
			"qualifier_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exact label for the identifier qualifier",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of environment in which to lookup identifier type",
			},
		},
	}
}

func dataSourceIdentifierTypeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	var identifiers []*muleb2b.IdentifierType
	var err error

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)

		identifiers, err = client.GetIdentifierTypesByName(name)
		if err != nil {
			return err
		}

		if len(identifiers) == 0 {
			return fmt.Errorf("no identifiers found with name (%s)", name)
		}
		if len(identifiers) > 1 {
			return fmt.Errorf("multiple results found with name (%s)", name)
		}
	} else if v, ok := d.GetOk("label"); ok {
		label := v.(string)

		identifiers, err = client.GetIdentifierTypesByLabel(label)
		if err != nil {
			return err
		}

		if len(identifiers) == 0 {
			return fmt.Errorf("no identifiers found with label (%s)", label)
		}
		if len(identifiers) > 1 {
			return fmt.Errorf("multiple results found with label (%s)", label)
		}
	} else {
		return fmt.Errorf("no identifier name or label specified")
	}

	if len((*identifiers[0]).Qualifiers) == 1 {
		qualifier := (*identifiers[0]).Qualifiers[0]
		d.SetId(*(*qualifier).Id)
		d.Set("qualifier_code", *(*qualifier).Code)
		d.Set("qualifier_label", *(*qualifier).Label)
		d.Set("name", *(*identifiers[0]).Name)
		d.Set("label", *(*identifiers[0]).Label)

		return nil
	}

	var qualifiers []*muleb2b.IdentifierTypeQualifier

	if v, ok := d.GetOk("qualifier_code"); ok {
		code := v.(string)

		qualifiers, err = identifiers[0].GetIdentifierTypeQualifiersByCode(code)
		if err != nil {
			return err
		} else if len(qualifiers) == 0 {
			return fmt.Errorf("no qualifiers found with code (%s)", code)
		} else if len(qualifiers) > 1 {
			return fmt.Errorf("multiple results found with code (%s)", code)
		}

	} else if v, ok := d.GetOk("qualifier_label"); ok {
		label := v.(string)

		qualifiers, err = identifiers[0].GetIdentifierTypeQualifiersByCode(label)
		if err != nil {
			return err
		} else if len(qualifiers) == 0 {
			return fmt.Errorf("no qualifiers found with label (%s)", label)
		} else if len(qualifiers) > 1 {
			return fmt.Errorf("multiple results found with label (%s)", label)
		}
	} else {
		return fmt.Errorf("multiple identifier qualifiers found and no qualifier code or label specified")
	}

	d.SetId(*(*qualifiers[0]).Id)
	d.Set("qualifier_code", *(*qualifiers[0]).Code)
	d.Set("qualifier_label", *(*qualifiers[0]).Label)
	d.Set("name", *(*identifiers[0]).Name)
	d.Set("label", *(*identifiers[0]).Label)

	return nil
}
