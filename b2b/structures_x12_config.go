package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
)

func readX12InboundConfig(data interface{}, x12 *muleb2b.X12) error {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse: %#v", raw)
		}

		if v, ok := cfg["character_encoding"]; ok {
			x12.CharacterSetAndEncoding.CharacterEncoding = muleb2b.String(v.(string))
		}
		if v, ok := cfg["character_set"]; ok {
			x12.CharacterSetAndEncoding.CharacterSet = muleb2b.String(v.(string))
		}

		if v, ok := cfg["acknowledgements"]; ok {
			err := readX12AcknowledgementsConfig(v, x12)
			if err != nil {
				return err
			}
		}
		if v, ok := cfg["validations"]; ok {
			err := readX12validationsConfig(v, x12)
			if err != nil {
				return err
			}
		}
		if v, ok := cfg["control_numbers"]; ok {
			err := readX12ControlNumbersConfig(v, x12)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readX12AcknowledgementsConfig(data interface{}, x12 *muleb2b.X12) error {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse: %#v", raw)
		}

		if v, ok := cfg["endpoint_id"]; ok {
			if v.(string) != "" {
				x12.ParserSettings.AckEndpointId = muleb2b.String(v.(string))
			} else {
				x12.ParserSettings.AckEndpointId = nil
			}
		}

		if v, ok := cfg["generate_ta1"]; ok {
			x12.ParserSettings.GenerateTA1 = muleb2b.Boolean(v.(bool))
		}

		if v, ok := cfg["failure_acknowledgement_type"]; ok {
			if v.(int) == 997 {
				x12.ParserSettings.Require997 = muleb2b.Boolean(true)
			} else if v.(int) == 999 {
				x12.ParserSettings.Generate999 = muleb2b.Boolean(true)
			}
		}
	}
	return nil
}

func readX12validationsConfig(data interface{}, x12 *muleb2b.X12) error {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse: %#v", raw)
		}

		if v, ok := cfg["fail_when_value_length_outside_allowed_range"]; ok {
			x12.ParserSettings.FailDocumentWhenValueLengthOutsideAllowedRange = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_when_unused_segments_included"]; ok {
			x12.ParserSettings.FailDocumentIfUnknownSegmentsAreUsed = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_when_too_many_repeats_of_segment"]; ok {
			x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_when_segments_out_of_order"]; ok {
			x12.ParserSettings.FailDocumentWhenSegmentsAreOutOfOrder = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_when_invalid_character_in_value"]; ok {
			x12.ParserSettings.FailDocumentWhenInvalidCharacterInValue = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_if_value_repeated_too_many_times"]; ok {
			x12.ParserSettings.FailDocumentIfValueIsRepeatedTooManyTimes = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["fail_if_unknown_segments_used"]; ok {
			x12.ParserSettings.FailDocumentIfUnknownSegmentsAreUsed = muleb2b.Boolean(v.(bool))
		}
	}
	return nil
}

func readX12ControlNumbersConfig(data interface{}, x12 *muleb2b.X12) error {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse: %#v", raw)
		}
		if v, ok := cfg["require_unique_interchange_number"]; ok {
			x12.ControlNumberSettings.RequireUniqueISAcontrolNumbersISA13 = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["require_unique_group_number"]; ok {
			x12.ControlNumberSettings.RequireUniqueGSControlNumbersGS06 = muleb2b.Boolean(v.(bool))
		}
		if v, ok := cfg["require_unique_transaction_set_number"]; ok {
			x12.ControlNumberSettings.RequireUniqueTransactionSetControlNumbersST02 = muleb2b.Boolean(v.(bool))
		}

	}
	return nil
}

func flattenX12InboundConfig(x12 *muleb2b.X12) []interface{} {
	m := make(map[string]interface{})
	if x12 != nil {
		if x12.CharacterSetAndEncoding.CharacterEncoding != nil {
			m["character_encoding"] = *x12.CharacterSetAndEncoding.CharacterEncoding
		}
		if x12.CharacterSetAndEncoding.CharacterSet != nil {
			m["character_set"] = *x12.CharacterSetAndEncoding.CharacterSet
		}
		m["acknowledgements"] = flattenX12Acknowledgements(x12.ParserSettings)
		m["validations"] = flattenX12Validations(x12.ParserSettings)
		m["control_numbers"] = flattenX12ControlNumbers(x12.ControlNumberSettings)
	}
	return []interface{}{m}
}

