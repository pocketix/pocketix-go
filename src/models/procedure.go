package models

import "encoding/json"

type Procedure struct {
	Name    string
	Program json.RawMessage
}
