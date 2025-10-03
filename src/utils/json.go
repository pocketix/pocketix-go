package utils

import "encoding/json"

func UnmarshalData[T any](data []byte) (*T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return &result, err
}
