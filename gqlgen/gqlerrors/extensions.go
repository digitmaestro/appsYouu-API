package gqlerrors

import (
	"encoding/json"
)

type extensions struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

// NewExtensions creates a new map with the default required extensions information
func NewExtensions(typ string, code string) (out map[string]interface{}) {
	ext := extensions{typ, code}

	return convertStructToMap(ext)
}

func convertStructToMap(in interface{}) (out map[string]interface{}) {
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		panic(err)
	}

	return
}