func flattenX12Acknowledgements(parserSettings *muleb2b.X12ParserSettings) []interface{} {
	m := make(map[string]interface{})
	if parserSettings != nil {
		if parserSettings.AckEndpointId != nil {
			m["endpoint_id"] = *parserSettings.AckEndpointId
		}
		if parserSettings.GenerateTA1 != nil {
			m["generate_ta1"] = *parserSettings.GenerateTA1
		}
		if parserSettings.Generate999 != nil && *parserSettings.Generate999 {
			m["failure_acknowledgement_type"] = muleb2b.Integer(999)
		} else if parserSettings.Require997 != nil && *parserSettings.Require997 {
			m["failure_acknowledgement_type"] = muleb2b.Integer(997)
		} else {
			m["failure_acknowledgement_type"] = muleb2b.Integer(0)
		}
	}
	return []interface{}{m}
}

func flattenX12Validations(parserSettings *muleb2b.X12ParserSettings) []interface{} {
	m := make(map[string]interface{})
	if parserSettings != nil {
		if parserSettings.FailDocumentWhenValueLengthOutsideAllowedRange != nil {
			m["fail_when_value_length_outside_allowed_range"] = *parserSettings.FailDocumentWhenValueLengthOutsideAllowedRange
		}
		if parserSettings.FailDocumentWhenUnusedSegmentsAreIncluded != nil {
			m["fail_when_unused_segments_included"] = *parserSettings.FailDocumentWhenUnusedSegmentsAreIncluded
		}
		if parserSettings.FailDocumentWhenTooManyRepeatsOfSegment != nil {
			m["fail_when_too_many_repeats_of_segment"] = *parserSettings.FailDocumentWhenTooManyRepeatsOfSegment
		}
		if parserSettings.FailDocumentWhenSegmentsAreOutOfOrder != nil {
			m["fail_when_segments_out_of_order"] = *parserSettings.FailDocumentWhenSegmentsAreOutOfOrder
		}
		if parserSettings.FailDocumentWhenInvalidCharacterInValue != nil {
			m["fail_when_invalid_character_in_value"] = *parserSettings.FailDocumentWhenInvalidCharacterInValue
		}
		if parserSettings.FailDocumentIfValueIsRepeatedTooManyTimes != nil {
			m["fail_if_value_repeated_too_many_times"] = *parserSettings.FailDocumentIfValueIsRepeatedTooManyTimes
		}
		if parserSettings.FailDocumentIfUnknownSegmentsAreUsed != nil {
			m["fail_if_unknown_segments_used"] = *parserSettings.FailDocumentIfUnknownSegmentsAreUsed
		}
	}
	return []interface{}{m}
}

func flattenX12ControlNumbers(settings *muleb2b.X12ControlNumberSettings) []interface{} {
	m := make(map[string]interface{})
	if settings != nil {
		if settings.RequireUniqueISAcontrolNumbersISA13 != nil {
			m["require_unique_interchange_number"] = *settings.RequireUniqueISAcontrolNumbersISA13
		}
		if settings.RequireUniqueGSControlNumbers != nil {
			m["require_unique_group_number"] = *settings.RequireUniqueGSControlNumbersGS06
		}
		if settings.RequireUniqueTransactionSetControlNumbersST02 != nil {
			m["require_unique_transaction_set_number"] = *settings.RequireUniqueTransactionSetControlNumbersST02
		}
	}
	return []interface{}{m}
}

func expandX12InboundConfig(d interface{}, x12 *muleb2b.X12) {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			if v, ok := configData["character_encoding"]; ok {
				x12.CharacterSetAndEncoding.CharacterEncoding = muleb2b.String(v.(string))
			}

			if v, ok := configData["character_set"]; ok {
				x12.CharacterSetAndEncoding.CharacterSet = muleb2b.String(v.(string))
			} else {
				x12.CharacterSetAndEncoding.CharacterSet = nil
			}

			expandX12Acknowledgements(configData["acknowledgements"], x12)
			expandX12Validations(configData["validations"], x12)
			expandX12ControlNumbers(configData["control_numbers"], x12)
		}
	}

}

func expandX12Acknowledgements(d interface{}, x12 *muleb2b.X12) {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			if v, ok := configData["endpoint_id"]; ok {
				if v.(string) != "" {
					x12.ParserSettings.AckEndpointId = muleb2b.String(v.(string))
				} else {
					x12.ParserSettings.AckEndpointId = nil
				}
			} else {
				x12.ParserSettings.AckEndpointId = nil
			}

			if v, ok := configData["generate_ta1"]; ok {
				x12.ParserSettings.GenerateTA1 = muleb2b.Boolean(v.(bool))
			}
			if v, ok := configData["failure_acknowledgement_type"]; ok {
				if v.(int) == 997 {
					x12.ParserSettings.Require997 = muleb2b.Boolean(true)
					x12.ParserSettings.Generate999 = muleb2b.Boolean(false)
				} else if v.(int) == 999 {
					x12.ParserSettings.Require997 = muleb2b.Boolean(false)
					x12.ParserSettings.Generate999 = muleb2b.Boolean(true)
				} else {
					x12.ParserSettings.Require997 = muleb2b.Boolean(false)
					x12.ParserSettings.Generate999 = muleb2b.Boolean(false)
				}
			} else {
				x12.ParserSettings.Require997 = muleb2b.Boolean(false)
				x12.ParserSettings.Generate999 = muleb2b.Boolean(false)
			}

		}
	}
}

