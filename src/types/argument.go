package types

import (
	"encoding/json"
)

type Argument struct {
	Reference string          `json:"reference"`
	Type      string          `json:"type"`
	Value     json.RawMessage `json:"value"`
}
