package types

type Block struct {
	Id        string     `json:"id"`
	Body      []Block    `json:"block"`
	Arguments []Argument `json:"arguments"`
	Devices   []Device   `json:"devices"`
}

type Device struct {
	ID     string `json:"deviceId"`
	Values []any  `json:"values"`
}