func expandX12Validations(d interface{}, x12 *muleb2b.X12) {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			x12.ParserSettings.FailDocumentWhenValueLengthOutsideAllowedRange = muleb2b.Boolean(configData["fail_when_value_length_outside_allowed_range"].(bool))
			x12.ParserSettings.FailDocumentWhenUnusedSegmentsAreIncluded = muleb2b.Boolean(configData["fail_when_unused_segments_included"].(bool))
			x12.ParserSettings.FailDocumentWhenTooManyRepeatsOfSegment = muleb2b.Boolean(configData["fail_when_too_many_repeats_of_segment"].(bool))
			x12.ParserSettings.FailDocumentWhenSegmentsAreOutOfOrder = muleb2b.Boolean(configData["fail_when_segments_out_of_order"].(bool))
			x12.ParserSettings.FailDocumentWhenInvalidCharacterInValue = muleb2b.Boolean(configData["fail_when_invalid_character_in_value"].(bool))
			x12.ParserSettings.FailDocumentIfValueIsRepeatedTooManyTimes = muleb2b.Boolean(configData["fail_if_value_repeated_too_many_times"].(bool))
			x12.ParserSettings.FailDocumentIfUnknownSegmentsAreUsed = muleb2b.Boolean(configData["fail_if_unknown_segments_used"].(bool))
		}
	}
}

func expandX12ControlNumbers(d interface{}, x12 *muleb2b.X12) {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			x12.ControlNumberSettings.RequireUniqueISAcontrolNumbersISA13 = muleb2b.Boolean(configData["require_unique_interchange_number"].(bool))
			x12.ControlNumberSettings.RequireUniqueGSControlNumbersGS06 = muleb2b.Boolean(configData["require_unique_group_number"].(bool))
			x12.ControlNumberSettings.RequireUniqueTransactionSetControlNumbersST02 = muleb2b.Boolean(configData["require_unique_transaction_set_number"].(bool))
		}
	}
}

func getDefaultInboundTemplate() *muleb2b.X12 {
	x12 := muleb2b.X12{
		ConfigType:      muleb2b.String("READ"),
		FormatType:      muleb2b.String("X12InboundConfig"),
		FormatTypeId:    muleb2b.String("25c1bc8a-801f-4337-a2a6-7721ef971460"),
		EnvelopeHeaders: &muleb2b.X12EnvelopeHeaders{},
		ParserSettings: &muleb2b.X12ParserSettings{
			FailDocumentWhenValueLengthOutsideAllowedRange: muleb2b.Boolean(true),
			FailDocumentWhenInvalidCharacterInValue:        muleb2b.Boolean(true),
			FailDocumentIfValueIsRepeatedTooManyTimes:      muleb2b.Boolean(true),
			FailDocumentIfUnknownSegmentsAreUsed:           muleb2b.Boolean(false),
			FailDocumentWhenSegmentsAreOutOfOrder:          muleb2b.Boolean(true),
			FailDocumentWhenTooManyRepeatsOfSegment:        muleb2b.Boolean(true),
			FailDocumentWhenUnusedSegmentsAreIncluded:      muleb2b.Boolean(false),
			Require997:    muleb2b.Boolean(false),
			Generate999:   muleb2b.Boolean(false),
			GenerateTA1:   muleb2b.Boolean(false),
			AckEndpointId: nil,
		},
		CharacterSetAndEncoding: &muleb2b.X12CharacterSetAndEncoding{
			CharacterSet:              muleb2b.String("EXTENDED"),
			CharacterEncoding:         nil,
			LineEndingBetweenSegments: nil,
		},
		ControlNumberSettings: &muleb2b.X12ControlNumberSettings{
			InitialInterchangeControlNumber:               muleb2b.String("00"),
			InitialGSControlNumber:                        muleb2b.String("00"),
			InitialTransactionSetControlNumber:            muleb2b.String("00"),
			RequireUniqueGSControlNumbers:                 nil,
			RequireUniqueTransactionSetControlNumber:      nil,
			RequireUniqueISAcontrolNumbersISA13:           muleb2b.Boolean(true),
			RequireUniqueGSControlNumbersGS06:             muleb2b.Boolean(false),
			RequireUniqueTransactionSetControlNumbersST02: muleb2b.Boolean(false),
		},
	}

	return &x12
}
