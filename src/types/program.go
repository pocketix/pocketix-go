package types

type Program struct {
	Header Header  `json:"header"`
	Blocks []Block `json:"block"`
}
