package models

import (
	"encoding/json"
)

type Argument struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}
