package utils

import "encoding/json"

func JsonPretty(data any) string {
	val, _ := json.MarshalIndent(data, "", "  ")

	return string(val)
}
