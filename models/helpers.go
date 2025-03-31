package models

import (
	"fmt"
	"slices"
)

// validateEnum checks if the value is in the allowed choices
func validateEnum(value any, field map[string]any, varName string,
	convertToInt func(any) (int, error)) error {

	switch choices := field["choices"].(type) {
	case []string:
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("input wrong type for enum %s: expected string", varName)
		}
		if !slices.Contains(choices, strVal) {
			return fmt.Errorf("input not in enum (%s): %s not in %v", varName, strVal, choices)
		}
	case []int:
		intVal, err := convertToInt(value)
		if err != nil {
			return err
		}
		if !slices.Contains(choices, intVal) {
			return fmt.Errorf("input not in enum (%s): %d not in %v", varName, intVal, choices)
		}
	default:
		return fmt.Errorf("unexpected type for choices (%s)", varName)
	}
	return nil
}

// validateType checks if the value matches the expected type
func validateType(value any, fieldType, varName string) error {
	switch fieldType {
	case "string", "string_enum":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("input wrong type for %s: expected string, got %T", varName, value)
		}
	case "int", "int_enum":
		switch v := value.(type) {
		case int:
			// Valid type
		case float64:
			// Check if it's an integer value
			if v != float64(int(v)) {
				return fmt.Errorf("input wrong type for %s: expected integer, got float %v", varName, v)
			}
		default:
			return fmt.Errorf("input wrong type for %s: expected int, got %T", varName, value)
		}
	case "decimal":
		switch value.(type) {
		case float64, int:
			// Valid types
		default:
			return fmt.Errorf("input wrong type for %s: expected decimal, got %T", varName, value)
		}
	default:
		return fmt.Errorf("unknown field type %s for %s", fieldType, varName)
	}
	return nil
}
