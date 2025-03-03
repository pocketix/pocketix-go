package models

import "github.com/pocketix/pocketix-go/src/interfaces"

type Program struct {
	Blocks []interfaces.Command `json:"block"`
}
