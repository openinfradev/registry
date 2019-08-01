package util

import "encoding/json"

// StringToMap returns map interface to convert json string
func StringToMap(str string) *map[string]interface{} {
	var raw map[string]interface{}
	json.Unmarshal([]byte(str), &raw)
	return &raw
}
