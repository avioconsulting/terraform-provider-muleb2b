package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Exact name of the environment",
			},
		},
	}
}

func dataSourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)

		env, err := client.GetEnvironmentByName(name)

		if err != nil {
			return err
		}

		if env == nil || env.Id == nil || *env.Id == "" {
			return fmt.Errorf("no environment found with name (%s)", name)
		}

		d.SetId(*env.Id)
		d.Set("name", *env.Name)

		return nil
	}

	return fmt.Errorf("no environment name specified")
}
