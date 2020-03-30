package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func readIdentifierConfig(data interface{}) ([]*muleb2b.Identifier, error) {
	config := data.(*schema.Set).List()
	var identifiers []*muleb2b.Identifier
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse: %#v", raw)
		}

		identifier := muleb2b.Identifier{
			IdentifierTypeQualifierId: muleb2b.String(cfg["identifier_type_id"].(string)),
			Status:                    muleb2b.String("ACTIVE"),
			Value:                     muleb2b.String(cfg["value"].(string)),
		}

		identifiers = append(identifiers, &identifier)
	}

	return identifiers, nil
}
func flattenIdentifiers(identifiers []*muleb2b.Identifier) []interface{} {
	var out = make([]map[string]interface{}, len(identifiers), len(identifiers))
	for i, v := range identifiers {
		m := make(map[string]interface{})
		m["id"] = *v.Id
		m["identifier_type_id"] = *v.IdentifierTypeQualifierId
		m["value"] = *v.Value
		m["status"] = *v.Status
		out[i] = m
	}
	return []interface{}{out}
}
func expandIdentifiers(d interface{}) []*muleb2b.Identifier {
	var identifiers []*muleb2b.Identifier
	if d != nil {
		configList := d.(*schema.Set).List()
		for i := range configList {
			configData := configList[i].(map[string]interface{})
			identifier := muleb2b.Identifier{
				Id:                        muleb2b.String(configData["id"].(string)),
				IdentifierTypeQualifierId: muleb2b.String(configData["identifier_type_id"].(string)),
				Status:                    muleb2b.String(configData["status"].(string)),
				Value:                     muleb2b.String(configData["value"].(string)),
			}
			identifiers = append(identifiers, &identifier)
		}
	}
	return identifiers
}

func readContactConfig(data interface{}) ([]*muleb2b.Contact, error) {
	config := data.(*schema.Set).List()
	var contacts []*muleb2b.Contact
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse: %#v", raw)
		}

		contact := muleb2b.Contact{
			Name:  muleb2b.String(cfg["name"].(string)),
			Email: muleb2b.String(cfg["email"].(string)),
		}

		if v, ok := cfg["phone"]; ok {
			contact.Phone = muleb2b.String(v.(string))
		}

		switch cfg["type"].(string) {
		case "business":
			contact.ContactType = muleb2b.GetBusinessContactType()
		case "technical":
			contact.ContactType = muleb2b.GetTechnicalContactType()
		default:
			contact.ContactType = muleb2b.GetOtherContactType()
		}

		contacts = append(contacts, &contact)
	}

	return contacts, nil
}
func flattenContacts(contacts []*muleb2b.Contact) []interface{} {
	var out = make([]map[string]interface{}, len(contacts), len(contacts))
	for i, v := range contacts {
		m := make(map[string]interface{})
		m["id"] = *v.Id
		m["name"] = *v.Name
		m["email"] = *v.Email
		m["status"] = *v.Status
		if v.Phone != nil && *v.Phone != "" {
			m["phone"] = *v.Phone
		}
		switch *v.ContactType.Name {
		case "Business":
			m["type"] = "business"
		case "Technical":
			m["type"] = "technical"
		default:
			m["type"] = "other"
		}
		out[i] = m
	}
	return []interface{}{out}
}
func expandContacts(d interface{}) []*muleb2b.Contact {
	var contacts []*muleb2b.Contact
	if d != nil {
		configList := d.(*schema.Set).List()
		for i := range configList {
			configData := configList[i].(map[string]interface{})
			contact := muleb2b.Contact{
				Name:  muleb2b.String(configData["name"].(string)),
				Email: muleb2b.String(configData["email"].(string)),
			}
			if v, ok := configData["id"]; ok && v.(string) != "" {
				contact.Id = muleb2b.String(v.(string))
			} else {
				contact.Id = nil
			}
			if v, ok := configData["status"]; ok && v.(string) != "" {
				contact.Status = muleb2b.String(v.(string))
			} else {
				contact.Status = nil
			}
			if v, ok := configData["phone"]; ok {
				contact.Phone = muleb2b.String(v.(string))
			}
			switch configData["type"].(string) {
			case "business":
				contact.ContactType = muleb2b.GetBusinessContactType()
			case "technical":
				contact.ContactType = muleb2b.GetTechnicalContactType()
			default:
				contact.ContactType = muleb2b.GetOtherContactType()
			}
			contacts = append(contacts, &contact)
		}
	}
	return contacts
}

func readAddressConfig(data interface{}) (*muleb2b.Address, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse: %#v", raw)
		}

		address := muleb2b.Address{
			Id:         nil,
			Addr1:      muleb2b.String(cfg["address_line_1"].(string)),
			City:       muleb2b.String(cfg["city"].(string)),
			State:      muleb2b.String(cfg["state"].(string)),
			Country:    muleb2b.String(cfg["country"].(string)),
			PostalCode: muleb2b.String(cfg["postal_code"].(string)),
		}

		if v, ok := cfg["address_line_2"]; ok && v.(string) != "" {
			address.Addr2 = muleb2b.String(v.(string))
		} else {
			address.Addr2 = nil
		}

		return &address, nil
	}
	return nil, nil
}
func flattenAddress(address *muleb2b.Address) []interface{} {
	m := make(map[string]interface{})
	if address != nil {
		m["id"] = *address.Id
		m["address_line_1"] = *address.Addr1
		if address.Addr2 != nil {
			m["address_line_2"] = *address.Addr2
		}
		m["city"] = *address.City
		m["state"] = *address.State
		m["country"] = *address.Country
		m["postal_code"] = *address.PostalCode
	}
	return []interface{}{m}
}
func expandAddress(d interface{}) *muleb2b.Address {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})

			address := muleb2b.Address{
				Id:         nil,
				Addr1:      muleb2b.String(configData["address_line_1"].(string)),
				Addr2:      muleb2b.String(configData["address_line_2"].(string)),
				City:       muleb2b.String(configData["city"].(string)),
				State:      muleb2b.String(configData["state"].(string)),
				Country:    muleb2b.String(configData["country"].(string)),
				PostalCode: muleb2b.String(configData["postal_code"].(string)),
			}

			return &address
		}
	}
	return nil
}
