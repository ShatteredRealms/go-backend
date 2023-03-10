package service

import (
	"bytes"
	"encoding/json"
)

func structToJson(s interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(s)
	return &buf, err
}
