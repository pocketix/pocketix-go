package types

import "encoding/json"

type Header struct {
	Variables  json.RawMessage `json:"userVariables"`
	Procedures json.RawMessage `json:"userProcedures"`
}
