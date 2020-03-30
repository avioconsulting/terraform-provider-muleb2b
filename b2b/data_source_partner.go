package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePartner() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePartnerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exact name of the partner",
			},
			"host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "True if the host should be retrieved, name will be ignored",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the environment to lookup Partner in",
			},
		},
	}
}

func dataSourcePartnerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string)
	client.SetEnvironment(envId)

	if v, ok := d.GetOk("host"); ok {
		if v.(bool) {
			partner, err := client.GetHostPartner()
			if err != nil {
				return err
			}
			if partner != nil {
				d.SetId(*partner.Id)
				d.Set("name", *partner.Name)
			}
		}
	}
	if v, ok := d.GetOk("name"); ok {
		name := v.(string)

		partner, err := client.GetPartnerByName(name)

		if err != nil {
			return err
		}

		if partner.Id == nil || *partner.Id == "" {
			return fmt.Errorf("no partner found with name (%s)", name)
		}

		d.SetId(*partner.Id)
		d.Set("name", *partner.Name)

		return nil
	}

	return fmt.Errorf("no partner name specified")
}
