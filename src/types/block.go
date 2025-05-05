package types

type Block struct {
	Id        string     `json:"id"`
	Body      []Block    `json:"block"`
	Arguments []Argument `json:"arguments"`
}
