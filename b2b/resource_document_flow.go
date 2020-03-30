package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceDocumentFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceDocumentFlowCreate,
		Read:   resourceDocumentFlowRead,
		Update: resourceDocumentFlowUpdate,
		Delete: resourceDocumentFlowDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name to give the document flow",
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDocumentFlowDirection,
				Description:  "The direction of the document flow. Only inbound is supported at this time.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the environment to place the document flow in",
			},
			"partner_from_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the partner that will initiate the document flow",
			},
			"partner_to_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the partner that will be on the receiving end of the document flow",
			},
			"config": {
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Description: "Low level configuration details of the document flow",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the document flow configuration - assigned after creation",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the document flow",
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version of the document flow configuration",
						},
						"preprocessing_endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"receiving_endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the endpoint that will initiate the flow",
						},
						"receiving_ack_endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the endpoint the acknowledgement should be sent to",
						},
						"target_endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the endpoint that will receive the document at the end of the flow",
						},
						"source_doc_type_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the document type that will be received from the sender",
						},
						"target_doc_type_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the document type that will be sent to the target",
						},
						"document_mapping": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the document mapping associated with the document flow",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the document mapping",
									},
									"file_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of mapping file",
									},
									"file_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Base64 encoded contents of the mapping file",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func validateDocumentFlowDirection(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
	}

	if v != "inbound" {
		errors = append(errors, fmt.Errorf("'inbound' is the only accepted value of %q at this time", key))
	}
	return warnings, errors
}

func resourceDocumentFlowCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	name := d.Get("name").(string)
	direction := d.Get("direction").(string)
	partnerFrom := d.Get("partner_from_id").(string)
	partnerTo := d.Get("partner_to_id").(string)

	dFlow := muleb2b.DocumentFlow{
		Name:          muleb2b.String(name),
		Direction:     muleb2b.String(strings.ToUpper(direction)),
		PartnerFromId: muleb2b.String(partnerFrom),
		PartnerToId:   muleb2b.String(partnerTo),
	}

	id, err := client.CreateDocumentFlow(&dFlow)
	if err != nil {
		return err
	}

	d.SetId(*id)

	newFlow, err := client.GetDocumentFlowById(*id)
	if err != nil {
		return err
	}

	err = updateFlowWithConfig(d.Get("config"), newFlow)
	if err != nil {
		return err
	}

	var mapping *muleb2b.Mapping = nil

	if len((*newFlow.Configurations[0]).DocumentMapping) > 0 {
		mapping = (*newFlow.Configurations[0]).DocumentMapping[0]
	}
	(*newFlow.Configurations[0]).DocumentMapping = []*muleb2b.Mapping{}

	udFlow, err := client.UpdateDocumentFlow(newFlow)
	if err != nil {
		return err
	}

	if mapping != nil {
		err = client.CreateMapping(*udFlow.Id, mapping)
		if err != nil {
			return err
		}
	}

	return resourceDocumentFlowRead(d, m)
}

func resourceDocumentFlowRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	dFlow, err := client.GetDocumentFlowById(id)
	if err != nil {
		return err
	}

	if dFlow != nil {
		d.Set("name", *(*dFlow).Name)
		d.Set("direction", strings.ToLower(*(*dFlow).Direction))
		d.Set("partner_from_id", *(*dFlow).PartnerFromId)
		d.Set("partner_to_id", *(*dFlow).PartnerToId)
		if len((*dFlow).Configurations) > 0 {
			if len((*(*dFlow).Configurations[0]).DocumentMapping) > 0 {
				mapping, err := client.GetMappingById(*dFlow.Id, *(*(*(*dFlow).Configurations[0]).DocumentMapping[0]).Id)
				if err != nil {
					return err
				}
				(*(*dFlow).Configurations[0]).DocumentMapping[0] = mapping

			}
			d.Set("config", flattenDocumentFlowConfig((*dFlow).Configurations[0]))
		} else {
			return fmt.Errorf("documentflow configuration is empty")
		}
	}
	return nil
}

func resourceDocumentFlowUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)
	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)
	cFlow, err := client.GetDocumentFlowById(d.Id())
	if err != nil {
		return err
	}
	dFlow := muleb2b.DocumentFlow{
		Id:            muleb2b.String(d.Id()),
		Name:          muleb2b.String(d.Get("name").(string)),
		Direction:     muleb2b.String(d.Get("direction").(string)),
		PartnerFromId: muleb2b.String(d.Get("partner_from_id").(string)),
		PartnerToId:   muleb2b.String(d.Get("partner_to_id").(string)),
	}

	dfConfig := expandDocumentFlowConfig(d.Get("config"), d.Id(), envId)

	if dfConfig != nil {
		if cFlow != nil && len((*cFlow).Configurations) > 0 {
			dfConfig.Id = (*(*cFlow).Configurations[0]).Id
			dfConfig.Version = (*(*cFlow).Configurations[0]).Version
			dfConfig.Status = (*(*cFlow).Configurations[0]).Status
		}
		dFlow.Configurations = []*muleb2b.DocumentFlowConfiguration{
			dfConfig,
		}
	}

	_, err = client.UpdateDocumentFlow(&dFlow)
	if err != nil {
		return err
	}

	return resourceDocumentFlowRead(d, m)
}

func resourceDocumentFlowDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	err := client.DeleteDocumentFlow(id)

	return err
}
