package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func updateFlowWithConfig(data interface{}, documentFlow *muleb2b.DocumentFlow) error {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse: %#v", raw)
		}

		if val, ok := cfg["preprocessing_endpoint_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).PreProcessingEndpointId = muleb2b.String(val.(string))
		}

		if val, ok := cfg["receiving_endpoint_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).ReceivingEndpointId = muleb2b.String(val.(string))
		}

		if val, ok := cfg["receiving_ack_endpoint_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).ReceivingAckEndpointId = muleb2b.String(val.(string))
		}

		if val, ok := cfg["target_endpoint_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).TargetEndpointId = muleb2b.String(val.(string))
		}

		if val, ok := cfg["source_doc_type_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).SourceDocTypeId = muleb2b.String(val.(string))
		}

		if val, ok := cfg["target_doc_type_id"]; ok && val != "" {
			(*documentFlow.Configurations[0]).TargetDocTypeId = muleb2b.String(val.(string))
		}

		dmCfg, ok := cfg["document_mapping"]
		if ok {
			mapping, err := readDocumentMapping(dmCfg)
			if err != nil {
				return err
			}
			if mapping != nil {
				(*documentFlow.Configurations[0]).DocumentMapping = []*muleb2b.Mapping{
					mapping,
				}
			} else {
				(*documentFlow.Configurations[0]).DocumentMapping = []*muleb2b.Mapping{}
			}
		}

		return nil
	}
	return nil
}

func readDocumentMapping(data interface{}) (*muleb2b.Mapping, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse: %#v", raw)
		}

		fileName := cfg["file_name"].(string)
		fileContents := cfg["file_content"].(string)

		mapping := muleb2b.Mapping{
			MappingType:      muleb2b.String("DWL_FILE"),
			MappingSourceRef: muleb2b.String(fileName),
			MappingContent:   muleb2b.String(muleb2b.MappingContentPrefix + fileContents),
		}

		return &mapping, nil
	}
	return nil, nil
}

func flattenDocumentFlowConfig(config *muleb2b.DocumentFlowConfiguration) []interface{} {
	m := make(map[string]interface{})

	if config != nil {
		m["id"] = *(*config).Id
		m["status"] = *(*config).Status
		m["version"] = *(*config).Version
		if (*config).PreProcessingEndpointId != nil {
			m["preprocessing_endpoint_id"] = *(*config).PreProcessingEndpointId
		}
		if (*config).ReceivingEndpointId != nil {
			m["receiving_endpoint_id"] = *(*config).ReceivingEndpointId
		}
		if (*config).ReceivingAckEndpointId != nil {
			m["receiving_ack_endpoint_id"] = *(*config).ReceivingAckEndpointId
		}
		if (*config).TargetEndpointId != nil {
			m["target_endpoint_id"] = *(*config).TargetEndpointId
		}
		if (*config).SourceDocTypeId != nil {
			m["source_doc_type_id"] = *(*config).SourceDocTypeId
		}
		if (*config).TargetDocTypeId != nil {
			m["target_doc_type_id"] = *(*config).TargetDocTypeId
		}
		if len((*config).DocumentMapping) > 0 {
			m["document_mapping"] = flattenDocumentFlowMapping((*config).DocumentMapping[0])
		}
	}
	return []interface{}{m}
}

func flattenDocumentFlowMapping(mapping *muleb2b.Mapping) []interface{} {
	m := make(map[string]interface{})
	if mapping != nil {
		m["id"] = *(*mapping).Id
		m["type"] = *(*mapping).MappingType
		m["file_name"] = *(*mapping).MappingSourceRef
		m["file_content"] = strings.TrimPrefix(*(*mapping).MappingContent, muleb2b.MappingContentPrefix)
	}
	return []interface{}{m}
}

func expandDocumentFlowConfig(d interface{}, docFlowId, envId string) *muleb2b.DocumentFlowConfiguration {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			dfConfig := muleb2b.DocumentFlowConfiguration{
				Id:             muleb2b.String(configData["id"].(string)),
				DocumentFlowId: muleb2b.String(docFlowId),
				EnvironmentId:  muleb2b.String(envId),
				Status:         muleb2b.String(configData["status"].(string)),
				Version:        muleb2b.Integer(configData["version"].(int)),
			}

			if val, ok := configData["preprocessing_endpoint_id"]; ok {
				dfConfig.PreProcessingEndpointId = muleb2b.String(val.(string))
			}

			if val, ok := configData["receiving_endpoint_id"]; ok {
				dfConfig.ReceivingEndpointId = muleb2b.String(val.(string))
			}

			if val, ok := configData["receiving_ack_endpoint_id"]; ok {
				dfConfig.ReceivingAckEndpointId = muleb2b.String(val.(string))
			}

			if val, ok := configData["target_endpoint_id"]; ok {
				dfConfig.TargetEndpointId = muleb2b.String(val.(string))
			}

			if val, ok := configData["source_doc_type_id"]; ok {
				dfConfig.SourceDocTypeId = muleb2b.String(val.(string))
			}

			if val, ok := configData["target_doc_type_id"]; ok {
				dfConfig.TargetDocTypeId = muleb2b.String(val.(string))
			}

			dfMapping := expandDocumentFlowMapping(configData["document_mapping"])

			if dfMapping != nil {
				dfConfig.DocumentMapping = []*muleb2b.Mapping{
					dfMapping,
				}
			} else {
				dfConfig.DocumentMapping = []*muleb2b.Mapping{}
			}
			return &dfConfig
		}
	}
	return nil
}

func expandDocumentFlowMapping(d interface{}) *muleb2b.Mapping {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			return &muleb2b.Mapping{
				Id:               muleb2b.String(configData["id"].(string)),
				MappingType:      muleb2b.String(configData["type"].(string)),
				MappingContent:   muleb2b.String(configData["file_content"].(string)),
				MappingSourceRef: muleb2b.String(configData["file_name"].(string)),
			}
		}
	}
	return nil
}
