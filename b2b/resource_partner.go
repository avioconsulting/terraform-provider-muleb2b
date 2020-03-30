package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePartner() *schema.Resource {
	return &schema.Resource{
		Create: resourcePartnerCreate,
		Read:   resourcePartnerRead,
		Update: resourcePartnerUpdate,
		Delete: resourcePartnerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name to give the partner",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the environment the partner will be created in",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the partner",
			},
			"website_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The URL of the partner's website",
			},
			"identifier": {
				Type:        schema.TypeSet,
				MinItems:    1,
				Required:    true,
				Description: "Set of identifiers for the provider",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the identifier",
						},
						"identifier_type_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the identifier type to use",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier value",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the identifier",
						},
					},
				},
			},
			"x12_inbound_config": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "X12 Inbound Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"character_encoding": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"character_set": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "EXTENDED",
						},
						"acknowledgements": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"generate_ta1": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"failure_acknowledgement_type": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validateX12FailureAckType,
									},
								},
							},
						},
						"validations": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fail_when_value_length_outside_allowed_range": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"fail_when_unused_segments_included": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"fail_when_too_many_repeats_of_segment": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"fail_when_segments_out_of_order": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"fail_when_invalid_character_in_value": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"fail_if_value_repeated_too_many_times": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"fail_if_unknown_segments_used": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"control_numbers": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"require_unique_interchange_number": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Require unique interchange control number (ISA13)",
									},
									"require_unique_group_number": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Require unique group control number (GS06)",
									},
									"require_unique_transaction_set_number": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Require unique transaction set control number (ST02)",
									},
								},
							},
						},
					},
				},
			},
			"contact": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the contact",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the contact",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Contact's full name",
						},
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Contact's email address",
						},
						"phone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Contact's phone number",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Contact type: business, technical, or other",
						},
					},
				},
			},
			"address": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the address",
						},
						"address_line_1": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Line 1 of address",
						},
						"address_line_2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Line 2 of address",
						},
						"country": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Company's country",
						},
						"state": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Company's state or province",
						},
						"city": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Company's city",
						},
						"postal_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Company's postal code",
						},
					},
				},
			},
		},
	}
}

func validateX12FailureAckType(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be integer", key))
	}

	if v != 997 && v != 999 {
		errors = append(errors, fmt.Errorf("expected %q to be a either 997 or 999", key))
	}
	return warnings, errors
}

func resourcePartnerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	partner := muleb2b.Partner{
		Name:          muleb2b.String(d.Get("name").(string)),
		EnvironmentId: muleb2b.String(envId),
		Description:   muleb2b.String(d.Get("description").(string)),
		WebsiteUrl:    muleb2b.String(d.Get("website_url").(string)),
	}

	id, err := client.CreatePartner(&partner)

	if err != nil {
		return err
	}

	d.SetId(*id)

	// Create Identifiers
	cfg := d.Get("identifier")
	identifiers, err := readIdentifierConfig(cfg)
	if err != nil {
		return err
	}
	for _, i := range identifiers {
		err = client.CreatePartnerIdentifier(*id, i)
	}

	// Modify X12 Inbound Config - it's automatically created with the partner
	x12 := getDefaultInboundTemplate()
	if x12Cfg, ok := d.GetOk("x12_inbound_config"); ok {
		err = readX12InboundConfig(x12Cfg, x12)
		if err != nil {
			return err
		}
	}
	x12.Id = nil
	x12.PartnerId = muleb2b.String(*id)
	x12.EnvelopeHeaders = &muleb2b.X12EnvelopeHeaders{}
	err = client.CreatePartnerX12Configuration(*id, x12)
	if err != nil {
		return err
	}

	// Create Contacts
	if contactCfg, ok := d.GetOk("contact"); ok {
		contacts, err := readContactConfig(contactCfg)
		if err != nil {
			return err
		}
		err = client.UpdatePartnerContacts(*id, contacts)
		if err != nil {
			return err
		}
	}

	// Create Address
	if addressCfg, ok := d.GetOk("address"); ok {
		address, err := readAddressConfig(addressCfg)
		if err != nil {
			return err
		}
		err = client.UpdatePartnerAddress(*id, address)
		if err != nil {
			return err
		}
	}

	return resourcePartnerRead(d, meta)
}

func resourcePartnerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	partner, err := client.GetPartner(id)

	if err != nil {
		return err
	}

	d.Set("name", *(*partner).Name)
	d.Set("environment_id", *(*partner).EnvironmentId)
	d.Set("description", *(*partner).Description)
	d.Set("website_url", *(*partner).WebsiteUrl)

	// Get Identifiers
	identifiers, err := client.ListPartnerIdentifiers(*(*partner).Id)
	d.Set("identifier", flattenIdentifiers(identifiers))

	// Get X12 Inbound Config
	currentConfig, err := client.GetPartnerInboundX12Configuration(id)
	d.Set("x12_inbound_config", flattenX12InboundConfig(currentConfig))

	// Get Contacts
	contacts, err := client.GetPartnerContacts(id)
	d.Set("contact", flattenContacts(contacts))

	// Get Address
	address, err := client.GetPartnerAddress(id)
	if !address.Empty() {
		d.Set("address", flattenAddress(address))
	}

	return nil
}

func resourcePartnerUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client := meta.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	if d.HasChange("description") || d.HasChange("website_url") {
		partner := muleb2b.Partner{
			Id:            muleb2b.String(d.Id()),
			Name:          muleb2b.String(d.Get("name").(string)),
			EnvironmentId: muleb2b.String(envId),
			Description:   muleb2b.String(d.Get("description").(string)),
			WebsiteUrl:    muleb2b.String(d.Get("website_url").(string)),
		}

		err := client.UpdatePartner(&partner)
		if err != nil {
			return err
		}

		d.SetPartial("description")
		d.SetPartial("website_url")
	}

	if d.HasChange("identifier") {
		o, n := d.GetChange("identifier")
		oldIdentifiers := expandIdentifiers(o)
		newIdentifiers := expandIdentifiers(n)

		// Get the IDs and Statuses for the existing identifiers
		currentIdentifiers, err := client.ListPartnerIdentifiers(d.Id())
		if err != nil {
			return nil
		}
		for _, identifier := range currentIdentifiers {
			for _, oldIdentifier := range oldIdentifiers {
				if identifier.QualifierIdAndValueEqual(oldIdentifier) {
					oldIdentifier.Id = identifier.Id
					oldIdentifier.Status = identifier.Status
					break
				}
			}
		}

		del := identifierDifference(oldIdentifiers, newIdentifiers)
		add := identifierDifference(newIdentifiers, oldIdentifiers)

		for _, identifier := range add {
			if identifier.Status == nil || *identifier.Status == "" {
				identifier.Status = muleb2b.String("ACTIVE")
			}
			err := client.CreatePartnerIdentifier(d.Id(), identifier)
			if err != nil {
				return err
			}
		}
		for _, identifier := range del {
			if identifier.Id == nil || *identifier.Id == "" {
				return fmt.Errorf("identifier.Id is nil! (%s)\n", identifier.String())
			}
			err := client.DeletePartnerIdentifier(d.Id(), *identifier.Id)
			if err != nil {
				return err
			}
		}
		d.SetPartial("identifier")
	}

	// Handle X12 changes
	if d.HasChange("x12_inbound_config") {
		currentX12, err := client.GetPartnerInboundX12Configuration(d.Id())
		if err != nil {
			return err
		}
		expandX12InboundConfig(d.Get("x12_inbound_config"), currentX12)
		currentX12.EnvelopeHeaders = &muleb2b.X12EnvelopeHeaders{}
		err = client.UpdatePartnerX12Configuration(d.Id(), currentX12)
		if err != nil {
			return err
		}
		d.SetPartial("x12_inbound_config")
	}

	// Handle Contact changes
	if d.HasChange("contact") {
		o, n := d.GetChange("contact")
		oldContacts := expandContacts(o)
		newContacts := expandContacts(n)

		currentContacts, err := client.GetPartnerContacts(d.Id())
		if err != nil {
			return err
		}

		// Populate ID and Status of contacts
		for _, contact := range currentContacts {
			for _, oldContact := range oldContacts {
				if *contact.Name == *oldContact.Name &&
					*contact.Email == *oldContact.Email &&
					*contact.ContactType.Id == *oldContact.ContactType.Id {
					oldContact.Id = contact.Id
					oldContact.Status = contact.Status
					break
				}
			}
			for _, newContact := range newContacts {
				if *contact.Name == *newContact.Name &&
					*contact.Email == *newContact.Email &&
					*contact.ContactType.Id == *newContact.ContactType.Id {
					newContact.Id = contact.Id
					newContact.Status = contact.Status
					break
				}
			}
		}

		// Delete contacts that must be deleted
		del := contactDifference(oldContacts, newContacts)
		for _, contact := range del {
			err = client.DeletePartnerContact(d.Id(), *contact.Id)
			if err != nil {
				return err
			}
		}

		// Update with new contacts
		err = client.UpdatePartnerContacts(d.Id(), newContacts)
		if err != nil {
			return err
		}
		d.SetPartial("contact")
	}

	if d.HasChange("address") {
		currentAddress, err := client.GetPartnerAddress(d.Id())
		if err != nil {
			return err
		}
		newAddress := expandAddress(d.Get("address"))
		if newAddress != nil {
			newAddress.Id = currentAddress.Id
			err = client.UpdatePartnerAddress(d.Id(), newAddress)
		} else {
			err = client.DeletePartnerAddress(d.Id())
		}
		if err != nil {
			return err
		}
		d.SetPartial("address")
	}

	d.Partial(false)
	return resourcePartnerRead(d, meta)
}

func identifierDifference(l1, l2 []*muleb2b.Identifier) []*muleb2b.Identifier {
	var diff []*muleb2b.Identifier
	for _, id1 := range l1 {
		found := false
		for _, id2 := range l2 {
			if id1.QualifierIdAndValueEqual(id2) {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, id1)
		}
	}
	return diff
}

func contactDifference(l1, l2 []*muleb2b.Contact) []*muleb2b.Contact {
	var diff []*muleb2b.Contact
	for _, contact1 := range l1 {
		found := false
		for _, contact2 := range l2 {
			if *contact1.Name == *contact2.Name &&
				*contact1.Email == *contact2.Email &&
				*contact1.ContactType.Id == *contact2.ContactType.Id {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, contact1)
		}
	}
	return diff
}

func resourcePartnerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	err := client.DeletePartnerById(muleb2b.String(id))

	if err != nil {
		return err
	}

	// No need to delete the Identifiers because they will be deleted with the partner
	// No need to delete the X12 configuration because it will be deleted with the partner
	// No need to delete the contacts or address either

	return nil
}
