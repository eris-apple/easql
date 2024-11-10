package where

import (
	"encoding/json"
	"strings"
)

type Where = map[string]interface{}

func NewWhereCondition(obj interface{}) (Where, error) {
	return structToMap(obj)
}

func structToMap(obj interface{}) (Where, error) {
	var unmarshal map[string]interface{}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &unmarshal)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	for key, value := range unmarshal {
		if value == nil || value == "" {
			continue
		}

		result[strings.ToLower(key)] = value
	}

	return result, nil
}
