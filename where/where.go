package where

import (
	"fmt"
	"reflect"
)

type Where = map[string]interface{}

func NewWhereCondition(obj interface{}, includeEmpty bool) (Where, error) {
	return structToMap(obj, includeEmpty)
}

// structToMap turns the structure into a [string]interface{} map using the `map` tags.
// includeEmpty determines whether to include fields with empty values in the result.
func structToMap(input interface{}, includeEmpty bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("map")
		if tag == "" {
			continue
		}

		if !includeEmpty && IsEmptyValue(value) {
			continue
		}

		if value.CanInterface() {
			result[tag] = value.Interface()
		}
	}

	return result, nil
}

// IsEmptyValue checks if the value is "empty" ("" for strings, 0 for numbers, nil for interfaces, etc.)
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
