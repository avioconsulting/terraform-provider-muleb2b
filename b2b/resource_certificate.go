package b2b

import (
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment to add certificate to",
			},
			"partner_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the partner to add the certificate to",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the certificate",
			},
			"certificate_body": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the partner to create the document under",
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partnerId := d.Get("partner_id").(string)
	name := d.Get("name").(string)
	certificateBody := d.Get("certificate_body").(string)

	id, err := client.CreatePartnerCertificate(partnerId, string(stripCR([]byte(certificateBody))), name, "PEM")
	if err != nil {
		return err
	}
	if id != nil {
		d.SetId(*id)
	}

	return resourceCertificateRead(d, m)
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partnerId := d.Get("partner_id").(string)

	cert, err := client.GetPartnerCertificate(partnerId, d.Id())
	if err != nil {
		return err
	}

	d.Set("name", *cert.Name)

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	var partnerId string
	if d.HasChange("partner_id") {
		o, _ := d.GetChange("partner_id")
		partnerId = o.(string)
	} else {
		partnerId = d.Get("partner_id").(string)
	}

	err := client.DeletePartnerCertificate(partnerId, d.Id())
	if err != nil {
		return err
	}
	return nil
}

// strip CRs from raw literals. Lifted from go/scanner/scanner.go
// See https://github.com/golang/go/blob/release-branch.go1.6/src/go/scanner/scanner.go#L479
func stripCR(b []byte) []byte {
	c := make([]byte, len(b))
	i := 0
	for _, ch := range b {
		if ch != '\r' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}
