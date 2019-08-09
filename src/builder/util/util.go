package util

import (
	"encoding/json"
	"strconv"
	"time"
)

// StringToMap returns map interface to convert json string
func StringToMap(str string) *map[string]interface{} {
	var raw map[string]interface{}
	json.Unmarshal([]byte(str), &raw)
	return &raw
}

// GetTimeMillisecond returns current time to millisecond string
func GetTimeMillisecond() string {
	now := time.Now()
	nano := now.UnixNano()
	ms := nano / 1000000
	return strconv.FormatInt(ms, 10)
}
